
#ifndef SKYSTRUCTS_H
#define SKYSTRUCTS_H

typedef signed char _GoInt8;
typedef unsigned char _GoUint8;
typedef short _GoInt16;
typedef unsigned short _GoUint16;
typedef int _GoInt32;
typedef unsigned int _GoUint32;
typedef long long _GoInt64;
typedef unsigned long long _GoUint64;
typedef _GoInt64 _GoInt;
typedef _GoUint64 _GoUint;
typedef __SIZE_TYPE__ _GoUintptr;
typedef float _GoFloat32;
typedef double _GoFloat64;
typedef float _Complex _GoComplex64;
typedef double _Complex _GoComplex128;

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of _GoInt.
*/
typedef char _check_for_64_bit_pointer_matching_GoInt[sizeof(void*)==64/8 ? 1:-1];

typedef struct { const char *p; _GoInt n; } _GoString;
typedef void *_GoMap;
typedef void *_GoChan;
typedef struct { void *t; void *v; } _GoInterface;
typedef struct { void *data; _GoInt len; _GoInt cap; } _GoSlice;
typedef unsigned char Ripemd160[20];

typedef struct {
	unsigned char Version;
	Ripemd160 Key;
} Address;

typedef unsigned char Checksum[4];

typedef unsigned char PubKey[33];
typedef unsigned char SecKey[32];
typedef unsigned char Checksum[4];

typedef struct {
	_GoString Addr;
	_GoInt64 Coins;
} SendAmount;

typedef _GoInt64 Handle;

typedef unsigned char SHA256[32];
typedef unsigned char Sig[65];

typedef struct {
	Address Address;
	_GoInt64 Coins;
	_GoInt64 Hours;
} TransactionOutput; 

typedef struct {
	_GoInt32 Length;
	_GoInt8  Type;
	SHA256  InnerHash;

	_GoSlice Sigs;
	_GoSlice In;
	_GoSlice Out;
} Transaction;

typedef struct {
	_GoMap Meta;
	_GoSlice Entries;
} Wallet;

typedef struct {
	Address Address;
	PubKey  Public;
	SecKey  Secret; 
} Entry;

typedef struct {
	SHA256   Hash;
	_GoInt64 BkSeq;
	Address  Address;
	_GoInt64 Coins;
	_GoInt64 Hours;
} UxBalance;

#endif

