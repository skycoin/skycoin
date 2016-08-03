package mesh

/*
Move to file
- Get working
- Fix key
*/

//TODO: doesnt need to exist as a struct, only needs function
type ChaChaCrypto struct {
	key [32]byte //key is not used
}

func (self *ChaChaCrypto) GetKey() []byte {
	return self.key[:] //key is not used
}

func (self *ChaChaCrypto) Encrypt(in []byte, peerKey []byte) []byte {
	out := make([]byte, len(in))
	//TODO, using fixed nonce
	chacha20.XORKeyStream(out, in, []byte("nonce123"), peerKey[:])
	return out
}

func (self *ChaChaCrypto) Decrypt(in []byte) []byte {
	out := make([]byte, len(in))
	//TODO, using fixed nonce
	chacha20.XORKeyStream(out, in, []byte("nonce123"), self.key[:])
	return out
}
