// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BeggingContract {
    // 合约所有者
    address private _owner;
    // 记录捐赠者的捐赠金额
    mapping(address => uint256) private _donations;
    // 记录捐赠金额排名前3的捐赠者
    address[3] private _topDonors;
    // 捐赠时间段
    uint256 private _startTime;
    uint256 private _endTime;

    event Donation(address indexed donator, uint256 amount);

    constructor(uint256 startTime, uint256 endTime) {
        _owner = msg.sender;
        _startTime = startTime;
        _endTime = endTime;
    }

    modifier onlyOwner() {
        require(msg.sender == _owner, "not the owner");
        _;
    }

    function blockTimestamp() public view returns (uint256) {
        return block.timestamp;
    }

    // 捐赠函数
    function donate() public payable {
        require(
            block.timestamp >= _startTime && block.timestamp <= _endTime,
            "not during the donation period"
        );
        require(msg.value > 0, "invalid donation amount");
        _donations[msg.sender] += msg.value;
        // 更新捐赠者前3名
        _updateTopDonors(msg.sender);
        emit Donation(msg.sender, msg.value);
    }

    // 合约所有者提取所有资金
    function withdraw() public onlyOwner {
        require(block.timestamp > _endTime, "the donation has not ended yet");
        payable(_owner).transfer(address(this).balance);
    }

    // 查询指定用户的捐赠金额
    function getDonation(address donator) public view returns (uint256) {
        return _donations[donator];
    }

    // 获取前3名捐赠者及其金额
    function getTopDonors()
        public
        view
        returns (address[3] memory, uint256[3] memory)
    {
        uint256[3] memory amounts;
        for (uint i = 0; i < 3; i++) {
            amounts[i] = _donations[_topDonors[i]];
        }
        return (_topDonors, amounts);
    }

    // 更新捐赠者前3名
    function _updateTopDonors(address donator) private {
        // 判断用户是否已经进入前三名
        for (uint i = 0; i < 3; i++) {
            // 用户已经进入前三名
            if (donator == _topDonors[i]) {
                _sortTopDonors();
                return;
            }
        }
        // 确认用户是否能进入前三
        if (_donations[donator] > _donations[_topDonors[2]]) {
            _topDonors[2] = donator;
            _sortTopDonors();
        }
    }

    // 对排名前三的捐赠者进行排序
    function _sortTopDonors() private {
        for (uint i = 0; i < 2; i++) {
            for (uint j = 0; j < 2 - i; j++) {
                if (_donations[_topDonors[j]] < _donations[_topDonors[j + 1]]) {
                    address tmp = _topDonors[j];
                    _topDonors[j] = _topDonors[j + 1];
                    _topDonors[j + 1] = tmp;
                }
            }
        }
    }
}
