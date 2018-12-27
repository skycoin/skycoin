/**
 * Integrity checksum, 4-bytes long.
 */
typedef GoUint8_  cipher__Checksum[4];

/**
 * Addresses of SKY accounts
 */
typedef struct{
    GoUint8_ Version;      ///< Address version identifier.
               ///< Used to differentiate testnet
                           ///< vs mainnet addresses, for ins
    cipher__Ripemd160 Key; ///< Address hash identifier.
} cipher__Address;
