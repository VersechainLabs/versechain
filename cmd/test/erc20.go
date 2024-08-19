package main

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	abigen "itachi/cmd/test/abi"
	"itachi/evm/ethrpc"
	"log"
	"math/big"
	"os"
	"time"
)

func testErc20DeployAndUse() {
	//_, contractAddr := createContract()
	contractAddr := common.HexToAddress("0x310b8685e3e69cb05b251a12f5ffab23001cda42")
	mintErc20(contractAddr)
	//sendEvent(contractAddr)
	//sendEmptyEvent(contractAddr)

	//testReadCache(contractAddr)

	//readBalance(contractAddr)

}
func testReadCache(contractAddr common.Address) {
	hasChange := false
	prevBlock := uint64(0)
	for i := 0; i < 4; i++ {
		currentBlock := readBlockNumber(contractAddr)
		if prevBlock == 0 {
			prevBlock = currentBlock
		} else if prevBlock != currentBlock {
			hasChange = true
		}

		time.Sleep(5 * time.Second)
	}

	log.Printf("BlockNumber changed = %v", hasChange)

	writeCounter(contractAddr, 100)
	currentBlock := readBlockNumber(contractAddr)
	log.Printf("BlockNumber changed = %v", currentBlock != prevBlock)

	//readCounter(contractAddr)
	//readCounter(contractAddr)
}

func createContract() (txHash string, contractAddr common.Address) {
	log.Println("================================")
	log.Println("Create Contract")
	log.Println("================================")
	contractBinByte, err := os.ReadFile("./cmd/test/abi/TestErc20.bin")
	if err != nil {
		log.Fatal("failed to load bin file: ", err)
	}

	// estimate gas
	data, err := hex.DecodeString(string(contractBinByte))
	if err != nil {
		log.Fatal(err)
	}

	gasLimit := estimateGas(data)
	nonce := getNonce()
	gasPrice := getGasPrice()

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       nil,
		Data:     data,
	})
	rawTx := SignTransaction(gethCfg, testWalletPrivateKeyStr, tx)
	requestBody := GenerateRequestBody("eth_sendRawTransaction", rawTx)
	response := SendRequest(requestBody)

	txHash = *ParseResponse[string](response)
	log.Printf("create contract's txHash: %s", txHash)
	receiptChan := waitForReceipt(txHash)
	receipt := <-receiptChan
	contractAddr = receipt.ContractAddress

	return
}

func mintErc20(contractAddr common.Address) {
	log.Println("================================")
	log.Println("Mint Erc20")
	log.Println("================================")
	abi, _ := abigen.TestErc20MetaData.GetAbi()
	data, _ := abi.Pack("mint", testWalletAddr, ether_100)

	dataHex := hexutil.Bytes(data)
	arg := ethrpc.TransactionArgs{
		From: &testWalletAddr,
		To:   &contractAddr,
		Data: &dataHex,
	}
	estimateGasBody := GenerateRequestBody("eth_estimateGas", arg)
	estimateGasResponse := SendRequest(estimateGasBody)
	gasLimit := ParseResponseAsBigInt(estimateGasResponse).Uint64()

	nonce := getNonce()
	gasPrice := getGasPrice()
	to := contractAddr

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Data:     data,
	})

	rawTx := SignTransaction(gethCfg, testWalletPrivateKeyStr, tx)
	requestBody := GenerateRequestBody("eth_sendRawTransaction", rawTx)
	response := SendRequest(requestBody)

	txHash := *ParseResponse[string](response)
	log.Printf("mint erc20 txHash: %s", txHash)
	receiptChan := waitForReceipt(txHash)
	receipt := <-receiptChan

	log.Printf("Mint Log Receipt: %v", ToJsonString(receipt))
	log.Printf("Mint Log: %v. Is nil = %v", receipt.Logs, receipt.Logs == nil)
}

func readBalance(contractAddr common.Address) {
	funcName := "Read Balance"
	log.Println("================================")
	log.Println(funcName)
	log.Println("================================")
	abi, _ := abigen.TestErc20MetaData.GetAbi()
	data, _ := abi.Pack("balanceOf", testWalletAddr)

	arg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}
	result, err := client.CallContract(context.Background(), arg, nil)
	outputs, err := abi.Unpack("balanceOf", result)
	if err != nil {
		log.Fatalf("Failed to unpack contract output: %v", err)
	}

	counter := outputs[0].(*big.Int)
	log.Printf("%v Counter value: %s", funcName, counter.String())
}

func sendEvent(contractAddr common.Address) {
	log.Println("================================")
	log.Println("Send Event")
	log.Println("================================")
	abi, _ := abigen.TestErc20MetaData.GetAbi()
	data, _ := abi.Pack("sendEvent")

	dataHex := hexutil.Bytes(data)
	arg := ethrpc.TransactionArgs{
		From: &testWalletAddr,
		To:   &contractAddr,
		Data: &dataHex,
	}
	estimateGasBody := GenerateRequestBody("eth_estimateGas", arg)
	estimateGasResponse := SendRequest(estimateGasBody)
	gasLimit := ParseResponseAsBigInt(estimateGasResponse).Uint64()

	nonce := getNonce()
	gasPrice := getGasPrice()
	to := contractAddr

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Data:     data,
	})

	rawTx := SignTransaction(gethCfg, testWalletPrivateKeyStr, tx)
	requestBody := GenerateRequestBody("eth_sendRawTransaction", rawTx)
	response := SendRequest(requestBody)

	txHash := *ParseResponse[string](response)
	log.Printf("Send Event txHash: %s", txHash)
	receiptChan := waitForReceipt(txHash)
	receipt := <-receiptChan

	log.Printf("Send Event Receipt: %v", ToJsonString(receipt))
	log.Printf("Send Event : %v. Is nil = %v", receipt.Logs, receipt.Logs == nil)
}

func sendEmptyEvent(contractAddr common.Address) {
	log.Println("================================")
	log.Println("Send Empty Event")
	log.Println("================================")
	abi, _ := abigen.TestErc20MetaData.GetAbi()
	data, _ := abi.Pack("sendEmptyEvent")

	dataHex := hexutil.Bytes(data)
	arg := ethrpc.TransactionArgs{
		From: &testWalletAddr,
		To:   &contractAddr,
		Data: &dataHex,
	}
	estimateGasBody := GenerateRequestBody("eth_estimateGas", arg)
	estimateGasResponse := SendRequest(estimateGasBody)
	gasLimit := ParseResponseAsBigInt(estimateGasResponse).Uint64()

	nonce := getNonce()
	gasPrice := getGasPrice()
	to := contractAddr

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Data:     data,
	})

	rawTx := SignTransaction(gethCfg, testWalletPrivateKeyStr, tx)
	requestBody := GenerateRequestBody("eth_sendRawTransaction", rawTx)
	response := SendRequest(requestBody)

	txHash := *ParseResponse[string](response)
	log.Printf("Send Empty Event txHash: %s", txHash)
	receiptChan := waitForReceipt(txHash)
	receipt := <-receiptChan

	log.Printf("Send Empty Event Receipt: %v", ToJsonString(receipt))
	log.Printf("Send Empty Event : %v. Is nil = %v", receipt.Logs, receipt.Logs == nil)
}

func writeCounter(contractAddr common.Address, newCounter int64) {
	funcName := "Write Counter"
	log.Println("================================")
	log.Println(funcName)
	log.Println("================================")
	abi, _ := abigen.TestErc20MetaData.GetAbi()
	data, _ := abi.Pack("writeCounter", big.NewInt(newCounter))

	dataHex := hexutil.Bytes(data)
	arg := ethrpc.TransactionArgs{
		From: &testWalletAddr,
		To:   &contractAddr,
		Data: &dataHex,
	}
	estimateGasBody := GenerateRequestBody("eth_estimateGas", arg)
	estimateGasResponse := SendRequest(estimateGasBody)
	gasLimit := ParseResponseAsBigInt(estimateGasResponse).Uint64()

	nonce := getNonce()
	gasPrice := getGasPrice()
	to := contractAddr

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Data:     data,
	})

	rawTx := SignTransaction(gethCfg, testWalletPrivateKeyStr, tx)
	requestBody := GenerateRequestBody("eth_sendRawTransaction", rawTx)
	response := SendRequest(requestBody)

	txHash := *ParseResponse[string](response)
	log.Printf("%v txHash: %s", funcName, txHash)
	receiptChan := waitForReceipt(txHash)
	receipt := <-receiptChan

	log.Printf("%v Receipt: %v", funcName, ToJsonString(receipt))
	log.Printf("%v : %v. Is nil = %v", funcName, receipt.Logs, receipt.Logs == nil)
}

func readCounter(contractAddr common.Address) {
	funcName := "Read Counter"
	log.Println("================================")
	log.Println(funcName)
	log.Println("================================")
	abi, _ := abigen.TestErc20MetaData.GetAbi()
	data, _ := abi.Pack("readCounter")

	arg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}
	result, err := client.CallContract(context.Background(), arg, nil)
	outputs, err := abi.Unpack("readCounter", result)
	if err != nil {
		log.Fatalf("Failed to unpack contract output: %v", err)
	}

	counter := outputs[0].(*big.Int)
	log.Printf("%v Counter value: %s", funcName, counter.String())
}

func readBlockNumber(contractAddr common.Address) uint64 {
	funcName := "Read BlockNumber"
	log.Println("================================")
	log.Println(funcName)
	log.Println("================================")
	abi, _ := abigen.TestErc20MetaData.GetAbi()
	data, _ := abi.Pack("currentBlockNumber")

	arg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}
	result, err := client.CallContract(context.Background(), arg, nil)
	outputs, err := abi.Unpack("currentBlockNumber", result)
	if err != nil {
		log.Fatalf("Failed to unpack contract output: %v", err)
	}

	blockNumber := outputs[0].(*big.Int)
	log.Printf("%v currentBlockNumber: %s", funcName, blockNumber.String())
	return blockNumber.Uint64()
}
