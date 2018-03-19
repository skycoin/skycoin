
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
    int i = 0;
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
