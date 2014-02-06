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
        sfee := r.FormValue("fee")
        fee, err := strconv.ParseUint(sfee, 10, 64)
        if err != nil {
            Error400(w, "Invalid \"fee\" value")
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
        m := rpc.Spend(visor.NewBalance(coins, hours), fee, dst)
        if SendJSON(w, m) != nil {
            Error500(w)
        }
    }
}

func walletSaveHandler(rpc *daemon.RPC) HTTPHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        err := rpc.SaveWallet()
        if err != nil {
            Error500(w, err.(error).Error())
        }
    }
}

func walletCreateAddressHandler(rpc *daemon.RPC) HTTPHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        addr := rpc.CreateAddress()
        if addr == nil {
            Error404(w)
        } else if SendJSON(w, addr) != nil {
            Error500(w)
        }
    }
}

func walletCreateHandler(rpc *daemon.RPC) HTTPHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        // TODO -- not clear to how to handle multiple wallets yet
    }
}

func walletHandler(rpc *daemon.RPC) HTTPHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        wallet := rpc.GetWallet()
        if wallet == nil {
            Error404(w)
        } else if SendJSON(w, wallet) != nil {
            Error500(w)
        }
    }
}

func RegisterWalletHandlers(mux *http.ServeMux, rpc *daemon.RPC) {
    mux.HandleFunc("/wallet", walletHandler(rpc))
    mux.HandleFunc("/wallet/balance", walletBalanceHandler(rpc))
    mux.HandleFunc("/wallet/spend", walletSpendHandler(rpc))
    mux.HandleFunc("/wallet/save", walletSaveHandler(rpc))
    mux.HandleFunc("/wallet/address/create", walletCreateAddressHandler(rpc))
    // Multiple wallets not supported
    // mux.HandleFunc("/wallet/create", walletCreateHandler(rpc))
    // History requires blockchain scans that will be very slow until
    // we have a more efficient datastructure
    // mux.HandleFunc("/wallet/history", walletHistoryHandler(rpc))
}
