package cases

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	TransferAmountPerTx = big.NewInt(1e16)
	DefaultGasLimit     = uint64(21000)
	BNBGasUsed          = 21000
	HighBNBGasPrice     = big.NewInt(1.1*1e9 + 1)
)

func RunValidCases(arg *BidCaseArg) {

	print("run case ")
	err := ValidBid_NilPayBidTx_1(arg)
	if err != nil {
		print(" failed: ", err.Error())
	} else {
		print("ValidBid_NilPayBidTx_1 succeed")
	}
	println()

}

// ValidBid_NilPayBidTx_1
// gasFee = 21000 * 1 * 0.0000001 BNB = 0.42/200 BNB
func ValidBid_NilPayBidTx_1(arg *BidCaseArg) error {
	// TODO: validator could ignore revert tx
	txs, _ := GenerateBNBTxs(arg, arg.SendAmount, arg.Data, arg.TxCount)

	for {
		time.Sleep(5000 * time.Millisecond)
		bid := GenerateValidBid(arg, txs)
		bidJson, _ := json.MarshalIndent(bid, "", "  ")
		println("---------------------")
		println(string(bidJson))
		//return nil
		_, err := arg.BidClient.SendBid(arg.Ctx, *bid)
		if err == nil {
			break
		}
		bidErr, ok := err.(rpc.Error)
		if ok && bidErr.ErrorCode() == types.MevNotInTurnError {
			continue
		} else {
			println("send bid failed: ", err.Error())
		}
	}
	fmt.Println("-----success------")
	// inTurn := true

	// ping := func() {
	// 	_, err := arg.BidClient.SendBid(arg.Ctx, *bid)
	// 	if err != nil {
	// 		// bidErr, ok := err.(rpc.Error)
	// 		// if ok && bidErr.ErrorCode() == types.InvalidBidParamError {
	// 		// 	inTurn = true
	// 		// }
	// 		println("send bid failed: ", err.Error())
	// 	}
	// }

	// ping()

	// for inTurn != true {
	// 	println("wait for in turn")
	// 	time.Sleep(500 * time.Millisecond)
	// 	ping()
	// }

	return nil
}

func GenerateValidBid(arg *BidCaseArg, txs []*types.Transaction) *types.BidArgs {
	txBytes := make([]hexutil.Bytes, 0)
	for _, tx := range txs {
		txByte, err := tx.MarshalBinary()
		if err != nil {
			fmt.Println("tx.MarshalBinary", "err", err)
		}
		txBytes = append(txBytes, txByte)
	}
	blockNum, err := arg.Client.BlockNumber(arg.Ctx)
	if err != nil {
		fmt.Println("BlockNumber", "err", err)
	}

	block, err := arg.Client.BlockByNumber(arg.Ctx, big.NewInt(int64(blockNum)))
	if err != nil {
		fmt.Println("BlockByNumber", "err", err)
	} else {
		fmt.Println("block", block.NumberU64())
	}

	rawBid := &types.RawBid{
		BlockNumber: block.NumberU64() + 1,
		ParentHash:  block.Hash(),
		Txs:         txBytes,
		GasUsed:     arg.GasLimit.Uint64() * uint64(arg.TxCount),
		GasFee:      big.NewInt(arg.GasLimit.Int64() * int64(arg.TxCount) * arg.GasPrice.Int64()),
	}

	bidArgs := arg.Builder.SignBid(rawBid)
	return bidArgs

}
