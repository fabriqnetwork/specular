// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Types} from "./Types.sol";
import {Hashing} from "./Hashing.sol";
import {RLPReader} from "./rlp/RLPReader.sol";

library Encoding {
    /**
     * @notice Decode the state root from an encoded block header.
     *
     * @param encodedBlockHeader Address to call on L2 execution.
     * @return blockHash The block hash.
     * @return stateRoot The state root.
     */
    function decodeStateRootFromEncodedBlockHeader(bytes memory encodedBlockHeader)
        internal
        pure
        returns (bytes32, bytes32)
    {
        RLPReader.RLPItem memory item = RLPReader.toRlpItem(encodedBlockHeader);
        RLPReader.RLPItem[] memory blockHeader = RLPReader.toList(item);
        bytes32 blockHash = keccak256(encodedBlockHeader);
        // stateRoot is the 4th element in the block header
        bytes32 stateRoot = bytes32(RLPReader.toUintStrict(blockHeader[3]));
        return (blockHash, stateRoot);
    }

    /**
     * @notice Decode the storage root from an encoded account.
     *
     * @param encodedAccount Address to call on L2 execution.
     * @return storageRoot The storage root.
     */
    function decodeStorageRootFromEncodedAccount(bytes memory encodedAccount) internal pure returns (bytes32) {
        RLPReader.RLPItem memory item = RLPReader.toRlpItem(encodedAccount);
        RLPReader.RLPItem[] memory account = RLPReader.toList(item);
        // storageRoot is the 3rd element in the account
        bytes32 storageRoot = bytes32(RLPReader.toUintStrict(account[2]));
        return storageRoot;
    }
}
