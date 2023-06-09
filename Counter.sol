// SPDX-License-Identifier: Unlicense

pragma solidity ^0.8.18;

contract Counter {
    uint256 public count;

    constructor() {}

    function increment() public {
        count += 1;
    }
}
