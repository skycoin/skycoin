%extend cipher__Address {
	int isEqual(cipher__Address* a){
		if( $self-> Version == a->Version ){
			return memcmp($self->Key, a->Key, sizeof(a->Key)) == 0;
		}
		return 0;
	}
}

%extend cipher_SecKeys {
	cipher_SecKey* getAt(int i){
		if( i < $self->count ){
			return &$self->data[i];
		}
		else 
			return NULL;
	}
}

