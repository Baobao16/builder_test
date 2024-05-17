package accuracy

import (
	"github.com/xkwang/testcase"
	"log"
	"testing"

	"github.com/xkwang/cases"
	"github.com/xkwang/conf"
	"github.com/xkwang/utils"
)

func Test_reset(t *testing.T) {
	utils.ResetContract(t, conf.Mylock, testcase.ResetData)
	// utils.SendLockMempool( conf.RootPk, conf.Mylock, lock_data, false)

}

// func Test_con(t *testing.T) {
// 	// defer utils.ResetContract(t, conf.Mylock, reset_data)
// 	// tx_0, _ := utils.SendLockMempool( conf.RootPk, conf.Mylock, lock_data, false)
// 	// time.Sleep(6 * time.Second)
// 	tx_01, _ := utils.SendLockMempool( conf.RootPk2, conf.Mylock, unlock_de_data, false)
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
		defer utils.ResetContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		txs, revertTxHashes := utils.SendLockMempool(conf.RootPk, conf.Mylock, testcase.LockData, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, testcase.UnlockMoreData, conf.High_gas)
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertTxHashes, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, testcase.UnlockStrData, conf.Low_gas)
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
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
		testcase.UsrList[1].Mined = true
		testcase.UsrList[1].Rst = conf.Txsucceed
		testcase.UsrList[2].Txs = txs_2
		testcase.UsrList[2].Mined = false
		testcase.UsrList[2].Rst = conf.Txfailed

		t.Log("[Step-5] Transaction order check .\n")
		utils.Verifytx(t, cbn, testcase.UsrList)
		// Expect [tx1,tx2] 校验链上交易顺序与 bundle-1 交易顺序一致

	})
	t.Run("backrun_tx1_failed", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, testcase.ResetData)
		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		txs, hs := utils.SendLockMempool(conf.RootPk, conf.Mylock, testcase.LockFData, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], none are allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, testcase.UnlockMoreData, conf.High_gas)
		txs_1, r1 := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		hs1 := append(r1, hs...)
		bundleArgs1 := utils.AddBundle(txs, txs_1, hs1, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], none are allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, testcase.UnlockStrData, conf.Low_gas)
		txs_2, r2 := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		hs2 := append(r2, hs...)
		bundleArgs2 := utils.AddBundle(txs, txs_2, hs2, 0)

		testcase.Args[0] = &usr1_arg
		testcase.Args[1] = &usr2_arg
		testcase.BundleArgs_lsit[0] = bundleArgs1
		testcase.BundleArgs_lsit[1] = bundleArgs2

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")
		cbn := utils.ConcurSendBundles(t, testcase.Args, testcase.BundleArgs_lsit)

		testcase.UsrList[0].Txs = txs
		testcase.UsrList[0].Mined = true
		testcase.UsrList[0].Rst = conf.Txfailed
		testcase.UsrList[1].Txs = txs_1
		testcase.UsrList[1].Mined = false
		testcase.UsrList[1].Rst = conf.Txfailed
		testcase.UsrList[2].Txs = txs_2
		testcase.UsrList[2].Mined = false
		testcase.UsrList[2].Rst = conf.Txfailed
		t.Log("[Step-5] Verify transaction .\n")
		utils.Verifytx(t, cbn, testcase.UsrList)
	})

}

func Test_p0_token_sniper(t *testing.T) {
	t.Run("tokenSniper_tx1_success", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		txs, revertHash := utils.SendLockMempool(conf.RootPk, conf.Mylock, testcase.LockData, true)

		t.Log("[Step-2] User 1 bundle [tx1, tx2], none are allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, testcase.FakelockMoreData, conf.High_gas)
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], none are allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, testcase.FakelockStrData, conf.Low_gas)
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(txs, txs_2, revertHash, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")

		testcase.Args[0] = &usr1_arg
		testcase.Args[1] = &usr2_arg
		testcase.BundleArgs_lsit[0] = bundleArgs1
		testcase.BundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, testcase.Args, testcase.BundleArgs_lsit)

		testcase.UsrList[0].Txs = txs
		testcase.UsrList[0].Mined = true
		testcase.UsrList[0].Rst = conf.Txsucceed
		testcase.UsrList[1].Txs = txs_1
		testcase.UsrList[1].Mined = true
		testcase.UsrList[1].Rst = conf.Txsucceed
		testcase.UsrList[2].Txs = txs_2
		testcase.UsrList[2].Mined = true
		testcase.UsrList[2].Rst = conf.Txsucceed

		t.Log("[Step-5] Verify transaction .\n")
		utils.Verifytx(t, cbn, testcase.UsrList)

		// 目标交易顺序 [tx1, tx2, tx3]
		target_txl := append(txs, txs_1...)
		target_txl = append(target_txl, txs_2...)
		for _, tx := range target_txl {
			log.Println(tx.Hash().Hex())
		}

	})
	t.Run("tokenSniper_tx1_failed", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		txs, revertHash := utils.SendLockMempool(conf.RootPk, conf.Mylock, testcase.LockFData, true) //conf.Lock56_lock0t "gasUsed":"0x342b" 13355

		t.Log("[Step-2] User 1 bundle [tx1, tx2], none are allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, testcase.FakelockMoreData, conf.High_gas) // unlock_long
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], none are allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, testcase.FakelockStrData, conf.Low_gas) // unlock_str
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(txs, txs_2, revertHash, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")

		testcase.Args[0] = &usr1_arg
		testcase.Args[1] = &usr2_arg
		testcase.BundleArgs_lsit[0] = bundleArgs1
		testcase.BundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, testcase.Args, testcase.BundleArgs_lsit)

		testcase.UsrList[0].Txs = txs
		testcase.UsrList[0].Mined = true
		testcase.UsrList[0].Rst = conf.Txfailed
		testcase.UsrList[1].Txs = txs_1
		testcase.UsrList[1].Mined = false
		testcase.UsrList[1].Rst = conf.Txfailed
		testcase.UsrList[2].Txs = txs_2
		testcase.UsrList[2].Mined = false
		testcase.UsrList[2].Rst = conf.Txfailed

		t.Log("[Step-5] Verify transaction .\n")
		utils.Verifytx(t, cbn, testcase.UsrList)

	})

}

func Test_p0_running_attack(t *testing.T) {
	t.Run("runningAttack_tx1_success", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mempool transaction tx0  tx1 \n") // 会上链

		tx_0, rh := utils.SendLockMempool(conf.RootPk4, conf.WBNB, conf.TransferWBNB_code, true)
		tx_01, revertHash := utils.SendLockMempool(conf.RootPk, conf.Mylock, testcase.LockData, true)
		tx_02 := append(tx_01, tx_0...)
		revertHash = append(revertHash, rh[0])

		t.Log("[Step-2] User 1 bundle [tx0, tx1, tx2], tx2 not allowed to revert.\n")              // 不会上链
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, testcase.UnlockStrData, conf.Low_gas) //   [tx2] "gasUsed":"0xdd65"
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(tx_02, txs_1, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")                    // 会上链
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, testcase.UnlockMoreData, conf.Med_gas) //  [tx3] "gasUsed":"0xe751"
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(tx_01, txs_2, revertHash, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")

		testcase.Args[0] = &usr1_arg
		testcase.Args[1] = &usr2_arg
		testcase.BundleArgs_lsit[0] = bundleArgs1
		testcase.BundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, testcase.Args, testcase.BundleArgs_lsit)
		// 在tx1成功执行的前提下，链上交易为[tx0,tx1,tx3]
		t.Log("[Step-5] check tx0 tx1 mined .\n")
		testcase.UsrList[0].Txs = tx_02
		testcase.UsrList[0].Mined = true
		testcase.UsrList[0].Rst = conf.Txsucceed
		testcase.UsrList[1].Txs = txs_1
		testcase.UsrList[1].Mined = false
		testcase.UsrList[1].Rst = conf.Txfailed
		testcase.UsrList[2].Txs = txs_2
		testcase.UsrList[2].Mined = true
		testcase.UsrList[2].Rst = conf.Txsucceed
		utils.Verifytx(t, cbn, testcase.UsrList)

		// 目标交易顺序 [ tx1, tx0, tx3]
		// target_txl := append(tx_02, txs_2...)
		// for _, tx := range target_txl {
		// 	log.Println(tx.Hash().Hex())
		// }

	})
	t.Run("runningAttack_tx1_failed", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mempool transaction tx0  tx1 \n")                       // 会上链
		tx_0, _ := utils.SendLockMempool(conf.RootPk4, conf.WBNB, conf.TransferWBNB_code, false) // [tx0]"gasUsed": "0x6323" :25379
		tx_01, _ := utils.SendLockMempool(conf.RootPk, conf.Mylock, testcase.LockFData, false)   // [tx1]conf.Lock56_lock0t "gasUsed":"0x342b" 13355
		tx_02 := append(tx_01, tx_0...)

		t.Log("[Step-2] User 1 bundle [tx0, tx1, tx2], tx2 not allowed to revert.\n")              // 不会上链
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, testcase.UnlockStrData, conf.Low_gas) // [tx2] unlock_str  "gasUsed":"0xbfd8" 49112
		txs_1, revertTxHashes := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(tx_02, txs_1, revertTxHashes, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")                     // 会上链
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, testcase.UnlockMoreData, conf.High_gas) // [tx3] unlock_long "gasUsed":"0xc944" 51524
		txs_2, revertTxHashes := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(tx_01, txs_2, revertTxHashes, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")
		testcase.Args[0] = &usr1_arg
		testcase.Args[1] = &usr2_arg
		testcase.BundleArgs_lsit[0] = bundleArgs1
		testcase.BundleArgs_lsit[1] = bundleArgs2
		utils.ConcurSendBundles(t, testcase.Args, testcase.BundleArgs_lsit)

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
		defer utils.ResetContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		txs, _ := utils.SendLockMempool(conf.RootPk, conf.Mylock, testcase.LockData, false) //[tx1]conf.Lock56_lock0t "gasUsed":"0x342b" 13355

		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, testcase.UnlockDeData, conf.Med_gas) // [tx2] unlock_long "gasUsed":"0xc944" 51524
		// tx2 会上链
		// gasfee(gasused * gasprice) 20w * 1e10
		// GasLimit * gasprice        20w * 1e10
		txs_1, revertTxHashes := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertTxHashes, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, testcase.UnlockDesData, conf.High_gas) // [tx3] unlock_str  "gasUsed":"0xe8b1"
		// tx3 不会上链
		// gasfee(gasused * gasprice) 3.5w * 1e10
		// GasLimit * gasprice        30w  * 1e10
		txs_2, revertTxHashes := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(txs, txs_2, revertTxHashes, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")

		testcase.Args[0] = &usr1_arg
		testcase.Args[1] = &usr2_arg
		testcase.BundleArgs_lsit[0] = bundleArgs1
		testcase.BundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, testcase.Args, testcase.BundleArgs_lsit)

		t.Log("[Step-5] Verify transaction .\n")
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

		// t.Log("[Step-5] Transaction order check .\n")
		// Expect [tx1,tx2] 校验链上交易顺序与 bundle-1 交易顺序一致

	})
	t.Run("gaslimitDeception_tx1_failed", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mempool transaction tx1\n")
		txs, revertHash := utils.SendLockMempool(conf.RootPk, conf.Mylock, testcase.LockFData, true) //conf.Lock56_lock0t "gasUsed":"0x342b" 13355

		t.Log("[Step-2] User 1 bundle [tx1, tx2], tx2 not allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, testcase.UnlockDeData, conf.Med_gas) // [tx2] unlock_long "gasUsed":"0xc944" 51524
		// tx2 会上链
		// gasfee(gasused * gasprice) 20w * 1e10
		// GasLimit * gasprice        20w * 1e10
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(txs, txs_1, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, testcase.UnlockDesData, conf.High_gas) // [tx3] unlock_str  "gasUsed":"0xe8b1"
		// tx3 不会上链
		// gasfee(gasused * gasprice) 3.5w * 1e10
		// GasLimit * gasprice        30w  * 1e10
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(txs, txs_2, revertHash, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles .\n")

		testcase.Args[0] = &usr1_arg
		testcase.Args[1] = &usr2_arg
		testcase.BundleArgs_lsit[0] = bundleArgs1
		testcase.BundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, testcase.Args, testcase.BundleArgs_lsit)

		t.Log("[Step-5] Verify transaction .\n")
		testcase.UsrList[0].Txs = txs
		testcase.UsrList[0].Mined = true
		testcase.UsrList[0].Rst = conf.Txfailed
		testcase.UsrList[1].Txs = txs_1
		testcase.UsrList[1].Mined = false
		testcase.UsrList[1].Rst = conf.Txfailed
		testcase.UsrList[2].Txs = txs_2
		testcase.UsrList[2].Mined = false
		testcase.UsrList[2].Rst = conf.Txfailed
		utils.Verifytx(t, cbn, testcase.UsrList)

	})

}

func Test_p0_sandwich(t *testing.T) {

	t.Run("sandwich_ol1", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mempool transaction tx0  tx1 \n")
		tx_0, revertHash_0 := utils.SendLockMempool(conf.RootPk2, conf.Mylock, testcase.LockData, true)
		tx_1, revertHash_1 := utils.SendLockMempool(conf.RootPk, conf.Mylock, testcase.FakelockStrData, true)

		tx_01 := append(tx_0, tx_1...)
		revertHash := append(revertHash_0, revertHash_1...)

		t.Log("[Step-2] User 1 bundle [tx0, tx1, tx2], tx2 not allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, testcase.UnlockMoreData, conf.High_gas)
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(tx_01, txs_1, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3], tx3 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, testcase.UnlockStrData, conf.Low_gas)
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(tx_1, txs_2, revertHash_1, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles . Expect :mined Tx_list : [tx0, tx1, tx2]\n")

		testcase.Args[0] = &usr1_arg
		testcase.Args[1] = &usr2_arg
		testcase.BundleArgs_lsit[0] = bundleArgs1
		testcase.BundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, testcase.Args, testcase.BundleArgs_lsit)

		// 在tx1成功执行的前提下，链上交易为[tx0,tx1,tx2]
		t.Log("[Step-5] check tx0 tx1 mined .\n")
		testcase.UsrList[0].Txs = tx_01
		testcase.UsrList[0].Mined = true
		testcase.UsrList[0].Rst = conf.Txsucceed
		testcase.UsrList[1].Txs = txs_1
		testcase.UsrList[1].Mined = true
		testcase.UsrList[1].Rst = conf.Txsucceed
		testcase.UsrList[2].Txs = txs_2
		testcase.UsrList[2].Mined = false
		testcase.UsrList[2].Rst = conf.Txfailed
		utils.Verifytx(t, cbn, testcase.UsrList)
		// tx0,tx1,tx2,tx3
	})

	t.Run("sandwich_both", func(t *testing.T) {
		defer utils.ResetContract(t, conf.Mylock, testcase.ResetData)

		t.Log("[Step-1] Root User Expose mempool transaction tx0  tx1 \n")
		tx_0, revertHash_0 := utils.SendLockMempool(conf.RootPk2, conf.Mylock, testcase.LockData, false)
		tx_1, revertHash_1 := utils.SendLockMempool(conf.RootPk, conf.Mylock, testcase.UnlockStrData, true)

		tx_01 := append(tx_0, tx_1...)
		revertHash := append(revertHash_0, revertHash_1...)

		t.Log("[Step-2] User 1 bundle [tx0, tx1, tx2], all not allowed to revert.\n")
		usr1_arg := utils.User_tx(conf.RootPk2, conf.Mylock, testcase.ResetData, conf.High_gas)
		txs_1, _ := cases.GenerateBNBTxs(&usr1_arg, usr1_arg.SendAmount, usr1_arg.Data, 1)
		bundleArgs1 := utils.AddBundle(tx_01, txs_1, revertHash, 0)

		t.Log("[Step-3] User 2 bundle [tx1, tx3],      tx1 not allowed to revert.\n")
		usr2_arg := utils.User_tx(conf.RootPk3, conf.Mylock, testcase.ResetData, conf.Low_gas)
		txs_2, _ := cases.GenerateBNBTxs(&usr2_arg, usr2_arg.SendAmount, usr2_arg.Data, 1)
		bundleArgs2 := utils.AddBundle(tx_1, txs_2, revertHash_1, 0)

		t.Log("[Step-4] User 1 and User 2 send bundles . Expect :mined Tx_list : [tx0, tx1, tx2, tx3]\n")

		testcase.Args[0] = &usr1_arg
		testcase.Args[1] = &usr2_arg
		testcase.BundleArgs_lsit[0] = bundleArgs1
		testcase.BundleArgs_lsit[1] = bundleArgs2
		cbn := utils.ConcurSendBundles(t, testcase.Args, testcase.BundleArgs_lsit)

		// 在tx1成功执行的前提下，链上交易为[tx0,tx1,tx2]
		t.Log("[Step-5] check tx0, tx1, tx2, tx3 mined .\n")
		testcase.UsrList[0].Txs = tx_01
		testcase.UsrList[0].Mined = true
		testcase.UsrList[0].Rst = conf.Txsucceed
		testcase.UsrList[1].Txs = txs_1
		testcase.UsrList[1].Mined = true
		testcase.UsrList[1].Rst = conf.Txsucceed
		testcase.UsrList[2].Txs = txs_2
		testcase.UsrList[2].Mined = true
		testcase.UsrList[2].Rst = conf.Txsucceed
		utils.Verifytx(t, cbn, testcase.UsrList)
		// tx0,tx1,tx2,tx3
	})

}
