// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initiaalizable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract StakingSystem is
    Initializable,
    AccessControlUpgradeable,
    UUPSUpgradeable,
{
    struct Pool {
        // 质押代币的地址
        address stTokenAddress;
        // 质押池的权重，影响奖励的分配
        uint256 poolWeight;
        // 最后一次计算奖励的区块号
        uint256 lastRewardBlock;
        // 每一个质押代币累积的 MetaNode 数量
        uint256 accMetaNodePerST;
        // 总计质押代币的数量
        uint256 stTokenAmount;
        // 最小质押金额
        uint256 minDepositAmount;
        // 解除质押时锁定的区块数量
        uint256 unstakeLockedBlocks;
    }

    struct UnstakeRequest {
        // 解除质押的代币数量
        uint256 amount;
        // 可提取代币的区块高度
        uint256 unlockBlocks;
    }

    struct User {
        // 用户质押的代币数量
        uint256 stAmount;
        // 已分配的 MNODE 数量
        uint256 finishedMetaNode;
        // 待领取的 MNODE 数量
        uint256 pendingMetaNode;
        // 接触质押的请求列表
        UnstakeRequest[] requests;
    }

    /**
     * @notice 设置 MNODE 代币地址
     */
    function initialize (
        address _mnodeAddress,
        uint256 _startBlock,
        uint256 _endBlock,
        uint256 _metaNodePerBlock
    ) public initializer {
        require(_startBlock <= _endBlock && _metaNodePerBlock > 0, "invalid parameters");

        __A
    }
}
