package evm

import (
	// "github.com/yu-org/yu/common/yerror"

	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/yu-org/yu/common/yerror"
	"itachi/evm/config"
	"math/big"
	"net/http"

	"github.com/sirupsen/logrus"
	yu_common "github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	yu_types "github.com/yu-org/yu/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"

	"github.com/NethermindEth/juno/jsonrpc"
	"github.com/holiman/uint256"
	"time"
)

type Solidity struct {
	*tripod.Tripod
	ethState    *EthState
	cfg         *GethConfig
	stateConfig *config.Config
}

func newEVM(cfg *GethConfig) *vm.EVM {
	txContext := vm.TxContext{
		Origin:     cfg.Origin,
		GasPrice:   cfg.GasPrice,
		BlobHashes: cfg.BlobHashes,
		BlobFeeCap: cfg.BlobFeeCap,
	}
	blockContext := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     cfg.GetHashFn,
		Coinbase:    cfg.Coinbase,
		BlockNumber: cfg.BlockNumber,
		Time:        cfg.Time,
		Difficulty:  cfg.Difficulty,
		GasLimit:    cfg.GasLimit,
		BaseFee:     cfg.BaseFee,
		BlobBaseFee: cfg.BlobBaseFee,
		Random:      cfg.Random,
	}

	return vm.NewEVM(blockContext, txContext, cfg.State, cfg.ChainConfig, cfg.EVMConfig)
}

type GethConfig struct {
	ChainConfig *params.ChainConfig

	// BlockContext provides the EVM with auxiliary information. Once provided
	// it shouldn't be modified.
	GetHashFn   func(n uint64) common.Hash
	Coinbase    common.Address `toml:"coinbase"`
	GasLimit    uint64
	BlockNumber *big.Int
	Time        uint64
	Difficulty  *big.Int
	BaseFee     *big.Int
	BlobBaseFee *big.Int
	Random      *common.Hash

	// TxContext provides the EVM with information about a transaction.
	// All fields can change between transactions.
	Origin     common.Address
	GasPrice   *big.Int
	BlobHashes []common.Hash
	BlobFeeCap *big.Int

	// StateDB gives access to the underlying state
	State *state.StateDB

	// Unknown
	Value     *big.Int
	Debug     bool
	EVMConfig vm.Config

	// Global config
	EnableEthRPC bool   `toml:"enable_eth_rpc"`
	EthHost      string `toml:"eth_host"`
	EthPort      string `toml:"eth_port"`

	GenesisContractCode     string `toml:"genesis_contract_code"`
	GenesisContractDeployer string `toml:"genesis_contract_deployer"`
	GenesisContractAddress  string
}

// sets defaults on the config
func SetDefaultGethConfig() *GethConfig {
	cfg := &GethConfig{
		ChainConfig: params.AllEthashProtocolChanges,
		Difficulty:  big.NewInt(1),
		Origin:      common.HexToAddress("0x0"),
		Coinbase:    common.HexToAddress("0xbaeFe32bc1636a90425AcBCC8cfAD1b0507eCdE1"),
		BlockNumber: big.NewInt(0),
		Time:        0,
		GasLimit:    8000000,
		GasPrice:    big.NewInt(1),
		Value:       big.NewInt(0),
		Debug:       false,
		EVMConfig:   vm.Config{},
		BaseFee:     big.NewInt(params.InitialBaseFee), // 1 gwei
		BlobBaseFee: big.NewInt(params.BlobTxMinBlobGasprice),
		BlobHashes:  []common.Hash{},
		BlobFeeCap:  big.NewInt(0),
		Random:      &common.Hash{},
		State:       nil,
		GetHashFn: func(n uint64) common.Hash {
			return common.BytesToHash(crypto.Keccak256([]byte(new(big.Int).SetUint64(n).String())))
		},
	}

	cfg.ChainConfig.ChainID = big.NewInt(50342)

	return cfg
}

func LoadEvmConfig(fpath string) *GethConfig {
	cfg := SetDefaultGethConfig()
	_, err := toml.DecodeFile(fpath, cfg)
	if err != nil {
		logrus.Fatalf("load config file failed: %v", err)
	}
	return cfg
}

func setDefaultEthStateConfig() *config.Config {
	return &config.Config{
		VMTrace:                 "",
		VMTraceConfig:           "",
		EnablePreimageRecording: false,
		Recovery:                false,
		NoBuild:                 false,
		SnapshotWait:            false,
		SnapshotCache:           128,              // Default cache size
		TrieCleanCache:          256,              // Default Trie cleanup cache size
		TrieDirtyCache:          256,              // Default Trie dirty cache size
		TrieTimeout:             60 * time.Second, // Default Trie timeout
		Preimages:               false,
		NoPruning:               false,
		NoPrefetch:              false,
		StateHistory:            0,                   // By default, there is no state history
		StateScheme:             "hash",              // Default state scheme
		DbPath:                  "verse_db",          // Default database path
		DbType:                  "pebble",            // Default database type
		NameSpace:               "eth/db/chaindata/", // Default namespace
		Ancient:                 "ancient",           // Default ancient data path
		Cache:                   512,                 // Default cache size
		Handles:                 64,                  // Default number of handles
	}
}

func (s *Solidity) InitChain(genesisBlock *yu_types.Block) {
	cfg := s.stateConfig
	genesis := DefaultGoerliGenesisBlock()

	logrus.Printf("Genesis GethConfig: %+v", genesis.Config)
	logrus.Println("Genesis Timestamp: ", genesis.Timestamp)
	logrus.Printf("Genesis ExtraData: %x", genesis.ExtraData)
	logrus.Println("Genesis GasLimit: ", genesis.GasLimit)
	logrus.Println("Genesis Difficulty: ", genesis.Difficulty.String())

	var lastStateRoot common.Hash
	block, err := s.GetCurrentBlock()
	if err != nil && err != yerror.ErrBlockNotFound {
		logrus.Fatal("get current block failed: ", err)
	}
	if block != nil {
		lastStateRoot = common.Hash(block.StateRoot)
	}

	ethState, err := NewEthState(cfg, lastStateRoot)
	if err != nil {
		logrus.Fatal("init NewEthState failed: ", err)
	}
	s.ethState = ethState
	s.cfg.State = ethState.stateDB

	chainConfig, _, err := SetupGenesisBlock(ethState, genesis)
	if err != nil {
		logrus.Fatal("SetupGenesisBlock failed: ", err)
	}

	// s.cfg.ChainConfig = chainConfig

	logrus.Println("Genesis SetupGenesisBlock chainConfig: ", chainConfig)
	logrus.Println("Genesis NewEthState cfg.DbPath: ", ethState.cfg.DbPath)
	logrus.Println("Genesis NewEthState ethState.cfg.NameSpace: ", ethState.cfg.NameSpace)
	logrus.Println("Genesis NewEthState ethState.StateDB.SnapshotCommits: ", ethState.stateDB)
	logrus.Println("Genesis NewEthState ethState.stateCache: ", ethState.stateCache)
	logrus.Println("Genesis NewEthState ethState.trieDB: ", ethState.trieDB)

	// commit genesis state
	genesisStateRoot, err := s.ethState.GenesisCommit()
	if err != nil {
		logrus.Fatal("genesis state commit failed: ", err)
	}

	genesisBlock.StateRoot = yu_common.Hash(genesisStateRoot)

	if s.cfg.GenesisContractCode != "" {
		initContract(s)
	}

	// Deploy Random Contact
	// TODO: load contract code from config
	initRandomContract(s)
}

func initContract(s *Solidity) {
	createContractInput := s.cfg.GenesisContractCode
	createContractInputByt, _ := hexutil.Decode(createContractInput)
	createContractTx := &TxRequest{
		Origin:   common.HexToAddress(s.cfg.GenesisContractDeployer),
		Input:    createContractInputByt,
		GasLimit: 10000000,
		Value:    big.NewInt(0),
		GasPrice: big.NewInt(2000000000),
	}
	_, contractAddr, _ := initRunTxReq(s, createContractTx)
	s.cfg.GenesisContractAddress = contractAddr.Hex()
}

func initRandomContract(s *Solidity) {
	createContractInput := "0x608060405234801561001057600080fd5b5060b28061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063aacc5a1714602d575b600080fd5b60336047565b604051603e91906063565b60405180910390f35b600090565b6000819050919050565b605d81604c565b82525050565b6000602082019050607660008301846056565b9291505056fea26469706673582212203e9972e2f120f68458fbfe81c2f6cd126ae8991b3c663c231980aef4a150bcf064736f6c63430008140033"
	createContractInputByt, _ := hexutil.Decode(createContractInput)
	createContractTx := &TxRequest{
		Origin:   common.HexToAddress(s.cfg.GenesisContractDeployer),
		Input:    createContractInputByt,
		GasLimit: 10000000,
		Value:    big.NewInt(0),
		GasPrice: big.NewInt(2000000000),
	}
	_, randomContractAddr, _ := initRunTxReq(s, createContractTx)

	logrus.Printf("[initRandomContract] Random Contract Addr = %v", randomContractAddr.Hex())
	s.cfg.ChainConfig.RandomContractAddr = *randomContractAddr
}

func initRunTxReq(s *Solidity, txReq *TxRequest) ([]byte, *common.Address, error) {
	vmenv := newEVM(s.cfg)
	//s.ethState.setTxContext()
	vmenv.StateDB = s.ethState.stateDB
	sender := vm.AccountRef(txReq.Origin)
	rules := s.cfg.ChainConfig.Rules(vmenv.Context.BlockNumber, vmenv.Context.Random != nil, vmenv.Context.Time)

	cfg := s.cfg
	if txReq.Address == nil {
		if cfg.EVMConfig.Tracer != nil && cfg.EVMConfig.Tracer.OnTxStart != nil {
			cfg.EVMConfig.Tracer.OnTxStart(vmenv.GetVMContext(), types.NewTx(&types.LegacyTx{Data: txReq.Input, Value: txReq.Value, Gas: txReq.GasLimit}), txReq.Origin)
		}

		s.ethState.Prepare(rules, cfg.Origin, cfg.Coinbase, nil, vm.ActivePrecompiles(rules), nil)

		code, address, leftOverGas, err := vmenv.Create(sender, txReq.Input, txReq.GasLimit, uint256.MustFromBig(txReq.Value))
		if err != nil {
			byt, _ := json.Marshal(txReq)
			logrus.Printf("[Execute Txn] Create contract Failed. err = %v. Request = %v", err, string(byt))
		}

		logrus.Printf("[Execute Txn] Create contract success. Address = %v, Left Gas = %v", address.Hex(), leftOverGas)

		return code, &address, err
	} else {
		if cfg.EVMConfig.Tracer != nil && cfg.EVMConfig.Tracer.OnTxStart != nil {
			cfg.EVMConfig.Tracer.OnTxStart(vmenv.GetVMContext(), types.NewTx(&types.LegacyTx{To: txReq.Address, Data: txReq.Input, Value: txReq.Value, Gas: txReq.GasLimit}), txReq.Origin)
		}

		s.ethState.Prepare(rules, cfg.Origin, cfg.Coinbase, txReq.Address, vm.ActivePrecompiles(rules), nil)

		code, leftOverGas, err := vmenv.Call(sender, *txReq.Address, txReq.Input, txReq.GasLimit, uint256.MustFromBig(txReq.Value))
		if err != nil {
			byt, _ := json.Marshal(txReq)
			logrus.Printf("[Execute Txn] SendTx Failed. err = %v. Request = %v", err, string(byt))
		}

		logrus.Printf("[Execute Txn] SendTx success. Hex Code = %v, Left Gas = %v", hex.EncodeToString(code), leftOverGas)
		return code, nil, err
	}
}

func NewSolidity(gethConfig *GethConfig) *Solidity {
	ethStateConfig := setDefaultEthStateConfig()

	solidity := &Solidity{
		Tripod:      tripod.NewTripod(),
		cfg:         gethConfig,
		stateConfig: ethStateConfig,
		// network:       utils.Network(cfg.Network),
	}

	solidity.SetWritings(solidity.ExecuteTxn)
	solidity.SetReadings(
		solidity.Call, solidity.GetReceipt, solidity.GetReceipts,
		// solidity.GetClass, solidity.GetClassAt,
		// 	solidity.GetClassHashAt, solidity.GetNonce, solidity.GetStorage,
		// 	solidity.GetTransaction, solidity.GetTransactionStatus,
		// 	solidity.SimulateTransactions,
		// 	solidity.GetBlockWithTxs, solidity.GetBlockWithTxHashes,
	)

	return solidity
}

// region ---- Tripod Api ----

func (s *Solidity) StartBlock(block *yu_types.Block) {
	s.cfg.BlockNumber = big.NewInt(int64(block.Height))
	//s.cfg.GasLimit =
	s.cfg.Time = block.Timestamp
	s.cfg.Difficulty = big.NewInt(int64(block.Difficulty))
}

func (s *Solidity) EndBlock(block *yu_types.Block) {
	// nothing
}

func (s *Solidity) FinalizeBlock(block *yu_types.Block) {
	// nothing
}

func (s *Solidity) PreHandleTxn(txn *yu_types.SignedTxn) error {
	var txReq TxRequest
	param := txn.GetParams()
	err := json.Unmarshal([]byte(param), &txReq)
	if err != nil {
		return err
	}

	yuHash, err := ConvertHashToYuHash(txReq.Hash)
	if err != nil {
		return err
	}
	txn.TxnHash = yuHash

	return nil
}

// ExecuteTxn executes the code using the input as call data during the execution.
// It returns the EVM's return value, the new state and an error if it failed.
//
// Execute sets up an in-memory, temporary, environment for the execution of
// the given code. It makes sure that it's restored to its original state afterwards.
func (s *Solidity) ExecuteTxn(ctx *context.WriteContext) error {
	txReq := new(TxRequest)
	err := ctx.BindJson(txReq)
	if err != nil {
		return err
	}

	origin := txReq.Origin
	gasLimit := txReq.GasLimit
	gasPrice := txReq.GasPrice
	value := txReq.Value

	cfg := s.cfg
	ethstate := s.ethState

	cfg.Origin = origin
	cfg.GasLimit = gasLimit
	cfg.GasPrice = gasPrice
	cfg.Value = value

	vmenv := newEVM(cfg)
	s.ethState.setTxContext(common.Hash(ctx.GetTxnHash()), ctx.TxnIndex)
	vmenv.StateDB = s.ethState.stateDB
	logrus.Printf("[StateDB] %v", s.ethState.stateDB == s.cfg.State)

	sender := vm.AccountRef(txReq.Origin)
	rules := cfg.ChainConfig.Rules(vmenv.Context.BlockNumber, vmenv.Context.Random != nil, vmenv.Context.Time)

	//Calculate whether the maximum consumption meets the balance
	err = canTransfer(gasLimit, txReq, gasPrice, s)
	if err != nil {
		_ = saveReceipt(ctx, vmenv, txReq, common.Address{}, 0, err)
		return err
	}

	if txReq.Address == nil {
		if cfg.EVMConfig.Tracer != nil && cfg.EVMConfig.Tracer.OnTxStart != nil {
			cfg.EVMConfig.Tracer.OnTxStart(vmenv.GetVMContext(), types.NewTx(&types.LegacyTx{Data: txReq.Input, Value: txReq.Value, Gas: txReq.GasLimit}), txReq.Origin)
		}

		ethstate.Prepare(rules, cfg.Origin, cfg.Coinbase, nil, vm.ActivePrecompiles(rules), nil)

		_, address, leftOverGas, err := vmenv.Create(sender, txReq.Input, txReq.GasLimit, uint256.MustFromBig(txReq.Value))
		if err != nil {
			byt, _ := json.Marshal(txReq)
			logrus.Printf("[Execute Txn] Create contract Failed. err = %v. Request = %v", err, string(byt))
			_ = saveReceipt(ctx, vmenv, txReq, address, leftOverGas, err)
			_ = calculateGasFee(gasLimit, leftOverGas, err, gasPrice, cfg, txReq, s)
			return err
		}

		logrus.Printf("[Execute Txn] Create contract. contractAddress = %v, Left Gas = %v", address.Hex(), leftOverGas)
		err = saveReceipt(ctx, vmenv, txReq, address, leftOverGas, err)
		if err != nil {
			return err
		}

		err = calculateGasFee(gasLimit, leftOverGas, err, gasPrice, cfg, txReq, s)
		if err != nil {
			return err
		}

	} else {
		if cfg.EVMConfig.Tracer != nil && cfg.EVMConfig.Tracer.OnTxStart != nil {
			cfg.EVMConfig.Tracer.OnTxStart(vmenv.GetVMContext(), types.NewTx(&types.LegacyTx{To: txReq.Address, Data: txReq.Input, Value: txReq.Value, Gas: txReq.GasLimit}), txReq.Origin)
		}

		ethstate.Prepare(rules, cfg.Origin, cfg.Coinbase, txReq.Address, vm.ActivePrecompiles(rules), nil)
		ethstate.SetNonce(txReq.Origin, ethstate.GetNonce(sender.Address())+1)

		logrus.Printf("before transfer: account %s balance %d \n", sender.Address(), ethstate.stateDB.GetBalance(sender.Address()))

		code, leftOverGas, err := vmenv.Call(sender, *txReq.Address, txReq.Input, txReq.GasLimit, uint256.MustFromBig(txReq.Value))
		logrus.Printf("after transfer: account %s balance %d \n", sender.Address(), ethstate.stateDB.GetBalance(sender.Address()))
		if err != nil {
			byt, _ := json.Marshal(txReq)
			logrus.Printf("[Execute Txn] SendTx Failed. err = %v. Request = %v", err, string(byt))
			_ = saveReceipt(ctx, vmenv, txReq, common.Address{}, leftOverGas, err)
			_ = calculateGasFee(gasLimit, leftOverGas, err, gasPrice, cfg, txReq, s)
			return err
		}

		logrus.Printf("[Execute Txn] SendTx. Hex Code = %v, Left Gas = %v", hex.EncodeToString(code), leftOverGas)
		err = saveReceipt(ctx, vmenv, txReq, common.Address{}, leftOverGas, err)
		if err != nil {
			return err
		}

		err = calculateGasFee(gasLimit, leftOverGas, err, gasPrice, cfg, txReq, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func calculateGasFee(gasLimit uint64, leftOverGas uint64, err error, gasPrice *big.Int, cfg *GethConfig, txReq *TxRequest, s *Solidity) error {
	gasUsed := gasLimit - leftOverGas
	usdtPricePerGasUnit, err := GetUSDTPricePerGasUnit()
	if err != nil {
		return err
	}
	gasfee := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice)
	gasfeeFloat := new(big.Float).SetInt(gasfee)
	usdtPricePerGasUnitFloat := new(big.Float).SetInt(usdtPricePerGasUnit)
	gasFeeFloat := new(big.Float).Quo(gasfeeFloat, usdtPricePerGasUnitFloat)
	gasFeeFloat = new(big.Float).SetPrec(18).Set(gasFeeFloat)

	logrus.Printf("[Execute Txn] gasfee = %v, gasUsed = %v,gasPrice = %v,usdtPricePerGasUnit = %v", gasFeeFloat, gasUsed, gasPrice, usdtPricePerGasUnit)
	tenPow18 := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	tenPow18Float := new(big.Float).SetInt(tenPow18)
	gasFeeInWeiFloat := new(big.Float).Mul(gasFeeFloat, tenPow18Float)
	gasFeeInWei := new(big.Int)
	gasFeeInWeiFloat.Int(gasFeeInWei)

	transferTx := constructTransferTxInput(cfg, gasLimit, gasPrice, gasFeeInWei, txReq.Origin, cfg.Coinbase)
	_, _, err = initRunTxReq(s, transferTx)
	if err != nil {
		logrus.Printf("[Execute Txn] Expend gas fail. cfg.Coinbase = %v, gasFeeInWei = %v,gasFeeInWeiFloat = %v", cfg.Coinbase, gasFeeInWei, gasFeeInWeiFloat)
		return err
	}

	logrus.Printf("[Execute Txn] Expend gas success. cfg.Coinbase = %v, gasFeeInWei = %v,gasFeeInWeiFloat = %v", cfg.Coinbase, gasFeeInWei, gasFeeInWeiFloat)
	return nil
}

func canTransfer(gasLimit uint64, txReq *TxRequest, gasPrice *big.Int, s *Solidity) error {
	balanceOfSelector := "balanceOf(address)"

	usdtPricePerGasUnit, err := GetUSDTPricePerGasUnit()
	if err != nil {
		return err
	}
	gasfee := new(big.Int).Mul(new(big.Int).SetUint64(gasLimit), gasPrice)
	gasfeeFloat := new(big.Float).SetInt(gasfee)
	usdtPricePerGasUnitFloat := new(big.Float).SetInt(usdtPricePerGasUnit)
	gasFeeFloat := new(big.Float).Quo(gasfeeFloat, usdtPricePerGasUnitFloat)
	gasFeeFloat = new(big.Float).SetPrec(18).Set(gasFeeFloat)

	tenPow18 := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	tenPow18Float := new(big.Float).SetInt(tenPow18)
	gasFeeInWeiFloat := new(big.Float).Mul(gasFeeFloat, tenPow18Float)
	gasFeeInWei := new(big.Int)
	gasFeeInWeiFloat.Int(gasFeeInWei)

	paddedAddress := padLeft(txReq.Origin.Hex()[2:], 64)

	hash := crypto.Keccak256Hash([]byte(balanceOfSelector))
	functionSelector := hash.Hex()[0:10]

	balanceOfInput := functionSelector + paddedAddress
	balanceOfInputByt, _ := hexutil.Decode(balanceOfInput)

	contractAddr := common.HexToAddress(s.cfg.GenesisContractAddress)
	balanceOfTx := &TxRequest{
		Origin:   txReq.Origin,
		Address:  &contractAddr,
		Value:    big.NewInt(0),
		Input:    balanceOfInputByt,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
	}
	code, _, err := initRunTxReq(s, balanceOfTx)
	if err != nil {
		logrus.Printf("[Execute Txn] Get balanceOf fail.")
	}
	balanceHex := hex.EncodeToString(code)
	balance, err := hexStringToBigInt(balanceHex)
	if err != nil {
		return fmt.Errorf("[Execute Txn] Failed to convert balance: %v", err)
	}
	logrus.Printf("[Execute Txn] Get balance: %v", balance.String())

	if balance.Cmp(gasFeeInWei) < 0 {
		return fmt.Errorf("[Execute Txn] Insufficient balance: balance = %s, required = %s", balance.String(), gasFeeInWei.String())
	}

	return nil
}

func hexStringToBigInt(hexStr string) (*big.Int, error) {
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %v", err)
	}

	result := new(big.Int).SetBytes(decoded)
	return result, nil
}

// Function to construct the ERC-20 transfer transaction input data
func constructTransferTxInput(cfg *GethConfig, gasLimit uint64, gasPrice *big.Int, gasFee *big.Int, originAddress common.Address, recipientAddress common.Address) *TxRequest {
	// 1. Define the function signature for ERC-20 `transfer` function
	functionSignature := "transfer(address,uint256)"

	// 2. Calculate the Keccak-256 hash and get the first 4 bytes (function selector)
	hash := crypto.Keccak256Hash([]byte(functionSignature))
	functionSelector := hash.Hex()[0:10]

	// 3. Define the recipient address
	recipient := recipientAddress.Hex()[2:]

	// 4. Define the amount to transfer (in Wei)
	amount := new(big.Int)
	amount.SetString(ConvertBigIntToUint256(gasFee).String(), 10)

	// 5. Encode the recipient and amount
	recipientPadded := padLeft(recipient, 64)
	amountPadded := padLeft(amount.Text(16), 64)

	// 6. Construct the transfer input data
	// 5. Construct the transfer input data
	transferTxInput := functionSelector + recipientPadded + amountPadded

	fmt.Println("Transfer Input Data:", transferTxInput)

	contractAddr := common.HexToAddress(cfg.GenesisContractAddress)
	transferTxInputByt, _ := hexutil.Decode(transferTxInput)
	transferTx := &TxRequest{
		Origin:   originAddress,
		Address:  &contractAddr,
		Value:    big.NewInt(0),
		Input:    transferTxInputByt,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
	}

	return transferTx
}

// Helper function to pad strings with leading zeros
func padLeft(str string, length int) string {
	return fmt.Sprintf("%0*s", length, str)
}

func saveReceipt(ctx *context.WriteContext, vmEvm *vm.EVM, txReq *TxRequest, contractAddr common.Address, leftOverGas uint64, err error) error {
	evmReceipt := makeEvmReceipt(vmEvm, ctx.Txn, ctx.Block, contractAddr, leftOverGas, err)
	receiptByt, err := json.Marshal(evmReceipt)
	if err != nil {
		txReqByt, _ := json.Marshal(txReq)
		logrus.Errorf("Receipt marshal err: %v. Tx: %v", err, string(txReqByt))
		return err
	}
	ctx.EmitExtra(receiptByt)
	return nil
}

// Call executes the code given by the contract's address. It will return the
// EVM's return value or an error if it failed.
func (s *Solidity) Call(ctx *context.ReadContext) {
	callReq := new(CallRequest)
	err := ctx.BindJson(callReq)
	if err != nil {
		ctx.Json(http.StatusBadRequest, &CallResponse{Err: jsonrpc.Err(jsonrpc.InvalidJSON, err.Error())})
		return
	}

	cfg := s.cfg
	address := callReq.Address
	input := callReq.Input
	origin := callReq.Origin
	gasLimit := callReq.GasLimit
	gasPrice := callReq.GasPrice
	value := callReq.Value

	cfg.Origin = origin
	cfg.GasLimit = gasLimit
	cfg.GasPrice = gasPrice
	cfg.Value = value

	var (
		vmenv    = newEVM(cfg)
		sender   = vm.AccountRef(origin)
		ethState = s.ethState
		rules    = cfg.ChainConfig.Rules(vmenv.Context.BlockNumber, vmenv.Context.Random != nil, vmenv.Context.Time)
	)

	logrus.Printf("[StateDB] %v", s.ethState.stateDB == s.cfg.State)
	vmenv.StateDB = s.ethState.stateDB

	if cfg.EVMConfig.Tracer != nil && cfg.EVMConfig.Tracer.OnTxStart != nil {
		cfg.EVMConfig.Tracer.OnTxStart(vmenv.GetVMContext(), types.NewTx(&types.LegacyTx{To: &address, Data: input, Value: value, Gas: gasLimit}), origin)
	}
	// Execute the preparatory steps for state transition which includes:
	// - prepare accessList(post-berlin)
	// - reset transient storage(eip 1153)
	ethState.Prepare(rules, origin, cfg.Coinbase, &address, vm.ActivePrecompiles(rules), nil)

	// Call the code with the given configuration.
	ret, leftOverGas, err := vmenv.Call(
		sender,
		address,
		input,
		gasLimit,
		uint256.MustFromBig(value),
	)

	logrus.Printf("[Call] Request from = %v, to = %v, gasLimit = %v, value = %v, input = %v", sender.Address().Hex(), address.Hex(), gasLimit, value.Uint64(), hex.EncodeToString(input))
	logrus.Printf("[Call] Response: Origin Code = %v, Hex Code = %v, String Code = %v, LeftOverGas = %v", ret, hex.EncodeToString(ret), new(big.Int).SetBytes(ret).String(), leftOverGas)

	if err != nil {
		ctx.Json(http.StatusInternalServerError, &CallResponse{Err: jsonrpc.Err(jsonrpc.InternalError, err.Error())})
		return
	}
	result := CallResponse{Ret: ret, LeftOverGas: leftOverGas}
	json, _ := json.Marshal(result)
	fmt.Printf("[ETH_CALL] eth return result is %v\n", string(json))
	ctx.JsonOk(&result)
}

func (s *Solidity) Commit(block *yu_types.Block) {
	blockNumber := uint64(block.Height)
	stateRoot, err := s.ethState.Commit(blockNumber)
	if err != nil {
		logrus.Errorf("Solidity commit failed on Block(%d), error: %v", blockNumber, err)
		return
	}
	block.StateRoot = AdaptHash(stateRoot)
}

func AdaptHash(ethHash common.Hash) yu_common.Hash {
	var yuHash yu_common.Hash
	copy(yuHash[:], ethHash[:])
	return yuHash
}

func makeEvmReceipt(vmEvm *vm.EVM, signedTx *yu_types.SignedTxn, block *yu_types.Block, address common.Address, leftOverGas uint64, err error) *types.Receipt {
	wrCallParams := signedTx.Raw.WrCall.Params
	var txReq = &TxRequest{}
	_ = json.Unmarshal([]byte(wrCallParams), txReq)

	txArgs := &TempTransactionArgs{}
	_ = json.Unmarshal(txReq.OriginArgs, txArgs)
	originTx := txArgs.ToTransaction(txReq.V, txReq.R, txReq.S)

	usedGas := originTx.Gas() - leftOverGas

	blockNumber := big.NewInt(int64(block.Height))
	txHash := common.Hash(signedTx.TxnHash)
	effectiveGasPrice := big.NewInt(1000000000) // 1 GWei

	status := types.ReceiptStatusFailed
	if err == nil {
		status = types.ReceiptStatusSuccessful
	}
	var root []byte
	stateDB := vmEvm.StateDB.(*state.StateDB)
	if vmEvm.ChainConfig().IsByzantium(blockNumber) {
		stateDB.Finalise(true)
	} else {
		root = stateDB.IntermediateRoot(vmEvm.ChainConfig().IsEIP158(blockNumber)).Bytes()
	}

	// TODO: 1. root is nil; 2. CumulativeGasUsed not; 3. Log is empty; 4. logBloom is empty

	receipt := &types.Receipt{
		Type:              originTx.Type(),
		Status:            status,
		PostState:         root,
		CumulativeGasUsed: leftOverGas,
		TxHash:            txHash,
		ContractAddress:   address,
		GasUsed:           usedGas,
		EffectiveGasPrice: effectiveGasPrice,
	}

	if originTx.Type() == types.BlobTxType {
		receipt.BlobGasUsed = uint64(len(originTx.BlobHashes()) * params.BlobTxBlobGasPerBlob)
		receipt.BlobGasPrice = vmEvm.Context.BlobBaseFee
	}

	receipt.Logs = stateDB.GetLogs(txHash, uint64(block.Height), common.Hash(block.Hash))
	receipt.Bloom = types.CreateBloom(types.Receipts{})
	receipt.BlockHash = common.Hash(block.Hash)
	receipt.BlockNumber = blockNumber
	receipt.TransactionIndex = uint(stateDB.TxIndex())

	logrus.Printf("[Receipt] uint64(block.Height) = %v", uint64(block.Height))
	logrus.Printf("[Receipt] txHash = %v", txHash)
	logrus.Printf("[Receipt] common.Hash(block.Hash) = %v", common.Hash(block.Hash))
	logrus.Printf("[Receipt] block.Hash = %v", block.Hash)
	logrus.Printf("[Receipt] log = %v", receipt.Logs)
	//spew.Dump("[Receipt] log = %v", stateDB.Logs())
	//logrus.Printf("[Receipt] log is nil = %v", receipt.Logs == nil)
	if receipt.Logs == nil {
		receipt.Logs = []*types.Log{}
	}

	for idx, txn := range block.Txns {
		if common.Hash(txn.TxnHash) == txHash {
			receipt.TransactionIndex = uint(idx)
		}
	}
	logrus.Printf("[Receipt] statedb txIndex = %v, actual txIndex = %v", stateDB.TxIndex(), receipt.TransactionIndex)

	return receipt
}

func (s *Solidity) StateAt(root common.Hash) (*state.StateDB, error) {
	return s.ethState.StateAt(root)
}

func (s *Solidity) GetEthDB() ethdb.Database {
	return s.ethState.ethDB
}

type ReceiptRequest struct {
	Hash common.Hash `json:"hash"`
}

type ReceiptResponse struct {
	Receipt *types.Receipt `json:"receipt"`
	Err     error          `json:"err"`
}

type ReceiptsRequest struct {
	Hashes []common.Hash `json:"hashes"`
}

type ReceiptsResponse struct {
	Receipts []*types.Receipt `json:"receipts"`
	Err      error            `json:"err"`
}

func (s *Solidity) GetReceipt(ctx *context.ReadContext) {
	var rq ReceiptRequest
	err := ctx.BindJson(&rq)
	if err != nil {
		ctx.Json(http.StatusBadRequest, &ReceiptResponse{Err: err})
		return
	}

	receipt, err := s.getReceipt(rq.Hash)
	if err != nil {
		ctx.Json(http.StatusInternalServerError, &ReceiptResponse{Err: err})
		return
	}

	ctx.JsonOk(&ReceiptResponse{Receipt: receipt})
}

func (s *Solidity) getReceipt(hash common.Hash) (*types.Receipt, error) {
	yuHash, err := ConvertHashToYuHash(hash)
	if err != nil {
		return nil, err
	}
	yuReceipt, err := s.TxDB.GetReceipt(yuHash)
	if err != nil {
		return nil, err
	}
	if yuReceipt == nil {
		return nil, ErrNotFoundReceipt
	}
	receipt := new(types.Receipt)
	err = json.Unmarshal(yuReceipt.Extra, receipt)
	return receipt, err
}

func (s *Solidity) GetReceipts(ctx *context.ReadContext) {
	var rq ReceiptsRequest
	err := ctx.BindJson(&rq)
	if err != nil {
		ctx.Json(http.StatusBadRequest, &ReceiptsResponse{Err: err})
		return
	}

	receipts := make([]*types.Receipt, 0, len(rq.Hashes))
	for _, hash := range rq.Hashes {
		receipt, err := s.getReceipt(hash)
		if err != nil {
			ctx.Json(http.StatusInternalServerError, &ReceiptsResponse{Err: err})
			return
		}

		receipts = append(receipts, receipt)
	}

	ctx.JsonOk(&ReceiptsResponse{Receipts: receipts})
}

// endregion ---- Tripod Api ----
