

%apply long long  {GoInt_, GoInt};
%apply long int  {ptrdiff_t};
%apply unsigned short  {GoUint16, GoUint16_};
%apply unsigned char  {GoUint8_, GoUint8};
%apply signed char  {GoInt8_, GoInt8};
%apply unsigned long long  {GoUint64, GoUint64_,GoUint_,GoUint};
%apply short {GoInt16_, GoInt16};
%apply int {GoInt32_, Go, GoInt32};
%apply unsigned int {GoUint32_, GoUint32, BOOL, error};
%apply long long  {GoInt64_, GoInt64, GoInt_, GoInt};
%apply float {GoFloat32_, GoFloat32};
%apply double {GoFloat64_, GoFloat64};
%apply long unsigned int {GoUintptr_, GoUintptr}

