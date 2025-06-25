// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/utils/math/Math.sol";

import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract StakingSystem is
    Initializable,
    AccessControlUpgradeable,
    PausableUpgradeable,
    UUPSUpgradeable
{
    using Math for uint256;

    bytes32 public constant ADMIN_ROLE = keccak256("admin_role");
    bytes32 public constant UPGRADE_ROLE = keccak256("upgrade_role");
    uint256 public constant ETH_PID = 0;

    IERC20 public metaNodeToken;
    uint256 public startBlock;
    uint256 public endBlock;
    uint256 public metaNodePerBlock;

    // 暂停状态变量
    bool public withdrawPaused;
    bool public claimPaused;

    // 总质押池权重 / 所有质押池权重之和
    uint256 public totalPoolWeight;
    Pool[] public pool;
    // pool id => user address => user info
    mapping(uint256 => mapping(address => User)) public user;

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
    function initialize(
        address _mnodeAddress,
        uint256 _startBlock,
        uint256 _endBlock,
        uint256 _metaNodePerBlock
    ) public initializer {
        require(
            _startBlock <= _endBlock && _metaNodePerBlock > 0,
            "invalid parameters"
        );

        __AccessControl_init();
        __UUPSUpgradeable_init();
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADE_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);
        setMetaNode(_mnodeAddress);
        startBlock = _startBlock;
        endBlock = _endBlock;
        metaNodePerBlock = _metaNodePerBlock;
    }

    function _authorizeUpgrade(
        address newImplementation
    ) internal override onlyRole(UPGRADE_ROLE) {}

    function setMetaNode(address _metaNode) public onlyRole(ADMIN_ROLE) {
        metaNodeToken = IERC20(_metaNode);
    }

    /**
     * @notice 添加新的质押池
     * @param _stTokenAddress 质押代币地址
     * @param _poolWeight 质押池权重
     * @param _minDepositAmount 质押最小金额
     * @param _unstackLockedBlocks 解除质押时锁定的区块数量
     * @param _withUpdate 是否更新
     */
    function addPool(
        address _stTokenAddress,
        uint256 _poolWeight,
        uint256 _minDepositAmount,
        uint256 _unstackLockedBlocks,
        bool _withUpdate
    ) public onlyRole(ADMIN_ROLE) {
        // 默认第一个质押池是 ETH 质押池
        if (pool.length > 0) {
            require(_stTokenAddress != address(0), "invalid _stTokenAddress");
        } else {
            require(_stTokenAddress == address(0), "invalid _stTokenAddress");
        }
        require(_minDepositAmount > 0, "invalid _minDepositAmount");
        require(_unstackLockedBlocks > 0, "invalid _unstackLockedBlocks");
        require(block.number < endBlock, "staking is over");

        if (_withUpdate) {
            massUpdatePools();
        }

        uint256 lastRewardBlock = block.number > startBlock
            ? block.number
            : startBlock;
        totalPoolWeight = totalPoolWeight + _poolWeight;

        pool.push(
            Pool({
                stTokenAddress: _stTokenAddress,
                poolWeight: _poolWeight,
                lastRewardBlock: lastRewardBlock,
                accMetaNodePerST: 0,
                stTokenAmount: 0,
                minDepositAmount: _minDepositAmount,
                unstakeLockedBlocks: _unstackLockedBlocks
            })
        );
    }

    function massUpdatePools() public {
        uint256 length = pool.length;
        for (uint256 i = 0; i < length; i++) {
            updatePool(i);
        }
    }

    modifier checkPid(uint256 _pid) {
        require(_pid < pool.length, "invalid _pid");
        _;
    }

    function updatePool(uint256 _pid) public checkPid(_pid) {
        Pool storage pool_ = pool[_pid];

        if (block.number <= pool_.lastRewardBlock) {
            return;
        }

        // 按照 poolWeight 计算当前质押池应该获得的奖励
        (bool success1, uint256 totalMetaNode) = getMultiplier(
            pool_.lastRewardBlock,
            block.number
        ).tryMul(pool_.poolWeight);
        require(success1, "overflow");

        (success1, totalMetaNode) = totalMetaNode.tryDiv(totalPoolWeight);
        require(success1, "overflow");

        uint256 stSupply = pool_.stTokenAmount;
        if (stSupply > 0) {
            // 扩大 1 ether 倍数
            (bool success2, uint256 totalMetaNode_) = totalMetaNode.tryMul(
                1 ether
            );
            require(success2, "overflow");
            // 计算每一个单位的质押代币应该获得的奖励
            (success2, totalMetaNode_) = totalMetaNode_.tryDiv(stSupply);
            require(success2, "overflow");
            // 累计计算每一个质押代币的奖励
            (bool success3, uint256 accMetaNodePerST) = pool_
                .accMetaNodePerST
                .tryAdd(totalMetaNode_);
            require(success3, "overflow");
            pool_.accMetaNodePerST = accMetaNodePerST;
        }
        pool_.lastRewardBlock = block.number;
    }

    /**
     * @notice 获取乘数
     * @param _from 起始区块（包含）
     * @param _to 结束区块（不包含）
     */
    function getMultiplier(
        uint256 _from,
        uint256 _to
    ) public view returns (uint256 multiplier) {
        require(_from <= _to, "invalid block");

        if (_from < startBlock) {
            _from = startBlock;
        }
        if (_to > endBlock) {
            _to = endBlock;
        }

        bool success;
        (success, multiplier) = (_to - _from).tryMul(metaNodePerBlock);
        require(success, "overflow");
    }

    function depositETH() public payable whenNotPaused {
        Pool storage pool_ = pool[ETH_PID];
        require(
            pool_.stTokenAddress == address(0x0),
            "invalid staking token address"
        );

        uint256 _amount = msg.value;
        require(
            _amount >= pool_.minDepositAmount,
            "deposit amount is too small"
        );

        _deposit(ETH_PID, _amount);
    }

    function deposit(
        uint256 _pid,
        uint256 _amount
    ) public whenNotPaused checkPid(_pid) {
        require(_pid != 0, "deposit not support ETH staking");
        Pool storage pool_ = pool[_pid];
        require(
            _amount >= pool_.minDepositAmount,
            "deposit amount is too small"
        );

        if (_amount > 0) {
            IERC20(pool_.stTokenAddress).transferFrom(
                msg.sender,
                address(this),
                _amount
            );
        }

        _deposit(_pid, _amount);
    }

    function _deposit(uint256 _pid, uint256 _amount) internal {
        Pool storage pool_ = pool[_pid];
        User storage user_ = user[_pid][msg.sender];

        updatePool(_pid);

        // 根据用户所持有的份额计算用户新增的奖励
        if (user_.stAmount > 0) {
            (bool success1, uint256 accST) = user_.stAmount.tryMul(
                pool_.accMetaNodePerST
            );
            require(success1, "user stAmount mul accMetaNodePerST overflow");
            (success1, accST) = accST.tryDiv(1 ether);
            require(success1, "accST div 1 ether overflow");

            (bool success2, uint256 pendingMetaNode_) = accST.trySub(
                user_.finishedMetaNode
            );
            require(success2, "accST sub finishedMetaNode overflow");

            if (pendingMetaNode_ > 0) {
                (bool success3, uint256 _pendingMetaNode) = user_
                    .pendingMetaNode
                    .tryAdd(pendingMetaNode_);
                require(success3, "user pendingMetaNode overflow");
                user_.pendingMetaNode = _pendingMetaNode;
            }
        }

        if (_amount > 0) {
            (bool success4, uint256 stAmount) = user_.stAmount.tryAdd(_amount);
            require(success4, "user stAmount overflow");
            user_.stAmount = stAmount;
        }

        (bool success5, uint256 stTokenAmount) = pool_.stTokenAmount.tryAdd(
            _amount
        );
        require(success5, "pool stTokenAmount overflow");
        pool_.stTokenAmount = stTokenAmount;

        // 隔离新增份额可获得的奖励（因为新增的份额不应该获得历史奖励）
        (bool success6, uint256 finishedMetaNode) = user_.stAmount.tryMul(
            pool_.accMetaNodePerST
        );
        require(success6, "user stAmount mul accMetaNodePerST overflow");

        (success6, finishedMetaNode) = finishedMetaNode.tryDiv(1 ether);
        require(success6, "finishedMetaNode div 1 ether overflow");

        user_.finishedMetaNode = finishedMetaNode;
    }

    function unstack(
        uint256 _pid,
        uint256 _amount
    ) public whenNotPaused checkPid(_pid) {
        Pool storage pool_ = pool[_pid];
        User storage user_ = user[_pid][msg.sender];

        require(user_.stAmount >= _amount, "not enough staked amount");

        // 更新质押池的奖励变量
        updatePool(_pid);

        // 计算用户可领取的奖励
        uint256 pendingMetaNode = ((user_.stAmount * pool_.accMetaNodePerST) /
            1 ether) - user_.finishedMetaNode;
        if (pendingMetaNode > 0) {
            user_.pendingMetaNode = user_.pendingMetaNode + pendingMetaNode;
        }

        if (_amount > 0) {
            user_.stAmount = user_.stAmount - _amount;
            user_.requests.push(
                UnstakeRequest({
                    amount: _amount,
                    unlockBlocks: block.number + pool_.unstakeLockedBlocks
                })
            );
            pool_.stTokenAmount = pool_.stTokenAmount - _amount;
        }

        // 在减去用户质押的份额后重新计算用户的累计奖励
        user_.finishedMetaNode =
            (user_.stAmount * pool_.accMetaNodePerST) /
            1 ether;
    }

    function withdraw(uint256 _pid) public whenNotPaused checkPid(_pid) {
        Pool storage pool_ = pool[_pid];
        User storage user_ = user[_pid][msg.sender];

        uint256 pendingWithdrawAmount_;
        uint256 popNum_;
        for (uint i = 0; i < user_.requests.length; i++) {
            // 校验提取请求是否已经全部解锁
            if (user_.requests[i].unlockBlocks > block.number) {
                break;
            }
            pendingWithdrawAmount_ += user_.requests[i].amount;
            popNum_++;
        }

        // 清理已经提取的请求
        for (uint i = 0; i < user_.requests.length - popNum_; i++) {
            user_.requests[i] = user_.requests[i + popNum_];
        }
        for (uint i = 0; i < popNum_; i++) {
            user_.requests.pop();
        }

        // 将质押代币转账给用户
        if (pendingWithdrawAmount_ > 0) {
            // eth 原生代币
            if (pool_.stTokenAddress == address(0)) {
                _safeETHTransfer(msg.sender, pendingWithdrawAmount_);
            } else {
                IERC20(pool_.stTokenAddress).transfer(
                    msg.sender,
                    pendingWithdrawAmount_
                );
            }
        }
    }

    function _safeETHTransfer(address _to, uint256 _amount) internal {
        (bool success, bytes memory data) = address(_to).call{value: _amount}(
            ""
        );
        require(success, "transfer failed");
        if (data.length > 0) {
            require(abi.decode(data, (bool)), "transfer failed");
        }
    }

    function claim(uint256 _pid) public whenNotPaused checkPid(_pid) {
        Pool storage pool_ = pool[_pid];
        User storage user_ = user[_pid][msg.sender];

        updatePool(_pid);

        uint256 pendingMetaNode_ = (user_.stAmount * pool_.accMetaNodePerST) /
            (1 ether) -
            user_.finishedMetaNode +
            user_.pendingMetaNode;

        if (pendingMetaNode_ > 0) {
            user_.pendingMetaNode = 0;
            _safeMetaNodeTransfer(msg.sender, pendingMetaNode_);
        }

        user_.finishedMetaNode =
            (user_.stAmount * pool_.accMetaNodePerST) /
            (1 ether);
    }

    function _safeMetaNodeTransfer(address _to, uint256 _amount) internal {
        uint256 metaNodeBalance_ = metaNodeToken.balanceOf(address(this));
        if (_amount > metaNodeBalance_) {
            metaNodeToken.transfer(_to, metaNodeBalance_);
        } else {
            metaNodeToken.transfer(_to, _amount);
        }
    }
}
