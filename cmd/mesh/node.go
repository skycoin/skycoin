package main

import (
    "encoding/json"
    "os"
    "log"
    "fmt"
    "io/ioutil"
    "flag"
    "time"
    "sync"
    "reflect"
)

import (
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/skycoin/skycoin/src/daemon/gnet"
    "github.com/skycoin/skycoin/src/mesh"
    "github.com/satori/go.uuid"
)

var l_err = log.New(os.Stderr, "", 0)

var host_service = flag.String("host", "", "Host builtin service, choices: 'proxy', 'tun'")
var config_path = flag.String("config", "./config.json", "Configuration file path.")

var tcp_pool *gnet.ConnectionPool

const STDIO_CHANLEN = 100

var stdoutQueue = make(chan interface{}, STDIO_CHANLEN)

var map_lock = &sync.Mutex{}
var configs_by_conn = make(map[*gnet.Connection]TCPOutgoingConnectionConfig)
var pub_keys_by_conn = make(map[*gnet.Connection]cipher.PubKey)
type ConnectionSet map[*gnet.Connection]bool 
var conns_by_pubkey = make(map[cipher.PubKey]ConnectionSet)

var node_impl *mesh.Node

type ConnectAnnouncementMessage struct {
    MyPubKey cipher.PubKey
}
var ConnectAnnouncementMessagePrefix = gnet.MessagePrefix{0,0,0,1}
func (self *ConnectAnnouncementMessage) Handle(context *gnet.MessageContext, x interface{}) error {
    map_lock.Lock()
    pub_keys_by_conn[context.Conn] = self.MyPubKey
    if conns_by_pubkey[self.MyPubKey] == nil {
        conns_by_pubkey[self.MyPubKey] = make(ConnectionSet)
    }
    conns_by_pubkey[self.MyPubKey][context.Conn] = true
    map_lock.Unlock()
    return nil
}

type NodeMessage struct {
    Contents []byte
}
var NodeMessagePrefix = gnet.MessagePrefix{0,0,0,2}
func (self *NodeMessage) Handle(context *gnet.MessageContext, x interface{}) error {
    map_lock.Lock()
    pub_key, exists := pub_keys_by_conn[context.Conn]
    if exists {
        node_impl.MessagesIn <- mesh.PhysicalMessage{pub_key, self.Contents}
    } else {
        l_err.Printf("Dropping NodeMessage from unknown connection")
    }
    map_lock.Unlock()
    return nil
}

func connectToPeerViaTCP(config TCPOutgoingConnectionConfig) {
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
        // Does not block, unless the message queue is full
        tcp_pool.SendMessage(conn, &ConnectAnnouncementMessage{node_impl.Config.MyPubKey})
        break
    }
}

// Does not block, unless the message queue is full
func doSendMessage(to_send mesh.PhysicalMessage) {
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
       tcp_pool.SendMessage(conn, &NodeMessage{to_send.Contents})
    }
}

func getConnectedPeerKey(Conn *gnet.Connection) cipher.PubKey {
    map_lock.Lock()
    key, exists := pub_keys_by_conn[Conn]
    if !exists {
        panic("Internal consistency failure: getConnectedPeerKey() called for unknown connection")
    }
    map_lock.Unlock()
    return key
}

func onStdInMessage(msg interface{}) {
    if reflect.TypeOf(msg) == reflect.TypeOf(Stdin_SendMessage{}) {
        msg_cast := msg.(Stdin_SendMessage)
        node_impl.SendMessage(msg_cast.RouteId, msg_cast.Contents)
    } else if reflect.TypeOf(msg) == reflect.TypeOf(Stdin_SendBack{}) {
        msg_cast := msg.(Stdin_SendBack)
        node_impl.SendReply(msg_cast.ReplyTo, msg_cast.Contents)
    } else {
        panic("Unknown message type in onStdInMessage")
    }
}

func sendRoutes() {
    route_names := make([]string, len(node_impl.Config.Routes))
    route_ids := make([]uuid.UUID, len(node_impl.Config.Routes))    
    for i, route_config := range node_impl.Config.Routes {
        route_names[i] = route_config.Name
        route_ids[i] = route_config.Id
    }
    stdoutQueue <- Stdout_RoutesChanged{route_names, route_ids}
}

func main() {
    flag.Parse()

    switch *host_service {
        case "":
            break
        case "proxy":
            HostProxy()
            return
        case "tun":
            HostTun()
            return
        default:
            panic(fmt.Sprintf("Unknown service to host: %v", *host_service))
    }

    gnet.RegisterMessage(ConnectAnnouncementMessagePrefix, ConnectAnnouncementMessage{})
    gnet.RegisterMessage(NodeMessagePrefix, NodeMessage{})

    stdio_serializer := NewStdioSerializer()

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
    config.Node.RouteEstablishmentCB = func(RouteId uuid.UUID, HopIdx int) {
        stdoutQueue <- Stdout_RouteEstablishment{RouteId, (uint32)(HopIdx)}
    }
    config.Node.RouteEstablishedCB = func(route mesh.EstablishedRoute) {
        stdoutQueue <- Stdout_EstablishedRoute{route.RouteId}
    }
    node_impl = mesh.NewNode(config.Node)

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
            connectToPeerViaTCP(config)
        }
    }
    config_cb.ConnectCallback = func (conn *gnet.Connection, solicited bool) {
        tcp_pool.SendMessage(conn, &ConnectAnnouncementMessage{node_impl.Config.MyPubKey})
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
        go connectToPeerViaTCP(conn_config)
    }

    // Send messages from queue
    go func() {
        for {
            to_send := <- node_impl.MessagesOut
            doSendMessage(to_send)
        }
    }()

    // Output messages received
    go func() {
        for {
            stdoutQueue <- Stdout_RecvMessage{<- node_impl.MeshMessagesIn}
        }
    }()

    // Send static routes
    sendRoutes();

    // Send config url
    stdoutQueue <- Stdout_StaticConfig{"about:test"}

    var stdinQueue = make(chan interface{}, STDIO_CHANLEN)

    // We are the subprocess
    go RunStdioSerializer(stdio_serializer, stdinQueue, os.Stdin, stdoutQueue, os.Stdout)
    
    go func() {
        for {
            onStdInMessage(<- stdinQueue)
        }
    }()

    // Run Node Implementation (blocks)
    node_impl.Run();
}
