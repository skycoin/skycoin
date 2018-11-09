
#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"
#include "skytxn.h"

GoUint32_ zeroFeeCalculator(Transaction__Handle handle, GoUint64_ *pFee, void* context){
  *pFee = 0;
  return SKY_OK;
}

int makeKeysAndAddress(cipher__PubKey* ppubkey, cipher__SecKey* pseckey, cipher__Address* paddress){
  int result;
  result = SKY_cipher_GenerateKeyPair(ppubkey, pseckey);
  cr_assert(result == SKY_OK, "SKY_cipher_GenerateKeyPair failed");
  result = SKY_cipher_AddressFromPubKey( ppubkey, paddress );
  cr_assert(result == SKY_OK, "SKY_cipher_AddressFromPubKey failed");
  return result;
}

int makeUxBodyWithSecret(coin__UxBody* puxBody, cipher__SecKey* pseckey){
  cipher__PubKey pubkey;
  cipher__Address address;
  int result;

  memset( puxBody, 0, sizeof(coin__UxBody) );
  puxBody->Coins = 1000000;
  puxBody->Hours = 100;

  result = SKY_cipher_GenerateKeyPair(&pubkey, pseckey);
  cr_assert(result == SKY_OK, "SKY_cipher_GenerateKeyPair failed");

  GoSlice slice;
  memset(&slice, 0, sizeof(GoSlice));
  cipher__SHA256 hash;

  result = SKY_cipher_RandByte( 128, (coin__UxArray*)&slice );
  registerMemCleanup( slice.data );
  cr_assert(result == SKY_OK, "SKY_cipher_RandByte failed");
  result = SKY_cipher_SumSHA256( slice, &puxBody->SrcTransaction );
  cr_assert(result == SKY_OK, "SKY_cipher_SumSHA256 failed");

  result = SKY_cipher_AddressFromPubKey( &pubkey, &puxBody->Address );
  cr_assert(result == SKY_OK, "SKY_cipher_AddressFromPubKey failed");
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

int makeUxBody(coin__UxBody* puxBody){
  cipher__SecKey seckey;
  return makeUxBodyWithSecret(puxBody, &seckey);
}

int makeUxOut(coin__UxOut* puxOut){
  cipher__SecKey seckey;
  return makeUxOutWithSecret(puxOut, &seckey);
}

int makeAddress(cipher__Address* paddress){
  cipher__PubKey pubkey;
  cipher__SecKey seckey;
  cipher__Address address;
  int result;

  result = SKY_cipher_GenerateKeyPair(&pubkey, &seckey);
  cr_assert(result == SKY_OK, "SKY_cipher_GenerateKeyPair failed");

  result = SKY_cipher_AddressFromPubKey( &pubkey, paddress );
  cr_assert(result == SKY_OK, "SKY_cipher_AddressFromPubKey failed");
  return result;
}

coin__Transaction* makeTransactionFromUxOut(coin__UxOut* puxOut, cipher__SecKey* pseckey, Transaction__Handle* handle ){
  int result;
  coin__Transaction* ptransaction = NULL;
  result  = SKY_coin_Create_Transaction(handle);
  cr_assert(result == SKY_OK, "SKY_coin_Create_Transaction failed");
  registerHandleClose(*handle);
  result = SKY_coin_GetTransactionObject( *handle, &ptransaction );
  cr_assert(result == SKY_OK, "SKY_coin_GetTransactionObject failed");
  cipher__SHA256 sha256;
  result = SKY_coin_UxOut_Hash(puxOut, &sha256);
  cr_assert(result == SKY_OK, "SKY_coin_UxOut_Hash failed");
  GoUint16 r;
  result = SKY_coin_Transaction_PushInput(*handle, &sha256, &r);
  cr_assert(result == SKY_OK, "SKY_coin_Transaction_PushInput failed");

  cipher__Address address1, address2;
  result = makeAddress(&address1);
  cr_assert(result == SKY_OK, "makeAddress failed");
  result = makeAddress(&address2);
  cr_assert(result == SKY_OK, "makeAddress failed");

  result = SKY_coin_Transaction_PushOutput(*handle, &address1, 1000000, 50);
  cr_assert(result == SKY_OK, "SKY_coin_Transaction_PushOutput failed");
  result = SKY_coin_Transaction_PushOutput(*handle, &address2, 5000000, 50);
  cr_assert(result == SKY_OK, "SKY_coin_Transaction_PushOutput failed");

  GoSlice secKeys = { pseckey, 1, 1 };
  result = SKY_coin_Transaction_SignInputs( *handle, secKeys );
  cr_assert(result == SKY_OK, "SKY_coin_Transaction_SignInputs failed");
  result = SKY_coin_Transaction_UpdateHeader( *handle );
  cr_assert(result == SKY_OK, "SKY_coin_Transaction_UpdateHeader failed");
  return ptransaction;
}

coin__Transaction* makeTransaction(Transaction__Handle* handle){
  int result;
  coin__UxOut uxOut;
  cipher__SecKey seckey;

  coin__Transaction* ptransaction = NULL;

  result = makeUxOutWithSecret( &uxOut, &seckey );
  cr_assert(result == SKY_OK, "makeUxOutWithSecret failed");
  return makeTransactionFromUxOut( &uxOut, &seckey, handle );
}

coin__Transaction* makeEmptyTransaction(Transaction__Handle* handle){
  int result;
  coin__Transaction* ptransaction = NULL;
  result  = SKY_coin_Create_Transaction(handle);
  cr_assert(result == SKY_OK, "SKY_coin_Create_Transaction failed");
  registerHandleClose(*handle);
  result = SKY_coin_GetTransactionObject( *handle, &ptransaction );
  cr_assert(result == SKY_OK, "SKY_coin_GetTransactionObject failed");
  return ptransaction;
}


int makeTransactions(int n, Transactions__Handle* handle){
  int result = SKY_coin_Create_Transactions(handle);
  cr_assert(result == SKY_OK);
  registerHandleClose(*handle);
  for(int i = 0; i < n; i++){
    Transaction__Handle thandle;
    makeTransaction(&thandle);
    registerHandleClose(thandle);
    result = SKY_coin_Transactions_Add(*handle, thandle);
    cr_assert(result == SKY_OK);
  }
  return result;
}

typedef struct{
  cipher__SHA256 hash;
  Transaction__Handle handle;
} TransactionObjectHandle;

int sortTransactions(Transactions__Handle txns_handle, Transactions__Handle* sorted_txns_handle){
  int result = SKY_coin_Create_Transactions(sorted_txns_handle);
  cr_assert(result == SKY_OK);
  registerHandleClose(*sorted_txns_handle);
  GoInt n, i, j;
  result = SKY_coin_Transactions_Length(txns_handle, &n);
  cr_assert(result == SKY_OK);
  TransactionObjectHandle* pTrans = malloc( n * sizeof(TransactionObjectHandle));
  cr_assert(pTrans != NULL);
  registerMemCleanup(pTrans);
  memset(pTrans, 0, n * sizeof(TransactionObjectHandle));
  int* indexes = malloc( n * sizeof(int) );
  cr_assert(indexes != NULL);
  registerMemCleanup(indexes);
  for( i = 0; i < n; i ++){
    indexes[i] = i;
    result = SKY_coin_Transactions_GetAt(txns_handle, i, &pTrans[i].handle);
    cr_assert(result == SKY_OK);
    registerHandleClose(pTrans[i].handle);
    result = SKY_coin_Transaction_Hash(pTrans[i].handle, &pTrans[i].hash);
    cr_assert(result == SKY_OK);
  }

  //Swap sort.
  cipher__SHA256 hash1, hash2;
  for(i = 0; i < n - 1; i++){
    for(j = i + 1; j < n; j++){
      int cmp = memcmp(&pTrans[indexes[i]].hash, &pTrans[indexes[j]].hash, sizeof(cipher__SHA256));
      if(cmp > 0){
        //Swap
        int tmp = indexes[i];
        indexes[i] = indexes[j];
        indexes[j] = tmp;
      }
    }
  }
  for( i = 0; i < n; i ++){
    result = SKY_coin_Transactions_Add(*sorted_txns_handle, pTrans[indexes[i]].handle);
    cr_assert(result == SKY_OK);
  }
  return result;
}

coin__Transaction* copyTransaction(Transaction__Handle handle, Transaction__Handle* handle2){
  coin__Transaction* ptransaction = NULL;
  int result = 0;
  result = SKY_coin_Transaction_Copy(handle, handle2);
  cr_assert(result == SKY_OK);
  registerHandleClose(*handle2);
  result = SKY_coin_GetTransactionObject( *handle2, &ptransaction );
  cr_assert(result == SKY_OK, "SKY_coin_GetTransactionObject failed");
  return ptransaction;
}

void makeRandHash(cipher__SHA256* phash){
  GoSlice slice;
  memset(&slice, 0, sizeof(GoSlice));

  int result = SKY_cipher_RandByte( 128, (coin__UxArray*)&slice );
  cr_assert(result == SKY_OK, "SKY_cipher_RandByte failed");
  registerMemCleanup( slice.data );
  result = SKY_cipher_SumSHA256( slice, phash );
  cr_assert(result == SKY_OK, "SKY_cipher_SumSHA256 failed");
}

int makeUxArray(coin__UxArray* parray, int n){
  parray->data = malloc( sizeof(coin__UxOut) * n );
  if(!parray->data)
    return SKY_ERROR;
  registerMemCleanup( parray->data );
  parray->cap = parray->len = n;
  coin__UxOut* p = (coin__UxOut*)parray->data;
  int result = SKY_OK;
  for(int i = 0; i < n; i++){
    result = makeUxOut(p);
    if( result != SKY_OK )
      break;
    p++;
  }
  return result;
}
