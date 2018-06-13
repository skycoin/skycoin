
typedef struct{
	GoUint8 data[33];
} cipher_PubKey;

typedef struct{
	GoUint8 data[32];
} cipher_SecKey;

typedef struct{
	GoUint8 data[20];
} cipher_Ripemd160;

typedef struct{
	GoUint8 data[65];
} cipher_Sig;

typedef struct{
	GoUint8 data[32];
} cipher_SHA256;

typedef struct{
	GoUint8 data[4];
} cipher_Checksum;

typedef struct{
	cipher_SecKey* data;
	int count;
} cipher_SecKeys;

typedef struct{
	cipher_PubKey* data;
	int count;
} cipher_PubKeys;

