package utils

import (
	"errors"
	"github.com/NethermindEth/juno/core/felt"
	"log"
)

var ErrUnknownNetwork = errors.New("unknown network (known: mainnet, goerli, goerli2, integration)")

type Network int

const (
	Mainnet Network = iota + 1
	Goerli
	Goerli2
	Integration
	Sepolia
	SepoliaIntegration
)

func (n Network) String() string {
	switch n {
	case Mainnet:
		return "mainnet"
	case Goerli:
		return "goerli"
	case Goerli2:
		return "goerli2"
	case Integration:
		return "integration"
	case Sepolia:
		return "sepolia"
	case SepoliaIntegration:
		return "sepolia-integration"
	default:
		// Should not happen.
		panic(ErrUnknownNetwork)
	}
}

func (n Network) ChainIDString() string {
	log.Print("enter ChainIDString()")
	log.Print(n)

	switch n {
	case Goerli, Integration:
		return "SN_GOERLI"
	case Mainnet:
		return "SN_MAIN"
	case Goerli2:
		return "SN_GOERLI2"
	case Sepolia:
		return "SN_SEPOLIA"
	case SepoliaIntegration:
		return "SN_INTEGRATION_SEPOLIA"
	default:
		// Should not happen.
		panic(ErrUnknownNetwork)
	}
}

func (n Network) ChainID() *felt.Felt {
	log.Print("enter ChainID()")
	return new(felt.Felt).SetBytes([]byte(n.ChainIDString()))
}
