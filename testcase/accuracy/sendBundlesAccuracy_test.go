package accuracy

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/xkwang/conf"
	"github.com/xkwang/sendBundle"
	"github.com/xkwang/testcase"
	"github.com/xkwang/utils"
	"log"
	"math/big"
	"strconv"
	"testing"
	"time"
)

func Test_reset(t *testing.T) {
	utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)
	// utils.SendLockMempool( conf.RootPk5, conf.Mylock, lock_data, false)

}

// func Test_con(t *testing.T) {
// 	// defer utils.ResetLockContract(t, conf.Mylock, reset_data)
// 	// tx_0, _ := utils.SendLockMempool( conf.RootPk5, conf.Mylock, lock_data, false)
// 	// time.Sleep(6 * time.Second)
// 	tx01, _ := utils.SendLockMempool( conf.RootPk2, conf.Mylock, unlock_de_data, false)
// 	// print(tx_0[0])
// 	print(tx01[0])
// 	time.Sleep(6 * time.Second)
// 	rcp := utils.GetTransactionReceipt(*tx01[0])
// 	println("rcp.Result.GasUsed")
// 	if rcp.Result.Status == "0x1" {
// 		println(rcp.Result.GasUsed)
// 	} else {
// 		println("Continue")
// 	}

// }

// SendBundles sends the transaction bundles concurrently.

// ExposeMempoolTransactions handles the mempool exposure for multiple transactions.
func ExposeMempoolTransactions(pk1 string, data1 []byte, pk2 string, data2 []byte, gas *big.Int) (types.Transactions, []common.Hash) {
	tx1, revertHash1 := utils.SendLockMempool(pk1, conf.Mylock, data1, gas, big.NewInt(conf.MinGasPrice), true, true)
	tx2, revertHash2 := utils.SendLockMempool(pk2, conf.Mylock, data2, gas, big.NewInt(conf.MinGasPrice), true, true)
	return append(tx1, tx2...), append(revertHash1, revertHash2...)
}

func TestMyFunction(t *testing.T) {
	const repeatCount = 20 // 设置要重复的次数

	for i := 0; i < repeatCount; i++ {
		t.Run("Iteration"+strconv.Itoa(i), func(t *testing.T) {
			// 你的测试逻辑在这里
			// 例如，调用你的函数并断言结果
			Test_P1_value(t)
			//Test_p0_back_run(t)
			//expected := "expected result"
			//if result != expected {
			//	t.Errorf("iteration %d: expected %s, got %s", i, expected, result)
			//}
		})
	}
}

func Test_P1_value(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		t.Log("[Step-1] Root User1 Expose Mem_pool transaction tx1 .\n")
		tx1, _ := utils.SendLockMempool(conf.RootPk5, conf.WBNB, conf.TransferWBNBCode, conf.LowGas, big.NewInt(conf.MinGasPrice), false, true)
		time.Sleep(100 * time.Millisecond)
		t.Log("[Step-2] Root User2 SendBundle [tx1, tx2] .\n")
		bundleArgs, arg, tx2 := testcase.AddUserBundle(conf.RootPk2, conf.WBNB, conf.TransferWBNBCode, conf.SendA, conf.LowGas, tx1, nil, 0)
		cbn := utils.SendBundlesMined(t, *arg, bundleArgs)
		utils.WaitMined(tx2, cbn)
		testcase.CheckTransactionIndex(t, *tx1[0], "0x0")
		testcase.CheckTransactionIndex(t, *tx2[0], "0x1")
	})
}

func chooseCurrentTime() {
	usr := utils.UserTx(conf.RootPk5, conf.WBNB, conf.TransferWBNBCode, conf.HighGas, big.NewInt(conf.MinGasPrice))
	initialBlockNumber, err := usr.Client.BlockNumber(usr.Ctx)
	if err != nil {
		log.Fatalf("Failed to get current block number: %v", err)
	}
	log.Printf("Initial block number: %d", initialBlockNumber)

	for {
		// 获取当前的区块高度
		currentBlockNumber, err := usr.Client.BlockNumber(usr.Ctx)
		if err != nil {
			log.Printf("Failed to get current block number: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// 检查区块高度是否增加
		if currentBlockNumber > initialBlockNumber {
			log.Println("block number increase, begin to sendBundle")
			break
		}
		// 等待一段时间，然后再次检查区块高度
		time.Sleep(1 * time.Second)
	}

}

func Test_p0_back_run(t *testing.T) {
	t.Run("back_run_tx1_success", func(t *testing.T) {
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)
		chooseCurrentTime()
		t.Log("[Step-1] Root User Expose mem-pool transaction tx1\n")
		txs, hs := utils.SendLockMempool(conf.RootPk5, conf.Mylock, testcase.LockData, conf.LowGas, big.NewInt(conf.MinGasPrice), false, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.\n")
		bundleArgs1, usr1Arg, tx1 := testcase.AddUserBundle(conf.RootPk2, conf.Mylock, testcase.UnlockMoreData, conf.SendA, conf.HighGas, txs, hs, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		bundleArgs2, usr2Arg, tx2 := testcase.AddUserBundle(conf.RootPk3, conf.Mylock, testcase.UnlockStrData, conf.SendA, conf.LowGas, txs, hs, 0)

		testcase.Args[0] = usr1Arg
		testcase.Args[1] = usr2Arg
		testcase.BundleargsLsit[0] = bundleArgs1
		testcase.BundleargsLsit[1] = bundleArgs2

		t.Log("[Step-4] User 1 and User 2 send bundles.\n")
		cbn := testcase.SendBundles(t, usr1Arg, usr2Arg, bundleArgs1, bundleArgs2)

		testcase.UpdateUsrList(0, txs, true, conf.TxSucceed)
		testcase.UpdateUsrList(1, tx1, true, conf.TxSucceed)
		testcase.UpdateUsrList(2, tx2, false, conf.TxFailed)

		t.Log("[Step-5] Transaction order check.\n")
		utils.VerifyTx(t, cbn, testcase.UsrList)
		// Expect [tx1, tx2] 校验链上交易顺序与 bundle-1 交易顺序一致
	})

	t.Run("back_run_tx1_failed", func(t *testing.T) {
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)
		chooseCurrentTime()

		t.Log("[Step-1] Root User Expose mem_pool transaction tx1")
		txs, hs := utils.SendLockMempool(conf.RootPk5, conf.Mylock, testcase.LockFData, conf.LowGas, big.NewInt(conf.MinGasPrice), true, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], none are allowed to revert.")
		bundleArgs1, usr1Arg, txs1 := testcase.AddUserBundle(conf.RootPk2, conf.Mylock, testcase.UnlockMoreData, conf.SendA, conf.HighGas, txs, hs, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], none are allowed to revert.")
		bundleArgs2, usr2Arg, txs2 := testcase.AddUserBundle(conf.RootPk3, conf.Mylock, testcase.UnlockStrData, conf.SendA, conf.LowGas, txs, hs, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles.")
		cbn := testcase.SendBundles(t, usr1Arg, usr2Arg, bundleArgs1, bundleArgs2)

		testcase.UpdateUsrList(0, txs, true, conf.TxFailed)
		testcase.UpdateUsrList(1, txs1, false, conf.TxFailed)
		testcase.UpdateUsrList(2, txs2, false, conf.TxFailed)

		t.Log("[Step-5] Verify transactions.")
		utils.VerifyTx(t, cbn, testcase.UsrList)
	})

	t.Run("sendBundle", func(t *testing.T) {
		// tx1 在revertList 且会revert tx2随意
		arg := utils.UserTx(conf.RootPk5, conf.WBNB, conf.TransferWBNBCode, conf.HighGas, big.NewInt(conf.MinGasPrice))
		arg.TxCount = 2
		arg.RevertList = []int{0}
		var txBlk = make([]string, 0)
		// 发送并验证交易
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)
		utils.WaitMined(txs, cbn)
		blkN1 := utils.CheckBundleTx(t, *txs[0], true, conf.TxFailed)
		txBlk = append(txBlk, blkN1)

		blkN2 := utils.CheckBundleTx(t, *txs[1], true, conf.TxSucceed)
		txBlk = append(txBlk, blkN2)

		if !utils.TxInSameBlk(txBlk) {
			t.Fatalf("Transactions are not in the same block")
		}
	})
}

func Test_p0_token_sniper(t *testing.T) {
	t.Run("tokenSniper_tx1_success", func(t *testing.T) {
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mem_pool transaction tx1")
		txs, revertHash := utils.SendLockMempool(conf.RootPk5, conf.Mylock, testcase.LockData, conf.LowGas, big.NewInt(conf.MinGasPrice), true, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], none are allowed to revert.")
		bundleArgs1, usr1Arg, txs1 := testcase.AddUserBundle(conf.RootPk2, conf.Mylock, testcase.FakelockMoreData, conf.SendA, conf.HighGas, txs, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], none are allowed to revert.")
		bundleArgs2, usr2Arg, txs2 := testcase.AddUserBundle(conf.RootPk3, conf.Mylock, testcase.FakelockStrData, conf.SendA, conf.LowGas, txs, revertHash, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles.")
		cbn := testcase.SendBundles(t, usr1Arg, usr2Arg, bundleArgs1, bundleArgs2)

		testcase.UpdateUsrList(0, txs, true, conf.TxSucceed)
		testcase.UpdateUsrList(1, txs1, true, conf.TxSucceed)
		testcase.UpdateUsrList(2, txs2, true, conf.TxSucceed)

		t.Log("[Step-5] Verify transactions.")
		utils.VerifyTx(t, cbn, testcase.UsrList)

		// Log target transaction sequence
		targetTxs := append(txs, txs1...)
		targetTxs = append(targetTxs, txs2...)
		for _, tx := range targetTxs {
			log.Println(tx.Hash().Hex())
		}
	})

	t.Run("tokenSniper_tx1_failed", func(t *testing.T) {
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mem_pool transaction tx1")
		txs, revertHash := utils.SendLockMempool(conf.RootPk5, conf.Mylock, testcase.LockFData, conf.LowGas, big.NewInt(conf.MinGasPrice), true, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], none are allowed to revert.")
		bundleArgs1, usr1Arg, txs1 := testcase.AddUserBundle(conf.RootPk2, conf.Mylock, testcase.FakelockMoreData, conf.SendA, conf.HighGas, txs, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], none are allowed to revert.")
		bundleArgs2, usr2Arg, txs2 := testcase.AddUserBundle(conf.RootPk3, conf.Mylock, testcase.FakelockStrData, conf.SendA, conf.LowGas, txs, revertHash, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles.")
		cbn := testcase.SendBundles(t, usr1Arg, usr2Arg, bundleArgs1, bundleArgs2)

		testcase.UpdateUsrList(0, txs, true, conf.TxFailed)
		testcase.UpdateUsrList(1, txs1, false, conf.TxFailed)
		testcase.UpdateUsrList(2, txs2, false, conf.TxFailed)

		t.Log("[Step-5] Verify transactions.")
		utils.VerifyTx(t, cbn, testcase.UsrList)
	})
}

func Test_p0_running_attack(t *testing.T) {
	t.Run("runningAttack_tx1_success", func(t *testing.T) {
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mem_pool transaction tx0 tx1")
		tx0, rh := utils.SendLockMempool(conf.RootPk4, conf.Mylock, testcase.UseGas, conf.MedGas, big.NewInt(conf.MinGasPrice), true, true)
		tx01, revertHash := utils.SendLockMempool(conf.RootPk5, conf.Mylock, testcase.LockData, conf.LowGas, big.NewInt(conf.MinGasPrice), true, true)
		tx02 := append(tx01, tx0...)
		revertHash = append(revertHash, rh[0])

		t.Log("[Step-2] User 1 bundle [tx0, tx1, tx2], tx2 not allowed to revert.")
		bundleArgs1, usr1Arg, txs1 := testcase.AddUserBundle(conf.RootPk2, conf.Mylock, testcase.UnlockStrData, conf.SendA, conf.MedGas, tx02, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.")
		bundleArgs2, usr2Arg, txs2 := testcase.AddUserBundle(conf.RootPk3, conf.Mylock, testcase.UnlockMoreData, conf.SendA, conf.HighGas, tx01, revertHash, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles.")
		cbn := testcase.SendBundles(t, usr1Arg, usr2Arg, bundleArgs1, bundleArgs2)

		testcase.UpdateUsrList(0, tx02, true, conf.TxSucceed)
		testcase.UpdateUsrList(1, txs1, false, conf.TxFailed)
		testcase.UpdateUsrList(2, txs2, true, conf.TxSucceed)
		utils.VerifyTx(t, cbn, testcase.UsrList)
	})

	t.Run("runningAttack_tx1_failed", func(t *testing.T) {
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mem_pool transaction tx0 tx1")
		tx0, _ := utils.SendLockMempool(conf.RootPk4, conf.WBNB, conf.TransferWBNBCode, conf.LowGas, big.NewInt(conf.MinGasPrice), false, true)
		tx01, _ := utils.SendLockMempool(conf.RootPk5, conf.Mylock, testcase.LockFData, conf.LowGas, big.NewInt(conf.MinGasPrice), false, true)
		tx02 := append(tx01, tx0...)

		t.Log("[Step-2] User 1 bundle [tx0, tx1, tx2], tx2 not allowed to revert.")
		bundleArgs1, usr1Arg, txs1 := testcase.AddUserBundle(conf.RootPk2, conf.Mylock, testcase.UnlockStrData, conf.SendA, conf.LowGas, tx02, nil, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.")
		bundleArgs2, usr2Arg, txs2 := testcase.AddUserBundle(conf.RootPk3, conf.Mylock, testcase.UnlockMoreData, conf.SendA, conf.HighGas, tx01, nil, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles.")
		testcase.SendBundles(t, usr1Arg, usr2Arg, bundleArgs1, bundleArgs2)

		t.Log("[Step-5] Verify transaction .\n")
		// for _, tx := range tx_02 {
		// 	utils.CheckBundleTx(t, *tx, true, conf.TxSucceed)
		// }
		// tx0 success
		// tx1 failed

		for _, tx := range txs1 {
			// 依次检查bundle中的交易是否成功上链
			log.Println("bundle 1 not mined")
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}

		for _, tx := range txs2 {
			log.Println("bundle 2 mined")
			// 依次检查bundle中的交易是否成功上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}

	})
}

func Test_p0_gasLimit_deception(t *testing.T) {
	t.Run("gasLimitDeception_tx1_success", func(t *testing.T) {
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mem_pool transaction tx1")
		txs, _ := utils.SendLockMempool(conf.RootPk5, conf.Mylock, testcase.LockData, conf.LowGas, big.NewInt(conf.MinGasPrice), false, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.")
		bundleArgs1, usr1Arg, txs1 := testcase.AddUserBundle(conf.RootPk2, conf.Mylock, testcase.UnlockDeMoreData, conf.SendA, conf.MedGas, txs, nil, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.")
		bundleArgs2, usr2Arg, txs2 := testcase.AddUserBundle(conf.RootPk3, conf.Mylock, testcase.UnlockDeStrData, conf.SendA, conf.HighGas, txs, nil, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles.")
		cbn := testcase.SendBundles(t, usr1Arg, usr2Arg, bundleArgs1, bundleArgs2)

		t.Log("[Step-5] Verify transactions.")
		testcase.UpdateUsrList(0, txs, true, conf.TxSucceed)
		testcase.UpdateUsrList(1, txs1, true, conf.TxSucceed)
		testcase.UpdateUsrList(2, txs2, false, conf.TxFailed)
		utils.VerifyTx(t, cbn, testcase.UsrList)
	})

	t.Run("gasLimitDeception_tx1_failed", func(t *testing.T) {
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mem_pool transaction tx1")
		txs, revertHash := utils.SendLockMempool(conf.RootPk5, conf.Mylock, testcase.LockFData, conf.LowGas, big.NewInt(conf.MinGasPrice), true, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.")
		bundleArgs1, usr1Arg, txs1 := testcase.AddUserBundle(conf.RootPk2, conf.Mylock, testcase.UnlockDeMoreData, conf.SendA, conf.LowGas, txs, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.")
		bundleArgs2, usr2Arg, txs2 := testcase.AddUserBundle(conf.RootPk3, conf.Mylock, testcase.UnlockDeStrData, conf.SendA, conf.HighGas, txs, revertHash, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles.")
		cbn := testcase.SendBundles(t, usr1Arg, usr2Arg, bundleArgs1, bundleArgs2)

		t.Log("[Step-5] Verify transactions.")
		testcase.UpdateUsrList(0, txs, true, conf.TxFailed)
		testcase.UpdateUsrList(1, txs1, false, conf.TxFailed)
		testcase.UpdateUsrList(2, txs2, false, conf.TxFailed)
		utils.VerifyTx(t, cbn, testcase.UsrList)
	})
}

func Test_p0_sandwich(t *testing.T) {

	t.Run("sandwich_ol1", func(t *testing.T) {
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mem_pool transaction tx0 tx1")
		tx01, revertHash := ExposeMempoolTransactions(conf.RootPk2, testcase.LockData, conf.RootPk5, testcase.FakelockStrData, conf.LowGas)

		t.Log("[Step-2] User 1 bundle [tx0, tx1, tx2], tx2 not allowed to revert.")
		bundleArgs1, usr1Arg, txs1 := testcase.AddUserBundle(conf.RootPk2, conf.Mylock, testcase.UnlockMoreData, conf.SendA, conf.HighGas, tx01, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.")
		bundleArgs2, usr2Arg, txs2 := testcase.AddUserBundle(conf.RootPk3, conf.Mylock, testcase.UnlockStrData, conf.SendA, conf.LowGas, tx01[1:], revertHash[1:], 0)

		t.Log("[Step-4] User 1 and User 2 send bundles. Expect mined Tx_list: [tx0, tx1, tx2]")
		cbn := testcase.SendBundles(t, usr1Arg, usr2Arg, bundleArgs1, bundleArgs2)

		t.Log("[Step-5] Check tx0, tx1, and tx2 mined.")
		testcase.UpdateUsrList(0, tx01, true, conf.TxSucceed)
		testcase.UpdateUsrList(1, txs1, true, conf.TxSucceed)
		testcase.UpdateUsrList(2, txs2, false, conf.TxFailed)
		utils.VerifyTx(t, cbn, testcase.UsrList)
	})

	t.Run("sandwich_both", func(t *testing.T) {
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mem_pool transaction tx0 tx1")
		tx01, revertHash := ExposeMempoolTransactions(conf.RootPk2, testcase.LockData, conf.RootPk5, testcase.UnlockStrData, conf.LowGas)

		t.Log("[Step-2] User 1 bundle [tx0, tx1, tx2], all not allowed to revert.")
		bundleArgs1, usr1Arg, txs1 := testcase.AddUserBundle(conf.RootPk2, conf.Mylock, testcase.ResetData, conf.SendA, conf.HighGas, tx01, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx1 not allowed to revert.")
		bundleArgs2, usr2Arg, txs2 := testcase.AddUserBundle(conf.RootPk3, conf.Mylock, testcase.ResetData, conf.SendA, conf.LowGas, tx01[1:], revertHash[1:], 0)

		t.Log("[Step-4] User 1 and User 2 send bundles. Expect mined Tx_list: [tx0, tx1, tx2, tx3]")
		cbn := testcase.SendBundles(t, usr1Arg, usr2Arg, bundleArgs1, bundleArgs2)

		t.Log("[Step-5] Check tx0, tx1, tx2, and tx3 mined.")
		testcase.UpdateUsrList(0, tx01, true, conf.TxSucceed)
		testcase.UpdateUsrList(1, txs1, true, conf.TxSucceed)
		testcase.UpdateUsrList(2, txs2, true, conf.TxSucceed)
		utils.VerifyTx(t, cbn, testcase.UsrList)
	})
}

func Test_p1_conflict_mb(t *testing.T) {
	t.Run("only mem_pool tx in bundle", func(t *testing.T) {
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mem_pool transaction tx0")
		tx0, revertTx := utils.SendLockMempool(conf.RootPk5, conf.Mylock, testcase.LockData, conf.LowGas, big.NewInt(conf.MinGasPrice), false, true)
		var txs types.Transactions
		t.Log("[Step-2] User 1 bundle [tx0].")
		userArg := utils.UserTx(conf.RootPk5, conf.Mylock, testcase.LockData, conf.HighGas, big.NewInt(conf.MinGasPrice))
		bundleArgs := utils.AddBundle(tx0, txs, revertTx, 0)
		err := userArg.BuilderClient.SendBundle(userArg.Ctx, bundleArgs)
		if err != nil {
			log.Println("failed: ", err.Error())
		}

		time.Sleep(6 * time.Second)
		t.Log("[Step-3] Verify the transaction.")
		response := utils.GetTransactionReceipt(*tx0[0])
		assert.Equal(t, response.Result.Status, conf.TxSucceed)
	})
	testCases := []struct {
		send     bool
		tx1Index string
		tx2Index string
	}{{true, "0x0", "0x1"}, {false, "0x1", "0x0"}}
	//	 1. private 不发expected: [tx2 tx1]
	//	 2. Mem_pool里有tx1 tx2则贵的先上
	for index, tc := range testCases {
		t.Run("mem_pool txs in bundle order check"+strconv.Itoa(index), func(t *testing.T) {
			// 生效需要注释掉SendLockMempool 的send
			defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)

			t.Log("[Step-2] Root User1 Expose mem_pool transaction  tx1 \n")
			tx1, revertHash := utils.SendLockMempool(conf.RootPk5, conf.Mylock, testcase.LockData, conf.HighGas, big.NewInt(conf.MinGasPrice), false, tc.send)

			t.Log("[Step-2] Root User2 Expose mem_pool transaction  tx2 \n")
			tx2, _ := utils.SendLockMempool(conf.RootPk2, conf.Mylock, testcase.LockData, conf.LowGas, big.NewInt(conf.MinGasPrice), false, true)

			t.Log("[Step-3] User3 send bundle [tx2, tx1].\n")
			usr1Arg := utils.UserTx(conf.RootPk3, conf.Mylock, testcase.UnlockMoreData, conf.HighGas, big.NewInt(conf.MinGasPrice))
			bundleArgs1 := utils.AddBundle(tx2, tx1, revertHash, 0)
			err := usr1Arg.BuilderClient.SendBundle(usr1Arg.Ctx, bundleArgs1)
			if err != nil {
				log.Println(" failed: ", err.Error())
			}
			time.Sleep(6 * time.Second)
			testcase.CheckTransactionIndex(t, *tx1[0], tc.tx1Index)
			testcase.CheckTransactionIndex(t, *tx2[0], tc.tx2Index)

		})
	}
}
