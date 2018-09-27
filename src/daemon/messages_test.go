package daemon

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon/gnet"
)

func setupMsgEncoding() {
	gnet.EraseMessages()
	var messagesConfig = NewMessagesConfig()
	messagesConfig.Register()
}

func TestIntroductionMessage(t *testing.T) {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	pubkey, _ := cipher.GenerateKeyPair()
	pubkey2, _ := cipher.GenerateKeyPair()

	type mirrorPortResult struct {
		port  uint16
		exist bool
	}

	type daemonMockValue struct {
		version                    uint32
		mirror                     uint32
		isDefaultConnection        bool
		isMaxConnectionsReached    bool
		isMaxConnectionsReachedErr error
		setHasIncomingPortErr      error
		getMirrorPortResult        mirrorPortResult
		recordMessageEventErr      error
		pubkey                     cipher.PubKey
		disconnectReason           gnet.DisconnectReason
		disconnectErr              error
		addPeerArg                 string
		addPeerErr                 error
	}

	tt := []struct {
		name      string
		addr      string
		mockValue daemonMockValue
		intro     *IntroductionMessage
		err       error
	}{
		{
			name: "INTR message without extra bytes",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:  10000,
				version: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Port:    6000,
				Version: 1,
				valid:   true,
			},
			err: nil,
		},
		{
			name: "INTR message with pubkey",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:  10000,
				version: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey: pubkey,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Port:    6000,
				Version: 1,
				valid:   true,
				Extra:   pubkey[:],
			},
			err: nil,
		},
		{
			name: "INTR message with pubkey",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:  10000,
				version: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey: pubkey,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Port:    6000,
				Version: 1,
				valid:   true,
				Extra:   pubkey[:],
			},
			err: nil,
		},
		{
			name: "INTR message with pubkey and additional data",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:  10000,
				version: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey: pubkey,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Port:    6000,
				Version: 1,
				valid:   true,
				Extra:   append(pubkey[:], []byte("additional data")...),
			},
			err: nil,
		},
		{
			name: "INTR message with different pubkey",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:  10000,
				version: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey:           pubkey,
				disconnectReason: ErrDisconnectBlockchainPubkeyNotMatched,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Port:    6000,
				Version: 1,
				valid:   true,
				Extra:   pubkey2[:],
			},
			err: ErrDisconnectBlockchainPubkeyNotMatched,
		},
		{
			name: "INTR message with invalid pubkey",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:  10000,
				version: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey:           pubkey,
				disconnectReason: ErrDisconnectInvalidExtraData,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Port:    6000,
				Version: 1,
				valid:   true,
				Extra:   []byte("invalid extra data"),
			},
			err: ErrDisconnectInvalidExtraData,
		},
		{
			name: "Disconnect self connection",
			mockValue: daemonMockValue{
				mirror:           10000,
				disconnectReason: ErrDisconnectSelf,
			},
			intro: &IntroductionMessage{
				Mirror: 10000,
			},
			err: ErrDisconnectSelf,
		},
		{
			name: "Invalid version",
			mockValue: daemonMockValue{
				mirror:           10000,
				version:          1,
				disconnectReason: ErrDisconnectInvalidVersion,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Version: 0,
			},
			err: ErrDisconnectInvalidVersion,
		},
		{
			name: "Invalid address",
			addr: "121.121.121.121",
			mockValue: daemonMockValue{
				mirror:           10000,
				version:          1,
				disconnectReason: ErrDisconnectOtherError,
				pubkey:           pubkey,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Version: 1,
				Port:    6000,
			},
			err: ErrDisconnectOtherError,
		},
		{
			name: "incomming connection",
			addr: "121.121.121.121:12345",
			mockValue: daemonMockValue{
				mirror:                  10000,
				version:                 1,
				isDefaultConnection:     true,
				isMaxConnectionsReached: true,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey:     pubkey,
				addPeerArg: "121.121.121.121:6000",
				addPeerErr: nil,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Version: 1,
				Port:    6000,
				valid:   true,
			},
		},
		{
			name: "Connect twice",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:              10000,
				version:             1,
				isDefaultConnection: true,
				getMirrorPortResult: mirrorPortResult{
					exist: true,
				},
				pubkey:           pubkey,
				addPeerArg:       "121.121.121.121:6000",
				addPeerErr:       nil,
				disconnectReason: ErrDisconnectConnectedTwice,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Version: 1,
				Port:    6000,
			},
			err: ErrDisconnectConnectedTwice,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mc := &gnet.MessageContext{Addr: tc.addr}
			tc.intro.c = mc

			d := &MockDaemoner{}
			d.On("DaemonConfig").Return(DaemonConfig{Version: int32(tc.mockValue.version)})
			d.On("Mirror").Return(tc.mockValue.mirror)
			d.On("IsDefaultConnection", tc.addr).Return(tc.mockValue.isDefaultConnection)
			d.On("SetHasIncomingPort", tc.addr).Return(tc.mockValue.setHasIncomingPortErr)
			d.On("GetMirrorPort", tc.addr, tc.intro.Mirror).Return(tc.mockValue.getMirrorPortResult.port, tc.mockValue.getMirrorPortResult.exist)
			d.On("RecordMessageEvent", tc.intro, mc).Return(tc.mockValue.recordMessageEventErr)
			d.On("ResetRetryTimes", tc.addr)
			d.On("BlockchainPubkey").Return(tc.mockValue.pubkey)
			d.On("Disconnect", tc.addr, tc.mockValue.disconnectReason).Return(tc.mockValue.disconnectErr)
			d.On("IncreaseRetryTimes", tc.addr)
			d.On("RemoveFromExpectingIntroductions", tc.addr)
			d.On("IsMaxDefaultConnectionsReached").Return(tc.mockValue.isMaxConnectionsReached, tc.mockValue.isMaxConnectionsReachedErr)
			d.On("AddPeer", tc.mockValue.addPeerArg).Return(tc.mockValue.addPeerErr)

			err := tc.intro.Handle(mc, d)
			require.Equal(t, tc.err, err)
		})
	}
}
