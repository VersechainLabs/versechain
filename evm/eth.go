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

	initContract(s)
}

func initContract(s *Solidity) {
	createContractInput := "0x608060405234801561001057600080fd5b506040518060400160405280600981526020016805465737445726332360bc1b8152506040518060400160405280600381526020016215115560ea1b815250816003908161005e9190610114565b50600461006b8282610114565b5050506101d3565b634e487b7160e01b600052604160045260246000fd5b600181811c9082168061009d57607f821691505b6020821081036100bd57634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111561010f576000816000526020600020601f850160051c810160208610156100ec5750805b601f850160051c820191505b8181101561010b578281556001016100f8565b5050505b505050565b81516001600160401b0381111561012d5761012d610073565b6101418161013b8454610089565b846100c3565b602080601f831160018114610176576000841561015e5750858301515b600019600386901b1c1916600185901b17855561010b565b600085815260208120601f198616915b828110156101a557888601518255948401946001909101908401610186565b50858210156101c35787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6108d2806101e26000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c806340c10f191161009757806398c41f691161006657806398c41f69146101cb578063a9059cbb146101de578063b8b085f2146101f1578063dd62ed3e146101f957600080fd5b806340c10f191461017f5780635e863fa71461019257806370a082311461019a57806395d89b41146101c357600080fd5b806323b872dd116100d357806323b872dd1461014d578063313ce5671461016057806332b7a7611461016f578063378ec23b1461017957600080fd5b806306fdde03146100fa578063095ea7b31461011857806318160ddd1461013b575b600080fd5b610102610232565b60405161010f9190610702565b60405180910390f35b61012b61012636600461076d565b6102c4565b604051901515815260200161010f565b6002545b60405190815260200161010f565b61012b61015b366004610797565b6102de565b6040516012815260200161010f565b610177610302565b005b4361013f565b61017761018d36600461076d565b61033e565b61017761034c565b61013f6101a83660046107d3565b6001600160a01b031660009081526020819052604090205490565b610102610377565b6101776101d93660046107f5565b610386565b61012b6101ec36600461076d565b6103cb565b60055461013f565b61013f61020736600461080e565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b60606003805461024190610841565b80601f016020809104026020016040519081016040528092919081815260200182805461026d90610841565b80156102ba5780601f1061028f576101008083540402835291602001916102ba565b820191906000526020600020905b81548152906001019060200180831161029d57829003601f168201915b5050505050905090565b6000336102d28185856103d9565b60019150505b92915050565b6000336102ec8582856103eb565b6102f785858561046e565b506001949350505050565b60405143808252339182907fe7a5a2788502413da105ba7eebbd320341bcf15413708e8261856f535531eabd9060200160405180910390a25050565b61034882826104cd565b5050565b6040517fcf16a92280c1bbb43f72d31126b724d508df2877835849e8744017ab36a9b47f90600090a1565b60606004805461024190610841565b600580549082905560408051828152602081018490527f28bab7182b1f3ed84e1006088a7e959b8999f63af651809d61a39e97b026fcde910160405180910390a15050565b6000336102d281858561046e565b6103e68383836001610503565b505050565b6001600160a01b038381166000908152600160209081526040808320938616835292905220546000198114610468578181101561045957604051637dc7a0d960e11b81526001600160a01b038416600482015260248101829052604481018390526064015b60405180910390fd5b61046884848484036000610503565b50505050565b6001600160a01b03831661049857604051634b637e8f60e11b815260006004820152602401610450565b6001600160a01b0382166104c25760405163ec442f0560e01b815260006004820152602401610450565b6103e68383836105d8565b6001600160a01b0382166104f75760405163ec442f0560e01b815260006004820152602401610450565b610348600083836105d8565b6001600160a01b03841661052d5760405163e602df0560e01b815260006004820152602401610450565b6001600160a01b03831661055757604051634a1406b160e11b815260006004820152602401610450565b6001600160a01b038085166000908152600160209081526040808320938716835292905220829055801561046857826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516105ca91815260200190565b60405180910390a350505050565b6001600160a01b0383166106035780600260008282546105f8919061087b565b909155506106759050565b6001600160a01b038316600090815260208190526040902054818110156106565760405163391434e360e21b81526001600160a01b03851660048201526024810182905260448101839052606401610450565b6001600160a01b03841660009081526020819052604090209082900390555b6001600160a01b038216610691576002805482900390556106b0565b6001600160a01b03821660009081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516106f591815260200190565b60405180910390a3505050565b60006020808352835180602085015260005b8181101561073057858101830151858201604001528201610714565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b038116811461076857600080fd5b919050565b6000806040838503121561078057600080fd5b61078983610751565b946020939093013593505050565b6000806000606084860312156107ac57600080fd5b6107b584610751565b92506107c360208501610751565b9150604084013590509250925092565b6000602082840312156107e557600080fd5b6107ee82610751565b9392505050565b60006020828403121561080757600080fd5b5035919050565b6000806040838503121561082157600080fd5b61082a83610751565b915061083860208401610751565b90509250929050565b600181811c9082168061085557607f821691505b60208210810361087557634e487b7160e01b600052602260045260246000fd5b50919050565b808201808211156102d857634e487b7160e01b600052601160045260246000fdfea2646970667358221220e9e543bc723021d52e011a4d87a13a4acf4c74bcee5c6af32278901304b9b82464736f6c63430008180033"
	createContractInputByt, _ := hexutil.Decode(createContractInput)
	createContractTx := &TxRequest{
		Origin:   common.HexToAddress("0x7Bd36074b61Cfe75a53e1B9DF7678C96E6463b02"),
		Input:    createContractInputByt,
		GasLimit: 10000000,
		Value:    big.NewInt(0),
		GasPrice: big.NewInt(2000000000),
	}
	initRunTxReq(s, createContractTx)

	contractAddr := common.HexToAddress("0x310b8685e3e69cb05b251a12f5ffab23001cda42")
	mintTxInput := "0x40c10f190000000000000000000000007bd36074b61cfe75a53e1b9df7678c96e6463b020000000000000000000000000000000000000000000000056bc75e2d63100000"
	mintTxInputByt, _ := hexutil.Decode(mintTxInput)
	mintTx := &TxRequest{
		Origin:   common.HexToAddress("0x7Bd36074b61Cfe75a53e1B9DF7678C96E6463b02"),
		Address:  &contractAddr,
		Value:    big.NewInt(0),
		Input:    mintTxInputByt,
		GasLimit: 10000000,
		GasPrice: big.NewInt(2000000000),
	}
	initRunTxReq(s, mintTx)

}

func initRunTxReq(s *Solidity, txReq *TxRequest) {
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

		_, address, leftOverGas, err := vmenv.Create(sender, txReq.Input, txReq.GasLimit, uint256.MustFromBig(txReq.Value))
		if err != nil {
			byt, _ := json.Marshal(txReq)
			logrus.Printf("[Execute Txn] Create contract Failed. err = %v. Request = %v", err, string(byt))
		}

		logrus.Printf("[Execute Txn] Create contract success. Address = %v, Left Gas = %v", address.Hex(), leftOverGas)
	} else {
		if cfg.EVMConfig.Tracer != nil && cfg.EVMConfig.Tracer.OnTxStart != nil {
			cfg.EVMConfig.Tracer.OnTxStart(vmenv.GetVMContext(), types.NewTx(&types.LegacyTx{To: txReq.Address, Data: txReq.Input, Value: txReq.Value, Gas: txReq.GasLimit}), txReq.Origin)
		}

		s.ethState.Prepare(rules, cfg.Origin, cfg.Coinbase, txReq.Address, vm.ActivePrecompiles(rules), nil)
		s.ethState.SetNonce(txReq.Origin, 1)

		code, leftOverGas, err := vmenv.Call(sender, *txReq.Address, txReq.Input, txReq.GasLimit, uint256.MustFromBig(txReq.Value))
		if err != nil {
			byt, _ := json.Marshal(txReq)
			logrus.Printf("[Execute Txn] SendTx Failed. err = %v. Request = %v", err, string(byt))
		}

		logrus.Printf("[Execute Txn] SendTx success. Hex Code = %v, Left Gas = %v", hex.EncodeToString(code), leftOverGas)
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
			return err
		}

		logrus.Printf("[Execute Txn] Create contract success. Address = %v, Left Gas = %v", address.Hex(), leftOverGas)
		err = saveReceipt(ctx, vmenv, txReq, address, leftOverGas, err)
		if err != nil {
			return err
		}

		gasUsed := gasLimit - leftOverGas
		gasfee := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice)
		logrus.Printf("[Execute Txn] gasfee = %v, gasUsed = %v,gasPrice = %v", gasfee, gasUsed, gasPrice)
		cTransfer := ethstate.CanTransfer(sender.Address(), ConvertBigIntToUint256(gasfee))
		if !cTransfer {
			logrus.Printf("[Execute Txn] Insufficient Balance.sender balance : %v,", ethstate.stateDB.GetBalance(sender.Address()))
			return nil
		}
		ethstate.Transfer(sender.Address(), cfg.Coinbase, ConvertBigIntToUint256(gasfee))
		logrus.Printf("[Execute Txn] Create contract success. cfg.Coinbase = %v, gasfee = %v", cfg.Coinbase, gasfee)

		// 1. Define the function signature for ERC-20 `transfer` function
		functionSignature := "transfer(address,uint256)"

		// 2. Calculate the Keccak-256 hash and get the first 4 bytes (function selector)
		hash := crypto.Keccak256Hash([]byte(functionSignature))
		functionSelector := hash.Hex()[0:10]

		// 3. Define the recipient address
		recipient := cfg.Coinbase.Hex()[2:]

		// 4. Define the amount to transfer (in Wei)
		amount := new(big.Int)
		amount.SetString(ConvertBigIntToUint256(gasfee).String(), 10)

		// 5. Encode the recipient and amount
		recipientPadded := padLeft(recipient, 64)
		amountPadded := padLeft(amount.Text(16), 64)

		// 6. Construct the transfer input data
		transferTxInput := functionSelector + recipientPadded + amountPadded

		fmt.Println("Transfer Input Data:", transferTxInput)

		contractAddr := common.HexToAddress("0x310b8685e3e69cb05b251a12f5ffab23001cda42")
		transferTxInputByt, _ := hexutil.Decode(transferTxInput)
		transferTx := &TxRequest{
			Origin:   common.HexToAddress("0x7Bd36074b61Cfe75a53e1B9DF7678C96E6463b02"),
			Address:  &contractAddr,
			Value:    big.NewInt(0),
			Input:    transferTxInputByt,
			GasLimit: gasLimit,
			GasPrice: gasPrice,
		}
		initRunTxReq(s, transferTx)

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
			return err
		}

		logrus.Printf("[Execute Txn] SendTx success. Hex Code = %v, Left Gas = %v", hex.EncodeToString(code), leftOverGas)
		err = saveReceipt(ctx, vmenv, txReq, common.Address{}, leftOverGas, err)
		if err != nil {
			return err
		}

		gasUsed := gasLimit - leftOverGas
		logrus.Printf("[Execute Txn] gasLimit = %v, leftOverGas = %v", gasLimit, leftOverGas)
		gasfee := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice)
		logrus.Printf("[Execute Txn] gasfee = %v, gasUsed = %v,gasPrice = %v", gasfee, gasUsed, gasPrice)
		cTransfer := ethstate.CanTransfer(sender.Address(), ConvertBigIntToUint256(gasfee))
		if !cTransfer {
			logrus.Printf("[Execute Txn] Insufficient Balance.sender balance : %v,", ethstate.stateDB.GetBalance(sender.Address()))
			return nil
		}
		ethstate.Transfer(sender.Address(), cfg.Coinbase, ConvertBigIntToUint256(gasfee))
		logrus.Printf("[Execute Txn] SendTx success. cfg.Coinbase = %v, gasfee = %v", cfg.Coinbase, gasfee)
	}

	return nil
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
