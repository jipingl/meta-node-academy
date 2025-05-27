// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RomaToInt {
    // 存储罗马数字到数值的映射
    mapping(bytes1 => uint32) private romaValueMap;

    constructor() {
        romaValueMap["I"] = 1;
        romaValueMap["V"] = 5;
        romaValueMap["X"] = 10;
        romaValueMap["L"] = 50;
        romaValueMap["C"] = 100;
        romaValueMap["D"] = 500;
        romaValueMap["M"] = 1000;
    }

    function romaToInt(string memory _roma) public view returns (uint32) {
        bytes memory romaBytes = bytes(_roma);
        uint256 len = romaBytes.length;
        uint32 preValue;
        uint32 res;
        // 从右向左遍历罗马数字
        for (uint256 i = 0; i < len; i++) {
            // 判断当前的值是否小于之前的值（处理IV=4这种情况）
            uint32 currentValue = romaValueMap[romaBytes[len - 1 - i]];
            if (currentValue < preValue) {
                res -= currentValue;
            } else {
                res += currentValue;
            }
            // 记录当前的值
            preValue = currentValue;
        }
        return res;
    }
}
