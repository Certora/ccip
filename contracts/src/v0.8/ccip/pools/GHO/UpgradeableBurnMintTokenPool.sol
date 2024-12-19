// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import {UpgradeableBurnMintTokenPoolAbstract} from "./UpgradeableBurnMintTokenPoolAbstract.sol";
import {UpgradeableTokenPool} from "./UpgradeableTokenPool.sol";

import {Initializable} from "solidity-utils/contracts/transparent-proxy/Initializable.sol";
import {IRouter} from "../../interfaces/IRouter.sol";

/// @title UpgradeableBurnMintTokenPool
/// @author Aave Labs
/// @notice Upgradeable version of Chainlink's CCIP BurnMintTokenPool
/// @dev Contract adaptations:
/// - Implementation of Initializable to allow upgrades
/// - Move of allowlist and router definition to initialization stage
/// - Add GHO-Specific onlyOwner `transferLiquidity` which mints liquidity to the old pool
/// - Remove i_token decimal check in UpgradeableTokenPool constructor

/// @dev Pool whitelisting mode is set in the constructor and cannot be modified later.
/// It either accepts any address as originalSender, or only accepts whitelisted originalSender.
/// The only way to change whitelisting mode is to deploy a new pool.
/// If that is expected, please make sure the token's burner/minter roles are adjustable.
/// @dev This contract is a variant of BurnMintTokenPool that uses `burn(amount)`.
contract UpgradeableBurnMintTokenPool is Initializable, UpgradeableBurnMintTokenPoolAbstract, ITypeAndVersion {
  string public constant override typeAndVersion = "BurnMintTokenPool 1.5.1";

  constructor(
    address token,
    uint8 localTokenDecimals,
    address rmnProxy,
    bool allowListEnabled
  ) UpgradeableTokenPool(IBurnMintERC20(token), localTokenDecimals, rmnProxy, allowListEnabled) {}

  function initialize(address owner_, address[] memory allowlist, address router) external initializer {
    if (router == address(0) || owner_ == address(0)) revert ZeroAddressNotAllowed();

    _transferOwnership(owner_);
    s_router = IRouter(router);
    if (i_allowlistEnabled) _applyAllowListUpdates(new address[](0), allowlist);
  }

  /// @notice This function allows the owner to mint `amount` tokens on behalf of the pool and transfer them to `to`.
  /// This is GHO-Specific and is called to match the facilitator level of the new pool with the old pool such that
  /// it can burn the bridged supply once the old pool is deprecated. The old pool is then expected to burn `amount` of tokens
  /// so that it can be removed as a facilitator on GHO.
  /// @dev This is only called while offboarding an old token pool (or facilitator) in favor of this pool.
  /// @param to The address to which the minted tokens will be transferred. This needs to be the old token pool,
  /// or the facilitator being offboarded.
  /// @param amount The amount of tokens to mint and transfer to old pool.
  function mintAndTransferLiquidity(address to, uint256 amount) external onlyOwner {
    IBurnMintERC20(address(i_token)).mint(to, amount);
  }

  /// @inheritdoc UpgradeableBurnMintTokenPoolAbstract
  function _burn(uint256 amount) internal virtual override {
    IBurnMintERC20(address(i_token)).burn(amount);
  }
}
