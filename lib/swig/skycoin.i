%module skycoin
%include "typemaps.i"
%{
	#define SWIG_FILE_WITH_INIT
	#include "libskycoin.h"
	#include "include/extras.h"
%}

%include "golang.cgo.i"
%include "skycoin.mem.i"
%include "custom.i"
%include "structs_typemaps.i"

/* Find the modified copy of libskycoin */
%include "../../../swig/include/libskycoin.h"
%include "structs.i"
