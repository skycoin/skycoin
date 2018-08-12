%typecheck(SWIG_TYPECHECK_STRING_ARRAY) cipher_SHA256s* {
  $1 = PyList_Check($input) ? 1 : 0;
}

/*cipher_SHA256s* input parameter */
%typemap(in) (cipher_SHA256s* __in_hashes) (cipher_SHA256s temp) {
	int i;
	$1 = &temp;
	$1->count = PyList_Size($input);
	$1->data = malloc(sizeof(cipher_SHA256) * $1->count);
	cipher_SHA256* pdata = $1->data;
	for(i = 0; i < $1->count; i++){
		PyObject *o = PyList_GetItem($input, i);
		void *argp = 0;
		int res = SWIG_ConvertPtr(o, &argp, SWIGTYPE_p_cipher_SHA256, 0 | 0);
		if (!SWIG_IsOK(res))
			SWIG_exception_fail(SWIG_TypeError, "expecting type cipher_SHA256");
		cipher_SHA256* p = (cipher_SHA256*)argp;
		memcpy(pdata, p, sizeof(cipher_SHA256));
		pdata++;
	}
}

%typemap(freearg) (cipher_SHA256s* __in_hashes) {
  if ($1->data) free($1->data);
}

/*cipher_SHA256s* parameter to return as a list */
%typemap(in, numinputs=0) (cipher_SHA256s*  __out_hashes) (cipher_SHA256s temp) {
	temp.data = NULL;
	temp.count = 0;
	$1 = &temp;
}

/*cipher_SHA256s* as function return typemap*/
%typemap(argout) (cipher_SHA256s* __out_hashes) {
	int i;
	PyObject *list = PyList_New(0);
	for (i = 0; i < $1->count; i++) {
		cipher_SHA256* key = &($1->data[i]);
		cipher_SHA256* newKey = malloc(sizeof(cipher_SHA256));
		memcpy(newKey, key, sizeof(cipher_SHA256));
		PyObject *o = SWIG_NewPointerObj(SWIG_as_voidptr(newKey), SWIGTYPE_p_cipher_SHA256, SWIG_POINTER_OWN );
		PyList_Append(list, o);
		Py_DECREF(o);
	}
	if( $1->data != NULL)
		free( (void*)$1->data );
	%append_output( list );
}



