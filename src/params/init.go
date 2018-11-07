package params

import (
	"fmt"
	"os"
	"strconv"
)

func init() {
	loadCoinHourBurnFactor()
	loadMaxUserTransactionSize()

	// Compute maxDropletDivisor from precision
	maxDropletDivisor = calculateDivisor(MaxDropletPrecision)

	sanityCheck()
}

func sanityCheck() {
	if UserBurnFactor <= 1 {
		panic("UserBurnFactor must be > 1")
	}

	if MaxUserTransactionSize < 1024 {
		panic("MaxUserTransactionSize must be >= 1024")
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

func loadCoinHourBurnFactor() {
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

	UserBurnFactor = uint32(x)
}

func loadMaxUserTransactionSize() {
	xs := os.Getenv("MAX_USER_TXN_SIZE")
	if xs == "" {
		return
	}

	x, err := strconv.ParseUint(xs, 10, 32)
	if err != nil {
		panic(fmt.Sprintf("Invalid MAX_USER_TXN_SIZE %q: %v", xs, err))
	}

	if x < 1024 {
		panic("MAX_USER_TXN_SIZE must be >= 1024")
	}

	MaxUserTransactionSize = uint32(x)
}
