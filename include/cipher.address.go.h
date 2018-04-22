/**
 * Integrity checksum, 4-bytes long.
 */
typedef unsigned char cipher_Checksum[4];

/**
 * Addresses of SKY accounts
 */
typedef struct {
	unsigned char Version;  ///< Address version identifier.
                          ///< Used to differentiate testnet
                          ///< vs mainnet addresses, for instance.
	cipher_Ripemd160 Key;   ///< Address hash identifier.
} cipher_Address;
