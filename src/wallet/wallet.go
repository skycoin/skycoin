package wallet

import (
	"fmt"
	"time"

	"encoding/hex"

	logging "github.com/op/go-logging"
	"github.com/skycoin/skycoin/src/cipher"
	//"math/rand"
)

var (
	logger = logging.MustGetLogger("skycoin.visor")
)

const WalletExt = "wlt"
const WalletTimestampFormat = "2006_01_02"

//check for collisions and retry if failure
func NewWalletFilename() string {
	timestamp := time.Now().Format(WalletTimestampFormat)
	//should read in wallet files and make sure does not exist
	padding := hex.EncodeToString((cipher.RandByte(2)))
	return fmt.Sprintf("%s_%s.%s", timestamp, padding, WalletExt)
}
