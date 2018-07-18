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
	gnet "github.com/skycoin/skycoin/src/daemon/gnet"
	pex "github.com/skycoin/skycoin/src/daemon/pex"
	visor "github.com/skycoin/skycoin/src/visor"
)

// DaemonerMock mock
type DaemonerMock struct {
	mock.Mock
}

func NewDaemonerMock() *DaemonerMock {
	return &DaemonerMock{}
}

// AddPeer mocked method
func (m *DaemonerMock) AddPeer(p0 string) error {

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

// AddPeers mocked method
func (m *DaemonerMock) AddPeers(p0 []string) int {

	ret := m.Called(p0)

	var r0 int
	switch res := ret.Get(0).(type) {
	case nil:
	case int:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// AnnounceAllTxns mocked method
func (m *DaemonerMock) AnnounceAllTxns() error {

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

// BlockchainPubkey mocked method
func (m *DaemonerMock) BlockchainPubkey() cipher.PubKey {

	ret := m.Called()

	var r0 cipher.PubKey
	switch res := ret.Get(0).(type) {
	case nil:
	case cipher.PubKey:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// BroadcastMessage mocked method
func (m *DaemonerMock) BroadcastMessage(p0 gnet.Message) error {

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

// DaemonConfig mocked method
func (m *DaemonerMock) DaemonConfig() DaemonConfig {

	ret := m.Called()

	var r0 DaemonConfig
	switch res := ret.Get(0).(type) {
	case nil:
	case DaemonConfig:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// Disconnect mocked method
func (m *DaemonerMock) Disconnect(p0 string, p1 gnet.DisconnectReason) error {

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

// ExecuteSignedBlock mocked method
func (m *DaemonerMock) ExecuteSignedBlock(p0 coin.SignedBlock) error {

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

// GetMirrorPort mocked method
func (m *DaemonerMock) GetMirrorPort(p0 string, p1 uint32) (uint16, bool) {

	ret := m.Called(p0, p1)

	var r0 uint16
	switch res := ret.Get(0).(type) {
	case nil:
	case uint16:
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

	return r0, r1

}

// GetSignedBlocksSince mocked method
func (m *DaemonerMock) GetSignedBlocksSince(p0 uint64, p1 uint64) ([]coin.SignedBlock, error) {

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

// GetUnconfirmedKnown mocked method
func (m *DaemonerMock) GetUnconfirmedKnown(p0 []cipher.SHA256) (coin.Transactions, error) {

	ret := m.Called(p0)

	var r0 coin.Transactions
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.Transactions:
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

// GetUnconfirmedUnknown mocked method
func (m *DaemonerMock) GetUnconfirmedUnknown(p0 []cipher.SHA256) ([]cipher.SHA256, error) {

	ret := m.Called(p0)

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

// HeadBkSeq mocked method
func (m *DaemonerMock) HeadBkSeq() (uint64, bool, error) {

	ret := m.Called()

	var r0 uint64
	switch res := ret.Get(0).(type) {
	case nil:
	case uint64:
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

// IncreaseRetryTimes mocked method
func (m *DaemonerMock) IncreaseRetryTimes(p0 string) {

	m.Called(p0)

}

// InjectTransaction mocked method
func (m *DaemonerMock) InjectTransaction(p0 coin.Transaction) (bool, *visor.ErrTxnViolatesSoftConstraint, error) {

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

// IsDefaultConnection mocked method
func (m *DaemonerMock) IsDefaultConnection(p0 string) bool {

	ret := m.Called(p0)

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

// IsMaxDefaultConnectionsReached mocked method
func (m *DaemonerMock) IsMaxDefaultConnectionsReached() (bool, error) {

	ret := m.Called()

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

// Mirror mocked method
func (m *DaemonerMock) Mirror() uint32 {

	ret := m.Called()

	var r0 uint32
	switch res := ret.Get(0).(type) {
	case nil:
	case uint32:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// PexConfig mocked method
func (m *DaemonerMock) PexConfig() pex.Config {

	ret := m.Called()

	var r0 pex.Config
	switch res := ret.Get(0).(type) {
	case nil:
	case pex.Config:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// RandomExchangeable mocked method
func (m *DaemonerMock) RandomExchangeable(p0 int) pex.Peers {

	ret := m.Called(p0)

	var r0 pex.Peers
	switch res := ret.Get(0).(type) {
	case nil:
	case pex.Peers:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// RecordConnectionMirror mocked method
func (m *DaemonerMock) RecordConnectionMirror(p0 string, p1 uint32) error {

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

// RecordMessageEvent mocked method
func (m *DaemonerMock) RecordMessageEvent(p0 AsyncMessage, p1 *gnet.MessageContext) error {

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

// RecordPeerHeight mocked method
func (m *DaemonerMock) RecordPeerHeight(p0 string, p1 uint64) {

	m.Called(p0, p1)

}

// RemoveFromExpectingIntroductions mocked method
func (m *DaemonerMock) RemoveFromExpectingIntroductions(p0 string) {

	m.Called(p0)

}

// RequestBlocksFromAddr mocked method
func (m *DaemonerMock) RequestBlocksFromAddr(p0 string) error {

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

// ResetRetryTimes mocked method
func (m *DaemonerMock) ResetRetryTimes(p0 string) {

	m.Called(p0)

}

// SendMessage mocked method
func (m *DaemonerMock) SendMessage(p0 string, p1 gnet.Message) error {

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

// SetHasIncomingPort mocked method
func (m *DaemonerMock) SetHasIncomingPort(p0 string) error {

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
