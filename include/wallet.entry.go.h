
/**
 * Wallet entry.
 */
typedef struct {
	cipher__Address Address;    ///< Wallet address.
	cipher__PubKey  Public;     ///< Public key used to generate address.
	cipher__SecKey  Secret;     ///< Secret key used to generate address.
} wallet__Entry;
