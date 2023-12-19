// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {FeeVault} from "src/libraries/FeeVault.sol";

/// @custom:proxied
/// @custom:predeploy 0x2A00000000000000000000000000000000000020
/// @title L1FeeVault
/// @notice The L1FeeVault accumulates the L1 fee component of L2 transaction fees.
contract L1FeeVault is FeeVault {}
