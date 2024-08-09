package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gorilla/websocket"
	"itachi/evm/ethrpc"
	MEVless "itachi/mev-less"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	mevlessWsAddr   = "localhost:9071"
	mevlessHttpAddr = "http://localhost:7999/api/writing"
)

type Response struct {
	BlockNumber int                    `json:"block_number"`
	Sequences   map[string]interface{} `json:"sequences"`
}

func mevlessTest() {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			// 1. Connect Websocket
			u := url.URL{Scheme: "ws", Host: mevlessWsAddr, Path: "/mev_less"}
			c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Fatal("dial:", err)
			}
			defer c.Close()

			// 2. Generate Request
			_, randomAddressStr := GeneratePrivateKey()
			randomAddress := common.HexToAddress(randomAddressStr)
			value := hexutil.Big(*ether)
			gasArgs := ethrpc.TransactionArgs{
				From:  &testWalletAddr,
				To:    &randomAddress,
				Value: &value,
			}
			gasLimit := estimateGas(&gasArgs)
			nonce := getNonce()
			gasPrice := getGasPrice()

			tx := types.NewTx(&types.LegacyTx{
				Nonce:    nonce,
				GasPrice: gasPrice,
				Gas:      gasLimit,
				To:       gasArgs.To,
				Value:    gasArgs.Value.ToInt(),
			})
			rawTx, signedTx := SignTransaction(gethCfg, testWalletPrivateKeyStr, tx)

			mevless_params := fmt.Sprintf("%s%s", MEVless.Prefix, signedTx.Hash().Hex())

			// 3. Send Order Request
			requestBody, err := json.Marshal(map[string]interface{}{
				"call": map[string]string{
					"tripod_name": "mevless",
					"func_name":   "OrderTx",
					"params":      mevless_params,
				},
			})

			fmt.Printf("mevless request mevless_params %s\n", mevless_params)
			fmt.Printf("mevless request requestBody %s\n", requestBody)
			resp, err := http.Post(mevlessHttpAddr, "application/json", bytes.NewBuffer(requestBody))
			defer resp.Body.Close()

			// 4. listen ws message
			_, message, err := c.ReadMessage()
			var response Response
			if err := json.Unmarshal(message, &response); err != nil {
				log.Fatalf("Error unmarshalling WS response: %v", err)
			}
			for num, hash := range response.Sequences {
				if hash == signedTx.Hash().Hex() {
					fmt.Printf("Client %d's request hash %s order at %v\n", clientID, hash, num)
				}
			}

			sleepDuration := time.Millisecond * time.Duration(rnd.Intn(200)+1)
			time.Sleep(sleepDuration)

			sendTxRequestBody := GenerateRequestBody("eth_sendRawTransaction", rawTx)
			_ = SendRequest(sendTxRequestBody)

		}(i + 1)
	}
	wg.Wait()
}
