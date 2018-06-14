/*cipher_SecKeys* parameter to return as a list */
%typemap(in, numinputs=0) (cipher_SecKeys*  __out_secKeys) (cipher_SecKeys temp) {
	temp.data = NULL;
	temp.count = 0;
	$1 = &temp;
}

/*cipher_SecKeys* as function return typemap*/
%typemap(argout) (cipher_SecKeys* __out_secKeys) {
	int i;
	PyObject *list = PyList_New($1->count);
	for (i = 0; i < $1->count; i++) {
		cipher_SecKey* key = &($1->data[i]);
		PyObject *o = SWIG_NewPointerObj(SWIG_as_voidptr(key), SWIGTYPE_p_cipher_SecKey, SWIG_POINTER_NOSHADOW |  0 );
		PyList_SetItem(list,i,o);
	}
	if( $1->data != NULL)
		free( (void*)$1->data );
	%append_output( list );
}


%rename(SKY_cipher_GenerateDeterministicKeyPairs) wrap_SKY_cipher_GenerateDeterministicKeyPairs;
%inline {
	GoUint32 wrap_SKY_cipher_GenerateDeterministicKeyPairs(GoSlice seed, GoInt n, cipher_SecKeys* __out_secKeys){
		GoSlice_ data;
		data.data = NULL;
		data.len = 0;
		data.cap = 0;
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
		GoSlice_ data;
		data.data = NULL;
		data.len = 0;
		data.cap = 0;
		GoUint32 result = SKY_cipher_GenerateDeterministicKeyPairsSeed(seed, n, newSeed, &data);
		if( result == 0){
			__out_secKeys->data = data.data;
			__out_secKeys->count = data.len;
		}
		return result;
	}
}

