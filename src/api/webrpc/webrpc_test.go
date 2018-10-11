package webrpc

import (
	"encoding/hex"
	"errors"
	"testing"

	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

func setupWebRPC(t *testing.T) *WebRPC {
	rpc, err := New(&fakeGateway{})
	require.NoError(t, err)
	return rpc
}

type fakeGateway struct {
	transactions         map[string]string
	injectRawTxMap       map[string]bool // key: transaction hash, value indicates whether the injectTransaction should return error.
	injectedTransactions map[string]string
	uxouts               []coin.UxOut
}

func (fg fakeGateway) GetLastBlocks(num uint64) ([]coin.SignedBlock, error) { // nolint: unparam
	return makeTestBlocksWithErr()
}

func (fg fakeGateway) GetBlocksInRange(start, end uint64) ([]coin.SignedBlock, error) {
	if start > end {
		return nil, nil
	}

	return makeTestBlocksWithErr()
}

func (fg fakeGateway) GetBlocks(vs []uint64) ([]coin.SignedBlock, error) {
	return nil, nil
}

func (fg fakeGateway) GetUnspentOutputsSummary(filters []visor.OutputsFilter) (*visor.UnspentOutputsSummary, error) {
	outs := []coin.UxOut{}
	for _, f := range filters {
		outs = f(fg.uxouts)
	}

	headTime := uint64(time.Now().UTC().Unix())

	rbOuts, err := visor.NewUnspentOutputs(outs, headTime)
	if err != nil {
		return nil, err
	}

	return &visor.UnspentOutputsSummary{
		HeadBlock: &coin.SignedBlock{},
		Confirmed: rbOuts,
	}, nil
}

func (fg fakeGateway) GetTransaction(txid cipher.SHA256) (*visor.Transaction, error) {
	str, ok := fg.transactions[txid.Hex()]
	if ok {
		return decodeRawTransaction(str), nil
	}
	return nil, nil
}

func (fg *fakeGateway) InjectBroadcastTransaction(txn coin.Transaction) error {
	if _, v := fg.injectRawTxMap[txn.Hash().Hex()]; v {
		if fg.injectedTransactions == nil {
			fg.injectedTransactions = make(map[string]string)
		}
		fg.injectedTransactions[txn.Hash().Hex()] = hex.EncodeToString(txn.Serialize())
		return nil
	}

	return errors.New("fake gateway inject transaction failed")
}

func (fg fakeGateway) GetSpentOutputsForAddresses(addr []cipher.Address) ([][]historydb.UxOut, error) {
	return make([][]historydb.UxOut, len(addr)), nil
}

func Test_rpcHandler_HandlerFunc(t *testing.T) {
	rpc := setupWebRPC(t)
	err := rpc.HandleFunc("get_status", getStatusHandler)
	testutil.RequireError(t, err, "get_status method already exist")
}
