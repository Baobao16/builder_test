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
    uint256 x;

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

    function useGas() internal {
        while(gasleft() > 20000) {
            x = x + 1;
        }
    }


    function unlock(uint256 v, string calldata str) public payable returns (string memory) {

        // 1.检查锁定状态：如果 lk 的值不为1（已锁定状态），则回滚（revert）交易；在 lk 的值为 1 时，执行锁定操作
        // 2.将lk标识位置为 1
        if (keccak256(abi.encodePacked(str)) == keccak256(abi.encodePacked("more"))) {
            increaseGasUsed();
        }
        if(lk == 1){
            val = v;
            lk = 2;
            owner = msg.sender;

        }else{
            revert("unlock error: locked");
        }
        useGas();
        return str;
    }

    function unlock_de(uint256 v, string calldata str) public payable returns (string memory) {

        // 1.检查锁定状态：如果 lk 的值不为1（已锁定状态），则回滚（revert）交易；在 lk 的值为 1 时，执行锁定操作
        // 2.将lk标识位置为 1
        // GaseUsed 不受GasLimit影响
        if (keccak256(abi.encodePacked(str)) == keccak256(abi.encodePacked("more"))) {
            increaseGasUsed();
        }
        if(lk == 1){
            val = v;
            lk = 2;
            owner = msg.sender;

        }else{
            revert("unlock error: locked");
        }
        return str;
    }

    function fakelock(uint256 v, string calldata str) public payable  returns (string memory) {

        if (keccak256(abi.encodePacked(str)) == keccak256(abi.encodePacked("more"))) {
            increaseGasUsed();
        }
        // 1.检查锁定状态：如果 lk 的值不为1（已锁定状态），则回滚（revert）交易；在 lk 的值为 1 时，执行锁定操作
        // 2.不会修改lk标识位
        if(lk == 1){
            val = v;
            owner = msg.sender;

        }else{
            revert("fakelock error: locked");
        }
        useGas();
        return str;
    }

    function lock(uint v, bool t) public returns (uint256){

        // 修改 lk 标识值
        if (t==true) {
            lk = v;

        }else{
            revert("lock error: locked");
        }
        useGas();
        return lk;
    }

    function reset() public payable  returns (uint256) {

        // 调用合约 将lk标识位恢复至初始状态
        lk = 2;
        useGas();
        return lk;
    }

}