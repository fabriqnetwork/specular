// SPDX-License-Identifier: Apache-2.0

/*
 * Modifications Copyright 2022, Specular contributors
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

import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

contract Faucet is UUPSUpgradeable, OwnableUpgradeable {
    function initialize() public initializer {
        __Ownable_init();
    }

    function _authorizeUpgrade(address) internal override onlyOwner {}

    // amountAllowed is initialized in the genesis.json
    // amountAllowed is set to 0.01 ETH
    uint256 public amountAllowed;

    mapping(address => uint256) public lockTime;

    event DepositReceived(address, uint256);
    event RequestFunds(address, uint256);

    receive() external payable {
        emit DepositReceived(msg.sender, msg.value);
    }

    function retrieve() external onlyOwner {
        payable(msg.sender).transfer(address(this).balance);
    }

    function requestFunds(address payable requestor) public payable onlyOwner {
        require(block.timestamp > lockTime[requestor], "Lock time has not expired.");
        require(address(this).balance > amountAllowed, "Not enough funds in faucet.");

        lockTime[requestor] = block.timestamp + 1 days;
        requestor.transfer(amountAllowed);

        emit RequestFunds(requestor, amountAllowed);
    }
}
