package pkg

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"itachi/evm"
)

const (
	GenesisPrivateKey = "32e3b56c9f2763d2332e6e4188e4755815ac96441e899de121969845e343c2ff"
)

type EthWallet struct {
	PK      string `json:"pk"`
	Address string `json:"address"`
}

func (e *EthWallet) Copy() *EthWallet {
	return &EthWallet{
		PK:      e.PK,
		Address: e.Address,
	}
}

type WalletManager struct {
	cfg         *evm.GethConfig
	hostAddress string
}

func NewWalletManager(cfg *evm.GethConfig, hostAddress string) *WalletManager {
	return &WalletManager{
		cfg:         cfg,
		hostAddress: hostAddress,
	}
}

func (m *WalletManager) GenerateRandomWallet(count int, initialEthCount int64) ([]*EthWallet, error) {
	wallets := make([]*EthWallet, 0)
	for i := 0; i < count; i++ {
		wallet, err := m.createEthWallet(initialEthCount)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}
	return wallets, nil
}

func (m *WalletManager) createEthWallet(initialEthCount int64) (*EthWallet, error) {
	privateKey, address := generatePrivateKey()
	if err := m.transferEth(GenesisPrivateKey, privateKey, initialEthCount); err != nil {
		return nil, err
	}
	time.Sleep(3 * time.Second)
	return &EthWallet{PK: privateKey, Address: address}, nil
}

func (m *WalletManager) TransferEth(from, to *EthWallet, amount int64) error {
	log.Println(fmt.Sprintf("transfer %v eth from %v to %v", amount, from.PK, to.Address))
	if err := m.transferEth(from.PK, to.Address, amount); err != nil {
		return err
	}
	return nil
}

func (m *WalletManager) transferEth(privateKeyHex string, toAddress string, amount int64) error {
	nonce := uint64(0)
	to := common.HexToAddress(toAddress)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(0)
	var data []byte

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Value:    big.NewInt(amount),
		Data:     data,
	})

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	chainID := m.cfg.ChainConfig.ChainID
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	rawTxBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		log.Fatal(err)
	}

	requestBody := fmt.Sprintf(
		`	{
		"jsonrpc": "2.0",
		"id": 0,
		"method": "eth_sendRawTransaction",
		"params": ["0x%x"] 
	}`, rawTxBytes)
	return sendRequest(m.hostAddress, requestBody)
}
