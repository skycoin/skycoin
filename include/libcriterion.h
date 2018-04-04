#ifndef LIBCRITERION_H
#define LIBCRITERION_H





#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skytest.h"


// // TODO: Write like this cr_assert(eq(type(Address), addr1, addr2))
extern int cr_user_Address_eq(Address *addr1, Address *addr2);

extern char *cr_user_Address_tostr(Address *addr1);
// // TODO: Write like this cr_assert(not(eq(type(Address), addr1, addr2)))
extern int cr_user_Address_noteq(Address *addr1, Address *addr2);

extern int cr_user_GoString_eq(GoString *string1, GoString *string2);

extern char *cr_user_GoString_tostr(GoString *string);

extern char *cr_user_GoString__tostr(GoString_ *string) ;


// // TODO: Write like this cr_assert(eq(type(SecKey), seckey1, seckey2))
extern int cr_user_SecKey_eq(SecKey *seckey1, SecKey *seckey2);

extern char *cr_user_SecKey_tostr(SecKey *seckey1);


extern int cr_user_Ripemd160_noteq(Ripemd160 *rp1, Ripemd160 *rp2);

extern int cr_user_Ripemd160_eq(Ripemd160 *rp1, Ripemd160 *rp2);

extern char *cr_user_Ripemd160_tostr(Ripemd160 *rp1);


// TODO: Write like this cr_assert(eq(type(GoSlice), slice1, slice2))
extern int cr_user_GoSlice_eq(GoSlice *slice1, GoSlice *slice2);

extern char *cr_user_GoSlice_tostr(GoSlice *slice1);
// // TODO: Write like this cr_assert(not(eq(type(GoSlice), slice1, slice2)))
extern int cr_user_GoSlice_noteq(GoSlice *slice1, GoSlice *slice2);


extern int cr_user_SHA256_noteq(SHA256 *sh1, SHA256 *sh2);

extern int cr_user_SHA256_eq(SHA256 *sh1, SHA256 *sh2);

extern char *cr_user_SHA256_tostr(SHA256 *sh1);


extern int cr_user_char_eq(unsigned char *string1, unsigned char *string2);

extern int cr_user_char_noteq(unsigned char *string1, unsigned char *string2);

extern char *cr_user_char_tostr(unsigned char *string1);



#endif //LIBCRITERION_H