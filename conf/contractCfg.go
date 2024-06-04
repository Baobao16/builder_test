package conf

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// WBNB：0xE5454b639B241c07Fc0d55b23690F9CeE18b7E4f
// RouterV2: 0xE1f45ef433b2ADF7583917974543a2df2161Dd6c
// Token: 0x429B2BEa55c0F2a30318d21D029EDc847977344F
var (
	Mylock              = common.HexToAddress("0x87b2a2fed66eb5ad57d5c02483e812d731e44da7") // 0514new
	WBNB                = common.HexToAddress("0xE5454b639B241c07Fc0d55b23690F9CeE18b7E4f")
	RouterV2            = common.HexToAddress("0xE1f45ef433b2ADF7583917974543a2df2161Dd6c")
	Token               = common.HexToAddress("0x429B2BEa55c0F2a30318d21D029EDc847977344F")
	SpecialOp           = common.HexToAddress("0x05A9d51810475F47C914b97268ac53198dA89D68")
	ValueCp             = common.HexToAddress("0x035b6E463A445aF6d12Dbf9b2D0150c15Be5b357")
	TransferBNBCode     = common.Hex2Bytes("a6f9dae10000000000000000000000007b09bb26c9fef574ea980a33fc71c184405a4023")
	TBalanceOfWBNBCode  = common.Hex2Bytes("70a08231000000000000000000000000e5454b639b241c07fc0d55b23690f9cee18b7e4f")
	TotallySplWBNBCode  = common.Hex2Bytes("18160ddd")
	AllowanceRouterCode = common.Hex2Bytes("3e5beab9000000000000000000000000e1f45ef433b2adf7583917974543a2df2161dd6c")

	RcvAddress   = common.HexToAddress("0x33Af2388136bf65b4b6413A1951391F89663c644") //blkrz收账地址
	C48Address   = common.HexToAddress("0x12AE9700eD0C8BEC37162b8a6883d097C8AbEc34") //48club账户地址
	MidAddress   = common.HexToAddress("0x116D2f846ada0dBd7Fb5BEdc80BADA210b55B911") //中间账户地址
	BribeAddress = common.HexToAddress("0x11c40ecf278CB259696b1f1E359f8682eE425522")
	SysAddress   = common.HexToAddress("0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE")
	// Coinbase = common.HexToAddress("0x7d83033eFaE53d3250cff2d9e39E4a63fdEd9712")

	MinGasPrice = int64(1e9)  // 小于1Gwei 的gasPrice会失败
	MedGasPrice = int64(1e10) // 小于1Gwei 的gasPrice会失败
	SendA       = big.NewInt(0)
	HighGas     = big.NewInt(3e6)
	MedGas      = big.NewInt(2e6)
	LowGas      = big.NewInt(1e6)
	LockPath    = "../../abi/ugLock.json"
	SpePath     = "../../abi/specialOp.json"
	ValueCpPath = "../../abi/ValueCp.json"

	SpecialOpBb = common.Hex2Bytes("1c6dc3c0")
	SpecialOpCb = common.Hex2Bytes("e6f1f7510000000000000000000000007d83033efae53d3250cff2d9e39e4a63fded9712")
	// SpecialopCb        = common.Hex2Bytes("e6f1f7510000000000000000000000000000000000000000000000000000000000000000")
	TransferWBNBCode   = common.Hex2Bytes("1a695230000000000000000000000000e5454b639b241c07fc0d55b23690f9cee18b7e4f") //"0x6e13" 28179
	TbalanceOfWBNBCode = common.Hex2Bytes("70a08231000000000000000000000000e5454b639b241c07fc0d55b23690f9cee18b7e4f")

	TxSucceed = "0x1"
	TxFailed  = "0x0"

	TxType = map[string]string{
		"Transfer":            "0x0",
		"Contract_Creation":   "0x1",
		"Contract_Invocation": "0x2",
	}

	// MissTx Error msg
	MissTx         = "bundle missing txs"
	InvalidTx      = "non-reverting tx in bundle failed"
	LargeTx        = "413 Request Entity Too Large: content length too large"
	TxCountLimit   = "only allow a maximum of 50 transactions"
	BundleConflict = "bundle already exist"
	// MaxBlockNumberL maxBlockNumber最多设为当前区块号+100
	MaxBlockNumberL = "the maxBlockNumber should not be lager than currentBlockNum + 100"
	MaxBlockNumberC = "maxBlockNumber should not be smaller than currentBlockNum"

	// TimestampTop maxTimestamp最多设为当前区块号+5minutes
	TimestampTop = "the minTimestamp/maxTimestamp should not be later than currentBlockTimestamp + 5 minutes"
	TimestampMM  = "the maxTimestamp should not be less than minTimestamp"
	TimestampMC  = "the maxTimestamp should not be less than currentBlockTimestamp"
)

// Contract:      common.HexToAddress("0x7b09bb26c9fef574ea980a33fc71c184405a4023"),
// Contract:   common.HexToAddress("0xb0b10B09780aa6A315158EF724404aa1497e9E6E"), // momo
