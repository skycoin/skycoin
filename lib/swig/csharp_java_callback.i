// %{
// GoUint32_ CSharp_skycoin_FeeCalculator(Transaction__Handle handle, GoUint64_* pFee, void* context){
// 	void* feeCalc = (void*)context;
// 	void *result = PyObject_CallFunctionObjArgs(feeCalc, PyLong_FromLong(handle), NULL);
// 	GoUint32_ error = 0;
// 	if(PyTuple_Check(result)){
// 		void* objerror = PyTuple_GetItem(result, 0);
// 		error = PyLong_AsLong(objerror);
// 		result = PyTuple_GetItem(result, 1);
// 	}
// 	if(error != 0)
// 		return error;
// 	GoUint64_ ret = PyLong_AsLong(result);
//   	Py_DECREF(result);
// 	if(pFee){
// 		*pFee = ret;
// 		return 0;
// 	}
// 	else
// 		return 1;
// }
// %}

// %typemap(in) FeeCalculator* (FeeCalculator temp) {
//   if (!PyCallable_Check($input)) ;
//   temp.callback = CSharp_skycoin_FeeCalculator;
//   temp.context = $input;
//   $1 = &temp;
// }
