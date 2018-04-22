#ifndef LIBCRITERION_H
#define LIBCRITERION_H

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"

extern int cr_user_cipher_Address_eq(cipher_Address *addr1, cipher_Address *addr2);
extern char *cr_user_cipher_Address_tostr(cipher_Address *addr1);
extern int cr_user_cipher_Address_noteq(cipher_Address *addr1, cipher_Address *addr2);

extern int cr_user_GoString_eq(GoString *string1, GoString *string2);
extern int cr_user_GoString__eq(GoString_ *string1, GoString_ *string2);

extern char *cr_user_GoString_tostr(GoString *string);
extern char *cr_user_GoString__tostr(GoString_ *string) ;

extern int cr_user_cipher_SecKey_eq(cipher_SecKey *seckey1, cipher_SecKey *seckey2);
extern char *cr_user_cipher_SecKey_tostr(cipher_SecKey *seckey1);

extern int cr_user_cipher_Ripemd160_noteq(cipher_Ripemd160 *rp1, cipher_Ripemd160 *rp2);
extern int cr_user_cipher_Ripemd160_eq(cipher_Ripemd160 *rp1, cipher_Ripemd160 *rp2);
extern char *cr_user_cipher_Ripemd160_tostr(cipher_Ripemd160 *rp1);

extern int cr_user_GoSlice_eq(GoSlice *slice1, GoSlice *slice2);
extern char *cr_user_GoSlice_tostr(GoSlice *slice1);
extern int cr_user_GoSlice_noteq(GoSlice *slice1, GoSlice *slice2);

extern int cr_user_cipher_SHA256_noteq(cipher_SHA256 *sh1, cipher_SHA256 *sh2);
extern int cr_user_cipher_SHA256_eq(cipher_SHA256 *sh1, cipher_SHA256 *sh2);
extern char *cr_user_cipher_SHA256_tostr(cipher_SHA256 *sh1);


#endif //LIBCRITERION_H
