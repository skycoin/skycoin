/**
 * Integrity checksum, 4-bytes long.
 */
typedef unsigned char cipher__Checksum[4];

/**
 * Addresses of SKY accounts
 */
typedef struct {
	unsigned char Version;  ///< Address version identifier.
                          ///< Used to differentiate testnet
                          ///< vs mainnet addresses, for instance.
	cipher__Ripemd160 Key;   ///< Address hash identifier.
} cipher__Address;
