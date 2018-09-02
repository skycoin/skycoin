%module skycoin
%include "typemaps.i"
%{
	#define SWIG_FILE_WITH_INIT
	#include "libskycoin.h"
	#include "swig.h"
%}

//Apply strictly to python
//Not for other languages
#if defined(SWIGPYTHON)
%include "python_skycoin.cipher.crypto.i"
%include "python_uxarray.i"
%include "python_sha256s.i"
%include "python_skycoin.coin.i"
%include "python_skycoin.callback.i"
%include "golang.cgo.i"
%include "structs_typemaps.i"
%include "skycoin.mem.i"
#endif

%include "skycoin.cipher.crypto.i"
%include "skycoin.coin.i"
#if defined(SWIGCSHARP)
%include "csharp_java_basic.i"
%include "csharp_java_typemap.i"
%include "csharp_java_skyerror.i"
#endif

%include "swig.h"
/* Find the modified copy of libskycoin */
%include "libskycoin.h"
%include "structs.i"
