// Wallet-related information for the GUI
package gui

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/daemon"
    "github.com/skycoin/skycoin/src/visor"
    "net/http"
    "strconv"
)

func walletBalanceHandler(gateway *daemon.Gateway) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        saddr := r.FormValue("addr")
        predicted := r.FormValue("predicted")
        var m interface{}
        if saddr == "" {
            m = gateway.GetTotalBalance(predicted != "")
        } else {
            addr, err := coin.DecodeBase58Address(saddr)
            if err != nil {
                Error400(w, "Invalid address")
                return
            }
            m = gateway.GetBalance(addr, predicted != "")
        }
        SendOr404(w, m)
    }
}

func walletSpendHandler(gateway *daemon.Gateway) http.HandlerFunc {
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
        SendOr404(w, gateway.Spend(visor.NewBalance(coins, hours), fee, dst))
    }
}

func walletSaveHandler(gateway *daemon.Gateway) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := gateway.SaveWallet()
        if err != nil {
            Error500(w, err.(error).Error())
        }
    }
}

func walletCreateAddressHandler(gateway *daemon.Gateway) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        SendOr404(w, gateway.CreateAddress())
    }
}

func walletCreateHandler(gateway *daemon.Gateway) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // TODO -- not clear to how to handle multiple wallets yet
    }
}

func walletHandler(gateway *daemon.Gateway) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        SendOr404(w, gateway.GetWallet())
    }
}

func walletTransactionResendHandler(gateway *daemon.Gateway) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        hash, err := coin.SHA256FromHex(r.FormValue("hash"))
        if err != nil {
            Error404(w)
            return
        }
        SendOr404(w, gateway.ResendTransaction(hash))
    }
}

func walletAddressTransactionsHandler(gateway *daemon.Gateway) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        saddr := r.FormValue("addr")
        addr, err := coin.DecodeBase58Address(saddr)
        if err != nil {
            Error404(w)
            return
        }
        SendOr404(w, gateway.GetAddressTransactions(addr))
    }
}

func walletTransactionHandler(gateway *daemon.Gateway) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        hash, err := coin.SHA256FromHex(r.FormValue("hash"))
        if err != nil {
            Error404(w)
            return
        }
        SendOr404(w, gateway.GetTransaction(hash))
    }
}

func RegisterWalletHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
    mux.HandleFunc("/wallet", walletHandler(gateway))
    mux.HandleFunc("/wallet/balance", walletBalanceHandler(gateway))
    mux.HandleFunc("/wallet/spend", walletSpendHandler(gateway))
    mux.HandleFunc("/wallet/save", walletSaveHandler(gateway))
    mux.HandleFunc("/wallet/transaction", walletTransactionHandler(gateway))
    mux.HandleFunc("/wallet/address/create",
        walletCreateAddressHandler(gateway))
    mux.HandleFunc("/wallet/address/transactions",
        walletAddressTransactionsHandler(gateway))
    mux.HandleFunc("/wallet/transaction/resend",
        walletTransactionResendHandler(gateway))
    // Multiple wallets not supported
    // mux.HandleFunc("/wallet/create", walletCreateHandler(gateway))
}
