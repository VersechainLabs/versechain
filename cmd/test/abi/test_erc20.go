// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abigen

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// TestErc20MetaData contains all meta data concerning the TestErc20 contract.
var TestErc20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"prevCounter\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newCounter\",\"type\":\"uint256\"}],\"name\":\"CounterUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EmptyEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"UserWriteEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"readCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sendEmptyEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sendEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newCounter\",\"type\":\"uint256\"}],\"name\":\"writeCounter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TestErc20ABI is the input ABI used to generate the binding from.
// Deprecated: Use TestErc20MetaData.ABI instead.
var TestErc20ABI = TestErc20MetaData.ABI

// TestErc20 is an auto generated Go binding around an Ethereum contract.
type TestErc20 struct {
	TestErc20Caller     // Read-only binding to the contract
	TestErc20Transactor // Write-only binding to the contract
	TestErc20Filterer   // Log filterer for contract events
}

// TestErc20Caller is an auto generated read-only Go binding around an Ethereum contract.
type TestErc20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestErc20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type TestErc20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestErc20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestErc20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestErc20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestErc20Session struct {
	Contract     *TestErc20        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestErc20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestErc20CallerSession struct {
	Contract *TestErc20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TestErc20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestErc20TransactorSession struct {
	Contract     *TestErc20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TestErc20Raw is an auto generated low-level Go binding around an Ethereum contract.
type TestErc20Raw struct {
	Contract *TestErc20 // Generic contract binding to access the raw methods on
}

// TestErc20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestErc20CallerRaw struct {
	Contract *TestErc20Caller // Generic read-only contract binding to access the raw methods on
}

// TestErc20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestErc20TransactorRaw struct {
	Contract *TestErc20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTestErc20 creates a new instance of TestErc20, bound to a specific deployed contract.
func NewTestErc20(address common.Address, backend bind.ContractBackend) (*TestErc20, error) {
	contract, err := bindTestErc20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestErc20{TestErc20Caller: TestErc20Caller{contract: contract}, TestErc20Transactor: TestErc20Transactor{contract: contract}, TestErc20Filterer: TestErc20Filterer{contract: contract}}, nil
}

// NewTestErc20Caller creates a new read-only instance of TestErc20, bound to a specific deployed contract.
func NewTestErc20Caller(address common.Address, caller bind.ContractCaller) (*TestErc20Caller, error) {
	contract, err := bindTestErc20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestErc20Caller{contract: contract}, nil
}

// NewTestErc20Transactor creates a new write-only instance of TestErc20, bound to a specific deployed contract.
func NewTestErc20Transactor(address common.Address, transactor bind.ContractTransactor) (*TestErc20Transactor, error) {
	contract, err := bindTestErc20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestErc20Transactor{contract: contract}, nil
}

// NewTestErc20Filterer creates a new log filterer instance of TestErc20, bound to a specific deployed contract.
func NewTestErc20Filterer(address common.Address, filterer bind.ContractFilterer) (*TestErc20Filterer, error) {
	contract, err := bindTestErc20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestErc20Filterer{contract: contract}, nil
}

// bindTestErc20 binds a generic wrapper to an already deployed contract.
func bindTestErc20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestErc20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestErc20 *TestErc20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestErc20.Contract.TestErc20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestErc20 *TestErc20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestErc20.Contract.TestErc20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestErc20 *TestErc20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestErc20.Contract.TestErc20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestErc20 *TestErc20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestErc20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestErc20 *TestErc20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestErc20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestErc20 *TestErc20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestErc20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestErc20 *TestErc20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestErc20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestErc20 *TestErc20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TestErc20.Contract.Allowance(&_TestErc20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestErc20 *TestErc20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TestErc20.Contract.Allowance(&_TestErc20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestErc20 *TestErc20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestErc20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestErc20 *TestErc20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _TestErc20.Contract.BalanceOf(&_TestErc20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestErc20 *TestErc20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TestErc20.Contract.BalanceOf(&_TestErc20.CallOpts, account)
}

// CurrentBlockNumber is a free data retrieval call binding the contract method 0x378ec23b.
//
// Solidity: function currentBlockNumber() view returns(uint256)
func (_TestErc20 *TestErc20Caller) CurrentBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestErc20.contract.Call(opts, &out, "currentBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentBlockNumber is a free data retrieval call binding the contract method 0x378ec23b.
//
// Solidity: function currentBlockNumber() view returns(uint256)
func (_TestErc20 *TestErc20Session) CurrentBlockNumber() (*big.Int, error) {
	return _TestErc20.Contract.CurrentBlockNumber(&_TestErc20.CallOpts)
}

// CurrentBlockNumber is a free data retrieval call binding the contract method 0x378ec23b.
//
// Solidity: function currentBlockNumber() view returns(uint256)
func (_TestErc20 *TestErc20CallerSession) CurrentBlockNumber() (*big.Int, error) {
	return _TestErc20.Contract.CurrentBlockNumber(&_TestErc20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestErc20 *TestErc20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TestErc20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestErc20 *TestErc20Session) Decimals() (uint8, error) {
	return _TestErc20.Contract.Decimals(&_TestErc20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestErc20 *TestErc20CallerSession) Decimals() (uint8, error) {
	return _TestErc20.Contract.Decimals(&_TestErc20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestErc20 *TestErc20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestErc20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestErc20 *TestErc20Session) Name() (string, error) {
	return _TestErc20.Contract.Name(&_TestErc20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestErc20 *TestErc20CallerSession) Name() (string, error) {
	return _TestErc20.Contract.Name(&_TestErc20.CallOpts)
}

// ReadCounter is a free data retrieval call binding the contract method 0xb8b085f2.
//
// Solidity: function readCounter() view returns(uint256)
func (_TestErc20 *TestErc20Caller) ReadCounter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestErc20.contract.Call(opts, &out, "readCounter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ReadCounter is a free data retrieval call binding the contract method 0xb8b085f2.
//
// Solidity: function readCounter() view returns(uint256)
func (_TestErc20 *TestErc20Session) ReadCounter() (*big.Int, error) {
	return _TestErc20.Contract.ReadCounter(&_TestErc20.CallOpts)
}

// ReadCounter is a free data retrieval call binding the contract method 0xb8b085f2.
//
// Solidity: function readCounter() view returns(uint256)
func (_TestErc20 *TestErc20CallerSession) ReadCounter() (*big.Int, error) {
	return _TestErc20.Contract.ReadCounter(&_TestErc20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestErc20 *TestErc20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestErc20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestErc20 *TestErc20Session) Symbol() (string, error) {
	return _TestErc20.Contract.Symbol(&_TestErc20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestErc20 *TestErc20CallerSession) Symbol() (string, error) {
	return _TestErc20.Contract.Symbol(&_TestErc20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestErc20 *TestErc20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestErc20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestErc20 *TestErc20Session) TotalSupply() (*big.Int, error) {
	return _TestErc20.Contract.TotalSupply(&_TestErc20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestErc20 *TestErc20CallerSession) TotalSupply() (*big.Int, error) {
	return _TestErc20.Contract.TotalSupply(&_TestErc20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TestErc20 *TestErc20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TestErc20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TestErc20 *TestErc20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TestErc20.Contract.Approve(&_TestErc20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_TestErc20 *TestErc20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _TestErc20.Contract.Approve(&_TestErc20.TransactOpts, spender, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_TestErc20 *TestErc20Transactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestErc20.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_TestErc20 *TestErc20Session) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestErc20.Contract.Mint(&_TestErc20.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_TestErc20 *TestErc20TransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestErc20.Contract.Mint(&_TestErc20.TransactOpts, to, amount)
}

// SendEmptyEvent is a paid mutator transaction binding the contract method 0x5e863fa7.
//
// Solidity: function sendEmptyEvent() returns()
func (_TestErc20 *TestErc20Transactor) SendEmptyEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestErc20.contract.Transact(opts, "sendEmptyEvent")
}

// SendEmptyEvent is a paid mutator transaction binding the contract method 0x5e863fa7.
//
// Solidity: function sendEmptyEvent() returns()
func (_TestErc20 *TestErc20Session) SendEmptyEvent() (*types.Transaction, error) {
	return _TestErc20.Contract.SendEmptyEvent(&_TestErc20.TransactOpts)
}

// SendEmptyEvent is a paid mutator transaction binding the contract method 0x5e863fa7.
//
// Solidity: function sendEmptyEvent() returns()
func (_TestErc20 *TestErc20TransactorSession) SendEmptyEvent() (*types.Transaction, error) {
	return _TestErc20.Contract.SendEmptyEvent(&_TestErc20.TransactOpts)
}

// SendEvent is a paid mutator transaction binding the contract method 0x32b7a761.
//
// Solidity: function sendEvent() returns()
func (_TestErc20 *TestErc20Transactor) SendEvent(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestErc20.contract.Transact(opts, "sendEvent")
}

// SendEvent is a paid mutator transaction binding the contract method 0x32b7a761.
//
// Solidity: function sendEvent() returns()
func (_TestErc20 *TestErc20Session) SendEvent() (*types.Transaction, error) {
	return _TestErc20.Contract.SendEvent(&_TestErc20.TransactOpts)
}

// SendEvent is a paid mutator transaction binding the contract method 0x32b7a761.
//
// Solidity: function sendEvent() returns()
func (_TestErc20 *TestErc20TransactorSession) SendEvent() (*types.Transaction, error) {
	return _TestErc20.Contract.SendEvent(&_TestErc20.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_TestErc20 *TestErc20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TestErc20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_TestErc20 *TestErc20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TestErc20.Contract.Transfer(&_TestErc20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_TestErc20 *TestErc20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TestErc20.Contract.Transfer(&_TestErc20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_TestErc20 *TestErc20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TestErc20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_TestErc20 *TestErc20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TestErc20.Contract.TransferFrom(&_TestErc20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_TestErc20 *TestErc20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _TestErc20.Contract.TransferFrom(&_TestErc20.TransactOpts, from, to, value)
}

// WriteCounter is a paid mutator transaction binding the contract method 0x98c41f69.
//
// Solidity: function writeCounter(uint256 newCounter) returns()
func (_TestErc20 *TestErc20Transactor) WriteCounter(opts *bind.TransactOpts, newCounter *big.Int) (*types.Transaction, error) {
	return _TestErc20.contract.Transact(opts, "writeCounter", newCounter)
}

// WriteCounter is a paid mutator transaction binding the contract method 0x98c41f69.
//
// Solidity: function writeCounter(uint256 newCounter) returns()
func (_TestErc20 *TestErc20Session) WriteCounter(newCounter *big.Int) (*types.Transaction, error) {
	return _TestErc20.Contract.WriteCounter(&_TestErc20.TransactOpts, newCounter)
}

// WriteCounter is a paid mutator transaction binding the contract method 0x98c41f69.
//
// Solidity: function writeCounter(uint256 newCounter) returns()
func (_TestErc20 *TestErc20TransactorSession) WriteCounter(newCounter *big.Int) (*types.Transaction, error) {
	return _TestErc20.Contract.WriteCounter(&_TestErc20.TransactOpts, newCounter)
}

// TestErc20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TestErc20 contract.
type TestErc20ApprovalIterator struct {
	Event *TestErc20Approval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TestErc20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestErc20Approval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TestErc20Approval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TestErc20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestErc20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestErc20Approval represents a Approval event raised by the TestErc20 contract.
type TestErc20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestErc20 *TestErc20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TestErc20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TestErc20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TestErc20ApprovalIterator{contract: _TestErc20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestErc20 *TestErc20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TestErc20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TestErc20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestErc20Approval)
				if err := _TestErc20.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestErc20 *TestErc20Filterer) ParseApproval(log types.Log) (*TestErc20Approval, error) {
	event := new(TestErc20Approval)
	if err := _TestErc20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestErc20CounterUpdatedIterator is returned from FilterCounterUpdated and is used to iterate over the raw logs and unpacked data for CounterUpdated events raised by the TestErc20 contract.
type TestErc20CounterUpdatedIterator struct {
	Event *TestErc20CounterUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TestErc20CounterUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestErc20CounterUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TestErc20CounterUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TestErc20CounterUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestErc20CounterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestErc20CounterUpdated represents a CounterUpdated event raised by the TestErc20 contract.
type TestErc20CounterUpdated struct {
	PrevCounter *big.Int
	NewCounter  *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterCounterUpdated is a free log retrieval operation binding the contract event 0x28bab7182b1f3ed84e1006088a7e959b8999f63af651809d61a39e97b026fcde.
//
// Solidity: event CounterUpdated(uint256 prevCounter, uint256 newCounter)
func (_TestErc20 *TestErc20Filterer) FilterCounterUpdated(opts *bind.FilterOpts) (*TestErc20CounterUpdatedIterator, error) {

	logs, sub, err := _TestErc20.contract.FilterLogs(opts, "CounterUpdated")
	if err != nil {
		return nil, err
	}
	return &TestErc20CounterUpdatedIterator{contract: _TestErc20.contract, event: "CounterUpdated", logs: logs, sub: sub}, nil
}

// WatchCounterUpdated is a free log subscription operation binding the contract event 0x28bab7182b1f3ed84e1006088a7e959b8999f63af651809d61a39e97b026fcde.
//
// Solidity: event CounterUpdated(uint256 prevCounter, uint256 newCounter)
func (_TestErc20 *TestErc20Filterer) WatchCounterUpdated(opts *bind.WatchOpts, sink chan<- *TestErc20CounterUpdated) (event.Subscription, error) {

	logs, sub, err := _TestErc20.contract.WatchLogs(opts, "CounterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestErc20CounterUpdated)
				if err := _TestErc20.contract.UnpackLog(event, "CounterUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCounterUpdated is a log parse operation binding the contract event 0x28bab7182b1f3ed84e1006088a7e959b8999f63af651809d61a39e97b026fcde.
//
// Solidity: event CounterUpdated(uint256 prevCounter, uint256 newCounter)
func (_TestErc20 *TestErc20Filterer) ParseCounterUpdated(log types.Log) (*TestErc20CounterUpdated, error) {
	event := new(TestErc20CounterUpdated)
	if err := _TestErc20.contract.UnpackLog(event, "CounterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestErc20EmptyEventIterator is returned from FilterEmptyEvent and is used to iterate over the raw logs and unpacked data for EmptyEvent events raised by the TestErc20 contract.
type TestErc20EmptyEventIterator struct {
	Event *TestErc20EmptyEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TestErc20EmptyEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestErc20EmptyEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TestErc20EmptyEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TestErc20EmptyEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestErc20EmptyEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestErc20EmptyEvent represents a EmptyEvent event raised by the TestErc20 contract.
type TestErc20EmptyEvent struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEmptyEvent is a free log retrieval operation binding the contract event 0xcf16a92280c1bbb43f72d31126b724d508df2877835849e8744017ab36a9b47f.
//
// Solidity: event EmptyEvent()
func (_TestErc20 *TestErc20Filterer) FilterEmptyEvent(opts *bind.FilterOpts) (*TestErc20EmptyEventIterator, error) {

	logs, sub, err := _TestErc20.contract.FilterLogs(opts, "EmptyEvent")
	if err != nil {
		return nil, err
	}
	return &TestErc20EmptyEventIterator{contract: _TestErc20.contract, event: "EmptyEvent", logs: logs, sub: sub}, nil
}

// WatchEmptyEvent is a free log subscription operation binding the contract event 0xcf16a92280c1bbb43f72d31126b724d508df2877835849e8744017ab36a9b47f.
//
// Solidity: event EmptyEvent()
func (_TestErc20 *TestErc20Filterer) WatchEmptyEvent(opts *bind.WatchOpts, sink chan<- *TestErc20EmptyEvent) (event.Subscription, error) {

	logs, sub, err := _TestErc20.contract.WatchLogs(opts, "EmptyEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestErc20EmptyEvent)
				if err := _TestErc20.contract.UnpackLog(event, "EmptyEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEmptyEvent is a log parse operation binding the contract event 0xcf16a92280c1bbb43f72d31126b724d508df2877835849e8744017ab36a9b47f.
//
// Solidity: event EmptyEvent()
func (_TestErc20 *TestErc20Filterer) ParseEmptyEvent(log types.Log) (*TestErc20EmptyEvent, error) {
	event := new(TestErc20EmptyEvent)
	if err := _TestErc20.contract.UnpackLog(event, "EmptyEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestErc20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TestErc20 contract.
type TestErc20TransferIterator struct {
	Event *TestErc20Transfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TestErc20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestErc20Transfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TestErc20Transfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TestErc20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestErc20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestErc20Transfer represents a Transfer event raised by the TestErc20 contract.
type TestErc20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestErc20 *TestErc20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TestErc20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestErc20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TestErc20TransferIterator{contract: _TestErc20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestErc20 *TestErc20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TestErc20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestErc20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestErc20Transfer)
				if err := _TestErc20.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestErc20 *TestErc20Filterer) ParseTransfer(log types.Log) (*TestErc20Transfer, error) {
	event := new(TestErc20Transfer)
	if err := _TestErc20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestErc20UserWriteEventIterator is returned from FilterUserWriteEvent and is used to iterate over the raw logs and unpacked data for UserWriteEvent events raised by the TestErc20 contract.
type TestErc20UserWriteEventIterator struct {
	Event *TestErc20UserWriteEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TestErc20UserWriteEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestErc20UserWriteEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TestErc20UserWriteEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TestErc20UserWriteEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestErc20UserWriteEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestErc20UserWriteEvent represents a UserWriteEvent event raised by the TestErc20 contract.
type TestErc20UserWriteEvent struct {
	User        common.Address
	BlockNumber *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUserWriteEvent is a free log retrieval operation binding the contract event 0xe7a5a2788502413da105ba7eebbd320341bcf15413708e8261856f535531eabd.
//
// Solidity: event UserWriteEvent(address indexed user, uint256 blockNumber)
func (_TestErc20 *TestErc20Filterer) FilterUserWriteEvent(opts *bind.FilterOpts, user []common.Address) (*TestErc20UserWriteEventIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TestErc20.contract.FilterLogs(opts, "UserWriteEvent", userRule)
	if err != nil {
		return nil, err
	}
	return &TestErc20UserWriteEventIterator{contract: _TestErc20.contract, event: "UserWriteEvent", logs: logs, sub: sub}, nil
}

// WatchUserWriteEvent is a free log subscription operation binding the contract event 0xe7a5a2788502413da105ba7eebbd320341bcf15413708e8261856f535531eabd.
//
// Solidity: event UserWriteEvent(address indexed user, uint256 blockNumber)
func (_TestErc20 *TestErc20Filterer) WatchUserWriteEvent(opts *bind.WatchOpts, sink chan<- *TestErc20UserWriteEvent, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TestErc20.contract.WatchLogs(opts, "UserWriteEvent", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestErc20UserWriteEvent)
				if err := _TestErc20.contract.UnpackLog(event, "UserWriteEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUserWriteEvent is a log parse operation binding the contract event 0xe7a5a2788502413da105ba7eebbd320341bcf15413708e8261856f535531eabd.
//
// Solidity: event UserWriteEvent(address indexed user, uint256 blockNumber)
func (_TestErc20 *TestErc20Filterer) ParseUserWriteEvent(log types.Log) (*TestErc20UserWriteEvent, error) {
	event := new(TestErc20UserWriteEvent)
	if err := _TestErc20.contract.UnpackLog(event, "UserWriteEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
