%extend cipher__Address {
	char getVersion(){
		return $self->Version;
    }
    	void setVersion(char pValue){
		$self->Version = pValue;
    }
}

%extend cipher_Sig{

GoSlice toGoSlice(){
	GoSlice slice;
slice.len = sizeof(cipher_Sig);
slice.cap = sizeof(cipher_Sig)+1;
slice.data = (cipher_Sig*)&$self;
return slice;
	}
}
%extend GoSlice {
	int isEqual(GoSlice *slice){
		return (($self->len == slice->len)) && (memcmp($self->data,slice->data, sizeof(GoSlice_))==0 );
	}

	void convertString(_GoString_ data){
		$self->data = data.p;
		$self->len = strlen(data.p);
		$self->cap = $self->len;
	}

	void setcoin_UxOut(coin__UxOut ux,int p){
		coin__UxOut *tmpOut = (coin__UxOut *)$self->data;
		for (int j = 0; j < $self->len; j++) {
      if(p == j){
      tmpOut->Head.Time = ux.Head.Time;
      tmpOut->Body.Coins = ux.Body.Coins;
      tmpOut->Body.Hours = ux.Body.Hours;
    }
	      tmpOut++;}
	}

	
}

%extend _GoString_ {
	int SetString(char * str){
		$self->p = str;
		$self->n = strlen(str);
	}
}

%extend cipher_SHA256 {
    	_GoString_ getStr(){
		_GoString_ str;
		str.p = (const char*)$self->data;
		str.n = strlen(str.p);
		return str;
    }
}

%extend coin__Transaction {
    int	setInnerHash(cipher_SHA256 h){
		memset($self->InnerHash, 0, sizeof(cipher__SHA256));
			cipher_SHA256_assignFrom($self->InnerHash,&h);
			return 0;
    }
}

%extend Transaction__Handle {
	char getVersion(){
		return $self->Version;
    }
    	void setVersion(char pValue){
		$self->Version = pValue;
    }
}