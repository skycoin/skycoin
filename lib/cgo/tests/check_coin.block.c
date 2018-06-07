
#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

TestSuite(coin_block, .init = setup, .fini = teardown);

int newBlock(cipher__SHA256* uxHash, coin__Block* newBlock){
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
