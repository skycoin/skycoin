
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


int cr_user_cipher__Ripemd160_noteq(Ripemd160 *rp1, Ripemd160 *rp2){
  return memcmp((void *)rp1,(void *)rp2, sizeof(Ripemd160)) != 0;
}

int cr_user_cipher__Ripemd160_eq(Ripemd160 *rp1, Ripemd160 *rp2){
  return memcmp((void *)rp1,(void *)rp2, sizeof(Ripemd160)) == 0;
}

char *cr_user_cipher__Ripemd160_tostr(Ripemd160 *rp1)
{
  char *out;
  char hexdump[101];

  strnhex((unsigned char *)rp1, hexdump, sizeof(cipher__Ripemd160));
  cr_asprintf(&out, "(cipher__Ripemd160) { %s }", hexdump );
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

int cr_user_GoSlice_eq(GoSlice *slice1, GoSlice *slice2){
	return ((slice1->len == slice2->len)) && 
		(memcmp(slice1->data,slice2->data, sizeof(GoSlice))==0);
}

int cr_user_GoSlice_noteq(GoSlice *slice1, GoSlice *slice2){
	return !(((slice1->len == slice2->len)) && 
		(memcmp(slice1->data,slice2->data, sizeof(GoSlice))==0));
}

char *cr_user_GoSlice_tostr(GoSlice *slice1) {
  char *out;
  cr_asprintf(&out, "(GoSlice) { .data %s, .len %d, .cap %d }", slice1->data,slice1->len,slice1->cap);
  return out;
}

int cr_user_GoSlice__eq(GoSlice_ *slice1, GoSlice_ *slice2){

  return ((slice1->len == slice2->len)) && (memcmp(slice1->data,slice2->data, sizeof(GoSlice_))==0 );

}

char *cr_user_GoSlice__tostr(GoSlice_ *slice1) {
  char *out;
  cr_asprintf(&out, "(GoSlice_) { .data %s, .len %d, .cap %d }", slice1->data,slice1->len,slice1->cap);
  return out;
}

int cr_user_secp256k1go__Field_eq(secp256k1go__Field* f1, secp256k1go__Field* f2){
 for( int i = 0; i < 10; i++){
  if( f1->n[i] != f2->n[i])
   return 0;
}
return 1;
}
