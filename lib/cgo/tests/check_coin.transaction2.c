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

Test(coin_outputs, TestTransactionsTruncateBytesTo){
    //SKY_coin_Transactions_TruncateBytesTo
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

Test(coin_outputs, TestVerifyTransactionCoinsSpending){
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

Test(coin_outputs, TestVerifyTransactionHoursSpending){
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
