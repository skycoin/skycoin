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
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

// UnconfirmedTxnPoolerMock mock
type UnconfirmedTxnPoolerMock struct {
	mock.Mock
}

func NewUnconfirmedTxnPoolerMock() *UnconfirmedTxnPoolerMock {
	return &UnconfirmedTxnPoolerMock{}
}

// ForEach mocked method
func (m *UnconfirmedTxnPoolerMock) ForEach(p0 *dbutil.Tx, p1 func(cipher.SHA256, UnconfirmedTxn) error) error {

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

// Get mocked method
func (m *UnconfirmedTxnPoolerMock) Get(p0 *dbutil.Tx, p1 cipher.SHA256) (*UnconfirmedTxn, error) {

	ret := m.Called(p0, p1)

	var r0 *UnconfirmedTxn
	switch res := ret.Get(0).(type) {
	case nil:
	case *UnconfirmedTxn:
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

// GetIncomingOutputs mocked method
func (m *UnconfirmedTxnPoolerMock) GetIncomingOutputs(p0 *dbutil.Tx, p1 coin.BlockHeader) (coin.UxArray, error) {

	ret := m.Called(p0, p1)

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

// GetKnown mocked method
func (m *UnconfirmedTxnPoolerMock) GetKnown(p0 *dbutil.Tx, p1 []cipher.SHA256) (coin.Transactions, error) {

	ret := m.Called(p0, p1)

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

// GetTxHashes mocked method
func (m *UnconfirmedTxnPoolerMock) GetTxHashes(p0 *dbutil.Tx, p1 func(tx UnconfirmedTxn) bool) ([]cipher.SHA256, error) {

	ret := m.Called(p0, p1)

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

// GetTxns mocked method
func (m *UnconfirmedTxnPoolerMock) GetTxns(p0 *dbutil.Tx, p1 func(tx UnconfirmedTxn) bool) ([]UnconfirmedTxn, error) {

	ret := m.Called(p0, p1)

	var r0 []UnconfirmedTxn
	switch res := ret.Get(0).(type) {
	case nil:
	case []UnconfirmedTxn:
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

// GetUnknown mocked method
func (m *UnconfirmedTxnPoolerMock) GetUnknown(p0 *dbutil.Tx, p1 []cipher.SHA256) ([]cipher.SHA256, error) {

	ret := m.Called(p0, p1)

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

// GetUnspentsOfAddr mocked method
func (m *UnconfirmedTxnPoolerMock) GetUnspentsOfAddr(p0 *dbutil.Tx, p1 cipher.Address) (coin.UxArray, error) {

	ret := m.Called(p0, p1)

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

// InjectTransaction mocked method
func (m *UnconfirmedTxnPoolerMock) InjectTransaction(p0 *dbutil.Tx, p1 Blockchainer, p2 coin.Transaction, p3 int) (bool, *ErrTxnViolatesSoftConstraint, error) {

	ret := m.Called(p0, p1, p2, p3)

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

// Len mocked method
func (m *UnconfirmedTxnPoolerMock) Len(p0 *dbutil.Tx) (uint64, error) {

	ret := m.Called(p0)

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

// RawTxns mocked method
func (m *UnconfirmedTxnPoolerMock) RawTxns(p0 *dbutil.Tx) (coin.Transactions, error) {

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

// RecvOfAddresses mocked method
func (m *UnconfirmedTxnPoolerMock) RecvOfAddresses(p0 *dbutil.Tx, p1 coin.BlockHeader, p2 []cipher.Address) (coin.AddressUxOuts, error) {

	ret := m.Called(p0, p1, p2)

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

// Refresh mocked method
func (m *UnconfirmedTxnPoolerMock) Refresh(p0 *dbutil.Tx, p1 Blockchainer, p2 int) ([]cipher.SHA256, error) {

	ret := m.Called(p0, p1, p2)

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

// RemoveInvalid mocked method
func (m *UnconfirmedTxnPoolerMock) RemoveInvalid(p0 *dbutil.Tx, p1 Blockchainer) ([]cipher.SHA256, error) {

	ret := m.Called(p0, p1)

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

// RemoveTransactions mocked method
func (m *UnconfirmedTxnPoolerMock) RemoveTransactions(p0 *dbutil.Tx, p1 []cipher.SHA256) error {

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

// SetTxnsAnnounced mocked method
func (m *UnconfirmedTxnPoolerMock) SetTxnsAnnounced(p0 *dbutil.Tx, p1 map[cipher.SHA256]int64) error {

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
