

%extend cipher__Address{
	int __eq__(cipher__Address* a){
		if( $self->Version == a->Version ){
			return memcmp($self->Key, a->Key, sizeof(a->Key)) == 0;
		}
		return 0;
	}
	PyObject* toStr(){
		return PyBytes_FromStringAndSize((const char*)$self->Key, sizeof($self->Key));
	}
}

%extend cipher_PubKey {
	int __eq__(cipher_PubKey* a){
		return memcmp($self->data, a->data, sizeof(a->data)) == 0;
	}
	int compareToString(PyObject * str){
		char* s = SWIG_Python_str_AsChar(str);
		int result = memcmp(s, $self->data, sizeof($self->data));
		SWIG_Python_str_DelForPy3(s);
		return result;
	}
	PyObject* toStr(){
		return PyBytes_FromStringAndSize((const char*)$self->data, sizeof($self->data));
	}
	void assignFrom(void* data){
		memcpy(&$self->data, data, sizeof($self->data));
	}
	void assignTo(void* data){
		memcpy(data, &$self->data, sizeof($self->data));
	}
}

%extend cipher_SecKey {
	int __eq__(cipher_SecKey* a){
		return memcmp($self->data, a->data, sizeof(a->data)) == 0;
	}
	int compareToString(PyObject * str){
		char* s = SWIG_Python_str_AsChar(str);
		int result = memcmp(s, $self->data, sizeof($self->data));
		SWIG_Python_str_DelForPy3(s);
		return result;
	}
	PyObject* toStr(){
		return PyBytes_FromStringAndSize((const char*)$self->data, sizeof($self->data));
	}
	void assignFrom(void* data){
		memcpy(&$self->data, data, sizeof($self->data));
	}
	void assignTo(void* data){
		memcpy(data, &$self->data, sizeof($self->data));
	}
}

%extend cipher_Ripemd160 {
	int __eq__(cipher_Ripemd160* a){
		return memcmp($self->data, a->data, sizeof(a->data)) == 0;
	}
	int compareToString(PyObject * str){
		char* s = SWIG_Python_str_AsChar(str);
		int result = memcmp(s, $self->data, sizeof($self->data));
		SWIG_Python_str_DelForPy3(s);
		return result;
	}
	PyObject* toStr(){
		return PyBytes_FromStringAndSize((const char*)$self->data, sizeof($self->data));
	}
	void assignFrom(void* data){
		memcpy(&$self->data, data, sizeof($self->data));
	}
	void assignTo(void* data){
		memcpy(data, &$self->data, sizeof($self->data));
	}
}

%extend cipher_Sig {
	int __eq__(cipher_Sig* a){
		return memcmp($self->data, a->data, sizeof(a->data)) == 0;
	}
	int compareToString(PyObject * str){
		char* s = SWIG_Python_str_AsChar(str);
		int result = memcmp(s, $self->data, sizeof($self->data));
		SWIG_Python_str_DelForPy3(s);
		return result;
	}
	PyObject* toStr(){
		return PyBytes_FromStringAndSize((const char*)$self->data, sizeof($self->data));
	}
	void assignFrom(void* data){
		memcpy(&$self->data, data, sizeof($self->data));
	}
	void assignTo(void* data){
		memcpy(data, &$self->data, sizeof($self->data));
	}
}

%extend cipher_SHA256 {
	int __eq__(cipher_SHA256* a){
		return memcmp($self->data, a->data, sizeof(a->data)) == 0;
	}
	int compareToString(PyObject * str){
		char* s = SWIG_Python_str_AsChar(str);
		int result = memcmp(s, $self->data, sizeof($self->data));
		SWIG_Python_str_DelForPy3(s);
		return result;
	}
	PyObject* toStr(){
		return PyBytes_FromStringAndSize((const char*)$self->data, sizeof($self->data));
	}
	void assignFrom(void* data){
		memcpy(&$self->data, data, sizeof($self->data));
	}
	void assignTo(void* data){
		memcpy(data, &$self->data, sizeof($self->data));
	}
}

%extend cipher_Checksum {
	int __eq__(cipher_Checksum* a){
		return memcmp($self->data, a->data, sizeof(a->data)) == 0;
	}
	int compareToString(PyObject * str){
		char* s = SWIG_Python_str_AsChar(str);
		int result = memcmp(s, $self->data, sizeof($self->data));
		SWIG_Python_str_DelForPy3(s);
		return result;
	}
	PyObject* toStr(){
		return PyBytes_FromStringAndSize((const char*)$self->data, sizeof($self->data));
	}
	void assignFrom(void* data){
		memcpy(&$self->data, data, sizeof($self->data));
	}
	void assignTo(void* data){
		memcpy(data, &$self->data, sizeof($self->data));
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
	
	int setAt(int i, cipher_SecKey* seckey){
		if( i < $self->count){
			memcpy(&$self->data[i], seckey, sizeof(*seckey));
			return i;
		} else {
			return -1;
		}
	}
	
	int __eq__(cipher_SecKeys* a){
		return $self->count == a->count && memcmp($self->data, a->data, sizeof(cipher_SecKey) * $self->count) == 0;
	}
	
	void allocate(int n){
		$self->data = malloc(n * sizeof(*($self->data)));
		$self->count = n;
	}
	
	void release(){
		destroy_cipher_SecKeys($self);
	}
}

%inline{
	void destroy_cipher_SecKeys(cipher_SecKeys* p){
		if( p != NULL ){
			if( p->data != NULL ){
				free( p->data );
			}
		}
	}
}

%extend cipher_PubKeys {
	cipher_PubKey* getAt(int i){
		if( i < $self->count ){
			return &$self->data[i];
		}
		else
			return NULL;
	}
	
	int setAt(int i, cipher_PubKey* pubkey){
		if( i < $self->count){
			memcpy(&self->data[i], pubkey, sizeof(*pubkey));
			return i;
		} else {
			return -1;
		}
	}
	
	int __eq__(cipher_PubKeys* a){
		return $self->count == a->count && memcmp($self->data, a->data, sizeof(cipher_PubKey) * $self->count) == 0;
	}
	
	void allocate(int n){
		$self->data = malloc(n * sizeof(*($self->data)));
		$self->count = n;
	}
	
	void release(){
		destroy_cipher_PubKeys($self);
	}
}

%extend cipher_SHA256s {
	cipher_SHA256* getAt(int i){
		if( i < $self->count ){
			return &$self->data[i];
		}
		else
			return NULL;
	}
	
	int setAt(int i, cipher_SHA256* hash){
		if( i < $self->count){
			memcpy(&self->data[i], hash, sizeof(*hash));
			return i;
		} else {
			return -1;
		}
	}
	
	int __eq__(cipher_SHA256s* a){
		return $self->count == a->count && memcmp($self->data, a->data, sizeof(cipher_SHA256) * $self->count) == 0;
	}
	
	void allocate(int n){
		$self->data = malloc(n * sizeof(*($self->data)));
		$self->count = n;
	}
	
	void release(){
		if($self->data != NULL) free($self->data);
	}
}

%inline{
	void destroy_cipher_PubKeys(cipher_PubKeys* p){
		if( p != NULL ){
			if( p->data != NULL ){
				free( p->data );
			}
		}
	}
}


%extend cipher__BitcoinAddress{
	int __eq__(cipher__BitcoinAddress* a){
		if( $self->Version == a->Version ){
			return memcmp($self->Key, a->Key, sizeof(a->Key)) == 0;
		}
		return 0;
	}
	PyObject* toStr(){
		return PyBytes_FromStringAndSize((const char*)$self->Key, sizeof($self->Key));
	}
}
