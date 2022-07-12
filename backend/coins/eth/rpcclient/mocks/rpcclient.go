// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/eth/erc20"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/eth/rpcclient"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"sync"
)

// Ensure, that InterfaceMock does implement rpcclient.Interface.
// If this is not the case, regenerate this file with moq.
var _ rpcclient.Interface = &InterfaceMock{}

// InterfaceMock is a mock implementation of rpcclient.Interface.
//
// 	func TestSomethingThatUsesInterface(t *testing.T) {
//
// 		// make and configure a mocked rpcclient.Interface
// 		mockedInterface := &InterfaceMock{
// 			BalanceFunc: func(ctx context.Context, account common.Address) (*big.Int, error) {
// 				panic("mock out the Balance method")
// 			},
// 			BlockNumberFunc: func(ctx context.Context) (*big.Int, error) {
// 				panic("mock out the BlockNumber method")
// 			},
// 			ERC20BalanceFunc: func(account common.Address, erc20Token *erc20.Token) (*big.Int, error) {
// 				panic("mock out the ERC20Balance method")
// 			},
// 			EstimateGasFunc: func(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
// 				panic("mock out the EstimateGas method")
// 			},
// 			PendingNonceAtFunc: func(ctx context.Context, account common.Address) (uint64, error) {
// 				panic("mock out the PendingNonceAt method")
// 			},
// 			SendTransactionFunc: func(ctx context.Context, tx *types.Transaction) error {
// 				panic("mock out the SendTransaction method")
// 			},
// 			SuggestGasPriceFunc: func(ctx context.Context) (*big.Int, error) {
// 				panic("mock out the SuggestGasPrice method")
// 			},
// 			TransactionByHashFunc: func(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
// 				panic("mock out the TransactionByHash method")
// 			},
// 			TransactionReceiptWithBlockNumberFunc: func(ctx context.Context, hash common.Hash) (*rpcclient.RPCTransactionReceipt, error) {
// 				panic("mock out the TransactionReceiptWithBlockNumber method")
// 			},
// 		}
//
// 		// use mockedInterface in code that requires rpcclient.Interface
// 		// and then make assertions.
//
// 	}
type InterfaceMock struct {
	// BalanceFunc mocks the Balance method.
	BalanceFunc func(ctx context.Context, account common.Address) (*big.Int, error)

	// BlockNumberFunc mocks the BlockNumber method.
	BlockNumberFunc func(ctx context.Context) (*big.Int, error)

	// ERC20BalanceFunc mocks the ERC20Balance method.
	ERC20BalanceFunc func(account common.Address, erc20Token *erc20.Token) (*big.Int, error)

	// EstimateGasFunc mocks the EstimateGas method.
	EstimateGasFunc func(ctx context.Context, call ethereum.CallMsg) (uint64, error)

	// PendingNonceAtFunc mocks the PendingNonceAt method.
	PendingNonceAtFunc func(ctx context.Context, account common.Address) (uint64, error)

	// SendTransactionFunc mocks the SendTransaction method.
	SendTransactionFunc func(ctx context.Context, tx *types.Transaction) error

	// SuggestGasPriceFunc mocks the SuggestGasPrice method.
	SuggestGasPriceFunc func(ctx context.Context) (*big.Int, error)

	// TransactionByHashFunc mocks the TransactionByHash method.
	TransactionByHashFunc func(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)

	// TransactionReceiptWithBlockNumberFunc mocks the TransactionReceiptWithBlockNumber method.
	TransactionReceiptWithBlockNumberFunc func(ctx context.Context, hash common.Hash) (*rpcclient.RPCTransactionReceipt, error)

	// calls tracks calls to the methods.
	calls struct {
		// Balance holds details about calls to the Balance method.
		Balance []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Account is the account argument value.
			Account common.Address
		}
		// BlockNumber holds details about calls to the BlockNumber method.
		BlockNumber []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// ERC20Balance holds details about calls to the ERC20Balance method.
		ERC20Balance []struct {
			// Account is the account argument value.
			Account common.Address
			// Erc20Token is the erc20Token argument value.
			Erc20Token *erc20.Token
		}
		// EstimateGas holds details about calls to the EstimateGas method.
		EstimateGas []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Call is the call argument value.
			Call ethereum.CallMsg
		}
		// PendingNonceAt holds details about calls to the PendingNonceAt method.
		PendingNonceAt []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Account is the account argument value.
			Account common.Address
		}
		// SendTransaction holds details about calls to the SendTransaction method.
		SendTransaction []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Tx is the tx argument value.
			Tx *types.Transaction
		}
		// SuggestGasPrice holds details about calls to the SuggestGasPrice method.
		SuggestGasPrice []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// TransactionByHash holds details about calls to the TransactionByHash method.
		TransactionByHash []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Hash is the hash argument value.
			Hash common.Hash
		}
		// TransactionReceiptWithBlockNumber holds details about calls to the TransactionReceiptWithBlockNumber method.
		TransactionReceiptWithBlockNumber []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Hash is the hash argument value.
			Hash common.Hash
		}
	}
	lockBalance                           sync.RWMutex
	lockBlockNumber                       sync.RWMutex
	lockERC20Balance                      sync.RWMutex
	lockEstimateGas                       sync.RWMutex
	lockPendingNonceAt                    sync.RWMutex
	lockSendTransaction                   sync.RWMutex
	lockSuggestGasPrice                   sync.RWMutex
	lockTransactionByHash                 sync.RWMutex
	lockTransactionReceiptWithBlockNumber sync.RWMutex
}

// Balance calls BalanceFunc.
func (mock *InterfaceMock) Balance(ctx context.Context, account common.Address) (*big.Int, error) {
	if mock.BalanceFunc == nil {
		panic("InterfaceMock.BalanceFunc: method is nil but Interface.Balance was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Account common.Address
	}{
		Ctx:     ctx,
		Account: account,
	}
	mock.lockBalance.Lock()
	mock.calls.Balance = append(mock.calls.Balance, callInfo)
	mock.lockBalance.Unlock()
	return mock.BalanceFunc(ctx, account)
}

// BalanceCalls gets all the calls that were made to Balance.
// Check the length with:
//     len(mockedInterface.BalanceCalls())
func (mock *InterfaceMock) BalanceCalls() []struct {
	Ctx     context.Context
	Account common.Address
} {
	var calls []struct {
		Ctx     context.Context
		Account common.Address
	}
	mock.lockBalance.RLock()
	calls = mock.calls.Balance
	mock.lockBalance.RUnlock()
	return calls
}

// BlockNumber calls BlockNumberFunc.
func (mock *InterfaceMock) BlockNumber(ctx context.Context) (*big.Int, error) {
	if mock.BlockNumberFunc == nil {
		panic("InterfaceMock.BlockNumberFunc: method is nil but Interface.BlockNumber was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockBlockNumber.Lock()
	mock.calls.BlockNumber = append(mock.calls.BlockNumber, callInfo)
	mock.lockBlockNumber.Unlock()
	return mock.BlockNumberFunc(ctx)
}

// BlockNumberCalls gets all the calls that were made to BlockNumber.
// Check the length with:
//     len(mockedInterface.BlockNumberCalls())
func (mock *InterfaceMock) BlockNumberCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockBlockNumber.RLock()
	calls = mock.calls.BlockNumber
	mock.lockBlockNumber.RUnlock()
	return calls
}

// ERC20Balance calls ERC20BalanceFunc.
func (mock *InterfaceMock) ERC20Balance(account common.Address, erc20Token *erc20.Token) (*big.Int, error) {
	if mock.ERC20BalanceFunc == nil {
		panic("InterfaceMock.ERC20BalanceFunc: method is nil but Interface.ERC20Balance was just called")
	}
	callInfo := struct {
		Account    common.Address
		Erc20Token *erc20.Token
	}{
		Account:    account,
		Erc20Token: erc20Token,
	}
	mock.lockERC20Balance.Lock()
	mock.calls.ERC20Balance = append(mock.calls.ERC20Balance, callInfo)
	mock.lockERC20Balance.Unlock()
	return mock.ERC20BalanceFunc(account, erc20Token)
}

// ERC20BalanceCalls gets all the calls that were made to ERC20Balance.
// Check the length with:
//     len(mockedInterface.ERC20BalanceCalls())
func (mock *InterfaceMock) ERC20BalanceCalls() []struct {
	Account    common.Address
	Erc20Token *erc20.Token
} {
	var calls []struct {
		Account    common.Address
		Erc20Token *erc20.Token
	}
	mock.lockERC20Balance.RLock()
	calls = mock.calls.ERC20Balance
	mock.lockERC20Balance.RUnlock()
	return calls
}

// EstimateGas calls EstimateGasFunc.
func (mock *InterfaceMock) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	if mock.EstimateGasFunc == nil {
		panic("InterfaceMock.EstimateGasFunc: method is nil but Interface.EstimateGas was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Call ethereum.CallMsg
	}{
		Ctx:  ctx,
		Call: call,
	}
	mock.lockEstimateGas.Lock()
	mock.calls.EstimateGas = append(mock.calls.EstimateGas, callInfo)
	mock.lockEstimateGas.Unlock()
	return mock.EstimateGasFunc(ctx, call)
}

// EstimateGasCalls gets all the calls that were made to EstimateGas.
// Check the length with:
//     len(mockedInterface.EstimateGasCalls())
func (mock *InterfaceMock) EstimateGasCalls() []struct {
	Ctx  context.Context
	Call ethereum.CallMsg
} {
	var calls []struct {
		Ctx  context.Context
		Call ethereum.CallMsg
	}
	mock.lockEstimateGas.RLock()
	calls = mock.calls.EstimateGas
	mock.lockEstimateGas.RUnlock()
	return calls
}

// PendingNonceAt calls PendingNonceAtFunc.
func (mock *InterfaceMock) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if mock.PendingNonceAtFunc == nil {
		panic("InterfaceMock.PendingNonceAtFunc: method is nil but Interface.PendingNonceAt was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Account common.Address
	}{
		Ctx:     ctx,
		Account: account,
	}
	mock.lockPendingNonceAt.Lock()
	mock.calls.PendingNonceAt = append(mock.calls.PendingNonceAt, callInfo)
	mock.lockPendingNonceAt.Unlock()
	return mock.PendingNonceAtFunc(ctx, account)
}

// PendingNonceAtCalls gets all the calls that were made to PendingNonceAt.
// Check the length with:
//     len(mockedInterface.PendingNonceAtCalls())
func (mock *InterfaceMock) PendingNonceAtCalls() []struct {
	Ctx     context.Context
	Account common.Address
} {
	var calls []struct {
		Ctx     context.Context
		Account common.Address
	}
	mock.lockPendingNonceAt.RLock()
	calls = mock.calls.PendingNonceAt
	mock.lockPendingNonceAt.RUnlock()
	return calls
}

// SendTransaction calls SendTransactionFunc.
func (mock *InterfaceMock) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if mock.SendTransactionFunc == nil {
		panic("InterfaceMock.SendTransactionFunc: method is nil but Interface.SendTransaction was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Tx  *types.Transaction
	}{
		Ctx: ctx,
		Tx:  tx,
	}
	mock.lockSendTransaction.Lock()
	mock.calls.SendTransaction = append(mock.calls.SendTransaction, callInfo)
	mock.lockSendTransaction.Unlock()
	return mock.SendTransactionFunc(ctx, tx)
}

// SendTransactionCalls gets all the calls that were made to SendTransaction.
// Check the length with:
//     len(mockedInterface.SendTransactionCalls())
func (mock *InterfaceMock) SendTransactionCalls() []struct {
	Ctx context.Context
	Tx  *types.Transaction
} {
	var calls []struct {
		Ctx context.Context
		Tx  *types.Transaction
	}
	mock.lockSendTransaction.RLock()
	calls = mock.calls.SendTransaction
	mock.lockSendTransaction.RUnlock()
	return calls
}

// SuggestGasPrice calls SuggestGasPriceFunc.
func (mock *InterfaceMock) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if mock.SuggestGasPriceFunc == nil {
		panic("InterfaceMock.SuggestGasPriceFunc: method is nil but Interface.SuggestGasPrice was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockSuggestGasPrice.Lock()
	mock.calls.SuggestGasPrice = append(mock.calls.SuggestGasPrice, callInfo)
	mock.lockSuggestGasPrice.Unlock()
	return mock.SuggestGasPriceFunc(ctx)
}

// SuggestGasPriceCalls gets all the calls that were made to SuggestGasPrice.
// Check the length with:
//     len(mockedInterface.SuggestGasPriceCalls())
func (mock *InterfaceMock) SuggestGasPriceCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockSuggestGasPrice.RLock()
	calls = mock.calls.SuggestGasPrice
	mock.lockSuggestGasPrice.RUnlock()
	return calls
}

// TransactionByHash calls TransactionByHashFunc.
func (mock *InterfaceMock) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	if mock.TransactionByHashFunc == nil {
		panic("InterfaceMock.TransactionByHashFunc: method is nil but Interface.TransactionByHash was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Hash common.Hash
	}{
		Ctx:  ctx,
		Hash: hash,
	}
	mock.lockTransactionByHash.Lock()
	mock.calls.TransactionByHash = append(mock.calls.TransactionByHash, callInfo)
	mock.lockTransactionByHash.Unlock()
	return mock.TransactionByHashFunc(ctx, hash)
}

// TransactionByHashCalls gets all the calls that were made to TransactionByHash.
// Check the length with:
//     len(mockedInterface.TransactionByHashCalls())
func (mock *InterfaceMock) TransactionByHashCalls() []struct {
	Ctx  context.Context
	Hash common.Hash
} {
	var calls []struct {
		Ctx  context.Context
		Hash common.Hash
	}
	mock.lockTransactionByHash.RLock()
	calls = mock.calls.TransactionByHash
	mock.lockTransactionByHash.RUnlock()
	return calls
}

// TransactionReceiptWithBlockNumber calls TransactionReceiptWithBlockNumberFunc.
func (mock *InterfaceMock) TransactionReceiptWithBlockNumber(ctx context.Context, hash common.Hash) (*rpcclient.RPCTransactionReceipt, error) {
	if mock.TransactionReceiptWithBlockNumberFunc == nil {
		panic("InterfaceMock.TransactionReceiptWithBlockNumberFunc: method is nil but Interface.TransactionReceiptWithBlockNumber was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Hash common.Hash
	}{
		Ctx:  ctx,
		Hash: hash,
	}
	mock.lockTransactionReceiptWithBlockNumber.Lock()
	mock.calls.TransactionReceiptWithBlockNumber = append(mock.calls.TransactionReceiptWithBlockNumber, callInfo)
	mock.lockTransactionReceiptWithBlockNumber.Unlock()
	return mock.TransactionReceiptWithBlockNumberFunc(ctx, hash)
}

// TransactionReceiptWithBlockNumberCalls gets all the calls that were made to TransactionReceiptWithBlockNumber.
// Check the length with:
//     len(mockedInterface.TransactionReceiptWithBlockNumberCalls())
func (mock *InterfaceMock) TransactionReceiptWithBlockNumberCalls() []struct {
	Ctx  context.Context
	Hash common.Hash
} {
	var calls []struct {
		Ctx  context.Context
		Hash common.Hash
	}
	mock.lockTransactionReceiptWithBlockNumber.RLock()
	calls = mock.calls.TransactionReceiptWithBlockNumber
	mock.lockTransactionReceiptWithBlockNumber.RUnlock()
	return calls
}
