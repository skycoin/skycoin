
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

/**
 * RIPEMD-160 hash.
 */
typedef unsigned char Ripemd160[20];

/**
 * Addresses of SKY accounts
 */
typedef struct {
	unsigned char Version;  ///< Address version identifier.
                          ///< Used to differentiate testnet
                          ///< vs mainnet addresses, for instance.
	Ripemd160 Key;          ///< Address hash identifier.
} Address;

/**
 * Public key, 33-bytes long.
 */
typedef unsigned char PubKey[33];
/**
 * Container type suitable for storing a variable number of
 * public keys.
 */
typedef GoSlice_ PubKeySlice;
/**
 * Secret key, 32 bytes long.
 */
typedef unsigned char SecKey[32];
/**
 * Integrity checksum, 4-bytes long.
 */
typedef unsigned char Checksum[4];

/**
 * Structure used to specify amounts transferred in a transaction.
 */
typedef struct {
	GoString_ Addr; ///< Sender / receipient address.
	GoInt64_ Coins; ///< Amount transferred (e.g. measured in SKY)
} SendAmount;

/**
 * Memory handles returned back to the caller and manipulated
 * internally by API functions. Usually used to avoid type dependencies
 * with internal implementation types.
 */
typedef GoInt64_ Handle;

/**
 * Hash obtained using SHA256 algorithm, 32 bytes long.
 */
typedef unsigned char SHA256[32];
/**
 * Hash signed using a secret key, 65 bytes long.
 */
typedef unsigned char Sig[65];

/**
 * Skycoin transaction output.
 *
 * Instances are integral part of transactions included in blocks.
 */
typedef struct {
	Address Address;  ///< Receipient address.
	GoInt64_ Coins;   ///< Amount sent to the receipient address.
	GoInt64_ Hours;   ///< Amount of Coin Hours sent to the receipient address.
} TransactionOutput;

/**
 * Skycoin transaction.
 *
 * Instances of this struct are included in blocks.
 */
typedef struct {
	GoInt32_ Length;    ///< Current transaction's length expressed in bytes.
	GoInt8_  Type;      ///< Transaction's version. When a node tries to process a transaction, it must verify whether it supports the transaction's type. This is intended to provide a way to update skycoin clients and servers without crashing the network. If the transaction is not compatible with the node, it should not process it.
	SHA256  InnerHash;  ///< It's a SHA256 hash of the inputs and outputs of the transaction. It is used to protect against transaction mutability. This means that the transaction cannot be altered after its creation.

	GoSlice_ Sigs;      ///< A list of digital signiatures generated by the skycoin client using the private key. It is used by Skycoin servers to verify the authenticy of the transaction. Each input requires a different signature.
	GoSlice_ In;        ///< A list of references to unspent transaction outputs. Unlike other cryptocurrencies, such as Bitcoin, Skycoin unspent transaction outputs (UX) and Skycoin transactions (TX) are separated in the blockchain protocol, allowing for lighter transactions, thus reducing the broadcasting costs across the network.
	GoSlice_ Out;       ///< Outputs: A list of outputs created by the client, that will be recorded in the blockchain if transactions are confirmed. An output consists of a data structure representing an UTXT, which is composed by a Skycoin address to be sent to, the amount in Skycoin to be sent, and the amount of Coin Hours to be sent, and the SHA256 hash of the previous fields.
} Transaction;

/**
 * Internal representation of a Skycoin wallet.
 */
typedef struct {
	GoMap_ Meta;        ///< Records items that are not deterministic, like filename, lable, wallet type, secrets, etc.
	GoSlice_ Entries;   ///< Entries field stores the address entries that are deterministically generated from seed.
} Wallet;

/**
 * Wallet entry.
 */
typedef struct {
	Address Address;    ///< Wallet address.
	PubKey  Public;     ///< Public key used to generate address.
	SecKey  Secret;     ///< Secret key used to generate address.
} Entry;

/**
 * Intermediate representation of a UxOut for sorting and spend choosing.
 */
typedef struct {
	SHA256   Hash;      ///< Hash of underlying UxOut.
	GoInt64_ BkSeq;     ///< Block height corresponding to the
                      ///< moment balance calculation is performed at.
	Address  Address;   ///< Account holder address.
	GoInt64_ Coins;     ///< Coins amount (e.g. in SKY).
	GoInt64_ Hours;     ///< Balance of Coin Hours generated by underlying UxOut, depending on UxOut's head time.
} UxBalance;

#endif

