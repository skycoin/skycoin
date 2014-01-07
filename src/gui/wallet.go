package sb_gui

import (
	"net/http"
)

import "common"

/*
   Wallet Page
*/

//todo, add support for tonal number system

type WalletAddressEntry struct {
	Id      int
	Address string
	Balance string
}

type WalletPage struct {
	Title     string
	Addresses []WalletAddressEntry
}

func (g *GUIServer) WalletPageHandler(w http.ResponseWriter, req *http.Request) {
	var p WalletPage
	sb.ShowTemplate(w, "wallet.html", p)
}
