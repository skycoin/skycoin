package connection

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"time"

	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
	"github.com/SkycoinProject/skycoin/src/daemon"
	"github.com/SkycoinProject/skycoin/src/daemon/gnet"
	"github.com/SkycoinProject/skycoin/src/util/iputil"
)

const (
	// Byte size of the length prefix in message, sizeof(int32)
	messagePrefixLength = 4
)

var (
	maxMessageLength = gnet.NewConfig().MaxIncomingMessageLength
)

// Connection is a connection to the peer
type Connection struct {
	IP             string
	Port           uint16
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	conn           net.Conn
}

func init() {
	gnet.RegisterMessage(gnet.MessagePrefixFromString("INTR"), daemon.IntroductionMessage{})
}

// NewConnection constructs new connection
func NewConnection(addr string, connectTimeout, readTimeout time.Duration) (*Connection, error) {
	ip, port, err := iputil.SplitAddr(addr)
	if err != nil {
		return nil, err
	}

	return &Connection{
		IP:             ip,
		Port:           port,
		ConnectTimeout: connectTimeout,
		ReadTimeout:    readTimeout,
	}, nil
}

// Connect tries to connect to the node
func (c *Connection) Connect() error {
	conn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%v", c.IP, c.Port), c.ConnectTimeout)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

// TryReadIntroductionMessage tries to read the introduction message from the peer. If
// it succeeds - returns the received message for further processing
func (c *Connection) TryReadIntroductionMessage() (*daemon.IntroductionMessage, error) {
	reader := bufio.NewReader(c.conn)
	buf := make([]byte, 1024)
	buffer := &bytes.Buffer{}

	for {
		if err := c.conn.SetReadDeadline(time.Now().Add(c.ReadTimeout)); err != nil {
			return nil, err
		}

		bytesWritten, err := writeToBuffer(buffer, reader, buf)
		if err != nil {
			return nil, err
		}
		if bytesWritten == 0 {
			continue
		}

		// decode data
		msgBytes, err := fetchNextMessage(buffer, maxMessageLength)
		if err != nil {
			return nil, err
		}

		if len(msgBytes) > 0 {
			m, err := convertToMessage(msgBytes)
			if err != nil {
				return nil, err
			}

			introductionMessage, ok := m.(*daemon.IntroductionMessage)
			if !ok {
				// TODO: fix error message
				return nil, fmt.Errorf("wrong message")
			}

			return introductionMessage, nil
		}
	}
}

// Disconnect disconnects from the peer
func (c *Connection) Disconnect() error {
	return c.conn.Close()
}

// writeToBuffer reads from `src` and writes the received data to `dst`. Reuses `buf`
// to avoid reallocations
func writeToBuffer(dst *bytes.Buffer, src io.Reader, buf []byte) (int, error) {
	c, err := src.Read(buf)
	if err != nil {
		return 0, err
	}
	if c == 0 {
		return 0, nil
	}

	return dst.Write(buf[:c])
}

// convertToMessage converts the byte array to the gnet.Message
func convertToMessage(msg []byte) (gnet.Message, error) {
	msgID := [4]byte{}
	if len(msg) < len(msgID) {
		return nil, gnet.ErrDisconnectTruncatedMessageID
	}

	copy(msgID[:], msg[:len(msgID)])

	msg = msg[len(msgID):]
	t, ok := gnet.MessageIDReverseMap[msgID]
	if !ok {
		return nil, gnet.ErrDisconnectUnknownMessage
	}

	v := reflect.New(t)
	m, ok := (v.Interface()).(gnet.Message)
	if !ok {
		return nil, errors.New("MessageIdMaps contain non-Message")
	}

	used, err := deserializeMessage(msg, v)
	if err != nil {
		return nil, gnet.ErrDisconnectMalformedMessage
	}

	if used != uint64(len(msg)) {
		return nil, gnet.ErrDisconnectMessageDecodeUnderflow
	}

	return m, nil
}

// deserializeMessage wraps Serializer.Decode and traps panics as an error
func deserializeMessage(msg []byte, v reflect.Value) (n uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("message deserialization failed")
			}
		}
	}()

	iface := v.Interface()
	x, ok := iface.(gnet.Serializer)
	if !ok {
		return 0, errors.New("deserializeMessage object does not have Serializer interface")
	}

	return x.Decode(msg)
}

// fetchNextMessage retrieves single message from `buf` if any
func fetchNextMessage(buf *bytes.Buffer, maxMsgLength int) ([]byte, error) {
	var data []byte

	if buf.Len() > messagePrefixLength {
		prefix := buf.Bytes()[:messagePrefixLength]
		// decode message length
		tmpLength, _, err := encoder.DeserializeUint32(prefix)
		if err != nil {
			return nil, err
		}

		length := int(tmpLength)

		// Disconnect if we received an invalid length
		if length < messagePrefixLength || length > maxMsgLength {
			return []byte{}, gnet.ErrDisconnectInvalidMessageLength
		}

		if buf.Len()-messagePrefixLength < length {
			return []byte{}, nil
		}

		buf.Next(messagePrefixLength) // strip the length prefix
		data = make([]byte, length)
		_, err = buf.Read(data)
		if err != nil {
			return []byte{}, err
		}
	}

	return data, nil
}
