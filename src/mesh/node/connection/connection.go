package connection

import (
	"fmt"
	"os"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/domain"
	"github.com/skycoin/skycoin/src/mesh/serialize"
	"github.com/skycoin/skycoin/src/mesh/transport"

	uuid "github.com/satori/go.uuid"
)

func init() {
	ConnectionManager = Connection{
		serializer: serialize.NewSerializer(),
	}

	ConnectionManager.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{1}, domain.UserMessage{})
	ConnectionManager.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{2}, domain.SetRouteMessage{})
	ConnectionManager.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{3}, domain.RefreshRouteMessage{})
	ConnectionManager.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{4}, domain.DeleteRouteMessage{})
	ConnectionManager.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{5}, domain.SetRouteReply{})

	ConnectionManager.serializer.RegisterMessageForSerialization(serialize.MessagePrefix{6}, domain.AddNodeMessage{})
}

var ConnectionManager Connection

type Connection struct {
	serializer *serialize.Serializer
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func (self *Connection) GetMaximumContentLength(toPeer cipher.PubKey, transport transport.Transport) uint64 {
	transportSize := transport.GetMaximumMessageSizeToPeer(toPeer)
	empty := domain.UserMessage{}
	emptySerialized := self.serializer.SerializeMessage(empty)
	if (uint)(len(emptySerialized)) >= transportSize {
		return 0
	}
	return (uint64)(transportSize) - (uint64)(len(emptySerialized))
}

func (self *Connection) FragmentMessage(fullContents []byte, toPeer cipher.PubKey, transport transport.Transport, base domain.MessageBase) []domain.UserMessage {
	ret_noCount := make([]domain.UserMessage, 0)
	maxContentLength := self.GetMaximumContentLength(toPeer, transport)
	fmt.Fprintf(os.Stdout, "MaxContentLength: %v\n", maxContentLength)
	remainingBytes := fullContents[:]
	messageId := (domain.MessageID)(uuid.NewV4())
	for len(remainingBytes) > 0 {
		nBytesThisMessage := min(maxContentLength, (uint64)(len(remainingBytes)))
		bytesThisMessage := remainingBytes[:nBytesThisMessage]
		remainingBytes = remainingBytes[nBytesThisMessage:]
		message := domain.UserMessage{
			MessageBase: base,
			MessageID:   messageId,
			Index:       (uint64)(len(ret_noCount)),
			Count:       0,
			Contents:    bytesThisMessage,
		}
		ret_noCount = append(ret_noCount, message)
	}
	ret := make([]domain.UserMessage, 0)
	for _, message := range ret_noCount {
		message.Count = (uint64)(len(ret_noCount))
		ret = append(ret, message)
	}
	fmt.Fprintf(os.Stdout, "Message fragmented in %v packets.\n", len(ret))
	return ret
}

func (self *Connection) DeserializeMessage(msg []byte) (interface{}, error) {
	return self.serializer.UnserializeMessage(msg)
}

func (self *Connection) SerializeMessage(msg interface{}) []byte {
	return self.serializer.SerializeMessage(msg)
}
