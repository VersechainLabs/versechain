package app

import (
	yuconfig "github.com/yu-org/yu/config"
	"itachi/cairo"
	"itachi/cairo/config"
	"itachi/cairo/l1"
	"itachi/cairo/starknetrpc"
	"itachi/evm"
	"itachi/evm/ethrpc"
	"itachi/utils"

	"github.com/common-nighthawk/go-figure"
	"github.com/yu-org/nine-tripods/consensus/poa"
	"github.com/yu-org/yu/core/kernel"
	"github.com/yu-org/yu/core/startup"
)

func StartUpChain(kernelCfg *yuconfig.KernelConf, poaCfg *poa.PoaConfig, crCfg *config.Config, evmCfg *evm.GethConfig) {
	figure.NewColorFigure("Itachi", "big", "green", false).Print()

	chain := InitItachi(kernelCfg, poaCfg, crCfg, evmCfg)

	// Starknet RPC server
	rpcSrv := starknetrpc.StartUpStarknetRPC(chain, crCfg)

	ethrpc.StartupEthRPC(chain, evmCfg)
	// Subscribe to L1
	l1.StartupL1(chain, crCfg, rpcSrv)

	utils.StartUpPprof(crCfg)

	chain.Startup()

}

func InitItachi(kernelCfg *yuconfig.KernelConf, poaCfg *poa.PoaConfig, crCfg *config.Config, evmCfg *evm.GethConfig) *kernel.Kernel {
	poaTri := poa.NewPoa(poaCfg)
	cairoTri := cairo.NewCairo(crCfg)
	solidityTri := evm.NewSolidity(evmCfg)
	chain := startup.InitDefaultKernel(
		kernelCfg, poaTri, cairoTri, solidityTri,
	)
	return chain
}
