typedef struct{
	char* p;
	int   n;
} GoString;

typedef struct{
	char* p;
	int   n;
} GoString_;


/*GoString in typemap*/
%typemap(in) GoString {
	char* buffer = 0;
	size_t size = 0;
	int res = SWIG_AsCharPtrAndSize( $input, &buffer, &size, 0 );
	if (!SWIG_IsOK(res)) {
		%argument_fail(res, "(TYPEMAP, SIZE)", $symname, $argnum);
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
