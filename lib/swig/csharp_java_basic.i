%include "arrays_csharp.i"

%apply int INOUT[] {int *array1}
%apply int INOUT[] {int *array2}
%apply int FIXED[] {int *sourceArray}
%apply int FIXED[] {int *targetArray}

%include cpointer.i
%pointer_functions(cipher_PubKey, cipher_PubKeyp);
%pointer_functions(cipher_SecKey, cipher_SecKeyp);
// %pointer_functions(GoSlice, GoSlicep);
// %pointer_functions(GoString_, GoStringp_);
