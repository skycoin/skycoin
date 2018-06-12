%extend cipher__Address {
	int isEqual(cipher__Address* a){
		if( $self-> Version == a->Version ){
			return memcmp($self->Key, a->Key, sizeof(a->Key)) == 0;
		}
		return 0;
	}
}

%extend cipher_SecKeys {
	cipher_SecKey* getAt(int i){
		if( i < $self->count ){
			return &$self->data[i];
		}
		else
			return NULL;
	}
}

typedef GoUint8_ cipher__PubKey[33];

typedef GoUint8_ cipher__Ripemd160[20];

typedef GoUint8_ cipher__SecKey[32];

typedef GoUint8_ cipher__Sig[65];

typedef GoUint8_ cipher__SHA256[32];

typedef GoUint8_  cipher__Checksum[4];

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
	cipher_SecKey data[MAX_ARRAY_LENGTH_WRAP];
	int count;
} cipher_SecKeys;

typedef struct{
	cipher_PubKey data[MAX_ARRAY_LENGTH_WRAP];
	int count;
} cipher_PubKeys;

typedef struct{
    GoUint8_ Version;      ///< Address version identifier.
						   ///< Used to differentiate testnet
                           ///< vs mainnet addresses, for ins
    cipher__Ripemd160 Key; ///< Address hash identifier.
} cipher__Address;
