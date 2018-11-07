%inline {
	int equalSlices(GoSlice* slice1, GoSlice* slice2, int elem_size){
	  if(slice1->len != slice2->len)
		return 0;
	  return memcmp(slice1->data, slice2->data, slice1->len * elem_size) == 0;
	}
	int equalTransactions(coin__Transaction* t1, coin__Transaction* t2){
		if( t1->Length != t2->Length || t1->Type != t2->Type ){
			return 0;
		}
		if( memcmp(&t1->InnerHash, &t2->InnerHash, sizeof(cipher__SHA256)) != 0 )
			return 0;
		if(!equalSlices((GoSlice*)&t1->Sigs, (GoSlice*)&t2->Sigs, sizeof(cipher__Sig)))
			return 0;
		if(!equalSlices((GoSlice*)&t1->In, (GoSlice*)&t2->In, sizeof(cipher__SHA256)))
			return 0;
		if(!equalSlices((GoSlice*)&t1->Out, (GoSlice*)&t2->Out, sizeof(coin__TransactionOutput)))
			return 0;
		return 1;
	}
	int equalTransactionsArrays(coin__Transactions* pTxs1, coin__Transactions* pTxs2){
		if( pTxs1->len != pTxs2->len )
			return 0;
		coin__Transaction* pTx1 = pTxs1->data;
		coin__Transaction* pTx2 = pTxs2->data;
		int i;
		for(i = 0; i < pTxs1->len; i++){
			if(!equalTransactions(pTx1, pTx2))
				return 0;
			pTx1++;
			pTx2++;
		}
		return 1;
	}
	int equalBlockHeaders(coin__BlockHeader* bh1, coin__BlockHeader* bh2){
		if( bh1->Version != bh2->Version || bh1->Time != bh2->Time || 
			bh1->BkSeq != bh2->BkSeq || bh1->Fee != bh2->Fee)
			return 0;
		if( memcmp( &bh1->PrevHash, bh2->PrevHash, sizeof(bh2->PrevHash) ) != 0 )
			return 0;
		if( memcmp( &bh1->BodyHash, bh2->PrevHash, sizeof(bh2->BodyHash) ) != 0 )
			return 0;
		if( memcmp( &bh1->UxHash, bh2->PrevHash, sizeof(bh2->UxHash) ) != 0 )
			return 0;
		return 1;
	}
}
