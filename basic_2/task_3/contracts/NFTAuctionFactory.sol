// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./NFTAuction.sol";

contract NFTAuctionFactory {
    // 所有拍卖合约地址列表
    address[] public allAuctions;

    // NFT合约 => TokenID => 拍卖合约地址
    mapping(address => mapping(uint256 => address)) public auctions;

    // 创建新拍卖
    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 duration
    ) external returns (address auction) {
        require(
            auctions[nftContract][tokenId] == address(0),
            "Auction already exists"
        );

        // 使用create2生成确定性地址
        bytes memory bytecode = type(NFTAuction).creationCode;
        bytes32 salt = keccak256(abi.encodePacked(nftContract, tokenId));
        assembly {
            auction := create2(0, add(bytecode, 32), mload(bytecode), salt)
        }
        NFTAuction(auction).initialize(
            msg.sender,
            nftContract,
            tokenId,
            duration
        );

        auctions[nftContract][tokenId] = auction;
        allAuctions.push(auction);
    }
}
