
#ifndef SKYSTRUCTS_H
#define SKYSTRUCTS_H

typedef signed char GoInt8_;
typedef unsigned char GoUint8_;
typedef short GoInt16_;
typedef unsigned short GoUint16_;
typedef int GoInt32_;
typedef unsigned int GoUint32_;
typedef long long GoInt64_;
typedef unsigned long long GoUint64_;
typedef GoInt64_ GoInt_;
typedef GoUint64_ GoUint_;
typedef __SIZE_TYPE__ GoUintptr_;
typedef float GoFloat32_;
typedef double GoFloat64_;
typedef float _Complex GoComplex64_;
typedef double _Complex GoComplex128_;

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt._
*/
typedef char _check_for_64_bit_pointer_matchingGoInt[sizeof(void*)==64/8 ? 1:-1];

typedef struct { const char *p; GoInt_ n; } GoString_;
typedef void *GoMap_;
typedef void *GoChan_;
typedef struct { void *t; void *v; } GoInterface_;
typedef struct { void *data; GoInt_ len; GoInt_ cap; } GoSlice_;
typedef unsigned char Ripemd160[20];

typedef struct {
	unsigned char Version;
	Ripemd160 Key;
} Address;

typedef unsigned char PubKey[33];
typedef unsigned char SecKey[32];
typedef unsigned char Checksum[4];

typedef struct {
	GoString_ Addr;
	GoInt64_ Coins;
} SendAmount;

typedef GoInt64_ Handle;

typedef unsigned char SHA256[32];
typedef unsigned char Sig[65];

typedef struct {
	Address Address;
	GoInt64_ Coins;
	GoInt64_ Hours;
} TransactionOutput; 

typedef struct {
	GoInt32_ Length;
	GoInt8_  Type;
	SHA256  InnerHash;

	GoSlice_ Sigs;
	GoSlice_ In;
	GoSlice_ Out;
} Transaction;

typedef struct {
	GoMap_ Meta;
	GoSlice_ Entries;
} Wallet;

typedef struct {
	Address Address;
	PubKey  Public;
	SecKey  Secret; 
} Entry;

typedef struct {
	SHA256   Hash;
	GoInt64_ BkSeq;
	Address  Address;
	GoInt64_ Coins;
	GoInt64_ Hours;
} UxBalance;

#endif

