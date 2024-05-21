package simulate

// bundleArgs := &SendBundleArgs{
// 	//MaxBlockNumber:    9,
// 	Txs:               txBytes,
// 	RevertingTxHashes: revertTxHashes,
// 	SimXYZ:            true,
// }

// bidJson, _ := json.MarshalIndent(bundleArgs, "", "  ")
// println(string(bidJson))

// err := arg.BuilderClient.Client().CallContext(arg.Ctx, nil, "eth_sendBundle", bundleArgs) //替换sendBundle  返回的是bundle哈希
// simulate Bundle 接口并发测试
// 预期 nonce 不生效，均能发送成功，不上链
