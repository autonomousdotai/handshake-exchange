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

// ExchangeHandshakeABI is the input ABI used to generate the binding from.
const ExchangeHandshakeABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ex\",\"outputs\":[{\"name\":\"coinOwner\",\"type\":\"address\"},{\"name\":\"cashOwner\",\"type\":\"address\"},{\"name\":\"exchanger\",\"type\":\"address\"},{\"name\":\"adrFeeRefund\",\"type\":\"address\"},{\"name\":\"fee\",\"type\":\"uint256\"},{\"name\":\"feeRefund\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"state\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"exchanger\",\"type\":\"address\"},{\"name\":\"adrFeeRefund\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"initByCoinOwner\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"exchanger\",\"type\":\"address\"},{\"name\":\"adrFeeRefund\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"initByCashOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"}],\"name\":\"getState\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"accept\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"f\",\"type\":\"uint256\"},{\"name\":\"fr\",\"type\":\"uint256\"}],\"name\":\"setFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"reject\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"closeByCoinOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"shake\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"closeByCashOwner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"cancel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"feeRefund\",\"type\":\"uint256\"}],\"name\":\"__setFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"coinOwner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__initByCoinOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"cashOwner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__initByCashOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__shake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__withdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__reject\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__accept\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__cancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__closeByCoinOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__closeByCashOwner\",\"type\":\"event\"}]"

// ExchangeHandshake is an auto generated Go binding around an Ethereum contract.
type ExchangeHandshake struct {
	ExchangeHandshakeCaller     // Read-only binding to the contract
	ExchangeHandshakeTransactor // Write-only binding to the contract
	ExchangeHandshakeFilterer   // Log filterer for contract events
}

// ExchangeHandshakeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExchangeHandshakeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeHandshakeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExchangeHandshakeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeHandshakeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExchangeHandshakeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeHandshakeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExchangeHandshakeSession struct {
	Contract     *ExchangeHandshake // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ExchangeHandshakeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExchangeHandshakeCallerSession struct {
	Contract *ExchangeHandshakeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ExchangeHandshakeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExchangeHandshakeTransactorSession struct {
	Contract     *ExchangeHandshakeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ExchangeHandshakeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExchangeHandshakeRaw struct {
	Contract *ExchangeHandshake // Generic contract binding to access the raw methods on
}

// ExchangeHandshakeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExchangeHandshakeCallerRaw struct {
	Contract *ExchangeHandshakeCaller // Generic read-only contract binding to access the raw methods on
}

// ExchangeHandshakeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExchangeHandshakeTransactorRaw struct {
	Contract *ExchangeHandshakeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExchangeHandshake creates a new instance of ExchangeHandshake, bound to a specific deployed contract.
func NewExchangeHandshake(address common.Address, backend bind.ContractBackend) (*ExchangeHandshake, error) {
	contract, err := bindExchangeHandshake(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshake{ExchangeHandshakeCaller: ExchangeHandshakeCaller{contract: contract}, ExchangeHandshakeTransactor: ExchangeHandshakeTransactor{contract: contract}, ExchangeHandshakeFilterer: ExchangeHandshakeFilterer{contract: contract}}, nil
}

// NewExchangeHandshakeCaller creates a new read-only instance of ExchangeHandshake, bound to a specific deployed contract.
func NewExchangeHandshakeCaller(address common.Address, caller bind.ContractCaller) (*ExchangeHandshakeCaller, error) {
	contract, err := bindExchangeHandshake(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeCaller{contract: contract}, nil
}

// NewExchangeHandshakeTransactor creates a new write-only instance of ExchangeHandshake, bound to a specific deployed contract.
func NewExchangeHandshakeTransactor(address common.Address, transactor bind.ContractTransactor) (*ExchangeHandshakeTransactor, error) {
	contract, err := bindExchangeHandshake(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeTransactor{contract: contract}, nil
}

// NewExchangeHandshakeFilterer creates a new log filterer instance of ExchangeHandshake, bound to a specific deployed contract.
func NewExchangeHandshakeFilterer(address common.Address, filterer bind.ContractFilterer) (*ExchangeHandshakeFilterer, error) {
	contract, err := bindExchangeHandshake(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeFilterer{contract: contract}, nil
}

// bindExchangeHandshake binds a generic wrapper to an already deployed contract.
func bindExchangeHandshake(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangeHandshakeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExchangeHandshake *ExchangeHandshakeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ExchangeHandshake.Contract.ExchangeHandshakeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExchangeHandshake *ExchangeHandshakeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.ExchangeHandshakeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExchangeHandshake *ExchangeHandshakeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.ExchangeHandshakeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExchangeHandshake *ExchangeHandshakeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ExchangeHandshake.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExchangeHandshake *ExchangeHandshakeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExchangeHandshake *ExchangeHandshakeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.contract.Transact(opts, method, params...)
}

// Ex is a free data retrieval call binding the contract method 0x1089f215.
//
// Solidity: function ex( uint256) constant returns(coinOwner address, cashOwner address, exchanger address, adrFeeRefund address, fee uint256, feeRefund uint256, value uint256, state uint8)
func (_ExchangeHandshake *ExchangeHandshakeCaller) Ex(opts *bind.CallOpts, arg0 *big.Int) (struct {
	CoinOwner    common.Address
	CashOwner    common.Address
	Exchanger    common.Address
	AdrFeeRefund common.Address
	Fee          *big.Int
	FeeRefund    *big.Int
	Value        *big.Int
	State        uint8
}, error) {
	ret := new(struct {
		CoinOwner    common.Address
		CashOwner    common.Address
		Exchanger    common.Address
		AdrFeeRefund common.Address
		Fee          *big.Int
		FeeRefund    *big.Int
		Value        *big.Int
		State        uint8
	})
	out := ret
	err := _ExchangeHandshake.contract.Call(opts, out, "ex", arg0)
	return *ret, err
}

// Ex is a free data retrieval call binding the contract method 0x1089f215.
//
// Solidity: function ex( uint256) constant returns(coinOwner address, cashOwner address, exchanger address, adrFeeRefund address, fee uint256, feeRefund uint256, value uint256, state uint8)
func (_ExchangeHandshake *ExchangeHandshakeSession) Ex(arg0 *big.Int) (struct {
	CoinOwner    common.Address
	CashOwner    common.Address
	Exchanger    common.Address
	AdrFeeRefund common.Address
	Fee          *big.Int
	FeeRefund    *big.Int
	Value        *big.Int
	State        uint8
}, error) {
	return _ExchangeHandshake.Contract.Ex(&_ExchangeHandshake.CallOpts, arg0)
}

// Ex is a free data retrieval call binding the contract method 0x1089f215.
//
// Solidity: function ex( uint256) constant returns(coinOwner address, cashOwner address, exchanger address, adrFeeRefund address, fee uint256, feeRefund uint256, value uint256, state uint8)
func (_ExchangeHandshake *ExchangeHandshakeCallerSession) Ex(arg0 *big.Int) (struct {
	CoinOwner    common.Address
	CashOwner    common.Address
	Exchanger    common.Address
	AdrFeeRefund common.Address
	Fee          *big.Int
	FeeRefund    *big.Int
	Value        *big.Int
	State        uint8
}, error) {
	return _ExchangeHandshake.Contract.Ex(&_ExchangeHandshake.CallOpts, arg0)
}

// GetState is a free data retrieval call binding the contract method 0x44c9af28.
//
// Solidity: function getState(hid uint256) constant returns(uint8)
func (_ExchangeHandshake *ExchangeHandshakeCaller) GetState(opts *bind.CallOpts, hid *big.Int) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _ExchangeHandshake.contract.Call(opts, out, "getState", hid)
	return *ret0, err
}

// GetState is a free data retrieval call binding the contract method 0x44c9af28.
//
// Solidity: function getState(hid uint256) constant returns(uint8)
func (_ExchangeHandshake *ExchangeHandshakeSession) GetState(hid *big.Int) (uint8, error) {
	return _ExchangeHandshake.Contract.GetState(&_ExchangeHandshake.CallOpts, hid)
}

// GetState is a free data retrieval call binding the contract method 0x44c9af28.
//
// Solidity: function getState(hid uint256) constant returns(uint8)
func (_ExchangeHandshake *ExchangeHandshakeCallerSession) GetState(hid *big.Int) (uint8, error) {
	return _ExchangeHandshake.Contract.GetState(&_ExchangeHandshake.CallOpts, hid)
}

// Accept is a paid mutator transaction binding the contract method 0x5203dcda.
//
// Solidity: function accept(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactor) Accept(opts *bind.TransactOpts, hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.contract.Transact(opts, "accept", hid, offchain)
}

// Accept is a paid mutator transaction binding the contract method 0x5203dcda.
//
// Solidity: function accept(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeSession) Accept(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.Accept(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// Accept is a paid mutator transaction binding the contract method 0x5203dcda.
//
// Solidity: function accept(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactorSession) Accept(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.Accept(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// Cancel is a paid mutator transaction binding the contract method 0xeafb64d5.
//
// Solidity: function cancel(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactor) Cancel(opts *bind.TransactOpts, hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.contract.Transact(opts, "cancel", hid, offchain)
}

// Cancel is a paid mutator transaction binding the contract method 0xeafb64d5.
//
// Solidity: function cancel(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeSession) Cancel(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.Cancel(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// Cancel is a paid mutator transaction binding the contract method 0xeafb64d5.
//
// Solidity: function cancel(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactorSession) Cancel(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.Cancel(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// CloseByCashOwner is a paid mutator transaction binding the contract method 0xbf1e459f.
//
// Solidity: function closeByCashOwner(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactor) CloseByCashOwner(opts *bind.TransactOpts, hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.contract.Transact(opts, "closeByCashOwner", hid, offchain)
}

// CloseByCashOwner is a paid mutator transaction binding the contract method 0xbf1e459f.
//
// Solidity: function closeByCashOwner(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeSession) CloseByCashOwner(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.CloseByCashOwner(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// CloseByCashOwner is a paid mutator transaction binding the contract method 0xbf1e459f.
//
// Solidity: function closeByCashOwner(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactorSession) CloseByCashOwner(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.CloseByCashOwner(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// CloseByCoinOwner is a paid mutator transaction binding the contract method 0x9a8c21bf.
//
// Solidity: function closeByCoinOwner(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactor) CloseByCoinOwner(opts *bind.TransactOpts, hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.contract.Transact(opts, "closeByCoinOwner", hid, offchain)
}

// CloseByCoinOwner is a paid mutator transaction binding the contract method 0x9a8c21bf.
//
// Solidity: function closeByCoinOwner(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeSession) CloseByCoinOwner(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.CloseByCoinOwner(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// CloseByCoinOwner is a paid mutator transaction binding the contract method 0x9a8c21bf.
//
// Solidity: function closeByCoinOwner(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactorSession) CloseByCoinOwner(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.CloseByCoinOwner(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// InitByCashOwner is a paid mutator transaction binding the contract method 0x3ebad1a2.
//
// Solidity: function initByCashOwner(exchanger address, adrFeeRefund address, value uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactor) InitByCashOwner(opts *bind.TransactOpts, exchanger common.Address, adrFeeRefund common.Address, value *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.contract.Transact(opts, "initByCashOwner", exchanger, adrFeeRefund, value, offchain)
}

// InitByCashOwner is a paid mutator transaction binding the contract method 0x3ebad1a2.
//
// Solidity: function initByCashOwner(exchanger address, adrFeeRefund address, value uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeSession) InitByCashOwner(exchanger common.Address, adrFeeRefund common.Address, value *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.InitByCashOwner(&_ExchangeHandshake.TransactOpts, exchanger, adrFeeRefund, value, offchain)
}

// InitByCashOwner is a paid mutator transaction binding the contract method 0x3ebad1a2.
//
// Solidity: function initByCashOwner(exchanger address, adrFeeRefund address, value uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactorSession) InitByCashOwner(exchanger common.Address, adrFeeRefund common.Address, value *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.InitByCashOwner(&_ExchangeHandshake.TransactOpts, exchanger, adrFeeRefund, value, offchain)
}

// InitByCoinOwner is a paid mutator transaction binding the contract method 0x1178a48a.
//
// Solidity: function initByCoinOwner(exchanger address, adrFeeRefund address, value uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactor) InitByCoinOwner(opts *bind.TransactOpts, exchanger common.Address, adrFeeRefund common.Address, value *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.contract.Transact(opts, "initByCoinOwner", exchanger, adrFeeRefund, value, offchain)
}

// InitByCoinOwner is a paid mutator transaction binding the contract method 0x1178a48a.
//
// Solidity: function initByCoinOwner(exchanger address, adrFeeRefund address, value uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeSession) InitByCoinOwner(exchanger common.Address, adrFeeRefund common.Address, value *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.InitByCoinOwner(&_ExchangeHandshake.TransactOpts, exchanger, adrFeeRefund, value, offchain)
}

// InitByCoinOwner is a paid mutator transaction binding the contract method 0x1178a48a.
//
// Solidity: function initByCoinOwner(exchanger address, adrFeeRefund address, value uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactorSession) InitByCoinOwner(exchanger common.Address, adrFeeRefund common.Address, value *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.InitByCoinOwner(&_ExchangeHandshake.TransactOpts, exchanger, adrFeeRefund, value, offchain)
}

// Reject is a paid mutator transaction binding the contract method 0x6be1320b.
//
// Solidity: function reject(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactor) Reject(opts *bind.TransactOpts, hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.contract.Transact(opts, "reject", hid, offchain)
}

// Reject is a paid mutator transaction binding the contract method 0x6be1320b.
//
// Solidity: function reject(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeSession) Reject(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.Reject(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// Reject is a paid mutator transaction binding the contract method 0x6be1320b.
//
// Solidity: function reject(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactorSession) Reject(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.Reject(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// SetFee is a paid mutator transaction binding the contract method 0x52f7c988.
//
// Solidity: function setFee(f uint256, fr uint256) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactor) SetFee(opts *bind.TransactOpts, f *big.Int, fr *big.Int) (*types.Transaction, error) {
	return _ExchangeHandshake.contract.Transact(opts, "setFee", f, fr)
}

// SetFee is a paid mutator transaction binding the contract method 0x52f7c988.
//
// Solidity: function setFee(f uint256, fr uint256) returns()
func (_ExchangeHandshake *ExchangeHandshakeSession) SetFee(f *big.Int, fr *big.Int) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.SetFee(&_ExchangeHandshake.TransactOpts, f, fr)
}

// SetFee is a paid mutator transaction binding the contract method 0x52f7c988.
//
// Solidity: function setFee(f uint256, fr uint256) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactorSession) SetFee(f *big.Int, fr *big.Int) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.SetFee(&_ExchangeHandshake.TransactOpts, f, fr)
}

// Shake is a paid mutator transaction binding the contract method 0xb09b2f85.
//
// Solidity: function shake(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactor) Shake(opts *bind.TransactOpts, hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.contract.Transact(opts, "shake", hid, offchain)
}

// Shake is a paid mutator transaction binding the contract method 0xb09b2f85.
//
// Solidity: function shake(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeSession) Shake(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.Shake(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// Shake is a paid mutator transaction binding the contract method 0xb09b2f85.
//
// Solidity: function shake(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactorSession) Shake(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.Shake(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// Withdraw is a paid mutator transaction binding the contract method 0xa8d2021a.
//
// Solidity: function withdraw(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactor) Withdraw(opts *bind.TransactOpts, hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.contract.Transact(opts, "withdraw", hid, offchain)
}

// Withdraw is a paid mutator transaction binding the contract method 0xa8d2021a.
//
// Solidity: function withdraw(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeSession) Withdraw(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.Withdraw(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// Withdraw is a paid mutator transaction binding the contract method 0xa8d2021a.
//
// Solidity: function withdraw(hid uint256, offchain bytes32) returns()
func (_ExchangeHandshake *ExchangeHandshakeTransactorSession) Withdraw(hid *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _ExchangeHandshake.Contract.Withdraw(&_ExchangeHandshake.TransactOpts, hid, offchain)
}

// ExchangeHandshakeAcceptIterator is returned from FilterAccept and is used to iterate over the raw logs and unpacked data for Accept events raised by the ExchangeHandshake contract.
type ExchangeHandshakeAcceptIterator struct {
	Event *ExchangeHandshakeAccept // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeAcceptIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeAccept)
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
		it.Event = new(ExchangeHandshakeAccept)
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
func (it *ExchangeHandshakeAcceptIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeAcceptIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeAccept represents a Accept event raised by the ExchangeHandshake contract.
type ExchangeHandshakeAccept struct {
	Hid      *big.Int
	Offchain [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAccept is a free log retrieval operation binding the contract event 0xbda7bc7c8123a85aa855c777d3191b2dc42bec38c45638643006fb84e76abf7a.
//
// Solidity: e __accept(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) FilterAccept(opts *bind.FilterOpts) (*ExchangeHandshakeAcceptIterator, error) {

	logs, sub, err := _ExchangeHandshake.contract.FilterLogs(opts, "__accept")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeAcceptIterator{contract: _ExchangeHandshake.contract, event: "__accept", logs: logs, sub: sub}, nil
}

// WatchAccept is a free log subscription operation binding the contract event 0xbda7bc7c8123a85aa855c777d3191b2dc42bec38c45638643006fb84e76abf7a.
//
// Solidity: e __accept(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) WatchAccept(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeAccept) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshake.contract.WatchLogs(opts, "__accept")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeAccept)
				if err := _ExchangeHandshake.contract.UnpackLog(event, "__accept", log); err != nil {
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

// ExchangeHandshakeCancelIterator is returned from FilterCancel and is used to iterate over the raw logs and unpacked data for Cancel events raised by the ExchangeHandshake contract.
type ExchangeHandshakeCancelIterator struct {
	Event *ExchangeHandshakeCancel // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeCancelIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeCancel)
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
		it.Event = new(ExchangeHandshakeCancel)
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
func (it *ExchangeHandshakeCancelIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeCancelIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeCancel represents a Cancel event raised by the ExchangeHandshake contract.
type ExchangeHandshakeCancel struct {
	Hid      *big.Int
	Offchain [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCancel is a free log retrieval operation binding the contract event 0xcb720f90f098c425fad2e9df556017c076f9bf7aefa096d1c904a06027ae0460.
//
// Solidity: e __cancel(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) FilterCancel(opts *bind.FilterOpts) (*ExchangeHandshakeCancelIterator, error) {

	logs, sub, err := _ExchangeHandshake.contract.FilterLogs(opts, "__cancel")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeCancelIterator{contract: _ExchangeHandshake.contract, event: "__cancel", logs: logs, sub: sub}, nil
}

// WatchCancel is a free log subscription operation binding the contract event 0xcb720f90f098c425fad2e9df556017c076f9bf7aefa096d1c904a06027ae0460.
//
// Solidity: e __cancel(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) WatchCancel(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeCancel) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshake.contract.WatchLogs(opts, "__cancel")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeCancel)
				if err := _ExchangeHandshake.contract.UnpackLog(event, "__cancel", log); err != nil {
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

// ExchangeHandshakeCloseByCashOwnerIterator is returned from FilterCloseByCashOwner and is used to iterate over the raw logs and unpacked data for CloseByCashOwner events raised by the ExchangeHandshake contract.
type ExchangeHandshakeCloseByCashOwnerIterator struct {
	Event *ExchangeHandshakeCloseByCashOwner // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeCloseByCashOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeCloseByCashOwner)
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
		it.Event = new(ExchangeHandshakeCloseByCashOwner)
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
func (it *ExchangeHandshakeCloseByCashOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeCloseByCashOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeCloseByCashOwner represents a CloseByCashOwner event raised by the ExchangeHandshake contract.
type ExchangeHandshakeCloseByCashOwner struct {
	Hid      *big.Int
	Offchain [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCloseByCashOwner is a free log retrieval operation binding the contract event 0x3d2fb2dfea1905d22bfa4e006af321688700e6bf97ce5156cfd9a10dc4735dc6.
//
// Solidity: e __closeByCashOwner(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) FilterCloseByCashOwner(opts *bind.FilterOpts) (*ExchangeHandshakeCloseByCashOwnerIterator, error) {

	logs, sub, err := _ExchangeHandshake.contract.FilterLogs(opts, "__closeByCashOwner")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeCloseByCashOwnerIterator{contract: _ExchangeHandshake.contract, event: "__closeByCashOwner", logs: logs, sub: sub}, nil
}

// WatchCloseByCashOwner is a free log subscription operation binding the contract event 0x3d2fb2dfea1905d22bfa4e006af321688700e6bf97ce5156cfd9a10dc4735dc6.
//
// Solidity: e __closeByCashOwner(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) WatchCloseByCashOwner(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeCloseByCashOwner) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshake.contract.WatchLogs(opts, "__closeByCashOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeCloseByCashOwner)
				if err := _ExchangeHandshake.contract.UnpackLog(event, "__closeByCashOwner", log); err != nil {
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

// ExchangeHandshakeCloseByCoinOwnerIterator is returned from FilterCloseByCoinOwner and is used to iterate over the raw logs and unpacked data for CloseByCoinOwner events raised by the ExchangeHandshake contract.
type ExchangeHandshakeCloseByCoinOwnerIterator struct {
	Event *ExchangeHandshakeCloseByCoinOwner // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeCloseByCoinOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeCloseByCoinOwner)
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
		it.Event = new(ExchangeHandshakeCloseByCoinOwner)
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
func (it *ExchangeHandshakeCloseByCoinOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeCloseByCoinOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeCloseByCoinOwner represents a CloseByCoinOwner event raised by the ExchangeHandshake contract.
type ExchangeHandshakeCloseByCoinOwner struct {
	Hid      *big.Int
	Offchain [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCloseByCoinOwner is a free log retrieval operation binding the contract event 0xdeefbe59da61b7848e0f3b1556ba1c6b10a35470f9916e5b59ec8dc7197d500b.
//
// Solidity: e __closeByCoinOwner(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) FilterCloseByCoinOwner(opts *bind.FilterOpts) (*ExchangeHandshakeCloseByCoinOwnerIterator, error) {

	logs, sub, err := _ExchangeHandshake.contract.FilterLogs(opts, "__closeByCoinOwner")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeCloseByCoinOwnerIterator{contract: _ExchangeHandshake.contract, event: "__closeByCoinOwner", logs: logs, sub: sub}, nil
}

// WatchCloseByCoinOwner is a free log subscription operation binding the contract event 0xdeefbe59da61b7848e0f3b1556ba1c6b10a35470f9916e5b59ec8dc7197d500b.
//
// Solidity: e __closeByCoinOwner(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) WatchCloseByCoinOwner(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeCloseByCoinOwner) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshake.contract.WatchLogs(opts, "__closeByCoinOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeCloseByCoinOwner)
				if err := _ExchangeHandshake.contract.UnpackLog(event, "__closeByCoinOwner", log); err != nil {
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

// ExchangeHandshakeInitByCashOwnerIterator is returned from FilterInitByCashOwner and is used to iterate over the raw logs and unpacked data for InitByCashOwner events raised by the ExchangeHandshake contract.
type ExchangeHandshakeInitByCashOwnerIterator struct {
	Event *ExchangeHandshakeInitByCashOwner // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeInitByCashOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeInitByCashOwner)
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
		it.Event = new(ExchangeHandshakeInitByCashOwner)
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
func (it *ExchangeHandshakeInitByCashOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeInitByCashOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeInitByCashOwner represents a InitByCashOwner event raised by the ExchangeHandshake contract.
type ExchangeHandshakeInitByCashOwner struct {
	Hid       *big.Int
	CashOwner common.Address
	Offchain  [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterInitByCashOwner is a free log retrieval operation binding the contract event 0xe4601f4f89d9b0b9deb2ee9ecb6f89b91ef95b0ae1f73a28e49b1159aeec729f.
//
// Solidity: e __initByCashOwner(hid uint256, cashOwner address, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) FilterInitByCashOwner(opts *bind.FilterOpts) (*ExchangeHandshakeInitByCashOwnerIterator, error) {

	logs, sub, err := _ExchangeHandshake.contract.FilterLogs(opts, "__initByCashOwner")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeInitByCashOwnerIterator{contract: _ExchangeHandshake.contract, event: "__initByCashOwner", logs: logs, sub: sub}, nil
}

// WatchInitByCashOwner is a free log subscription operation binding the contract event 0xe4601f4f89d9b0b9deb2ee9ecb6f89b91ef95b0ae1f73a28e49b1159aeec729f.
//
// Solidity: e __initByCashOwner(hid uint256, cashOwner address, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) WatchInitByCashOwner(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeInitByCashOwner) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshake.contract.WatchLogs(opts, "__initByCashOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeInitByCashOwner)
				if err := _ExchangeHandshake.contract.UnpackLog(event, "__initByCashOwner", log); err != nil {
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

// ExchangeHandshakeInitByCoinOwnerIterator is returned from FilterInitByCoinOwner and is used to iterate over the raw logs and unpacked data for InitByCoinOwner events raised by the ExchangeHandshake contract.
type ExchangeHandshakeInitByCoinOwnerIterator struct {
	Event *ExchangeHandshakeInitByCoinOwner // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeInitByCoinOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeInitByCoinOwner)
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
		it.Event = new(ExchangeHandshakeInitByCoinOwner)
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
func (it *ExchangeHandshakeInitByCoinOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeInitByCoinOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeInitByCoinOwner represents a InitByCoinOwner event raised by the ExchangeHandshake contract.
type ExchangeHandshakeInitByCoinOwner struct {
	Hid       *big.Int
	CoinOwner common.Address
	Offchain  [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterInitByCoinOwner is a free log retrieval operation binding the contract event 0x2fd154f5b4d50fb27d0311b078ca858c2a8d52d40cad8d62881f416b0b7a469f.
//
// Solidity: e __initByCoinOwner(hid uint256, coinOwner address, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) FilterInitByCoinOwner(opts *bind.FilterOpts) (*ExchangeHandshakeInitByCoinOwnerIterator, error) {

	logs, sub, err := _ExchangeHandshake.contract.FilterLogs(opts, "__initByCoinOwner")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeInitByCoinOwnerIterator{contract: _ExchangeHandshake.contract, event: "__initByCoinOwner", logs: logs, sub: sub}, nil
}

// WatchInitByCoinOwner is a free log subscription operation binding the contract event 0x2fd154f5b4d50fb27d0311b078ca858c2a8d52d40cad8d62881f416b0b7a469f.
//
// Solidity: e __initByCoinOwner(hid uint256, coinOwner address, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) WatchInitByCoinOwner(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeInitByCoinOwner) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshake.contract.WatchLogs(opts, "__initByCoinOwner")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeInitByCoinOwner)
				if err := _ExchangeHandshake.contract.UnpackLog(event, "__initByCoinOwner", log); err != nil {
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

// ExchangeHandshakeRejectIterator is returned from FilterReject and is used to iterate over the raw logs and unpacked data for Reject events raised by the ExchangeHandshake contract.
type ExchangeHandshakeRejectIterator struct {
	Event *ExchangeHandshakeReject // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeRejectIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeReject)
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
		it.Event = new(ExchangeHandshakeReject)
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
func (it *ExchangeHandshakeRejectIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeRejectIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeReject represents a Reject event raised by the ExchangeHandshake contract.
type ExchangeHandshakeReject struct {
	Hid      *big.Int
	Offchain [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReject is a free log retrieval operation binding the contract event 0xae76720f3a5d319b91bc94d8a6c2e3096a4f3554c8cb897e3aedfced5824a10a.
//
// Solidity: e __reject(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) FilterReject(opts *bind.FilterOpts) (*ExchangeHandshakeRejectIterator, error) {

	logs, sub, err := _ExchangeHandshake.contract.FilterLogs(opts, "__reject")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeRejectIterator{contract: _ExchangeHandshake.contract, event: "__reject", logs: logs, sub: sub}, nil
}

// WatchReject is a free log subscription operation binding the contract event 0xae76720f3a5d319b91bc94d8a6c2e3096a4f3554c8cb897e3aedfced5824a10a.
//
// Solidity: e __reject(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) WatchReject(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeReject) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshake.contract.WatchLogs(opts, "__reject")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeReject)
				if err := _ExchangeHandshake.contract.UnpackLog(event, "__reject", log); err != nil {
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

// ExchangeHandshakeSetFeeIterator is returned from FilterSetFee and is used to iterate over the raw logs and unpacked data for SetFee events raised by the ExchangeHandshake contract.
type ExchangeHandshakeSetFeeIterator struct {
	Event *ExchangeHandshakeSetFee // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeSetFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeSetFee)
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
		it.Event = new(ExchangeHandshakeSetFee)
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
func (it *ExchangeHandshakeSetFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeSetFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeSetFee represents a SetFee event raised by the ExchangeHandshake contract.
type ExchangeHandshakeSetFee struct {
	Fee       *big.Int
	FeeRefund *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSetFee is a free log retrieval operation binding the contract event 0xda33b66207af71514f4eb8f9fee7b74ba441e293e937e20f099781f435c2786b.
//
// Solidity: e __setFee(fee uint256, feeRefund uint256)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) FilterSetFee(opts *bind.FilterOpts) (*ExchangeHandshakeSetFeeIterator, error) {

	logs, sub, err := _ExchangeHandshake.contract.FilterLogs(opts, "__setFee")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeSetFeeIterator{contract: _ExchangeHandshake.contract, event: "__setFee", logs: logs, sub: sub}, nil
}

// WatchSetFee is a free log subscription operation binding the contract event 0xda33b66207af71514f4eb8f9fee7b74ba441e293e937e20f099781f435c2786b.
//
// Solidity: e __setFee(fee uint256, feeRefund uint256)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) WatchSetFee(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeSetFee) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshake.contract.WatchLogs(opts, "__setFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeSetFee)
				if err := _ExchangeHandshake.contract.UnpackLog(event, "__setFee", log); err != nil {
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

// ExchangeHandshakeShakeIterator is returned from FilterShake and is used to iterate over the raw logs and unpacked data for Shake events raised by the ExchangeHandshake contract.
type ExchangeHandshakeShakeIterator struct {
	Event *ExchangeHandshakeShake // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeShakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeShake)
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
		it.Event = new(ExchangeHandshakeShake)
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
func (it *ExchangeHandshakeShakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeShakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeShake represents a Shake event raised by the ExchangeHandshake contract.
type ExchangeHandshakeShake struct {
	Hid      *big.Int
	Offchain [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterShake is a free log retrieval operation binding the contract event 0x6a5ee55c9df2daa4375d2b5e4ec8b9e5662f1863207bcbe6e38c6f5fe3c24300.
//
// Solidity: e __shake(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) FilterShake(opts *bind.FilterOpts) (*ExchangeHandshakeShakeIterator, error) {

	logs, sub, err := _ExchangeHandshake.contract.FilterLogs(opts, "__shake")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeShakeIterator{contract: _ExchangeHandshake.contract, event: "__shake", logs: logs, sub: sub}, nil
}

// WatchShake is a free log subscription operation binding the contract event 0x6a5ee55c9df2daa4375d2b5e4ec8b9e5662f1863207bcbe6e38c6f5fe3c24300.
//
// Solidity: e __shake(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) WatchShake(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeShake) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshake.contract.WatchLogs(opts, "__shake")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeShake)
				if err := _ExchangeHandshake.contract.UnpackLog(event, "__shake", log); err != nil {
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

// ExchangeHandshakeWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the ExchangeHandshake contract.
type ExchangeHandshakeWithdrawIterator struct {
	Event *ExchangeHandshakeWithdraw // Event containing the contract specifics and raw log

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
func (it *ExchangeHandshakeWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeHandshakeWithdraw)
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
		it.Event = new(ExchangeHandshakeWithdraw)
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
func (it *ExchangeHandshakeWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeHandshakeWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeHandshakeWithdraw represents a Withdraw event raised by the ExchangeHandshake contract.
type ExchangeHandshakeWithdraw struct {
	Hid      *big.Int
	Offchain [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x86e48a712d324364a79089a18be31e993ed6bda36550b789f988e7aaf9ed7cf8.
//
// Solidity: e __withdraw(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) FilterWithdraw(opts *bind.FilterOpts) (*ExchangeHandshakeWithdrawIterator, error) {

	logs, sub, err := _ExchangeHandshake.contract.FilterLogs(opts, "__withdraw")
	if err != nil {
		return nil, err
	}
	return &ExchangeHandshakeWithdrawIterator{contract: _ExchangeHandshake.contract, event: "__withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x86e48a712d324364a79089a18be31e993ed6bda36550b789f988e7aaf9ed7cf8.
//
// Solidity: e __withdraw(hid uint256, offchain bytes32)
func (_ExchangeHandshake *ExchangeHandshakeFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *ExchangeHandshakeWithdraw) (event.Subscription, error) {

	logs, sub, err := _ExchangeHandshake.contract.WatchLogs(opts, "__withdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeHandshakeWithdraw)
				if err := _ExchangeHandshake.contract.UnpackLog(event, "__withdraw", log); err != nil {
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
