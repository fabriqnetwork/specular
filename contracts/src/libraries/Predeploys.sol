// SPDX-License-Identifier: Apache-2.0

/*
 * Copyright 2022, Specular contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

pragma solidity ^0.8.0;

library Predeploys {
    // @notice Address of the L1Oracle predeploy.
    address internal constant L1_ORACLE = 0x2a00000000000000000000000000000000000010;

    // @notice Address of the L2Portal predeploy.
    address internal constant L2_PORTAL = 0x2A00000000000000000000000000000000000011;

    // @notice Address of the L2StandardBridge predeploy.
    address internal constant L2_STANDARD_BRIDGE = 0x2a00000000000000000000000000000000000012;
}
