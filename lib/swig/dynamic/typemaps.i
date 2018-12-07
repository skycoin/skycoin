/*GoInt64* as function return typemap*/
%typemap(argout) GoInt64* {
	%append_output( SWIG_From_long( *$1 ) );
}

%apply int* OUTPUT {GoInt*}
%apply int* OUTPUT {GoUint*}
%apply int* OUTPUT {GoUint8*}
%apply int* OUTPUT {GoInt8*}
%apply int* OUTPUT {GoUint16*}
%apply int* OUTPUT {GoInt16*}
%apply int* OUTPUT {GoUint32*}
%apply int* OUTPUT {GoInt32*}
