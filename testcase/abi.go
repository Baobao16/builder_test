package testcase

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/xkwang/conf"
	"github.com/xkwang/sendBundle"
	"github.com/xkwang/utils"
	"log"
	"math/big"
	"testing"
)

var LockABIfile = utils.GeneABI(conf.LockPath)

// step-1 读取abi文件
var LockABI = utils.Contract{ABI: *LockABIfile}
var ResetData = utils.GeneEncodedData(LockABI, "reset")
var LockData = utils.GeneEncodedData(LockABI, "lock", 1, true) // "gasUsed":"0x6ab1"
var LockFData = utils.GeneEncodedData(LockABI, "lock", 1, false)

var UnlockStrData = utils.GeneEncodedData(LockABI, "unlock", 1, "str")
var UnlockMoreData = utils.GeneEncodedData(LockABI, "unlock", 1, "more")
var UnlockDeMoreData = utils.GeneEncodedData(LockABI, "unlock_de", 1, "more") //"0x1d58c7"1923271
var UnlockDeStrData = utils.GeneEncodedData(LockABI, "unlock_de", 1, "str")
var FakelockStrData = utils.GeneEncodedData(LockABI, "fakelock", 1, "str")
var FakelockMoreData = utils.GeneEncodedData(LockABI, "fakelock", 1, "more")
var UseGas = utils.GeneEncodedData(LockABI, "increaseGasUsed")
var (
	Args            = make([]*sendBundle.BidCaseArg, 2)
	Args3           = make([]*sendBundle.BidCaseArg, 3)
	UsrList         = make([]utils.TxStatus, 3)
	UsrList6        = make([]utils.TxStatus, 6)
	BundleargsLsit  = make([]*types.SendBundleArgs, 2)
	BundleargsLsit3 = make([]*types.SendBundleArgs, 3)
)

func AddUserBundle(pk string, lock common.Address, data []byte, SendAmount *big.Int, gas *big.Int, existingTxs types.Transactions, revertTx []common.Hash, MaxBN uint64) (*types.SendBundleArgs, *sendBundle.BidCaseArg, types.Transactions) {
	userArg := utils.UserTx(pk, lock, data, gas, big.NewInt(conf.MedGasPrice))
	newTxs, revertTxHashes := sendBundle.GenerateBNBTxs(&userArg, SendAmount, userArg.Data, 1)
	revertTx = append(revertTxHashes, revertTx...)
	bundleArgs := utils.AddBundle(existingTxs, newTxs, revertTx, MaxBN)
	return bundleArgs, &userArg, newTxs
}

func UpdateUsrList(index int, txs types.Transactions, mined bool, result string) {
	UsrList[index].Txs = txs
	UsrList[index].Mined = mined
	UsrList[index].Rst = result
}

func UpdateUsrList6(index int, txs types.Transactions, mined bool, result string) {
	UsrList6[index].Txs = txs
	UsrList6[index].Mined = mined
	UsrList6[index].Rst = result
}

func SendBundles(t *testing.T, usr1Arg *sendBundle.BidCaseArg, usr2Arg *sendBundle.BidCaseArg, bundleArgs1 *types.SendBundleArgs, bundleArgs2 *types.SendBundleArgs) uint64 {
	Args[0] = usr1Arg
	Args[1] = usr2Arg
	BundleargsLsit[0] = bundleArgs1
	BundleargsLsit[1] = bundleArgs2
	return utils.ConcurSendBundles(t, Args, BundleargsLsit)
}
func SendBundlesTri(t *testing.T, usr1Arg *sendBundle.BidCaseArg, usr2Arg *sendBundle.BidCaseArg, usr3Arg *sendBundle.BidCaseArg, bundleArgs1 *types.SendBundleArgs, bundleArgs2 *types.SendBundleArgs, bundleArgs3 *types.SendBundleArgs) uint64 {
	Args3[0] = usr1Arg
	Args3[1] = usr2Arg
	Args3[2] = usr3Arg
	BundleargsLsit3[0] = bundleArgs1
	BundleargsLsit3[1] = bundleArgs2
	BundleargsLsit3[2] = bundleArgs3
	return utils.ConcurSendBundles(t, Args3, BundleargsLsit3)
}

func CheckTransactionIndex(t *testing.T, tx types.Transaction, expectedIndex string) {
	response := utils.GetTransactionReceipt(tx)
	txIndex := response.Result.TransactionIndex
	assert.Equal(t, txIndex, expectedIndex)
	log.Printf("Transaction %v index: %v,gasUsed %v", tx.Hash().Hex(), txIndex, response.Result.GasUsed)
}

func GetTxIndex(tx types.Transaction) string {
	response := utils.GetTransactionReceipt(tx)
	txIndex := response.Result.TransactionIndex
	log.Printf("Transaction %v index: %v,gasUsed %v", tx.Hash().Hex(), txIndex, response.Result.GasUsed)
	return txIndex
}
