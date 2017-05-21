package main

import (
	"fmt"
	"log"
	"time"

	gnet "github.com/skycoin/skycoin/src/aether/gnet"
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

type ServiceConnectMessage struct {
	LocalChannel  uint16 //channel of service on sender
	RemoteChannel uint16 //channel of service on receiver
	Originating   uint32 //peer originating requests sets to 1
	ErrorMessage  []byte //fail if error len != 0
}

func (self *ServiceConnectMessage) Handle(context *gnet.MessageContext,
	state interface{}) error {
	server := state.(*SkywireDaemon) //service server state

	//message from remote for connection
	if self.Originating == 1 {
		service, ok := server.ServiceManager.Services[self.RemoteChannel]
		if ok == false {
			//server does not exist
			log.Printf("local service does not exist on channel %d \n", self.RemoteChannel)

			//failure message
			var scm ServiceConnectMessage
			scm.LocalChannel = self.RemoteChannel
			scm.RemoteChannel = self.LocalChannel
			scm.Originating = 0
			scm.ErrorMessage = []byte("no service on channel")
			server.Service.Send(context.Conn, &scm) //channel 0
			return nil
		} else {
			//service exists, send success message
			var scm ServiceConnectMessage
			scm.LocalChannel = self.RemoteChannel
			scm.RemoteChannel = self.LocalChannel
			scm.Originating = 0
			scm.ErrorMessage = []byte("")
			server.Service.Send(context.Conn, &scm) //channel 0
			//trigger connection event
			service.ConnectionEvent(context.Conn, self.LocalChannel)
			return nil
		}
	}
	//message response from remote for connection
	if self.Originating == 0 {
		if len(self.ErrorMessage) != 0 {
			log.Printf("Service Connection Failed: addr= %s, LocalChannel= %d, Remotechannel= %d \n",
				context.Conn.Addr(), self.LocalChannel, self.RemoteChannel)
			return nil
		}

		service, ok := server.ServiceManager.Services[self.RemoteChannel]

		if ok == false {
			log.Printf("service does not exist on local, LocalChannel= %d from addr= %s \n",
				self.RemoteChannel, context.Conn.Addr())
		}

		service.ConnectionEvent(context.Conn, self.LocalChannel)
		return nil
	}
	return nil
}

//Daemon on channel 0
//The channel 0 service manages exposing service metainformation and
//server setup and teardown
type SkywireDaemon struct {
	Service        *gnet.Service //service for daemon
	ServiceManager *gnet.ServiceManager
}

// TODO:
// - add request packet for service list
// - add connection packet for service
// - move into daemon

func NewSkywireDaemon(sm *gnet.ServiceManager) *SkywireDaemon {
	var swd SkywireDaemon
	swd.ServiceManager = sm

	//associate service with channel 0
	swd.Service = sm.AddService([]byte("Skywire Daemon"), []byte(""), 0, &swd)

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
		"SCON": ServiceConnectMessage{}, //connect to service
	}
	d.RegisterMessages(messageMap)

}

//Daemon on Channel 1

//define message we want to be able to handle
type TestMessage struct {
	Text []byte
}

func (self *TestMessage) Handle(context *gnet.MessageContext, state interface{}) error {
	server := state.(*TestServiceServer) //service server state

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

	service_id := []byte("id short")

	cpool1 := SpawnConnectionPool(6060)   //connection pool
	sm1 := gnet.NewServiceManager(cpool1) //service manager
	//add services
	swd1 := NewSkywireDaemon(sm1) //server

	tss1 := NewTestServiceServer()
	sm1.AddService(
		service_id,
		[]byte("id long"), 1, tss1)

	cpool2 := SpawnConnectionPool(6061)
	sm2 := gnet.NewServiceManager(cpool2)
	//add services
	swd2 := NewSkywireDaemon(sm2)
	//sm2.AddService([]byte("Skywire Daemon"), 0, swd2)

	tss2 := NewTestServiceServer()
	sm2.AddService(
		service_id,
		[]byte("id long"), 1, tss2)

	//TODO: do need servive level connect function?
	//connect to peer
	con, err := cpool1.Connect("127.0.0.1:6061")
	_ = con

	if err != nil {
		log.Panic(err)
	}

	//connection attempt message
	scm := ServiceConnectMessage{}
	scm.LocalChannel = 1  //channel of local service
	scm.RemoteChannel = 1 //channel of remote service
	scm.Originating = 1
	scm.ErrorMessage = []byte("")
	//send connection intiation
	swd1.Service.Send(con, &scm)

	//create a message to send
	//tm := TestMessage{Text: []byte("Message test")}
	//d1.SendMessage(con, 3, &tm)

	_ = swd1
	_ = swd2

	time.Sleep(time.Second * 10)
}
