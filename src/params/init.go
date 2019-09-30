package params

import (
	"fmt"
	"os"
	"strconv"

	"github.com/SkycoinProject/skycoin/src/util/droplet"
)

func init() {
	loadUserBurnFactor()
	loadUserMaxTransactionSize()
	loadUserMaxDecimals()
	sanityCheck()
}

func sanityCheck() {
	if err := UserVerifyTxn.Validate(); err != nil {
		panic(err)
	}

	MainNetDistribution.MustValidate()
}

func loadUserBurnFactor() {
	xs := os.Getenv("USER_BURN_FACTOR")
	if xs == "" {
		return
	}

	x, err := strconv.ParseUint(xs, 10, 32)
	if err != nil {
		panic(fmt.Sprintf("Invalid USER_BURN_FACTOR %q: %v", xs, err))
	}

	if x < uint64(MinBurnFactor) {
		panic(fmt.Sprintf("USER_BURN_FACTOR must be >= %d", MinBurnFactor))
	}

	UserVerifyTxn.BurnFactor = uint32(x)
}

func loadUserMaxTransactionSize() {
	xs := os.Getenv("USER_MAX_TXN_SIZE")
	if xs == "" {
		return
	}

	x, err := strconv.ParseUint(xs, 10, 32)
	if err != nil {
		panic(fmt.Sprintf("Invalid USER_MAX_TXN_SIZE %q: %v", xs, err))
	}

	if x < uint64(MinTransactionSize) {
		panic(fmt.Sprintf("USER_MAX_TXN_SIZE must be >= %d", MinTransactionSize))
	}

	UserVerifyTxn.MaxTransactionSize = uint32(x)
}

func loadUserMaxDecimals() {
	xs := os.Getenv("USER_MAX_DECIMALS")
	if xs == "" {
		return
	}

	x, err := strconv.ParseUint(xs, 10, 8)
	if err != nil {
		panic(fmt.Sprintf("Invalid USER_MAX_DECIMALS %q: %v", xs, err))
	}

	if x > uint64(droplet.Exponent) {
		panic(fmt.Sprintf("USER_MAX_DECIMALS must be <= %d", droplet.Exponent))
	}

	UserVerifyTxn.MaxDropletPrecision = uint8(x)
}
