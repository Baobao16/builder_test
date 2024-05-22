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

type ResultB struct {
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
	JSONRPC string  `json:"jsonrpc"`
	ID      int     `json:"id"`
	Result  ResultB `json:"result"`
}

type ResponseNum struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

type ResponseBlk struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Difficulty       string            `json:"difficulty"`
		ExtraData        string            `json:"extraData"`
		GasLimit         string            `json:"gasLimit"`
		GasUsed          string            `json:"gasUsed"`
		Hash             string            `json:"hash"`
		LogsBloom        string            `json:"logsBloom"`
		Miner            string            `json:"miner"`
		MixHash          string            `json:"mixHash"`
		Nonce            string            `json:"nonce"`
		Number           string            `json:"number"`
		ParentHash       string            `json:"parentHash"`
		ReceiptsRoot     string            `json:"receiptsRoot"`
		Sha3Uncles       string            `json:"sha3Uncles"`
		Size             string            `json:"size"`
		StateRoot        string            `json:"stateRoot"`
		Timestamp        string            `json:"timestamp"`
		TotalDifficulty  string            `json:"totalDifficulty"`
		Transactions     []json.RawMessage `json:"transactions"`
		TransactionsRoot string            `json:"transactionsRoot"`
		Uncles           []string          `json:"uncles"`
	}
}

func sendRequest(url string, payload *strings.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	return res, nil
}

func GetLatestBlockNumber() string {
	payload := strings.NewReader("{\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"eth_blockNumber\"}")
	res, err := sendRequest(conf.Url, payload)
	if err != nil {
		log.Fatalf("Error in sending request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var response ResponseNum
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("GetLatestBlockNumber Failed to decode JSON: %v", err)
	}

	latestBN, err := strconv.ParseInt(strings.TrimPrefix(response.Result, "0x"), 16, 64)
	if err != nil {
		log.Fatalf("Failed to parse block number: %v", err)
	}
	log.Printf("Current block height is: %v (%v)", response.Result, latestBN)

	return response.Result
}

func GetBlockMsg(arg string) string {
	blkNum := GetLatestBlockNumber()
	payload := fmt.Sprintf(`{"id":1,"jsonrpc":"2.0","params":["%s",true],"method":"eth_getBlockByNumber"}`, blkNum)
	fmt.Printf("%v", payload)

	res, err := sendRequest(conf.Url, strings.NewReader(payload))
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	//fmt.Printf("%v", body)
	var response ResponseBlk
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		log.Fatalf("GetBlockMsg Failed to decode JSON: %v", err)
	}

	switch arg {
	case "testTimestamp", "testTimestampEq":
		return response.Result.Timestamp
	case "testBlockHash":
		return response.Result.Hash
	default:
		log.Printf("Return Block Number by default")
		return response.Result.Number
	}
}

func GetTransactionReceipt(tx types.Transaction) Response {
	payload := strings.NewReader(fmt.Sprintf(`{"id":1,"jsonrpc":"2.0","params":["%s"],"method":"eth_getTransactionReceipt"}`, tx.Hash().Hex()))
	res, err := sendRequest(conf.Url, payload)
	if err != nil {
		log.Fatalf("Error in sending request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("GetTransactionReceipt Failed to decode JSON: %v", err)
	}

	return response
}

func WaitMined(txs types.Transactions, cbn uint64) {
	log.Printf("[WaitMined] Current block number: %v", cbn)
	waitTime := 5 * time.Second

	for _, tx := range txs {
		start := time.Now()
		for {
			response := GetTransactionReceipt(*tx)

			if response.Result.BlockHash != "" {
				log.Printf("[WaitMined] Transaction %v included in block: %v", tx.Hash().Hex(), response.Result.BlockNumber)
				break
			} else if time.Since(start) >= waitTime {
				log.Printf("[WaitMined] Timeout: Transaction %v not included in a block within %v seconds", tx.Hash().Hex(), waitTime.Seconds())
				break
			} else {
				log.Printf("[WaitMined] Transaction %v not included in a block yet. Waiting...", tx.Hash().Hex())
				time.Sleep(1 * time.Second) // 短暂休眠后重试
			}
		}
	}
}

func TxInSameBlk(blk []string) bool {
	if len(blk) <= 1 {
		return true // 如果列表为空或只有一个元素，则认为全部相同
	}
	firstElement := blk[0]
	for _, str := range blk {
		log.Printf("blk %v", str)
		if str != firstElement {
			return false // 如果有任何一个元素与第一个元素不相同，则返回false
		}
	}
	return true // 如果所有元素都与第一个元素相同，则返回true
}

func CheckBundleTx(t *testing.T, tx types.Transaction, valid bool, status string) string {
	log.Printf("Check Tx :%v \n", tx.Hash().Hex())
	response := GetTransactionReceipt(tx)

	if valid {
		log.Printf("Transaction executed %v transactionIndex is: 【 %v 】", response.Result.Status, response.Result.TransactionIndex)
		if status == conf.Txsucceed {
			log.Printf("%v gasUsed %v \n", tx.Hash().Hex(), response.Result.GasUsed)
		} else if response.Result.Status == conf.Txsucceed {
			log.Printf("Transaction %v blockNumber is :%v ", tx.Hash().Hex(), response.Result.BlockNumber)
			log.Printf("Transaction %v [expected] failed to be mined gasUsed %v transactionIndex is: 【 %v 】  \n ", tx.Hash().Hex(), response.Result.GasUsed, response.Result.TransactionIndex)
		}
		assert.Equal(t, status, response.Result.Status, "Tx %v tx_Status wrong , [expected] status is  %v  ", tx.Hash().Hex(), status)
		return response.Result.BlockNumber
	} else {
		assert.Empty(t, response.Result.Type, "Transaction %v [expected] failed to be mined gasUsed %v transactionIndex is: 【 %v 】 \n ", tx.Hash().Hex(), response.Result.GasUsed, response.Result.TransactionIndex)
		log.Printf("Transaction %v failed to be mined ", tx.Hash().Hex())
		return "Transaction failed to be mined"
	}
}

func BlockHeightIncreased(t *testing.T) {
	// 等待链上高度增加
	response1 := GetLatestBlockNumber()
	t.Log("Waiting for Bundle committed \n")
	time.Sleep(6 * time.Second)
	for {
		response := GetLatestBlockNumber()
		if response > response1 {
			log.Println("current height increased")
			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}

}

func JsonToMap(jsonStr string) map[string]interface{} {
	// json转map
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &mapResult)
	if err != nil {
		log.Printf("JsonToMap err:", err)
	}
	return mapResult
}

func GetLatestBlkMsg(t *testing.T, path string, arg string, add uint64) []byte {
	// 根据最新区块信息调用SpecialOp合约方法
	myABI := GeneABI(path)
	var vl interface{}
	blkHash := GetBlockMsg(arg)
	log.Printf("初始值：%v", blkHash)
	if add != 0 {
		decNum, _ := strconv.ParseUint(strings.TrimPrefix(blkHash, "0x"), 16, 64)
		decNum += add
		blkHash = fmt.Sprintf("%x", decNum)
		log.Printf("加 上add %v 后的结果：%v\n", add, blkHash)
		var bigInt big.Int
		_, success := bigInt.SetString(blkHash, 16)
		if !success {
			log.Printf("解析失败")
		}
		vl = &bigInt

	} else {
		vl = common.HexToHash(blkHash)
	}
	// 编码函数调用数据
	encodedData, err := myABI.Pack(arg, vl)
	if err != nil {
		log.Printf("Error encoding data:", err)
	}

	// 打印编码数据
	log.Printf("Encoded data: %x\n", encodedData)
	return encodedData
}

func GeneABI(path string) *abi.ABI {
	// 读取合约ABI文件
	abiData, err := os.ReadFile(path) //, method string, args ...interface{}
	if err != nil {
		log.Printf("无法读取 ABI 文件:", err)
	}
	abiJSON := string(abiData)

	myABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		log.Printf("解析 ABI 数据时出错:", err)
	}

	return &myABI
}
