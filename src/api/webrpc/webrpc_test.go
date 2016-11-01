package webrpc

import (
	"encoding/json"

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
	return []visor.ReadableOutput{}
}

func (fg fakeGateway) GetUnspentByHashes(hashes []string) []visor.ReadableOutput {
	return []visor.ReadableOutput{}
}
