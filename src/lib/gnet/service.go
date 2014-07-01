package gnet

import (
	"log"
)

//deal with to and from and connection handlings
type ServiceManager struct {
	ConnectionPool  *ConnectionPool
	DispatchManager *DispatcherManager

	Services map[uint16]*Service //channel to service map
}

func NewServiceManager(pool *ConnectionPool) *ServiceManager {

	var sm ServiceManager

	sm.Services = make(map[uint16]*Service)

	sm.ConnectionPool = pool
	sm.DispatchManager = NewDispatcherManager()

	//handle messages
	pool.Config.MessageCallback = sm.DispatchManager.OnMessage
	//connect/disconnect callbacks
	pool.Config.DisconnectCallback = sm.OnDisconnect
	pool.Config.ConnectCallback = sm.OnConnect
	return &sm
}

func (sm *ServiceManager) AddService(id []byte, idLong []byte, channel uint16, server ServiceServer) *Service {

	if _, ok := sm.Services[channel]; ok != false {
		log.Panic("duplicate service channels")
	}

	if len(idLong) > 140 {
		log.Panic("Service Identifier must not be longer than 140 characters")
	}

	if len(id) > 20 {
		log.Panic("ServiceManager: ID must be 20 bytes or less")
	}
	var s Service
	copy(s.Id[0:20], id[:])

	s.IdLong = idLong

	s.Channel = channel
	//need to pass in object
	s.Dispatcher = sm.DispatchManager.NewDispatcher(sm.ConnectionPool, channel, server)
	s.Server = server
	s.Connections = make(map[*Connection]uint16)

	server.RegisterMessages(s.Dispatcher) //register server messages

	sm.Services[channel] = &s
	return &s
}

//connection level connect
func (sm *ServiceManager) OnConnect(c *Connection, solicited bool) {

	//channel 0 gets all connection/disconnect events
	if _, ok := sm.Services[0]; ok != false {
		sm.Services[0].ConnectionEvent(c, 0)
	} else {
		log.Panic("channel 0 service not defined")
	}
}

//connection level disconnect
func (sm *ServiceManager) OnDisconnect(c *Connection,
	reason DisconnectReason) {

	for _, service := range sm.Services {
		if _, ok := service.Connections[c]; ok != false {
			service.DisconnectEvent(c)
		}
	}
}

//return service by ID or return null
func (sm *ServiceManager) ServiceById(Id [20]byte) *Service {
	for _, service := range sm.Services {
		if service.Id == Id {
			return service
		}
	}
	return nil
}

//func (sm *ServiceManager) OnMessage(c *Connection, channel uint16,
//	msg []byte) {
//}

type Service struct {
	//Name             []byte
	Id     [20]byte
	IdLong []byte

	Channel     uint16                 //channel for receiving
	Connections map[*Connection]uint16 //outgoing channel for connection
	Dispatcher  *Dispatcher

	Server ServiceServer //server implementing service
}

//send to single peer of service
func (self *Service) Send(c *Connection, msg Message) {
	channel, ok := self.Connections[c]
	if ok != true {
		log.Panic("service not connected")
	}
	self.Dispatcher.SendMessage(c, channel, msg)
}

//broadcast to all peers on service
func (self *Service) Broadcast(msg Message) {
	for c, channel := range self.Connections {
		self.Dispatcher.SendMessage(c, channel, msg)
	}
}

//service level connection event
func (self *Service) ConnectionEvent(c *Connection, channel uint16) {
	if _, ok := self.Connections[c]; ok != false {
		log.Panic("already connected; duplicate")
	}
	self.Connections[c] = channel
	//TODO: notify object?
	self.Server.OnConnect(c)
}

//service level disconnection event
func (self *Service) DisconnectEvent(c *Connection) {
	if _, ok := self.Connections[c]; ok == false {
		log.Panic("connection does not exist")
	}
	delete(self.Connections, c)
	//TODO: notify object?
	self.Server.OnDisconnect(c)
}

//implements a service
type ServiceServer interface {
	OnConnect(c *Connection)
	OnDisconnect(c *Connection)
	RegisterMessages(d *Dispatcher)
}
