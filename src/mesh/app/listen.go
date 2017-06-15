package app

import (
	"io"
	"net"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

func getFullMessage(conn net.Conn) ([]byte, error) {
	sizeMessage := make([]byte, 8)

	_, err := conn.Read(sizeMessage)
	if err != nil {
		return nil, err
	}

	size := messages.BytesToNum(sizeMessage)

	message := make([]byte, size)

	_, err = io.ReadFull(conn, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}
