// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";

import {SafeCall} from "../libraries/SafeCall.sol";
import {Types} from "../libraries/Types.sol";
import {Hashing} from "../libraries/Hashing.sol";
import {Encoding} from "../libraries/Hashing.sol";
import {MerkleTrie} from "../libraries/trie/MerkleTrie.sol";
import {SecureMerkleTrie} from "../libraries/trie/SecureMerkleTrie.sol";
import {Predeploys} from "../libraries/Predeploys.sol";
import {AddressAliasHelper} from "../vendor/AddressAliasHelper.sol";
import {IL2Portal} from "./IL2Portal.sol";
import {L1Oracle} from "./L1Oracle.sol";

import "../libraries/Errors.sol";

abstract contract L2PortalDeterministicStorage {
    /**
     * @notice A list of initiated withdrawal hashes.
     */
    mapping(bytes32 => bool) public initiatedWithdrawals;
}

contract L2Portal is
    L2PortalDeterministicStorage,
    IL2Portal,
    UUPSUpgradeable,
    OwnableUpgradeable,
    PausableUpgradeable
{
    /**
     * @notice Value used to reset the l1Sender, this is more efficient than setting it to zero.
     */
    address internal constant DEFAULT_L1_SENDER = 0x000000000000000000000000000000000000dEaD;

    /**
     * @notice The L1 gas limit set when eth is withdrawn using the receive() function.
     */
    uint256 internal constant RECEIVE_DEFAULT_GAS_LIMIT = 100_000;

    /**
     * @notice Additional gas reserved for clean up after finalizing a transaction deposit.
     */
    uint256 internal constant FINALIZE_GAS_BUFFER = 20_000;

    /**
     * @notice Address of the L2Portal deployed on L1.
     */
    address public l1PortalAddress; // TODO: store the hash instead

    /**
     * @notice Address of the L1 account which initiated a deposit in this transaction. If the
     *         of this variable is the default L1 sender address, then we are NOT inside of a call
     *         to finalizeDepositTransaction.
     */
    address public l1Sender;

    /**
     * @notice A unique value hashed with each withdrawal.
     */
    uint256 public nonce;

    /**
     * @notice A list of deposit hashes which have been successfully finalized.
     */
    mapping(bytes32 => bool) public finalizedDeposits;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initializer;
     */
    function initialize(address _l1PortalAddress) public initializer {
        if (_l1PortalAddress == address(0)) {
            revert ZeroAddress();
        }

        l1PortalAddress = _l1PortalAddress;
        l1Sender = DEFAULT_L1_SENDER;

        __Ownable_init();
        __UUPSUpgradeable_init();
    }

    function pause() public onlyOwner {
        _pause();
    }

    function unpause() public onlyOwner {
        _unpause();
    }

    function _authorizeUpgrade(address) internal override onlyOwner whenPaused {}

    /**
     * @notice Allows users to withdraw ETH by sending directly to this contract.
     */
    receive() external payable {
        initiateWithdrawal(msg.sender, RECEIVE_DEFAULT_GAS_LIMIT, bytes(""));
    }

    /**
     * @notice Accepts ETH value without triggering a withdrawal to L1.
     */
    function donateETH() external payable {
        // Intentionally empty.
    }

    /**
     * @notice Sends a message from L2 to L1.
     *
     * @param _target   Address to call on L1 execution.
     * @param _gasLimit Minimum gas limit for executing the message on L1.
     * @param _data     Data to forward to L1 target.
     */
    function initiateWithdrawal(address _target, uint256 _gasLimit, bytes memory _data) public payable whenNotPaused {
        bytes32 withdrawalHash = Hashing.hashCrossDomainMessage(
            Types.CrossDomainMessage({
                version: 0,
                nonce: nonce,
                sender: msg.sender,
                target: _target,
                value: msg.value,
                gasLimit: _gasLimit,
                data: _data
            })
        );

        initiatedWithdrawals[withdrawalHash] = true;

        emit WithdrawalInitiated(nonce, msg.sender, _target, msg.value, _gasLimit, _data, withdrawalHash);
        unchecked {
            ++nonce;
        }
    }

    // @inheritdoc IL2Portal
    function finalizeDepositTransaction(
        uint256 l1BlockNumber,
        Types.CrossDomainMessage memory depositTx,
        bytes[] calldata depositAccountProof,
        bytes[] calldata depositProof
    ) external whenNotPaused {
        // TODO: re-add `onlyProxy`
        // Prevent nested deposits within deposits.
        require(l1Sender == DEFAULT_L1_SENDER, "L2Portal: can only trigger one deposit per transaction");

        // Prevent users from creating a deposit transaction where this address is the message
        // sender on L2.
        // In the context of the proxy delegate calling to this implementation,
        // address(this) will return the address of the proxy.
        require(depositTx.target != address(this), "L2Portal: you cannot send messages to the portal contract");

        // All deposits have a unique hash, we'll use this as the identifier for the deposit
        // and to prevent replay attacks.
        bytes32 depositHash = Hashing.hashCrossDomainMessage(depositTx);

        // Check that this deposit has not already been finalized, this is replay protection.
        require(finalizedDeposits[depositHash] == false, "L2Portal: deposit has already been finalized");

        bytes32 stateRoot = L1Oracle(Predeploys.L1_ORACLE).prevStateRoots(uint8(l1BlockNumber % 256));

        // Verify the account proof.
        bytes32 storageRoot = _verifyAccountInclusion(l1PortalAddress, stateRoot, depositAccountProof);

        // Verify that the hash of this deposit was stored in the L1Portal contract on L1.
        // If this is true, then we know that this deposit was actually triggered on L1
        // and can therefore be relayed on L2.
        require(
            _verifyDepositInclusion(depositHash, storageRoot, depositProof), "L2Portal: invalid deposit inclusion proof"
        );

        // Mark the deposit as finalized so it can't be replayed.
        finalizedDeposits[depositHash] = true;

        // We want to maintain the property that the amount of gas supplied to the call to the
        // target contract is at least the gas limit specified by the user. We can do this by
        // enforcing that, at this point in time, we still have gaslimit + buffer gas available.
        require(gasleft() >= depositTx.gasLimit + FINALIZE_GAS_BUFFER, "L2Portal: insufficient gas to finalize deposit");

        // Set the l2Sender so contracts know who triggered this deposit on L2.
        l1Sender = depositTx.sender;

        // Trigger the call to the target contract. We use SafeCall because we don't
        // care about the returndata and we don't want target contracts to be able to force this
        // call to run out of gas via a returndata bomb.
        bool success = SafeCall.call(depositTx.target, depositTx.gasLimit, depositTx.value, depositTx.data);

        require(success, "L2Portal: call to target contract reverted");

        // Reset the l2Sender back to the default value.
        l1Sender = DEFAULT_L1_SENDER;

        emit DepositFinalized(depositHash, success);
    }

    /**
     * @notice Verifies a Merkle Trie inclusion proof that an account is present in
     *         the world state and extract its storage root.
     *
     * @param account      Account address to verify.
     * @param stateRoot    Root of the world state root of L2.
     * @param proof        Inclusion proof of the account in the storage root.
     */
    function _verifyAccountInclusion(address account, bytes32 stateRoot, bytes[] memory proof)
        internal
        pure
        returns (bytes32)
    {
        (bool exists, bytes memory encodedAccount) = SecureMerkleTrie.get(abi.encodePacked(account), proof, stateRoot);
        require(exists, "L2Portal: invalid account proof");

        return Encoding.decodeStorageRootFromEncodedAccount(encodedAccount);
    }

    /**
     * @notice Verifies a Merkle Trie inclusion proof that a given deposit hash is present in
     *         the storage of the L2ToL1MessagePasser contract.
     *
     * @param _depositHash    Hash of the deposit to verify.
     * @param _storageRoot    Root of the storage of the L1Portal contract.
     * @param _proof          Inclusion proof of the deposit hash in the storage root.
     */
    function _verifyDepositInclusion(bytes32 _depositHash, bytes32 _storageRoot, bytes[] memory _proof)
        internal
        pure
        returns (bool)
    {
        bytes32 storageKey = keccak256(
            abi.encode(
                _depositHash,
                uint256(0) // The deposits mapping is at the first slot in the layout.
            )
        );

        return SecureMerkleTrie.verifyInclusionProof(abi.encode(storageKey), hex"01", _proof, _storageRoot);
    }
}
