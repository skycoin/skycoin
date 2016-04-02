package main

import (
    "encoding/json"
    "os"
    "log"
    "io/ioutil"
    "flag"
)

import (
    "github.com/skycoin/skycoin/src/daemon/gnet"
)

var config_path = flag.String("config", "./config.json", "Configuration file path.")

func main() {
    RegisterTCPMessages()
    flag.Parse()

    l_err := log.New(os.Stderr, "", 0)

	file, e := ioutil.ReadFile(*config_path)
    if e != nil {
        l_err.Printf("Config file open error: %v\n", e)
        os.Exit(1)
    }

	var config Config
	e_parse := json.Unmarshal(file, &config)
    if e_parse != nil {
        l_err.Printf("Config parse error: %v\n", e_parse)
        os.Exit(1)
    }

    // Start listening for incoming connections via TCP
    var pool *gnet.ConnectionPool
    if config.TCPServer != nil {
        config_mod := *config.TCPServer
        pool = gnet.NewConnectionPool(config_mod, nil)

        pool.StartListen()
        go pool.AcceptConnections()

        go func() {
            for {
                pool.HandleMessages();
            } 
        }()
    }
    // Connect to other nodes
    for _ , conn_config := range config.Connections {
        l_err.Println("Connections %v", conn_config)

        // gnet.sendMessage
    }

    // Wait forever
    select{}
}
