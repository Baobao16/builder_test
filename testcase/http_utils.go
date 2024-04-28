package newtestcases

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
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
	} `json:"result"`
}

type Response_num struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

// 获取当前最新区块高度
func getLatestBlockNumber() int64 {
	payload := strings.NewReader("{\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"eth_blockNumber\"}")
	req, _ := http.NewRequest("POST", url, payload)

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
		return 0
	} else {
		// fmt.Printf("变量 x 的类型：%T\n", hex.EncodeToString([]byte(response.Result)))
		latsetBN, _ := strconv.ParseInt(strings.TrimPrefix(response.Result, "0x"), 16, 64)
		log.Printf("current height is : %v", latsetBN)

		return latsetBN
	}
}

// 查询交易回执
func getTransactionReceipt(tx types.Transaction) []byte {
	// 根据交易哈希查询回执

	payload := strings.NewReader("{\"id\":1,\"jsonrpc\":\"2.0\",\"params\":[\"" + string(tx.Hash().Hex()) + "\"],\"method\":\"eth_getTransactionReceipt\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	// log.Println(string(body))
	return body
}

// 确认bundle已包含交易
func checkBundleTx(t *testing.T, tx types.Transaction, valid bool, status string, tx_type string) {
	body := getTransactionReceipt(tx)
	// 将 JSON 字符串解码到 Response 结构体中
	var response Response
	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("decode JSON failed:", err)
		return
	}
	if valid {
		// t.Log("查询交易已上链")
		log.Printf("Transaction %v mined:", tx.Hash().Hex())
		// 校验交易类型
		if status == Txsucceed {
			t.Logf("Type: %v Transaction executed successfully", tx_type)
		} else if status == Txfailed {
			t.Logf("Type: %v Transaction executed failed", tx_type)
		} else {
			t.Logf("Type :%v Transaction executed error!", tx_type)
		}
		assert.Equal(t, txType[tx_type], response.Result.Type)
		// 校验交易执行状态
		assert.Equal(t, status, response.Result.Status)

	} else {
		// log.Println("查询交易未上链")
		//"removed" 字段为 false 表示该日志未被移除，仍然有效。如果为 true，则表示该日志已被移除或撤销，
		log.Printf("Transaction %v failed to be mined ", tx.Hash().Hex())
		assert.Empty(t, response.Result.Type)
	}

}

// 等待链上高度增加
func BlockheightIncreased(t *testing.T) {
	response_1 := getLatestBlockNumber()
	t.Log("Waiting for Bundle commited \n")
	time.Sleep(6 * time.Second)
	for {
		response := getLatestBlockNumber()
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
