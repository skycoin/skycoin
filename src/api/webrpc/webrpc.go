package webrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	errCodeParseError     = -32700 // Parse error	Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.
	errCodeInvalidRequest = -32600 // Invalid Request	The JSON sent is not a valid Request object.
	errCodeMethodNotFound = -32601 // Method not found	The method does not exist / is not available.
	errCodeInvalidParams  = -32602 // Invalid params	Invalid method parameter(s).
	errCodeInternalError  = -32603 // Internal error	Internal JSON-RPC error.

	errMsgParseError     = "Parse error"
	errMsgMethodNotFound = "Method not found"
	errMsgInvalidParams  = "Invalid params"
	errMsgInternalError  = "Internal error"
	errMsgNotPost        = "only support http POST"
	errMsgInvalidJsonrpc = "invalid jsonrpc"

	// -32000 to -32099	Server error	Reserved for implementation-defined server-errors.

	jsonRPC = "2.0"

	logger = logging.MustGetLogger("webrpc")
)

const (
	defaultReadTimeout  = time.Second * 10
	defaultWriteTimeout = time.Second * 60
	defaultIdleTimeout  = time.Second * 120
	defaultWorkerNum    = 5
	defaultChanBuffSize = 1000
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
	rlt, _ := json.Marshal(result)
	return Response{
		ID:      &id,
		Result:  rlt,
		Jsonrpc: jsonRPC,
	}
}

func makeErrorResponse(code int, msg string, msgs ...string) Response {
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
	Addr         string // service address
	Gateway      Gatewayer
	WorkerNum    uint
	ChanBuffSize uint // size of ops channel

	ops      chan operation // request channel
	mux      *http.ServeMux
	server   *http.Server
	handlers map[string]HandlerFunc
	listener net.Listener
	quit     chan struct{}
}

// Config configures the WebRPC
type Config struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	WorkerNum    uint
	ChanBuffSize uint // size of ops channel
}

// New returns a new WebRPC object
func New(addr string, c Config, gw Gatewayer) (*WebRPC, error) {
	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaultReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = defaultIdleTimeout
	}
	if c.WorkerNum == 0 {
		c.WorkerNum = defaultWorkerNum
	}
	if c.ChanBuffSize == 0 {
		c.ChanBuffSize = defaultChanBuffSize
	}

	mux := http.NewServeMux()

	rpc := &WebRPC{
		Addr:         addr,
		Gateway:      gw,
		WorkerNum:    c.WorkerNum,
		ChanBuffSize: c.ChanBuffSize,
		quit:         make(chan struct{}),
		mux:          mux,
		server: &http.Server{
			Handler:      mux,
			ReadTimeout:  c.ReadTimeout,
			WriteTimeout: c.WriteTimeout,
			IdleTimeout:  c.IdleTimeout,
		},
		handlers: make(map[string]HandlerFunc),
	}

	mux.Handle("/webrpc", wh.HostCheck(logger, addr, http.HandlerFunc(rpc.Handler)))

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
func (rpc *WebRPC) Run() error {
	if rpc.WorkerNum < 1 {
		return errors.New("rpc.WorkerNum must be > 0")
	}

	if rpc.ChanBuffSize < 1 {
		return errors.New("rpc.ChanBuffSize must be > 0")
	}

	logger.Infof("Start webrpc on http://%s", rpc.Addr)
	defer logger.Info("Webrpc service closed")

	var err error
	if rpc.listener, err = net.Listen("tcp", rpc.Addr); err != nil {
		return err
	}

	rpc.ops = make(chan operation, rpc.ChanBuffSize)

	for i := uint(0); i < rpc.WorkerNum; i++ {
		go rpc.workerThread(i)
	}

	errC := make(chan error, 1)
	go func() {
		if err := rpc.server.Serve(rpc.listener); err != nil {
			select {
			case <-rpc.quit:
				errC <- nil
			default:
				// the webrpc service failed unexpectedly
				logger.Info("webrpc.Run, http.Serve error:", err)
				errC <- err
			}
		}
	}()

	return <-errC
}

// Shutdown close the webrpc service
func (rpc *WebRPC) Shutdown() error {
	if rpc.quit != nil {
		close(rpc.quit)
	}

	if rpc.listener != nil {
		return rpc.listener.Close()
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
	// only support post.
	if r.Method != http.MethodPost {
		res := makeErrorResponse(errCodeInvalidRequest, errMsgNotPost)
		wh.SendJSONOr500(logger, w, &res)
		return
	}

	// decoder request.
	req := Request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res := makeErrorResponse(errCodeParseError, errMsgParseError)
		wh.SendJSONOr500(logger, w, &res)
		return
	}

	if req.Jsonrpc != jsonRPC {
		res := makeErrorResponse(errCodeInvalidParams, errMsgInvalidJsonrpc)
		wh.SendJSONOr500(logger, w, &res)
		return
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
			logger.Info("webrpc handling method: %v", req.Method)
			resC <- handler(req, rpc.Gateway)
		} else {
			resC <- makeErrorResponse(errCodeMethodNotFound, errMsgMethodNotFound)
		}
	}

	res := <-resC
	wh.SendJSONOr500(logger, w, &res)
}

func (rpc *WebRPC) workerThread(seq uint) {
	for {
		select {
		case <-rpc.quit:
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
}
