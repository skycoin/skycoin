package main

import (
	"fmt"
	"github.com/skycoin/skywire/src/lib/gnet"
	"log"
	"time"
)

/*
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
*/

func SpawnConnectionPool(Port int) *gnet.ConnectionPool {
	config := gnet.NewConfig()
	config.Port = uint16(Port) //set listening port
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

//array of messages to register
//var messageMap map[string](interface{}) = map[string](interface{}){
//	"id01": TestMessage{}, //message id, message type
//}

//Daemon on channel 0
//The channel 0 service manages exposing service metainformation and
//server setup and teardown
type SkywireDaemon struct {
	ServiceManager *gnet.ServiceManager
}

func NewSkywireDaemon(sm *gnet.ServiceManager) *SkywireDaemon {
	var swd SkywireDaemon
	swd.ServiceManager = sm
	return &swd
}

func (sd *SkywireDaemon) OnConnect(c *gnet.Connection) {
	fmt.Printf("SkywireDaemon: OnConnect, addr= %s \n", c.Addr())
}

func (sd *SkywireDaemon) OnDisconnect(c *gnet.Connection) {
	fmt.Printf("SkywireDaemon: OnDisconnect, addr= %s \n", c.Addr())
}

func (sd *SkywireDaemon) RegisterMessages(d *gnet.Dispatcher) {
	fmt.Printf("SkywireDaemon: RegisterMessages \n")

	var messageMap map[string](interface{}) = map[string](interface{}){
	//put messages here
	//"id01": TestMessage{}, //message id, message type
	}
	d.RegisterMessages(messageMap)

}

//Daemon on Channel 1

//define message we want to be able to handle
type TestMessage struct {
	Text []byte
}

func (self *TestMessage) Handle(context *gnet.MessageContext, state interface{}) error {
	server := state.(TestServiceServer) //service server state

	fmt.Printf("TestMessage Handle: ServerName= %s, Text= %s \n", string(server.Name), string(self.Text))
	return nil
}

type TestServiceServer struct {
	Name []byte
}

func NewTestServiceServer() *TestServiceServer {
	var tss TestServiceServer
	tss.Name = []byte("Server1")
	return &tss
}

func (sd *TestServiceServer) OnConnect(c *gnet.Connection) {
	fmt.Printf("TestServiceServer: OnConnect, addr= %s \n", c.Addr())
}

func (sd *TestServiceServer) OnDisconnect(c *gnet.Connection) {
	fmt.Printf("TestServiceServer: OnDisconnect, addr= %s \n", c.Addr())
}

func (sd *TestServiceServer) RegisterMessages(d *gnet.Dispatcher) {
	fmt.Printf("TestServiceServer: RegisterMessages \n")

	var messageMap map[string](interface{}) = map[string](interface{}){
		//put messages here
		"id01": TestMessage{}, //message id, message type
	}
	d.RegisterMessages(messageMap)

}

/*
TODO:
	- expose server meta-information through channel 0
	- handle connection setup and teardown through channel 0


*/

//create connection pool and tests
func main() {

	cpool1 := SpawnConnectionPool(6060)   //connection pool
	sm1 := gnet.NewServiceManager(cpool1) //service manager
	//add services
	swd1 := NewSkywireDaemon(sm1) //server
	sm1.AddService([]byte("Skywire Daemon"), 0, swd1)
	tss1 := NewTestServiceServer()
	sm1.AddService([]byte("TestServiceServer"), 1, tss1)

	cpool2 := SpawnConnectionPool(6061)
	sm2 := gnet.NewServiceManager(cpool2)
	//add services
	swd2 := NewSkywireDaemon(sm2)
	sm2.AddService([]byte("Skywire Daemon"), 0, swd2)
	tss2 := NewTestServiceServer()
	sm2.AddService([]byte("TestServiceServer"), 1, tss2)

	//TODO: do need servive level connect function?
	//connect to peer
	con, err := cpool1.Connect("127.0.0.1:6061")
	_ = con

	if err != nil {
		log.Panic(err)
	}

	//create a message to send
	//tm := TestMessage{Text: []byte("Message test")}
	//d1.SendMessage(con, 3, &tm)

	time.Sleep(time.Second * 10)
}
