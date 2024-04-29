package newtestcases

import (
	"github.com/ethereum/go-ethereum/common"
)

// WBNB：0xE5454b639B241c07Fc0d55b23690F9CeE18b7E4f
// RouterV2: 0xE1f45ef433b2ADF7583917974543a2df2161Dd6c
// Token: 0x429B2BEa55c0F2a30318d21D029EDc847977344F

var (
	Lock     = common.HexToAddress("0x7EF2e0048f5bAeDe046f6BF797943daF4ED8CB47")
	Mylock   = common.HexToAddress("0xd9145CCE52D386f254917e481eB44e9943F39138")
	WBNB     = common.HexToAddress("0xE5454b639B241c07Fc0d55b23690F9CeE18b7E4f")
	RouterV2 = common.HexToAddress("0xE1f45ef433b2ADF7583917974543a2df2161Dd6c")
	Token    = common.HexToAddress("0x429B2BEa55c0F2a30318d21D029EDc847977344F")
	// test_Token = common.HexToAddress("0x199e3Bfb54f4aAa9D67d1BB56429c5ef9D1A2A91")

	WBNB_gas = int64(50000)
	Test_gas = int64(210000)

	Lock_set_code       = common.Hex2Bytes("2801617e0000000000000000000000007ef2e0048f5baede046f6bf797943daf4ed8cb47")
	Lock_lock_code      = common.Hex2Bytes("f435f5a70000000000000000000000007ef2e0048f5baede046f6bf797943daf4ed8cb47")
	Lock_fakeshort_code = common.Hex2Bytes("992222fb0000000000000000000000007ef2e0048f5baede046f6bf797943daf4ed8cb47000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000036161610000000000000000000000000000000000000000000000000000000000")
	Lock_fakelong_code  = common.Hex2Bytes("992222fb0000000000000000000000007ef2e0048f5baede046f6bf797943daf4ed8cb47000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000146161616161616161616161616161616161616161000000000000000000000000")
	Lock_fakel_code     = common.Hex2Bytes("6da565f90000000000000000000000007ef2e0048f5baede046f6bf797943daf4ed8cb47")
	Mylock_lock_code    = common.Hex2Bytes("0df97361000000000000000000000000d9145cce52d386f254917e481eb44e9943f39138") //"0x53b8" lock 21432
	Mylock_unlock_code  = common.Hex2Bytes("600d647a000000000000000000000000d9145cce52d386f254917e481eb44e9943f39138") //"0x53b8" unlock 21432
	TransferToken_code  = common.Hex2Bytes("2d339b1e000000000000000000000000429b2bea55c0f2a30318d21d029edc847977344f")
	TransferBNB_code    = common.Hex2Bytes("a6f9dae10000000000000000000000007b09bb26c9fef574ea980a33fc71c184405a4023")
	TransferWBNB_code   = common.Hex2Bytes("1a695230000000000000000000000000e5454b639b241c07fc0d55b23690f9cee18b7e4f") //"0x6e13" 28179
	TBalanceOfWBNB_code = common.Hex2Bytes("70a08231000000000000000000000000e5454b639b241c07fc0d55b23690f9cee18b7e4f")
	TotallysplWBNB_code = common.Hex2Bytes("18160ddd") // "0x533d"

	Txsucceed = "0x1"
	Txfailed  = "0x0"

	txType = map[string]string{
		"Transfer":           "0x0",
		"Contract_Creation":  "0x1",
		"Contract_Invocatio": "0x2",
	}

	// Error msg

	MissTx    = "bundle missing txs"
	InvalidTx = "no valid sim result"
	LargeTx   = "413 Request Entity Too Large: content length too large"

	BundleConflict = "bundle already exist"
	// maxBlockNumber最多设为当前区块号+100
	maxBlockNumberL = "the maxBlockNumber should not be lager than currentBlockNum + 100"
	maxBlockNumberC = "maxBlockNumber should not be smaller than currentBlockNum"

	// maxTimestamp最多设为当前区块号+5minutes
	TimestampTop = "the minTimestamp/maxTimestamp should not be later than currentBlockTimestamp + 5 minutes"
	TimestampMM  = "the maxTimestamp should not be less than minTimestamp"
	TimestampMC  = "the maxTimestamp should not be less than currentBlockTimestamp"
)

// Contract:      common.HexToAddress("0x7b09bb26c9fef574ea980a33fc71c184405a4023"),
// Contract:   common.HexToAddress("0xb0b10B09780aa6A315158EF724404aa1497e9E6E"), // momo
