
%rename(SKY_cipher_GenerateDeterministicKeyPairs) wrap_SKY_cipher_GenerateDeterministicKeyPairs;
%inline {
	GoUint32 wrap_SKY_cipher_GenerateDeterministicKeyPairs(GoSlice seed, GoInt n, cipher_SecKeys* secKeys){
		if( n > MAX_ARRAY_LENGTH_WRAP )
			return -1;
		GoSlice_ data;
		data.data = secKeys->data;
		data.len = 0;
		data.cap = MAX_ARRAY_LENGTH_WRAP;
		GoUint32 result = SKY_cipher_GenerateDeterministicKeyPairs(seed, n, &data);
		return result;
	}
}

%rename(SKY_cipher_GenerateDeterministicKeyPairsSeed) wrap_SKY_cipher_GenerateDeterministicKeyPairsSeed;
%inline {
	GoUint32 wrap_SKY_cipher_GenerateDeterministicKeyPairsSeed(GoSlice seed, GoInt n, coin__UxArray* newSeed, cipher_SecKeys* secKeys){
		if( n > MAX_ARRAY_LENGTH_WRAP )
			return -1;
		GoSlice_ data;
		data.data = secKeys->data;
		data.len = 0;
		data.cap = MAX_ARRAY_LENGTH_WRAP;
		GoUint32 result = SKY_cipher_GenerateDeterministicKeyPairsSeed(seed, n, newSeed, &data);
		return result;
	}
}

%rename(SKY_cipher_PubKeySlice_Len) wrap_SKY_cipher_PubKeySlice_Len;
%inline {
	//[]PubKey
	GoUint32 wrap_SKY_cipher_PubKeySlice_Len(cipher_PubKeys* pubKeys){
		GoSlice_ data;
		data.data = secKeys->data;
		data.len = secKeys->count;
		data.cap = MAX_ARRAY_LENGTH_WRAP;
		GoUint32 result = SKY_cipher_PubKeySlice_Len(&data);
		return result;
	}
}

%rename(SKY_cipher_PubKeySlice_Less) wrap_SKY_cipher_PubKeySlice_Less;
%inline {
	GoUint32 wrap_SKY_cipher_PubKeySlice_Less(cipher_PubKeys* pubKeys, GoInt p1, GoInt p2){
		GoSlice_ data;
		data.data = secKeys->data;
		data.len = secKeys->count;
		data.cap = MAX_ARRAY_LENGTH_WRAP;
		GoUint32 result = SKY_cipher_PubKeySlice_Less(&data, p1, p2);
		return result;
	}
}

%rename(SKY_cipher_PubKeySlice_Swap) wrap_SKY_cipher_PubKeySlice_Swap;
%inline {
	GoUint32 wrap_SKY_cipher_PubKeySlice_Swap(cipher_PubKeys* pubKeys, GoInt p1, GoInt p2){
		GoSlice_ data;
		data.data = secKeys->data;
		data.len = secKeys->count;
		data.cap = MAX_ARRAY_LENGTH_WRAP;
		GoUint32 result = SKY_cipher_PubKeySlice_Swap(&data, p1, p2);
		return result;
	}
}
