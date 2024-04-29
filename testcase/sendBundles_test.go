package newtestcases

import (
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/xkwang/cases"
)

/*
sendbundle 接口测试
*/

func setup() (*cases.BidCaseArg, string) {
	tx_type := "Transfer" // 默认为转账交易

	client, err := ethclient.Dial(url)
	if err != nil {
		log.Println("node ethclient.DialOptions", "err", err)
	}

	client2, err := ethclient.Dial(url_1)
	if err != nil {
		log.Println("client2 bidclient ethclient.DialOptions", "err", err)
	}

	client3, err := ethclient.Dial(url)
	if err != nil {
		log.Println("client3 bidclient ethclient.DialOptions", "err", err)
	}

	// query chainID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Printf("err %v\n", err)
	} else {
		log.Printf("==========获取当前链chainID ========== %v", chainID)
	}

	arg := &cases.BidCaseArg{
		Ctx:           ctx,
		Client:        client,  //客户端
		ChainID:       chainID, //client.ChainID
		RootPk:        rootPk,  //root Private Key
		BobPk:         bobPk,
		Builder:       cases.NewAccount(builderPk),
		Validators:    []common.Address{common.HexToAddress(*validator)},
		BidClient:     client2,
		BuilderClient: client3,
		TxCount:       5,
		Contract:      WBNB, // 调用合约的地址
		Data:          TransferWBNB_code,
		GasPrice:      big.NewInt(500),
		GasLimit:      big.NewInt(WBNB_gas),
		SendAmount:    big.NewInt(500),
		RevertList:    []int{},
		RevertListAdd: []int{},
		// 调用非转账合约 1）更新Data字段 2）SendAmount置为0
		// 确保提供的 Nonce 值是发送账户的下一个有效值
	}
	// t.Log(arg.Builder.Address.Hex())

	return arg, tx_type
}

// 正常sendBundle
func Test_p0_sendbundle(t *testing.T) {
	arg, tx_type := setup()
	t.Run("sendvalidbundle_tx", func(t *testing.T) {
		// bundle 中均为合法转账交易
		t.Log("Start sendBundle \n")
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		assert.Nil(t, err)
		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, true, Txsucceed, tx_type)
		}
	})
}

// sendBundle包含 revert交易
func Test_p0_sendbundle_revert(t *testing.T) {
	arg, tx_type := setup()
	t.Run("sendvalidbundle_all_revert", func(t *testing.T) {
		// revert 交易均在revertList中记录
		t.Log("generate revert transaction \n")
		arg.TxCount = 3
		arg.RevertList = []int{0, 1, 2}
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		assert.Nil(t, err)
		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, true, Txfailed, tx_type)
		}
	})

	t.Run("sendvalidbundle_part_revert", func(t *testing.T) {
		// revert 交易部分在revertList中记录
		t.Log("generate revert transaction \n")
		arg.TxCount = 10
		arg.RevertList = []int{0, 1, 2}
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		assert.Nil(t, err)
		BlockheightIncreased(t)
		for index, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			if index < 3 {
				checkBundleTx(t, *tx, true, Txfailed, tx_type)
			} else {
				checkBundleTx(t, *tx, true, Txsucceed, tx_type)
			}
		}
	})
}

// Todo:sendBundle包含 revert交易_异常情况
func Test_p1_sendbundle_revert(t *testing.T) {
	arg, tx_type := setup()
	arg.RevertList = []int{0, 1, 2}
	msg := InvalidTx
	// 存在未记录在revertList中的 revert交易
	t.Run("sendvalidbundle_revert", func(t *testing.T) {
		// bundle中均为 revert交易 && RevertList不为空
		t.Log("generate revert transaction \n")
		arg.TxCount = 4
		arg.RevertListAdd = []int{3}
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))

		}
		BlockheightIncreased(t)
		for _, tx := range txs {
			// bundule中的交易均不能上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}
	})
	t.Run("sendvalidbundle_miss_revert", func(t *testing.T) {
		// bundle中包含 revert交易 和 合法交易
		t.Log("generate revert transaction \n")
		arg.TxCount = 10
		arg.RevertListAdd = []int{3, 4, 5}
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))

		}
		BlockheightIncreased(t)
		for _, tx := range txs {
			// bundule中的交易均不能上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}
	})
	t.Run("sendvalidbundle_ilegal_revert", func(t *testing.T) {
		t.Log("generate revert transaction \n")
		// bundle中均为 revert交易 && RevertList为空
		arg.TxCount = 3
		arg.RevertList = []int{}
		arg.Data = TotallysplWBNB_code
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))

		}
		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}
	})
}

// todo: 一个bundle包含交易的数量 待产品确认，目前无限制
func Test_p2_sendbundle_arg_tx(t *testing.T) {
	arg, _ := setup()
	msg := ""
	txCountLists := []int{0, 30, 3000, 99999}
	for _, count := range txCountLists {
		t.Run("txCounts_"+strconv.Itoa(count), func(t *testing.T) {
			arg.TxCount = count
			if count == 0 {
				msg = MissTx
			} else if count == 99999 {
				msg = LargeTx
			} else if count == 3000 {
				msg = InvalidTx
			}
			_, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
			err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
			if err != nil {
				log.Println(" failed: ", err.Error())
				assert.True(t, strings.Contains(err.Error(), msg))

			} else {
				log.Println("SendBundle succeed")
			}
			BlockheightIncreased(t)
			// for _, tx := range txs {
			// 	// 依次检查bundle中的交易是否成功上链
			// 	checkBundleTx(t, *tx, valid, Txsucceed, tx_type)
			// }
		})
	}
}

// different accounts
func Test_p0_sendbundle_batch(t *testing.T) {
	arg, _ := setup()
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Println("node ethclient.DialOptions", "err", err)
	}

	client2, err := ethclient.Dial(url_1)
	if err != nil {
		log.Println("client2 bidclient ethclient.DialOptions", "err", err)
	}

	client3, err := ethclient.Dial(url)
	if err != nil {
		log.Println("client3 bidclient ethclient.DialOptions", "err", err)
	}
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Printf("err %v\n", err)
	} else {
		log.Printf("==========获取当前链chainID ========== %v", chainID)
	}
	t.Run("sendvalidbundle_batch", func(t *testing.T) {
		args := make([]*cases.BidCaseArg, 2)
		args[0] = arg
		args[1] = &cases.BidCaseArg{
			Ctx:           ctx,
			Client:        client,
			ChainID:       chainID,
			RootPk:        rootPk2,
			BobPk:         rootPk2,
			Builder:       cases.NewAccount(builderPk),
			Validators:    []common.Address{common.HexToAddress(*validator)},
			BidClient:     client2,
			BuilderClient: client3,
			TxCount:       10,
			Data:          TotallysplWBNB_code,
			Contract:      WBNB,
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
				_, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, args[i])
				err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
				assert.Nil(t, err)
				wg.Done()
			}(i)
		}
		wg.Wait()

	})
}

// same account
func Test_p1_sendbundle_conflict(t *testing.T) {
	arg, tx_type := setup()
	t.Run("sendvalidbundle_conflict", func(t *testing.T) {
		args := make([]*cases.BidCaseArg, 2)
		args[0] = arg
		args[1] = arg
		wg := sync.WaitGroup{}
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(i int) {
				time.Sleep(time.Duration(i) * time.Second)
				txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, args[i])
				err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
				if i == 1 {
					assert.True(t, strings.Contains(err.Error(), BundleConflict))
				} else {
					assert.Nil(t, err)
				}
				time.Sleep(5 * time.Second)
				for _, tx := range txs {
					// 依次检查bundle中的交易是否成功上链
					checkBundleTx(t, *tx, true, Txsucceed, tx_type)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
	})
}

// sendbundle_arg_maxBN
func Test_p1_sendbundle_maxBN(t *testing.T) {
	// maxBlockNumber最多设为当前区块号+100
	arg, tx_type := setup()
	t.Run("sendvalidbundle_arg_maxBN_large", func(t *testing.T) {
		blockNum, err := arg.Client.BlockNumber(arg.Ctx)
		if err != nil {
			fmt.Println("BlockNumber", "err", err)
		}
		arg.MaxBN = blockNum + 101
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err = arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), maxBlockNumberL))
		}
		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}
	})
	// 过期的区块号
	maxBNLists := []int{0, 10, 999}
	for _, count := range maxBNLists {
		t.Run("sendvalidbundle_arg_maxBN_small", func(t *testing.T) {
			arg.MaxBN = uint64(count)
			txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
			err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
			if err != nil {
				log.Println(" failed: ", err.Error())
				assert.True(t, strings.Contains(err.Error(), maxBlockNumberC))
			}
			if count != 0 {
				for _, tx := range txs {
					checkBundleTx(t, *tx, false, Txfailed, tx_type)
				}
			}
		})
	}

}

func Test_p1_sendbundle_maxTS(t *testing.T) {
	arg, tx_type := setup()
	//区块号合法
	t.Run("sendvalidbundle_arg_maxTS_legal", func(t *testing.T) {
		currentTime := time.Now().Unix()
		futureTime := currentTime + int64(rand.Intn(300))
		convertedTime := uint64(futureTime)
		arg.MaxTS = &convertedTime
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		assert.Nil(t, err)
		t.Log(txs)
		t.Log(tx_type)
		// BlockheightIncreased(t)
		// for _, tx := range txs {
		// 	// 依次检查bundle中的交易是否成功上链
		// 	// 根据 交易index 确定校验的 tx_type
		// 	checkBundleTx(t, *tx, true, Txsucceed, tx_type)
		// }
	})
	//区块号合法
	t.Run("sendvalidbundle_arg_maxTS_maxBN", func(t *testing.T) {
		currentTime := time.Now().Unix()
		// 指定bundle 最大为当前区块的下一个块，bundle时间为两个块后
		// expected 不会上链
		convertedTime := uint64(currentTime + 10)
		blockNum, err := arg.Client.BlockNumber(arg.Ctx)
		if err != nil {
			fmt.Println("BlockNumber", "err", err)
		}
		arg.MaxBN = blockNum + 2
		arg.MinTS = &convertedTime
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err = arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		assert.Nil(t, err)
		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}
	})
}

func Test_p2_sendbundle_maxTS(t *testing.T) {
	arg, tx_type := setup()
	var currentTime int64 = time.Now().Unix()
	t.Run("sendvalidbundle_arg_maxTS_large", func(t *testing.T) {
		msg := TimestampTop
		convertedTime := uint64(currentTime + 301)
		arg.MaxTS = &convertedTime
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		// assert.Nil(t, err)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))
		}
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}
	})

	t.Run("sendvalidbundle_arg_maxTS_small", func(t *testing.T) {
		msg := TimestampMC
		convertedTime := uint64(30000)
		arg.MaxTS = &convertedTime
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))
		}
		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}

	})
}
func Test_p1_sendbundle_minTS(t *testing.T) {
	// maxTimestamp最多设为当前区块号+5minutes
	arg, tx_type := setup()
	t.Run("sendvalidbundle_arg_minT", func(t *testing.T) {
		var currentTime int64 = time.Now().Unix()
		convertedTime := uint64(currentTime + 2)
		convertedTime1 := uint64(currentTime + 3)
		arg.MinTS = &convertedTime
		arg.MaxTS = &convertedTime1
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		assert.Nil(t, err)
		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}
	})
	t.Run("sendvalidbundle_arg_min_drop", func(t *testing.T) {
		currentTime := time.Now().Unix()
		convertedTime := uint64(currentTime - 3)
		convertedTime1 := uint64(currentTime + 1)
		arg.MinTS = &convertedTime
		arg.MaxTS = &convertedTime1
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		assert.Nil(t, err)

		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}
	})
}

func Test_p2_sendbundle_minTS(t *testing.T) {
	arg, tx_type := setup()
	currentTime := time.Now().Unix()

	t.Run("sendvalidbundle_arg_no_drop", func(t *testing.T) {
		convertedTime := uint64(currentTime + 7)
		convertedTime1 := uint64(currentTime + 8)
		arg.MinTS = &convertedTime
		arg.MaxTS = &convertedTime1
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		assert.Nil(t, err)

		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}
		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, true, Txsucceed, tx_type)
		}
	})
	//  maxTimestamp 超出5*60
	t.Run("sendvalidbundle_arg_minTS_large", func(t *testing.T) {
		convertedTime := uint64(currentTime + 301)
		arg.MinTS = &convertedTime
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), TimestampMM))
		}
		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}
	})

	//  maxTimestamp 为5*60
	t.Run("sendvalidbundle_arg_maxTS_limit", func(t *testing.T) {
		currentTime = time.Now().Unix()
		convertedTime := uint64(currentTime + 300)
		convertedTime1 := uint64(currentTime + 301)

		arg.MinTS = &(convertedTime)
		arg.MaxTS = &convertedTime1
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), TimestampTop), "Expected result to be %d, but got %d", TimestampTop, err.Error())
		}
		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}
		// 300后可上链
	})

}

func Test_p2_sendbundle_mmTS(t *testing.T) {
	arg, tx_type := setup()
	var currentTime int64 = time.Now().Unix()
	msg := TimestampMM
	convertedTime := uint64(currentTime + int64(rand.Intn(300)))

	//maxTimestamp 等于 minTimestamp
	t.Run("sendvalidbundle_arg_maxminTS", func(t *testing.T) {
		arg.MaxTS = &convertedTime
		arg.MinTS = &convertedTime
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))
		}
		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}

	})

	//maxTimestamp 小于 minTimestamp
	t.Run("sendvalidbundle_arg_minmaxTS", func(t *testing.T) {
		convertedTime := uint64(currentTime + 3)
		convertedTime1 := uint64(currentTime + 1)
		arg.MinTS = &convertedTime
		arg.MaxTS = &convertedTime1
		txs, bundleArgs, _ := cases.ValidBundle_NilPayBidTx_1(t, arg)
		err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
		if err != nil {
			log.Println(" failed: ", err.Error())
			assert.True(t, strings.Contains(err.Error(), msg))
		}
		BlockheightIncreased(t)
		for _, tx := range txs {
			// 依次检查bundle中的交易是否成功上链
			checkBundleTx(t, *tx, false, Txfailed, tx_type)
		}

	})
}

// 压力测试函数
// 持续发送bundles
func BenchmarkSendbundless(b *testing.B) {
	arg, _ := setup()
	// 循环执行测试函数
	for i := 0; i < b.N; i++ {
		// 在每次迭代中调用接口
		log.Println("run case")
		txs, err := cases.ValidBundle_NilPayBidTx_2(arg)
		if err != nil {
			log.Println(" failed: ", err.Error())
		} else {
			log.Println("ValidBundle_NilPayBidTx_1 succeed ")
		}
		println(txs)
		if err != nil {
			b.Fatalf("call failed: %v", err)
		}
	}
}
