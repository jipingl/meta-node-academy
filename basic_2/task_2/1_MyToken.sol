// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract MyToken {
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(
        address indexed owner,
        address indexed spender,
        uint256 value
    );

    // 合约所有者
    address public owner;
    // 代币发行量
    uint256 public totalSupply;

    // 存储账户余额
    mapping(address account => uint256) private _balances;
    // 存储定量限额授权
    mapping(address owner => mapping(address spender => uint256))
        private _allowances;

    constructor() {
        owner = msg.sender;
    }

    // 铸造代币
    function mint(address account, uint256 value) public {
        require(msg.sender == owner, "only owner can mint");
        _balances[account] += value;
        totalSupply += value;
    }

    // 查询账户余额
    function balanceOf(address account) public view returns (uint256) {
        return _balances[account];
    }

    // 转账
    function transfer(address to, uint256 value) public {
        // 检查余额
        require(_balances[msg.sender] >= value, "insufficient balance");
        // 减去余额
        _balances[msg.sender] -= value;
        // 目标账户增加余额
        _balances[to] += value;
        // 触发转账事件
        emit Transfer(msg.sender, to, value);
    }

    // 授权
    function approve(address spender, uint256 value) public {
        _allowances[msg.sender][spender] = value;
        emit Approval(msg.sender, spender, value);
    }

    // 代扣转账
    function transferFrom(address from, address to, uint256 value) public {
        address spender = msg.sender;
        // 检查授权余额
        require(_allowances[from][spender] >= value, "insufficient allowance");
        // 扣减授权余额
        _allowances[from][spender] -= value;
        // 转账
        require(_balances[from] >= value, "insufficient balance");
        _balances[from] -= value;
        _balances[to] += value;
    }
}
