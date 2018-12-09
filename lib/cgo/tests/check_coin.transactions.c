
#include <criterion/criterion.h>
#include <criterion/new/assert.h>
#include <signal.h>
#include <stdio.h>
#include <string.h>

#include "libskycoin.h"
#include "skycriterion.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"
#include "skytxn.h"

TestSuite(coin_transaction, .init = setup, .fini = teardown);

GoUint64 Million = 1000000;

Test(coin_transaction, TestTransactionVerify) {
  unsigned long long MaxUint64 = 0xFFFFFFFFFFFFFFFF;
  unsigned int MaxUint16 = 0xFFFF;
  int result;
  coin__Transaction *ptx;
  Transaction__Handle handle;
  // Mismatch header hash
  ptx = makeTransaction(&handle);
  memset(ptx->InnerHash, 0, sizeof(cipher__SHA256));
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert(result == SKY_ERROR);

  // No inputs
  ptx = makeTransaction(&handle);
  result = SKY_coin_Transaction_ResetInputs(handle, 0);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert(result == SKY_ERROR);

  // No outputs
  ptx = makeTransaction(&handle);
  result = SKY_coin_Transaction_ResetOutputs(handle, 0);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert(result == SKY_ERROR);

  // Invalid number of Sigs
  ptx = makeTransaction(&handle);
  result = SKY_coin_Transaction_ResetSignatures(handle, 0);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert(result == SKY_ERROR);
  result = SKY_coin_Transaction_ResetSignatures(handle, 20);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert(result == SKY_ERROR);

  // Too many sigs & inputs
  ptx = makeTransaction(&handle);
  result = SKY_coin_Transaction_ResetSignatures(handle, MaxUint16);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_ResetInputs(handle, MaxUint16);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert(result == SKY_ERROR);

  // Duplicate inputs
  coin__UxOut ux;
  cipher__SecKey seckey;
  cipher__SHA256 sha256;
  makeUxOutWithSecret(&ux, &seckey);
  ptx = makeTransactionFromUxOut(&ux, &seckey, &handle);
  memcpy(&sha256, ptx->In.data, sizeof(cipher__SHA256));
  GoUint16 r;
  result = SKY_coin_Transaction_PushInput(handle, &sha256, &r);
  result = SKY_coin_Transaction_ResetSignatures(handle, 0);
  cr_assert(result == SKY_OK);
  GoSlice seckeys;
  seckeys.data = malloc(sizeof(cipher__SecKey) * 2);
  cr_assert(seckeys.data != NULL);
  registerMemCleanup(seckeys.data);
  seckeys.len = seckeys.cap = 2;
  memcpy(seckeys.data, &seckey, sizeof(cipher__SecKey));
  memcpy(((cipher__SecKey *)seckeys.data) + 1, &seckey, sizeof(cipher__SecKey));
  result = SKY_coin_Transaction_SignInputs(handle, seckeys);
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert(result == SKY_ERROR);

  // Duplicate outputs
  ptx = makeTransaction(&handle);
  coin__TransactionOutput *pOutput = ptx->Out.data;
  cipher__Address addr;
  memcpy(&addr, &pOutput->Address, sizeof(cipher__Address));
  result = SKY_coin_Transaction_PushOutput(handle, &addr, pOutput->Coins,
                                           pOutput->Hours);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert(result == SKY_ERROR);

  // Invalid signature, empty
  ptx = makeTransaction(&handle);
  memset(ptx->Sigs.data, 0, sizeof(cipher__Sig));
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert(result == SKY_ErrInvalidSigPubKeyRecovery);

  // Output coins are 0
  ptx = makeTransaction(&handle);
  pOutput = ptx->Out.data;
  pOutput->Coins = 0;
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert(result == SKY_ERROR);

  // Output coin overflow
  ptx = makeTransaction(&handle);
  pOutput = ptx->Out.data;
  pOutput->Coins = MaxUint64 - 3000000;
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert(result == SKY_ERROR);

  // Output coins are not multiples of 1e6 (valid, decimal restriction is not
  // enforced here)
  ptx = makeTransaction(&handle);
  pOutput = ptx->Out.data;
  pOutput->Coins += 10;
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_ResetSignatures(handle, 0);
  cr_assert(result == SKY_OK);
  cipher__PubKey pubkey;
  result = SKY_cipher_GenerateKeyPair(&pubkey, &seckey);
  cr_assert(result == SKY_OK);
  seckeys.data = &seckey;
  seckeys.len = 1;
  seckeys.cap = 1;
  result = SKY_coin_Transaction_SignInputs(handle, seckeys);
  cr_assert(result == SKY_OK);
  cr_assert(pOutput->Coins % 1000000 != 0);
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert(result == SKY_OK);

  // Valid
  ptx = makeTransaction(&handle);
  pOutput = ptx->Out.data;
  pOutput->Coins = 10000000;
  pOutput++;
  pOutput->Coins = 1000000;
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_Verify(handle);
  cr_assert(result == SKY_OK);
}

Test(coin_transaction, TestTransactionPushInput, SKY_ABORT) {
  unsigned long long MaxUint64 = 0xFFFFFFFFFFFFFFFF;
  unsigned int MaxUint16 = 0xFFFF;
  int result;
  Transaction__Handle handle;
  coin__Transaction *ptx;
  coin__UxOut ux;
  ptx = makeEmptyTransaction(&handle);
  makeUxOut(&ux);
  cipher__SHA256 hash;
  result = SKY_coin_UxOut_Hash(&ux, &hash);
  cr_assert(result == SKY_OK);
  GoUint16 r;
  result = SKY_coin_Transaction_PushInput(handle, &hash, &r);
  cr_assert(result == SKY_OK);
  cr_assert(r == 0);
  cr_assert(ptx->In.len == 1);
  cipher__SHA256 *pIn = ptx->In.data;
  cr_assert(eq(u8[sizeof(cipher__SHA256)], hash, *pIn));

  int len = ptx->In.len;
  void *data = malloc(len * sizeof(cipher__SHA256));
  cr_assert(data != NULL);
  registerMemCleanup(data);
  memcpy(data, ptx->In.data, len * sizeof(cipher__SHA256));
  result = SKY_coin_Transaction_ResetInputs(handle, MaxUint16 + len);
  cr_assert(result == SKY_OK);
  memcpy(ptx->In.data, data, len * sizeof(cipher__Sig));
  freeRegisteredMemCleanup(data);
  makeUxOut(&ux);
  result = SKY_coin_UxOut_Hash(&ux, &hash);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_PushInput(handle, &hash, &r);
  cr_assert(result == SKY_ERROR);
}

Test(coin_transaction, TestTransactionPushOutput) {
  int result;
  Transaction__Handle handle;
  coin__Transaction *ptx;
  ptx = makeEmptyTransaction(&handle);

  cipher__Address addr;
  makeAddress(&addr);
  result = SKY_coin_Transaction_PushOutput(handle, &addr, 100, 150);
  cr_assert(result == SKY_OK);
  cr_assert(ptx->Out.len == 1);
  coin__TransactionOutput *pOutput = ptx->Out.data;
  coin__TransactionOutput output;
  memcpy(&output.Address, &addr, sizeof(cipher__Address));
  output.Coins = 100;
  output.Hours = 150;
  cr_assert(eq(type(coin__TransactionOutput), output, *pOutput));
  for (int i = 1; i < 20; i++) {
    makeAddress(&addr);
    result = SKY_coin_Transaction_PushOutput(handle, &addr, i * 100, i * 50);
    cr_assert(result == SKY_OK);
    cr_assert(ptx->Out.len == i + 1);
    pOutput = ptx->Out.data;
    pOutput += i;
    memcpy(&output.Address, &addr, sizeof(cipher__Address));
    output.Coins = i * 100;
    output.Hours = i * 50;
    cr_assert(eq(type(coin__TransactionOutput), output, *pOutput));
  }
}

Test(coin_transaction, TestTransactionHash) {
  int result;
  Transaction__Handle handle;
  coin__Transaction *ptx;
  ptx = makeEmptyTransaction(&handle);

  cipher__SHA256 nullHash, hash1, hash2;
  memset(&nullHash, 0, sizeof(cipher__SHA256));
  result = SKY_coin_Transaction_Hash(handle, &hash1);
  cr_assert(result == SKY_OK);
  cr_assert(not(eq(u8[sizeof(cipher__SHA256)], nullHash, hash1)));
  result = SKY_coin_Transaction_HashInner(handle, &hash2);
  cr_assert(result == SKY_OK);
  cr_assert(not(eq(u8[sizeof(cipher__SHA256)], hash2, hash1)));
}

Test(coin_transaction, TestTransactionUpdateHeader) {
  int result;
  Transaction__Handle handle;
  coin__Transaction *ptx;
  ptx = makeTransaction(&handle);
  cipher__SHA256 hash, nullHash, hashInner;
  memcpy(&hash, &ptx->InnerHash, sizeof(cipher__SHA256));
  memset(&ptx->InnerHash, 0, sizeof(cipher__SHA256));
  memset(&nullHash, 0, sizeof(cipher__SHA256));
  result = SKY_coin_Transaction_UpdateHeader(handle);
  cr_assert(not(eq(u8[sizeof(cipher__SHA256)], ptx->InnerHash, nullHash)));
  cr_assert(eq(u8[sizeof(cipher__SHA256)], hash, ptx->InnerHash));
  result = SKY_coin_Transaction_HashInner(handle, &hashInner);
  cr_assert(result == SKY_OK);
  cr_assert(eq(u8[sizeof(cipher__SHA256)], hashInner, ptx->InnerHash));
}

Test(coin_transaction, TestTransactionsSize) {
  int result;
  Transactions__Handle txns;
  result = makeTransactions(10, &txns);
  cr_assert(result == SKY_OK);
  GoInt size = 0;
  for (size_t i = 0; i < 10; i++) {
    Transaction__Handle handle;
    result = SKY_coin_Transactions_GetAt(txns, i, &handle);
    registerHandleClose(handle);
    cr_assert(result == SKY_OK);
    GoSlice p1 = {NULL, 0, 0};
    result = SKY_coin_Transaction_Serialize(handle, (GoSlice_ *)&p1);
    cr_assert(result == SKY_OK, "SKY_coin_Transaction_Serialize");
    size += p1.len;
    cr_assert(result == SKY_OK, "SKY_coin_Transaction_Size");
  }
  GoUint32 sizeTransactions;
  result = SKY_coin_Transactions_Size(txns, &sizeTransactions);
  cr_assert(size != 0);
  cr_assert(sizeTransactions == size);
}

Test(coin_transactions, TestTransactionVerifyInput, SKY_ABORT) {
  int result;
  Transaction__Handle handle;
  coin__Transaction *ptx;
  ptx = makeTransaction(&handle);
  result = SKY_coin_Transaction_VerifyInput(handle, NULL);
  cr_assert(result == SKY_ERROR);
  coin__UxArray ux;
  memset(&ux, 0, sizeof(coin__UxArray));
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert(result == SKY_ERROR);
  memset(&ux, 0, sizeof(coin__UxArray));
  ux.data = malloc(3 * sizeof(coin__UxOut));
  cr_assert(ux.data != NULL);
  registerMemCleanup(ux.data);
  ux.len = 3;
  ux.cap = 3;
  memset(ux.data, 0, 3 * sizeof(coin__UxOut));
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert(result == SKY_ERROR);

  coin__UxOut uxOut;
  cipher__SecKey seckey;
  cipher__Sig sig;
  cipher__SHA256 hash;

  result = makeUxOutWithSecret(&uxOut, &seckey);
  cr_assert(result == SKY_OK);
  ptx = makeTransactionFromUxOut(&uxOut, &seckey, &handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_ResetSignatures(handle, 0);
  cr_assert(result == SKY_OK);
  ux.data = &uxOut;
  ux.len = 1;
  ux.cap = 1;
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert(result == SKY_ERROR);

  memset(&sig, 0, sizeof(cipher__Sig));
  result = makeUxOutWithSecret(&uxOut, &seckey);
  cr_assert(result == SKY_OK);
  ptx = makeTransactionFromUxOut(&uxOut, &seckey, &handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_ResetSignatures(handle, 1);
  cr_assert(result == SKY_OK);
  memcpy(ptx->Sigs.data, &sig, sizeof(cipher__Sig));
  ux.data = &uxOut;
  ux.len = 1;
  ux.cap = 1;
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert(result == SKY_ERROR);

  // Invalid Tx Inner Hash
  result = makeUxOutWithSecret(&uxOut, &seckey);
  cr_assert(result == SKY_OK);
  ptx = makeTransactionFromUxOut(&uxOut, &seckey, &handle);
  cr_assert(result == SKY_OK);
  memset(ptx->InnerHash, 0, sizeof(cipher__SHA256));
  ux.data = &uxOut;
  ux.len = 1;
  ux.cap = 1;
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert(result == SKY_ERROR);

  // Ux hash mismatch
  result = makeUxOutWithSecret(&uxOut, &seckey);
  cr_assert(result == SKY_OK);
  ptx = makeTransactionFromUxOut(&uxOut, &seckey, &handle);
  cr_assert(result == SKY_OK);
  memset(&uxOut, 0, sizeof(coin__UxOut));
  ux.data = &uxOut;
  ux.len = 1;
  ux.cap = 1;
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert(result == SKY_ERROR);

  // Invalid signature
  result = makeUxOutWithSecret(&uxOut, &seckey);
  cr_assert(result == SKY_OK);
  ptx = makeTransactionFromUxOut(&uxOut, &seckey, &handle);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_ResetSignatures(handle, 1);
  cr_assert(result == SKY_OK);
  memset(ptx->Sigs.data, 0, sizeof(cipher__Sig));
  ux.data = &uxOut;
  ux.len = 1;
  ux.cap = 1;
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert(result == SKY_ERROR);

  // Valid
  result = makeUxOutWithSecret(&uxOut, &seckey);
  cr_assert(result == SKY_OK);
  ptx = makeTransactionFromUxOut(&uxOut, &seckey, &handle);
  cr_assert(result == SKY_OK);
  ux.data = &uxOut;
  ux.len = 1;
  ux.cap = 1;
  result = SKY_coin_Transaction_VerifyInput(handle, &ux);
  cr_assert(result == SKY_OK);
}

Test(coin_transactions, TestTransactionSignInputs, SKY_ABORT) {
  int result;
  coin__Transaction *ptx;
  Transaction__Handle handle;
  coin__UxOut ux, ux2;
  cipher__SecKey seckey, seckey2;
  cipher__SHA256 hash, hash2;
  cipher__Address addr, addr2;
  cipher__PubKey pubkey;
  GoUint16 r;
  GoSlice keys;

  // Error if txns already signed
  ptx = makeEmptyTransaction(&handle);
  result = SKY_coin_Transaction_ResetSignatures(handle, 1);
  cr_assert(result == SKY_OK);

  memset(&seckey, 0, sizeof(cipher__SecKey));
  keys.data = &seckey;
  keys.len = 1;
  keys.cap = 1;
  result = SKY_coin_Transaction_SignInputs(handle, keys);
  cr_assert(result == SKY_ERROR);

  // Panics if not enough keys
  ptx = makeEmptyTransaction(&handle);
  memset(&seckey, 0, sizeof(cipher__SecKey));
  memset(&seckey2, 0, sizeof(cipher__SecKey));
  result = makeUxOutWithSecret(&ux, &seckey);
  cr_assert(result == SKY_OK);
  result = SKY_coin_UxOut_Hash(&ux, &hash);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_PushInput(handle, &hash, &r);
  cr_assert(result == SKY_OK);
  result = makeUxOutWithSecret(&ux2, &seckey2);
  cr_assert(result == SKY_OK);
  result = SKY_coin_UxOut_Hash(&ux2, &hash2);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_PushInput(handle, &hash2, &r);
  cr_assert(result == SKY_OK);
  makeAddress(&addr);
  result = SKY_coin_Transaction_PushOutput(handle, &addr, 40, 80);
  cr_assert(result == SKY_OK);
  cr_assert(ptx->Sigs.len == 0);
  keys.data = &seckey;
  keys.len = 1;
  keys.cap = 1;
  result = SKY_coin_Transaction_SignInputs(handle, keys);
  cr_assert(result == SKY_ERROR);
  cr_assert(ptx->Sigs.len == 0);

  // Valid signing
  result = SKY_coin_Transaction_HashInner(handle, &hash);
  cr_assert(result == SKY_OK);
  keys.data = malloc(2 * sizeof(cipher__SecKey));
  cr_assert(keys.data != NULL);
  registerMemCleanup(keys.data);
  keys.len = keys.cap = 2;
  memcpy(keys.data, &seckey, sizeof(cipher__SecKey));
  memcpy(((cipher__SecKey *)keys.data) + 1, &seckey2, sizeof(cipher__SecKey));
  result = SKY_coin_Transaction_SignInputs(handle, keys);
  cr_assert(result == SKY_OK);
  cr_assert(ptx->Sigs.len == 2);
  result = SKY_coin_Transaction_HashInner(handle, &hash2);
  cr_assert(result == SKY_OK);
  cr_assert(eq(u8[sizeof(cipher__SHA256)], hash, hash2));

  result = SKY_cipher_PubKeyFromSecKey(&seckey, &pubkey);
  cr_assert(result == SKY_OK);
  result = SKY_cipher_AddressFromPubKey(&pubkey, &addr);
  cr_assert(result == SKY_OK);
  result = SKY_cipher_PubKeyFromSecKey(&seckey2, &pubkey);
  cr_assert(result == SKY_OK);
  result = SKY_cipher_AddressFromPubKey(&pubkey, &addr2);
  cr_assert(result == SKY_OK);

  cipher__SHA256 addHash, addHash2;
  result = SKY_cipher_AddSHA256(&hash, (cipher__SHA256 *)ptx->In.data, &addHash);
  cr_assert(result == SKY_OK);
  result = SKY_cipher_AddSHA256(&hash, ((cipher__SHA256 *)ptx->In.data) + 1, &addHash2);
  cr_assert(result == SKY_OK);
  result = SKY_cipher_VerifyAddressSignedHash(&addr, (cipher__Sig *)ptx->Sigs.data, &addHash);
  cr_assert(result == SKY_OK);
  result = SKY_cipher_VerifyAddressSignedHash(&addr2, ((cipher__Sig *)ptx->Sigs.data) + 1, &addHash2);
  cr_assert(result == SKY_OK);
  result = SKY_cipher_VerifyAddressSignedHash(&addr, ((cipher__Sig *)ptx->Sigs.data) + 1, &hash);
  cr_assert(result == SKY_ERROR);
  result = SKY_cipher_VerifyAddressSignedHash(&addr2, (cipher__Sig *)ptx->Sigs.data, &hash);
  cr_assert(result == SKY_ERROR);
}

Test(coin_transactions, TestTransactionHashInner) {
  int result;
  Transaction__Handle handle1 = 0, handle2 = 0;
  coin__Transaction *ptx = NULL;
  coin__Transaction *ptx2 = NULL;
  ptx = makeTransaction(&handle1);
  cipher__SHA256 hash, nullHash;
  result = SKY_coin_Transaction_HashInner(handle1, &hash);
  cr_assert(result == SKY_OK);
  memset(&nullHash, 0, sizeof(cipher__SHA256));
  cr_assert(not(eq(u8[sizeof(cipher__SHA256)], nullHash, hash)));

  // If tx.In is changed, hash should change
  ptx2 = copyTransaction(handle1, &handle2);
  cr_assert(eq(type(coin__Transaction), *ptx, *ptx2));
  cr_assert(ptx != ptx2);
  cr_assert(ptx2->In.len > 0);
  coin__UxOut uxOut;
  makeUxOut(&uxOut);
  cipher__SHA256 *phash = ptx2->In.data;
  result = SKY_coin_UxOut_Hash(&uxOut, phash);
  cr_assert(result == SKY_OK);
  cr_assert(not(eq(type(coin__Transaction), *ptx, *ptx2)));
  cipher__SHA256 hash1, hash2;
  result = SKY_coin_Transaction_HashInner(handle1, &hash1);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_HashInner(handle2, &hash2);
  cr_assert(result == SKY_OK);
  cr_assert(not(eq(u8[sizeof(cipher__SHA256)], hash1, hash2)));

  // If tx.Out is changed, hash should change
  handle2 = 0;
  ptx2 = copyTransaction(handle1, &handle2);
  cr_assert(ptx != ptx2);
  cr_assert(eq(type(coin__Transaction), *ptx, *ptx2));
  coin__TransactionOutput *output = ptx2->Out.data;
  cipher__Address addr;
  makeAddress(&addr);
  memcpy(&output->Address, &addr, sizeof(cipher__Address));
  cr_assert(not(eq(type(coin__Transaction), *ptx, *ptx2)));
  cr_assert(eq(type(cipher__Address), addr, output->Address));
  result = SKY_coin_Transaction_HashInner(handle1, &hash1);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_HashInner(handle2, &hash2);
  cr_assert(result == SKY_OK);
  cr_assert(not(eq(u8[sizeof(cipher__SHA256)], hash1, hash2)));

  // If tx.Head is changed, hash should not change
  ptx2 = copyTransaction(handle1, &handle2);
  int len = ptx2->Sigs.len;
  cipher__Sig *newSigs = malloc((len + 1) * sizeof(cipher__Sig));
  cr_assert(newSigs != NULL);
  registerMemCleanup(newSigs);
  memcpy(newSigs, ptx2->Sigs.data, len * sizeof(cipher__Sig));
  result = SKY_coin_Transaction_ResetSignatures(handle2, len + 1);
  cr_assert(result == SKY_OK);
  memcpy(ptx2->Sigs.data, newSigs, len * sizeof(cipher__Sig));
  newSigs += len;
  memset(newSigs, 0, sizeof(cipher__Sig));
  result = SKY_coin_Transaction_HashInner(handle1, &hash1);
  cr_assert(result == SKY_OK);
  result = SKY_coin_Transaction_HashInner(handle2, &hash2);
  cr_assert(result == SKY_OK);
  cr_assert(eq(u8[sizeof(cipher__SHA256)], hash1, hash2));
}

Test(coin_transactions, TestTransactionSerialization) {
  int result;
  coin__Transaction *ptx;
  Transaction__Handle handle;
  ptx = makeTransaction(&handle);
  GoSlice_ data;
  memset(&data, 0, sizeof(GoSlice_));
  result = SKY_coin_Transaction_Serialize(handle, &data);
  cr_assert(result == SKY_OK);
  registerMemCleanup(data.data);
  coin__Transaction *ptx2;
  Transaction__Handle handle2;
  GoSlice d = {data.data, data.len, data.cap};
  result = SKY_coin_TransactionDeserialize(d, &handle2);
  cr_assert(result == SKY_OK);
  result = SKY_coin_GetTransactionObject(handle2, &ptx2);
  cr_assert(result == SKY_OK);
  cr_assert(eq(type(coin__Transaction), *ptx, *ptx2));
}

Test(coin_transactions, TestTransactionOutputHours) {
  coin__Transaction *ptx;
  Transaction__Handle handle;
  ptx = makeEmptyTransaction(&handle);
  cipher__Address addr;
  makeAddress(&addr);
  int result;
  result = SKY_coin_Transaction_PushOutput(handle, &addr, 1000000, 100);
  cr_assert(result == SKY_OK);
  makeAddress(&addr);
  result = SKY_coin_Transaction_PushOutput(handle, &addr, 1000000, 200);
  cr_assert(result == SKY_OK);
  makeAddress(&addr);
  result = SKY_coin_Transaction_PushOutput(handle, &addr, 1000000, 500);
  cr_assert(result == SKY_OK);
  makeAddress(&addr);
  result = SKY_coin_Transaction_PushOutput(handle, &addr, 1000000, 0);
  cr_assert(result == SKY_OK);
  GoUint64 hours;
  result = SKY_coin_Transaction_OutputHours(handle, &hours);
  cr_assert(result == SKY_OK);
  cr_assert(hours == 800);
  makeAddress(&addr);
  result = SKY_coin_Transaction_PushOutput(handle, &addr, 1000000,
                                           0xFFFFFFFFFFFFFFFF - 700);
  result = SKY_coin_Transaction_OutputHours(handle, &hours);
  cr_assert(result == SKY_ERROR);
}

Test(coin_transactions, TestTransactionsHashes) {
  int result;
  GoSlice_ hashes = {NULL, 0, 0};
  Transactions__Handle hTxns;
  result = makeTransactions(4, &hTxns);
  cr_assert(result == SKY_OK);

  result = SKY_coin_Transactions_Hashes(hTxns, &hashes);
  cr_assert(result == SKY_OK, "SKY_coin_Transactions_Hashes failed");
  registerMemCleanup(hashes.data);
  cr_assert(hashes.len == 4);
  cipher__SHA256 *ph = hashes.data;
  cipher__SHA256 hash;
  for (int i = 0; i < 4; i++) {
    Transaction__Handle handle;
    result = SKY_coin_Transactions_GetAt(hTxns, i, &handle);
    cr_assert(result == SKY_OK);
    result = SKY_coin_Transaction_Hash(handle, &hash);
    cr_assert(result == SKY_OK, "SKY_coin_Transaction_Hash failed");
    cr_assert(eq(u8[sizeof(cipher__SHA256)], *ph, hash));
    ph++;
  }
}

Test(coin_transactions, TestTransactionsTruncateBytesTo) {
  int result;
  Transactions__Handle h1, h2;
  result = makeTransactions(10, &h1);
  cr_assert(result == SKY_OK);
  GoInt length;
  result = SKY_coin_Transactions_Length(h1, &length);
  cr_assert(result == SKY_OK);
  int trunc = 0;
  GoUint32 size;
  for (int i = 0; i < length / 2; i++)
  {
    Transaction__Handle handle;
    result = SKY_coin_Transactions_GetAt(h1, i, &handle);
    registerHandleClose(handle);
    cr_assert(result == SKY_OK);
    result = SKY_coin_Transaction_Size(handle, &size);
    trunc += size;
    cr_assert(result == SKY_OK, "SKY_coin_Transaction_Size failed");
  }
  result = SKY_coin_Transactions_TruncateBytesTo(h1, trunc, &h2);
  cr_assert(result == SKY_OK, "SKY_coin_Transactions_TruncateBytesTo failed");
  registerHandleClose(h2);

  GoInt length2;
  result = SKY_coin_Transactions_Length(h2, &length2);
  cr_assert(result == SKY_OK);
  cr_assert(length2 == length / 2);
  result = SKY_coin_Transactions_Size(h2, &size);
  cr_assert(result == SKY_OK, "SKY_coin_Transactions_Size failed");
  cr_assert(trunc == size);

  trunc++;
  result = SKY_coin_Transactions_TruncateBytesTo(h1, trunc, &h2);
  cr_assert(result == SKY_OK, "SKY_coin_Transactions_TruncateBytesTo failed");
  registerHandleClose(h2);

  // Stepping into next boundary has same cutoff, must exceed
  result = SKY_coin_Transactions_Length(h2, &length2);
  cr_assert(result == SKY_OK);
  cr_assert(length2 == length / 2);
  result = SKY_coin_Transactions_Size(h2, &size);
  cr_assert(result == SKY_OK, "SKY_coin_Transactions_Size failed");
  cr_assert(trunc - 1 == size);
}

typedef struct {
  GoUint64 coins;
  GoUint64 hours;
} test_ux;

typedef struct {
  test_ux *inUxs;
  test_ux *outUxs;
  int sizeIn;
  int sizeOut;
  GoUint64 headTime;
  int failure;
} test_case;

int makeTestCaseArrays(test_ux *elems, int size, coin__UxArray *pArray) {
  if (size <= 0) {
    pArray->len = 0;
    pArray->cap = 0;
    pArray->data = NULL;
    return SKY_OK;
  }
  int elems_size = sizeof(coin__UxOut);
  void *data;
  data = malloc(size * elems_size);
  if (data == NULL)
    return SKY_ERROR;
  registerMemCleanup(data);
  memset(data, 0, size * elems_size);
  pArray->data = data;
  pArray->len = size;
  pArray->cap = size;
  coin__UxOut *p = data;
  for (int i = 0; i < size; i++) {
    p->Body.Coins = elems[i].coins;
    p->Body.Hours = elems[i].hours;
    p++;
  }
  return SKY_OK;
}

Test(coin_transactions, TestVerifyTransactionCoinsSpending) {
  unsigned long long MaxUint64 = 0xFFFFFFFFFFFFFFFF;
  unsigned int MaxUint16 = 0xFFFF;
  // Input coins overflow
  test_ux in1[] = {{MaxUint64 - Million + 1, 10}, {Million, 0}};

  // Output coins overflow
  test_ux in2[] = {{10 * Million, 10}};
  test_ux out2[] = {{MaxUint64 - 10 * Million + 1, 0}, {20 * Million, 1}};

  // Insufficient coins
  test_ux in3[] = {{10 * Million, 10}, {15 * Million, 10}};
  test_ux out3[] = {{20 * Million, 1}, {10 * Million, 1}};

  // Destroyed coins
  test_ux in4[] = {{10 * Million, 10}, {15 * Million, 10}};
  test_ux out4[] = {{5 * Million, 1}, {10 * Million, 1}};

  // Valid
  test_ux in5[] = {{10 * Million, 10}, {15 * Million, 10}};
  test_ux out5[] = {{10 * Million, 11}, {10 * Million, 1}, {5 * Million, 0}};

  test_case tests[] = {
      {in1, NULL, 2, 0, 0, 1}, // Input coins overflow
      {in2, out2, 1, 2, 0, 1}, // Output coins overflow
      {in3, out3, 2, 2, 0, 1}, // Destroyed coins
      {in4, out4, 1, 1, Million,
       1}, // Invalid (coin hours overflow when adding earned hours, which is
           // treated as 0, and now enough coin hours)
      {in5, out5, 2, 3, 0, 0} // Valid
  };

  coin__UxArray inArray;
  coin__UxArray outArray;
  int result;
  int count = sizeof(tests) / sizeof(tests[0]);
  for (int i = 0; i < count; i++) {
    result = makeTestCaseArrays(tests[i].inUxs, tests[i].sizeIn, &inArray);
    cr_assert(result == SKY_OK);
    result = makeTestCaseArrays(tests[i].outUxs, tests[i].sizeOut, &outArray);
    cr_assert(result == SKY_OK);
    result = SKY_coin_VerifyTransactionCoinsSpending(&inArray, &outArray);
    if (tests[i].failure)
      cr_assert(result == SKY_ERROR, "VerifyTransactionCoinsSpending succeeded %d", i + 1);
    else
      cr_assert(result == SKY_OK, "VerifyTransactionCoinsSpending failed %d", i + 1);
  }
}

Test(coin_transactions, TestVerifyTransactionHoursSpending) {

  GoUint64 Million = 1000000;
  unsigned long long MaxUint64 = 0xFFFFFFFFFFFFFFFF;
  unsigned int MaxUint16 = 0xFFFF;
  // Input hours overflow
  test_ux in1[] = {{3 * Million, MaxUint64 - Million + 1}, {Million, Million}};

  // Insufficient coin hours
  test_ux in2[] = {{10 * Million, 10}, {15 * Million, 10}};

  test_ux out2[] = {{15 * Million, 10}, {10 * Million, 11}};

  // coin hours time calculation overflow
  test_ux in3[] = {{10 * Million, 10}, {15 * Million, 10}};

  test_ux out3[] = {{10 * Million, 11}, {10 * Million, 1}, {5 * Million, 0}};

  // Invalid (coin hours overflow when adding earned hours, which is treated as
  // 0, and now enough coin hours)
  test_ux in4[] = {{10 * Million, MaxUint64}};

  test_ux out4[] = {{10 * Million, 1}};

  // Valid (coin hours overflow when adding earned hours, which is treated as 0,
  // but not sending any hours)
  test_ux in5[] = {{10 * Million, MaxUint64}};

  test_ux out5[] = {{10 * Million, 0}};

  // Valid (base inputs have insufficient coin hours, but have sufficient after
  // adjusting coinhours by headTime)
  test_ux in6[] = {{10 * Million, 10}, {15 * Million, 10}};

  test_ux out6[] = {{15 * Million, 10}, {10 * Million, 11}};

  // valid
  test_ux in7[] = {{10 * Million, 10}, {15 * Million, 10}};

  test_ux out7[] = {{10 * Million, 11}, {10 * Million, 1}, {5 * Million, 0}};

  test_case tests[] = {
      {in1, NULL, 2, 0, 0, 1},         // Input hours overflow
      {in2, out2, 2, 2, 0, 1},         // Insufficient coin hours
      {in3, out3, 2, 3, MaxUint64, 1}, // coin hours time calculation overflow
      {in4, out4, 1, 1, Million,
       1}, // Invalid (coin hours overflow when adding earned hours, which is
           // treated as 0, and now enough coin hours)
      {in5, out5, 1, 1, 0,
       0}, // Valid (coin hours overflow when adding earned hours, which is
           // treated as 0, but not sending any hours)
      {in6, out6, 2, 2, 1492707255,
       0}, // Valid (base inputs have insufficient coin hours, but have
           // sufficient after adjusting coinhours by headTime)
      {in7, out7, 2, 3, 0, 0}, // Valid
  };
  coin__UxArray inArray;
  coin__UxArray outArray;
  int result;
  int count = sizeof(tests) / sizeof(tests[0]);
  for (int i = 0; i < count; i++) {
    result = makeTestCaseArrays(tests[i].inUxs, tests[i].sizeIn, &inArray);
    cr_assert(result == SKY_OK);
    result = makeTestCaseArrays(tests[i].outUxs, tests[i].sizeOut, &outArray);
    cr_assert(result == SKY_OK);
    result = SKY_coin_VerifyTransactionHoursSpending(tests[i].headTime,
                                                     &inArray, &outArray);
    if (tests[i].failure)
      cr_assert(result == SKY_ERROR,
                "SKY_coin_VerifyTransactionHoursSpending succeeded %d", i + 1);
    else
      cr_assert(result == SKY_OK,
                "SKY_coin_VerifyTransactionHoursSpending failed %d", i + 1);
  }
}

GoUint32_ fix1FeeCalculator(Transaction__Handle handle, GoUint64_ *pFee, void *context) {
  *pFee = 1;
  return SKY_OK;
}

GoUint32_ badFeeCalculator(Transaction__Handle handle, GoUint64_ *pFee, void *context) {
  return SKY_ERROR;
}

GoUint32_ overflowFeeCalculator(Transaction__Handle handle, GoUint64_ *pFee, void *context) {
  *pFee = 0xFFFFFFFFFFFFFFFF;
  return SKY_OK;
}

Test(coin_transactions, TestTransactionsFees) {
  GoUint64 fee;
  int result;
  Transactions__Handle transactionsHandle = 0;
  Transaction__Handle transactionHandle = 0;

  // Nil txns
  makeTransactions(0, &transactionsHandle);
  FeeCalculator f1 = {fix1FeeCalculator, NULL};
  result = SKY_coin_Transactions_Fees(transactionsHandle, &f1, &fee);
  cr_assert(result == SKY_OK);
  cr_assert(fee == 0);

  makeEmptyTransaction(&transactionHandle);
  result = SKY_coin_Transactions_Add(transactionsHandle, transactionHandle);
  cr_assert(result == SKY_OK);
  makeEmptyTransaction(&transactionHandle);
  result = SKY_coin_Transactions_Add(transactionsHandle, transactionHandle);
  cr_assert(result == SKY_OK);
  // 2 transactions, calc() always returns 1
  result = SKY_coin_Transactions_Fees(transactionsHandle, &f1, &fee);
  cr_assert(result == SKY_OK);
  cr_assert(fee == 2);

  // calc error
  FeeCalculator badFee = {badFeeCalculator, NULL};
  result = SKY_coin_Transactions_Fees(transactionsHandle, &badFee, &fee);
  cr_assert(result == SKY_ERROR);

  // summing of calculated fees overflows
  FeeCalculator overflow = {overflowFeeCalculator, NULL};
  result = SKY_coin_Transactions_Fees(transactionsHandle, &overflow, &fee);
  cr_assert(result == SKY_ERROR);
}

GoUint32_ feeCalculator1(Transaction__Handle handle, GoUint64_ *pFee, void *context) {
  coin__Transaction *pTx;
  int result = SKY_coin_GetTransactionObject(handle, &pTx);
  if (result == SKY_OK) {
    coin__TransactionOutput *pOutput = pTx->Out.data;
    *pFee = 100 * Million - pOutput->Hours;
  }
  return result;
}

GoUint32_ feeCalculator2(Transaction__Handle handle, GoUint64_ *pFee, void *context) {
  *pFee = 100 * Million;
  return SKY_OK;
}

void assertTransactionsHandleEqual(Transaction__Handle h1,
                                   Transaction__Handle h2, char *testName) {
  coin__Transaction *pTx1;
  coin__Transaction *pTx2;
  int result;
  result = SKY_coin_GetTransactionObject(h1, &pTx1);
  cr_assert(result == SKY_OK);
  result = SKY_coin_GetTransactionObject(h2, &pTx2);
  cr_assert(result == SKY_OK);
  cr_assert(eq(type(coin__Transaction), *pTx1, *pTx2), "Failed SortTransactions test \"%s\"", testName);
}

void testTransactionSorting(Transactions__Handle hTrans, int *original_indexes,
                            int original_indexes_count, int *expected_indexes,
                            int expected_indexes_count, FeeCalculator *feeCalc,
                            char *testName) {

  int result;
  Transactions__Handle transactionsHandle, sortedTxnsHandle;
  Transaction__Handle handle;
  makeTransactions(0, &transactionsHandle);
  for (int i = 0; i < original_indexes_count; i++) {
    result = SKY_coin_Transactions_GetAt(hTrans, original_indexes[i], &handle);
    cr_assert(result == SKY_OK);
    registerHandleClose(handle);
    result = SKY_coin_Transactions_Add(transactionsHandle, handle);
    cr_assert(result == SKY_OK);
  }
  result = SKY_coin_SortTransactions(transactionsHandle, feeCalc, &sortedTxnsHandle);
  cr_assert(result == SKY_OK, "SKY_coin_SortTransactions");
  registerHandleClose(sortedTxnsHandle);
  Transaction__Handle h1, h2;
  for (int i = 0; i < expected_indexes_count; i++) {
    int expected_index = expected_indexes[i];
    result = SKY_coin_Transactions_GetAt(sortedTxnsHandle, i, &h1);
    cr_assert(result == SKY_OK);
    registerHandleClose(h1);
    result = SKY_coin_Transactions_GetAt(hTrans, expected_index, &h2);
    cr_assert(result == SKY_OK);
    registerHandleClose(h2);
    assertTransactionsHandleEqual(h1, h2, testName);
  }
}

GoUint32_ feeCalculator3(Transaction__Handle handle, GoUint64_ *pFee, void *context) {
  cipher__SHA256 *thirdHash = (cipher__SHA256 *)context;
  cipher__SHA256 hash;
  unsigned long long MaxUint64 = 0xFFFFFFFFFFFFFFFF;
  unsigned int MaxUint16 = 0xFFFF;
  int result = SKY_coin_Transaction_Hash(handle, &hash);
  if (result == SKY_OK &&
      (memcmp(&hash, thirdHash, sizeof(cipher__SHA256)) == 0)) {
    *pFee = MaxUint64 / 2;
  } else {
    coin__Transaction *pTx;
    result = SKY_coin_GetTransactionObject(handle, &pTx);
    if (result == SKY_OK) {
      coin__TransactionOutput *pOutput = pTx->Out.data;
      *pFee = 100 * Million - pOutput->Hours;
    }
  }
  return result;
}

GoUint32_ feeCalculator4(Transaction__Handle handle, GoUint64_ *pFee, void *context) {
  cipher__SHA256 hash;
  cipher__SHA256 *thirdHash = (cipher__SHA256 *)context;

  int result = SKY_coin_Transaction_Hash(handle, &hash);
  if (result == SKY_OK &&
      (memcmp(&hash, thirdHash, sizeof(cipher__SHA256)) == 0)) {
    *pFee = 0;
    result = SKY_ERROR;
  } else {
    coin__Transaction *pTx;
    result = SKY_coin_GetTransactionObject(handle, &pTx);
    if (result == SKY_OK) {
      coin__TransactionOutput *pOutput = pTx->Out.data;
      *pFee = 100 * Million - pOutput->Hours;
    }
  }
  return result;
}

Test(coin_transactions, TestSortTransactions) {
  int n = 6;
  int i;
  int result;

  Transactions__Handle transactionsHandle = 0;
  Transactions__Handle transactionsHandle2 = 0;
  Transactions__Handle hashSortedTxnsHandle = 0;
  Transactions__Handle sortedTxnsHandle = 0;
  Transaction__Handle transactionHandle = 0;
  cipher__Address addr;
  makeTransactions(0, &transactionsHandle);
  cipher__SHA256 thirdHash;
  for (i = 0; i < 6; i++) {
    makeEmptyTransaction(&transactionHandle);
    makeAddress(&addr);
    result = SKY_coin_Transaction_PushOutput(transactionHandle, &addr, 1000000, i * 1000);
    cr_assert(result == SKY_OK);
    result = SKY_coin_Transaction_UpdateHeader(transactionHandle);
    cr_assert(result == SKY_OK);
    result = SKY_coin_Transactions_Add(transactionsHandle, transactionHandle);
    cr_assert(result == SKY_OK);
    if (i == 2) {
      result = SKY_coin_Transaction_Hash(transactionHandle, &thirdHash);
      cr_assert(result == SKY_OK);
    }
  }
  sortTransactions(transactionsHandle, &hashSortedTxnsHandle);

  int index1[] = {0, 1};
  int expec1[] = {0, 1};
  FeeCalculator fc1 = {feeCalculator1, NULL};
  testTransactionSorting(transactionsHandle, index1, 2, expec1, 2, &fc1, "Already sorted");

  int index2[] = {1, 0};
  int expec2[] = {0, 1};
  testTransactionSorting(transactionsHandle, index2, 2, expec2, 2, &fc1, "reverse sorted");

  FeeCalculator fc2 = {feeCalculator2, NULL};
  testTransactionSorting(hashSortedTxnsHandle, index2, 2, expec2, 2, &fc2, "hash tiebreaker");

  int index3[] = {1, 2, 0};
  int expec3[] = {2, 0, 1};
  FeeCalculator f3 = {feeCalculator3, &thirdHash};
  testTransactionSorting(transactionsHandle, index3, 3, expec3, 3, &f3, "invalid fee multiplication is capped");

  int index4[] = {1, 2, 0};
  int expec4[] = {0, 1};
  FeeCalculator f4 = {feeCalculator4, &thirdHash};
  testTransactionSorting(transactionsHandle, index4, 3, expec4, 2, &f4, "failed fee calc is filtered");
}
