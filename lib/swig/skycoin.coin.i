%include "cmp.i"

%extend coin__BlockHeader {
	int isEqual(coin__BlockHeader* bh){
		return equalBlockHeaders($self, bh);
	}
}

%extend coin__Transaction {
	int isEqual(coin__Transaction* t){
		return equalTransactions($self, t);
	}
}

%extend coin__BlockBody {
	int isEqual(coin__BlockBody* b){
		return equalTransactionsArrays(&$self->Transactions, &b->Transactions);
	}
}

%extend coin__UxOut {
	int isEqual(coin__UxOut* u){
		if($self->Head.Time != u->Head.Time)
			return 0;
		if($self->Head.BkSeq != u->Head.BkSeq)
			return 0;
		if($self->Body.Coins != u->Body.Coins)
			return 0;
		if($self->Body.Hours != u->Body.Hours)
			return 0;
		if(memcmp(&$self->Body.Address, &u->Body.Address, sizeof(cipher__Address)) != 0)
			return 0;
		if(memcmp(&$self->Body.SrcTransaction, &u->Body.SrcTransaction, sizeof(cipher__SHA256)) != 0)
			return 0;
		return 1;
	}
}

%extend coin_UxOutArray {
	coin__UxOut* getAt(int i){
		if( i < $self->count ){
			return &$self->data[i];
		}
		else
			return NULL;
	}
	
	int setAt(int i, coin__UxOut* uxout){
		if( i < $self->count){
			memcpy(&self->data[i], uxout, sizeof(*uxout));
			return i;
		} else {
			return -1;
		}
	}
	
	int isEqual(coin_UxOutArray* a){
		return $self->count == a->count && memcmp($self->data, a->data, sizeof(coin__UxOut) * $self->count) == 0;
	}
	
	void allocate(int n){
		$self->data = malloc(n * sizeof(*($self->data)));
		$self->count = n;
	}
	
	void release(){
		if($self->data != NULL)
			free($self->data);
	}
}

%extend coin__TransactionOutput {
	int isEqual(coin__TransactionOutput* t){
		if( $self->Coins != t->Coins ||
			$self->Hours != t->Hours ){
			return 0;
	  	}

	  	if(memcmp(&$self->Address, &t->Address, sizeof(cipher__Address)) != 0)
			return 0;
	 	return 1;
	}
}

