
#ifndef SKY_ERRORS_H
#define SKY_ERRORS_H

#define SKY_OK            0
#define SKY_ERROR         0xFFFFFFFF

#define SKY_PKG_API       0x01000000
#define SKY_PKG_CIPHER    0x02000001
#define SKY_PKG_CLI       0x03000002
#define SKY_PKG_COIN      0x04000003
#define SKY_PKG_CONSENSUS 0x05000004
#define SKY_PKG_DAEMON    0x06000005
#define SKY_PKG_GUI       0x07000006
#define SKY_PKG_SKYCOIN   0x08000007
#define SKY_PKG_UTIL      0x09000008
#define SKY_PKG_VISOR     0x0A000009
#define SKY_PKG_WALLET    0x0B00000A

#define SKY_ErrInvalidLength     0x02000000
#define SKY_ErrInvalidChecksum   0x02000001
#define SKY_ErrInvalidVersion    0x02000002
#define SKY_ErrInvalidPubKey     0x02000003
#define SKY_ErrInvalidFirstByte  0x02000004
#define SKY_ErrInvalidLastByte   0x02000005
#define SKY_ErrBufferUnderflow   0x02000006
#define SKY_ErrInvalidOmitEmpty  0x02000007

#define SKY_ErrTemporaryInsufficientBalance   0x03000000
#define SKY_ErrAddress                        0x03000001
#define SKY_ErrWalletName                     0x03000002
#define SKY_ErrJSONMarshal                    0x03000003
#define SKY_WalletLoadError                   0x03000004
#define SKY_WalletSaveError                   0x03000005

#define SKY_ErrAddEarnedCoinHoursAdditionOverflow 0x04000000
#define SKY_ErrUint64MultOverflow                 0x04000001
#define SKY_ErrUint64AddOverflow                  0x04000002
#define SKY_ErrUint32AddOverflow                  0x04000003
#define SKY_ErrUint64OverflowsInt64               0x04000004
#define SKY_ErrInt64UnderflowsUint64              0x04000005

#define SKY_ErrPeerlistFull                               0x06000000
#define SKY_ErrInvalidAddress                             0x06000001
#define SKY_ErrNoLocalhost                                0x06000002
#define SKY_ErrNotExternalIP                              0x06000003
#define SKY_ErrPortTooLow                                 0x06000004
#define SKY_ErrBlacklistedAddress                         0x06000005
#define SKY_ErrDisconnectReadFailed                       0x06000006
#define SKY_ErrDisconnectWriteFailed                      0x06000007
#define SKY_ErrDisconnectSetReadDeadlineFailed            0x06000008
#define SKY_ErrDisconnectInvalidMessageLength             0x06000009
#define SKY_ErrDisconnectMalformedMessage                 0x0600000A
#define SKY_ErrDisconnectUnknownMessage                   0x0600000B
#define SKY_ErrDisconnectUnexpectedError                  0x0600000C
#define SKY_ErrConnectionPoolClosed                       0x0600000D
#define SKY_ErrWriteQueueFull                             0x0600000E
#define SKY_ErrNoReachableConnections                     0x0600000F
#define SKY_ErrMaxDefaultConnectionsReached               0x06000010
#define SKY_ErrDisconnectInvalidVersion                   0x06000011
#define SKY_ErrDisconnectIntroductionTimeout              0x06000012
#define SKY_ErrDisconnectVersionSendFailed                0x06000013
#define SKY_ErrDisconnectIsBlacklisted                    0x06000014
#define SKY_ErrDisconnectSelf                             0x06000015
#define SKY_ErrDisconnectConnectedTwice                   0x06000016
#define SKY_ErrDisconnectIdle                             0x06000017
#define SKY_ErrDisconnectNoIntroduction                   0x06000018
#define SKY_ErrDisconnectIPLimitReached                   0x06000019
#define SKY_ErrDisconnectOtherError                       0x0600001A
#define SKY_ErrDisconnectMaxDefaultConnectionReached      0x0600001B
#define SKY_ErrDisconnectMaxOutgoingConnectionsReached    0x0600001C
#define SKY_ConnectionError                               0x0600001D

#define SKY_ErrTxnNoFee                   0x09000000 
#define SKY_ErrTxnInsufficientFee         0x09000001 
#define SKY_ErrTxnInsufficientCoinHours   0x09000002 
#define SKY_ErrNegativeValue              0x09000003 
#define SKY_ErrTooManyDecimals            0x09000004 
#define SKY_ErrTooLarge                   0x09000005 
#define SKY_ErrEmptyDirectoryName         0x09000006 
#define SKY_ErrDotDirectoryName           0x09000007 

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

#define SKY_ErrInsufficientBalance            0x0B000000
#define SKY_ErrInsufficientHours              0x0B000001
#define SKY_ErrZeroSpend                      0x0B000002
#define SKY_ErrSpendingUnconfirmed            0x0B000003
#define SKY_ErrInvalidEncryptedField          0x0B000004
#define SKY_ErrWalletEncrypted                0x0B000005
#define SKY_ErrWalletNotEncrypted             0x0B000006
#define SKY_ErrMissingPassword                0x0B000007
#define SKY_ErrMissingEncrypt                 0x0B000008
#define SKY_ErrInvalidPassword                0x0B000009
#define SKY_ErrMissingSeed                    0x0B00000A
#define SKY_ErrMissingAuthenticated           0x0B00000B
#define SKY_ErrWrongCryptoType                0x0B00000C
#define SKY_ErrWalletNotExist                 0x0B00000D
#define SKY_ErrSeedUsed                       0x0B00000E
#define SKY_ErrWalletAPIDisabled              0x0B00000F
#define SKY_ErrSeedAPIDisabled                0x0B000010
#define SKY_ErrWalletNameConflict             0x0B000011
#define SKY_ErrInvalidHoursSelectionMode      0x0B000012
#define SKY_ErrInvalidHoursSelectionType      0x0B000013
#define SKY_ErrUnknownAddress                 0x0B000014
#define SKY_ErrUnknownUxOut                   0x0B000015
#define SKY_ErrNoUnspents                     0x0B000016
#define SKY_ErrNullChangeAddress              0x0B000017
#define SKY_ErrMissingTo                      0x0B000018
#define SKY_ErrZeroCoinsTo                    0x0B000019
#define SKY_ErrNullAddressTo                  0x0B00001A
#define SKY_ErrDuplicateTo                    0x0B00001B
#define SKY_ErrMissingWalletID                0x0B00001C
#define SKY_ErrIncludesNullAddress            0x0B00001D
#define SKY_ErrDuplicateAddresses             0x0B00001E
#define SKY_ErrZeroToHoursAuto                0x0B00001F
#define SKY_ErrMissingModeAuto                0x0B000020
#define SKY_ErrInvalidHoursSelMode            0x0B000021
#define SKY_ErrInvalidModeManual              0x0B000022
#define SKY_ErrInvalidHoursSelType            0x0B000023
#define SKY_ErrMissingShareFactor             0x0B000024
#define SKY_ErrInvalidShareFactor             0x0B000025
#define SKY_ErrShareFactorOutOfRange          0x0B000026
#define SKY_ErrWalletConstraint               0x0B000027
#define SKY_ErrDuplicateUxOuts                0x0B000028
#define SKY_ErrUnknownWalletID                0x0B000029

#endif
