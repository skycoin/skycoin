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
)

// historyerMock mock
type historyerMock struct {
	mock.Mock
}

func newHistoryerMock() *historyerMock {
	return &historyerMock{}
}

// ForEach mocked method
func (m *historyerMock) ForEach(p0 func(tx *historydb.Transaction) error) error {

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

// GetAddrTxns mocked method
func (m *historyerMock) GetAddrTxns(p0 cipher.Address) ([]historydb.Transaction, error) {

	ret := m.Called(p0)

	var r0 []historydb.Transaction
	switch res := ret.Get(0).(type) {
	case nil:
	case []historydb.Transaction:
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
func (m *historyerMock) GetAddrUxOuts(p0 cipher.Address) ([]*historydb.UxOut, error) {

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

// GetTransaction mocked method
func (m *historyerMock) GetTransaction(p0 cipher.SHA256) (*historydb.Transaction, error) {

	ret := m.Called(p0)

	var r0 *historydb.Transaction
	switch res := ret.Get(0).(type) {
	case nil:
	case *historydb.Transaction:
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

// GetUxout mocked method
func (m *historyerMock) GetUxout(p0 cipher.SHA256) (*historydb.UxOut, error) {

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

// ParseBlock mocked method
func (m *historyerMock) ParseBlock(p0 *coin.Block) error {

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

// ParsedHeight mocked method
func (m *historyerMock) ParsedHeight() int64 {

	ret := m.Called()

	var r0 int64
	switch res := ret.Get(0).(type) {
	case nil:
	case int64:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// ResetIfNeed mocked method
func (m *historyerMock) ResetIfNeed() error {

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
