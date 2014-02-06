// Blockchain related information for the GUI
package gui

import (
    "fmt"
    "github.com/skycoin/skycoin/src/daemon"
    "net/http"
    "strconv"
)

func blockchainHandler(rpc *daemon.RPC) HTTPHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        bc := rpc.GetBlockchainMetadata()
        if bc == nil {
            Error404(w)
        } else if SendJSON(w, bc) != nil {
            Error500(w)
        }
    }
}

func blockchainBlockHandler(rpc *daemon.RPC) HTTPHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        sseq := r.FormValue("seq")
        seq, err := strconv.ParseUint(sseq, 10, 64)
        if err != nil {
            Error400(w, fmt.Sprintf("Invalid seq value \"%s\"", sseq))
            return
        }
        block := rpc.GetBlock(seq)
        if block == nil {
            Error404(w)
        } else if SendJSON(w, block) != nil {
            Error500(w)
        }
    }
}

func blockchainBlocksHandler(rpc *daemon.RPC) HTTPHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        sstart := r.FormValue("start")
        start, err := strconv.ParseUint(sstart, 10, 64)
        if err != nil {
            Error400(w, fmt.Sprintf("Invalid start value \"%s\"", sstart))
            return
        }
        send := r.FormValue("end")
        end, err := strconv.ParseUint(send, 10, 64)
        if err != nil {
            Error400(w, fmt.Sprintf("Invalid end value \"%s\"", send))
            return
        }
        blocks := rpc.GetBlocks(start, end)
        if blocks == nil {
            Error404(w)
        } else if SendJSON(w, blocks) != nil {
            Error500(w)
        }
    }
}

func RegisterBlockchainHandlers(mux *http.ServeMux, rpc *daemon.RPC) {
    mux.HandleFunc("/blockchain", blockchainHandler(rpc))
    mux.HandleFunc("/blockchain/block", blockchainBlockHandler(rpc))
    mux.HandleFunc("/blockchain/blocks", blockchainBlocksHandler(rpc))
}
