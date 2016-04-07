package main

import (
    "encoding/json"
    "os"
    "log"
    "io/ioutil"
    "flag"
    "time"
    "sync"
)

import (
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/skycoin/skycoin/src/daemon/gnet"
    "github.com/skycoin/skycoin/src/mesh"
)

var l_err = log.New(os.Stderr, "", 0)

var config_path = flag.String("config", "./config.json", "Configuration file path.")

var tcp_pool *gnet.ConnectionPool

var map_lock = &sync.Mutex{}
var configs_by_conn = make(map[*gnet.Connection]TCPOutgoingConnectionConfig)
type ConnectionSet map[*gnet.Connection]bool 
var conns_by_pubkey = make(map[cipher.PubKey]ConnectionSet)

var node_impl *mesh.Node

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
        map_lock.Lock()
        configs_by_conn[conn] = config
        conns_by_pubkey[config.PeerPubKey][conn] = true
        map_lock.Unlock()
        break
    }
}

// Does not block, unless the message queue is full
func doSendMessage(to_send mesh.OutgoingMessage) {
    defer func() {
        // recover from panic if one occured. Set err to nil otherwise.
        err := recover()
        if err != nil {
            l_err.Printf("doSendMessage panic %v", err)
        }
    }()

    map_lock.Lock()
    conns := conns_by_pubkey[to_send.ConnectedPeerPubKey]
    var conn *gnet.Connection = nil
    // For now we just choose the first one, and don't monitor health or anything
    // TODO: Choose a connection intelligently?
    for it_conn, _ := range conns {
        conn = it_conn
        break
    }
    map_lock.Unlock()
    if conn == nil {
        l_err.Printf("Warning: Send requested with no connections, dropping message to %v\n", to_send.ConnectedPeerPubKey)
    } else {
        // Does not block, unless the message queue is full
       tcp_pool.SendMessage(conn, WrapMessage(to_send.Message))
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

    config_node := config.Node
    config_node.ConnectedPeers = []cipher.PubKey{}
    for _ , conn_config := range config.TCPConnections {
        config_node.ConnectedPeers = append(config_node.ConnectedPeers, conn_config.PeerPubKey)
    }
    node_impl = mesh.NewNode(config_node)

    // Start listening for incoming connections via TCP
    config_cb := config.TCPConfig
    config_cb.DisconnectCallback = func(c *gnet.Connection, reason gnet.DisconnectReason) {
        map_lock.Lock()
        config, exists := configs_by_conn[c]
        if exists {
            delete(conns_by_pubkey[config.PeerPubKey], c)
            delete(configs_by_conn, c)
        }
        map_lock.Unlock()
        if exists {
            time.Sleep(config.RetryDelay)
            ConnectToPeerViaTCP(config)
        }
    }
    tcp_pool = gnet.NewConnectionPool(config_cb, node_impl)

    tcp_pool.StartListen()
    go tcp_pool.AcceptConnections()

    // Run connection pool
    go func() {
        for {
            disc := <- tcp_pool.DisconnectQueue
            tcp_pool.HandleDisconnectEvent(disc)
        }
    }()

    go func() {
        for {
            tcp_pool.HandleMessages();
        }
    }()

    // Connect to other nodes
    for _ , conn_config := range config.TCPConnections {
        map_lock.Lock()
        conns_by_pubkey[conn_config.PeerPubKey] = make(ConnectionSet)
        map_lock.Unlock()
        go ConnectToPeerViaTCP(conn_config)
    }

    // TODO: Pipe data

    // Send messages from queue
    go func() {
        for {
            to_send := <- node_impl.MessagesOut
            doSendMessage(to_send)
        }
    }()

    // Run Node Implementation (blocks)
    node_impl.Run();
}
