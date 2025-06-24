// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@uniswap/v2-core/contracts/interfaces/IUniswapV2Factory.sol";
import "@uniswap/v2-periphery/contracts/interfaces/IUniswapV2Router02.sol";

contract MyToken is ERC20, Ownable {
    // —— Fee Parameters ——
    /// @notice Total fee in basis points (1 BP = 0.01%). Max allowed 1000 (10%).
    uint16 public totalFeeBP;

    /// @notice Sub‑fee for liquidity (BP)
    uint16 public liqFeeBP;
    /// @notice Sub‑fee for marketing (BP)
    uint16 public mktFeeBP;
    /// @notice Sub‑fee for burn (BP)
    uint16 public burnFeeBP;

    // —— Fee Receivers ——
    address public liquidityReceiver;
    address public marketingReceiver;
    address public burnReceiver;

    // —— Exemptions ——
    mapping(address => bool) public isExcludedFromFee;

    // —— Uniswap V2 Router & Pair ——
    IUniswapV2Router02 public immutable uniswapV2Router;
    address public immutable uniswapV2Pair;

    // —— Transaction Limits ——
    uint256 public maxTxAmount;
    uint256 public maxDailyTxCount;
    struct TxData {
        uint256 count;
        uint256 windowStart;
    }
    mapping(address => TxData) public txData;

    // —— Events ——
    event FeesDistributed(uint256 liqAmt, uint256 mktAmt, uint256 burnAmt);
    event FeeParamsUpdated(
        uint16 totalBP,
        uint16 liqBP,
        uint16 mktBP,
        uint16 burnBP
    );
    event FeeReceiversUpdated(
        address indexed liqReceiver,
        address indexed mktReceiver,
        address indexed burnReceiver
    );
    event ExcludedFromFeeUpdated(address indexed account, bool excluded);
    event TxLimitUpdated(uint256 maxTxAmount, uint256 maxDailyTxCount);
    event TxCountReset(address indexed account);

    constructor(
        address router_,
        address _liq,
        address _mkt,
        address _burn
    ) ERC20("MyToken", "MTK") Ownable(msg.sender) {
        // set uniswap router
        uniswapV2Router = IUniswapV2Router02(router_);
        // create a uniswap pair for this token and WETH
        uniswapV2Pair = IUniswapV2Factory(uniswapV2Router.factory()).createPair(
                address(this),
                uniswapV2Router.WETH()
            );

        // set initial receivers
        liquidityReceiver = _liq;
        marketingReceiver = _mkt;
        burnReceiver = _burn;

        // exempt deployer and this contract from fees
        isExcludedFromFee[msg.sender] = true;
        isExcludedFromFee[address(this)] = true;

        // initial fee parameters
        setFeeParams(500, 167, 167, 166);

        // set default tx amount
        uint256 initialSupply = 1_000_000_000_000 * 10 ** decimals();
        maxTxAmount = (initialSupply * 5) / 1000;
        maxDailyTxCount = 10;
        emit TxLimitUpdated(maxTxAmount, maxDailyTxCount);

        // mint 1 000 000 000 000 * 10^decimals() to deployer
        _mint(msg.sender, initialSupply);
    }

    function _update(
        address from,
        address to,
        uint256 value
    ) internal override {
        // transaction limit (skip for exempted)
        if (!isExcludedFromFee[from] && !isExcludedFromFee[to]) {
            _enforceTxLimit(from, value);
        }

        // if fee disabled or either side is exempt → normal transfer
        if (
            totalFeeBP == 0 ||
            isExcludedFromFee[from] == true ||
            isExcludedFromFee[to] == true
        ) {
            super._update(from, to, value);
            return;
        }

        // calculate fee
        uint256 feeAmount = (value * totalFeeBP) / 10_000;
        uint256 sendAmount = value - feeAmount;

        // allocate sub‑fees
        uint256 liqAmount = (feeAmount * liqFeeBP) / totalFeeBP;
        uint256 mktAmount = (feeAmount * mktFeeBP) / totalFeeBP;
        uint256 burnAmount = feeAmount - liqAmount - mktAmount;

        // distribute fees
        super._update(from, liquidityReceiver, liqAmount);
        super._update(from, marketingReceiver, mktAmount);
        super._update(from, burnReceiver, burnAmount);

        // final transfer to recipient
        super._update(from, to, sendAmount);

        emit FeesDistributed(liqAmount, mktAmount, burnAmount);
    }

    function _enforceTxLimit(address sender, uint256 amount) internal {
        require(amount <= maxTxAmount, "amount exceeds maxTxAmount");

        TxData storage data = txData[sender];
        uint256 nowTS = block.timestamp;
        // if 24h window expired, reset
        if (nowTS >= data.windowStart + 1 days) {
            data.windowStart = nowTS;
            data.count = 0;
            emit TxCountReset(sender);
        }
        // increment & check count
        data.count += 1;
        require(data.count <= maxDailyTxCount, "exceeds daily tx count");
    }

    function setFeeParams(
        uint16 _totalFeeBP,
        uint16 _liqBP,
        uint16 _mktBP,
        uint16 _burnBP
    ) public onlyOwner {
        require(_totalFeeBP <= 1000, "totalFeeBP > 10%");
        require(_totalFeeBP == _liqBP + _mktBP + _burnBP, "sub-fees mismatch");

        totalFeeBP = _totalFeeBP;
        liqFeeBP = _liqBP;
        mktFeeBP = _mktBP;
        burnFeeBP = _burnBP;
        emit FeeParamsUpdated(_totalFeeBP, _liqBP, _mktBP, _burnBP);
    }

    /// @notice Update fee receivers
    function setFeeReceivers(
        address _liqReceiver,
        address _mktReceiver,
        address _burnReceiver
    ) external onlyOwner {
        liquidityReceiver = _liqReceiver;
        marketingReceiver = _mktReceiver;
        burnReceiver = _burnReceiver;
        emit FeeReceiversUpdated(_liqReceiver, _mktReceiver, _burnReceiver);
    }

    /// @notice Exempt or include an address from fee
    function setExcludedFromFee(
        address account,
        bool excluded
    ) external onlyOwner {
        isExcludedFromFee[account] = excluded;
        emit ExcludedFromFeeUpdated(account, excluded);
    }

    /// @notice Update max tx amount and daily tx count
    function setTxLimits(
        uint256 _maxTxAmount,
        uint256 _maxDailyTxCount
    ) external onlyOwner {
        maxTxAmount = _maxTxAmount;
        maxDailyTxCount = _maxDailyTxCount;
        emit TxLimitUpdated(_maxTxAmount, _maxDailyTxCount);
    }

    /// @notice Add liquidity to the DEX pool by supplying token+ETH.
    function addLiquidity(
        uint256 tokenAmount,
        uint256 amountTokenMin,
        uint256 amountETHMin
    ) external payable {
        require(msg.value > 0, "Must send ETH");
        // transfer tokens from caller
        _update(msg.sender, address(this), tokenAmount);
        // approve router to spend tokens
        _approve(address(this), address(uniswapV2Router), tokenAmount);

        // add liquidity, send LP tokens to caller
        uniswapV2Router.addLiquidityETH{value: msg.value}(
            address(this),
            tokenAmount,
            amountTokenMin,
            amountETHMin,
            msg.sender,
            block.timestamp
        );
    }

    /// @notice Remove liquidity from the DEX pool, burn or return LP tokens
    /// @dev Caller must approve this contract to spend `liquidity` LP tokens beforehand.
    function removeLiquidity(
        uint256 liquidity,
        uint256 amountTokenMin,
        uint256 amountETHMin
    ) external {
        IERC20(uniswapV2Pair).transferFrom(
            msg.sender,
            address(this),
            liquidity
        );
        IERC20(uniswapV2Pair).approve(address(uniswapV2Router), liquidity);

        uniswapV2Router.removeLiquidityETH(
            address(this),
            liquidity,
            amountTokenMin,
            amountETHMin,
            msg.sender,
            block.timestamp
        );
    }
}
