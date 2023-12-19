// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {FeeVault} from "src/libraries/FeeVault.sol";

/// @custom:proxied
/// @custom:predeploy 0x2A00000000000000000000000000000000000021
/// @title L2BaseFeeVault
/// @notice The L2BaseFeeVault accumulates the base fee component of L2 transaction fees.
contract L2BaseFeeVault is FeeVault {}
