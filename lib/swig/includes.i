//Apply strictly to python
//Not for other languages
%include "cmp.i"
#if !defined(SWIGPYTHON)
%include "skycoin.cipher.crypto.i"
%include "skycoin.coin.i"
#endif
