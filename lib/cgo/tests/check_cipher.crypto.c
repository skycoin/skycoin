
#include <stdlib.h>
#include <time.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"

#define ALPHANUM "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
#define ALPHANUM_LEN 62

void randBytes(unsigned char* bytes, size_t n) {
  size_t i;
  unsigned char *ptr;
  for (i = 0, ptr = bytes; i < n; ++i, ++ptr) {
    *ptr = ALPHANUM[rand() % ALPHANUM_LEN];
  } 
}

void setup(void) {
  srand ((unsigned int) time (NULL));
}

Test(asserts, TestNewPubKey) {
  // buffer big enough to hold all kind of data needed by test cases
  unsigned char buff[101];
  GoSlice slice;
  PubKey pk;

  slice.data = buff;
  slice.cap = 101;

  randBytes(buff, 31);
  slice.len = 31;
  unsigned int errcode = SKY_cipher_NewPubKey(slice, &pk);
  cr_assert(errcode == SKY_ERROR, "31 random bytes");

  randBytes(buff, 32);
  slice.len = 32;
  errcode = SKY_cipher_NewPubKey(slice, &pk);
  cr_assert(errcode == SKY_ERROR, "32 random bytes");

  randBytes(buff, 34);
  slice.len = 34;
  errcode = SKY_cipher_NewPubKey(slice, &pk);
  cr_assert(errcode == SKY_ERROR, "34 random bytes");

  slice.len = 0;
  errcode = SKY_cipher_NewPubKey(slice, &pk);
  cr_assert(errcode == SKY_ERROR, "0 random bytes");

  randBytes(buff, 100);
  slice.len = 100;
  errcode = SKY_cipher_NewPubKey(slice, &pk);
  cr_assert(errcode == SKY_ERROR, "100 random bytes");

  randBytes(buff, 33);
  slice.len = 33;
  errcode = SKY_cipher_NewPubKey(slice, &pk);
  cr_assert(errcode == SKY_OK, "33 random bytes");

  cr_assert(eq(u8[33], pk, buff));
}

