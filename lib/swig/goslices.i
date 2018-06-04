/*GoSlice in typemap*/
%typemap(in) GoSlice { 
	char* buffer = 0;
	size_t size = 0;
	int res = SWIG_AsCharPtrAndSize( $input, &buffer, &size, 0 );
	if (!SWIG_IsOK(res)) {
		%argument_fail(res, "(TYPEMAP, SIZE)", $symname, $argnum);
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

