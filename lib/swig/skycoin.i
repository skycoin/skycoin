%module skycoin
%include "typemaps.i"
%{
	#define SWIG_FILE_WITH_INIT
	#include "libskycoin.h"
	#include "swig.h"
%}


//Apply typemaps only for Python
//It can be applied to other languages that fit in
#if defined(SWIGPYTHON)
%include "golang.cgo.i"
%include "skycoin.mem.i"
%include "python_skycoin.cipher.crypto.i"
%include "python_skycoin.coin.i"
%include "structs_typemaps.i"
#else
%include "skycoin.cipher.crypto.i"
%include "skycoin.coin.i"
#endif

%include "swig.h"
/* Find the modified copy of libskycoin */
%include "libskycoin.h"
%include "structs.i"
