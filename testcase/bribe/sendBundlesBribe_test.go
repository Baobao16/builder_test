package bribe

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/xkwang/testcase"
	"log"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/xkwang/conf"
	"github.com/xkwang/sendBundle"
	"github.com/xkwang/utils"
)

var ValueCpABI = utils.GeneABI(conf.ValueCpPath)
var OwnerABI = utils.GeneABI(conf.OwnerPath)
var valueCpABI = utils.Contract{ABI: *ValueCpABI}
var ownerABI = utils.Contract{ABI: *OwnerABI}

var betT = utils.GeneEncodedData(valueCpABI, "bet", true)
var betF = utils.GeneEncodedData(valueCpABI, "bet", false)

//var changeOwner = utils.GeneEncodedData(ownerABI, "changeOwner", conf.RootPk5)

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
		//{betF, betF, true, false, false}, // revert [tx1] tx2 预期上链
	}
	for index, tc := range testCases {
		t.Run("backRun_value_preservation"+strconv.Itoa(index), func(t *testing.T) {
			rh := make([]common.Hash, 0)
			t.Logf("[Step-1] User 1 bundle [tx1], tx1 Contract bet_%v .\n", tc.aContract)
			bundleArgs1, usr1Arg, txs1 := testcase.AddUserBundle(conf.RootPk, conf.ValueCp, tc.a, bribeFee1, conf.HighGas, txs, rh, 0)

			t.Logf("[Step-2] User 2 bundle [tx2], tx2 Contract bet_%v .\n", tc.bMinted)
			blockNum, _ := usr1Arg.Client.BlockNumber(usr1Arg.Ctx)
			t.Logf("Current Blk_num : %v .\n", blockNum)
			MaxBN := blockNum + 1
			bundleArgs2, usr2Arg, txs2 := testcase.AddUserBundle(conf.RootPk2, conf.ValueCp, tc.b, bribeFee2, conf.HighGas, txs, rh, MaxBN)

			t.Log("[Step-3] User 1 and User 2 send bundles .\n")
			testcase.UpdateUsrList(0, txs1, true, conf.TxSucceed)
			testcase.UpdateUsrList(1, txs2, false, conf.TxFailed)
			testcase.SendBundles(t, usr1Arg, usr2Arg, bundleArgs1, bundleArgs2)

			t.Log("[Step-4] Check tx Minted .\n")
			//blk, _ := strconv.ParseInt(strings.TrimPrefix(response1.Result.BlockNumber, "0x"), 16, 64)
			//t.Logf("Tx1 in Blk : %v .\n", blk)
			//assert.Equal(t, blk, int64(MaxBN))
			if tc.aMinted {
				testcase.CheckTransactionIndex(t, *txs1[0], "0x0")
			}
			if tc.bMinted {
				//assert.Equal(t, response1.Result.BlockNumber, response2.Result.BlockNumber, "tx1 tx2 in diff Block")
				testcase.CheckTransactionIndex(t, *txs2[0], "0x1")
			}

		})
	}
}

func Test_p0_bribe(t *testing.T) {

	t.Run("bribe_success", func(t *testing.T) {
		/*
			tx1 gasLimit 30w, gasPrice 1Gwei,
			tx2 gasLimit 30w, 高于 tx3  >10w,
			tx3 gasLimit 20w, gasPrice都是1Gwei,
			tx4 gasLimit 30w, gasPrice 1Gwei , 贿赂 SendAmount = 0.00015 * 1Gwei【贿赂成功】
		*/
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose Mem_pool transaction tx1\n")
		lockData := utils.GeneEncodedData(testcase.LockABI, "lock", 1, true)
		txs, _ := utils.SendLockMempool(conf.RootPk, conf.Mylock, lockData, conf.MedGas, big.NewInt(conf.MinGasPrice), false, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.\n")
		bundleArgs1, usr1Arg, txs1 := testcase.AddUserBundle(conf.RootPk, conf.ValueCp, testcase.UnlockMoreData, conf.SendA, conf.HighGas, txs, nil, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2Arg := utils.UserTx(conf.RootPk3, conf.Mylock, testcase.UnlockStrData, conf.MedGas, big.NewInt(conf.MinGasPrice))
		txs2, _ := sendBundle.GenerateBNBTxs(&usr2Arg, usr2Arg.SendAmount, usr2Arg.Data, 1)

		//  Bribe Transaction 【private tx】
		arg := utils.UserTx(conf.RootPk4, conf.SpecialOp, conf.SpecialOpBb, conf.HighGas, big.NewInt(conf.MinGasPrice))
		bribeFee := big.NewInt(1500000 * 1e9)
		log.Printf("bribe price is %v", bribeFee)
		tmp := arg.Contract
		arg.Contract = conf.SysAddress
		txb, revertTxHashes := sendBundle.GenerateBNBTxs(&arg, bribeFee, arg.Data, 1)
		arg.Contract = tmp

		txs = append(txb, txs...)
		bundleArgs2 := utils.AddBundle(txs, txs2, revertTxHashes, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")
		cbn := testcase.SendBundles(t, usr1Arg, &usr2Arg, bundleArgs1, bundleArgs2)

		testcase.UpdateUsrList(0, txs, true, conf.TxSucceed)
		testcase.UpdateUsrList(1, txs1, false, conf.TxFailed)
		testcase.UpdateUsrList(2, txs2, true, conf.TxSucceed)
		utils.VerifyTx(t, cbn, testcase.UsrList)
		// 交易顺序符合预期
		// tx1,tx3

	})
	t.Run("bribe_failed", func(t *testing.T) {
		/*
			tx1 gasLimit 30w, gasPrice 1Gwei,
			tx2 gasLimit 30w, 高于 tx3  >10w,
			tx3 gasLimit 20w, gasPrice都是1Gwei,
			tx4 gasLimit 30w, gasPrice 1Gwei , 贿赂 SendAmount = 0.00005 * 1Gwei【贿赂失败】
		*/
		t.Log("[Step-1] Root User Expose mem_pool transaction  tx1\n")
		lockData := utils.GeneEncodedData(testcase.LockABI, "lock", 1, true)
		txs, _ := utils.SendLockMempool(conf.RootPk, conf.Mylock, lockData, conf.MedGas, big.NewInt(conf.MinGasPrice), false, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.\n")
		bundleArgs1, usr1Arg, txs1 := testcase.AddUserBundle(conf.RootPk2, conf.Mylock, testcase.UnlockMoreData, conf.SendA, conf.HighGas, txs, nil, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2Arg := utils.UserTx(conf.RootPk3, conf.Mylock, testcase.UnlockStrData, conf.MedGas, big.NewInt(conf.MinGasPrice))
		txs2, _ := sendBundle.GenerateBNBTxs(&usr2Arg, usr2Arg.SendAmount, usr2Arg.Data, 1)
		//  Bribe Transaction 【private tx】
		arg := utils.UserTx(conf.RootPk4, conf.SpecialOp, conf.SpecialOpBb, conf.HighGas, big.NewInt(conf.MinGasPrice))
		bribeFee := big.NewInt(50000 * 1e9)
		log.Printf("bribe price is %v", bribeFee)
		tmp := arg.Contract
		arg.Contract = conf.SysAddress
		txb, revertTxHashes := sendBundle.GenerateBNBTxs(&arg, bribeFee, arg.Data, 1)
		arg.Contract = tmp

		txs0 := append(txb, txs...)
		bundleArgs2 := utils.AddBundle(txs0, txs2, revertTxHashes, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")
		cbn := testcase.SendBundles(t, usr1Arg, &usr2Arg, bundleArgs1, bundleArgs2)

		testcase.UpdateUsrList(0, txs, true, conf.TxSucceed)
		testcase.UpdateUsrList(1, txs1, true, conf.TxSucceed)
		testcase.UpdateUsrList(2, txs2, false, conf.TxFailed)
		utils.VerifyTx(t, cbn, testcase.UsrList)
		// 交易顺序符合预期
		// tx1,tx2

	})
}

func Test_p0_SpecialOp(t *testing.T) {
	t.Run("testCoinbase", func(t *testing.T) {
		t.Log("send bundles testCoinbase \n")
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)

		arg := utils.UserTx(conf.RootPk, conf.SpecialOp, conf.SpecialOpCb, conf.HighGas, big.NewInt(conf.MinGasPrice))
		arg.TxCount = 1
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)
		utils.WaitMined(txs, cbn)
		utils.CheckBundleTx(t, *txs[0], true, conf.TxSucceed)

	})

	t.Run("testTimestamp", func(t *testing.T) {
		t.Log("send bundles testTimestamp \n")
		blkTS := utils.GetLatestBlkMsg(t, conf.SpePath, "testTimestamp", 5)
		arg := utils.UserTx(conf.RootPk2, conf.SpecialOp, blkTS, conf.HighGas, big.NewInt(conf.MinGasPrice))
		arg.TxCount = 1
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)
		utils.WaitMined(txs, cbn)
		utils.CheckBundleTx(t, *txs[0], true, conf.TxSucceed)
	})

	t.Run("testTimestampEq_parent", func(t *testing.T) {
		t.Log("send bundles testTimestamp parent \n")
		blkTS := utils.GetLatestBlkMsg(t, conf.SpePath, "testTimestampEq", 5)
		arg := utils.UserTx(conf.RootPk3, conf.SpecialOp, blkTS, conf.HighGas, big.NewInt(conf.MinGasPrice))
		arg.TxCount = 1
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)
		utils.WaitMined(txs, cbn)
		utils.CheckBundleTx(t, *txs[0], true, conf.TxSucceed)

	})

	t.Run("blkHash", func(t *testing.T) {
		defer utils.ResetLockContract(t, conf.Mylock, testcase.ResetData)
		t.Log("send bundles blkHash \n")
		blkHash := utils.GetLatestBlkMsg(t, conf.SpePath, "testBlockHash", 0)
		arg := utils.UserTx(conf.RootPk4, conf.SpecialOp, blkHash, conf.HighGas, big.NewInt(conf.MinGasPrice))
		arg.TxCount = 1
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)
		utils.WaitMined(txs, cbn)
		utils.CheckBundleTx(t, *txs[0], true, conf.TxSucceed)

	})

}

func Test_p0_BundleBribe(t *testing.T) {
	// 贿赂地址Balance永远是0
	testCases := []struct {
		bribe1          *big.Int
		bribe2          *big.Int
		add1            common.Address
		add2            common.Address
		txOrder         []string
		balanceIncrease *big.Int
	}{
		{big.NewInt(15 * 1e10), big.NewInt(1e11), conf.SysAddress, conf.BribeAddress, []string{"0x1", "0x0"}, big.NewInt(1e10)},
		{big.NewInt(21 * 1e10), big.NewInt(1e11), conf.SysAddress, conf.BribeAddress, []string{"0x0", "0x1"}, big.NewInt(1e10)},
		{big.NewInt(11 * 1e10), big.NewInt(1e11), conf.BribeAddress, conf.BribeAddress, []string{"0x0", "0x1"}, big.NewInt(21e9)},
	}

	for index, tc := range testCases {
		t.Run("backRun_value_preservation_case"+strconv.Itoa(index), func(t *testing.T) {
			utils.GetAccBalance(conf.BribeAddress)
			Balance1 := utils.GetAccBalance(conf.RcvAddress)
			utils.GetAccBalance(conf.C48Address)
			utils.GetAccBalance(conf.MidAddress)
			var txs types.Transactions
			t.Log("[Step-1]  User1 SendBundle transaction tx1 \n")
			usr1Arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferWBNBCode, conf.HighGas, big.NewInt(conf.MinGasPrice))
			tmp := usr1Arg.Contract
			usr1Arg.Contract = tc.add1
			txs1, revertTxHashes := sendBundle.GenerateBNBTxs(&usr1Arg, tc.bribe1, usr1Arg.Data, 1)
			bundleArgs1 := utils.AddBundle(txs, txs1, revertTxHashes, 0)
			usr1Arg.Contract = tmp

			t.Log("[Step-2]  User2 SendBundle transaction tx2 \n")
			usr2Arg := utils.UserTx(conf.RootPk2, conf.WBNB, conf.TransferWBNBCode, conf.HighGas, big.NewInt(conf.MinGasPrice))
			tmp = usr2Arg.Contract
			usr2Arg.Contract = tc.add2
			txs2, revertTxHashes := sendBundle.GenerateBNBTxs(&usr2Arg, tc.bribe2, usr2Arg.Data, 1)
			bundleArgs2 := utils.AddBundle(txs, txs2, revertTxHashes, 0)
			usr2Arg.Contract = tmp

			t.Log("[Step-3] User 1 and User 2 send bundles .\n")
			//
			testcase.SendBundles(t, &usr1Arg, &usr2Arg, bundleArgs1, bundleArgs2)
			time.Sleep(3 * time.Second)
			testcase.CheckTransactionIndex(t, *txs1[0], tc.txOrder[0])
			testcase.CheckTransactionIndex(t, *txs2[0], tc.txOrder[1])
			utils.GetAccBalance(conf.BribeAddress)
			Balance2 := utils.GetAccBalance(conf.RcvAddress)
			result := new(big.Int)
			result.Sub(Balance2, Balance1)
			assert.Equal(t, result, tc.balanceIncrease)
			utils.GetAccBalance(conf.C48Address)
			utils.GetAccBalance(conf.MidAddress)
		})
	}

}

func Test_p0_BundleLedger(t *testing.T) {
	utils.GetAccBalance(conf.C48Address)
	utils.GetAccBalance(conf.MidAddress)
	Balance1 := utils.GetAccBalance(conf.RcvAddress)
	t.Log("[Step-1] Blk1 Root User Expose Mem_pool transaction  gasFee 100Gwei . \n")
	tx1, _ := utils.SendLockMempool(conf.RootPk, conf.WBNB, conf.TransferWBNBCode, conf.MedGas, big.NewInt(conf.MinGasPrice), false, true)
	utils.BlockHeightIncreased(t)
	utils.CheckBundleTx(t, *tx1[0], true, conf.TxSucceed)

	t.Log("[Step-2] Blk2 bundle with 0.1	ether bribe .\n")
	arg := utils.UserTx(conf.RootPk4, conf.SpecialOp, conf.SpecialOpBb, conf.HighGas, big.NewInt(conf.MinGasPrice))
	bribeFee := big.NewInt(0.01 * 1e18)
	tmp := arg.Contract
	arg.Contract = conf.SysAddress
	txs1 := make([]*types.Transaction, 0)
	txb, revertTxHashes := sendBundle.GenerateBNBTxs(&arg, bribeFee, arg.Data, 1)
	arg.Contract = tmp
	bundleArgs2 := utils.AddBundle(txb, txs1, revertTxHashes, 0)
	utils.SendBundlesMined(t, arg, bundleArgs2)

	Balance2 := utils.GetAccBalance(conf.RcvAddress)
	testcase.CheckTransactionIndex(t, *txb[0], "0x0")
	result := new(big.Int)
	result.Sub(Balance2, Balance1)
	log.Printf("result %v", result)
	utils.GetAccBalance(conf.C48Address)
	utils.GetAccBalance(conf.MidAddress)

}

func Test_P1_ChooseBd(t *testing.T) {

	t.Run("choose_In", func(t *testing.T) {
		/*
			tx1 gasLimit 30w, private transfer,
			tx2 gasLimit 30w, private transfer,
			tx3 gasLimit 30w, mempool transfer,
		*/
		t.Log("[Step-1] Root User Expose mem_pool transaction tx0 tx1")
		tx3, _ := utils.SendLockMempool(conf.RootPk4, conf.WBNB, conf.TransferWBNBCode, conf.LowGas, big.NewInt(conf.MinGasPrice), false, true)

		t.Log("[Step-2] User 1 bundle [tx1].")
		bundleArgs1, usr1Arg, txs1 := testcase.AddUserBundle(conf.RootPk2, conf.WBNB, conf.TransferWBNBCode, conf.SendA, conf.LowGas, nil, nil, 0)

		t.Log("[Step-3] User 2 bundle [tx2].")
		bundleArgs2, usr2Arg, txs2 := testcase.AddUserBundle(conf.RootPk3, conf.WBNB, conf.TransferWBNBCode, conf.SendA, conf.LowGas, nil, nil, 0)

		t.Log("[Step-4] User 3 bundle [tx1, tx2, tx3], tx3 is in mem_pool.")
		txs2 = append(tx3, txs2...)
		bundleArgs3 := utils.AddBundle(txs1, txs2, nil, 0)

		t.Log("[Step-5] User 1 and User 2 and User 3 send bundles.")
		usr3Arg := utils.UserTx(conf.RootPk, conf.Mylock, testcase.UnlockStrData, conf.MedGas, big.NewInt(conf.MinGasPrice))
		cbn := testcase.SendBundlesTri(t, usr1Arg, usr2Arg, &usr3Arg, bundleArgs1, bundleArgs2, bundleArgs3)

		testcase.UpdateUsrList(0, tx3, true, conf.TxSucceed)
		testcase.UpdateUsrList(1, txs1, true, conf.TxSucceed)
		testcase.UpdateUsrList(2, txs2, true, conf.TxSucceed)
		utils.VerifyTx(t, cbn, testcase.UsrList)

	})

}

func Test_P1_comValue(t *testing.T) {
	t.Run("tx2", func(t *testing.T) {
		/*
			tx1 gasLimit 30w, private transfer,
			tx2 gasLimit 30w, private transfer,
			tx3 gasLimit 30w, mempool transfer,
		*/
		//RootPk6 0x199e3Bfb54f4aAa9D67d1BB56429c5ef9D1A2A91 部署
		//pk3 0x6c85F133fa06Fe5eb185743FB6c79f4a7cb9C076
		//pk2 0xb0b10B09780aa6A315158EF724404aa1497e9E6E
		t.Log("[Step-2] User Expose mem_pool transaction tx2")
		tx2, _ := utils.SendLockMempool(conf.RootPk3, conf.Owner, conf.ChangeOwner_other, big.NewInt(1e5), big.NewInt(5e9), false, true)

		t.Log("[Step-3] User Expose mem_pool transaction tx3")
		tx3, _ := utils.SendLockMempool(conf.RootPk2, conf.Owner, conf.ChangeOwner_deployer, big.NewInt(1e6), big.NewInt(1e9), false, true)

		t.Log("[Step-1 User 1 bundle [tx1].")
		t.Log("Root User reset Contract lock")
		usrArg := utils.UserTx(conf.RootPk5, conf.Mylock, testcase.ResetData, conf.LowGas, big.NewInt(conf.MinGasPrice))
		usrArg.TxCount = 1
		tx1, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&usrArg)
		err := usrArg.BuilderClient.SendBundle(usrArg.Ctx, bundleArgs)
		if err != nil {
			return
		}
		t.Log("[Step-3] User 1 and User 2 send transaction .\n")
		time.Sleep(time.Second * 6)
		cbn, _ := usrArg.Client.BlockNumber(usrArg.Ctx)
		testcase.UpdateUsrList(0, tx1, true, conf.TxSucceed)
		testcase.UpdateUsrList(1, tx2, true, conf.TxSucceed)
		testcase.UpdateUsrList(2, tx3, true, conf.TxSucceed)
		utils.VerifyTx(t, cbn, testcase.UsrList)

	})

}
