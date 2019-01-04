
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

typedef struct{
	cipher_SHA256* data;
	int count;
} cipher_SHA256s;

typedef struct{
	coin__UxOut* data;
	int count;
} coin_UxOutArray;

typedef struct{
	cipher__Address* data;
	int count;
} cipher_Addresses;

typedef GoUint32_ (*FeeCalcFunc)(Transaction__Handle handle, unsigned long long * pFee, void* context);

typedef struct {
  FeeCalcFunc callback;
  void* context;
} Fee_Calculator ;