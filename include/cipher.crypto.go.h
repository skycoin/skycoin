/**
 * Hash signed using a secret key, 65 bytes long.
 */
typedef unsigned char cipher_Sig[65];

/**
 * Public key, 33-bytes long.
 */
typedef unsigned char cipher_PubKey[33];

/**
 * Container type suitable for storing a variable number of
 * public keys.
 */
typedef GoSlice_ cipher_PubKeySlice;

/**
 * Secret key, 32 bytes long.
 */
typedef unsigned char cipher_SecKey[32];
