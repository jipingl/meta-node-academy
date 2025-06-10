// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyToken is ERC20, Ownable {
    uint256 public constant RATE = 100; // 1ETH => 100 RDT

    constructor() ERC20("MyToken", "MYT") Ownable(msg.sender) {}

    function mint() public payable {
        _mint(msg.sender, msg.value * RATE);
    }

    function withdrawETH() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No ETH to withdraw");
        payable(owner()).transfer(balance);
    }

    receive() external payable {
        mint();
    }
}
