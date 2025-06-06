// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// OutputValidityProof is an auto generated low-level Go binding around an user-defined struct.
type OutputValidityProof struct {
	OutputIndex          uint64
	OutputHashesSiblings [][32]byte
}

// ApplicationMetaData contains all meta data concerning the Application contract.
var ApplicationMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"outputsMerkleRootValidator\",\"type\":\"address\",\"internalType\":\"contractIOutputsMerkleRootValidator\"},{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"templateHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"dataAvailability\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"executeOutput\",\"inputs\":[{\"name\":\"output\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structOutputValidityProof\",\"components\":[{\"name\":\"outputIndex\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"outputHashesSiblings\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getDataAvailability\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDeploymentBlockNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOutputsMerkleRootValidator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIOutputsMerkleRootValidator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTemplateHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrateToOutputsMerkleRootValidator\",\"inputs\":[{\"name\":\"newOutputsMerkleRootValidator\",\"type\":\"address\",\"internalType\":\"contractIOutputsMerkleRootValidator\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onERC1155BatchReceived\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onERC1155Received\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onERC721Received\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validateOutput\",\"inputs\":[{\"name\":\"output\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structOutputValidityProof\",\"components\":[{\"name\":\"outputIndex\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"outputHashesSiblings\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validateOutputHash\",\"inputs\":[{\"name\":\"outputHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"proof\",\"type\":\"tuple\",\"internalType\":\"structOutputValidityProof\",\"components\":[{\"name\":\"outputIndex\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"outputHashesSiblings\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"wasOutputExecuted\",\"inputs\":[{\"name\":\"outputIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"OutputExecuted\",\"inputs\":[{\"name\":\"outputIndex\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"output\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OutputsMerkleRootValidatorChanged\",\"inputs\":[{\"name\":\"newOutputsMerkleRootValidator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIOutputsMerkleRootValidator\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InsufficientFunds\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidOutputHashesSiblingsArrayLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidOutputsMerkleRoot\",\"inputs\":[{\"name\":\"outputsMerkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"OutputNotExecutable\",\"inputs\":[{\"name\":\"output\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"OutputNotReexecutable\",\"inputs\":[{\"name\":\"output\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
}

// ApplicationABI is the input ABI used to generate the binding from.
// Deprecated: Use ApplicationMetaData.ABI instead.
var ApplicationABI = ApplicationMetaData.ABI

// Application is an auto generated Go binding around an Ethereum contract.
type Application struct {
	ApplicationCaller     // Read-only binding to the contract
	ApplicationTransactor // Write-only binding to the contract
	ApplicationFilterer   // Log filterer for contract events
}

// ApplicationCaller is an auto generated read-only Go binding around an Ethereum contract.
type ApplicationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ApplicationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ApplicationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ApplicationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ApplicationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ApplicationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ApplicationSession struct {
	Contract     *Application      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ApplicationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ApplicationCallerSession struct {
	Contract *ApplicationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ApplicationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ApplicationTransactorSession struct {
	Contract     *ApplicationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ApplicationRaw is an auto generated low-level Go binding around an Ethereum contract.
type ApplicationRaw struct {
	Contract *Application // Generic contract binding to access the raw methods on
}

// ApplicationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ApplicationCallerRaw struct {
	Contract *ApplicationCaller // Generic read-only contract binding to access the raw methods on
}

// ApplicationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ApplicationTransactorRaw struct {
	Contract *ApplicationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewApplication creates a new instance of Application, bound to a specific deployed contract.
func NewApplication(address common.Address, backend bind.ContractBackend) (*Application, error) {
	contract, err := bindApplication(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Application{ApplicationCaller: ApplicationCaller{contract: contract}, ApplicationTransactor: ApplicationTransactor{contract: contract}, ApplicationFilterer: ApplicationFilterer{contract: contract}}, nil
}

// NewApplicationCaller creates a new read-only instance of Application, bound to a specific deployed contract.
func NewApplicationCaller(address common.Address, caller bind.ContractCaller) (*ApplicationCaller, error) {
	contract, err := bindApplication(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ApplicationCaller{contract: contract}, nil
}

// NewApplicationTransactor creates a new write-only instance of Application, bound to a specific deployed contract.
func NewApplicationTransactor(address common.Address, transactor bind.ContractTransactor) (*ApplicationTransactor, error) {
	contract, err := bindApplication(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ApplicationTransactor{contract: contract}, nil
}

// NewApplicationFilterer creates a new log filterer instance of Application, bound to a specific deployed contract.
func NewApplicationFilterer(address common.Address, filterer bind.ContractFilterer) (*ApplicationFilterer, error) {
	contract, err := bindApplication(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ApplicationFilterer{contract: contract}, nil
}

// bindApplication binds a generic wrapper to an already deployed contract.
func bindApplication(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ApplicationMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Application *ApplicationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Application.Contract.ApplicationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Application *ApplicationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Application.Contract.ApplicationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Application *ApplicationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Application.Contract.ApplicationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Application *ApplicationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Application.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Application *ApplicationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Application.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Application *ApplicationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Application.Contract.contract.Transact(opts, method, params...)
}

// GetDataAvailability is a free data retrieval call binding the contract method 0xf02478de.
//
// Solidity: function getDataAvailability() view returns(bytes)
func (_Application *ApplicationCaller) GetDataAvailability(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Application.contract.Call(opts, &out, "getDataAvailability")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetDataAvailability is a free data retrieval call binding the contract method 0xf02478de.
//
// Solidity: function getDataAvailability() view returns(bytes)
func (_Application *ApplicationSession) GetDataAvailability() ([]byte, error) {
	return _Application.Contract.GetDataAvailability(&_Application.CallOpts)
}

// GetDataAvailability is a free data retrieval call binding the contract method 0xf02478de.
//
// Solidity: function getDataAvailability() view returns(bytes)
func (_Application *ApplicationCallerSession) GetDataAvailability() ([]byte, error) {
	return _Application.Contract.GetDataAvailability(&_Application.CallOpts)
}

// GetDeploymentBlockNumber is a free data retrieval call binding the contract method 0xb3a1acd8.
//
// Solidity: function getDeploymentBlockNumber() view returns(uint256)
func (_Application *ApplicationCaller) GetDeploymentBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Application.contract.Call(opts, &out, "getDeploymentBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDeploymentBlockNumber is a free data retrieval call binding the contract method 0xb3a1acd8.
//
// Solidity: function getDeploymentBlockNumber() view returns(uint256)
func (_Application *ApplicationSession) GetDeploymentBlockNumber() (*big.Int, error) {
	return _Application.Contract.GetDeploymentBlockNumber(&_Application.CallOpts)
}

// GetDeploymentBlockNumber is a free data retrieval call binding the contract method 0xb3a1acd8.
//
// Solidity: function getDeploymentBlockNumber() view returns(uint256)
func (_Application *ApplicationCallerSession) GetDeploymentBlockNumber() (*big.Int, error) {
	return _Application.Contract.GetDeploymentBlockNumber(&_Application.CallOpts)
}

// GetOutputsMerkleRootValidator is a free data retrieval call binding the contract method 0xa94dfc5a.
//
// Solidity: function getOutputsMerkleRootValidator() view returns(address)
func (_Application *ApplicationCaller) GetOutputsMerkleRootValidator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Application.contract.Call(opts, &out, "getOutputsMerkleRootValidator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOutputsMerkleRootValidator is a free data retrieval call binding the contract method 0xa94dfc5a.
//
// Solidity: function getOutputsMerkleRootValidator() view returns(address)
func (_Application *ApplicationSession) GetOutputsMerkleRootValidator() (common.Address, error) {
	return _Application.Contract.GetOutputsMerkleRootValidator(&_Application.CallOpts)
}

// GetOutputsMerkleRootValidator is a free data retrieval call binding the contract method 0xa94dfc5a.
//
// Solidity: function getOutputsMerkleRootValidator() view returns(address)
func (_Application *ApplicationCallerSession) GetOutputsMerkleRootValidator() (common.Address, error) {
	return _Application.Contract.GetOutputsMerkleRootValidator(&_Application.CallOpts)
}

// GetTemplateHash is a free data retrieval call binding the contract method 0x61b12c66.
//
// Solidity: function getTemplateHash() view returns(bytes32)
func (_Application *ApplicationCaller) GetTemplateHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Application.contract.Call(opts, &out, "getTemplateHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetTemplateHash is a free data retrieval call binding the contract method 0x61b12c66.
//
// Solidity: function getTemplateHash() view returns(bytes32)
func (_Application *ApplicationSession) GetTemplateHash() ([32]byte, error) {
	return _Application.Contract.GetTemplateHash(&_Application.CallOpts)
}

// GetTemplateHash is a free data retrieval call binding the contract method 0x61b12c66.
//
// Solidity: function getTemplateHash() view returns(bytes32)
func (_Application *ApplicationCallerSession) GetTemplateHash() ([32]byte, error) {
	return _Application.Contract.GetTemplateHash(&_Application.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Application *ApplicationCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Application.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Application *ApplicationSession) Owner() (common.Address, error) {
	return _Application.Contract.Owner(&_Application.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Application *ApplicationCallerSession) Owner() (common.Address, error) {
	return _Application.Contract.Owner(&_Application.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Application *ApplicationCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Application.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Application *ApplicationSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Application.Contract.SupportsInterface(&_Application.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Application *ApplicationCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Application.Contract.SupportsInterface(&_Application.CallOpts, interfaceId)
}

// ValidateOutput is a free data retrieval call binding the contract method 0xe88d39c0.
//
// Solidity: function validateOutput(bytes output, (uint64,bytes32[]) proof) view returns()
func (_Application *ApplicationCaller) ValidateOutput(opts *bind.CallOpts, output []byte, proof OutputValidityProof) error {
	var out []interface{}
	err := _Application.contract.Call(opts, &out, "validateOutput", output, proof)

	if err != nil {
		return err
	}

	return err

}

// ValidateOutput is a free data retrieval call binding the contract method 0xe88d39c0.
//
// Solidity: function validateOutput(bytes output, (uint64,bytes32[]) proof) view returns()
func (_Application *ApplicationSession) ValidateOutput(output []byte, proof OutputValidityProof) error {
	return _Application.Contract.ValidateOutput(&_Application.CallOpts, output, proof)
}

// ValidateOutput is a free data retrieval call binding the contract method 0xe88d39c0.
//
// Solidity: function validateOutput(bytes output, (uint64,bytes32[]) proof) view returns()
func (_Application *ApplicationCallerSession) ValidateOutput(output []byte, proof OutputValidityProof) error {
	return _Application.Contract.ValidateOutput(&_Application.CallOpts, output, proof)
}

// ValidateOutputHash is a free data retrieval call binding the contract method 0x08eb89ab.
//
// Solidity: function validateOutputHash(bytes32 outputHash, (uint64,bytes32[]) proof) view returns()
func (_Application *ApplicationCaller) ValidateOutputHash(opts *bind.CallOpts, outputHash [32]byte, proof OutputValidityProof) error {
	var out []interface{}
	err := _Application.contract.Call(opts, &out, "validateOutputHash", outputHash, proof)

	if err != nil {
		return err
	}

	return err

}

// ValidateOutputHash is a free data retrieval call binding the contract method 0x08eb89ab.
//
// Solidity: function validateOutputHash(bytes32 outputHash, (uint64,bytes32[]) proof) view returns()
func (_Application *ApplicationSession) ValidateOutputHash(outputHash [32]byte, proof OutputValidityProof) error {
	return _Application.Contract.ValidateOutputHash(&_Application.CallOpts, outputHash, proof)
}

// ValidateOutputHash is a free data retrieval call binding the contract method 0x08eb89ab.
//
// Solidity: function validateOutputHash(bytes32 outputHash, (uint64,bytes32[]) proof) view returns()
func (_Application *ApplicationCallerSession) ValidateOutputHash(outputHash [32]byte, proof OutputValidityProof) error {
	return _Application.Contract.ValidateOutputHash(&_Application.CallOpts, outputHash, proof)
}

// WasOutputExecuted is a free data retrieval call binding the contract method 0x71891db0.
//
// Solidity: function wasOutputExecuted(uint256 outputIndex) view returns(bool)
func (_Application *ApplicationCaller) WasOutputExecuted(opts *bind.CallOpts, outputIndex *big.Int) (bool, error) {
	var out []interface{}
	err := _Application.contract.Call(opts, &out, "wasOutputExecuted", outputIndex)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// WasOutputExecuted is a free data retrieval call binding the contract method 0x71891db0.
//
// Solidity: function wasOutputExecuted(uint256 outputIndex) view returns(bool)
func (_Application *ApplicationSession) WasOutputExecuted(outputIndex *big.Int) (bool, error) {
	return _Application.Contract.WasOutputExecuted(&_Application.CallOpts, outputIndex)
}

// WasOutputExecuted is a free data retrieval call binding the contract method 0x71891db0.
//
// Solidity: function wasOutputExecuted(uint256 outputIndex) view returns(bool)
func (_Application *ApplicationCallerSession) WasOutputExecuted(outputIndex *big.Int) (bool, error) {
	return _Application.Contract.WasOutputExecuted(&_Application.CallOpts, outputIndex)
}

// ExecuteOutput is a paid mutator transaction binding the contract method 0x33137b76.
//
// Solidity: function executeOutput(bytes output, (uint64,bytes32[]) proof) returns()
func (_Application *ApplicationTransactor) ExecuteOutput(opts *bind.TransactOpts, output []byte, proof OutputValidityProof) (*types.Transaction, error) {
	return _Application.contract.Transact(opts, "executeOutput", output, proof)
}

// ExecuteOutput is a paid mutator transaction binding the contract method 0x33137b76.
//
// Solidity: function executeOutput(bytes output, (uint64,bytes32[]) proof) returns()
func (_Application *ApplicationSession) ExecuteOutput(output []byte, proof OutputValidityProof) (*types.Transaction, error) {
	return _Application.Contract.ExecuteOutput(&_Application.TransactOpts, output, proof)
}

// ExecuteOutput is a paid mutator transaction binding the contract method 0x33137b76.
//
// Solidity: function executeOutput(bytes output, (uint64,bytes32[]) proof) returns()
func (_Application *ApplicationTransactorSession) ExecuteOutput(output []byte, proof OutputValidityProof) (*types.Transaction, error) {
	return _Application.Contract.ExecuteOutput(&_Application.TransactOpts, output, proof)
}

// MigrateToOutputsMerkleRootValidator is a paid mutator transaction binding the contract method 0xbf8abff8.
//
// Solidity: function migrateToOutputsMerkleRootValidator(address newOutputsMerkleRootValidator) returns()
func (_Application *ApplicationTransactor) MigrateToOutputsMerkleRootValidator(opts *bind.TransactOpts, newOutputsMerkleRootValidator common.Address) (*types.Transaction, error) {
	return _Application.contract.Transact(opts, "migrateToOutputsMerkleRootValidator", newOutputsMerkleRootValidator)
}

// MigrateToOutputsMerkleRootValidator is a paid mutator transaction binding the contract method 0xbf8abff8.
//
// Solidity: function migrateToOutputsMerkleRootValidator(address newOutputsMerkleRootValidator) returns()
func (_Application *ApplicationSession) MigrateToOutputsMerkleRootValidator(newOutputsMerkleRootValidator common.Address) (*types.Transaction, error) {
	return _Application.Contract.MigrateToOutputsMerkleRootValidator(&_Application.TransactOpts, newOutputsMerkleRootValidator)
}

// MigrateToOutputsMerkleRootValidator is a paid mutator transaction binding the contract method 0xbf8abff8.
//
// Solidity: function migrateToOutputsMerkleRootValidator(address newOutputsMerkleRootValidator) returns()
func (_Application *ApplicationTransactorSession) MigrateToOutputsMerkleRootValidator(newOutputsMerkleRootValidator common.Address) (*types.Transaction, error) {
	return _Application.Contract.MigrateToOutputsMerkleRootValidator(&_Application.TransactOpts, newOutputsMerkleRootValidator)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_Application *ApplicationTransactor) OnERC1155BatchReceived(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _Application.contract.Transact(opts, "onERC1155BatchReceived", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_Application *ApplicationSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _Application.Contract.OnERC1155BatchReceived(&_Application.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_Application *ApplicationTransactorSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _Application.Contract.OnERC1155BatchReceived(&_Application.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_Application *ApplicationTransactor) OnERC1155Received(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _Application.contract.Transact(opts, "onERC1155Received", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_Application *ApplicationSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _Application.Contract.OnERC1155Received(&_Application.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_Application *ApplicationTransactorSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _Application.Contract.OnERC1155Received(&_Application.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_Application *ApplicationTransactor) OnERC721Received(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _Application.contract.Transact(opts, "onERC721Received", arg0, arg1, arg2, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_Application *ApplicationSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _Application.Contract.OnERC721Received(&_Application.TransactOpts, arg0, arg1, arg2, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_Application *ApplicationTransactorSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _Application.Contract.OnERC721Received(&_Application.TransactOpts, arg0, arg1, arg2, arg3)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Application *ApplicationTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Application.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Application *ApplicationSession) RenounceOwnership() (*types.Transaction, error) {
	return _Application.Contract.RenounceOwnership(&_Application.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Application *ApplicationTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Application.Contract.RenounceOwnership(&_Application.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Application *ApplicationTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Application.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Application *ApplicationSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Application.Contract.TransferOwnership(&_Application.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Application *ApplicationTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Application.Contract.TransferOwnership(&_Application.TransactOpts, newOwner)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Application *ApplicationTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Application.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Application *ApplicationSession) Receive() (*types.Transaction, error) {
	return _Application.Contract.Receive(&_Application.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Application *ApplicationTransactorSession) Receive() (*types.Transaction, error) {
	return _Application.Contract.Receive(&_Application.TransactOpts)
}

// ApplicationOutputExecutedIterator is returned from FilterOutputExecuted and is used to iterate over the raw logs and unpacked data for OutputExecuted events raised by the Application contract.
type ApplicationOutputExecutedIterator struct {
	Event *ApplicationOutputExecuted // Event containing the contract specifics and raw log

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
func (it *ApplicationOutputExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ApplicationOutputExecuted)
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
		it.Event = new(ApplicationOutputExecuted)
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
func (it *ApplicationOutputExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ApplicationOutputExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ApplicationOutputExecuted represents a OutputExecuted event raised by the Application contract.
type ApplicationOutputExecuted struct {
	OutputIndex uint64
	Output      []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOutputExecuted is a free log retrieval operation binding the contract event 0xcad1f361c6e84664e892230291c8e8eb9555683e0a6a5ce8ea7b204ac0ac3676.
//
// Solidity: event OutputExecuted(uint64 outputIndex, bytes output)
func (_Application *ApplicationFilterer) FilterOutputExecuted(opts *bind.FilterOpts) (*ApplicationOutputExecutedIterator, error) {

	logs, sub, err := _Application.contract.FilterLogs(opts, "OutputExecuted")
	if err != nil {
		return nil, err
	}
	return &ApplicationOutputExecutedIterator{contract: _Application.contract, event: "OutputExecuted", logs: logs, sub: sub}, nil
}

// WatchOutputExecuted is a free log subscription operation binding the contract event 0xcad1f361c6e84664e892230291c8e8eb9555683e0a6a5ce8ea7b204ac0ac3676.
//
// Solidity: event OutputExecuted(uint64 outputIndex, bytes output)
func (_Application *ApplicationFilterer) WatchOutputExecuted(opts *bind.WatchOpts, sink chan<- *ApplicationOutputExecuted) (event.Subscription, error) {

	logs, sub, err := _Application.contract.WatchLogs(opts, "OutputExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ApplicationOutputExecuted)
				if err := _Application.contract.UnpackLog(event, "OutputExecuted", log); err != nil {
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

// ParseOutputExecuted is a log parse operation binding the contract event 0xcad1f361c6e84664e892230291c8e8eb9555683e0a6a5ce8ea7b204ac0ac3676.
//
// Solidity: event OutputExecuted(uint64 outputIndex, bytes output)
func (_Application *ApplicationFilterer) ParseOutputExecuted(log types.Log) (*ApplicationOutputExecuted, error) {
	event := new(ApplicationOutputExecuted)
	if err := _Application.contract.UnpackLog(event, "OutputExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ApplicationOutputsMerkleRootValidatorChangedIterator is returned from FilterOutputsMerkleRootValidatorChanged and is used to iterate over the raw logs and unpacked data for OutputsMerkleRootValidatorChanged events raised by the Application contract.
type ApplicationOutputsMerkleRootValidatorChangedIterator struct {
	Event *ApplicationOutputsMerkleRootValidatorChanged // Event containing the contract specifics and raw log

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
func (it *ApplicationOutputsMerkleRootValidatorChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ApplicationOutputsMerkleRootValidatorChanged)
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
		it.Event = new(ApplicationOutputsMerkleRootValidatorChanged)
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
func (it *ApplicationOutputsMerkleRootValidatorChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ApplicationOutputsMerkleRootValidatorChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ApplicationOutputsMerkleRootValidatorChanged represents a OutputsMerkleRootValidatorChanged event raised by the Application contract.
type ApplicationOutputsMerkleRootValidatorChanged struct {
	NewOutputsMerkleRootValidator common.Address
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterOutputsMerkleRootValidatorChanged is a free log retrieval operation binding the contract event 0x6ad3188ba8f430fba0656cb0a7e839ab2020d5586ba11a1477d18f7092f8bece.
//
// Solidity: event OutputsMerkleRootValidatorChanged(address newOutputsMerkleRootValidator)
func (_Application *ApplicationFilterer) FilterOutputsMerkleRootValidatorChanged(opts *bind.FilterOpts) (*ApplicationOutputsMerkleRootValidatorChangedIterator, error) {

	logs, sub, err := _Application.contract.FilterLogs(opts, "OutputsMerkleRootValidatorChanged")
	if err != nil {
		return nil, err
	}
	return &ApplicationOutputsMerkleRootValidatorChangedIterator{contract: _Application.contract, event: "OutputsMerkleRootValidatorChanged", logs: logs, sub: sub}, nil
}

// WatchOutputsMerkleRootValidatorChanged is a free log subscription operation binding the contract event 0x6ad3188ba8f430fba0656cb0a7e839ab2020d5586ba11a1477d18f7092f8bece.
//
// Solidity: event OutputsMerkleRootValidatorChanged(address newOutputsMerkleRootValidator)
func (_Application *ApplicationFilterer) WatchOutputsMerkleRootValidatorChanged(opts *bind.WatchOpts, sink chan<- *ApplicationOutputsMerkleRootValidatorChanged) (event.Subscription, error) {

	logs, sub, err := _Application.contract.WatchLogs(opts, "OutputsMerkleRootValidatorChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ApplicationOutputsMerkleRootValidatorChanged)
				if err := _Application.contract.UnpackLog(event, "OutputsMerkleRootValidatorChanged", log); err != nil {
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

// ParseOutputsMerkleRootValidatorChanged is a log parse operation binding the contract event 0x6ad3188ba8f430fba0656cb0a7e839ab2020d5586ba11a1477d18f7092f8bece.
//
// Solidity: event OutputsMerkleRootValidatorChanged(address newOutputsMerkleRootValidator)
func (_Application *ApplicationFilterer) ParseOutputsMerkleRootValidatorChanged(log types.Log) (*ApplicationOutputsMerkleRootValidatorChanged, error) {
	event := new(ApplicationOutputsMerkleRootValidatorChanged)
	if err := _Application.contract.UnpackLog(event, "OutputsMerkleRootValidatorChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ApplicationOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Application contract.
type ApplicationOwnershipTransferredIterator struct {
	Event *ApplicationOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ApplicationOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ApplicationOwnershipTransferred)
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
		it.Event = new(ApplicationOwnershipTransferred)
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
func (it *ApplicationOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ApplicationOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ApplicationOwnershipTransferred represents a OwnershipTransferred event raised by the Application contract.
type ApplicationOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Application *ApplicationFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ApplicationOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Application.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ApplicationOwnershipTransferredIterator{contract: _Application.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Application *ApplicationFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ApplicationOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Application.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ApplicationOwnershipTransferred)
				if err := _Application.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Application *ApplicationFilterer) ParseOwnershipTransferred(log types.Log) (*ApplicationOwnershipTransferred, error) {
	event := new(ApplicationOwnershipTransferred)
	if err := _Application.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
