package blockchain

import (
    "fmt"
    "log"
    "time"
    //"net/http"
)

type ServerConfig struct {
}

func NewServerConfig() ServerConfig {
    return ServerConfig{}
}

type Server struct {
    Config     ServerConfig
    Blockchain Blockchain
}

func NewServer(c ServerConfig) *Server {
    return &Server{
        Config:     NewServerConfig(),
        Blockchain: NewLocalBlockchain(),
    }
}

/*
func handler(w http.ResponseWriter, r *http.Request) {
}

func StartTransactionServer() {
    fmt.Printf("starting transaction server on port 666")
    http.HandleFunc("/injectTransaction", handler)
    http.ListenAndServe(":666", nil)

}
*/

//Start server as master
func (self *Server) Start() {

    t := time.Now().Unix()

    for true {

        //wait 15 seconds between blocks
        if t+15 < time.Now().Unix() {
            time.Sleep(50)
            continue
        } else {
            t = time.Now().Unix() //update time
        }

        if self.Blockchain.PendingTransactions() == false {
            continue
        }

        //create block
        block, err := self.Blockchain.CreateBlock()
        if err != nil {
            fmt.Printf("Create Block Error: %s \n", err)
            continue
        }
        //sign block
        signedBlock := self.Blockchain.SignBlock(block)

        //inject block/execute
        err = self.Blockchain.InjectBlock(signedBlock)
        if err != nil {
            log.Panic(err)
        }
        //prune unconfirmed transactions
        self.Blockchain.RefreshUnconfirmed()
    }
}

func (self *Server) StartSlave() {
    for true {
        time.Sleep(500)
        self.Blockchain.RefreshUnconfirmed()
    }
}

// Closes the block chain server, saving blockchain to disk
func (self *Server) Shutdown() {

    /*
       bcFile := self.Config.Config.BlockchainFile
       err := self.Blockchain.SaveBlockchain()
       if err == nil {
           logger.Info("Saved blockchain to \"%s\"", bcFile)
       } else {
           logger.Critical("Failed to save blockchain to \"%s\"", bcFile)
       }
       bsFile := self.Config.Config.BlockSigsFile
       err = self.Blockchain.SaveBlockSigs()
       if err == nil {
           logger.Info("Saved block sigs to \"%s\"", bsFile)
       } else {
           logger.Critical("Failed to save block sigs to \"%s\"", bsFile)
       }
    */
}
