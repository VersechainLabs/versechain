package main

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"itachi/evm"
	"itachi/evm/ethrpc"
)

func testTx() {
	rawTx := "0xf9016a0a84773594008329fa2d94949f543ae523d7ab89a70a397bd7dffbf514c9f480b90104e8e33700000000000000000000000000f881bdc361532bdb11aca0f5d75702bbabb5b179000000000000000000000000bb99f4bc5df3c6260ab79afeca4e05c70a6224f8000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007bd36074b61cfe75a53e1b9df7678c96e6463b020000000000000000000000000000000000000000000000000000000066b3905326a0454ff6d3bebdca04df54f981e9b36bce2c30fb5dc7fa94c4ce4bc3f11c36a020a005656e85853a4eea02293b97d65f5774c3f4fd1ef3146049db7a80b446339da0"
	input, _ := hexutil.Decode(rawTx)
	signedTx := new(types.Transaction)
	if err := signedTx.UnmarshalBinary(input); err != nil {
		logrus.Fatal(err)
	}

	logrus.Printf("Tx Type = %v, Hash = %v", signedTx.Type(), signedTx.Hash().Hex())
	jsonBytes, _ := json.Marshal(signedTx)
	logrus.Printf("%v", string(jsonBytes))

	v, r, s := signedTx.RawSignatureValues()
	args := ethrpc.NewTxArgsFromTx(signedTx)
	argsByt, _ := json.Marshal(args)
	txReq := &evm.TxRequest{
		Input:    signedTx.Data(),
		Origin:   common.HexToAddress(testWalletAddrStr),
		Address:  signedTx.To(),
		GasLimit: signedTx.Gas(),
		GasPrice: signedTx.GasPrice(),
		Value:    signedTx.Value(),
		Hash:     signedTx.Hash(),
		Nonce:    signedTx.Nonce(),
		V:        v,
		R:        r,
		S:        s,

		OriginArgs: argsByt,
	}
	byt, _ := json.Marshal(txReq)
	logrus.Println(string(byt))
	logrus.Printf("----")

	newTxReq := &evm.TxRequest{}
	json.Unmarshal(byt, newTxReq)
	newArgs := &ethrpc.TransactionArgs{}
	json.Unmarshal(txReq.OriginArgs, newArgs)
	tx := newArgs.ToTransaction(txReq.V, txReq.R, txReq.S)
	logrus.Printf("Tx Type = %v, Hash = %v", tx.Type(), tx.Hash().Hex())
	newJsonBytes, _ := json.Marshal(tx)
	logrus.Printf("%v", string(newJsonBytes))

	//newTx := args.ToTransaction(signedTx.RawSignatureValues())
	//logrus.Printf("Tx Type = %v, Hash = %v", newTx.Type(), newTx.Hash().Hex())
	//jsonBytes, _ = json.Marshal(newTx)
	//logrus.Printf("%v", string(jsonBytes))

	return
}
