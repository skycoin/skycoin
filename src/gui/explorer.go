package gui

import (
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/util/droplet"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
	"github.com/skycoin/skycoin/src/visor"
)

// CoinSupply records the coin supply info
type CoinSupply struct {
	// Coins distributed beyond the project:
	CurrentSupply string `json:"current_supply"`
	// TotalSupply is CurrentSupply plus coins held by the distribution addresses that are spendable
	TotalSupply string `json:"total_supply"`
	// MaxSupply is the maximum number of coins to be distributed ever
	MaxSupply string `json:"max_supply"`
	// CurrentCoinHourSupply is coins hours in non distribution addresses
	CurrentCoinHourSupply string `json:"current_coinhour_supply"`
	// TotalCoinHourSupply is coin hours in all addresses including unlocked distribution addresses
	TotalCoinHourSupply string `json:"total_coinhour_supply"`
	// Distribution addresses which count towards total supply
	UnlockedAddresses []string `json:"unlocked_distribution_addresses"`
	// Distribution addresses which are locked and do not count towards total supply
	LockedAddresses []string `json:"locked_distribution_addresses"`
}

func getCoinSupply(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		supply := coinSupply(gateway, w, r)
		if supply != nil {
			wh.SendJSONOr500(logger, w, supply)
		}
	}
}

func coinSupply(gateway Gatewayer, w http.ResponseWriter, r *http.Request) *CoinSupply {
	if r.Method != http.MethodGet {
		wh.Error405(w)
		return nil
	}

	allUnspents, err := gateway.GetUnspentOutputs()
	if err != nil {
		logger.Error("gateway.GetUnspentOutputs error: %v", err)
		wh.Error500(w)
		return nil
	}

	unlockedAddrs := visor.GetUnlockedDistributionAddresses()
	// Search map of unlocked addresses
	// used to filter unspents
	unlockedAddrMap := daemon.MakeSearchMap(unlockedAddrs)

	var unlockedSupply uint64
	// check confirmed unspents only
	for _, u := range allUnspents.HeadOutputs {
		// check if address is an unlocked distribution address
		if _, ok := unlockedAddrMap[u.Address]; ok {
			coins, err := droplet.FromString(u.Coins)
			if err != nil {
				logger.Error("Invalid unlocked output balance string %s: %v", u.Coins, err)
				wh.Error500(w)
				return nil
			}
			unlockedSupply += coins
		}
	}

	// "total supply" is the number of coins unlocked.
	// Each distribution address was allocated visor.DistributionAddressInitialBalance coins.
	totalSupply := uint64(len(unlockedAddrs)) * visor.DistributionAddressInitialBalance
	totalSupply *= droplet.Multiplier

	// "current supply" is the number of coins distributed from the unlocked pool
	currentSupply := totalSupply - unlockedSupply

	currentSupplyStr, err := droplet.ToString(currentSupply)
	if err != nil {
		logger.Error("Failed to convert coins to string: %v", err)
		wh.Error500(w)
		return nil
	}

	totalSupplyStr, err := droplet.ToString(totalSupply)
	if err != nil {
		logger.Error("Failed to convert coins to string: %v", err)
		wh.Error500(w)
		return nil
	}

	maxSupplyStr, err := droplet.ToString(visor.MaxCoinSupply * droplet.Multiplier)
	if err != nil {
		logger.Error("Failed to convert coins to string: %v", err)
		wh.Error500(w)
		return nil
	}

	// locked distribution addresses
	lockedAddrs := visor.GetLockedDistributionAddresses()
	lockedAddrMap := daemon.MakeSearchMap(lockedAddrs)

	// get total coins hours which excludes locked distribution addresses
	var totalCoinHours uint64
	for _, out := range allUnspents.HeadOutputs {
		if _, ok := lockedAddrMap[out.Address]; !ok {
			totalCoinHours += out.CalculatedHours
		}
	}

	// get current coin hours which excludes all distribution addresses
	var currentCoinHours uint64
	for _, out := range allUnspents.HeadOutputs {
		// check if address not in locked distribution addresses
		if _, ok := lockedAddrMap[out.Address]; !ok {
			// check if address not in unlocked distribution addresses
			if _, ok := unlockedAddrMap[out.Address]; !ok {
				currentCoinHours += out.CalculatedHours
			}
		}
	}

	if err != nil {
		logger.Errorf("Failed to get total coinhours: %v", err.Error())
		wh.Error500(w)
		return nil
	}

	cs := CoinSupply{
		CurrentSupply:         currentSupplyStr,
		TotalSupply:           totalSupplyStr,
		MaxSupply:             maxSupplyStr,
		CurrentCoinHourSupply: strconv.FormatUint(currentCoinHours, 10),
		TotalCoinHourSupply:   strconv.FormatUint(totalCoinHours, 10),
		UnlockedAddresses:     unlockedAddrs,
		LockedAddresses:       visor.GetLockedDistributionAddresses(),
	}

	return &cs
}

// method: GET
// url: /explorer/address?address=${address}
func getTransactionsForAddress(gateway Gatewayer) http.HandlerFunc {
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
					logger.Error("uxout of %v does not exist in history db", id.Hex())
					wh.Error500(w)
					return
				}

				tIn, err := visor.NewReadableTransactionInput(tx.Transaction.In[i], uxout.Out.Body.Address.String(), uxout.Out.Body.Coins, uxout.Out.Body.Hours)
				if err != nil {
					wh.Error500(w)
					return
				}

				in[i] = *tIn
			}

			resTxs = append(resTxs, NewReadableTransaction(tx, in))
		}

		wh.SendJSONOr500(logger, w, &resTxs)
	}
}

// Richlist is the API response for /richlist, contains top address balances
type Richlist struct {
	Richlist visor.Richlist `json:"richlist"`
}

// method: GET
// url: /richlist?n=${number}&include-distribution=${bool}
func getRichlist(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		var topn int
		topnStr := r.FormValue("n")
		if topnStr == "" {
			topn = 20
		} else {
			var err error
			topn, err = strconv.Atoi(topnStr)
			if err != nil {
				wh.Error400(w, "invalid n")
				return
			}
		}

		var includeDistribution bool
		includeDistributionStr := r.FormValue("include-distribution")
		if includeDistributionStr == "" {
			includeDistribution = false
		} else {
			var err error
			includeDistribution, err = strconv.ParseBool(includeDistributionStr)
			if err != nil {
				wh.Error400(w, "invalid include-distribution")
				return
			}
		}

		richlist, err := gateway.GetRichlist(includeDistribution)
		if err != nil {
			logger.Error(err.Error())
			wh.Error500(w)
			return
		}

		if topn > 0 && topn < len(richlist) {
			richlist = richlist[:topn]
		}

		wh.SendJSONOr500(logger, w, Richlist{
			Richlist: richlist,
		})
	}
}

// method: GET
// url: /addresscount
func getAddressCount(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		addrCount, err := gateway.GetAddressCount()
		if err != nil {
			logger.Error(err.Error())
			wh.Error500(w)
			return
		}

		wh.SendJSONOr500(logger, w, &map[string]uint64{"count": addrCount})
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
