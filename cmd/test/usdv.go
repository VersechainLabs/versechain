package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	abigen "itachi/cmd/test/abi"
	"itachi/evm/ethrpc"
	"log"
	"math/big"
	"os"
	"strings"
)

func deployAndMintUsdv() {
	//_, contractAddr := createUsdvContract()
	contractAddr := common.HexToAddress("0x310b8685e3e69cb05b251a12f5ffab23001cda42")
	testWalletAddr_9527 := common.HexToAddress("0x2Efe24c33f049Ffec693ec1D809A45Fff14e9527")
	// mintErc20(contractAddr)
	//sendEvent(contractAddr)
	//sendEmptyEvent(contractAddr)

	//testReadCache(contractAddr)

	readUsdvBalance(contractAddr, testWalletAddr_9527)
	testTransferUsdv(contractAddr)
	readUsdvBalance(contractAddr, testWalletAddr_9527)

}

func createUsdvContract() (txHash string, contractAddr common.Address) {
	log.Println("================================")
	log.Println("Create Contract")
	log.Println("================================")
	contractBinByte, err := os.ReadFile("./cmd/test/abi/USDV.bin")
	if err != nil {
		log.Fatal("failed to load bin file: ", err)
	}

	abiFile, err := os.ReadFile("./cmd/test/abi/USDV.abi")
	if err != nil {
		log.Fatal("failed to load abi file: ", err)
	}
	contractABI, err := abi.JSON(strings.NewReader(string(abiFile)))
	if err != nil {
		log.Fatal("failed to parse abi file: ", err)
	}

	initialSupply := ether_100
	name := "USDV"
	symbol := "USDV"
	decimals := big.NewInt(8)
	packedData, err := contractABI.Pack("", initialSupply, name, symbol, decimals)
	if err != nil {
		log.Fatal("failed to pack data: ", err)
	}

	fmt.Printf("packedData: %x\n", packedData)
	contractData, err := hex.DecodeString(string(contractBinByte))
	data := append(contractData, packedData[4:]...)

	//fmt.Println("1")
	//fmt.Println(string(contractBinByte))
	//fmt.Println("2")
	//fmt.Println(string(packedData))
	//fmt.Println("3")
	//fmt.Println(hex.DecodeString(string(contractBinByte)))
	//fmt.Println("4")
	//fmt.Println(hex.DecodeString(string(packedData)))
	// estimate gas
	//data, err := hex.DecodeString(string(fullData))

	//if err != nil {
	//	log.Fatal("failed to encode data: ", err)
	//}

	//gasLimit := estimateGas(data)
	gasLimit := uint64(5000000000)
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

func testTransferUsdv(contractAddr common.Address) {
	abi, _ := abigen.TestErc20MetaData.GetAbi()
	data, _ := abi.Pack("transfer", common.HexToAddress("0x2Efe24c33f049Ffec693ec1D809A45Fff14e9527"), ether_100)
	log.Printf("Data: %x\n", data)

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
	log.Printf("Transfer Log Receipt: %v", ToJsonString(receipt))

}

func mintUsdv(contractAddr common.Address) {
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

func readUsdvBalance(contractAddr common.Address, walletAddr common.Address) {
	funcName := "Read Balance"
	log.Println("================================")
	log.Println(funcName)
	log.Println("================================")
	abiFile, err := os.ReadFile("./cmd/test/abi/USDV.abi")
	if err != nil {
		log.Fatal("failed to load abi file: ", err)
	}
	contractABI, err := abi.JSON(strings.NewReader(string(abiFile)))
	if err != nil {
		log.Fatal("failed to parse abi file: ", err)
	}

	data, _ := contractABI.Pack("balanceOf", walletAddr)

	arg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}
	result, err := client.CallContract(context.Background(), arg, nil)
	outputs, err := contractABI.Unpack("balanceOf", result)
	if err != nil {
		log.Fatalf("Failed to unpack contract output: %v", err)
	}

	counter := outputs[0].(*big.Int)
	log.Printf("%v Counter value: %s", funcName, counter.String())
}
