package visor

import (
	"fmt"

	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

// BlockchainParser parses the blockchain and stores the data into historydb.
type BlockchainParser struct {
	historyDB    *historydb.HistoryDB
	parsedHeight uint64
	blkC         chan coin.Block
	closing      chan chan struct{}
	bc           *Blockchain

	startC  chan struct{}
	isStart bool
}

// NewBlockchainParser create and init the parser instance.
func NewBlockchainParser(hisDB *historydb.HistoryDB, bc *Blockchain) *BlockchainParser {
	bp := &BlockchainParser{
		bc:           bc,
		historyDB:    hisDB,
		closing:      make(chan chan struct{}),
		blkC:         make(chan coin.Block, 10),
		startC:       make(chan struct{}),
		parsedHeight: 0,
	}

	bp.run()
	return bp
}

// Start the parsing process
func (bcp *BlockchainParser) Start() {
	bcp.startC <- struct{}{}
}

// BlockListener when new block appended to blockchain, this method will b invoked
func (bcp *BlockchainParser) BlockListener(b coin.Block) {
	bcp.blkC <- b
}

// Start starts blockchain parser
func (bcp *BlockchainParser) run() {
	go func() {
		for {
			select {
			case cc := <-bcp.closing:
				cc <- struct{}{}
				return
			case <-bcp.startC:
				bcp.isStart = true
				b := bcp.bc.Head()
				if b != nil {
					bcp.blkC <- *(bcp.bc.Head())
				}
			case b := <-bcp.blkC:
				if bcp.isStart {
					if err := bcp.parseTo(b.Head.BkSeq); err != nil {
						logger.Fatal(err)
					}
				}
			}
		}
	}()
}

// Stop close the block parsing process.
func (bcp *BlockchainParser) Stop() {
	cc := make(chan struct{})
	bcp.closing <- cc
	<-cc
	logger.Debug("blockchain parser stopped")
}

func (bcp *BlockchainParser) parseTo(bcHeight uint64) error {
	if bcp.parsedHeight == 0 {
		// logger.Critical("historydb parse %d/%d", bcp.parsedHeight, bcHeight)
		for i := uint64(0); i <= bcHeight-bcp.parsedHeight; i++ {
			b := bcp.bc.GetBlockInDepth(bcp.parsedHeight + i)
			if b == nil {
				return fmt.Errorf("no block exist in depth:%d", bcp.parsedHeight+i)
			}

			if err := bcp.historyDB.ProcessBlock(b); err != nil {
				return err
			}
		}
		// logger.Critical("historydb parse %d/%d", bcHeight, bcHeight)
	} else {
		// logger.Critical("historydb parse %d/%d", bcp.parsedHeight, bcHeight)
		for i := uint64(0); i < bcHeight-bcp.parsedHeight; i++ {
			b := bcp.bc.GetBlockInDepth(bcp.parsedHeight + i + 1)
			if b == nil {
				return fmt.Errorf("no block exist in depth:%d", bcp.parsedHeight+i+1)
			}

			if err := bcp.historyDB.ProcessBlock(b); err != nil {
				return err
			}
		}
		// logger.Critical("historydb parse %d/%d", bcHeight, bcHeight)
	}
	bcp.parsedHeight = bcHeight
	return nil
}
