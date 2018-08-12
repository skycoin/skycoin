%typecheck(SWIG_TYPECHECK_STRING_ARRAY) coin_UxOutArray* {
  $1 = PyList_Check($input) ? 1 : 0;
}

/*coin_UxOutArray* input parameter */
%typemap(in) (coin_UxOutArray* __uxIn) (coin_UxOutArray temp) {
	int i;
	$1 = &temp;
	$1->count = PyList_Size($input);
	$1->data = malloc(sizeof(coin__UxOut) * $1->count);
	coin__UxOut* pdata = $1->data;
	for(i = 0; i < $1->count; i++){
		PyObject *o = PyList_GetItem($input, i);
		void *argp = 0;
		int res = SWIG_ConvertPtr(o, &argp, SWIGTYPE_p_coin__UxOut, 0 | 0);
		if (!SWIG_IsOK(res))
			SWIG_exception_fail(SWIG_TypeError, "expecting type UxOut");
		coin__UxOut* p = (coin__UxOut*)argp;
		memcpy(pdata, p, sizeof(coin__UxOut));
		pdata++;
	}
}

%typemap(freearg) (coin_UxOutArray* __uxIn) {
  if ($1->data) free($1->data);
}

%apply (coin_UxOutArray* __uxIn) {(coin_UxOutArray* __uxOut), (coin_UxOutArray* __uxIn2)}

/*coin_UxOutArray* parameter to return as a list */
%typemap(in, numinputs=0) (coin_UxOutArray*  __return_Ux) (coin_UxOutArray temp) {
	temp.data = NULL;
	temp.count = 0;
	$1 = &temp;
}

/*coin_UxOutArray* as function return typemap*/
%typemap(argout) (coin_UxOutArray* __return_Ux) {
	int i;
	PyObject *list = PyList_New(0);
	for (i = 0; i < $1->count; i++) {
		coin__UxOut* key = &($1->data[i]);
		coin__UxOut* newKey = malloc(sizeof(coin__UxOut));
		memcpy(newKey, key, sizeof(coin__UxOut));
		PyObject *o = SWIG_NewPointerObj(SWIG_as_voidptr(newKey), SWIGTYPE_p_coin__UxOut, SWIG_POINTER_OWN );
		PyList_Append(list, o);
		Py_DECREF(o);
	}
	if( $1->data != NULL)
		free( (void*)$1->data );
	%append_output( list );
}
