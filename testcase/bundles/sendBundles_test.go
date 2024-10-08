package bundles

import (
	"fmt"
	"github.com/xkwang/testcase"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/xkwang/conf"
	"github.com/xkwang/sendBundle"
	"github.com/xkwang/utils"
)

var txBlk = make([]string, 0)
var timeLimit = 300

/*
sendbundle 接口测试
*/

// 正常sendBundle
func Test_p0_sendBundle(t *testing.T) {
	arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferwbnbCode, conf.HighGas)
	t.Run("sendValidBundle_tx", func(t *testing.T) {
		// bundle 中均为合法转账交易
		t.Log("Start sendBundle \n")
		//conf.Mylock, testcase.UnlockDeMoreData, conf.SendA, conf.Low_gas,
		arg.Contract = conf.Mylock
		//arg.Data = testcase.LockData
		arg.Data = testcase.UnlockDeMoreData
		arg.TxCount = 1
		arg.GasLimit = conf.LowGas
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)
		utils.WaitMined(txs, cbn)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			blkN := utils.CheckBundleTx(t, *tx, true, conf.TxSucceed)
			txBlk = append(txBlk, blkN)
		}
		if !utils.TxInSameBlk(txBlk) {
			t.Fatalf("Transactions are not in the same block")
		}
	})
}

/*
sendBundle包含 revert交易
1.Bundle_all_revert  - in revertList
2.Bundle_part_revert - in revertList
3.Bundle_part_revert - not in revertList
4.Bundle_all_revert  - not in revertList
5.Bundle_no_revert   - in revertList
*/
func Test_p0_sendBundle_revert(t *testing.T) {
	arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferwbnbCode, conf.HighGas)
	t.Run("sendValidBundle_all_revert", func(t *testing.T) {
		// revert 交易均在revertList中记录
		t.Log("generate revert transaction \n")
		// 设置参数
		arg.TxCount = 3
		arg.RevertList = []int{0, 1, 2}

		// 发送并验证交易
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)
		utils.WaitMined(txs, cbn)
		for _, tx := range txs {
			blkN := utils.CheckBundleTx(t, *tx, true, conf.TxFailed)
			txBlk = append(txBlk, blkN)
		}
		if !utils.TxInSameBlk(txBlk) {
			t.Fatalf("Transactions are not in the same block")
		}
	})

	t.Run("sendValidBundle_part_revert", func(t *testing.T) {
		t.Log("generate revert transaction \n")

		// 设置参数
		arg.TxCount = 10
		arg.RevertList = []int{0, 1, 2}

		// 发送并验证交易
		txs, bundleArgs, err := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		if err != nil {
			t.Fatalf("Failed to generate valid bundle: %v", err)
		}

		cbn := utils.SendBundlesMined(t, arg, bundleArgs)
		utils.WaitMined(txs, cbn)

		var txBlk []string
		for index, tx := range txs {
			// 检查交易是否成功上链
			expectedStatus := conf.TxSucceed
			if index < 3 {
				expectedStatus = conf.TxFailed
			}
			blkN := utils.CheckBundleTx(t, *tx, true, expectedStatus)
			txBlk = append(txBlk, blkN)
		}

		if !utils.TxInSameBlk(txBlk) {
			t.Fatalf("Transactions are not in the same block")
		}
	})
}

func Test_p1_sendBundle_revert(t *testing.T) {
	arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferwbnbCode, conf.HighGas)
	arg.RevertList = []int{0, 1, 2}
	msg := conf.InvalidTx
	// 存在未记录在revertList中的 revert交易
	t.Run("sendValidBundle_revert", func(t *testing.T) {
		// bundle中均为 revert交易 && RevertList不为空
		t.Log("generate revert transaction \n")
		arg.TxCount = 4
		arg.RevertListAdd = []int{3}
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))

		}
		utils.BlockHeightIncreased(t)
		for _, tx := range txs {
			// bundle中的交易均不能上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}
	})
	t.Run("sendValidBundle_miss_revert", func(t *testing.T) {
		// bundle中包含 revert交易 和 合法交易
		t.Log("generate revert transaction \n")
		arg.TxCount = 10
		arg.RevertListAdd = []int{3, 4, 5}
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))

		}
		utils.BlockHeightIncreased(t)
		for _, tx := range txs {
			// bundle中的交易均不能上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}
	})
	t.Run("sendValidBundle_illegal_revert", func(t *testing.T) {
		t.Log("generate revert transaction \n")
		// bundle中均为 revert交易 && RevertList为空
		arg.TxCount = 3
		arg.RevertList = []int{}
		arg.Contract = conf.Mylock
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))

		}
		utils.BlockHeightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}
	})
}

/*
sendBundle 参数 - TxCount
*/
func Test_p2_sendBundle_txCount(t *testing.T) {
	arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferwbnbCode, conf.HighGas)
	msg := ""
	txCountLists := []int{0, 30, 3000, 99999}
	for _, count := range txCountLists {
		t.Run("txCounts_"+strconv.Itoa(count), func(t *testing.T) {
			arg.TxCount = count
			if count == 0 {
				msg = conf.MissTx
			} else if count == 99999 {
				msg = conf.LargeTx
			} else if count == 3000 {
				msg = conf.TxCountLimit
			}
			_, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
			err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
			if err != nil {
				log.Println(" failed: ", err.Error())
				assert.True(t, strings.Contains(err.Error(), msg))

			} else {
				log.Println("SendBundle succeed")
			}
			utils.BlockHeightIncreased(t)
			// for _, tx := range txs {
			// 	// 依次检查bundle中的交易是否成功上链
			// 	CheckBundleTx(t, *tx, valid, conf.TxSucceed)
			// }
		})
	}
}

/*
sendBundle 参数 - maxBlockNumber
*/
func Test_p1_sendBundle_maxBN(t *testing.T) {
	arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferwbnbCode, conf.HighGas)
	// maxBlockNumber最多设为当前区块号+100
	t.Run("maxBlockNumber_large", func(t *testing.T) {
		blockNum, err := arg.Client.BlockNumber(arg.Ctx)
		if err != nil {
			fmt.Println("BlockNumber", "err", err)
		}
		arg.MaxBN = blockNum + 101
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		err = arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), conf.MaxBlockNumberL))
		}
		utils.BlockHeightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}
	})
	// 过期的区块号
	maxBNLists := []int{0, 10, 999}
	for _, count := range maxBNLists {
		t.Run("maxBlockNumber_less_than_currentBlk", func(t *testing.T) {
			arg.MaxBN = uint64(count)
			txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
			err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
			if err != nil {
				log.Println(" failed: ", err.Error())
				assert.True(t, strings.Contains(err.Error(), conf.MaxBlockNumberC))
			}
			if count != 0 {
				for _, tx := range txs {
					utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
				}
			}
		})
	}

}

/*
sendBundle 参数 - maxTimeStamp
*/
func Test_p1_sendBundle_maxTS(t *testing.T) {
	arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferwbnbCode, conf.HighGas)
	t.Run("maxTimestamp_equal_current+300", func(t *testing.T) {
		currentTime := time.Now().Unix()
		futureTime := currentTime + int64(rand.Intn(timeLimit))
		convertedTime := uint64(futureTime)
		arg.MaxTS = &convertedTime
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		assert.Nil(t, err)
		t.Log(txs)
		// t.Log(tx_type)
		// BlockHeightIncreased(t)
		// for _, tx := range txs {
		// 	// 依次检查bundle中的交易是否成功上链
		// 	// 根据 交易index 确定校验的 tx_type
		// 	CheckBundleTx(t, *tx, true, conf.TxSucceed)
		// }
	})

	t.Run("maxTimestamp_legal_maxBlockNumber_work", func(t *testing.T) {
		currentTime := time.Now().Unix()
		// 指定bundle 最大为当前区块的下一个块，bundle时间为两个块后，expected 不会上链
		convertedTime := uint64(currentTime + 11)
		blockNum, err := arg.Client.BlockNumber(arg.Ctx)
		if err != nil {
			fmt.Println("BlockNumber", "err", err)
		}
		arg.MaxBN = blockNum + 1
		arg.MinTS = &convertedTime
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)
		utils.WaitMined(txs, cbn)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}
	})
}

func Test_p2_sendBundle_maxTS(t *testing.T) {
	arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferwbnbCode, conf.HighGas)
	var currentTime int64 = time.Now().Unix()
	t.Run("maxTS_more_than_current+300", func(t *testing.T) {
		msg := conf.TimestampTop
		convertedTime := uint64(currentTime + 301)
		arg.MaxTS = &convertedTime
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		// assert.Nil(t, err)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))
		}
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}
	})

	t.Run("maxTS_less_than_current", func(t *testing.T) {
		msg := conf.TimestampMC
		convertedTime := uint64(30000)
		arg.MaxTS = &convertedTime
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))
		}
		utils.BlockHeightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}

	})
}

/*
sendBundle 参数 - minTimeStamp
*/
func Test_p1_sendBundle_minTS(t *testing.T) {
	arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferwbnbCode, conf.HighGas)
	t.Run("sendValidBundle_arg_minT", func(t *testing.T) {
		//minTimestamp为未来时间 未超过设置时间 5*60 限制MaxBN生效
		var currentTime int64 = time.Now().Unix()
		convertedTime := uint64(currentTime + 2)
		convertedTime1 := uint64(currentTime + 3)
		blockNum, err := arg.Client.BlockNumber(arg.Ctx)
		if err != nil {
			fmt.Println("BlockNumber", "err", err)
		}
		arg.MaxBN = blockNum + 1
		arg.MinTS = &convertedTime
		arg.MaxTS = &convertedTime1
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)
		utils.WaitMined(txs, cbn)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}
	})
	t.Run("minTimeStamp_less_than_currentTime", func(t *testing.T) {
		currentTime := time.Now().Unix()
		convertedTime := uint64(currentTime - 20)
		convertedTime1 := uint64(currentTime + 1)
		arg.MaxBN = 0
		arg.MinTS = &convertedTime
		arg.MaxTS = &convertedTime1
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		cbn := utils.SendBundlesMined(t, arg, bundleArgs)
		utils.WaitMined(txs, cbn)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}
	})
}

func Test_p2_sendBundle_minTS(t *testing.T) {
	arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferwbnbCode, conf.HighGas)
	currentTime := time.Now().Unix()

	t.Run("sendValidBundle_arg_minTS_no_drop", func(t *testing.T) {
		//minTimestamp为未来时间 未超过设置时间 5*60 不限制MaxBN 生效
		convertedTime := uint64(currentTime + 6)
		convertedTime1 := uint64(currentTime + 7)
		arg.MinTS = &convertedTime
		arg.MaxTS = &convertedTime1
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		assert.Nil(t, err)

		utils.BlockHeightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}
		utils.BlockHeightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			blkN := utils.CheckBundleTx(t, *tx, true, conf.TxSucceed)
			txBlk = append(txBlk, blkN)
		}
		if !utils.TxInSameBlk(txBlk) {
			t.Fatalf("Transactions are not in the same block")
		}
	})

	//  minTimestamp 超出5*60
	t.Run("minTimestamp_later_than_Limit", func(t *testing.T) {
		convertedTime := uint64(currentTime + 301)
		arg.MinTS = &convertedTime
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), conf.TimestampMM))
		}
		utils.BlockHeightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			blkN := utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
			txBlk = append(txBlk, blkN)
		}
	})

	//  minTimestamp 为5*60 maxTimestamp 为5*60 +1
	t.Run("minTimestamp_equal_Limit", func(t *testing.T) {
		currentTime = time.Now().Unix()
		convertedTime := uint64(currentTime + int64(rand.Intn(timeLimit)))
		convertedTime1 := uint64(currentTime + int64(rand.Intn(timeLimit+1)))

		arg.MinTS = &convertedTime
		arg.MaxTS = &convertedTime1
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), conf.TimestampTop), "Expected result to be %d, but got %d", conf.TimestampTop, err.Error())
		}
		utils.BlockHeightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}
		// 300后可上链
	})

}

/*
sendBundle 参数 - maxTimeStamp&minTimeStamp
*/
func Test_p2_sendBundle_mmTS(t *testing.T) {
	arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferwbnbCode, conf.HighGas)
	var currentTime int64 = time.Now().Unix()
	msg := conf.TimestampMM
	convertedTime := uint64(currentTime + int64(rand.Intn(timeLimit)))

	//maxTimestamp 等于 minTimestamp
	t.Run("maxTimestamp_equal_minTimestamp", func(t *testing.T) {
		arg.MaxTS = &convertedTime
		arg.MinTS = &convertedTime
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))
		}
		utils.BlockHeightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}

	})

	//maxTimestamp 小于 minTimestamp
	t.Run("maxTimestamp_less_than_minTimestamp", func(t *testing.T) {
		convertedTime := uint64(currentTime + 3)
		convertedTime1 := uint64(currentTime + 1)
		arg.MinTS = &convertedTime
		arg.MaxTS = &convertedTime1
		txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(&arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))
		}
		utils.BlockHeightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			utils.CheckBundleTx(t, *tx, false, conf.TxFailed)
		}

	})
}

/*
sendBundle 并发
*/
// diff account
func Test_p0_sendBundle_batch(t *testing.T) {
	arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferwbnbCode, conf.HighGas)
	client := utils.CreateClient(conf.Url)
	client2 := utils.CreateClient(conf.Url_1)
	client3 := utils.CreateClient(conf.Url)

	chainID, err := client.ChainID(conf.Ctx)
	if err != nil {
		log.Printf("err %v\n", err)
	} else {
		log.Printf("==========获取当前链chainID ========== %v", chainID)
	}
	t.Run("sendValidBundle_batch", func(t *testing.T) {
		args := make([]*sendBundle.BidCaseArg, 2)
		args[0] = &arg
		args[1] = &sendBundle.BidCaseArg{
			Ctx:           conf.Ctx,
			Client:        client,
			ChainID:       chainID,
			RootPk:        conf.RootPk2,
			BobPk:         conf.RootPk2,
			Builder:       sendBundle.NewAccount(conf.BuilderPk),
			Validators:    []common.Address{common.HexToAddress(*conf.Validator)},
			BidClient:     client2,
			BuilderClient: client3,
			TxCount:       10,
			Data:          conf.TotallysplwbnbCode,
			Contract:      conf.WBNB,
			GasPrice:      big.NewInt(10e9),
			GasLimit:      big.NewInt(24000),
			SendAmount:    big.NewInt(0),
			RevertList:    []int{0},
		}

		wg := sync.WaitGroup{}
		for i := 0; i < 2; i++ {
			log.Println(args[i].RootPk)
			wg.Add(1)
			go func(i int) {
				time.Sleep(time.Duration(i) * time.Second)
				_, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(args[i])
				err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
				assert.Nil(t, err)
				wg.Done()
			}(i)
		}
		wg.Wait()

	})
}

// same account
func Test_p1_sendBundle_conflict(t *testing.T) {
	arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferwbnbCode, conf.HighGas)
	t.Run("sendValidBundle_conflict", func(t *testing.T) {
		args := make([]*sendBundle.BidCaseArg, 2)
		args[0] = &arg
		args[1] = &arg
		wg := sync.WaitGroup{}
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(i int) {
				time.Sleep(time.Duration(i) * time.Second)
				txs, bundleArgs, _ := sendBundle.ValidBundle_NilPayBidTx_1(args[i])
				err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
				if err != nil {
					assert.True(t, strings.Contains(err.Error(), conf.BundleConflict))
				}
				time.Sleep(5 * time.Second)
				for _, tx := range txs {
					// 依次检查bundle中的交易是否成功上链
					utils.CheckBundleTx(t, *tx, true, conf.TxSucceed)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
	})
}
