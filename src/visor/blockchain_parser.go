package visor

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

// ParserOption option type which will be used when creating parser instance
type ParserOption func(*BlockchainParser)

// BlockchainParser parses the blockchain and stores the data into historydb.
type BlockchainParser struct {
	historyDB *historydb.HistoryDB
	blkC      chan coin.Block
	quit      chan struct{}
	done      chan struct{}
	bc        *Blockchain

	isStart bool
}

// NewBlockchainParser create and init the parser instance.
func NewBlockchainParser(hisDB *historydb.HistoryDB, bc *Blockchain, ops ...ParserOption) *BlockchainParser {
	bp := &BlockchainParser{
		bc:        bc,
		historyDB: hisDB,
		quit:      make(chan struct{}),
		done:      make(chan struct{}),
		blkC:      make(chan coin.Block, 10),
	}

	for _, op := range ops {
		op(bp)
	}

	return bp
}

// FeedBlock feeds block to the parser
func (bcp *BlockchainParser) FeedBlock(b coin.Block) {
	bcp.blkC <- b
}

// Run starts blockchain parser
func (bcp *BlockchainParser) Run(tx *bolt.Tx) error {
	logger.Info("Blockchain parser start")
	defer logger.Info("Blockchain parser closed")
	defer close(bcp.done)

	if err := bcp.historyDB.ResetIfNeed(tx); err != nil {
		return err
	}

	// parse to the blockchain head
	headSeq := bcp.bc.HeadSeq()
	if err := bcp.parseTo(tx, headSeq); err != nil {
		return err
	}

	for {
		select {
		case <-bcp.quit:
			return nil
		case b := <-bcp.blkC:
			if err := bcp.historyDB.ParseBlock(tx, &b); err != nil {
				return err
			}
		}
	}
}

// Shutdown close the block parsing process.
func (bcp *BlockchainParser) Shutdown() {
	close(bcp.quit)
	<-bcp.done
}

func (bcp *BlockchainParser) parseTo(tx *bolt.Tx, bcHeight uint64) error {
	parsedHeight, err := bcp.historyDB.ParsedHeight(tx)
	if err != nil {
		return err
	}

	for i := int64(0); i < int64(bcHeight)-parsedHeight; i++ {
		b, err := bcp.bc.store.GetSignedBlockBySeq(tx, uint64(parsedHeight+i+1))
		if err != nil {
			return err
		}

		if b == nil {
			return fmt.Errorf("no block exist in depth:%d", parsedHeight+i+1)
		}

		if err := bcp.historyDB.ParseBlock(tx, &b.Block); err != nil {
			return err
		}
	}

	return nil
}
