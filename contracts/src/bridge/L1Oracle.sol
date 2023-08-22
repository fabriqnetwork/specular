// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/PausableUpgradeable.sol";

contract L1Oracle is Initializable, UUPSUpgradeable, OwnableUpgradeable, PausableUpgradeable {
    /**
     * @notice Emitted when the L1 stateRoot is updated.
     */
    event L1OracleValuesUpdated(uint256 blockNumber, bytes32 stateRoot, uint256 baseFee);

    /**
     * @notice The latest L1 block number known by the L2 system.
     */
    uint256 public blockNumber;

    /**
     * @notice The latest L1 stateRoot known by the L2 system.
     */
    bytes32 public stateRoot;

    /**
     * @notice The latest L1 base fee known by the L2 system.
     */
    uint256 public baseFee;

    /**
     * @notice The address of the L2 sequencer.
     */
    address public sequencer;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initializer;
     * @param _sequencer The address of the L2 sequencer.
     */
    function initialize(address _sequencer) public initializer {
        sequencer = _sequencer;
        __Ownable_init();
        __Pausable_init();
        __UUPSUpgradeable_init();
    }

    function _authorizeUpgrade(address) internal override onlyOwner {}

    /**
     * @notice Updates the L1 block values.
     *
     * @param _blockNumber L1 blockNumber.
     * @param _stateRoot L1 stateRoot.
     * @param _baseFee L1 baseFee.
     */
    function setL1OracleValues(uint256 _blockNumber, bytes32 _stateRoot, uint256 _baseFee) external onlySequencer {
        blockNumber = _blockNumber;
        stateRoot = _stateRoot;
        baseFee = _baseFee;
        emit L1OracleValuesUpdated(blockNumber, _stateRoot, _baseFee);
    }

    /**
     * @notice Updates the L2 sequencer address.
     *
     * @param _sequencer L2 sequencer address.
     */
    function setSequencer(address _sequencer) external onlyOwner {
        sequencer = _sequencer;
    }

    /**
     * @notice Modifier to check if the caller is the sequencer.
     */
    modifier onlySequencer() {
        require(msg.sender == sequencer, "Only the sequencer can call this function.");
        _;
    }
}
