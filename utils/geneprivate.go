package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func generatePrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func new() {
	numKeys := 5 // 生成私钥的数量

	for i := 0; i < numKeys; i++ {
		privateKey, err := generatePrivateKey()
		if err != nil {
			fmt.Println("Error generating private key:", err)
			return
		}

		privateKeyBytes := privateKey.D.Bytes()
		fmt.Printf("Private Key %d: %s\n", i+1, hex.EncodeToString(privateKeyBytes))
	}
}
