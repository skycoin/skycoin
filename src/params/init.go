package params

import (
	"fmt"
	"os"
	"strconv"
)

func init() {
	loadCoinHourBurnFactor()

	// Compute maxDropletDivisor from precision
	maxDropletDivisor = calculateDivisor(MaxDropletPrecision)

	sanityCheck()
}

func sanityCheck() {
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
	xs := os.Getenv("COINHOUR_BURN_FACTOR")
	if xs == "" {
		return
	}

	x, err := strconv.ParseUint(xs, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Invalid COINHOUR_BURN_FACTOR %q: %v", xs, err))
	}

	if x <= 1 {
		panic(fmt.Sprintf("CoinHourBurnFactor must be > 1"))
	}

	CoinHourBurnFactor = x
}
