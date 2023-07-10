// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

/**
 * @title Types
 * @notice Contains various types used throughout the Specular contract system.
 */
library Types {
    /**
     * @notice Struct representing a cross domain message.
     */
    struct CrossDomainMessage {
        uint256 nonce;
        address sender;
        address target;
        uint256 value;
        uint256 gasLimit;
        bytes data;
    }
}
