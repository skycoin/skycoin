package blockchain



type ServerConfig struct {
    //Config Blockchain.ServerConfig
    // Disabled the Blockchain completely
    //Disabled bool

    // How often to request blocks from peers
    //BlocksRequestRate time.Duration
    
    // How often to announce our blocks to peers
    //BlocksAnnounceRate time.Duration
    
    // How many blocks to respond with to a GetBlocksMessage
    //BlocksResponseCount uint64
    // How often to rebroadcast txns that we are a party to
    //TransactionRebroadcastRate time.Duration
}

func NewServerConfig() ServerConfig {
    return ServerConfig{
        //Config:                     Blockchain.NewServerConfig(),
        //Disabled:                   false,
        //MasterKeysFile:             "",
        //BlocksRequestRate:          time.Minute * 5,
        //BlocksAnnounceRate:         time.Minute * 15,
        //BlocksResponseCount:        20,
        //TransactionRebroadcastRate: time.Minute * 5,
    }
}

type Server struct {
    Config ServerConfig
    Blockchain  *Blockchain
}

func NewServer(c ServerConfig) *Blockchain {
    var v *Blockchain.Blockchain = nil
    if !c.Disabled {
        v = Blockchain.NewBlockchain(c.Config)
    }
    return &Blockchain{
        Config:            c,
        Blockchain:             v,
        blockchainLengths: make(map[string]uint64),
    }
}

//start running
func (self *Server) Start() {

	t := time.Now.Unix()


	for true {

		if t + 15 < time.Now.Unix() {
			time.Sleep(50)
		}


	}
}

// Closes the block chain server, saving blockchain to disk
func (self *Server) Shutdown() {

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
}
