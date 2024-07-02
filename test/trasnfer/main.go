package main

import (
	"encoding/json"
	"fmt"

	"itachi/evm"
	"itachi/test/trasnfer/pkg"
)

const (
	HostAddress     = "localhost:9092"
	cfgPath         = "../../conf/evm_cfg.toml"
	generateCount   = 10
	initialEthCount = 100 * 100
	testSteps       = 5
)

func main() {
	cfg := evm.LoadEvmConfig(cfgPath)
	walletsManager := pkg.NewWalletManager(cfg, HostAddress)
	wallets, err := walletsManager.GenerateRandomWallet(generateCount, initialEthCount)
	if err != nil {
		panic(err)
	}
	testStepsManager := pkg.NewTransferManager()
	tc := testStepsManager.GenerateTransferSteps(testSteps, pkg.GenerateCaseWallets(initialEthCount, wallets))
	b, err := json.Marshal(tc)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
