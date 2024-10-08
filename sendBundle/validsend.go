package sendBundle

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// func ValidSend_NilPayBidTx_1(arg *BidCaseArg) error {
// 	txs := GenerateBNBTxs(arg, TransferAmountPerTx, nil, 1)
// 	err := arg.Client.SendTransaction(arg.Ctx, txs[0])
// 	if err != nil {
// 		fmt.Println("failed to send bundle", "err", err)
// 	}
// 	return nil
// }

func ValidSend_ContractTx_1(arg *BidCaseArg) (types.Transaction, error) {

	txs, _ := GenerateBNBTxs(arg, arg.SendAmount, arg.Data, 1)

	txBytes := make([]hexutil.Bytes, 0)
	for _, tx := range txs {
		log.Println(tx.Nonce())
		txByte, err := tx.MarshalBinary()
		fmt.Printf("txhash %v\n", tx.Hash().Hex())
		if err != nil {
			log.Println("tx.MarshalBinary", "err", err)
		}
		txBytes = append(txBytes, txByte)
	}

	err := arg.Client.SendTransaction(arg.Ctx, txs[0])
	if err != nil {
		fmt.Println("failed to send single Transaction", "err", err)
	}
	return *txs[0], nil
}
func ValidSend_ContractTx(arg *BidCaseArg) (types.Transaction, error) {

	txs, _ := GenerateBNBTxs(arg, arg.SendAmount, arg.Data, 1)

	txBytes := make([]hexutil.Bytes, 0)
	for _, tx := range txs {
		log.Println(tx.Nonce())
		txByte, err := tx.MarshalBinary()
		fmt.Printf("txhash %v\n", tx.Hash().Hex())
		if err != nil {
			log.Println("tx.MarshalBinary", "err", err)
		}
		txBytes = append(txBytes, txByte)
	}
	return *txs[0], nil
}

func SendRaw(arg *BidCaseArg) error {
	var tx types.Transaction
	tx.UnmarshalBinary(common.Hex2Bytes("f86514843b9aca0083033450947b09bb26c9fef574ea980a33fc71c184405a4023808081e5a0d5e3d792a94528787a4ea713aa9a97f0459ced82878347231bb3f5110e37ac86a07c51444a65829ec2c0bd7c342cfe6586651dcdc1cb6ee75402491966759cfeb6"))
	fmt.Printf("txhash %v\n", tx.Hash().Hex())
	return arg.Client.SendTransaction(arg.Ctx, &tx)
}

func RunValidSendCases(arg *BidCaseArg) {
	print("run case ")
	_, err := ValidSend_ContractTx_1(arg)
	if err != nil {
		print(" failed: ", err.Error())
	} else {
		print("ValidSend_ContractTx_1 succeed")
	}
	println()

}
