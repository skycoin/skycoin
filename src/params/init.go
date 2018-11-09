package params

import (
	"fmt"
	"os"
	"strconv"

	"github.com/skycoin/skycoin/src/util/droplet"
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

	if InitialUnlockedCount > DistributionAddressesTotal {
		panic("unlocked addresses > total distribution addresses")
	}

	if uint64(len(distributionAddresses)) != DistributionAddressesTotal {
		panic("available distribution addresses > total allowed distribution addresses")
	}

	if DistributionAddressInitialBalance*DistributionAddressesTotal > MaxCoinSupply {
		panic("total balance in distribution addresses > max coin supply")
	}

	if MaxCoinSupply%DistributionAddressesTotal != 0 {
		panic("MaxCoinSupply should be perfectly divisible by DistributionAddressesTotal")
	}
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

	if x <= 1 {
		panic("USER_BURN_FACTOR must be > 1")
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

	if x < 1024 {
		panic("USER_MAX_TXN_SIZE must be >= 1024")
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

	if x > droplet.Exponent {
		panic("USER_MAX_DECIMALS must be >= droplet.Exponent")
	}

	UserVerifyTxn.MaxDropletPrecision = uint8(x)
}
