
#include <string.h>
#include "skycriterion.h"
#include "skystring.h"

int cr_user_cipher__Address_eq(cipher__Address *addr1, cipher__Address *addr2){
  if(addr1->Version != addr2->Version)
    return 0;
  for (int i = 0; i < sizeof(cipher__Ripemd160); ++i) {
    if(addr1->Key[i] != addr2->Key[i])
      return 0;
  }
  return 1;
}

char *cr_user_cipher__Address_tostr(cipher__Address *addr1)
{
  char *out;

  cr_asprintf(&out, "(cipher__Address) { .Key = %s, .Version = %llu }", addr1->Key, (unsigned long long) addr1->Version);
  return out;
}

int cr_user_cipher__Address_noteq(cipher__Address *addr1, cipher__Address *addr2){
  if(addr1->Version != addr2->Version)
    return 0;
  for (int i = 0; i < sizeof(cipher__Ripemd160); ++i) {
    if(addr1->Key[i] != addr2->Key[i])
      return 0;
  }
  return 1;
}

int cr_user_GoString_eq(GoString *string1, GoString *string2){
  return (string1->n == string2->n) &&
    (strcmp( (char *) string1->p, (char *) string2->p) == 0);
}

char *cr_user_GoString_tostr(GoString *string)
{
  char *out;
  cr_asprintf(&out, "(GoString) { .Data = %s, .Length = %llu }",
      string->p, (unsigned long long) string->n);
  return out;
}

int cr_user_GoString__eq(GoString_ *string1, GoString_ *string2){
  return cr_user_GoString_eq((GoString *) &string1, (GoString *) &string2);
}

char *cr_user_GoString__tostr(GoString_ *string) {
  return cr_user_GoString_tostr((GoString *)string);
}

int cr_user_cipher__SecKey_eq(cipher__SecKey *seckey1, cipher__SecKey *seckey2){
  return memcmp((void *)seckey1,(void *)seckey2, sizeof(cipher__SecKey)) == 0;
}

char *cr_user_cipher__SecKey_tostr(cipher__SecKey *seckey1)
{
  char *out;
  char hexdump[101];

  strnhex((unsigned char *)seckey1, hexdump, sizeof(cipher__SecKey));
  cr_asprintf(&out, "(cipher__SecKey) { %s }", hexdump);
  return out;
}


int cr_user_cipher__Ripemd160_noteq(cipher__Ripemd160 *rp1, cipher__Ripemd160 *rp2){
  return memcmp((void *)rp1,(void *)rp2, sizeof(cipher__Ripemd160)) != 0;
}

int cr_user_cipher__Ripemd160_eq(cipher__Ripemd160 *rp1, cipher__Ripemd160 *rp2){
  return memcmp((void *)rp1,(void *)rp2, sizeof(cipher__Ripemd160)) == 0;
}

char *cr_user_cipher__Ripemd160_tostr(cipher__Ripemd160 *rp1)
{
  char *out;
  char hexdump[101];

  strnhex((unsigned char *)rp1, hexdump, sizeof(cipher__Ripemd160));
  cr_asprintf(&out, "(cipher__Ripemd160) { %s }", hexdump );
  return out;
}

int cr_user_cipher__SHA256_noteq(cipher__SHA256 *sh1, cipher__SHA256 *sh2){
  return memcmp((void *)sh1,(void *)sh1, sizeof(cipher__SHA256)) != 0;
}

int cr_user_cipher__SHA256_eq(cipher__SHA256 *sh1, cipher__SHA256 *sh2){
  return memcmp((void *)sh1,(void *)sh1, sizeof(cipher__SHA256)) == 0;
}

char *cr_user_cipher__SHA256_tostr(cipher__SHA256 *sh1) {
  char *out;
  char hexdump[101];

  strnhex((unsigned char *)sh1, hexdump, sizeof(cipher__SHA256));
  cr_asprintf(&out, "(cipher__SHA256) { %s }", hexdump);
  return out;
}

