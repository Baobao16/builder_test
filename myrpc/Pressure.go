package new

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/xkwang/conf"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xkwang/sendBundle"
)

var (
	priKey = os.Getenv("PRIVATE_KEY")

	// setting: root bnb&abc boss
	rootPrivateKey = flag.String("rootpk",
		strings.TrimPrefix(priKey, "0x"),
		"private key of root account")

	root2PrivateKey = flag.String("rootpk2",
		"08560159b1ad148be0c695b511a6bdaa91562158f36ea88ef48ea8355ee64755",
		"private key of root2 account")

	builderPrivateKey = flag.String("builderpk",
		"7b94e64fc431b0daa238d6ed8629f3747782b8bc10fb8a41619c5fb2ba55f4e3",
		"private key of builder account")

	validator = flag.String("validator", "0xF474Cf03ccEfF28aBc65C9cbaE594F725c80e12d", "validator address")
)

// 定义一个结构体表示用户
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// 定义一个处理函数，处理API请求
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// 创建一个用户
	// user := User{
	// 	ID:   1,
	// 	Name: "John Doe",
	// }
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
		TxCount:       rand.Intn(20),
		// Contract:      common.HexToAddress("0x7b09bb26c9fef574ea980a33fc71c184405a4023"),
		Contract:   common.HexToAddress("0x199e3Bfb54f4aAa9D67d1BB56429c5ef9D1A2A91"), // 合约地址
		Data:       conf.TransferBNBCode,                                              //合约方法
		GasPrice:   big.NewInt(2e9),
		GasLimit:   big.NewInt(22000),
		SendAmount: big.NewInt(1e18), //非转账交易需设置为0
		// RevertList: []int{1},
	}

	txs, err := sendBundle.RunValidBundleCases(arg)

	txBytes := make([]hexutil.Bytes, 0)
	for _, tx := range txs {
		txByte, _ := tx.MarshalBinary()
		txBytes = append(txBytes, txByte)
	}

	// 将用户数据编码为JSON格式
	jsonData, err := json.Marshal(txBytes)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 写入响应数据
	w.Write(jsonData)
}

func new() {
	// 定义路由，指定处理函数
	http.HandleFunc("/user", getUserHandler)

	// 启动HTTP服务器，监听端口8080
	log.Println("Server started on port 8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
