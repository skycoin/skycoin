/*GoSlice in typemap*/
%typemap(in) GoSlice {
	char* buffer = 0;
	size_t size = 0;
	int res = SWIG_AsCharPtrAndSize( $input, &buffer, &size, 0 );
	if (!SWIG_IsOK(res)) {
		SWIG_exception_fail(SWIG_TypeError, "in method '$symname', expecting string");
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
	int res = SWIG_AsCharPtrAndSize( $input, &buffer, &size, 0 );
	if (!SWIG_IsOK(res)) {
		SWIG_exception_fail(SWIG_TypeError, "in method '$symname', expecting string");
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

/*GoUint64* parameter as reference */
%typemap(in, numinputs=0) GoUint64* (GoUint64 temp) {
	temp = 0;
	$1 = &temp;
}

/*GoUint64* as function return typemap*/
%typemap(argout) GoUint64* {
	%append_output( SWIG_From_long( *$1 ) );
}

/*GoInt64* parameter as reference */
%typemap(in, numinputs=0) GoInt64* (GoInt64 temp) {
	temp = 0;
	$1 = &temp;
}



%apply GoString {GoString_}
%apply GoString* {GoString_*}


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
