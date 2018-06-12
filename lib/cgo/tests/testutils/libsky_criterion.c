
#include <string.h>
#include "skycriterion.h"
#include "skystring.h"

int cr_user_cipher__Address_eq(cipher__Address *addr1, cipher__Address *addr2){
  if(addr1->Version != addr2->Version)
    return 0;
  return memcmp((void*)addr1, (void*) addr2, sizeof(cipher__Address)) == 0;
}

char *cr_user_cipher__Address_tostr(cipher__Address *addr1)
{
  char *out;

  cr_asprintf(&out, "(cipher__Address) { .Key = %s, .Version = %llu }", addr1->Key, (unsigned long long) addr1->Version);
  return out;
}

int cr_user_cipher__Address_noteq(cipher__Address *addr1, cipher__Address *addr2){
  if(addr1->Version == addr2->Version)
    return 0;
  return memcmp((void*)addr1, (void*) addr2, sizeof(cipher__Address)) != 0;
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
  return cr_user_GoString_eq((GoString *) string1, (GoString *) string2);
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
	return
		(slice1->len == slice2->len) &&
		(memcmp(slice1->data, slice2->data, slice1->len)==0);
}

int cr_user_GoSlice_noteq(GoSlice *slice1, GoSlice *slice2){
	return !(((slice1->len == slice2->len)) &&
		(memcmp(slice1->data,slice2->data, slice1->len)==0));
}

char *cr_user_GoSlice_tostr(GoSlice *slice1) {
  char *out;
  cr_asprintf(&out, "(GoSlice) { .data %s, .len %d, .cap %d }", (char*)slice1->data, slice1->len, slice1->cap);
  return out;
}

int cr_user_GoSlice__eq(GoSlice_ *slice1, GoSlice_ *slice2){
	return ((slice1->len == slice2->len)) && (memcmp(slice1->data,slice2->data, slice1->len)==0 );
}

char *cr_user_GoSlice__tostr(GoSlice_ *slice1) {
  char *out;
  cr_asprintf(&out, "(GoSlice_) { .data %s, .len %d, .cap %d }", (char*)slice1->data, slice1->len, slice1->cap);
  return out;
}

int cr_user_secp256k1go__Field_eq(secp256k1go__Field* f1, secp256k1go__Field* f2){
 for( int i = 0; i < 10; i++){
  if( f1->n[i] != f2->n[i])
   return 0;
}
return 1;
}

int cr_user_coin__Transactions_eq(coin__Transactions *slice1, coin__Transactions *slice2){
	return
		(slice1->len == slice2->len) &&
		(memcmp(slice1->data, slice2->data, slice1->len)==0);
}

int cr_user_coin__Transactions_noteq(coin__Transactions *slice1, coin__Transactions *slice2){
	return
		!((slice1->len == slice2->len) &&
		(memcmp(slice1->data, slice2->data, slice1->len)==0));
}

char *cr_user_coin__Transactions_tostr(coin__Transactions *slice1) {
  char *out;
  cr_asprintf(&out, "(coin__Transactions) { .data %s, .len %d, .cap %d }", (char*)slice1->data, slice1->len, slice1->cap);
  return out;
}

int cr_user_coin__BlockBody_eq(coin__BlockBody *b1, coin__BlockBody *b2){
	return
		cr_user_GoSlice__eq((GoSlice_*)&(b1->Transactions), (GoSlice_*)&(b2->Transactions));
}

int cr_user_coin__BlockBody_noteq(coin__BlockBody *b1, coin__BlockBody *b2){
	return
		!cr_user_GoSlice__eq((GoSlice_*)&(b1->Transactions), (GoSlice_*)&(b2->Transactions));
}

char *cr_user_coin__BlockBody_tostr(coin__BlockBody *b) {
  char *out;
  cr_asprintf(&out, "(coin__BlockBody) { .data %s, .len %d, .cap %d }", (char*)b->Transactions.data, b->Transactions.len, b->Transactions.cap);
  return out;
}

int cr_user_coin__UxOut_eq(coin__UxOut *x1, coin__UxOut *x2){
	return memcmp(x1, x2, sizeof(coin__UxOut)) == 0;
}

int cr_user_coin__UxOut_noteq(coin__UxOut *x1, coin__UxOut *x2){
	return memcmp(x1, x2, sizeof(coin__UxOut)) != 0;
}

char* cr_user_coin__UxOut_tostr(coin__UxOut *x1){
  char *out;
  cr_asprintf(&out, "(coin__UxOut) { %s }", (char*)x1);
  return out;
}

int cr_user_coin__Transaction_eq(coin__Transaction *x1, coin__Transaction *x2){
	if( x1->Length != x2->Length ||
      x1->Type != x2->Type ){
      return 0;
  }
  if(!cr_user_cipher__SHA256_eq(&x1->InnerHash, &x2->InnerHash))
    return 0;
  if(!cr_user_GoSlice__eq(&x1->Sigs, &x2->Sigs) )
    return 0;
  if(!cr_user_GoSlice__eq(&x1->In, &x2->In) )
    return 0;
  if(!cr_user_GoSlice__eq(&x1->Out, &x2->Out) )
    return 0;
  return 1;
}

int cr_user_coin__Transaction_noteq(coin__Transaction *x1, coin__Transaction *x2){
	return !cr_user_coin__Transaction_eq(x1, x2);
}

char* cr_user_coin__Transaction_tostr(coin__Transaction *x1){
  char *out;
  cr_asprintf(&out, "(coin__Transaction) { Length : %d }", x1->Length);
  return out;
}

int cr_user_coin__TransactionOutput_eq(coin__TransactionOutput *x1, coin__TransactionOutput *x2){
	if( x1->Coins != x2->Coins ||
      x1->Hours != x2->Hours ){
      return 0;
  }

  if(!cr_user_cipher__Address_eq(&x1->Address, &x2->Address))
    return 0;
  return 1;
}

int cr_user_coin__TransactionOutput_noteq(coin__TransactionOutput *x1, coin__TransactionOutput *x2){
	return !cr_user_coin__TransactionOutput_eq(x1, x2);
}

char* cr_user_coin__TransactionOutput_tostr(coin__TransactionOutput *x1){
  char *out;
  cr_asprintf(&out, "(coin__TransactionOutput) { Coins : %d, Hours: %d, Address: %s }", x1->Coins, x1->Hours, x1->Address);
  return out;
}
