package conf

import (
	"context"
	"flag"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	Url       = "http://10.1.8.114:28545"
	Url_1     = "http://10.1.8.114:18545"
	Tx_type   = "Transfer"
	Ctx       = context.Background()
	RootPk    = *rootPrivateKey
	RootPk2   = *root2PrivateKey
	RootPk3   = *root3PrivateKey
	RootPk4   = *root4PrivateKey
	BobPk     = *rootPrivateKey
	BuilderPk = *BuilderPrivateKey
	PriKey    = os.Getenv("PRIVATE_KEY")

	// setting: root bnb&abc boss
	rootPrivateKey = flag.String("rootpk",
		strings.TrimPrefix(PriKey, "0x"),
		"private key of root account")



	Validator = flag.String("validator", "0xF474Cf03ccEfF28aBc65C9cbaE594F725c80e12d", "validator address")
)

// 模拟 bundle 不会真的上链
type SendBundleArgs struct {
	Txs               []hexutil.Bytes `json:"txs"`
	MaxBlockNumber    uint64          `json:"maxBlockNumber"`
	MinTimestamp      *uint64         `json:"minTimestamp"`
	MaxTimestamp      *uint64         `json:"maxTimestamp"`
	RevertingTxHashes []common.Hash   `json:"revertingTxHashes"`
	SimXYZ            bool            `json:"simXYZ"`
}

// bundleArgs := &SendBundleArgs{
// 	//MaxBlockNumber:    9,
// 	Txs:               txBytes,
// 	RevertingTxHashes: revertTxHashes,
// 	SimXYZ:            true,
// }

// bidJson, _ := json.MarshalIndent(bundleArgs, "", "  ")
// println(string(bidJson))

// err := arg.BuilderClient.Client().CallContext(arg.Ctx, nil, "eth_sendBundle", bundleArgs) //替换sendBundle  返回的是bundle哈希
