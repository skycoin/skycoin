%apply unsigned short  {GoUint16, GoUint16_};
%apply unsigned long  {GoUintptr, __SIZE_TYPE__};
%apply short  {GoInt16, GoInt16_};
%apply unsigned char  {GoUint8_, GoUint8};
%apply unsigned int  {GoUint32_, GoUint32};
%apply signed char  {GoInt8_, GoInt8};
%apply unsigned long long  {GoUint64, GoUint64_,GoUint,GoUint_};
%apply long long  {GoInt64, GoInt64_,GoInt_, GoInt };
%apply GoSlice_* {coin__UxArray*,GoSlice_**};
%apply int {GoInt32,GoInt32_,ptrdiff_t};
%apply int* {GoInt32*,GoInt32_*,ptrdiff_t*};
%apply float {GoFloat32};
%apply double {GoFloat64};

%typemap(freearg) (cipher_PubKeys* __in_pubKeys) {
  if ($1->data) free($1->data);
}