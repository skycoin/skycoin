package main

import (
	"fmt"
	"log"
	"time"

	"github.com/skycoin/skycoin/src/aether/daemon"
	gnet "github.com/skycoin/skycoin/src/aether/gnet"
)

//move into daemon

//Daemon on channel 0
//The channel 0 service manages exposing service metainformation and
//server setup and teardown

func NewDaemon(port int) *daemon.Daemon {

	config := daemon.NewConfig()
	//config.Daemon.LocalhostOnly = true
	config.Daemon.Port = port
	config.DHT.Disabled = true
	daemon := daemon.NewDaemon(config)
	return daemon
	//var swd SkywireDaemon
	//swd.ServiceManager = sm
	//associate service with channel 0
	//swd.Service = sm.AddService([]byte("Skywire Daemon"), 0, &swd)

}

//Daemon on Channel 1

//define message we want to be able to handle
type TestMessage struct {
	Text []byte
}

func (self *TestMessage) Handle(context *gnet.MessageContext, state interface{}) error {
	server := state.(*TestServiceServer) //service server state

	fmt.Printf("TestMessage Handle: ServiceIdLong= %s, Text= %s \n", string(server.Service.IdLong), string(self.Text))
	return nil
}

type TestServiceServer struct {
	Service *gnet.Service // Service
}

func NewTestServiceServer() *TestServiceServer {
	var tss TestServiceServer
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

	var quit1 chan int
	var quit2 chan int

	d1 := NewDaemon(6060) //server
	tss1 := NewTestServiceServer()
	tss1.Service = d1.ServiceManager.AddService(
		[]byte("test service"),
		[]byte("{service=\"test service\"}"), 1, tss1)

	//start daemon mainloop
	go d1.Start(quit1)

	//add services
	d2 := NewDaemon(6061)
	tss2 := NewTestServiceServer()
	tss2.Service = d2.ServiceManager.AddService(
		[]byte("test service"),
		[]byte("{service=\"test service\"}"), 1, tss2)

	go d2.Start(quit2) //start daemon main loop

	//sm2.AddService([]byte("Skywire Daemon"), 0, swd2)
	//TODO: do need servive level connect function?
	//connect to peer

	//use daemon method?
	time.Sleep(time.Second * 1)

	con, err := d1.Pool.Connect("127.0.0.1:6061")
	_ = con

	if err != nil {
		log.Panic(err)
	}

	//connect to service
	d1.ConnectToService(con, tss1.Service)
	time.Sleep(time.Second * 1)

	tm1 := TestMessage{Text: []byte("Message test 1")}
	tss1.Service.Send(con, &tm1)

	time.Sleep(time.Second * 1)
	tm2 := TestMessage{Text: []byte("Message test 2")}

	time.Sleep(time.Second * 1)
	tss2.Service.Broadcast(&tm2)

	time.Sleep(time.Second * 10)

	quit1 <- 1
	quit2 <- 2

}
