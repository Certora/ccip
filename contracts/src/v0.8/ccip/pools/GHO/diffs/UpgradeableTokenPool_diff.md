```diff
diff --git a/src/v0.8/ccip/pools/TokenPool.sol b/src/v0.8/ccip/pools/GHO/UpgradeableTokenPool.sol
index cd3096f4ef..8c0965a67f 100644
--- a/src/v0.8/ccip/pools/TokenPool.sol
+++ b/src/v0.8/ccip/pools/GHO/UpgradeableTokenPool.sol
@@ -1,26 +1,29 @@
 // SPDX-License-Identifier: BUSL-1.1
-pragma solidity 0.8.24;
+pragma solidity ^0.8.0;

-import {IPoolV1} from "../interfaces/IPool.sol";
-import {IRMN} from "../interfaces/IRMN.sol";
-import {IRouter} from "../interfaces/IRouter.sol";
+import {IPoolV1} from "../../interfaces/IPool.sol";
+import {IRMN} from "../../interfaces/IRMN.sol";
+import {IRouter} from "../../interfaces/IRouter.sol";

-import {Ownable2StepMsgSender} from "../../shared/access/Ownable2StepMsgSender.sol";
-import {Pool} from "../libraries/Pool.sol";
-import {RateLimiter} from "../libraries/RateLimiter.sol";
+import {Ownable2StepMsgSender} from "../../../shared/access/Ownable2StepMsgSender.sol";
+import {Pool} from "../../libraries/Pool.sol";
+import {RateLimiter} from "../../libraries/RateLimiter.sol";

-import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
-import {IERC20Metadata} from
-  "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol";
-import {IERC165} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";
-import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";
+import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
+import {IERC165} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";
+import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";

+/// @title UpgradeableTokenPool
+/// @author Aave Labs
+/// @notice Upgradeable version of Chainlink's CCIP TokenPool
 /// @dev This pool supports different decimals on different chains but using this feature could impact the total number
 /// of tokens in circulation. Since all of the tokens are locked/burned on the source, and a rounded amount is minted/released on the
 /// destination, the number of tokens minted/released could be less than the number of tokens burned/locked. This is because the source
 /// chain does not know about the destination token decimals. This is not a problem if the decimals are the same on both
 /// chains.
-///
+/// @dev Contract adaptations:
+///  - Remove i_token decimal check in constructor.
+///  - Add storage `__gap` for future upgrades.
 /// Example:
 /// Assume there is a token with 6 decimals on chain A and 3 decimals on chain B.
 /// - 1.234567 tokens are burned on chain A.
@@ -29,7 +32,7 @@ import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts
 /// 0.000567 tokens.
 /// In the case of a burnMint pool on chain A, these funds are burned in the pool on chain A.
 /// In the case of a lockRelease pool on chain A, these funds accumulate in the pool on chain A.
-abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
+abstract contract UpgradeableTokenPool is IPoolV1, Ownable2StepMsgSender {
   using EnumerableSet for EnumerableSet.Bytes32Set;
   using EnumerableSet for EnumerableSet.AddressSet;
   using EnumerableSet for EnumerableSet.UintSet;
@@ -117,34 +120,18 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
   /// @dev Can be address(0) if none is configured.
   address internal s_rateLimitAdmin;

-  constructor(IERC20 token, uint8 localTokenDecimals, address[] memory allowlist, address rmnProxy, address router) {
-    if (address(token) == address(0) || router == address(0) || rmnProxy == address(0)) revert ZeroAddressNotAllowed();
+  constructor(IERC20 token, uint8 localTokenDecimals, address rmnProxy, bool allowListEnabled) {
+    if (address(token) == address(0) || rmnProxy == address(0)) revert ZeroAddressNotAllowed();
     i_token = token;
     i_rmnProxy = rmnProxy;
-
-    try IERC20Metadata(address(token)).decimals() returns (uint8 actualTokenDecimals) {
-      if (localTokenDecimals != actualTokenDecimals) {
-        revert InvalidDecimalArgs(localTokenDecimals, actualTokenDecimals);
-      }
-    } catch {
-      // The decimals function doesn't exist, which is possible since it's optional in the ERC20 spec. We skip the check and
-      // assume the supplied token decimals are correct.
-    }
     i_tokenDecimals = localTokenDecimals;

-    s_router = IRouter(router);
-
     // Pool can be set as permissioned or permissionless at deployment time only to save hot-path gas.
-    i_allowlistEnabled = allowlist.length > 0;
-    if (i_allowlistEnabled) {
-      _applyAllowListUpdates(new address[](0), allowlist);
-    }
+    i_allowlistEnabled = allowListEnabled;
   }

   /// @inheritdoc IPoolV1
-  function isSupportedToken(
-    address token
-  ) public view virtual returns (bool) {
+  function isSupportedToken(address token) public view virtual returns (bool) {
     return token == address(i_token);
   }

@@ -168,9 +155,7 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {

   /// @notice Sets the pool's Router
   /// @param newRouter The new Router
-  function setRouter(
-    address newRouter
-  ) public onlyOwner {
+  function setRouter(address newRouter) public onlyOwner {
     if (newRouter == address(0)) revert ZeroAddressNotAllowed();
     address oldRouter = address(s_router);
     s_router = IRouter(newRouter);
@@ -179,11 +164,11 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
   }

   /// @notice Signals which version of the pool interface is supported
-  function supportsInterface(
-    bytes4 interfaceId
-  ) public pure virtual override returns (bool) {
-    return interfaceId == Pool.CCIP_POOL_V1 || interfaceId == type(IPoolV1).interfaceId
-      || interfaceId == type(IERC165).interfaceId;
+  function supportsInterface(bytes4 interfaceId) public pure virtual override returns (bool) {
+    return
+      interfaceId == Pool.CCIP_POOL_V1 ||
+      interfaceId == type(IPoolV1).interfaceId ||
+      interfaceId == type(IERC165).interfaceId;
   }

   // ================================================================
@@ -199,9 +184,7 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
   /// @param lockOrBurnIn The input to validate.
   /// @dev This function should always be called before executing a lock or burn. Not doing so would allow
   /// for various exploits.
-  function _validateLockOrBurn(
-    Pool.LockOrBurnInV1 calldata lockOrBurnIn
-  ) internal {
+  function _validateLockOrBurn(Pool.LockOrBurnInV1 calldata lockOrBurnIn) internal {
     if (!isSupportedToken(lockOrBurnIn.localToken)) revert InvalidToken(lockOrBurnIn.localToken);
     if (IRMN(i_rmnProxy).isCursed(bytes16(uint128(lockOrBurnIn.remoteChainSelector)))) revert CursedByRMN();
     _checkAllowList(lockOrBurnIn.originalSender);
@@ -219,9 +202,7 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
   /// @param releaseOrMintIn The input to validate.
   /// @dev This function should always be called before executing a release or mint. Not doing so would allow
   /// for various exploits.
-  function _validateReleaseOrMint(
-    Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
-  ) internal {
+  function _validateReleaseOrMint(Pool.ReleaseOrMintInV1 calldata releaseOrMintIn) internal {
     if (!isSupportedToken(releaseOrMintIn.localToken)) revert InvalidToken(releaseOrMintIn.localToken);
     if (IRMN(i_rmnProxy).isCursed(bytes16(uint128(releaseOrMintIn.remoteChainSelector)))) revert CursedByRMN();
     _onlyOffRamp(releaseOrMintIn.remoteChainSelector);
@@ -247,9 +228,7 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
     return abi.encode(i_tokenDecimals);
   }

-  function _parseRemoteDecimals(
-    bytes memory sourcePoolData
-  ) internal view virtual returns (uint8) {
+  function _parseRemoteDecimals(bytes memory sourcePoolData) internal view virtual returns (uint8) {
     // Fallback to the local token decimals if the source pool data is empty. This allows for backwards compatibility.
     if (sourcePoolData.length == 0) {
       return i_tokenDecimals;
@@ -304,9 +283,7 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
   /// @notice Gets the pool address on the remote chain.
   /// @param remoteChainSelector Remote chain selector.
   /// @dev To support non-evm chains, this value is encoded into bytes
-  function getRemotePools(
-    uint64 remoteChainSelector
-  ) public view returns (bytes[] memory) {
+  function getRemotePools(uint64 remoteChainSelector) public view returns (bytes[] memory) {
     bytes32[] memory remotePoolHashes = s_remoteChainConfigs[remoteChainSelector].remotePools.values();

     bytes[] memory remotePools = new bytes[](remotePoolHashes.length);
@@ -327,9 +304,7 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
   /// @notice Gets the token address on the remote chain.
   /// @param remoteChainSelector Remote chain selector.
   /// @dev To support non-evm chains, this value is encoded into bytes
-  function getRemoteToken(
-    uint64 remoteChainSelector
-  ) public view returns (bytes memory) {
+  function getRemoteToken(uint64 remoteChainSelector) public view returns (bytes memory) {
     return s_remoteChainConfigs[remoteChainSelector].remoteTokenAddress;
   }

@@ -358,9 +333,7 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
   }

   /// @inheritdoc IPoolV1
-  function isSupportedChain(
-    uint64 remoteChainSelector
-  ) public view returns (bool) {
+  function isSupportedChain(uint64 remoteChainSelector) public view returns (bool) {
     return s_remoteChainSelectors.contains(remoteChainSelector);
   }

@@ -379,8 +352,7 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
   /// @notice Sets the permissions for a list of chains selectors. Actual senders for these chains
   /// need to be allowed on the Router to interact with this pool.
   /// @param remoteChainSelectorsToRemove A list of chain selectors to remove.
-  /// @param chainsToAdd A list of chains and their new permission status & rate limits. Rate limits
-  /// are only used when the chain is being added through `allowed` being true.
+  /// @param chainsToAdd A list of chains and their new permission status & rate limits.
   /// @dev Only callable by the owner
   function applyChainUpdates(
     uint64[] calldata remoteChainSelectorsToRemove,
@@ -495,9 +467,7 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
   /// @notice Sets the rate limiter admin address.
   /// @dev Only callable by the owner.
   /// @param rateLimitAdmin The new rate limiter admin address.
-  function setRateLimitAdmin(
-    address rateLimitAdmin
-  ) external onlyOwner {
+  function setRateLimitAdmin(address rateLimitAdmin) external onlyOwner {
     s_rateLimitAdmin = rateLimitAdmin;
     emit RateLimitAdminSet(rateLimitAdmin);
   }
@@ -566,18 +536,14 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {

   /// @notice Checks whether remote chain selector is configured on this contract, and if the msg.sender
   /// is a permissioned onRamp for the given chain on the Router.
-  function _onlyOnRamp(
-    uint64 remoteChainSelector
-  ) internal view {
+  function _onlyOnRamp(uint64 remoteChainSelector) internal view {
     if (!isSupportedChain(remoteChainSelector)) revert ChainNotAllowed(remoteChainSelector);
     if (!(msg.sender == s_router.getOnRamp(remoteChainSelector))) revert CallerIsNotARampOnRouter(msg.sender);
   }

   /// @notice Checks whether remote chain selector is configured on this contract, and if the msg.sender
   /// is a permissioned offRamp for the given chain on the Router.
-  function _onlyOffRamp(
-    uint64 remoteChainSelector
-  ) internal view {
+  function _onlyOffRamp(uint64 remoteChainSelector) internal view {
     if (!isSupportedChain(remoteChainSelector)) revert ChainNotAllowed(remoteChainSelector);
     if (!s_router.isOffRamp(remoteChainSelector, msg.sender)) revert CallerIsNotARampOnRouter(msg.sender);
   }
@@ -586,9 +552,7 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
   // │                          Allowlist                           │
   // ================================================================

-  function _checkAllowList(
-    address sender
-  ) internal view {
+  function _checkAllowList(address sender) internal view {
     if (i_allowlistEnabled) {
       if (!s_allowlist.contains(sender)) {
         revert SenderNotAllowed(sender);
@@ -635,4 +599,8 @@ abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
       }
     }
   }
+
+  /// @dev This empty reserved space is put in place to allow future versions to add new
+  /// variables without shifting down storage in the inheritance chain.
+  uint256[42] private __gap;
 }

```
