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
)

// UnspentPoolerMock mock
type UnspentPoolerMock struct {
	mock.Mock
}

func NewUnspentPoolerMock() *UnspentPoolerMock {
	return &UnspentPoolerMock{}
}

// AddressCount mocked method
func (m *UnspentPoolerMock) AddressCount(p0 *dbutil.Tx) (uint64, error) {

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

// Contains mocked method
func (m *UnspentPoolerMock) Contains(p0 *dbutil.Tx, p1 cipher.SHA256) (bool, error) {

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

// Get mocked method
func (m *UnspentPoolerMock) Get(p0 *dbutil.Tx, p1 cipher.SHA256) (*coin.UxOut, error) {

	ret := m.Called(p0, p1)

	var r0 *coin.UxOut
	switch res := ret.Get(0).(type) {
	case nil:
	case *coin.UxOut:
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

// GetAll mocked method
func (m *UnspentPoolerMock) GetAll(p0 *dbutil.Tx) (coin.UxArray, error) {

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

// GetArray mocked method
func (m *UnspentPoolerMock) GetArray(p0 *dbutil.Tx, p1 []cipher.SHA256) (coin.UxArray, error) {

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

// GetUnspentsOfAddrs mocked method
func (m *UnspentPoolerMock) GetUnspentsOfAddrs(p0 *dbutil.Tx, p1 []cipher.Address) (coin.AddressUxOuts, error) {

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

// GetUxHash mocked method
func (m *UnspentPoolerMock) GetUxHash(p0 *dbutil.Tx) (cipher.SHA256, error) {

	ret := m.Called(p0)

	var r0 cipher.SHA256
	switch res := ret.Get(0).(type) {
	case nil:
	case cipher.SHA256:
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
func (m *UnspentPoolerMock) Len(p0 *dbutil.Tx) (uint64, error) {

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

// MaybeBuildIndexes mocked method
func (m *UnspentPoolerMock) MaybeBuildIndexes(p0 *dbutil.Tx, p1 uint64) error {

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

// ProcessBlock mocked method
func (m *UnspentPoolerMock) ProcessBlock(p0 *dbutil.Tx, p1 *coin.SignedBlock) error {

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
