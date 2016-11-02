package webrpc

import (
	"encoding/json"
	"errors"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
)

func setup() (*rpcHandler, func()) {
	c := make(chan struct{})
	f := func() {
		close(c)
	}

	return makeRPC(1, 1, &fakeGateway{}, c), f
}

type fakeGateway struct {
	transactions   map[string]string
	injectRawTxMap map[string]bool // key: transacion hash, value indicates whether the injectTransaction should return error.
}

func (fg fakeGateway) GetLastBlocks(num uint64) *visor.ReadableBlocks {
	var blocks visor.ReadableBlocks
	if err := json.Unmarshal([]byte(blockString), &blocks); err != nil {
		panic(err)
	}

	return &blocks
}

func (fg fakeGateway) GetBlocks(start, end uint64) *visor.ReadableBlocks {
	var blocks visor.ReadableBlocks
	if start > end {
		return &blocks
	}

	if err := json.Unmarshal([]byte(blockString), &blocks); err != nil {
		panic(err)
	}

	return &blocks
}

func (fg fakeGateway) GetUnspentByAddrs(addrs []string) []visor.ReadableOutput {
	addrMap := make(map[string]bool)
	for _, a := range addrs {
		addrMap[a] = true
	}

	return filterOut(decodeOutputStr(outputStr), func(out visor.ReadableOutput) bool {
		_, ok := addrMap[out.Address]
		return ok
	})
}

func (fg fakeGateway) GetUnspentByHashes(hashes []string) []visor.ReadableOutput {
	return []visor.ReadableOutput{}
}

func (fg fakeGateway) GetTransaction(txid cipher.SHA256) (*visor.TransactionResult, error) {
	str, ok := fg.transactions[txid.Hex()]
	if ok {
		return decodeTransaction(str), nil
	}
	return nil, nil
}

func (fg fakeGateway) InjectTransaction(txn coin.Transaction) (coin.Transaction, error) {
	if _, v := fg.injectRawTxMap[txn.Hash().Hex()]; v {
		return txn, nil
	}

	return txn, errors.New("inject transaction failed")
}
