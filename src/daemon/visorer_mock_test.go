/*
* CODE GENERATED AUTOMATICALLY WITH github.com/ernesto-jimenez/goautomock
* THIS FILE MUST NEVER BE EDITED MANUALLY
 */

package daemon

import (
	"fmt"

	mock "github.com/stretchr/testify/mock"

	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	visor "github.com/skycoin/skycoin/src/visor"
	wallet "github.com/skycoin/skycoin/src/wallet"
)

// VisorerMock mock
type VisorerMock struct {
	mock.Mock
}

func NewVisorerMock() *VisorerMock {
	return &VisorerMock{}
}

// broadcastBlock mocked method
func (m *VisorerMock) broadcastBlock(p0 coin.SignedBlock, p1 *Pool) error {

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

// broadcastTransaction mocked method
func (m *VisorerMock) broadcastTransaction(p0 coin.Transaction, p1 *Pool) error {

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

// AnnounceAllTxns mocked method
func (m *VisorerMock) AnnounceAllTxns(p0 *Pool) error {

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

// AnnounceBlocks mocked method
func (m *VisorerMock) AnnounceBlocks(p0 *Pool) error {

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

// AnnounceTxns mocked method
func (m *VisorerMock) AnnounceTxns(p0 *Pool, p1 []cipher.SHA256) error {

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

// CreateAndPublishBlock mocked method
func (m *VisorerMock) CreateAndPublishBlock(p0 *Pool) (coin.SignedBlock, error) {

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

// EstimateBlockchainHeight mocked method
func (m *VisorerMock) EstimateBlockchainHeight() uint64 {

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

// GetConfig mocked method
func (m *VisorerMock) GetConfig() VisorConfig {

	ret := m.Called()

	var r0 VisorConfig
	switch res := ret.Get(0).(type) {
	case nil:
	case VisorConfig:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetPeerBlockchainHeights mocked method
func (m *VisorerMock) GetPeerBlockchainHeights() []PeerBlockchainHeight {

	ret := m.Called()

	var r0 []PeerBlockchainHeight
	switch res := ret.Get(0).(type) {
	case nil:
	case []PeerBlockchainHeight:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetSignedBlock mocked method
func (m *VisorerMock) GetSignedBlock(p0 uint64) (*coin.SignedBlock, error) {

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

// GetVisor mocked method
func (m *VisorerMock) GetVisor() visor.Visorer {

	ret := m.Called()

	var r0 visor.Visorer
	switch res := ret.Get(0).(type) {
	case nil:
	case visor.Visorer:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

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

// InjectBroadcastTransaction mocked method
func (m *VisorerMock) InjectBroadcastTransaction(p0 coin.Transaction, p1 *Pool) error {

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

// InjectTransaction mocked method
func (m *VisorerMock) InjectTransaction(p0 coin.Transaction) (bool, *visor.ErrTxnViolatesSoftConstraint, error) {

	ret := m.Called(p0)

	var r0 bool
	switch res := ret.Get(0).(type) {
	case nil:
	case bool:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 *visor.ErrTxnViolatesSoftConstraint
	switch res := ret.Get(1).(type) {
	case nil:
	case *visor.ErrTxnViolatesSoftConstraint:
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

// RecordBlockchainHeight mocked method
func (m *VisorerMock) RecordBlockchainHeight(p0 string, p1 uint64) {

	m.Called(p0, p1)

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

// RemoveConnection mocked method
func (m *VisorerMock) RemoveConnection(p0 string) {

	m.Called(p0)

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

// RequestBlocks mocked method
func (m *VisorerMock) RequestBlocks(p0 *Pool) error {

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

// RequestBlocksFromAddr mocked method
func (m *VisorerMock) RequestBlocksFromAddr(p0 *Pool, p1 string) error {

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

// ResendTransaction mocked method
func (m *VisorerMock) ResendTransaction(p0 cipher.SHA256, p1 *Pool) error {

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

// ResendUnconfirmedTxns mocked method
func (m *VisorerMock) ResendUnconfirmedTxns(p0 *Pool) []cipher.SHA256 {

	ret := m.Called(p0)

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

// SetTxnsAnnounced mocked method
func (m *VisorerMock) SetTxnsAnnounced(p0 []cipher.SHA256) {

	m.Called(p0)

}

// Shutdown mocked method
func (m *VisorerMock) Shutdown() {

	m.Called()

}

// UnConfirmFilterKnown mocked method
func (m *VisorerMock) UnConfirmFilterKnown(p0 []cipher.SHA256) []cipher.SHA256 {

	ret := m.Called(p0)

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

// UnConfirmKnow mocked method
func (m *VisorerMock) UnConfirmKnow(p0 []cipher.SHA256) coin.Transactions {

	ret := m.Called(p0)

	var r0 coin.Transactions
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.Transactions:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}
