
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
  coin__Transaction tx;

  // Mismatch header hash
  makeTransaction(&tx);
  memset(&tx.InnerHash, 0, sizeof(cipher__SHA256));
  result = SKY_coin_Transaction_Verify(&tx);
  cr_assert( result != SKY_OK );

  // No inputs
  makeTransaction(&tx);
  tx.In.len = 0;
  result = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(&tx);
  cr_assert( result != SKY_OK );

  // No outputs
  makeTransaction(&tx);
  tx.Out.len = 0;
  result = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(&tx);
  cr_assert( result != SKY_OK );

  //Invalid number of Sigs
  makeTransaction(&tx);
  tx.Sigs.len = 0;
  result = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(&tx);
  cr_assert( result != SKY_OK );
  tx.Sigs.data = malloc(20 * sizeof(cipher__Sig));
  cr_assert( tx.Sigs.data != NULL );
  registerMemCleanup( tx.Sigs.data );
  memset( tx.Sigs.data, 0, 20 * sizeof(cipher__Sig) );
  tx.Sigs.len = 20; tx.Sigs.cap = 20;
  result = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(&tx);
  cr_assert( result != SKY_OK );

  int MaxUint16 = 0xFFFF;
  // Too many sigs & inputs
  makeTransaction(&tx);
  tx.Sigs.data = malloc(MaxUint16 * sizeof(cipher__Sig));
  cr_assert( tx.Sigs.data != NULL );
  registerMemCleanup( tx.Sigs.data );
  memset(tx.Sigs.data, 0, MaxUint16 * sizeof(cipher__Sig));
  tx.Sigs.len = tx.Sigs.cap = MaxUint16;
  tx.In.data = malloc( MaxUint16 * sizeof(cipher__SHA256) );
  cr_assert( tx.In.data != NULL );
  registerMemCleanup( tx.In.data );
  tx.In.len = tx.In.cap = MaxUint16;
  result = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(&tx);
  cr_assert( result != SKY_OK );

  freeRegisteredMemCleanup(tx.Sigs.data); //Too much memory (~8MB) free ASAP
  freeRegisteredMemCleanup(tx.In.data);

  // Duplicate inputs
  coin__UxOut ux;
  cipher__SecKey seckey;
  cipher__SHA256 sha256;
  makeUxOutWithSecret( &ux, &seckey );
  makeTransactionFromUxOut( &ux, &seckey, &tx );
  memcpy(&sha256, tx.In.data, sizeof(cipher__SHA256));
  GoUint16 r;
  result = SKY_coin_Transaction_PushInput(&tx, &sha256, &r);
  tx.Sigs.len = 0;
  GoSlice seckeys;
  seckeys.data = malloc(sizeof(cipher__SecKey) * 2);
  cr_assert( seckeys.data != NULL );
  registerMemCleanup( seckeys.data );
  seckeys.len = seckeys.cap = 2;
  memcpy( seckeys.data, &seckey, sizeof(cipher__SecKey) );
  memcpy( ((cipher__SecKey*)seckeys.data) + 1, &seckey, sizeof(cipher__SecKey) );
  result = SKY_coin_Transaction_SignInputs( &tx, seckeys );
  result = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(&tx);
  cr_assert( result != SKY_OK );

  //Duplicate outputs
  makeTransaction(&tx);
  coin__TransactionOutput* pOutput = tx.Out.data;
  cipher__Address addr;
  memcpy(&addr, &pOutput->Address, sizeof(cipher__Address));
  result = SKY_coin_Transaction_PushOutput(&tx, &addr, pOutput->Coins, pOutput->Hours);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(&tx);
  cr_assert( result != SKY_OK );

  // Invalid signature, empty
  makeTransaction(&tx);
  memset(tx.Sigs.data, 0, sizeof(cipher__Sig));
  result = SKY_coin_Transaction_Verify(&tx);
  cr_assert( result != SKY_OK );

  // Output coins are 0
  makeTransaction(&tx);
  pOutput = tx.Out.data;
  pOutput->Coins = 0;
  result = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(&tx);
  cr_assert( result != SKY_OK );

  GoUint64 MaxUint64 = 0xFFFFFFFFFFFFFFFF;
  // Output coin overflow
  makeTransaction(&tx);
  pOutput = tx.Out.data;
  pOutput->Coins = MaxUint64 - 3000000;
  result = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(&tx);
  cr_assert( result != SKY_OK );

  // Output coins are not multiples of 1e6 (valid, decimal restriction is not enforced here)
  makeTransaction(&tx);
  pOutput = tx.Out.data;
  pOutput->Coins += 10;
  result = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert( result == SKY_OK );
  tx.Sigs.data = NULL; tx.Sigs.len = 0; tx.Sigs.cap = 0;
  cipher__PubKey pubkey;
  result = SKY_cipher_GenerateKeyPair(&pubkey, &seckey);
  cr_assert( result == SKY_OK );
  seckeys.data = &seckey; seckeys.len = 1; seckeys.cap = 1;
  result = SKY_coin_Transaction_SignInputs(&tx, seckeys);
  cr_assert( result == SKY_OK );
  cr_assert( pOutput->Coins % 1000000 != 0 );
  result = SKY_coin_Transaction_Verify(&tx);
  cr_assert( result == SKY_OK );

  //Valid
  makeTransaction(&tx);
  pOutput = tx.Out.data;
  pOutput->Coins = 10000000;
  pOutput++;
  pOutput->Coins = 1000000;
  result = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_Verify(&tx);
  cr_assert( result == SKY_OK );
}


Test(coin_transaction, TestTransactionPushInput)
{
  int result;
  coin__Transaction tx;
  coin__UxOut ux;
  memset( &tx, 0, sizeof(coin__Transaction) );
  makeUxOut( &ux );
  cipher__SHA256 hash;
  result = SKY_coin_UxOut_Hash( &ux, &hash );
  cr_assert( result == SKY_OK );
  GoUint16 r;
  result = SKY_coin_Transaction_PushInput(&tx, &hash, &r);
  cr_assert( result == SKY_OK );
  cr_assert( r == 0 );
  cr_assert( tx.In.len == 1 );
  cipher__SHA256* pIn = tx.In.data;
  cr_assert( eq( u8[sizeof(cipher__SHA256)], hash, *pIn) );

  GoUint16 MaxUint16 = 0xFFFF;
  void* data = malloc( (tx.In.len + MaxUint16) * sizeof(cipher__SHA256) );
  cr_assert( data != NULL );
  registerMemCleanup( data );
  memset( data, 0, (tx.In.len + MaxUint16) * sizeof(cipher__SHA256) );
  memcpy( data, tx.In.data,  (tx.In.len) * sizeof(cipher__SHA256) );
  tx.In.len += MaxUint16;
  tx.In.cap = tx.In.len;
  tx.In.data = data;
  makeUxOut( &ux );
  result = SKY_coin_UxOut_Hash( &ux, &hash );
  cr_assert( result == SKY_OK );
  result = SKY_coin_Transaction_PushInput(&tx, &hash, &r);
  cr_assert( result != SKY_OK );

  freeRegisteredMemCleanup( data );
}


Test(coin_transaction, TestTransactionPushOutput)
{
  int result;
  coin__Transaction tx;
  cipher__Address addr;
  memset( &tx, 0, sizeof(coin__Transaction) );
  makeAddress( &addr );
  result = SKY_coin_Transaction_PushOutput( &tx, &addr, 100, 150 );
  cr_assert( result == SKY_OK );
  cr_assert( tx.Out.len == 1 );
  coin__TransactionOutput* pOutput = tx.Out.data;
  coin__TransactionOutput output;
  memcpy(&output.Address, &addr, sizeof(cipher__Address));
  output.Coins = 100;
  output.Hours = 150;
  cr_assert( eq( type(coin__TransactionOutput), output, *pOutput ) );
  for(int i = 1; i < 20; i++){
    makeAddress( &addr );
    result = SKY_coin_Transaction_PushOutput( &tx, &addr, i * 100, i * 50 );
    cr_assert( result == SKY_OK );
    cr_assert( tx.Out.len == i + 1 );
    pOutput = tx.Out.data;
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
  coin__Transaction tx;
  makeTransaction(&tx);
  cipher__SHA256 nullHash, hash1, hash2;
  memset( &nullHash, 0, sizeof(cipher__SHA256) );
  result = SKY_coin_Transaction_Hash( &tx, &hash1 );
  cr_assert( result == SKY_OK );
  cr_assert( not ( eq( u8[sizeof(cipher__SHA256)], nullHash, hash1) ) );
  result = SKY_coin_Transaction_HashInner( &tx, &hash2 );
  cr_assert( result == SKY_OK );
  cr_assert( not ( eq( u8[sizeof(cipher__SHA256)], hash2, hash1) ) );
}

Test(coin_transaction, TestTransactionUpdateHeader)
{
  int result;
  coin__Transaction tx;
  makeTransaction(&tx);
  cipher__SHA256 hash, nullHash, hashInner;
  memcpy(&hash, &tx.InnerHash, sizeof(cipher__SHA256));
  memset(&tx.InnerHash, 0, sizeof(cipher__SHA256));
  memset(&nullHash, 0, sizeof(cipher__SHA256));
  result = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert( not ( eq( u8[sizeof(cipher__SHA256)], tx.InnerHash, nullHash) ) );
  cr_assert( eq( u8[sizeof(cipher__SHA256)], hash, tx.InnerHash) );
  result = SKY_coin_Transaction_HashInner( &tx, &hashInner );
  cr_assert( result == SKY_OK );
  cr_assert( eq( u8[sizeof(cipher__SHA256)], hashInner, tx.InnerHash) );
}
