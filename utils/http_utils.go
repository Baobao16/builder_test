package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/xkwang/conf"
)

type Result_b struct {
	BlockHash         string        `json:"blockHash"`
	BlockNumber       string        `json:"blockNumber"`
	ContractAddress   string        `json:"contractAddress"`
	CumulativeGasUsed string        `json:"cumulativeGasUsed"`
	EffectiveGasPrice string        `json:"effectiveGasPrice"`
	From              string        `json:"from"`
	GasUsed           string        `json:"gasUsed"`
	Logs              []interface{} `json:"logs"`
	LogsBloom         string        `json:"logsBloom"`
	Status            string        `json:"status"`
	To                string        `json:"to"`
	TransactionHash   string        `json:"transactionHash"`
	TransactionIndex  string        `json:"transactionIndex"`
	Type              string        `json:"type"`
}

type Response struct {
	JSONRPC string   `json:"jsonrpc"`
	ID      int      `json:"id"`
	Result  Result_b `json:"result"`
}

type Response_num struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

type Response_blk struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Difficulty       string   `json:"difficulty"`
		ExtraData        string   `json:"extraData"`
		GasLimit         string   `json:"gasLimit"`
		GasUsed          string   `json:"gasUsed"`
		Hash             string   `json:"hash"`
		LogsBloom        string   `json:"logsBloom"`
		Miner            string   `json:"miner"`
		MixHash          string   `json:"mixHash"`
		Nonce            string   `json:"nonce"`
		Number           string   `json:"number"`
		ParentHash       string   `json:"parentHash"`
		ReceiptsRoot     string   `json:"receiptsRoot"`
		Sha3Uncles       string   `json:"sha3Uncles"`
		Size             string   `json:"size"`
		StateRoot        string   `json:"stateRoot"`
		Timestamp        string   `json:"timestamp"`
		TotalDifficulty  string   `json:"totalDifficulty"`
		Transactions     []string `json:"transactions"`
		TransactionsRoot string   `json:"transactionsRoot"`
		Uncles           []string `json:"uncles"`
	}
}

// 根据区块号获取区块hash
func GetBlockMsg(arg string) string {
	blknum := GetLatestBlockNumber()
	payload := strings.NewReader("{\"id\":1,\"jsonrpc\":\"2.0\",\"params\":[\"" + blknum + "\",true],\"method\":\"eth_getBlockByNumber\"}")
	req, _ := http.NewRequest("POST", conf.Url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	// log.Println(string(body))

	var response1 Response_blk
	err := json.Unmarshal([]byte(body), &response1)
	if err != nil {
		fmt.Println("decode JSON Failed:", err)
	}
	if arg == "testTimestamp" {
		return response1.Result.Timestamp

	} else if arg == "testTimestampEq" {
		return response1.Result.Timestamp

	} else if arg == "testBlockHash" {
		return response1.Result.Hash

	} else {
		log.Printf("Return Block Number by Default")
		return response1.Result.Number
	}

}

// 获取当前最新区块高度
func GetLatestBlockNumber() string {
	payload := strings.NewReader("{\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"eth_blockNumber\"}")
	req, _ := http.NewRequest("POST", conf.Url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	// log.Println(string(body))

	// 将 JSON 字符串解码到 Response 结构体中
	var response Response_num
	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("decode JSON Failed:", err)
		return ""
	} else {
		// fmt.Printf("变量 x 的类型：%T\n", hex.EncodeToString([]byte(response.Result)))
		latsetBN, _ := strconv.ParseInt(strings.TrimPrefix(response.Result, "0x"), 16, 64)
		log.Printf("%v current height is : %v", response.Result, latsetBN)

		return response.Result
	}
}

// 查询交易回执
func GetTransactionReceipt(tx types.Transaction) Response {
	// 根据交易哈希查询回执

	payload := strings.NewReader("{\"id\":1,\"jsonrpc\":\"2.0\",\"params\":[\"" + string(tx.Hash().Hex()) + "\"],\"method\":\"eth_getTransactionReceipt\"}")

	req, _ := http.NewRequest("POST", conf.Url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	// log.Println(string(body))
	// return body
	var response Response
	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("decode JSON failed:", err)
		// return "decode JSON failed"
	}
	return response
}

func WaitMined(txs types.Transactions, cbn uint64) {
	log.Printf("[WaitMined] current blk_num : %v", cbn)
	var response Response
	for _, tx := range txs {
		response = GetTransactionReceipt(*tx)
	}
	var ept bool
	if response.Result.BlockHash == "" {
		ept = false
	} else {
		ept = true
	}

	log.Printf("[WaitMined] tx included blk_num is: %v", ept)
	if ept {
		log.Printf("[WaitMined] tx included blk_num : %v", cbn)

	} else {
		time.Sleep(5 * time.Second)

		log.Printf("[WaitMined] tx included blk_num : %v", cbn+1)
	}

}

func TxinSameBlk(blk []string) bool {
	if len(blk) <= 1 {
		return true // 如果列表为空或只有一个元素，则认为全部相同
	}
	firstElement := blk[0]
	for _, str := range blk {
		if str != firstElement {
			return false // 如果有任何一个元素与第一个元素不相同，则返回false
		}
	}
	return true // 如果所有元素都与第一个元素相同，则返回true

}

// 确认bundle已包含交易
func CheckBundleTx(t *testing.T, tx types.Transaction, valid bool, status string) string {
	log.Printf("Check Tx :%v \n", tx.Hash().Hex())
	response := GetTransactionReceipt(tx)

	if valid {
		// 预期会上链的交易
		log.Printf("Type: %v Transaction executed %v transactionIndex is: 【 %v 】", response.Result.Status, response.Result.TransactionIndex)
		if status == conf.Txsucceed {
			// 打印gasUsed
			log.Printf("%v gasUsed %v \n", tx.Hash().Hex(), response.Result.GasUsed)
		} else if response.Result.Status == conf.Txsucceed {
			log.Printf("Transaction %v blockNumber is :%v ", tx.Hash().Hex(), response.Result.BlockNumber)
			log.Printf("Transaction %v [expected] failed to be mined gasUsed %v transactionIndex is: 【 %v 】  \n ", tx.Hash().Hex(), response.Result.GasUsed, response.Result.TransactionIndex)
		}
		// 校验交易执行状态
		assert.Equal(t, status, response.Result.Status, "Tx %v tx_Status wrong , [expected] status is  %v  ", tx.Hash().Hex(), status)
		return response.Result.BlockNumber
	} else {
		// 预期不上链的交易
		//"removed" 字段为 false 表示该日志未被移除，仍然有效。如果为 true，则表示该日志已被移除或撤销，
		assert.Empty(t, response.Result.Type, "Transaction %v [expected] failed to be mined gasUsed %v transactionIndex is: 【 %v 】 \n ", tx.Hash().Hex(), response.Result.GasUsed, response.Result.TransactionIndex)
		log.Printf("Transaction %v failed to be mined ", tx.Hash().Hex())
		return "Transaction failed to be mined"
	}
	// if response.Result.Type == "" {
	// 	log.Printf("Transaction %v failed to be mined ", tx.Hash().Hex())
	// 	panic("Transaction [expected] success")
	// 	return "Transaction failed to be mined"
	// } else

}

// 等待链上高度增加
func BlockheightIncreased(t *testing.T) {
	response_1 := GetLatestBlockNumber()
	t.Log("Waiting for Bundle commited \n")
	time.Sleep(6 * time.Second)
	for {
		response := GetLatestBlockNumber()
		if response > response_1 {
			log.Println("current height increased")
			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}

}

// json转map
func JsonToMap(jsonStr string) map[string]interface{} {
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &mapResult)
	if err != nil {
		fmt.Println("JsonToMap err:", err)
	}
	return mapResult
}

// 根据最新区块信息调用SpecialOp合约方法
func GetLatestBlkMsg(t *testing.T, path string, arg string, add uint64) []byte {
	myABI := GeneABI(path)

	var vl interface{}

	blkHash := GetBlockMsg(arg)
	fmt.Println("初始值：", blkHash)
	if add != 0 {
		decNum, _ := strconv.ParseUint(strings.TrimPrefix(blkHash, "0x"), 16, 64)
		decNum += add
		blkHash = fmt.Sprintf("%x", decNum)
		fmt.Printf("加 上add %v 后的结果：%v\n", add, blkHash)
		var bigInt big.Int
		_, success := bigInt.SetString(blkHash, 16)
		if !success {
			fmt.Println("解析失败")
		}
		// fmt.Println("解析值：", &bigInt)
		vl = &bigInt

	} else {
		vl = common.HexToHash(blkHash)
	}
	// 编码函数调用数据
	encodedData, err := myABI.Pack(arg, vl)
	if err != nil {
		fmt.Println("Error encoding data:", err)
	}

	// 打印编码数据
	fmt.Printf("Encoded data: %x\n", encodedData)
	return encodedData
}

// 通用：根据ABI文件生成Encode
// func CreateContractEncoded(myABI abi.ABI, method string) {

// }

// 读取合约ABI文件
func GeneABI(path string) *abi.ABI {
	abiData, err := os.ReadFile(path) //, method string, args ...interface{}
	if err != nil {
		fmt.Println("无法读取 ABI 文件:", err)
	}
	abiJSON := string(abiData)

	myABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		fmt.Println("解析 ABI 数据时出错:", err)
	}

	return &myABI
}
