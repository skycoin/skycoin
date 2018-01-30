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

// BlockchainerMock mock
type BlockchainerMock struct {
	mock.Mock
}

func NewBlockchainerMock() *BlockchainerMock {
	return &BlockchainerMock{}
}

// BindListener mocked method
func (m *BlockchainerMock) BindListener(p0 BlockListener) {

	m.Called(p0)

}

// ExecuteBlockWithTx mocked method
func (m *BlockchainerMock) ExecuteBlockWithTx(p0 *bolt.Tx, p1 *coin.SignedBlock) error {

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
func (m *BlockchainerMock) GetBlockByHash(p0 cipher.SHA256) (*coin.SignedBlock, error) {

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
func (m *BlockchainerMock) GetBlockBySeq(p0 uint64) (*coin.SignedBlock, error) {

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
func (m *BlockchainerMock) GetBlocks(p0 uint64, p1 uint64) []coin.SignedBlock {

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
func (m *BlockchainerMock) GetGenesisBlock() *coin.SignedBlock {

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
func (m *BlockchainerMock) GetLastBlocks(p0 uint64) []coin.SignedBlock {

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
func (m *BlockchainerMock) Head() (*coin.SignedBlock, error) {

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
func (m *BlockchainerMock) HeadSeq() uint64 {

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
func (m *BlockchainerMock) Len() uint64 {

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
func (m *BlockchainerMock) NewBlock(p0 coin.Transactions, p1 uint64) (*coin.Block, error) {

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
func (m *BlockchainerMock) Notify(p0 coin.Block) {

	m.Called(p0)

}

// Time mocked method
func (m *BlockchainerMock) Time() uint64 {

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
func (m *BlockchainerMock) TransactionFee(p0 *coin.Transaction) (uint64, error) {

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
func (m *BlockchainerMock) Unspent() blockdb.UnspentPool {

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
func (m *BlockchainerMock) UpdateDB(p0 func(tx *bolt.Tx) error) error {

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

// VerifySingleTxnAllConstraints mocked method
func (m *BlockchainerMock) VerifySingleTxnAllConstraints(p0 coin.Transaction, p1 int) error {

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

// VerifySingleTxnHardConstraints mocked method
func (m *BlockchainerMock) VerifySingleTxnHardConstraints(p0 coin.Transaction) error {

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

// VerifyBlockTxnConstraints mocked method
func (m *BlockchainerMock) VerifyBlockTxnConstraints(p0 coin.Transaction) error {

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
