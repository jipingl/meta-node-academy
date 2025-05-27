// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

contract Voting {
    // 存储合约创建者
    address public owner;
    // 存储所有的候选人
    string[] public candidates;
    // 存储候选人的得票数
    mapping(string => uint32) public voteCount;
    // 存储投票的人
    address[] public voters;
    // 存储投票的账号
    mapping(address => bool) public voterMap;

    constructor() {
        owner = msg.sender;
    }

    // 投票函数
    function vote(string calldata candidate) public {
        require(!voterMap[msg.sender], "You have already cast your vote");
        voters.push(msg.sender);
        voterMap[msg.sender] = true;
        voteCount[candidate] += 1;
        candidates.push(candidate);
    }

    // 获取候选人的得票数
    function getVotes(string calldata candidate) public view returns (uint32) {
        return voteCount[candidate];
    }

    // 重置
    function resetVotes() public {
        // 只有合约创建者可以重置
        require(owner == msg.sender, "You are not allowed");
        // 清空投票数
        for (uint i = 0; i < candidates.length; ++i) {
            voteCount[candidates[i]] = 0;
        }
        // 清空投票人
        for (uint i = 0; i < voters.length; ++i) {
            voterMap[voters[i]] = false;
        }
    }
}
