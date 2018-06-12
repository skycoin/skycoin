
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

Test(coin_transaction, TestTransactionVerify)
{
  int result;
  coin__Transaction* ptx;
  Transaction__Handle handle;
  // Mismatch header hash
  ptx = makeTransaction(&handle);
  memset(ptx->InnerHash, 0, sizeof(cipher__SHA256));
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert( result != SKY_OK );

  // No inputs
  ptx = makeTransaction(&handle);
  result = SKY_coin_Transaction_ResetInputs(handle, 0);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert( result != SKY_OK );

  // No outputs
  ptx = makeTransaction(&handle);
  result = SKY_coin_Transaction_ResetOutputs(handle, 0);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert( result != SKY_OK );

  //Invalid number of Sigs
  ptx = makeTransaction(&handle);
  result = SKY_coin_Transaction_ResetSignatures(handle, 0);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert( result != SKY_OK );
  result = SKY_coin_Transaction_ResetSignatures(handle, 20);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert( result != SKY_OK );

  int MaxUint16 = 0xFFFF;
  // Too many sigs & inputs
  ptx = makeTransaction(&handle);
  result = SKY_coin_Transaction_ResetSignatures(handle, MaxUint16);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_ResetInputs(handle, MaxUint16);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert( result != SKY_OK );


  // Duplicate inputs
  coin__UxOut ux;
  cipher__SecKey seckey;
  cipher__SHA256 sha256;
  makeUxOutWithSecret( &ux, &seckey );
  ptx = makeTransactionFromUxOut( &ux, &seckey, &handle);
  memcpy(&sha256, ptx->In.data, sizeof(cipher__SHA256));
  GoUint16 r;
  result = SKY_coin_Transaction_PushInput(handle, &sha256, &r);
  result = SKY_coin_Transaction_ResetSignatures(handle, 0);
  cr_assert( result == SKY_OK );
  GoSlice seckeys;
  seckeys.data = malloc(sizeof(cipher__SecKey) * 2);
  cr_assert( seckeys.data != NULL );
  registerMemCleanup( seckeys.data );
  seckeys.len = seckeys.cap = 2;
  memcpy( seckeys.data, &seckey, sizeof(cipher__SecKey) );
  memcpy( ((cipher__SecKey*)seckeys.data) + 1, &seckey, sizeof(cipher__SecKey) );
  result = SKY_coin_Transaction_SignInputs( handle, seckeys );
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert( result != SKY_OK );

  //Duplicate outputs
  ptx = makeTransaction(&handle);
  coin__TransactionOutput* pOutput = ptx->Out.data;
  cipher__Address addr;
  memcpy(&addr, &pOutput->Address, sizeof(cipher__Address));
  result = SKY_coin_Transaction_PushOutput(handle, &addr, pOutput->Coins, pOutput->Hours);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert( result != SKY_OK );

  // Invalid signature, empty
  ptx = makeTransaction(&handle);
  memset(ptx->Sigs.data, 0, sizeof(cipher__Sig));
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert( result != SKY_OK );

  // Output coins are 0
  ptx = makeTransaction(&handle);
  pOutput = ptx->Out.data;
  pOutput->Coins = 0;
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert( result != SKY_OK );

  GoUint64 MaxUint64 = 0xFFFFFFFFFFFFFFFF;
  // Output coin overflow
  ptx = makeTransaction(&handle);
  pOutput = ptx->Out.data;
  pOutput->Coins = MaxUint64 - 3000000;
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert( result != SKY_OK );

  // Output coins are not multiples of 1e6 (valid, decimal restriction is not enforced here)
  ptx = makeTransaction(&handle);
  pOutput = ptx->Out.data;
  pOutput->Coins += 10;
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_ResetSignatures(handle, 0);
  cr_assert( result == SKY_OK );
  cipher__PubKey pubkey;
  result = SKY_cipher_GenerateKeyPair(&pubkey, &seckey);
  cr_assert( result == SKY_OK );
  seckeys.data = &seckey; seckeys.len = 1; seckeys.cap = 1;
  result = SKY_coin_Transaction_SignInputs(handle, seckeys);
  cr_assert( result == SKY_OK );
  cr_assert( pOutput->Coins % 1000000 != 0 );
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert( result == SKY_OK );

  //Valid
  ptx = makeTransaction(&handle);
  pOutput = ptx->Out.data;
  pOutput->Coins = 10000000;
  pOutput++;
  pOutput->Coins = 1000000;
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert( result == SKY_OK );
}


Test(coin_transaction, TestTransactionPushInput)
{
  int result;
  Transaction__Handle handle;
  coin__Transaction* ptx;
  coin__UxOut ux;
  ptx = makeEmptyTransaction(&handle);
  makeUxOut( &ux );
  cipher__SHA256 hash;
  result = SKY_coin_UxOut_Hash( &ux, &hash );
  cr_assert( result == SKY_OK );
  GoUint16 r;
  result = SKY_coin_Transaction_PushInput(handle, &hash, &r);
  cr_assert( result == SKY_OK );
  cr_assert( r == 0 );
  cr_assert( ptx->In.len == 1 );
  cipher__SHA256* pIn = ptx->In.data;
  cr_assert( eq( u8[sizeof(cipher__SHA256)], hash, *pIn) );

  GoUint16 MaxUint16 = 0xFFFF;
  int len = ptx->In.len;
  void* data = malloc(len * sizeof(cipher__SHA256));
  cr_assert(data != NULL);
  registerMemCleanup(data);
  memcpy(data, ptx->In.data, len * sizeof(cipher__SHA256) );
  result = SKY_coin_Transaction_ResetInputs(handle, MaxUint16 + len);
  cr_assert( result == SKY_OK );
  memcpy(ptx->In.data, data, len * sizeof(cipher__Sig));
  freeRegisteredMemCleanup(data);
  makeUxOut( &ux );
  result = SKY_coin_UxOut_Hash( &ux, &hash );
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_PushInput(handle, &hash, &r);
  cr_assert( result != SKY_OK );

}


Test(coin_transaction, TestTransactionPushOutput)
{
  int result;
  Transaction__Handle handle;
  coin__Transaction* ptx;
  ptx = makeEmptyTransaction(&handle);

  cipher__Address addr;
  makeAddress( &addr );
  result = SKY_coin_Transaction_PushOutput( handle, &addr, 100, 150 );
  cr_assert( result == SKY_OK );
  cr_assert( ptx->Out.len == 1 );
  coin__TransactionOutput* pOutput = ptx->Out.data;
  coin__TransactionOutput output;
  memcpy(&output.Address, &addr, sizeof(cipher__Address));
  output.Coins = 100;
  output.Hours = 150;
  cr_assert( eq( type(coin__TransactionOutput), output, *pOutput ) );
  for(int i = 1; i < 20; i++){
    makeAddress( &addr );
    result = SKY_coin_Transaction_PushOutput( handle, &addr, i * 100, i * 50 );
    cr_assert( result == SKY_OK );
    cr_assert( ptx->Out.len == i + 1 );
    pOutput = ptx->Out.data;
    pOutput += i;
    memcpy(&output.Address, &addr, sizeof(cipher__Address));
    output.Coins = i * 100;
    output.Hours = i * 50;
    cr_assert( eq( type(coin__TransactionOutput), output, *pOutput ) );
  }
}

Test(coin_transaction, TestTransactionHash)
{
  int result;
  Transaction__Handle handle;
  coin__Transaction* ptx;
  ptx = makeEmptyTransaction(&handle);

  cipher__SHA256 nullHash, hash1, hash2;
  memset( &nullHash, 0, sizeof(cipher__SHA256) );
  result = SKY_coin_Transaction_Hash( handle, &hash1 );
  cr_assert( result == SKY_OK );
  cr_assert( not ( eq( u8[sizeof(cipher__SHA256)], nullHash, hash1) ) );
  result = SKY_coin_Transaction_HashInner( handle, &hash2 );
  cr_assert( result == SKY_OK );
  cr_assert( not ( eq( u8[sizeof(cipher__SHA256)], hash2, hash1) ) );
}


Test(coin_transaction, TestTransactionUpdateHeader)
{
  int result;
  Transaction__Handle handle;
  coin__Transaction* ptx;
  ptx = makeTransaction(&handle);
  cipher__SHA256 hash, nullHash, hashInner;
  memcpy(&hash, &ptx->InnerHash, sizeof(cipher__SHA256));
  memset(&ptx->InnerHash, 0, sizeof(cipher__SHA256));
  memset(&nullHash, 0, sizeof(cipher__SHA256));
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert( not ( eq( u8[sizeof(cipher__SHA256)], ptx->InnerHash, nullHash) ) );
  cr_assert( eq( u8[sizeof(cipher__SHA256)], hash, ptx->InnerHash) );
  result = SKY_coin_Transaction_HashInner( handle, &hashInner );
  cr_assert( result == SKY_OK );
  cr_assert( eq( u8[sizeof(cipher__SHA256)], hashInner, ptx->InnerHash) );
}
