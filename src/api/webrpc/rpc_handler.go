package webrpc

import (
	"encoding/json"
	"net/http"

	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/visor"
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

// Gatewayer provides interfaces for getting skycoin related info.
type Gatewayer interface {
	GetLastBlocks(num uint64) *visor.ReadableBlocks
	GetBlocks(start, end uint64) *visor.ReadableBlocks
}

type jobHandler func(req Request, gateway Gatewayer) Response

type rpcHandler struct {
	workerNum int
	reqChan   chan job // request channel
	close     chan struct{}
	mux       *http.ServeMux
	handlers  map[string]jobHandler
	gateway   Gatewayer
}

// create rpc handler instance.
func newRPCHandler(queueSize int, workerNum int, gateway Gatewayer, close chan struct{}) *rpcHandler {
	rpc := &rpcHandler{
		workerNum: workerNum,
		reqChan:   make(chan job, queueSize),
		close:     close,
		mux:       http.NewServeMux(),
		handlers:  make(map[string]jobHandler),
		gateway:   gateway,
	}

	rpc.mux.HandleFunc("/webrpc", rpc.Handler)
	rpc.dispatch()
	return rpc
}

func (rh *rpcHandler) HandlerFunc(method string, jh jobHandler) {
	if _, ok := rh.handlers[method]; ok {
		logger.Fatalf("%s method already exist", method)
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
			res = makeErrorResponse("", &RPCError{
				Code:    errCodeInvalidRequest,
				Message: errMsgNotPost,
			})
			break
		}

		// deocder request.
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			res = makeErrorResponse("", &RPCError{
				Code:    errCodeParseError,
				Message: errMsgParseError,
			})
			break
		}

		if req.Jsonrpc != jsonRPC {
			res = makeErrorResponse("", &RPCError{
				Code:    errCodeInvalidParams,
				Message: errMsgInvalidJsonrpc,
			})
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
	for i := 0; i < rh.workerNum; i++ {
		go func(seq int) {
			var (
				handler jobHandler
				ok      bool
				res     Response
			)

			for {
				select {
				case <-rh.close:
					// logger.Infof("[%d]rpc job handler quit", seq)
					return
				case jb := <-rh.reqChan:
					logger.Debugf("[%d] got job", seq)
					if handler, ok = rh.handlers[jb.Req.Method]; ok {
						jb.ResC <- handler(jb.Req, rh.gateway)
						logger.Debugf("[%d] job done", seq)
						continue
					}

					res.ID = jb.Req.ID
					res.Jsonrpc = jsonRPC
					res.Error = &RPCError{
						Code:    errCodeMethodNotFound,
						Message: errMsgMethodNotFound,
					}
					jb.ResC <- res
					logger.Debugf("[%d] job done", seq)
				}
			}
		}(i)
	}
}
