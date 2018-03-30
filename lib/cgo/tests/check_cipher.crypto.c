
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"

#define ALPHANUM "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
#define ALPHANUM_LEN 62

void randBytes(GoSlice *bytes, size_t n) {
  size_t i = 0;
  unsigned char *ptr = (unsigned char *) bytes->data;
  for (; i < n; ++i, ++ptr) {
    *ptr = ALPHANUM[rand() % ALPHANUM_LEN];
  } 
  bytes->len = (GoInt) n;
}

void setup(void) {
  srand ((unsigned int) time (NULL));
}

Test(asserts, TestNewPubKey) {
  unsigned char buff[101];
  GoSlice slice;
  PubKey pk;

  slice.data = buff;
  slice.cap = 101;

  randBytes(&slice, 31);
  slice.len = 31;
  unsigned int errcode = SKY_cipher_NewPubKey(slice, &pk);
  cr_assert(errcode == SKY_ERROR, "31 random bytes");

  randBytes(&slice, 32);
  errcode = SKY_cipher_NewPubKey(slice, &pk);
  cr_assert(errcode == SKY_ERROR, "32 random bytes");

  randBytes(&slice, 34);
  errcode = SKY_cipher_NewPubKey(slice, &pk);
  cr_assert(errcode == SKY_ERROR, "34 random bytes");

  slice.len = 0;
  errcode = SKY_cipher_NewPubKey(slice, &pk);
  cr_assert(errcode == SKY_ERROR, "0 random bytes");

  randBytes(&slice, 100);
  errcode = SKY_cipher_NewPubKey(slice, &pk);
  cr_assert(errcode == SKY_ERROR, "100 random bytes");

  randBytes(&slice, 33);
  errcode = SKY_cipher_NewPubKey(slice, &pk);
  cr_assert(errcode == SKY_OK, "33 random bytes");

  cr_assert(eq(u8[33], pk, buff));
}

#define SIZE_ALL -1

// TODO: Move to libsky_string.c
void strnhex(unsigned char* buf, char *str, int n){
    unsigned char * pin = buf;
    const char * hex = "0123456789ABCDEF";
    char * pout = str;
    for(; *pin && n; --n){
        *pout++ = hex[(*pin>>4)&0xF];
        *pout++ = hex[(*pin++)&0xF];
    }
    *pout = 0;
}

// TODO: Move to libsky_string.c
void strhex(unsigned char* buf, char *str){
  strnhex(buf, str, SIZE_ALL);
}

Test(asserts, TestPubKeyFromHex) {
  PubKey p, p1;
  GoString s;
  unsigned char buff[50];
  char sbuff[101];
  GoSlice slice;
  unsigned int errcode;

  slice.data = (void *)buff;
  slice.cap = 51;
  slice.len = 0;

	// Invalid hex
  s.n = 0;
  errcode = SKY_cipher_PubKeyFromHex(s, &p1);
  cr_assert(errcode == SKY_ERROR, "TestPubKeyFromHex: Invalid hex. Empty string");

  s.p = "cascs";
  s.n = strlen(s.p);
  errcode = SKY_cipher_PubKeyFromHex(s, &p1);
  cr_assert(errcode == SKY_ERROR, "TestPubKeyFromHex: Invalid hex. Bad chars");

	// Invalid hex length
  randBytes(&slice, 33);
  errcode = SKY_cipher_NewPubKey(slice, &p);
  cr_assert(errcode == SKY_OK);
  strnhex(&p[0], sbuff, slice.len / 2);
  s.p = sbuff;
  s.n = strlen(s.p);
  errcode = SKY_cipher_PubKeyFromHex(s, &p1);
  cr_assert(errcode == SKY_ERROR, "TestPubKeyFromHex: Invalid hex length");

	// Valid
  strhex(&p[0], sbuff);
  s.p = sbuff;
  s.n = strlen(s.p);
  errcode = SKY_cipher_PubKeyFromHex(s, &p1);
  cr_assert(errcode == SKY_OK, "TestPubKeyFromHex: Valid. No panic.");
  cr_assert(eq(u8[33], p, p1));
}

Test(asserts, TestPubKeyHex) {
  PubKey p, p2;
  GoString s, s2;
  unsigned char buff[50];
  GoSlice slice;
  unsigned int errcode;

  slice.data = buff;
  slice.len = 0;
  slice.cap = 50;

  randBytes(&slice, 33);
  errcode = SKY_cipher_NewPubKey(slice, &p);
  cr_assert(errcode == SKY_OK);
  s.p = SKY_cipher_PubKey_Hex(&p);
  s.n = strlen(s.p);
  errcode = SKY_cipher_PubKeyFromHex(s, &p2);
  cr_assert(errcode == SKY_OK);
  cr_assert(eq(u8[33], p, p2));

  s2.p = SKY_cipher_PubKey_Hex(&p2);
  s2.n = strlen(s2.p);
  // TODO: Write like this cr_assert(eq(type(struct GoString), &s, &s2))
  cr_assert(eq(int, s.n, s2.n));
  cr_assert(eq(str, (char *) s.p, (char *) s2.p));
  if (s.p != NULL) {
    free((void *) s.p);
  }
  if (s2.p != NULL) {
    free((void *) s2.p);
  }
}

Test(asserts, TestPubKeyVerify) {
  PubKey p;
  unsigned char buff[50];
  GoSlice slice;
  unsigned int errcode;

  slice.data = buff;
  slice.len = 0;
  slice.cap = 50;

  int i = 0;
  for (; i < 10; i++) {
    randBytes(&slice, 33);
    errcode = SKY_cipher_NewPubKey(slice, &p);
    cr_assert(errcode == SKY_OK);
    errcode = SKY_cipher_PubKey_Verify(&p);
    cr_assert(errcode == SKY_ERROR);
  } 
}

Test(asserts, TestPubKeyVerifyNil) {
  PubKey p = {
    0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
    0, 0, 0
  };
  unsigned int errcode;

  errcode = SKY_cipher_PubKey_Verify(&p);
  cr_assert(errcode == SKY_ERROR);
}

Test(asserts, TestPubKeyVerifyDefault1) {
  PubKey p;
  SecKey s;

  fprintf(stderr, "p1 %p %p\n", &p, &s);
  SKY_cipher_GenerateKeyPair(&p, &s);
  fprintf(stderr, "p2\n");
  unsigned int errcode = SKY_cipher_PubKey_Verify(&p);
  fprintf(stderr, "p3\n");
  cr_assert(errcode == SKY_OK);
}

Test(asserts, TestPubKeyVerifyDefault2) {
  PubKey p;
  SecKey s;
  int i;

  for (i = 0; i < 1024; ++i) {
    fprintf(stderr, "p1 %p %p\n", &p, &s);
    SKY_cipher_GenerateKeyPair(&p, &s);
    fprintf(stderr, "p2\n");
    unsigned int errcode = SKY_cipher_PubKey_Verify(&p);
    fprintf(stderr, "p3\n");
    cr_assert(errcode == SKY_OK);
  }
}

Test(asserts, TestPubKeyToAddressHash) {
  PubKey p;
  SecKey s;
  Ripemd160 h;

  SKY_cipher_GenerateKeyPair(&p, &s);
  SKY_cipher_PubKey_ToAddressHash(&p, &h);
  /* TODO: Translate code snippet
   *
   * x := sha256.Sum256(p[:])
   * x = sha256.Sum256(x[:])
   * rh := ripemd160.New()
   * rh.Write(x[:])
   * y := rh.Sum(nil)
   * assert.True(t, bytes.Equal(h[:], y))
   *
   */
}

Test(asserts, TestPubKeyToAddress) {
  PubKey p;
  SecKey s;
  Address addr;
  Ripemd160 h;
  int errcode;

  SKY_cipher_GenerateKeyPair(&p, &s);
  SKY_cipher_AddressFromPubKey(&p, &addr);
  errcode = SKY_cipher_Address_Verify(&addr, &p);
  cr_assert(errcode == SKY_OK);
}

Test(asserts, TestPubKeyToAddress2) {
  PubKey p;
  SecKey s;
  Address addr;
  GoString_ addrStr;
  int i, errcode;

	for (i = 0; i < 1024; i++) {
		SKY_cipher_GenerateKeyPair(&p, &s);
		SKY_cipher_AddressFromPubKey(&p, &addr);
		//func (self Address) Verify(key PubKey) error {
		errcode = SKY_cipher_Address_Verify(&addr, &p);
    cr_assert(errcode == SKY_OK);
		SKY_cipher_Address_String(&addr, &addrStr);
		errcode = SKY_cipher_DecodeBase58Address(
        *((GoString*)(GoString_*)&addrStr), &addr);
		//func DecodeBase58Address(addr string) (Address, error) {
    cr_assert(errcode == SKY_OK);
	}
}

Test(asserts, TestMustNewSecKey) {
  unsigned char buff[101];
  GoSlice b;
  SecKey sk;
  int errcode;

  b.data = buff;
  b.cap = 101;

  randBytes(&b, 31);
  errcode = SKY_cipher_NewSecKey(b, &sk);
  cr_assert(errcode == SKY_ERROR);

  randBytes(&b, 33);
  errcode = SKY_cipher_NewSecKey(b, &sk);
  cr_assert(errcode == SKY_ERROR);

  randBytes(&b, 34);
  errcode = SKY_cipher_NewSecKey(b, &sk);
  cr_assert(errcode == SKY_ERROR);

  randBytes(&b, 0);
  errcode = SKY_cipher_NewSecKey(b, &sk);
  cr_assert(errcode == SKY_ERROR);

  randBytes(&b, 100);
  errcode = SKY_cipher_NewSecKey(b, &sk);
  cr_assert(errcode == SKY_ERROR);

  randBytes(&b, 32);
  errcode = SKY_cipher_NewSecKey(b, &sk);
  cr_assert(errcode == SKY_OK);
  cr_assert(eq(u8[32], sk, buff));
}

Test(asserts, TestMustSecKeyFromHex) {
  GoString str;
  SecKey sk, sk1;
  unsigned int buff[50];
  GoSlice b;
  char strBuff[101];
  GoString s;
	int errcode;

  // Invalid hex
  s.p = "";
  s.n = strlen(s.p);
  errcode = SKY_cipher_SecKeyFromHex(s, &sk);
  cr_assert(errcode == SKY_ERROR);

  s.p = "cascs";
  s.n = strlen(s.p);
  errcode = SKY_cipher_SecKeyFromHex(s, &sk);
  cr_assert(errcode == SKY_ERROR);

	// Invalid hex length
  b.data = buff;
  b.cap = 50;
  randBytes(&b, 32);
  SKY_cipher_NewSecKey(b, &sk);
  strnhex(sk, strBuff, 32);
  s.p = strBuff;
  s.n = strlen(strBuff) >> 1;
  errcode = SKY_cipher_SecKeyFromHex(s, &sk1);
  cr_assert(errcode == SKY_ERROR);

	// Valid
  s.p = strBuff;
  s.n = strlen(strBuff);
  errcode = SKY_cipher_SecKeyFromHex(s, &sk1);
  cr_assert(errcode == SKY_OK);
  cr_assert(eq(u8[32], sk, sk1));
}

Test(asserts, TestSecKeyHex) {
  SecKey sk, sk2;
  unsigned char buff[101];
  char strBuff[50];
  GoSlice b;
  GoString str, h;
  int errcode;

  b.data = buff;
  b.cap = 50;
  h.p = strBuff;
  h.n = 0;

  randBytes(&b, 32);
  SKY_cipher_NewSecKey(b, &sk);
  SKY_cipher_SecKey_Hex(&sk, (GoString_ *)&str);

  // Copy early to ensure memory is released
  strncpy((char *) h.p, str.p, str.n);
  h.n = str.n;
  free((void *) str.p);

  errcode = SKY_cipher_SecKeyFromHex(h, &sk2);
  cr_assert(errcode == SKY_OK);
  cr_assert(eq(u8[32], sk, sk2));
}

Test(asserts, TestSecKeyVerify) {
  SecKey sk;
  PubKey pk;
  int errcode;

	// Empty secret key should not be valid
  memset(&sk, 0, 32);
  errcode = SKY_cipher_SecKey_Verify(&sk);
  cr_assert(errcode == SKY_OK);

	// Generated sec key should be valid
  SKY_cipher_GenerateKeyPair(&pk, &sk);
  errcode = SKY_cipher_SecKey_Verify(&sk);
  cr_assert(errcode == SKY_OK);

	// Random bytes are usually valid
}

Test(asserts, TestECDHonce) {
  PubKey pub1, pub2;
  SecKey sec1, sec2;
  unsigned char buff1[50], buff2[50];
  GoSlice_ buf1, buf2;

  buf1.data = buff1;
  buf1.len = 0;
  buf1.cap = 50;
  buf2.data = buff2;
  buf2.len = 0;
  buf2.cap = 50;

  SKY_cipher_GenerateKeyPair(&pub1, &sec1);
  SKY_cipher_GenerateKeyPair(&pub2, &sec2);

  SKY_cipher_ECDH(&pub2, &sec1, &buf1);
  SKY_cipher_ECDH(&pub1, &sec2, &buf2);

  // ECDH shared secrets are 32 bytes SHA256 hashes in the end
  cr_assert(eq(u8[32], buff1, buff2));
}

Test(asserts, TestECDHloop) {
  int i;
  PubKey pub1, pub2;
  SecKey sec1, sec2;
  unsigned char buff1[50], buff2[50];
  GoSlice_ buf1, buf2;

  buf1.data = buff1;
  buf1.len = 0;
  buf1.cap = 50;
  buf2.data = buff2;
  buf2.len = 0;
  buf2.cap = 50;

	for (i = 0; i < 128; i++) {
    SKY_cipher_GenerateKeyPair(&pub1, &sec1);
    SKY_cipher_GenerateKeyPair(&pub2, &sec2);
		SKY_cipher_ECDH(&pub2, &sec1, &buf1);
		SKY_cipher_ECDH(&pub1, &sec2, &buf2);
    cr_assert(eq(u8[32], buff1, buff2));
	}
}

Test(asserts, TestNewSig) {
  unsigned char buff[101];
  GoSlice b;
  Sig s;
  int errcode;

  b.data = buff;
  b.len = 0;
  b.cap = 101;

  randBytes(&b, 64);
  SKY_cipher_NewSig(b, &s);
  cr_assert(errcode == SKY_ERROR);

  randBytes(&b, 66);
  errcode = SKY_cipher_NewSig(b, &s);
  cr_assert(errcode == SKY_ERROR);

  randBytes(&b, 67);
  errcode = SKY_cipher_NewSig(b, &s);
  cr_assert(errcode == SKY_ERROR);

  randBytes(&b, 0);
  errcode = SKY_cipher_NewSig(b, &s);
  cr_assert(errcode == SKY_ERROR);

  randBytes(&b, 100);
  errcode = SKY_cipher_NewSig(b, &s);
  cr_assert(errcode == SKY_ERROR);

  randBytes(&b, 65);
  errcode = SKY_cipher_NewSig(b, &s);
  cr_assert(errcode == SKY_OK);
  cr_assert(eq(u8[65], buff, s));
}

Test(asserts, TestMustSigFromHex) {
  unsigned char buff[101];
  char strBuff[101];
  GoSlice b;
  GoString str;
  Sig s;
  int errcode;

  b.data = buff;
  b.len = 0;
  b.cap = 101;

	// Invalid hex
  str.p = "";
  str.n = strlen(str.p);
  errcode = SKY_cipher_SigFromHex(str, &s);
  cr_assert(errcode == SKY_ERROR);

  str.p = "cascs";
  str.n = strlen(str.p);
  errcode = SKY_cipher_SigFromHex(str, &s);
  cr_assert(errcode == SKY_ERROR);

	// Invalid hex length
  randBytes(&b, 65);
  errcode = SKY_cipher_NewSig(b, &s);
  strnhex(buff, (char *) str.p, b.len >> 1);
  str.n = strlen(str.p);

  errcode = SKY_cipher_SigFromHex(str, &s);
  cr_assert(errcode == SKY_ERROR);

	// Valid
  strnhex(buff, (char *)str.p, b.len);
  str.n = strlen(str.p);
  errcode = SKY_cipher_SigFromHex(str, &s);
  cr_assert(errcode == SKY_OK);
  cr_assert(eq(u8[65], buff, s));
}

Test(asserts, TestSigHex) {
  cr_fatal("Not implemented");
  /*
	b := randBytes(t, 65)
	p := NewSig(b)
	h := p.Hex()
	p2 := MustSigFromHex(h)
	assert.Equal(t, p2, p)
	assert.Equal(t, p2.Hex(), h)
  */
}

Test(asserts, TestChkSig) {
  cr_fatal("Not implemented");
  /*
	p, s := GenerateKeyPair()
	assert.Nil(t, p.Verify())
	assert.Nil(t, s.Verify())
	a := AddressFromPubKey(p)
	assert.Nil(t, a.Verify(p))
	b := randBytes(t, 256)
	h := SumSHA256(b)
	sig := SignHash(h, s)
	assert.Nil(t, ChkSig(a, h, sig))
	// Empty sig should be invalid
	assert.NotNil(t, ChkSig(a, h, Sig{}))
	// Random sigs should not pass
	for i := 0; i < 100; i++ {
		assert.NotNil(t, ChkSig(a, h, NewSig(randBytes(t, 65))))
	}
	// Sig for one hash does not work for another hash
	h2 := SumSHA256(randBytes(t, 256))
	sig2 := SignHash(h2, s)
	assert.Nil(t, ChkSig(a, h2, sig2))
	assert.NotNil(t, ChkSig(a, h, sig2))
	assert.NotNil(t, ChkSig(a, h2, sig))

	// Different secret keys should not create same sig
	p2, s2 := GenerateKeyPair()
	a2 := AddressFromPubKey(p2)
	h = SHA256{}
	sig = SignHash(h, s)
	sig2 = SignHash(h, s2)
	assert.Nil(t, ChkSig(a, h, sig))
	assert.Nil(t, ChkSig(a2, h, sig2))
	assert.NotEqual(t, sig, sig2)
	h = SumSHA256(randBytes(t, 256))
	sig = SignHash(h, s)
	sig2 = SignHash(h, s2)
	assert.Nil(t, ChkSig(a, h, sig))
	assert.Nil(t, ChkSig(a2, h, sig2))
	assert.NotEqual(t, sig, sig2)

	// Bad address should be invalid
	assert.NotNil(t, ChkSig(a, h, sig2))
	assert.NotNil(t, ChkSig(a2, h, sig))
  */
}

Test(asserts, TestSignHash) {
  cr_fatal("Not implemented");
  /*
	p, s := GenerateKeyPair()
	a := AddressFromPubKey(p)
	h := SumSHA256(randBytes(t, 256))
	sig := SignHash(h, s)
	assert.NotEqual(t, sig, Sig{})
	assert.Nil(t, ChkSig(a, h, sig))
  */
}

Test(asserts, TestPubKeyFromSecKey) {
  cr_fatal("Not implemented");
  /*
	p, s := GenerateKeyPair()
	assert.Equal(t, PubKeyFromSecKey(s), p)
	assert.Panics(t, func() { PubKeyFromSecKey(SecKey{}) })
	assert.Panics(t, func() { PubKeyFromSecKey(NewSecKey(randBytes(t, 99))) })
	assert.Panics(t, func() { PubKeyFromSecKey(NewSecKey(randBytes(t, 31))) })
  */
}

Test(asserts, TestPubKeyFromSig) {
  cr_fatal("Not implemented");
  /*
	p, s := GenerateKeyPair()
	h := SumSHA256(randBytes(t, 256))
	sig := SignHash(h, s)
	p2, err := PubKeyFromSig(sig, h)
	assert.Equal(t, p, p2)
	assert.Nil(t, err)
	_, err = PubKeyFromSig(Sig{}, h)
	assert.NotNil(t, err)
  */
}

Test(asserts, TestVerifySignature) {
  cr_fatal("Not implemented");
  /*
	p, s := GenerateKeyPair()
	h := SumSHA256(randBytes(t, 256))
	h2 := SumSHA256(randBytes(t, 256))
	sig := SignHash(h, s)
	assert.Nil(t, VerifySignature(p, sig, h))
	assert.NotNil(t, VerifySignature(p, Sig{}, h))
	assert.NotNil(t, VerifySignature(p, sig, h2))
	p2, _ := GenerateKeyPair()
	assert.NotNil(t, VerifySignature(p2, sig, h))
	assert.NotNil(t, VerifySignature(PubKey{}, sig, h))
  */
}

Test(asserts, TestGenerateKeyPair) {
  cr_fatal("Not implemented");
  /*
	p, s := GenerateKeyPair()
	assert.Nil(t, p.Verify())
	assert.Nil(t, s.Verify())
  */
}

Test(asserts, TestGenerateDeterministicKeyPair) {
  cr_fatal("Not implemented");
  /*
	// TODO -- deterministic key pairs are useless as is because we can't
	// generate pair n+1, only pair 0
	seed := randBytes(t, 32)
	p, s := GenerateDeterministicKeyPair(seed)
	assert.Nil(t, p.Verify())
	assert.Nil(t, s.Verify())
	p, s = GenerateDeterministicKeyPair(seed)
	assert.Nil(t, p.Verify())
	assert.Nil(t, s.Verify())
  */
}

Test(asserts, TestSecKeTest) {
  cr_fatal("Not implemented");
  /*
	_, s := GenerateKeyPair()
	assert.Nil(t, TestSecKey(s))
	assert.NotNil(t, TestSecKey(SecKey{}))
  */
}

Test(asserts, TestSecKeyHashTest) {
  cr_fatal("Not implemented");
  /*
	_, s := GenerateKeyPair()
	h := SumSHA256(randBytes(t, 256))
	assert.Nil(t, TestSecKeyHash(s, h))
	assert.NotNil(t, TestSecKeyHash(SecKey{}, h))
  */
}
