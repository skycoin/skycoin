%{

GoUint32_ _WrapperFeeCalculator(Transaction__Handle handle, GoUint64_* pFee, void* context){
	PyObject* feeCalc = (PyObject*)context;
	PyObject *result = PyObject_CallFunctionObjArgs(feeCalc, PyLong_FromLong(handle), NULL);
	GoUint32_ error = 0;
	if(PyTuple_Check(result)){
		PyObject* objerror = PyTuple_GetItem(result, 0);
		error = PyLong_AsLong(objerror);
		result = PyTuple_GetItem(result, 1);
	}
	if(error != 0)
		return error;
	GoUint64_ ret = PyLong_AsLong(result);
  	Py_DECREF(result);
	if(pFee){
		*pFee = ret;
		return 0;
	}
	else
		return 1;
}
%}

%typemap(in) FeeCalculator* (FeeCalculator temp) {
  if (!PyCallable_Check($input)) SWIG_fail;
  temp.callback = _WrapperFeeCalculator;
  temp.context = $input;
  $1 = &temp;
}
