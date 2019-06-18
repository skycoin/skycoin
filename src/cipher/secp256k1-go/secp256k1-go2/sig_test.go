package secp256k1go

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestSigForceLowS(t *testing.T) {
	// forceLowS was a hardcoded parameter that forced the Signature's S point
	// to be "low", i.e. below the half-order of the curve.
	// This is necessary to break signature malleability and should always be on,
	// so the forceLowS parameter was removed, and the current code is equivalent
	// to forceLowS=true.

	// Check that forceLowS forces the S point to the lower half of the curve
	var sec, msg, non Number
	sec.SetHex("7A642C99F7719F57D8F4BEB11A303AFCD190243A51CED8782CA6D3DBE014D146")
	msg.SetHex("DD72CBF2203C1A55A411EEC4404AF2AFB2FE942C434B23EFE46E9F04DA8433CA")
	non.SetHex("9F3CD9AB0F32911BFDE39AD155F527192CE5ED1F51447D63C4F154C118DA598E")

	// The signature when forceLowS is true (not malleable)
	sigHexLowS := "8c20a668be1b5a910205de46095023fe4823a3757f4417114168925f28193bff520ce833da9313d726f2a4d481e3195a5dd8e935a6c7f4dc260ed4c66ebe6da7"
	// The signature when forceLowS is false (malleable)
	// "8c20a668be1b5a910205de46095023fe4823a3757f4417114168925f28193bffadf317cc256cec28d90d5b2b7e1ce6a45cd5f3b10880ab5f99c389c66177d39a"

	var sig Signature
	var recid int
	res := sig.Sign(&sec, &msg, &non, &recid)
	if res != 1 {
		t.Error("res failed", res)
		return
	}

	if recid != 0 {
		t.Error("recid should be 0 because of forceLowS")
	}
	if sigHexLowS != hex.EncodeToString(sig.Bytes()) {
		t.Error("forceLowS did not modify the S point as expected")
	}
}

func TestSigRecover(t *testing.T) {
	cases := []struct {
		r     string
		s     string
		msg   string
		recid int
		x     string
		y     string
	}{
		{
			r:     "6028b9e3a31c9e725fcbd7d5d16736aaaafcc9bf157dfb4be62bcbcf0969d488",
			s:     "036d4a36fa235b8f9f815aa6f5457a607f956a71a035bf0970d8578bf218bb5a",
			msg:   "9cff3da1a4f86caf3683f865232c64992b5ed002af42b321b8d8a48420680487",
			recid: 0,
			x:     "56dc5df245955302893d8dda0677cc9865d8011bc678c7803a18b5f6faafec08",
			y:     "54b5fbdcd8fac6468dac2de88fadce6414f5f3afbb103753e25161bef77705a6",
		},
		{
			r:     "b470e02f834a3aaafa27bd2b49e07269e962a51410f364e9e195c31351a05e50",
			s:     "560978aed76de9d5d781f87ed2068832ed545f2b21bf040654a2daff694c8b09",
			msg:   "9ce428d58e8e4caf619dc6fc7b2c2c28f0561654d1f80f322c038ad5e67ff8a6",
			recid: 1,
			x:     "15b7e7d00f024bffcd2e47524bb7b7d3a6b251e23a3a43191ed7f0a418d9a578",
			y:     "bf29a25e2d1f32c5afb18b41ae60112723278a8af31275965a6ec1d95334e840",
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s,%s", tc.r, tc.s), func(t *testing.T) {
			var sig Signature
			var pubkey, expectedPubkey XY
			var msg Number

			sig.R.SetHex(tc.r)
			sig.S.SetHex(tc.s)
			msg.SetHex(tc.msg)

			expectedPubkey.X.SetHex(tc.x)
			expectedPubkey.Y.SetHex(tc.y)

			if sig.Recover(&pubkey, &msg, tc.recid) {
				if !expectedPubkey.X.Equals(&pubkey.X) {
					t.Error("X mismatch")
				}
				if !expectedPubkey.Y.Equals(&pubkey.Y) {
					t.Error("Y mismatch")
				}
			} else {
				t.Error("sig.Recover failed")
			}
		})
	}
}

func TestSigRecover2(t *testing.T) {
	cases := []struct {
		msg          string
		sig          string
		pubkey       string
		recoverFails bool
	}{
		{
			msg:    "016b81623cf98f45879f3a48fa34af77dde44b2ffa0ddd2bf9edb386f76ec0ef",
			sig:    "d2a8ec2b29ce3cf3e6048296188adff4b5dfcb337c1d1157f28654e445bb940b4e47d6b0c7ba43d072bf8618775f123a435e8d1a150cb39bbb1aa80da8c57ea100",
			pubkey: "03c0b0e24d55255f7aefe3da7a947a63028b573f45356a9c22e9a3c103fd00c3d1",
		},

		{
			msg:    "176b81623cf98f45879f3a48fa34af77dde44b2ffa0ddd2bf9edb386f76ec0ef",
			sig:    "d2a8ec2b20ce3cf3e6048296188adff4b5dfcb337c1d1157f28654e445bb940b4e47d6b0c7ba43d072bf8618775f123a435e8d1a150cb39bbb1aa80da8c57ea100",
			pubkey: "03cee91b6d329e00c344ad5d67cfd00d885ec36e8975b5d9097738939cb8c08b31",
		},
		{
			msg:          "176b81623cf98f45879f3a48fa34af77dde44b2ffa0ddd2bf9edb386f76ec0ef",
			sig:          "d201ec2b29ce3cf3e6048296188adff4b5dfcb337c1d1157f28654e445bb940b4e47d6b0c7ba43d072bf8618775f123a435e8d1a150cb39bbb1aa80da8c57ea100",
			recoverFails: true,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s-%s", tc.msg, tc.sig), func(t *testing.T) {
			var sig Signature
			var pubkey XY
			var msg Number

			sigBytes, err := hex.DecodeString(tc.sig)
			if err != nil {
				t.Error("invalid sig hex")
			}
			recid := int(sigBytes[64])
			sig.R.SetBytes(sigBytes[:32])
			sig.S.SetBytes(sigBytes[32:64])
			msg.SetHex(tc.msg)

			if sig.Recover(&pubkey, &msg, recid) {
				if tc.recoverFails {
					t.Error("sig.Recover expected to fail")
				}

				pubkeyHex := hex.EncodeToString(pubkey.Bytes())
				if tc.pubkey != pubkeyHex {
					t.Errorf("pubkey does not match %s != %s", tc.pubkey, pubkeyHex)
				}
			} else {
				if !tc.recoverFails {
					t.Error("sig.Recover failed")
				}
			}
		})
	}
}

func TestSigVerify(t *testing.T) {
	var msg Number
	var sig Signature
	var key XY

	//// len(65) keys are rejected now, this test case is invalid:
	// msg.SetHex("3382219555ddbb5b00e0090f469e590ba1eae03c7f28ab937de330aa60294ed6")
	// sig.R.SetHex("fe00e013c244062847045ae7eb73b03fca583e9aa5dbd030a8fd1c6dfcf11b10")
	// sig.S.SetHex("7d0d04fed8fa1e93007468d5a9e134b0a7023b6d31db4e50942d43a250f4d07c")
	// xy, _ := hex.DecodeString("040eaebcd1df2df853d66ce0e1b0fda07f67d1cabefde98514aad795b86a6ea66dbeb26b67d7a00e2447baeccc8a4cef7cd3cad67376ac1c5785aeebb4f6441c16")
	// key.ParsePubkey(xy)
	// if !sig.Verify(&key, &msg) {
	// 	t.Error("sig.Verify 0")
	// }

	msg.SetHex("D474CBF2203C1A55A411EEC4404AF2AFB2FE942C434B23EFE46E9F04DA8433CA")
	sig.R.SetHex("98F9D784BA6C5C77BB7323D044C0FC9F2B27BAA0A5B0718FE88596CC56681980")
	sig.S.SetHex("E3599D551029336A745B9FB01566624D870780F363356CEE1425ED67D1294480")
	key.X.SetHex("7d709f85a331813f9ae6046c56b3a42737abf4eb918b2e7afee285070e968b93")
	key.Y.SetHex("26150d1a63b342986c373977b00131950cb5fc194643cad6ea36b5157eba4602")
	if !sig.Verify(&key, &msg) {
		t.Error("sig.Verify 1")
	}

	msg.SetHex("2c43a883f4edc2b66c67a7a355b9312a565bb3d33bb854af36a06669e2028377")
	sig.R.SetHex("6b2fa9344462c958d4a674c2a42fbedf7d6159a5276eb658887e2e1b3915329b")
	sig.S.SetHex("eddc6ea7f190c14a0aa74e41519d88d2681314f011d253665f301425caf86b86")
	xy, err := hex.DecodeString("02a60d70cfba37177d8239d018185d864b2bdd0caf5e175fd4454cc006fd2d75ac")
	if err != nil {
		t.Fail()
	}

	if err := key.ParsePubkey(xy); err != nil {
		t.Errorf("ParsePubkey failed: %v", err)
	}
	if !sig.Verify(&key, &msg) {
		t.Error("sig.Verify 2")
	}
}

func TestSigSign(t *testing.T) {
	var sec, msg, non Number
	var sig Signature
	var recid int
	sec.SetHex("73641C99F7719F57D8F4BEB11A303AFCD190243A51CED8782CA6D3DBE014D146")
	msg.SetHex("D474CBF2203C1A55A411EEC4404AF2AFB2FE942C434B23EFE46E9F04DA8433CA")
	non.SetHex("9E3CD9AB0F32911BFDE39AD155F527192CE5ED1F51447D63C4F154C118DA598E")
	res := sig.Sign(&sec, &msg, &non, &recid)
	if res != 1 {
		t.Error("res failed", res)
	}
	if recid != 0 {
		t.Error("recid failed", recid)
	}
	non.SetHex("98f9d784ba6c5c77bb7323d044c0fc9f2b27baa0a5b0718fe88596cc56681980")
	if sig.R.Cmp(&non.Int) != 0 {
		t.Error("R failed", sig.R.String())
	}
	non.SetHex("1ca662aaefd6cc958ba4604fea999db133a75bf34c13334dabac7124ff0cfcc1")
	if sig.S.Cmp(&non.Int) != 0 {
		t.Error("S failed", sig.S.String())
	}
	expectSig := "98f9d784ba6c5c77bb7323d044c0fc9f2b27baa0a5b0718fe88596cc566819801ca662aaefd6cc958ba4604fea999db133a75bf34c13334dabac7124ff0cfcc1"
	if expectSig != hex.EncodeToString(sig.Bytes()) {
		t.Error("signature doesnt match")
	}
}

func BenchmarkVerify(b *testing.B) {
	var msg Number
	var sig Signature
	var key XY
	msg.SetHex("D474CBF2203C1A55A411EEC4404AF2AFB2FE942C434B23EFE46E9F04DA8433CA")
	sig.R.SetHex("98F9D784BA6C5C77BB7323D044C0FC9F2B27BAA0A5B0718FE88596CC56681980")
	sig.S.SetHex("E3599D551029336A745B9FB01566624D870780F363356CEE1425ED67D1294480")
	key.X.SetHex("7d709f85a331813f9ae6046c56b3a42737abf4eb918b2e7afee285070e968b93")
	key.Y.SetHex("26150d1a63b342986c373977b00131950cb5fc194643cad6ea36b5157eba4602")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !sig.Verify(&key, &msg) {
			b.Fatal("sig_verify failed")
		}
	}
}

func BenchmarkSign(b *testing.B) {
	var sec, msg, non Number
	var sig Signature
	var recid int
	sec.SetHex("73641C99F7719F57D8F4BEB11A303AFCD190243A51CED8782CA6D3DBE014D146")
	msg.SetHex("D474CBF2203C1A55A411EEC4404AF2AFB2FE942C434B23EFE46E9F04DA8433CA")
	non.SetHex("9E3CD9AB0F32911BFDE39AD155F527192CE5ED1F51447D63C4F154C118DA598E")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sig.Sign(&sec, &msg, &non, &recid)
	}
}
