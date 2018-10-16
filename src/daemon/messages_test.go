package daemon

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
)

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
		protocolVersion            uint32
		minProtocolVersion         uint32
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
				mirror:          10000,
				protocolVersion: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
			},
			intro: &IntroductionMessage{
				Mirror:          10001,
				Port:            6000,
				Version:         1,
				validationError: nil,
			},
			err: nil,
		},
		{
			name: "INTR message with pubkey",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:          10000,
				protocolVersion: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey: pubkey,
			},
			intro: &IntroductionMessage{
				Mirror:          10001,
				Port:            6000,
				Version:         1,
				validationError: nil,
			},
			err: nil,
		},
		{
			name: "INTR message with pubkey",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:          10000,
				protocolVersion: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey: pubkey,
			},
			intro: &IntroductionMessage{
				Mirror:          10001,
				Port:            6000,
				Version:         1,
				validationError: nil,
			},
			err: nil,
		},
		{
			name: "INTR message with pubkey and additional data",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:          10000,
				protocolVersion: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey: pubkey,
			},
			intro: &IntroductionMessage{
				Mirror:          10001,
				Port:            6000,
				Version:         1,
				validationError: nil,
			},
			err: nil,
		},
		{
			name: "INTR message with different pubkey",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:          10000,
				protocolVersion: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey:           pubkey,
				disconnectReason: ErrDisconnectBlockchainPubkeyNotMatched,
			},
			intro: &IntroductionMessage{
				Mirror:          10001,
				Port:            6000,
				Version:         1,
				validationError: nil,
			},
			err: ErrDisconnectBlockchainPubkeyNotMatched,
		},
		{
			name: "INTR message with invalid pubkey",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:          10000,
				protocolVersion: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey:           pubkey,
				disconnectReason: ErrDisconnectInvalidExtraData,
			},
			intro: &IntroductionMessage{
				Mirror:          10001,
				Port:            6000,
				Version:         1,
				validationError: nil,
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
			name: "Version below minimum supported version",
			mockValue: daemonMockValue{
				mirror:             10000,
				protocolVersion:    1,
				minProtocolVersion: 2,
				disconnectReason:   ErrDisconnectVersionNotSupported,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Version: 0,
			},
			err: ErrDisconnectVersionNotSupported,
		},
		{
			name: "Invalid address",
			addr: "121.121.121.121",
			mockValue: daemonMockValue{
				mirror:           10000,
				protocolVersion:  1,
				disconnectReason: ErrDisconnectIncomprehensibleError,
				pubkey:           pubkey,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Version: 1,
				Port:    6000,
			},
			err: ErrDisconnectIncomprehensibleError,
		},
		{
			name: "incomming connection",
			addr: "121.121.121.121:12345",
			mockValue: daemonMockValue{
				mirror:                  10000,
				protocolVersion:         1,
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
				Mirror:          10001,
				Version:         1,
				Port:            6000,
				validationError: nil,
			},
		},
		{
			name: "Connect twice",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:              10000,
				protocolVersion:     1,
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
			d.On("DaemonConfig").Return(DaemonConfig{
				ProtocolVersion:    int32(tc.mockValue.protocolVersion),
				MinProtocolVersion: int32(tc.mockValue.minProtocolVersion),
			})
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

func TestMessageEncodeDecode(t *testing.T) {
	update := false

	introPubKey := cipher.MustPubKeyFromHex("03cd7dfcd8c3452d1bb5d9d9e34dd95d6848cb9f66c2aad127b60578f4be7498f2")

	cases := []struct {
		goldenFile string
		obj        interface{}
		msg        interface{}
	}{
		{
			goldenFile: "intro-msg.golden",
			obj:        &IntroductionMessage{},
			msg: &IntroductionMessage{
				Mirror:  99998888,
				Port:    8888,
				Version: 12341234,
			},
		},
		{
			goldenFile: "intro-msg-pubkey.golden",
			obj:        &IntroductionMessage{},
			msg: &IntroductionMessage{
				Mirror:  99998888,
				Port:    8888,
				Version: 12341234,
				Extra:   introPubKey[:],
			},
		},
		{
			goldenFile: "get-peers-msg.golden",
			obj:        &GetPeersMessage{},
			msg:        &GetPeersMessage{},
		},
		{
			goldenFile: "give-peers-msg.golden",
			obj:        &GivePeersMessage{},
			msg: &GivePeersMessage{
				Peers: []IPAddr{
					{
						IP:   12345678,
						Port: 1234,
					},
					{
						IP:   87654321,
						Port: 4321,
					},
				},
			},
		},
		{
			goldenFile: "ping-msg.golden",
			obj:        &PingMessage{},
			msg:        &PingMessage{},
		},
		{
			goldenFile: "pong-msg.golden",
			obj:        &PongMessage{},
			msg:        &PongMessage{},
		},
		{
			goldenFile: "get-blocks-msg.golden",
			obj:        &GetBlocksMessage{},
			msg: &GetBlocksMessage{
				LastBlock:       999988887777,
				RequestedBlocks: 888899997777,
			},
		},
		{
			goldenFile: "give-blocks-msg.golden",
			obj:        &GiveBlocksMessage{},
			msg: &GiveBlocksMessage{
				Blocks: []coin.SignedBlock{
					{
						Sig: cipher.MustSigFromHex("8cf145e9ef4a4a5254bc57798a7a61dfed238768f94edc5635175c6b91bccd8ec1555da603c5e31b018e135b82b1525be8a92973c468a74b5b40b8da189cb465eb"),
						Block: coin.Block{
							Head: coin.BlockHeader{
								Version:  1,
								Time:     1538036613,
								BkSeq:    9999999999,
								Fee:      1234123412341234,
								PrevHash: cipher.MustSHA256FromHex("59cb7d0e2ce8a03d1054afcc28a22fe864a8813460d241db38c59d10e7c29132"),
								BodyHash: cipher.MustSHA256FromHex("6d421469409591f0c3112884c8cf10f8bca5d8ab87c9c30dea2ea73b6751bbf9"),
								UxHash:   cipher.MustSHA256FromHex("6ea6a972cf06d25908b29953aeddb68c3b6f3a9903e8f964dc89b0abc0645dea"),
							},
							Body: coin.BlockBody{
								Transactions: coin.Transactions{
									{
										Length:    43214321,
										Type:      1,
										InnerHash: cipher.MustSHA256FromHex("cbedf8ef0bda91afc6a180eea0dddf8e3a986b6b6f87f70e8bffc63c6fbaa4e6"),
										Sigs: []cipher.Sig{
											cipher.MustSigFromHex("1cfd7a4db3a52a85d2a86708695112b6520acc8dc83c86e8da67915199fdf04964c168543598ab07c2b99c292899890891950364c2bf66f1aaa6d6a66a5c9a73ff"),
											cipher.MustSigFromHex("442167c6b3d13957bc32f83182c7f4fda0bb6bde893a41a6a04cdd8eecee0048d03a57eb2af04ea6050e1f418769c94c7f12fad9287dc650e6b307fdfce6b42a59"),
										},
										In: []cipher.SHA256{
											cipher.MustSHA256FromHex("536f0a1a915fadfa3a2720a0615641827ff67394d2b2149d6db63b8c619e14af"),
											cipher.MustSHA256FromHex("64ba5f01f90f97f84999f13aeaa75fed8d5b3e4a3a4a093dedf4795969e8bd27"),
										},
										Out: []coin.TransactionOutput{
											{
												Address: cipher.MustDecodeBase58Address("23FF4fshzD8tZk2d88P22WATfzUpNQF1x85"),
												Coins:   987987987,
												Hours:   789789789,
											},
											{
												Address: cipher.MustDecodeBase58Address("29V2iRpZAqHiFZHHRqaZLArZZuTcZM5owqT"),
												Coins:   123123,
												Hours:   321321,
											},
										},
									},
									{
										Length:    98769876,
										Type:      0,
										InnerHash: cipher.MustSHA256FromHex("46856af925fde9a1652d39eea479dd92589a741451a0228402e399fae02f8f3d"),
										Sigs: []cipher.Sig{
											cipher.MustSigFromHex("92e289792200518df9a82cf9dddd1f334bf0d47fb0ed4ff70c25403f39577af5ab24ef2d02a11cf6b76e6bd017457ad60d6ca85c0567c21f5c62599c93ee98e18c"),
											cipher.MustSigFromHex("e995da86ed87640ecb44e624074ba606b781aa0cbeb24e8c27ff30becf7181175479c0d74d93fe1e8692bba628b5cf532ca80fed4135148d84e6ecc2a762a10b19"),
										},
										In: []cipher.SHA256{
											cipher.MustSHA256FromHex("69b14a7ee184f24b95659d6887101ef7c921fa7977d95c73fbc0c4d0d22671bc"),
											cipher.MustSHA256FromHex("3a050b4ec33ec9ad2c789f24655ab1c8f7691d3a1c3d0e05cc14b022b4c360ea"),
										},
										Out: []coin.TransactionOutput{
											{
												Address: cipher.MustDecodeBase58Address("XvvjeyGcTBVXDXmfJoTUseFiqHvm12C6oQ"),
												Coins:   15,
												Hours:   1237882,
											},
											{
												Address: cipher.MustDecodeBase58Address("fQXVLq9fbCC9XVxDKLAGDLXmQYjPy59j22"),
												Coins:   2102123,
												Hours:   1003,
											},
										},
									},
								},
							},
						},
					},
					{
						Sig: cipher.MustSigFromHex("8015c8776de577d89c29d1cbd1d558ba4855dec94ba58f6c67d55ece5c85708b9906bd0b72b451e27008f3938fcec42c1a28ddac336ae8206d8e6443b95dde966c"),
						Block: coin.Block{
							Head: coin.BlockHeader{
								Version:  0,
								Time:     1427248825,
								BkSeq:    100,
								Fee:      120939323123,
								PrevHash: cipher.MustSHA256FromHex("04d40b5d27c539ab9d98934628604baef7dbfb1c35ddf9c0f96a67f6b061fa26"),
								BodyHash: cipher.MustSHA256FromHex("9a67fbb00216ae99f334d4efa2c9c42a25aac5d1a5bbb2058fe5705cfe0e30ea"),
								UxHash:   cipher.MustSHA256FromHex("58981d30da11be3c8e9dd8fdb7b51b48ba13dc0214cf211251308985bf089f76"),
							},
							Body: coin.BlockBody{
								Transactions: coin.Transactions{
									{
										Length:    128,
										Type:      99,
										InnerHash: cipher.MustSHA256FromHex("e943fd54a8071bb0ae92800c23c5a26443b5e5bf9b9321cefcdd9e80f518c37e"),
										Sigs: []cipher.Sig{
											cipher.MustSigFromHex("cff49d1d450db812d42748d4f7001e03a1dd2b98afcbb62eca1b3b1fa137e5095a0368250aabd3976008afe61471ecd31ed99185c3df49269d9aada4ca1dd2eecb"),
											cipher.MustSigFromHex("1313e5a80d6d9386fe2dffa13afba7277402f029d411e60f99b3806fee547d6157ca2d8d6407df3e858d6f3f58902f460412611282a0dec2468e41a2c5a39cc93e"),
										},
										In: []cipher.SHA256{
											cipher.MustSHA256FromHex("6a76c83b7b75075e2e34405e21d5e8d37adb69e4e6487a6179944ea7e04bc7db"),
											cipher.MustSHA256FromHex("a7555179a255e6a7dddb6121bd4c2259f75ebc321345be26b690f34094012f95"),
										},
										Out: []coin.TransactionOutput{
											{
												Address: cipher.MustDecodeBase58Address("2RmSTGbj5qaFT1WvKGz4SobaT4xSb9GvaCh"),
												Coins:   12301923233,
												Hours:   39932,
											},
											{
												Address: cipher.MustDecodeBase58Address("uA8XQnNzS4kit9DFzybyVSpWDEDy62MXur"),
												Coins:   9945924,
												Hours:   9030300895893902,
											},
										},
									},
									{
										Length:    1304,
										Type:      255,
										InnerHash: cipher.MustSHA256FromHex("d92057e9a4874aa876b7fd20074d78a4d890c2d3af483a10206f243308586763"),
										Sigs: []cipher.Sig{
											cipher.MustSigFromHex("394d53cc0bfeef11cc94bf39316d555549cf1a1afd14920be7d065e7940cc60752b8ade8c37991307a5681b06e0445c1c19ceb0e6611fd4593dcc65d18975c87be"),
											cipher.MustSigFromHex("50ad670bc672558c235653b6396135352bfbc8eec575de3cffce65d5a07076082f9694880eb6b1e708eb8fb39d21a96dd99615b5759fc917c3fdd4d9845489119b"),
										},
										In: []cipher.SHA256{
											cipher.MustSHA256FromHex("f37057407a6b5b103218abdfc5b5527f8abcc229256c912ec81ac6d72b68454e"),
											cipher.MustSHA256FromHex("9cd1fccddb5895ab77cd419802430e16a1e05f0f796d026fc69961c5c308b766"),
										},
										Out: []coin.TransactionOutput{
											{
												Address: cipher.MustDecodeBase58Address("MNf67cWXYmSizin4XUtGnFfQQzxkvNqCEH"),
												Coins:   1,
												Hours:   1,
											},
											{
												Address: cipher.MustDecodeBase58Address("HEkH8R1Uc58mAjZqGM15cqF4QMqG4mu4ry"),
												Coins:   1,
												Hours:   0,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			goldenFile: "announce-blocks-msg.golden",
			obj:        &AnnounceBlocksMessage{},
			msg: &AnnounceBlocksMessage{
				MaxBkSeq: 50000,
			},
		},
		{
			goldenFile: "announce-txns-msg.golden",
			obj:        &AnnounceTxnsMessage{},
			msg: &AnnounceTxnsMessage{
				Transactions: []cipher.SHA256{
					cipher.MustSHA256FromHex("23dc4b68c0fc790989bb82f04b9d5174baab6f0f6808ed35be9b93cb73c69108"),
					cipher.MustSHA256FromHex("2be4b0155c1ab9613007fe522e3b12bac4be79800a19bc8cd8ca343868caa583"),
				},
			},
		},
		{
			goldenFile: "get-txns-msg.golden",
			obj:        &GetTxnsMessage{},
			msg: &GetTxnsMessage{
				Transactions: []cipher.SHA256{
					cipher.MustSHA256FromHex("335b63b0f335c6aee5e7e1b3c62dd09bb6074e38b48e2469e294a019d5ae5aa1"),
					cipher.MustSHA256FromHex("619a367f4e5dee741348366899237ddc920335fc847ccafdf2d32ed57bb7b385"),
				},
			},
		},
		{
			goldenFile: "give-txns-msg.golden",
			obj:        &GiveTxnsMessage{},
			msg: &GiveTxnsMessage{
				Transactions: coin.Transactions{
					{
						Length:    256,
						Type:      0,
						InnerHash: cipher.MustSHA256FromHex("1773d8901df96bba4c6d65499e11e6ec73a9978c611d1463898ffbc2b49773fc"),
						Sigs: []cipher.Sig{
							cipher.MustSigFromHex("a711880ae54d1b6b9adade2ef1e743d6d539a78b0cecf1af08107e467956de80ef1d49fb5e896c9d0870ef8bf8a4d328ca0ecf7c1956866867ec56064e68f8a374"),
							cipher.MustSigFromHex("f9890ddd93f9479e364261ebc647326d2fd57e50b7728795adbf507c956f9eb44f77207b528700c4cef338290cdfc17f814dc3d94e3d711e92492ecc7b8abef808"),
						},
						In: []cipher.SHA256{
							cipher.MustSHA256FromHex("703f84ee0702b44fc89ce573a239d5fbf185bf5d4e7fc8f4930262bcda1e8fb0"),
							cipher.MustSHA256FromHex("c9e904862da01f2d7676c12c4342dde36d9a9a9d25be5351e2b57fae6f426bb9"),
						},
						Out: []coin.TransactionOutput{
							{
								Address: cipher.MustDecodeBase58Address("29VEn56iRr2TpVVpPoPxUJPfFWuhbLSBRdU"),
								Coins:   1111111111111111111,
								Hours:   9999999999999999999,
							},
							{
								Address: cipher.MustDecodeBase58Address("2bqs99tysFtfs8QPT81kpZWnzTT1rWd8xtQ"),
								Coins:   9922581002,
								Hours:   9932900022223334,
							},
						},
					},
					{
						Length:    13043,
						Type:      128,
						InnerHash: cipher.MustSHA256FromHex("a9da3e4acb1892a000c1b658a64d4e420d0c381862928ab820fb3f3a534a9674"),
						Sigs: []cipher.Sig{
							cipher.MustSigFromHex("7bbbdfd58c0533aed95f18d9413e0e0517892350eaf132eadf7a9a03d4a974ca0bc074abc001f86a34cf66c10f832dbcca20c2c67b5e8517f4ff0e1d0123fecb21"),
							cipher.MustSigFromHex("68732b78ac3a4e2fe146b8819c8b1c0b126a0188008c9c7c98fee965beba039778010ff7b0379dadeeadbbc42f9541ce4ad3c8cec12108d3aa58aca583bddd0df0"),
						},
						In: []cipher.SHA256{
							cipher.MustSHA256FromHex("766d6f6ed56599a91759c75466e3f09b9d6d5995b58dd5bbfba5af10b1a8cdea"),
							cipher.MustSHA256FromHex("2c7989f47524721bb2c7a7f967208c9b1c01829c9a55addf22d066e5c55ab3ac"),
						},
						Out: []coin.TransactionOutput{
							{
								Address: cipher.MustDecodeBase58Address("24iFsYHzVfYXo8cvWg1jhetpTMNvHH7j6AX"),
								Coins:   1123103123,
								Hours:   123000,
							},
							{
								Address: cipher.MustDecodeBase58Address("JV5xJ33po1Bj3dXZT3SYA3ZmnTibREFxxd"),
								Coins:   999999,
								Hours:   9043285343,
							},
						},
					},
				},
			},
		},
	}

	if update {
		for _, tc := range cases {
			t.Run(tc.goldenFile, func(t *testing.T) {
				fn := filepath.Join("testdata/", tc.goldenFile)

				f, err := os.Create(fn)
				require.NoError(t, err)
				defer f.Close()

				b := encoder.Serialize(tc.msg)
				_, err = f.Write(b)
				require.NoError(t, err)
			})
		}
	}

	for _, tc := range cases {
		t.Run(tc.goldenFile, func(t *testing.T) {
			fn := filepath.Join("testdata/", tc.goldenFile)

			f, err := os.Open(fn)
			require.NoError(t, err)
			defer f.Close()

			d, err := ioutil.ReadAll(f)
			require.NoError(t, err)

			err = encoder.DeserializeRaw(d, tc.obj)
			require.NoError(t, err)

			require.Equal(t, tc.msg, tc.obj)

			d2 := encoder.Serialize(tc.msg)
			require.Equal(t, d, d2)
		})
	}
}
