#include <criterion/criterion.h>
#include <criterion/new/assert.h>
#include <signal.h>
#include <stdio.h>
#include <string.h>

#include "libskycoin.h"
#include "math.h"
#include "skycriterion.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"
#include "skytxn.h"

TestSuite(util_fee, .init = setup, .fini = teardown);
#define BUFFER_SIZE 1024
#define BurnFactor 2
unsigned long long MaxUint64 = 0xFFFFFFFFFFFFFFFF;
unsigned int MaxUint16 = 0xFFFF;
typedef struct {
  GoInt64 inputHours;
  GoInt64 outputHours;
  GoInt64 err;
} verifyTxFeeTestCase;

verifyTxFeeTestCase burnFactor2verifyTxFeeTestCase[] = {
    {0, 0, SKY_ErrTxnNoFee},
    {1, 0, SKY_OK},
    {1, 1, SKY_ErrTxnNoFee},
    {2, 0, SKY_OK},
    {2, 1, SKY_OK},
    {2, 2, SKY_ErrTxnNoFee},
    {3, 0, SKY_OK},
    {3, 1, SKY_OK},
    {3, 2, SKY_ErrTxnInsufficientFee},
    {3, 3, SKY_ErrTxnNoFee},
    {4, 0, SKY_OK},
    {4, 1, SKY_OK},
    {4, 2, SKY_OK},
    {4, 3, SKY_ErrTxnInsufficientFee},
    {4, 4, SKY_ErrTxnNoFee},
};

verifyTxFeeTestCase burnFactor3verifyTxFeeTestCase[] = {
    {0, 0, SKY_ErrTxnNoFee}, {1, 0, SKY_OK},
    {1, 1, SKY_ErrTxnNoFee}, {2, 0, SKY_OK},
    {2, 1, SKY_OK},          {2, 2, SKY_ErrTxnNoFee},
    {3, 0, SKY_OK},          {3, 1, SKY_OK},
    {3, 2, SKY_OK},          {3, 3, SKY_ErrTxnNoFee},
    {4, 0, SKY_OK},          {4, 1, SKY_OK},
    {4, 2, SKY_OK},          {4, 3, SKY_ErrTxnInsufficientFee},
    {4, 4, SKY_ErrTxnNoFee}, {5, 0, SKY_OK},
    {5, 1, SKY_OK},          {5, 2, SKY_OK},
    {5, 3, SKY_OK},          {5, 4, SKY_ErrTxnInsufficientFee},
    {5, 5, SKY_ErrTxnNoFee},
};
#define cases burnFactor2verifyTxFeeTestCase
Test(util_fee, TestVerifyTransactionFee) {
  Transaction__Handle *emptyTxn;
  makeEmptyTransaction(&emptyTxn);
  GoUint64 hours;
  GoUint64 err = SKY_coin_Transaction_OutputHours(emptyTxn, &hours);
  cr_assert(err == SKY_OK);
  cr_assert(hours == 0);

  // A txn with no outputs hours and no coinhours burn fee is valid
  err = SKY_fee_VerifyTransactionFee(emptyTxn, 0,2);
  cr_assert(err == SKY_ErrTxnNoFee);

  // A txn with no outputs hours but with a coinhours burn fee is valid
  err = SKY_fee_VerifyTransactionFee(emptyTxn, 100,2);
  cr_assert(err == SKY_OK);

  Transaction__Handle txn;
  makeEmptyTransaction(&txn);
  coin__TransactionOutput *txnOut;
  cipher__Address addr;
  makeAddress(&addr);
  err = SKY_coin_Transaction_PushOutput(txn, &addr, 0, 1000000);
  cr_assert(err == SKY_OK);
  err = SKY_coin_Transaction_PushOutput(txn, &addr, 0, 3000000);
  cr_assert(err == SKY_OK);

  err = SKY_coin_Transaction_OutputHours(txn, &hours);
  cr_assert(err == SKY_OK);
  cr_assert(hours == 4000000);

  // A txn with insufficient net coinhours burn fee is invalid
  err = SKY_fee_VerifyTransactionFee(txn, 0,2);
  cr_assert(err == SKY_ErrTxnNoFee);
  err = SKY_fee_VerifyTransactionFee(txn, 1,2);
  cr_assert(err == SKY_ErrTxnInsufficientFee);

  // A txn with sufficient net coinhours burn fee is valid
  err = SKY_coin_Transaction_OutputHours(txn, &hours);
  cr_assert(err == SKY_OK);
  err = SKY_fee_VerifyTransactionFee(txn, hours,2);
  cr_assert(err == SKY_OK);
  err = SKY_coin_Transaction_OutputHours(txn, &hours);
  cr_assert(err == SKY_OK);
  err = SKY_fee_VerifyTransactionFee(txn, (hours * 10),2);
  cr_assert(err == SKY_OK);

  // fee + hours overflows
  err = SKY_fee_VerifyTransactionFee(txn, (MaxUint64 - 3000000),2);
  cr_assert(err == SKY_ERROR);

  // txn has overflowing output hours
  err = SKY_coin_Transaction_PushOutput(txn, &addr, 0,
                                        (MaxUint64 - 1000000 - 3000000 + 1));
  cr_assert(err == SKY_OK);
  err = SKY_fee_VerifyTransactionFee(txn, 10,2);
  cr_assert(err == SKY_ERROR);

  int len = (sizeof(cases) / sizeof(verifyTxFeeTestCase));

  for (int i = 0; i < len; i++) {
    makeEmptyTransaction(&txn);
    verifyTxFeeTestCase tc = cases[i];
    err = SKY_coin_Transaction_PushOutput(txn, &addr, 0, tc.outputHours);
    cr_assert(err == SKY_OK);
    cr_assert(tc.inputHours >= tc.outputHours);
    err = SKY_fee_VerifyTransactionFee(txn, (tc.inputHours - tc.outputHours),2);
    cr_assert(tc.err == err, "Iter %d is %x != %x", i, tc.err, err);
  }
}

typedef struct {
  GoInt64 hours;
  GoInt64 fee;
} requiredFeeTestCase;

requiredFeeTestCase burnFactor2RequiredFeeTestCases[] = {
    {0, 0}, {1, 1}, {2, 1},     {3, 2},     {4, 2},      {5, 3},
    {6, 3}, {7, 4}, {998, 499}, {999, 500}, {1000, 500}, {1001, 501},
};
#define cases1 burnFactor2RequiredFeeTestCases
Test(util_fee, TestRequiredFee) {
  int len = (sizeof(cases1) / sizeof(requiredFeeTestCase));
  for (int i = 0; i < len; i++) {
    requiredFeeTestCase tc = cases1[i];
    GoUint64 fee;
    GoUint64 err = SKY_fee_RequiredFee(tc.hours,2, &fee);
    cr_assert(err == SKY_OK);
    cr_assert(tc.fee == fee);

    GoUint64 remainingHours;
    err = SKY_fee_RemainingHours(tc.hours,2, &remainingHours);
    cr_assert(err == SKY_OK);
    cr_assert(eq(ullong, (tc.hours - fee), remainingHours));
  }
}

Test(util_fee, TestTransactionFee) {
  GoUint64 headTime = 1000;
  GoUint64 nextTime = headTime + 3600;

  typedef struct {
    GoUint64 times;
    GoUint64 coins;
    GoUint64 hours;
  } uxInput;

  typedef struct {
    GoString name;
    GoUint64 out[BUFFER_SIZE];
    uxInput in[BUFFER_SIZE];
    GoUint64 headTime;
    GoUint64 fee;
    GoUint64 err;
    int lens[2];

  } tmpstruct;

  tmpstruct cases[] = {
      // Test case with multiple outputs, multiple inputs
      {{" ", 1}, {5}, {{headTime, 10000000, 10}}, headTime, 5, .err = 0, {1, 1}

      },
      // Test case with multiple outputs, multiple inputs
      {{"", 1},
       .out = {5, 7, 3},
       .in =
           {
               {headTime, 10000000, 10},
               {headTime, 10000000, 5},
           },
       .headTime = headTime,
       .fee = 0,
       .err = 0,
       .lens = {2, 3}},
      // Test case with multiple outputs, multiple inputs, and some inputs have
      // more CoinHours once adjusted for HeadTime
      {{" ", 1},
       .out = {5, 10},
       .in =
           {
               {nextTime, 10000000, 10},
               {headTime, 8000000, 5},
           },
       nextTime,
       8,
       0,
       {2, 2}},
      // // Test case with insufficient coin hours
      {
          .err = SKY_ErrTxnInsufficientCoinHours,
          .out = {5, 10, 1},
          .in = {{headTime, 10000000, 10}, {headTime, 8000000, 5}},
          .headTime = headTime,
          .lens = {2, 3},
          .fee = 1,
      },
      // // Test case with overflowing input hours
      {.err = SKY_ERROR,
       .out = {0},
       .in = {{headTime, 10000000, 10}, {headTime, 10000000, (MaxUint64 - 9)}},
       .headTime = headTime,
       .lens = {2, 1}},
      // // Test case with overflowing output hours
      {.err = SKY_ERROR,
       .out = {0, 10, MaxUint64 - 9},
       .in = {{headTime, 10000000, 10}, {headTime, 10000000, 100}},
       .headTime = headTime,
       .lens = {2, 3}

      }

  };
  int len = (sizeof(cases) / sizeof(tmpstruct));
  GoUint64 err;
  cipher__Address addr;
  makeAddress(&addr);
  for (int i = 0; i < len; i++) {
    tmpstruct tc = cases[i];
    Transaction__Handle tx;
    makeEmptyTransaction(&tx);
    for (int k = 0; k < tc.lens[1]; k++) {
      GoInt64 h = tc.out[k];
      err = SKY_coin_Transaction_PushOutput(tx, &addr, 0, h);
      cr_assert(err == SKY_OK);
    }

    coin__UxArray inUxs;
    makeUxArray(&inUxs, tc.lens[0]);
    coin__UxOut *tmpOut = (coin__UxOut *)inUxs.data;
    for (int j = 0; j < tc.lens[0]; j++) {
      uxInput b = tc.in[j];
      tmpOut->Head.Time = b.times;
      tmpOut->Body.Coins = b.coins;
      tmpOut->Body.Hours = b.hours;
      tmpOut++;
    }
    GoUint64 fee;
    err = SKY_fee_TransactionFee(tx, tc.headTime, &inUxs, &fee);
    cr_assert(err == tc.err);
    if (err != SKY_OK) {
      cr_assert(fee != 0, "Failed %d != %d in Iter %d", fee, tc.fee, i);
    } else {
      cr_assert(fee == tc.fee, "Failed %d != %d in Iter %d", fee, tc.fee, i);
    }
  }
}