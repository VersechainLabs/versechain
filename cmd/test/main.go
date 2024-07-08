package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	abigen "itachi/cmd/test/abi"
	"itachi/evm"
	"log"
	"math/big"
	"os"
	"strconv"
)

const (
	testWalletPrivateKeyStr = "32e3b56c9f2763d2332e6e4188e4755815ac96441e899de121969845e343c2ff"
	testWalletAddr          = "0x7Bd36074b61Cfe75a53e1B9DF7678C96E6463b02"
)

var (
	rpcId   = 0
	gethCfg = evm.LoadEvmConfig("./conf/evm_cfg.toml")
)

func main() {
	testGetBalance()
	//testCreateContract()
	testMintErc20()
}

func testGetBalance() {
	checkBalanceParam := []interface{}{testWalletAddr, "latest"}
	checkBalanceBody := GenerateRequestBody("eth_getBalance", checkBalanceParam...)
	log.Println(checkBalanceBody)
	response := SendRequest(checkBalanceBody)
	result, _ := strconv.ParseInt(ParseResponse(response), 16, 64)
	log.Println(fmt.Sprintf("Response: %s, Balance: %d", ParseResponse(response), result))
}

func testMintErc20() {
	contractAddress := common.HexToAddress("0x310B8685e3E69Cb05b251A12f5FFAb23001CdA42")
	client, _ := ethclient.Dial("http://localhost:9092")
	privateKey, _ := crypto.HexToECDSA(testWalletPrivateKeyStr)

	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, gethCfg.ChainConfig.ChainID)
	auth.Value = big.NewInt(0)
	auth.Nonce = big.NewInt(0)
	auth.GasLimit = uint64(210000)
	auth.GasPrice = big.NewInt(0)

	contract, _ := abigen.NewTestErc20(contractAddress, client)
	tx, _ := contract.Mint(auth, common.HexToAddress(testWalletAddr), big.NewInt(100000000000000))
	log.Printf("Tx sent: %s\n", tx.Hash().Hex())
}

func testTransferErc20() {
	nonce := uint64(0)
	to := common.HexToAddress("0x2Efe24c33f049Ffec693ec1D809A45Fff14e9527")
	amount := big.NewInt(1)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(0)
	data := []byte{}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Value:    amount,
		Data:     data,
	})

	rawTx := SignTransaction(gethCfg, testWalletPrivateKeyStr, tx)
	requestBody := GenerateRequestBody("eth_sendRawTransaction", rawTx)
	log.Println(requestBody)
	response := SendRequest(requestBody)
	log.Println(ParseResponse(response))
}

func testTransferEth() {
	nonce := uint64(0)
	to := common.HexToAddress("0x2Efe24c33f049Ffec693ec1D809A45Fff14e9527")
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(0)
	var data []byte

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Value:    big.NewInt(100),
		Data:     data,
	})

	rawTx := SignTransaction(gethCfg, testWalletPrivateKeyStr, tx)
	requestBody := GenerateRequestBody("eth_sendRawTransaction", rawTx)
	log.Println(requestBody)
	response := SendRequest(requestBody)
	log.Println(ParseResponse(response))
}

func testCreateContract() {
	contractBinByte, err := os.ReadFile("./cmd/test/abi/TestErc20.bin")
	if err != nil {
		log.Fatal("failed to load bin file: ", err)
	}

	nonce := uint64(0)
	amount := big.NewInt(0)
	gasLimit := uint64(21000000)
	gasPrice := big.NewInt(0)
	data, err := hex.DecodeString(string(contractBinByte))
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       nil,
		Value:    amount,
		Data:     data,
	})

	rawTx := SignTransaction(gethCfg, testWalletPrivateKeyStr, tx)
	requestBody := GenerateRequestBody("eth_sendRawTransaction", rawTx)
	log.Println(requestBody)
	response := SendRequest(requestBody)
	log.Println(ParseResponse(response))
}
