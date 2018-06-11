
#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"
#include "transutil.h"
#include "skycriterion.h"
#include "time.h"

TestSuite(coin_block, .init = setup, .fini = teardown);

Transactions__Handle makeTestTransactions(){
  Transactions__Handle transactions;
  Transaction__Handle transaction;

  int result = SKY_coin_Create_Transactions(&transactions);
  cr_assert(result == SKY_OK, "SKY_coin_Create_Transactions failed");
  registerHandleClose( transactions );
  result = SKY_coin_Create_Transaction(&transaction);
  cr_assert(result == SKY_OK, "SKY_coin_Create_Transaction failed");
  registerHandleClose( transaction );
  result = SKY_coin_Transactions_Add(transactions, transaction);
  cr_assert(result == SKY_OK, "SKY_coin_Transactions_Add failed");
  return transactions;
}

int makeNewBlock(cipher__SHA256* uxHash, Block__Handle* newBlock){
  int result;
  cipher__SHA256 bodyhash;
  BlockBody__Handle block;
  Transactions__Handle transactions = makeTestTransactions();

  result = SKY_coin_NewEmptyBlock(transactions, &block);
  cr_assert(result == SKY_OK, "SKY_coin_NewEmptyBlock failed");
  registerHandleClose( block );
  coin__Block* pBlock;
  result = SKY_coin_GetBlockObject(block, &pBlock);
  cr_assert(result == SKY_OK, "SKY_coin_Get_Block_Object failed");

  pBlock->Head.Version = 0x02;
  pBlock->Head.Time = 100;
  pBlock->Head.BkSeq = 0;
  pBlock->Head.Fee = 10;
  BlockBody__Handle body;
  result = SKY_coin_GetBlockBody(block, &body);
  cr_assert(result == SKY_OK, "SKY_coin_Get_Block_Body failed");
  result = SKY_coin_BlockBody_Hash(body, &bodyhash);
  cr_assert(result == SKY_OK, "SKY_coin_BlockBody_Hash failed");
  result = SKY_coin_NewBlock(block, 100 + 200, uxHash, transactions, 0, newBlock);
  cr_assert(result == SKY_OK, "SKY_coin_NewBlock failed");
  registerHandleClose( *newBlock );
  return result;
}

Test(coin_block, TestNewBlock) {
  Block__Handle prevBlock = 0;
  Block__Handle newBlock = 0;
  coin__Block* pPrevBlock = NULL;
  coin__Block* pNewBlock = NULL;
  int result = 0;

  Transactions__Handle transactions = makeTestTransactions();
  result = SKY_coin_NewEmptyBlock(transactions, &prevBlock);
  cr_assert(result == SKY_OK, "SKY_coin_NewEmptyBlock failed");
  registerHandleClose( prevBlock );
  coin__Block* pBlock;
  result = SKY_coin_GetBlockObject(prevBlock, &pPrevBlock);
  cr_assert(result == SKY_OK, "SKY_coin_GetBlockObject failed");

  pPrevBlock->Head.Version = 0x02;
  pPrevBlock->Head.Time = 100;
  pPrevBlock->Head.BkSeq = 98;


  GoSlice slice;
  memset(&slice, 0, sizeof(GoSlice));
  cipher__SHA256 hash;

  result = SKY_cipher_RandByte( 128, (coin__UxArray*)&slice );
  cr_assert(result == SKY_OK, "SKY_cipher_RandByte failed");
  registerMemCleanup( slice.data );
  result = SKY_cipher_SumSHA256( slice, &hash );
  cr_assert(result == SKY_OK, "SKY_cipher_SumSHA256 failed");

  result = SKY_coin_NewBlock(prevBlock, 133, &hash, 0, 0, &newBlock);
  cr_assert(result != SKY_OK, "SKY_coin_NewBlock has to fail with no transactions");
  registerHandleClose( newBlock );

  transactions = 0;
  Transaction__Handle tx = 0;
  result = SKY_coin_Create_Transactions(&transactions);
  cr_assert(result == SKY_OK, "SKY_coin_Create_Transactions failed");
  registerHandleClose(transactions);
  makeEmptyTransaction(&tx);
  registerHandleClose(tx);
  result = SKY_coin_Transactions_Add(transactions, tx);
  cr_assert(result == SKY_OK, "SKY_coin_Transactions_Add failed");

  GoUint64 fee = 121;
  GoUint64 currentTime = 133;

  result = SKY_coin_NewBlock(prevBlock, currentTime, &hash, transactions, fee, &newBlock);
  cr_assert(result == SKY_OK, "SKY_coin_NewBlock failed");
  registerHandleClose(newBlock);
  result = SKY_coin_GetBlockObject(newBlock, &pNewBlock);
  cr_assert(result == SKY_OK, "SKY_coin_GetBlockObject failed");
  coin__Transactions* pTransactions = NULL;
  SKY_coin_Get_Transactions_Object(transactions, &pTransactions);
  cr_assert( eq( type(GoSlice), *((GoSlice*)&pNewBlock->Body.Transactions), *((GoSlice*)pTransactions)) );
  cr_assert( eq(pNewBlock->Head.Fee, fee * (GoUint64)( pTransactions->len )));
  cr_assert( eq(pNewBlock->Head.Time, currentTime));
  cr_assert( eq(pNewBlock->Head.BkSeq, pPrevBlock->Head.BkSeq + 1));
  cr_assert( eq( u8[sizeof(cipher__SHA256)], pNewBlock->Head.UxHash, hash) );
}


Test(coin_block, TestBlockHashHeader){
  int result;
  Block__Handle block = 0;
  coin__Block* pBlock = NULL;
  GoSlice slice;
  memset(&slice, 0, sizeof(GoSlice));
  cipher__SHA256 hash;

  result = SKY_cipher_RandByte( 128, (coin__UxArray*)&slice );
  cr_assert(result == SKY_OK, "SKY_cipher_RandByte failed");
  registerMemCleanup( slice.data );
  result = SKY_cipher_SumSHA256( slice, &hash );
  cr_assert(result == SKY_OK, "SKY_cipher_SumSHA256 failed");
  result = makeNewBlock( &hash, &block );
  cr_assert(result == SKY_OK, "makeNewBlock failed");
  result = SKY_coin_GetBlockObject(block, &pBlock);
  cr_assert(result == SKY_OK, "SKY_coin_GetBlockObject failed, block handle : %d", block);

  cipher__SHA256 hash1, hash2;
  result = SKY_coin_Block_HashHeader(block, &hash1);
  cr_assert(result == SKY_OK, "SKY_coin_Block_HashHeader failed");
  result = SKY_coin_BlockHeader_Hash(&pBlock->Head, &hash2);
  cr_assert(result == SKY_OK, "SKY_coin_BlockHeader_Hash failed");
  cr_assert( eq( u8[sizeof(cipher__SHA256)],hash1, hash2) );
  memset(&hash2, 0, sizeof(cipher__SHA256));
  cr_assert( not( eq( u8[sizeof(cipher__SHA256)],hash1, hash2) ) );
}


Test(coin_block, TestBlockHashBody){
  int result;
  Block__Handle block = 0;
  GoSlice slice;
  memset(&slice, 0, sizeof(GoSlice));
  cipher__SHA256 hash;

  result = SKY_cipher_RandByte( 128, (coin__UxArray*)&slice );
  cr_assert(result == SKY_OK, "SKY_cipher_RandByte failed");
  registerMemCleanup( slice.data );
  result = SKY_cipher_SumSHA256( slice, &hash );
  cr_assert(result == SKY_OK, "SKY_cipher_SumSHA256 failed");
  result = makeNewBlock( &hash, &block );
  cr_assert(result == SKY_OK, "makeNewBlock failed");

  cipher__SHA256 hash1, hash2;
  result = SKY_coin_Block_HashBody(block, &hash1);
  cr_assert(result == SKY_OK, "SKY_coin_BlockBody_Hash failed");
  BlockBody__Handle blockBody = 0;
  result = SKY_coin_GetBlockBody(block, &blockBody);
  cr_assert(result == SKY_OK, "SKY_coin_GetBlockBody failed");
  result = SKY_coin_BlockBody_Hash(blockBody, &hash2);
  cr_assert(result == SKY_OK, "SKY_coin_BlockBody_Hash failed");
  cr_assert( eq( u8[sizeof(cipher__SHA256)], hash1, hash2) );
}

Test(coin_block, TestNewGenesisBlock){
  cipher__PubKey pubkey;
  cipher__SecKey seckey;
  cipher__Address address;
  GoUint64 genTime = 1000;
  GoUint64 genCoins = 1000 * 1000 * 1000;
  GoUint64 genCoinHours = 1000 * 1000;
  Block__Handle block = 0;
  coin__Block* pBlock = NULL;

  int result = makeKeysAndAddress(&pubkey, &seckey, &address);
  cr_assert(result == SKY_OK, "makeKeysAndAddress failed");
  result = SKY_coin_NewGenesisBlock(&address, genCoins, genTime, &block);
  cr_assert(result == SKY_OK, "SKY_coin_NewGenesisBlock failed");
  result = SKY_coin_GetBlockObject(block, &pBlock);
  cr_assert(result == SKY_OK, "SKY_coin_GetBlockObject failed");

  cipher__SHA256 nullHash;
  memset(&nullHash, 0, sizeof(cipher__SHA256));
  cr_assert( eq( u8[sizeof(cipher__SHA256)], nullHash, pBlock->Head.PrevHash) );
  cr_assert( genTime == pBlock->Head.Time );
  cr_assert( 0 == pBlock->Head.BkSeq );
  cr_assert( 0 == pBlock->Head.Version );
  cr_assert( 0 == pBlock->Head.Fee );
  cr_assert( eq( u8[sizeof(cipher__SHA256)], nullHash, pBlock->Head.UxHash) );

  cr_assert( 1 == pBlock->Body.Transactions.len );
  coin__Transaction* ptransaction = (coin__Transaction*)pBlock->Body.Transactions.data;
  cr_assert( 0 == ptransaction->In.len);
  cr_assert( 0 == ptransaction->Sigs.len);
  cr_assert( 1 == ptransaction->Out.len);

  coin__TransactionOutput* poutput = (coin__TransactionOutput*)ptransaction->Out.data;
  cr_assert( eq( type(cipher__Address), address, poutput->Address ) );
  cr_assert( genCoins == poutput->Coins );
  cr_assert( genCoins == poutput->Hours );
}

typedef struct {
  int index;
  int failure;
} testcase_unspent;

Test(coin_block, TestCreateUnspent){
  cipher__PubKey pubkey;
  cipher__SecKey seckey;
  cipher__Address address;
  int result = makeKeysAndAddress(&pubkey, &seckey, &address);

  cipher__SHA256 hash;
  coin__Transaction* ptx;
  Transaction__Handle handle;
  ptx = makeEmptyTransaction(&handle);
  result = SKY_coin_Transaction_PushOutput(handle, &address, 11000000, 255);
  cr_assert(result == SKY_OK, "SKY_coin_Transaction_PushOutput failed");
  coin__BlockHeader bh;
  memset(&bh, 0,  sizeof(coin__BlockHeader));
  bh.Time = time(0);
  bh.BkSeq = 1;

  testcase_unspent t[] = {
    {0, 0}, {10, 1},
  };
  coin__UxOut ux;
  int tests_count = sizeof(t) / sizeof(testcase_unspent);
  for( int i = 0; i <  tests_count; i++){
    memset(&ux, 0, sizeof(coin__UxOut));
    result = SKY_coin_CreateUnspent( &bh, handle, t[i].index, &ux );
    if( t[i].failure ){
      cr_assert( result == SKY_ERROR, "SKY_coin_CreateUnspent should have failed" );
      continue;
    } else {
      cr_assert( result == SKY_OK, "SKY_coin_CreateUnspent failed" );
    }
    cr_assert( bh.Time == ux.Head.Time );
    cr_assert( bh.BkSeq == ux.Head.BkSeq );
    result = SKY_coin_Transaction_Hash( handle, &hash );
    cr_assert( result == SKY_OK, "SKY_coin_Transaction_Hash failed" );
    cr_assert( eq( u8[sizeof(cipher__SHA256)], hash, ux.Body.SrcTransaction) );
    cr_assert( t[i].index < ptx->Out.len);
    coin__TransactionOutput* poutput = (coin__TransactionOutput*)ptx->Out.data;
    cr_assert( eq( type(cipher__Address), ux.Body.Address, poutput->Address ) );
    cr_assert( ux.Body.Coins == poutput->Coins );
    cr_assert( ux.Body.Hours == poutput->Hours );
  }
}

Test(coin_block, TestCreateUnspents){
  cipher__PubKey pubkey;
  cipher__SecKey seckey;
  cipher__Address address;
  int result = makeKeysAndAddress(&pubkey, &seckey, &address);

  cipher__SHA256 hash;
  coin__Transaction* ptx;
  Transaction__Handle handle;
  ptx = makeEmptyTransaction(&handle);
  result = SKY_coin_Transaction_PushOutput(handle, &address, 11000000, 255);
  cr_assert(result == SKY_OK, "SKY_coin_Transaction_PushOutput failed");
  coin__BlockHeader bh;
  memset(&bh, 0,  sizeof(coin__BlockHeader));
  bh.Time = time(0);
  bh.BkSeq = 1;

  coin__UxArray uxs = {NULL, 0, 0};
  result = SKY_coin_CreateUnspents(&bh, handle, &uxs);
  cr_assert( result == SKY_OK, "SKY_coin_CreateUnspents failed" );
  registerMemCleanup( uxs.data );
  cr_assert( uxs.len == 1 );
  cr_assert( uxs.len == ptx->Out.len );
  coin__UxOut* pout = (coin__UxOut*)uxs.data;
  coin__TransactionOutput* ptxout = (coin__TransactionOutput*)ptx->Out.data;
  for(int i = 0; i < uxs.len; i++){
    cr_assert( bh.Time == pout->Head.Time );
    cr_assert( bh.BkSeq == pout->Head.BkSeq );
    result = SKY_coin_Transaction_Hash( handle, &hash );
    cr_assert( result == SKY_OK, "SKY_coin_Transaction_Hash failed" );
    cr_assert( eq( u8[sizeof(cipher__SHA256)], hash, pout->Body.SrcTransaction) );
    cr_assert( eq( type(cipher__Address), pout->Body.Address, ptxout->Address ) );
    cr_assert( pout->Body.Coins == ptxout->Coins );
    cr_assert( pout->Body.Hours == ptxout->Hours );
    pout++;
    ptxout++;
  }
}
