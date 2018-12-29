package main

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/base58"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/cipher/encrypt"
	"github.com/skycoin/skycoin/src/cli"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"
)

const (
	// SKY_PKG_LIBCGO package prefix for internal API errors
	SKY_PKG_LIBCGO = 0x7F000000 // nolint megacheck
	// SKY_OK error code is used to report success
	SKY_OK = 0
	// SKY_ERROR generic error condition
	SKY_ERROR = 0x7FFFFFFF
)

const (
	// SKY_BAD_HANDLE invalid handle argument
	SKY_BAD_HANDLE = SKY_PKG_LIBCGO + iota + 1
	// SKY_INVALID_TIMESTRING invalid time value
	SKY_INVALID_TIMESTRING
)

// Package prefixes for error codes
//nolint megacheck
const (
	// Error code prefix for api package
	SKY_PKG_API = (1 + iota) << 24 // nolint megacheck
	// Error code prefix for cipher package and subpackages
	SKY_PKG_CIPHER
	// Error code prefix for cli package
	SKY_PKG_CLI
	// Error code prefix for coin package
	SKY_PKG_COIN
	// Error code prefix for consensus package
	SKY_PKG_CONSENSUS // nolint megacheck
	// Error code prefix for daemon package
	SKY_PKG_DAEMON
	// Error code prefix for gui package
	SKY_PKG_GUI // nolint megacheck
	// Error code prefix for skycoin package
	SKY_PKG_SKYCOIN // nolint megacheck
	// Error code prefix for util package
	SKY_PKG_UTIL
	// Error code prefix for visor package
	SKY_PKG_VISOR
	// Error code prefix for wallet package
	SKY_PKG_WALLET
	// Error code prefix for params package
	SKY_PKG_PARAMS
)

// Error codes defined in cipher package
//nolint megacheck
const (
	// SKY_ErrAddressInvalidLength Unexpected size of address bytes buffer
	SKY_ErrAddressInvalidLength = SKY_PKG_CIPHER + iota
	// SKY_ErrAddressInvalidChecksum Computed checksum did not match expected value
	SKY_ErrAddressInvalidChecksum
	// SKY_ErrAddressInvalidVersion Unsupported address version value
	SKY_ErrAddressInvalidVersion
	// SKY_ErrAddressInvalidPubKey Public key invalid for address
	SKY_ErrAddressInvalidPubKey
	// SKY_ErrAddressInvalidFirstByte Invalid first byte in wallet import format string
	SKY_ErrAddressInvalidFirstByte
	// SKY_ErrAddressInvalidLastByte 33rd byte in wallet import format string is invalid
	SKY_ErrAddressInvalidLastByte
	// SKY_ErrBufferUnderflow bytes in input buffer not enough to deserialize expected type
	SKY_ErrBufferUnderflow
	// SKY_ErrInvalidOmitEmpty field tagged with omitempty and it's not last one in struct
	SKY_ErrInvalidOmitEmpty
	// SKY_ErrInvalidLengthPubKey  Invalid public key length
	SKY_ErrInvalidLengthPubKey
	// SKY_ErrPubKeyFromNullSecKey PubKeyFromSecKey, attempt to load null seckey, unsafe
	SKY_ErrPubKeyFromNullSecKey
	// SKY_ErrPubKeyFromBadSecKey  PubKeyFromSecKey, pubkey recovery failed. Function
	SKY_ErrPubKeyFromBadSecKey
	// SKY_ErrInvalidLengthSecKey Invalid secret key length
	SKY_ErrInvalidLengthSecKey
	// SKY_ErrECHDInvalidPubKey   ECDH invalid pubkey input
	SKY_ErrECHDInvalidPubKey
	// SKY_ErrECHDInvalidSecKey   ECDH invalid seckey input
	SKY_ErrECHDInvalidSecKey
	// SKY_ErrInvalidLengthSig    Invalid signature length
	SKY_ErrInvalidLengthSig
	// SKY_ErrInvalidLengthRipemd160 Invalid ripemd160 length
	SKY_ErrInvalidLengthRipemd160
	// SKY_ErrInvalidLengthSHA256 Invalid sha256 length
	SKY_ErrInvalidLengthSHA256
	// SKY_ErrInvalidBase58Char   Invalid base58 character
	SKY_ErrInvalidBase58Char
	// SKY_ErrInvalidBase58String Invalid base58 string
	SKY_ErrInvalidBase58String
	// SKY_ErrInvalidBase58Length Invalid base58 length
	SKY_ErrInvalidBase58Length
	// SKY_ErrInvalidHexLength       Invalid hex length
	SKY_ErrInvalidHexLength
	// SKY_ErrInvalidBytesLength     Invalid bytes length
	SKY_ErrInvalidBytesLength
	// SKY_ErrInvalidPubKey       Invalid public key
	SKY_ErrInvalidPubKey
	// SKY_ErrInvalidSecKey       Invalid public key
	SKY_ErrInvalidSecKey
	// SKY_ErrInvalidSigPubKeyRecovery Invalig sig: PubKey recovery failed
	SKY_ErrInvalidSigPubKeyRecovery
	// SKY_ErrInvalidSecKeyHex    Invalid SecKey: not valid hex
	SKY_ErrInvalidSecKeyHex // nolint megacheck
	// SKY_ErrInvalidAddressForSig Invalid sig: address does not match output address
	SKY_ErrInvalidAddressForSig
	// SKY_ErrInvalidHashForSig   Signature invalid for hash
	SKY_ErrInvalidHashForSig
	// SKY_ErrPubKeyRecoverMismatch Recovered pubkey does not match pubkey
	SKY_ErrPubKeyRecoverMismatch
	// SKY_ErrInvalidSigInvalidPubKey VerifySignedHash, secp256k1.VerifyPubkey failed
	SKY_ErrInvalidSigInvalidPubKey
	// SKY_ErrInvalidSigValidity  VerifySignedHash, VerifySignatureValidity failed
	SKY_ErrInvalidSigValidity
	// SKY_ErrInvalidSigForMessage Invalid signature for this message
	SKY_ErrInvalidSigForMessage
	// SKY_ErrInvalidSecKyVerification Seckey secp256k1 verification failed
	SKY_ErrInvalidSecKyVerification
	// SKY_ErrNullPubKeyFromSecKey Impossible error, CheckSecKey, nil pubkey recovered
	SKY_ErrNullPubKeyFromSecKey
	// SKY_ErrInvalidDerivedPubKeyFromSecKey impossible error, CheckSecKey, Derived Pubkey verification failed
	SKY_ErrInvalidDerivedPubKeyFromSecKey
	// SKY_ErrInvalidPubKeyFromHash Recovered pubkey does not match signed hash
	SKY_ErrInvalidPubKeyFromHash
	// SKY_ErrPubKeyFromSecKeyMismatch impossible error CheckSecKey, pubkey does not match recovered pubkey
	SKY_ErrPubKeyFromSecKeyMismatch
	// SKY_ErrInvalidLength Unexpected size of string or bytes buffer
	SKY_ErrInvalidLength
	// SKY_ErrBitcoinWIFInvalidFirstByte Unexpected value (!= 0x80) of first byte in Bitcoin Wallet Import Format
	SKY_ErrBitcoinWIFInvalidFirstByte
	// SKY_ErrBitcoinWIFInvalidSuffix Unexpected value (!= 0x01) of 33rd byte in Bitcoin Wallet Import Format
	SKY_ErrBitcoinWIFInvalidSuffix
	// SKY_ErrBitcoinWIFInvalidChecksum Invalid Checksum in Bitcoin WIF address
	SKY_ErrBitcoinWIFInvalidChecksum
	// SKY_ErrEmptySeed Seed input is empty
	SKY_ErrEmptySeed
	// SKY_ErrInvalidSig Invalid signature
	SKY_ErrInvalidSig
	// SKY_ErrMissingPassword missing password
	SKY_ErrMissingPassword
	// SKY_SKY_ErrDataTooLarge data length overflowed, it must <= math.MaxUint32(4294967295)
	SKY_ErrDataTooLarge
	// SKY_ErrInvalidChecksumLength invalid checksum length
	SKY_ErrInvalidChecksumLength
	// SKY_ErrInvalidChecksum invalid data, checksum is not matched
	SKY_ErrInvalidChecksum
	// SKY_ErrInvalidNonceLength invalid nonce length
	SKY_ErrInvalidNonceLength
	// SKY_ErrInvalidBlockSize invalid block size, must be multiple of 32 bytes
	SKY_ErrInvalidBlockSize
	// SKY_ErrReadDataHashFailed read data hash failed: read length != 32
	SKY_ErrReadDataHashFailed
	// SKY_ErrInvalidPassword invalid password SHA256or
	SKY_ErrInvalidPassword
	// SKY_ErrReadDataLengthFailed read data length failed
	SKY_ErrReadDataLengthFailed
	// SKY_ErrInvalidDataLength invalid data length
	SKY_ErrInvalidDataLength
)

// Error codes defined in cli package
// nolint megacheck
const (
	// SKY_ErrTemporaryInsufficientBalance is returned if a wallet does not have
	// enough balance for a spend, but will have enough after unconfirmed transactions confirm
	SKY_ErrTemporaryInsufficientBalance = SKY_PKG_CLI + iota
	// SKY_ErrAddress is returned if an address is invalid
	SKY_ErrAddress
	// ErrWalletName is returned if the wallet file name is invalid
	SKY_ErrWalletName
	// ErrJSONMarshal is returned if JSON marshaling failed
	SKY_ErrJSONMarshal
	// WalletLoadError is returned if a wallet could not be loaded
	SKY_WalletLoadError
	// WalletSaveError is returned if a wallet could not be saved
	SKY_WalletSaveError
)

// Error codes defined in coin package
// nolint megacheck
const (
	// ErrAddEarnedCoinHoursAdditionOverflow is returned by UxOut.CoinHours()
	// if during the addition of base coin
	// hours to additional earned coin hours, the value would overflow a uint64.
	// Callers may choose to ignore this errors and use 0 as the coinhours value instead.
	// This affects one existing spent output, spent in block 13277.
	SKY_ErrAddEarnedCoinHoursAdditionOverflow = SKY_PKG_COIN + iota
	// SKY_ErrUint64MultOverflow is returned when multiplying uint64 values would overflow uint64
	SKY_ErrUint64MultOverflow
	// SKY_ErrUint64AddOverflow is returned when adding uint64 values would overflow uint64
	SKY_ErrUint64AddOverflow
	// SKY_ErrUint32AddOverflow is returned when adding uint32 values would overflow uint32
	SKY_ErrUint32AddOverflow
	// SKY_ErrUint64OverflowsInt64 is returned when converting a uint64 to an int64 would overflow int64
	SKY_ErrUint64OverflowsInt64
	// SKY_ErrInt64UnderflowsUint64 is returned when converting an int64 to a uint64 would underflow uint64
	SKY_ErrInt64UnderflowsUint64
	// SKY_ErrIntUnderflowsUint32 is returned if when converting an int to a uint32 would underflow uint32
	SKY_ErrIntUnderflowsUint32
	// SKY_ErrIntOverflowsUint32 is returned if when converting an int to a uint32 would overflow uint32
	SKY_ErrIntOverflowsUint32
)

// Error codes defined in daemon package
// nolint megacheck
const (
	// SKY_ErrPeerlistFull is returned when the Pex is at a maximum
	SKY_ErrPeerlistFull = SKY_PKG_DAEMON + iota
	// SKY_ErrInvalidAddress is returned when an address appears malformed
	SKY_ErrInvalidAddress
	// SKY_ErrNoLocalhost is returned if a localhost addresses are not allowed
	SKY_ErrNoLocalhost
	// SKY_ErrNotExternalIP is returned if an IP address is not a global unicast address
	SKY_ErrNotExternalIP
	// SKY_ErrPortTooLow is returned if a port is less than 1024
	SKY_ErrPortTooLow
	// SKY_ErrBlacklistedAddress returned when attempting to add a blacklisted peer
	SKY_ErrBlacklistedAddress
	// // SKY_ErrDisconnectReadFailed also includes a remote closed socket
	// SKY_ErrDisconnectReadFailed
	// SKY_ErrDisconnectWriteFailed write faile
	// SKY_ErrDisconnectWriteFailed
	// SKY_ErrDisconnectSetReadDeadlineFailed set read deadline failed
	SKY_ErrDisconnectSetReadDeadlineFailed
	// SKY_ErrDisconnectInvalidMessageLength invalid message length
	SKY_ErrDisconnectInvalidMessageLength
	// SKY_ErrDisconnectMalformedMessage malformed message
	SKY_ErrDisconnectMalformedMessage
	// SKY_ErrDisconnectUnknownMessage unknow message
	SKY_ErrDisconnectUnknownMessage
	// SKY_ErrConnectionPoolClosed error message indicates the connection pool is closed
	SKY_ErrConnectionPoolClosed
	// SKY_ErrWriteQueueFull write queue is full
	SKY_ErrWriteQueueFull
	// SKY_ErrNoReachableConnections when broadcasting a message, no connections were available to send a message to
	SKY_ErrNoReachableConnections
	// SKY_ErrMaxDefaultConnectionsReached returns when maximum number of default connections is reached
	SKY_ErrMaxDefaultConnectionsReached // nolint megacheck
	// SKY_ErrDisconnectReasons invalid version
	SKY_ErrDisconnectVersionNotSupported
	// SKY_ErrDisconnectIntroductionTimeout timeout
	SKY_ErrDisconnectIntroductionTimeout
	// SKY_ErrDisconnectIsBlacklisted is blacklisted
	SKY_ErrDisconnectIsBlacklisted
	// SKY_ErrDisconnectSelf self connnect
	SKY_ErrDisconnectSelf
	// SKY_ErrDisconnectConnectedTwice connect twice
	SKY_ErrDisconnectConnectedTwice
	// SKY_ErrDisconnectIdle idle
	SKY_ErrDisconnectIdle
	// SKY_ErrDisconnectNoIntroduction no introduction
	SKY_ErrDisconnectNoIntroduction
	// SKY_ErrDisconnectIPLimitReached ip limit reached
	SKY_ErrDisconnectIPLimitReached
	// SKY_ErrDisconnectMaxDefaultConnectionReached Maximum number of default connections was reached
	SKY_ErrDisconnectMaxDefaultConnectionReached // nolint megacheck
	// SKY_ErrDisconnectMaxOutgoingConnectionsReached is returned when connection pool size is greater than the maximum allowed
	SKY_ErrDisconnectMaxOutgoingConnectionsReached
	// SKY_ConnectionError represent a failure to connect/dial a connection, with context
	SKY_ConnectionError // nolint megacheck
)

// Error codes defined in util package
// nolint megacheck
const (
	// ErrTxnNoFee is returned if a transaction has no coinhour fee
	SKY_ErrTxnNoFee = SKY_PKG_UTIL + iota
	// ErrTxnInsufficientFee is returned if a transaction's coinhour burn fee is not enough
	SKY_ErrTxnInsufficientFee
	// ErrTxnInsufficientCoinHours is returned if a transaction has more coinhours in its outputs than its inputs
	SKY_ErrTxnInsufficientCoinHours
	// ErrNegativeValue is returned if a balance string is a negative number
	SKY_ErrNegativeValue
	// ErrTooManyDecimals is returned if a balance string has more than 6 decimal places
	SKY_ErrTooManyDecimals
	// ErrTooLarge is returned if a balance string is greater than math.MaxInt64
	SKY_ErrTooLarge
	// ErrEmptyDirectoryName is returned by constructing the full path
	SKY_ErrEmptyDirectoryName
	// ErrDotDirectoryName is returned by constructing the full path of
	SKY_ErrDotDirectoryName
)

// Error codes defined in visor package
// nolint megacheck
const (
	// SKY_ErrHistoryDBCorrupted Internal format error in HistoryDB database
	SKY_ErrHistoryDBCorrupted = SKY_PKG_VISOR + iota
	// SKY_ErrUxOutNotExist is returned if an uxout is not found in historydb
	SKY_ErrUxOutNotExist
	// ErrNoHeadBlock is returned when calling Blockchain.Head() when no head block exists
	SKY_ErrNoHeadBlock
	// SKY_ErrMissingSignature is returned if a block in the db does not have a corresponding signature in the db
	SKY_ErrMissingSignature
	// SKY_ErrUnspentNotExist is returned if an unspent is not found in the pool
	SKY_ErrUnspentNotExist
	// ErrVerifyStopped is returned when database verification is interrupted
	SKY_ErrVerifyStopped
	// SKY_ErrCreateBucketFailed is returned if creating a bolt.DB bucket fails
	SKY_ErrCreateBucketFailed
	// SKY_ErrBucketNotExist is returned if a bolt.DB bucket does not exist
	SKY_ErrBucketNotExist
	// SKY_ErrTxnViolatesHardConstraint is returned when a transaction violates hard constraints
	SKY_ErrTxnViolatesHardConstraint
	// SKY_ErrTxnViolatesSoftConstraint is returned when a transaction violates soft constraints
	SKY_ErrTxnViolatesSoftConstraint
	// SKY_ErrTxnViolatesUserConstraint is returned when a transaction violates user constraints
	SKY_ErrTxnViolatesUserConstraint
)

// Error codes defined in wallet package
// nolint megacheck
const (
	// SKY_ErrInsufficientBalance is returned if a wallet does not have enough balance for a spend
	SKY_ErrInsufficientBalance = SKY_PKG_WALLET + iota
	// SKY_ErrInsufficientHours is returned if a wallet does not have enough hours for a spend with requested hours
	SKY_ErrInsufficientHours
	// SKY_ErrZeroSpend is returned if a transaction is trying to spend 0 coins
	SKY_ErrZeroSpend
	// SKY_ErrSpendingUnconfirmed is returned if caller attempts to spend unconfirmed outputs
	SKY_ErrSpendingUnconfirmed
	// SKY_ErrInvalidEncryptedField is returned if a wallet's Meta.encrypted value is invalid.
	SKY_ErrInvalidEncryptedField
	// SKY_ErrWalletEncrypted is returned when trying to generate addresses or sign tx in encrypted wallet
	SKY_ErrWalletEncrypted
	// SKY_ErrWalletNotEncrypted is returned when trying to decrypt unencrypted wallet
	SKY_ErrWalletNotEncrypted
	// SKY_ErrWalletMissingPassword is returned when trying to create wallet with encryption, but password is not provided.
	SKY_ErrWalletMissingPassword
	// SKY_ErrMissingEncrypt is returned when trying to create wallet with password, but options.Encrypt is not set.
	SKY_ErrMissingEncrypt
	// SKY_ErrWalletInvalidPassword is returned if decrypts secrets failed
	SKY_ErrWalletInvalidPassword
	// SKY_ErrMissingSeed is returned when trying to create wallet without a seed
	SKY_ErrMissingSeed
	// SKY_ErrMissingAuthenticated is returned if try to decrypt a scrypt chacha20poly1305 encrypted wallet, and find no authenticated metadata.
	SKY_ErrMissingAuthenticated
	// SKY_ErrWrongCryptoType is returned when decrypting wallet with wrong crypto method
	SKY_ErrWrongCryptoType
	// SKY_ErrWalletNotExist is returned if a wallet does not exist
	SKY_ErrWalletNotExist
	// SKY_ErrSeedUsed is returned if a wallet already exists with the same seed
	SKY_ErrSeedUsed
	// SKY_ErrWalletAPIDisabled is returned when trying to do wallet actions while the EnableWalletAPI option is false
	SKY_ErrWalletAPIDisabled
	// SKY_ErrSeedAPIDisabled is returned when trying to get seed of wallet while the EnableWalletAPI or EnableSeedAPI is false
	SKY_ErrSeedAPIDisabled
	// SKY_ErrWalletNameConflict represents the wallet name conflict error
	SKY_ErrWalletNameConflict
	// SKY_ErrInvalidHoursSelectionMode for invalid HoursSelection mode values
	SKY_ErrInvalidHoursSelectionMode
	// SKY_ErrInvalidHoursSelectionType for invalid HoursSelection type values
	SKY_ErrInvalidHoursSelectionType
	// SKY_ErrUnknownAddress is returned if an address is not found in a wallet
	SKY_ErrUnknownAddress
	// SKY_ErrUnknownUxOut is returned if a uxout is not owned by any address in a wallet
	SKY_ErrUnknownUxOut
	// SKY_ErrNoUnspents is returned if a wallet has no unspents to spend
	SKY_ErrNoUnspents
	// SKY_ErrNullChangeAddress ChangeAddress must not be the null address
	SKY_ErrNullChangeAddress
	// SKY_ErrMissingTo To is required
	SKY_ErrMissingTo
	// SKY_ErrZeroCoinsTo To.Coins must not be zero
	SKY_ErrZeroCoinsTo
	// SKY_ErrNullAddressTo To.Address must not be the null address
	SKY_ErrNullAddressTo
	// SKY_ErrDuplicateTo To contains duplicate values
	SKY_ErrDuplicateTo
	// SKY_ErrMissingWalletID Wallet.ID is required
	SKY_ErrMissingWalletID
	// SKY_ErrIncludesNullAddress Wallet.Addresses must not contain the null address
	SKY_ErrIncludesNullAddress
	// SKY_ErrDuplicateAddresses Wallet.Addresses contains duplicate values
	SKY_ErrDuplicateAddresses
	// SKY_ErrZeroToHoursAuto To.Hours must be zero for auto type hours selection
	SKY_ErrZeroToHoursAuto
	// SKY_ErrMissingModeAuto HoursSelection.Mode is required for auto type hours selection
	SKY_ErrMissingModeAuto
	// SKY_ErrInvalidHoursSelMode Invalid HoursSelection.Mode
	SKY_ErrInvalidHoursSelMode
	// SKY_ErrInvalidModeManual HoursSelection.Mode cannot be used for manual type hours selection
	SKY_ErrInvalidModeManual
	// SKY_ErrInvalidHoursSelType Invalid HoursSelection.Type
	SKY_ErrInvalidHoursSelType
	// SKY_ErrMissingShareFactor HoursSelection.ShareFactor must be set for share mode
	SKY_ErrMissingShareFactor
	// SKY_ErrInvalidShareFactor HoursSelection.ShareFactor can only be used for share mode
	SKY_ErrInvalidShareFactor
	// SKY_ErrShareFactorOutOfRange HoursSelection.ShareFactor must be >= 0 and <= 1
	SKY_ErrShareFactorOutOfRange
	// SKY_ErrWalletParamsConflict Wallet.UxOuts and Wallet.Addresses cannot be combined
	SKY_ErrWalletParamsConflict
	// SKY_ErrDuplicateUxOuts Wallet.UxOuts contains duplicate values
	SKY_ErrDuplicateUxOuts
	// SKY_ErrUnknownWalletID params.Wallet.ID does not match wallet
	SKY_ErrUnknownWalletID
	// SKY_ErrVerifySignatureInvalidInputsNils VerifySignature, ERROR: invalid input, nils
	SKY_ErrVerifySignatureInvalidInputsNils
	// SKY_ErrVerifySignatureInvalidSigLength
	SKY_ErrVerifySignatureInvalidSigLength
	// SKY_ErrVerifySignatureInvalidPubkeysLength
	SKY_ErrVerifySignatureInvalidPubkeysLength
)

// Error codes defined in params package
// nolint megacheck
const (
	// SKY_ErrInvalidDecimals is returned by DropletPrecisionCheck if a coin amount has an invalid number of decimal places
	SKY_ErrInvalidDecimals = SKY_PKG_PARAMS + iota
)

var (
	// ErrorBadHandle invalid handle value
	ErrorBadHandle = errors.New("Invalid or unknown handle value")
	// ErrorUnknown unexpected error
	ErrorUnknown = errors.New("Unexpected error")
	// ErrorInvalidTimeString time string does not match expected time format
	// More precise errors conditions can be found in the logs
	ErrorInvalidTimeString = errors.New("Invalid time value")

	codeToErrorMap = make(map[uint32]error)
	errorToCodeMap = map[error]uint32{
		// libcgo
		ErrorBadHandle:         SKY_BAD_HANDLE,
		ErrorUnknown:           SKY_ERROR,
		ErrorInvalidTimeString: SKY_INVALID_TIMESTRING,
		// cipher
		cipher.ErrAddressInvalidLength:     SKY_ErrAddressInvalidLength,
		cipher.ErrAddressInvalidChecksum:   SKY_ErrAddressInvalidChecksum,
		cipher.ErrAddressInvalidVersion:    SKY_ErrAddressInvalidVersion,
		cipher.ErrAddressInvalidPubKey:     SKY_ErrAddressInvalidPubKey,
		cipher.ErrAddressInvalidFirstByte:  SKY_ErrAddressInvalidFirstByte,
		cipher.ErrAddressInvalidLastByte:   SKY_ErrAddressInvalidLastByte,
		encoder.ErrBufferUnderflow:         SKY_ErrBufferUnderflow,
		encoder.ErrInvalidOmitEmpty:        SKY_ErrInvalidOmitEmpty,
		cipher.ErrInvalidLengthPubKey:      SKY_ErrInvalidLengthPubKey,
		cipher.ErrPubKeyFromNullSecKey:     SKY_ErrPubKeyFromNullSecKey,
		cipher.ErrPubKeyFromBadSecKey:      SKY_ErrPubKeyFromBadSecKey,
		cipher.ErrInvalidLengthSecKey:      SKY_ErrInvalidLengthSecKey,
		cipher.ErrECHDInvalidPubKey:        SKY_ErrECHDInvalidPubKey,
		cipher.ErrECHDInvalidSecKey:        SKY_ErrECHDInvalidSecKey,
		cipher.ErrInvalidLengthSig:         SKY_ErrInvalidLengthSig,
		cipher.ErrInvalidLengthRipemd160:   SKY_ErrInvalidLengthRipemd160,
		cipher.ErrInvalidLengthSHA256:      SKY_ErrInvalidLengthSHA256,
		base58.ErrInvalidBase58Char:        SKY_ErrInvalidBase58Char,
		base58.ErrInvalidBase58String:      SKY_ErrInvalidBase58String,
		base58.ErrInvalidBase58Length:      SKY_ErrInvalidBase58Length,
		cipher.ErrInvalidHexLength:         SKY_ErrInvalidHexLength,
		cipher.ErrInvalidBytesLength:       SKY_ErrInvalidBytesLength,
		cipher.ErrInvalidPubKey:            SKY_ErrInvalidPubKey,
		cipher.ErrInvalidSecKey:            SKY_ErrInvalidSecKey,
		cipher.ErrInvalidSigPubKeyRecovery: SKY_ErrInvalidSigPubKeyRecovery,
		// Removed in ea0aafbffb76
		// cipher.ErrInvalidSecKeyHex:               SKY_ErrInvalidSecKeyHex,
		cipher.ErrInvalidAddressForSig:           SKY_ErrInvalidAddressForSig,
		cipher.ErrInvalidHashForSig:              SKY_ErrInvalidHashForSig,
		cipher.ErrPubKeyRecoverMismatch:          SKY_ErrPubKeyRecoverMismatch,
		cipher.ErrInvalidSigInvalidPubKey:        SKY_ErrInvalidSigInvalidPubKey,
		cipher.ErrInvalidSigValidity:             SKY_ErrInvalidSigValidity,
		cipher.ErrInvalidSigForMessage:           SKY_ErrInvalidSigForMessage,
		cipher.ErrInvalidSecKyVerification:       SKY_ErrInvalidSecKyVerification,
		cipher.ErrNullPubKeyFromSecKey:           SKY_ErrNullPubKeyFromSecKey,
		cipher.ErrInvalidDerivedPubKeyFromSecKey: SKY_ErrInvalidDerivedPubKeyFromSecKey,
		cipher.ErrInvalidPubKeyFromHash:          SKY_ErrInvalidPubKeyFromHash,
		cipher.ErrPubKeyFromSecKeyMismatch:       SKY_ErrPubKeyFromSecKeyMismatch,
		cipher.ErrInvalidLength:                  SKY_ErrInvalidLength,
		cipher.ErrBitcoinWIFInvalidFirstByte:     SKY_ErrBitcoinWIFInvalidFirstByte,
		cipher.ErrBitcoinWIFInvalidSuffix:        SKY_ErrBitcoinWIFInvalidSuffix,
		cipher.ErrBitcoinWIFInvalidChecksum:      SKY_ErrBitcoinWIFInvalidChecksum,
		cipher.ErrEmptySeed:                      SKY_ErrEmptySeed,
		cipher.ErrInvalidSig:                     SKY_ErrInvalidSig,
		encrypt.ErrMissingPassword:               SKY_ErrMissingPassword,
		encrypt.ErrDataTooLarge:                  SKY_ErrDataTooLarge,
		encrypt.ErrInvalidChecksumLength:         SKY_ErrInvalidChecksumLength,
		encrypt.ErrInvalidChecksum:               SKY_ErrInvalidChecksum,
		encrypt.ErrInvalidNonceLength:            SKY_ErrInvalidNonceLength,
		encrypt.ErrInvalidBlockSize:              SKY_ErrInvalidBlockSize,
		encrypt.ErrReadDataHashFailed:            SKY_ErrReadDataHashFailed,
		encrypt.ErrInvalidPassword:               SKY_ErrInvalidPassword,
		encrypt.ErrReadDataLengthFailed:          SKY_ErrReadDataLengthFailed,
		encrypt.ErrInvalidDataLength:             SKY_ErrInvalidDataLength,

		// cli
		cli.ErrTemporaryInsufficientBalance: SKY_ErrTemporaryInsufficientBalance,
		cli.ErrAddress:                      SKY_ErrAddress,
		cli.ErrWalletName:                   SKY_ErrWalletName,
		cli.ErrJSONMarshal:                  SKY_ErrJSONMarshal,
		// coin
		coin.ErrAddEarnedCoinHoursAdditionOverflow: SKY_ErrAddEarnedCoinHoursAdditionOverflow,
		coin.ErrUint64MultOverflow:                 SKY_ErrUint64MultOverflow,
		coin.ErrUint64AddOverflow:                  SKY_ErrUint64AddOverflow,
		coin.ErrUint32AddOverflow:                  SKY_ErrUint32AddOverflow,
		coin.ErrUint64OverflowsInt64:               SKY_ErrUint64OverflowsInt64,
		coin.ErrInt64UnderflowsUint64:              SKY_ErrInt64UnderflowsUint64,
		coin.ErrIntUnderflowsUint32:                SKY_ErrIntUnderflowsUint32,
		coin.ErrIntOverflowsUint32:                 SKY_ErrIntOverflowsUint32,
		// daemon
		// Removed in 34ad39ddb350
		// gnet.ErrMaxDefaultConnectionsReached:           SKY_ErrMaxDefaultConnectionsReached,
		pex.ErrPeerlistFull:       SKY_ErrPeerlistFull,
		pex.ErrInvalidAddress:     SKY_ErrInvalidAddress,
		pex.ErrNoLocalhost:        SKY_ErrNoLocalhost,
		pex.ErrNotExternalIP:      SKY_ErrNotExternalIP,
		pex.ErrPortTooLow:         SKY_ErrPortTooLow,
		pex.ErrBlacklistedAddress: SKY_ErrBlacklistedAddress,
		// gnet.ErrDisconnectReadFailed:              SKY_ErrDisconnectReadFailed,
		// gnet.ErrDisconnectWriteFailed:             SKY_ErrDisconnectWriteFailed,
		gnet.ErrDisconnectSetReadDeadlineFailed: SKY_ErrDisconnectSetReadDeadlineFailed,
		gnet.ErrDisconnectInvalidMessageLength:  SKY_ErrDisconnectInvalidMessageLength,
		gnet.ErrDisconnectMalformedMessage:      SKY_ErrDisconnectMalformedMessage,
		gnet.ErrDisconnectUnknownMessage:        SKY_ErrDisconnectUnknownMessage,
		gnet.ErrConnectionPoolClosed:            SKY_ErrConnectionPoolClosed,
		gnet.ErrWriteQueueFull:                  SKY_ErrWriteQueueFull,
		gnet.ErrNoReachableConnections:          SKY_ErrNoReachableConnections,
		daemon.ErrDisconnectVersionNotSupported: SKY_ErrDisconnectVersionNotSupported,
		daemon.ErrDisconnectIntroductionTimeout: SKY_ErrDisconnectIntroductionTimeout,
		daemon.ErrDisconnectIsBlacklisted:       SKY_ErrDisconnectIsBlacklisted,
		daemon.ErrDisconnectSelf:                SKY_ErrDisconnectSelf,
		daemon.ErrDisconnectConnectedTwice:      SKY_ErrDisconnectConnectedTwice,
		daemon.ErrDisconnectIdle:                SKY_ErrDisconnectIdle,
		daemon.ErrDisconnectNoIntroduction:      SKY_ErrDisconnectNoIntroduction,
		daemon.ErrDisconnectIPLimitReached:      SKY_ErrDisconnectIPLimitReached,
		// Removed
		//		daemon.ErrDisconnectMaxDefaultConnectionReached:   SKY_ErrDisconnectMaxDefaultConnectionReached,
		daemon.ErrDisconnectMaxOutgoingConnectionsReached: SKY_ErrDisconnectMaxOutgoingConnectionsReached,
		// util
		fee.ErrTxnNoFee:                 SKY_ErrTxnNoFee,
		fee.ErrTxnInsufficientFee:       SKY_ErrTxnInsufficientFee,
		fee.ErrTxnInsufficientCoinHours: SKY_ErrTxnInsufficientCoinHours,
		droplet.ErrNegativeValue:        SKY_ErrNegativeValue,
		droplet.ErrTooManyDecimals:      SKY_ErrTooManyDecimals,
		droplet.ErrTooLarge:             SKY_ErrTooLarge,
		file.ErrEmptyDirectoryName:      SKY_ErrEmptyDirectoryName,
		file.ErrDotDirectoryName:        SKY_ErrDotDirectoryName,
		// visor
		blockdb.ErrNoHeadBlock: SKY_ErrNoHeadBlock,
		visor.ErrVerifyStopped: SKY_ErrVerifyStopped,
		// wallet
		wallet.ErrInsufficientBalance:       SKY_ErrInsufficientBalance,
		wallet.ErrInsufficientHours:         SKY_ErrInsufficientHours,
		wallet.ErrZeroSpend:                 SKY_ErrZeroSpend,
		wallet.ErrSpendingUnconfirmed:       SKY_ErrSpendingUnconfirmed,
		wallet.ErrInvalidEncryptedField:     SKY_ErrInvalidEncryptedField,
		wallet.ErrWalletEncrypted:           SKY_ErrWalletEncrypted,
		wallet.ErrWalletNotEncrypted:        SKY_ErrWalletNotEncrypted,
		wallet.ErrMissingPassword:           SKY_ErrWalletMissingPassword,
		wallet.ErrMissingEncrypt:            SKY_ErrMissingEncrypt,
		wallet.ErrInvalidPassword:           SKY_ErrWalletInvalidPassword,
		wallet.ErrMissingSeed:               SKY_ErrMissingSeed,
		wallet.ErrMissingAuthenticated:      SKY_ErrMissingAuthenticated,
		wallet.ErrWrongCryptoType:           SKY_ErrWrongCryptoType,
		wallet.ErrWalletNotExist:            SKY_ErrWalletNotExist,
		wallet.ErrSeedUsed:                  SKY_ErrSeedUsed,
		wallet.ErrWalletAPIDisabled:         SKY_ErrWalletAPIDisabled,
		wallet.ErrSeedAPIDisabled:           SKY_ErrSeedAPIDisabled,
		wallet.ErrWalletNameConflict:        SKY_ErrWalletNameConflict,
		wallet.ErrInvalidHoursSelectionMode: SKY_ErrInvalidHoursSelectionMode,
		wallet.ErrInvalidHoursSelectionType: SKY_ErrInvalidHoursSelectionType,
		wallet.ErrUnknownAddress:            SKY_ErrUnknownAddress,
		wallet.ErrUnknownUxOut:              SKY_ErrUnknownUxOut,
		wallet.ErrNoUnspents:                SKY_ErrNoUnspents,
		wallet.ErrNullChangeAddress:         SKY_ErrNullChangeAddress,
		wallet.ErrMissingTo:                 SKY_ErrMissingTo,
		wallet.ErrZeroCoinsTo:               SKY_ErrZeroCoinsTo,
		wallet.ErrNullAddressTo:             SKY_ErrNullAddressTo,
		wallet.ErrDuplicateTo:               SKY_ErrDuplicateTo,
		wallet.ErrMissingWalletID:           SKY_ErrMissingWalletID,
		wallet.ErrIncludesNullAddress:       SKY_ErrIncludesNullAddress,
		wallet.ErrDuplicateAddresses:        SKY_ErrDuplicateAddresses,
		wallet.ErrZeroToHoursAuto:           SKY_ErrZeroToHoursAuto,
		wallet.ErrMissingModeAuto:           SKY_ErrMissingModeAuto,
		wallet.ErrInvalidHoursSelMode:       SKY_ErrInvalidHoursSelMode,
		wallet.ErrInvalidModeManual:         SKY_ErrInvalidModeManual,
		wallet.ErrInvalidHoursSelType:       SKY_ErrInvalidHoursSelType,
		wallet.ErrMissingShareFactor:        SKY_ErrMissingShareFactor,
		wallet.ErrInvalidShareFactor:        SKY_ErrInvalidShareFactor,
		wallet.ErrShareFactorOutOfRange:     SKY_ErrShareFactorOutOfRange,
		wallet.ErrWalletParamsConflict:      SKY_ErrWalletParamsConflict,
		wallet.ErrDuplicateUxOuts:           SKY_ErrDuplicateUxOuts,
		wallet.ErrUnknownWalletID:           SKY_ErrUnknownWalletID,
		// params
		params.ErrInvalidDecimals: SKY_ErrInvalidDecimals,
	}
)

func libErrorCode(err error) uint32 {
	if err == nil {
		return SKY_OK
	}
	if errcode, isKnownError := errorToCodeMap[err]; isKnownError {
		return errcode
	}
	switch err.(type) {
	case cli.WalletLoadError:
		return SKY_WalletLoadError
	case cli.WalletSaveError:
		return SKY_WalletSaveError
		//	case daemon.ConnectionError:
		//		return SKY_ConnectionError
	case historydb.ErrHistoryDBCorrupted:
		return SKY_ErrHistoryDBCorrupted
	case historydb.ErrUxOutNotExist:
		return SKY_ErrUxOutNotExist
	case blockdb.ErrUnspentNotExist:
		return SKY_ErrUnspentNotExist
	case blockdb.ErrMissingSignature:
		return SKY_ErrMissingSignature
	case dbutil.ErrCreateBucketFailed:
		return SKY_ErrCreateBucketFailed
	case dbutil.ErrBucketNotExist:
		return SKY_ErrBucketNotExist
	case visor.ErrTxnViolatesHardConstraint:
		return SKY_ErrTxnViolatesHardConstraint
	case visor.ErrTxnViolatesSoftConstraint:
		return SKY_ErrTxnViolatesSoftConstraint
	case visor.ErrTxnViolatesUserConstraint:
		return SKY_ErrTxnViolatesUserConstraint
	}
	return SKY_ERROR
}

func errorFromLibCode(errcode uint32) error {
	if err, exists := codeToErrorMap[errcode]; exists {
		return err
	}

	// FIXME: Be more specific and encode type, sub-error in error code
	err := errors.New("libskycoin error")
	if errcode == SKY_WalletLoadError {
		return cli.WalletLoadError{}
	}
	if errcode == SKY_WalletSaveError {
		return cli.WalletSaveError{}
	}
	if errcode == SKY_ErrHistoryDBCorrupted {
		return historydb.NewErrHistoryDBCorrupted(err)
	}
	if errcode == SKY_ErrUxOutNotExist {
		return historydb.ErrUxOutNotExist{UxID: ""}
	}
	if errcode == SKY_ErrUnspentNotExist {
		return blockdb.ErrUnspentNotExist{UxID: ""}
	}
	if errcode == SKY_ErrMissingSignature {
		return blockdb.NewErrMissingSignature(nil)
	}
	if errcode == SKY_ErrCreateBucketFailed {
		return dbutil.ErrCreateBucketFailed{Bucket: "", Err: nil}
	}
	if errcode == SKY_ErrBucketNotExist {
		return dbutil.ErrBucketNotExist{Bucket: ""}
	}
	if errcode == SKY_ErrTxnViolatesHardConstraint {
		return visor.ErrTxnViolatesHardConstraint{Err: err}
	}
	if errcode == SKY_ErrTxnViolatesSoftConstraint {
		return visor.ErrTxnViolatesSoftConstraint{Err: err}
	}
	if errcode == SKY_ErrTxnViolatesUserConstraint {
		return visor.ErrTxnViolatesUserConstraint{Err: err}
	}
	return nil
}

func init() {
	// Init reverse error code map
	for _err := range errorToCodeMap {
		codeToErrorMap[errorToCodeMap[_err]] = _err
	}
}
