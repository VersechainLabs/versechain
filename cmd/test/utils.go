package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"itachi/evm"
	"log"
	"net/http"
)

func GenerateRequestBody(method string, params ...interface{}) string {
	rpcId = rpcId + 1
	requestMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      rpcId,
		"method":  method,
	}

	if len(params) > 0 {
		requestMap["params"] = params
	} else {
		requestMap["params"] = []interface{}{}
	}

	requestBodyBytes, err := json.Marshal(requestMap)
	if err != nil {
		fmt.Println("Error marshalling the request body:", err)
		return ""
	}

	return string(requestBodyBytes)
}

func SendRequest(dataString string) string {
	req, err := http.NewRequest("POST", "http://localhost:9092", bytes.NewBuffer([]byte(dataString)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("curl --location 'localhost:9092' --header 'Content-Type: application/json' --data '%s'\n", dataString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	return string(body)
}

func GeneratePrivateKey() (string, string) {
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

type JSONResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

func ParseResponse(response string) string {
	var resp JSONResponse
	err := json.Unmarshal([]byte(response), &resp)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return resp.Result
}

func SignTransaction(gethCfg *evm.GethConfig, privateKeyStr string, tx *types.Transaction) string {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatal(err)
	}

	chainID := gethCfg.ChainConfig.ChainID
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	rawTxBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("0x%x", rawTxBytes)
}
