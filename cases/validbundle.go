package cases

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func ValidBundle_NilPayBidTx_1(t *testing.T, arg *BidCaseArg) (types.Transactions, *types.SendBundleArgs, error) {
	txs, revertTxHashes := GenerateBNBTxs(arg, arg.SendAmount, arg.Data, arg.TxCount) // []common.Hash,
	txBytes := make([]hexutil.Bytes, 0)
	for _, tx := range txs {
		// log.Println(tx.Nonce())
		txByte, err := tx.MarshalBinary()
		// fmt.Printf("txhash %v\n", tx.Hash().Hex())
		if err != nil {
			log.Println("tx.MarshalBinary", "err", err)
		}
		txBytes = append(txBytes, txByte)
	}

	bundleArgs := &types.SendBundleArgs{
		Txs:               txBytes,
		RevertingTxHashes: revertTxHashes,
		MaxBlockNumber:    arg.MaxBN,
		MinTimestamp:      arg.MinTS,
		MaxTimestamp:      arg.MaxTS,
	}
	bidJson, _ := json.MarshalIndent(bundleArgs, "", "  ")
	log.Println(string(bidJson))

	return txs, bundleArgs, nil

}

func ValidBundle_NilPayBidTx_2(arg *BidCaseArg) (types.Transactions, error) {
	txs, revertTxHashes := GenerateBNBTxs(arg, arg.SendAmount, arg.Data, arg.TxCount) // []common.Hash,
	txBytes := make([]hexutil.Bytes, 0)
	for _, tx := range txs {
		// log.Println(tx.Nonce())
		txByte, err := tx.MarshalBinary()
		// fmt.Printf("txhash %v\n", tx.Hash().Hex())
		if err != nil {
			log.Println("tx.MarshalBinary", "err", err)
		}
		txBytes = append(txBytes, txByte)
	}

	bundleArgs := &types.SendBundleArgs{
		Txs:               txBytes,
		RevertingTxHashes: revertTxHashes,
		MaxBlockNumber:    arg.MaxBN,
		MinTimestamp:      arg.MaxTS,
		MaxTimestamp:      arg.MinTS,
	}
	// bidJson, _ := json.MarshalIndent(bundleArgs, "", "  ")
	// log.Println(string(bidJson))

	err := arg.BuilderClient.SendBundle(arg.Ctx, bundleArgs)
	if err != nil {
		log.Println("failed to send bundle", "err", err)
	}

	return txs, nil

}

func RunValidBundleCases(arg *BidCaseArg) (types.Transactions, error) {
	log.Println("run case \n")
	txs, err := ValidBundle_NilPayBidTx_2(arg)
	if err != nil {
		log.Println(" failed: ", err.Error())
	} else {
		log.Println("ValidBundle_NilPayBidTx_2 succeed \n")
	}
	return txs, nil

}
