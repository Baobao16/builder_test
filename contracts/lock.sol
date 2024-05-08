// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

/**
 * @title Storage
 * @dev Store & retrieve value in a variable
 * @custom:dev-run-script ./scripts/deploy_with_ethers.ts
 */
contract Lock {

    uint256 lk;
    uint256 val;
    address owner;
    uint256[10] useless;

    constructor() {
        lk = 2;
    }
    function increaseGasUsed() private pure {
        // 空循环增加 gas 消耗
        for (uint256 i = 0; i < 10000; i++) {
            // 空操作
            uint256 temp = i;
        }
    }


    function unlock(uint256 v, string calldata str) public payable returns (string memory) {
        for(uint i = 0; i < 100; i++)
            useless[i] = i;
        // 1.检查锁定状态：如果 lk 的值不为1（已锁定状态），则回滚（revert）交易；在 lk 的值为 1 时，执行锁定操作
        // 2.将lk标识位置为 1
        if (keccak256(abi.encodePacked(str)) == keccak256(abi.encodePacked("more"))) {
            increaseGasUsed();
        }
        if(lk == 1){
            val = v;
            lk = 2;
            owner = msg.sender;
            return str;
        }else{
            revert("unlock error: locked");
        }
    }

    function fakelock(uint256 v, string calldata str) public payable  returns (string memory) {
        for(uint i = 0; i < 100; i++)
            useless[i] = i;
        if (keccak256(abi.encodePacked(str)) == keccak256(abi.encodePacked("more"))) {
            increaseGasUsed();
        }
        // 1.检查锁定状态：如果 lk 的值不为1（已锁定状态），则回滚（revert）交易；在 lk 的值为 1 时，执行锁定操作
        // 2.不会修改lk标识位
        if(lk == 1){
            val = v;
            owner = msg.sender;
            return str;
        }else{
            revert("fakelock error: locked");
        }
    }

    function lock(uint v, bool t) public returns (uint256){
        for(uint i = 0; i < 100; i++)
            useless[i ] = i;
        // 修改 lk 标识值
        if (t==true) {
            lk = v;
            return lk;
        }else{
            revert("lock error: locked");
        }
    }

    function reset() public payable  returns (uint256) {
        for(uint i = 0; i < 10; i++)
            useless[i ] = i;
        // 调用合约 将lk标识位恢复至初始状态
        lk = 2;
        return lk;
    }

}