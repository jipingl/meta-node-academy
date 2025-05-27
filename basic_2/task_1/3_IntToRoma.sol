// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract IntToRoma {
    // 存储罗马数字与对应值的关系
    struct Roma {
        uint32 value;
        string symbol;
    }

    Roma[] private romaArr;

    constructor() {
        romaArr.push(Roma({value: 1000, symbol: "M"}));
        romaArr.push(Roma({value: 900, symbol: "CM"}));
        romaArr.push(Roma({value: 500, symbol: "D"}));
        romaArr.push(Roma({value: 400, symbol: "CD"}));
        romaArr.push(Roma({value: 100, symbol: "C"}));
        romaArr.push(Roma({value: 90, symbol: "XC"}));
        romaArr.push(Roma({value: 50, symbol: "L"}));
        romaArr.push(Roma({value: 40, symbol: "XL"}));
        romaArr.push(Roma({value: 10, symbol: "X"}));
        romaArr.push(Roma({value: 9, symbol: "IX"}));
        romaArr.push(Roma({value: 5, symbol: "V"}));
        romaArr.push(Roma({value: 4, symbol: "IV"}));
        romaArr.push(Roma({value: 1, symbol: "I"}));
    }

    // 整数转罗马数字
    function intToRoma(uint32 _num) public view returns (string memory) {
        string memory res;
        uint32 remaining = _num;
        for (uint8 i = 0; i < romaArr.length; i++) {
            while (remaining >= romaArr[i].value) {
                res = string(abi.encodePacked(res, romaArr[i].symbol));
                remaining -= romaArr[i].value;
            }
        }
        return res;
    }
}
