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
	Block__Handle, SignedBlock__Handle, BlockBody__Handle, BuildInfo_Handle, Number_Handle, Signature_Handle,AddressUxOuts_Handle
	}

%apply Handle* { Wallet__Handle*, Options__Handle*, ReadableEntry__Handle*, ReadableWallet__Handle*, WebRpcClient__Handle*,
	WalletResponse__Handle*, Client__Handle*, Strings__Handle*, Wallets__Handle*, Config__Handle*,
	App__Handle*, Context__Handle*, GoStringMap_*, PasswordReader__Handle*,
	Transaction__Handle*, Transactions__Handle*, CreatedTransaction__Handle*,
	CreatedTransactionOutput__Handle*, CreatedTransactionInput__Handle*, CreateTransactionResponse__Handle*,
	Block__Handle*, SignedBlock__Handle*, BlockBody__Handle*, BuildInfo_Handle*, Number_Handle*, Signature_Handle*,AddressUxOuts_Handle*
	}

%typecheck(SWIG_TYPECHECK_INTEGER) Transaction__Handle {
  $1 = PyInt_Check($input) ? 1 : 0;
}

%typecheck(SWIG_TYPECHECK_INTEGER) Transactions__Handle {
  $1 = PyInt_Check($input) ? 1 : 0;
}

%typecheck(SWIG_TYPECHECK_INTEGER) AddressUxOuts_Handle {
  $1 = PyInt_Check($input) ? 1 : 0;
}

#if defined(SWIGPYTHON)
	%include "python_seckeys.i"
	%include "python_pubkeys.i"
	%include "python_uxarray.i"
	%include "python_addresses.i"
#endif

%rename(SKY_coin_Transaction_SignInputs) wrap_SKY_coin_Transaction_SignInputs;
%inline{
	GoUint32 wrap_SKY_coin_Transaction_SignInputs(Transaction__Handle handle, cipher_SecKeys* __in_secKeys){
		GoSlice data;
		data.data = __in_secKeys->data;
		data.len = __in_secKeys->count;
		data.cap = __in_secKeys->count;
		return SKY_coin_Transaction_SignInputs(handle, data);
	}
}


%rename(SKY_cipher_GenerateDeterministicKeyPairs) wrap_SKY_cipher_GenerateDeterministicKeyPairs;
%inline {
	GoUint32 wrap_SKY_cipher_GenerateDeterministicKeyPairs(GoSlice seed, GoInt n, cipher_SecKeys* __out_secKeys){
		__out_secKeys->data = NULL;
		__out_secKeys->count = 0;
		GoSlice_ data;
		data.data = malloc(sizeof(cipher_SecKey) * n);
		data.len = n;
		data.cap = n;
		GoUint32 result = SKY_cipher_GenerateDeterministicKeyPairs(seed, n, &data);
		if( result == 0){
			__out_secKeys->data = data.data;
			__out_secKeys->count = data.len;
		}
		return result;
	}
}

%inline {
	GoUint32 wrap_SKY_cipher_GenerateDeterministicKeyPairsSeed(GoSlice seed, GoInt n, coin__UxArray* newSeed, cipher_SecKeys* __out_secKeys){
		__out_secKeys->data = NULL;
		__out_secKeys->count = 0;
		GoSlice_ data;
		data.data = malloc(sizeof(cipher_SecKey) * n);
		data.len = n;
		data.cap = n;
		GoUint32 result = SKY_cipher_GenerateDeterministicKeyPairsSeed(seed, n, newSeed, &data);
		if( result == 0){
			__out_secKeys->data = data.data;
			__out_secKeys->count = data.len;
		}
		return result;
	}
}

// %rename(SKY_cipher_PubKeySlice_Len) wrap_SKY_cipher_PubKeySlice_Len;
// %inline {
// 	GoUint32 wrap_SKY_cipher_PubKeySlice_Len(cipher_PubKeys* __in_pubKeys){
// 		GoSlice_ data;
// 		data.data = __in_pubKeys->data;
// 		data.len = __in_pubKeys->count;
// 		data.cap = __in_pubKeys->count;
// 		GoUint32 result = SKY_cipher_PubKeySlice_Len(&data);
// 		return result;
// 	}
// }

// %rename(SKY_cipher_PubKeySlice_Less) wrap_SKY_cipher_PubKeySlice_Less;
// %inline {
// 	GoUint32 wrap_SKY_cipher_PubKeySlice_Less(cipher_PubKeys* __in_pubKeys, GoInt p1, GoInt p2){
// 		GoSlice_ data;
// 		data.data = __in_pubKeys->data;
// 		data.len = __in_pubKeys->count;
// 		data.cap = __in_pubKeys->count;
// 		GoUint32 result = SKY_cipher_PubKeySlice_Less(&data, p1, p2);
// 		return result;
// 	}
// }

// %rename(SKY_cipher_PubKeySlice_Swap) wrap_SKY_cipher_PubKeySlice_Swap;
// %inline {
// 	GoUint32 wrap_SKY_cipher_PubKeySlice_Swap(cipher_PubKeys* __in_pubKeys, GoInt p1, GoInt p2){
// 		GoSlice_ data;
// 		data.data = __in_pubKeys->data;
// 		data.len = __in_pubKeys->count;
// 		data.cap = __in_pubKeys->count;
// 		GoUint32 result = SKY_cipher_PubKeySlice_Swap(&data, p1, p2);
// 		return result;
// 	}
// }

%rename(SKY_coin_VerifyTransactionCoinsSpending) wrap_SKY_coin_VerifyTransactionCoinsSpending;
%inline {
	GoUint32 wrap_SKY_coin_VerifyTransactionCoinsSpending(coin_UxOutArray* __uxIn, coin_UxOutArray* __uxOut){
		GoSlice_ dataIn;
		dataIn.data = __uxIn->data;
		dataIn.len = __uxIn->count;
		dataIn.cap = __uxIn->count;
		GoSlice_ dataOut;
		dataOut.data = __uxOut->data;
		dataOut.len = __uxOut->count;
		dataOut.cap = __uxOut->count;
		GoUint32 result = SKY_coin_VerifyTransactionCoinsSpending(&dataIn, &dataOut);
		return result;
	};
}

%rename(SKY_coin_VerifyTransactionHoursSpending) wrap_SKY_coin_VerifyTransactionHoursSpending;
%inline {
	GoUint32 wrap_SKY_coin_VerifyTransactionHoursSpending(GoUint64 _headTime , coin_UxOutArray* __uxIn, coin_UxOutArray* __uxOut){
		GoSlice_ dataIn;
		dataIn.data = __uxIn->data;
		dataIn.len = __uxIn->count;
		dataIn.cap = __uxIn->count;
		GoSlice_ dataOut;
		dataOut.data = __uxOut->data;
		dataOut.len = __uxOut->count;
		dataOut.cap = __uxOut->count;
		GoUint32 result = SKY_coin_VerifyTransactionHoursSpending(_headTime, &dataIn, &dataOut);
		return result;
	};
}

%rename(SKY_coin_CreateUnspents) wrap_SKY_coin_CreateUnspents;
%inline {
	GoUint32 wrap_SKY_coin_CreateUnspents(coin__BlockHeader* bh, Transaction__Handle t, coin_UxOutArray* __return_Ux){
		__return_Ux->data = NULL;
		__return_Ux->count = 0;
		GoSlice_ data;
		data.data = NULL;
		data.len = 0;
		data.cap = 0;
		GoUint32 result = SKY_coin_CreateUnspents(bh, t, &data);
		if( result == 0){
			__return_Ux->data = data.data;
			__return_Ux->count = data.len;
		}
		return result;
	}
}

%rename(SKY_coin_Transaction_VerifyInput) wrap_SKY_coin_Transaction_VerifyInput;
%inline{
	GoUint32 wrap_SKY_coin_Transaction_VerifyInput(Transaction__Handle handle, coin_UxOutArray* __uxIn){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_Transaction_VerifyInput(handle, &data);
	}
}

%rename(SKY_coin_UxArray_HasDupes) wrap_SKY_coin_UxArray_HasDupes;
%inline{
	GoUint32 wrap_SKY_coin_UxArray_HasDupes(coin_UxOutArray* __uxIn, GoUint8* p1){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_UxArray_HasDupes(&data, p1);
	}
}

%rename(SKY_coin_UxArray_Coins) wrap_SKY_coin_UxArray_Coins;
%inline{
	GoUint32 wrap_SKY_coin_UxArray_Coins(coin_UxOutArray* __uxIn, GoUint64* p1){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_UxArray_Coins(&data, p1);
	}
}

%rename(SKY_coin_UxArray_CoinHours) wrap_SKY_coin_UxArray_CoinHours;
%inline{
	GoUint32 wrap_SKY_coin_UxArray_CoinHours(coin_UxOutArray* __uxIn, GoUint64 p1, GoUint64* p2){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_UxArray_CoinHours(&data, p1, p2);
	}
}

%rename(SKY_coin_UxArray_Less) wrap_SKY_coin_UxArray_Less;
%inline{
	GoUint32 wrap_SKY_coin_UxArray_Less(coin_UxOutArray* __uxIn, GoInt p1, GoInt p2, GoUint8* p3){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_UxArray_Less(&data, p1, p2, p3);
	}
}

%rename(SKY_coin_UxArray_Swap) wrap_SKY_coin_UxArray_Swap;
%inline{
	GoUint32 wrap_SKY_coin_UxArray_Swap(coin_UxOutArray* __uxIn, GoInt p1, GoInt p2){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_UxArray_Swap(&data, p1, p2);
	}
}

%rename(SKY_coin_UxArray_Sub) wrap_SKY_coin_UxArray_Sub;
%inline{
	GoUint32 wrap_SKY_coin_UxArray_Sub(coin_UxOutArray* __uxIn, coin_UxOutArray* __uxIn2, coin_UxOutArray* __return_Ux){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		GoSlice_ data2;
		data2.data = __uxIn2->data;
		data2.len = __uxIn2->count;
		data2.cap = __uxIn2->count;
		GoSlice_ data3;
		data3.data = NULL;
		data3.len = 0;
		data3.cap = 0;
		GoUint32 result = SKY_coin_UxArray_Sub(&data, &data2, &data3);
		if( result == 0){
			__return_Ux->data = data3.data;
			__return_Ux->count = data3.len;
		}
		return result;
	}
}

%rename(SKY_coin_UxArray_Add) wrap_SKY_coin_UxArray_Add;
%inline{
	GoUint32 wrap_SKY_coin_UxArray_Add(coin_UxOutArray* __uxIn, coin_UxOutArray* __uxIn2, coin_UxOutArray* __return_Ux){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		GoSlice_ data2;
		data2.data = __uxIn2->data;
		data2.len = __uxIn2->count;
		data2.cap = __uxIn2->count;
		GoSlice_ data3;
		data3.data = NULL;
		data3.len = 0;
		data3.cap = 0;
		GoUint32 result = SKY_coin_UxArray_Add(&data, &data2, &data3);
		if( result == 0){
			__return_Ux->data = data3.data;
			__return_Ux->count = data3.len;
		}
		return result;
	}
}

%rename(SKY_coin_NewAddressUxOuts) wrap_SKY_coin_NewAddressUxOuts;
%inline{ 
	GoUint32 wrap_SKY_coin_NewAddressUxOuts(coin_UxOutArray* __uxIn,  AddressUxOuts_Handle* p1){
		coin__UxArray data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_NewAddressUxOuts(&data, p1);
	}
}

%rename(SKY_coin_UxArray_Hashes) wrap_SKY_coin_UxArray_Hashes;
%inline{ 
	GoUint32 wrap_SKY_coin_UxArray_Hashes(coin_UxOutArray* __uxIn,  cipher_SHA256s* __out_hashes){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		GoSlice_ dataOut;
		dataOut.data = NULL;
		dataOut.len = 0;
		dataOut.cap = 0;
		GoUint32 result = SKY_coin_UxArray_Hashes(&data, &dataOut);
		if(result == 0){
			__out_hashes->data = dataOut.data;
			__out_hashes->count = dataOut.len;
		}
		return result;
	}
}

%rename(SKY_coin_AddressUxOuts_Flatten) wrap_SKY_coin_AddressUxOuts_Flatten;
%inline{ 
	GoUint32 wrap_SKY_coin_AddressUxOuts_Flatten(AddressUxOuts_Handle p0, coin_UxOutArray* __return_Ux){
		GoSlice_ data;
		data.data = NULL;
		data.len = 0;
		data.cap = 0;
		GoUint32 result = SKY_coin_AddressUxOuts_Flatten(p0, &data);
		if( result == 0 ){
			__return_Ux->data = data.data;
			__return_Ux->count = data.len;
		}
		return result;
	}
}

%rename(SKY_coin_AddressUxOuts_Get) wrap_SKY_coin_AddressUxOuts_Get;
%inline{ 
	GoUint32 wrap_SKY_coin_AddressUxOuts_Get(AddressUxOuts_Handle p0, cipher__Address* p1, coin_UxOutArray* __return_Ux){
		GoSlice_ data;
		data.data = NULL;
		data.len = 0;
		data.cap = 0;
		GoUint32 result = SKY_coin_AddressUxOuts_Get(p0, p1, &data);
		if( result == 0 ){
			__return_Ux->data = data.data;
			__return_Ux->count = data.len;
		}
		return result;
	}
}

%rename(SKY_coin_AddressUxOuts_Set) wrap_SKY_coin_AddressUxOuts_Set;
%inline{ 
	GoUint32 wrap_SKY_coin_AddressUxOuts_Set(AddressUxOuts_Handle p0, cipher__Address* p1, coin_UxOutArray* __uxIn){
		coin__UxArray data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_AddressUxOuts_Set(p0, p1, &data);
	}
}

%rename(SKY_coin_AddressUxOuts_Keys) wrap_SKY_coin_AddressUxOuts_Keys;
%inline{ 
	GoUint32 wrap_SKY_coin_AddressUxOuts_Keys(AddressUxOuts_Handle p0, cipher_Addresses* __out_addresses){
		coin__UxArray data;
		data.data = NULL;
		data.len = 0;
		data.cap = 0;
		GoUint32 result = SKY_coin_AddressUxOuts_Keys(p0, &data);
		if( result == 0){
			__out_addresses->data = data.data;
			__out_addresses->count = data.len;
		}
		return result;
	}
}

%rename(SKY_coin_Transactions_Hashes) wrap_SKY_coin_Transactions_Hashes;
%inline{
	GoUint32 wrap_SKY_coin_Transactions_Hashes(Transactions__Handle p0, cipher_SHA256s* __out_hashes){
		GoSlice_ data;
		data.data = NULL;
		data.len = 0;
		data.cap = 0;
		GoUint32 result = SKY_coin_Transactions_Hashes(p0, &data);
		if( result == 0){
			__out_hashes->data = data.data;
			__out_hashes->count = data.len;
		}
		return result;
	}
}

%rename(SKY_fee_TransactionFee) wrap_SKY_fee_TransactionFee;
%inline{
	GoUint32 wrap_SKY_fee_TransactionFee(Transaction__Handle __txn, GoUint64 __p1, coin_UxOutArray*  __uxIn, GoUint64  *__return_fee ){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_fee_TransactionFee(__txn,__p1, &data,__return_fee);
	}
}
