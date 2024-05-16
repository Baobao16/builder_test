package newtestcases

import (
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xkwang/cases"
	"github.com/xkwang/conf"
	"github.com/xkwang/utils"
)

var usrList = make([]utils.TxStatus, 3)
var args = make([]*cases.BidCaseArg, 2)
var bundleArgs_lsit = make([]*types.SendBundleArgs, 2)
var LockABI = utils.GeneABI(conf.Lock_path)

// step-1 读取abi文件
var lockABI = utils.Contract{ABI: *LockABI}
var reset_data = utils.GeneEncodedData(lockABI, "reset")
var lock_data = utils.GeneEncodedData(lockABI, "lock", 1, true) // "gasUsed":"0x6ab1"
var lock_f_data = utils.GeneEncodedData(lockABI, "lock", 1, false)

var unlock_str_data = utils.GeneEncodedData(lockABI, "unlock", 1, "str")
var unlock_more_data = utils.GeneEncodedData(lockABI, "unlock", 1, "more")
var unlock_de_data = utils.GeneEncodedData(lockABI, "unlock_de", 1, "more") //"0x1d58c7"1923271
var unlock_des_data = utils.GeneEncodedData(lockABI, "unlock_de", 1, "str")
var fakelock_str_data = utils.GeneEncodedData(lockABI, "fakelock", 1, "str")
var fakelock_more_data = utils.GeneEncodedData(lockABI, "fakelock", 1, "more")

func Test_reset(t *testing.T) {
	utils.ResetContract(t, conf.Mylock, reset_data)
	// utils.SendLockMempool(t, conf.RootPk, conf.Mylock, lock_data, false)

}

// func Test_con(t *testing.T) {
// 	// defer utils.ResetContract(t, conf.Mylock, reset_data)
// 	// tx_0, _ := utils.SendLockMempool(t, conf.RootPk, conf.Mylock, lock_data, false)
// 	// time.Sleep(6 * time.Second)
// 	tx_01, _ := utils.SendLockMempool(t, conf.RootPk2, conf.Mylock, unlock_de_data, false)
// 	// print(tx_0[0])
// 	print(tx_01[0])
// 	time.Sleep(6 * time.Second)
// 	rcp := utils.GetTransactionReceipt(*tx_01[0])
// 	println("rcp.Result.GasUsed")
// 	if rcp.Result.Status == "0x1" {
// 		println(rcp.Result.GasUsed)
// 	} else {
// 		println("Continue")
// 	}

// }

func Test_p0_backrun(t *testing.T) {
	t.Run("backrun_tx1_success", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, reset_data)

		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		txs, revertTxHashes := utils.SendLockMempool(t, conf.RootPk, conf.Mylock, lock_data, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, unlock_more_data, conf.High_gas)
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertTxHashes, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, unlock_str_data, conf.Low_gas)
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
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
		usrList[1].Mined = true
		usrList[1].Rst = conf.Txsucceed
		usrList[2].Txs = txs_2
		usrList[2].Mined = false
		usrList[2].Rst = conf.Txfailed

		t.Log("[Step-5] Transaction order check .\n")
		utils.Verifytx(t, cbn, usrList)
		// Expect [tx1,tx2] 校验链上交易顺序与 bundle-1 交易顺序一致

	})
	t.Run("backrun_tx1_failed", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, reset_data)
		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		txs, hs := utils.SendLockMempool(t, conf.RootPk, conf.Mylock, lock_f_data, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], none are allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, unlock_more_data, conf.High_gas)
		txs_1, r1 := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		hs1 := append(r1, hs...)
		bundleArgs1 := utils.AddBundle(txs, txs_1, hs1, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], none are allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, unlock_str_data, conf.Low_gas)
		txs_2, r2 := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		hs2 := append(r2, hs...)
		bundleArgs2 := utils.AddBundle(txs, txs_2, hs2, 0)

		args[0] = &usr1_arg
		args[1] = &usr2_arg
		bundleArgs_lsit[0] = bundleArgs1
		bundleArgs_lsit[1] = bundleArgs2

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")
		cbn := utils.ConcurSendBundles(t, args, bundleArgs_lsit)

		usrList[0].Txs = txs
		usrList[0].Mined = true
		usrList[0].Rst = conf.Txfailed
		usrList[1].Txs = txs_1
		usrList[1].Mined = false
		usrList[1].Rst = conf.Txfailed
		usrList[2].Txs = txs_2
		usrList[2].Mined = false
		usrList[2].Rst = conf.Txfailed
		t.Log("[Step-5] Verify transaction .\n")
		utils.Verifytx(t, cbn, usrList)
	})

}

func Test_p0_token_sniper(t *testing.T) {
	t.Run("tokenSniper_tx1_success", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, reset_data)

		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		txs, revertHash := utils.SendLockMempool(t, conf.RootPk, conf.Mylock, lock_data, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], none are allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, fakelock_more_data, conf.High_gas)
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], none are allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, fakelock_str_data, conf.Low_gas)
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(txs, txs_2, revertHash, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")

		args[0] = &usr1_arg
		args[1] = &usr2_arg
		bundleArgs_lsit[0] = bundleArgs1
		bundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, args, bundleArgs_lsit)

		usrList[0].Txs = txs
		usrList[0].Mined = true
		usrList[0].Rst = conf.Txsucceed
		usrList[1].Txs = txs_1
		usrList[1].Mined = true
		usrList[1].Rst = conf.Txsucceed
		usrList[2].Txs = txs_2
		usrList[2].Mined = true
		usrList[2].Rst = conf.Txsucceed

		t.Log("[Step-5] Verify transaction .\n")
		utils.Verifytx(t, cbn, usrList)

		// 目标交易顺序 [tx1, tx2, tx3]
		target_txl := append(txs, txs_1...)
		target_txl = append(target_txl, txs_2...)
		for _, tx := range target_txl {
			log.Println(tx.Hash().Hex())
		}

	})
	t.Run("tokenSniper_tx1_failed", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, reset_data)

		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		txs, revertHash := utils.SendLockMempool(t, conf.RootPk, conf.Mylock, lock_f_data, true) //conf.Lock56_lock0t "gasUsed":"0x342b" 13355

		t.Log("[Step-2] User 1 bundle [tx1, tx2], none are allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, fakelock_more_data, conf.High_gas) // unlock_long
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], none are allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, fakelock_str_data, conf.Low_gas) // unlock_str
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(txs, txs_2, revertHash, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")

		args[0] = &usr1_arg
		args[1] = &usr2_arg
		bundleArgs_lsit[0] = bundleArgs1
		bundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, args, bundleArgs_lsit)

		usrList[0].Txs = txs
		usrList[0].Mined = true
		usrList[0].Rst = conf.Txfailed
		usrList[1].Txs = txs_1
		usrList[1].Mined = false
		usrList[1].Rst = conf.Txfailed
		usrList[2].Txs = txs_2
		usrList[2].Mined = false
		usrList[2].Rst = conf.Txfailed

		t.Log("[Step-5] Verify transaction .\n")
		utils.Verifytx(t, cbn, usrList)

	})

}

func Test_p0_running_attack(t *testing.T) {
	t.Run("runningAttack_tx1_success", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, reset_data)

		t.Log("[Step-1] Root User Expose mempool transaction tx0  tx1 \n") // 会上链

		tx_0, rh := utils.SendLockMempool(t, conf.RootPk4, conf.WBNB, conf.TransferWBNB_code, true)
		tx_01, revertHash := utils.SendLockMempool(t, conf.RootPk, conf.Mylock, lock_data, true)
		tx_02 := append(tx_01, tx_0...)
		revertHash = append(revertHash, rh[0])

		t.Log("[Step-2] User 1 bundle [tx0, tx1, tx2], tx2 not allowed to revert.\n")       // 不会上链
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, unlock_str_data, conf.Low_gas) //   [tx2] "gasUsed":"0xdd65"
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(tx_02, txs_1, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")             // 会上链
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, unlock_more_data, conf.Med_gas) //  [tx3] "gasUsed":"0xe751"
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(tx_01, txs_2, revertHash, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")

		args[0] = &usr1_arg
		args[1] = &usr2_arg
		bundleArgs_lsit[0] = bundleArgs1
		bundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, args, bundleArgs_lsit)
		// 在tx1成功执行的前提下，链上交易为[tx0,tx1,tx3]
		t.Log("[Step-5] check tx0 tx1 mined .\n")
		usrList[0].Txs = tx_02
		usrList[0].Mined = true
		usrList[0].Rst = conf.Txsucceed
		usrList[1].Txs = txs_1
		usrList[1].Mined = false
		usrList[1].Rst = conf.Txfailed
		usrList[2].Txs = txs_2
		usrList[2].Mined = true
		usrList[2].Rst = conf.Txsucceed
		utils.Verifytx(t, cbn, usrList)

		// 目标交易顺序 [ tx1, tx0, tx3]
		// target_txl := append(tx_02, txs_2...)
		// for _, tx := range target_txl {
		// 	log.Println(tx.Hash().Hex())
		// }

	})
	t.Run("runningAttack_tx1_failed", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, reset_data)

		t.Log("[Step-1] Root User Expose mempool transaction tx0  tx1 \n")                          // 会上链
		tx_0, _ := utils.SendLockMempool(t, conf.RootPk4, conf.WBNB, conf.TransferWBNB_code, false) // [tx0]"gasUsed": "0x6323" :25379
		tx_01, _ := utils.SendLockMempool(t, conf.RootPk, conf.Mylock, lock_f_data, false)          // [tx1]conf.Lock56_lock0t "gasUsed":"0x342b" 13355
		tx_02 := append(tx_01, tx_0...)

		t.Log("[Step-2] User 1 bundle [tx0, tx1, tx2], tx2 not allowed to revert.\n")       // 不会上链
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, unlock_str_data, conf.Low_gas) // [tx2] unlock_str  "gasUsed":"0xbfd8" 49112
		txs_1, revertTxHashes := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(tx_02, txs_1, revertTxHashes, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")              // 会上链
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, unlock_more_data, conf.High_gas) // [tx3] unlock_long "gasUsed":"0xc944" 51524
		txs_2, revertTxHashes := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(tx_01, txs_2, revertTxHashes, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")
		args[0] = &usr1_arg
		args[1] = &usr2_arg
		bundleArgs_lsit[0] = bundleArgs1
		bundleArgs_lsit[1] = bundleArgs2
		utils.ConcurSendBundles(t, args, bundleArgs_lsit)

		t.Log("[Step-5] Verify transaction .\n")
		// for _, tx := range tx_02 {
		// 	utils.CheckBundleTx(t, *tx, true, conf.Txsucceed)
		// }
		// tx0 success
		// tx1 failed

		for _, tx := range txs_1 {
			// 依次检查bundle中的交易是否成功上链
			log.Println("bundle 1 not mined")
			utils.CheckBundleTx(t, *tx, false, conf.Txfailed)
		}

		for _, tx := range txs_2 {
			log.Println("bundle 2 mined")
			// 依次检查bundle中的交易是否成功上链
			utils.CheckBundleTx(t, *tx, false, conf.Txfailed)
		}

	})

}

func Test_p0_gaslimit_deception(t *testing.T) {
	t.Run("gaslimitDeception_tx1_success", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, reset_data)

		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		txs, _ := utils.SendLockMempool(t, conf.RootPk, conf.Mylock, lock_data, false) //[tx1]conf.Lock56_lock0t "gasUsed":"0x342b" 13355

		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, unlock_de_data, conf.Med_gas) // [tx2] unlock_long "gasUsed":"0xc944" 51524
		// tx2 会上链
		// gasfee(gasused * gasprice) 20w * 1e10
		// GasLimit * gasprice        20w * 1e10
		txs_1, revertTxHashes := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertTxHashes, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, unlock_des_data, conf.High_gas) // [tx3] unlock_str  "gasUsed":"0xe8b1"
		// tx3 不会上链
		// gasfee(gasused * gasprice) 3.5w * 1e10
		// GasLimit * gasprice        30w  * 1e10
		txs_2, revertTxHashes := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(txs, txs_2, revertTxHashes, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")

		args[0] = &usr1_arg
		args[1] = &usr2_arg
		bundleArgs_lsit[0] = bundleArgs1
		bundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, args, bundleArgs_lsit)

		t.Log("[Step-5] Verify transaction .\n")
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

		// t.Log("[Step-5] Transaction order check .\n")
		// Expect [tx1,tx2] 校验链上交易顺序与 bundle-1 交易顺序一致

	})
	t.Run("gaslimitDeception_tx1_failed", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, reset_data)

		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		txs, revertHash := utils.SendLockMempool(t, conf.RootPk, conf.Mylock, lock_f_data, true) //conf.Lock56_lock0t "gasUsed":"0x342b" 13355

		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, unlock_de_data, conf.Med_gas) // [tx2] unlock_long "gasUsed":"0xc944" 51524
		// tx2 会上链
		// gasfee(gasused * gasprice) 20w * 1e10
		// GasLimit * gasprice        20w * 1e10
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, unlock_des_data, conf.High_gas) // [tx3] unlock_str  "gasUsed":"0xe8b1"
		// tx3 不会上链
		// gasfee(gasused * gasprice) 3.5w * 1e10
		// GasLimit * gasprice        30w  * 1e10
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(txs, txs_2, revertHash, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")

		args[0] = &usr1_arg
		args[1] = &usr2_arg
		bundleArgs_lsit[0] = bundleArgs1
		bundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, args, bundleArgs_lsit)

		t.Log("[Step-5] Verify transaction .\n")
		usrList[0].Txs = txs
		usrList[0].Mined = true
		usrList[0].Rst = conf.Txfailed
		usrList[1].Txs = txs_1
		usrList[1].Mined = false
		usrList[1].Rst = conf.Txfailed
		usrList[2].Txs = txs_2
		usrList[2].Mined = false
		usrList[2].Rst = conf.Txfailed
		utils.Verifytx(t, cbn, usrList)

	})

}

func Test_p0_sandwich(t *testing.T) {

	t.Run("sandwich_ol1", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, reset_data)

		t.Log("[Step-1] Root User Expose mempool transaction tx0  tx1 \n")
		tx_0, revertHash_0 := utils.SendLockMempool(t, conf.RootPk2, conf.Mylock, lock_data, true)
		tx_1, revertHash_1 := utils.SendLockMempool(t, conf.RootPk, conf.Mylock, fakelock_str_data, true)

		tx_01 := append(tx_0, tx_1...)
		revertHash := append(revertHash_0, revertHash_1...)

		t.Log("[Step-2] User 1 bundle [tx0, tx1, tx2], tx2 not allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, unlock_more_data, conf.High_gas)
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(tx_01, txs_1, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, unlock_str_data, conf.Low_gas)
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(tx_1, txs_2, revertHash_1, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles . Expect :mined Tx_list : [tx0, tx1, tx2]\n")

		args[0] = &usr1_arg
		args[1] = &usr2_arg
		bundleArgs_lsit[0] = bundleArgs1
		bundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, args, bundleArgs_lsit)

		// 在tx1成功执行的前提下，链上交易为[tx0,tx1,tx2]
		t.Log("[Step-5] check tx0 tx1 mined .\n")
		usrList[0].Txs = tx_01
		usrList[0].Mined = true
		usrList[0].Rst = conf.Txsucceed
		usrList[1].Txs = txs_1
		usrList[1].Mined = true
		usrList[1].Rst = conf.Txsucceed
		usrList[2].Txs = txs_2
		usrList[2].Mined = false
		usrList[2].Rst = conf.Txfailed
		utils.Verifytx(t, cbn, usrList)
		// tx0,tx1,tx2,tx3
	})

	t.Run("sandwich_both", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, reset_data)

		t.Log("[Step-1] Root User Expose mempool transaction tx0  tx1 \n")
		tx_0, revertHash_0 := utils.SendLockMempool(t, conf.RootPk2, conf.Mylock, lock_data, false)
		tx_1, revertHash_1 := utils.SendLockMempool(t, conf.RootPk, conf.Mylock, unlock_str_data, true)

		tx_01 := append(tx_0, tx_1...)
		revertHash := append(revertHash_0, revertHash_1...)

		t.Log("[Step-2] User 1 bundle [tx0, tx1, tx2], all not allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, reset_data, conf.High_gas)
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(tx_01, txs_1, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3],      tx1 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, reset_data, conf.Low_gas)
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(tx_1, txs_2, revertHash_1, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles . Expect :mined Tx_list : [tx0, tx1, tx2, tx3]\n")

		args[0] = &usr1_arg
		args[1] = &usr2_arg
		bundleArgs_lsit[0] = bundleArgs1
		bundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, args, bundleArgs_lsit)

		// 在tx1成功执行的前提下，链上交易为[tx0,tx1,tx2]
		t.Log("[Step-5] check tx0, tx1, tx2, tx3 mined .\n")
		usrList[0].Txs = tx_01
		usrList[0].Mined = true
		usrList[0].Rst = conf.Txsucceed
		usrList[1].Txs = txs_1
		usrList[1].Mined = true
		usrList[1].Rst = conf.Txsucceed
		usrList[2].Txs = txs_2
		usrList[2].Mined = true
		usrList[2].Rst = conf.Txsucceed
		utils.Verifytx(t, cbn, usrList)
		// tx0,tx1,tx2,tx3
	})

}
