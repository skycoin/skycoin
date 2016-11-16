package webrpc

import (
	"net/http"

	logging "github.com/op/go-logging"

	"encoding/json"

	"bytes"
)

var (
	errCodeParseError     = -32700 // Parse error	Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.
	errCodeInvalidRequest = -32600 // Invalid Request	The JSON sent is not a valid Request object.
	errCodeMethodNotFound = -32601 // Method not found	The method does not exist / is not available.
	errCodeInvalidParams  = -32602 // Invalid params	Invalid method parameter(s).
	errCodeInternalError  = -32603 // Internal error	Internal JSON-RPC error.

	errMsgParseError     = "Parse error"
	errMsgInvalidRequest = "Invalid Request"
	errMsgMethodNotFound = "Method not found"
	errMsgInvalidParams  = "Invalid params"
	errMsgInternalError  = "Internal error"

	errMsgNotPost = "only support http POST"

	errMsgInvalidJsonrpc = "invalid jsonrpc"

	// -32000 to -32099	Server error	Reserved for implementation-defined server-errors.

	jsonRPC = "2.0"
)

var logger = logging.MustGetLogger("webrpc")

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
	rlt, _ := json.Marshal(result)
	return Response{
		ID:      &id,
		Result:  rlt,
		Jsonrpc: jsonRPC,
	}
}

func makeErrorResponse(code int, message string) Response {
	return Response{
		Error:   &RPCError{Code: code, Message: message},
		Jsonrpc: jsonRPC,
	}
}

// Start start the webrpc service.
func Start(addr string, queueSize uint, workerNum uint, gateway Gatewayer, c chan struct{}) {
	rpc := makeRPC(queueSize, workerNum, gateway, c)
	for {
		select {
		case <-c:
			logger.Info("webrpc quit")
			return
		default:
			logger.Infof("start webrpc on http://%s", addr)
			logger.Fatal(http.ListenAndServe(addr, rpc))
		}
	}
}

func makeRPC(queueSize uint, workerNum uint, gateway Gatewayer, c chan struct{}) *rpcHandler {
	rpc := newRPCHandler(queueSize, workerNum, gateway, c)

	// register handlers
	rpc.HandlerFunc("get_status", getStatusHandler)
	rpc.HandlerFunc("get_lastblocks", getLastBlocksHandler)
	rpc.HandlerFunc("get_blocks", getBlocksHandler)
	rpc.HandlerFunc("get_outputs", getOutputsHandler)
	rpc.HandlerFunc("get_transaction", getTransactionHandler)
	rpc.HandlerFunc("inject_transaction", injectTransactionHandler)

	return rpc
}
