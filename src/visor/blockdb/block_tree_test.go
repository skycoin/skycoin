package blockdb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
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
	db, close := prepareDB(t)
	defer close()

	btree := &blockTree{}
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

		err := db.Update("", func(tx *dbutil.Tx) error {
			switch d.Action {
			case "add":
				err := btree.AddBlock(tx, &b)
				require.Equal(t, d.Err, err, "expect err:%v, but get err:%v", d.Err, err)

				if err == nil {
					b1, err := btree.GetBlock(tx, b.HashHeader())
					require.NoError(t, err)
					require.Equal(t, b, *b1)
				}
			case "remove":
				err := btree.RemoveBlock(tx, &b)
				require.Equal(t, d.Err, err, "expect err:%v, but get err:%v", d.Err, err)
				if err == nil {
					b1, err := btree.GetBlock(tx, b.HashHeader())
					require.NoError(t, err)
					require.Nil(t, b1)
				}
			}

			return nil
		})

		require.NoError(t, err)
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
	db, teardown := prepareDB(t)
	defer teardown()

	bc := &blockTree{}
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

	err := db.Update("", func(tx *dbutil.Tx) error {
		err := bc.AddBlock(tx, &blocks[0])
		require.NoError(t, err)

		blocks[1].Head.PrevHash = blocks[0].HashHeader()
		err = bc.AddBlock(tx, &blocks[1])
		require.NoError(t, err)

		blocks[2].Head.PrevHash = blocks[0].HashHeader()
		err = bc.AddBlock(tx, &blocks[2])
		require.NoError(t, err)

		return nil
	})

	require.NoError(t, err)

	var block *coin.Block
	err = db.View("", func(tx *dbutil.Tx) error {
		var err error
		block, err = bc.GetBlockInDepth(tx, 1, func(tx *dbutil.Tx, hps []coin.HashPair) (cipher.SHA256, bool) {
			for _, hp := range hps {
				b, err := bc.GetBlock(tx, hp.Hash)
				require.NoError(t, err)
				if b.Time() == 2 {
					return b.HashHeader(), true
				}
			}
			return cipher.SHA256{}, false
		})
		require.NoError(t, err)
		return err
	})

	require.NoError(t, err)

	require.NotNil(t, block)
	require.Equal(t, blocks[2], *block)
}
