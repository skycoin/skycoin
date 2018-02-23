/* Created by "go tool cgo" - DO NOT EDIT. */

/* package command-line-arguments */

/* Start of preamble from import "C" comments.  */


#line 3 "/Users/olemis/Documents/workspace/work/go/src/github.com/skycoin/skycoin/lib/cgo/cipher.address.go"

#include <string.h>
#include <stdlib.h>

#include "../../include/skytypes.h"


#line 1 "cgo-generated-wrapper"

#line 3 "/Users/olemis/Documents/workspace/work/go/src/github.com/skycoin/skycoin/lib/cgo/cli.transaction.go"

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


extern unsigned int SKY_Cipher_DecodeBase58Address(GoString p0, Address* p1);

extern void SKY_Cipher_AddressFromPubKey(PubKey* p0, Address* p1);

extern void SKY_Cipher_AddressFromSecKey(SecKey* p0, Address* p1);

extern unsigned int SKY_Cipher_BitcoinDecodeBase58Address(GoString p0, Address* p1);

extern void SKY_Cipher_Address_Bytes(Address* p0, unsigned char* p1);

extern void SKY_Cipher_Address_BitcoinBytes(Address* p0, unsigned char* p1);

extern unsigned int SKY_Cipher_Address_Verify(Address* p0, PubKey* p1);

extern GoString SKY_Cipher_Address_String(Address* p0);

extern GoString SKY_Cipher_Address_BitcoinString(Address* p0);

extern void SKY_Cipher_Address_Checksum(Address* p0, Checksum* p1);

extern void SKY_Cipher_Address_BitcoinChecksum(Address* p0, Checksum* p1);

extern GoString SKY_Cipher_BitcoinAddressFromPubkey(PubKey* p0);

extern GoString SKY_Cipher_BitcoinWalletImportFormatFromSeckey(SecKey* p0);

extern unsigned int SKY_Cipher_BitcoinAddressFromBytes(unsigned char* p0, size_t p1, Address* p2);

extern unsigned int SKY_Cipher_SecKeyFromWalletImportFormat(GoString p0, SecKey* p1);

extern unsigned int SKY_CLI_CreateRawTxFromWallet(Handle p0, GoString p1, GoString p2, _GoSlice* p3, Transaction* p4);

extern unsigned int SKY_CLI_CreateRawTxFromAddress(Handle p0, GoString p1, GoString p2, GoString p3, _GoSlice p4, Transaction* p5);

extern void SKY_CLI_CreateRawTx(Handle p0, Wallet* p1, _GoSlice p2, GoString p3, _GoSlice p4, Transaction* p5);

extern void SKY_CLI_NewTransaction(_GoSlice* p0, GoSlice p1, GoSlice p2, Transaction* p3);

#ifdef __cplusplus
}
#endif
