package webrpc

import (
	"net/http"

	logging "github.com/op/go-logging"
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
	errMsgInternalErr    = "Internal error"

	errMsgNotPost = "only support http POST"

	errMsgInvalidJsonrpc = "invalid jsonrpc"

	// -32000 to -32099	Server error	Reserved for implementation-defined server-errors.

	jsonRPC = "2.0"
)

var logger = logging.MustGetLogger("skycoin.webrpc")

// Request rpc request struct
type Request struct {
	Jsonrpc string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  map[string]string `json:"params"`
	ID      string            `json:"id"`
}

// RPCError response error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// Response rpc response struct
type Response struct {
	Jsonrpc string    `json:"jsonrpc"`
	Error   *RPCError `json:"error,omitempty"`
	Result  string    `json:"result,omitempty"`
	ID      string    `json:"id"`
}

// NewRequest create new webrpc request.
func NewRequest(method string, params map[string]string, id string) *Request {
	return &Request{
		Jsonrpc: jsonRPC,
		Method:  method,
		Params:  params,
		ID:      id,
	}
}

func makeSuccessResponse(id, result string) Response {
	return Response{
		ID:      id,
		Result:  result,
		Jsonrpc: jsonRPC,
	}
}

func makeErrorResponse(id string, err *RPCError) Response {
	return Response{
		ID:      id,
		Error:   err,
		Jsonrpc: jsonRPC,
	}
}

// Start start the webrpc service.
func Start(addr string, queueSize int, workerNum int, gateway Gatewayer, c chan struct{}) {
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

func makeRPC(queueSize int, workerNum int, gateway Gatewayer, c chan struct{}) *rpcHandler {
	rpc := newRPCHandler(queueSize, workerNum, gateway, c)

	// register handlers
	rpc.HandlerFunc("get_status", getStatus)
	rpc.HandlerFunc("get_lastblocks", getLastBlocks)
	return rpc
}
