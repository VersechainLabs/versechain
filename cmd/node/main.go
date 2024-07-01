package main

import (
	"itachi/cairo/config"
	"itachi/cmd/node/app"
	"itachi/evm"

	"github.com/yu-org/yu/apps/poa"
	"github.com/yu-org/yu/core/startup"
)

func main() {
	startup.InitDefaultKernelConfig()
	poaCfg := poa.DefaultCfg(0)
	cairoCfg := config.LoadCairoCfg("./conf/cairo_cfg.toml")
	gethCfg := evm.LoadEvmConfig("./conf/evm_cfg.toml")

	//// 1. Generate random wallet with 100 eth
	//go func() {
	//	CreateRandomWallet(gethCfg, 10)
	//}()

	//2. Transfer eth to another addr
	//requestBody := GenerateTransferEthRequest(gethCfg, "337746e3ff5cd4833088439cd5b695a11cdd185818ddd1c8cf5135c95d643125", "0x8fE1407582B7FA3B76611875a044Cc16533aFeb1", 1)
	//fmt.Printf("---- Transfer Eth Request Body ----\n%s\n---------\n", requestBody)

	app.StartUpChain(poaCfg, cairoCfg, gethCfg)
}
