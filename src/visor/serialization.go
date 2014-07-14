package visor

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/skycoin/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
)

type SerializedBlockchain struct {
	Blocks   []coin.Block
	Unspents coin.UxArray
}

func NewSerializedBlockchain(bc *coin.Blockchain) *SerializedBlockchain {
	return &SerializedBlockchain{
		Blocks:   bc.Blocks,
		Unspents: bc.Unspent.Array(),
	}
}

func (self *SerializedBlockchain) Save(filename string) error {
	data := encoder.Serialize(self)
	return util.SaveBinary(filename, data, 0644)
}

func (self SerializedBlockchain) ToBlockchain() *coin.Blockchain {
	bc := &coin.Blockchain{}
	bc.Blocks = self.Blocks
	pool := coin.NewUnspentPool()
	pool.Rebuild(self.Unspents)
	bc.Unspent = pool
	return bc
}

func LoadSerializedBlockchain(filename string) (*SerializedBlockchain, error) {
	sbc := &SerializedBlockchain{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = encoder.DeserializeRaw(data, sbc)
	if err != nil {
		return nil, err
	}
	return sbc, nil
}

// Saves blockchain to disk
func SaveBlockchain(bc *coin.Blockchain, filename string) error {
	sbc := NewSerializedBlockchain(bc)
	return sbc.Save(filename)
}

// Loads a coin.Blockchain from disk
func LoadBlockchain(filename string) (*coin.Blockchain, error) {
	if sbc, err := LoadSerializedBlockchain(filename); err == nil {
		logger.Info("Loaded serialized blockchain from \"%s\"", filename)
		return sbc.ToBlockchain(), nil
	} else {
		return nil, err
	}
}

// Loads a blockchain but subdues errors into the logger, or panics.
// If no blockchain is found, it creates a new empty one
func loadBlockchain(filename string, genAddr cipher.Address) *coin.Blockchain {
	bc := &coin.Blockchain{}
	created := false
	if filename != "" {
		var err error
		bc, err = LoadBlockchain(filename)
		if err == nil {
			if len(bc.Blocks) == 0 {
				log.Panic("Loaded empty blockchain")
			}
			loadedGenAddr := bc.Blocks[0].Body.Transactions[0].Out[0].Address
			if loadedGenAddr != genAddr {
				log.Panic("Configured genesis address does not match the " +
					"address in the blockchain")
			}
			created = true
		} else {
			if os.IsNotExist(err) {
				logger.Info("No blockchain file, will create a new blockchain")
			} else {
				log.Panicf("Failed to load blockchain file \"%s\": %v",
					filename, err)
			}
		}
	}
	if !created {
		bc = coin.NewBlockchain()
	}
	return bc
}
