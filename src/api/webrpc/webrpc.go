package webrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	// ErrCodeParseError parse JSON error
	ErrCodeParseError = -32700 // Parse error	Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.
	// ErrCodeInvalidRequest invalid JSON object format
	ErrCodeInvalidRequest = -32600 // Invalid Request	The JSON sent is not a valid Request object.
	// ErrCodeMethodNotFound unknown method
	ErrCodeMethodNotFound = -32601 // Method not found	The method does not exist / is not available.
	// ErrCodeInvalidParams invalid method parameters
	ErrCodeInvalidParams = -32602 // Invalid params	Invalid method parameter(s).
	// ErrCodeInternalError internal error
	ErrCodeInternalError = -32603 // Internal error	Internal JSON-RPC error.

	// ErrMsgParseError parse error message
	ErrMsgParseError = "Parse error"
	// ErrMsgMethodNotFound method not found message
	ErrMsgMethodNotFound = "Method not found"
	// ErrMsgInvalidParams invalid params message
	ErrMsgInvalidParams = "Invalid params"
	// ErrMsgInternalError internal error message
	ErrMsgInternalError = "Internal error"
	// ErrMsgNotPost not an HTTP POST request message
	ErrMsgNotPost = "only support http POST"
	// ErrMsgInvalidJsonrpc invalid jsonrpc message
	ErrMsgInvalidJsonrpc = "invalid jsonrpc"

	// -32000 to -32099	Server error	Reserved for implementation-defined server-errors.

	jsonRPC = "2.0"

	logger = logging.MustGetLogger("webrpc")
)

// Request rpc request struct
type Request struct {
	ID      string          `json:"id"`
	Jsonrpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// RPCError response error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func (e RPCError) Error() string {
	return fmt.Sprintf("%s [code: %d]", e.Message, e.Code)
}

// Response rpc response struct
type Response struct {
	ID      *string         `json:"id"`
	Jsonrpc string          `json:"jsonrpc"`
	Error   *RPCError       `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

// NewRequest create new webrpc request.
func NewRequest(method string, params interface{}, id string) (*Request, error) {
	var p json.RawMessage
	if params != nil {
		var err error
		p, err = json.Marshal(params)
		if err != nil {
			return nil, err
		}
	}

	return &Request{
		Jsonrpc: jsonRPC,
		Method:  method,
		Params:  p,
		ID:      id,
	}, nil
}

// DecodeParams decodes request params to specific value.
func (r *Request) DecodeParams(v interface{}) error {
	return json.NewDecoder(bytes.NewBuffer(r.Params)).Decode(v)
}

func makeSuccessResponse(id string, result interface{}) Response {
	// TODO -- don't ignore error
	rlt, _ := json.Marshal(result)
	return Response{
		ID:      &id,
		Result:  rlt,
		Jsonrpc: jsonRPC,
	}
}

// MakeErrorResponse creates an error Response
func MakeErrorResponse(code int, msg string, msgs ...string) Response {
	msg = strings.Join(append([]string{msg}, msgs[:]...), "\n")
	return Response{
		Error:   &RPCError{Code: code, Message: msg},
		Jsonrpc: jsonRPC,
	}
}

type operation func(rpc *WebRPC)

// HandlerFunc represents the function type for processing the request
type HandlerFunc func(req Request, gateway Gatewayer) Response

// WebRPC manage the web rpc state and handles
type WebRPC struct {
	Gateway  Gatewayer
	handlers map[string]HandlerFunc
}

// New returns a new WebRPC object
func New(gw Gatewayer) (*WebRPC, error) {
	rpc := &WebRPC{
		Gateway:  gw,
		handlers: make(map[string]HandlerFunc),
	}

	if err := rpc.initHandlers(); err != nil {
		return nil, err
	}

	return rpc, nil
}

// initHandlers initialize webrpc handlers
func (rpc *WebRPC) initHandlers() error {
	handles := map[string]HandlerFunc{
		// get service status
		"get_status": getStatusHandler,
		// get blocks by seq
		"get_blocks_by_seq": getBlocksBySeqHandler,
		// get last N blocks
		"get_lastblocks": getLastBlocksHandler,
		// get blocks in specific seq range
		"get_blocks": getBlocksHandler,
		// get unspent outputs of address
		"get_outputs": getOutputsHandler,
		// get transaction by txid
		"get_transaction": getTransactionHandler,
		// broadcast transaction
		"inject_transaction": injectTransactionHandler,
		// get address affected uxouts
		"get_address_uxouts": getAddrUxOutsHandler,
	}

	// register handlers
	for path, handle := range handles {
		if err := rpc.HandleFunc(path, handle); err != nil {
			return err
		}
	}

	return nil
}

// HandleFunc registers handler function
func (rpc *WebRPC) HandleFunc(method string, h HandlerFunc) error {
	if _, ok := rpc.handlers[method]; ok {
		return fmt.Errorf("%s method already exist", method)
	}

	rpc.handlers[method] = h
	return nil
}

// Handler processes the http request
func (rpc *WebRPC) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res := MakeErrorResponse(ErrCodeInvalidRequest, ErrMsgNotPost)
		logger.Error("Only POST is allowed")
		wh.SendJSONOr500(logger, w, &res)
		return
	}

	req := Request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res := MakeErrorResponse(ErrCodeParseError, ErrMsgParseError)
		logger.WithError(err).Error("Invalid request body")
		wh.SendJSONOr500(logger, w, &res)
		return
	}

	if req.Jsonrpc != jsonRPC {
		res := MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidJsonrpc)
		logger.Error("Invalid JSON-RPC version")
		wh.SendJSONOr500(logger, w, &res)
		return
	}

	var res Response
	if handler, ok := rpc.handlers[req.Method]; ok {
		logger.Infof("Handling method: %s", req.Method)
		res = handler(req, rpc.Gateway)
	} else {
		res = MakeErrorResponse(ErrCodeMethodNotFound, ErrMsgMethodNotFound)
	}

	if res.Error != nil {
		logger.Errorf("%d %s", res.Error.Code, res.Error.Message)
	}

	wh.SendJSONOr500(logger, w, &res)
}
