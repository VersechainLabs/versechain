package app

import (
	"itachi/cairo"
	"itachi/cairo/config"
	"itachi/cairo/l1"
	"itachi/cairo/starknetrpc"
	"itachi/evm"
	"itachi/evm/ethrpc"
	mev_less "itachi/mev-less"
	mev_less_poa "itachi/mev-less/poa"
	"itachi/utils"

	"github.com/common-nighthawk/go-figure"
	"github.com/yu-org/yu/apps/poa"
	"github.com/yu-org/yu/core/kernel"
	"github.com/yu-org/yu/core/startup"
)

func StartUpChain(poaCfg *poa.PoaConfig, crCfg *config.Config, evmCfg *evm.GethConfig, mevLessPoaCfg *mev_less_poa.PoaConfig, mevLessCfg *mev_less.Config) {
	figure.NewColorFigure("Itachi", "big", "green", false).Print()

	chain := InitItachi(poaCfg, crCfg, evmCfg, mevLessPoaCfg, mevLessCfg)

	// Starknet RPC server
	rpcSrv := starknetrpc.StartUpStarknetRPC(chain, crCfg)

	ethrpc.StartupEthRPC(chain, evmCfg)
	// Subscribe to L1
	l1.StartupL1(chain, crCfg, rpcSrv)

	utils.StartUpPprof(crCfg)

	chain.Startup()

}

func InitItachi(poaCfg *poa.PoaConfig, crCfg *config.Config, evmCfg *evm.GethConfig, mevLessPoaCfg *mev_less_poa.PoaConfig, mevLessCfg *mev_less.Config) *kernel.Kernel {
	poaTri := poa.NewPoa(poaCfg)
	cairoTri := cairo.NewCairo(crCfg)
	solidityTri := evm.NewSolidity(evmCfg)
	mevLessPoaTri := mev_less_poa.NewPoa(mevLessPoaCfg)
	mevLessTri := mev_less.NewMEVless(mevLessCfg)

	chain := startup.InitDefaultKernel(
		poaTri, mevLessPoaTri, mevLessTri, cairoTri, solidityTri,
	)
	return chain
}
