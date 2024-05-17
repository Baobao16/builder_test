package testcase

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xkwang/cases"
	"github.com/xkwang/conf"
	"github.com/xkwang/utils"
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
	Args            = make([]*cases.BidCaseArg, 2)
	UsrList         = make([]utils.TxStatus, 3)
	BundleArgs_lsit = make([]*types.SendBundleArgs, 2)
)
