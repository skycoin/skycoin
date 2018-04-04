
#include "libcriterion.h"

// // TODO: Write like this cr_assert(eq(type(Address), addr1, addr2))
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
// // TODO: Write like this cr_assert(not(eq(type(Address), addr1, addr2)))
int cr_user_Address_noteq(Address *addr1, Address *addr2){
  if(addr1->Version != addr2->Version)
    return SKY_OK;
  for (int i = 0; i < sizeof(Ripemd160); ++i) {
    if(addr1->Key[i] != addr2->Key[i])
      return SKY_OK;
  }
  return SKY_ERROR;
}

int cr_user_GoString_eq(GoString *string1, GoString *string2){

if (strlen(string1->p) != strlen(string2->p) ) return SKY_ERROR;

  if(  strcmp( (unsigned char *) &string1->p, (unsigned char *) &string2->p) != 0 )
  {
    return SKY_ERROR;
  } else {
    return SKY_OK;
  }
}

char *cr_user_GoString_tostr(GoString *string)
{
  char *out;
  cr_asprintf(&out, "(GoString) { .Data = %s, .Length = %llu }", (unsigned char *)&string->p, (unsigned long long) &string->n);
  return out;
}

int cr_user_GoString__eq(GoString_ *string1, GoString_ *string2){
  return cr_user_GoString_eq((GoString *) &string1, (GoString *) &string2);
}

char *cr_user_GoString__tostr(GoString_ *string) {
  return cr_user_GoString_tostr((GoString *)string);
}


// // TODO: Write like this cr_assert(eq(type(SecKey), seckey1, seckey2))
int cr_user_SecKey_eq(SecKey *seckey1, SecKey *seckey2){
if (strcmp((unsigned char *)seckey1,(unsigned char *)seckey2) != 0)
{
  return SKY_ERROR;
}else {
  return SKY_OK;
}
}

char *cr_user_SecKey_tostr(SecKey *seckey1)
{
  char *out;

  cr_asprintf(&out, "(SecKey) { .SecKey = %s,}", &seckey1);
  return out;
}


int cr_user_Ripemd160_noteq(Ripemd160 *rp1, Ripemd160 *rp2){

  if( strcmp((char *)rp1,(char *)rp2) == 0 ) {
    return SKY_ERROR;
  }else
  return SKY_OK;
}

int cr_user_Ripemd160_eq(Ripemd160 *rp1, Ripemd160 *rp2){

    if( strcmp((char *)rp1,(char *)rp2) == 0 ) {

    return SKY_OK;
  }else
  return SKY_ERROR;
}

char *cr_user_Ripemd160_tostr(Ripemd160 *rp1)
{
  char *out;
  cr_asprintf(&out, "(Ripemd160) { %s }", (unsigned char *)&rp1);
  return out;
}


// TODO: Write like this cr_assert(eq(type(GoSlice), slice1, slice2))
int cr_user_GoSlice_eq(GoSlice *slice1, GoSlice *slice2){
  if(slice1->len != slice1->len)
    return SKY_ERROR;

  if( strcmp(slice1->data,slice2->data) == 0){
    return SKY_OK;
  }
  else{
  return SKY_ERROR;}
}

char *cr_user_GoSlice_tostr(GoSlice *slice1)
{
  char *out;

  cr_asprintf(&out, "(GoSlice) { .data = %s, .len = %llu, .cap = %llu }", slice1->data, (unsigned long long) slice1->len, (unsigned long long)slice1->cap);
  return out;
}
// // TODO: Write like this cr_assert(not(eq(type(GoSlice), slice1, slice2)))
int cr_user_GoSlice_noteq(GoSlice *slice1, GoSlice *slice2){
  if(slice1->len != slice1->len)
    return SKY_OK;

  if( strcmp(slice1->data,slice2->data) == 0){
    return SKY_ERROR;
  }
  else{
  return SKY_OK;}
}


int cr_user_SHA256_noteq(SHA256 *sh1, SHA256 *sh2){

  if( strcmp((char *)sh1,(char *)sh1) == 0 ) {
    return SKY_ERROR;
  }else
  return SKY_OK;
}

int cr_user_SHA256_eq(SHA256 *sh1, SHA256 *sh2){

    if( strcmp((char *)sh1,(char *)sh2) == 0 ) {
    return SKY_OK;
  }else
  return SKY_ERROR;
}

char *cr_user_SHA256_tostr(SHA256 *sh1)
{
  char *out;
  cr_asprintf(&out, "(SHA256) { %s }", &sh1);
  return out;
}


int cr_user_char_eq(unsigned char *string1, unsigned char *string2){

  if( strlen(string1) != strlen(string2) ) return SKY_ERROR;

  if (strcmp(string1,string2) == 0)
  {
    return SKY_OK;
  }

  return SKY_ERROR;
}

int cr_user_char_noteq(unsigned char *string1, unsigned char *string2){

  if( strlen(string1) != strlen(string2) ) return SKY_OK;

  if (strcmp(string1,string2) == 0)
  {
    return SKY_ERROR;
  }
  return SKY_OK;
}

char *cr_user_char_tostr(unsigned char *string1)
{
  char *out;

  cr_asprintf(&out, "(CHAR) {  %s }", string1);
  return out;
}