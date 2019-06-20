package secp256k1go

import (
	"bytes"
	"log"
)

// Signature represents the signature
type Signature struct {
	R, S Number
}

// Verify verify the signature
func (sig *Signature) Verify(pubkey *XY, message *Number) bool {
	var r2 Number
	return sig.recompute(&r2, pubkey, message) && sig.R.Cmp(&r2.Int) == 0
}

func (sig *Signature) recompute(r2 *Number, pubkey *XY, message *Number) bool {
	var sn, u1, u2 Number

	sn.modInv(&sig.S, &TheCurve.Order)
	u1.modMul(&sn, message, &TheCurve.Order)
	u2.modMul(&sn, &sig.R, &TheCurve.Order)

	var pr, pubkeyj XYZ
	pubkeyj.SetXY(pubkey)

	pubkeyj.ECmult(&pr, &u2, &u1)
	if pr.IsInfinity() {
		return false
	}

	var xr Field
	pr.getX(&xr)
	xr.Normalize()
	var xrb [32]byte
	xr.GetB32(xrb[:])
	r2.SetBytes(xrb[:])
	r2.Mod(&r2.Int, &TheCurve.Order.Int)

	return true
}

/*
Reference code for Signature.Recover

https://github.com/bitcoin-core/secp256k1/blob/e541a90ef6461007d9c6a74b9f9a7fb8aa34aaa8/src/modules/recovery/main_impl.h

static int secp256k1_ecdsa_sig_recover(const secp256k1_ecmult_context *ctx, const secp256k1_scalar *sigr, const secp256k1_scalar* sigs, secp256k1_ge *pubkey, const secp256k1_scalar *message, int recid) {
    unsigned char brx[32];
    secp256k1_fe fx;
    secp256k1_ge x;
    secp256k1_gej xj;
    secp256k1_scalar rn, u1, u2;
    secp256k1_gej qj;
    int r;

    if (secp256k1_scalar_is_zero(sigr) || secp256k1_scalar_is_zero(sigs)) {
        return 0;
    }

    secp256k1_scalar_get_b32(brx, sigr);
    r = secp256k1_fe_set_b32(&fx, brx);
    (void)r;
    VERIFY_CHECK(r); // brx comes from a scalar, so is less than the order; certainly less than p
    if (recid & 2) {
        if (secp256k1_fe_cmp_var(&fx, &secp256k1_ecdsa_const_p_minus_order) >= 0) {
            return 0;
        }
        secp256k1_fe_add(&fx, &secp256k1_ecdsa_const_order_as_fe);
    }
    if (!secp256k1_ge_set_xo_var(&x, &fx, recid & 1)) {
        return 0;
    }
    secp256k1_gej_set_ge(&xj, &x);
    secp256k1_scalar_inverse_var(&rn, sigr);
    secp256k1_scalar_mul(&u1, &rn, message);
    secp256k1_scalar_negate(&u1, &u1);
    secp256k1_scalar_mul(&u2, &rn, sigs);
    secp256k1_ecmult(ctx, &qj, &xj, &u2, &u1);
    secp256k1_ge_set_gej_var(pubkey, &qj);
    return !secp256k1_gej_is_infinity(&qj);
}
*/

// Recover recovers a pubkey XY point given the message that was signed to create
// this signature.
func (sig *Signature) Recover(pubkey *XY, msg *Number, recid int) bool {
	var rx, rn, u1, u2 Number
	var fx Field
	var x XY
	var xj, qj XYZ

	if sig.R.Sign() == 0 || sig.S.Sign() == 0 {
		return false
	}

	rx.Set(&sig.R.Int)
	if (recid & 2) != 0 {
		rx.Add(&rx.Int, &TheCurve.Order.Int)
		if rx.Cmp(&TheCurve.p.Int) >= 0 {
			return false
		}
	}

	fx.SetB32(LeftPadBytes(rx.Bytes(), 32))

	x.SetXO(&fx, (recid&1) != 0)
	if !x.IsValid() {
		return false
	}

	xj.SetXY(&x)
	rn.modInv(&sig.R, &TheCurve.Order)
	u1.modMul(&rn, msg, &TheCurve.Order)
	u1.Sub(&TheCurve.Order.Int, &u1.Int)
	u2.modMul(&rn, &sig.S, &TheCurve.Order)
	xj.ECmult(&qj, &u2, &u1)
	pubkey.SetXYZ(&qj)
	return !qj.IsInfinity()
}

/*
Reference code for Signature.Sign

https://github.com/bitcoin-core/secp256k1/blob/master/src/ecdsa_impl.h

static int secp256k1_ecdsa_sig_sign(const secp256k1_ecmult_gen_context *ctx, secp256k1_scalar *sigr, secp256k1_scalar *sigs, const secp256k1_scalar *seckey, const secp256k1_scalar *message, const secp256k1_scalar *nonce, int *recid) {
    unsigned char b[32];
    secp256k1_gej rp;
    secp256k1_ge r;
    secp256k1_scalar n;
    int overflow = 0;

    secp256k1_ecmult_gen(ctx, &rp, nonce);
    secp256k1_ge_set_gej(&r, &rp);
    secp256k1_fe_normalize(&r.x);
    secp256k1_fe_normalize(&r.y);
    secp256k1_fe_get_b32(b, &r.x);
    secp256k1_scalar_set_b32(sigr, b, &overflow);
    // These two conditions should be checked before calling
    VERIFY_CHECK(!secp256k1_scalar_is_zero(sigr));
    VERIFY_CHECK(overflow == 0);

    if (recid) {
        // The overflow condition is cryptographically unreachable as hitting it requires finding the discrete log
        // of some P where P.x >= order, and only 1 in about 2^127 points meet this criteria.
        *recid = (overflow ? 2 : 0) | (secp256k1_fe_is_odd(&r.y) ? 1 : 0);
    }
    secp256k1_scalar_mul(&n, sigr, seckey);
    secp256k1_scalar_add(&n, &n, message);
    secp256k1_scalar_inverse(sigs, nonce);
    secp256k1_scalar_mul(sigs, sigs, &n);
    secp256k1_scalar_clear(&n);
    secp256k1_gej_clear(&rp);
    secp256k1_ge_clear(&r);
    if (secp256k1_scalar_is_zero(sigs)) {
        return 0;
    }
    if (secp256k1_scalar_is_high(sigs)) {
        secp256k1_scalar_negate(sigs, sigs);
        if (recid) {
            *recid ^= 1;
        }
    }
    return 1;
}
*/

// Sign signs the signature. Returns 1 on success, 0 on failure
func (sig *Signature) Sign(seckey, message, nonce *Number, recid *int) int {
	var r XY
	var n Number
	var b [32]byte

	// r = nonce*G
	rp := ECmultGen(*nonce)
	r.SetXYZ(&rp)
	r.X.Normalize()
	r.Y.Normalize()
	r.X.GetB32(b[:])
	sig.R.SetBytes(b[:])

	if sig.R.Sign() == 0 {
		log.Panic("sig R value should not be 0")
	}

	if recid != nil {
		*recid = 0
		// The overflow condition is cryptographically unreachable as hitting
		// it requires finding the discrete log of some P where P.x >= order,
		// and only 1 in about 2^127 points meet this criteria.
		if sig.R.Cmp(&TheCurve.Order.Int) >= 0 {
			*recid |= 2
		}
		if r.Y.IsOdd() {
			*recid |= 1
		}
	}

	sig.R.mod(&TheCurve.Order)
	n.modMul(&sig.R, seckey, &TheCurve.Order)
	n.Add(&n.Int, &message.Int)
	n.mod(&TheCurve.Order)
	sig.S.modInv(nonce, &TheCurve.Order)
	sig.S.modMul(&sig.S, &n, &TheCurve.Order)

	if sig.S.Sign() == 0 {
		return 0
	}

	// Break signature malleability
	if sig.S.Cmp(&TheCurve.halfOrder.Int) == 1 {
		sig.S.Sub(&TheCurve.Order.Int, &sig.S.Int)
		if recid != nil {
			*recid ^= 1
		}
	}

	return 1
}

// ParseBytes parses a serialized R||S pair to Signature.
// R and S should be in big-endian encoding.
func (sig *Signature) ParseBytes(v []byte) {
	if len(v) != 64 {
		log.Panic("Signature.ParseBytes requires 64 bytes")
	}
	sig.R.SetBytes(v[0:32])
	sig.S.SetBytes(v[32:64])
}

//secp256k1_num_get_bin(sig64, 32, &sig.r);
//secp256k1_num_get_bin(sig64 + 32, 32, &sig.s);

// Bytes serializes compressed signatures as bytes.
// The serialization format is R||S. R and S bytes are big-endian.
// R and S are left-padded with the NUL byte to ensure their length is 32 bytes.
func (sig *Signature) Bytes() []byte {
	r := sig.R.Bytes() // big-endian
	s := sig.S.Bytes() // big-endian

	for len(r) < 32 {
		r = append([]byte{0}, r...)
	}
	for len(s) < 32 {
		s = append([]byte{0}, s...)
	}

	if len(r) != 32 || len(s) != 32 {
		log.Panicf("signature size invalid: %d, %d", len(r), len(s))
	}

	res := new(bytes.Buffer)
	if _, err := res.Write(r); err != nil {
		log.Panic(err)
	}
	if _, err := res.Write(s); err != nil {
		log.Panic(err)
	}

	//test
	if true {
		ret := res.Bytes()
		var sig2 Signature
		sig2.ParseBytes(ret)
		if !bytes.Equal(sig.R.Bytes(), sig2.R.Bytes()) {
			log.Panic("serialization failed 1")
		}
		if !bytes.Equal(sig.S.Bytes(), sig2.S.Bytes()) {
			log.Panic("serialization failed 2")
		}
	}

	if len(res.Bytes()) != 64 {
		log.Panic("Signature.Bytes result bytes must be 64 bytes long")
	}
	return res.Bytes()
}
