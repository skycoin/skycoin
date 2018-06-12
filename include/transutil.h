
#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"
#include "skytypes.h"

int makeKeysAndAddress(cipher__PubKey* ppubkey, cipher__SecKey* pseckey, cipher__Address* paddress);

int makeUxBodyWithSecret(coin__UxBody* puxBody, cipher__SecKey* pseckey);

int makeUxOutWithSecret(coin__UxOut* puxOut, cipher__SecKey* pseckey);

int makeUxBody(coin__UxBody* puxBody);

int makeUxOut(coin__UxOut* puxOut);

int makeAddress(cipher__Address* paddress);

coin__Transaction* makeTransactionFromUxOut(coin__UxOut* puxOut, cipher__SecKey* pseckey, Transaction__Handle* handle);

coin__Transaction* makeTransaction(Transaction__Handle* handle);

coin__Transaction* makeEmptyTransaction(Transaction__Handle* handle);

int makeTransactions(int n, Transactions__Handle* handle);

coin__Transaction* copyTransaction(Transaction__Handle handle, Transaction__Handle* handle2);

void makeRandHash(cipher__SHA256* phash);

int makeUxArray(coin__UxArray* parray, int n);
