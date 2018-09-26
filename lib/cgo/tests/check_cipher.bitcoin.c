
#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

#define SKYCOIN_ADDRESS_VALID "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"

TestSuite(cipher_bitcoin, .init = setup, .fini = teardown);

// buffer big enough to hold all kind of data needed by test cases
unsigned char buff[1024];

Test (cipher, TestBitcoinAddress){
  cipher__SecKey seckey;
  cipher__PubKey pubkey;
  cipher__BitcoinAddress btcAddr;
  GoString secKeys[3] = {
    {"1111111111111111111111111111111111111111111111111111111111111111", 64},
    {"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd", 64},
    {"47f7616ea6f9b923076625b4488115de1ef1187f760e65f89eb6f4f7ff04b012", 64}
  };
  GoString pubKeys[3] = {
    {"034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa", 66},
    {"02ed83704c95d829046f1ac27806211132102c34e9ac7ffa1b71110658e5b9d1bd", 66},
    {"032596957532fc37e40486b910802ff45eeaa924548c0e1c080ef804e523ec3ed3", 66}
  };
  GoString addrs[3] = {
    {"1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9", 34},
    {"1NKRhS7iYUGTaAfaR5z8BueAJesqaTyc4a", 34},
    {"19ck9VKC6KjGxR9LJg4DNMRc45qFrJguvV", 34}
  };
  unsigned  int  error;

  GoString *secKeyStr = secKeys;
  GoString *pubKeyStr = pubKeys;
  GoString *addrStr = addrs;

  for (int i = 0; i < 3; ++i, ++secKeyStr, ++pubKeyStr, ++addrStr) {
    error = SKY_cipher_SecKeyFromHex(*secKeyStr, &seckey);
    cr_assert(error == SKY_OK, "Create SecKey from Hex"); // (seckeyFailMsg));
    error = SKY_cipher_PubKeyFromHex(*pubKeyStr, &pubkey);
    cr_assert(error == SKY_OK, "Create PubKey from Hex");

    GoString str = {NULL, 0};
    SKY_cipher_BitcoinAddressFromPubKey(&pubkey, &btcAddr);
    SKY_cipher_BitcoinAddress_String(&btcAddr, (GoString_ *) &str);
    registerMemCleanup((void *) str.p);
    cr_assert(eq(type(GoString), str, (*addrStr)));

    str.p = NULL; str.n = 0;
    error = SKY_cipher_BitcoinAddressFromSecKey(&seckey, &btcAddr);
    cr_assert(error == SKY_OK);
    SKY_cipher_BitcoinAddress_String(&btcAddr, (GoString_ *) &str);
    registerMemCleanup((void *) str.p);
    cr_assert(eq(type(GoString), str, *addrStr));
  }
}

Test(cipher_address, TestBitcoinWIFRoundTrip){
  cipher__SecKey seckey;
  cipher__PubKey pubkey;
  GoSlice slice;
  slice.data = buff;
  slice.cap = sizeof(buff);
  slice.len = 33;
  SKY_cipher_GenerateKeyPair(&pubkey,&seckey);
  GoString_ wip1;
  SKY_cipher_BitcoinWalletImportFormatFromSeckey(&seckey,&wip1);
  cipher__SecKey seckey2;
  unsigned int err;
  err = SKY_cipher_SecKeyFromBitcoinWalletImportFormat( (*((GoString *) &wip1)) ,&seckey2);
  GoString_ wip2;
  SKY_cipher_BitcoinWalletImportFormatFromSeckey(&seckey2,&wip2);
  cr_assert(err == SKY_OK);
  cr_assert(eq(u8[sizeof(cipher__SecKey)],seckey,seckey2));
  GoString_ seckeyhex1;
  GoString_ seckeyhex2;
  SKY_cipher_SecKey_Hex(&seckey,&seckeyhex1);
  SKY_cipher_SecKey_Hex(&seckey2,&seckeyhex2);
  cr_assert(eq(type(GoString), (*(GoString*)&seckeyhex1),(*(GoString*)&seckeyhex2) ));
  cr_assert(eq(type(GoString), (*(GoString*)&wip1),(*(GoString*)&wip2) ));
}

Test(cipher_address, TestBitcoinWIF ){
  //wallet input format string
  GoString wip[3];
  wip[0].p = "KwntMbt59tTsj8xqpqYqRRWufyjGunvhSyeMo3NTYpFYzZbXJ5Hp";
  wip[1].p = "L4ezQvyC6QoBhxB4GVs9fAPhUKtbaXYUn8YTqoeXwbevQq4U92vN";
  wip[2].p = "KydbzBtk6uc7M6dXwEgTEH2sphZxSPbmDSz6kUUHi4eUpSQuhEbq";
  wip[0].n = 52;
  wip[1].n = 52;
  wip[2].n = 52;
  //   // //the expected pubkey to generate
  GoString_ pub[3];
  pub[0].p="034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa";
  pub[1].p="02ed83704c95d829046f1ac27806211132102c34e9ac7ffa1b71110658e5b9d1bd";
  pub[2].p="032596957532fc37e40486b910802ff45eeaa924548c0e1c080ef804e523ec3ed3";
  pub[0].n = 66;
  pub[1].n = 66;
  pub[2].n = 66;
  // //the expected addrss to generate
  GoString addr[3];
  addr[0].p="1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9";
  addr[1].p="1NKRhS7iYUGTaAfaR5z8BueAJesqaTyc4a";
  addr[2].p="19ck9VKC6KjGxR9LJg4DNMRc45qFrJguvV";
  addr[0].n =34;
  addr[1].n=34;
  addr[2].n=34;
  for (int i = 0; i < 3; i++)
  {
    cipher__SecKey seckey;
    unsigned int err;
    err = SKY_cipher_SecKeyFromBitcoinWalletImportFormat(wip[i],&seckey);
    cr_assert(err==SKY_OK);
    cipher__PubKey pubkey;
    SKY_cipher_PubKeyFromSecKey(&seckey,&pubkey);
    unsigned char * pubkeyhextmp;
    GoString_ string;
    SKY_cipher_PubKey_Hex(&pubkey,&string);
    cr_assert(eq(type(GoString), (*(GoString*)&string),(*(GoString*)&pub[i]) ));
    cipher__BitcoinAddress bitcoinAddr;
    GoString bitcoinAddrStr;
    SKY_cipher_BitcoinAddressFromPubKey(&pubkey, &bitcoinAddr);
    SKY_cipher_BitcoinAddress_String(&bitcoinAddr, (GoString_ *)&bitcoinAddrStr);
    cr_assert(eq(type(GoString),addr[i],bitcoinAddrStr));
  }
}
