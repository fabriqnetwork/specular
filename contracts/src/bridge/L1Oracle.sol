// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";

/**
 * @custom:proxied
 * @title L1Oracle
 * @notice The L1Oracle stores latests known L1 block values.
 *         Should be updated by the sequencer only.
 */
contract L1Oracle is Initializable, UUPSUpgradeable, OwnableUpgradeable, PausableUpgradeable {
    /**
     * @notice The latest L1 block number known by the L2 system.
     */
    uint256 public number;

    /**
     * @notice The latest L1 block timestamp known by the L2 system.
     */
    uint256 public timestamp;

    /**
     * @notice The latest L1 base fee known by the L2 system.
     */
    uint256 public baseFee;

    /**
     * @notice The latest L1 block hash known by the L2 system.
     */
    bytes32 public hash;

    /**
     * @notice The latest 256 L1 stateRoots known by the L2 system.
     */
    mapping(uint8 => bytes32) public stateRoots;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initializer;
     */
    function initialize() public initializer {
        __Ownable_init();
        __Pausable_init();
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
     * @notice Updates the L1 block values.
     *
     * @param _number L1 block number.
     * @param _timestamp L1 timestamp.
     * @param _baseFee L1 baseFee.
     * @param _hash L1 block hash.
     * @param _stateRoot L1 stateRoot.
     */
    function setL1OracleValues(uint256 _number, uint256 _timestamp, uint256 _baseFee, bytes32 _hash, bytes32 _stateRoot)
        external
        onlyCoinbase
        whenNotPaused
    {
        require(number < _number, "Block number must be greater than the current block number.");
        number = _number;
        timestamp = _timestamp;
        baseFee = _baseFee;
        hash = _hash;
        stateRoots[uint8(_number % 256)] = _stateRoot;
    }

    /**
     * @notice Modifier to check if the caller is the coinbase.
     */
    modifier onlyCoinbase() {
        require(msg.sender == block.coinbase, "Only the coinbase can call this function.");
        _;
    }
}
