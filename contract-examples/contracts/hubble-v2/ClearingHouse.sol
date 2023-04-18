// SPDX-License-Identifier: BUSL-1.1

pragma solidity 0.8.9;

import { SafeCast } from "@openzeppelin/contracts/utils/math/SafeCast.sol";

contract ClearingHouse {
    using SafeCast for uint256;
    using SafeCast for int256;

    function getUnderlyingPrice() public pure returns(uint[] memory prices) {
        prices = new uint[](1);
        prices[0] = 10000000;
    }
}
