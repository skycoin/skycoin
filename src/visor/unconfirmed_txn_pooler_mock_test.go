/*
* CODE GENERATED AUTOMATICALLY WITH github.com/ernesto-jimenez/goautomock
* THIS FILE MUST NEVER BE EDITED MANUALLY
 */

package visor

import (
	"fmt"

	mock "github.com/stretchr/testify/mock"

	time "time"

	bolt "github.com/boltdb/bolt"
	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	blockdb "github.com/skycoin/skycoin/src/visor/blockdb"
)

// UnconfirmedTxnPoolerMock mock
type UnconfirmedTxnPoolerMock struct {
	mock.Mock
}

func NewUnconfirmedTxnPoolerMock() *UnconfirmedTxnPoolerMock {
	return &UnconfirmedTxnPoolerMock{}
}

// FilterKnown mocked method
func (m *UnconfirmedTxnPoolerMock) FilterKnown(p0 []cipher.SHA256) []cipher.SHA256 {

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

// ForEach mocked method
func (m *UnconfirmedTxnPoolerMock) ForEach(p0 func(cipher.SHA256, *UnconfirmedTxn) error) error {

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

// Get mocked method
func (m *UnconfirmedTxnPoolerMock) Get(p0 cipher.SHA256) (*UnconfirmedTxn, bool) {

	ret := m.Called(p0)

	var r0 *UnconfirmedTxn
	switch res := ret.Get(0).(type) {
	case nil:
	case *UnconfirmedTxn:
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

// GetIncomingOutputs mocked method
func (m *UnconfirmedTxnPoolerMock) GetIncomingOutputs(p0 coin.BlockHeader) coin.UxArray {

	ret := m.Called(p0)

	var r0 coin.UxArray
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.UxArray:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetKnown mocked method
func (m *UnconfirmedTxnPoolerMock) GetKnown(p0 []cipher.SHA256) coin.Transactions {

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

// GetSpendingOutputs mocked method
func (m *UnconfirmedTxnPoolerMock) GetSpendingOutputs(p0 blockdb.UnspentPool) (coin.UxArray, error) {

	ret := m.Called(p0)

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

// GetTxHashes mocked method
func (m *UnconfirmedTxnPoolerMock) GetTxHashes(p0 func(tx UnconfirmedTxn) bool) []cipher.SHA256 {

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

// GetTxns mocked method
func (m *UnconfirmedTxnPoolerMock) GetTxns(p0 func(tx UnconfirmedTxn) bool) []UnconfirmedTxn {

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

// GetUnspentsOfAddr mocked method
func (m *UnconfirmedTxnPoolerMock) GetUnspentsOfAddr(p0 cipher.Address) coin.UxArray {

	ret := m.Called(p0)

	var r0 coin.UxArray
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.UxArray:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// InjectTxn mocked method
func (m *UnconfirmedTxnPoolerMock) InjectTxn(p0 Blockchainer, p1 coin.Transaction) (bool, error) {

	ret := m.Called(p0, p1)

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

// Len mocked method
func (m *UnconfirmedTxnPoolerMock) Len() int {

	ret := m.Called()

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

// RawTxns mocked method
func (m *UnconfirmedTxnPoolerMock) RawTxns() coin.Transactions {

	ret := m.Called()

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

// RecvOfAddresses mocked method
func (m *UnconfirmedTxnPoolerMock) RecvOfAddresses(p0 coin.BlockHeader, p1 []cipher.Address) (coin.AddressUxOuts, error) {

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

// Refresh mocked method
func (m *UnconfirmedTxnPoolerMock) Refresh(p0 Blockchainer, p1 int) ([]cipher.SHA256, error) {

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

// RemoveInvalid mocked method
func (m *UnconfirmedTxnPoolerMock) RemoveInvalid(p0 Blockchainer) ([]cipher.SHA256, error) {

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

// RemoveTransactions mocked method
func (m *UnconfirmedTxnPoolerMock) RemoveTransactions(p0 []cipher.SHA256) error {

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

// RemoveTransactionsWithTx mocked method
func (m *UnconfirmedTxnPoolerMock) RemoveTransactionsWithTx(p0 *bolt.Tx, p1 []cipher.SHA256) {

	m.Called(p0, p1)

}

// SetAnnounced mocked method
func (m *UnconfirmedTxnPoolerMock) SetAnnounced(p0 cipher.SHA256, p1 time.Time) error {

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

// SpendsOfAddresses mocked method
func (m *UnconfirmedTxnPoolerMock) SpendsOfAddresses(p0 []cipher.Address, p1 blockdb.UnspentGetter) (coin.AddressUxOuts, error) {

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

// InjectTransaction mocked method
func (m *UnconfirmedTxnPoolerMock) InjectTransaction(p0 Blockchainer, p1 coin.Transaction, p2 int) (bool, *ErrTxnViolatesSoftConstraint, error) {

	ret := m.Called(p0, p1, p2)

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
