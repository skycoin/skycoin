/*
* CODE GENERATED AUTOMATICALLY WITH github.com/ernesto-jimenez/goautomock
* THIS FILE MUST NEVER BE EDITED MANUALLY
 */

package webrpc

import (
	"fmt"
	mock "github.com/stretchr/testify/mock"

	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	visor "github.com/skycoin/skycoin/src/visor"
	historydb "github.com/skycoin/skycoin/src/visor/historydb"
)

// GatewayerMock mock
type GatewayerMock struct {
	mock.Mock
}

func NewGatewayerMock() *GatewayerMock {
	return &GatewayerMock{}
}

// GetAddrUxOuts mocked method
func (m *GatewayerMock) GetAddrUxOuts(p0 cipher.Address) ([]*historydb.UxOutJSON, error) {

	ret := m.Called(p0)

	var r0 []*historydb.UxOutJSON
	switch res := ret.Get(0).(type) {
	case nil:
	case []*historydb.UxOutJSON:
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
func (m *GatewayerMock) GetBlocks(p0 uint64, p1 uint64) *visor.ReadableBlocks {

	ret := m.Called(p0, p1)

	var r0 *visor.ReadableBlocks
	switch res := ret.Get(0).(type) {
	case nil:
	case *visor.ReadableBlocks:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetBlocksInDepth mocked method
func (m *GatewayerMock) GetBlocksInDepth(p0 []uint64) *visor.ReadableBlocks {

	ret := m.Called(p0)

	var r0 *visor.ReadableBlocks
	switch res := ret.Get(0).(type) {
	case nil:
	case *visor.ReadableBlocks:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetLastBlocks mocked method
func (m *GatewayerMock) GetLastBlocks(p0 uint64) *visor.ReadableBlocks {

	ret := m.Called(p0)

	var r0 *visor.ReadableBlocks
	switch res := ret.Get(0).(type) {
	case nil:
	case *visor.ReadableBlocks:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetTimeNow mocked method
func (m *GatewayerMock) GetTimeNow() uint64 {

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

// GetTransaction mocked method
func (m *GatewayerMock) GetTransaction(p0 cipher.SHA256) (*visor.TransactionResult, error) {

	ret := m.Called(p0)

	var r0 *visor.TransactionResult
	switch res := ret.Get(0).(type) {
	case nil:
	case *visor.TransactionResult:
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

// GetUnspentByAddrs mocked method
func (m *GatewayerMock) GetUnspentByAddrs(p0 []string) []visor.ReadableOutput {

	ret := m.Called(p0)

	var r0 []visor.ReadableOutput
	switch res := ret.Get(0).(type) {
	case nil:
	case []visor.ReadableOutput:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetUnspentByHashes mocked method
func (m *GatewayerMock) GetUnspentByHashes(p0 []string) []visor.ReadableOutput {

	ret := m.Called(p0)

	var r0 []visor.ReadableOutput
	switch res := ret.Get(0).(type) {
	case nil:
	case []visor.ReadableOutput:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// InjectTransaction mocked method
func (m *GatewayerMock) InjectTransaction(p0 coin.Transaction) (coin.Transaction, error) {

	ret := m.Called(p0)

	var r0 coin.Transaction
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.Transaction:
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
