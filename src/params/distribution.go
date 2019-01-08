package params

import "github.com/skycoin/skycoin/src/cipher"

// distributionAddressesDecoded is initialized in init.go from params.go's distributionAddresses
var distributionAddressesDecoded []cipher.Address

// GetDistributionAddresses returns a copy of the hardcoded distribution addresses array.
// Each address has 1,000,000 coins. There are 100 addresses.
func GetDistributionAddresses() []string {
	addrs := make([]string, len(distributionAddresses))
	for i := range distributionAddresses {
		addrs[i] = distributionAddresses[i]
	}
	return addrs
}

// GetUnlockedDistributionAddresses returns distribution addresses that are unlocked, i.e. they have spendable outputs
func GetUnlockedDistributionAddresses() []string {
	// The first InitialUnlockedCount (25) addresses are unlocked by default.
	// Subsequent addresses will be unlocked at a rate of UnlockAddressRate (5) per year,
	// after the InitialUnlockedCount (25) addresses have no remaining balance.
	// The unlock timer will be enabled manually once the
	// InitialUnlockedCount (25) addresses are distributed.

	// NOTE: To have automatic unlocking, transaction verification would have
	// to be handled in visor rather than in coin.Transactions.Visor(), because
	// the coin package is agnostic to the state of the blockchain and cannot reference it.
	// Instead of automatic unlocking, we can hardcode the timestamp at which the first 30%
	// is distributed, then compute the unlocked addresses easily here.

	addrs := make([]string, InitialUnlockedCount)
	copy(addrs[:], distributionAddresses[:InitialUnlockedCount])
	return addrs
}

// GetLockedDistributionAddresses returns distribution addresses that are locked, i.e. they have unspendable outputs
func GetLockedDistributionAddresses() []string {
	// TODO -- once we reach 30% distribution, we can hardcode the
	// initial timestamp for releasing more coins
	addrs := make([]string, DistributionAddressesTotal-InitialUnlockedCount)
	for i := range distributionAddresses[InitialUnlockedCount:] {
		addrs[i] = distributionAddresses[InitialUnlockedCount+uint64(i)]
	}
	return addrs
}

// GetDistributionAddressesDecoded returns a copy of the hardcoded distribution addresses array.
// Each address has 1,000,000 coins. There are 100 addresses.
func GetDistributionAddressesDecoded() []cipher.Address {
	addrs := make([]cipher.Address, len(distributionAddressesDecoded))
	for i := range distributionAddressesDecoded {
		addrs[i] = distributionAddressesDecoded[i]
	}
	return addrs
}

// GetUnlockedDistributionAddressesDecoded returns distribution addresses that are unlocked, i.e. they have spendable outputs
func GetUnlockedDistributionAddressesDecoded() []cipher.Address {
	// The first InitialUnlockedCount (25) addresses are unlocked by default.
	// Subsequent addresses will be unlocked at a rate of UnlockAddressRate (5) per year,
	// after the InitialUnlockedCount (25) addresses have no remaining balance.
	// The unlock timer will be enabled manually once the
	// InitialUnlockedCount (25) addresses are distributed.

	// NOTE: To have automatic unlocking, transaction verification would have
	// to be handled in visor rather than in coin.Transactions.Visor(), because
	// the coin package is agnostic to the state of the blockchain and cannot reference it.
	// Instead of automatic unlocking, we can hardcode the timestamp at which the first 30%
	// is distributed, then compute the unlocked addresses easily here.

	addrs := make([]cipher.Address, InitialUnlockedCount)
	copy(addrs[:], distributionAddressesDecoded[:InitialUnlockedCount])
	return addrs
}

// GetLockedDistributionAddressesDecoded returns distribution addresses that are locked, i.e. they have unspendable outputs
func GetLockedDistributionAddressesDecoded() []cipher.Address {
	// TODO -- once we reach 30% distribution, we can hardcode the
	// initial timestamp for releasing more coins
	addrs := make([]cipher.Address, DistributionAddressesTotal-InitialUnlockedCount)
	for i := range distributionAddressesDecoded[InitialUnlockedCount:] {
		addrs[i] = distributionAddressesDecoded[InitialUnlockedCount+uint64(i)]
	}
	return addrs
}
