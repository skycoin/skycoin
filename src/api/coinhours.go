package api

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/fee"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/wallet"
)

// depositCoinhoursRequest is sent to /wallet/coinhours/deposit
type depositCoinhoursRequest struct {
	IgnoreUnconfirmed bool                          `json:"ignore_unconfirmed"`
	Wallet            depositCoinhoursRequestWallet `json:"wallet"`
	Amount            uint64                        `json:"amount"`
	To                wh.Address                    `json:"to"`
}

// depositCoinhoursRequestWallet defines a wallet to send coinhours from and from which address in the wallet as well
type depositCoinhoursRequestWallet struct {
	ID       string     `json:"id"`
	Password string     `json:"password"`
	Address  wh.Address `json:"address"`
}

// Deposit coinhours into the coinhour bank
// URI: /coinhourbank/deposit
// Method: POST
// Args: JSON body
func chbDepositHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		var params depositCoinhoursRequest
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			logger.WithError(err).Error("Invalid deposit coinhours request")
			wh.Error400(w, err.Error())
			return
		}

		outputSummary, err := gateway.GetUnspentOutputsSummary([]visor.OutputsFilter{visor.FbyAddresses([]cipher.Address{params.To.Address})})
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		// convert unspentoutputs to uxarray
		var uxa []coin.UxOut
		for _, ux := range outputSummary.Confirmed {
			uxa = append(uxa, ux.UxOut)
		}

		wlt, err := gateway.GetWallet(params.Wallet.ID)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		err = gateway.DepositCoinhours(params.Amount, params.To.String(), uxa, wlt)
		if err != nil {
			switch err.(type) {
			case wallet.Error:
				switch err {
				case wallet.ErrWalletAPIDisabled:
					wh.Error403(w, "")
				case wallet.ErrWalletNotExist:
					wh.Error404(w, err.Error())
				default:
					wh.Error400(w, err.Error())
				}
			case blockdb.ErrUnspentNotExist:
				wh.Error400(w, err.Error())
			default:
				switch err {
				case fee.ErrTxnNoFee,
					fee.ErrTxnInsufficientCoinHours,
					wallet.ErrSpendingUnconfirmed:
					wh.Error400(w, err.Error())
				default:
					wh.Error500(w, err.Error())
				}
			}
		}
	}
}