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
	blockdb "github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

// BlockchainerMock mock
type BlockchainerMock struct {
	mock.Mock
}

func NewBlockchainerMock() *BlockchainerMock {
	return &BlockchainerMock{}
}

// ExecuteBlock mocked method
func (m *BlockchainerMock) ExecuteBlock(p0 *dbutil.Tx, p1 *coin.SignedBlock) error {

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

// GetBlocks mocked method
func (m *BlockchainerMock) GetBlocks(p0 *dbutil.Tx, p1 uint64, p2 uint64) ([]coin.SignedBlock, error) {

	ret := m.Called(p0, p1, p2)

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

// GetGenesisBlock mocked method
func (m *BlockchainerMock) GetGenesisBlock(p0 *dbutil.Tx) (*coin.SignedBlock, error) {

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

// GetLastBlocks mocked method
func (m *BlockchainerMock) GetLastBlocks(p0 *dbutil.Tx, p1 uint64) ([]coin.SignedBlock, error) {

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

// GetSignedBlockByHash mocked method
func (m *BlockchainerMock) GetSignedBlockByHash(p0 *dbutil.Tx, p1 cipher.SHA256) (*coin.SignedBlock, error) {

	ret := m.Called(p0, p1)

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

// GetSignedBlockBySeq mocked method
func (m *BlockchainerMock) GetSignedBlockBySeq(p0 *dbutil.Tx, p1 uint64) (*coin.SignedBlock, error) {

	ret := m.Called(p0, p1)

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

// Head mocked method
func (m *BlockchainerMock) Head(p0 *dbutil.Tx) (*coin.SignedBlock, error) {

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

// HeadSeq mocked method
func (m *BlockchainerMock) HeadSeq(p0 *dbutil.Tx) (uint64, bool, error) {

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

// Len mocked method
func (m *BlockchainerMock) Len(p0 *dbutil.Tx) (uint64, error) {

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

// NewBlock mocked method
func (m *BlockchainerMock) NewBlock(p0 *dbutil.Tx, p1 coin.Transactions, p2 uint64) (*coin.Block, error) {

	ret := m.Called(p0, p1, p2)

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

// Time mocked method
func (m *BlockchainerMock) Time(p0 *dbutil.Tx) (uint64, error) {

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

// TransactionFee mocked method
func (m *BlockchainerMock) TransactionFee(p0 *dbutil.Tx, p1 uint64) coin.FeeCalculator {

	ret := m.Called(p0, p1)

	var r0 coin.FeeCalculator
	switch res := ret.Get(0).(type) {
	case nil:
	case coin.FeeCalculator:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// Unspent mocked method
func (m *BlockchainerMock) Unspent() blockdb.UnspentPooler {

	ret := m.Called()

	var r0 blockdb.UnspentPooler
	switch res := ret.Get(0).(type) {
	case nil:
	case blockdb.UnspentPooler:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// VerifyBlockTxnConstraints mocked method
func (m *BlockchainerMock) VerifyBlockTxnConstraints(p0 *dbutil.Tx, p1 coin.Transaction) error {

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

// VerifySingleTxnSoftHardConstraints mocked method
func (m *BlockchainerMock) VerifySingleTxnSoftHardConstraints(p0 *dbutil.Tx, p1 coin.Transaction, p2 int) error {

	ret := m.Called(p0, p1, p2)

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
func (m *BlockchainerMock) VerifySingleTxnHardConstraints(p0 *dbutil.Tx, p1 coin.Transaction) error {

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
