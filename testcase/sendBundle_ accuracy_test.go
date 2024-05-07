package newtestcases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/xkwang/cases"
)

var tx_type = "Transfer"

func user_tx(root_name string) cases.BidCaseArg {
	ctx := context.Background()

	rootPk := root_name
	bobPk := root_name
	builderPk := *builderPrivateKey

	client, err := ethclient.Dial(url)
	if err != nil {
		fmt.Println("node ethclient.DialOptions", "err", err)
	}

	client2, err := ethclient.Dial(url_1)
	if err != nil {
		fmt.Println("client2 bidclient ethclient.DialOptions", "err", err)
	}

	client3, err := ethclient.Dial(url)
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
		Validators:    []common.Address{common.HexToAddress(*validator)},
		BidClient:     client2,
		BuilderClient: client3,
		TxCount:       3,
		Contract:      WBNB,
		//Data:          common.Hex2Bytes("0df97361000000000000000000000000d9145cce52d386f254917e481eb44e9943f39138"), // lock 21432
		Data:       cases.TransferWBNB_code,
		GasPrice:   big.NewInt(350000),
		GasLimit:   big.NewInt(350000),
		SendAmount: big.NewInt(0),
	}
	return *arg
}

func addBundle(txs types.Transactions, txs_new types.Transactions, revertTxHashes []common.Hash) *types.SendBundleArgs {
	// 构造新的bundle，包含Mempool交易tx1
	txBytes := make([]hexutil.Bytes, 0)
	txByte, _ := txs[0].MarshalBinary()
	txBytes = append(txBytes, txByte)
	for _, tx := range txs_new {
		txByte, err := tx.MarshalBinary()
		fmt.Printf("txhash %v\n", tx.Hash().Hex())
		if err != nil {
			log.Println("tx.MarshalBinary", "err", err)
		}
		txBytes = append(txBytes, txByte)
	}

	bundleArgs := &types.SendBundleArgs{
		Txs:               txBytes,
		RevertingTxHashes: revertTxHashes,
	}
	bidJson, _ := json.MarshalIndent(bundleArgs, "", "  ")
	log.Println(string(bidJson))
	return bundleArgs

}

func tearDown(t *testing.T) {
	// 执行测试后的清理工作:调用reset合约重置lock
	t.Log("Root User reset Contract lock\n")
	usr_arg := user_tx(rootPk)
	usr_arg.Contract = ugMylock
	usr_arg.GasLimit = big.NewInt(ugunlock_gas)
	usr_arg.GasPrice = big.NewInt(ugunlock_gas)
	usr_arg.TxCount = 1
	usr_arg.Data = mylock56_reset
	txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, &usr_arg)
	err := usr_arg.BuilderClient.SendBundle(usr_arg.Ctx, bundleArgs)
	assert.Nil(t, err)
	BlockheightIncreased(t)
	checkBundleTx(t, *txs[0], true, Txsucceed, tx_type)

}

func sendLockMempool(t *testing.T, contract common.Address, data []byte) types.Transactions {
	t.Log("Root User mempool transaction[tx1 ] Contract lock with [0,true]\n")
	usr_arg := user_tx(rootPk)
	usr_arg.Contract = contract
	usr_arg.GasLimit = big.NewInt(ugunlock_gas)
	usr_arg.GasPrice = big.NewInt(ugunlock_gas)
	usr_arg.TxCount = 1
	usr_arg.Data = data

	txs, _ := cases.GenerateBNBTxs(&usr_arg, usr_arg.SendAmount, usr_arg.Data, 1)
	err := usr_arg.Client.SendTransaction(usr_arg.Ctx, txs[0])
	fmt.Printf("txhash %v\n", txs[0].Hash())
	if err != nil {
		fmt.Println("failed to send single Transaction", "err", err)
	}
	return txs

}
func Test_p0_backrun(t *testing.T) {
	defer tearDown(t)
	t.Run("sendvalidbundle_tx1", func(t *testing.T) {
		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		txs := sendLockMempool(t, ugMylock, mylock56_lock0t) //mylock56_lock0t //"gasUsed":"0x342b"
		var err error

		t.Log("[Step-2] User 1 bundle [tx1, tx2], none are allowed to revert.\n")
		usr1_arg := user_tx(rootPk2)
		usr1_arg.Contract = ugMylock
		usr1_arg.Data = ugMylock_unlock_long_code //"gasUsed":"0xc944"

		txs_1, revertTxHashes := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := addBundle(txs, txs_1, revertTxHashes)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], none are allowed to revert.\n")
		usr2_arg := user_tx(rootPk3)
		usr2_arg.Contract = ugMylock
		// usr2_arg.RevertList = []int{0}
		usr2_arg.Data = ugMylock_unlock_str_code //"gasUsed":"0xbfd8"
		txs_2, revertTxHashes := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := addBundle(txs, txs_2, revertTxHashes)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")
		args := make([]*cases.BidCaseArg, 2)
		bundleArgs_lsit := make([]*types.SendBundleArgs, 2)
		args[0] = &usr1_arg
		args[1] = &usr2_arg
		bundleArgs_lsit[0] = bundleArgs1
		bundleArgs_lsit[1] = bundleArgs2
		wg := sync.WaitGroup{}
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(i int) {
				err = args[i].BuilderClient.SendBundle(args[i].Ctx, bundleArgs_lsit[i])
				if err != nil {
					log.Println(" failed: ", err.Error())
					assert.True(t, strings.Contains(err.Error(), InvalidTx))
				}
				wg.Done()
			}(i)
		}
		wg.Wait()

		BlockheightIncreased(t)
		// 在tx1成功执行的前提下，tx2和tx3只有一个能够成功
		checkBundleTx(t, *txs[0], true, Txsucceed, tx_type)
		for _, tx := range txs_1 {
			// 依次检查bundle中的交易是否成功上链
			println(tx)
			checkBundleTx(t, *tx, true, Txsucceed, tx_type)
			// 交易gasfee
		}
		// todo: tx2的gasfee(gasused * gasprice) 高于 tx3
		for _, tx := range txs_2 {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}
		checkBundleTx(t, *txs[0], true, Txsucceed, tx_type)

	})

}
func Test_p0_token_sniper(t *testing.T) {
	t.Run("sendvalidbundle_token_sniper", func(t *testing.T) {
		// tx2,tx3 不允许revert，tx1可以revert
		t.Log("Root User Expose mempool transaction tx1\n")
		usr_arg := user_tx(rootPk)
		usr_arg.Contract = Lock
		usr_arg.Data = Lock_fakel_code //
		usr_arg.SendAmount = big.NewInt(0)

		// 生成 user-1 发起的【公开mempool交易tx1】
		txs, _ := cases.GenerateBNBTxs(&usr_arg, usr_arg.SendAmount, usr_arg.Data, 1)
		err := usr_arg.Client.SendTransaction(usr_arg.Ctx, txs[0])
		fmt.Printf("txhash %v\n", txs[0].Hash())
		if err != nil {
			fmt.Println("failed to send single Transaction", "err", err)
		}
		t.Log("User 1 sends bundle [tx1, tx2], none are allowed to revert.\n")
		usr1_arg := user_tx(rootPk2)
		usr1_arg.Contract = Lock
		usr1_arg.Data = Lock_lockshort_code
		txs_1, revertTxHashes := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := addBundle(txs, txs_1, revertTxHashes)

		t.Log("User 2 sends bundle [tx1, tx2], none are allowed to revert.\n")
		usr2_arg := user_tx(rootPk3)
		usr2_arg.Contract = Lock
		usr2_arg.Data = Lock_locklong_code // gas高于short
		txs_2, revertTxHashes := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := addBundle(txs, txs_2, revertTxHashes)

		args := make([]*cases.BidCaseArg, 2)
		bundleArgs_lsit := make([]*types.SendBundleArgs, 2)
		args[0] = &usr1_arg
		args[1] = &usr2_arg
		bundleArgs_lsit[0] = bundleArgs1
		bundleArgs_lsit[1] = bundleArgs2
		wg := sync.WaitGroup{}
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(i int) {
				time.Sleep(time.Duration(i) * time.Second)
				err = args[i].BuilderClient.SendBundle(args[i].Ctx, bundleArgs_lsit[i])
				// assert.Nil(t, err)
				if err != nil {
					log.Println(" failed: ", err.Error())
					assert.True(t, strings.Contains(err.Error(), InvalidTx))
				}
				wg.Done()
			}(i)
		}
		wg.Wait()

		BlockheightIncreased(t)
		// 在tx1成功执行的前提下，tx2和tx3只有一个能够成功
		checkBundleTx(t, *txs[0], true, Txsucceed, tx_type)
		for _, tx := range txs_1[1:] {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, true, Txsucceed, tx_type)
		}
		// todo: tx2的gasfee(gasused * gasprice) 高于 tx3
		for _, tx := range txs_2 {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}
	})

}
func Test_p0_running_attack(t *testing.T) {
	t.Run("sendvalidbundle_Running_Attack", func(t *testing.T) {})

}
func Test_p0_gaslimit_deception(t *testing.T) {
	t.Run("sendvalidbundle_gaslimit_deception", func(t *testing.T) {})

}
func Test_p0_sandwich(t *testing.T) {
	t.Run("sandwich_both", func(t *testing.T) {})
	t.Run("sandwich_only1", func(t *testing.T) {})

}
