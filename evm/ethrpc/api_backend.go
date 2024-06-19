package ethrpc

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/gasprice"
	"github.com/ethereum/go-ethereum/rpc"
)

//	type EthAPIBackend struct {
//		allowUnprotectedTxs bool
//		ethChainCfg         *params.ChainConfig
//		chain               *kernel.Kernel
//	}
//
// EthAPIBackend implements ethapi.Backend and tracers.Backend for full nodes
type EthAPIBackend struct {
	extRPCEnabled       bool
	allowUnprotectedTxs bool
	eth                 *Ethereum
	gpo                 *gasprice.Oracle
}

//	func (b *EthAPIBackend) SyncProgress() ethereum.SyncProgress {
//		//TODO implement me
//		panic("implement me")
//		prog := b.eth.Downloader().Progress()
//		if txProg, err := b.eth.blockchain.TxIndexProgress(); err == nil {
//			prog.TxIndexFinishedBlocks = txProg.Indexed
//			prog.TxIndexRemainingBlocks = txProg.Remaining
//		}
//		return prog
//	}
//
//	func (e *EthAPIBackend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
//		//TODO implement me
//		//panic("implement me")
//		return e.gpo.SuggestTipCap(ctx)
//	}
//
//	func (e *EthAPIBackend) FeeHistory(ctx context.Context, blockCount uint64, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*big.Int, [][]*big.Int, []*big.Int, []float64, []*big.Int, []float64, error) {
//		//TODO implement me
//		//panic("implement me")
//		return e.gpo.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
//	}
//
//	func (e *EthAPIBackend) BlobBaseFee(ctx context.Context) *big.Int {
//		//TODO implement me
//		//panic("implement me")
//		if excess := e.CurrentHeader().ExcessBlobGas; excess != nil {
//			return eip4844.CalcBlobFee(*excess)
//		}
//		return nil
//	}
//
//	func (s *EthAPIBackend) ChainDb() ethdb.Database {
//		//TODO implement me
//		//panic("implement me")
//		return e.eth.ChainDb()
//	}
//
//	func (e *EthAPIBackend) AccountManager() *accounts.Manager {
//		//TODO implement me
//		return nil
//		//panic("implement me")
//	}
//
//	func (e *EthAPIBackend) ExtRPCEnabled() bool {
//		//TODO implement me
//		//panic("implement me")
//		return e.extRPCEnabled
//	}
//
//	func (e *EthAPIBackend) RPCGasCap() uint64 {
//		//return 50000000
//		return e.eth.config.RPCGasCap
//	}
//
//	func (e *EthAPIBackend) RPCEVMTimeout() time.Duration {
//		//return 5 * time.Second
//		return e.eth.config.RPCEVMTimeout
//	}
//
//	func (e *EthAPIBackend) RPCTxFeeCap() float64 {
//		//return 1
//		return e.eth.config.RPCTxFeeCap
//	}
//
//	func (e *EthAPIBackend) UnprotectedAllowed() bool {
//		return e.allowUnprotectedTxs
//	}
//
//	func (e *EthAPIBackend) SetHead(number uint64) {
//		//TODO implement me
//		panic("implement me")
//	}
func (b *EthAPIBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	//TODO implement me
	//panic("implement me")
	// Pending block is only known by the miner
	if number == rpc.PendingBlockNumber {
		block, _, _ := b.eth.miner.Pending()
		if block == nil {
			return nil, errors.New("pending block is not available")
		}
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if number == rpc.LatestBlockNumber {
		return b.eth.blockchain.CurrentBlock(), nil
	}
	if number == rpc.FinalizedBlockNumber {
		block := b.eth.blockchain.CurrentFinalBlock()
		if block == nil {
			return nil, errors.New("finalized block not found")
		}
		return block, nil
	}
	if number == rpc.SafeBlockNumber {
		block := b.eth.blockchain.CurrentSafeBlock()
		if block == nil {
			return nil, errors.New("safe block not found")
		}
		return block, nil
	}
	return b.eth.blockchain.GetHeaderByNumber(uint64(number)), nil
}

//
//func (b *EthAPIBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
//	//TODO implement me
//	panic("implement me")
//	return b.eth.blockchain.GetHeaderByHash(hash), nil
//}

func (b *EthAPIBackend) HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error) {
	//TODO implement me
	//panic("implement me")
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.HeaderByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header := b.eth.blockchain.GetHeaderByHash(hash)
		if header == nil {
			return nil, errors.New("header for hash not found")
		}
		if blockNrOrHash.RequireCanonical && b.eth.blockchain.GetCanonicalHash(header.Number.Uint64()) != hash {
			return nil, errors.New("hash is not currently canonical")
		}
		return header, nil
	}
	return nil, errors.New("invalid arguments; neither block nor hash specified")
}

//
//func (e *EthAPIBackend) CurrentHeader() *types.Header {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) CurrentBlock() *types.Header {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (b *EthAPIBackend) StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
//	//TODO implement me
//	panic("implement me")
//	// Pending state is only known by the miner
//	if number == rpc.PendingBlockNumber {
//		block, _, state := b.eth.miner.Pending()
//		if block == nil || state == nil {
//			return nil, nil, errors.New("pending state is not available")
//		}
//		return state, block.Header(), nil
//	}
//	// Otherwise resolve the block number and return its state
//	header, err := b.HeaderByNumber(ctx, number)
//	if err != nil {
//		return nil, nil, err
//	}
//	if header == nil {
//		return nil, nil, errors.New("header not found")
//	}
//	stateDb, err := b.eth.BlockChain().StateAt(header.Root)
//	if err != nil {
//		return nil, nil, err
//	}
//	return stateDb, header, nil
//}
//
//func (e *EthAPIBackend) StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*state.StateDB, *types.Header, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) Pending() (*types.Block, types.Receipts, *state.StateDB) {
//	//TODO implement me
//	panic("implement me")
//	return e.eth.miner.Pending()
//}
//
//func (e *EthAPIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) GetTd(ctx context.Context, hash common.Hash) *big.Int {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) GetEVM(ctx context.Context, msg *core.Message, state *state.StateDB, header *types.Header, vmConfig *vm.Config, blockCtx *vm.BlockContext) *vm.EVM {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
//	signer := types.NewEIP155Signer(e.ethChainCfg.ChainID)
//	sender, err := types.Sender(signer, signedTx)
//	if err != nil {
//		return err
//	}
//	txReq := &evm.TxRequest{
//		Input:    signedTx.Data(),
//		Origin:   sender,
//		Address:  *signedTx.To(),
//		GasLimit: signedTx.Gas(),
//		GasPrice: signedTx.GasPrice(),
//		Value:    signedTx.Value(),
//	}
//	byt, err := json.Marshal(txReq)
//	logrus.Printf("SendTx, Request=%+v\n", string(byt))
//	if err != nil {
//		return err
//	}
//	signedWrCall := &yucore.SignedWrCall{
//		Call: &yucommon.WrCall{
//			TripodName: SolidityTripod,
//			FuncName:   "ExecuteTxn",
//			Params:     string(byt),
//		},
//	}
//	return e.chain.HandleTxn(signedWrCall)
//}
//
//func (e *EthAPIBackend) GetTransaction(ctx context.Context, txHash common.Hash) (bool, *types.Transaction, common.Hash, uint64, uint64, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) GetPoolTransactions() (types.Transactions, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) GetPoolTransaction(txHash common.Hash) *types.Transaction {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) Stats() (pending int, queued int) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) TxPoolContent() (map[common.Address][]*types.Transaction, map[common.Address][]*types.Transaction) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) TxPoolContentFrom(addr common.Address) ([]*types.Transaction, []*types.Transaction) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) SubscribeNewTxsEvent(events chan<- core.NewTxsEvent) event.Subscription {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) ChainConfig() *params.ChainConfig {
//	return e.ethChainCfg
//}
//
//func (e *EthAPIBackend) Engine() consensus.Engine {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (b *EthAPIBackend) GetBody(ctx context.Context, hash common.Hash, number rpc.BlockNumber) (*types.Body, error) {
//	//TODO implement me
//	//panic("implement me")
//	if number < 0 || hash == (common.Hash{}) {
//		return nil, errors.New("invalid arguments; expect hash and no special block numbers")
//	}
//	if body := b.eth.blockchain.GetBody(hash); body != nil {
//		return body, nil
//	}
//	return nil, errors.New("block body not found")
//}
//
//func (e *EthAPIBackend) GetLogs(ctx context.Context, blockHash common.Hash, number uint64) ([][]*types.Log, error) {
//	//TODO implement me
//	panic("implement me")
//	return rawdb.ReadLogs(e.eth.chainDb, hash, number), nil
//}
//
//func (e *EthAPIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (e *EthAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (b *EthAPIBackend) BloomStatus() (uint64, uint64) {
//	//TODO implement me
//	//panic("implement me")
//	sections, _, _ := b.eth.bloomIndexer.Sections()
//	return params.BloomBitsBlocks, sections
//}
//
//func (b *EthAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
//	//TODO implement me
//	//panic("implement me")
//	// 需重新实现bloombits service => backend.go里的Ethereum core 代码（Ethereum struct）
//	for i := 0; i < bloomFilterThreads; i++ {
//		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.eth.bloomRequests)
//	}
//}
