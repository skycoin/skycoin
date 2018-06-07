
#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

#define SKYCOIN_ADDRESS_VALID "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"

TestSuite(cipher_address, .init = setup, .fini = teardown);

// buffer big enough to hold all kind of data needed by test cases
unsigned char buff[1024];

Test(cipher_address, TestDecodeBase58Address) {

  GoString strAddr = {
    SKYCOIN_ADDRESS_VALID,
    35
  };
  cipher__Address addr;

  cr_assert( SKY_cipher_DecodeBase58Address(strAddr, &addr) == SKY_OK, "accept valid address");

  char tempStr[50];

  // preceding whitespace is invalid
  strcpy(tempStr, " ");
  strcat(tempStr, SKYCOIN_ADDRESS_VALID);
  strAddr.p = tempStr;
  strAddr.n = strlen(tempStr);
  cr_assert( SKY_cipher_DecodeBase58Address(strAddr, &addr) == SKY_ERROR, "preceding whitespace is invalid");

  // preceding zeroes are invalid
  strcpy(tempStr, "000");
  strcat(tempStr, SKYCOIN_ADDRESS_VALID);
  strAddr.p = tempStr;
  strAddr.n = strlen(tempStr);
  cr_assert( SKY_cipher_DecodeBase58Address(strAddr, &addr) == SKY_ERROR, "leading zeroes prefix are invalid");

  // trailing whitespace is invalid
  strcpy(tempStr, SKYCOIN_ADDRESS_VALID);
  strcat(tempStr, " ");
  strAddr.p = tempStr;
  strAddr.n = strlen(tempStr);
  cr_assert( SKY_cipher_DecodeBase58Address(strAddr, &addr) == SKY_ERROR, " trailing whitespace is invalid");

  // trailing zeroes are invalid
  strcpy(tempStr, SKYCOIN_ADDRESS_VALID);
  strcat(tempStr, "000");
  strAddr.p = tempStr;
  strAddr.n = strlen(tempStr);
  cr_assert( SKY_cipher_DecodeBase58Address(strAddr, &addr) == SKY_ERROR, " trailing zeroes suffix are invalid");

}

Test(cipher_address, TestAddressFromBytes){
  GoString strAddr = {
    SKYCOIN_ADDRESS_VALID,
    35
  };
  cipher__Address addr, addr2;
  GoSlice bytes;

  bytes.data = buff;
  bytes.len = 0;
  bytes.cap = sizeof(buff);

  SKY_cipher_DecodeBase58Address(strAddr, &addr);
  SKY_cipher_Address_BitcoinBytes(&addr, (GoSlice_ *)&bytes);
  cr_assert(bytes.len > 0, "address bytes written");
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2) == SKY_OK, "convert bytes to SKY address");

  cr_assert(eq(type(cipher__Address), addr, addr2));

  int bytes_len = bytes.len;

  bytes.len = bytes.len - 2;
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2) == SKY_ERROR, "no SKY address due to short bytes length");

  bytes.len = bytes_len;
  ((char *) bytes.data)[bytes.len - 1] = '2';
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2) == SKY_ERROR, "no SKY address due to corrupted bytes");
}

Test(cipher_address, TestAddressVerify){

  cipher__PubKey pubkey;
  cipher__SecKey seckey;
  cipher__PubKey pubkey2;
  cipher__SecKey seckey2;
  cipher__Address addr;

  SKY_cipher_GenerateKeyPair(&pubkey,&seckey);
  SKY_cipher_AddressFromPubKey(&pubkey,&addr);

  // Valid pubkey+address
  cr_assert( SKY_cipher_Address_Verify(&addr,&pubkey) == SKY_OK ,"Valid pubkey + address");

  SKY_cipher_GenerateKeyPair(&pubkey,&seckey2);
  //   // Invalid pubkey
  cr_assert( SKY_cipher_Address_Verify(&addr,&pubkey) == SKY_ERROR," Invalid pubkey");

  // Bad version
  addr.Version = 0x01;
  cr_assert( SKY_cipher_Address_Verify(&addr,&pubkey) == SKY_ERROR,"  Bad version");
}

Test(cipher_address,TestAddressString){

}

Test (cipher, TestBitcoinAddress1){

  cipher__SecKey seckey;
  cipher__PubKey pubkey;

  GoString str = {
    "1111111111111111111111111111111111111111111111111111111111111111",
    64
  }, s1, s2;

  unsigned  int  error;
  error = SKY_cipher_SecKeyFromHex(str, &seckey);
  cr_assert(error == SKY_OK, "Create SecKey from Hex");
  error = SKY_cipher_PubKeyFromSecKey(&seckey,&pubkey);
  cr_assert(error == SKY_OK, "Create PubKey from SecKey");

  GoString pubkeyStr = { "034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa", 66 };

  SKY_cipher_PubKey_Hex(&pubkey, (GoString_ *) &s1);
  registerMemCleanup((void *) s1.p);
  cr_assert(eq(type(GoString), pubkeyStr, s1));

  GoString bitcoinStr = {"1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9",34};
  SKY_cipher_BitcoinAddressFromPubkey(&pubkey, (GoString_ *) &s2);
  registerMemCleanup((void *) s2.p);
  cr_assert(eq(type(GoString), bitcoinStr, s2));
}

Test (cipher, TestBitcoinAddress2){

  cipher__SecKey seckey;
  cipher__PubKey pubkey  ;
  GoString str = {
    "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
    64
  }, s1, s2;

  unsigned  int error;
  error = SKY_cipher_SecKeyFromHex(str, &seckey);
  cr_assert(error == SKY_OK, "Create SecKey from Hex");
  error = SKY_cipher_PubKeyFromSecKey(&seckey,&pubkey);
  cr_assert(error == SKY_OK, "Create PubKey from SecKey");

  char strBuff[101];
  GoString pubkeyStr = {
    "02ed83704c95d829046f1ac27806211132102c34e9ac7ffa1b71110658e5b9d1bd",
    66
  };
  SKY_cipher_PubKey_Hex(&pubkey, (GoString_ *) &s1);
  registerMemCleanup((void *) s1.p);
  cr_assert(eq(type(GoString), pubkeyStr, s1));

  GoString bitcoinStr = {"1NKRhS7iYUGTaAfaR5z8BueAJesqaTyc4a",34};
  SKY_cipher_BitcoinAddressFromPubkey(&pubkey, (GoString_ *) &s2);
  registerMemCleanup((void *) s2.p);
  cr_assert(eq(type(GoString), bitcoinStr, s2));

}

Test (cipher, TestBitcoinAddress3){

  cipher__SecKey seckey;
  cipher__PubKey pubkey;
  GoString str = {
    "47f7616ea6f9b923076625b4488115de1ef1187f760e65f89eb6f4f7ff04b012",
    64
  };

  unsigned  int error;
  error = SKY_cipher_SecKeyFromHex(str, &seckey);
  cr_assert(error == SKY_OK, "Create SecKey from Hex");
  error = SKY_cipher_PubKeyFromSecKey(&seckey,&pubkey);
  cr_assert(error == SKY_OK, "Create PubKey from SecKey");

  char strBuff[101];
  GoString pubkeyStr = {
    "032596957532fc37e40486b910802ff45eeaa924548c0e1c080ef804e523ec3ed3",
    66
  }, s1, s2;

  SKY_cipher_PubKey_Hex(&pubkey, (GoString_ *)&s1);
  registerMemCleanup((void *) s1.p);
  cr_assert(eq(type(GoString), pubkeyStr, s1));

  GoString bitcoinStr = {"19ck9VKC6KjGxR9LJg4DNMRc45qFrJguvV",34};
  SKY_cipher_BitcoinAddressFromPubkey(&pubkey, (GoString_ *)&s2);
  registerMemCleanup((void *) s2.p);
  cr_assert(eq(type(GoString), bitcoinStr, s2));

}

Test(cipher_address, TestBitcoinWIPRoundTrio){

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

  err = SKY_cipher_SecKeyFromWalletImportFormat( (*((GoString *) &wip1)) ,&seckey2);

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


Test(cipher_address, TestBitcoinWIP ){

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

    err = SKY_cipher_SecKeyFromWalletImportFormat(wip[i],&seckey);
    cr_assert(err==SKY_OK);

    cipher__PubKey pubkey;

    SKY_cipher_PubKeyFromSecKey(&seckey,&pubkey);

    unsigned char * pubkeyhextmp;
    GoString_ string;

    SKY_cipher_PubKey_Hex(&pubkey,&string);
    cr_assert(eq(type(GoString), (*(GoString*)&string),(*(GoString*)&pub[i]) ));
    GoString bitcoinAddr;
    SKY_cipher_BitcoinAddressFromPubkey(&pubkey, (GoString_ *)&bitcoinAddr);
    cr_assert(eq(type(GoString),addr[i],bitcoinAddr));

  }
}

Test(cipher_address, TestAddressBulk){

  unsigned char buff[50];
  GoSlice slice = { buff, 0, 50 };

  for (int i = 0; i < 1024; ++i)
  {
    randBytes(&slice,32);
    cipher__PubKey pubkey;
    cipher__SecKey seckey;
    SKY_cipher_GenerateDeterministicKeyPair( slice,&pubkey,&seckey);
    cipher__Address addr;
    SKY_cipher_AddressFromPubKey(&pubkey,&addr);
    unsigned int err;
    err = SKY_cipher_Address_Verify(&addr,&pubkey);
    cr_assert(err == SKY_OK);
    GoString strAddr;
    SKY_cipher_Address_String(&addr, (GoString_ *)&strAddr);
    registerMemCleanup((void *) strAddr.p);
    cipher__Address addr2;

    err = SKY_cipher_DecodeBase58Address(strAddr,&addr2);
    cr_assert(err == SKY_OK);
    cr_assert(eq(type(cipher__Address),addr,addr2));
  }

}

