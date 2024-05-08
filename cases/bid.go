package cases

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	WBNB                 = common.HexToAddress("0xE5454b639B241c07Fc0d55b23690F9CeE18b7E4f")
	TransferBNB_code     = common.Hex2Bytes("a6f9dae10000000000000000000000007b09bb26c9fef574ea980a33fc71c184405a4023")
	TransferWBNB_code    = common.Hex2Bytes("1a695230000000000000000000000000e5454b639b241c07fc0d55b23690f9cee18b7e4f")
	TBalanceOfWBNB_code  = common.Hex2Bytes("70a08231000000000000000000000000e5454b639b241c07fc0d55b23690f9cee18b7e4f")
	TotallysplWBNB_code  = common.Hex2Bytes("18160ddd")
	AllowanceRouter_code = common.Hex2Bytes("3e5beab9000000000000000000000000e1f45ef433b2adf7583917974543a2df2161dd6c")
)

type Account struct {
	Address    common.Address
	privateKey *ecdsa.PrivateKey
	Nonce      uint64
}

func (a *Account) TransferBNB(nonce uint64, toAddress common.Address, data []byte, chainID *big.Int, amount *big.Int, gasprice *big.Int, gaslimit *big.Int) (*types.Transaction, error) {
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    amount,
		GasPrice: gasprice,
		Gas:      gaslimit.Uint64(),
		Data:     data,
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), a.privateKey)
	if err != nil {
		fmt.Println("failed to sign tx", "err", err)
		return nil, err
	}

	return signedTx, nil
}

func (a *Account) SignBid(rawBid *types.RawBid) *types.BidArgs {
	data, err := rlp.EncodeToBytes(rawBid)
	if err != nil {
		fmt.Println("failed to encode raw bid", "err", err)
	}

	sig, err := crypto.Sign(crypto.Keccak256(data), a.privateKey)
	if err != nil {
		fmt.Println("failed to sign raw bid", "err", err)
	}

	bidArgs := types.BidArgs{
		RawBid:    rawBid,
		Signature: sig,
	}

	return &bidArgs
}

func NewAccount(privateKey string) *Account {
	privateECDSAKey, address := PriKeyToAddress(privateKey)

	return &Account{
		Address:    address,
		privateKey: privateECDSAKey,
	}
}

func PriKeyToAddress(privateKey string) (*ecdsa.PrivateKey, common.Address) {
	privateECDSAKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		fmt.Println("failed to parse private key", "err", err)
	}

	publicKey, ok := privateECDSAKey.Public().(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("failed to get public key", "err", err)
	}

	selfAddress := crypto.PubkeyToAddress(*publicKey)

	return privateECDSAKey, selfAddress
}

type BidCaseArg struct {
	Ctx              context.Context
	Client           *ethclient.Client
	BidClient        *ethclient.Client
	BuilderClient    *ethclient.Client
	ChainID          *big.Int
	RootPk, BobPk    string
	Contract         common.Address
	Builder          *Account
	Validators       []common.Address
	TxCount          int
	Data             []byte
	GasLimit         *big.Int
	GasPrice         *big.Int
	SendAmount       *big.Int
	RevertList       []int
	RevertListAdd    []int
	RevertListnormal []int
	MaxBN            uint64
	MinTS            *uint64
	MaxTS            *uint64
}

func GenerateBNBTxs(arg *BidCaseArg, amountPerTx *big.Int, data []byte, txcount int) (types.Transactions, []common.Hash) {
	//bundleFactory := NewBidFactory(arg.Ctx, arg.Client, arg.RootPk, arg.BobPk, arg.Abc)
	revertTxHashes := make([]common.Hash, 0)
	rootAccount := NewAccount(arg.RootPk)
	var err error
	rootAccount.Nonce, err = arg.Client.PendingNonceAt(arg.Ctx, rootAccount.Address)
	if err != nil {
		fmt.Println("get nonce", "err", err)
	}
	revertTxs := make(map[int]string)

	for _, v := range arg.RevertListnormal {
		revertTxs[v] = "RevertListnormal"
	}
	for _, v := range arg.RevertList {
		revertTxs[v] = "RevertList"
	}
	for _, v := range arg.RevertListAdd {
		revertTxs[v] = "RevertListAdd"
	}

	txs := make([]*types.Transaction, 0)
	for i := 0; i < txcount; i++ {
		var (
			bundle *types.Transaction
			err    error
		)
		if revertTxs[i] == "RevertListnormal" {
			// 把正常交易加入revertlist中
			bundle, err = rootAccount.TransferBNB(rootAccount.Nonce, arg.Contract, data, arg.ChainID, amountPerTx, arg.GasPrice, arg.GasLimit)
			revertTxHashes = append(revertTxHashes, bundle.Hash())
			// fmt.Printf("List revert txhash %v\n", bundle.Hash().Hex())

		} else if revertTxs[i] == "RevertList" {
			// RevertList 中的交易设置为会revert
			bundle, err = rootAccount.TransferBNB(rootAccount.Nonce, WBNB, TotallysplWBNB_code, arg.ChainID, amountPerTx, arg.GasPrice, arg.GasLimit)
			// fmt.Printf("noList revert txhash %v\n", bundle.Hash().Hex())
			revertTxHashes = append(revertTxHashes, bundle.Hash())

		} else if revertTxs[i] == "RevertListAdd" {
			// 发送会revert的交易，但不加入revertList
			bundle, err = rootAccount.TransferBNB(rootAccount.Nonce, WBNB, TotallysplWBNB_code, arg.ChainID, amountPerTx, arg.GasPrice, arg.GasLimit)
			// fmt.Printf("noList revert txhash %v\n", bundle.Hash().Hex())

		} else {
			// arg.GasLimit = big.NewInt(int64(int(BNBGasUsed) + rand.Intn(3300)))
			bundle, err = rootAccount.TransferBNB(rootAccount.Nonce, arg.Contract, data, arg.ChainID, amountPerTx, arg.GasPrice, arg.GasLimit)
			// fmt.Printf(" validtx txhash %v\n", bundle.Hash().Hex())

		}
		if err != nil {
			fmt.Println("fail to sign tx TransferBNB", "err", err)
		}
		log.Printf("Txhash %v in bundle [gasPrice: %v ,gasLimit: %v ,SendAmount :%v]\n", bundle.Hash().Hex(), arg.GasPrice, arg.GasLimit, amountPerTx)
		txs = append(txs, bundle)
		rootAccount.Nonce = rootAccount.Nonce + 1
	}

	return txs, revertTxHashes
}
