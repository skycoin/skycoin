
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

int makeNewBlock(cipher__SHA256* uxHash, coin__Block* newBlock){
  int result;
  cipher__SHA256 bodyhash;
  coin__Block block;
  coin__Transactions transactions;

  memset(&block, 0, sizeof(coin__Block));
  memset(&transactions, 0, sizeof(coin__Transactions));
  memset(newBlock, 0, sizeof(coin__Block));

  block.Head.Version = 0x02;
  block.Head.Time = 100;
  block.Head.BkSeq = 0;
  block.Head.Fee = 10;
  result = SKY_coin_BlockBody_Hash(&block.Body, &bodyhash);
  cr_assert(result == SKY_OK, "SKY_coin_BlockBody_Hash failed");
  result = SKY_coin_NewBlock(&block, 100 + 200, uxHash, &transactions, 0, newBlock);
  cr_assert(result == SKY_OK, "SKY_coin_NewBlock failed");
  return result;
}

Test(coin_block, TestNewBlock) {
  coin__Block prevBlock;
  coin__Block newBlock;
  coin__Transactions transactions;
  memset(&prevBlock, 0, sizeof(coin__Block));
  memset(&newBlock, 0, sizeof(coin__Block));
  memset(&transactions, 0, sizeof(coin__Transactions));

  prevBlock.Head.Version = 0x02;
  prevBlock.Head.Time = 100;
  prevBlock.Head.BkSeq = 98;
  int result;

  GoSlice slice;
  memset(&slice, 0, sizeof(GoSlice));
  cipher__SHA256 hash;

  result = SKY_cipher_RandByte( 128, (coin__UxArray*)&slice );
  cr_assert(result == SKY_OK, "SKY_cipher_RandByte failed");
  registerMemCleanup( slice.data );
  result = SKY_cipher_SumSHA256( slice, &hash );
  cr_assert(result == SKY_OK, "SKY_cipher_SumSHA256 failed");

  result = SKY_coin_NewBlock(&prevBlock, 133, &hash, NULL, 0, &newBlock);
  cr_assert(result != SKY_OK, "SKY_coin_NewBlock has to fail with no transactions");
  transactions.data = malloc(sizeof(coin__Transaction));
  registerMemCleanup(transactions.data);
  transactions.len = 1;
  transactions.cap = 1;
  memset(transactions.data, 0, sizeof(coin__Transaction));
  GoUint64 fee = 121;
  GoUint64 currentTime = 133;

  result = SKY_coin_NewBlock(&prevBlock, currentTime, &hash, &transactions, fee, &newBlock);
  cr_assert(result == SKY_OK, "SKY_coin_NewBlock failed");
  cr_assert( eq( type(GoSlice), *((GoSlice*)&newBlock.Body.Transactions), *((GoSlice*)&transactions)) );
  cr_assert( eq(newBlock.Head.Fee, fee * (GoUint64)( transactions.len )));
  cr_assert( eq(newBlock.Head.Time, currentTime));
  cr_assert( eq(newBlock.Head.BkSeq, prevBlock.Head.BkSeq + 1));
  cr_assert( eq( u8[sizeof(cipher__SHA256)], newBlock.Head.UxHash, hash) );

  coin__BlockBody body;
  memset(&body, 0,  sizeof(coin__BlockBody));
  body.Transactions.data = transactions.data;
  body.Transactions.len = transactions.len;
  body.Transactions.cap = transactions.cap;
  cr_assert( eq(type(coin__BlockBody), newBlock.Body, body) );
}

Test(coin_block, TestBlockHashHeader){
  int result;
  coin__Block block;
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
  result = SKY_coin_Block_HashHeader(&block, &hash1);
  cr_assert(result == SKY_OK, "SKY_coin_Block_HashHeader failed");
  result = SKY_coin_BlockHeader_Hash(&block.Head, &hash2);
  cr_assert(result == SKY_OK, "SKY_coin_BlockHeader_Hash failed");
  cr_assert( eq( u8[sizeof(cipher__SHA256)],hash1, hash2) );
  memset(&hash2, 0, sizeof(cipher__SHA256));
  cr_assert( not( eq( u8[sizeof(cipher__SHA256)],hash1, hash2) ) );
}

Test(coin_block, TestBlockHashBody){
  int result;
  coin__Block block;
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
  result = SKY_coin_Block_HashBody(&block, &hash1);
  cr_assert(result == SKY_OK, "SKY_coin_BlockBody_Hash failed");
  result = SKY_coin_BlockBody_Hash(&block.Body, &hash2);
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
  coin__Block block;

  memset(&block, 0, sizeof(coin__Block));
  int result = makeKeysAndAddress(&pubkey, &seckey, &address);
  cr_assert(result == SKY_OK, "makeKeysAndAddress failed");
  result = SKY_coin_NewGenesisBlock(&address, genCoins, genTime, &block);
  cr_assert(result == SKY_OK, "SKY_coin_NewGenesisBlock failed");

  cipher__SHA256 nullHash;
  memset(&nullHash, 0, sizeof(cipher__SHA256));
  cr_assert( eq( u8[sizeof(cipher__SHA256)], nullHash, block.Head.PrevHash) );
  cr_assert( genTime == block.Head.Time );
  cr_assert( 0 == block.Head.BkSeq );
  cr_assert( 0 == block.Head.Version );
  cr_assert( 0 == block.Head.Fee );
  cr_assert( eq( u8[sizeof(cipher__SHA256)], nullHash, block.Head.UxHash) );

  cr_assert( 1 == block.Body.Transactions.len );
  coin__Transaction* ptransaction = (coin__Transaction*)block.Body.Transactions.data;
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
  coin__Transaction tx;
  memset( &tx, 0, sizeof(coin__Transaction) );
  result = SKY_coin_Transaction_PushOutput(&tx, &address, 11000000, 255);
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
    result = SKY_coin_CreateUnspent( &bh, &tx, t[i].index, &ux );
    if( t[i].failure ){
      cr_assert( result == SKY_ERROR, "SKY_coin_CreateUnspent should have failed" );
      continue;
    } else {
      cr_assert( result == SKY_OK, "SKY_coin_CreateUnspent failed" );
    }
    cr_assert( bh.Time == ux.Head.Time );
    cr_assert( bh.BkSeq == ux.Head.BkSeq );
    result = SKY_coin_Transaction_Hash( &tx, &hash );
    cr_assert( result == SKY_OK, "SKY_coin_Transaction_Hash failed" );
    cr_assert( eq( u8[sizeof(cipher__SHA256)], hash, ux.Body.SrcTransaction) );
    cr_assert( t[i].index < tx.Out.len);
    coin__TransactionOutput* poutput = (coin__TransactionOutput*)tx.Out.data;
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
  coin__Transaction tx;
  memset( &tx, 0, sizeof(coin__Transaction) );
  result = SKY_coin_Transaction_PushOutput(&tx, &address, 11000000, 255);
  cr_assert(result == SKY_OK, "SKY_coin_Transaction_PushOutput failed");
  coin__BlockHeader bh;
  memset(&bh, 0,  sizeof(coin__BlockHeader));
  bh.Time = time(0);
  bh.BkSeq = 1;

  coin__UxArray uxs = {NULL, 0, 0};
  result = SKY_coin_CreateUnspents(&bh, &tx, &uxs);
  cr_assert( result == SKY_OK, "SKY_coin_CreateUnspents failed" );
  registerMemCleanup( uxs.data );
  cr_assert( uxs.len == 1 );
  cr_assert( uxs.len == tx.Out.len );
  coin__UxOut* pout = (coin__UxOut*)uxs.data;
  coin__TransactionOutput* ptxout = (coin__TransactionOutput*)tx.Out.data;
  for(int i = 0; i < uxs.len; i++){
    cr_assert( bh.Time == pout->Head.Time );
    cr_assert( bh.BkSeq == pout->Head.BkSeq );
    result = SKY_coin_Transaction_Hash( &tx, &hash );
    cr_assert( result == SKY_OK, "SKY_coin_Transaction_Hash failed" );
    cr_assert( eq( u8[sizeof(cipher__SHA256)], hash, pout->Body.SrcTransaction) );
    cr_assert( eq( type(cipher__Address), pout->Body.Address, ptxout->Address ) );
    cr_assert( pout->Body.Coins == ptxout->Coins );
    cr_assert( pout->Body.Hours == ptxout->Hours );
    pout++;
    ptxout++;
  }
}
