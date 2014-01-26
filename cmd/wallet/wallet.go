package main

import (
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"github.com/skycoin/skycoin/src/keyring"
	//"io"
	"net/http"
	//"os"
	"path/filepath"
	//"log"
)

type Profile struct {
	Name    string
	Hobbies []string
}

type walletData struct {
	Seed      string
	Addresses []string
	History   []string
}

var (
	logger = logging.MustGetLogger("skycoin.gui")
)

var walletFile = walletData{}

func main() {
	static_path, _ := filepath.Abs("../../static/app/")
	logger.Debug("Serving %s", static_path)

	http.Handle("/", http.FileServer(http.Dir(static_path)))

	http.HandleFunc("/api/newAdress", newAddress)

	http.ListenAndServe(":3003", nil)
}

func newAddress(w http.ResponseWriter, r *http.Request) {

	logger.Debug("Serving %s", r)

	//js, err := json.Marshal(profile)
	addr := keyring.GenerateAddress()

	//walletFile.Addresses = append(walletFile.Addresses, addr)
	fmt.Printf("address= %s \n", addr.Address.String())

	js, err := json.Marshal(addr.Address.String())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
