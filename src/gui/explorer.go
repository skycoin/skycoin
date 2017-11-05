package gui

import (
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/util/droplet"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
	"github.com/skycoin/skycoin/src/visor"
)

// RegisterExplorerHandlers register explorer handlers
func RegisterExplorerHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	// get set of pending transactions
	mux.HandleFunc("/explorer/address", getTransactionsForAddress(gateway))

	mux.HandleFunc("/explorer/getEffectiveOutputs", getEffectiveOutputs(gateway))

	mux.HandleFunc("/coinSupply", getCoinSupply(gateway))
}

// DeprecatedCoinSupply records the coin supply info
type DeprecatedCoinSupply struct {
	CoinSupply
	DeprecatedCurrentSupply                           uint64   `json:"coinSupply"`
	DeprecatedCoinCap                                 uint64   `json:"coinCap"`
	DeprecatedUndistributedLockedCoinBalance          uint64   `json:"UndistributedLockedCoinBalance"`
	DeprecatedUndistributedLockedCoinHoldingAddresses []string `json:"UndistributedLockedCoinHoldingAddresses"`
}

// CoinSupply records the coin supply info
type CoinSupply struct {
	// Coins distributed beyond the project:
	CurrentSupply string `json:"current_supply"`
	// TotalSupply is CurrentSupply plus coins held by the distribution addresses that are spendable
	TotalSupply string `json:"total_supply"`
	// MaxSupply is the maximum number of coins to be distributed ever
	MaxSupply string `json:"max_supply"`
	// Distribution addresses which count towards total supply
	UnlockedAddresses []string `json:"unlocked_distribution_addresses"`
	// Distribution addresses which are locked and do not count towards total supply
	LockedAddresses []string `json:"locked_distribution_addresses"`
}

func getCoinSupply(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		supply, _ := coinSupply(gateway, w, r)
		if supply != nil {
			wh.SendOr404(w, supply)
		}
	}
}

// TODO: DEPRECATED Remove for v21 release
func getEffectiveOutputs(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, oldSupply := coinSupply(gateway, w, r)
		if oldSupply != nil {
			wh.SendOr404(w, oldSupply)
		}
	}
}

func coinSupply(gateway *daemon.Gateway, w http.ResponseWriter, r *http.Request) (*CoinSupply, *DeprecatedCoinSupply) {
	if r.Method != http.MethodGet {
		wh.Error405(w)
		return nil, nil
	}

	unlockedAddrs := visor.GetUnlockedDistributionAddresses()

	filterInUnlocked := []daemon.OutputsFilter{}
	filterInUnlocked = append(filterInUnlocked, daemon.FbyAddresses(unlockedAddrs))
	unlockedOutputs, err := gateway.GetUnspentOutputs(filterInUnlocked...)
	if err != nil {
		wh.Error500(w)
		return nil, nil
	}

	var unlockedSupply uint64
	for _, u := range unlockedOutputs.HeadOutputs {
		coins, err := droplet.FromString(u.Coins)
		if err != nil {
			logger.Error("Invalid unlocked output balance string %s: %v", u.Coins, err)
			wh.Error500(w)
			return nil, nil
		}
		unlockedSupply += coins
	}

	// "total supply" is the number of coins unlocked.
	// Each distribution address was allocated visor.DistributionAddressInitialBalance coins.
	totalSupply := uint64(len(unlockedAddrs)) * visor.DistributionAddressInitialBalance
	totalSupply *= droplet.Multiplier

	// "current supply" is the number of coins distribution from the unlocked pool
	currentSupply := totalSupply - unlockedSupply

	currentSupplyStr, err := droplet.ToString(currentSupply)
	if err != nil {
		logger.Error("Failed to convert coins to string: %v", err)
		wh.Error500(w)
		return nil, nil
	}

	totalSupplyStr, err := droplet.ToString(totalSupply)
	if err != nil {
		logger.Error("Failed to convert coins to string: %v", err)
		wh.Error500(w)
		return nil, nil
	}

	maxSupplyStr, err := droplet.ToString(visor.MaxCoinSupply * droplet.Multiplier)
	if err != nil {
		logger.Error("Failed to convert coins to string: %v", err)
		wh.Error500(w)
		return nil, nil
	}

	cs := CoinSupply{
		CurrentSupply:     currentSupplyStr,
		TotalSupply:       totalSupplyStr,
		MaxSupply:         maxSupplyStr,
		UnlockedAddresses: unlockedAddrs,
		LockedAddresses:   visor.GetLockedDistributionAddresses(),
	}

	return &cs, &DeprecatedCoinSupply{
		CoinSupply:                                        cs,
		DeprecatedCurrentSupply:                           currentSupply,
		DeprecatedCoinCap:                                 visor.MaxCoinSupply,
		DeprecatedUndistributedLockedCoinBalance:          unlockedSupply,
		DeprecatedUndistributedLockedCoinHoldingAddresses: visor.GetDistributionAddresses(),
	}
}

// method: GET
// url: /explorer/address?address=${address}
func getTransactionsForAddress(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		addr := r.FormValue("address")
		if addr == "" {
			wh.Error400(w, "address is empty")
			return
		}

		cipherAddr, err := cipher.DecodeBase58Address(addr)
		if err != nil {
			wh.Error400(w, "invalid address")
			return
		}

		txns, err := gateway.GetAddressTxns(cipherAddr)
		if err != nil {
			logger.Error("Get address transactions failed: %v", err)
			wh.Error500(w)
			return
		}

		resTxs := make([]ReadableTransaction, 0, len(txns.Txns))

		for _, tx := range txns.Txns {
			in := make([]visor.ReadableTransactionInput, len(tx.Transaction.In))
			for i := range tx.Transaction.In {
				id, err := cipher.SHA256FromHex(tx.Transaction.In[i])
				if err != nil {
					logger.Error("%v", err)
					wh.Error500(w)
					return
				}

				uxout, err := gateway.GetUxOutByID(id)
				if err != nil {
					logger.Error("%v", err)
					wh.Error500(w)
					return
				}

				if uxout == nil {
					logger.Error("uxout of %d does not exist in history db", id)
					wh.Error500(w)
					return
				}

				in[i] = visor.NewReadableTransactionInput(tx.Transaction.In[i], uxout.Out.Body.Address.String())
			}

			resTxs = append(resTxs, NewReadableTransaction(tx, in))
		}

		wh.SendOr404(w, &resTxs)
	}
}

// ReadableTransaction represents readable address transaction
type ReadableTransaction struct {
	Status    visor.TransactionStatus `json:"status"`
	Length    uint32                  `json:"length"`
	Type      uint8                   `json:"type"`
	Hash      string                  `json:"txid"`
	InnerHash string                  `json:"inner_hash"`
	Timestamp uint64                  `json:"timestamp,omitempty"`

	Sigs []string                          `json:"sigs"`
	In   []visor.ReadableTransactionInput  `json:"inputs"`
	Out  []visor.ReadableTransactionOutput `json:"outputs"`
}

// NewReadableTransaction creates readable address transaction
func NewReadableTransaction(t visor.TransactionResult, inputs []visor.ReadableTransactionInput) ReadableTransaction {
	return ReadableTransaction{
		Status:    t.Status,
		Length:    t.Transaction.Length,
		Type:      t.Transaction.Type,
		Hash:      t.Transaction.Hash,
		InnerHash: t.Transaction.InnerHash,
		Timestamp: t.Time,

		Sigs: t.Transaction.Sigs,
		In:   inputs,
		Out:  t.Transaction.Out,
	}
}
