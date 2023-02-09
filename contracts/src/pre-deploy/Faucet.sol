// SPDX-License-Identifier: UNLICENSED

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

import "@openzeppelin/contracts/access/Ownable.sol";

contract Faucet is Ownable {
    uint256 public amountAllowed;

    event LogDepositReceived(address, uint256);
    event LogRequestFunds(address, uint256);

    mapping(address => uint256) public lockTime;

    constructor() payable {}

    receive() external payable {
        emit LogDepositReceived(msg.sender, msg.value);
    }

    function retrieve() external onlyOwner {
        payable(msg.sender).transfer(address(this).balance);
    }

    function requestFunds(address payable _requestor) public payable onlyOwner {
        require(block.timestamp > lockTime[_requestor], "Lock time has not expired.");
        require(address(this).balance > amountAllowed, "Not enough funds in faucet.");

        _requestor.transfer(amountAllowed);

        lockTime[_requestor] = block.timestamp + 1 days;
        emit LogRequestFunds(_requestor, amountAllowed);
    }
}
