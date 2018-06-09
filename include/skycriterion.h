#ifndef LIBCRITERION_H
#define LIBCRITERION_H

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"

extern int cr_user_cipher__Address_eq(cipher__Address *addr1, cipher__Address *addr2);
extern char *cr_user_cipher__Address_tostr(cipher__Address *addr1);
extern int cr_user_cipher__Address_noteq(cipher__Address *addr1, cipher__Address *addr2);

extern int cr_user_GoString_eq(GoString *string1, GoString *string2);
extern int cr_user_GoString__eq(GoString_ *string1, GoString_ *string2);

extern char *cr_user_GoString_tostr(GoString *string);
extern char *cr_user_GoString__tostr(GoString_ *string) ;

extern int cr_user_cipher__SecKey_eq(cipher__SecKey *seckey1, cipher__SecKey *seckey2);
extern char *cr_user_cipher__SecKey_tostr(cipher__SecKey *seckey1);

extern int cr_user_cipher__Ripemd160_noteq(Ripemd160 *rp1, Ripemd160 *rp2);
extern int cr_user_cipher__Ripemd160_eq(Ripemd160 *rp1, Ripemd160 *rp2);
extern char *cr_user_cipher__Ripemd160_tostr(Ripemd160 *rp1);

extern int cr_user_GoSlice_eq(GoSlice *slice1, GoSlice *slice2);
extern char *cr_user_GoSlice_tostr(GoSlice *slice1);
extern int cr_user_GoSlice_noteq(GoSlice *slice1, GoSlice *slice2);


extern int cr_user_GoSlice__eq(GoSlice_ *slice1, GoSlice_ *slice2);
extern char *cr_user_GoSlice__tostr(GoSlice_ *slice1);
extern int cr_user_GoSlice__noteq(GoSlice_ *slice1, GoSlice_ *slice2);

extern int cr_user_cipher__SHA256_noteq(cipher__SHA256 *sh1, cipher__SHA256 *sh2);
extern int cr_user_cipher__SHA256_eq(cipher__SHA256 *sh1, cipher__SHA256 *sh2);
extern char *cr_user_cipher__SHA256_tostr(cipher__SHA256 *sh1);

extern int cr_user_secp256k1go__Field_eq(secp256k1go__Field* f1, secp256k1go__Field* f2);

extern int cr_user_coin__BlockBody_eq(coin__BlockBody *b1, coin__BlockBody *b2);
extern int cr_user_coin__BlockBody_noteq(coin__BlockBody *b1, coin__BlockBody *b2);
extern char *cr_user_coin__BlockBody_tostr(coin__BlockBody *b);

extern int cr_user_coin__UxOut_eq(coin__UxOut *x1, coin__UxOut *x2);
extern int cr_user_coin__UxOut_noteq(coin__UxOut *x1, coin__UxOut *x2);
extern char* cr_user_coin__UxOut_tostr(coin__UxOut *x1);

extern int cr_user_coin__Transaction_eq(coin__Transaction *x1, coin__Transaction *x2);
extern int cr_user_coin__Transaction_noteq(coin__Transaction *x1, coin__Transaction *x2);
extern char* cr_user_coin__Transaction_tostr(coin__Transaction *x1);

extern int cr_user_coin__TransactionOutput_eq(coin__TransactionOutput *x1, coin__TransactionOutput *x2);
extern int cr_user_coin__TransactionOutput_noteq(coin__TransactionOutput *x1, coin__TransactionOutput *x2);
extern char* cr_user_coin__TransactionOutput_tostr(coin__TransactionOutput *x1);

#endif //LIBCRITERION_H
