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

// CreditATMABI is the input ABI used to generate the binding from.
const CreditATMABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ex\",\"outputs\":[{\"name\":\"sender\",\"type\":\"address\"},{\"name\":\"percentage\",\"type\":\"uint256\"},{\"name\":\"escrow\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"offchain\",\"type\":\"bytes32\"},{\"name\":\"percentage\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"},{\"name\":\"customer\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"releasePartialFund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hid\",\"type\":\"uint256\"}],\"name\":\"getDepositList\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"stationOwner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"percentage\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"hid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"customer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"offchain\",\"type\":\"bytes32\"}],\"name\":\"__releasePartialFund\",\"type\":\"event\"}]"

// CreditATM is an auto generated Go binding around an Ethereum contract.
type CreditATM struct {
	CreditATMCaller     // Read-only binding to the contract
	CreditATMTransactor // Write-only binding to the contract
	CreditATMFilterer   // Log filterer for contract events
}

// CreditATMCaller is an auto generated read-only Go binding around an Ethereum contract.
type CreditATMCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CreditATMTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CreditATMTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CreditATMFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CreditATMFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CreditATMSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CreditATMSession struct {
	Contract     *CreditATM        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CreditATMCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CreditATMCallerSession struct {
	Contract *CreditATMCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// CreditATMTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CreditATMTransactorSession struct {
	Contract     *CreditATMTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// CreditATMRaw is an auto generated low-level Go binding around an Ethereum contract.
type CreditATMRaw struct {
	Contract *CreditATM // Generic contract binding to access the raw methods on
}

// CreditATMCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CreditATMCallerRaw struct {
	Contract *CreditATMCaller // Generic read-only contract binding to access the raw methods on
}

// CreditATMTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CreditATMTransactorRaw struct {
	Contract *CreditATMTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCreditATM creates a new instance of CreditATM, bound to a specific deployed contract.
func NewCreditATM(address common.Address, backend bind.ContractBackend) (*CreditATM, error) {
	contract, err := bindCreditATM(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CreditATM{CreditATMCaller: CreditATMCaller{contract: contract}, CreditATMTransactor: CreditATMTransactor{contract: contract}, CreditATMFilterer: CreditATMFilterer{contract: contract}}, nil
}

// NewCreditATMCaller creates a new read-only instance of CreditATM, bound to a specific deployed contract.
func NewCreditATMCaller(address common.Address, caller bind.ContractCaller) (*CreditATMCaller, error) {
	contract, err := bindCreditATM(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CreditATMCaller{contract: contract}, nil
}

// NewCreditATMTransactor creates a new write-only instance of CreditATM, bound to a specific deployed contract.
func NewCreditATMTransactor(address common.Address, transactor bind.ContractTransactor) (*CreditATMTransactor, error) {
	contract, err := bindCreditATM(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CreditATMTransactor{contract: contract}, nil
}

// NewCreditATMFilterer creates a new log filterer instance of CreditATM, bound to a specific deployed contract.
func NewCreditATMFilterer(address common.Address, filterer bind.ContractFilterer) (*CreditATMFilterer, error) {
	contract, err := bindCreditATM(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CreditATMFilterer{contract: contract}, nil
}

// bindCreditATM binds a generic wrapper to an already deployed contract.
func bindCreditATM(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CreditATMABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CreditATM *CreditATMRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CreditATM.Contract.CreditATMCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CreditATM *CreditATMRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CreditATM.Contract.CreditATMTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CreditATM *CreditATMRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CreditATM.Contract.CreditATMTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CreditATM *CreditATMCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CreditATM.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CreditATM *CreditATMTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CreditATM.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CreditATM *CreditATMTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CreditATM.Contract.contract.Transact(opts, method, params...)
}

// Ex is a free data retrieval call binding the contract method 0x1089f215.
//
// Solidity: function ex( uint256) constant returns(sender address, percentage uint256, escrow uint256)
func (_CreditATM *CreditATMCaller) Ex(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Sender     common.Address
	Percentage *big.Int
	Escrow     *big.Int
}, error) {
	ret := new(struct {
		Sender     common.Address
		Percentage *big.Int
		Escrow     *big.Int
	})
	out := ret
	err := _CreditATM.contract.Call(opts, out, "ex", arg0)
	return *ret, err
}

// Ex is a free data retrieval call binding the contract method 0x1089f215.
//
// Solidity: function ex( uint256) constant returns(sender address, percentage uint256, escrow uint256)
func (_CreditATM *CreditATMSession) Ex(arg0 *big.Int) (struct {
	Sender     common.Address
	Percentage *big.Int
	Escrow     *big.Int
}, error) {
	return _CreditATM.Contract.Ex(&_CreditATM.CallOpts, arg0)
}

// Ex is a free data retrieval call binding the contract method 0x1089f215.
//
// Solidity: function ex( uint256) constant returns(sender address, percentage uint256, escrow uint256)
func (_CreditATM *CreditATMCallerSession) Ex(arg0 *big.Int) (struct {
	Sender     common.Address
	Percentage *big.Int
	Escrow     *big.Int
}, error) {
	return _CreditATM.Contract.Ex(&_CreditATM.CallOpts, arg0)
}

// GetDepositList is a free data retrieval call binding the contract method 0x99058690.
//
// Solidity: function getDepositList(hid uint256) constant returns(address, uint256, uint256)
func (_CreditATM *CreditATMCaller) GetDepositList(opts *bind.CallOpts, hid *big.Int) (common.Address, *big.Int, *big.Int, error) {
	var (
		ret0 = new(common.Address)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _CreditATM.contract.Call(opts, out, "getDepositList", hid)
	return *ret0, *ret1, *ret2, err
}

// GetDepositList is a free data retrieval call binding the contract method 0x99058690.
//
// Solidity: function getDepositList(hid uint256) constant returns(address, uint256, uint256)
func (_CreditATM *CreditATMSession) GetDepositList(hid *big.Int) (common.Address, *big.Int, *big.Int, error) {
	return _CreditATM.Contract.GetDepositList(&_CreditATM.CallOpts, hid)
}

// GetDepositList is a free data retrieval call binding the contract method 0x99058690.
//
// Solidity: function getDepositList(hid uint256) constant returns(address, uint256, uint256)
func (_CreditATM *CreditATMCallerSession) GetDepositList(hid *big.Int) (common.Address, *big.Int, *big.Int, error) {
	return _CreditATM.Contract.GetDepositList(&_CreditATM.CallOpts, hid)
}

// Deposit is a paid mutator transaction binding the contract method 0x1de26e16.
//
// Solidity: function deposit(offchain bytes32, percentage uint256) returns()
func (_CreditATM *CreditATMTransactor) Deposit(opts *bind.TransactOpts, offchain [32]byte, percentage *big.Int) (*types.Transaction, error) {
	return _CreditATM.contract.Transact(opts, "deposit", offchain, percentage)
}

// Deposit is a paid mutator transaction binding the contract method 0x1de26e16.
//
// Solidity: function deposit(offchain bytes32, percentage uint256) returns()
func (_CreditATM *CreditATMSession) Deposit(offchain [32]byte, percentage *big.Int) (*types.Transaction, error) {
	return _CreditATM.Contract.Deposit(&_CreditATM.TransactOpts, offchain, percentage)
}

// Deposit is a paid mutator transaction binding the contract method 0x1de26e16.
//
// Solidity: function deposit(offchain bytes32, percentage uint256) returns()
func (_CreditATM *CreditATMTransactorSession) Deposit(offchain [32]byte, percentage *big.Int) (*types.Transaction, error) {
	return _CreditATM.Contract.Deposit(&_CreditATM.TransactOpts, offchain, percentage)
}

// ReleasePartialFund is a paid mutator transaction binding the contract method 0x63d26f2f.
//
// Solidity: function releasePartialFund(hid uint256, customer address, amount uint256, offchain bytes32) returns()
func (_CreditATM *CreditATMTransactor) ReleasePartialFund(opts *bind.TransactOpts, hid *big.Int, customer common.Address, amount *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _CreditATM.contract.Transact(opts, "releasePartialFund", hid, customer, amount, offchain)
}

// ReleasePartialFund is a paid mutator transaction binding the contract method 0x63d26f2f.
//
// Solidity: function releasePartialFund(hid uint256, customer address, amount uint256, offchain bytes32) returns()
func (_CreditATM *CreditATMSession) ReleasePartialFund(hid *big.Int, customer common.Address, amount *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _CreditATM.Contract.ReleasePartialFund(&_CreditATM.TransactOpts, hid, customer, amount, offchain)
}

// ReleasePartialFund is a paid mutator transaction binding the contract method 0x63d26f2f.
//
// Solidity: function releasePartialFund(hid uint256, customer address, amount uint256, offchain bytes32) returns()
func (_CreditATM *CreditATMTransactorSession) ReleasePartialFund(hid *big.Int, customer common.Address, amount *big.Int, offchain [32]byte) (*types.Transaction, error) {
	return _CreditATM.Contract.ReleasePartialFund(&_CreditATM.TransactOpts, hid, customer, amount, offchain)
}

// CreditATMDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the CreditATM contract.
type CreditATMDepositIterator struct {
	Event *CreditATMDeposit // Event containing the contract specifics and raw log

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
func (it *CreditATMDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CreditATMDeposit)
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
		it.Event = new(CreditATMDeposit)
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
func (it *CreditATMDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CreditATMDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CreditATMDeposit represents a Deposit event raised by the CreditATM contract.
type CreditATMDeposit struct {
	Hid          *big.Int
	StationOwner common.Address
	Value        *big.Int
	Percentage   *big.Int
	Offchain     [32]byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x2667b1cd4b4faaa3decf032654c25acd379781b7ec53b8494b5c8fe87327a547.
//
// Solidity: e __deposit(hid uint256, stationOwner address, value uint256, percentage uint256, offchain bytes32)
func (_CreditATM *CreditATMFilterer) FilterDeposit(opts *bind.FilterOpts) (*CreditATMDepositIterator, error) {

	logs, sub, err := _CreditATM.contract.FilterLogs(opts, "__deposit")
	if err != nil {
		return nil, err
	}
	return &CreditATMDepositIterator{contract: _CreditATM.contract, event: "__deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x2667b1cd4b4faaa3decf032654c25acd379781b7ec53b8494b5c8fe87327a547.
//
// Solidity: e __deposit(hid uint256, stationOwner address, value uint256, percentage uint256, offchain bytes32)
func (_CreditATM *CreditATMFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *CreditATMDeposit) (event.Subscription, error) {

	logs, sub, err := _CreditATM.contract.WatchLogs(opts, "__deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CreditATMDeposit)
				if err := _CreditATM.contract.UnpackLog(event, "__deposit", log); err != nil {
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

// CreditATMReleasePartialFundIterator is returned from FilterReleasePartialFund and is used to iterate over the raw logs and unpacked data for ReleasePartialFund events raised by the CreditATM contract.
type CreditATMReleasePartialFundIterator struct {
	Event *CreditATMReleasePartialFund // Event containing the contract specifics and raw log

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
func (it *CreditATMReleasePartialFundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CreditATMReleasePartialFund)
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
		it.Event = new(CreditATMReleasePartialFund)
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
func (it *CreditATMReleasePartialFundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CreditATMReleasePartialFundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CreditATMReleasePartialFund represents a ReleasePartialFund event raised by the CreditATM contract.
type CreditATMReleasePartialFund struct {
	Hid      *big.Int
	Customer common.Address
	Amount   *big.Int
	Offchain [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReleasePartialFund is a free log retrieval operation binding the contract event 0xa8de681159300a7673a54446f9e99457e22b0ba45217abd287e6847a8e13eb3c.
//
// Solidity: e __releasePartialFund(hid uint256, customer address, amount uint256, offchain bytes32)
func (_CreditATM *CreditATMFilterer) FilterReleasePartialFund(opts *bind.FilterOpts) (*CreditATMReleasePartialFundIterator, error) {

	logs, sub, err := _CreditATM.contract.FilterLogs(opts, "__releasePartialFund")
	if err != nil {
		return nil, err
	}
	return &CreditATMReleasePartialFundIterator{contract: _CreditATM.contract, event: "__releasePartialFund", logs: logs, sub: sub}, nil
}

// WatchReleasePartialFund is a free log subscription operation binding the contract event 0xa8de681159300a7673a54446f9e99457e22b0ba45217abd287e6847a8e13eb3c.
//
// Solidity: e __releasePartialFund(hid uint256, customer address, amount uint256, offchain bytes32)
func (_CreditATM *CreditATMFilterer) WatchReleasePartialFund(opts *bind.WatchOpts, sink chan<- *CreditATMReleasePartialFund) (event.Subscription, error) {

	logs, sub, err := _CreditATM.contract.WatchLogs(opts, "__releasePartialFund")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CreditATMReleasePartialFund)
				if err := _CreditATM.contract.UnpackLog(event, "__releasePartialFund", log); err != nil {
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
