package newtestcases

import (
	"log"
	"math/big"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/xkwang/cases"
	"github.com/xkwang/conf"
	"github.com/xkwang/utils"
)

var ValueCpABI = utils.GeneABI(conf.ValueCp_path)

var valuecpABI = utils.Contract{ABI: *ValueCpABI}

// var SpeABI = utils.GeneABI(conf.Spe_path)

var bet_t = utils.GeneEncodedData(valuecpABI, "bet", true)
var bet_f = utils.GeneEncodedData(valuecpABI, "bet", false)

func Test_p0_value_perservation(t *testing.T) {
	var txs types.Transactions
	bribe_fee_1 := big.NewInt(0.01 * 1e18)
	bribe_fee_2 := big.NewInt(0.005 * 1e18)
	testCases := []struct {
		a        []byte
		b        []byte
		a_minted bool //
		b_minted bool
	}{
		{bet_t, bet_t, true, true},  // 链上交易顺序：[tx1,tx2]
		{bet_t, bet_f, true, false}, //链上交易顺序：[tx1]
		{bet_f, bet_t, true, true},  // 链上交易顺序：[tx1,tx2]
		{bet_f, bet_f, true, false}, //链上交易顺序：[tx1]
	}
	for index, tc := range testCases {
		t.Run("backrun_value_perservation"+strconv.Itoa(index), func(t *testing.T) {
			// defer utils.ResetContract(t, conf.Mylock, reset_data)

			t.Log("[Step-1] User 1 bundle [tx1], tx1 not allowed to revert.\n")
			usr1_arg := utils.User_tx(conf.RootPk2, conf.ValueCp, tc.a)
			txs_1, revertTxHashes := cases.GenerateBNBTxs(&usr1_arg, bribe_fee_1, usr1_arg.Data, 1)
			bundleArgs1 := utils.AddBundle(txs, txs_1, revertTxHashes, 0)

			t.Log("[Step-2] User 2 bundle [tx2], tx2 not allowed to revert.\n")
			usr2_arg := utils.User_tx(conf.RootPk3, conf.ValueCp, tc.b)
			blockNum, _ := usr2_arg.Client.BlockNumber(usr2_arg.Ctx)
			MaxBN := blockNum + 1
			txs_2, revertTxHashes := cases.GenerateBNBTxs(&usr2_arg, bribe_fee_2, usr2_arg.Data, 1)
			bundleArgs2 := utils.AddBundle(txs, txs_2, revertTxHashes, MaxBN)

			args[0] = &usr1_arg
			args[1] = &usr2_arg
			bundleArgs_lsit[0] = bundleArgs1
			bundleArgs_lsit[1] = bundleArgs2
			t.Log("[Step-4] User 1 and User 2 send bundles .\n")
			utils.ConcurSendBundles(t, args, bundleArgs_lsit)

			response1 := utils.GetTransactionReceipt(*txs_1[0])
			tx1_index := response1.Result.TransactionIndex
			response2 := utils.GetTransactionReceipt(*txs_2[0])
			tx2_index := response2.Result.TransactionIndex
			log.Println(tx1_index, tx2_index)
			if tc.a_minted {
				assert.Equal(t, tx1_index, "0x0")
			}
			if tc.b_minted {
				assert.Equal(t, tx2_index, "0x1")
			} else {
				assert.Equal(t, tx2_index, "")
			}

		})
	}
}

func Test_bribe(t *testing.T) {

	t.Run("bribe_success", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, reset_data)

		t.Log("[Step-1] Root User Expose mempoool transaction tx1\n")
		lock_data := utils.GeneEncodedData(lockABI, "lock", 1, true)
		txs, _ := utils.SendLockMempool(t, conf.RootPk, conf.Mylock, lock_data, false)
		// tx1 gaslimit 30w, gasprice 1Gwei,
		// tx2 gasLimit 30w, 高于 tx3  >10w,
		// tx3 gasLimit 20w, gasPrice都是1Gwei,
		// tx4 gaslimit 30w, gasprice 1Gwei , 贿赂 SendAmount = 0.00015 * 1Gwei【贿赂成功】
		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.\n")

		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, unlock_more_data)
		usr1_arg.GasLimit = big.NewInt(3e6)

		txs_1, revertTxHashes := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertTxHashes, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, unlock_str_data)
		usr2_arg.GasLimit = big.NewInt(2e6)
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)

		//  Bribe Transaction 【private tx】
		arg := utils.User_tx(conf.RootPk4, conf.SpecialOp, conf.SpecialOp_Bb)
		bribe_fee := big.NewInt(1500000 * 1e9)
		log.Printf("bribe price is %v", bribe_fee)
		tmp := arg.Contract
		arg.Contract = conf.System_add
		txb, revertTxHashes := cases.GenerateBNBTxs(&arg, bribe_fee, arg.Data, 1)
		arg.Contract = tmp

		txs = append(txb, txs...)
		bundleArgs2 := utils.AddBundle(txs, txs_2, revertTxHashes, 0)

		args[0] = &usr1_arg
		args[1] = &usr2_arg
		bundleArgs_lsit[0] = bundleArgs1
		bundleArgs_lsit[1] = bundleArgs2

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")
		cbn := utils.ConcurSendBundles(t, args, bundleArgs_lsit)

		usrList[0].Txs = txs
		usrList[0].Mined = true
		usrList[0].Rst = conf.Txsucceed

		usrList[1].Txs = txs_1
		usrList[1].Mined = false
		usrList[1].Rst = conf.Txfailed

		usrList[2].Txs = txs_2
		usrList[2].Mined = true
		usrList[2].Rst = conf.Txsucceed

		utils.Verifytx(t, cbn, usrList)
		// 交易顺序符合预期
		// tx1,tx3

	})
	t.Run("bribe_failed", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, reset_data)

		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		lock_data := utils.GeneEncodedData(lockABI, "lock", 1, true)
		txs, _ := utils.SendLockMempool(t, conf.RootPk, conf.Mylock, lock_data, false)
		// tx1 gaslimit 30w, gasprice 1Gwei,
		// tx2 gasLimit 30w, 高于 tx3  >10w,
		// tx3 gasLimit 20w, gasPrice都是1Gwei,
		// tx4 gaslimit 30w, gasprice 1Gwei , 贿赂 SendAmount = 0.00005 * 1Gwei【贿赂失败】
		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.\n")

		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, unlock_more_data)
		usr1_arg.GasLimit = big.NewInt(3e6)

		txs_1, revertTxHashes := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertTxHashes, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, unlock_str_data)
		usr2_arg.GasLimit = big.NewInt(2e6)
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)

		//  Bribe Transaction 【private tx】
		arg := utils.User_tx(conf.RootPk4, conf.SpecialOp, conf.SpecialOp_Bb)
		bribe_fee := big.NewInt(50000 * 1e9)
		log.Printf("bribe price is %v", bribe_fee)
		tmp := arg.Contract
		arg.Contract = conf.System_add
		txb, revertTxHashes := cases.GenerateBNBTxs(&arg, bribe_fee, arg.Data, 1)
		arg.Contract = tmp

		txs0 := append(txb, txs...)
		bundleArgs2 := utils.AddBundle(txs0, txs_2, revertTxHashes, 0)

		args[0] = &usr1_arg
		args[1] = &usr2_arg
		bundleArgs_lsit[0] = bundleArgs1
		bundleArgs_lsit[1] = bundleArgs2

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")
		cbn := utils.ConcurSendBundles(t, args, bundleArgs_lsit)

		usrList[0].Txs = txs
		usrList[0].Mined = true
		usrList[0].Rst = conf.Txsucceed

		usrList[1].Txs = txs_1
		usrList[1].Mined = true
		usrList[1].Rst = conf.Txsucceed

		usrList[2].Txs = txs_2
		usrList[2].Mined = false
		usrList[2].Rst = conf.Txfailed

		utils.Verifytx(t, cbn, usrList)
		// 交易顺序符合预期
		// tx1,tx2

	})

}

func Test_p0_SpecialOp(t *testing.T) {
	// 	t.Run("testCoinbase", func(t *testing.T) {
	// 		t.Log("send bundles testCoinbase \n")
	// 		arg := utils.User_tx(conf.RootPk, conf.SpecialOp, conf.SpecialOp_Cb)
	// 		arg.TxCount = 1
	// 		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, &arg)
	// 		cbn := utils.SendBundlesMined(t, arg, bundleArgs)

	// 		utils.WaitMined(txs, cbn)

	// 		utils.CheckBundleTx(t, *txs[0], true, conf.Txsucceed)

	// 	})

	// 	// 意义不大
	// 	// t.Skip("testTimestamp", func(t *testing.T) {
	// 	// 	t.Log("send bundles testTimestamp \n")
	// 	// 	blkTims := utils.GetLatestBlkMsg(t,conf.Spe_path, "testTimestamp", 5)
	// 	// 	arg := utils.User_tx(conf.RootPk2, conf.SpecialOp, blkTims)
	// 	// 	arg.TxCount = 1
	// 	// 	txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, &arg)
	// 	// 	cbn := utils.SendBundlesMined(t, arg, bundleArgs)

	// 	// 	utils.WaitMined(txs, cbn)

	// 	// 	utils.CheckBundleTx(t, *txs[0], true, conf.Txsucceed)

	// 	// })

	t.Run("testTimestampEq_parent", func(t *testing.T) {
		t.Log("send bundles testTimestamp parent \n")
		blkTims := utils.GetLatestBlkMsg(t, conf.Spe_path, "testTimestampEq", 5)

		arg := utils.User_tx(conf.RootPk, conf.SpecialOp, blkTims)
		arg.TxCount = 1
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, &arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)

		utils.WaitMined(txs, cbn)

		utils.CheckBundleTx(t, *txs[0], true, conf.Txsucceed)

	})

	t.Run("blkHash", func(t *testing.T) {
		t.Log("send bundles blkHash \n")
		blkHash := utils.GetLatestBlkMsg(t, conf.Spe_path, "testBlockHash", 0)
		arg := utils.User_tx(conf.RootPk2, conf.SpecialOp, blkHash)
		arg.TxCount = 1
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, &arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)

		utils.WaitMined(txs, cbn)

		utils.CheckBundleTx(t, *txs[0], true, conf.Txsucceed)

	})

}
