#include <signal.h>

#if __APPLE__
#include "TargetConditionals.h"
#endif

#if __linux__
#define SKY_ABORT .signal = SIGABRT
#elif __APPLE__
#if TARGET_OS_MAC
#define SKY_ABORT .exit_code = 2
#endif
#endif

#ifndef SKY_ERRORS_H
#define SKY_ERRORS_H

// Generic error conditions
#define SKY_OK                    0
#define SKY_ERROR                 0x7FFFFFFF
#define SKY_BAD_HANDLE            0x7F000001
#define SKY_INVALID_TIMESTRING    0x7F000002

// Package error code prefix list
#define SKY_PKG_API       0x01000000
#define SKY_PKG_CIPHER    0x02000000
#define SKY_PKG_CLI       0x03000000
#define SKY_PKG_COIN      0x04000000
#define SKY_PKG_CONSENSUS 0x05000000
#define SKY_PKG_DAEMON    0x06000000
#define SKY_PKG_GUI       0x07000000
#define SKY_PKG_SKYCOIN   0x08000000
#define SKY_PKG_UTIL      0x09000000
#define SKY_PKG_VISOR     0x0A000000
#define SKY_PKG_WALLET    0x0B000000
#define SKY_PKG_PARAMS    0x0C000000
#define SKY_PKG_LIBCGO    0x7F000000

#define SKY_PKG_LIBCGO    0x7F000000

// libcgo error codes
#define SKY_BAD_HANDLE    0x7F000001

// cipher error codes
#define SKY_ErrAddressInvalidLength                 0x02000000
#define SKY_ErrAddressInvalidChecksum               0x02000001
#define SKY_ErrAddressInvalidVersion                0x02000002
#define SKY_ErrAddressInvalidPubKey                 0x02000003
#define SKY_ErrAddressInvalidFirstByte              0x02000004
#define SKY_ErrAddressInvalidLastByte               0x02000005
#define SKY_ErrBufferUnderflow                      0x02000006
#define SKY_ErrInvalidOmitEmpty                     0x02000007
#define SKY_ErrInvalidLengthPubKey                  0x02000008
#define SKY_ErrPubKeyFromNullSecKey                 0x02000009
#define SKY_ErrPubKeyFromBadSecKey                  0x0200000A
#define SKY_ErrInvalidLengthSecKey                  0x0200000B
#define SKY_ErrECHDInvalidPubKey                    0x0200000C
#define SKY_ErrECHDInvalidSecKey                    0x0200000D
#define SKY_ErrInvalidLengthSig                     0x0200000E
#define SKY_ErrInvalidLengthRipemd160               0x0200000F
#define SKY_ErrInvalidLengthSHA256                  0x02000010
#define SKY_ErrInvalidBase58Char                    0x02000011
#define SKY_ErrInvalidBase58String                  0x02000012
#define SKY_ErrInvalidBase58Length                  0x02000013
#define SKY_ErrInvalidHexLength                     0x02000014
#define SKY_ErrInvalidBytesLength                   0x02000015
#define SKY_ErrInvalidPubKey                        0x02000016
#define SKY_ErrInvalidSecKey                        0x02000017
#define SKY_ErrInvalidSigPubKeyRecovery             0x02000018
#define SKY_ErrInvalidSecKeyHex                     0x02000019
#define SKY_ErrInvalidAddressForSig                 0x0200001A
#define SKY_ErrInvalidHashForSig                    0x0200001B
#define SKY_ErrPubKeyRecoverMismatch                0x0200001C
#define SKY_ErrInvalidSigInvalidPubKey              0x0200001D
#define SKY_ErrInvalidSigValidity                   0x0200001E
#define SKY_ErrInvalidSigForMessage                 0x0200001F
#define SKY_ErrInvalidSecKyVerification             0x02000020
#define SKY_ErrNullPubKeyFromSecKey                 0x02000021
#define SKY_ErrInvalidDerivedPubKeyFromSecKey       0x02000022
#define SKY_ErrInvalidPubKeyFromHash                0x02000023
#define SKY_ErrPubKeyFromSecKeyMismatch             0x02000024
#define SKY_ErrInvalidLength                        0x02000025
#define SKY_ErrBitcoinWIFInvalidFirstByte           0x02000026
#define SKY_ErrBitcoinWIFInvalidSuffix              0x02000027
#define SKY_ErrBitcoinWIFInvalidChecksum            0x02000028
#define SKY_ErrEmptySeed                            0x02000029
#define SKY_ErrInvalidSig                           0x0200002A
#define SKY_ErrMissingPassword                      0x0200002B
#define SKY_ErrDataTooLarge                         0x0200002C
#define SKY_ErrInvalidChecksumLength                0x0200002D
#define SKY_ErrInvalidChecksum                      0x0200002E
#define SKY_ErrInvalidNonceLength                   0x0200002F
#define SKY_ErrInvalidBlockSize                     0x02000030
#define SKY_ErrReadDataHashFailed                   0x02000031
#define SKY_ErrInvalidPassword                      0x02000032
#define SKY_ErrReadDataLengthFailed                 0x02000033
#define SKY_ErrInvalidDataLength                    0x02000034

// cli error codes
#define SKY_ErrTemporaryInsufficientBalance   0x03000000
#define SKY_ErrAddress                        0x03000001
#define SKY_ErrWalletName                     0x03000002
#define SKY_ErrJSONMarshal                    0x03000003
#define SKY_WalletLoadError                   0x03000004
#define SKY_WalletSaveError                   0x03000005

// coin error codes
#define SKY_ErrAddEarnedCoinHoursAdditionOverflow 0x04000000
#define SKY_ErrUint64MultOverflow                 0x04000001
#define SKY_ErrUint64AddOverflow                  0x04000002
#define SKY_ErrUint32AddOverflow                  0x04000003
#define SKY_ErrUint64OverflowsInt64               0x04000004
#define SKY_ErrInt64UnderflowsUint64              0x04000005
#define SKY_ErrIntUnderflowsUint32                0x04000006
#define SKY_ErrIntOverflowsUint32                 0x04000007

// daemon error codes
#define SKY_ErrPeerlistFull                               0x06000000
#define SKY_ErrInvalidAddress                             0x06000001
#define SKY_ErrNoLocalhost                                0x06000002
#define SKY_ErrNotExternalIP                              0x06000003
#define SKY_ErrPortTooLow                                 0x06000004
#define SKY_ErrBlacklistedAddress                         0x06000005
// #define SKY_ErrDisconnectReadFailed                       0x06000006
#define SKY_ErrDisconnectWriteFailed                      0x06000007
#define SKY_ErrDisconnectSetReadDeadlineFailed            0x06000008
#define SKY_ErrDisconnectInvalidMessageLength             0x06000009
#define SKY_ErrDisconnectMalformedMessage                 0x0600000A
#define SKY_ErrDisconnectUnknownMessage                   0x0600000B
#define SKY_ErrConnectionPoolClosed                       0x0600000D
#define SKY_ErrWriteQueueFull                             0x0600000E
#define SKY_ErrNoReachableConnections                     0x0600000F
#define SKY_ErrMaxDefaultConnectionsReached               0x06000010
#define SKY_ErrDisconnectVersionNotSupported              0x06000011
#define SKY_ErrDisconnectIntroductionTimeout              0x06000012
#define SKY_ErrDisconnectIsBlacklisted                    0x06000014
#define SKY_ErrDisconnectSelf                             0x06000015
#define SKY_ErrDisconnectConnectedTwice                   0x06000016
#define SKY_ErrDisconnectIdle                             0x06000017
#define SKY_ErrDisconnectNoIntroduction                   0x06000018
#define SKY_ErrDisconnectIPLimitReached                   0x06000019
#define SKY_ErrDisconnectMaxDefaultConnectionReached      0x0600001B
#define SKY_ErrDisconnectMaxOutgoingConnectionsReached    0x0600001C
#define SKY_ConnectionError                               0x0600001D

// util error codes
#define SKY_ErrTxnNoFee                   0x09000000
#define SKY_ErrTxnInsufficientFee         0x09000001
#define SKY_ErrTxnInsufficientCoinHours   0x09000002
#define SKY_ErrNegativeValue              0x09000003
#define SKY_ErrTooManyDecimals            0x09000004
#define SKY_ErrTooLarge                   0x09000005
#define SKY_ErrEmptyDirectoryName         0x09000006
#define SKY_ErrDotDirectoryName           0x09000007

// visor error codes
#define SKY_ErrHistoryDBCorrupted         0x0A000000
#define SKY_ErrUxOutNotExist              0x0A000001
#define SKY_ErrNoHeadBlock                0x0A000002
#define SKY_ErrMissingSignature           0x0A000003
#define SKY_ErrUnspentNotExist            0x0A000004
#define SKY_ErrVerifyStopped              0x0A000005
#define SKY_ErrCreateBucketFailed         0x0A000000
#define SKY_ErrBucketNotExist             0x0A000006
#define SKY_ErrTxnViolatesHardConstraint  0x0A000007
#define SKY_ErrTxnViolatesSoftConstraint  0x0A000008
#define SKY_ErrTxnViolatesUserConstraint  0x0A000009

// wallet error codes
#define SKY_ErrInsufficientBalance                    0x0B000000
#define SKY_ErrInsufficientHours                      0x0B000001
#define SKY_ErrZeroSpend                              0x0B000002
#define SKY_ErrSpendingUnconfirmed                    0x0B000003
#define SKY_ErrInvalidEncryptedField                  0x0B000004
#define SKY_ErrWalletEncrypted                        0x0B000005
#define SKY_ErrWalletNotEncrypted                     0x0B000006
#define SKY_ErrWalletMissingPassword                  0x0B000007
#define SKY_ErrMissingEncrypt                         0x0B000008
#define SKY_ErrWalletInvalidPassword                  0x0B000009
#define SKY_ErrMissingSeed                            0x0B00000A
#define SKY_ErrMissingAuthenticated                   0x0B00000B
#define SKY_ErrWrongCryptoType                        0x0B00000C
#define SKY_ErrWalletNotExist                         0x0B00000D
#define SKY_ErrSeedUsed                               0x0B00000E
#define SKY_ErrWalletAPIDisabled                      0x0B00000F
#define SKY_ErrSeedAPIDisabled                        0x0B000010
#define SKY_ErrWalletNameConflict                     0x0B000011
#define SKY_ErrInvalidHoursSelectionMode              0x0B000012
#define SKY_ErrInvalidHoursSelectionType              0x0B000013
#define SKY_ErrUnknownAddress                         0x0B000014
#define SKY_ErrUnknownUxOut                           0x0B000015
#define SKY_ErrNoUnspents                             0x0B000016
#define SKY_ErrNullChangeAddress                      0x0B000017
#define SKY_ErrMissingTo                              0x0B000018
#define SKY_ErrZeroCoinsTo                            0x0B000019
#define SKY_ErrNullAddressTo                          0x0B00001A
#define SKY_ErrDuplicateTo                            0x0B00001B
#define SKY_ErrMissingWalletID                        0x0B00001C
#define SKY_ErrIncludesNullAddress                    0x0B00001D
#define SKY_ErrDuplicateAddresses                     0x0B00001E
#define SKY_ErrZeroToHoursAuto                        0x0B00001F
#define SKY_ErrMissingModeAuto                        0x0B000020
#define SKY_ErrInvalidHoursSelMode                    0x0B000021
#define SKY_ErrInvalidModeManual                      0x0B000022
#define SKY_ErrInvalidHoursSelType                    0x0B000023
#define SKY_ErrMissingShareFactor                     0x0B000024
#define SKY_ErrInvalidShareFactor                     0x0B000025
#define SKY_ErrShareFactorOutOfRange                  0x0B000026
#define SKY_ErrWalletConstraint                       0x0B000027
#define SKY_ErrDuplicateUxOuts                        0x0B000028
#define SKY_ErrUnknownWalletID                        0x0B000029
#define SKY_ErrVerifySignatureInvalidInputsNils       0x0B000033
#define SKY_ErrVerifySignatureInvalidSigLength        0x0B000034
#define SKY_ErrVerifySignatureInvalidPubkeysLength    0x0B000035

// daemon error codes
#define SKY_ErrInvalidDecimals                        0x0C000000

#endif
