%{

GoUint32_ _WrapperFeeCalculator(Transaction__Handle handle, GoUint64_* pFee, void* context){
	FeeCalcFunc* feeCalc = (FeeCalcFunc*)context;
	int *result = callFeeCalculator(feeCalc, handle, pFee);
	GoUint32_ error = 0;
	if(result != 0)
		return error;
  return 0;
}
%}

%typemap(in) FeeCalculator_* (FeeCalculator temp) {
  if (!PyCallable_Check($input)) return ;
  temp.callback = _WrapperFeeCalculator;
  temp.context = $input;
  $1 = &temp;
}


%define %cs_callback(TYPE, CSTYPE) 
        %typemap(ctype) TYPE, TYPE& "void*" 
        %typemap(in) TYPE  %{ $1 = (TYPE)$input; %} 
        %typemap(in) TYPE& %{ $1 = (TYPE*)&$input; %} 
        %typemap(imtype, out="IntPtr") TYPE, TYPE& "CSTYPE" 
        %typemap(cstype, out="IntPtr") TYPE, TYPE& "CSTYPE" 
        %typemap(csin) TYPE, TYPE& "$csinput" 
%enddef 
%define %cs_callback2(TYPE, CTYPE, CSTYPE) 
        %typemap(ctype) TYPE "CTYPE" 
        %typemap(in) TYPE %{ $1 = (TYPE)$input; %} 
        %typemap(imtype, out="IntPtr") TYPE "CSTYPE" 
        %typemap(cstype, out="IntPtr") TYPE "CSTYPE" 
        %typemap(csin) TYPE "$csinput" 
%enddef 

%cs_callback2(FeeCalcFunc*,FeeCalcFunc*,FeeCalcFunc*)