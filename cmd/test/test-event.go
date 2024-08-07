package main

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	abigen "itachi/cmd/test/abi"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	TestErc20Contract *abigen.TestErc20
)

func TestGetLog() {
	// 1. Deploy Contract
	_, contractAddr := createContract()
	auth, _ := bind.NewKeyedTransactorWithChainID(testPrivateKey, gethCfg.ChainConfig.ChainID)
	auth.Value = big.NewInt(0)

	TestErc20Contract, _ = abigen.NewTestErc20(contractAddr, client)
	tx, err := TestErc20Contract.Mint(auth, testWalletAddr, ether_100)
	if err != nil {
		log.Printf("failed to do mint. err = %v", err)
		return
	}
	log.Printf("Tx sent: %s\n", tx.Hash().Hex())

	pollEvent(contractAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}

func pollEvent(contractAddr common.Address) {
	startBlock := big.NewInt(1)
	toBlock := startBlock
	pollInterval := time.Second * 5
	go func() {
		for {
			select {
			case <-time.After(pollInterval):
				startBlock = big.NewInt(0).Add(toBlock, big.NewInt(1))
				endBlock, err := client.BlockNumber(context.Background())
				if err != nil {
					log.Printf("Failed to get latest block number: %v", err)
					continue
				}
				toBlock = big.NewInt(int64(endBlock))
				// if startBlock larger than toBlock, wait next poll
				if toBlock.Cmp(startBlock) < 0 {
					continue
				}

				log.Printf("Indexer: poll event from %v to %v", startBlock, toBlock)

				logs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{
					FromBlock: startBlock,
					ToBlock:   toBlock,
					Addresses: []common.Address{contractAddr},
				})
				if err != nil {
					log.Printf("Failed to filter logs: %v", err)
					continue
				}

				for _, vLog := range logs {
					processLog(vLog)
				}
			}
		}

	}()

	return
}

func processLog(vlog types.Log) {
	transferEventHash := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	var firstTopicHash common.Hash
	if len(vlog.Topics) > 0 {
		firstTopicHash = vlog.Topics[0]
	}
	isTransferEvent := false
	if firstTopicHash == common.HexToHash(transferEventHash) {
		isTransferEvent = true
	}
	log.Printf("Log detail: isTransfer=%v, address=%v, topic count=%v, first topic=%v", isTransferEvent, vlog.Address, len(vlog.Topics), firstTopicHash.Hex())
	if !isTransferEvent {
		return
	}

	event, err := TestErc20Contract.ParseTransfer(vlog)
	if err != nil {
		log.Printf("Failed to parse event: %v", err)
		return
	}

	log.Printf("Parsed Log: %+v", event)
}
