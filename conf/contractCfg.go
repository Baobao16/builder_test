package conf

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// WBNB：0xE5454b639B241c07Fc0d55b23690F9CeE18b7E4f
// RouterV2: 0xE1f45ef433b2ADF7583917974543a2df2161Dd6c
// Token: 0x429B2BEa55c0F2a30318d21D029EDc847977344F
var (
	Mylock    = common.HexToAddress("0xd324788cb01def44a4f89d775718e577bb4553c1") // 0514new
	WBNB      = common.HexToAddress("0xE5454b639B241c07Fc0d55b23690F9CeE18b7E4f")
	RouterV2  = common.HexToAddress("0xE1f45ef433b2ADF7583917974543a2df2161Dd6c")
	Token     = common.HexToAddress("0x429B2BEa55c0F2a30318d21D029EDc847977344F")
	SpecialOp = common.HexToAddress("0x05A9d51810475F47C914b97268ac53198dA89D68")
	ValueCp   = common.HexToAddress("0x035b6E463A445aF6d12Dbf9b2D0150c15Be5b357")

	// RcvAddress 贿赂地址
	RcvAddress   = common.HexToAddress("0x33Af2388136bf65b4b6413A1951391F89663c644")
	BribeAddress = common.HexToAddress("0x11c40ecf278CB259696b1f1E359f8682eE425522")
	SysAddress   = common.HexToAddress("0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE")
	// Coinbase = common.HexToAddress("0x7d83033eFaE53d3250cff2d9e39E4a63fdEd9712")

	WBNB_gas     = int64(50000)
	Min_gasPrice = int64(1e9)     // 小于1Gwei 的gasPrice会失败
	Max_gasLimit = int64(3000000) // 大于3000000 的gasLimit会失败
	SendA        = big.NewInt(0)
	High_gas     = big.NewInt(3000000)
	Med_gas      = big.NewInt(2000000)
	Low_gas      = big.NewInt(1000000)
	Lock_path    = "../../abi/ugLock.json"
	Spe_path     = "../../abi/specialOp.json"
	ValueCp_path = "../../abi/ValueCp.json"

	SpecialOp_Bb = common.Hex2Bytes("1c6dc3c0")
	SpecialOp_Cb = common.Hex2Bytes("e6f1f7510000000000000000000000007d83033efae53d3250cff2d9e39e4a63fded9712")
	// SpecialOp_Cb        = common.Hex2Bytes("e6f1f7510000000000000000000000000000000000000000000000000000000000000000")
	SpecialOp_ts        = common.Hex2Bytes("dd48b86d00000000000000000000000000000000000000000000000000000000661f7e59")
	SpecialOp_bh        = common.Hex2Bytes("e3533697146127548443eb5810584a021ee1b11893c267d85b88fd955f69c777f06ebe6a")
	TransferToken_code  = common.Hex2Bytes("2d339b1e000000000000000000000000429b2bea55c0f2a30318d21d029edc847977344f")
	TransferBNB_code    = common.Hex2Bytes("a6f9dae10000000000000000000000007b09bb26c9fef574ea980a33fc71c184405a4023")
	TransferWBNB_code   = common.Hex2Bytes("1a695230000000000000000000000000e5454b639b241c07fc0d55b23690f9cee18b7e4f") //"0x6e13" 28179
	TBalanceOfWBNB_code = common.Hex2Bytes("70a08231000000000000000000000000e5454b639b241c07fc0d55b23690f9cee18b7e4f")
	TotallysplWBNB_code = common.Hex2Bytes("18160ddd") // "0x533d"

	Txsucceed = "0x1"
	Txfailed  = "0x0"

	TxType = map[string]string{
		"Transfer":           "0x0",
		"Contract_Creation":  "0x1",
		"Contract_Invocatio": "0x2",
	}

	// Error msg
	// Tx error
	MissTx         = "bundle missing txs"
	InvalidTx      = "non-reverting tx in bundle failed"
	LargeTx        = "413 Request Entity Too Large: content length too large"
	TxCountLimit   = "only allow a maximum of 50 transactions"
	BundleConflict = "bundle already exist"
	// maxBlockNumber最多设为当前区块号+100
	MaxBlockNumberL = "the maxBlockNumber should not be lager than currentBlockNum + 100"
	MaxBlockNumberC = "maxBlockNumber should not be smaller than currentBlockNum"

	// maxTimestamp最多设为当前区块号+5minutes
	TimestampTop = "the minTimestamp/maxTimestamp should not be later than currentBlockTimestamp + 5 minutes"
	TimestampMM  = "the maxTimestamp should not be less than minTimestamp"
	TimestampMC  = "the maxTimestamp should not be less than currentBlockTimestamp"
)

// Contract:      common.HexToAddress("0x7b09bb26c9fef574ea980a33fc71c184405a4023"),
// Contract:   common.HexToAddress("0xb0b10B09780aa6A315158EF724404aa1497e9E6E"), // momo
