// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ReverseString {
    // 反转字符串
    function reverseString(
        string memory _str
    ) public pure returns (string memory) {
        bytes memory strBytes = bytes(_str);
        bytes memory reversedBytes = new bytes(strBytes.length);
        for (uint32 i = 0; i < strBytes.length; i++) {
            reversedBytes[strBytes.length - 1 - i] = strBytes[i];
        }
        return string(reversedBytes);
    }
}
