/* Handle out typemap. */
%typemap(argout) Handle* {
	%append_output( SWIG_From_long(*$1) );
}

/* Handle not as pointer is input. */
%typemap(in) Handle {
	SWIG_AsVal_long($input, (long*)&$1);
}

%apply Handle { Wallet__Handle, Options__Handle, ReadableEntry__Handle, ReadableWallet__Handle, WebRpcClient__Handle,
	WalletResponse__Handle, Client__Handle, Strings__Handle, Wallets__Handle, Config__Handle, App__Handle, Context__Handle,
	GoStringMap, PasswordReader__Handle_,
	Transaction__Handle, Transactions__Handle, CreatedTransaction__Handle,
	CreatedTransactionOutput__Handle, CreatedTransactionInput__Handle, CreateTransactionResponse__Handle,
	Block__Handle, SignedBlock__Handle, BlockBody__Handle, BuildInfo_Handle, Number_Handle, Signature_Handle,AddressUxOuts_Handle
	}

%apply Handle* { Wallet__Handle*, Options__Handle*, ReadableEntry__Handle*, ReadableWallet__Handle*, WebRpcClient__Handle*,
	WalletResponse__Handle*, Client__Handle*, Strings__Handle*, Wallets__Handle*, Config__Handle*,
	App__Handle*, Context__Handle*, GoStringMap_*, PasswordReader__Handle*,
	Transaction__Handle*, Transactions__Handle*, CreatedTransaction__Handle*,
	CreatedTransactionOutput__Handle*, CreatedTransactionInput__Handle*, CreateTransactionResponse__Handle*,
	Block__Handle*, SignedBlock__Handle*, BlockBody__Handle*, BuildInfo_Handle*, Number_Handle*, Signature_Handle*,AddressUxOuts_Handle*
	}

	/*
cipher__PubKey* input typemap
*/
%typemap(in) cipher__PubKey* {
	void *argp = 0;
	int res = SWIG_ConvertPtr($input, &argp, cipher__PubKey, 0 | 0);
	// if (!SWIG_IsOK(res))
		// SWIG_exception_fail(SWIG_TypeError, "expecting type PubKey");
	cipher_PubKey* p = (cipher_PubKey*)argp;
	$1 = &p->data;
}


/*
cipher__SecKey* input typemap
*/
%typemap(in) cipher__SecKey*{
	void *argp = 0;
	// int res = SWIG_ConvertPtr($input, &argp, SWIGTYPE_p_cipher_SecKey, 0 | 0);
	// if (!SWIG_IsOK(res))
		// SWIG_exception_fail(SWIG_TypeError, "expecting type SecKey");
	cipher_SecKey* p = (cipher_SecKey*)argp;
	$1 = &p->data;
}

%typemap(in) cipher__Ripemd160* {
	void *argp = 0;
	// int res = SWIG_ConvertPtr($input, &argp, SWIGTYPE_p_cipher_Ripemd160, 0 | 0);
	// if (!SWIG_IsOK(res))
		// SWIG_exception_fail(SWIG_TypeError, "expecting type Ripemd160");
	cipher_Ripemd160* p = (cipher_Ripemd160*)argp;
	$1 = &p->data;
}

%typemap(in) cipher__Sig* {
	void *argp = 0;
	// int res = SWIG_ConvertPtr($input, &argp, SWIGTYPE_p_cipher_Sig, 0 | 0);
	// if (!SWIG_IsOK(res))
		// SWIG_exception_fail(SWIG_TypeError, "expecting type Sig");
	cipher_Sig* p = (cipher_Sig*)argp;
	$1 = &p->data;
}



%typemap(in) cipher__SHA256* {
	void *argp = 0;
	// int res = SWIG_ConvertPtr($input, &argp, cipher__SHA256, 0 | 0);
	// if (!SWIG_IsOK(res))
		// SWIG_exception_fail(SWIG_TypeError, "expecting type SHA256");
	cipher_SHA256* p = (cipher_SHA256*)argp;
	$1 = &p->data;
}

%typemap(in) cipher__Checksum* {
	void *argp = 0;
	// int res = SWIG_ConvertPtr($input, &argp, SWIGTYPE_p_cipher_Checksum, 0 | 0);
	// if (!SWIG_IsOK(res))
		// SWIG_exception_fail(SWIG_TypeError, "expecting type Checksum");
	cipher_Checksum* p = (cipher_Checksum*)argp;
	$1 = &p->data;
}