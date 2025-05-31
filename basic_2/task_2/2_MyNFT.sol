// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

// 测试网合约地址
// 0xfcd24d14df96eb1c91987b88f5e824e45ee91898
// 3#代币JSON文件地址
// ipfs://bafkreibxha5tp65alskyo3icdnwjnmn3xuin7t6ibmhoi3q2uranea55hy
contract MyNFT is ERC721, Ownable {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIds;

    string private _name;
    string private _symbol;
    mapping(uint256 tokenId => string tokenURI) private _tokenURIs;

    constructor(
        string memory name_,
        string memory symbol_
    ) ERC721(name_, symbol_) Ownable(msg.sender) {
        _name = name_;
        _symbol = symbol_;
    }

    function mintNFT(string memory _tokenURI) public returns (uint256) {
        _tokenIds.increment();
        uint256 newTokenId = _tokenIds.current();

        _mint(msg.sender, newTokenId);
        _tokenURIs[newTokenId] = _tokenURI;

        return newTokenId;
    }

    function tokenURI(
        uint256 tokenId
    ) public view override(ERC721) returns (string memory) {
        string memory _tokenURI = _tokenURIs[tokenId];
        require(
            bytes(_tokenURI).length > 0 &&
                keccak256(bytes(_tokenURI)) != keccak256(bytes("")),
            "token does not exist"
        );
        return _tokenURI;
    }
}
