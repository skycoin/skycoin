/**
 * Addresses of Bitcoin accounts
 */
typedef struct {
  GoUint8_ Version;  ///< Address version identifier.
                          ///< Used to differentiate testnet
                          ///< vs mainnet addresses, for instance.
  cipher__Ripemd160 Key;   ///< Address hash identifier.
} cipher__BitcoinAddress;
