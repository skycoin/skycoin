// Wallet-related information for the GUI
package gui

import (
    "net/http"
)

type WalletAddressEntry struct {
    Id      int
    Address string
    Balance string
}

type WalletPage struct {
    Title     string
    Addresses []WalletAddressEntry
}

// TODO -- strictly json api
func walletPageHandler(w http.ResponseWriter, req *http.Request) {
    var p WalletPage
    //fmt.Printf("S= %v \n", len(p.Folders) );

    /*
       for i, FM := range FL.FileList {
           var fpe FilePageEntry;
           fpe.Name = FM.Name;
           fpe.Size = FileSizeString(FM.Size)
           fpe.Hash = base64.URLEncoding.EncodeToString(FM.Hash[:])
           fpe.FileId = i
           p.Folders = append(p.Folders, fpe)
       }
    */
    //title := r.URL.Path[1:]
    ShowTemplate(w, "wallet.html", p)
}

func RegisterWalletHandlers(mux *http.ServeMux) {
    mux.HandleFunc("/wallet", walletPageHandler)
}
