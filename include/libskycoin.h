/* Created by "go tool cgo" - DO NOT EDIT. */

/* package command-line-arguments */

/* Start of preamble from import "C" comments.  */


#line 3 "/Users/olemis/Documents/workspace/work/go/macos/src/github.com/skycoin/skycoin/lib/cgo/api.cli.create_rawtx.go"

#include <string.h>
#include <stdlib.h>

#include "../../include/skytypes.h"


#line 1 "cgo-generated-wrapper"

#line 3 "/Users/olemis/Documents/workspace/work/go/macos/src/github.com/skycoin/skycoin/lib/cgo/cipher.address.go"

#include <string.h>
#include <stdlib.h>

#include "../../include/skytypes.h"


#line 1 "cgo-generated-wrapper"

#line 9 "/Users/olemis/Documents/workspace/work/go/macos/src/github.com/skycoin/skycoin/lib/cgo/cipher.crypto.go"


  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"

#line 1 "cgo-generated-wrapper"

#line 9 "/Users/olemis/Documents/workspace/work/go/macos/src/github.com/skycoin/skycoin/lib/cgo/cipher.hash.go"


  #include <string.h>
  #include <stdlib.h>

  #include "../../include/skytypes.h"

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


extern GoUint32 SKY_cli_CreateRawTxFromWallet(Handle p0, GoString p1, GoString p2, GoSlice p3, Transaction* p4);

extern GoUint32 SKY_cli_CreateRawTxFromAddress(Handle p0, GoString p1, GoString p2, GoString p3, GoSlice p4, Transaction* p5);

extern void SKY_cli_CreateRawTx(Handle p0, Wallet* p1, GoSlice p2, GoString p3, GoSlice p4, Transaction* p5);

extern void SKY_cli_NewTransaction(GoSlice p0, GoSlice p1, GoSlice p2, Transaction* p3);

extern GoUint32 SKY_cipher_DecodeBase58Address(GoString p0, Address* p1);

extern void SKY_cipher_AddressFromPubKey(PubKey* p0, Address* p1);

extern void SKY_cipher_AddressFromSecKey(SecKey* p0, Address* p1);

extern GoUint32 SKY_cipher_BitcoinDecodeBase58Address(GoString p0, Address* p1);

extern void SKY_cipher_Address_Bytes(Address* p0, PubKeySlice* p1);

extern void SKY_cipher_Address_BitcoinBytes(Address* p0, PubKeySlice* p1);

extern GoUint32 SKY_cipher_Address_Verify(Address* p0, PubKey* p1);

extern GoString SKY_cipher_Address_String(Address* p0);

extern GoString SKY_cipher_Address_BitcoinString(Address* p0);

extern void SKY_cipher_Address_Checksum(Address* p0, Checksum* p1);

extern void SKY_cipher_Address_BitcoinChecksum(Address* p0, Checksum* p1);

extern GoString SKY_cipher_BitcoinAddressFromPubkey(PubKey* p0);

extern GoString SKY_cipher_BitcoinWalletImportFormatFromSeckey(SecKey* p0);

extern GoUint32 SKY_cipher_BitcoinAddressFromBytes(GoSlice p0, Address* p1);

extern GoUint32 SKY_cipher_SecKeyFromWalletImportFormat(GoString p0, SecKey* p1);

extern GoInt SKY_cipher_PubKeySlice_Len(PubKeySlice* p0);

extern GoUint8 SKY_cipher_PubKeySlice_Less(PubKeySlice* p0, GoInt p1, GoInt p2);

extern void SKY_cipher_PubKeySlice_Swap(PubKeySlice* p0, GoInt p1, GoInt p2);

extern void SKY_cipher_RandByte(GoInt p0, PubKeySlice* p1);

extern GoUint32 SKY_cipher_NewPubKey(GoSlice p0, PubKey* p1);

extern GoUint32 SKY_cipher_PubKeyFromHex(GoString p0, PubKey* p1);

extern GoUint32 SKY_cipher_PubKeyFromSecKey(SecKey* p0, PubKey* p1);

extern GoUint32 SKY_cipher_PubKeyFromSig(Sig* p0, SHA256* p1, PubKey* p2);

extern GoUint32 SKY_cipher_PubKey_Verify(PubKey* p0);

extern char* SKY_cipher_PubKey_Hex(PubKey* p0);

extern void SKY_cipher_PubKey_ToAddressHash(PubKey* p0, Ripemd160* p1);

extern void SKY_cipher_NewSecKey(GoSlice p0, SecKey* p1);

extern GoUint32 SKY_cipher_SecKeyFromHex(GoString p0, SecKey* p1);

extern GoUint32 SKY_cipher_SecKey_Verify(SecKey* p0);

extern GoString SKY_cipher_SecKey_Hex(SecKey* p0);

extern void SKY_cipher_ECDH(PubKey* p0, SecKey* p1, PubKeySlice* p2);

extern void SKY_cipher_NewSig(GoSlice p0, Sig* p1);

extern GoUint32 SKY_cipher_SigFromHex(GoString p0, Sig* p1);

extern GoString SKY_cipher_Sig_Hex(Sig* p0);

extern void SKY_cipher_SignHash(SHA256* p0, SecKey* p1, Sig* p2);

extern GoUint32 SKY_cipher_ChkSig(Address* p0, SHA256* p1, Sig* p2);

extern GoUint32 SKY_cipher_VerifySignedHash(Sig* p0, SHA256* p1);

extern GoUint32 SKY_cipher_VerifySignature(PubKey* p0, Sig* p1, SHA256* p2);

extern void SKY_cipher_GenerateKeyPair(PubKey* p0, SecKey* p1);

extern void SKY_cipher_GenerateDeterministicKeyPair(GoSlice p0, PubKey* p1, SecKey* p2);

extern void SKY_cipher_DeterministicKeyPairIterator(GoSlice p0, PubKeySlice* p1, PubKey* p2, SecKey* p3);

extern void SKY_cipher_GenerateDeterministicKeyPairs(GoSlice p0, GoInt p1, PubKeySlice* p2);

extern void SKY_cipher_GenerateDeterministicKeyPairsSeed(GoSlice p0, GoInt p1, PubKeySlice* p2, PubKeySlice* p3);

extern GoUint32 SKY_cipher_TestSecKey(SecKey* p0);

extern GoUint32 SKY_cipher_TestSecKeyHash(SecKey* p0, SHA256* p1);

extern GoUint32 SKY_cipher_Ripemd160_Set(Ripemd160* p0, GoSlice p1);

extern void SKY_cipher_HashRipemd160(GoSlice p0, Ripemd160* p1);

extern GoUint32 SKY_cipher_SHA256_Set(SHA256* p0, GoSlice p1);

extern GoString SKY_cipher_SHA256_Hex(SHA256* p0);

extern void SKY_cipher_SHA256_Xor(SHA256* p0, SHA256* p1, SHA256* p2);

extern void SKY_cipher_SumSHA256(GoSlice p0, SHA256* p1);

extern GoUint32 SKY_cipher_SHA256FromHex(GoString p0, SHA256* p1);

extern void SKY_cipher_DoubleSHA256(GoSlice p0, SHA256* p1);

extern void SKY_cipher_AddSHA256(SHA256* p0, SHA256* p1, SHA256* p2);

extern void SKY_cipher_Merkle(PubKeySlice* p0, SHA256* p1);

#ifdef __cplusplus
}
#endif
