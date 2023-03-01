// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

contract L1Oracle is Initializable, UUPSUpgradeable, OwnableUpgradeable {
    /**
     * @notice The latest L1 stateRoot known by the L2 system.
     */
    bytes32 public stateRoot;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initializer;
     */
    function initialize() public initializer {
        __Ownable_init();
        __UUPSUpgradeable_init();
    }

    function _authorizeUpgrade(address) internal override onlyOwner {}

    /**
     * @notice Updates the L1 block values.
     *
     * @param _stateRoot   L1 stateRoot.
     */
    function setL1OracleValues(bytes32 _stateRoot) external onlyOwner {
        stateRoot = _stateRoot;
    }
}
