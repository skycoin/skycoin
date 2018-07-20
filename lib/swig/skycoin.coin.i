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
		return memcmp(&$self, u, sizeof(coin__UxOut)) == 0;
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
