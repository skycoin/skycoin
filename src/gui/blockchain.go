// Blockchain related information for the GUI
package gui

import (
    "fmt"
    "github.com/skycoin/skycoin/src/daemon"
    "net/http"
    "strconv"
)

func blockchainHandler(rpc *daemon.RPC) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        SendOr404(w, rpc.GetBlockchainMetadata())
    }
}

func blockchainBlockHandler(rpc *daemon.RPC) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        sseq := r.FormValue("seq")
        seq, err := strconv.ParseUint(sseq, 10, 64)
        if err != nil {
            Error400(w, fmt.Sprintf("Invalid seq value \"%s\"", sseq))
            return
        }
        SendOr404(w, rpc.GetBlock(seq))
    }
}

func blockchainBlocksHandler(rpc *daemon.RPC) http.HandlerFunc {
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
        SendOr404(w, rpc.GetBlocks(start, end))
    }
}

func blockchainProgressHandler(rpc *daemon.RPC) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        SendOr404(w, rpc.GetBlockchainProgress())
    }
}

func RegisterBlockchainHandlers(mux *http.ServeMux, rpc *daemon.RPC) {
    mux.HandleFunc("/blockchain", blockchainHandler(rpc))
    mux.HandleFunc("/blockchain/block", blockchainBlockHandler(rpc))
    mux.HandleFunc("/blockchain/blocks", blockchainBlocksHandler(rpc))
    mux.HandleFunc("/blockchain/progress", blockchainProgressHandler(rpc))
}
