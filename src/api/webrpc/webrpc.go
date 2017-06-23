package webrpc

import (
	"fmt"
	"net"
	"net/http"

	"encoding/json"

	"github.com/skycoin/skycoin/src/util"
	wh "github.com/skycoin/skycoin/src/util/http"

	"bytes"
	"strings"
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

var logger = util.MustGetLogger("webrpc")

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

func makeErrorResponse(code int, msgs ...string) Response {
	msg := strings.Join(msgs[:], "\n")
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
	addr      string // service address
	workerNum uint
	ops       chan operation // request channel
	close     chan struct{}
	mux       *http.ServeMux
	handlers  map[string]HandlerFunc
	gateway   Gatewayer
}

// Option is the argument type for creating webrpc instance.
type Option func(*WebRPC)

// New creates webrpc instance
func New(addr string, ops ...Option) (*WebRPC, error) {
	rpc := &WebRPC{
		addr: addr,
	}

	for _, opt := range ops {
		opt(rpc)
	}

	rpc.handlers = make(map[string]HandlerFunc)
	rpc.mux = http.NewServeMux()

	rpc.mux.HandleFunc("/webrpc", rpc.Handler)

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

// Run starts the webrpc service.
func (rpc *WebRPC) Run(quit chan struct{}) {
	logger.Infof("start webrpc on http://%s", rpc.addr)

	l, err := net.Listen("tcp", rpc.addr)
	if err != nil {
		logger.Error("%v", err)
		close(quit)
		return
	}

	c := make(chan struct{})
	q := make(chan struct{}, 1)
	go func() {
		if err := http.Serve(l, rpc); err != nil {
			select {
			case <-c:
				return
			default:
				// the webrpc service failed unexpectly, notify the
				logger.Error("%v", err)
				q <- struct{}{}
			}
		}
	}()

	select {
	case <-quit:
		close(c)
		l.Close()
	case <-q:
		close(quit)
	}
	logger.Info("webrpc quit")
	return
}

// HandleFunc registers handler function
func (rpc *WebRPC) HandleFunc(method string, h HandlerFunc) error {
	if _, ok := rpc.handlers[method]; ok {
		return fmt.Errorf("%s method already exist", method)
	}

	rpc.handlers[method] = h
	return nil
}

// ServHTTP implements the interface of http.Handler
func (rpc *WebRPC) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rpc.dispatch()
	rpc.mux.ServeHTTP(w, r)
}

// Handler processes the http request
func (rpc *WebRPC) Handler(w http.ResponseWriter, r *http.Request) {
	var (
		req Request
		res Response
	)

	for {
		// only support post.
		if r.Method != "POST" {
			res = makeErrorResponse(errCodeInvalidRequest, errMsgNotPost)
			break
		}

		// deocder request.
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			res = makeErrorResponse(errCodeParseError, errMsgParseError)
			break
		}

		if req.Jsonrpc != jsonRPC {
			res = makeErrorResponse(errCodeInvalidParams, errMsgInvalidJsonrpc)
			break
		}

		resC := make(chan Response)
		rpc.ops <- func(rpc *WebRPC) {
			defer func() {
				if r := recover(); r != nil {
					logger.Critical(fmt.Sprintf("%v", r))
					resC <- makeErrorResponse(errCodeInternalError, errMsgInternalError)
				}
			}()
			if handler, ok := rpc.handlers[req.Method]; ok {
				logger.Info("method: %v", req.Method)
				resC <- handler(req, rpc.gateway)
				return
			}
			resC <- makeErrorResponse(errCodeMethodNotFound, errMsgMethodNotFound)
		}
		res = <-resC
		break
	}

	wh.SendOr404(w, &res)
}

// dispatch will create numbers of goroutines, each routine will
func (rpc *WebRPC) dispatch() {
	for i := uint(0); i < rpc.workerNum; i++ {
		go func(seq uint) {
			for {
				select {
				case <-rpc.close:
					// logger.Infof("[%d]rpc job handler quit", seq)
					return
				case op := <-rpc.ops:
					func() {
						defer func() {
							if r := recover(); r != nil {
								logger.Error("recover: %v", r)
							}
						}()

						op(rpc)
					}()
				}
			}
		}(i)
	}
}

// ChanBuffSize set request channel buffer size
func ChanBuffSize(n uint) Option {
	return func(rpc *WebRPC) {
		rpc.ops = make(chan operation, n)
	}
}

// ThreadNum set concurrent request processor number
func ThreadNum(n uint) Option {
	return func(rpc *WebRPC) {
		rpc.workerNum = n
	}
}

// Gateway set gateway
func Gateway(gateway Gatewayer) Option {
	return func(rpc *WebRPC) {
		rpc.gateway = gateway
	}
}

// Quit set closing channel
func Quit(c chan struct{}) Option {
	return func(rpc *WebRPC) {
		rpc.close = c
	}
}
