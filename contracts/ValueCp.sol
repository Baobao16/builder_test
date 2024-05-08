// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

/**
 * @title Storage
 * @dev Store & retrieve value in a variable
 * @custom:dev-run-script ./scripts/deploy_with_ethers.ts
 */
contract ValueCapture {

    address immutable bribe = 0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE;
 
    function bet(bool valid) public payable {
        if(bribe.balance == 0) {
            bribe.call{value: msg.value}("");
        }else{
            if(valid){
                uint val = msg.value * 90 /100;
                bribe.call{value: val}("");
            }else{
                uint val = msg.value * 70 /100;
                bribe.call{value: val}("");
            }
        }
    }

}