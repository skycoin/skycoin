package wallet

import (
	"fmt"
	"time"

	"github.com/op/go-logging"
	"github.com/skycoin/skycoin/src/cipher"
	"math/rand"
)

var (
	logger = logging.MustGetLogger("skycoin.visor")
)

const WalletExt = "wlt"
const WalletTimestampFormat = "2006_01_02"

type WalletID string
type AddressSet map[cipher.Address]byte

func (self AddressSet) Update(other AddressSet) AddressSet {
	for k, v := range other {
		self[k] = v
	}
	return self
}

//type WalletConstructor func() Wallet

//check for collisions and retry if failure
func NewWalletFilename() string {
	timestamp := time.Now().Format(WalletTimestampFormat)
	id := rand.Int() % 9999 //should read in wallet files and make sure does not exist
	return fmt.Sprintf("%s_%04d.%s", timestamp, id, WalletExt)
}
