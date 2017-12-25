// +build ignore
// NOTE: the build directive was originally "+build cgo", but this code is not meant to be built
// for skycoin and exists only as a reference implementation

int xor_key_stream(
	unsigned char *out,
	const unsigned char *in,
	unsigned long long inlen,
	const unsigned char *n,
	const unsigned char *k,
	const unsigned int c[],
	const unsigned int rounds
);
