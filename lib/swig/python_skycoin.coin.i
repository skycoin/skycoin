%extend coin__BlockHeader {
	int __eq__(coin__BlockHeader* a){
		if( $self->Version != a->Version || $self->Time != a->Time || 
			$self->BkSeq != a->BkSeq || $self->Fee != a->Fee)
			return 0;
		if( memcmp( &$self->PrevHash, a->PrevHash, sizeof(a->PrevHash) ) != 0 )
			return 0;
		if( memcmp( &$self->BodyHash, a->PrevHash, sizeof(a->BodyHash) ) != 0 )
			return 0;
		if( memcmp( &$self->UxHash, a->PrevHash, sizeof(a->UxHash) ) != 0 )
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




