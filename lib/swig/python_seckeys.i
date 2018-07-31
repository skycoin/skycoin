%typecheck(SWIG_TYPECHECK_STRING_ARRAY) cipher_SecKeys* {
  $1 = PyList_Check($input) ? 1 : 0;
}

/*cipher_SecKeys* input parameter */
%typemap(in) (cipher_SecKeys* __in_secKeys) (cipher_SecKeys temp) {
	int i;
	$1 = &temp;
	$1->count = PyList_Size($input);
	$1->data = malloc(sizeof(cipher_SecKey) * $1->count);
	cipher_SecKey* pdata = $1->data;
	for(i = 0; i < $1->count; i++){
		PyObject *o = PyList_GetItem($input, i);
		void *argp = 0;
		int res = SWIG_ConvertPtr(o, &argp, SWIGTYPE_p_cipher_SecKey, 0 | 0);
		if (!SWIG_IsOK(res))
			SWIG_exception_fail(SWIG_TypeError, "expecting type SecKey");
		cipher_SecKey* p = (cipher_SecKey*)argp;
		memcpy(pdata, p, sizeof(cipher_SecKey));
		pdata++;
	}
}

%typemap(freearg) (cipher_SecKeys* __in_secKeys) {
  if ($1->data) free($1->data);
}

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
		Py_DECREF(o);
	}
	if( $1->data != NULL)
		free( (void*)$1->data );
	%append_output( list );
}
