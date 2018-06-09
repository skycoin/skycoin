
#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

Test(coin_transaction, TestTransactionVerify)
{
  char bufferSHA[1024];
  char bufferSHA_[1024];
  int error;

  // Mismatch header hash
  coin__Transaction tx;
  coin__TransactionOutput* to_void;
  makeTransaction(&tx);
  memset(&tx.InnerHash, 0, sizeof(cipher__SHA256));
  GoUint32 errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode != SKY_OK, "Invalid header hash");

  // No inputs
  errcode = makeTransaction(&tx);
  cr_assert(errcode == SKY_OK);
  memset(&tx.In, 0, sizeof(GoSlice));
  errcode = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert(errcode == SKY_OK);
  errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode != SKY_OK, "No inputs");

  // No outputs
  errcode = makeTransaction(&tx);
  cr_assert(errcode == SKY_OK);
  memset(&tx.Out, 0, sizeof(GoSlice));
  errcode = SKY_coin_Transaction_UpdateHeader(&tx);
  cr_assert(errcode == SKY_OK);
  errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode != SKY_OK, "No outputs");

  // Invalid number of sigs
  errcode = makeTransaction(&tx);
  cr_assert(errcode == SKY_OK);
  memset(&tx.Sigs, 0, sizeof(cipher__Sig));
  errcode = SKY_coin_Transaction_UpdateHeader(&tx);
  errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode != SKY_OK, "Invalid number of signatures");

  errcode = makeTransaction(&tx);
  cr_assert(errcode == SKY_OK);
  GoSlice slice = { NULL, 20, 20 };
  memset(&tx.Sigs, 0, sizeof(slice));
  errcode = SKY_coin_Transaction_UpdateHeader(&tx);
  errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode != SKY_OK, "Invalid number of signatures");

  errcode = makeTransaction(&tx);
  cr_assert(errcode == SKY_OK);
  GoSlice slice_sigs = { NULL, 32768, 32768 };
  GoSlice slice_in = { NULL, 32768, 32768 };
  memset(&tx.Sigs, 0, sizeof(slice_sigs));
  memset(&tx.In, 0, sizeof(slice_in));
  errcode = SKY_coin_Transaction_UpdateHeader(&tx);
  errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode != SKY_OK, "Too many signatures and inputs");

  // Duplicate inputs
  coin__UxOut* ux;
  cipher__SecKey* s;
  errcode = makeUxOutWithSecret(&ux, &s);
  cr_assert(errcode == SKY_OK);
  makeTransactionFromUxOut(&ux, &s, &tx);
  cipher__SHA256 sha256;
  SKY_coin_Transaction_Hash(&tx, &sha256);
  SKY_coin_Transaction_PushInput(&tx, &sha256, 0);
  tx.Sigs.data = NULL;
  GoSlice slice_duplicate;
  copySlice(((GoSlice*)&slice_duplicate.data),
            (GoSlice*)&tx.In.data,
            sizeof(cipher__SecKey));
  SKY_coin_Transaction_SignInputs(&tx, slice_duplicate);
  errcode = SKY_coin_Transaction_UpdateHeader(&tx);
  errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode != SKY_OK, "Duplicate spend");

  // Duplicate outputs
  makeTransaction(&tx);
  coin__TransactionOutput to;
  to = (*(coin__TransactionOutput*)&tx.Out);
  errcode = SKY_coin_Transaction_PushOutput(
    &to, ((cipher__Address*)&to.Address), to.Coins, to.Hours);
  errcode = SKY_coin_Transaction_UpdateHeader(&tx);
  errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode != SKY_OK, "Duplicate output in transaction");

  // Invalid signature, empty
  makeTransaction(&tx);
  memset(&tx.Sigs, 0, sizeof(cipher__Sig));
  errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode != SKY_OK, "Duplicate spend");

  // Output coins are 0

  makeTransaction(&tx);
  memset(&tx.Out, 0, sizeof(coin__TransactionOutput));
  to_void = ((coin__TransactionOutput*)&tx.Out);
  to_void->Coins = 0;
  errcode = SKY_coin_Transaction_UpdateHeader(&tx);
  errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode != SKY_OK, "Zero coin output");

  // Output coin overflow
  makeTransaction(&tx);
  memset(&tx.Out, 0, sizeof(coin__TransactionOutput));
  to_void = ((coin__TransactionOutput*)&tx.Out);
  to_void->Coins = 9223372036851775808;
  errcode = SKY_coin_Transaction_UpdateHeader(&tx);
  errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode != SKY_OK, "Output coins overflow");

  // Output coins are not multiples of 1e6 (valid, decimal restriction is not
  // enforced here)

  memset(&tx.Out, 0, sizeof(coin__TransactionOutput));
  makeTransaction(&tx);
  to = (*(coin__TransactionOutput*)&tx.Out);
  to.Coins += 10;
  tx.Sigs.data = NULL;
  errcode = SKY_coin_Transaction_UpdateHeader(&tx);
  SKY_coin_Transaction_PushInput(&tx, &sha256, 0);
  GoSlice slice_decimal;
  memset(&slice_decimal.data, 0, sizeof(cipher__SecKey));
  SKY_coin_Transaction_SignInputs(&tx, slice_decimal);
  cr_assert(0 != (to.Coins % ((GoUint64_)1.000000e+006)));
  errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode != SKY_OK);

  // Valid
  memset(&tx.Out, 0, sizeof(coin__TransactionOutput));
  makeTransaction(&tx);
  to_void = ((coin__TransactionOutput*)&tx.Out);
  to_void->Coins = 1.000000e+007;
  to_void++;
  to_void->Coins = 1.000000e+006;
  errcode = SKY_coin_Transaction_UpdateHeader(&tx);
  errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode == SKY_OK);
}

// Test(coin_transaction, TestTransactionVerifyInput)
// {
// coin__Transaction tx;
// GoUint64_ errcode;
// coin__UxArray uxArray;

// // Invalid uxIn args
// makeTransaction(&tx);
// cli__PasswordFromBytes seckey;
// SKY_coin_UxArray_Coins(&seckey, NULL);
// errcode = SKY_coin_Transaction_VerifyInput(&tx, &seckey);
// cr_assert(errcode != SKY_OK, "tx.In != uxIn");
// SKY_coin_UxArray_Coins(&seckey, 0);
// errcode = SKY_coin_Transaction_VerifyInput(&tx, &seckey);
// cr_assert(errcode != SKY_OK, "tx.In != uxIn");
// SKY_coin_UxArray_Coins(&seckey, 3);
// errcode = SKY_coin_Transaction_VerifyInput(&tx, &seckey);
// cr_assert(errcode != SKY_OK, "tx.In != uxIn");

// // 	// tx.In != tx.Sigs
// // ux, s := makeUxOutWithSecret(t)
// // tx = makeTransactionFromUxOut(ux, s)
// // tx.Sigs = []cipher.Sig{}
// // _require.PanicsWithLogMessage(t, "tx.In != tx.Sigs", func() {
// // 	tx.VerifyInput(UxArray{ux})
// // })

// // ux, s = makeUxOutWithSecret(t)
// // tx = makeTransactionFromUxOut(ux, s)
// // tx.Sigs = append(tx.Sigs, cipher.Sig{})
// // _require.PanicsWithLogMessage(t, "tx.In != tx.Sigs", func() {
// // 	tx.VerifyInput(UxArray{ux})
// // })

// // tx.InnerHash != tx.HashInner()
// coin__UxOut ux;
// cipher__SecKey s;
// errcode = makeUxOutWithSecret(&ux, &s);
// cr_assert(errcode == SKY_OK);
// errcode = makeTransactionFromUxOut(&ux, &s, &tx);
// cr_assert(errcode == SKY_OK);
// memset(&tx.Sigs, 0, sizeof(cipher__Sig));
// memset(&uxArray, 0, sizeof(coin__UxArray));
// uxArray.data = &ux;
// SKY_coin_UxArray_Coins(&uxArray, 1);
// errcode = SKY_coin_Transaction_VerifyInput(&tx, &uxArray);
// cr_assert(errcode != SKY_OK, "tx.In != tx.Sigs");

// errcode = makeUxOutWithSecret(&ux, &s);
// cr_assert(errcode == SKY_OK);
// errcode = makeTransactionFromUxOut(&ux, &s, &tx);
// cr_assert(errcode == SKY_OK);

// // coin__UxOut uxo;
// // coin__UxArray uxa;
// // uxa.data = uxo.;
// }

Test(coin_transaction, TestTransactionPushInput)
{
  coin__Transaction tx;
  memset(&tx, 0, sizeof(coin__Transaction));
  coin__UxOut ux;
  memset(&ux, 0, sizeof(coin__UxOut));
  GoUint64_ errcode = makeUxOut(&ux);
  GoUint16 value;
  SKY_coin_Transaction_PushInput(&tx, &ux, &value);
  cr_assert(value == 0);
  cr_assert(tx.In.len == 1);
  errcode = memcmp(((cipher__SHA256*)&tx.In.data), &ux, sizeof(cipher__SHA256));
  cr_assert(errcode > 0);

  cipher__SHA256* cipher;

  cipher = ((cipher__SHA256*)&tx.In.data);
  makeRandHash(&cipher);
  for (int i = 0; i < (1 << 16 - 1); ++i) {
    cipher++;
    makeRandHash(&cipher);
  }
  errcode = makeUxOut(&ux);
  errcode = SKY_coin_Transaction_PushInput(&tx, &ux, &value);

  cr_assert(errcode == SKY_OK);
}

Test(coin_transaction, TestTransactionPushOutput)
{
  coin__Transaction tx;
  cipher__Address a;
  GoUint32 errcode = makeAddress(a);
  cr_assert(errcode == SKY_OK);
  memset(&tx, 0, sizeof(tx));
  errcode = SKY_coin_Transaction_PushOutput(&tx, &a, 100, 500);
  cr_assert(tx.Out.len == 1);
  coin__TransactionOutput value = { a, 100, 150 };
  GoSlice_ slice = { &value, 1, 1 };
  cr_assert(eq(type(GoSlice_), slice, tx.Out));

  memset(&tx, 0, sizeof(coin__Transaction));
  coin__TransactionOutput* slice_void =
    ((coin__TransactionOutput*)&tx.Out.data);
  for (int i = 0; i < 20; i++) {
    cipher__Address address;
    errcode = makeAddress(&address);
    cr_assert(errcode == SKY_OK);
    GoUint16 coins = ((GoUint16)(i * 100));
    GoUint16 hours = ((GoUint16)(i * 50));
    errcode = SKY_coin_Transaction_PushOutput(&tx, &address, coins, hours);
    cr_assert(tx.Out.len == (i + 1));
    // slice_void++;
  }
}



Test(coin_transaction, TestTransactionHash)
{
  coin__Transaction tx;
  cleanupMem(&tx);
  makeTransaction(&tx);
  cipher__SHA256 sha;
  cipher__SHA256 shatmp;
  cipher__SHA256 shainner;
  memset(&sha, 0, sizeof(cipher__SHA256));
  memset(&shatmp, 0, sizeof(cipher__SHA256));
  memset(&shainner, 0, sizeof(cipher__SHA256));
  SKY_coin_Transaction_Hash(&tx, &shatmp);
  SKY_coin_Transaction_HashInner(&tx, &shainner);

  cr_assert(not(eq(u8[sizeof(cipher__SHA256)], sha, shatmp)));

  cr_assert(not(eq(u8[sizeof(cipher__SHA256)], sha, shainner)));
}

Test(coin_transaction, TestTransactionUpdateHeader)
{
  coin__Transaction tx;
  cleanupMem(&tx);
  makeTransaction(&tx);
  cipher__SHA256 h;
  cipher__SHA256 shatmp;
  cipher__SHA256 sha_clean;
  cipher__SHA256 shainner;
  memset(&h, 0, sizeof(cipher__SHA256));
  memset(&sha_clean, 0, sizeof(cipher__SHA256));

  SKY_coin_Transaction_HashInner(&tx, &shainner);
  //
  memcpy(&h, &tx.InnerHash, sizeof(cipher__SHA256));
  memset(&tx.InnerHash, 0, sizeof(cipher__SHA256));

  SKY_coin_Transaction_UpdateHeader(&tx);
  SKY_coin_Transaction_HashInner(&tx, &shatmp);

  cr_assert(not(eq(u8[sizeof(cipher__SHA256)], tx.InnerHash, sha_clean)));
  cr_assert(eq(u8[sizeof(cipher__SHA256)], tx.InnerHash, h));
  cr_assert(eq(u8[sizeof(cipher__SHA256)], tx.InnerHash, shainner));
}
