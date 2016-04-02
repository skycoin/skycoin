package main

import (
    "encoding/json"
    "os"
    "log"
    "io/ioutil"
    "flag"
    "time"
)

import (
    "github.com/skycoin/skycoin/src/daemon/gnet"
)

var l_err = log.New(os.Stderr, "", 0)

var config_path = flag.String("config", "./config.json", "Configuration file path.")

var tcp_pool *gnet.ConnectionPool
var configs_by_conn = make(map[*gnet.Connection]TCPOutgoingConnectionConfig)

func ConnectToPeerViaTCP(config TCPOutgoingConnectionConfig) {
    for {
        conn, err := tcp_pool.Connect(config.Endpoint)
        if err != nil {
            l_err.Printf("Error connecting to %v(%s): %v", config.PeerPubKey, config.Endpoint, err)
            l_err.Printf("Retrying in %v second(s)...", (int)(config.RetryDelay / 1000000000))
            time.Sleep(config.RetryDelay)
            continue
        }
        l_err.Printf("Connected to %v(%s)", config.PeerPubKey, config.Endpoint)
        configs_by_conn[conn] = config
        break
    }
}

func main() {
    RegisterTCPMessages()
    flag.Parse()

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
    config_cb := config.TCPConfig
    config_cb.DisconnectCallback = func(c *gnet.Connection, reason gnet.DisconnectReason) {
        config, exists := configs_by_conn[c]
        if exists {
            time.Sleep(config.RetryDelay)
            ConnectToPeerViaTCP(config)
        }
    }
    tcp_pool = gnet.NewConnectionPool(config_cb, nil)

    tcp_pool.StartListen()
    go tcp_pool.AcceptConnections()

    go func() {
        for {
            tcp_pool.HandleMessages();
        } 
    }()

    // Connect to other nodes
    for _ , conn_config := range config.TCPConnections {
        go ConnectToPeerViaTCP(conn_config)
    }

    // TODO: Establish route
{
    /*
        // TODO: Temp
        var test_m *SendMessageWrapper = &SendMessageWrapper{}
        test_m.SendMessage.SendId = uint32(time.Now().Second() % 1000)
        test_m.SendMessage.Message = []byte{4, 5, 6}

        tcp_pool.SendMessage(conn, test_m)
        select{ }
    */
}

    // TODO: Pipe data


    // Wait forever
    select{}
}
