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
// Use goswagger y otras cosas test
func outputsHandler(gateway Gatewayer) http.HandlerFunc {

	// swagger:operation GET /api/v1/outputs outputsGet
	//
	// If neither addrs nor hashes are specificed, return all unspent outputs. If only one filter is specified, then return outputs match the filter. Both filters cannot be specified.
	//
	// ---
	// parameters:
	// - name: address
	//   in: query
	//   required: false
	//   type: string
	// - name: hash
	//   in: query
	//   required: false
	//   type: string
	//
	// produces:
	// - application/json
	//
	// responses:
	//   200:
	//     description: UnspentOutputsSummary records unspent outputs in different status.
	//     schema:
	//       properties:
	//         head:
	//           type: object
	//           properties:
	//             seq:
	//               type: integer
	//               format: int64
	//             hash:
	//               type: string
	//             previous_block_hash:
	//               type: string
	//             tx_body_hash:
	//               type: string
	//             ux_hash:
	//               type: string
	//             timestamp:
	//               type: integer
	//               format: int64
	//             fee:
	//               type: integer
	//               format: int64
	//             version:
	//               type: integer
	//               format: int64
	//         head_outputs:
	//           description: HeadOutputs are unspent outputs confirmed in the blockchain
	//           type: array
	//           items:
	//             properties:
	//               hash:
	//                 type: string
	//               src_tx:
	//                 type: string
	//               address:
	//                 type: string
	//               coins:
	//                 type: string
	//               time:
	//                 type: integer
	//                 format: int64
	//               hours:
	//                 type: integer
	//                 format: int64
	//               calculated_hours:
	//                 type: integer
	//                 format: int64
	//               block_seq:
	//                 type: integer
	//                 format: int64
	//         outgoing_outputs:
	//           description: OutgoingOutputs are unspent outputs being spent in unconfirmed transactions
	//           type: array
	//           items:
	//             properties:
	//               hash:
	//                 type: string
	//               src_tx:
	//                 type: string
	//               address:
	//                 type: string
	//               coins:
	//                 type: string
	//               time:
	//                 type: integer
	//                 format: int64
	//               hours:
	//                 type: integer
	//                 format: int64
	//               calculated_hours:
	//                 type: integer
	//                 format: int64
	//               block_seq:
	//                 type: integer
	//                 format: int64
	//         incoming_outputs:
	//           description: IncomingOutputs are unspent outputs being created by unconfirmed transactions
	//           type: array
	//           items:
	//             properties:
	//               hash:
	//                 type: string
	//               src_tx:
	//                 type: string
	//               address:
	//                 type: string
	//               coins:
	//                 type: string
	//               time:
	//                 type: integer
	//                 format: int64
	//               hours:
	//                 type: integer
	//                 format: int64
	//               calculated_hours:
	//                 type: integer
	//                 format: int64
	//               block_seq:
	//                 type: integer
	//                 format: int64
	//   default:
	//	   $ref: '#/responses/genericError'

	// swagger:operation POST /api/v1/outputs outputsPost
	//
	// If neither addrs nor hashes are specificed, return all unspent outputs. If only one filter is specified, then return outputs match the filter. Both filters cannot be specified.
	//
	// ---
	// parameters:
	// - name: address
	//   in: query
	//   required: false
	//   type: string
	// - name: hash
	//   in: query
	//   required: false
	//   type: string
	//
	// produces:
	// - application/json
	//
	// responses:
	//   200:
	//     description: UnspentOutputsSummary records unspent outputs in different status.
	//     schema:
	//       properties:
	//         head:
	//           type: object
	//           properties:
	//             seq:
	//               type: integer
	//               format: int64
	//             hash:
	//               type: string
	//             previous_block_hash:
	//               type: string
	//             tx_body_hash:
	//               type: string
	//             ux_hash:
	//               type: string
	//             timestamp:
	//               type: integer
	//               format: int64
	//             fee:
	//               type: integer
	//               format: int64
	//             version:
	//               type: integer
	//               format: int64
	//         head_outputs:
	//           description: HeadOutputs are unspent outputs confirmed in the blockchain
	//           type: array
	//           items:
	//             properties:
	//               hash:
	//                 type: string
	//               src_tx:
	//                 type: string
	//               address:
	//                 type: string
	//               coins:
	//                 type: string
	//               time:
	//                 type: integer
	//                 format: int64
	//               hours:
	//                 type: integer
	//                 format: int64
	//               calculated_hours:
	//                 type: integer
	//                 format: int64
	//               block_seq:
	//                 type: integer
	//                 format: int64
	//         outgoing_outputs:
	//           description: OutgoingOutputs are unspent outputs being spent in unconfirmed transactions
	//           type: array
	//           items:
	//             properties:
	//               hash:
	//                 type: string
	//               src_tx:
	//                 type: string
	//               address:
	//                 type: string
	//               coins:
	//                 type: string
	//               time:
	//                 type: integer
	//                 format: int64
	//               hours:
	//                 type: integer
	//                 format: int64
	//               calculated_hours:
	//                 type: integer
	//                 format: int64
	//               block_seq:
	//                 type: integer
	//                 format: int64
	//         incoming_outputs:
	//           description: IncomingOutputs are unspent outputs being created by unconfirmed transactions
	//           type: array
	//           items:
	//             properties:
	//               hash:
	//                 type: string
	//               src_tx:
	//                 type: string
	//               address:
	//                 type: string
	//               coins:
	//                 type: string
	//               time:
	//                 type: integer
	//                 format: int64
	//               hours:
	//                 type: integer
	//                 format: int64
	//               calculated_hours:
	//                 type: integer
	//                 format: int64
	//               block_seq:
	//                 type: integer
	//                 format: int64
	//   default:
	//	   $ref: '#/responses/genericError'

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}
		addrStr := r.FormValue("hashes")
		hashStr := r.FormValue("addrs")

		if addrStr != "" && hashStr != "" {
			wh.Error400(w, "addrs and hashes cannot be specified together")
			return
		}

		var filters []visor.OutputsFilter

		if addrStr != "" {
			addrs , err := parseAddressesFromStr(addrStr)
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
