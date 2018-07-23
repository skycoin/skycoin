%typemap(in) FeeCalc {
  PyObject *callback;
  GoUint32_ _WrapperFeeCalculator(Transaction__Handle handle, GoUint64_* pFee){
	PyObject *result = PyObject_CallFunctionObjArgs(callback, PyLong_FromLong(handle), NULL);
	const GoUint64_ ret = PyLong_AsLong(result);
  	Py_DECREF(result);
	if(pFee){
		*pFee = ret;
		return 0;
	}
	else{
		return -1;
	}	
  }	
  if (!PyCallable_Check($input)) SWIG_fail;
  $1 = _WrapperFeeCalculator;
  callback = $input;
}
