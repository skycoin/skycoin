package coin

import (
	"fmt"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
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

func testAddBlock(t *testing.T, cases []blockCase) {
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
	}

	testCase(t, testData)
}

func testCase(t *testing.T, cases []blockCase) {
	_, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}

	defer teardown()

	btree := NewBlockTree()
	blocks := make([]Block, len(cases))
	for i, d := range cases {
		var preHash cipher.SHA256
		if d.BInfo.Pre != -1 {
			preHash = blocks[d.BInfo.Pre].HashHeader()
		}

		b := Block{
			Head: BlockHeader{
				BkSeq:    d.BInfo.Seq,
				Time:     d.BInfo.Time,
				Fee:      d.BInfo.Fee,
				PrevHash: preHash,
			},
		}
		blocks[i] = b

		switch d.Action {
		case "add":
			err := btree.AddBlock(b)
			if err != d.Err {
				t.Fatal(fmt.Sprintf("expect err:%v, but get err:%v", d.Err, err))
			}

			if err == nil {
				b1 := btree.GetBlock(b.HashHeader())
				assert.Equal(t, *b1, b)
			}
		case "remove":
			err := btree.RemoveBlock(b)
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

// func TestAddBlock(t *testing.T) {
// 	_, teardown, err := setup(t)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	defer teardown()

// 	btree := NewBlockTree()

// 	assert.NotNil(t, btree.blocks)
// 	assert.NotNil(t, btree.tree)

// 	b := Block{
// 		Head: BlockHeader{
// 			BkSeq: 0,
// 		},
// 	}
// 	if err := btree.AddBlock(b); err != nil {
// 		t.Fatal(err)
// 	}

// 	// check block in bucket
// 	key := b.HashHeader()
// 	bin := btree.blocks.Get(key[:])
// 	assert.NotNil(t, bin)
// 	b1 := Block{}
// 	if err := encoder.DeserializeRaw(bin, &b1); err != nil {
// 		t.Fatal(err)
// 	}

// 	assert.Equal(t, b, b1)

// 	// get block hash pair in depth.
// 	depKey := itob(0)
// 	pairsBin := btree.tree.Get(depKey)
// 	pairs := []HashPair{}
// 	if err := encoder.DeserializeRaw(pairsBin, &pairs); err != nil {
// 		t.Fatal(err)
// 	}
// 	assert.Equal(t, len(pairs), 1)
// 	assert.Equal(t, pairs[0].Hash, b.HashHeader())
// 	assert.Equal(t, pairs[0].PreHash, cipher.SHA256{})

// 	// add one block in depth 1
// 	b2 := Block{
// 		Head: BlockHeader{
// 			BkSeq:    1,
// 			PrevHash: b.HashHeader(),
// 		},
// 	}

// 	if err := btree.AddBlock(b2); err != nil {
// 		t.Fatal(err)
// 	}
// 	b2Dep := itob(1)
// 	b2PairsBin := btree.tree.Get(b2Dep)
// 	b2Pairs := []HashPair{}
// 	if err := encoder.DeserializeRaw(b2PairsBin, &b2Pairs); err != nil {
// 		t.Fatal(err)
// 	}
// 	assert.Equal(t, b2Pairs[0].Hash, b2.HashHeader())
// 	assert.Equal(t, b2Pairs[0].PreHash, b.HashHeader())

// 	// add one more block in depth 1
// 	b3 := Block{
// 		Head: BlockHeader{
// 			BkSeq:    1,
// 			PrevHash: b.HashHeader(),
// 			Fee:      100,
// 		},
// 	}

// 	if err := btree.AddBlock(b3); err != nil {
// 		t.Fatal(err)
// 	}

// 	// get block pairs.
// 	b3PairsBin := btree.tree.Get(b2Dep)
// 	b3Pairs := []HashPair{}
// 	if err := encoder.DeserializeRaw(b3PairsBin, &b3Pairs); err != nil {
// 		t.Fatal(err)
// 	}

// 	assert.Equal(t, len(b3Pairs), 2)
// 	assert.Equal(t, b3Pairs[1].Hash, b3.HashHeader())
// 	assert.Equal(t, b3Pairs[1].PreHash, b.HashHeader())
// }

// func TestRemoveBlock(t *testing.T) {
// 	_, teardown, err := setup(t)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	defer teardown()

// 	btree := NewBlockTree()
// 	b := Block{
// 		Head: BlockHeader{
// 			BkSeq: 0,
// 		},
// 	}

// 	if err := btree.AddBlock(b); err != nil {
// 		t.Fatal(err)
// 	}

//     b1 := Block {
//         Head: BlockHeader {
//             BkSeq: 1,
//         }
//     }

// 	if err := btree.RemoveBlock(b); err != nil {
// 		t.Fatal(err)
// 	}

// 	bKey := b.HashHeader()
// 	v := btree.blocks.Get(bKey[:])
// 	if v != nil {
// 		t.Fatal("remove block failed")
// 	}

// 	bDep := itob(0)
// 	v = btree.tree.Get(bDep)
// 	if v != nil {
// 		t.Fatal("remove block in tree failed")
// 	}
// }

// type prepareBlock struct {
//     ID uint64
//     BkSeq uint64
//     Time uint64
// }

// type removeCases struct {
//    PrepareBlocks []prepareBlock
//    RemoveID uint64
// }

// func removeBlockTestCases(t *testing.T, cases )
