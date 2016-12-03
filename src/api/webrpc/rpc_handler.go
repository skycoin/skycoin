package webrpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	wh "github.com/skycoin/skycoin/src/util/http"
)

type operation func(rpc *rpcHandler)

type jobHandler func(req Request, gateway Gatewayer) Response

type rpcHandler struct {
	workerNum uint
	ops       chan operation // request channel
	// reqChan chan
	close    chan struct{}
	mux      *http.ServeMux
	handlers map[string]jobHandler
	gateway  Gatewayer
}

// Arg is the argument type for creating webrpc instance.
type Arg func(*rpcHandler)

func newRPCHandler(args ...Arg) *rpcHandler {
	rpc := &rpcHandler{}
	for _, arg := range args {
		arg(rpc)
	}

	rpc.handlers = make(map[string]jobHandler)
	rpc.mux = http.NewServeMux()

	rpc.mux.HandleFunc("/webrpc", rpc.Handler)
	rpc.dispatch()
	return rpc
}

func (rh *rpcHandler) HandlerFunc(method string, jh jobHandler) {
	if _, ok := rh.handlers[method]; ok {
		logger.Panicf("%s method already exist", method)
	}
	rh.handlers[method] = jh
}

func (rh *rpcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rh.mux.ServeHTTP(w, r)
}

func (rh *rpcHandler) Handler(w http.ResponseWriter, r *http.Request) {
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
		rh.ops <- func(rpc *rpcHandler) {
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
func (rh *rpcHandler) dispatch() {
	for i := uint(0); i < rh.workerNum; i++ {
		go func(seq uint) {
			for {
				select {
				case <-rh.close:
					// logger.Infof("[%d]rpc job handler quit", seq)
					return
				case op := <-rh.ops:
					op(rh)
				}
			}
		}(i)
	}
}

// ChanBuffSize set request channel buffer size
func ChanBuffSize(n uint) Arg {
	return func(rpc *rpcHandler) {
		rpc.ops = make(chan operation, n)
	}
}

// ThreadNum set concurrent request processor number
func ThreadNum(n uint) Arg {
	return func(rpc *rpcHandler) {
		if n == 0 {
			panic("thread num must > 0")
		}
		rpc.workerNum = n
	}
}

// Gateway set gateway
func Gateway(gateway Gatewayer) Arg {
	return func(rpc *rpcHandler) {
		rpc.gateway = gateway
	}
}

// Quit set closing channel
func Quit(c chan struct{}) Arg {
	return func(rpc *rpcHandler) {
		rpc.close = c
	}
}
