package webrpc

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"

	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
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
	addrRecvUxOuts       []*historydb.UxOut
	addrSpentUxOUts      []*historydb.UxOut
	uxouts               []coin.UxOut
}

func (fg fakeGateway) GetLastBlocks(num uint64) (*visor.ReadableBlocks, error) { // nolint: unparam
	var blocks visor.ReadableBlocks
	if err := json.Unmarshal([]byte(blockString), &blocks); err != nil {
		return nil, err
	}

	return &blocks, nil
}

func (fg fakeGateway) GetBlocks(start, end uint64) (*visor.ReadableBlocks, error) {
	var blocks visor.ReadableBlocks
	if start > end {
		return nil, nil
	}

	if err := json.Unmarshal([]byte(blockString), &blocks); err != nil {
		return nil, err
	}

	return &blocks, nil
}

func (fg fakeGateway) GetBlocksInDepth(vs []uint64) (*visor.ReadableBlocks, error) {
	return nil, nil
}

func (fg fakeGateway) GetUnspentOutputs(filters ...daemon.OutputsFilter) (*visor.ReadableOutputSet, error) {
	outs := []coin.UxOut{}
	for _, f := range filters {
		outs = f(fg.uxouts)
	}

	headTime := uint64(time.Now().UTC().Unix())

	rbOuts, err := visor.NewReadableOutputs(headTime, outs)
	if err != nil {
		return nil, err
	}

	return &visor.ReadableOutputSet{
		HeadOutputs: rbOuts,
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

func (fg fakeGateway) GetAddrUxOuts(addr []cipher.Address) ([]*historydb.UxOut, error) {
	return nil, nil
}

func (fg fakeGateway) GetTimeNow() uint64 {
	return 0
}

func Test_rpcHandler_HandlerFunc(t *testing.T) {
	rpc := setupWebRPC(t)
	rpc.HandleFunc("get_status", getStatusHandler)
	err := rpc.HandleFunc("get_status", getStatusHandler)
	require.Error(t, err)
}
