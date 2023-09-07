// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";

import {SafeCall} from "../libraries/SafeCall.sol";
import {Types} from "../libraries/Types.sol";
import {Hashing} from "../libraries/Hashing.sol";
import {Encoding} from "../libraries/Hashing.sol";
import {MerkleTrie} from "../libraries/trie/MerkleTrie.sol";
import {SecureMerkleTrie} from "../libraries/trie/SecureMerkleTrie.sol";
import {AddressAliasHelper} from "../vendor/AddressAliasHelper.sol";
import {IL1Portal} from "./IL1Portal.sol";
import {IRollup} from "../IRollup.sol";

import "../libraries/Errors.sol";

abstract contract L1PortalDeterministicStorage {
    /**
     * @notice A list of initiated deposit hashes.
     */
    mapping(bytes32 => bool) public initiatedDeposits;
}

/**
 * @custom:proxied
 * @title L1Portal
 * @notice The L1Portal is a low-level contract responsible for passing messages between L1
 *         and L2. Messages sent directly to the L1Portal have no form of replayability.
 *         Users are encouraged to use the L1CrossDomainMessenger for a higher-level interface.
 */
contract L1Portal is
    L1PortalDeterministicStorage,
    IL1Portal,
    Initializable,
    UUPSUpgradeable,
    OwnableUpgradeable,
    PausableUpgradeable
{
    /**
     * @notice Value used to reset the l2Sender, this is more efficient than setting it to zero.
     */
    address internal constant DEFAULT_L2_SENDER = 0x000000000000000000000000000000000000dEaD;

    /**
     * @notice The L2 gas limit set when eth is deposited using the receive() function.
     */
    uint64 internal constant RECEIVE_DEFAULT_GAS_LIMIT = 100_000;

    /**
     * @notice Additional gas reserved for clean up after finalizing a transaction withdrawal.
     */
    uint256 internal constant FINALIZE_GAS_BUFFER = 20_000;

    /**
     * @notice The Rollup.
     */
    IRollup public rollup;

    /**
     * @notice Address of the L2Portal deployed on L2.
     */
    address public l2PortalAddress; // TODO: store the hash instead

    /**
     * @notice Address of the L2 account which initiated a withdrawal in this transaction. If the
     *         of this variable is the default L2 sender address, then we are NOT inside of a call
     *         to finalizeWithdrawalTransaction.
     */
    address public l2Sender;

    /**
     * @notice A unique value hashed with each deposit.
     */
    uint256 public nonce;

    /**
     * @notice A list of withdrawal hashes which have been successfully finalized.
     */
    mapping(bytes32 => bool) public finalizedWithdrawals;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initializer;
     */
    function initialize(address _rollup) public initializer {
        if (_rollup == address(0)) {
            revert ZeroAddress();
        }

        rollup = IRollup(_rollup);
        l2Sender = DEFAULT_L2_SENDER;

        __Ownable_init();
        __Pausable_init();
        __UUPSUpgradeable_init();
    }

    function setL2PortalAddress(address _l2PortalAddress) external onlyOwner {
        l2PortalAddress = _l2PortalAddress;
    }

    function pause() public onlyOwner {
        _pause();
    }

    function unpause() public onlyOwner {
        _unpause();
    }

    function _authorizeUpgrade(address) internal override onlyOwner whenPaused {}

    /**
     * @notice Accepts value so that users can send ETH directly to this contract and have the
     *         funds be deposited to their address on L2. This is intended as a convenience
     *         function for EOAs. Contracts should call the initiateDeposit() function directly
     *         otherwise any deposited funds will be lost due to address aliasing.
     */
    receive() external payable whenNotPaused {
        initiateDeposit(msg.sender, RECEIVE_DEFAULT_GAS_LIMIT, bytes(""));
    }

    /**
     * @notice Accepts ETH value without triggering a deposit to L2.
     */
    function donateETH() external payable {
        // Intentionally empty.
    }

    // @inheritdoc IL1Portal
    function finalizeWithdrawalTransaction(
        Types.CrossDomainMessage memory withdrawalTx,
        uint256 assertionID,
        // bytes calldata encodedBlockHeader,
        bytes[] calldata withdrawalAccountProof,
        bytes[] calldata withdrawalProof
    ) external override L2Deployed onlyProxy whenNotPaused {
        // Prevent nested withdrawals within withdrawals.
        require(l2Sender == DEFAULT_L2_SENDER, "L1Portal: can only trigger one withdrawal per transaction");

        // Prevent users from creating a deposit transaction where this address is the message
        // sender on L2.
        // In the context of the proxy delegate calling to this implementation,
        // address(this) will return the address of the proxy.
        require(withdrawalTx.target != address(this), "L1Portal: you cannot send messages to the portal contract");

        // Get the L2 assertion claimed to include this withdrawal.
        IRollup.Assertion memory assertion = rollup.getAssertion(assertionID);

        // Ensure that the assertion is confirmed.
        require(_isAssertionConfirmed(assertionID, assertion.stateHash), "L1Portal: assertion not confirmed");

        // All withdrawals have a unique hash, we'll use this as the identifier for the withdrawal
        // and to prevent replay attacks.
        bytes32 withdrawalHash = Hashing.hashCrossDomainMessage(withdrawalTx);

        // Check that this withdrawal has not already been finalized, this is replay protection.
        require(finalizedWithdrawals[withdrawalHash] == false, "L1Portal: withdrawal has already been finalized");

        // Avoid stack too deep
        {
            // (bytes32 blockHash, bytes32 stateRoot) = Encoding.decodeStateRootFromEncodedBlockHeader(encodedBlockHeader);

            // Verify that the block hash is the assertion's stateHash.
            // require(blockHash == assertion.stateHash, "L1Portal: invalid block");

            // Verify the account proof.
            bytes32 storageRoot = _verifyAccountInclusion(l2PortalAddress, assertion.stateHash, withdrawalAccountProof);

            // Verify that the hash of this withdrawal was stored in the L2Portal contract on L2.
            // If this is true, then we know that this withdrawal was actually triggered on L2
            // and can therefore be relayed on L1.
            require(
                _verifyWithdrawalInclusion(withdrawalHash, storageRoot, withdrawalProof),
                "L1Portal: invalid withdrawal inclusion proof"
            );
        }

        // Mark the withdrawal as finalized so it can't be replayed.
        finalizedWithdrawals[withdrawalHash] = true;

        // We want to maintain the property that the amount of gas supplied to the call to the
        // target contract is at least the gas limit specified by the user. We can do this by
        // enforcing that, at this point in time, we still have gaslimit + buffer gas available.
        require(
            gasleft() >= withdrawalTx.gasLimit + FINALIZE_GAS_BUFFER,
            "L1Portal: insufficient gas to finalize withdrawal"
        );

        // Set the l2Sender so contracts know who triggered this withdrawal on L2.
        l2Sender = withdrawalTx.sender;

        // Trigger the call to the target contract. We use SafeCall because we don't
        // care about the returndata and we don't want target contracts to be able to force this
        // call to run out of gas via a returndata bomb.
        bool success = SafeCall.call(withdrawalTx.target, withdrawalTx.gasLimit, withdrawalTx.value, withdrawalTx.data);

        // Reset the l2Sender back to the default value.
        l2Sender = DEFAULT_L2_SENDER;

        // All withdrawals are immediately finalized. Replayability can
        // be achieved through contracts built on top of this contract
        emit WithdrawalFinalized(withdrawalHash, success);
    }

    /// @inheritdoc IL1Portal
    function isAssertionConfirmed(uint256 assertionID) public view override returns (bool) {
        IRollup.Assertion memory assertion = rollup.getAssertion(assertionID);
        return _isAssertionConfirmed(assertionID, assertion.stateHash);
    }

    /// @inheritdoc IL1Portal
    function initiateDeposit(address target, uint256 gasLimit, bytes memory data)
        public
        payable
        override
        L2Deployed
        onlyProxy
        whenNotPaused
    {
        // Transform the from-address to its alias if the caller is a contract.
        address from = msg.sender;
        if (msg.sender != tx.origin) {
            from = AddressAliasHelper.applyL1ToL2Alias(msg.sender);
        }

        bytes32 depositHash = Hashing.hashCrossDomainMessage(
            Types.CrossDomainMessage({
                version: 0,
                nonce: nonce,
                sender: from,
                target: target,
                value: msg.value,
                gasLimit: gasLimit,
                data: data
            })
        );

        initiatedDeposits[depositHash] = true;
        emit DepositInitiated(nonce, from, target, msg.value, gasLimit, data, depositHash);

        // Increment the nonce so that the next deposit will have a different hash.
        unchecked {
            nonce++;
        }
    }

    /**
     * @notice Determine if a given assertion is finalized and confirmed.
     *
     * @param assertionID The ID of the assertion.
     * @param stateHash   The stateHash field of the assertion.
     */
    function _isAssertionConfirmed(uint256 assertionID, bytes32 stateHash) internal view returns (bool) {
        // Must be finalized.
        if (assertionID > rollup.getLastConfirmedAssertionID()) {
            return false;
        }

        // Must be confirmed.
        if (stateHash == bytes32(0)) {
            return false;
        }

        return true;
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
        require(exists, "L1Portal: invalid account proof");

        return Encoding.decodeStorageRootFromEncodedAccount(encodedAccount);
    }

    /**
     * @notice Verifies a Merkle Trie inclusion proof that a given withdrawal hash is present in
     *         the storage of the L2ToL1MessagePasser contract.
     *
     * @param withdrawalHash Hash of the withdrawal to verify.
     * @param storageRoot    Root of the storage of the L2ToL1MessagePasser contract.
     * @param proof          Inclusion proof of the withdrawal hash in the storage root.
     */
    function _verifyWithdrawalInclusion(bytes32 withdrawalHash, bytes32 storageRoot, bytes[] memory proof)
        internal
        pure
        returns (bool)
    {
        bytes32 storageKey = keccak256(
            abi.encode(
                withdrawalHash,
                uint256(0) // The withdrawals mapping is at the first slot in the layout.
            )
        );

        return SecureMerkleTrie.verifyInclusionProof(abi.encode(storageKey), hex"01", proof, storageRoot);
    }

    /**
     * Modifiers
     */

    modifier L2Deployed() {
        require(l2PortalAddress != address(0), "L1Portal: L2 Portal not deployed");
        _;
    }
}
