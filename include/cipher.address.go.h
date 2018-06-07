typedef GoUint8_  cipher__Checksum[4];
typedef struct{
    GoUint8_ Version;      ///< Address version identifier.
						   ///< Used to differentiate testnet
                           ///< vs mainnet addresses, for ins
    cipher__Ripemd160 Key; ///< Address hash identifier.
} cipher__Address;
