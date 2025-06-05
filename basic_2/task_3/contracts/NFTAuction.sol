// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract NFTAuction is Ownable, Initializable {
    struct AuctionDetails {
        address seller;
        address nftContract;
        uint256 tokenId;
        address erc20Token; // 当竞拍者使用ERC20代币竞拍时使用
        uint256 startTime;
        uint256 endTime;
        uint256 highestBid;
        address highestBidder;
        bool ended;
    }

    AuctionDetails public auction;
    mapping(address => address) public priceFeeds;

    // 设置价格源
    function setPriceFeed(
        address _erc20Token,
        address _priceFeed
    ) external onlyOwner {
        priceFeeds[_erc20Token] = _priceFeed;
    }

    // 将代币金额转换成USD（结果会被放大10^18倍）
    function convertToUSD(
        address _erc20Token, // address(0) 代表ETH
        uint256 _amount
    ) public view returns (uint256) {
        require(priceFeeds[_erc20Token] != address(0), "Price feed not set");
        // 获取 ERC20 代币的精度
        uint8 tokenDecimals = _erc20Token == address(0)
            ? 18
            : IERC20Metadata(_erc20Token).decimals();
        // 获取喂价合约
        AggregatorV3Interface dataFeed = AggregatorV3Interface(
            priceFeeds[_erc20Token]
        );
        // prettier-ignore
        (
            /* uint80 roundId */,
            int256 answer,
            /*uint256 startedAt*/,
            /*uint256 updatedAt*/,
            /*uint80 answeredInRound*/
        ) = dataFeed.latestRoundData();
        // 将结果统一放大10^18倍确保结果为整数
        return (uint256(_amount) *
            uint256(answer) *
            (10 ** (18 - tokenDecimals)));
    }

    // 创建拍卖
    function initialize(
        address _seller,
        address _nftContract,
        uint256 _tokenId,
        uint256 _duration
    ) external initializer {
        require(_duration > 0, "Duration must be greater than 0");

        // 转移NFT到拍卖合约（转移前需要NFT所有者向拍卖合约授权）
        IERC721(_nftContract).transferFrom(msg.sender, address(this), _tokenId);

        // 创建新的拍卖
        auction = AuctionDetails({
            seller: _seller,
            nftContract: _nftContract,
            tokenId: _tokenId,
            erc20Token: address(0), // 默认为ETH
            startTime: block.timestamp,
            endTime: block.timestamp + _duration,
            highestBid: 0,
            highestBidder: address(0),
            ended: false
        });
    }

    // 出价
    function placeBid(
        address _erc20Token, // 当竞拍者使用ERC20代币竞拍时使用
        uint256 _amount
    ) external payable {
        require(block.timestamp >= auction.startTime, "Auction not started");
        require(block.timestamp < auction.endTime, "Auction ended");
        require(!auction.ended, "Auction already ended");

        // 出价金额
        uint256 priceInUSD;
        if (_erc20Token == address(0)) {
            require(msg.value == _amount, "ETH amount does not match");
            priceInUSD = convertToUSD(address(0), msg.value);
        } else {
            IERC20(_erc20Token).transferFrom(
                msg.sender,
                address(this),
                _amount
            );
            priceInUSD = convertToUSD(_erc20Token, _amount);
        }
        // 历史出价金额
        uint256 oldPriceInUSD;
        if (auction.highestBidder != address(0)) {
            if (auction.erc20Token == address(0)) {
                oldPriceInUSD = convertToUSD(address(0), auction.highestBid);
            } else {
                oldPriceInUSD = convertToUSD(
                    auction.erc20Token,
                    auction.highestBid
                );
            }
        }
        require(priceInUSD > oldPriceInUSD, "Bid too low");

        // 退还前一个最高出价者的资金
        if (auction.highestBidder != address(0)) {
            if (auction.erc20Token == address(0)) {
                // 如果是ETH，直接退还
                payable(auction.highestBidder).transfer(auction.highestBid);
            } else {
                // 退还前一个最高出价者的ERC20
                IERC20 token = IERC20(auction.erc20Token);
                require(
                    token.transfer(auction.highestBidder, auction.highestBid),
                    "ERC20 refund failed"
                );
            }
        }

        auction.erc20Token = _erc20Token;
        auction.highestBidder = msg.sender;
        if (_erc20Token == address(0)) {
            auction.highestBid = msg.value;
        } else {
            auction.highestBid = _amount;
        }
    }

    // 结束拍卖
    function endAuction() external onlyOwner {
        require(block.timestamp >= auction.endTime, "Auction not ended");
        require(!auction.ended, "Auction already ended");

        auction.ended = true;

        if (auction.highestBidder != address(0)) {
            // 转移NFT给最高出价者
            IERC721(auction.nftContract).transferFrom(
                address(this),
                auction.highestBidder,
                auction.tokenId
            );

            // 转移资金给卖家
            if (auction.erc20Token == address(0)) {
                payable(auction.seller).transfer(auction.highestBid);
            } else {
                IERC20(auction.erc20Token).transfer(
                    auction.seller,
                    auction.highestBid
                );
            }
        } else {
            // 无人出价，退回NFT给卖家
            IERC721(auction.nftContract).transferFrom(
                address(this),
                auction.seller,
                auction.tokenId
            );
        }
    }
}
