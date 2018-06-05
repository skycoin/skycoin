/*
* CODE GENERATED AUTOMATICALLY WITH github.com/ernesto-jimenez/goautomock
* THIS FILE MUST NEVER BE EDITED MANUALLY
 */

package api

import (
	"fmt"

	mock "github.com/stretchr/testify/mock"

	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	daemon "github.com/skycoin/skycoin/src/daemon"
	visor "github.com/skycoin/skycoin/src/visor"
	historydb "github.com/skycoin/skycoin/src/visor/historydb"
	wallet "github.com/skycoin/skycoin/src/wallet"
)

// GatewayerMock mock
type GatewayerMock struct {
	mock.Mock
}

func NewGatewayerMock() *GatewayerMock {
	return &GatewayerMock{}
}

// CreateTransaction mocked method
func (m *GatewayerMock) CreateTransaction(p0 wallet.CreateTransactionParams) (*coin.Transaction, []wallet.UxBalance, error) {

	ret := m.Called(p0)

	var r0 *coin.Transaction
	switch res := ret.Get(0).(type) {
	case nil:
	case *coin.Transaction:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 []wallet.UxBalance
	switch res := ret.Get(1).(type) {
	case nil:
	case []wallet.UxBalance:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r2 error
	switch res := ret.Get(2).(type) {
	case nil:
	case error:
		r2 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1, r2

}

// CreateWallet mocked method
func (m *GatewayerMock) CreateWallet(p0 string, p1 wallet.Options) (*wallet.Wallet, error) {

	ret := m.Called(p0, p1)

	var r0 *wallet.Wallet
	switch res := ret.Get(0).(type) {
	case nil:
	case *wallet.Wallet:
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

// DecryptWallet mocked method
func (m *GatewayerMock) DecryptWallet(p0 string, p1 []byte) (*wallet.Wallet, error) {

	ret := m.Called(p0, p1)

	var r0 *wallet.Wallet
	switch res := ret.Get(0).(type) {
	case nil:
	case *wallet.Wallet:
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

// EncryptWallet mocked method
func (m *GatewayerMock) EncryptWallet(p0 string, p1 []byte) (*wallet.Wallet, error) {

	ret := m.Called(p0, p1)

	var r0 *wallet.Wallet
	switch res := ret.Get(0).(type) {
	case nil:
	case *wallet.Wallet:
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

// GetAddrUxOuts mocked method
func (m *GatewayerMock) GetAddrUxOuts(p0 []cipher.Address) ([]*historydb.UxOut, error) {

	ret := m.Called(p0)

	var r0 []*historydb.UxOut
	switch res := ret.Get(0).(type) {
	case nil:
	case []*historydb.UxOut:
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

// GetAddressCount mocked method
func (m *GatewayerMock) GetAddressCount() (uint64, error) {

	ret := m.Called()

	var r0 uint64
	switch res := ret.Get(0).(type) {
	case nil:
	case uint64:
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

// GetAllUnconfirmedTxns mocked method
func (m *GatewayerMock) GetAllUnconfirmedTxns() ([]visor.UnconfirmedTxn, error) {

	ret := m.Called()

	var r0 []visor.UnconfirmedTxn
	switch res := ret.Get(0).(type) {
	case nil:
	case []visor.UnconfirmedTxn:
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

// GetBalanceOfAddrs mocked method
func (m *GatewayerMock) GetBalanceOfAddrs(p0 []cipher.Address) ([]wallet.BalancePair, error) {

	ret := m.Called(p0)

	var r0 []wallet.BalancePair
	switch res := ret.Get(0).(type) {
	case nil:
	case []wallet.BalancePair:
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
func (m *GatewayerMock) GetBlockchainMetadata() (*visor.BlockchainMetadata, error) {

	ret := m.Called()

	var r0 *visor.BlockchainMetadata
	switch res := ret.Get(0).(type) {
	case nil:
	case *visor.BlockchainMetadata:
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

// GetBlockchainProgress mocked method
func (m *GatewayerMock) GetBlockchainProgress() (*daemon.BlockchainProgress, error) {

	ret := m.Called()

	var r0 *daemon.BlockchainProgress
	switch res := ret.Get(0).(type) {
	case nil:
	case *daemon.BlockchainProgress:
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

// GetBlocks mocked method
func (m *GatewayerMock) GetBlocks(p0 uint64, p1 uint64) (*visor.ReadableBlocks, error) {

	ret := m.Called(p0, p1)

	var r0 *visor.ReadableBlocks
	switch res := ret.Get(0).(type) {
	case nil:
	case *visor.ReadableBlocks:
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

// GetBuildInfo mocked method
func (m *GatewayerMock) GetBuildInfo() visor.BuildInfo {

	ret := m.Called()

	var r0 visor.BuildInfo
	switch res := ret.Get(0).(type) {
	case nil:
	case visor.BuildInfo:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetConnection mocked method
func (m *GatewayerMock) GetConnection(p0 string) *daemon.Connection {

	ret := m.Called(p0)

	var r0 *daemon.Connection
	switch res := ret.Get(0).(type) {
	case nil:
	case *daemon.Connection:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetConnections mocked method
func (m *GatewayerMock) GetConnections() *daemon.Connections {

	ret := m.Called()

	var r0 *daemon.Connections
	switch res := ret.Get(0).(type) {
	case nil:
	case *daemon.Connections:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetDefaultConnections mocked method
func (m *GatewayerMock) GetDefaultConnections() []string {

	ret := m.Called()

	var r0 []string
	switch res := ret.Get(0).(type) {
	case nil:
	case []string:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetExchgConnection mocked method
func (m *GatewayerMock) GetExchgConnection() []string {

	ret := m.Called()

	var r0 []string
	switch res := ret.Get(0).(type) {
	case nil:
	case []string:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetHealth mocked method
func (m *GatewayerMock) GetHealth() (*daemon.Health, error) {

	ret := m.Called()

	var r0 *daemon.Health
	switch res := ret.Get(0).(type) {
	case nil:
	case *daemon.Health:
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

// GetLastBlocks mocked method
func (m *GatewayerMock) GetLastBlocks(p0 uint64) (*visor.ReadableBlocks, error) {

	ret := m.Called(p0)

	var r0 *visor.ReadableBlocks
	switch res := ret.Get(0).(type) {
	case nil:
	case *visor.ReadableBlocks:
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

// GetRichlist mocked method
func (m *GatewayerMock) GetRichlist(p0 bool) (visor.Richlist, error) {

	ret := m.Called(p0)

	var r0 visor.Richlist
	switch res := ret.Get(0).(type) {
	case nil:
	case visor.Richlist:
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

// GetSignedBlockByHash mocked method
func (m *GatewayerMock) GetSignedBlockByHash(p0 cipher.SHA256) (*coin.SignedBlock, error) {

	ret := m.Called(p0)

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

// GetSignedBlockBySeq mocked method
func (m *GatewayerMock) GetSignedBlockBySeq(p0 uint64) (*coin.SignedBlock, error) {

	ret := m.Called(p0)

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

// GetTransaction mocked method
func (m *GatewayerMock) GetTransaction(p0 cipher.SHA256) (*visor.Transaction, error) {

	ret := m.Called(p0)

	var r0 *visor.Transaction
	switch res := ret.Get(0).(type) {
	case nil:
	case *visor.Transaction:
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

// GetTransactions mocked method
func (m *GatewayerMock) GetTransactions(p0 ...visor.TxFilter) ([]visor.Transaction, error) {

	ret := m.Called(p0)

	var r0 []visor.Transaction
	switch res := ret.Get(0).(type) {
	case nil:
	case []visor.Transaction:
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

// GetTransactionsForAddress mocked method
func (m *GatewayerMock) GetTransactionsForAddress(p0 cipher.Address) ([]daemon.ReadableTransaction, error) {

	ret := m.Called(p0)

	var r0 []daemon.ReadableTransaction
	switch res := ret.Get(0).(type) {
	case nil:
	case []daemon.ReadableTransaction:
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

// GetTrustConnections mocked method
func (m *GatewayerMock) GetTrustConnections() []string {

	ret := m.Called()

	var r0 []string
	switch res := ret.Get(0).(type) {
	case nil:
	case []string:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetUnspentOutputs mocked method
func (m *GatewayerMock) GetUnspentOutputs(p0 ...daemon.OutputsFilter) (*visor.ReadableOutputSet, error) {

	ret := m.Called(p0)

	var r0 *visor.ReadableOutputSet
	switch res := ret.Get(0).(type) {
	case nil:
	case *visor.ReadableOutputSet:
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

// GetUxOutByID mocked method
func (m *GatewayerMock) GetUxOutByID(p0 cipher.SHA256) (*historydb.UxOut, error) {

	ret := m.Called(p0)

	var r0 *historydb.UxOut
	switch res := ret.Get(0).(type) {
	case nil:
	case *historydb.UxOut:
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

// GetWallet mocked method
func (m *GatewayerMock) GetWallet(p0 string) (*wallet.Wallet, error) {

	ret := m.Called(p0)

	var r0 *wallet.Wallet
	switch res := ret.Get(0).(type) {
	case nil:
	case *wallet.Wallet:
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

// GetWalletBalance mocked method
func (m *GatewayerMock) GetWalletBalance(p0 string) (wallet.BalancePair, wallet.AddressBalance, error) {

	ret := m.Called(p0)

	var r0 wallet.BalancePair
	switch res := ret.Get(0).(type) {
	case nil:
	case wallet.BalancePair:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 wallet.AddressBalance
	switch res := ret.Get(1).(type) {
	case nil:
	case wallet.AddressBalance:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r2 error
	switch res := ret.Get(2).(type) {
	case nil:
	case error:
		r2 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1, r2

}

// GetWalletDir mocked method
func (m *GatewayerMock) GetWalletDir() (string, error) {

	ret := m.Called()

	var r0 string
	switch res := ret.Get(0).(type) {
	case nil:
	case string:
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

// GetWalletSeed mocked method
func (m *GatewayerMock) GetWalletSeed(p0 string, p1 []byte) (string, error) {

	ret := m.Called(p0, p1)

	var r0 string
	switch res := ret.Get(0).(type) {
	case nil:
	case string:
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

// GetWalletUnconfirmedTxns mocked method
func (m *GatewayerMock) GetWalletUnconfirmedTxns(p0 string) ([]visor.UnconfirmedTxn, error) {

	ret := m.Called(p0)

	var r0 []visor.UnconfirmedTxn
	switch res := ret.Get(0).(type) {
	case nil:
	case []visor.UnconfirmedTxn:
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
func (m *GatewayerMock) GetWallets() (wallet.Wallets, error) {

	ret := m.Called()

	var r0 wallet.Wallets
	switch res := ret.Get(0).(type) {
	case nil:
	case wallet.Wallets:
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

// InjectBroadcastTransaction mocked method
func (m *GatewayerMock) InjectBroadcastTransaction(p0 coin.Transaction) error {

	ret := m.Called(p0)

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

// IsWalletAPIEnabled mocked method
func (m *GatewayerMock) IsWalletAPIEnabled() bool {

	ret := m.Called()

	var r0 bool
	switch res := ret.Get(0).(type) {
	case nil:
	case bool:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// NewAddresses mocked method
func (m *GatewayerMock) NewAddresses(p0 string, p1 []byte, p2 uint64) ([]cipher.Address, error) {

	ret := m.Called(p0, p1, p2)

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

// ResendUnconfirmedTxns mocked method
func (m *GatewayerMock) ResendUnconfirmedTxns() (*daemon.ResendResult, error) {

	ret := m.Called()

	var r0 *daemon.ResendResult
	switch res := ret.Get(0).(type) {
	case nil:
	case *daemon.ResendResult:
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

// Spend mocked method
func (m *GatewayerMock) Spend(p0 string, p1 []byte, p2 uint64, p3 cipher.Address) (*coin.Transaction, error) {

	ret := m.Called(p0, p1, p2, p3)

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

// UnloadWallet mocked method
func (m *GatewayerMock) UnloadWallet(p0 string) error {

	ret := m.Called(p0)

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
func (m *GatewayerMock) UpdateWalletLabel(p0 string, p1 string) error {

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

// VerifyTxnVerbose mocked method
func (m *GatewayerMock) VerifyTxnVerbose(p0 *coin.Transaction) ([]wallet.UxBalance, bool, error) {

	ret := m.Called(p0)

	var r0 []wallet.UxBalance
	switch res := ret.Get(0).(type) {
	case nil:
	case []wallet.UxBalance:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 bool
	switch res := ret.Get(1).(type) {
	case nil:
	case bool:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r2 error
	switch res := ret.Get(2).(type) {
	case nil:
	case error:
		r2 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1, r2

}
