package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

// pendingTxnsHandler returns pending (unconfirmed) transactions
// Method: GET
// URI: /api/v1/pendingTxs
// Args:
//	verbose: [bool] include verbose transaction input data
func pendingTxnsHandler(gateway Gatewayer) http.HandlerFunc {

	// swagger:operation GET /api/v1/pendingTxs pendingTxs
	//
	// Returns pending (unconfirmed) transactions
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: verbose
	//   in: query
	//   default: true
	//   description: include verbose transaction input data
	//   required: false
	//   type: boolean
	// responses:
	//   200:
	//     description: Returns a transaction identified by its txid hash.
	//     schema:
	//       type: array
	//       items:
	//         properties:
	//           transaction:
	//             type: object
	//             description: BlockTransactionVerbose has readable transaction data for transactions inside a block. It differs from Transaction in that it includes metadata for transaction inputs and the calculated coinhour fee spent by the block
	//             properties:
	//               length:
	//                 type: integer
	//                 format: int32
	//               type:
	//                 type: integer
	//                 format: int32
	//               hash:
	//                 type: string
	//               inner_hash:
	//                 type: string
	//               fee:
	//                 type: integer
	//                 format: int32
	//               sigs:
	//                 type: array
	//                 items:
	//                   type: string
	//               inputs:
	//                 type: array
	//                 items:
	//                   properties:
	//                     uxid:
	//                       type: string
	//                     dst:
	//                       type: string
	//                     coins:
	//                       type: string
	//                     hours:
	//                       type: integer
	//                       format: int64
	//                     calculated_hours:
	//                       type: integer
	//                       format: int64
	//               outputs:
	//                 type: array
	//                 items:
	//                   properties:
	//                     uxid:
	//                       type: string
	//                     dst:
	//                       type: string
	//                     coins:
	//                       type: string
	//                     hours:
	//                       type: integer
	//                       format: int64
	//           received:
	//             type: string
	//           checked:
	//             type: string
	//           announced:
	//             type: string
	//           is_valid:
	//             type: boolean
	//   default:
	//     $ref: '#/responses/genericError'

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		verbose, err := parseBoolFlag(r.FormValue("verbose"))
		if err != nil {
			wh.Error400(w, "Invalid value for verbose")
			return
		}

		if verbose {
			txns, inputs, err := gateway.GetAllUnconfirmedTransactionsVerbose()
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			vb, err := readable.NewUnconfirmedTransactionsVerbose(txns, inputs)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, vb)
		} else {
			txns, err := gateway.GetAllUnconfirmedTransactions()
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			ret, err := readable.NewUnconfirmedTransactions(txns)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, ret)
		}
	}
}

// TransactionEncodedResponse represents the data struct of the response to /api/v1/transaction?encoded=1
type TransactionEncodedResponse struct {
	Status             readable.TransactionStatus `json:"status"`
	Time               uint64                     `json:"time"`
	EncodedTransaction string                     `json:"encoded_transaction"`
}

// transactionHandler returns a transaction identified by its txid hash
// Method: GET
// URI: /api/v1/transaction
// Args:
//	txid: transaction hash
//	verbose: [bool] include verbose transaction input data
//  encoded: [bool] return as a raw encoded transaction
func transactionHandler(gateway Gatewayer) http.HandlerFunc {

	// swagger:operation GET /api/v1/transaction transaction
	//
	// Returns a transaction identified by its txid hash
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: txid
	//   in: query
	//   description: transaction hash
	//   required: true
	//   type: string
	// - name: encoded
	//   in: query
	//   description: return as a raw encoded transaction.
	//   required: false
	//   type: boolean
	// - name: verbose
	//   in: query
	//   default: true
	//   description: include verbose transaction input data
	//   required: false
	//   type: boolean
	// responses:
	//   200:
	//     description: Returns a transaction identified by its txid hash.
	//     schema:
	//       properties:
	//           status:
	//             type: object
	//             properties:
	//               confirmed:
	//                 type: boolean
	//               unconfirmed:
	//                 type: boolean
	//               height:
	//                 description: If confirmed, how many blocks deep in the chain it is. Will be at least 1 if confirmed
	//                 type: integer
	//                 format: int64
	//               block_seq:
	//                 description: If confirmed, the sequence of the block in which the transaction was executed
	//                 type: integer
	//                 format: int64
	//           time:
	//             type: integer
	//             format: int64
	//           txn:
	//             description: TransactionVerbose has readable transaction data. It adds TransactionStatus to a BlockTransactionVerbose
	//             type: object
	//             properties:
	//               status:
	//                 type: object
	//                 properties:
	//                   confirmed:
	//                      type: boolean
	//                   unconfirmed:
	//                     type: boolean
	//                   height:
	//                     description: If confirmed, how many blocks deep in the chain it is. Will be at least 1 if confirmed
	//                     type: integer
	//                     format: int64
	//                   block_seq:
	//                     description: If confirmed, the sequence of the block in which the transaction was executed
	//                     type: integer
	//                     format: int64
	//               timestamp:
	//                 type: integer
	//                 format: int64
	//               length:
	//                 type: integer
	//                 format: int32
	//               type:
	//                 type: integer
	//                 format: int32
	//               hash:
	//                 type: string
	//               inner_hash:
	//                 type: string
	//               fee:
	//                 type: integer
	//                 format: int32
	//               sigs:
	//                 type: array
	//                 items:
	//                   type: string
	//               inputs:
	//                 type: array
	//                 items:
	//                   properties:
	//                     uxid:
	//                       type: string
	//                     dst:
	//                       type: string
	//                     coins:
	//                       type: string
	//                     hours:
	//                       type: integer
	//                       format: int64
	//                     calculated_hours:
	//                       type: integer
	//                       format: int64
	//               outputs:
	//                 type: array
	//                 items:
	//                   properties:
	//                     uxid:
	//                       type: string
	//                     dst:
	//                       type: string
	//                     coins:
	//                       type: string
	//                     hours:
	//                       type: integer
	//                       format: int64
	//   default:
	//     $ref: '#/responses/genericError'

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		txid := r.FormValue("txid")
		if txid == "" {
			wh.Error400(w, "txid is empty")
			return
		}

		verbose, err := parseBoolFlag(r.FormValue("verbose"))
		if err != nil {
			wh.Error400(w, "Invalid value for verbose")
			return
		}

		encoded, err := parseBoolFlag(r.FormValue("encoded"))
		if err != nil {
			wh.Error400(w, "Invalid value for encoded")
			return
		}

		if verbose && encoded {
			wh.Error400(w, "verbose and encoded cannot be combined")
			return
		}

		h, err := cipher.SHA256FromHex(txid)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		if verbose {
			txn, inputs, err := gateway.GetTransactionVerbose(h)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}
			if txn == nil {
				wh.Error404(w, "")
				return
			}

			rTxn, err := readable.NewTransactionWithStatusVerbose(txn, inputs)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, rTxn)
			return
		}

		txn, err := gateway.GetTransaction(h)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}
		if txn == nil {
			wh.Error404(w, "")
			return
		}

		if encoded {
			txnStr := hex.EncodeToString(txn.Transaction.Serialize())

			wh.SendJSONOr500(logger, w, TransactionEncodedResponse{
				EncodedTransaction: txnStr,
				Status:             readable.NewTransactionStatus(txn.Status),
				Time:               txn.Time,
			})
			return
		}

		rTxn, err := readable.NewTransactionWithStatus(txn)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, rTxn)
	}
}

// TransactionsWithStatus array of transaction results
type TransactionsWithStatus struct {
	Transactions []readable.TransactionWithStatus `json:"txns"`
}

// Sort sorts transactions chronologically, using txid for tiebreaking
func (r TransactionsWithStatus) Sort() {
	sort.Slice(r.Transactions, func(i, j int) bool {
		a := r.Transactions[i]
		b := r.Transactions[j]

		if a.Time == b.Time {
			return strings.Compare(a.Transaction.Hash, b.Transaction.Hash) < 0
		}

		return a.Time < b.Time
	})
}

// NewTransactionsWithStatus converts []Transaction to TransactionsWithStatus
func NewTransactionsWithStatus(txns []visor.Transaction) (*TransactionsWithStatus, error) {
	txnRlts := make([]readable.TransactionWithStatus, 0, len(txns))
	for _, txn := range txns {
		rTxn, err := readable.NewTransactionWithStatus(&txn)
		if err != nil {
			return nil, err
		}
		txnRlts = append(txnRlts, *rTxn)
	}

	return &TransactionsWithStatus{
		Transactions: txnRlts,
	}, nil
}

// TransactionsWithStatusVerbose array of transaction results
type TransactionsWithStatusVerbose struct {
	Transactions []readable.TransactionWithStatusVerbose `json:"txns"`
}

// Sort sorts transactions chronologically, using txid for tiebreaking
func (r TransactionsWithStatusVerbose) Sort() {
	sort.Slice(r.Transactions, func(i, j int) bool {
		a := r.Transactions[i]
		b := r.Transactions[j]

		if a.Time == b.Time {
			return strings.Compare(a.Transaction.Hash, b.Transaction.Hash) < 0
		}

		return a.Time < b.Time
	})
}

// NewTransactionsWithStatusVerbose converts []Transaction to []TransactionsWithStatusVerbose
func NewTransactionsWithStatusVerbose(txns []visor.Transaction, inputs [][]visor.TransactionInput) (*TransactionsWithStatusVerbose, error) {
	if len(txns) != len(inputs) {
		return nil, errors.New("NewTransactionsWithStatusVerbose: len(txns) != len(inputs)")
	}

	txnRlts := make([]readable.TransactionWithStatusVerbose, len(txns))
	for i, txn := range txns {
		rTxn, err := readable.NewTransactionWithStatusVerbose(&txn, inputs[i])
		if err != nil {
			return nil, err
		}
		txnRlts[i] = *rTxn
	}

	return &TransactionsWithStatusVerbose{
		Transactions: txnRlts,
	}, nil
}

// Returns transactions that match the filters.
// Method: GET, POST
// URI: /api/v1/transactions
// Args:
//     addrs: Comma separated addresses [optional, returns all transactions if no address provided]
//     confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
//	   verbose: [bool] include verbose transaction input data
func transactionsHandler(gateway Gatewayer) http.HandlerFunc {

	// swagger:operation GET /api/v1/transactions transactionsGet
	//
	// Returns transactions that match the filters.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: addrs
	//   in: query
	//   description: command separated list of addresses
	//   required: false
	//   type: string
	// - name: confirmed
	//   in: query
	//   description: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
	//   required: false
	//   type: string
	// - name: verbose
	//   in: query
	//   default: true
	//   description: include verbose transaction input data
	//   required: false
	//   type: boolean
	// responses:
	//   200:
	//     description: Returns transactions that match the filters.
	//     schema:
	//       type: array
	//       items:
	//         properties:
	//             status:
	//               type: object
	//               properties:
	//                 confirmed:
	//                   type: boolean
	//                 unconfirmed:
	//                   type: boolean
	//                 height:
	//                   description: If confirmed, how many blocks deep in the chain it is. Will be at least 1 if confirmed
	//                   type: integer
	//                   format: int64
	//                 block_seq:
	//                   description: If confirmed, the sequence of the block in which the transaction was executed
	//                   type: integer
	//                   format: int64
	//             time:
	//               type: integer
	//               format: int64
	//             txn:
	//               description: TransactionVerbose has readable transaction data. It adds TransactionStatus to a BlockTransactionVerbose
	//               type: object
	//               properties:
	//                 status:
	//                   type: object
	//                   properties:
	//                     confirmed:
	//                       type: boolean
	//                     unconfirmed:
	//                       type: boolean
	//                     height:
	//                       description: If confirmed, how many blocks deep in the chain it is. Will be at least 1 if confirmed
	//                       type: integer
	//                       format: int64
	//                     block_seq:
	//                       description: If confirmed, the sequence of the block in which the transaction was executed
	//                       type: integer
	//                       format: int64
	//                 timestamp:
	//                   type: integer
	//                   format: int64
	//                 length:
	//                   type: integer
	//                   format: int32
	//                 type:
	//                   type: integer
	//                   format: int32
	//                 hash:
	//                   type: string
	//                 inner_hash:
	//                   type: string
	//                 fee:
	//                   type: integer
	//                   format: int32
	//                 sigs:
	//                   type: array
	//                   items:
	//                     type: string
	//                 inputs:
	//                   type: array
	//                   items:
	//                     properties:
	//                       uxid:
	//                         type: string
	//                       dst:
	//                         type: string
	//                       coins:
	//                         type: string
	//                       hours:
	//                         type: integer
	//                         format: int64
	//                       calculated_hours:
	//                         type: integer
	//                         format: int64
	//                 outputs:
	//                   type: array
	//                   items:
	//                     properties:
	//                       uxid:
	//                         type: string
	//                       dst:
	//                         type: string
	//                       coins:
	//                         type: string
	//                       hours:
	//                         type: integer
	//                         format: int64
	//   default:
	//     $ref: '#/responses/genericError'

	// swagger:operation POST /api/v1/transactions transactionsPost
	//
	// Returns transactions that match the filters.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: addrs
	//   in: query
	//   description: command separated list of addresses
	//   required: false
	//   type: string
	// - name: confirmed
	//   in: query
	//   description: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
	//   required: false
	//   type: string
	// - name: verbose
	//   in: query
	//   default: true
	//   description: include verbose transaction input data
	//   required: false
	//   type: boolean
	// responses:
	//   200:
	//     description: Returns transactions that match the filters.
	//     schema:
	//       type: array
	//       items:
	//         properties:
	//             status:
	//               type: object
	//               properties:
	//                 confirmed:
	//                   type: boolean
	//                 unconfirmed:
	//                   type: boolean
	//                 height:
	//                   description: If confirmed, how many blocks deep in the chain it is. Will be at least 1 if confirmed
	//                   type: integer
	//                   format: int64
	//                 block_seq:
	//                   description: If confirmed, the sequence of the block in which the transaction was executed
	//                   type: integer
	//                   format: int64
	//             time:
	//               type: integer
	//               format: int64
	//             txn:
	//               description: TransactionVerbose has readable transaction data. It adds TransactionStatus to a BlockTransactionVerbose
	//               type: object
	//               properties:
	//                 status:
	//                   type: object
	//                   properties:
	//                     confirmed:
	//                       type: boolean
	//                     unconfirmed:
	//                       type: boolean
	//                     height:
	//                       description: If confirmed, how many blocks deep in the chain it is. Will be at least 1 if confirmed
	//                       type: integer
	//                       format: int64
	//                     block_seq:
	//                       description: If confirmed, the sequence of the block in which the transaction was executed
	//                       type: integer
	//                       format: int64
	//                 timestamp:
	//                   type: integer
	//                   format: int64
	//                 length:
	//                   type: integer
	//                   format: int32
	//                 type:
	//                   type: integer
	//                   format: int32
	//                 hash:
	//                   type: string
	//                 inner_hash:
	//                   type: string
	//                 fee:
	//                   type: integer
	//                   format: int32
	//                 sigs:
	//                   type: array
	//                   items:
	//                     type: string
	//                 inputs:
	//                   type: array
	//                   items:
	//                     properties:
	//                       uxid:
	//                         type: string
	//                       dst:
	//                         type: string
	//                       coins:
	//                         type: string
	//                       hours:
	//                         type: integer
	//                         format: int64
	//                       calculated_hours:
	//                         type: integer
	//                         format: int64
	//                 outputs:
	//                   type: array
	//                   items:
	//                     properties:
	//                       uxid:
	//                         type: string
	//                       dst:
	//                         type: string
	//                       coins:
	//                         type: string
	//                       hours:
	//                         type: integer
	//                         format: int64
	//   default:
	//     $ref: '#/responses/genericError'

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		verbose, err := parseBoolFlag(r.FormValue("verbose"))
		if err != nil {
			wh.Error400(w, "Invalid value for verbose")
			return
		}

		// Gets 'addrs' parameter value
		addrs, err := parseAddressesFromStr(r.FormValue("addrs"))
		if err != nil {
			wh.Error400(w, fmt.Sprintf("parse parameter: 'addrs' failed: %v", err))
			return
		}

		// Initialize transaction filters
		flts := []visor.TxFilter{visor.NewAddrsFilter(addrs)}

		// Gets the 'confirmed' parameter value
		confirmedStr := r.FormValue("confirmed")
		if confirmedStr != "" {
			confirmed, err := strconv.ParseBool(confirmedStr)
			if err != nil {
				wh.Error400(w, fmt.Sprintf("invalid 'confirmed' value: %v", err))
				return
			}

			flts = append(flts, visor.NewConfirmedTxFilter(confirmed))
		}

		if verbose {
			txns, inputs, err := gateway.GetTransactionsVerbose(flts)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			rTxns, err := NewTransactionsWithStatusVerbose(txns, inputs)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			rTxns.Sort()

			wh.SendJSONOr500(logger, w, rTxns.Transactions)
		} else {
			txns, err := gateway.GetTransactions(flts)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			rTxns, err := NewTransactionsWithStatus(txns)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			rTxns.Sort()

			wh.SendJSONOr500(logger, w, rTxns.Transactions)
		}
	}
}

// RawTxnData used in requests and responses including raw transaction data
type RawTxnData struct {
	Rawtx string `json:"rawtx"`
}

// URI: /api/v1/injectTransaction
// Method: POST
// Content-Type: application/json
// Body: {"rawtx": "<hex encoded transaction>"}
// Response:
//      200 - ok, returns the transaction hash in hex as string
//      400 - bad transaction
//		500 - other error
//      503 - network unavailable for broadcasting transaction
func injectTransactionHandler(gateway Gatewayer, forAPIVersion2 bool) http.HandlerFunc {

	// swagger:operation POST /api/v1/injectTransaction injectTransaction
	//
	// Broadcast a hex-encoded, serialized transaction to the network.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: rawtx
	//   in: header
	//   description: hex-encoded serialized transaction string.
	//   required: true
	//   type: string
	// responses:
	//   200:
	//     description: This endpoint a hex-encoded transaction to the network.
	//     type: string
	//   default:
	//     $ref: '#/responses/genericError'

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			if forAPIVersion2 {
				resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
				writeHTTPResponse(w, resp)
			} else {
				wh.Error405(w)
			}
			return
		}
		// get the rawtransaction
		var v RawTxnData

		if forAPIVersion2 && r.Header.Get("Content-Type") != ContentTypeJSON {
			resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
			writeHTTPResponse(w, resp)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			if forAPIVersion2 {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
				writeHTTPResponse(w, resp)
			} else {
				wh.Error400(w, err.Error())
			}
			return
		}

		b, err := hex.DecodeString(v.Rawtx)
		if err != nil {
			if forAPIVersion2 {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
				writeHTTPResponse(w, resp)
			} else {
				wh.Error400(w, err.Error())
			}
			return
		}

		txn, err := coin.TransactionDeserialize(b)
		if err != nil {
			if forAPIVersion2 {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
				writeHTTPResponse(w, resp)
			} else {
				wh.Error400(w, err.Error())
			}
			return
		}

		if err := gateway.InjectBroadcastTransaction(txn); err != nil {
			if daemon.IsBroadcastFailure(err) {
				if forAPIVersion2 {
					resp := NewHTTPErrorResponse(http.StatusForbidden, "")
					writeHTTPResponse(w, resp)
				} else {
					wh.Error503(w, err.Error())
				}
			} else {
				if forAPIVersion2 {
					resp := NewHTTPErrorResponse(http.StatusInternalServerError, "")
					writeHTTPResponse(w, resp)
				} else {
					wh.Error500(w, err.Error())
				}
			}
			return
		}

		if forAPIVersion2 {
			var resp HTTPResponse
			rTxn, err := readable.NewTransaction(txn, false)
			if err != nil {
				resp = NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
			resp.Data = rTxn
			writeHTTPResponse(w, resp)
		} else {
			wh.SendJSONOr500(logger, w, txn.Hash().Hex())
		}
	}
}

// ResendResult the result of rebroadcasting transaction
// swagger:response resendResult
type ResendResult struct {
	Txids []string `json:"txids"`
}

// NewResendResult creates a ResendResult from a list of transaction ID hashes
func NewResendResult(hashes []cipher.SHA256) ResendResult {
	txids := make([]string, len(hashes))
	for i, h := range hashes {
		txids[i] = h.Hex()
	}
	return ResendResult{
		Txids: txids,
	}
}

// URI: /api/v1/resendUnconfirmedTxns
// Method: POST
// Broadcasts all unconfirmed transactions from the unconfirmed transaction pool
// Response:
//      200 - ok, returns the transaction hashes that were resent
//      405 - method not POST
//		500 - other error
//      503 - network unavailable for broadcasting transaction
func resendUnconfirmedTxnsHandler(gateway Gatewayer) http.HandlerFunc {

	// swagger:route POST /api/v1/resendUnconfirmedTxns resendUnconfirmedTxns
	//
	// Broadcasts all unconfirmed transactions from the unconfirmed transaction pool
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Responses:
	//       default: genericError
	//       200: resendResult

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		hashes, err := gateway.ResendUnconfirmedTxns()
		if err != nil {
			if daemon.IsBroadcastFailure(err) {
				wh.Error503(w, err.Error())
			} else {
				wh.Error500(w, err.Error())
			}
			return
		}

		wh.SendJSONOr500(logger, w, NewResendResult(hashes))
	}
}

// URI: /api/v1/rawtx
// Method: GET
// Args:
//	txid: transaction ID hash
//	forAPIVersion2: return rawdata in JSON object (instead of plain/text string)
// Returns the hex-encoded byte serialization of a transaction.
// The transaction may be confirmed or unconfirmed.
func rawTxnHandler(gateway Gatewayer, forAPIVersion2 bool) http.HandlerFunc {

	// swagger:operation GET /api/v1/rawtx rawtx
	//
	// Returns the hex-encoded byte serialization of a transaction. The transaction may be confirmed or unconfirmed.
	//
	// ---
	//
	// produces:
	// - application/json
	// parameters:
	// - name: txid
	//   in: query
	//   description: Transaction id hash
	//   type: string
	//   x-go-name: txid
	//
	// responses:
	//   200:
	//     description: Returns the hex-encoded byte serialization of a transaction
	//     properties:
	//       type: string
	//   default:
	//     $ref: '#/responses/genericError'

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			if forAPIVersion2 {
				resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
				writeHTTPResponse(w, resp)
			} else {
				wh.Error405(w)
			}
			return
		}

		if forAPIVersion2 && r.Header.Get("Content-Type") != ContentTypeJSON {
			resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
			writeHTTPResponse(w, resp)
			return
		}

		txid := r.FormValue("txid")
		if txid == "" {
			if forAPIVersion2 {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, "txid is empty")
				writeHTTPResponse(w, resp)
			} else {
				wh.Error400(w, "txid is empty")
			}
			return
		}

		h, err := cipher.SHA256FromHex(txid)
		if err != nil {
			if forAPIVersion2 {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
				writeHTTPResponse(w, resp)
			} else {
				wh.Error400(w, err.Error())
			}
			return
		}

		txn, err := gateway.GetTransaction(h)
		if err != nil {
			if forAPIVersion2 {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
				writeHTTPResponse(w, resp)
			} else {
				wh.Error400(w, err.Error())
			}
			return
		}

		if txn == nil {
			if forAPIVersion2 {
				resp := NewHTTPErrorResponse(http.StatusNotFound, err.Error())
				writeHTTPResponse(w, resp)
			} else {
				wh.Error404(w, "")
			}
			return
		}

		d := txn.Transaction.Serialize()
		d2hex := hex.EncodeToString(d)
		if forAPIVersion2 {
			var resp HTTPResponse
			resp.Data = RawTxnData{
				Rawtx: d2hex,
			}
			writeHTTPResponse(w, resp)
		} else {
			wh.SendJSONOr500(logger, w, d2hex)
		}
	}
}

// VerifyTxnRequest represents the data struct of the request for /api/v2/transaction/verify
type VerifyTxnRequest struct {
	EncodedTransaction string `json:"encoded_transaction"`
}

// VerifyTxnResponse the response data struct for /api/v2/transaction/verify
type VerifyTxnResponse struct {
	Confirmed   bool               `json:"confirmed"`
	Transaction CreatedTransaction `json:"transaction"`
}

// Decode and verify an encoded transaction
// Method: POST
// URI: /api/v2/transaction/verify
func verifyTxnHandler(gateway Gatewayer) http.HandlerFunc {

	// TODO For v3

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		if r.Header.Get("Content-Type") != ContentTypeJSON {
			resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req VerifyTxnRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		txn, err := decodeTxn(req.EncodedTransaction)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("decode transaction failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		var resp HTTPResponse
		inputs, isTxnConfirmed, err := gateway.VerifyTxnVerbose(txn)
		if err != nil {
			switch err.(type) {
			case visor.ErrTxnViolatesSoftConstraint,
				visor.ErrTxnViolatesHardConstraint,
				visor.ErrTxnViolatesUserConstraint:
				resp.Error = &HTTPError{
					Code:    http.StatusUnprocessableEntity,
					Message: err.Error(),
				}
			default:
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
		}

		verifyTxnResp := VerifyTxnResponse{
			Confirmed: isTxnConfirmed,
		}

		if len(inputs) != len(txn.In) {
			inputs = nil
		}
		verboseTxn, err := newCreatedTransactionFuzzy(txn, inputs)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		verifyTxnResp.Transaction = *verboseTxn

		resp.Data = verifyTxnResp

		if isTxnConfirmed && resp.Error == nil {
			resp.Error = &HTTPError{
				Code:    http.StatusUnprocessableEntity,
				Message: "transaction has been spent",
			}
		}

		writeHTTPResponse(w, resp)
	}
}

func decodeTxn(encodedTxn string) (*coin.Transaction, error) {
	var txn coin.Transaction
	b, err := hex.DecodeString(encodedTxn)
	if err != nil {
		return nil, err
	}

	txn, err = coin.TransactionDeserialize(b)
	if err != nil {
		return nil, err
	}

	return &txn, nil
}

// newCreatedTransactionFuzzy creates a CreatedTransaction but accommodates possibly invalid txn input
func newCreatedTransactionFuzzy(txn *coin.Transaction, inputs []wallet.UxBalance) (*CreatedTransaction, error) {
	if len(txn.In) != len(inputs) && len(inputs) != 0 {
		return nil, errors.New("len(txn.In) != len(inputs)")
	}

	var outputHours uint64
	var feeInvalid bool
	for _, o := range txn.Out {
		var err error
		outputHours, err = coin.AddUint64(outputHours, o.Hours)
		if err != nil {
			feeInvalid = true
		}
	}

	var inputHours uint64
	for _, i := range inputs {
		var err error
		inputHours, err = coin.AddUint64(inputHours, i.Hours)
		if err != nil {
			feeInvalid = true
		}
	}

	if inputHours < outputHours {
		feeInvalid = true
	}

	var fee uint64
	if !feeInvalid {
		fee = inputHours - outputHours
	}

	sigs := make([]string, len(txn.Sigs))
	for i, s := range txn.Sigs {
		sigs[i] = s.Hex()
	}

	txid := txn.Hash()
	out := make([]CreatedTransactionOutput, len(txn.Out))
	for i, o := range txn.Out {
		co, err := NewCreatedTransactionOutput(o, txid)
		if err != nil {
			logger.WithError(err).Error("NewCreatedTransactionOutput failed")
			continue
		}
		out[i] = *co
	}

	in := make([]CreatedTransactionInput, len(txn.In))
	if len(inputs) == 0 {
		for i, h := range txn.In {
			in[i] = CreatedTransactionInput{
				UxID: h.Hex(),
			}
		}
	} else {
		for i, o := range inputs {
			ci, err := NewCreatedTransactionInput(o)
			if err != nil {
				logger.WithError(err).Error("NewCreatedTransactionInput failed")
				continue
			}
			in[i] = *ci
		}
	}

	return &CreatedTransaction{
		Length:    txn.Length,
		Type:      txn.Type,
		TxID:      txid.Hex(),
		InnerHash: txn.InnerHash.Hex(),
		Fee:       fmt.Sprint(fee),

		Sigs: sigs,
		In:   in,
		Out:  out,
	}, nil
}
