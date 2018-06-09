%module skycoin
%include "typemaps.i"
%{
	#define SWIG_FILE_WITH_INIT	
	#include "libskycoin.h"
	#include "include/extras.h"
%}

%include "simpletypes.i"
%include "handletypemaps.i"
%include "gostrings.i"
%include "goslices.i"
%include "goslices2.i"
%include "structs.i"
%include "extend.i"
%include "structs_typemaps.i"
/* Find the modified copy of libskycoin */
%include "../../../swig/include/libskycoin.h"




