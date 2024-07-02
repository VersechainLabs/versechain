package pkg

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func generatePrivateKey() (string, string) {
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return "", ""
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return hexutil.Encode(privateKeyBytes)[2:], address
}

func sendRequest(hostAddress string, dataString string) error {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s", hostAddress), bytes.NewBuffer([]byte(dataString)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("curl --location 'localhost:9092' --header 'Content-Type: application/json' --data '%s'\n", dataString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code err: %v", resp.StatusCode)
	}
	return nil
}
