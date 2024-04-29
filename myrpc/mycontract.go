package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/contracts/MyContract" // 你的合约包
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 连接到以太坊节点
	client, err := ethclient.Dial("http://10.2.66.75:28545")
	if err != nil {
		log.Fatal(err)
	}

	// 加载合约 ABI
	contractAddress := common.HexToAddress("0xd9145CCE52D386f254917e481eB44e9943F39138") // 合约地址
	contractInstance, err := MyContract.NewMyContract(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// 获取账户地址和私钥
	privateKey, err := crypto.HexToECDSA("61bfe9aea17bec5de54a86ad6cb0418f678a2fc8b746cc3901687eaebe1da809")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 创建交易选项
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = nil                       // nil 会自动填充 nonce
	auth.Value = big.NewInt(0)             // 转账金额
	auth.GasLimit = uint64(300000)         // 交易所能消耗的 gas 上限
	auth.GasPrice = big.NewInt(1000000000) // gas 价格

	// 调用合约的 setValue 函数并传入参数
	tx, err := contractInstance.SetValue(auth, big.NewInt(123))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Transaction hash: %s\n", tx.Hash().Hex())
}
