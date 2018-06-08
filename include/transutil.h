
#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

int makeKeysAndAddress(cipher__PubKey* ppubkey, cipher__SecKey* pseckey, cipher__Address* paddress);

int makeUxBodyWithSecret(coin__UxBody* puxBody, cipher__SecKey* pseckey);

int makeUxOutWithSecret(coin__UxOut* puxOut, cipher__SecKey* pseckey);

int makeUxBody(coin__UxBody* puxBody);

int makeUxOut(coin__UxOut* puxOut);

int makeAddress(cipher__Address* paddress);

int makeTransactionFromUxOut(coin__UxOut* puxOut, cipher__SecKey* pseckey,
                          coin__Transaction* ptransaction);

int makeTransaction(coin__Transaction* ptransaction);

int makeTransactions(GoSlice* transactions, int n);

void copyTransaction(coin__Transaction* pt1, coin__Transaction* pt2);

void makeRandHash(cipher__SHA256* phash);

int makeUxArray(coin__UxArray* parray, int n);
