package cipher

// Encryptor is the interface that provides encryption method
type Encryptor interface {
	// Encrypt encrypts data with password
	Encrypt(data, password []byte) ([]byte, error)
}

// Decryptor is the interface that provides decryption method
type Decryptor interface {
	// Decrypt decrypts the data with password
	Decrypt(data, password []byte) ([]byte, error)
}

// EncryptorDecryptor is the interface that provides encrypt and decrypt methods
type EncryptorDecryptor interface {
	Encryptor
	Decryptor
}
