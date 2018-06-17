#if defined(SWIGPYTHON)
	%include "python_seckeys.i"
	%include "python_pubkeys.i"
#else
%rename(SKY_cipher_GenerateDeterministicKeyPairs) wrap_SKY_cipher_GenerateDeterministicKeyPairs;
%inline {
	GoUint32 wrap_SKY_cipher_GenerateDeterministicKeyPairs(GoSlice seed, GoInt n, cipher_SecKeys* secKeys){
		GoSlice_ data;
		data.data = secKeys->data;
		data.len = secKeys->count;
		data.cap = secKeys->count;
		GoUint32 result = SKY_cipher_GenerateDeterministicKeyPairs(seed, n, &data);
		return result;
	}
}

%rename(SKY_cipher_GenerateDeterministicKeyPairsSeed) wrap_SKY_cipher_GenerateDeterministicKeyPairsSeed;
%inline {
	GoUint32 wrap_SKY_cipher_GenerateDeterministicKeyPairsSeed(GoSlice seed, GoInt n, coin__UxArray* newSeed, cipher_SecKeys* secKeys){
		GoSlice_ data;
		data.data = secKeys->data;
		data.len = secKeys->count;
		data.cap = secKeys->count;
		GoUint32 result = SKY_cipher_GenerateDeterministicKeyPairsSeed(seed, n, newSeed, &data);
		return result;
	}
}

%rename(SKY_cipher_PubKeySlice_Len) wrap_SKY_cipher_PubKeySlice_Len;
%inline {
	GoUint32 wrap_SKY_cipher_PubKeySlice_Len(cipher_PubKeys* pubKeys){
		GoSlice_ data;
		data.data = pubKeys->data;
		data.len = pubKeys->count;
		data.cap = pubKeys->count;
		GoUint32 result = SKY_cipher_PubKeySlice_Len(&data);
		return result;
	}
}

%rename(SKY_cipher_PubKeySlice_Less) wrap_SKY_cipher_PubKeySlice_Less;
%inline {
	GoUint32 wrap_SKY_cipher_PubKeySlice_Less(cipher_PubKeys* pubKeys, GoInt p1, GoInt p2){
		GoSlice_ data;
		data.data = pubKeys->data;
		data.len = pubKeys->count;
		data.cap = pubKeys->count;
		GoUint32 result = SKY_cipher_PubKeySlice_Less(&data, p1, p2);
		return result;
	}
}

%rename(SKY_cipher_PubKeySlice_Swap) wrap_SKY_cipher_PubKeySlice_Swap;
%inline {
	GoUint32 wrap_SKY_cipher_PubKeySlice_Swap(cipher_PubKeys* pubKeys, GoInt p1, GoInt p2){
		GoSlice_ data;
		data.data = pubKeys->data;
		data.len = pubKeys->count;
		data.cap = pubKeys->count;
		GoUint32 result = SKY_cipher_PubKeySlice_Swap(&data, p1, p2);
		return result;
	}
}

#endif




/**
*
* typemaps for Handles
*
**/


/* Handle reference typemap. */
%typemap(in, numinputs=0) Handle* (Handle temp) {
	$1 = &temp;
}

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
	Block__Handle, SignedBlock__Handle, BlockBody__Handle
	}

%apply Handle* { Wallet__Handle*, Options__Handle*, ReadableEntry__Handle*, ReadableWallet__Handle*, WebRpcClient__Handle*,
	WalletResponse__Handle*, Client__Handle*, Strings__Handle*, Wallets__Handle*, Config__Handle*,
	App__Handle*, Context__Handle*, GoStringMap_*, PasswordReader__Handle*,
	Transaction__Handle*, Transactions__Handle*, CreatedTransaction__Handle*,
	CreatedTransactionOutput__Handle*, CreatedTransactionInput__Handle*, CreateTransactionResponse__Handle*,
	Block__Handle*, SignedBlock__Handle*, BlockBody__Handle*
	}
