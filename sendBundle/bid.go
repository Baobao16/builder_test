package sendBundle

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/xkwang/conf"
	"log"
	"math/big"
)

var ()

type Account struct {
	Address    common.Address
	privateKey *ecdsa.PrivateKey
	Nonce      uint64
	//Nonce atomic.Int32
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
	RevertList       []int //会revert的交易，加入revertList
	RevertListAdd    []int //会revert的交易，但不加入revertList
	RevertListNormal []int //把正常交易加入revertList
	MaxBN            uint64
	MinTS            *uint64
	MaxTS            *uint64
}

func (a *Account) TransferBNB(nonce uint64, toAddress common.Address, data []byte, chainID *big.Int, amount *big.Int, gasPrice *big.Int, gaslimit *big.Int) (*types.Transaction, error) {
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    amount,
		GasPrice: gasPrice,
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

func GenerateBNBTxs(arg *BidCaseArg, amountPerTx *big.Int, data []byte, txCount int, nonce int) (types.Transactions, []common.Hash) {
	//bundleFactory := NewBidFactory(arg.Ctx, arg.Client, arg.RootPk, arg.BobPk, arg.Abc)
	revertTxHashes := make([]common.Hash, 0)
	rootAccount := NewAccount(arg.RootPk)
	var err error
	rootAccount.Nonce, err = arg.Client.PendingNonceAt(arg.Ctx, rootAccount.Address)
	if nonce != 0 {
		rootAccount.Nonce = rootAccount.Nonce - uint64(nonce)
	}
	if err != nil {
		fmt.Println("get nonce", "err", err)
	}
	revertTxs := make(map[int]string)

	for _, v := range arg.RevertListNormal {
		revertTxs[v] = "RevertListNormal"
	}
	for _, v := range arg.RevertList {
		revertTxs[v] = "RevertList"
	}
	for _, v := range arg.RevertListAdd {
		revertTxs[v] = "RevertListAdd"
	}

	txs := make([]*types.Transaction, 0)
	for i := 0; i < txCount; i++ {
		var (
			bundle *types.Transaction
			err    error
		)
		rvtData := common.Hex2Bytes("346c94cf000000000000000000000000000000000000000000000000000000000000006f0000000000000000000000000000000000000000000000000000000000000000")
		if revertTxs[i] == "RevertListNormal" {
			// 把正常交易加入revertList中
			bundle, err = rootAccount.TransferBNB(rootAccount.Nonce, arg.Contract, data, arg.ChainID, amountPerTx, arg.GasPrice, arg.GasLimit)
			revertTxHashes = append(revertTxHashes, bundle.Hash())
			// fmt.Printf("List revert txHash %v\n", bundle.Hash().Hex())

		} else if revertTxs[i] == "RevertList" {
			// RevertList 中的交易设置为会revert
			bundle, err = rootAccount.TransferBNB(rootAccount.Nonce, conf.Mylock, rvtData, arg.ChainID, amountPerTx, arg.GasPrice, arg.GasLimit)
			fmt.Printf("noList revert txHash %v\n", bundle.Hash().Hex())
			revertTxHashes = append(revertTxHashes, bundle.Hash())

		} else if revertTxs[i] == "RevertListAdd" {
			// 发送会revert的交易，但不加入revertList
			bundle, err = rootAccount.TransferBNB(rootAccount.Nonce, conf.Mylock, rvtData, arg.ChainID, amountPerTx, arg.GasPrice, arg.GasLimit)
			// fmt.Printf("noList revert txHash %v\n", bundle.Hash().Hex())

		} else {
			// arg.GasLimit = big.NewInt(int64(int(BNBGasUsed) + rand.Intn(3300)))
			bundle, err = rootAccount.TransferBNB(rootAccount.Nonce, arg.Contract, data, arg.ChainID, amountPerTx, arg.GasPrice, arg.GasLimit)
			// fmt.Printf(" validtx txHash %v\n", bundle.Hash().Hex())

		}
		if err != nil {
			fmt.Println("fail to sign tx TransferBNB", "err", err)
		}
		log.Printf("Generate txHash %v  [gasPrice: %v ,gasLimit: %v ,SendAmount :%v]\n", bundle.Hash().Hex(), arg.GasPrice, arg.GasLimit, amountPerTx)
		txs = append(txs, bundle)
		rootAccount.Nonce = rootAccount.Nonce + 1
		//rootAccount.Nonce.Add(1)
	}

	return txs, revertTxHashes
}
