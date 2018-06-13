%module skycoin
%include "typemaps.i"
%{
	#define SWIG_FILE_WITH_INIT
	#include "libskycoin.h"
	#include "swig.h"
%}

%include "golang.cgo.i"
%include "skycoin.mem.i"
%include "skycoin.cipher.crypto.i"
%include "structs_typemaps.i"

%include "swig.h"
/* Find the modified copy of libskycoin */
%include "libskycoin.h"
%include "structs.i"
