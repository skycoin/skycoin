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

TestSuite(coin_outputs, .init = setup, .fini = teardown);

Test(coin_outputs, TestUxBodyHash){
  int result;
  coin__UxBody uxbody;
  result = makeUxBody(&uxbody);
  cr_assert( result == SKY_OK, "makeUxBody failed" );
  cipher__SHA256 hash, nullHash;
  result = SKY_coin_UxBody_Hash(&uxbody, &hash);
  cr_assert( result == SKY_OK, "SKY_coin_UxBody_Hash failed" );
  memset(&nullHash, 0, sizeof(cipher__SHA256));
  cr_assert( not( eq( u8[sizeof(cipher__SHA256)], nullHash, hash) ) );
}

Test(coin_outputs, TestUxOutHash){
  int result;
  coin__UxBody uxbody;
  result = makeUxBody(&uxbody);
  cr_assert( result == SKY_OK, "makeUxBody failed" );

  coin__UxOut uxout;
  memset(&uxout, 0, sizeof(coin__UxOut));
  memcpy(&uxout.Body, &uxbody, sizeof(coin__UxBody));

  cipher__SHA256 hashBody, hashOut;
  result = SKY_coin_UxBody_Hash(&uxbody, &hashBody);
  cr_assert( result == SKY_OK, "SKY_coin_UxBody_Hash failed" );
  result = SKY_coin_UxOut_Hash(&uxout, &hashOut);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_Hash failed" );
  cr_assert( eq( u8[sizeof(cipher__SHA256)], hashBody, hashOut) );

  //Head should not affect hash
  uxout.Head.Time = 0;
  uxout.Head.BkSeq = 1;
  result = SKY_coin_UxOut_Hash(&uxout, &hashOut);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_Hash failed" );
  cr_assert( eq( u8[sizeof(cipher__SHA256)], hashBody, hashOut) );
}

Test(coin_outputs, TestUxOutSnapshotHash){
  int result;
  coin__UxOut uxout, uxout2;
  result = makeUxOut(&uxout);
  cr_assert( result == SKY_OK, "makeUxOut failed" );
  cipher__SHA256 hash1, hash2;
  result = SKY_coin_UxOut_SnapshotHash(&uxout, &hash1);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_SnapshotHash failed" );

  memcpy( &uxout2, &uxout, sizeof(coin__UxOut) );
  uxout2.Head.Time = 20;
  result = SKY_coin_UxOut_SnapshotHash(&uxout2, &hash2);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_SnapshotHash failed" );
  cr_assert( not( eq( u8[sizeof(cipher__SHA256)], hash1, hash2) ) , "Snapshot hash must be different");

  memcpy( &uxout2, &uxout, sizeof(coin__UxOut) );
  uxout2.Head.BkSeq = 4;
  result = SKY_coin_UxOut_SnapshotHash(&uxout2, &hash2);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_SnapshotHash failed" );
  cr_assert( not( eq( u8[sizeof(cipher__SHA256)], hash1, hash2) ) , "Snapshot hash must be different");

  memcpy( &uxout2, &uxout, sizeof(coin__UxOut) );
  makeRandHash(&uxout2.Body.SrcTransaction);
  result = SKY_coin_UxOut_SnapshotHash(&uxout2, &hash2);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_SnapshotHash failed" );
  cr_assert( not( eq( u8[sizeof(cipher__SHA256)], hash1, hash2) ) , "Snapshot hash must be different");

  memcpy( &uxout2, &uxout, sizeof(coin__UxOut) );
  makeAddress(&uxout2.Body.Address);
  result = SKY_coin_UxOut_SnapshotHash(&uxout2, &hash2);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_SnapshotHash failed" );
  cr_assert( not( eq( u8[sizeof(cipher__SHA256)], hash1, hash2) ) , "Snapshot hash must be different");

  memcpy( &uxout2, &uxout, sizeof(coin__UxOut) );
  uxout2.Body.Coins = uxout.Body.Coins * 2;
  result = SKY_coin_UxOut_SnapshotHash(&uxout2, &hash2);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_SnapshotHash failed" );
  cr_assert( not( eq( u8[sizeof(cipher__SHA256)], hash1, hash2) ) , "Snapshot hash must be different");

  memcpy( &uxout2, &uxout, sizeof(coin__UxOut) );
  uxout2.Body.Hours = uxout.Body.Hours * 2;
  result = SKY_coin_UxOut_SnapshotHash(&uxout2, &hash2);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_SnapshotHash failed" );
  cr_assert( not( eq( u8[sizeof(cipher__SHA256)], hash1, hash2) ) , "Snapshot hash must be different");
}

Test(coin_outputs, TestUxOutCoinHours){
  GoUint64 _genCoins = 1000000000;
  GoUint64 _genCoinHours = 1000 * 1000;

  int result;
  coin__UxOut ux;
  result = makeUxOut(&ux);
  cr_assert( result == SKY_OK, "makeUxOut failed" );

  GoUint64 now, hours;

  //Less than an hour passed
  now = ux.Head.Time + 100;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( hours == ux.Body.Hours );

  //An hour passed
  now = ux.Head.Time + 3600;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( hours == ux.Body.Hours + ux.Body.Coins / 1000000  );

  //6 hours passed
  now = ux.Head.Time + 3600 * 6;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( hours == ux.Body.Hours + (ux.Body.Coins / 1000000) * 6  );

  //Time is backwards (treated as no hours passed)
  now = ux.Head.Time / 2;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( hours == ux.Body.Hours  );

  //1 hour has passed, output has 1.5 coins, should gain 1 coinhour
  ux.Body.Coins = 1500000;
  now = ux.Head.Time + 3600;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( hours == ux.Body.Hours + 1  );

  //2 hours have passed, output has 1.5 coins, should gain 3 coin hours
  ux.Body.Coins = 1500000;
  now = ux.Head.Time + 3600 * 2;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( hours == ux.Body.Hours + 3  );

  //1 second has passed, output has 3600 coins, should gain 1 coin hour
  ux.Body.Coins = 3600000000;
  now = ux.Head.Time + 1;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( hours == ux.Body.Hours + 1  );

  //1000000 hours minus 1 second have passed, output has 1 droplet, should gain 0 coin hour
  ux.Body.Coins = 1;
  now = ux.Head.Time + (GoUint64)(1000000)*(GoUint64)(3600)-1;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( hours == ux.Body.Hours );

  //1000000 hours have passed, output has 1 droplet, should gain 1 coin hour
  ux.Body.Coins = 1;
  now = ux.Head.Time + (GoUint64)(1000000)*(GoUint64)(3600);
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( hours == ux.Body.Hours + 1 );

  // No hours passed, using initial coin hours
  ux.Body.Coins = _genCoins;
  ux.Body.Hours = _genCoinHours;
  now = ux.Head.Time;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( hours == ux.Body.Hours );

  // One hour passed, using initial coin hours
  now = ux.Head.Time + 3600;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( hours == ux.Body.Hours + _genCoins / 1000000 );

  // No hours passed and no hours to begin with0
  ux.Body.Hours = 0;
  now = ux.Head.Time;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( hours == 0 );

  // Centuries have passed, time-based calculation overflows uint64
	// when calculating the whole coin seconds
  ux.Body.Coins = 2000000;
  now = 0xFFFFFFFFFFFFFFFF;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_ERROR, "SKY_coin_UxOut_CoinHours should fail" );

  // Centuries have passed, time-based calculation overflows uint64
	// when calculating the droplet seconds
  ux.Body.Coins = 1500000;
  now = 0xFFFFFFFFFFFFFFFF;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_ERROR, "SKY_coin_UxOut_CoinHours should fail" );

  // Output would overflow if given more hours, has reached its limit
  ux.Body.Coins = 3600000000;
  now = 0xFFFFFFFFFFFFFFFE;
  result = SKY_coin_UxOut_CoinHours(&ux, now, &hours);
  cr_assert( result == SKY_ERROR, "SKY_coin_UxOut_CoinHours should fail" );
}

Test(coin_outputs, TestUxArrayCoins){
  coin__UxArray uxs;
  int result = makeUxArray(&uxs, 4);
  cr_assert( result == SKY_OK, "makeUxArray failed" );
  GoUint64 coins;
  result = SKY_coin_UxArray_Coins( &uxs, &coins );
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_Coins failed" );
  cr_assert( coins == 4000000 );
  coin__UxOut* p = (coin__UxOut*)uxs.data;
  p += 2;
  p->Body.Coins = 0xFFFFFFFFFFFFFFFF - 1000000;
  result = SKY_coin_UxArray_Coins( &uxs, &coins );
  cr_assert( result != SKY_OK, "SKY_coin_UxArray_Coins should fail with overflow" );
}

Test(coin_outputs, TestUxArrayCoinHours){
  coin__UxArray uxs;
  int result = makeUxArray(&uxs, 4);
  cr_assert( result == SKY_OK, "makeUxArray failed" );
  coin__UxOut* p = (coin__UxOut*)uxs.data;
  GoUint64 n;

  result = SKY_coin_UxArray_CoinHours(&uxs, p->Head.Time, &n);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( n == 400 );

  result = SKY_coin_UxArray_CoinHours(&uxs, p->Head.Time + 3600, &n);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( n == 404 );

  result = SKY_coin_UxArray_CoinHours(&uxs, p->Head.Time + 3600 + 4600, &n);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_CoinHours failed" );
  cr_assert( n == 408 );

  p[2].Body.Hours = 0xFFFFFFFFFFFFFFFF - 100;
  result = SKY_coin_UxArray_CoinHours(&uxs, p->Head.Time, &n);
  cr_assert( result != SKY_OK, "SKY_coin_UxOut_CoinHours should have fail with overflow" );

  result = SKY_coin_UxArray_CoinHours(&uxs, p->Head.Time * (GoUint64)1000000000000, &n);
  cr_assert( result != SKY_OK, "SKY_coin_UxOut_CoinHours should have fail with overflow" );
}

Test(coin_outputs, TestUxArrayHashArray){
  coin__UxArray uxs;
  int result = makeUxArray(&uxs, 4);
  cr_assert( result == SKY_OK, "makeUxArray failed" );
  coin__UxOut* p = (coin__UxOut*)uxs.data;

  GoSlice_ hashes = {NULL, 0, 0};
  result = SKY_coin_UxArray_Hashes(&uxs, &hashes);
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_Hashes failed" );
  registerMemCleanup( hashes.data );
  cr_assert(hashes.len == uxs.len);
  coin__UxOut* pux = (coin__UxOut*)uxs.data;
  cipher__SHA256* ph = (cipher__SHA256*)hashes.data;
  cipher__SHA256 hash;
  for(int i = 0; i < hashes.len; i++){
    result = SKY_coin_UxOut_Hash(pux, &hash);
    cr_assert( result == SKY_OK, "SKY_coin_UxOut_Hash failed" );
    cr_assert( eq( u8[sizeof(cipher__SHA256)], hash, *ph) );
    pux++;
    ph++;
  }
}

Test(coin_outputs, TestUxArrayHasDupes){
  coin__UxArray uxs;
  int result = makeUxArray(&uxs, 4);
  cr_assert( result == SKY_OK, "makeUxArray failed" );
  GoUint8 hasDupes;
  result = SKY_coin_UxArray_HasDupes(&uxs, &hasDupes);
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_HasDupes failed" );
  cr_assert( hasDupes == 0 );
  coin__UxOut* p = (coin__UxOut*)uxs.data;
  p++;
  memcpy(uxs.data, p, sizeof(coin__UxOut));
  result = SKY_coin_UxArray_HasDupes(&uxs, &hasDupes);
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_HasDupes failed" );
  cr_assert( hasDupes != 0 );
}

Test(coin_outputs, TestUxArraySub){

  int result;
  coin__UxArray uxa, uxb, uxc, uxd;
  coin__UxArray t1, t2, t3, t4;

  int arraySize = sizeof(coin__UxArray);
  memset(&uxa, 0, arraySize); memset(&uxb, 0, arraySize);
  memset(&uxc, 0, arraySize); memset(&uxd, 0, arraySize);
  memset(&t1, 0, arraySize); memset(&t2, 0, arraySize);
  memset(&t3, 0, arraySize); memset(&t4, 0, arraySize);

  result = makeUxArray(&uxa, 4);
  cr_assert( result == SKY_OK, "makeUxArray failed" );
  result = makeUxArray(&uxb, 4);
  cr_assert( result == SKY_OK, "makeUxArray failed" );

  int elems_size = sizeof(coin__UxOut);
  cutSlice(&uxa, 0, 1, elems_size, &t1);
  cr_assert( result == SKY_OK, "cutSlice failed" );
  result = concatSlices( &t1, &uxb, elems_size, &t2 );
  cr_assert( result == SKY_OK, "concatSlices failed" );
  result = cutSlice(&uxa, 1, 2, elems_size, &t3);
  cr_assert( result == SKY_OK, "cutSlice failed" );
  result = concatSlices( &t2, &t3, elems_size, &uxc );
  cr_assert( result == SKY_OK, "concatSlices failed" );

  memset(&uxd, 0, arraySize);
  result = SKY_coin_UxArray_Sub(&uxc, &uxa, &uxd);
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_Sub failed" );
  registerMemCleanup( uxd.data );
  cr_assert( eq( type(GoSlice), *((GoSlice*)&uxd), *((GoSlice*)&uxb)) );

  memset(&uxd, 0, arraySize);
  result = SKY_coin_UxArray_Sub(&uxc, &uxb, &uxd);
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_Sub failed" );
  registerMemCleanup( uxd.data );
  cr_assert( uxd.len == 2, "uxd length must be 2 and it is: %s", uxd.len );
  cutSlice(&uxa, 0, 2, elems_size, &t1);
  cr_assert( eq( type(GoSlice), *((GoSlice*)&uxd), *((GoSlice*)&t1)) );

  // No intersection
  memset(&t1, 0, arraySize); memset(&t2, 0, arraySize);
  result = SKY_coin_UxArray_Sub(&uxa, &uxb, &t1);
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_Sub failed" );
  registerMemCleanup( t1.data );
  result = SKY_coin_UxArray_Sub(&uxb, &uxa, &t2);
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_Sub failed" );
  registerMemCleanup( t2.data );
  cr_assert( eq( type(GoSlice), *((GoSlice*)&uxa), *((GoSlice*)&t1)) );
  cr_assert( eq( type(GoSlice), *((GoSlice*)&uxb), *((GoSlice*)&t2)) );
}

int isUxArraySorted(coin__UxArray* uxa){
  int n = uxa->len;
  coin__UxOut* prev = uxa->data;
  coin__UxOut* current = prev;
  current++;
  cipher__SHA256 hash1, hash2;
  cipher__SHA256* prevHash = NULL;
  cipher__SHA256* currentHash = NULL;

  int result;
  for(int i = 1; i < n; i++){
    if(prevHash == NULL){
      result = SKY_coin_UxOut_Hash(prev, &hash1);
      cr_assert( result == SKY_OK, "SKY_coin_UxOut_Hash failed" );
      prevHash = &hash1;
    }
    if(currentHash == NULL)
      currentHash = &hash2;
    result = SKY_coin_UxOut_Hash(current, currentHash);
    cr_assert( result == SKY_OK, "SKY_coin_UxOut_Hash failed" );
    if( memcmp(prevHash, currentHash, sizeof(cipher__SHA256)) > 0)
      return 0; //Array is not sorted
    if(i % 2 != 0){
      prevHash = &hash2;
      currentHash = &hash1;
    } else {
      prevHash = &hash1;
      currentHash = &hash2;
    }
    prev++;
    current++;
  }
  return 1;
}

Test(coin_outputs, TestUxArraySorting){

  int result;
  coin__UxArray uxa;
  result = makeUxArray(&uxa, 4);
  cr_assert( result == SKY_OK, "makeUxArray failed" );
  int isSorted = isUxArraySorted(&uxa);
  if( isSorted ){ //If already sorted then break the order
    coin__UxOut temp;
    coin__UxOut* p = uxa.data;
    memcpy(&temp, p, sizeof(coin__UxOut));
    memcpy(p, p + 1, sizeof(coin__UxOut));
    memcpy(p + 1, &temp, sizeof(coin__UxOut));
  }
  isSorted = isUxArraySorted(&uxa);
  cr_assert( isSorted == 0);
  result = SKY_coin_UxArray_Sort( &uxa );
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_Sort failed" );
  isSorted = isUxArraySorted(&uxa);
  cr_assert( isSorted == 1);
}

Test(coin_outputs, TestUxArrayLen){
  int result;
  coin__UxArray uxa;
  result = makeUxArray(&uxa, 4);
  cr_assert( result == SKY_OK, "makeUxArray failed" );
  GoInt len;
  result = SKY_coin_UxArray_Len(&uxa, &len);
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_Len failed" );
  cr_assert( len == uxa.len );
  cr_assert( len == 4 );
}

Test(coin_outputs, TestUxArrayLess){
  int result;
  coin__UxArray uxa;
  result = makeUxArray(&uxa, 2);
  cr_assert( result == SKY_OK, "makeUxArray failed" );
  cipher__SHA256 hashes[2];
  coin__UxOut* p = uxa.data;
  result = SKY_coin_UxOut_Hash(p, &hashes[0]);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_Hash failed" );
  p++;
  result = SKY_coin_UxOut_Hash(p, &hashes[1]);
  cr_assert( result == SKY_OK, "SKY_coin_UxOut_Hash failed" );
  GoUint8 lessResult1, lessResult2;
  int memcmpResult;
  result = SKY_coin_UxArray_Less(&uxa, 0, 1, &lessResult1);
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_Less failed" );
  result = SKY_coin_UxArray_Less(&uxa, 1, 0, &lessResult2);
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_Less failed" );
  memcmpResult = memcmp( &hashes[0], &hashes[1], sizeof(cipher__SHA256) );
  int r;
  r = (lessResult1 == 1) == (memcmpResult < 0);
  cr_assert(r != 0);
  r = (lessResult2 == 1) == (memcmpResult > 0);
  cr_assert(r != 0);
}

Test(coin_outputs, TestUxArraySwap){
  int result;
  coin__UxArray uxa;
  result = makeUxArray(&uxa, 2);
  cr_assert( result == SKY_OK, "makeUxArray failed" );
  coin__UxOut uxx, uxy;
  coin__UxOut* p = uxa.data;
  memcpy(&uxx, p, sizeof(coin__UxOut));
  memcpy(&uxy, p + 1, sizeof(coin__UxOut));

  result = SKY_coin_UxArray_Swap(&uxa, 0, 1);
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_Swap failed" );
  cr_assert( eq(type(coin__UxOut), uxy, *p) );
  cr_assert( eq(type(coin__UxOut), uxx, *(p+1)) );

  result = SKY_coin_UxArray_Swap(&uxa, 0, 1);
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_Swap failed" );
  cr_assert( eq(type(coin__UxOut), uxy, *(p+1)) );
  cr_assert( eq(type(coin__UxOut), uxx, *p) );

  result = SKY_coin_UxArray_Swap(&uxa, 1, 0);
  cr_assert( result == SKY_OK, "SKY_coin_UxArray_Swap failed" );
  cr_assert( eq(type(coin__UxOut), uxy, *p) );
  cr_assert( eq(type(coin__UxOut), uxx, *(p+1)) );
}

/**********************************************************
* 6 Tests involvinf AddressUxOuts were not done
* because the corresponding functions were not exported
*************************************************************/
