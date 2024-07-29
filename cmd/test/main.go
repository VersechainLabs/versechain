package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	abigen "itachi/cmd/test/abi"
	"itachi/evm"
	"itachi/evm/ethrpc"
	"log"
	"math/big"
	"os"
	"time"
)

const (
	testWalletPrivateKeyStr = "32e3b56c9f2763d2332e6e4188e4755815ac96441e899de121969845e343c2ff"
	testWalletAddrStr       = "0x7Bd36074b61Cfe75a53e1B9DF7678C96E6463b02"
)

var (
	rpcId          = 0
	gethCfg        = evm.LoadEvmConfig("./conf/evm_cfg.toml")
	testWalletAddr = common.HexToAddress(testWalletAddrStr)
	testPrivateKey *ecdsa.PrivateKey
	client         *ethclient.Client

	ether     = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	ether_100 = new(big.Int).Mul(big.NewInt(100), ether)
)

func init() {
	testPrivateKey, _ = crypto.HexToECDSA(testWalletPrivateKeyStr)
	client, _ = ethclient.Dial("http://localhost:9092")
}

func main() {
	//testTransferEth()
	//testGetBalance()
	testCreateContract()
	//testMintErc20()
	//testErc20DeployAndUse()
}

func testGetBalance() {
	checkBalanceParam := []interface{}{testWalletAddrStr, "latest"}
	checkBalanceBody := GenerateRequestBody("eth_getBalance", checkBalanceParam...)
	log.Println(checkBalanceBody)
	response := SendRequest(checkBalanceBody)
	result := ParseResponseAsBigInt(response)
	log.Println(fmt.Sprintf("Balance: %d (%d ether)", result, new(big.Int).Div(result, ether)))
}

func testTransferEth() {
	nonce := uint64(0)
	to := common.HexToAddress("0x756c2d2bb2b4a2b82939861072eb687c6b6f5d93")
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(1)
	var data []byte

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Value:    ether,
		Data:     data,
	})

	printTxDetail(tx)

	rawTx := SignTransaction(gethCfg, testWalletPrivateKeyStr, tx)
	requestBody := GenerateRequestBody("eth_sendRawTransaction", rawTx)
	log.Println(requestBody)
	response := SendRequest(requestBody)
	log.Println(ParseResponse[string](response))
}

func testCreateContract() {
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

	txHash := *ParseResponse[string](response)
	log.Printf("create contract's txHash: %s", txHash)
	receiptChan := waitForReceipt(txHash)
	receipt := <-receiptChan
	contractAddr := receipt.ContractAddress

	log.Printf("create contract's contractAddr: %s", contractAddr)
}

func testMintErc20() {
	abi, _ := abigen.TestErc20MetaData.GetAbi()
	data, _ := abi.Pack("mint", testWalletAddr, ether_100)
	log.Printf("Data: %x\n", data)

	contractAddress := common.HexToAddress("0x310B8685e3E69Cb05b251A12f5FFAb23001CdA42")
	nonce := uint64(0)
	to := contractAddress
	amount := big.NewInt(0)
	gasLimit := uint64(210000)
	gasPrice := big.NewInt(0)

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
	log.Println(ParseResponse[string](response))
}

func testMintErc20ByAbi() {
	contractAddress := common.HexToAddress("0x310B8685e3E69Cb05b251A12f5FFAb23001CdA42")

	auth, _ := bind.NewKeyedTransactorWithChainID(testPrivateKey, gethCfg.ChainConfig.ChainID)
	auth.Value = big.NewInt(0)
	auth.Nonce = big.NewInt(0)
	auth.GasLimit = uint64(210000)
	auth.GasPrice = big.NewInt(0)

	contract, _ := abigen.NewTestErc20(contractAddress, client)
	tx, _ := contract.Mint(auth, testWalletAddr, ether_100)
	log.Printf("Tx sent: %s\n", tx.Hash().Hex())
}

func testTransferErc20() {
	abi, _ := abigen.TestErc20MetaData.GetAbi()
	data, _ := abi.Pack("transfer", common.HexToAddress("0x2Efe24c33f049Ffec693ec1D809A45Fff14e9527"), ether)
	log.Printf("Data: %x\n", data)

	contractAddress := common.HexToAddress("0x310B8685e3E69Cb05b251A12f5FFAb23001CdA42")
	nonce := uint64(0)
	to := contractAddress
	amount := big.NewInt(0)
	gasLimit := uint64(210000)
	gasPrice := big.NewInt(0)

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
	log.Println(ParseResponse[string](response))
}

func printTxDetail(tx *types.Transaction) {
	log.Println("---- Tx Detail Start ----")
	originTxByte, _ := json.Marshal(tx)
	log.Printf("[TxDetail] Befroe sign = %v", string(originTxByte))

	privateKey, err := crypto.HexToECDSA(testWalletPrivateKeyStr)
	if err != nil {
		log.Fatal(err)
	}
	chainID := gethCfg.ChainConfig.ChainID
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	signedTxByte, _ := json.Marshal(signedTx)
	log.Printf("[TxDetail] After sign = %v", string(signedTxByte))

	rawTxBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		log.Fatal(err)
	}
	rawTx := fmt.Sprintf("0x%x", rawTxBytes)
	rawTxByte, _ := json.Marshal(rawTx)
	log.Printf("[TxDetail] Raw = %v", string(rawTxByte))
	log.Println("---- Tx Detail End ----")
}

func estimateGas(data hexutil.Bytes) uint64 {
	arg := ethrpc.TransactionArgs{
		From: &testWalletAddr,
		To:   nil,
		Data: &data,
	}
	estimateGasBody := GenerateRequestBody("eth_estimateGas", arg)
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

func getTransactionReceipt(txHash string) (receipt *types.Receipt, err error) {
	requestBody := GenerateRequestBody("eth_getTransactionReceipt", txHash)
	response := SendRequest(requestBody)

	receipt = ParseResponse[types.Receipt](response)
	return receipt, nil
}

func waitForReceipt(txHash string) <-chan *types.Receipt {
	receiptChan := make(chan *types.Receipt)

	go func() {
		defer close(receiptChan)

		for {
			receipt, err := getTransactionReceipt(txHash)
			if err != nil {
				log.Printf("Error getting transaction receipt: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if receipt == nil {
				//log.Println("Transaction receipt is not available yet (no receipt). Waiting...")
				time.Sleep(1 * time.Second)
				continue
			}

			if receipt.Status != types.ReceiptStatusSuccessful {
				//log.Println("Receipt status is not successful. Waiting...")
				time.Sleep(1 * time.Second)
				continue
			}

			receiptChan <- receipt
			return
		}
	}()

	return receiptChan
}
