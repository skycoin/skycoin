
%inline {
	void  __copyToSecKeys(cipher_SecKeys* secKeys, cipher__SecKey* data, int count){
		secKeys->count = count;
		for(int i = 0; i < count; i++){
			memcpy( secKeys->data[i].data, data, sizeof(secKeys->data[i].data) );
			data++;
		}
	}
	
	void  __copyToPubKeys(cipher_PubKeys* pubKeys, cipher__PubKey* data, int count){
		pubKeys->count = count;
		for(int i = 0; i < count; i++){
			memcpy( pubKeys->data[i].data, data, sizeof(pubKeys->data[i].data) );
			data++;
		}
	}
	
	void __copyFromPubKeys(cipher_PubKeys* pubKeys, GoSlice_* slice){
		slice->data = malloc( pubKeys->count * 33 * sizeof(GoUint8) );
		slice->len = pubKeys->count;
		slice->cap = slice->len;
		cipher__PubKey* data = slice->data;
		for(int i = 0; i < slice->len; i++){
			memcpy( data, pubKeys->data[i].data, sizeof(pubKeys->data[i].data) );
			data++;
		}
	}
	
}

%rename(SKY_cipher_GenerateDeterministicKeyPairs) wrap_SKY_cipher_GenerateDeterministicKeyPairs;
%inline {
	GoUint32 wrap_SKY_cipher_GenerateDeterministicKeyPairs(GoSlice seed, GoInt n, cipher_SecKeys* secKeys){
		if( n > MAX_ARRAY_LENGTH_WRAP )
			return -1;
		GoSlice_ data;
		data.data = NULL;
		data.len = 0;
		data.cap = 0;
		GoUint32 result = SKY_cipher_GenerateDeterministicKeyPairs(seed, n, &data);
		__copyToSecKeys(secKeys, data.data, data.len);
		free(data.data);
		return result;
	}
}




%rename(SKY_cipher_GenerateDeterministicKeyPairsSeed) wrap_SKY_cipher_GenerateDeterministicKeyPairsSeed;
%inline {
	GoUint32 wrap_SKY_cipher_GenerateDeterministicKeyPairsSeed(GoSlice seed, GoInt n, coin__UxArray* newSeed, cipher_SecKeys* secKeys){
		if( n > MAX_ARRAY_LENGTH_WRAP )
			return -1;
		GoSlice_ data;
		data.data = NULL;
		data.len = 0;
		data.cap = 0;
		GoUint32 result = SKY_cipher_GenerateDeterministicKeyPairsSeed(seed, n, newSeed, &data);
		__copyToSecKeys(secKeys, data.data, data.len);
		free(data.data);
		return result;
	}
}

%rename(SKY_cipher_PubKeySlice_Len) wrap_SKY_cipher_PubKeySlice_Len;
%inline {
	//[]PubKey
	GoUint32 wrap_SKY_cipher_PubKeySlice_Len(cipher_PubKeys* pubKeys){
		GoSlice_ data;
		__copyFromPubKeys(pubKeys, &data);
		GoUint32 result = SKY_cipher_PubKeySlice_Len(&data);
		free(data.data);
		return result;
	}
}

%rename(SKY_cipher_PubKeySlice_Less) wrap_SKY_cipher_PubKeySlice_Less;
%inline {
	GoUint32 wrap_SKY_cipher_PubKeySlice_Less(cipher_PubKeys* pubKeys, GoInt p1, GoInt p2){
		GoSlice_ data;
		__copyFromPubKeys(pubKeys, &data);
		GoUint32 result = SKY_cipher_PubKeySlice_Less(&data, p1, p2);
		free(data.data);
		return result;
	}
}

%rename(SKY_cipher_PubKeySlice_Swap) wrap_SKY_cipher_PubKeySlice_Swap;
%inline {
	GoUint32 wrap_SKY_cipher_PubKeySlice_Swap(cipher_PubKeys* pubKeys, GoInt p1, GoInt p2){
		GoSlice_ data;
		__copyFromPubKeys(pubKeys, &data);
		GoUint32 result = SKY_cipher_PubKeySlice_Swap(&data, p1, p2);
		free(data.data);
		return result;
	}
}


