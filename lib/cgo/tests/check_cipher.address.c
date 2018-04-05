
#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skytest.h"
#include "libsky_string.h"
#include "libcriterion.h"

#define SKYCOIN_ADDRESS_VALID "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"

// buffer big enough to hold all kind of data needed by test cases
unsigned char buff[1024];

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

Test(cipher, TestAddressVerify){

  PubKey pubkey;
  SecKey seckey;
  PubKey pubkey2;
  SecKey seckey2;
  Address addr;

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

Test(cipher,TestAddressString){

PubKey pubkey;
SecKey seckey;
SKY_cipher_GenerateKeyPair(&pubkey,&seckey);
Address a;
SKY_cipher_AddressFromPubKey(&pubkey,&a);
GoString_ s;
SKY_cipher_Address_String(&a,&s);
Address a2;
int error;

GoString tmp_s;
toGoString(&s,&tmp_s);
error = SKY_cipher_DecodeBase58Address( tmp_s,&a2 );
cr_assert(error == SKY_OK);

cr_assert(eq(type(Address),a,a2));

GoString_ s2;
SKY_cipher_Address_String(&a2,&s2);
GoString tmp_s2;
toGoString(&s,&tmp_s2);
Address a3;
error = SKY_cipher_DecodeBase58Address( tmp_s2,&a3 );
cr_assert(error == SKY_OK);
cr_assert(eq(type(Address),a2,a3));
}

// Test (cipher, TestBitcoinAddress1){

//   SecKey seckey;
//   PubKey pubkey;

//     GoString str = {
//     "1111111111111111111111111111111111111111111111111111111111111111",
//     64
//   };
//   unsigned  int  error;
//   error = SKY_cipher_SecKeyFromHex(str, &seckey);

//   cr_assert(error == SKY_OK);
//   error = SKY_cipher_PubKeyFromSecKey(&seckey,&pubkey);
//   cr_assert(error == SKY_OK);
//   GoString pubkeyStr = { "034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa", 66 };
//   GoString_ s1;
//   SKY_cipher_PubKey_Hex(&pubkey, &s1);
//   GoString tmp_s1;
//   toGoString(&s1,&tmp_s1);
//   cr_assert(eq(type(GoString), pubkeyStr, tmp_s1));

//   GoString bitcoinStr = {"1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9",34};
//   GoString_ s2;
//   SKY_cipher_BitcoinAddressFromPubkey(&pubkey,&s2);
//   GoString tmp_s2;
//   toGoString(&s2,&tmp_s2);
//   cr_assert(eq(type(GoString), bitcoinStr, tmp_s2));
// }

// func TestAddressRoundtrip(t *testing.T) {
//   p, _ := GenerateKeyPair()
//   a := AddressFromPubKey(p)
//   a2, err := addressFromBytes(a.Bytes())
//   require.NoError(t, err)
//   require.Equal(t, a, a2)
//   require.Equal(t, a.String(), a2.String())
// }
Test(cipher, TestBitcoinWIPRoundTrio){

  SecKey seckey;
  PubKey pubkey;
GoSlice slice;
  slice.data = buff;
  slice.cap = sizeof(buff);
  slice.len = 33;

  SKY_cipher_GenerateKeyPair(&pubkey,&seckey);

  GoString_ wip1;

  SKY_cipher_BitcoinWalletImportFormatFromSeckey(&seckey,&wip1);

  SecKey seckey2;

  unsigned int err;

  err = SKY_cipher_SecKeyFromWalletImportFormat( (*((GoString *) &wip1)) ,&seckey2);

  GoString_ wip2;

  SKY_cipher_BitcoinWalletImportFormatFromSeckey(&seckey2,&wip2);

  cr_assert(err == SKY_OK);

  // cr_assert(eq(type(char),(unsigned char *)seckey,(unsigned char *)seckey2));

  GoString_ seckeyhex1;
  GoString_ seckeyhex2;

  SKY_cipher_SecKey_Hex(&seckey,&seckeyhex1);

  SKY_cipher_SecKey_Hex(&seckey2,&seckeyhex2);
  cr_assert(eq(type(GoString_), seckeyhex1, seckeyhex2));

  cr_assert(eq(type(GoString_), wip1, wip2));

}



Test(cipher, TestBitcoinWIP ){

  //wallet input format string
GoString_ wip[3];

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

GoString_ addr[3];

addr[0].p="1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9";
addr[1].p="1NKRhS7iYUGTaAfaR5z8BueAJesqaTyc4a";
addr[2].p="19ck9VKC6KjGxR9LJg4DNMRc45qFrJguvV";

addr[0].n =34;
addr[1].n=34;
addr[2].n=34;


for (int i = 0; i < 3; ++i)
{
  SecKey seckey;

  unsigned int err;

  err = SKY_cipher_SecKeyFromWalletImportFormat( (*((GoString *) &wip[i])),&seckey);

  cr_assert(err==SKY_OK);

  PubKey pubkey;

SKY_cipher_PubKeyFromSecKey(&seckey,&pubkey);

unsigned char * pubkeyhextmp;

GoString_ string;

SKY_cipher_PubKey_Hex(&pubkey,&string);

cr_assert(eq(type(GoString_),string,pub[i]));

GoString_ bitcoinAddr;

SKY_cipher_BitcoinAddressFromPubkey(&pubkey,&bitcoinAddr);

cr_assert(eq(type(GoString_),addr[i],bitcoinAddr));

}
}

Test(cipher, TestAddressBulk){

for (int i = 0; i < 1024; ++i)
{
 GoSlice slice;
randBytes(&slice,32);
PubKey pubkey;
SecKey seckey;
//  SKY_cipher_GenerateDeterministicKeyPair( slice,&pubkey,&seckey);
   Address addr;
  SKY_cipher_AddressFromPubKey(&pubkey,&addr);
  unsigned int err;
  err = SKY_cipher_Address_Verify(&addr,&pubkey);
  cr_assert(err == SKY_OK);
  GoString_ strAddr;
  SKY_cipher_Address_String(&addr,&strAddr);
Address addr2;

err = SKY_cipher_DecodeBase58Address((*((GoString *) &strAddr)),&addr2);
 cr_assert(err == SKY_OK);
 cr_assert(eq(type(Address),addr,addr2));
}

}







