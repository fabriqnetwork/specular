// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Testing utilities
import {Test, StdUtils} from "forge-std/Test.sol";
import "forge-std/console.sol";

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {L1Oracle} from "../../src/bridge/L1Oracle.sol";
import {L1Portal} from "../../src/bridge/L1Portal.sol";
import {L2Portal} from "../../src/bridge/L2Portal.sol";
import {L1StandardBridge} from "../../src/bridge/L1StandardBridge.sol";
import {L2StandardBridge} from "../../src/bridge/L2StandardBridge.sol";
import {Predeploys} from "../../src/libraries/Predeploys.sol";
import {MintableERC20} from "../../src/bridge/mintable/MintableERC20.sol";
import {MintableERC20Factory} from "../../src/bridge/mintable/MintableERC20Factory.sol";

contract UUPSProxy is ERC1967Proxy {
    constructor(address _implementation, bytes memory _data) ERC1967Proxy(_implementation, _data) {}
}

contract CommonTest is Test {
    address alice = address(1);
    address bob = address(2);

    function setUp() public virtual {
        vm.deal(alice, 1 << 16);
        vm.deal(bob, 1 << 16);
    }
}

contract L1Oracle_Initializer is CommonTest {
    address sequencerAddress = address(8);
    address l1OracleAddress = Predeploys.L1_ORACLE;

    L1Oracle oracle;

    function setUp() public virtual override {
        super.setUp();

        vm.etch(l1OracleAddress, address(new L1Oracle()).code);
        vm.label(l1OracleAddress, "L1Oracle");

        oracle = L1Oracle(l1OracleAddress);
        oracle.initialize();
    }
}

contract Portal_Initializer is L1Oracle_Initializer {
    address payable l1PortalAddress;
    address payable l2PortalAddress;

    L1Portal l1Portal;
    L2Portal l2Portal;

    UUPSProxy l1PortalProxy;
    UUPSProxy l2PortalProxy;

    event DepositInitiated(
        uint256 indexed nonce,
        address indexed sender,
        address indexed target,
        uint256 value,
        uint256 gasLimit,
        bytes data,
        bytes32 depositHash
    );

    event WithdrawalInitiated(
        uint256 indexed nonce,
        address indexed sender,
        address indexed target,
        uint256 value,
        uint256 gasLimit,
        bytes data,
        bytes32 withdrawalHash
    );

    event WithdrawalFinalized(bytes32 indexed withdrawalHash, bool success);

    function setUp() public virtual override {
        super.setUp();

        l1PortalProxy = new UUPSProxy(address(new L1Portal()), "");
        l1PortalAddress = payable(address(l1PortalProxy));
        vm.label(l1PortalAddress, "L1Portal");
        l1Portal = L1Portal(l1PortalAddress);

        // dummy rollup address
        l1Portal.initialize(address(42));

        // the L2Portal needs to be deployed at the correct predeploy address
        // deploy implementation contract
        address l2PortalImplAddress = address(new L2Portal());

        // deploy Proxy at random address
        l2PortalProxy = new UUPSProxy(l2PortalImplAddress, "");

        // etch proxy to the predeploy address
        l2PortalAddress = payable(Predeploys.L2_PORTAL);
        vm.etch(l2PortalAddress, address(l2PortalProxy).code);
        vm.label(l2PortalAddress, "L2Portal");

        // store implementation address at the correct storage slot
        // the slot is defined in https://eips.ethereum.org/EIPS/eip-1967
        // bytes32(uint256(keccak256('eip1967.proxy.implementation')) - 1)
        bytes32 implSlot = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
        vm.store(l2PortalAddress, implSlot, bytes32(uint256(uint160(l2PortalImplAddress))));

        l2Portal = L2Portal(l2PortalAddress);
        l2Portal.initialize(l1PortalAddress);
    }
}

contract StandardBridge_Initializer is Portal_Initializer {
    address payable l1StandardBridgeAddress = payable(address(128));
    address payable l2StandardBridgeAddress = payable(Predeploys.L2_STANDARD_BRIDGE);

    L1StandardBridge l1StandardBridge;
    L2StandardBridge l2StandardBridge;

    MintableERC20Factory l1TokenFactory;
    MintableERC20Factory l2TokenFactory;
    ERC20 l1Token;
    ERC20 badL1Token;
    MintableERC20 l2Token;
    ERC20 nativeL2Token;
    ERC20 badL2Token;
    MintableERC20 remoteL1Token;

    event ETHBridgeInitiated(address indexed from, address indexed to, uint256 amount, bytes extraData);

    event ERC20BridgeInitiated(
        address indexed localToken,
        address indexed remoteToken,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
    );

    event ETHBridgeFinalized(address indexed from, address indexed to, uint256 amount, bytes extraData);

    event ERC20BridgeFinalized(
        address indexed localToken,
        address indexed remoteToken,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
    );

    function setUp() public virtual override {
        super.setUp();

        // TODO: deploy StandardBridge contracts behind proxy
        vm.etch(l1StandardBridgeAddress, address(new L1StandardBridge()).code);
        vm.label(l1StandardBridgeAddress, "L1StandardBridge");

        l1StandardBridge = L1StandardBridge(payable(l1StandardBridgeAddress));
        l1StandardBridge.initialize(l1PortalAddress);

        vm.etch(l2StandardBridgeAddress, address(new L2StandardBridge()).code);
        vm.label(l2StandardBridgeAddress, "L2StandardBridge");

        l2StandardBridge = L2StandardBridge(l2StandardBridgeAddress);
        l2StandardBridge.initialize(l1StandardBridgeAddress);

        // depoly ERC20 tokens for testing
        l2TokenFactory = new MintableERC20Factory(
            l2StandardBridgeAddress
        );

        l1Token = new ERC20("Native L1 Token", "L1T");

        l2Token = MintableERC20(l2TokenFactory.createMintableERC20(address(l1Token), "Remote L1 Token", "rL1T"));

        badL2Token = MintableERC20(l2TokenFactory.createMintableERC20(address(1), "Bad L2 Token", "bL2T"));
    }
}
