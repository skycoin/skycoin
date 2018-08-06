%typecheck(SWIG_TYPECHECK_STRING_ARRAY) cipher_Addresses* {
  $1 = PyList_Check($input) ? 1 : 0;
}

/*cipher_Addresses* input parameter */
%typemap(in) (cipher_Addresses* __in_addresses) (cipher_Addresses temp) {
	int i;
	$1 = &temp;
	$1->count = PyList_Size($input);
	$1->data = malloc(sizeof(cipher__Address) * $1->count);
	cipher__Address* pdata = $1->data;
	for(i = 0; i < $1->count; i++){
		PyObject *o = PyList_GetItem($input, i);
		void *argp = 0;
		int res = SWIG_ConvertPtr(o, &argp, SWIGTYPE_p_cipher__Address, 0 | 0);
		if (!SWIG_IsOK(res))
			SWIG_exception_fail(SWIG_TypeError, "expecting type cipher__Address");
		cipher__Address* p = (cipher__Address*)argp;
		memcpy(pdata, p, sizeof(cipher__Address));
		pdata++;
	}
}

%typemap(freearg) (cipher_Addresses* __in_addresses) {
  if ($1->data) free($1->data);
}

/*cipher_Addresses* parameter to return as a list */
%typemap(in, numinputs=0) (cipher_Addresses*  __out_addresses) (cipher_Addresses temp) {
	temp.data = NULL;
	temp.count = 0;
	$1 = &temp;
}

/*cipher_Addresses* as function return typemap*/
%typemap(argout) (cipher_Addresses* __out_addresses) {
	int i;
	PyObject *list = PyList_New(0);
	for (i = 0; i < $1->count; i++) {
		cipher__Address* key = &($1->data[i]);
		cipher__Address* newKey = malloc(sizeof(cipher__Address));
		memcpy(newKey, key, sizeof(cipher__Address));
		PyObject *o = SWIG_NewPointerObj(SWIG_as_voidptr(newKey), SWIGTYPE_p_cipher__Address, SWIG_POINTER_OWN );
		PyList_Append(list, o);
		Py_DECREF(o);
	}
	if( $1->data != NULL)
		free( (void*)$1->data );
	%append_output( list );
}



