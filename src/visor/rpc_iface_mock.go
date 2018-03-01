/*
* CODE GENERATED AUTOMATICALLY WITH github.com/ernesto-jimenez/goautomock
* THIS FILE MUST NEVER BE EDITED MANUALLY
 */

package visor

import (
	"fmt"

	mock "github.com/stretchr/testify/mock"

	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	blockdb "github.com/skycoin/skycoin/src/visor/blockdb"
	wallet "github.com/skycoin/skycoin/src/wallet"
)

// RPCIfaceMock mock
type RPCIfaceMock struct {
	mock.Mock
}

func NewRPCIfaceMock() *RPCIfaceMock {
	return &RPCIfaceMock{}
}

// CreateAndSignTransaction mocked method
func (m *RPCIfaceMock) CreateAndSignTransaction(p0 string, p1 wallet.Validator, p2 blockdb.UnspentGetter, p3 uint64, p4 uint64, p5 cipher.Address) (*coin.Transaction, error) {

	ret := m.Called(p0, p1, p2, p3, p4, p5)

	var r0 *coin.Transaction
	switch res := ret.Get(0).(type) {
	case nil:
	case *coin.Transaction:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// CreateWallet mocked method
func (m *RPCIfaceMock) CreateWallet(p0 string, p1 wallet.Options) (wallet.Wallet, error) {

	ret := m.Called(p0, p1)

	var r0 wallet.Wallet
	switch res := ret.Get(0).(type) {
	case nil:
	case wallet.Wallet:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// GetAddressTxns mocked method
func (m *RPCIfaceMock) GetAddressTxns(p0 Visorer, p1 cipher.Address) ([]Transaction, error) {

	ret := m.Called(p0, p1)

	var r0 []Transaction
	switch res := ret.Get(0).(type) {
	case nil:
	case []Transaction:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// GetBlock mocked method
func (m *RPCIfaceMock) GetBlock(p0 Visorer, p1 uint64) (*coin.SignedBlock, error) {

	ret := m.Called(p0, p1)

	var r0 *coin.SignedBlock
	switch res := ret.Get(0).(type) {
	case nil:
	case *coin.SignedBlock:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// GetBlockBySeq mocked method
func (m *RPCIfaceMock) GetBlockBySeq(p0 Visorer, p1 uint64) (*coin.SignedBlock, error) {

	ret := m.Called(p0, p1)

	var r0 *coin.SignedBlock
	switch res := ret.Get(0).(type) {
	case nil:
	case *coin.SignedBlock:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// GetBlockchainMetadata mocked method
func (m *RPCIfaceMock) GetBlockchainMetadata(p0 Visorer) *BlockchainMetadata {

	ret := m.Called(p0)

	var r0 *BlockchainMetadata
	switch res := ret.Get(0).(type) {
	case nil:
	case *BlockchainMetadata:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetBlocks mocked method
func (m *RPCIfaceMock) GetBlocks(p0 Visorer, p1 uint64, p2 uint64) []coin.SignedBlock {

	ret := m.Called(p0, p1, p2)

	var r0 []coin.SignedBlock
	switch res := ret.Get(0).(type) {
	case nil:
	case []coin.SignedBlock:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetBuildInfo mocked method
func (m *RPCIfaceMock) GetBuildInfo() BuildInfo {

	ret := m.Called()

	var r0 BuildInfo
	switch res := ret.Get(0).(type) {
	case nil:
	case BuildInfo:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetLastBlocks mocked method
func (m *RPCIfaceMock) GetLastBlocks(p0 Visorer, p1 uint64) []coin.SignedBlock {

	ret := m.Called(p0, p1)

	var r0 []coin.SignedBlock
	switch res := ret.Get(0).(type) {
	case nil:
	case []coin.SignedBlock:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetTransaction mocked method
func (m *RPCIfaceMock) GetTransaction(p0 Visorer, p1 cipher.SHA256) (*Transaction, error) {

	ret := m.Called(p0, p1)

	var r0 *Transaction
	switch res := ret.Get(0).(type) {
	case nil:
	case *Transaction:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// GetUnconfirmedReceiving mocked method
func (m *RPCIfaceMock) GetUnconfirmedReceiving(p0 Visorer, p1 []cipher.Address) (coin.AddressUxOuts, error) {

	ret := m.Called(p0, p1)

	var r0 coin.AddressUxOuts
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.AddressUxOuts:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// GetUnconfirmedSpends mocked method
func (m *RPCIfaceMock) GetUnconfirmedSpends(p0 Visorer, p1 []cipher.Address) (coin.AddressUxOuts, error) {

	ret := m.Called(p0, p1)

	var r0 coin.AddressUxOuts
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.AddressUxOuts:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// GetUnconfirmedTxns mocked method
func (m *RPCIfaceMock) GetUnconfirmedTxns(p0 Visorer, p1 []cipher.Address) []UnconfirmedTxn {

	ret := m.Called(p0, p1)

	var r0 []UnconfirmedTxn
	switch res := ret.Get(0).(type) {
	case nil:
	case []UnconfirmedTxn:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetUnspent mocked method
func (m *RPCIfaceMock) GetUnspent(p0 Visorer) blockdb.UnspentPool {

	ret := m.Called(p0)

	var r0 blockdb.UnspentPool
	switch res := ret.Get(0).(type) {
	case nil:
	case blockdb.UnspentPool:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetWallet mocked method
func (m *RPCIfaceMock) GetWallet(p0 string) (wallet.Wallet, error) {

	ret := m.Called(p0)

	var r0 wallet.Wallet
	switch res := ret.Get(0).(type) {
	case nil:
	case wallet.Wallet:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// GetWalletAddresses mocked method
func (m *RPCIfaceMock) GetWalletAddresses(p0 string) ([]cipher.Address, error) {

	ret := m.Called(p0)

	var r0 []cipher.Address
	switch res := ret.Get(0).(type) {
	case nil:
	case []cipher.Address:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// GetWallets mocked method
func (m *RPCIfaceMock) GetWallets() wallet.Wallets {

	ret := m.Called()

	var r0 wallet.Wallets
	switch res := ret.Get(0).(type) {
	case nil:
	case wallet.Wallets:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// NewAddresses mocked method
func (m *RPCIfaceMock) NewAddresses(p0 string, p1 uint64) ([]cipher.Address, error) {

	ret := m.Called(p0, p1)

	var r0 []cipher.Address
	switch res := ret.Get(0).(type) {
	case nil:
	case []cipher.Address:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// ReloadWallets mocked method
func (m *RPCIfaceMock) ReloadWallets() error {

	ret := m.Called()

	var r0 error
	switch res := ret.Get(0).(type) {
	case nil:
	case error:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// UpdateWalletLabel mocked method
func (m *RPCIfaceMock) UpdateWalletLabel(p0 string, p1 string) error {

	ret := m.Called(p0, p1)

	var r0 error
	switch res := ret.Get(0).(type) {
	case nil:
	case error:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}
