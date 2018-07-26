/*cipher_PubKeys* input parameter */
%typemap(in) (cipher_PubKeys* __in_pubKeys) (cipher_PubKeys temp) {
	int i;
	$1 = &temp;
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
		memcpy(pdata, p, sizeof(cipher_PubKey));
		pdata++;
	}
}

%typemap(freearg) (cipher_PubKeys* __in_pubKeys) {
  if ($1->data) free($1->data);
}

