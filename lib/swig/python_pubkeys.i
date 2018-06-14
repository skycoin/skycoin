/*cipher_PubKeys* input parameter */
%typemap(in) (cipher_PubKeys* __in_pubKeys) {
	int i;
	$1->count = PyList_Size($input);
	$1->data = malloc(sizeof(cipher_PubKey) * $1->count);
	cipher_PubKey* pdata = $1->data;
	for(i = 0; i < $1->count; i++){
		PyObject *o = PyList_GetItem($input, i);
		void *argp = 0;
		int res = SWIG_ConvertPtr(o, &argp, SWIGTYPE_p_cipher_PubKey, 0 | 0);
		if (!SWIG_IsOK(res))
			SWIG_exception_fail(SWIG_TypeError, "expecting type PubKey");
		cipher_PubKey* p = (cipher_PubKey*)argp;
		memcpy(p, pdata, sizeof(cipher_PubKey));
		pdata++;
	}
}

%typemap(freearg) (cipher_PubKeys* __in_pubKeys) {
  if ($1->data) free($1->data);
}

%rename(SKY_cipher_PubKeySlice_Len) wrap_SKY_cipher_PubKeySlice_Len;
%inline {
	GoUint32 wrap_SKY_cipher_PubKeySlice_Len(cipher_PubKeys* __in_pubKeys){
		GoSlice_ data;
		data.data = __in_pubKeys->data;
		data.len = __in_pubKeys->count;
		data.cap = __in_pubKeys->count;
		GoUint32 result = SKY_cipher_PubKeySlice_Len(&data);
		return result;
	}
}

%rename(SKY_cipher_PubKeySlice_Less) wrap_SKY_cipher_PubKeySlice_Less;
%inline {
	GoUint32 wrap_SKY_cipher_PubKeySlice_Less(cipher_PubKeys* __in_pubKeys, GoInt p1, GoInt p2){
		GoSlice_ data;
		data.data = __in_pubKeys->data;
		data.len = __in_pubKeys->count;
		data.cap = __in_pubKeys->count;
		GoUint32 result = SKY_cipher_PubKeySlice_Less(&data, p1, p2);
		return result;
	}
}

%rename(SKY_cipher_PubKeySlice_Swap) wrap_SKY_cipher_PubKeySlice_Swap;
%inline {
	GoUint32 wrap_SKY_cipher_PubKeySlice_Swap(cipher_PubKeys* __in_pubKeys, GoInt p1, GoInt p2){
		GoSlice_ data;
		data.data = __in_pubKeys->data;
		data.len = __in_pubKeys->count;
		data.cap = __in_pubKeys->count;
		GoUint32 result = SKY_cipher_PubKeySlice_Swap(&data, p1, p2);
		return result;
	}
}

