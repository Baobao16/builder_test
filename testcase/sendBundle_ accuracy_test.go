package newtestcases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xkwang/cases"
)

// func setup() (*cases.BidCaseArg, string) {
// 	tx_type := "Transfer" // 默认为转账交易

// 	client, err := ethclient.Dial(url)
// 	if err != nil {
// 		log.Println("node ethclient.DialOptions", "err", err)
// 	}

// 	client2, err := ethclient.Dial(url_1)
// 	if err != nil {
// 		log.Println("client2 bidclient ethclient.DialOptions", "err", err)
// 	}

// 	client3, err := ethclient.Dial(url)
// 	if err != nil {
// 		log.Println("client3 bidclient ethclient.DialOptions", "err", err)
// 	}

// 	// query chainID
// 	chainID, err := client.ChainID(ctx)
// 	if err != nil {
// 		log.Printf("err %v\n", err)
// 	} else {
// 		log.Printf("==========获取当前链chainID ========== %v", chainID)
// 	}

// 	arg := &cases.BidCaseArg{
// 		Ctx:           ctx,
// 		Client:        client,  //客户端
// 		ChainID:       chainID, //client.ChainID
// 		RootPk:        rootPk,  //root Private Key
// 		BobPk:         bobPk,
// 		Builder:       cases.NewAccount(builderPk),
// 		Validators:    []common.Address{common.HexToAddress(*validator)},
// 		BidClient:     client2,
// 		BuilderClient: client3,
// 		TxCount:       5,
// 		Contract:      WBNB, // 调用合约的地址
// 		Data:          TransferWBNB_code,
// 		GasPrice:      big.NewInt(500),
// 		GasLimit:      big.NewInt(WBNB_gas),
// 		SendAmount:    big.NewInt(500),
// 		RevertList:    []int{},
// 		RevertListAdd: []int{},
// 		// 调用非转账合约 1）更新Data字段 2）SendAmount置为0
// 		// 确保提供的 Nonce 值是发送账户的下一个有效值
// 	}
// 	// t.Log(arg.Builder.Address.Hex())

// 	return arg, tx_type
// }

// 正常sendBundle
func Test_p0_backrun(t *testing.T) {
	arg, tx_type := setup()
	t.Run("sendvalidbundle_tx", func(t *testing.T) {
		//构造一笔交易 tx1
		t.Log("Expose mempool transaction tx1\n")

		t.Log("User 1 sends bundle [tx1, tx2], none are allowed to revert.\n")

		t.Log("User 2 sends bundle [tx1, tx2], none are allowed to revert.\n")

		t.Log("Expose mempool transaction tx1\n")

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
