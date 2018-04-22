/**
 * Wallet entry.
 */
typedef struct {
	cipher_Address Address;    ///< Wallet address.
	cipher_PubKey  Public;     ///< Public key used to generate address.
	cipher_SecKey  Secret;     ///< Secret key used to generate address.
} wallet_Entry;
