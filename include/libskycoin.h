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


extern unsigned int SKY_cli_CreateRawTxFromWallet(Handle p0, GoString p1, GoString p2, GoSlice p3, Transaction* p4);

extern unsigned int SKY_cli_CreateRawTxFromAddress(Handle p0, GoString p1, GoString p2, GoString p3, GoSlice p4, Transaction* p5);

extern void SKY_cli_CreateRawTx(Handle p0, Wallet* p1, GoSlice p2, GoString p3, GoSlice p4, Transaction* p5);

extern void SKY_cli_NewTransaction(GoSlice p0, GoSlice p1, GoSlice p2, Transaction* p3);

extern unsigned int SKY_cipher_DecodeBase58Address(GoString p0, Address* p1);

extern void SKY_cipher_AddressFromPubKey(PubKey* p0, Address* p1);

extern void SKY_cipher_AddressFromSecKey(SecKey* p0, Address* p1);

extern unsigned int SKY_cipher_BitcoinDecodeBase58Address(GoString p0, Address* p1);

extern void SKY_cipher_Address_Bytes(Address* p0, GoSlice_* p1);

extern void SKY_cipher_Address_BitcoinBytes(Address* p0, GoSlice_* p1);

extern unsigned int SKY_cipher_Address_Verify(Address* p0, PubKey* p1);

extern GoString SKY_cipher_Address_String(Address* p0);

extern GoString SKY_cipher_Address_BitcoinString(Address* p0);

extern void SKY_cipher_Address_Checksum(Address* p0, Checksum* p1);

extern void SKY_cipher_Address_BitcoinChecksum(Address* p0, Checksum* p1);

extern GoString SKY_cipher_BitcoinAddressFromPubkey(PubKey* p0);

extern GoString SKY_cipher_BitcoinWalletImportFormatFromSeckey(SecKey* p0);

extern unsigned int SKY_cipher_BitcoinAddressFromBytes(GoSlice p0, Address* p1);

extern unsigned int SKY_cipher_SecKeyFromWalletImportFormat(GoString p0, SecKey* p1);

#ifdef __cplusplus
}
#endif
