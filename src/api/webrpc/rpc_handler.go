package webrpc

import (
	"encoding/json"
	"net/http"

	wh "github.com/skycoin/skycoin/src/util/http"
)

type job struct {
	Req  Request
	ResC chan Response
}

func makeJob(req Request) job {
	return job{
		Req:  req,
		ResC: make(chan Response),
	}
}

type jobHandler func(req Request, gateway Gatewayer) Response

type rpcHandler struct {
	workerNum uint
	reqChan   chan job // request channel
	close     chan struct{}
	mux       *http.ServeMux
	handlers  map[string]jobHandler
	gateway   Gatewayer
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

		// make job
		jb := makeJob(req)
		// push job to handler channel.
		rh.reqChan <- jb
		// get response from channel.
		res = <-jb.ResC
		break
	}

	wh.SendOr404(w, &res)
}

// dispatch will create numbers of goroutines, each routine will
//
func (rh *rpcHandler) dispatch() {
	for i := uint(0); i < rh.workerNum; i++ {
		go func(seq uint) {
			var (
				handler jobHandler
				ok      bool
			)

			for {
				select {
				case <-rh.close:
					// logger.Infof("[%d]rpc job handler quit", seq)
					return
				case jb := <-rh.reqChan:
					logger.Debugf("[%d] got job", seq)
					if handler, ok = rh.handlers[jb.Req.Method]; ok {
						logger.Info("method: %v", jb.Req.Method)
						jb.ResC <- handler(jb.Req, rh.gateway)
						logger.Debugf("[%d] job done", seq)
						continue
					}

					jb.ResC <- makeErrorResponse(errCodeMethodNotFound, errMsgMethodNotFound)
					logger.Debugf("[%d] job done", seq)
				}
			}
		}(i)
	}
}

// ChanBuffSize set request channel buffer size
func ChanBuffSize(n uint) Arg {
	return func(rpc *rpcHandler) {
		rpc.reqChan = make(chan job, n)
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
