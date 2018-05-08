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
	dbutil "github.com/skycoin/skycoin/src/visor/dbutil"
	historydb "github.com/skycoin/skycoin/src/visor/historydb"
)

// historyerMock mock
type historyerMock struct {
	mock.Mock
}

func newHistoryerMock() *historyerMock {
	return &historyerMock{}
}

// Erase mocked method
func (m *historyerMock) Erase(p0 *dbutil.Tx) error {

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

// ForEachTxn mocked method
func (m *historyerMock) ForEachTxn(p0 *dbutil.Tx, p1 func(cipher.SHA256, *historydb.Transaction) error) error {

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

// GetAddrTxns mocked method
func (m *historyerMock) GetAddrTxns(p0 *dbutil.Tx, p1 cipher.Address) ([]historydb.Transaction, error) {

	ret := m.Called(p0, p1)

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
func (m *historyerMock) GetAddrUxOuts(p0 *dbutil.Tx, p1 cipher.Address) ([]*historydb.UxOut, error) {

	ret := m.Called(p0, p1)

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
func (m *historyerMock) GetTransaction(p0 *dbutil.Tx, p1 cipher.SHA256) (*historydb.Transaction, error) {

	ret := m.Called(p0, p1)

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

// GetUxOut mocked method
func (m *historyerMock) GetUxOut(p0 *dbutil.Tx, p1 cipher.SHA256) (*historydb.UxOut, error) {

	ret := m.Called(p0, p1)

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

// NeedsReset mocked method
func (m *historyerMock) NeedsReset(p0 *dbutil.Tx) (bool, error) {

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

// ParseBlock mocked method
func (m *historyerMock) ParseBlock(p0 *dbutil.Tx, p1 coin.Block) error {

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

// ParsedHeight mocked method
func (m *historyerMock) ParsedHeight(p0 *dbutil.Tx) (uint64, bool, error) {

	ret := m.Called(p0)

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
