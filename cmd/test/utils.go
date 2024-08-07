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
	"math/big"
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
	req, err := http.NewRequest("POST", Endpoint, bytes.NewBuffer([]byte(dataString)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("curl --location '%v' --header 'Content-Type: application/json' --data '%s'\n", Endpoint, dataString)

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

type JSONResponse[T any] struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  T      `json:"result,omitempty"`
}

func ParseResponse[T any](response string) *T {
	var resp JSONResponse[*T]
	err := json.Unmarshal([]byte(response), &resp)
	if err != nil {
		fmt.Println(response)
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return resp.Result
}

func ParseResponseAsBigInt(response string) *big.Int {
	responseStr := ParseResponse[string](response)
	result, _ := big.NewInt(0).SetString(*responseStr, 0)
	return result
}

func SignTransaction(gethCfg *evm.GethConfig, privateKeyStr string, tx *types.Transaction) string {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatal(err)
	}

	signer := types.LatestSigner(gethCfg.ChainConfig)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	rawTxBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("0x%x", rawTxBytes)
}

func ToJsonString(v interface{}) string {
	byt, _ := json.Marshal(v)
	return string(byt)
}
