package main


import (
  "net/http"
  "io/ioutil"
  wh "github.com/skycoin/skycoin/src/util/http"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {

  resp, err := http.Get("http://127.0.0.1:6420/outputs")
  if err != nil {
    wh.Error500(w,"Unable to respond back")
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  wh.SendOr500(w,body)
}

func getBlocks(w http.ResponseWriter, r *http.Request) {
  startBlock := r.URL.Query().Get("start")
  endBlock := r.URL.Query().Get("end")
  resp, err := http.Get("http://127.0.0.1:6420/blocks?start="+startBlock+"&end="+endBlock)
  if err != nil {
    wh.Error500(w,"Unable to respond back")
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  w.Write(body)
}

func getSupply(w http.ResponseWriter, r *http.Request) {
  resp, err := http.Get("http://127.0.0.1:6420/explorer/getEffectiveOutputs")
  if err != nil {
    wh.Error500(w,"Unable to respond back")
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  w.Write(body)
}


func getBlockChainMetaData(w http.ResponseWriter, r *http.Request) {
  resp, err := http.Get("http://127.0.0.1:6420/blockchain/metadata")
  if err != nil {
    wh.Error500(w,"Unable to respond back")
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  w.Write(body)
}

func getAddress(w http.ResponseWriter, r *http.Request) {
  address := r.URL.Query().Get("address")
  resp, err := http.Get("http://127.0.0.1:6420/explorer/address?address="+address)
  if err != nil {
    wh.Error500(w,"Unable to respond back")
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  w.Write(body)
}

func getCurrentBalance(w http.ResponseWriter, r *http.Request) {
  address := r.URL.Query().Get("address")
  resp, err := http.Get("http://127.0.0.1:6420/outputs?addrs="+address)
  if err != nil {
    wh.Error500(w,"Unable to respond back")
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  w.Write(body)
}


func getUxID(w http.ResponseWriter, r *http.Request) {
  uxid := r.URL.Query().Get("uxid")
  resp, err := http.Get("http://127.0.0.1:6420/uxout?uxid="+uxid)
  if err != nil {
    wh.Error500(w,"Unable to respond back")
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  w.Write(body)
}

func getTransaction(w http.ResponseWriter, r *http.Request) {
  txid := r.URL.Query().Get("txid")
  resp, err := http.Get("http://127.0.0.1:6420/transaction?txid="+txid)
  if err != nil {
    wh.Error500(w,"Unable to respond back")
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  w.Write(body)
}

func getBlock(w http.ResponseWriter, r *http.Request) {
  hash := r.URL.Query().Get("hash")
  resp, err := http.Get("http://127.0.0.1:6420/block?hash="+hash)
  if err != nil {
    wh.Error500(w,"Unable to respond back")
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  w.Write(body)
}
func redirectToRoot(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {

  //http.Handle("/", http.FileServer(http.Dir("./dist")))


  http.Handle("/fonts/*", http.FileServer(http.Dir("./dist/fonts")))

  http.HandleFunc("/fonts/roboto/Roboto-Light.woff2", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./dist/fonts/roboto/Roboto-Light.woff2")
  })
  http.HandleFunc("/fonts/roboto/Roboto-Medium.woff2", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./dist/fonts/roboto/Roboto-Medium.woff2")
  })
  http.HandleFunc("/fonts/roboto/Roboto-Light.woff", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./dist/fonts/roboto/Roboto-Light.woff")
  })
  http.HandleFunc("/fonts/roboto/Roboto-Medium.woff", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./dist/fonts/roboto/Roboto-Medium.woff")
  })

  http.HandleFunc("/assets/materialize.css", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./dist/assets/materialize.css")
  })

  http.HandleFunc("/assets/styles.css", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./dist/assets/styles.css")
  })

  http.HandleFunc("/inline.bundle.js", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./dist/inline.bundle.js")
  })

  http.HandleFunc("/scripts.bundle.js", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./dist/scripts.bundle.js")
  })

  http.HandleFunc("/main.bundle.js", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./dist/main.bundle.js")
  })

  http.HandleFunc("/vendor.bundle.js", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./dist/vendor.bundle.js")
  })

  http.HandleFunc("/skycoin.ico", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./dist/skycoin.ico")
  })

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./dist/index.html")
  })


  http.HandleFunc("/api/hello", helloWorld)
  http.HandleFunc("/api/blocks", getBlocks)
  http.HandleFunc("/api/blockchain/metadata", getBlockChainMetaData)
  http.HandleFunc("/api/address", getAddress)
  http.HandleFunc("/api/currentBalance",getCurrentBalance)
  http.HandleFunc("/api/uxout", getUxID)
  http.HandleFunc("/api/transaction", getTransaction)
  http.HandleFunc("/api/block", getBlock)
  http.HandleFunc("/api/coinSupply", getSupply)
  http.ListenAndServe(":8001", nil)
}
