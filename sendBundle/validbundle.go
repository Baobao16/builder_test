package sendBundle

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xkwang/conf"
	"log"
)

func ValidBundle_NilPayBidTx_1(arg *BidCaseArg) (types.Transactions, *types.SendBundleArgs, error) {
	txs, revertTxHashes := GenerateBNBTxs(arg, arg.SendAmount, arg.Data, arg.TxCount)
	txBytes := make([]hexutil.Bytes, 0, len(txs))
	for _, tx := range txs {
		txByte, err := tx.MarshalBinary()
		if err != nil {
			log.Println("tx.MarshalBinary", "err", err)
			return nil, nil, err
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

func ValidBundle_NilPayBidTx_2(arg *BidCaseArg, sim bool) (types.Transactions, error) {
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
	if sim {
		bundleArgs := &types.SendBundleArgs{
			Txs:               txBytes,
			RevertingTxHashes: revertTxHashes,
			MaxBlockNumber:    arg.MaxBN,
			MinTimestamp:      arg.MaxTS,
			MaxTimestamp:      arg.MinTS,
		}
		_, err := arg.BuilderClient.SendBundle(arg.Ctx, *bundleArgs)
		if err != nil {
			log.Println("failed to send bundle", "err", err)
		}
	} else {
		bundleArgs := &conf.SendBundleArgs{
			//MaxBlockNumber:    9,
			Txs:               txBytes,
			RevertingTxHashes: revertTxHashes,
			SimXYZ:            true,
		}
		err := arg.BuilderClient.Client().CallContext(arg.Ctx, nil, "eth_sendBundle", bundleArgs) //替换sendBundle  返回的是bundle哈希
		if err != nil {
			log.Println("failed to send bundle", "err", err)
		}
	}

	// bidJson, _ := json.MarshalIndent(bundleArgs, "", "  ")
	// log.Println(string(bidJson))

	return txs, nil

}

func RunValidBundleCases(arg *BidCaseArg) (types.Transactions, error) {
	log.Println("run case \n")
	txs, err := ValidBundle_NilPayBidTx_2(arg, true)
	if err != nil {
		log.Println(" failed: ", err.Error())
	} else {
		log.Println("ValidBundle_NilPayBidTx_2 succeed \n")
	}
	return txs, nil

}
