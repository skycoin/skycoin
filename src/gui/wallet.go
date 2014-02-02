// Wallet-related information for the GUI
package gui

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/daemon"
    "github.com/skycoin/skycoin/src/visor"
    "net/http"
    "strconv"
)

func walletBalanceHandler(rpc *daemon.RPC) HTTPHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        saddr := r.FormValue("addr")
        var m interface{}
        if saddr == "" {
            m = rpc.GetTotalBalance()
        } else {
            addr, err := coin.DecodeBase58Address(saddr)
            if err != nil {
                Error400(w, "Invalid address")
                return
            }
            m = rpc.GetBalance(addr)
        }
        balance := m.(*visor.Balance)
        if balance == nil {
            Error404(w)
        } else if SendJSON(w, m) != nil {
            Error500(w)
        }
    }
}

func walletSpendHandler(rpc *daemon.RPC) HTTPHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        sdst := r.FormValue("dst")
        if sdst == "" {
            Error400(w, "Missing destination address \"dst\"")
            return
        }
        dst, err := coin.DecodeBase58Address(sdst)
        if err != nil {
            Error400(w, "Invalid destination address")
            return
        }

        scoins := r.FormValue("coins")
        shours := r.FormValue("hours")
        coins, err := strconv.ParseUint(scoins, 10, 64)
        if err != nil {
            Error400(w, "Invalid \"coins\" value")
            return
        }
        hours, err := strconv.ParseUint(shours, 10, 64)
        if err != nil {
            Error400(w, "Invalid \"hours\" value")
            return
        }
        m := rpc.Spend(visor.NewBalance(coins, hours), dst)
        if SendJSON(w, m) != nil {
            Error500(w)
        }
    }
}

func RegisterWalletHandlers(mux *http.ServeMux, rpc *daemon.RPC) {
    mux.HandleFunc("/wallet/balance", walletBalanceHandler(rpc))
    mux.HandleFunc("/wallet/spend", walletSpendHandler(rpc))
}
