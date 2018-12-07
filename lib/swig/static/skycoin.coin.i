%extend coin__BlockHeader {
	int isEqual(coin__BlockHeader* bh){
		return equalBlockHeaders($self, bh);
	}
}

%extend coin__Transaction {
	int isEqual(coin__Transaction* t){
		return equalTransactions($self, t);
	}
	cipher_SHA256 GetInnerHash(){
		cipher_SHA256 h;
		cipher_SHA256_assignFrom(&h,$self->InnerHash);
		return h;
	}
}

%extend coin__BlockBody {
	int isEqual(coin__BlockBody* b){
		return equalTransactionsArrays(&$self->Transactions, &b->Transactions);
	}
}

%extend coin__UxOut {
	int isEqual(coin__UxOut* u){
		return memcmp(&$self, u, sizeof(coin__UxOut)) == 0;
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

	void append(coin__UxOut* uxout){
		int n = $self->count+1;
				$self->data = malloc(n * sizeof(*($self->data)));
		$self->count =n ;
		memcpy(&self->data[n-1], uxout, sizeof(*uxout));

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

%extend coin__UxBody {
	void SetSrcTransaction(cipher_SHA256 *o){
			cipher_SHA256* p = (cipher_SHA256*)o;
			memcpy( &$self->SrcTransaction, &p->data, sizeof(cipher__SHA256));
		}
	
}
