package api

// APIs for blockchain related information (v2)

import (
	"fmt"
	"net/http"
	"strconv"

	wh "github.com/skycoin/skycoin/src/util/http"
)

func getBlocksV2(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}
		sstart := r.FormValue("start")
		start, err := strconv.ParseUint(sstart, 10, 64)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("Invalid start value \"%s\"", sstart))
			return
		}

		send := r.FormValue("end")
		end, err := strconv.ParseUint(send, 10, 64)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("Invalid end value \"%s\"", send))
			return
		}
		rb, err := gateway.GetBlocksV2(start, end)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("Get blocks failed: %v", err))
			return
		}
		wh.SendJSONOr500(logger, w, rb)
	}
}
