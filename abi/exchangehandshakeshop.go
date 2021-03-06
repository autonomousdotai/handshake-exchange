// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// ExchangeHandshakeShopABI is the input ABI used to generate the binding from.
const ExchangeHandshakeShopABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ex\",\"outputs\":[{\"name\":\"stationOwner\",\"type\":\"address\"},{\"name\":\"customer\",\"type\":\"address\"},{\"name\":\"escrow\",\"type\":\"uint256\"},{\"name\":\"state\",\"type\":\"uint8\"},{\"name\":\"exType\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"customer\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"offchainP\",\"type\":\"bytes32\"},{\"name\":\"offchainC\",\"type\":\"bytes32\"}],\"name\":\"releasePartialFund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"}],\"name\":\"getBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"finish\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"}],\"name\":\"getState\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"resetAllStation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"f\",\"type\":\"uint256\"}],\"name\":\"setFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"reject\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"addInventory\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"stationOwner\",\"type\":\"address\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"initByCustomer\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"initByStationOwner\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"cancel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"closeByStationOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"__setFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"stationOwner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__initByStationOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__closeByStationOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"customer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchainP\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"offchainC\",\"type\":\"bytes32\"}],\"name\":\"__releasePartialFund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__addInventory\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"customer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"stationOwner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__initByCustomer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__cancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__reject\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__finish\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"}],\"name\":\"__resetStation\",\"type\":\"event\"}]"

// ExchangeHandshakeShop is an auto generated Go binding around an Ethereum contract.
type ExchangeHandshakeShop struct {
	ExchangeHandshakeShopCaller     // Read-only binding to the contract
	ExchangeHandshakeShopTransactor // Write-only binding to the contract
	ExchangeHandshakeShopFilterer   // Log filterer for contract events
}

// ExchangeHandshakeShopCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExchangeHandshakeShopCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeHandshakeShopTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExchangeHandshakeShopTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeHandshakeShopFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExchangeHandshakeShopFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeHandshakeShopSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExchangeHandshakeShopSession struct {
	Contract     *ExchangeHandshakeShop // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ExchangeHandshakeShopCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExchangeHandshakeShopCallerSession struct {
	Contract *ExchangeHandshakeShopCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// ExchangeHandshakeShopTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExchangeHandshakeShopTransactorSession struct {
	Contract     *ExchangeHandshakeShopTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// ExchangeHandshakeShopRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExchangeHandshakeShopRaw struct {
	Contract *ExchangeHandshakeShop // Generic contract binding to access the raw methods on
}

// ExchangeHandshakeShopCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExchangeHandshakeShopCallerRaw struct {
	Contract *ExchangeHandshakeShopCaller // Generic read-only contract binding to access the raw methods on
}

// ExchangeHandshakeShopTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExchangeHandshakeShopTransactorRaw struct {
	Contract *ExchangeHandshakeShopTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExchangeHandshakeShop creates a new instance of ExchangeHandshakeShop, bound to a specific deployed contract.
func NewExchangeHandshakeShop(address common.Address, backend bind.ContractBackend) (*ExchangeHandshakeShop, error) {
	contract, err := bindExchangeHandshakeShop(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShop{ExchangeHandshakeShopCaller: ExchangeHandshakeShopCaller{contract: contract}, ExchangeHandshakeShopTransactor: ExchangeHandshakeShopTransactor{contract: contract}, ExchangeHandshakeShopFilterer: ExchangeHandshakeShopFilterer{contract: contract}}, nil
}

// NewExchangeHandshakeShopCaller creates a new read-only instance of ExchangeHandshakeShop, bound to a specific deployed contract.
func NewExchangeHandshakeShopCaller(address common.Address, caller bind.ContractCaller) (*ExchangeHandshakeShopCaller, error) {
	contract, err := bindExchangeHandshakeShop(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShopCaller{contract: contract}, nil
}

// NewExchangeHandshakeShopTransactor creates a new write-only instance of ExchangeHandshakeShop, bound to a specific deployed contract.
func NewExchangeHandshakeShopTransactor(address common.Address, transactor bind.ContractTransactor) (*ExchangeHandshakeShopTransactor, error) {
	contract, err := bindExchangeHandshakeShop(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShopTransactor{contract: contract}, nil
}

// NewExchangeHandshakeShopFilterer creates a new log filterer instance of ExchangeHandshakeShop, bound to a specific deployed contract.
func NewExchangeHandshakeShopFilterer(address common.Address, filterer bind.ContractFilterer) (*ExchangeHandshakeShopFilterer, error) {
	contract, err := bindExchangeHandshakeShop(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShopFilterer{contract: contract}, nil
}

// bindExchangeHandshakeShop binds a generic wrapper to an already deployed contract.
func bindExchangeHandshakeShop(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangeHandshakeShopABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExchangeHandshakeShop *ExchangeHandshakeShopRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ExchangeHandshakeShop.Contract.ExchangeHandshakeShopCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExchangeHandshakeShop *ExchangeHandshakeShopRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.ExchangeHandshakeShopTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExchangeHandshakeShop *ExchangeHandshakeShopRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.ExchangeHandshakeShopTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExchangeHandshakeShop *ExchangeHandshakeShopCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ExchangeHandshakeShop.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.contract.Transact(opts, method, params...)
}

// Ex is a free data retrieval call binding the contract method 0x1089f215.
//
// Solidity: function ex( uint256) constant returns(stationOwner address, customer address, escrow uint256, state uint8, exType uint8)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopCaller) Ex(opts *bind.CallOpts, arg0 *big.Int) (struct {
	StationOwner common.Address
	Customer     common.Address
	Escrow       *big.Int
	State        uint8
	ExType       uint8
}, error) {
	ret := new(struct {
		StationOwner common.Address
		Customer     common.Address
		Escrow       *big.Int
		State        uint8
		ExType       uint8
	})
	out := ret
	err := _ExchangeHandshakeShop.contract.Call(opts, out, "ex", arg0)
	return *ret, err
}

// Ex is a free data retrieval call binding the contract method 0x1089f215.
//
// Solidity: function ex( uint256) constant returns(stationOwner address, customer address, escrow uint256, state uint8, exType uint8)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopSession) Ex(arg0 *big.Int) (struct {
	StationOwner common.Address
	Customer     common.Address
	Escrow       *big.Int
	State        uint8
	ExType       uint8
}, error) {
	return _ExchangeHandshakeShop.Contract.Ex(&_ExchangeHandshakeShop.CallOpts, arg0)
}

// Ex is a free data retrieval call binding the contract method 0x1089f215.
//
// Solidity: function ex( uint256) constant returns(stationOwner address, customer address, escrow uint256, state uint8, exType uint8)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopCallerSession) Ex(arg0 *big.Int) (struct {
	StationOwner common.Address
	Customer     common.Address
	Escrow       *big.Int
	State        uint8
	ExType       uint8
}, error) {
	return _ExchangeHandshakeShop.Contract.Ex(&_ExchangeHandshakeShop.CallOpts, arg0)
}

// GetBalance is a free data retrieval call binding the contract method 0x1e010439.
//
// Solidity: function getBalance(hid uint256) constant returns(uint256)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopCaller) GetBalance(opts *bind.CallOpts, hid *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExchangeHandshakeShop.contract.Call(opts, out, "getBalance", hid)
	return *ret0, err
}

// GetBalance is a free data retrieval call binding the contract method 0x1e010439.
//
// Solidity: function getBalance(hid uint256) constant returns(uint256)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopSession) GetBalance(hid *big.Int) (*big.Int, error) {
	return _ExchangeHandshakeShop.Contract.GetBalance(&_ExchangeHandshakeShop.CallOpts, hid)
}

// GetBalance is a free data retrieval call binding the contract method 0x1e010439.
//
// Solidity: function getBalance(hid uint256) constant returns(uint256)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopCallerSession) GetBalance(hid *big.Int) (*big.Int, error) {
	return _ExchangeHandshakeShop.Contract.GetBalance(&_ExchangeHandshakeShop.CallOpts, hid)
}

// GetState is a free data retrieval call binding the contract method 0x44c9af28.
//
// Solidity: function getState(hid uint256) constant returns(uint8)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopCaller) GetState(opts *bind.CallOpts, hid *big.Int) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _ExchangeHandshakeShop.contract.Call(opts, out, "getState", hid)
	return *ret0, err
}

// GetState is a free data retrieval call binding the contract method 0x44c9af28.
//
// Solidity: function getState(hid uint256) constant returns(uint8)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopSession) GetState(hid *big.Int) (uint8, error) {
	return _ExchangeHandshakeShop.Contract.GetState(&_ExchangeHandshakeShop.CallOpts, hid)
}

// GetState is a free data retrieval call binding the contract method 0x44c9af28.
//
// Solidity: function getState(hid uint256) constant returns(uint8)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopCallerSession) GetState(hid *big.Int) (uint8, error) {
	return _ExchangeHandshakeShop.Contract.GetState(&_ExchangeHandshakeShop.CallOpts, hid)
}

// AddInventory is a paid mutator transaction binding the contract method 0x93341a0b.
//
// Solidity: function addInventory(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactor) AddInventory(opts *bind.TransactOpts, hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.contract.Transact(opts, "addInventory", hid, offchain)
}

// AddInventory is a paid mutator transaction binding the contract method 0x93341a0b.
//
// Solidity: function addInventory(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopSession) AddInventory(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.AddInventory(&_ExchangeHandshakeShop.TransactOpts, hid, offchain)
}

// AddInventory is a paid mutator transaction binding the contract method 0x93341a0b.
//
// Solidity: function addInventory(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactorSession) AddInventory(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.AddInventory(&_ExchangeHandshakeShop.TransactOpts, hid, offchain)
}

// Cancel is a paid mutator transaction binding the contract method 0xeafb64d5.
//
// Solidity: function cancel(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactor) Cancel(opts *bind.TransactOpts, hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.contract.Transact(opts, "cancel", hid, offchain)
}

// Cancel is a paid mutator transaction binding the contract method 0xeafb64d5.
//
// Solidity: function cancel(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopSession) Cancel(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.Cancel(&_ExchangeHandshakeShop.TransactOpts, hid, offchain)
}

// Cancel is a paid mutator transaction binding the contract method 0xeafb64d5.
//
// Solidity: function cancel(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactorSession) Cancel(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.Cancel(&_ExchangeHandshakeShop.TransactOpts, hid, offchain)
}

// CloseByStationOwner is a paid mutator transaction binding the contract method 0xef3ae612.
//
// Solidity: function closeByStationOwner(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactor) CloseByStationOwner(opts *bind.TransactOpts, hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.contract.Transact(opts, "closeByStationOwner", hid, offchain)
}

// CloseByStationOwner is a paid mutator transaction binding the contract method 0xef3ae612.
//
// Solidity: function closeByStationOwner(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopSession) CloseByStationOwner(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.CloseByStationOwner(&_ExchangeHandshakeShop.TransactOpts, hid, offchain)
}

// CloseByStationOwner is a paid mutator transaction binding the contract method 0xef3ae612.
//
// Solidity: function closeByStationOwner(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactorSession) CloseByStationOwner(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.CloseByStationOwner(&_ExchangeHandshakeShop.TransactOpts, hid, offchain)
}

// Finish is a paid mutator transaction binding the contract method 0x35e8b2df.
//
// Solidity: function finish(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactor) Finish(opts *bind.TransactOpts, hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.contract.Transact(opts, "finish", hid, offchain)
}

// Finish is a paid mutator transaction binding the contract method 0x35e8b2df.
//
// Solidity: function finish(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopSession) Finish(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.Finish(&_ExchangeHandshakeShop.TransactOpts, hid, offchain)
}

// Finish is a paid mutator transaction binding the contract method 0x35e8b2df.
//
// Solidity: function finish(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactorSession) Finish(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.Finish(&_ExchangeHandshakeShop.TransactOpts, hid, offchain)
}

// InitByCustomer is a paid mutator transaction binding the contract method 0xa3b2ea71.
//
// Solidity: function initByCustomer(stationOwner address, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactor) InitByCustomer(opts *bind.TransactOpts, stationOwner common.Address, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.contract.Transact(opts, "initByCustomer", stationOwner, offchain)
}

// InitByCustomer is a paid mutator transaction binding the contract method 0xa3b2ea71.
//
// Solidity: function initByCustomer(stationOwner address, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopSession) InitByCustomer(stationOwner common.Address, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.InitByCustomer(&_ExchangeHandshakeShop.TransactOpts, stationOwner, offchain)
}

// InitByCustomer is a paid mutator transaction binding the contract method 0xa3b2ea71.
//
// Solidity: function initByCustomer(stationOwner address, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactorSession) InitByCustomer(stationOwner common.Address, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.InitByCustomer(&_ExchangeHandshakeShop.TransactOpts, stationOwner, offchain)
}

// InitByStationOwner is a paid mutator transaction binding the contract method 0xb4bd5dbb.
//
// Solidity: function initByStationOwner(offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactor) InitByStationOwner(opts *bind.TransactOpts, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.contract.Transact(opts, "initByStationOwner", offchain)
}

// InitByStationOwner is a paid mutator transaction binding the contract method 0xb4bd5dbb.
//
// Solidity: function initByStationOwner(offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopSession) InitByStationOwner(offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.InitByStationOwner(&_ExchangeHandshakeShop.TransactOpts, offchain)
}

// InitByStationOwner is a paid mutator transaction binding the contract method 0xb4bd5dbb.
//
// Solidity: function initByStationOwner(offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactorSession) InitByStationOwner(offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.InitByStationOwner(&_ExchangeHandshakeShop.TransactOpts, offchain)
}

// Reject is a paid mutator transaction binding the contract method 0x6be1320b.
//
// Solidity: function reject(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactor) Reject(opts *bind.TransactOpts, hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.contract.Transact(opts, "reject", hid, offchain)
}

// Reject is a paid mutator transaction binding the contract method 0x6be1320b.
//
// Solidity: function reject(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopSession) Reject(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.Reject(&_ExchangeHandshakeShop.TransactOpts, hid, offchain)
}

// Reject is a paid mutator transaction binding the contract method 0x6be1320b.
//
// Solidity: function reject(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactorSession) Reject(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.Reject(&_ExchangeHandshakeShop.TransactOpts, hid, offchain)
}

// ReleasePartialFund is a paid mutator transaction binding the contract method 0x15d85cee.
//
// Solidity: function releasePartialFund(hid uint256, customer address, amount uint256, offchainP bytes32, offchainC bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactor) ReleasePartialFund(opts *bind.TransactOpts, hid *big.Int, customer common.Address, amount *big.Int, offchainP [32]byte, offchainC [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.contract.Transact(opts, "releasePartialFund", hid, customer, amount, offchainP, offchainC)
}

// ReleasePartialFund is a paid mutator transaction binding the contract method 0x15d85cee.
//
// Solidity: function releasePartialFund(hid uint256, customer address, amount uint256, offchainP bytes32, offchainC bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopSession) ReleasePartialFund(hid *big.Int, customer common.Address, amount *big.Int, offchainP [32]byte, offchainC [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.ReleasePartialFund(&_ExchangeHandshakeShop.TransactOpts, hid, customer, amount, offchainP, offchainC)
}

// ReleasePartialFund is a paid mutator transaction binding the contract method 0x15d85cee.
//
// Solidity: function releasePartialFund(hid uint256, customer address, amount uint256, offchainP bytes32, offchainC bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactorSession) ReleasePartialFund(hid *big.Int, customer common.Address, amount *big.Int, offchainP [32]byte, offchainC [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.ReleasePartialFund(&_ExchangeHandshakeShop.TransactOpts, hid, customer, amount, offchainP, offchainC)
}

// ResetAllStation is a paid mutator transaction binding the contract method 0x5dccb051.
//
// Solidity: function resetAllStation(offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactor) ResetAllStation(opts *bind.TransactOpts, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.contract.Transact(opts, "resetAllStation", offchain)
}

// ResetAllStation is a paid mutator transaction binding the contract method 0x5dccb051.
//
// Solidity: function resetAllStation(offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopSession) ResetAllStation(offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.ResetAllStation(&_ExchangeHandshakeShop.TransactOpts, offchain)
}

// ResetAllStation is a paid mutator transaction binding the contract method 0x5dccb051.
//
// Solidity: function resetAllStation(offchain bytes32) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactorSession) ResetAllStation(offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.ResetAllStation(&_ExchangeHandshakeShop.TransactOpts, offchain)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(f uint256) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactor) SetFee(opts *bind.TransactOpts, f *big.Int) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.contract.Transact(opts, "setFee", f)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(f uint256) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopSession) SetFee(f *big.Int) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.SetFee(&_ExchangeHandshakeShop.TransactOpts, f)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(f uint256) returns()
func (_ExchangeHandshakeShop *ExchangeHandshakeShopTransactorSession) SetFee(f *big.Int) (*types.Transaction, error) {
	return _ExchangeHandshakeShop.Contract.SetFee(&_ExchangeHandshakeShop.TransactOpts, f)
}

// ExchangeHandshakeShopAddInventoryIterator is returned from FilterAddInventory and is used to iterate over the raw logs and unpacked data for AddInventory events raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopAddInventoryIterator struct {
	Event *ExchangeHandshakeShopAddInventory // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeShopAddInventoryIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeShopAddInventory)
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
		it.Event = new(ExchangeHandshakeShopAddInventory)
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
func (it *ExchangeHandshakeShopAddInventoryIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeShopAddInventoryIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeShopAddInventory represents a AddInventory event raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopAddInventory struct {
	Hid      *big.Int
	Offchain [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAddInventory is a free log retrieval operation binding the contract event 0x414b45f3f079051d7402c185b2031c6033fd30f7ff54ea87dbc22ec69c626452.
//
// Solidity: e __addInventory(hid uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) FilterAddInventory(opts *bind.FilterOpts) (*ExchangeHandshakeShopAddInventoryIterator, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.FilterLogs(opts, "__addInventory")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShopAddInventoryIterator{contract: _ExchangeHandshakeShop.contract, event: "__addInventory", logs: logs, sub: sub}, nil
}

// WatchAddInventory is a free log subscription operation binding the contract event 0x414b45f3f079051d7402c185b2031c6033fd30f7ff54ea87dbc22ec69c626452.
//
// Solidity: e __addInventory(hid uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) WatchAddInventory(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeShopAddInventory) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.WatchLogs(opts, "__addInventory")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeShopAddInventory)
				if err := _ExchangeHandshakeShop.contract.UnpackLog(event, "__addInventory", log); err != nil {
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

// ExchangeHandshakeShopCancelIterator is returned from FilterCancel and is used to iterate over the raw logs and unpacked data for Cancel events raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopCancelIterator struct {
	Event *ExchangeHandshakeShopCancel // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeShopCancelIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeShopCancel)
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
		it.Event = new(ExchangeHandshakeShopCancel)
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
func (it *ExchangeHandshakeShopCancelIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeShopCancelIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeShopCancel represents a Cancel event raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopCancel struct {
	Hid      *big.Int
	Offchain [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCancel is a free log retrieval operation binding the contract event 0xcb720f90f098c425fad2e9df556017c076f9bf7aefa096d1c904a06027ae0460.
//
// Solidity: e __cancel(hid uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) FilterCancel(opts *bind.FilterOpts) (*ExchangeHandshakeShopCancelIterator, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.FilterLogs(opts, "__cancel")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShopCancelIterator{contract: _ExchangeHandshakeShop.contract, event: "__cancel", logs: logs, sub: sub}, nil
}

// WatchCancel is a free log subscription operation binding the contract event 0xcb720f90f098c425fad2e9df556017c076f9bf7aefa096d1c904a06027ae0460.
//
// Solidity: e __cancel(hid uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) WatchCancel(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeShopCancel) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.WatchLogs(opts, "__cancel")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeShopCancel)
				if err := _ExchangeHandshakeShop.contract.UnpackLog(event, "__cancel", log); err != nil {
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

// ExchangeHandshakeShopCloseByStationOwnerIterator is returned from FilterCloseByStationOwner and is used to iterate over the raw logs and unpacked data for CloseByStationOwner events raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopCloseByStationOwnerIterator struct {
	Event *ExchangeHandshakeShopCloseByStationOwner // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeShopCloseByStationOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeShopCloseByStationOwner)
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
		it.Event = new(ExchangeHandshakeShopCloseByStationOwner)
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
func (it *ExchangeHandshakeShopCloseByStationOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeShopCloseByStationOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeShopCloseByStationOwner represents a CloseByStationOwner event raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopCloseByStationOwner struct {
	Hid      *big.Int
	Offchain [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCloseByStationOwner is a free log retrieval operation binding the contract event 0xa4b2dfc35c48bbe07f2063bd9fcf0afa08891ebf8b417a983f957ac8708612a3.
//
// Solidity: e __closeByStationOwner(hid uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) FilterCloseByStationOwner(opts *bind.FilterOpts) (*ExchangeHandshakeShopCloseByStationOwnerIterator, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.FilterLogs(opts, "__closeByStationOwner")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShopCloseByStationOwnerIterator{contract: _ExchangeHandshakeShop.contract, event: "__closeByStationOwner", logs: logs, sub: sub}, nil
}

// WatchCloseByStationOwner is a free log subscription operation binding the contract event 0xa4b2dfc35c48bbe07f2063bd9fcf0afa08891ebf8b417a983f957ac8708612a3.
//
// Solidity: e __closeByStationOwner(hid uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) WatchCloseByStationOwner(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeShopCloseByStationOwner) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.WatchLogs(opts, "__closeByStationOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeShopCloseByStationOwner)
				if err := _ExchangeHandshakeShop.contract.UnpackLog(event, "__closeByStationOwner", log); err != nil {
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

// ExchangeHandshakeShopFinishIterator is returned from FilterFinish and is used to iterate over the raw logs and unpacked data for Finish events raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopFinishIterator struct {
	Event *ExchangeHandshakeShopFinish // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeShopFinishIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeShopFinish)
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
		it.Event = new(ExchangeHandshakeShopFinish)
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
func (it *ExchangeHandshakeShopFinishIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeShopFinishIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeShopFinish represents a Finish event raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopFinish struct {
	Hid      *big.Int
	Offchain [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterFinish is a free log retrieval operation binding the contract event 0xfda1b0d3f21a187df4a198a15b2361d3cc73501a41bf582d8bcadc9d266da83a.
//
// Solidity: e __finish(hid uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) FilterFinish(opts *bind.FilterOpts) (*ExchangeHandshakeShopFinishIterator, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.FilterLogs(opts, "__finish")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShopFinishIterator{contract: _ExchangeHandshakeShop.contract, event: "__finish", logs: logs, sub: sub}, nil
}

// WatchFinish is a free log subscription operation binding the contract event 0xfda1b0d3f21a187df4a198a15b2361d3cc73501a41bf582d8bcadc9d266da83a.
//
// Solidity: e __finish(hid uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) WatchFinish(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeShopFinish) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.WatchLogs(opts, "__finish")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeShopFinish)
				if err := _ExchangeHandshakeShop.contract.UnpackLog(event, "__finish", log); err != nil {
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

// ExchangeHandshakeShopInitByCustomerIterator is returned from FilterInitByCustomer and is used to iterate over the raw logs and unpacked data for InitByCustomer events raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopInitByCustomerIterator struct {
	Event *ExchangeHandshakeShopInitByCustomer // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeShopInitByCustomerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeShopInitByCustomer)
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
		it.Event = new(ExchangeHandshakeShopInitByCustomer)
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
func (it *ExchangeHandshakeShopInitByCustomerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeShopInitByCustomerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeShopInitByCustomer represents a InitByCustomer event raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopInitByCustomer struct {
	Hid          *big.Int
	Customer     common.Address
	StationOwner common.Address
	Value        *big.Int
	Offchain     [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterInitByCustomer is a free log retrieval operation binding the contract event 0xaa539cf6b39fbef365d4e0d2b9d8017bd6b1a36893a02b2a7f0324a705e83094.
//
// Solidity: e __initByCustomer(hid uint256, customer address, stationOwner address, value uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) FilterInitByCustomer(opts *bind.FilterOpts) (*ExchangeHandshakeShopInitByCustomerIterator, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.FilterLogs(opts, "__initByCustomer")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShopInitByCustomerIterator{contract: _ExchangeHandshakeShop.contract, event: "__initByCustomer", logs: logs, sub: sub}, nil
}

// WatchInitByCustomer is a free log subscription operation binding the contract event 0xaa539cf6b39fbef365d4e0d2b9d8017bd6b1a36893a02b2a7f0324a705e83094.
//
// Solidity: e __initByCustomer(hid uint256, customer address, stationOwner address, value uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) WatchInitByCustomer(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeShopInitByCustomer) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.WatchLogs(opts, "__initByCustomer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeShopInitByCustomer)
				if err := _ExchangeHandshakeShop.contract.UnpackLog(event, "__initByCustomer", log); err != nil {
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

// ExchangeHandshakeShopInitByStationOwnerIterator is returned from FilterInitByStationOwner and is used to iterate over the raw logs and unpacked data for InitByStationOwner events raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopInitByStationOwnerIterator struct {
	Event *ExchangeHandshakeShopInitByStationOwner // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeShopInitByStationOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeShopInitByStationOwner)
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
		it.Event = new(ExchangeHandshakeShopInitByStationOwner)
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
func (it *ExchangeHandshakeShopInitByStationOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeShopInitByStationOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeShopInitByStationOwner represents a InitByStationOwner event raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopInitByStationOwner struct {
	Hid          *big.Int
	StationOwner common.Address
	Value        *big.Int
	Offchain     [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterInitByStationOwner is a free log retrieval operation binding the contract event 0x04d411f10f50ce6fe597ddd1436d44c2118417e90cbb982e009b8215ee689c79.
//
// Solidity: e __initByStationOwner(hid uint256, stationOwner address, value uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) FilterInitByStationOwner(opts *bind.FilterOpts) (*ExchangeHandshakeShopInitByStationOwnerIterator, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.FilterLogs(opts, "__initByStationOwner")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShopInitByStationOwnerIterator{contract: _ExchangeHandshakeShop.contract, event: "__initByStationOwner", logs: logs, sub: sub}, nil
}

// WatchInitByStationOwner is a free log subscription operation binding the contract event 0x04d411f10f50ce6fe597ddd1436d44c2118417e90cbb982e009b8215ee689c79.
//
// Solidity: e __initByStationOwner(hid uint256, stationOwner address, value uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) WatchInitByStationOwner(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeShopInitByStationOwner) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.WatchLogs(opts, "__initByStationOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeShopInitByStationOwner)
				if err := _ExchangeHandshakeShop.contract.UnpackLog(event, "__initByStationOwner", log); err != nil {
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

// ExchangeHandshakeShopRejectIterator is returned from FilterReject and is used to iterate over the raw logs and unpacked data for Reject events raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopRejectIterator struct {
	Event *ExchangeHandshakeShopReject // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeShopRejectIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeShopReject)
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
		it.Event = new(ExchangeHandshakeShopReject)
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
func (it *ExchangeHandshakeShopRejectIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeShopRejectIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeShopReject represents a Reject event raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopReject struct {
	Hid      *big.Int
	Offchain [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReject is a free log retrieval operation binding the contract event 0xae76720f3a5d319b91bc94d8a6c2e3096a4f3554c8cb897e3aedfced5824a10a.
//
// Solidity: e __reject(hid uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) FilterReject(opts *bind.FilterOpts) (*ExchangeHandshakeShopRejectIterator, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.FilterLogs(opts, "__reject")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShopRejectIterator{contract: _ExchangeHandshakeShop.contract, event: "__reject", logs: logs, sub: sub}, nil
}

// WatchReject is a free log subscription operation binding the contract event 0xae76720f3a5d319b91bc94d8a6c2e3096a4f3554c8cb897e3aedfced5824a10a.
//
// Solidity: e __reject(hid uint256, offchain bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) WatchReject(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeShopReject) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.WatchLogs(opts, "__reject")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeShopReject)
				if err := _ExchangeHandshakeShop.contract.UnpackLog(event, "__reject", log); err != nil {
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

// ExchangeHandshakeShopReleasePartialFundIterator is returned from FilterReleasePartialFund and is used to iterate over the raw logs and unpacked data for ReleasePartialFund events raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopReleasePartialFundIterator struct {
	Event *ExchangeHandshakeShopReleasePartialFund // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeShopReleasePartialFundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeShopReleasePartialFund)
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
		it.Event = new(ExchangeHandshakeShopReleasePartialFund)
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
func (it *ExchangeHandshakeShopReleasePartialFundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeShopReleasePartialFundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeShopReleasePartialFund represents a ReleasePartialFund event raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopReleasePartialFund struct {
	Hid       *big.Int
	Customer  common.Address
	Amount    *big.Int
	OffchainP [32]byte
	OffchainC [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterReleasePartialFund is a free log retrieval operation binding the contract event 0xf83285d8d874832066930a6cd961e4fa4156c7a28dcdfb4bc9aaa13e0be2bfcc.
//
// Solidity: e __releasePartialFund(hid uint256, customer address, amount uint256, offchainP bytes32, offchainC bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) FilterReleasePartialFund(opts *bind.FilterOpts) (*ExchangeHandshakeShopReleasePartialFundIterator, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.FilterLogs(opts, "__releasePartialFund")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShopReleasePartialFundIterator{contract: _ExchangeHandshakeShop.contract, event: "__releasePartialFund", logs: logs, sub: sub}, nil
}

// WatchReleasePartialFund is a free log subscription operation binding the contract event 0xf83285d8d874832066930a6cd961e4fa4156c7a28dcdfb4bc9aaa13e0be2bfcc.
//
// Solidity: e __releasePartialFund(hid uint256, customer address, amount uint256, offchainP bytes32, offchainC bytes32)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) WatchReleasePartialFund(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeShopReleasePartialFund) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.WatchLogs(opts, "__releasePartialFund")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeShopReleasePartialFund)
				if err := _ExchangeHandshakeShop.contract.UnpackLog(event, "__releasePartialFund", log); err != nil {
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

// ExchangeHandshakeShopResetStationIterator is returned from FilterResetStation and is used to iterate over the raw logs and unpacked data for ResetStation events raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopResetStationIterator struct {
	Event *ExchangeHandshakeShopResetStation // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeShopResetStationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeShopResetStation)
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
		it.Event = new(ExchangeHandshakeShopResetStation)
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
func (it *ExchangeHandshakeShopResetStationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeShopResetStationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeShopResetStation represents a ResetStation event raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopResetStation struct {
	Offchain [32]byte
	Hid      *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterResetStation is a free log retrieval operation binding the contract event 0x73d8ef4b77a16350f9072178ca21a79e0d68298007b2a92c870e2b2c66f7868a.
//
// Solidity: e __resetStation(offchain bytes32, hid uint256)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) FilterResetStation(opts *bind.FilterOpts) (*ExchangeHandshakeShopResetStationIterator, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.FilterLogs(opts, "__resetStation")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShopResetStationIterator{contract: _ExchangeHandshakeShop.contract, event: "__resetStation", logs: logs, sub: sub}, nil
}

// WatchResetStation is a free log subscription operation binding the contract event 0x73d8ef4b77a16350f9072178ca21a79e0d68298007b2a92c870e2b2c66f7868a.
//
// Solidity: e __resetStation(offchain bytes32, hid uint256)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) WatchResetStation(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeShopResetStation) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.WatchLogs(opts, "__resetStation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeShopResetStation)
				if err := _ExchangeHandshakeShop.contract.UnpackLog(event, "__resetStation", log); err != nil {
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

// ExchangeHandshakeShopSetFeeIterator is returned from FilterSetFee and is used to iterate over the raw logs and unpacked data for SetFee events raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopSetFeeIterator struct {
	Event *ExchangeHandshakeShopSetFee // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeShopSetFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeShopSetFee)
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
		it.Event = new(ExchangeHandshakeShopSetFee)
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
func (it *ExchangeHandshakeShopSetFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeShopSetFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeShopSetFee represents a SetFee event raised by the ExchangeHandshakeShop contract.
type ExchangeHandshakeShopSetFee struct {
	Fee *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterSetFee is a free log retrieval operation binding the contract event 0x1941cda5868c57eb88923fbcebe63c45bf133819cf0720c99865546a2615f4a7.
//
// Solidity: e __setFee(fee uint256)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) FilterSetFee(opts *bind.FilterOpts) (*ExchangeHandshakeShopSetFeeIterator, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.FilterLogs(opts, "__setFee")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShopSetFeeIterator{contract: _ExchangeHandshakeShop.contract, event: "__setFee", logs: logs, sub: sub}, nil
}

// WatchSetFee is a free log subscription operation binding the contract event 0x1941cda5868c57eb88923fbcebe63c45bf133819cf0720c99865546a2615f4a7.
//
// Solidity: e __setFee(fee uint256)
func (_ExchangeHandshakeShop *ExchangeHandshakeShopFilterer) WatchSetFee(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeShopSetFee) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshakeShop.contract.WatchLogs(opts, "__setFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeShopSetFee)
				if err := _ExchangeHandshakeShop.contract.UnpackLog(event, "__setFee", log); err != nil {
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
