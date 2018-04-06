
#include <string.h>
#include "skycriterion.h"
#include "skystring.h"

int cr_user_Address_eq(Address *addr1, Address *addr2){
  if(addr1->Version != addr2->Version)
    return 0;
  for (int i = 0; i < sizeof(Ripemd160); ++i) {
    if(addr1->Key[i] != addr2->Key[i])
      return 0;
  }
  return 1;
}

char *cr_user_Address_tostr(Address *addr1)
{
  char *out;

  cr_asprintf(&out, "(Address) { .Key = %s, .Version = %llu }", addr1->Key, (unsigned long long) addr1->Version);
  return out;
}

int cr_user_Address_noteq(Address *addr1, Address *addr2){
  if(addr1->Version != addr2->Version)
    return 0;
  for (int i = 0; i < sizeof(Ripemd160); ++i) {
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

int cr_user_SecKey_eq(SecKey *seckey1, SecKey *seckey2){
  return memcmp((void *)seckey1,(void *)seckey2, sizeof(SecKey)) == 0;
}

char *cr_user_SecKey_tostr(SecKey *seckey1)
{
  char *out;
  char hexdump[101];

  strnhex((unsigned char *)seckey1, hexdump, sizeof(SecKey));
  cr_asprintf(&out, "(SecKey) { .SecKey = %s }", hexdump);
  return out;
}


int cr_user_Ripemd160_noteq(Ripemd160 *rp1, Ripemd160 *rp2){
  return memcmp((void *)rp1,(void *)rp2, sizeof(Ripemd160)) != 0;
}

int cr_user_Ripemd160_eq(Ripemd160 *rp1, Ripemd160 *rp2){
  return memcmp((void *)rp1,(void *)rp2, sizeof(Ripemd160)) == 0;
}

char *cr_user_Ripemd160_tostr(Ripemd160 *rp1)
{
  char *out;
  char hexdump[101];

  strnhex((unsigned char *)rp1, hexdump, sizeof(Ripemd160));
  cr_asprintf(&out, "(Ripemd160) { %s }", hexdump );
}

int cr_user_SHA256_noteq(SHA256 *sh1, SHA256 *sh2){
  return memcmp((void *)sh1,(void *)sh1, sizeof(SHA256)) != 0;
}

int cr_user_SHA256_eq(SHA256 *sh1, SHA256 *sh2){
  return memcmp((void *)sh1,(void *)sh1, sizeof(SHA256)) == 0;
}

char *cr_user_SHA256_tostr(SHA256 *sh1) {
  char *out;
  char hexdump[101];

  strnhex((unsigned char *)sh1, hexdump, sizeof(SHA256));
  cr_asprintf(&out, "(SHA256) { %s }", hexdump);
  return out;
}

