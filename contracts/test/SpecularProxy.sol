// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract SpecularProxy is ERC1967Proxy {
    constructor(address _implementationAddress, bytes memory _data) ERC1967Proxy(_implementationAddress, _data) {
        // No need to add any more logic here.
    }
}
