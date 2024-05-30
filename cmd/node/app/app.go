package app

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/yu-org/yu/apps/poa"
	"github.com/yu-org/yu/core/kernel"
	"github.com/yu-org/yu/core/startup"
	"itachi/cairo"
	"itachi/cairo/config"
	"itachi/cairo/starknetrpc"
	"itachi/evm"
	"itachi/evm/ethrpc"
	"itachi/utils"
)

func StartUpChain(poaCfg *poa.PoaConfig, crCfg *config.Config, evmCfg *evm.GethConfig) {
	figure.NewColorFigure("Itachi", "big", "green", false).Print()

	chain := InitItachi(poaCfg, crCfg, evmCfg)
	starknetrpc.StartUpStarknetRPC(chain, crCfg)
	ethrpc.StartupEthRPC(chain, evmCfg)
	utils.StartUpPprof(crCfg)
	chain.Startup()
}

func InitItachi(poaCfg *poa.PoaConfig, crCfg *config.Config, evmCfg *evm.GethConfig) *kernel.Kernel {
	poaTri := poa.NewPoa(poaCfg)
	cairoTri := cairo.NewCairo(crCfg)
	solidityTri := evm.NewSolidity(evmCfg)
	chain := startup.InitDefaultKernel(
		poaTri, cairoTri, solidityTri,
	)
	return chain
}
