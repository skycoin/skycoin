package blockdb

import (
	"fmt"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/stretchr/testify/assert"
)

type blockInfo struct {
	Seq  uint64
	Time uint64
	Fee  uint64
	Pre  int
}

type blockCase struct {
	BInfo  blockInfo
	Err    error
	Action string
}

func testCase(t *testing.T, cases []blockCase) {
	db, close := testutil.PrepareDB(t)
	defer close()

	btree, err := newBlockTree(db)
	assert.Nil(t, err)
	blocks := make([]coin.Block, len(cases))
	for i, d := range cases {
		var preHash cipher.SHA256
		if d.BInfo.Pre != -1 {
			preHash = blocks[d.BInfo.Pre].HashHeader()
		}

		b := coin.Block{
			Head: coin.BlockHeader{
				BkSeq:    d.BInfo.Seq,
				Time:     d.BInfo.Time,
				Fee:      d.BInfo.Fee,
				PrevHash: preHash,
			},
		}
		blocks[i] = b

		switch d.Action {
		case "add":
			err := btree.AddBlock(&b)
			if err != d.Err {
				t.Fatal(fmt.Sprintf("expect err:%v, but get err:%v", d.Err, err))
			}

			if err == nil {
				b1 := btree.GetBlock(b.HashHeader())
				assert.Equal(t, *b1, b)
			}
		case "remove":
			err := btree.RemoveBlock(&b)
			if err != d.Err {
				t.Fatal(fmt.Sprintf("expect err:%v, but get err:%v", d.Err, err))
			}
			if err == nil {
				b1 := btree.GetBlock(b.HashHeader())
				assert.Nil(t, b1)
			}
		}
	}
}

func TestAddBlock(t *testing.T) {
	testData := []blockCase{
		blockCase{
			BInfo:  blockInfo{Seq: 0, Time: 0, Fee: 0, Pre: -1},
			Err:    nil,
			Action: "add",
		},
		blockCase{
			BInfo:  blockInfo{Seq: 1, Time: 0, Fee: 0, Pre: 0},
			Err:    nil,
			Action: "add",
		},
		blockCase{
			BInfo:  blockInfo{Seq: 1, Time: 1, Fee: 0, Pre: 0},
			Err:    nil,
			Action: "add",
		},
		blockCase{
			BInfo:  blockInfo{Seq: 2, Time: 2, Fee: 0, Pre: 1},
			Err:    nil,
			Action: "add",
		},
		blockCase{
			BInfo:  blockInfo{Seq: 2, Time: 2, Fee: 0, Pre: 1},
			Err:    errBlockExist,
			Action: "add",
		},
		blockCase{
			BInfo:  blockInfo{Seq: 2, Time: 2, Fee: 0, Pre: 0},
			Err:    errWrongParent,
			Action: "add",
		},
		blockCase{
			BInfo:  blockInfo{Seq: 4, Time: 2, Fee: 0, Pre: 3},
			Err:    errWrongParent,
			Action: "add",
		},
		blockCase{
			BInfo:  blockInfo{Seq: 3, Time: 2, Fee: 0, Pre: -1},
			Err:    errNoParent,
			Action: "add",
		},
	}

	testCase(t, testData)
}

func TestRemoveBlock(t *testing.T) {
	testData := []blockCase{
		blockCase{
			BInfo:  blockInfo{Seq: 0, Time: 0, Fee: 0, Pre: -1},
			Err:    nil,
			Action: "add",
		},
		blockCase{
			BInfo:  blockInfo{Seq: 1, Time: 1, Fee: 0, Pre: 0},
			Err:    nil,
			Action: "add",
		},
		blockCase{
			BInfo:  blockInfo{Seq: 1, Time: 2, Fee: 0, Pre: 0},
			Err:    nil,
			Action: "add",
		},
		// remove block normally.
		blockCase{
			BInfo:  blockInfo{Seq: 1, Time: 2, Fee: 0, Pre: 0},
			Err:    nil,
			Action: "remove",
		},
		// remove genesis block, which has children.
		blockCase{
			BInfo:  blockInfo{Seq: 0, Time: 0, Fee: 0, Pre: -1},
			Err:    errHasChild,
			Action: "remove",
		},
		// remove the last block in depth 1.
		blockCase{
			BInfo:  blockInfo{Seq: 1, Time: 1, Fee: 0, Pre: 0},
			Err:    nil,
			Action: "remove",
		},
	}

	testCase(t, testData)
}

func TestGetBlockInDepth(t *testing.T) {
	db, teardown := testutil.PrepareDB(t)
	defer teardown()

	bc, err := newBlockTree(db)
	assert.Nil(t, err)
	blocks := []coin.Block{
		coin.Block{
			Head: coin.BlockHeader{
				BkSeq: 0,
				Time:  0,
				Fee:   0,
			},
		},
		coin.Block{
			Head: coin.BlockHeader{
				BkSeq: 1,
				Time:  1,
			},
		},
		coin.Block{
			Head: coin.BlockHeader{
				BkSeq: 1,
				Time:  2,
			},
		},
	}

	assert.Nil(t, bc.AddBlock(&blocks[0]))
	blocks[1].Head.PrevHash = blocks[0].HashHeader()
	assert.Nil(t, bc.AddBlock(&blocks[1]))
	blocks[2].Head.PrevHash = blocks[0].HashHeader()
	assert.Nil(t, bc.AddBlock(&blocks[2]))

	block := bc.GetBlockInDepth(1, func(hps []coin.HashPair) cipher.SHA256 {
		for _, hp := range hps {
			b := bc.GetBlock(hp.Hash)
			if b.Time() == 2 {
				return b.HashHeader()
			}
		}
		return cipher.SHA256{}
	})

	assert.Equal(t, *block, blocks[2])
}
