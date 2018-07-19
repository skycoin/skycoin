/*cipher_SecKeys* parameter to return as a list */
%typemap(in, numinputs=0) (cipher_SecKeys*  __out_secKeys) (cipher_SecKeys temp) {
	temp.data = NULL;
	temp.count = 0;
	$1 = &temp;
}

/*cipher_SecKeys* as function return typemap*/
%typemap(argout) (cipher_SecKeys* __out_secKeys) {
	int i;
	PyObject *list = PyList_New(0);
	for (i = 0; i < $1->count; i++) {
		cipher_SecKey* key = &($1->data[i]);
		cipher_SecKey* newKey = malloc(sizeof(cipher_SecKey));
		memcpy(newKey, key, sizeof(cipher_SecKey));
		PyObject *o = SWIG_NewPointerObj(SWIG_as_voidptr(newKey), SWIGTYPE_p_cipher_SecKey, SWIG_POINTER_OWN );
		PyList_Append(list, o);
		//Py_DECREF(o);
	}
	if( $1->data != NULL)
		free( (void*)$1->data );
	%append_output( list );
}


%rename(SKY_cipher_GenerateDeterministicKeyPairs) wrap_SKY_cipher_GenerateDeterministicKeyPairs;
%inline {
	GoUint32 wrap_SKY_cipher_GenerateDeterministicKeyPairs(GoSlice seed, GoInt n, cipher_SecKeys* __out_secKeys){
		__out_secKeys->data = NULL;
		__out_secKeys->count = 0;
		GoSlice_ data;
		data.data = malloc(sizeof(cipher_SecKey) * n);
		data.len = n;
		data.cap = n;
		GoUint32 result = SKY_cipher_GenerateDeterministicKeyPairs(seed, n, &data);
		if( result == 0){
			__out_secKeys->data = data.data;
			__out_secKeys->count = data.len;
		}
		return result;
	}
}

%inline {
	GoUint32 wrap_SKY_cipher_GenerateDeterministicKeyPairsSeed(GoSlice seed, GoInt n, coin__UxArray* newSeed, cipher_SecKeys* __out_secKeys){
		__out_secKeys->data = NULL;
		__out_secKeys->count = 0;
		GoSlice_ data;
		data.data = malloc(sizeof(cipher_SecKey) * n);
		data.len = n;
		data.cap = n;
		GoUint32 result = SKY_cipher_GenerateDeterministicKeyPairsSeed(seed, n, newSeed, &data);
		if( result == 0){
			__out_secKeys->data = data.data;
			__out_secKeys->count = data.len;
		}
		return result;
	}
}

