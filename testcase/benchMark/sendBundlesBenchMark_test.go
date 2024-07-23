package benchMark

import (
	"github.com/xkwang/conf"
	"github.com/xkwang/sendBundle"
	"github.com/xkwang/utils"
	"log"
	"math/big"
	"testing"
)

// 压力测试函数
// 持续发送bundles
func BenchmarkSendBundles(b *testing.B) {
	arg := utils.UserTx(conf.RootPk, conf.WBNB, conf.TransferWBNBCode, conf.MedGasLimit, big.NewInt(conf.MinGasPrice))
	// 循环执行测试函数
	for i := 0; i < b.N; i++ {
		// 在每次迭代中调用接口
		log.Println("run case")
		txs, err := sendBundle.ValidBundle_NilPayBidTx_2(&arg, true)
		if err != nil {
			log.Println(" failed: ", err.Error())
		} else {
			log.Println("ValidBundle_NilPayBidTx_1 succeed ")
		}
		println(txs)
		if err != nil {
			b.Fatalf("call failed: %v", err)
		}
	}
}
