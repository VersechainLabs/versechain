package main

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	abigen "itachi/cmd/test/abi"
	"itachi/evm/ethrpc"
	"log"
	"os"
)

func testErc20DeployAndUse() {
	_, contractAddr := createContract()
	mintErc20(contractAddr)
}

func createContract() (txHash string, contractAddr common.Address) {
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
	log.Println(requestBody)
	response := SendRequest(requestBody)

	txHash := *ParseResponse[string](response)
	log.Printf("ming erc20 txHash: %s", txHash)
	receiptChan := waitForReceipt(txHash)
	_ = <-receiptChan
}
