// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IPool} from "../../interfaces/IPool.sol";

import {Internal} from "../../libraries/Internal.sol";
import {EVM2EVMOffRamp} from "../../offRamp/EVM2EVMOffRamp.sol";
import {EVM2EVMOnRamp} from "../../onRamp/EVM2EVMOnRamp.sol";
import {BurnMintTokenPool} from "../../pools/BurnMintTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {BaseTest} from "../BaseTest.t.sol";
import {BurnMintSetup} from "./BurnMintSetup.t.sol";

contract BurnMintTokenPoolSetup is BurnMintSetup {
  BurnMintTokenPool internal s_pool;

  function setUp() public virtual override {
    BurnMintSetup.setUp();

    s_pool = new BurnMintTokenPool(s_burnMintERC677, new address[](0), address(s_mockARM), address(s_sourceRouter));
    s_burnMintERC677.grantMintAndBurnRoles(address(s_pool));

    _applyChainUpdates(address(s_pool));
  }
}

contract BurnMintTokenPool_lockOrBurn is BurnMintTokenPoolSetup {
  function test_Setup_Success() public view {
    assertEq(address(s_burnMintERC677), address(s_pool.getToken()));
    assertEq(address(s_mockARM), s_pool.getArmProxy());
    assertEq(false, s_pool.getAllowListEnabled());
    assertEq("BurnMintTokenPool 1.4.0", s_pool.typeAndVersion());
  }

  function test_PoolBurn_Success() public {
    uint256 burnAmount = 20_000e18;

    deal(address(s_burnMintERC677), address(s_pool), burnAmount);
    assertEq(s_burnMintERC677.balanceOf(address(s_pool)), burnAmount);

    vm.startPrank(s_burnMintOnRamp);

    vm.expectEmit();
    emit TokensConsumed(burnAmount);

    vm.expectEmit();
    emit Transfer(address(s_pool), address(0), burnAmount);

    vm.expectEmit();
    emit Burned(address(s_burnMintOnRamp), burnAmount);

    bytes4 expectedSignature = bytes4(keccak256("burn(uint256)"));
    vm.expectCall(address(s_burnMintERC677), abi.encodeWithSelector(expectedSignature, burnAmount));

    s_pool.lockOrBurn(OWNER, bytes(""), burnAmount, DEST_CHAIN_SELECTOR, bytes(""));

    assertEq(s_burnMintERC677.balanceOf(address(s_pool)), 0);
  }

  // Should not burn tokens if cursed.
  function test_PoolBurnRevertNotHealthy_Revert() public {
    s_mockARM.voteToCurse(bytes32(0));
    uint256 before = s_burnMintERC677.balanceOf(address(s_pool));
    vm.startPrank(s_burnMintOnRamp);

    vm.expectRevert(EVM2EVMOnRamp.BadARMSignal.selector);
    s_pool.lockOrBurn(OWNER, bytes(""), 1e5, DEST_CHAIN_SELECTOR, bytes(""));

    assertEq(s_burnMintERC677.balanceOf(address(s_pool)), before);
  }

  function test_ChainNotAllowed_Revert() public {
    uint64 wrongChainSelector = 8838833;

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, wrongChainSelector));
    s_pool.lockOrBurn(OWNER, bytes(""), 1, wrongChainSelector, bytes(""));
  }
}

contract BurnMintTokenPool_releaseOrMint is BurnMintTokenPoolSetup {
  function test_PoolMint_Success() public {
    uint256 amount = 1e19;

    vm.startPrank(s_burnMintOffRamp);

    vm.expectEmit();
    emit Transfer(address(0), OWNER, amount);
    s_pool.releaseOrMint(
      bytes(""),
      OWNER,
      amount,
      DEST_CHAIN_SELECTOR,
      IPool.SourceTokenData({
        sourcePoolAddress: abi.encode(s_remoteBurnMintPool),
        destPoolAddress: abi.encode(address(s_pool)),
        extraData: ""
      }),
      ""
    );

    assertEq(s_burnMintERC677.balanceOf(OWNER), amount);
  }

  function test_PoolMintNotHealthy_Revert() public {
    // Should not mint tokens if cursed.
    s_mockARM.voteToCurse(bytes32(0));
    uint256 before = s_burnMintERC677.balanceOf(OWNER);
    vm.startPrank(s_burnMintOffRamp);

    vm.expectRevert(EVM2EVMOffRamp.BadARMSignal.selector);
    s_pool.releaseOrMint(bytes(""), OWNER, 1e5, DEST_CHAIN_SELECTOR, generateSourceTokenData(), bytes(""));

    assertEq(s_burnMintERC677.balanceOf(OWNER), before);
  }

  function test_ChainNotAllowed_Revert() public {
    uint64 wrongChainSelector = 8838833;

    vm.expectRevert(abi.encodeWithSelector(TokenPool.ChainNotAllowed.selector, wrongChainSelector));
    s_pool.releaseOrMint(bytes(""), OWNER, 1, wrongChainSelector, generateSourceTokenData(), bytes(""));
  }
}
