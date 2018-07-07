%begin %{
#define SWIG_PYTHON_STRICT_BYTE_CHAR
#define SWIG_PYTHON_STRICT_UNICODE_WCHAR
%}


/*GoSlice in typemap*/
%typemap(in) GoSlice {
	char* buffer = 0;
	size_t size = 0;
	int res = SWIG_AsCharPtrAndSize( $input, &buffer, &size, 0 );
	if (!SWIG_IsOK(res)) {
		SWIG_exception_fail(SWIG_TypeError, "in method '$symname', expecting byte string");
	}
	$1.data = buffer;
	$1.len = size - 1;
	$1.cap = size;
}


%typecheck(SWIG_TYPECHECK_STRING) GoSlice {
  	char* buffer = 0;
	size_t size = 0;
	int res = SWIG_AsCharPtrAndSize( $input, &buffer, &size, 0 );
	$1 = SWIG_IsOK(res) ? 1 : 0;
}

/*GoSlice_* parameter as reference */
%typemap(in, numinputs=0) GoSlice_* (GoSlice_ temp) {
	temp.data = NULL;
	temp.len = 0;
	temp.cap = 0;
	$1 = ($1_type)&temp;
}

/*GoSlice_* as function return typemap*/
%typemap(argout) GoSlice_* {
	%append_output( SWIG_FromCharPtrAndSize( $1->data, $1->len  ) );
	free( (void*)$1->data );
}

%apply GoSlice_* {coin__UxArray*}


/*GoString in typemap*/
%typemap(in) GoString {
	char* buffer = 0;
	size_t size = 0;
	/*if (PyString_Check($input)){
		SWIG_exception_fail(SWIG_TypeError, "in method '$symname', type is not unicode string");
	}
	else if(PyUnicode_Check($input)){
		SWIG_exception_fail(SWIG_TypeError, "in method '$symname', type is unicode unicode");
	}
	else {
		SWIG_exception_fail(SWIG_TypeError, "in method '$symname', what the hell is this");
	}*/
	/*
	%#if PY_VERSION_HEX>=0x03000000
	%#if defined(SWIG_PYTHON_STRICT_BYTE_CHAR)
	  if (PyBytes_Check($input)){
	  	char *cstr; Py_ssize_t len;
		int ret = SWIG_OK;
		PyBytes_AsStringAndSize($input, &cstr, &len);
		buffer = PyBytes_AsString($input);
		if( !buffer )
			SWIG_exception_fail(SWIG_TypeError, "in method '$symname', PyBytes_AsString failed");
		else
			SWIG_exception_fail(SWIG_TypeError, "in method '$symname', PyBytes_AsString ok");
	  }
	%#else
	  if (PyUnicode_Check($input)){
	  	SWIG_exception_fail(SWIG_TypeError, "in method '$symname', we shouldn't be here");
	  }
	%#endif
	%#else  
	  if (PyString_Check($input)){
	  	char *cstr; Py_ssize_t len;
		int ret = SWIG_OK;
	  	PyString_AsStringAndSize($input, &cstr, &len);
	  	buffer = SWIG_Python_str_AsChar($input);
	  	if( !buffer )
			SWIG_exception_fail(SWIG_TypeError, "in method '$symname', SWIG_Python_str_AsChar failed");
		else
			SWIG_exception_fail(SWIG_TypeError, "in method '$symname', SWIG_Python_str_AsChar ok");
	  }
	%#endif
	*/
	int calloc = SWIG_OLDOBJ;
	int res = SWIG_AsCharPtrAndSize( $input, &buffer, &size, &calloc );
	if (!SWIG_IsOK(res)) {
		if( res == SWIG_TypeError)
			SWIG_exception_fail(SWIG_TypeError, "in method '$symname', expecting byte string");
		else
			SWIG_exception_fail(SWIG_TypeError, "in method '$symname', some unknown error");
	}
	$1.p = buffer;
	$1.n = size - 1;
}

/*GoString* parameter as reference */
%typemap(in, numinputs=0) GoString* (GoString temp) {
	temp.p = NULL;
	temp.n = 0;
	$1 = ($1_type)&temp;
}

/*GoString* as function return typemap*/
%typemap(argout) GoString* {
	%append_output( SWIG_FromCharPtrAndSize( $1->p, $1->n  ) );
	free( (void*)$1->p );
}

%apply GoString {GoString_}
%apply GoString* {GoString_*}

%apply int* OUTPUT {GoInt*}
%apply int* OUTPUT {GoUint*}
%apply int* OUTPUT {GoUint8*}
%apply int* OUTPUT {GoInt8*}
%apply int* OUTPUT {GoUint16*}
%apply int* OUTPUT {GoInt16*}
%apply int* OUTPUT {GoUint32*}
%apply int* OUTPUT {GoInt32*}
%apply int* OUTPUT {GoUint64*}
%apply int* OUTPUT {GoInt64*}

typedef GoInt GoInt_;
typedef GoUint GoUint_;
typedef GoInt8 GoInt8_;
typedef GoUint8 GoUint8_;
typedef GoInt16 GoInt16_;
typedef GoUint16 GoUint16_;
typedef GoInt32 GoInt32_;
typedef GoUint32 GoUint32_;
typedef GoInt64 GoInt64_;
typedef GoUint64 GoUint64_;
