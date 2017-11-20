package visor

import (
	"fmt"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

// ParserOption option type which will be used when creating parser instance
type ParserOption func(*BlockchainParser)

// BlockchainParser parses the blockchain and stores the data into historydb.
type BlockchainParser struct {
	db        *dbutil.DB
	historyDB *historydb.HistoryDB
	blkC      chan coin.Block
	quit      chan struct{}
	done      chan struct{}
	bc        *Blockchain

	isStart bool
}

// NewBlockchainParser create and init the parser instance.
func NewBlockchainParser(db *dbutil.DB, hisDB *historydb.HistoryDB, bc *Blockchain) *BlockchainParser {
	return &BlockchainParser{
		db:        db,
		bc:        bc,
		historyDB: hisDB,
		quit:      make(chan struct{}),
		done:      make(chan struct{}),
		blkC:      make(chan coin.Block, 10),
	}
}

// FeedBlock feeds block to the parser
func (bcp *BlockchainParser) FeedBlock(b coin.Block) {
	bcp.blkC <- b
}

// Init initializes blockchain parser
func (bcp *BlockchainParser) Init(tx *bolt.Tx) error {
	logger.Info("Blockchain parser initializing")
	defer logger.Info("Blockchain parser initialization completed")

	if err := bcp.historyDB.ResetIfNeed(tx); err != nil {
		return err
	}

	// parse to the blockchain head
	headSeq, ok, err := bcp.bc.HeadSeq(tx)
	if err != nil {
		return err
	} else if !ok {
		return nil
	}

	return bcp.parseTo(tx, headSeq)
}

// Run starts blockchain parser
func (bcp *BlockchainParser) Run() error {
	logger.Info("Blockchain parser start")
	defer logger.Info("Blockchain parser closed")
	defer close(bcp.done)

	for {
		select {
		case <-bcp.quit:
			return nil
		case b := <-bcp.blkC:
			if err := bcp.db.Update(func(tx *bolt.Tx) error {
				return bcp.historyDB.ParseBlock(tx, &b)
			}); err != nil {
				logger.Errorf("BlockchainParser.historyDB.ParseBlock failed: %v", err)
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
