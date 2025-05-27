// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BinarySearch {
    function binarySearch(
        int[] memory array,
        int target
    ) public pure returns (uint index) {
        uint left = 0;
        uint right = array.length - 1;
        while (left <= right) {
            // 计算中间位置
            uint mid = (left + right) / 2;
            // 检查目标值所处范围
            if (array[mid] == target) {
                return mid;
            } else if (array[mid] > target) {
                right = mid - 1;
            } else {
                left = mid + 1;
            }
        }
        return type(uint).max;
    }
}
