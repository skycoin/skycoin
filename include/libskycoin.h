/* Created by "go tool cgo" - DO NOT EDIT. */

/* package command-line-arguments */

/* Start of preamble from import "C" comments.  */


#line 10 "/lib/cgo/cipher.crypto.go"


  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"

#line 1 "cgo-generated-wrapper"

#line 3 "/lib/cgo/api.cli.create_rawtx.go"

#include <string.h>
#include <stdlib.h>

#include "skytypes.h"


#line 1 "cgo-generated-wrapper"

#line 3 "/lib/cgo/cipher.address.go"

#include <string.h>
#include <stdlib.h>

#include "skytypes.h"


#line 1 "cgo-generated-wrapper"

#line 10 "/lib/cgo/cipher.hash.go"


  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"

#line 1 "cgo-generated-wrapper"


/* End of preamble from import "C" comments.  */


/* Start of boilerplate cgo prologue.  */
#line 1 "cgo-gcc-export-header-prolog"

#ifndef GO_CGO_PROLOGUE_H
#define GO_CGO_PROLOGUE_H

typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef __SIZE_TYPE__ GoUintptr;
typedef float GoFloat32;
typedef double GoFloat64;
typedef float _Complex GoComplex64;
typedef double _Complex GoComplex128;

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt.
*/
typedef char _check_for_64_bit_pointer_matching_GoInt[sizeof(void*)==64/8 ? 1:-1];

typedef struct { const char *p; GoInt n; } GoString;
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;

#endif

/* End of boilerplate cgo prologue.  */

#ifdef __cplusplus
extern "C" {
#endif


extern GoInt SKY_cipher_PubKeySlice_Len(cipher__PubKeySlice* p0);

extern GoUint8 SKY_cipher_PubKeySlice_Less(cipher__PubKeySlice* p0, GoInt p1, GoInt p2);

extern void SKY_cipher_PubKeySlice_Swap(cipher__PubKeySlice* p0, GoInt p1, GoInt p2);

extern void SKY_cipher_RandByte(GoInt p0, cipher__PubKeySlice* p1);

extern GoUint32 SKY_cipher_NewPubKey(GoSlice p0, cipher__PubKey* p1);

extern GoUint32 SKY_cipher_PubKeyFromHex(GoString p0, cipher__PubKey* p1);

extern GoUint32 SKY_cipher_PubKeyFromSecKey(cipher__SecKey* p0, cipher__PubKey* p1);

extern GoUint32 SKY_cipher_PubKeyFromSig(cipher__Sig* p0, cipher__SHA256* p1, cipher__PubKey* p2);

extern GoUint32 SKY_cipher_PubKey_Verify(cipher__PubKey* p0);

extern void SKY_cipher_PubKey_Hex(cipher__PubKey* p0, GoString_* p1);

extern void SKY_cipher_PubKey_ToAddressHash(cipher__PubKey* p0, cipher__Ripemd160* p1);

extern GoUint32 SKY_cipher_NewSecKey(GoSlice p0, cipher__SecKey* p1);

extern GoUint32 SKY_cipher_SecKeyFromHex(GoString p0, cipher__SecKey* p1);

extern GoUint32 SKY_cipher_SecKey_Verify(cipher__SecKey* p0);

extern void SKY_cipher_SecKey_Hex(cipher__SecKey* p0, GoString_* p1);

extern void SKY_cipher_ECDH(cipher__PubKey* p0, cipher__SecKey* p1, cipher__PubKeySlice* p2);

extern GoUint32 SKY_cipher_NewSig(GoSlice p0, cipher__Sig* p1);

extern GoUint32 SKY_cipher_SigFromHex(GoString p0, cipher__Sig* p1);

extern void SKY_cipher_Sig_Hex(cipher__Sig* p0, GoString_* p1);

extern void SKY_cipher_SignHash(cipher__SHA256* p0, cipher__SecKey* p1, cipher__Sig* p2);

extern GoUint32 SKY_cipher_ChkSig(cipher__Address* p0, cipher__SHA256* p1, cipher__Sig* p2);

extern GoUint32 SKY_cipher_VerifySignedHash(cipher__Sig* p0, cipher__SHA256* p1);

extern GoUint32 SKY_cipher_VerifySignature(cipher__PubKey* p0, cipher__Sig* p1, cipher__SHA256* p2);

extern void SKY_cipher_GenerateKeyPair(cipher__PubKey* p0, cipher__SecKey* p1);

extern void SKY_cipher_GenerateDeterministicKeyPair(GoSlice p0, cipher__PubKey* p1, cipher__SecKey* p2);

extern void SKY_cipher_DeterministicKeyPairIterator(GoSlice p0, cipher__PubKeySlice* p1, cipher__PubKey* p2, cipher__SecKey* p3);

extern void SKY_cipher_GenerateDeterministicKeyPairs(GoSlice p0, GoInt p1, cipher__PubKeySlice* p2);

extern void SKY_cipher_GenerateDeterministicKeyPairsSeed(GoSlice p0, GoInt p1, cipher__PubKeySlice* p2, cipher__PubKeySlice* p3);

extern GoUint32 SKY_cipher_TestSecKey(cipher__SecKey* p0);

extern GoUint32 SKY_cipher_TestSecKeyHash(cipher__SecKey* p0, cipher__SHA256* p1);

extern GoUint32 SKY_cli_CreateRawTxFromWallet(Handle p0, GoString p1, GoString p2, GoSlice p3, coin__Transaction* p4);

extern GoUint32 SKY_cli_CreateRawTxFromAddress(Handle p0, GoString p1, GoString p2, GoString p3, GoSlice p4, coin__Transaction* p5);

extern void SKY_cli_CreateRawTx(Handle p0, wallet__Wallet* p1, GoSlice p2, GoString p3, GoSlice p4, coin__Transaction* p5);

extern void SKY_cli_NewTransaction(GoSlice p0, GoSlice p1, GoSlice p2, coin__Transaction* p3);

extern GoUint32 SKY_cipher_DecodeBase58Address(GoString p0, cipher__Address* p1);

extern void SKY_cipher_AddressFromPubKey(cipher__PubKey* p0, cipher__Address* p1);

extern void SKY_cipher_AddressFromSecKey(cipher__SecKey* p0, cipher__Address* p1);

extern GoUint32 SKY_cipher_BitcoinDecodeBase58Address(GoString p0, cipher__Address* p1);

extern void SKY_cipher_Address_Bytes(cipher__Address* p0, cipher__PubKeySlice* p1);

extern void SKY_cipher_Address_BitcoinBytes(cipher__Address* p0, cipher__PubKeySlice* p1);

extern GoUint32 SKY_cipher_Address_Verify(cipher__Address* p0, cipher__PubKey* p1);

extern void SKY_cipher_Address_String(cipher__Address* p0, GoString_* p1);

extern void SKY_cipher_Address_BitcoinString(cipher__Address* p0, GoString_* p1);

extern void SKY_cipher_Address_Checksum(cipher__Address* p0, cipher__Checksum* p1);

extern void SKY_cipher_Address_BitcoinChecksum(cipher__Address* p0, cipher__Checksum* p1);

extern void SKY_cipher_BitcoinAddressFromPubkey(cipher__PubKey* p0, GoString_* p1);

extern void SKY_cipher_BitcoinWalletImportFormatFromSeckey(cipher__SecKey* p0, GoString_* p1);

extern GoUint32 SKY_cipher_BitcoinAddressFromBytes(GoSlice p0, cipher__Address* p1);

extern GoUint32 SKY_cipher_SecKeyFromWalletImportFormat(GoString p0, cipher__SecKey* p1);

extern GoUint32 SKY_cipher_Ripemd160_Set(cipher__Ripemd160* p0, GoSlice p1);

extern void SKY_cipher_HashRipemd160(GoSlice p0, cipher__Ripemd160* p1);

extern GoUint32 SKY_cipher_SHA256_Set(cipher__SHA256* p0, GoSlice p1);

extern void SKY_cipher_SHA256_Hex(cipher__SHA256* p0, GoString_* p1);

extern void SKY_cipher_SHA256_Xor(cipher__SHA256* p0, cipher__SHA256* p1, cipher__SHA256* p2);

extern GoUint32 SKY_cipher_SumSHA256(GoSlice p0, cipher__SHA256* p1);

extern GoUint32 SKY_cipher_SHA256FromHex(GoString p0, cipher__SHA256* p1);

extern void SKY_cipher_DoubleSHA256(GoSlice p0, cipher__SHA256* p1);

extern void SKY_cipher_AddSHA256(cipher__SHA256* p0, cipher__SHA256* p1, cipher__SHA256* p2);

extern void SKY_cipher_Merkle(GoSlice* p0, cipher__SHA256* p1);

#ifdef __cplusplus
}
#endif
