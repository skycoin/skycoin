package connection

import ("github.com/skycoin/skycoin/src/cipher")

/*

An Endpoint, will be used for an application to send/receive []byte messages

An Endpoint connects to a remote server
An endpoint sends []byte messages
An endpoint receives []byte messages

An endpoint has a "window"

End-point connection

- takes public key to connect to
- queries route service to obtain a route through network
- attempts to connect to remote endpoint through a route

*/


/*
Record:
- how much data is sent
- how much data is received
- uint64 time, when packet fragment is sent
- uint64 time, when packet fragment ACK is received
*/


/*
Message Format
- 4 byte, length prefixed byte array
- 2 byte message type int16
*/

/*
Fragmentation
- a node accepts a length prefixed message to be send
- example: 4500 bytes
- the message must be fragmented and cut into smaller messages below the MTU
- the messages are reassembled on the other side
- when a message is received, the other side sends an ACK message

- messages are assigned ids, as they go out
- message fragments are assigned ids

- The output fragments are a fixed length (example: 512 bytes)
- 
*/

/*

Multiple Messages or commands can be wrapped in a messsage fragment
- one command is "message fragment"
- one command is "ACK"

Command Channel
- length prefixed messages are send to the node
- the messages are wrapped and unwrapped and rewrapped
*/

/*
Messages:



*/

struct PacketFragment {
	PacketId uint32 //for ack
	FragmentN uint16 //fragment number
	FragmentMax uint16 //number of fragments
	//Length uint32
	Data []byte
}

struct AckFragment {
	FragmentId uint32
	PacketId uint32
}

/*
Code:
*/

//put static settings here
struct EndpointConfig {
	Destination cipher.PubKey
	uint32 TransmissionUnit //size of packet fragments leaving
	uint MaxWindow // number of bytes outgoing pending, before send blocks

	bool InOrder //will only report packets in order of stream 
	//if InOrder is true, then only report packets in sequence for receipt
	//otherwise report packets to receiver as they come in

}

// keep track of messages in, pending processing
// keep track of 
struct Endpoint {
	Config EndpointConfig 
	//Route []cipher.PubKey


	Node *(chan []byte) //node that receives length prefixed messages, []byte
	MsgIn *(chan []byte) //this is where node will put incoming messages
}

//blocks until connected
func (self *Endpoint) Connect(Destination cipher.PubKey) {

}

//send a message over the channel
//blocks if window is full
func (self *Endpoint) Send(msg []byte) {


}

//check if there are any incoming messages
func (self *Endpoint) Poll() {
	//check if channel is empty
}