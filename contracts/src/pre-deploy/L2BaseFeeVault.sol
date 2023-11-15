// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {FeeVault} from "src/libraries/FeeVault.sol";

/// @custom:proxied
/// @custom:predeploy 0x2A00000000000000000000000000000000000020
/// @title L2BaseFeeVault
/// @notice The L2BaseFeeVault accumulates the base fee that is paid by transactions.
contract L2BaseFeeVault is FeeVault {}
