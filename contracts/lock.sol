// SPDX-License-Identifier: GPL-3.0

pragma solidity >0.5.10;

/**
 * @title Storage
 * @dev Store & retrieve value in a variable
 * @custom:dev-run-script ./scripts/deploy_with_ethers.ts
 */
contract Lock {

    uint256 lk;
    uint256 val;
    address owner;

    constructor() {
        lk = 1;
    }


    function unlock(uint256 v, string calldata str) public payable returns (string memory) {
        // 1.检查锁定状态：如果 lk 的值不为0（已锁定状态），则回滚（revert）交易；在 lk 的值为 0 时，执行锁定操作
        // 2.将lk标识位置为 1
        if(lk == 0){
            val = v;
            lk = 1;
            owner = msg.sender;
            return str;
        }else{
            revert("unlock error: locked");
        }
    }

    function fakelock(uint256 v, string calldata str) public payable  returns (string memory) {
        // 1.检查锁定状态：如果 lk 的值不为0（已锁定状态），则回滚（revert）交易；在 lk 的值为 0 时，执行锁定操作
        // 2.不会修改lk标识位
        if(lk == 0){
            val = v;
            owner = msg.sender;
            return str;
        }else{
            revert("fakelock error: locked");
        }
    }

    function lock(uint v) public returns (uint256){
        // 修改 lk 标识值
        lk = v;
        return lk;
    }

    function reset() public payable  returns (uint256) {
        // 调用合约 将lk标识位恢复至初始状态
        lk = 1;
        return lk;
    }


}