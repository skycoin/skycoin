// +build cgo

int xor_key_stream(
	unsigned char *out,
	const unsigned char *in,
	unsigned long long inlen,
	const unsigned char *n,
	const unsigned char *k,
	const unsigned int c[],
	const unsigned int rounds
);