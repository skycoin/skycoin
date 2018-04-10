package gui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/droplet"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor"
)

type AdvancedSpendResult struct {
	Transaction *visor.ReadableTransaction `json:"txn,omitempty"`
	Error       string                     `json:"error,omitempty"`
}
func advancedSpendHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		var advancedSpend wallet.AdvancedSpendRequest
		err := json.NewDecoder(r.Body).Decode(&advancedSpend)
		if err != nil {
			logger.Errorf("Invalid advanced spend request: %v", err)
			wh.Error400(w, fmt.Sprintf("Bad Request: %v", err))
			return
		}
		defer func() {
			if err := r.Body.Close(); err != nil {
				logger.Errorf("Failed to close response body: %v", err)
			}
		}()

		switch advancedSpend.HoursSelection.Type {
		case "manual", "auto":
		case "":
			logger.Error("missing hours selection type")
			wh.Error400(w, "missing hours selection type")
			return
		default:
			logger.Errorf("invalid hours selection type: %v", advancedSpend.HoursSelection.Type)
			wh.Error400(w, fmt.Sprintf("invalid hours selection type: %v", advancedSpend.HoursSelection.Type))
			return
		}

		if advancedSpend.HoursSelection.Type == "auto" {
			switch advancedSpend.HoursSelection.Mode {
			case "split_even":
				switch advancedSpend.HoursSelection.ShareFactor {
				case "":
					logger.Error("missing hours selection share factor when mode is `split_even`")
					wh.Error400(w, "missing hours selection share factor when mode is split_even")
					return
				default:
					shareFactor, err := strconv.ParseFloat(advancedSpend.HoursSelection.ShareFactor, 64)
					if err != nil {
						logger.Error(err)
						wh.Error400(w, fmt.Sprintf("invalid share factor: %v", shareFactor))
						return
					}

					if shareFactor < 0 {
						logger.Warning("negative share factor")
						wh.Error400(w, "share factor cannot be negative")
						return
					}
				}

			case "match_coins":
			case "":
				logger.Error("missing hours selection mode when type is auto")
				wh.Error400(w, "missing hours selection mode when type is auto")
				return
			default:
				logger.Errorf("invalid hours selection mode: %v", advancedSpend.HoursSelection.Mode)
				wh.Error400(w, fmt.Sprintf("invalid hours selection mode: %v", advancedSpend.HoursSelection.Mode))
				return
			}
		}

		if advancedSpend.ChangeAddress == "" {
			logger.Warning("missing change address")
			wh.Error400(w, "missing change address")
			return
		}

		changeAddr, err := cipher.DecodeBase58Address(advancedSpend.ChangeAddress)
		if err != nil {
			logger.Errorf("invalid change address: %v", err)
			wh.Error400(w, fmt.Sprintf("invalid change address: %v", err))
			return
		}

		// check whether destination addresses are correct or not
		destList := make([]coin.TransactionOutput, len(advancedSpend.To))
		for idx, to := range advancedSpend.To {
			// check that the address is valid
			toAddress, err := cipher.DecodeBase58Address(to.Address)
			if err != nil {
				logger.Errorf("invalid destination address %v", to.Address)
				wh.Error400(w, fmt.Sprintf("invalid destination address %v", to.Address))
				return
			}

			// convert coins to droplets
			coins, err := droplet.FromString(to.Coins)
			if err != nil {
				logger.Errorf("unable to convert coins to droplet: %v", err)
				wh.Error400(w, fmt.Sprintf("invalid coin amount %v", to.Coins))
				return
			}

			destList[idx] = coin.TransactionOutput{
				Address: toAddress,
				Coins:   coins,
			}

			// parse coinhours
			// when mode is auto the Hours field can be empty
			// when mode is manual Hours field cannot be empty
			if to.Hours != "" {
				hours, err := strconv.ParseUint(to.Hours, 10, 64)
				if err != nil {
					logger.Errorf("unable to parse coinhours: %v", err)
					wh.Error400(w, fmt.Sprintf("invalid coinhours %v", to.Hours))
					return
				}

				destList[idx].Hours = hours
			} else if advancedSpend.HoursSelection.Type == "manual" {
				logger.Errorf("coinhours value missing for %v when mode is manual", to.Address)
				wh.Error400(w, fmt.Sprintf("coinhours value missing for %v when mode is manual", to.Address))
				return
			}

		}

		// create a wltmap
		wltMap := make(map[string]struct{})
		for _, wlt := range advancedSpend.Wallets {
			wltMap[wlt] = struct{}{}
		}

		// fetch all wallets on the system
		wallets, err := gateway.GetWallets()
		if err != nil {
			logger.Error(err)
			wh.Error500(w)
			return
		}

		// entry map
		// stores all entries to be used in creating the transaction
		entryMap := make(map[cipher.Address]wallet.Entry)

		// addr map to check if provided input addresses
		// are in the wallets on the system
		addrMap := make(map[string]wallet.Entry)

		// iterate over all the wallets
		// look for wallets provided in the input
		// collect all entries of those wallets
		for _, wlt := range wallets {
			if _, ok := wltMap[wlt.Label()]; ok {
				for _, wltEntry := range wlt.Entries {
					entryMap[wltEntry.Address] = wltEntry
				}
			} else {
				for _, wltEntry := range wlt.Entries {
					addrMap[wltEntry.Address.String()] = wltEntry
				}
			}
		}

		// check that provided addresses are in the addrMap
		for _, addr := range advancedSpend.Addresses {
			// check that address is not an empty string
			if addr != "" {
				if _, ok := addrMap[addr]; !ok {
					logger.Errorf("address %v not found in any wallet", addr)
					wh.Error400(w, fmt.Sprintf("address %v not found in any wallet", addr))
					return
				}

				wltEntry := addrMap[addr]
				entryMap[wltEntry.Address] = wltEntry
			} else {
				logger.Warningf("empty sender address")
				wh.Error400(w, "empty sender address")
				return
			}
		}

		if len(entryMap) == 0 {
			logger.Error("no sender addresses found")
			wh.Error400(w, "no sender addresses found")
			return
		}

		tx, err :=gateway.AdvancedSpend(
			wallet.AdvancedSpend{
				HoursSelection: advancedSpend.HoursSelection,
				Entries:        entryMap,
				ChangeAddress:  changeAddr,
				To:             destList,
			},
		)

		switch err {
		case nil:
		case fee.ErrTxnNoFee, wallet.ErrSpendingUnconfirmed, wallet.ErrInsufficientBalance, wallet.ErrZeroSpend:
			wh.Error400(w, err.Error())
			return
		case wallet.ErrWalletNotExist:
			wh.Error404(w)
			return
		case wallet.ErrWalletAPIDisabled:
			wh.Error403(w)
			return
		default:
			wh.Error500Msg(w, err.Error())
			return
		}

		txStr, err := visor.TransactionToJSON(*tx)
		if err != nil {
			logger.Error(err)
			wh.SendJSONOr500(logger, w, SpendResult{
				Error: err.Error(),
			})
			return
		}

		logger.Infof("Spend: \ntx= \n %s \n", txStr)

		var ret AdvancedSpendResult
		ret.Transaction, err = visor.NewReadableTransaction(&visor.Transaction{Txn: *tx})
		if err != nil {
			err = fmt.Errorf("Creation of new readable transaction failed: %v", err)
			logger.Error(err)
			ret.Error = err.Error()
			wh.SendJSONOr500(logger, w, ret)
			return
		}

		wh.SendJSONOr500(logger, w, ret)
	}
}
