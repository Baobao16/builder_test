package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/xkwang/cases"
	"github.com/xkwang/conf"
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

func User_tx(root_name string, contract common.Address, data []byte, gasLimit *big.Int) cases.BidCaseArg {
	ctx := context.Background()

	rootPk := root_name
	bobPk := root_name
	builderPk := *conf.BuilderPrivateKey

	client, err := ethclient.Dial(conf.Url)
	if err != nil {
		fmt.Println("node ethclient.DialOptions", "err", err)
	}

	client2, err := ethclient.Dial(conf.Url_1)
	if err != nil {
		fmt.Println("client2 bidclient ethclient.DialOptions", "err", err)
	}

	client3, err := ethclient.Dial(conf.Url)
	if err != nil {
		fmt.Println("client3 bidclient ethclient.DialOptions", "err", err)
	}

	//query chainID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		fmt.Printf("err %v\n", err)
	} else {
		fmt.Printf("chainID %v\n", chainID)
	}

	arg := &cases.BidCaseArg{
		Ctx:           ctx,
		Client:        client,
		ChainID:       chainID,
		RootPk:        rootPk,
		BobPk:         bobPk,
		Builder:       cases.NewAccount(builderPk),
		Validators:    []common.Address{common.HexToAddress(*conf.Validator)},
		BidClient:     client2,
		BuilderClient: client3,
		TxCount:       3,
		Contract:      contract,
		Data:          data,
		GasPrice:      big.NewInt(conf.Min_gasPrice),
		GasLimit:      gasLimit,
		SendAmount:    big.NewInt(0),
	}
	return *arg
}

func AddBundle(txs types.Transactions, txs_new types.Transactions, revertTxHashes []common.Hash, MaxBN uint64) *types.SendBundleArgs {
	// 构造新的bundle，包含Mempool交易
	txBytes := make([]hexutil.Bytes, 0)

	for _, tx := range txs {
		txByte, err := tx.MarshalBinary()
		// fmt.Printf("txhash %v\n", tx.Hash().Hex())
		if err != nil {
			log.Println("tx.MarshalBinary", "err", err)
		}
		txBytes = append(txBytes, txByte)
	}
	for _, tx := range txs_new {
		txByte, err := tx.MarshalBinary()
		// log.Printf("Private txhash %v\n", tx.Hash().Hex())
		if err != nil {
			log.Println("tx.MarshalBinary", "err", err)
		}
		txBytes = append(txBytes, txByte)
	}

	bundleArgs := &types.SendBundleArgs{
		Txs:               txBytes,
		RevertingTxHashes: revertTxHashes,
		MaxBlockNumber:    MaxBN,
	}

	bidJson, _ := json.MarshalIndent(bundleArgs, "", "  ")
	log.Println(string(bidJson))
	return bundleArgs

}
func IsEmptyField(result Result_b) bool {
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

func SendLockMempool(t *testing.T, usr string, contract common.Address, data []byte, revert bool) (types.Transactions, []common.Hash) {

	usr_arg := User_tx(usr, contract, data, big.NewInt(2000000))
	// usr_arg.GasLimit = big.NewInt(conf.Max_gasLimit)
	// usr_arg.GasPrice = big.NewInt(conf.Min_gasPrice)
	if revert {
		log.Printf("Mempool transaction will in bundle RevertList . ")
		usr_arg.RevertListnormal = []int{0} // 当前交易被记入RevertList
	}
	log.Printf("Set Mempool transaction ")
	tx, revertHash := cases.GenerateBNBTxs(&usr_arg, usr_arg.SendAmount, usr_arg.Data, 1)
	// txBytes := make([]hexutil.Bytes, 0)
	// for _, tx := range tx {
	// 	txByte, err := tx.MarshalBinary()
	// 	fmt.Printf("sendLockMempool txhash %v\n", tx.Hash().Hex())
	// 	if err != nil {
	// 		log.Println("tx.MarshalBinary", "err", err)
	// 	}
	// 	txBytes = append(txBytes, txByte)
	// }
	// log.Printf("Set Mempool transaction %v [gasPrice: %v , gasLimit : %v] \n", tx[0].Hash(), usr_arg.GasPrice, usr_arg.GasLimit)
	err := usr_arg.Client.SendTransaction(usr_arg.Ctx, tx[0])

	if err != nil {
		fmt.Println("failed to send single Transaction", "err", err)
	}
	return tx, revertHash

}

func ConcurSendBundles(t *testing.T, args []*cases.BidCaseArg, bundleArgs_lsit []*types.SendBundleArgs) uint64 {
	cb, _ := args[0].Client.BlockNumber(args[0].Ctx)
	log.Println("sendBundel and waiting for blk_num increased")

	wg := sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			// time.Sleep(time.Duration(i) * time.Second)
			err := args[i].BuilderClient.SendBundle(args[i].Ctx, bundleArgs_lsit[i])
			// msg := "non-reverting tx in bundle failed"
			if err != nil {
				log.Println(" failed: ", err.Error())
				assert.True(t, strings.Contains(err.Error(), "non-reverting tx in bundle failed"))
				// assert.True(t, strings.Contains(err.Error(), conf.InvalidTx))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	time.Sleep(5 * time.Second)
	cbn, _ := args[0].Client.BlockNumber(args[0].Ctx)
	assert.True(t, cbn > cb, " blk_num not increased")
	return cbn
}

func Changetype(arg interface{}) int {
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
	// step-2 传入method和参数 返回Encoded Data
	var bigInt big.Int
	var contract_data []byte
	var err error
	if len(args) == 0 {
		contract_data, err = con.CallFunction(method)
	} else if len(args) == 2 {
		val := Changetype(args[0])
		log.Println(val)
		bigInt = *bigInt.SetInt64(int64(val))

		args2 := []interface{}{&bigInt, args[1]}

		contract_data, err = con.CallFunction(method, args2...)
	} else if len(args) == 1 {
		contract_data, err = con.CallFunction(method, args[0])
	} else {
		log.Fatalln("args nums wrong")
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("调用智能合约函数 %s, 传递参数值 %v\n", method, args) //, 交易数据: %x
	return contract_data
}

func ResetContract(t *testing.T, contract common.Address, data []byte) {
	// 执行测试后的清理工作:调用reset合约重置lock
	t.Log("Root User reset Contract lock\n")
	usr_arg := User_tx(conf.RootPk, contract, data, conf.High_gas) //"gasUsed":"0x5bb2"  23474
	usr_arg.TxCount = 1
	txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, &usr_arg)
	bn, _ := usr_arg.Client.BlockNumber(usr_arg.Ctx)
	log.Printf("Current Block height is: %v\n", bn)

	cbn := SendBundlesMined(t, usr_arg, bundleArgs)

	WaitMined(txs, cbn)

	blk_num := CheckBundleTx(t, *txs[0], true, conf.Txsucceed)
	log.Printf("bundle in Blk : %v \n", blk_num)

}

func SendBundlesMined(t *testing.T, usr cases.BidCaseArg, bundleArgs *types.SendBundleArgs) uint64 {
	cb, _ := usr.Client.BlockNumber(usr.Ctx)
	log.Println("sendBundel and waiting for blk_num increased")
	err := usr.BuilderClient.SendBundle(usr.Ctx, bundleArgs)
	assert.Nil(t, err)

	time.Sleep(6 * time.Second)
	cbn, _ := usr.Client.BlockNumber(usr.Ctx)
	assert.True(t, cbn > cb, " blk_num not increased")
	return cbn
}

func Verifytx(t *testing.T, cbn uint64, usrList []TxStatus) {
	WaitMined(usrList[0].Txs, cbn)
	tx_blk := make([]string, 0)
	t.Log("Verify menmpool tx .\n")
	for _, tx := range usrList[0].Txs {
		blk_num := CheckBundleTx(t, *tx, usrList[0].Mined, usrList[0].Rst)
		// 若交易上链，记录交易所在区块号
		if usrList[0].Mined {
			tx_blk = append(tx_blk, blk_num)
		}
	}
	// 依次检查bundle中的交易是否成功上链
	t.Log("Verify Bundle1 tx .\n")
	for _, tx := range usrList[1].Txs {
		blk_num := CheckBundleTx(t, *tx, usrList[1].Mined, usrList[1].Rst)
		// 若交易上链，记录交易所在区块号
		if usrList[1].Mined {
			tx_blk = append(tx_blk, blk_num)
		}
	}
	t.Log("Verify Bundle2 tx .\n")
	for _, tx := range usrList[2].Txs {
		blk_num := CheckBundleTx(t, *tx, usrList[2].Mined, usrList[2].Rst)
		if usrList[2].Mined {
			// 若交易上链，记录交易所在区块号
			tx_blk = append(tx_blk, blk_num)
		}
	}
	TxinSameBlk(tx_blk)

}
