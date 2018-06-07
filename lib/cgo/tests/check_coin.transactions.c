
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
  memset(&tx.Sigs.data, 0, sizeof(slice));
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
  SKY_coin_Transaction_PushInput(&tx, &sha256 , 0);
  tx.Sigs.data = NULL;
  GoSlice slice_duplicate;
  copySlice( ((GoSlice*)&slice_duplicate.data), (GoSlice*)&tx.In.data,sizeof(cipher__SecKey) );
  SKY_coin_Transaction_SignInputs(&tx, slice_duplicate);
  errcode = SKY_coin_Transaction_UpdateHeader(&tx);
  errcode = SKY_coin_Transaction_Verify(&tx);
  cr_assert(errcode != SKY_OK, "Duplicate spend");

}

// Test(coin_transaction, TestTransactionPushInput)
// {

//   coin__Transaction tx;
//   coin__UxOut ux;
//   int errcode = makeUxOut(&ux);
//   cr_assert(errcode == SKY_OK);
// }
