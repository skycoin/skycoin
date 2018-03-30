
#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"

#define SKYCOIN_ADDRESS_VALID "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"

// buffer big enough to hold all kind of data needed by test cases
unsigned char buff[1024];

// // TODO: Write like this cr_assert(eq(type(Address), addr1, addr2))
int cr_user_Address_eq(Address *addr1, Address *addr2){
  if(addr1->Version != addr2->Version)
    return 0;
  for (int i = 0; i < sizeof(Ripemd160); ++i) {
    if(addr1->Key[i] != addr2->Key[i])
      return 0;
  }
  return 1;
}

char *cr_user_Address_tostr(Address *addr1)
{
  char *out;

  cr_asprintf(&out, "(Address) { .Key = %s, .Version = %llu }", addr1->Key, (unsigned long long) addr1->Version);
  return out;
}
// // TODO: Write like this cr_assert(not(eq(type(Address), addr1, addr2)))
int cr_user_Address_noteq(Address *addr1, Address *addr2){
  if(addr1->Version != addr2->Version)
    return SKY_OK;
  for (int i = 0; i < sizeof(Ripemd160); ++i) {
    if(addr1->Key[i] != addr2->Key[i])
      return SKY_OK;
  }
  return SKY_ERROR;
}

int cr_user_GoString_eq(GoString *string1, GoString *string2){

  if(  strcmp(string1->p,string2->p) != 0 )
  {
    return SKY_ERROR;
  } else {
    return SKY_OK;
  }
}

char *cr_user_GoString_tostr(GoString *string)
{
  char *out;
  cr_asprintf(&out, "(GoString) { .Data = %s, .Length = %llu }", string->p, (unsigned long long) string->n);
  return out;
}

int cr_user_GoString__eq(GoString_ *string1, GoString_ *string2){
  return cr_user_GoString_eq((GoString *)string1, (GoString *)string2);
}

char *cr_user_GoString__tostr(GoString_ *string) {
  return cr_user_GoString_tostr((GoString *)string);
}

Test(cipher, TestDecodeBase58Address) {

 GoString strAddr = {
  SKYCOIN_ADDRESS_VALID,
  35
};
Address addr;

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

Test(cipher, TestAddressFromBytes){
  GoString strAddr = {
    SKYCOIN_ADDRESS_VALID,
    35
  };
  Address addr, addr2;
  GoSlice bytes;

  bytes.data = buff;
  bytes.len = 0;
  bytes.cap = sizeof(buff);

  SKY_cipher_DecodeBase58Address(strAddr, &addr);
  SKY_cipher_Address_BitcoinBytes(&addr, (GoSlice_ *)&bytes);
  cr_assert(bytes.len > 0, "address bytes written");
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2) == SKY_OK, "convert bytes to SKY address");

  cr_assert(eq(type(Address), addr, addr2));

  int bytes_len = bytes.len;

  bytes.len = bytes.len - 2;
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2) == SKY_ERROR, "no SKY address due to short bytes length");

  bytes.len = bytes_len;
  ((char *) bytes.data)[bytes.len - 1] = '2';
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2) == SKY_ERROR, "no SKY address due to corrupted bytes");
}

Test (cipher, TestBitcoinAddress1){

  SecKey seckey;
  PubKey pubkey;
  GoString str = {
    "1111111111111111111111111111111111111111111111111111111111111111",
    64
  };

  SKY_cipher_SecKeyFromHex(str, &seckey);
  unsigned  int  error;
  error = SKY_cipher_PubKeyFromSecKey(&seckey,&pubkey);
  cr_assert(error == SKY_OK, "Create PubKey from SecKey");

  char pubkeyStr[67];
  strcpy( pubkeyStr, "034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa");
  char  pubkeyhex[101];
  strcpy(pubkeyhex,SKY_cipher_PubKey_Hex(&pubkey));
  cr_assert( strcmp(pubkeyStr,pubkeyhex) == 0);

  GoString_ bitcoinAddr;

  GoString_ bitcoinStr = {"1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9",34};
  SKY_cipher_BitcoinAddressFromPubkey(&pubkey, &bitcoinAddr);
  cr_assert(eq(type(GoString_), bitcoinStr, bitcoinAddr));

}

Test (cipher, TestBitcoinAddress2){

  SecKey seckey;
  PubKey pubkey  ;
  GoString str = {
    "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
    64
  };

  SKY_cipher_SecKeyFromHex(str, &seckey);
  unsigned  int error;
  error = SKY_cipher_PubKeyFromSecKey(&seckey,&pubkey);

  cr_assert(error == SKY_OK, "Create PubKey from SecKey");

  char pubkeyStr[67];
  strcpy( pubkeyStr, "02ed83704c95d829046f1ac27806211132102c34e9ac7ffa1b71110658e5b9d1bd");
  char  pubkeyhex[31];
  strcpy(pubkeyhex,SKY_cipher_PubKey_Hex(&pubkey));
  cr_assert( strcmp(pubkeyStr,pubkeyhex) == 0);

  GoString_ bitcoinAddr;

  GoString_ bitcoinStr = {"1NKRhS7iYUGTaAfaR5z8BueAJesqaTyc4a",34};
  SKY_cipher_BitcoinAddressFromPubkey(&pubkey, &bitcoinAddr);
  cr_assert(eq(type(GoString_), bitcoinStr, bitcoinAddr));

}

Test (cipher, TestBitcoinAddress3){

  SecKey seckey;
  PubKey pubkey;
  GoString str = {
    "47f7616ea6f9b923076625b4488115de1ef1187f760e65f89eb6f4f7ff04b012",
    64
  };

  SKY_cipher_SecKeyFromHex(str, &seckey);
  unsigned  int error;
  error = SKY_cipher_PubKeyFromSecKey(&seckey,&pubkey);

  cr_assert(error == SKY_OK, "Create PubKey from SecKey");

  char pubkeyStr[67];
  strcpy( pubkeyStr, "032596957532fc37e40486b910802ff45eeaa924548c0e1c080ef804e523ec3ed3");
  char  pubkeyhex[31];
  strcpy(pubkeyhex,SKY_cipher_PubKey_Hex(&pubkey));
  cr_assert( strcmp(pubkeyStr,pubkeyhex) == 0);

  GoString_ bitcoinAddr;

  GoString_ bitcoinStr = {"19ck9VKC6KjGxR9LJg4DNMRc45qFrJguvV",34};
  SKY_cipher_BitcoinAddressFromPubkey(&pubkey, &bitcoinAddr);
  cr_assert(eq(type(GoString_), bitcoinStr, bitcoinAddr));

}

Test(cipher, TestAddressVerify){

  PubKey pubkey;
  PubKey pubkey2;
  GoSlice slice;
  GoSlice slice2;
  
  slice.data = buff;
  slice.cap = sizeof(buff);
  slice.len = 33;

  slice2.data = buff;
  slice2.cap = sizeof(buff);
  slice2.len = 33;
  Address addr;

  // SKY_cipher_RandByte(33,&slice);
  // SKY_cipher_RandByte(33,&slice2);

  SKY_cipher_NewPubKey(slice,&pubkey);
  SKY_cipher_NewPubKey(slice,&pubkey2);

  SKY_cipher_AddressFromPubKey(&pubkey,&addr);

  // Valid pubkey+address
  cr_assert( SKY_cipher_Address_Verify(&addr,&pubkey) == SKY_OK ,"Valid pubkey + address");

//   // Invalid pubkey
  cr_assert( SKY_cipher_Address_Verify(&addr,&pubkey2) == SKY_ERROR," Invalid pubkey");

  // Bad version
  addr.Version = 0x01;
  cr_assert( SKY_cipher_Address_Verify(&addr,&pubkey) == SKY_ERROR,"  Bad version");
}
