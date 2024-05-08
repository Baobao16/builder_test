// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.8.2 <0.9.0;

/**
 * @title Storage
 * @dev Store & retrieve value in a variable
 * @custom:dev-run-script ./scripts/deploy_with_ethers.ts
 */
contract SpecialOp {

    address immutable bribe = 0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE;
    function testTimestamp(uint timestamp) public view{
        if(timestamp <= block.timestamp){

        }else{
            revert();
        }
    }

        function testTimestampEq(uint timestamp) public view{
        if(timestamp == block.timestamp){

        }else{
            revert();
        }
    }

    function testCoinbase(address coinbase) public view {
        if(coinbase == block.coinbase){

        }else{
            revert();
        }
    }

    function testBlockHash(bytes32 input) public view {
        bytes32 hash = blockhash(block.number-1);
        if (compareBytes32(input, hash)) {

        }else{
            revert();
        }
    }
    function compareBytes32(bytes32 a, bytes32 b) public pure returns(bool) {
        return keccak256(abi.encodePacked(a)) == keccak256(abi.encodePacked(b));
    }

    function testBribe() payable public {
        bribe.call{value: msg.value}("");
    }


    function getCoinbase()  public view returns(address)   {
        return block.coinbase;
    }

}