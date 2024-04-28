package new

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func new() {

	url := "http://10.2.66.75:28545"

	payload := strings.NewReader("{\"id\":1,\"jsonrpc\":\"2.0\",\"params\":[\"0xaf557aae059cd9df66099425d7814d93fe025e8259cebb301cfd45a1d59cd306\"],\"method\":\"eth_getTransactionReceipt\"}")
	fmt.Println(payload)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	// fmt.Println(string(body))
	// fmt.Printf("变量 body 的类型：%T\n", body)

	// 将 JSON 字符串解码到 Response 结构体中
	var response Response
	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("解码 JSON 数据失败:", err)
		return
	}

	// 打印指定字段的值
	fmt.Println("blockHash 字段的值:", response.Result.BlockHash)

	// var data map[string]interface{}
	// err := json.NewDecoder(res.Body).Decode(&data)
	// if err != nil {
	// 	fmt.Println("解析响应失败:", err)
	// 	return
	// }

	// // 打印解析后的 JSON 数据
	// fmt.Println("解析后的 JSON 数据:")
	// for key, value := range data {
	// 	fmt.Printf("%s: %v\n", key, value)
	// }

}
