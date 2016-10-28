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
	errCodeRequirePost    = -31000 // Need post

	errMsgParseError     = "Parse error"
	errMsgInvalidRequest = "Invalid Request"
	errMsgMethodNotFound = "Method not found"
	errMsgInvalidParams  = "Invalid params"
	errMsgInternalErr    = "Internal error"
	errMsgRequirePost    = "Need http post"

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
	Jsonrpc string   `json:"jsonrpc"`
	Result  string   `json:"result,omitempty"`
	Error   RPCError `json:"error,omitempty"`
	ID      string   `json:"id"`
}

func makeSuccessResponse(id, result string) Response {
	return Response{
		ID:      id,
		Result:  result,
		Jsonrpc: jsonRPC,
	}
}

func makeErrorResponse(id string, err RPCError) Response {
	return Response{
		ID:      id,
		Error:   err,
		Jsonrpc: jsonRPC,
	}
}

// Start start the webrpc service.
func Start(addr string, c chan struct{}) {
	for {
		select {
		case <-c:
			logger.Info("webrpc quit")
			return
		default:
			mux := http.NewServeMux()
			mux.HandleFunc("/webrpc", rpcHandler)
			logger.Fatal(http.ListenAndServe(addr, mux))
		}
	}
}

// Request rpc request struct
type Request struct {
	Jsonrpc string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  map[string]string `json:"params"`
	ID      string            `json:"id"`
}

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	// var req Request
}
