// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MergeSortedArray {
    function mergeSortedArray(
        int[] memory _nums1,
        int[] memory _nums2
    ) public pure returns (int[] memory) {
        uint nums1Len = _nums1.length;
        uint nums2Len = _nums2.length;
        int[] memory mergedBytes = new int[](nums1Len + nums2Len);

        // 三个指针
        uint ptr1 = 0;
        uint ptr2 = 0;
        uint mPtr = 0;
        while (ptr1 < nums1Len && ptr2 < nums2Len) {
            if (_nums1[ptr1] <= _nums2[ptr2]) {
                mergedBytes[mPtr++] = _nums1[ptr1];
                ++ptr1;
            } else {
                mergedBytes[mPtr++] = _nums2[ptr2];
                ++ptr2;
            }
        }
        while (ptr1 < nums1Len) {
            // 把剩余的放到合并数组中
            mergedBytes[mPtr++] = _nums1[ptr1];
            ++ptr1;
        }
        while (ptr2 < nums2Len) {
            // 把剩余的放到合并数组中
            mergedBytes[mPtr++] = _nums2[ptr2];
            ++ptr2;
        }
        return mergedBytes;
    }
}
