package utils

import (
	"fmt"
	"time"
)

func main() {
	// 获取当前时间戳（自1970年1月1日以来的秒数）
	currentTime := time.Now().Unix()

	// 输出当前时间戳
	fmt.Println("当前时间戳：", currentTime)
}
