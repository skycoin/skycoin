%rename(SKY_cipher_SumSHA256) Java_skycoin_libjava_skycoinJNI_SKY_cipher_SumSHA256;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_SumSHA256(GoSlice seed, cipher_SHA256* sha){
		GoUint32 result = SKY_cipher_SumSHA256(seed,  sha);
		return result;
	}
}

%rename(SKY_cipher_SignHash) Java_skycoin_libjava_skycoinJNI_SKY_cipher_SignHash;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_SignHash(cipher_SHA256 *sha,cipher_SecKey *sec,cipher_Sig *s){
		GoUint32 result = SKY_cipher_SignHash(sha,sec,s);
		return result;
	}
}

%rename(SKY_cipher_ChkSig) Java_skycoin_libjava_skycoinJNI_SKY_cipher_ChkSig;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_ChkSig(cipher__Address *a,cipher_SHA256 *sha,cipher_Sig *s){
		GoUint32 result = SKY_cipher_ChkSig(a,sha,s);
		return result;
	}
}

%rename(SKY_cipher_PubKeyFromSig) Java_skycoin_libjava_skycoinJNI_SKY_cipher_PubKeyFromSig;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_PubKeyFromSig(cipher_Sig *sig,cipher_SHA256 *h,cipher_PubKey *p){
		GoUint32 result = SKY_cipher_PubKeyFromSig(sig,h,p);
		return result;
	}
}

%rename(SKY_cipher_VerifySignature) Java_skycoin_libjava_skycoinJNI_SKY_cipher_VerifySignature;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_VerifySignature(cipher_PubKey *p,cipher_Sig *sig,cipher_SHA256 *h){
		GoUint32 result = SKY_cipher_VerifySignature(p,sig,h);
		return result;
	}
}

// %rename(SKY_cipher_TestSecKeyHash) Java_skycoin_libjava_skycoinJNI_SKY_cipher_TestSecKeyHash;
// %inline {
// 	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_TestSecKeyHash(cipher_SecKey *s,cipher_SHA256 *h){
// 		GoUint32 result = SKY_cipher_TestSecKeyHash(s,h);
// 		return result;
// 	}
// }

%rename(SKY_cipher_SHA256_Set) Java_skycoin_libjava_skycoinJNI_SKY_cipher_SHA256_Set;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_SHA256_Set(cipher_SHA256 *h,GoSlice s){
		GoUint32 result = SKY_cipher_SHA256_Set(h,s);
		return result;
	}
}

%rename(SKY_cipher_SHA256_Hex) Java_skycoin_libjava_skycoinJNI_SKY_cipher_SHA256_Hex;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_SHA256_Hex(cipher_SHA256 *h,GoString_* s){
		GoUint32 result = SKY_cipher_SHA256_Hex(h,s);
		return result;
	}
}

%rename(SKY_cipher_SHA256FromHex) Java_skycoin_libjava_skycoinJNI_SKY_cipher_SHA256FromHex;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_SHA256FromHex(GoString s,cipher_SHA256 *h){
		GoUint32 result = SKY_cipher_SHA256FromHex(s,h);
		return result;
	}
}

%rename(SKY_coin_Transaction_HashInner) Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_HashInner;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_HashInner(Transaction__Handle tx,cipher_SHA256 *h){
		GoUint32 result = SKY_coin_Transaction_HashInner(tx,h);
		return result;
	}
}

%rename(SKY_coin_Transaction_Hash) Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_Hash;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_Hash(Transaction__Handle tx,cipher_SHA256 *h){
		GoUint32 result = SKY_coin_Transaction_Hash(tx,h);
		return result;
	}
}

%rename(SKY_coin_Transaction_SetInputAt) Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_SetInputAt;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_SetInputAt(Transaction__Handle tx,GoInt p1,cipher_SHA256 *h){
		GoUint32 result = SKY_coin_Transaction_SetInputAt(tx,p1,h);
		return result;
	}
}



%rename(SKY_coin_Transaction_GetInputAt) Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_GetInputAt;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_GetInputAt(Transaction__Handle tx, GoInt p1,cipher_SHA256 *h){
		GoUint32 result = SKY_coin_Transaction_GetInputAt(tx,p1,h);
		return result;
	}
}

%rename(SKY_coin_Transaction_PushInput) Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_PushInput;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_PushInput(Transaction__Handle tx, cipher_SHA256* h, GoUint16* p1){
		GoUint32 result = SKY_coin_Transaction_PushInput(tx,h,p1);
		return result;
	}
}

%rename(SKY_coin_Transaction_SignInputs) Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_SignInputs;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_SignInputs(Transaction__Handle handle, cipher_SecKeys* __in_pubKeys){
		GoSlice data;
		data.data = __in_pubKeys->data;
		data.len = __in_pubKeys->count;
		data.cap = __in_pubKeys->count;
		return SKY_coin_Transaction_SignInputs(handle, data);
	}
}

%rename(SKY_cipher_GenerateDeterministicKeyPairs) Java_skycoin_libjava_skycoinJNI_SKY_cipher_GenerateDeterministicKeyPairs;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_GenerateDeterministicKeyPairs(GoSlice seed, GoInt n, cipher_SecKeys* __out_secKeys){
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

%rename(SKY_cipher_GenerateDeterministicKeyPairsSeed) Java_skycoin_libjava_skycoinJNI_SKY_cipher_GenerateDeterministicKeyPairsSeed;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_GenerateDeterministicKeyPairsSeed(GoSlice seed, GoInt n, coin__UxArray* newSeed, cipher_SecKeys* __out_secKeys){
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


// %rename(SKY_cipher_PubKeySlice_Swap) Java_skycoin_libjava_skycoinJNI_SKY_cipher_PubKeySlice_Swap;
// %inline {
// 	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_PubKeySlice_Swap(cipher_PubKeys* __in_pubKeys, GoInt p1, GoInt p2){
// 		GoSlice_ data;
// 		data.data = __in_pubKeys->data;
// 		data.len = __in_pubKeys->count;
// 		data.cap = __in_pubKeys->count;
// 		GoUint32 result = SKY_cipher_PubKeySlice_Swap(&data, p1, p2);
// 		return result;
// 	}
// }

%rename(SKY_coin_VerifyTransactionCoinsSpending) Java_skycoin_libjava_skycoinJNI_SKY_coin_VerifyTransactionCoinsSpending;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_VerifyTransactionCoinsSpending(coin_UxOutArray* __uxIn, coin_UxOutArray* __uxOut){
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

%rename(SKY_coin_VerifyTransactionHoursSpending) Java_skycoin_libjava_skycoinJNI_SKY_coin_VerifyTransactionHoursSpending;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_VerifyTransactionHoursSpending(GoUint64 _headTime , coin_UxOutArray* __uxIn, coin_UxOutArray* __uxOut){
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

%rename(SKY_coin_CreateUnspents) Java_skycoin_libjava_skycoinJNI_SKY_coin_CreateUnspents;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_CreateUnspents(coin__BlockHeader* bh, Transaction__Handle t, coin_UxOutArray* __return_Ux){
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

%rename(SKY_coin_Transaction_VerifyInput) Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_VerifyInput;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_Transaction_VerifyInput(Transaction__Handle handle, coin_UxOutArray* __uxIn){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_Transaction_VerifyInput(handle, &data);
	}
}

%rename(SKY_coin_UxArray_HasDupes) Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_HasDupes;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_HasDupes(coin_UxOutArray* __uxIn, GoUint8* p1){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_UxArray_HasDupes(&data, p1);
	}
}

%rename(SKY_coin_UxArray_Coins) Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_Coins;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_Coins(coin_UxOutArray* __uxIn, GoUint64* p1){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_UxArray_Coins(&data, p1);
	}
}

%rename(SKY_coin_UxArray_CoinHours) Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_CoinHours;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_CoinHours(coin_UxOutArray* __uxIn, GoUint64 p1, GoUint64* p2){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_UxArray_CoinHours(&data, p1, p2);
	}
}

%rename(SKY_coin_UxArray_Less) Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_Less;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_Less(coin_UxOutArray* __uxIn, GoInt p1, GoInt p2, GoUint8* p3){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_UxArray_Less(&data, p1, p2, p3);
	}
}

%rename(SKY_coin_UxArray_Swap) Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_Swap;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_Swap(coin_UxOutArray* __uxIn, GoInt p1, GoInt p2){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_UxArray_Swap(&data, p1, p2);
	}
}

%rename(SKY_coin_UxArray_Sub) Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_Sub;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_Sub(coin_UxOutArray* __uxIn, coin_UxOutArray* __uxIn2, coin_UxOutArray* __return_Ux){
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

%rename(SKY_coin_UxArray_Add) Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_Add;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_Add(coin_UxOutArray* __uxIn, coin_UxOutArray* __uxIn2, coin_UxOutArray* __return_Ux){
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

%rename(SKY_coin_NewAddressUxOuts) Java_skycoin_libjava_skycoinJNI_SKY_coin_NewAddressUxOuts;
%inline{ 
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_NewAddressUxOuts(coin_UxOutArray* __uxIn,  AddressUxOuts_Handle* p1){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_NewAddressUxOuts(&data, p1);
	}
}

%rename(SKY_coin_UxArray_Hashes) Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_Hashes;
%inline{ 
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_UxArray_Hashes(coin_UxOutArray* __uxIn,  cipher_SHA256s* __out_hashes){
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

%rename(SKY_coin_AddressUxOuts_Flatten) Java_skycoin_libjava_skycoinJNI_SKY_coin_AddressUxOuts_Flatten;
%inline{ 
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_AddressUxOuts_Flatten(AddressUxOuts_Handle p0, coin_UxOutArray* __return_Ux){
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

%rename(SKY_coin_AddressUxOuts_Get) Java_skycoin_libjava_skycoinJNI_SKY_coin_AddressUxOuts_Get;
%inline{ 
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_AddressUxOuts_Get(AddressUxOuts_Handle p0, cipher__Address* p1, coin_UxOutArray* __return_Ux){
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

%rename(SKY_coin_AddressUxOuts_Set) Java_skycoin_libjava_skycoinJNI_SKY_coin_AddressUxOuts_Set;
%inline{ 
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_AddressUxOuts_Set(AddressUxOuts_Handle p0, cipher__Address* p1, coin_UxOutArray* __uxIn){
		GoSlice_ data;
		data.data = __uxIn->data;
		data.len = __uxIn->count;
		data.cap = __uxIn->count;
		return SKY_coin_AddressUxOuts_Set(p0, p1, &data);
	}
}

%rename(SKY_coin_AddressUxOuts_Keys) Java_skycoin_libjava_skycoinJNI_SKY_coin_AddressUxOuts_Keys;
%inline{ 
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_AddressUxOuts_Keys(AddressUxOuts_Handle p0, cipher_Addresses* __out_addr){
		GoSlice_ data;
		data.data = NULL;
		data.len = 0;
		data.cap = 0;
		GoUint32 result = SKY_coin_AddressUxOuts_Keys(p0, &data);
		if( result == 0){
			__out_addr->data = data.data;
			__out_addr->count = data.len;
		}
		return result;
	}
}

%rename(SKY_coin_Transactions_Hashes) Java_skycoin_libjava_skycoinJNI_SKY_coin_Transactions_Hashes;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_Transactions_Hashes(Transactions__Handle p0, cipher_SHA256s* __out_hashes){
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

%rename(SKY_coin_UxOut_Hash) Java_skycoin_libjava_skycoinJNI_SKY_coin_UxOut_Hash;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_UxOut_Hash(coin__UxOut* ux, cipher_SHA256* sha){
		GoUint32 result = SKY_coin_UxOut_Hash(ux, sha);
		return result;
	}
}

%rename(SKY_cipher_AddSHA256) Java_skycoin_libjava_skycoinJNI_SKY_cipher_AddSHA256;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_AddSHA256(cipher_SHA256* p0, cipher_SHA256* p1, cipher_SHA256* p2){
		GoUint32 result = SKY_cipher_AddSHA256(p0,p1,p2);
		return result;
	}
}

%rename(SKY_coin_GetTransactionObject) Java_skycoin_libjava_skycoinJNI_SKY_coin_GetTransactionObject;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_GetTransactionObject(Transaction__Handle tx, coin__Transaction *p1){
		GoUint32 result = SKY_coin_GetTransactionObject(tx,&p1);
		return result;
	}
}

%rename(SKY_coin_UxBody_Hash) Java_skycoin_libjava_skycoinJNI_SKY_coin_UxBody_Hash;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_UxBody_Hash(coin__UxBody* p0, cipher_SHA256* p1){
		GoUint32 result = SKY_coin_UxBody_Hash(p0,p1);
		return result;
	}
}

%rename(SKY_coin_UxOut_SnapshotHash) Java_skycoin_libjava_skycoinJNI_SKY_coin_UxOut_SnapshotHash;
%inline{
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_UxOut_SnapshotHash(coin__UxOut* p0, cipher_SHA256* p1){
		GoUint32 result = SKY_coin_UxOut_SnapshotHash(p0,p1);
		return result;
	}
}

%rename(SKY_fee_TransactionFee) Java_skycoin_libjava_skycoinJNI_SKY_fee_TransactionFee;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_fee_TransactionFee(Transaction__Handle handle , GoUint64 p1,coin_UxOutArray* __uxIn, GoUint64* p3){
		GoSlice_ dataIn;
		dataIn.data = __uxIn->data;
		dataIn.len = __uxIn->count;
		dataIn.cap = __uxIn->count;
		GoUint32 result = SKY_fee_TransactionFee(handle, p1,&dataIn,p3);
		return result;
	};
}


%rename(SKY_cipher_CheckSecKeyHash) Java_skycoin_libjava_skycoinJNI_SKY_cipher_CheckSecKeyHash;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_cipher_CheckSecKeyHash(cipher_SecKey *s, cipher_SHA256* sha){
		GoUint32 result = SKY_cipher_CheckSecKeyHash(s,  sha);
		return result;
	}
}

%rename(SKY_coin_NewBlock) Java_skycoin_libjava_skycoinJNI_SKY_coin_NewBlock;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_NewBlock(Block__Handle p0, GoUint64 p1, cipher_SHA256* p2, Transactions__Handle p3, FeeCalculator* p4, Block__Handle* p5){
		GoUint32 result = SKY_coin_NewBlock(p0,  p1,p2,p3,p4,p5);
		return result;
	}
}

%rename(SKY_coin_Block_HashHeader) Java_skycoin_libjava_skycoinJNI_SKY_coin_Block_HashHeader;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_Block_HashHeader(Block__Handle p0, cipher_SHA256* p1){
		GoUint32 result = SKY_coin_Block_HashHeader(p0,  p1);
		return result;
	}
}

%rename(SKY_coin_Block_PreHashHeader) Java_skycoin_libjava_skycoinJNI_SKY_coin_Block_PreHashHeader;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_Block_PreHashHeader(Block__Handle p0, cipher_SHA256* p1){
		GoUint32 result = SKY_coin_Block_PreHashHeader(p0,  p1);
		return result;
	}
}

%rename(SKY_coin_BlockBody_Hash) Java_skycoin_libjava_skycoinJNI_SKY_coin_BlockBody_Hash;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_BlockBody_Hash(BlockBody__Handle p0, cipher_SHA256* p1){
		GoUint32 result = SKY_coin_BlockBody_Hash(p0,  p1);
		return result;
	}
}

%rename(SKY_coin_BlockHeader_Hash) Java_skycoin_libjava_skycoinJNI_SKY_coin_BlockHeader_Hash;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_BlockHeader_Hash(coin__BlockHeader* p0, cipher_SHA256* p1){
		GoUint32 result = SKY_coin_BlockHeader_Hash(p0,  p1);
		return result;
	}
}

%rename(SKY_coin_Block_HashBody) Java_skycoin_libjava_skycoinJNI_SKY_coin_Block_HashBody;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_coin_Block_HashBody(Block__Handle p0, cipher_SHA256* p1){
		GoUint32 result = SKY_coin_Block_HashBody(p0,  p1);
		return result;
	}
}