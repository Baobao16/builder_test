package bribe

import (
	"github.com/xkwang/testcase"
	"log"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/xkwang/cases"
	"github.com/xkwang/conf"
	"github.com/xkwang/utils"
)

var ValueCpABI = utils.GeneABI(conf.ValueCp_path)
var valueCpABI = utils.Contract{ABI: *ValueCpABI}
var SpeABI = utils.GeneABI(conf.Spe_path)
var speABI = utils.Contract{ABI: *SpeABI}

var betT = utils.GeneEncodedData(valueCpABI, "bet", true)
var betF = utils.GeneEncodedData(valueCpABI, "bet", false)

//var getCoinbase = utils.GeneEncodedData(speABI, "getCoinbase")

func Test_p0_value_preservation(t *testing.T) {
	var txs types.Transactions
	bribeFee1 := big.NewInt(0.01 * 1e18)
	bribeFee2 := big.NewInt(0.005 * 1e18)
	testCases := []struct {
		a         []byte
		b         []byte
		aMinted   bool // tx1 是否上链
		bMinted   bool // tx2 是否上链
		aContract bool // tx1 合约参数
	}{
		{betT, betT, true, true, true},   // 链上交易顺序：[tx1,tx2]
		{betT, betF, true, false, true},  // 链上交易顺序：[tx1]
		{betF, betT, true, true, false},  // 链上交易顺序：[tx1,tx2]
		{betF, betF, true, false, false}, // 链上交易顺序：[tx1]
	}
	for index, tc := range testCases {
		t.Run("backRun_value_preservation"+strconv.Itoa(index), func(t *testing.T) {
			//defer utils.ResetContract(t, conf.Mylock, testcase.ResetData)

			t.Logf("[Step-1] User 1 bundle [tx1], tx1 Contract bet_%v .\n", tc.aContract)
			usr1Arg := utils.User_tx(conf.RootPk2, conf.ValueCp, tc.a, conf.High_gas)
			txs1, revertTxHashes := cases.GenerateBNBTxs(&usr1Arg, bribeFee1, usr1Arg.Data, 1)
			bundleArgs1 := utils.AddBundle(txs, txs1, revertTxHashes, 0)

			t.Logf("[Step-2] User 2 bundle [tx2], tx2 Contract bet_%v .\n", tc.bMinted)
			usr2Arg := utils.User_tx(conf.RootPk3, conf.ValueCp, tc.b, conf.High_gas)
			txs2, revertTxHashes := cases.GenerateBNBTxs(&usr2Arg, bribeFee2, usr2Arg.Data, 1)
			blockNum, _ := usr1Arg.Client.BlockNumber(usr1Arg.Ctx)
			t.Logf("Current Blk_num : %v .\n", blockNum)
			MaxBN := blockNum + 1
			bundleArgs2 := utils.AddBundle(txs, txs2, revertTxHashes, MaxBN)

			testcase.Args[0] = &usr1Arg
			testcase.Args[1] = &usr2Arg
			testcase.BundleArgs_lsit[0] = bundleArgs1
			testcase.BundleArgs_lsit[1] = bundleArgs2
			t.Log("[Step-3] User 1 and User 2 send bundles .\n")
			utils.ConcurSendBundles(t, testcase.Args, testcase.BundleArgs_lsit)

			t.Log("[Step-4] Check tx Minted .\n")

			response1 := utils.GetTransactionReceipt(*txs1[0])
			tx1Index := response1.Result.TransactionIndex

			blk, _ := strconv.ParseInt(strings.TrimPrefix(response1.Result.BlockNumber, "0x"), 16, 64)
			t.Logf("Tx1 in Blk : %v .\n", blk)
			assert.Equal(t, blk, int64(MaxBN))
			response2 := utils.GetTransactionReceipt(*txs2[0])
			tx2Index := response2.Result.TransactionIndex
			if tc.aMinted {
				assert.Equal(t, tx1Index, "0x0", "tx Index wrong", txs1[0].Hash().Hex())
			} else {
				assert.Equal(t, tx1Index, "")
			}
			if tc.bMinted {
				assert.Equal(t, response1.Result.BlockNumber, response2.Result.BlockNumber, "tx1 tx2 in diff Block")
				assert.Equal(t, tx2Index, "0x1")
			} else {
				assert.Equal(t, tx2Index, "")
			}
			//	237052

		})
	}
}

func Test_bribe(t *testing.T) {

	t.Run("bribe_success", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mempoool transaction tx1\n")
		lock_data := utils.GeneEncodedData(testcase.LockABI, "lock", 1, true)
		txs, _ := utils.SendLockMempool(conf.RootPk, conf.Mylock, lock_data, false)
		// tx1 gaslimit 30w, gasprice 1Gwei,
		// tx2 gasLimit 30w, 高于 tx3  >10w,
		// tx3 gasLimit 20w, gasPrice都是1Gwei,
		// tx4 gaslimit 30w, gasprice 1Gwei , 贿赂 SendAmount = 0.00015 * 1Gwei【贿赂成功】
		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.\n")

		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, testcase.UnlockMoreData, conf.High_gas)

		txs_1, revertTxHashes := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertTxHashes, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, testcase.UnlockStrData, conf.Med_gas)
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)

		//  Bribe Transaction 【private tx】
		arg := utils.User_tx(conf.RootPk4, conf.SpecialOp, conf.SpecialOp_Bb, conf.High_gas)
		bribe_fee := big.NewInt(1500000 * 1e9)
		log.Printf("bribe price is %v", bribe_fee)
		tmp := arg.Contract
		arg.Contract = conf.System_add
		txb, revertTxHashes := cases.GenerateBNBTxs(&arg, bribe_fee, arg.Data, 1)
		arg.Contract = tmp

		txs = append(txb, txs...)
		bundleArgs2 := utils.AddBundle(txs, txs_2, revertTxHashes, 0)

		testcase.Args[0] = &usr1_arg
		testcase.Args[1] = &usr2_arg
		testcase.BundleArgs_lsit[0] = bundleArgs1
		testcase.BundleArgs_lsit[1] = bundleArgs2

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")
		cbn := utils.ConcurSendBundles(t, testcase.Args, testcase.BundleArgs_lsit)

		testcase.UsrList[0].Txs = txs
		testcase.UsrList[0].Mined = true
		testcase.UsrList[0].Rst = conf.Txsucceed

		testcase.UsrList[1].Txs = txs_1
		testcase.UsrList[1].Mined = false
		testcase.UsrList[1].Rst = conf.Txfailed

		testcase.UsrList[2].Txs = txs_2
		testcase.UsrList[2].Mined = true
		testcase.UsrList[2].Rst = conf.Txsucceed

		utils.Verifytx(t, cbn, testcase.UsrList)
		// 交易顺序符合预期
		// tx1,tx3

	})
	t.Run("bribe_failed", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mem_pool transaction  tx1\n")
		lock_data := utils.GeneEncodedData(testcase.LockABI, "lock", 1, true)
		txs, _ := utils.SendLockMempool(conf.RootPk, conf.Mylock, lock_data, false)
		// tx1 gaslimit 30w, gasprice 1Gwei,
		// tx2 gasLimit 30w, 高于 tx3  >10w,
		// tx3 gasLimit 20w, gasPrice都是1Gwei,
		// tx4 gaslimit 30w, gasprice 1Gwei , 贿赂 SendAmount = 0.00005 * 1Gwei【贿赂失败】
		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.\n")

		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, testcase.UnlockMoreData, conf.High_gas)

		txs_1, revertTxHashes := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertTxHashes, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, testcase.UnlockStrData, conf.Med_gas)
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)

		//  Bribe Transaction 【private tx】
		arg := utils.User_tx(conf.RootPk4, conf.SpecialOp, conf.SpecialOp_Bb, conf.High_gas)
		bribe_fee := big.NewInt(50000 * 1e9)
		log.Printf("bribe price is %v", bribe_fee)
		tmp := arg.Contract
		arg.Contract = conf.System_add
		txb, revertTxHashes := cases.GenerateBNBTxs(&arg, bribe_fee, arg.Data, 1)
		arg.Contract = tmp

		txs0 := append(txb, txs...)
		bundleArgs2 := utils.AddBundle(txs0, txs_2, revertTxHashes, 0)

		testcase.Args[0] = &usr1_arg
		testcase.Args[1] = &usr2_arg
		testcase.BundleArgs_lsit[0] = bundleArgs1
		testcase.BundleArgs_lsit[1] = bundleArgs2

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")
		cbn := utils.ConcurSendBundles(t, testcase.Args, testcase.BundleArgs_lsit)

		testcase.UsrList[0].Txs = txs
		testcase.UsrList[0].Mined = true
		testcase.UsrList[0].Rst = conf.Txsucceed

		testcase.UsrList[1].Txs = txs_1
		testcase.UsrList[1].Mined = true
		testcase.UsrList[1].Rst = conf.Txsucceed

		testcase.UsrList[2].Txs = txs_2
		testcase.UsrList[2].Mined = false
		testcase.UsrList[2].Rst = conf.Txfailed

		utils.Verifytx(t, cbn, testcase.UsrList)
		// 交易顺序符合预期
		// tx1,tx2

	})

}

func Test_p0_SpecialOp(t *testing.T) {
	t.Run("testCoinbase", func(t *testing.T) {
		t.Log("send bundles testCoinbase \n")
		arg := utils.User_tx(conf.RootPk, conf.SpecialOp, conf.SpecialOp_Cb, conf.High_gas)
		arg.TxCount = 1
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(&arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)

		utils.WaitMined(txs, cbn)

		utils.CheckBundleTx(t, *txs[0], true, conf.Txsucceed)

	})

	// 	// 意义不大
	// 	// t.Skip("testTimestamp", func(t *testing.T) {
	// 	// 	t.Log("send bundles testTimestamp \n")
	// 	// 	blkTims := utils.GetLatestBlkMsg(t,conf.Spe_path, "testTimestamp", 5)
	// 	// 	arg := utils.User_tx(conf.RootPk2, conf.SpecialOp, blkTims)
	// 	// 	arg.TxCount = 1
	// 	// 	txs, bundletestcase.Args, _ := cases.ValidBundle_NilPayBidTx_1( &arg)
	// 	// 	cbn := utils.SendBundlesMined(t, arg, bundleArgs)

	// 	// 	utils.WaitMined(txs, cbn)

	// 	// 	utils.CheckBundleTx(t, *txs[0], true, conf.Txsucceed)

	// 	// })

	t.Run("testTimestampEq_parent", func(t *testing.T) {
		t.Log("send bundles testTimestamp parent \n")
		blkTims := utils.GetLatestBlkMsg(t, conf.Spe_path, "testTimestampEq", 5)

		arg := utils.User_tx(conf.RootPk, conf.SpecialOp, blkTims, conf.High_gas)
		arg.TxCount = 1
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(&arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)

		utils.WaitMined(txs, cbn)

		utils.CheckBundleTx(t, *txs[0], true, conf.Txsucceed)

	})

	t.Run("blkHash", func(t *testing.T) {
		t.Log("send bundles blkHash \n")
		blkHash := utils.GetLatestBlkMsg(t, conf.Spe_path, "testBlockHash", 0)
		arg := utils.User_tx(conf.RootPk2, conf.SpecialOp, blkHash, conf.High_gas)
		arg.TxCount = 1
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(&arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)

		utils.WaitMined(txs, cbn)

		utils.CheckBundleTx(t, *txs[0], true, conf.Txsucceed)

	})

}
