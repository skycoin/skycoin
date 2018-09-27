

%inline %{
#include "json.h"
	//Define function SKY_handle_close to avoid including libskycoin.h
void SKY_handle_close(Handle p0);

int MEMPOOLIDX = 0;
void *MEMPOOL[1024 * 256];

int JSONPOOLIDX = 0;
json_value* JSON_POOL[128];

int HANDLEPOOLIDX = 0;
Handle HANDLE_POOL[128];

typedef struct {
	Client__Handle client;
	WalletResponse__Handle wallet;
} wallet_register;

int WALLETPOOLIDX = 0;
wallet_register WALLET_POOL[64];

int stdout_backup;
int pipefd[2];

void * registerMemCleanup(void *p) {
	int i;
	for (i = 0; i < MEMPOOLIDX; i++) {
		if(MEMPOOL[i] == NULL){
			MEMPOOL[i] = p;
			return p;
		}
	}
	MEMPOOL[MEMPOOLIDX++] = p;
	return p;
}

void freeRegisteredMemCleanup(void *p){
	int i;
	for (i = 0; i < MEMPOOLIDX; i++) {
		if(MEMPOOL[i] == p){
			free(p);
			MEMPOOL[i] = NULL;
			break;
		}
	}
}

int registerJsonFree(void *p){
	int i;
	for (i = 0; i < JSONPOOLIDX; i++) {
		if(JSON_POOL[i] == NULL){
			JSON_POOL[i] = p;
			return i;
		}
	}
	JSON_POOL[JSONPOOLIDX++] = p;
	return JSONPOOLIDX-1;
}

void freeRegisteredJson(void *p){
	int i;
	for (i = 0; i < JSONPOOLIDX; i++) {
		if(JSON_POOL[i] == p){
			JSON_POOL[i] = NULL;
			json_value_free( (json_value*)p );
			break;
		}
	}
}

int registerWalletClean(Client__Handle clientHandle,
						WalletResponse__Handle walletHandle){
	int i;
	for (i = 0; i < WALLETPOOLIDX; i++) {
		if(WALLET_POOL[i].wallet == 0 && WALLET_POOL[i].client == 0){
			WALLET_POOL[i].wallet = walletHandle;
			WALLET_POOL[i].client = clientHandle;
			return i;
		}
	}
	WALLET_POOL[WALLETPOOLIDX].wallet = walletHandle;
	WALLET_POOL[WALLETPOOLIDX].client = clientHandle;
	return WALLETPOOLIDX++;
}

int registerHandleClose(Handle handle){
	int i;
	for (i = 0; i < HANDLEPOOLIDX; i++) {
		if(HANDLE_POOL[i] == 0){
			HANDLE_POOL[i] = handle;
			return i;
		}
	}
	HANDLE_POOL[HANDLEPOOLIDX++] = handle;
	return HANDLEPOOLIDX - 1;
}

void closeRegisteredHandle(Handle handle){
	int i;
	for (i = 0; i < HANDLEPOOLIDX; i++) {
		if(HANDLE_POOL[i] == handle){
			HANDLE_POOL[i] = 0;
			SKY_handle_close(handle);
			break;
		}
	}
}

void cleanupWallet(Client__Handle client, WalletResponse__Handle wallet){
	int result;
	GoString_ strWalletDir;
	GoString_ strFileName;
	memset(&strWalletDir, 0, sizeof(GoString_));
	memset(&strFileName, 0, sizeof(GoString_));


	result = SKY_api_Handle_Client_GetWalletDir(client, &strWalletDir);
	if( result != 0 ){
		return;
	}
	result = SKY_api_Handle_Client_GetWalletFileName(wallet, &strFileName);
	if( result != 0 ){
		free( (void*)strWalletDir.p );
		return;
	}
	char fullPath[128];
	if( strWalletDir.n + strFileName.n < 126){
		strcpy( fullPath, strWalletDir.p );
		if( fullPath[0] == 0 || fullPath[strlen(fullPath) - 1] != '/' )
			strcat(fullPath, "/");
		strcat( fullPath, strFileName.p );
		result = unlink( fullPath );
		if( strlen(fullPath) < 123 ){
			strcat( fullPath, ".bak" );
			result = unlink( fullPath );
		}
	}
	GoString str = { strFileName.p, strFileName.n };
	result = SKY_api_Client_UnloadWallet( client, str );
	GoString strFullPath = { fullPath, strlen(fullPath) };
	free( (void*)strWalletDir.p );
	free( (void*)strFileName.p );
}

void cleanRegisteredWallet(
			Client__Handle client,
			WalletResponse__Handle wallet){

	int i;
	for (i = 0; i < WALLETPOOLIDX; i++) {
		if(WALLET_POOL[i].wallet == wallet && WALLET_POOL[i].client == client){
			WALLET_POOL[i].wallet = 0;
			WALLET_POOL[i].client = 0;
			cleanupWallet( client, wallet );
			return;
		}
	}
}

void cleanupMem() {
	int i;

	for (i = 0; i < WALLETPOOLIDX; i++) {
		if(WALLET_POOL[i].client != 0 && WALLET_POOL[i].wallet != 0){
			cleanupWallet( WALLET_POOL[i].client, WALLET_POOL[i].wallet );
		}
	}

  void **ptr;
  for (i = MEMPOOLIDX, ptr = MEMPOOL; i; --i) {
	if( *ptr )
		free(*ptr);
	ptr++;
  }
  for (i = JSONPOOLIDX, ptr = (void*)JSON_POOL; i; --i) {
	if( *ptr )
		json_value_free(*ptr);
	ptr++;
  }
  for (i = 0; i < HANDLEPOOLIDX; i++) {
	  if( HANDLE_POOL[i] )
		SKY_handle_close(HANDLE_POOL[i]);
  }
}


void setup(void) {
	srand ((unsigned int) time (NULL));
}

void teardown(void) {
	cleanupMem();
}

// TODO: Move to libsky_io.c
void fprintbuff(FILE *f, void *buff, size_t n) {
  unsigned char *ptr = (unsigned char *) buff;
  fprintf(f, "[ ");
  for (; n; --n, ptr++) {
    fprintf(f, "%02d ", *ptr);
  }
  fprintf(f, "]");
}

int parseBoolean(const char* str, int length){
	int result = 0;
	if(length == 1){
		result = str[0] == '1' || str[0] == 't' || str[0] == 'T';
	} else {
		result = strncmp(str, "true", length) == 0 ||
			strncmp(str, "True", length) == 0 ||
			strncmp(str, "TRUE", length) == 0;
	}
	return result;
}

void toGoString(GoString_ *s, GoString *r){
GoString * tmp = r;

  *tmp = (*(GoString *) s);
}

int copySlice(GoSlice_* pdest, GoSlice_* psource, int elem_size){
  pdest->len = psource->len;
  pdest->cap = psource->len;
  int size = pdest->len * elem_size;
  pdest->data = malloc(size);
	if( pdest->data == NULL )
		return 1;
  registerMemCleanup( pdest->data );
  memcpy(pdest->data, psource->data, size );
	return 0;
}



int concatSlices(GoSlice_* slice1, GoSlice_* slice2, int elem_size, GoSlice_* result){
	int size1 = slice1->len;
	int size2 = slice2->len;
	int size = size1 + size2;
	if (size <= 0)
		return 1;
	void* data = malloc(size * elem_size);
	if( data == NULL )
		return 1;
	registerMemCleanup( data );
	result->data = data;
	result->len = size;
	result->cap = size;
	char* p = data;
	if(size1 > 0){
		memcpy( p, slice1->data, size1 * elem_size );
		p += (elem_size * size1);
	}
	if(size2 > 0){
		memcpy( p, slice2->data, size2 * elem_size );
	}
	return 0;
}
    void parseJsonMetaData(char *metadata, int *n, int *r, int *p, int *keyLen)
{
	*n = *r = *p = *keyLen = 0;
	int length = strlen(metadata);
	int openingQuote = -1;
	const char *keys[] = {"n", "r", "p", "keyLen"};
	int keysCount = 4;
	int keyIndex = -1;
	int startNumber = -1;
	for (int i = 0; i < length; i++)
	{
		if (metadata[i] == '\"')
		{
			startNumber = -1;
			if (openingQuote >= 0)
			{
				keyIndex = -1;
				metadata[i] = 0;
				for (int k = 0; k < keysCount; k++)
				{
					if (strcmp(metadata + openingQuote + 1, keys[k]) == 0)
					{
						keyIndex = k;
					}
				}
				openingQuote = -1;
			}
			else
			{
				openingQuote = i;
			}
		}
		else if (metadata[i] >= '0' && metadata[i] <= '9')
		{
			if (startNumber < 0)
				startNumber = i;
		}
		else if (metadata[i] == ',')
		{
			if (startNumber >= 0)
			{
				metadata[i] = 0;
				int number = atoi(metadata + startNumber);
				startNumber = -1;
				if (keyIndex == 0)
					*n = number;
				else if (keyIndex == 1)
					*r = number;
				else if (keyIndex == 2)
					*p = number;
				else if (keyIndex == 3)
					*keyLen = number;
			}
		}
		else
		{
			startNumber = -1;
		}
	}
}

int cutSlice(GoSlice_* slice, int start, int end, int elem_size, GoSlice_* result){
	int size = end - start;
	if( size <= 0)
		return 1;
	void* data = malloc(size * elem_size);
	if( data == NULL )
		return 1;
	registerMemCleanup( data );
	result->data = data;
	result->len = size;
	result->cap = size;
	char* p = slice->data;
	p += (elem_size * start);
	memcpy( data, p, elem_size * size );
	return 0;
}

coin__Transaction* makeEmptyTransaction(Transaction__Handle* handle){
  int result;
  coin__Transaction* ptransaction = NULL;
  result  = SKY_coin_Create_Transaction(handle);
   registerHandleClose(*handle);
  result = SKY_coin_GetTransactionObject( *handle, &ptransaction );
    return ptransaction;
}
int makeUxBodyWithSecret(coin__UxBody* puxBody, cipher__SecKey* pseckey){
  cipher__PubKey pubkey;
  cipher__Address address;
  int result;

  memset( puxBody, 0, sizeof(coin__UxBody) );
  puxBody->Coins = 1000000;
  puxBody->Hours = 100;

  result = SKY_cipher_GenerateKeyPair(&pubkey, pseckey);
  if(result != 0){ return 1;}

  GoSlice slice;
  memset(&slice, 0, sizeof(GoSlice));
  cipher__SHA256 hash;

  result = SKY_cipher_RandByte( 128, (coin__UxArray*)&slice );
  registerMemCleanup( slice.data );
  if(result != 0){ return 1;}
  result = SKY_cipher_SumSHA256( slice, &puxBody->SrcTransaction );
  if(result != 0){ return 1;}

  result = SKY_cipher_AddressFromPubKey( &pubkey, &puxBody->Address );
  if(result != 0){ return 1;}
  return result;
}
int makeUxOutWithSecret(coin__UxOut* puxOut, cipher__SecKey* pseckey){
  int result;
  memset( puxOut, 0, sizeof(coin__UxOut) );
  result = makeUxBodyWithSecret(&puxOut->Body, pseckey);
  puxOut->Head.Time = 100;
  puxOut->Head.BkSeq = 2;
  return result;
}
int makeUxOut(coin__UxOut* puxOut){
  cipher__SecKey seckey;
  return makeUxOutWithSecret(puxOut, &seckey);
}
int makeUxArray(GoSlice* parray, int n){
  parray->data = malloc( sizeof(coin__UxOut) * n );
  if(!parray->data)
    return 1;
  registerMemCleanup( parray->data );
  parray->cap = parray->len = n;
  coin__UxOut* p = (coin__UxOut*)parray->data;
  int result = 0;
  for(int i = 0; i < n; i++){
    result = makeUxOut(p);
    if( result != 0 )
      break;
    p++;
  }
  return result;
}
int makeAddress(cipher__Address* paddress){
  cipher__PubKey pubkey;
  cipher__SecKey seckey;
  cipher__Address address;
  int result;

  result = SKY_cipher_GenerateKeyPair(&pubkey, &seckey);
  if(result != 0) return 1;

  result = SKY_cipher_AddressFromPubKey( &pubkey, paddress );
  if(result != 0) return 1;
  return result;
}
coin__Transaction* makeTransactionFromUxOut(coin__UxOut* puxOut, cipher__SecKey* pseckey, Transaction__Handle* handle ){
  int result;
  coin__Transaction* ptransaction = NULL;
  result  = SKY_coin_Create_Transaction(handle);
//   cr_assert(result == SKY_OK, "SKY_coin_Create_Transaction failed");
  registerHandleClose(*handle);
  result = SKY_coin_GetTransactionObject( *handle, &ptransaction );
//   cr_assert(result == SKY_OK, "SKY_coin_GetTransactionObject failed");
  cipher__SHA256 sha256;
  result = SKY_coin_UxOut_Hash(puxOut, &sha256);
//   cr_assert(result == SKY_OK, "SKY_coin_UxOut_Hash failed");
  GoUint16 r;
  result = SKY_coin_Transaction_PushInput(*handle, &sha256, &r);
//   cr_assert(result == SKY_OK, "SKY_coin_Transaction_PushInput failed");

  cipher__Address address1, address2;
  result = makeAddress(&address1);
//   cr_assert(result == SKY_OK, "makeAddress failed");
  result = makeAddress(&address2);
//   cr_assert(result == SKY_OK, "makeAddress failed");

  result = SKY_coin_Transaction_PushOutput(*handle, &address1, 1000000, 50);
//   cr_assert(result == SKY_OK, "SKY_coin_Transaction_PushOutput failed");
  result = SKY_coin_Transaction_PushOutput(*handle, &address2, 5000000, 50);
//   cr_assert(result == SKY_OK, "SKY_coin_Transaction_PushOutput failed");

  GoSlice secKeys = { pseckey, 1, 1 };
  result = SKY_coin_Transaction_SignInputs( *handle, secKeys );
//   cr_assert(result == SKY_OK, "SKY_coin_Transaction_SignInputs failed");
  result = SKY_coin_Transaction_UpdateHeader( *handle );
//   cr_assert(result == SKY_OK, "SKY_coin_Transaction_UpdateHeader failed");
  return ptransaction;
}

coin__Transaction* makeTransaction(Transaction__Handle* handle){
  int result;
  coin__UxOut uxOut;
  cipher__SecKey seckey;
  coin__Transaction* ptransaction = NULL;
  result = makeUxOutWithSecret( &uxOut, &seckey );
  if(result != 0) return 1;
  ptransaction = makeTransactionFromUxOut( &uxOut, &seckey, handle );
  if(result != 0) return 1;
  return ptransaction;
}
    %}
