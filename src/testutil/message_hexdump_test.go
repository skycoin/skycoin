package testutil

import (
	"testing"
	"github.com/skycoin/skycoin/src/daemon"
)


func TestIntroductionMessage(t *testing.T) {
	var message = daemon.NewIntroductionMessage(1234,5,7890)
	HexDump(message)
}

func TestGetPeersMessage(t *testing.T) {

}

func TestGivePeersMessage(t *testing.T) {

}

func TestPingMessage(t *testing.T) {

}

func TestPongMessage(t *testing.T) {

}

func TestGetBlocksMessage(t *testing.T) {

}

func TestGiveBlocksMessage(t *testing.T) {

}

func TestAnnounceBlocksMessage(t *testing.T) {

}

func TestGetTxnsMessage(t *testing.T) {

}

func TestGiveTxnsMessage(t *testing.T) {

}

func TestAnnounceTxnsMessage(t *testing.T) {

}
