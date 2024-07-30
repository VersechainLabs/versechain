package main

import (
	"flag"
	"github.com/yu-org/yu/apps/poa"
	yuconfig "github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/startup"
	"itachi/cairo/config"
	"itachi/cmd/node/app"
	"itachi/evm"
	"os"
)

var isDebug bool

func init() {
	debugFlag := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if debugFlag != nil && *debugFlag {
		isDebug = true
	}

	if isDebug {
		_ = os.RemoveAll("./verse_db")
		_ = os.RemoveAll("./yu")
		_ = os.RemoveAll("./cairo_db")
	}
}

func main() {
	var kernelCfg *yuconfig.KernelConf

	if isDebug {
		kernelCfg = startup.InitKernelConfigFromPath("./conf/yu_debug_cfg.toml")
	} else {
		kernelCfg = startup.InitDefaultKernelConfig()
	}

	poaCfg := poa.DefaultCfg(0)
	cairoCfg := config.LoadCairoCfg("./conf/cairo_cfg.toml")
	gethCfg := evm.LoadEvmConfig("./conf/evm_cfg.toml")

	app.StartUpChain(kernelCfg, poaCfg, cairoCfg, gethCfg)
}
