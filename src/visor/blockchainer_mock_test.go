/*
* CODE GENERATED AUTOMATICALLY WITH github.com/ernesto-jimenez/goautomock
* THIS FILE MUST NEVER BE EDITED MANUALLY
 */

package visor

import (
	"fmt"

	mock "github.com/stretchr/testify/mock"

	bolt "github.com/boltdb/bolt"
	cipher "github.com/skycoin/skycoin/src/cipher"
	coin "github.com/skycoin/skycoin/src/coin"
	blockdb "github.com/skycoin/skycoin/src/visor/blockdb"
)

// blockchainerMock mock
type blockchainerMock struct {
	mock.Mock
}

func NewBlockchainerMock() *blockchainerMock {
	return &blockchainerMock{}
}

// BindListener mocked method
func (m *blockchainerMock) BindListener(p0 BlockListener) {

	m.Called(p0)

}

// ExecuteBlockWithTx mocked method
func (m *blockchainerMock) ExecuteBlockWithTx(p0 *bolt.Tx, p1 *coin.SignedBlock) error {

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

// GetBlockByHash mocked method
func (m *blockchainerMock) GetBlockByHash(p0 cipher.SHA256) (*coin.SignedBlock, error) {

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
func (m *blockchainerMock) GetBlockBySeq(p0 uint64) (*coin.SignedBlock, error) {

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

// GetBlocks mocked method
func (m *blockchainerMock) GetBlocks(p0 uint64, p1 uint64) []coin.SignedBlock {

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

// GetGenesisBlock mocked method
func (m *blockchainerMock) GetGenesisBlock() *coin.SignedBlock {

	ret := m.Called()

	var r0 *coin.SignedBlock
	switch res := ret.Get(0).(type) {
	case nil:
	case *coin.SignedBlock:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetLastBlocks mocked method
func (m *blockchainerMock) GetLastBlocks(p0 uint64) []coin.SignedBlock {

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

// Head mocked method
func (m *blockchainerMock) Head() (*coin.SignedBlock, error) {

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

// HeadSeq mocked method
func (m *blockchainerMock) HeadSeq() uint64 {

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

// Len mocked method
func (m *blockchainerMock) Len() uint64 {

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

// NewBlock mocked method
func (m *blockchainerMock) NewBlock(p0 coin.Transactions, p1 uint64) (*coin.Block, error) {

	ret := m.Called(p0, p1)

	var r0 *coin.Block
	switch res := ret.Get(0).(type) {
	case nil:
	case *coin.Block:
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

// Notify mocked method
func (m *blockchainerMock) Notify(p0 coin.Block) {

	m.Called(p0)

}

// Time mocked method
func (m *blockchainerMock) Time() uint64 {

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

// TransactionFee mocked method
func (m *blockchainerMock) TransactionFee(p0 *coin.Transaction) (uint64, error) {

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

// Unspent mocked method
func (m *blockchainerMock) Unspent() blockdb.UnspentPool {

	ret := m.Called()

	var r0 blockdb.UnspentPool
	switch res := ret.Get(0).(type) {
	case nil:
	case blockdb.UnspentPool:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// UpdateDB mocked method
func (m *blockchainerMock) UpdateDB(p0 func(tx *bolt.Tx) error) error {

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

// VerifyTransaction mocked method
func (m *blockchainerMock) VerifyTransaction(p0 coin.Transaction) error {

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
