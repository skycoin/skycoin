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

%typemap(in) FeeCalculator* (FeeCalculator temp) {
  if (!PyCallable_Check($input)) return ;
  temp.callback = _WrapperFeeCalculator;
  temp.context = $input;
  $1 = &temp;
}