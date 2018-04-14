package daemon

import (
	"github.com/skycoin/skycoin/src/daemon/gnet"
)

func setupMsgEncoding() {
	gnet.EraseMessages()
	var messagesConfig = NewMessagesConfig()
	messagesConfig.Register()
}
