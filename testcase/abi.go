package testcase

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xkwang/conf"
	"github.com/xkwang/sendBundle"
	"github.com/xkwang/utils"
	"math/big"
	"testing"
)

var LockABIfile = utils.GeneABI(conf.Lock_path)

// step-1 读取abi文件
var LockABI = utils.Contract{ABI: *LockABIfile}
var ResetData = utils.GeneEncodedData(LockABI, "reset")
var LockData = utils.GeneEncodedData(LockABI, "lock", 1, true) // "gasUsed":"0x6ab1"
var LockFData = utils.GeneEncodedData(LockABI, "lock", 1, false)

var UnlockStrData = utils.GeneEncodedData(LockABI, "unlock", 1, "str")
var UnlockMoreData = utils.GeneEncodedData(LockABI, "unlock", 1, "more")
var UnlockDeData = utils.GeneEncodedData(LockABI, "unlock_de", 1, "more") //"0x1d58c7"1923271
var UnlockDesData = utils.GeneEncodedData(LockABI, "unlock_de", 1, "str")
var FakelockStrData = utils.GeneEncodedData(LockABI, "fakelock", 1, "str")
var FakelockMoreData = utils.GeneEncodedData(LockABI, "fakelock", 1, "more")
var (
	Args           = make([]*sendBundle.BidCaseArg, 2)
	UsrList        = make([]utils.TxStatus, 3)
	BundleargsLsit = make([]*types.SendBundleArgs, 2)
)

func AddUserBundle(pk string, lock common.Address, data []byte, SendAmount *big.Int, gas *big.Int, existingTxs types.Transactions, revertTx []common.Hash, MaxBN uint64) (*types.SendBundleArgs, *sendBundle.BidCaseArg, types.Transactions) {
	userArg := utils.UserTx(pk, lock, data, gas)
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

func SendBundles(t *testing.T, usr1Arg *sendBundle.BidCaseArg, usr2Arg *sendBundle.BidCaseArg, bundleArgs1 *types.SendBundleArgs, bundleArgs2 *types.SendBundleArgs) uint64 {
	Args[0] = usr1Arg
	Args[1] = usr2Arg
	BundleargsLsit[0] = bundleArgs1
	BundleargsLsit[1] = bundleArgs2
	return utils.ConcurSendBundles(t, Args, BundleargsLsit)
}
