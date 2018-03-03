/*
* CODE GENERATED AUTOMATICALLY WITH github.com/ernesto-jimenez/goautomock
* THIS FILE MUST NEVER BE EDITED MANUALLY
 */

package blockdb

import (
	"fmt"

	mock "github.com/stretchr/testify/mock"

	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	bucket "github.com/skycoin/skycoin/src/visor/bucket"
)

// UnspentPoolMock mock
type UnspentPoolMock struct {
	mock.Mock
}

func NewUnspentPoolMock() *UnspentPoolMock {
	return &UnspentPoolMock{}
}

// Contains mocked method
func (m *UnspentPoolMock) Contains(p0 cipher.SHA256) bool {

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

// Get mocked method
func (m *UnspentPoolMock) Get(p0 cipher.SHA256) (coin.UxOut, bool) {

	ret := m.Called(p0)

	var r0 coin.UxOut
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.UxOut:
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

// GetAll mocked method
func (m *UnspentPoolMock) GetAll() (coin.UxArray, error) {

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

// GetArray mocked method
func (m *UnspentPoolMock) GetArray(p0 []cipher.SHA256) (coin.UxArray, error) {

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

// GetUnspentsOfAddrs mocked method
func (m *UnspentPoolMock) GetUnspentsOfAddrs(p0 []cipher.Address) coin.AddressUxOuts {

	ret := m.Called(p0)

	var r0 coin.AddressUxOuts
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.AddressUxOuts:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetUxHash mocked method
func (m *UnspentPoolMock) GetUxHash() cipher.SHA256 {

	ret := m.Called()

	var r0 cipher.SHA256
	switch res := ret.Get(0).(type) {
	case nil:
	case cipher.SHA256:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// Len mocked method
func (m *UnspentPoolMock) Len() uint64 {

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

// ProcessBlock mocked method
func (m *UnspentPoolMock) ProcessBlock(p0 *coin.SignedBlock) bucket.TxHandler {

	ret := m.Called(p0)

	var r0 bucket.TxHandler
	switch res := ret.Get(0).(type) {
	case nil:
	case bucket.TxHandler:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}
