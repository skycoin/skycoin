package api

import (
	"fmt"
	"net/http"

	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/visor"
)

// outputsHandler returns UxOuts filtered by a set of addresses or a set of hashes
// URI: /api/v1/outputs
// Method: GET, POST
// Args:
//    addrs: comma-separated list of addresses
//    hashes: comma-separated list of uxout hashes
// If neither addrs nor hashes are specificed, return all unspent outputs.
// If only one filter is specified, then return outputs match the filter.
// Both filters cannot be specified.
func outputsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		addrStr := r.FormValue("addrs")
		hashStr := r.FormValue("hashes")

		if addrStr != "" && hashStr != "" {
			wh.Error400(w, "addrs and hashes cannot be specified together")
			return
		}

		var filters []visor.OutputsFilter

		if addrStr != "" {
			addrs, err := parseAddressesFromStr(addrStr)
			if err != nil {
				wh.Error400(w, err.Error())
				return
			}

			if len(addrs) > 0 {
				filters = append(filters, visor.FbyAddresses(addrs))
			}
		}

		if hashStr != "" {
			hashes, err := parseHashesFromStr(hashStr)
			if err != nil {
				wh.Error400(w, err.Error())
				return
			}

			if len(hashes) > 0 {
				filters = append(filters, visor.FbyHashes(hashes))
			}
		}

		summary, err := gateway.GetUnspentOutputsSummary(filters)
		if err != nil {
			err = fmt.Errorf("gateway.GetUnspentOutputsSummary failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		rSummary, err := readable.NewUnspentOutputsSummary(summary)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, rSummary)
	}
}
