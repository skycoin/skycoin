package messages

import "errors"

var ERR_APP_TIMEOUT = errors.New("Application timeout")
var ERR_CONN_TIMEOUT = errors.New("Connection timeout")
var ERR_TRANSPORT_TIMEOUT = errors.New("Transport timeout")
var ERR_DISCONNECTED = errors.New("Connection if off")
var ERR_ROUTE_EXISTS = errors.New("Route already exists")
var ERR_ROUTE_DOESNT_EXIST = errors.New("Route doesn't exists")
var ERR_TRANSPORT_EXISTS = errors.New("Transport already exists")
var ERR_TRANSPORT_DOESNT_EXIST = errors.New("Transport doesn't exists")
var ERR_UNKNOWN_MESSAGE_TYPE = errors.New("Unknown message type")
var ERR_NO_TRANSPORT_TO_NODE = errors.New("No transport to node")
var ERR_ALREADY_CONNECTED = errors.New("Nodes are already connected")
var ERR_NO_ROUTE = errors.New("No route between nodes")
var ERR_NODE_NOT_FOUND = errors.New("Node not found")
var ERR_TOO_MANY_NODES = errors.New("Too many nodes, should be 100 or less")
var ERR_WRONG_NUMBER_ARGS = errors.New("Wrong number of arguments")
var ERR_NODE_NUM_OUT_OF_RANGE = errors.New("Node number is out of range")
var ERR_CONNECTED_TO_ITSELF = errors.New("Node cannot be connected to itself")
var ERR_NO_CLIENT_RESPONSE_CHANNEL = errors.New("Client response channel doesn't exist")
var ERR_NO_TRANSPORT_ACK_CHANNEL = errors.New("Transport ack channel doesn't exist")
