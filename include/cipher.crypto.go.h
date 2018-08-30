/**
 * Hash signed using a secret key, 65 bytes long.
 */
typedef GoUint8_ cipher__Sig[65];

/**
 * Public key, 33-bytes long.
 */
typedef GoUint8_ cipher__PubKey[33];

/**
 * Container type suitable for storing a variable number of
 * public keys.
 */
typedef GoSlice_ cipher__PubKeySlice;

/**
 * Secret key, 32 bytes long.
 */
typedef GoUint8_ cipher__SecKey[32];
