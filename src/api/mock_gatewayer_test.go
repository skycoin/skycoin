// Code generated by mockery v1.0.0. DO NOT EDIT.

package api

import cipher "github.com/skycoin/skycoin/src/cipher"
import coin "github.com/skycoin/skycoin/src/coin"
import daemon "github.com/skycoin/skycoin/src/daemon"
import historydb "github.com/skycoin/skycoin/src/visor/historydb"
import mock "github.com/stretchr/testify/mock"
import notes "github.com/skycoin/skycoin/src/notes"
import visor "github.com/skycoin/skycoin/src/visor"
import wallet "github.com/skycoin/skycoin/src/wallet"

// MockGatewayer is an autogenerated mock type for the Gatewayer type
type MockGatewayer struct {
	mock.Mock
}

// AddNote provides a mock function with given fields: _a0
func (_m *MockGatewayer) AddNote(_a0 notes.Note) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(notes.Note) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateTransaction provides a mock function with given fields: w
func (_m *MockGatewayer) CreateTransaction(w wallet.CreateTransactionParams) (*coin.Transaction, []wallet.UxBalance, error) {
	ret := _m.Called(w)

	var r0 *coin.Transaction
	if rf, ok := ret.Get(0).(func(wallet.CreateTransactionParams) *coin.Transaction); ok {
		r0 = rf(w)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coin.Transaction)
		}
	}

	var r1 []wallet.UxBalance
	if rf, ok := ret.Get(1).(func(wallet.CreateTransactionParams) []wallet.UxBalance); ok {
		r1 = rf(w)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]wallet.UxBalance)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(wallet.CreateTransactionParams) error); ok {
		r2 = rf(w)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CreateWallet provides a mock function with given fields: wltName, options
func (_m *MockGatewayer) CreateWallet(wltName string, options wallet.Options) (*wallet.Wallet, error) {
	ret := _m.Called(wltName, options)

	var r0 *wallet.Wallet
	if rf, ok := ret.Get(0).(func(string, wallet.Options) *wallet.Wallet); ok {
		r0 = rf(wltName, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wallet.Wallet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, wallet.Options) error); ok {
		r1 = rf(wltName, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DecryptWallet provides a mock function with given fields: wltID, password
func (_m *MockGatewayer) DecryptWallet(wltID string, password []byte) (*wallet.Wallet, error) {
	ret := _m.Called(wltID, password)

	var r0 *wallet.Wallet
	if rf, ok := ret.Get(0).(func(string, []byte) *wallet.Wallet); ok {
		r0 = rf(wltID, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wallet.Wallet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []byte) error); ok {
		r1 = rf(wltID, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EncryptWallet provides a mock function with given fields: wltID, password
func (_m *MockGatewayer) EncryptWallet(wltID string, password []byte) (*wallet.Wallet, error) {
	ret := _m.Called(wltID, password)

	var r0 *wallet.Wallet
	if rf, ok := ret.Get(0).(func(string, []byte) *wallet.Wallet); ok {
		r0 = rf(wltID, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wallet.Wallet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []byte) error); ok {
		r1 = rf(wltID, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAddressCount provides a mock function with given fields:
func (_m *MockGatewayer) GetAddressCount() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllNotes provides a mock function with given fields:
func (_m *MockGatewayer) GetAllNotes() []notes.Note {
	ret := _m.Called()

	var r0 []notes.Note
	if rf, ok := ret.Get(0).(func() []notes.Note); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]notes.Note)
		}
	}

	return r0
}

// GetAllUnconfirmedTransactions provides a mock function with given fields:
func (_m *MockGatewayer) GetAllUnconfirmedTransactions() ([]visor.UnconfirmedTransaction, error) {
	ret := _m.Called()

	var r0 []visor.UnconfirmedTransaction
	if rf, ok := ret.Get(0).(func() []visor.UnconfirmedTransaction); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]visor.UnconfirmedTransaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllUnconfirmedTransactionsVerbose provides a mock function with given fields:
func (_m *MockGatewayer) GetAllUnconfirmedTransactionsVerbose() ([]visor.UnconfirmedTransaction, [][]visor.TransactionInput, error) {
	ret := _m.Called()

	var r0 []visor.UnconfirmedTransaction
	if rf, ok := ret.Get(0).(func() []visor.UnconfirmedTransaction); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]visor.UnconfirmedTransaction)
		}
	}

	var r1 [][]visor.TransactionInput
	if rf, ok := ret.Get(1).(func() [][]visor.TransactionInput); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([][]visor.TransactionInput)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetBalanceOfAddrs provides a mock function with given fields: addrs
func (_m *MockGatewayer) GetBalanceOfAddrs(addrs []cipher.Address) ([]wallet.BalancePair, error) {
	ret := _m.Called(addrs)

	var r0 []wallet.BalancePair
	if rf, ok := ret.Get(0).(func([]cipher.Address) []wallet.BalancePair); ok {
		r0 = rf(addrs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]wallet.BalancePair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]cipher.Address) error); ok {
		r1 = rf(addrs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockchainMetadata provides a mock function with given fields:
func (_m *MockGatewayer) GetBlockchainMetadata() (*visor.BlockchainMetadata, error) {
	ret := _m.Called()

	var r0 *visor.BlockchainMetadata
	if rf, ok := ret.Get(0).(func() *visor.BlockchainMetadata); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*visor.BlockchainMetadata)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockchainProgress provides a mock function with given fields:
func (_m *MockGatewayer) GetBlockchainProgress() (*daemon.BlockchainProgress, error) {
	ret := _m.Called()

	var r0 *daemon.BlockchainProgress
	if rf, ok := ret.Get(0).(func() *daemon.BlockchainProgress); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*daemon.BlockchainProgress)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlocksInRange provides a mock function with given fields: start, end
func (_m *MockGatewayer) GetBlocksInRange(start uint64, end uint64) ([]coin.SignedBlock, error) {
	ret := _m.Called(start, end)

	var r0 []coin.SignedBlock
	if rf, ok := ret.Get(0).(func(uint64, uint64) []coin.SignedBlock); ok {
		r0 = rf(start, end)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]coin.SignedBlock)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, uint64) error); ok {
		r1 = rf(start, end)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlocksInRangeVerbose provides a mock function with given fields: start, end
func (_m *MockGatewayer) GetBlocksInRangeVerbose(start uint64, end uint64) ([]coin.SignedBlock, [][][]visor.TransactionInput, error) {
	ret := _m.Called(start, end)

	var r0 []coin.SignedBlock
	if rf, ok := ret.Get(0).(func(uint64, uint64) []coin.SignedBlock); ok {
		r0 = rf(start, end)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]coin.SignedBlock)
		}
	}

	var r1 [][][]visor.TransactionInput
	if rf, ok := ret.Get(1).(func(uint64, uint64) [][][]visor.TransactionInput); ok {
		r1 = rf(start, end)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([][][]visor.TransactionInput)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(uint64, uint64) error); ok {
		r2 = rf(start, end)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetConnection provides a mock function with given fields: addr
func (_m *MockGatewayer) GetConnection(addr string) (*daemon.Connection, error) {
	ret := _m.Called(addr)

	var r0 *daemon.Connection
	if rf, ok := ret.Get(0).(func(string) *daemon.Connection); ok {
		r0 = rf(addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*daemon.Connection)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(addr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDefaultConnections provides a mock function with given fields:
func (_m *MockGatewayer) GetDefaultConnections() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// GetExchgConnection provides a mock function with given fields:
func (_m *MockGatewayer) GetExchgConnection() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// GetHealth provides a mock function with given fields:
func (_m *MockGatewayer) GetHealth() (*daemon.Health, error) {
	ret := _m.Called()

	var r0 *daemon.Health
	if rf, ok := ret.Get(0).(func() *daemon.Health); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*daemon.Health)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastBlocks provides a mock function with given fields: num
func (_m *MockGatewayer) GetLastBlocks(num uint64) ([]coin.SignedBlock, error) {
	ret := _m.Called(num)

	var r0 []coin.SignedBlock
	if rf, ok := ret.Get(0).(func(uint64) []coin.SignedBlock); ok {
		r0 = rf(num)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]coin.SignedBlock)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(num)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastBlocksVerbose provides a mock function with given fields: num
func (_m *MockGatewayer) GetLastBlocksVerbose(num uint64) ([]coin.SignedBlock, [][][]visor.TransactionInput, error) {
	ret := _m.Called(num)

	var r0 []coin.SignedBlock
	if rf, ok := ret.Get(0).(func(uint64) []coin.SignedBlock); ok {
		r0 = rf(num)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]coin.SignedBlock)
		}
	}

	var r1 [][][]visor.TransactionInput
	if rf, ok := ret.Get(1).(func(uint64) [][][]visor.TransactionInput); ok {
		r1 = rf(num)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([][][]visor.TransactionInput)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(uint64) error); ok {
		r2 = rf(num)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetNoteByTxID provides a mock function with given fields: txID
func (_m *MockGatewayer) GetNoteByTxID(txID string) notes.Note {
	ret := _m.Called(txID)

	var r0 notes.Note
	if rf, ok := ret.Get(0).(func(string) notes.Note); ok {
		r0 = rf(txID)
	} else {
		r0 = ret.Get(0).(notes.Note)
	}

	return r0
}

// GetOutgoingConnections provides a mock function with given fields:
func (_m *MockGatewayer) GetOutgoingConnections() ([]daemon.Connection, error) {
	ret := _m.Called()

	var r0 []daemon.Connection
	if rf, ok := ret.Get(0).(func() []daemon.Connection); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]daemon.Connection)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRichlist provides a mock function with given fields: includeDistribution
func (_m *MockGatewayer) GetRichlist(includeDistribution bool) (visor.Richlist, error) {
	ret := _m.Called(includeDistribution)

	var r0 visor.Richlist
	if rf, ok := ret.Get(0).(func(bool) visor.Richlist); ok {
		r0 = rf(includeDistribution)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(visor.Richlist)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bool) error); ok {
		r1 = rf(includeDistribution)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSignedBlockByHash provides a mock function with given fields: hash
func (_m *MockGatewayer) GetSignedBlockByHash(hash cipher.SHA256) (*coin.SignedBlock, error) {
	ret := _m.Called(hash)

	var r0 *coin.SignedBlock
	if rf, ok := ret.Get(0).(func(cipher.SHA256) *coin.SignedBlock); ok {
		r0 = rf(hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coin.SignedBlock)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(cipher.SHA256) error); ok {
		r1 = rf(hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSignedBlockByHashVerbose provides a mock function with given fields: hash
func (_m *MockGatewayer) GetSignedBlockByHashVerbose(hash cipher.SHA256) (*coin.SignedBlock, [][]visor.TransactionInput, error) {
	ret := _m.Called(hash)

	var r0 *coin.SignedBlock
	if rf, ok := ret.Get(0).(func(cipher.SHA256) *coin.SignedBlock); ok {
		r0 = rf(hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coin.SignedBlock)
		}
	}

	var r1 [][]visor.TransactionInput
	if rf, ok := ret.Get(1).(func(cipher.SHA256) [][]visor.TransactionInput); ok {
		r1 = rf(hash)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([][]visor.TransactionInput)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(cipher.SHA256) error); ok {
		r2 = rf(hash)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetSignedBlockBySeq provides a mock function with given fields: seq
func (_m *MockGatewayer) GetSignedBlockBySeq(seq uint64) (*coin.SignedBlock, error) {
	ret := _m.Called(seq)

	var r0 *coin.SignedBlock
	if rf, ok := ret.Get(0).(func(uint64) *coin.SignedBlock); ok {
		r0 = rf(seq)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coin.SignedBlock)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(seq)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSignedBlockBySeqVerbose provides a mock function with given fields: seq
func (_m *MockGatewayer) GetSignedBlockBySeqVerbose(seq uint64) (*coin.SignedBlock, [][]visor.TransactionInput, error) {
	ret := _m.Called(seq)

	var r0 *coin.SignedBlock
	if rf, ok := ret.Get(0).(func(uint64) *coin.SignedBlock); ok {
		r0 = rf(seq)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coin.SignedBlock)
		}
	}

	var r1 [][]visor.TransactionInput
	if rf, ok := ret.Get(1).(func(uint64) [][]visor.TransactionInput); ok {
		r1 = rf(seq)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([][]visor.TransactionInput)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(uint64) error); ok {
		r2 = rf(seq)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetSpentOutputsForAddresses provides a mock function with given fields: addr
func (_m *MockGatewayer) GetSpentOutputsForAddresses(addr []cipher.Address) ([][]historydb.UxOut, error) {
	ret := _m.Called(addr)

	var r0 [][]historydb.UxOut
	if rf, ok := ret.Get(0).(func([]cipher.Address) [][]historydb.UxOut); ok {
		r0 = rf(addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]historydb.UxOut)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]cipher.Address) error); ok {
		r1 = rf(addr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransaction provides a mock function with given fields: txid
func (_m *MockGatewayer) GetTransaction(txid cipher.SHA256) (*visor.Transaction, error) {
	ret := _m.Called(txid)

	var r0 *visor.Transaction
	if rf, ok := ret.Get(0).(func(cipher.SHA256) *visor.Transaction); ok {
		r0 = rf(txid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*visor.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(cipher.SHA256) error); ok {
		r1 = rf(txid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionVerbose provides a mock function with given fields: txid
func (_m *MockGatewayer) GetTransactionVerbose(txid cipher.SHA256) (*visor.Transaction, []visor.TransactionInput, error) {
	ret := _m.Called(txid)

	var r0 *visor.Transaction
	if rf, ok := ret.Get(0).(func(cipher.SHA256) *visor.Transaction); ok {
		r0 = rf(txid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*visor.Transaction)
		}
	}

	var r1 []visor.TransactionInput
	if rf, ok := ret.Get(1).(func(cipher.SHA256) []visor.TransactionInput); ok {
		r1 = rf(txid)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]visor.TransactionInput)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(cipher.SHA256) error); ok {
		r2 = rf(txid)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTransactions provides a mock function with given fields: flts
func (_m *MockGatewayer) GetTransactions(flts []visor.TxFilter) ([]visor.Transaction, error) {
	ret := _m.Called(flts)

	var r0 []visor.Transaction
	if rf, ok := ret.Get(0).(func([]visor.TxFilter) []visor.Transaction); ok {
		r0 = rf(flts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]visor.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]visor.TxFilter) error); ok {
		r1 = rf(flts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionsVerbose provides a mock function with given fields: flts
func (_m *MockGatewayer) GetTransactionsVerbose(flts []visor.TxFilter) ([]visor.Transaction, [][]visor.TransactionInput, error) {
	ret := _m.Called(flts)

	var r0 []visor.Transaction
	if rf, ok := ret.Get(0).(func([]visor.TxFilter) []visor.Transaction); ok {
		r0 = rf(flts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]visor.Transaction)
		}
	}

	var r1 [][]visor.TransactionInput
	if rf, ok := ret.Get(1).(func([]visor.TxFilter) [][]visor.TransactionInput); ok {
		r1 = rf(flts)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([][]visor.TransactionInput)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func([]visor.TxFilter) error); ok {
		r2 = rf(flts)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTrustConnections provides a mock function with given fields:
func (_m *MockGatewayer) GetTrustConnections() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// GetUnspentOutputsSummary provides a mock function with given fields: filters
func (_m *MockGatewayer) GetUnspentOutputsSummary(filters []visor.OutputsFilter) (*visor.UnspentOutputsSummary, error) {
	ret := _m.Called(filters)

	var r0 *visor.UnspentOutputsSummary
	if rf, ok := ret.Get(0).(func([]visor.OutputsFilter) *visor.UnspentOutputsSummary); ok {
		r0 = rf(filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*visor.UnspentOutputsSummary)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]visor.OutputsFilter) error); ok {
		r1 = rf(filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUxOutByID provides a mock function with given fields: id
func (_m *MockGatewayer) GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error) {
	ret := _m.Called(id)

	var r0 *historydb.UxOut
	if rf, ok := ret.Get(0).(func(cipher.SHA256) *historydb.UxOut); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*historydb.UxOut)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(cipher.SHA256) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVerboseTransactionsForAddress provides a mock function with given fields: a
func (_m *MockGatewayer) GetVerboseTransactionsForAddress(a cipher.Address) ([]visor.Transaction, [][]visor.TransactionInput, error) {
	ret := _m.Called(a)

	var r0 []visor.Transaction
	if rf, ok := ret.Get(0).(func(cipher.Address) []visor.Transaction); ok {
		r0 = rf(a)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]visor.Transaction)
		}
	}

	var r1 [][]visor.TransactionInput
	if rf, ok := ret.Get(1).(func(cipher.Address) [][]visor.TransactionInput); ok {
		r1 = rf(a)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([][]visor.TransactionInput)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(cipher.Address) error); ok {
		r2 = rf(a)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetWallet provides a mock function with given fields: wltID
func (_m *MockGatewayer) GetWallet(wltID string) (*wallet.Wallet, error) {
	ret := _m.Called(wltID)

	var r0 *wallet.Wallet
	if rf, ok := ret.Get(0).(func(string) *wallet.Wallet); ok {
		r0 = rf(wltID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wallet.Wallet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(wltID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWalletBalance provides a mock function with given fields: wltID
func (_m *MockGatewayer) GetWalletBalance(wltID string) (wallet.BalancePair, wallet.AddressBalances, error) {
	ret := _m.Called(wltID)

	var r0 wallet.BalancePair
	if rf, ok := ret.Get(0).(func(string) wallet.BalancePair); ok {
		r0 = rf(wltID)
	} else {
		r0 = ret.Get(0).(wallet.BalancePair)
	}

	var r1 wallet.AddressBalances
	if rf, ok := ret.Get(1).(func(string) wallet.AddressBalances); ok {
		r1 = rf(wltID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(wallet.AddressBalances)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(wltID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetWalletDir provides a mock function with given fields:
func (_m *MockGatewayer) GetWalletDir() (string, error) {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWalletSeed provides a mock function with given fields: wltID, password
func (_m *MockGatewayer) GetWalletSeed(wltID string, password []byte) (string, error) {
	ret := _m.Called(wltID, password)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, []byte) string); ok {
		r0 = rf(wltID, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []byte) error); ok {
		r1 = rf(wltID, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWalletUnconfirmedTransactions provides a mock function with given fields: wltID
func (_m *MockGatewayer) GetWalletUnconfirmedTransactions(wltID string) ([]visor.UnconfirmedTransaction, error) {
	ret := _m.Called(wltID)

	var r0 []visor.UnconfirmedTransaction
	if rf, ok := ret.Get(0).(func(string) []visor.UnconfirmedTransaction); ok {
		r0 = rf(wltID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]visor.UnconfirmedTransaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(wltID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWalletUnconfirmedTransactionsVerbose provides a mock function with given fields: wltID
func (_m *MockGatewayer) GetWalletUnconfirmedTransactionsVerbose(wltID string) ([]visor.UnconfirmedTransaction, [][]visor.TransactionInput, error) {
	ret := _m.Called(wltID)

	var r0 []visor.UnconfirmedTransaction
	if rf, ok := ret.Get(0).(func(string) []visor.UnconfirmedTransaction); ok {
		r0 = rf(wltID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]visor.UnconfirmedTransaction)
		}
	}

	var r1 [][]visor.TransactionInput
	if rf, ok := ret.Get(1).(func(string) [][]visor.TransactionInput); ok {
		r1 = rf(wltID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([][]visor.TransactionInput)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string) error); ok {
		r2 = rf(wltID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetWallets provides a mock function with given fields:
func (_m *MockGatewayer) GetWallets() (wallet.Wallets, error) {
	ret := _m.Called()

	var r0 wallet.Wallets
	if rf, ok := ret.Get(0).(func() wallet.Wallets); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(wallet.Wallets)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InjectBroadcastTransaction provides a mock function with given fields: txn
func (_m *MockGatewayer) InjectBroadcastTransaction(txn coin.Transaction) error {
	ret := _m.Called(txn)

	var r0 error
	if rf, ok := ret.Get(0).(func(coin.Transaction) error); ok {
		r0 = rf(txn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewAddresses provides a mock function with given fields: wltID, password, n
func (_m *MockGatewayer) NewAddresses(wltID string, password []byte, n uint64) ([]cipher.Address, error) {
	ret := _m.Called(wltID, password, n)

	var r0 []cipher.Address
	if rf, ok := ret.Get(0).(func(string, []byte, uint64) []cipher.Address); ok {
		r0 = rf(wltID, password, n)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]cipher.Address)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []byte, uint64) error); ok {
		r1 = rf(wltID, password, n)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveNote provides a mock function with given fields: txID
func (_m *MockGatewayer) RemoveNote(txID string) error {
	ret := _m.Called(txID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(txID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResendUnconfirmedTxns provides a mock function with given fields:
func (_m *MockGatewayer) ResendUnconfirmedTxns() ([]cipher.SHA256, error) {
	ret := _m.Called()

	var r0 []cipher.SHA256
	if rf, ok := ret.Get(0).(func() []cipher.SHA256); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]cipher.SHA256)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Spend provides a mock function with given fields: wltID, password, coins, dest
func (_m *MockGatewayer) Spend(wltID string, password []byte, coins uint64, dest cipher.Address) (*coin.Transaction, error) {
	ret := _m.Called(wltID, password, coins, dest)

	var r0 *coin.Transaction
	if rf, ok := ret.Get(0).(func(string, []byte, uint64, cipher.Address) *coin.Transaction); ok {
		r0 = rf(wltID, password, coins, dest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coin.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []byte, uint64, cipher.Address) error); ok {
		r1 = rf(wltID, password, coins, dest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnloadWallet provides a mock function with given fields: id
func (_m *MockGatewayer) UnloadWallet(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateWalletLabel provides a mock function with given fields: wltID, label
func (_m *MockGatewayer) UpdateWalletLabel(wltID string, label string) error {
	ret := _m.Called(wltID, label)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(wltID, label)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyTxnVerbose provides a mock function with given fields: txn
func (_m *MockGatewayer) VerifyTxnVerbose(txn *coin.Transaction) ([]wallet.UxBalance, bool, error) {
	ret := _m.Called(txn)

	var r0 []wallet.UxBalance
	if rf, ok := ret.Get(0).(func(*coin.Transaction) []wallet.UxBalance); ok {
		r0 = rf(txn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]wallet.UxBalance)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(*coin.Transaction) bool); ok {
		r1 = rf(txn)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*coin.Transaction) error); ok {
		r2 = rf(txn)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
