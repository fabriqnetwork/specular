// SPDX-License-Identifier: MIT

// The following code is a derivative work of the code from the Optimism contributors
// https://github.com/ethereum-optimism/optimism/blob/develop/packages/contracts-bedrock/contracts/universal/IOptimismMintableERC20.sol
// commit hash: c93958755b4f6ab7f95cc0b2459f39ca95c06684

pragma solidity ^0.8.4;

import {IERC165} from "@openzeppelin/contracts/utils/introspection/IERC165.sol";

/// @title IMintableERC20
/// @notice This interface is available on the MintableERC20 contract.
///         We declare it as a separate interface so that it can be used in
///         custom implementations of MintableERC20.
interface IMintableERC20 is IERC165 {
    function REMOTE_TOKEN() external view returns (address);

    function BRIDGE() external returns (address);

    function mint(address _to, uint256 _amount) external;

    function burn(address _from, uint256 _amount) external;
}
