package visor

import (
	"fmt"
	"time"

	"github.com/skycoin/skycoin/src/visor/historydb"
)

const parseCycle = 1000 // millisecond

// BlockchainParser parses the blockchain and stores the data into historydb.
type BlockchainParser struct {
	bc           *Blockchain
	historyDB    *historydb.HistoryDB
	parsedHeight uint64
	closing      chan chan struct{}
}

// NewBlockchainParser create and init the parser instance.
func NewBlockchainParser(hisDB *historydb.HistoryDB, bc *Blockchain) *BlockchainParser {
	return &BlockchainParser{
		bc:           bc,
		historyDB:    hisDB,
		closing:      make(chan chan struct{}),
		parsedHeight: 0,
	}
}

// Start start to parse the blockchain.
func (bcp *BlockchainParser) Start() {
	go func() {
		logger.Debug("start blockchain parser")
		for {
			select {
			case cc := <-bcp.closing:
				cc <- struct{}{}
			default:
				bcHeight := bcp.bc.Head().Seq()
				if bcp.parsedHeight > bcHeight {
					logger.Fatal("parsedHeight must be <= blockchain height seq")
				}

				if bcp.parsedHeight == bcHeight {
					logger.Debug("no new block need to parse")
					time.Sleep(time.Duration(1000) * time.Millisecond)
					continue
				}

				if err := bcp.parseTo(bcHeight); err != nil {
					logger.Fatal(err)
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
	logger.Debug("blockchian parser stopped")
}

func (bcp *BlockchainParser) parseTo(bcHeight uint64) error {
	if bcp.parsedHeight == 0 {
		for i := uint64(0); i <= bcHeight-bcp.parsedHeight; i++ {
			b := bcp.bc.GetBlockInDepth(bcp.parsedHeight + i)
			if b == nil {
				return fmt.Errorf("no block exist in depth:%d", bcp.parsedHeight+i)
			}

			if err := bcp.historyDB.ProcessBlock(b); err != nil {
				return err
			}
		}
	} else {
		for i := uint64(0); i < bcHeight-bcp.parsedHeight; i++ {
			b := bcp.bc.GetBlockInDepth(bcp.parsedHeight + i + 1)
			if b == nil {
				return fmt.Errorf("no block exist in depth:%d", bcp.parsedHeight+i+1)
			}

			if err := bcp.historyDB.ProcessBlock(b); err != nil {
				return err
			}
		}
	}
	bcp.parsedHeight = bcHeight
	return nil
}
