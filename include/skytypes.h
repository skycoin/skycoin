
#ifndef SKYSTRUCTS_H
#define SKYSTRUCTS_H

/**
 * Go 8-bit signed integer values.
 */
typedef signed char GoInt8_;
/**
 * Go 8-bit unsigned integer values.
 */
typedef unsigned char GoUint8_;
/**
 * Go 16-bit signed integer values.
 */
typedef short GoInt16_;
/**
 * Go 16-bit unsigned integer values.
 */
typedef unsigned short GoUint16_;
/**
 * Go 32-bit signed integer values.
 */
typedef int GoInt32_;
/**
 * Go 32-bit unsigned integer values.
 */
typedef unsigned int GoUint32_;
/**
 * Go 64-bit signed integer values.
 */
typedef long long GoInt64_;
/**
 * Go 64-bit unsigned integer values.
 */
typedef unsigned long long GoUint64_;
/**
 * Go integer values aligned to the word size of the underlying architecture.
 */
typedef GoInt64_ GoInt_;
/**
 * Go unsigned integer values aligned to the word size of the underlying
 * architecture.
 */
typedef GoUint64_ GoUint_;
/**
 * Architecture-dependent type representing instances Go `uintptr` type.
 * Used as a generic representation of pointer types.
 */
typedef __SIZE_TYPE__ GoUintptr_;
/**
 * Go single precision 32-bits floating point values.
 */
typedef float GoFloat32_;
/**
 * Go double precision 64-bits floating point values.
 */
typedef double GoFloat64_;
/**
 * Instances of Go `complex` type.
 */
typedef float _Complex GoComplex64_;
/**
 * Instances of Go `complex` type.
 */
typedef double _Complex GoComplex128_;
typedef short bool;
typedef GoUint32_ error;

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt._
*/
typedef char _check_for_64_bit_pointer_matchingGoInt[sizeof(void*)==64/8 ? 1:-1];

/**
 * Instances of Go `string` type.
 */
typedef struct {
  const char *p;    ///< Pointer to string characters buffer.
  GoInt_ n;         ///< String size not counting trailing `\0` char
                    ///< if at all included.
} GoString_;
/**
 * Instances of Go `map` type.
 */
typedef void *GoMap_;
/**
 * Instances of Go `chan` channel types.
 */
typedef void *GoChan_;

/**
 * Memory handles returned back to the caller and manipulated
 * internally by API functions. Usually used to avoid type dependencies
 * with internal implementation types.
 */
typedef GoInt64_ Handle;

/**
 * Instances of Go interface types.
 */
typedef struct {
  void *t;      ///< Pointer to the information of the concrete Go type
                ///< bound to this interface reference.
  void *v;      ///< Pointer to the data corresponding to the value 
                ///< bound to this interface type.
} GoInterface_;
/**
 * Instances of Go slices
 */
typedef struct {
  void *data;   ///< Pointer to buffer containing slice data.
  GoInt_ len;   ///< Number of items stored in slice buffer
  GoInt_ cap;   ///< Maximum number of items that fits in this slice
                ///< considering allocated memory and item type's
                ///< size.
} GoSlice_;

typedef struct {	
	bool 		neg;
	GoSlice_ 	nat;
} Number;

/**
 * RIPEMD-160 hash.
 */
typedef unsigned char Ripemd160[20];

typedef struct {	
	//TODO: stdevEclipse Define Signature
	Number a;
	Number b;
} Signature;

#include "skytypes.gen.h"

/**
 * Internal representation of a Skycoin wallet.
 */
typedef struct {
	GoMap_ Meta;        ///< Records items that are not deterministic, like filename, lable, wallet type, secrets, etc.
	GoSlice_ Entries;   ///< Entries field stores the address entries that are deterministically generated from seed.
} Wallet;


#endif

