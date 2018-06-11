#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"
#include "skycriterion.h"
#include "transutil.h"

TestSuite(coin_transactions, .init = setup, .fini = teardown);

Test(coin_transactions, TestTransactionVerifyInput){
  int result;
  Transaction__Handle handle;
  coin__Transaction* ptx;
  ptx = makeTransaction(&handle);
  result = SKY_coin_Transaction_VerifyInput(handle, NULL);
  cr_assert( result != SKY_OK );
  coin__UxArray ux;
  memset(&ux, 0, sizeof(coin__UxArray));
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert( result != SKY_OK );
  memset(&ux, 0, sizeof(coin__UxArray));
  ux.data = malloc(3 * sizeof(coin__UxOut));
  cr_assert(ux.data != NULL);
  registerMemCleanup(ux.data);
  ux.len = 3;
  ux.cap = 3;
  memset(ux.data, 0, 3 * sizeof(coin__UxOut));
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert( result != SKY_OK );

  coin__UxOut uxOut;
  cipher__SecKey seckey;
  cipher__Sig sig;
  cipher__SHA256 hash;

  result = makeUxOutWithSecret(&uxOut, &seckey);
  cr_assert( result == SKY_OK );
  ptx = makeTransactionFromUxOut(&uxOut, &seckey, &handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_ResetSignatures(handle, 0);
  cr_assert( result == SKY_OK );
  ux.data = &uxOut;
  ux.len = 1; ux.cap = 1;
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert( result != SKY_OK );

  memset(&sig, 0, sizeof(cipher__Sig));
  result = makeUxOutWithSecret(&uxOut, &seckey);
  cr_assert( result == SKY_OK );
  ptx = makeTransactionFromUxOut(&uxOut, &seckey, &handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_ResetSignatures(handle, 1);
  cr_assert( result == SKY_OK );
  memcpy(ptx->Sigs.data, &sig, sizeof(cipher__Sig));
  ux.data = &uxOut;
  ux.len = 1; ux.cap = 1;
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert( result != SKY_OK );

  //Invalid Tx Inner Hash
  result = makeUxOutWithSecret(&uxOut, &seckey);
  cr_assert( result == SKY_OK );
  ptx = makeTransactionFromUxOut(&uxOut, &seckey, &handle);
  cr_assert( result == SKY_OK );
  memset( ptx->InnerHash, 0, sizeof(cipher__SHA256) );
  ux.data = &uxOut;
  ux.len = 1; ux.cap = 1;
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert( result != SKY_OK );

  //Ux hash mismatch
  result = makeUxOutWithSecret(&uxOut, &seckey);
  cr_assert( result == SKY_OK );
  ptx = makeTransactionFromUxOut(&uxOut, &seckey, &handle);
  cr_assert( result == SKY_OK );
  memset( &uxOut, 0, sizeof(coin__UxOut) );
  ux.data = &uxOut;
  ux.len = 1; ux.cap = 1;
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert( result != SKY_OK );

  //Invalid signature
  result = makeUxOutWithSecret(&uxOut, &seckey);
  cr_assert( result == SKY_OK );
  ptx = makeTransactionFromUxOut(&uxOut, &seckey, &handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_ResetSignatures(handle, 1);
  cr_assert( result == SKY_OK );
  memset(ptx->Sigs.data, 0, sizeof(cipher__Sig));
  ux.data = &uxOut;
  ux.len = 1; ux.cap = 1;
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert( result != SKY_OK );

  //Valid
  result = makeUxOutWithSecret(&uxOut, &seckey);
  cr_assert( result == SKY_OK );
  ptx = makeTransactionFromUxOut(&uxOut, &seckey, &handle);
  cr_assert( result == SKY_OK );
  ux.data = &uxOut;
  ux.len = 1; ux.cap = 1;
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert( result == SKY_OK );
}

Test(coin_transactions, TestTransactionSignInputs){
  int result;
  coin__Transaction* ptx;
  Transaction__Handle handle;
  coin__UxOut ux, ux2;
  cipher__SecKey seckey, seckey2;
  cipher__SHA256 hash, hash2;
  cipher__Address addr, addr2;
  cipher__PubKey pubkey;
  GoUint16 r;
  GoSlice keys;

  //Error if txns already signed
  ptx = makeEmptyTransaction(&handle);
  result = SKY_coin_Transaction_ResetSignatures(handle, 1);
  cr_assert( result == SKY_OK );

  memset( &seckey, 0, sizeof(cipher__SecKey) );
  keys.data = &seckey; keys.len = 1; keys.cap = 1;
  result = SKY_coin_Transaction_SignInputs(handle, keys);
  cr_assert( result != SKY_OK );

  // Panics if not enough keys
  ptx = makeEmptyTransaction(&handle);
  memset(&seckey, 0, sizeof(cipher__SecKey));
  memset(&seckey2, 0, sizeof(cipher__SecKey));
  result = makeUxOutWithSecret( &ux, &seckey );
  cr_assert( result == SKY_OK );
  result = SKY_coin_UxOut_Hash(&ux, &hash);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_PushInput(handle, &hash, &r);
  cr_assert( result == SKY_OK );
  result = makeUxOutWithSecret( &ux2, &seckey2 );
  cr_assert( result == SKY_OK );
  result = SKY_coin_UxOut_Hash(&ux2, &hash2);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_PushInput(handle, &hash2, &r);
  cr_assert( result == SKY_OK );
  makeAddress(&addr);
  result = SKY_coin_Transaction_PushOutput(handle, &addr, 40, 80);
  cr_assert( result == SKY_OK );
  cr_assert( ptx->Sigs.len == 0 );
  keys.data = &seckey; keys.len = 1; keys.cap = 1;
  result = SKY_coin_Transaction_SignInputs(handle, keys);
  cr_assert( result != SKY_OK );
  cr_assert( ptx->Sigs.len == 0 );

  // Valid signing
  result = SKY_coin_Transaction_HashInner( handle, &hash );
  cr_assert( result == SKY_OK );
  keys.data = malloc(2 * sizeof(cipher__SecKey));
  cr_assert( keys.data != NULL );
  registerMemCleanup( keys.data );
  keys.len = keys.cap = 2;
  memcpy(keys.data, &seckey, sizeof(cipher__SecKey));
  memcpy(((cipher__SecKey*)keys.data) + 1, &seckey2, sizeof(cipher__SecKey));
  result = SKY_coin_Transaction_SignInputs(handle, keys);
  cr_assert( result == SKY_OK );
  cr_assert(ptx->Sigs.len == 2);
  result = SKY_coin_Transaction_HashInner( handle, &hash2 );
  cr_assert( result == SKY_OK );
  cr_assert( eq( u8[sizeof(cipher__SHA256)], hash, hash2) );

  result = SKY_cipher_PubKeyFromSecKey(&seckey, &pubkey);
  cr_assert( result == SKY_OK );
  result = SKY_cipher_AddressFromPubKey( &pubkey, &addr );
  cr_assert( result == SKY_OK );
  result = SKY_cipher_PubKeyFromSecKey(&seckey2, &pubkey);
  cr_assert( result == SKY_OK );
  result = SKY_cipher_AddressFromPubKey( &pubkey, &addr2 );
  cr_assert( result == SKY_OK );

  cipher__SHA256 addHash, addHash2;
  result = SKY_cipher_AddSHA256(&hash, (cipher__SHA256*)ptx->In.data, &addHash);
  cr_assert( result == SKY_OK );
  result = SKY_cipher_AddSHA256(&hash, ((cipher__SHA256*)ptx->In.data) + 1, &addHash2);
  cr_assert( result == SKY_OK );
  result = SKY_cipher_ChkSig(&addr, &addHash, (cipher__Sig*)ptx->Sigs.data);
  cr_assert( result == SKY_OK );
  result = SKY_cipher_ChkSig(&addr2, &addHash2, ((cipher__Sig*)ptx->Sigs.data)+1);
  cr_assert( result == SKY_OK );
  result = SKY_cipher_ChkSig(&addr, &hash, ((cipher__Sig*)ptx->Sigs.data)+1);
  cr_assert( result != SKY_OK );
  result = SKY_cipher_ChkSig(&addr2, &hash, (cipher__Sig*)ptx->Sigs.data);
  cr_assert( result != SKY_OK );
}

Test(coin_transactions, TestTransactionHashInner){
  int result;
  Transaction__Handle handle1 = 0, handle2 = 0;
  coin__Transaction* ptx = NULL;
  coin__Transaction* ptx2 = NULL;
  ptx = makeTransaction(&handle1);
  cipher__SHA256 hash, nullHash;
  result = SKY_coin_Transaction_HashInner( handle1, &hash );
  cr_assert( result == SKY_OK );
  memset( &nullHash, 0, sizeof(cipher__SHA256) );
  cr_assert( not ( eq( u8[sizeof(cipher__SHA256)], nullHash, hash) ) );

  // If tx.In is changed, hash should change
  ptx2 = copyTransaction( handle1, &handle2 );
  cr_assert( eq(  type(coin__Transaction), *ptx, *ptx2 ) );
  cr_assert( ptx != ptx2 );
  cr_assert(ptx2->In.len > 0);
  coin__UxOut uxOut;
  makeUxOut( &uxOut );
  cipher__SHA256* phash = ptx2->In.data;
  result = SKY_coin_UxOut_Hash( &uxOut, phash );
  cr_assert( result == SKY_OK );
  cr_assert( not( eq(  type(coin__Transaction), *ptx, *ptx2 ) ) );
  cipher__SHA256 hash1, hash2;
  result = SKY_coin_Transaction_HashInner( handle1, &hash1 );
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_HashInner( handle2, &hash2 );
  cr_assert( result == SKY_OK );
  cr_assert( not ( eq( u8[sizeof(cipher__SHA256)], hash1, hash2) ) );

  // If tx.Out is changed, hash should change
  handle2 = 0;
  ptx2 = copyTransaction( handle1, &handle2 );
  cr_assert( ptx != ptx2 );
  cr_assert( eq(  type(coin__Transaction), *ptx, *ptx2 ) );
  coin__TransactionOutput* output = ptx2->Out.data;
  cipher__Address addr;
  makeAddress( &addr );
  memcpy( &output->Address, &addr, sizeof(cipher__Address) );
  cr_assert( not( eq(  type(coin__Transaction), *ptx, *ptx2 ) ) );
  cr_assert(eq(type(cipher__Address), addr, output->Address));
  result = SKY_coin_Transaction_HashInner( handle1, &hash1 );
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_HashInner( handle2, &hash2 );
  cr_assert( result == SKY_OK );
  cr_assert( not ( eq( u8[sizeof(cipher__SHA256)], hash1, hash2) ) );

  // If tx.Head is changed, hash should not change
  ptx2 = copyTransaction( handle1, &handle2 );
  int len = ptx2->Sigs.len;
  cipher__Sig* newSigs = malloc((len + 1) * sizeof(cipher__Sig));
  cr_assert( newSigs != NULL );
  registerMemCleanup( newSigs );
  memcpy( newSigs, ptx2->Sigs.data, len * sizeof(cipher__Sig));
  result = SKY_coin_Transaction_ResetSignatures(handle2, len + 1);
  cr_assert( result == SKY_OK );
  memcpy( ptx2->Sigs.data, newSigs, len * sizeof(cipher__Sig));
  newSigs += len;
  memset( newSigs, 0, sizeof(cipher__Sig) );
  result = SKY_coin_Transaction_HashInner( handle1, &hash1 );
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_HashInner( handle2, &hash2 );
  cr_assert( result == SKY_OK );
  cr_assert( eq( u8[sizeof(cipher__SHA256)], hash1, hash2) );
}

Test(coin_transactions, TestTransactionSerialization){
  int result;
  coin__Transaction* ptx;
  Transaction__Handle handle;
  ptx = makeTransaction(&handle);
  GoSlice_ data;
  memset( &data, 0, sizeof(GoSlice_) );
  result = SKY_coin_Transaction_Serialize( handle, &data );
  cr_assert( result == SKY_OK );
  registerMemCleanup( data.data );
  coin__Transaction* ptx2;
  Transaction__Handle handle2;
  GoSlice d = {data.data, data.len, data.cap};
  result = SKY_coin_TransactionDeserialize(d, &handle2);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Get_Transaction_Object(handle2, &ptx2);
  cr_assert( result == SKY_OK );
  cr_assert( eq( type(coin__Transaction), *ptx, *ptx2) );
}

Test(coin_transactions, TestTransactionOutputHours){
  coin__Transaction* ptx;
  Transaction__Handle handle;
  ptx = makeEmptyTransaction(&handle);
  cipher__Address addr;
  makeAddress(&addr);
  int result;
  result = SKY_coin_Transaction_PushOutput(handle, &addr, 1000000, 100);
  cr_assert( result == SKY_OK );
  makeAddress(&addr);
  result = SKY_coin_Transaction_PushOutput(handle, &addr, 1000000, 200);
  cr_assert( result == SKY_OK );
  makeAddress(&addr);
  result = SKY_coin_Transaction_PushOutput(handle, &addr, 1000000, 500);
  cr_assert( result == SKY_OK );
  makeAddress(&addr);
  result = SKY_coin_Transaction_PushOutput(handle, &addr, 1000000, 0);
  cr_assert( result == SKY_OK );
  GoUint64 hours;
  result = SKY_coin_Transaction_OutputHours(handle, &hours);
  cr_assert( result == SKY_OK );
  cr_assert( hours == 800 );
  makeAddress(&addr);
  result = SKY_coin_Transaction_PushOutput(handle, &addr, 1000000, 0xFFFFFFFFFFFFFFFF - 700);
  result = SKY_coin_Transaction_OutputHours(handle, &hours);
  cr_assert( result != SKY_OK );
}

Test(coin_transactions, TestTransactionsHashes){
  int result;
  GoSlice_ hashes = {NULL, 0, 0};
  Transactions__Handle hTxns;
  result = makeTransactions(4, &hTxns);
  cr_assert( result == SKY_OK );

  result = SKY_coin_Transactions_Hashes(hTxns, &hashes);
  cr_assert( result == SKY_OK, "SKY_coin_Transactions_Hashes failed" );
  registerMemCleanup( hashes.data );
  cr_assert( hashes.len == 4 );
  cipher__SHA256* ph = hashes.data;
  cipher__SHA256 hash;
  for(int i = 0; i < 4; i++){
    Transaction__Handle handle;
    result = SKY_coin_Transactions_GetAt(hTxns, i, &handle);
    cr_assert( result == SKY_OK );
    result = SKY_coin_Transaction_Hash(handle, &hash);
    cr_assert( result == SKY_OK, "SKY_coin_Transaction_Hash failed" );
    cr_assert( eq( u8[sizeof(cipher__SHA256)], *ph, hash) );
    ph++;
  }
}


Test(coin_transactions, TestTransactionsTruncateBytesTo){
  int result;
  Transactions__Handle h1, h2;
  result = makeTransactions(10, &h1);
  cr_assert( result == SKY_OK );
  GoInt length;
  result = SKY_coin_Transactions_Length(h1, &length);
  cr_assert( result == SKY_OK );
  int trunc = 0;
  GoInt size;
  for(int i = 0; i < length / 2; i++){
    Transaction__Handle handle;
    result = SKY_coin_Transactions_GetAt(h1, i, &handle);
    registerHandleClose(handle);
    cr_assert( result == SKY_OK );
    result = SKY_coin_Transaction_Size(handle, &size);
    trunc += size;
    cr_assert( result == SKY_OK, "SKY_coin_Transaction_Size failed" );
  }
  result = SKY_coin_Transactions_TruncateBytesTo(h1, trunc, &h2);
  cr_assert( result == SKY_OK, "SKY_coin_Transactions_TruncateBytesTo failed" );
  registerHandleClose(h2);

  GoInt length2;
  result = SKY_coin_Transactions_Length(h2, &length2);
  cr_assert( result == SKY_OK );
  cr_assert( length2 == length / 2 );
  result = SKY_coin_Transactions_Size( h2, &size );
  cr_assert( result == SKY_OK, "SKY_coin_Transactions_Size failed" );
  cr_assert( trunc == size );

  trunc++;
  result = SKY_coin_Transactions_TruncateBytesTo(h1, trunc, &h2);
  cr_assert( result == SKY_OK, "SKY_coin_Transactions_TruncateBytesTo failed" );
  registerHandleClose(h2);

  // Stepping into next boundary has same cutoff, must exceed
  result = SKY_coin_Transactions_Length(h2, &length2);
  cr_assert( result == SKY_OK );
  cr_assert( length2 == length / 2 );
  result = SKY_coin_Transactions_Size( h2, &size );
  cr_assert( result == SKY_OK, "SKY_coin_Transactions_Size failed" );
  cr_assert( trunc - 1 == size );
}

typedef struct {
  GoUint64 coins;
  GoUint64 hours;
} test_ux;

typedef struct {
  test_ux* inUxs;
  test_ux* outUxs;
  int sizeIn;
  int sizeOut;
  GoUint64 headTime;
  int      failure;
} test_case;

int makeTestCaseArrays(test_ux* elems, int size, coin__UxArray* pArray){
  if(size <= 0){
    pArray->len = 0;
    pArray->cap = 0;
    pArray->data = NULL;
    return SKY_OK;
  }
  int elems_size = sizeof(coin__UxOut);
  void* data;
  data = malloc(size * elems_size);
  if( data == NULL)
    return SKY_ERROR;
  registerMemCleanup( data );
  pArray->data = data;
  pArray->len = size;
  pArray->cap = size;
  coin__UxOut* p = data;
  for(int i = 0; i < size; i++){
    p->Body.Coins = elems[i].coins;
    p->Body.Hours = elems[i].hours;
    p++;
  }
  return SKY_OK;
}

Test(coin_transactions, TestVerifyTransactionCoinsSpending){
  GoUint64 MaxUint64 = 0xFFFFFFFFFFFFFFFF;
  GoUint64 Million = 1000000;

  //Input coins overflow
  test_ux in1[] = {
      {MaxUint64 - Million + 1, 10},
      {Million, 0}
  };

  //Output coins overflow
  test_ux in2[] = {
      {10 * Million, 10}
  };
  test_ux out2[] = {
      {MaxUint64 - 10 * Million + 1, 0},
      {20 * Million, 1}
  };

  //Insufficient coins
  test_ux in3[] = {
    {10 * Million, 10},
    {15 * Million, 10}
  };
  test_ux out3[] = {
    {20 * Million, 1},
    {10 * Million, 1}
  };

  //Destroyed coins
  test_ux in4[] = {
    {10 * Million, 10},
    {15 * Million, 10}
  };
  test_ux out4[] = {
    {5 * Million, 1},
    {10 * Million, 1}
  };

  //Valid
  test_ux in5[] = {
    {10 * Million, 10},
    {15 * Million, 10}
  };
  test_ux out5[] = {
    {10 * Million, 11},
    {10 * Million, 1},
    {5 * Million, 0}
  };

  test_case tests[] = {
      {in1, NULL, 2, 0, 0, 1},    //Input coins overflow
      {in2, out2, 1, 2, 0, 1},    //Output coins overflow
      {in3, out3, 2, 2, 0, 1},    //Destroyed coins
      {in4, out4, 1, 1, Million, 1},    //Invalid (coin hours overflow when adding earned hours, which is treated as 0, and now enough coin hours)
      {in5, out5, 2, 3, 0, 0}   //Valid
  };

  coin__UxArray inArray;
  coin__UxArray outArray;
  int result;
  int count = sizeof(tests) / sizeof(tests[0]);
  for( int i = 0; i < count; i++){
    result = makeTestCaseArrays(tests[i].inUxs, tests[i].sizeIn, &inArray);
    cr_assert(result == SKY_OK);
    result = makeTestCaseArrays(tests[i].outUxs, tests[i].sizeOut, &outArray);
    cr_assert(result == SKY_OK);
    result = SKY_coin_VerifyTransactionCoinsSpending(&inArray, &outArray);
    if( tests[i].failure )
      cr_assert( result != SKY_OK, "VerifyTransactionCoinsSpending succeeded %d", i+1 );
    else
      cr_assert( result == SKY_OK, "VerifyTransactionCoinsSpending failed %d", i+1 );
  }
}

Test(coin_transactions, TestVerifyTransactionHoursSpending){
  GoUint64 MaxUint64 = 0xFFFFFFFFFFFFFFFF;
  GoUint64 Million = 1000000;

  //Input hours overflow
  test_ux in1[] = {
      {3 * Million, MaxUint64 - Million + 1},
      {Million, Million}
  };

  //Insufficient coin hours
  test_ux in2[] = {
      {10 * Million, 10},
      {15 * Million, 10}
  };


  test_ux out2[] = {
      {15 * Million, 10},
      {10 * Million, 11}
  };

  //coin hours time calculation overflow
  test_ux in3[] = {
      {10 * Million, 10},
      {15 * Million, 10}
  };


  test_ux out3[] = {
      {10 * Million, 11},
      {10 * Million, 1},
      {5 * Million, 0}
  };

  //Invalid (coin hours overflow when adding earned hours, which is treated as 0, and now enough coin hours)
  test_ux in4[] = {
      {10 * Million, MaxUint64}
  };

  test_ux out4[] = {
      {10 * Million, 1}
  };

  //Valid (coin hours overflow when adding earned hours, which is treated as 0, but not sending any hours)
  test_ux in5[] = {
      {10 * Million, MaxUint64}
  };

  test_ux out5[] = {
      {10 * Million, 0}
  };

  //Valid (base inputs have insufficient coin hours, but have sufficient after adjusting coinhours by headTime)
  test_ux in6[] = {
      {10 * Million, 10},
      {15 * Million, 10}
  };

  test_ux out6[] = {
      {15 * Million, 10},
      {10 * Million, 11}
  };

  //valid
  test_ux in7[] = {
      {10 * Million, 10},
      {15 * Million, 10}
  };

  test_ux out7[] = {
      {10 * Million, 11},
      {10 * Million, 1},
      {5 * Million, 0}
  };

  test_case tests[] = {
      {in1, NULL, 2, 0, 0, 1},    //Input hours overflow
      {in2, out2, 2, 2, 0, 1},    //Insufficient coin hours
      {in3, out3, 2, 3, MaxUint64, 1}, //coin hours time calculation overflow
      {in4, out4, 1, 1, Million, 1},    //Invalid (coin hours overflow when adding earned hours, which is treated as 0, and now enough coin hours)
      {in5, out5, 1, 1, 0, 0},    //Valid (coin hours overflow when adding earned hours, which is treated as 0, but not sending any hours)
      {in6, out6, 2, 2, 1492707255, 0},    //Valid (base inputs have insufficient coin hours, but have sufficient after adjusting coinhours by headTime)
      {in7, out7, 2, 3, 0, 0},  //Valid
  };
  coin__UxArray inArray;
  coin__UxArray outArray;
  int result;
  int count = sizeof(tests) / sizeof(tests[0]);
  for( int i = 0; i < count; i++){
    result = makeTestCaseArrays(tests[i].inUxs, tests[i].sizeIn, &inArray);
    cr_assert(result == SKY_OK);
    result = makeTestCaseArrays(tests[i].outUxs, tests[i].sizeOut, &outArray);
    cr_assert(result == SKY_OK);
    result = SKY_coin_VerifyTransactionHoursSpending(tests[i].headTime, &inArray, &outArray);
    if( tests[i].failure )
      cr_assert( result != SKY_OK, "SKY_coin_VerifyTransactionHoursSpending succeeded %d", i+1 );
    else
      cr_assert( result == SKY_OK, "SKY_coin_VerifyTransactionHoursSpending failed %d", i+1 );
  }
}


/*******************************************************
*  Tests not done because of wrapper functions were not created:
*  TestTransactionsFees
*  TestSortTransactions
********************************************************/
