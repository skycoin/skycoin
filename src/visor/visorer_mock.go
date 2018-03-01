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
	historydb "github.com/skycoin/skycoin/src/visor/historydb"
	wallet "github.com/skycoin/skycoin/src/wallet"
)

// VisorerMock mock
type VisorerMock struct {
	mock.Mock
}

func NewVisorerMock() *VisorerMock {
	return &VisorerMock{}
}

// getTransactionsOfAddrs mocked method
func (m *VisorerMock) getTransactionsOfAddrs(p0 []cipher.Address) (map[cipher.Address][]Transaction, error) {

	ret := m.Called(p0)

	var r0 map[cipher.Address][]Transaction
	switch res := ret.Get(0).(type) {
	case nil:
	case map[cipher.Address][]Transaction:
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

// maybeCreateGenesisBlock mocked method
func (m *VisorerMock) maybeCreateGenesisBlock() error {

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

// traverseTxns mocked method
func (m *VisorerMock) traverseTxns(p0 ...TxFilter) ([]Transaction, error) {

	ret := m.Called(p0)

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

// AddressBalance mocked method
func (m *VisorerMock) AddressBalance(p0 coin.AddressUxOuts) (uint64, uint64, error) {

	ret := m.Called(p0)

	var r0 uint64
	switch res := ret.Get(0).(type) {
	case nil:
	case uint64:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 uint64
	switch res := ret.Get(1).(type) {
	case nil:
	case uint64:
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

// CreateAndExecuteBlock mocked method
func (m *VisorerMock) CreateAndExecuteBlock() (coin.SignedBlock, error) {

	ret := m.Called()

	var r0 coin.SignedBlock
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.SignedBlock:
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

// CreateBlock mocked method
func (m *VisorerMock) CreateBlock(p0 uint64) (coin.SignedBlock, error) {

	ret := m.Called(p0)

	var r0 coin.SignedBlock
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.SignedBlock:
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

// ExecuteSignedBlock mocked method
func (m *VisorerMock) ExecuteSignedBlock(p0 coin.SignedBlock) error {

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

// GenesisPreconditions mocked method
func (m *VisorerMock) GenesisPreconditions() {

	m.Called()

}

// GetAddrUxOuts mocked method
func (m *VisorerMock) GetAddrUxOuts(p0 cipher.Address) ([]*historydb.UxOut, error) {

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

// GetAddressTxns mocked method
func (m *VisorerMock) GetAddressTxns(p0 cipher.Address) ([]Transaction, error) {

	ret := m.Called(p0)

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

// GetAllUnconfirmedTxns mocked method
func (m *VisorerMock) GetAllUnconfirmedTxns() []UnconfirmedTxn {

	ret := m.Called()

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

// GetAllValidUnconfirmedTxHashes mocked method
func (m *VisorerMock) GetAllValidUnconfirmedTxHashes() []cipher.SHA256 {

	ret := m.Called()

	var r0 []cipher.SHA256
	switch res := ret.Get(0).(type) {
	case nil:
	case []cipher.SHA256:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetBalanceOfAddrs mocked method
func (m *VisorerMock) GetBalanceOfAddrs(p0 []cipher.Address) ([]wallet.BalancePair, error) {

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

// GetBlock mocked method
func (m *VisorerMock) GetBlock(p0 uint64) (*coin.SignedBlock, error) {

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

// GetBlockByHash mocked method
func (m *VisorerMock) GetBlockByHash(p0 cipher.SHA256) (*coin.SignedBlock, error) {

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

// GetBlockBySeq mocked method
func (m *VisorerMock) GetBlockBySeq(p0 uint64) (*coin.SignedBlock, error) {

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

// GetBlockchain mocked method
func (m *VisorerMock) GetBlockchain() Blockchainer {

	ret := m.Called()

	var r0 Blockchainer
	switch res := ret.Get(0).(type) {
	case nil:
	case Blockchainer:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetBlockchainMetadata mocked method
func (m *VisorerMock) GetBlockchainMetadata() BlockchainMetadata {

	ret := m.Called()

	var r0 BlockchainMetadata
	switch res := ret.Get(0).(type) {
	case nil:
	case BlockchainMetadata:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetBlocks mocked method
func (m *VisorerMock) GetBlocks(p0 uint64, p1 uint64) []coin.SignedBlock {

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

// GetConfig mocked method
func (m *VisorerMock) GetConfig() Config {

	ret := m.Called()

	var r0 Config
	switch res := ret.Get(0).(type) {
	case nil:
	case Config:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetHeadBlock mocked method
func (m *VisorerMock) GetHeadBlock() (*coin.SignedBlock, error) {

	ret := m.Called()

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

// GetLastBlocks mocked method
func (m *VisorerMock) GetLastBlocks(p0 uint64) []coin.SignedBlock {

	ret := m.Called(p0)

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

// GetLastTxs mocked method
func (m *VisorerMock) GetLastTxs() ([]*Transaction, error) {

	ret := m.Called()

	var r0 []*Transaction
	switch res := ret.Get(0).(type) {
	case nil:
	case []*Transaction:
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

// GetSignedBlocksSince mocked method
func (m *VisorerMock) GetSignedBlocksSince(p0 uint64, p1 uint64) ([]coin.SignedBlock, error) {

	ret := m.Called(p0, p1)

	var r0 []coin.SignedBlock
	switch res := ret.Get(0).(type) {
	case nil:
	case []coin.SignedBlock:
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
func (m *VisorerMock) GetTransaction(p0 cipher.SHA256) (*Transaction, error) {

	ret := m.Called(p0)

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

// GetTransactions mocked method
func (m *VisorerMock) GetTransactions(p0 ...TxFilter) ([]Transaction, error) {

	ret := m.Called(p0)

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

// GetUnconfirmed mocked method
func (m *VisorerMock) GetUnconfirmed() UnconfirmedTxnPooler {

	ret := m.Called()

	var r0 UnconfirmedTxnPooler
	switch res := ret.Get(0).(type) {
	case nil:
	case UnconfirmedTxnPooler:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetUnconfirmedTxns mocked method
func (m *VisorerMock) GetUnconfirmedTxns(p0 func(UnconfirmedTxn) bool) []UnconfirmedTxn {

	ret := m.Called(p0)

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

// GetUnspentOutputs mocked method
func (m *VisorerMock) GetUnspentOutputs() ([]coin.UxOut, error) {

	ret := m.Called()

	var r0 []coin.UxOut
	switch res := ret.Get(0).(type) {
	case nil:
	case []coin.UxOut:
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
func (m *VisorerMock) GetUxOutByID(p0 cipher.SHA256) (*historydb.UxOut, error) {

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

// HeadBkSeq mocked method
func (m *VisorerMock) HeadBkSeq() uint64 {

	ret := m.Called()

	var r0 uint64
	switch res := ret.Get(0).(type) {
	case nil:
	case uint64:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// InjectTransaction mocked method
func (m *VisorerMock) InjectTransaction(p0 coin.Transaction) (bool, *ErrTxnViolatesSoftConstraint, error) {

	ret := m.Called(p0)

	var r0 bool
	switch res := ret.Get(0).(type) {
	case nil:
	case bool:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 *ErrTxnViolatesSoftConstraint
	switch res := ret.Get(1).(type) {
	case nil:
	case *ErrTxnViolatesSoftConstraint:
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

// InjectTransactionStrict mocked method
func (m *VisorerMock) InjectTransactionStrict(p0 coin.Transaction) (bool, error) {

	ret := m.Called(p0)

	var r0 bool
	switch res := ret.Get(0).(type) {
	case nil:
	case bool:
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

// RefreshUnconfirmed mocked method
func (m *VisorerMock) RefreshUnconfirmed() ([]cipher.SHA256, error) {

	ret := m.Called()

	var r0 []cipher.SHA256
	switch res := ret.Get(0).(type) {
	case nil:
	case []cipher.SHA256:
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

// RemoveInvalidUnconfirmed mocked method
func (m *VisorerMock) RemoveInvalidUnconfirmed() ([]cipher.SHA256, error) {

	ret := m.Called()

	var r0 []cipher.SHA256
	switch res := ret.Get(0).(type) {
	case nil:
	case []cipher.SHA256:
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

// Run mocked method
func (m *VisorerMock) Run() error {

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

// ScanAheadWalletAddresses mocked method
func (m *VisorerMock) ScanAheadWalletAddresses(p0 string, p1 uint64) (wallet.Wallet, error) {

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

// Shutdown mocked method
func (m *VisorerMock) Shutdown() {

	m.Called()

}

// SignBlock mocked method
func (m *VisorerMock) SignBlock(p0 coin.Block) coin.SignedBlock {

	ret := m.Called(p0)

	var r0 coin.SignedBlock
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.SignedBlock:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// UnconfirmedIncomingOutputs mocked method
func (m *VisorerMock) UnconfirmedIncomingOutputs() (coin.UxArray, error) {

	ret := m.Called()

	var r0 coin.UxArray
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.UxArray:
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

// UnconfirmedSpendingOutputs mocked method
func (m *VisorerMock) UnconfirmedSpendingOutputs() (coin.UxArray, error) {

	ret := m.Called()

	var r0 coin.UxArray
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.UxArray:
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

// Wallets mocked method
func (m *VisorerMock) Wallets() *wallet.Service {

	ret := m.Called()

	var r0 *wallet.Service
	switch res := ret.Get(0).(type) {
	case nil:
	case *wallet.Service:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}
