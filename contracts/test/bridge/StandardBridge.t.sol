// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {StandardBridge} from "../../src/bridge/StandardBridge.sol";
import {MintableERC20} from "../../src/bridge/mintable/MintableERC20.sol";
import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/// @title StandardBridgeTester
/// @notice Simple wrapper around the StandardBridge contract that exposes
///         internal functions so they can be more easily tested directly.
contract StandardBridgeTester is StandardBridge {
    constructor(address payable _messenger, address payable _otherBridge) {
        __StandardBridge_init(_messenger, _otherBridge);
    }

    modifier onlyOtherBridge() override {
        _;
    }

    function isCorrectTokenPair(address _mintableToken, address _otherToken) external view returns (bool) {
        return _isNonNativeTokenPair(_mintableToken, _otherToken);
    }

    receive() external payable override {}
}

/// @title StandardBridge_Stateless_Test
/// @notice Tests internal functions that require no existing state or contract
///         interactions with the messenger.
contract StandardBridge_Stateless_Test is Test {
    StandardBridgeTester internal bridge;
    MintableERC20 internal mintable;
    ERC20 internal erc20;

    function setUp() public {
        bridge = new StandardBridgeTester({
            _messenger: payable(address(0)),
            _otherBridge: payable(address(0))
        });

        mintable = new MintableERC20({
            _bridge: address(0),
            _remoteToken: address(0),
            _name: "Stonks",
            _symbol: "STONK"
        });

        erc20 = new ERC20("Altcoin", "ALT");
    }

    /// @notice Test coverage of isCorrectTokenPair under different types of
    ///         tokens.
    function test_isCorrectTokenPair_succeeds() external {
        // known to be correct remote token
        assertTrue(bridge.isCorrectTokenPair(address(mintable), mintable.REMOTE_TOKEN()));
        // known to be incorrect remote token
        assertTrue(mintable.REMOTE_TOKEN() != address(0x20));
        assertFalse(bridge.isCorrectTokenPair(address(mintable), address(0x20)));
        // A token that doesn't support the mintable interface will revert
        assertFalse(bridge.isCorrectTokenPair(address(erc20), address(1)));
    }
}
