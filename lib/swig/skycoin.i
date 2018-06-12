%module skycoin
%include "typemaps.i"
%{
	#define SWIG_FILE_WITH_INIT
	#include "libskycoin.h"
	#include "include/extras.h"
%}

%include "golang.cgo.i"
%include "skycoin.mem.i"
%include "skycoin.cipher.crypto.i"
%include "skycoin.cipher.encrypt.i"
%include "skycoin.cipher.secp256k1go.i"
%include "skycoin.coin.transaction.i"
%include "skycoin.wallet.i"
%include "structs_typemaps.i"
/* Find the modified copy of libskycoin */
%include "../../../swig/include/libskycoin.h"
