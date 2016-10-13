package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	gnet "github.com/skycoin/skycoin/src/aether"
	"github.com/skycoin/skycoin/src/cipher"
)

//this is called when client connects
func onConnect(c *gnet.Connection, solicited bool) {
	fmt.Printf("Event Callback: connnect event \n")
}

//this is called when client disconnects
func onDisconnect(c *gnet.Connection,
	reason gnet.DisconnectReason) {
	fmt.Printf("Event Callback: disconnect event \n")
}

//this is called when a message is received
func onMessage(c *gnet.Connection, channel uint16,
	msg []byte) error {
	fmt.Printf("Event Callback: message event: addr= %s, channel %v, msg= %s \n", c.Addr(), channel, msg)
	return nil
}

func SpawnConnectionPool(Port int) *gnet.ConnectionPool {
	config := gnet.NewConfig()
	config.Port = uint16(Port)               //set listening port
	config.DisconnectCallback = onDisconnect //disconnect callback
	config.ConnectCallback = onConnect       //connect callback
	config.MessageCallback = onMessage       //message callback

	//create pool
	cpool := gnet.NewConnectionPool(config)
	//open lsitening port
	if err := cpool.StartListen(); err != nil {
		log.Panic(err)
	}

	//listen for connections in new goroutine
	go cpool.AcceptConnections()

	//handle income data in new goroutine
	go func() {
		for true {
			time.Sleep(time.Millisecond * 100)
			cpool.HandleMessages()
		}
	}()

	return cpool
}

//MESSAGES

//define message we want to be able to handle
type TestMessage struct {
	Text []byte
}

func (self *TestMessage) Handle(context *gnet.MessageContext, state interface{}) error {
	pbcn := state.(*PublicBroadcastChannelNode) //the node receiving the message
	_ = pbcn

	fmt.Printf("TestMessage Handle: Text= %s \n", string(self.Text))
	return nil
}

//announce hash to peer
type AnnounceHashMessage struct {
	Hash   cipher.SHA256
	Pubkey cipher.PubKey
	Seq    uint64
}

func (self *AnnounceHashMessage) Handle(context *gnet.MessageContext, state interface{}) error {
	pbcn := state.(*PublicBroadcastChannelNode) //the node receiving the message
	_ = pbcn

	fmt.Printf("AnnounceHashMessage: Hash= %s \n", self.Hash[:])
	return nil
}

//request data
type GetHashMessage struct {
	Hash cipher.SHA256
	//Pubkey cipher.PubKey
	//Seq uint64
}

func (self *GetHashMessage) Handle(context *gnet.MessageContext, state interface{}) error {
	pbcn := state.(*PublicBroadcastChannelNode) //the node receiving the message
	_ = pbcn

	fmt.Printf("GetHashMessageMessage: Hash= %s \n", self.Hash[:])
	return nil
}

//request data
type GiveHashMessage struct {
	Hash   cipher.SHA256
	Pubkey cipher.PubKey
	//Seq uint64
	Data []byte
}

func (self *GiveHashMessage) Handle(context *gnet.MessageContext, state interface{}) error {
	pbcn := state.(*PublicBroadcastChannelNode) //the node receiving the message
	_ = pbcn

	fmt.Printf("GiveHashMessage: Hash= %s \n", self.Hash[:])
	return nil
}

//array of messages to register
var messageMap map[string](interface{}) = map[string](interface{}){
	"id01": TestMessage{}, //message id, message type
	"id02": AnnounceHashMessage{},
	"id03": GetHashMessage{},
	"id04": GiveHashMessage{},
}

//create connection pool and tests

type PublicBroadcastChannelNode struct {
	ConnectionPool *gnet.ConnectionPool
	Dispatcher     *gnet.Dispatcher
}

//new object
func NewPublicBroadcastChannelNode() *PublicBroadcastChannelNode {
	node := PublicBroadcastChannelNode{}
	return &node
}

func (self *PublicBroadcastChannelNode) InitConnectionPool(port int) {
	cpool := SpawnConnectionPool(port)

	//new dispatch manager
	dm := gnet.NewDispatcherManager()
	cpool.Config.MessageCallback = dm.OnMessage //set message handler
	//new dispatcher for handling messages on channel 1
	//d := dm.NewDispatcher(cpool1, 1, nil) //dispatcher, channel 1
	d := dm.NewDispatcher(cpool, 1, self) //dispatcher, channel 1, self as receiving object
	d.RegisterMessages(messageMap)

	self.ConnectionPool = cpool
	self.Dispatcher = d
}

//connect to ip:port
func (self *PublicBroadcastChannelNode) AddConnection(addr string) (*gnet.Connection, error) {
	con, err := self.ConnectionPool.Connect("127.0.0.1:6061")
	_ = con

	if err != nil {
		log.Panic(err)
		return nil, err
	}

	return con, nil
}

//broadcast to every connected peer
func (self *PublicBroadcastChannelNode) BroadcastMessage(message gnet.Message) {
	//create a message to send
	//tm := TestMessage{Text: []byte("Message test")}
	//Node1.ConnectionPool.SendMessage(con, 1, &tm)
	//Node1.Dispatcher.SendMessage(con, 1, &tm)
	for _, con := range self.ConnectionPool.Pool {
		self.Dispatcher.SendMessage(con, 1, message)
	}
}

//spawn network or array, the randomly connect them to each other
func SpawnNetwork(n int) []*PublicBroadcastChannelNode {

	var StartPort = 6060
	var NodeList []*PublicBroadcastChannelNode = make([]*PublicBroadcastChannelNode, n)

	for i := 0; i < n; i++ {
		NodeList[i] = NewPublicBroadcastChannelNode()
		NodeList[i].InitConnectionPool(StartPort + n)
	}

	//connect random nodes toget
	var ConnectionsPerNode = 3
	for i := 0; i < n; i++ {
		for j := 0; j < ConnectionsPerNode; j++ {
			var x int = rand.Int() % n

			for x == i {
				x = rand.Int() % n //avoid connection to self?
			}

			addr := fmt.Sprintf("128.0.0.1:%d", StartPort+x)
			fmt.Printf("addr= %s", addr)
			NodeList[i].AddConnection(addr)

			con, err := NodeList[i].AddConnection("127.0.0.1:6061")
			_ = con

			if err != nil {
				log.Panic(err)
			}

		}
	}

	return NodeList
}

func main() {
	Node1 := NewPublicBroadcastChannelNode()
	Node1.InitConnectionPool(6060)

	Node2 := NewPublicBroadcastChannelNode()
	Node2.InitConnectionPool(6061)

	//connect to peer
	con, err := Node1.AddConnection("127.0.0.1:6061")
	_ = con

	if err != nil {
		log.Panic(err)
	}

	//create a message to send
	tm := TestMessage{Text: []byte("Message test")}
	//Node1.ConnectionPool.SendMessage(con, 1, &tm)
	Node1.Dispatcher.SendMessage(con, 1, &tm)

	//d1.BroadcastMessage(3, &tm)

	time.Sleep(time.Second * 10)
}
