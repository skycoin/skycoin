package sb_gui

import (
	"fmt"
	"net/http"
)

import "common"

var GuiServerPort uint32 = 7999

type GUIServer struct {
	Port uint32
}

func (self *GUIServer) Run(mux *http.ServeMux, errchan chan<- error) {
	if mux == nil {
		mux = http.NewServeMux()
	}
	mux.Handle("/static/css/", http.StripPrefix("/static/css/",
		http.FileServer(http.Dir("./static/css/"))))

	mux.HandleFunc("/wallet", self.WalletPageHandler)

	self.Port = GuiServerPort
	address := fmt.Sprintf("localhost:%d", self.Port)
	fmt.Printf("GUI server: running on http://%s\n", address)

	go sb.ListenAndServeBackground(address, mux, errchan)
}
