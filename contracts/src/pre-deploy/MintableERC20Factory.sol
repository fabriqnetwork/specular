// SPDX-License-Identifier: MIT

// The following code is a derivative work of the code from the Optimism contributors
// https://github.com/ethereum-optimism/optimism/blob/develop/packages/contracts-bedrock/contracts/universal/IOptimismMintableERC20.sol
// commit hash: c93958755b4f6ab7f95cc0b2459f39ca95c06684

pragma solidity ^0.8.4;

import {MintableERC20} from "../bridge/mintable/MintableERC20.sol";

/// @title MintableERC20Factory
/// @notice MintableERC20Factory is a factory contract that generates MintableERC20
///         contracts on the network it's deployed to. Simplifies the deployment process for users
///         who may be less familiar with deploying smart contracts. Designed to be backwards
///         compatible with the older StandardL2ERC20Factory contract.
contract MintableERC20Factory {
    /// @notice Address of the StandardBridge on this chain.
    address public immutable BRIDGE;

    /// @notice Emitted whenever a new MintableERC20 is created.
    /// @param localToken  Address of the created token on the local chain.
    /// @param remoteToken Address of the corresponding token on the remote chain.
    /// @param deployer    Address of the account that deployed the token.
    event MintableERC20Created(address indexed localToken, address indexed remoteToken, address deployer);

    /// @notice The semver MUST be bumped any time that there is a change in
    ///         the MintableERC20 token contract since this contract
    ///         is responsible for deploying MintableERC20 contracts.
    /// @param _bridge Address of the StandardBridge on this chain.
    constructor(address _bridge) {
        BRIDGE = _bridge;
    }

    /// @notice Creates an instance of the MintableERC20 contract.
    /// @param _remoteToken Address of the token on the remote chain.
    /// @param _name        ERC20 name.
    /// @param _symbol      ERC20 symbol.
    /// @return Address of the newly created token.
    function createMintableERC20(address _remoteToken, string memory _name, string memory _symbol)
        public
        returns (address)
    {
        require(_remoteToken != address(0), "MintableERC20Factory: must provide remote token address");

        address localToken = address(new MintableERC20(BRIDGE, _remoteToken, _name, _symbol));

        // Emit the updated event. The arguments here differ from the legacy event, but
        // are consistent with the ordering used in StandardBridge events.
        emit MintableERC20Created(localToken, _remoteToken, msg.sender);

        return localToken;
    }
}
