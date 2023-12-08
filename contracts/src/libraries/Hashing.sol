// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {Types} from "./Types.sol";
import {Encoding} from "./Encoding.sol";

/**
 * @title Hashing
 * @notice Hashing handles Specular's various different hashing schemes.
 */
library Hashing {
    /**
     * @notice Derives the withdrawal hash according to the encoding in the L2 Withdrawer contract
     *
     * @param _tx Withdrawal transaction to hash.
     *
     * @return Hashed withdrawal transaction.
     */
    function hashCrossDomainMessage(Types.CrossDomainMessage memory _tx) internal pure returns (bytes32) {
        return keccak256(abi.encode(_tx.nonce, _tx.sender, _tx.target, _tx.value, _tx.gasLimit, _tx.data));
    }

    bytes32 public constant STATE_COMMITMENT_V0 = bytes32(0);

    /**
     * @notice creates a versioned state commitment
     *
     * @param l2VmHash l2 vm hash
     */
    function createStateCommitmentV0(bytes32 l2VmHash) internal pure returns (bytes32) {
        // state commitment v0 format is keccak256(version || vmHash)
        bytes memory stateCommitment = bytes.concat(STATE_COMMITMENT_V0, l2VmHash);
        return keccak256(stateCommitment);
    }
}
