%include "cmp.i"

%extend coin__BlockHeader {
	int __eq__(coin__BlockHeader* bh){
		return equalBlockHeaders($self, bh);
	}
}

%extend coin__Transaction {
	int __eq__(coin__Transaction* t){
		return equalTransactions($self, t);
	}
}

%extend coin__BlockBody {
	int __eq__(coin__BlockBody* b){
		return equalTransactionsArrays(&$self->Transactions, &b->Transactions);
	}
}

%extend coin__UxOut {
	int __eq__(coin__UxOut* u){
		return memcmp(&$self, u, sizeof(coin__UxOut)) == 0;
	}
}

%extend coin__TransactionOutput {
	int __eq__(coin__TransactionOutput* t){
		if( $self->Coins != t->Coins ||
			$self->Hours != t->Hours ){
			return 0;
	  	}

	  	if(memcmp(&$self->Address, &t->Address, sizeof(cipher__Address)) != 0)
			return 0;
	 	return 1;
	}
}

%typemap(in, numinputs=0) (coin__Transaction**) (coin__Transaction* temp) {
	temp = NULL;
	$1 = &temp;
}

/*Return a pointer created with own = 0 because
Python will not own the object
 */
%typemap(argout) (coin__Transaction**) {
	%append_output( SWIG_NewPointerObj(SWIG_as_voidptr(*$1), SWIGTYPE_p_coin__Transaction, 0 ) );
}

%typemap(in, numinputs=0) (coin__Block**) (coin__Block* temp) {
	temp = NULL;
	$1 = &temp;
}

/*Return a pointer created with own = 0 because
Python will not own the object
 */
%typemap(argout) (coin__Block**) {
	%append_output( SWIG_NewPointerObj(SWIG_as_voidptr(*$1), SWIGTYPE_p_coin__Block, 0 ) );
}


%typemap(in, numinputs=0) (coin__UxArray**) (coin__UxArray* temp){
	temp = NULL;
	$1 = &temp;
}

%typemap(argout) (coin__UxArray**) {
	int i;
	PyObject *list = PyList_New(0);
	coin__Transaction* pTrans = (*$1)->data;
	i = (*$1)->len;
	for (; i > 0; i--) {
		PyObject *o = SWIG_NewPointerObj(SWIG_as_voidptr(pTrans), SWIGTYPE_p_coin__Transaction, 0 );
		PyList_Append(list, o);
		pTrans++;
	}
	%append_output( list );
}




