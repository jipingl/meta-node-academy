// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MetaNodeToken is ERC20 {
    constructor() ERC20("MetaNode", "MNODE") {
        _mint(msg.sender, 10_000_000 * 10 ** decimals());
    }
}
