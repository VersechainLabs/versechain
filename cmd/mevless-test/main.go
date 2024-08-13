package main

import (
	"github.com/ethereum/go-ethereum/common"
	"itachi/evm"
	"itachi/evm/ethrpc"
	"math/big"
)

const (
	testWalletPrivateKeyStr = "32e3b56c9f2763d2332e6e4188e4755815ac96441e899de121969845e343c2ff"
	testWalletAddrStr       = "0x7Bd36074b61Cfe75a53e1B9DF7678C96E6463b02"
)

var (
	rpcId          = 0
	gethCfg        = evm.LoadEvmConfig("./conf/evm_cfg.toml")
	testWalletAddr = common.HexToAddress(testWalletAddrStr)

	ether = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
)

func main() {
	mevlessTest()
}

func estimateGas(arg *ethrpc.TransactionArgs) uint64 {
	estimateGasBody := GenerateRequestBody("eth_estimateGas", *arg)
	estimateGasResponse := SendRequest(estimateGasBody)
	return ParseResponseAsBigInt(estimateGasResponse).Uint64()
}

func getNonce() uint64 {
	getNonceRequest := GenerateRequestBody("eth_getTransactionCount", testWalletAddrStr, "latest")
	getNonceResponse := SendRequest(getNonceRequest)
	return ParseResponseAsBigInt(getNonceResponse).Uint64()
}

func getGasPrice() *big.Int {
	request := GenerateRequestBody("eth_gasPrice")
	response := SendRequest(request)
	return ParseResponseAsBigInt(response)
}
