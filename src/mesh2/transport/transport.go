package transport

import ()

type Transport struct {
	IncomingChannel chan ([]byte)
}

func (self *Node) New() {
	self.IncomingChannel = make(chan []byte, 1024)
}

func (self *Node) Shutdown() {
	close(self.IncomingChannel)
}

//move node forward on tick, process events
func (self *Node) Tick() {
	//process incoming messages
	for msg := range self.IncomingChannel {
		//process our incoming messages
		//fmt.Println(msg)

		switch messsages.GetMessageType(msg) {

		//InRouteMessage is the only message coming in to node from transports
		case messages.OutRouteMessage:

			var m1 messages.OutRouteMessage
			messages.Deserialize(msg, m1)
			//get message and put into the queue to be sent out
		}

	}
}

//inject an incoming message from the transport
func (self *Node) InjectNodeMessage([]byte) {

}
