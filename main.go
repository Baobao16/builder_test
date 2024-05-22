package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xkwang/sendBundle"
)

// func main() {
// 	client, err := ethclient.Dial("https://bsc-testnet-rpc.publicnode.com")
// 	if err != nil {
// 		panic(err)
// 	}
// 	ctx := context.Background()
// 	blockNumber, err := client.BlockNumber(ctx)
// 	if err != nil {
// 		fmt.Printf("err %v", err)
// 	} else {
// 		fmt.Printf("blockNumber %v", blockNumber)
// 	}

// }

var (
	priKey = os.Getenv("PRIVATE_KEY")

	// setting: root bnb&abc boss
	rootPrivateKey = flag.String("rootpk",
		strings.TrimPrefix(priKey, "0x"),
		"private key of root account")

	root2PrivateKey = flag.String("rootpk2",
		// "08560159b1ad148be0c695b511a6bdaa91562158f36ea88ef48ea8355ee64755",
		"61bfe9aea17bec5de54a86ad6cb0418f678a2fc8b746cc3901687eaebe1da809",
		"private key of root2 account")

	// root3PrivateKey = flag.String("rootpk3",
	// 	"eb1ee3f15d54f3afcc735ddac56ef8498a006c0bb999a9c267bbf99414698f11",
	// 	"private key of root3 account")

	// root4PrivateKey = flag.String("rootpk4",
	// 	"7540900d280a6df50c6bcaeda216d97df23afb444f82ad840321de853b6bfe9c",
	// 	"private key of root4 account")

	// root5PrivateKey = flag.String("rootpk5",
	// 	"446cdc7ef45999fb635dcbf18acaccd4a796cb7c4fd560b3a6c39b87723e4fc8",
	// 	"private key of root5 account")

	// root6PrivateKey = flag.String("rootpk6",
	// 	"50b9bb6c14ad320ec12b3e21e16296a446059a2453bb9b323a00eb2e051c5eb5",
	// 	"private key of root6 account")

	// root7PrivateKey = flag.String("rootpk7",
	// 	"8fb1b911b16cc94cb2edb8b707c782121c2cf70cd71f2adf2e8bb52bb967a2c4",
	// 	"private key of root7 account")

	builderPrivateKey = flag.String("builderpk",
		"7b94e64fc431b0daa238d6ed8629f3747782b8bc10fb8a41619c5fb2ba55f4e3",
		"private key of builder account")

	validator = flag.String("validator", "0xF474Cf03ccEfF28aBc65C9cbaE594F725c80e12d", "validator address")
)

func main() {

	ctx := context.Background()

	rootPk := *rootPrivateKey
	// rootPk2 := *root2PrivateKey
	bobPk := *rootPrivateKey
	builderPk := *builderPrivateKey

	// client, err := ethclient.Dial("https://bsc-testnet-rpc.publicnode.com")
	client, err := ethclient.Dial("http://10.2.66.75:28545")
	if err != nil {
		fmt.Println("node ethclient.DialOptions", "err", err)
	}

	// client2, err := ethclient.Dial("https://bsc-testnet-elbrus.bnbchain.org")
	// client2, err := ethclient.Dial("https://bsc-testnet-ararat.bnbchain.org")
	// client2, err := ethclient.Dial("http://localhost:28545")
	client2, err := ethclient.Dial("http://10.2.66.75:18545")
	if err != nil {
		fmt.Println("client2 bidclient ethclient.DialOptions", "err", err)
	}

	// client3, err := ethclient.Dial("https://bsc-testnet-builder.bnbchain.org")
	client3, err := ethclient.Dial("http://10.2.66.75:28545")
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

	arg := &sendBundle.BidCaseArg{
		Ctx:           ctx,
		Client:        client,
		ChainID:       chainID,
		RootPk:        rootPk,
		BobPk:         bobPk,
		Builder:       sendBundle.NewAccount(builderPk),
		Validators:    []common.Address{common.HexToAddress(*validator)},
		BidClient:     client2,
		BuilderClient: client3,
		TxCount:       3,
		// Contract:      common.HexToAddress("0x7b09bb26c9fef574ea980a33fc71c184405a4023"),
		Contract:   sendBundle.WBNB,
		Data:       sendBundle.TransferWBNB_code,
		GasPrice:   big.NewInt(500),
		GasLimit:   big.NewInt(50000),
		SendAmount: big.NewInt(50),
		// RevertList: []int{1},
	}

	fmt.Println(arg.Builder.Address.Hex())

	//单发一个bid
	// sendBundle.RunValidCases(arg)

	//单发一个bundle

	// sendBundle.RunValidBundleCases(arg)

	//单发一个转账交易
	sendBundle.RunValidSendCases(arg)

	//单发一个裸交易
	sendBundle.SendRaw(arg)

	//发两个bid

	// args := make([]*sendBundle.BidCaseArg, 2)
	// args[0] = arg
	// args[1] = &sendBundle.BidCaseArg{
	// 	Ctx:           ctx,
	// 	Client:        client,
	// 	ChainID:       chainID,
	// 	RootPk:        rootPk2,
	// 	BobPk:         bobPk,
	// 	Builder:       sendBundle.NewAccount(builderPk),
	// 	Validators:    []common.Address{common.HexToAddress(*validator)},
	// 	BidClient:     client2,
	// 	BuilderClient: client3,
	// 	TxCount:       5,
	// 	Contract:      sendBundle.WBNB,
	// 	Data:          sendBundle.TransferWBNB_code,
	// 	GasPrice:      big.NewInt(500),
	// 	GasLimit:      big.NewInt(50000),
	// 	SendAmount:    big.NewInt(50),
	// 	RevertList:    []int{0},
	// }

	// // //发两个bid
	// wg := sync.WaitGroup{}
	// for i := 0; i < 2; i++ {
	// 	wg.Add(1)
	// 	go func(i int) {
	// 		time.Sleep(time.Duration(i) * time.Second)
	// 		sendBundle.RunValidBundleCases(args[i])
	// 		wg.Done()
	// 	}(i)
	// }
	// wg.Wait()
	// println("done!!!!!")

}
