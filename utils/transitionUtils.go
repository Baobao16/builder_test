package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/xkwang/conf"
	"github.com/xkwang/sendBundle"
)

type TxStatus struct {
	Usr   string
	Mined bool
	Rst   string
	Txs   types.Transactions
}

type Contract struct {
	ABI abi.ABI
}

func (c *Contract) CallFunction(method string, args ...interface{}) ([]byte, error) {
	input, err := c.ABI.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func CreateClient(url string) *ethclient.Client {
	// 创建以太坊客户端的通用函数
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %v", url, err)
	}
	return client
}

func UserTx(rootName string, contract common.Address, data []byte, gasLimit *big.Int, gasPrice *big.Int) sendBundle.BidCaseArg {
	ctx := context.Background()

	rootPk := rootName
	bobPk := rootName
	builderPk := *conf.BuilderPrivateKey

	client := CreateClient(conf.Url)
	client2 := CreateClient(conf.Url_1)
	client3 := CreateClient(conf.Url)

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	return sendBundle.BidCaseArg{
		Ctx:           ctx,
		Client:        client,
		ChainID:       chainID,
		RootPk:        rootPk,
		BobPk:         bobPk,
		Builder:       sendBundle.NewAccount(builderPk),
		Validators:    []common.Address{common.HexToAddress(*conf.Validator)},
		BidClient:     client2,
		BuilderClient: client3,
		TxCount:       3,
		Contract:      contract,
		Data:          data,
		GasPrice:      gasPrice,
		//GasPrice:      big.NewInt(conf.MinGasPrice),
		GasLimit:   gasLimit,
		SendAmount: big.NewInt(0),
	}
}

func serializeTxs(txs types.Transactions, txBytes []hexutil.Bytes) []hexutil.Bytes {
	// 定义一个函数来处理交易的序列化
	for _, tx := range txs {
		txByte, err := tx.MarshalBinary()
		log.Printf("txhash %v will be send \n", tx.Hash().Hex())
		if err != nil {
			log.Printf("Failed to marshal tx %v: %v", tx.Hash().Hex(), err)
			continue
		}
		txBytes = append(txBytes, txByte)
	}
	return txBytes
}

func AddBundle(txs types.Transactions, txsNew types.Transactions, revertTxHashes []common.Hash, MaxBN uint64) *types.SendBundleArgs {
	// 构造新的bundle，包含Mem-pool交易
	txBytes := make([]hexutil.Bytes, 0)

	// 序列化交易
	txBytes = serializeTxs(txs, txBytes)
	txBytes = serializeTxs(txsNew, txBytes)

	// 构建bundle参数
	bundleArgs := &types.SendBundleArgs{
		Txs:               txBytes,
		RevertingTxHashes: revertTxHashes,
		MaxBlockNumber:    MaxBN,
	}

	// 打印bundle参数的JSON表示
	if bidJson, err := json.MarshalIndent(bundleArgs, "", "  "); err == nil {
		log.Println(string(bidJson))
	} else {
		log.Printf("Failed to marshal bundleArgs: %v", err)
	}

	return bundleArgs
}

func IsEmptyField(result ResultB) bool {
	v := reflect.ValueOf(result)
	t := reflect.TypeOf(result)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name
		fieldValue := field.Interface()
		fmt.Printf("Field %s: %v\n", fieldName, fieldValue)
		switch field.Kind() {
		case reflect.String, reflect.Slice:
			if field.Len() == 0 {
				return true
			}
		}
	}
	return false
}

func SendLockMempool(usr string, contract common.Address, data []byte, gasLimit *big.Int, gasPrice *big.Int, revert bool, send bool) (types.Transactions, []common.Hash) {
	usrArg := UserTx(usr, contract, data, gasLimit, gasPrice)
	if revert {
		log.Printf("mem_pool transaction  will in bundle RevertList . ")
		usrArg.RevertListNormal = []int{0} // 当前交易被记入RevertList
	}
	//log.Printf("Set mem_pool transaction  ")
	tx, revertHash := sendBundle.GenerateBNBTxs(&usrArg, usrArg.SendAmount, usrArg.Data, 1)
	txBytes := make([]hexutil.Bytes, 0)
	serializeTxs(tx, txBytes)
	if send {
		err := usrArg.Client.SendTransaction(usrArg.Ctx, tx[0])
		log.Printf("Send Mem_pool transaction  %v [gasPrice: %v , gasLimit : %v] \n", tx[0].Hash(), usrArg.GasPrice, usrArg.GasLimit)
		if err != nil {
			fmt.Println("failed to send single Transaction", "err", err)
		}
	}

	return tx, revertHash

}
func ConcurSendBundles(t *testing.T, args []*sendBundle.BidCaseArg, bundleArgsList []*types.SendBundleArgs) uint64 {
	// 获取当前的区块号

	log.Println("Sending bundles and waiting for block number to increase")
	currentBlockNumber, _ := args[0].Client.BlockNumber(args[0].Ctx)

	var wg sync.WaitGroup

	sendBundle := func(i int) {
		defer wg.Done()
		_, err := args[i].BuilderClient.SendBundle(args[i].Ctx, *bundleArgsList[i])
		if err != nil {
			log.Printf("Sending bundle %d failed: %v", i, err)
			assert.Contains(t, err.Error(), conf.InvalidTx)
		}
	}

	for i := 0; i < len(args); i++ {
		wg.Add(1)
		go sendBundle(i)
	}

	wg.Wait()

	// 等待一段时间以确保区块号增加
	time.Sleep(5 * time.Second)

	// 获取新的区块号
	newBlockNumber, err := args[0].Client.BlockNumber(args[0].Ctx)
	if err != nil {
		t.Fatalf("Failed to get new block number: %v", err)
	}

	// 断言区块号确实增加了
	assert.True(t, newBlockNumber > currentBlockNumber, "Block number did not increase")

	return newBlockNumber
}

func ChangeArgType(arg interface{}) int {
	// println(arg)
	if val, ok := arg.(int); ok {
		fmt.Println("参数转换结果:", val)
		return val
	} else {
		fmt.Println("转换失败")
		return -1
	}
}

func GeneEncodedData(con Contract, method string, args ...interface{}) []byte {
	// 构造调用智能合约函数的数据
	encodeArgs := func(args ...interface{}) ([]byte, error) {
		if len(args) == 0 {
			return con.CallFunction(method)
		} else if len(args) == 1 {
			return con.CallFunction(method, args[0])
		} else if len(args) == 2 {
			val := ChangeArgType(args[0])
			log.Printf("Converted value: %v", val)
			var bigInt big.Int
			bigInt.SetInt64(int64(val))
			args2 := []interface{}{&bigInt, args[1]}
			return con.CallFunction(method, args2...)
		} else {
			return nil, fmt.Errorf("wrong number of arguments")
		}
	}

	contractData, err := encodeArgs(args...)
	if err != nil {
		log.Fatalf("Error encoding data: %v", err)
	}

	log.Printf("Calling smart contract function %s with arguments %v", method, args)
	return contractData
}

func ResetLockContract(t *testing.T, contract common.Address, data []byte) {
	t.Log("Root User reset Contract lock")
	usrArg := UserTx(conf.RootPk5, contract, data, conf.HighGas, big.NewInt(conf.MinGasPrice))
	usrArg.TxCount = 1

	txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&usrArg)

	bn, err := usrArg.Client.BlockNumber(usrArg.Ctx)
	if err != nil {
		t.Fatalf("Failed to get current block number: %v", err)
	}
	log.Printf("Current Block height is: %v", bn)

	cbn := SendBundlesMined(t, usrArg, bundleArgs)
	WaitMined(txs, cbn)

	blkNum := CheckBundleTx(t, *txs[0], true, conf.TxSucceed)
	log.Printf("Bundle included in block: %v", blkNum)
}

func SendBundlesMined(t *testing.T, usr sendBundle.BidCaseArg, bundleArgs *types.SendBundleArgs) uint64 {
	// 获取当前的区块号
	currentBlockNumber, err := usr.Client.BlockNumber(usr.Ctx)
	if err != nil {
		t.Fatalf("Failed to get current block number: %v", err)
	}

	log.Println("Sending bundle and waiting for block number to increase")

	// 发送捆绑交易
	_, err = usr.BuilderClient.SendBundle(usr.Ctx, *bundleArgs)
	if err != nil {
		t.Fatalf("Failed to send bundle: %v", err)
	}

	// 等待一段时间以确保区块号增加
	time.Sleep(6 * time.Second)

	// 获取新的区块号
	newBlockNumber, err := usr.Client.BlockNumber(usr.Ctx)
	if err != nil {
		t.Fatalf("Failed to get new block number: %v", err)
	}

	// 断言区块号确实增加了
	if newBlockNumber <= currentBlockNumber {
		t.Fatalf("Block number did not increase")
	}

	return newBlockNumber
}

func VerifyTx(t *testing.T, cbn uint64, usrList []TxStatus) {
	WaitMined(usrList[0].Txs, cbn)
	txBlk := make([]string, 0)

	verifyTransactions := func(txs types.Transactions, mined bool, result string, logMsg string) {
		t.Log(logMsg)
		for _, tx := range txs {
			blkNum := CheckBundleTx(t, *tx, mined, result)
			if mined {
				txBlk = append(txBlk, blkNum)
			}
		}
	}

	t.Log("Verify Mem_pool tx.")
	verifyTransactions(usrList[0].Txs, usrList[0].Mined, usrList[0].Rst, "Verifying Mem_pool transactions.")

	t.Log("Verify Bundle1 tx.")
	verifyTransactions(usrList[1].Txs, usrList[1].Mined, usrList[1].Rst, "Verifying Bundle1 transactions.")

	t.Log("Verify Bundle2 tx.")
	verifyTransactions(usrList[2].Txs, usrList[2].Mined, usrList[2].Rst, "Verifying Bundle2 transactions.")

	TxInSameBlk(txBlk)
}

func VerifyTx6(t *testing.T, cbn uint64, usrList6 []TxStatus) {
	WaitMined(usrList6[0].Txs, cbn)
	txBlk := make([]string, 0)

	verifyTransactions := func(txs types.Transactions, mined bool, result string, logMsg string) {
		t.Log(logMsg)
		for _, tx := range txs {
			blkNum := CheckBundleTx(t, *tx, mined, result)
			if mined {
				txBlk = append(txBlk, blkNum)
			}
		}
	}

	t.Log("Verify Mem_pool tx.")
	verifyTransactions(usrList6[0].Txs, usrList6[0].Mined, usrList6[0].Rst, "Verifying Mem_pool transactions.")

	t.Log("Verify Bundle1 tx.")
	verifyTransactions(usrList6[1].Txs, usrList6[1].Mined, usrList6[1].Rst, "Verifying Bundle1 transactions.")

	t.Log("Verify Bundle2 tx.")
	verifyTransactions(usrList6[2].Txs, usrList6[2].Mined, usrList6[2].Rst, "Verifying Bundle2 transactions.")
	t.Log("Verify Mem_pool tx.")
	verifyTransactions(usrList6[3].Txs, usrList6[3].Mined, usrList6[3].Rst, "Verifying Mem_pool transactions.")

	t.Log("Verify Bundle1 tx.")
	verifyTransactions(usrList6[4].Txs, usrList6[4].Mined, usrList6[4].Rst, "Verifying Bundle1 transactions.")

	t.Log("Verify Bundle2 tx.")
	verifyTransactions(usrList6[5].Txs, usrList6[5].Mined, usrList6[5].Rst, "Verifying Bundle2 transactions.")

	TxInSameBlk(txBlk)
}

func GetAccBalance(address common.Address) *big.Int {
	// 连接到 BSD 节点
	client, err := ethclient.Dial(conf.Url)
	if err != nil {
		log.Fatalf("Failed to connect to the BSD node: %v", err)
	}

	// 指定账户地址
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatalf("Failed to retrieve account balance: %v", err)
	}
	// 打印账户余额
	fmt.Printf("Address %v Balance: %s \n", address, balance.String())
	return balance

}
