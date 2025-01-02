```diff
diff --git a/src/v0.8/ccip/pools/BurnMintTokenPoolAbstract.sol b/src/v0.8/ccip/pools/GHO/UpgradeableBurnMintTokenPoolAbstract.sol
index b3bbf4ff5e..2e90c6d4ea 100644
--- a/src/v0.8/ccip/pools/BurnMintTokenPoolAbstract.sol
+++ b/src/v0.8/ccip/pools/GHO/UpgradeableBurnMintTokenPoolAbstract.sol
@@ -1,18 +1,16 @@
 // SPDX-License-Identifier: BUSL-1.1
-pragma solidity 0.8.24;
+pragma solidity ^0.8.0;

-import {IBurnMintERC20} from "../../shared/token/ERC20/IBurnMintERC20.sol";
+import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

-import {Pool} from "../libraries/Pool.sol";
-import {TokenPool} from "./TokenPool.sol";
+import {Pool} from "../../libraries/Pool.sol";
+import {UpgradeableTokenPool} from "./UpgradeableTokenPool.sol";

-abstract contract BurnMintTokenPoolAbstract is TokenPool {
+abstract contract UpgradeableBurnMintTokenPoolAbstract is UpgradeableTokenPool {
   /// @notice Contains the specific burn call for a pool.
   /// @dev overriding this method allows us to create pools with different burn signatures
   /// without duplicating the underlying logic.
-  function _burn(
-    uint256 amount
-  ) internal virtual;
+  function _burn(uint256 amount) internal virtual;

   /// @notice Burn the token in the pool
   /// @dev The _validateLockOrBurn check is an essential security check
@@ -25,10 +23,11 @@ abstract contract BurnMintTokenPoolAbstract is TokenPool {

     emit Burned(msg.sender, lockOrBurnIn.amount);

-    return Pool.LockOrBurnOutV1({
-      destTokenAddress: getRemoteToken(lockOrBurnIn.remoteChainSelector),
-      destPoolData: _encodeLocalDecimals()
-    });
+    return
+      Pool.LockOrBurnOutV1({
+        destTokenAddress: getRemoteToken(lockOrBurnIn.remoteChainSelector),
+        destPoolData: _encodeLocalDecimals()
+      });
   }

   /// @notice Mint tokens from the pool to the recipient
@@ -39,8 +38,10 @@ abstract contract BurnMintTokenPoolAbstract is TokenPool {
     _validateReleaseOrMint(releaseOrMintIn);

     // Calculate the local amount
-    uint256 localAmount =
-      _calculateLocalAmount(releaseOrMintIn.amount, _parseRemoteDecimals(releaseOrMintIn.sourcePoolData));
+    uint256 localAmount = _calculateLocalAmount(
+      releaseOrMintIn.amount,
+      _parseRemoteDecimals(releaseOrMintIn.sourcePoolData)
+    );

     // Mint to the receiver
     IBurnMintERC20(address(i_token)).mint(releaseOrMintIn.receiver, localAmount);
```
