
#include <stdio.h>
#include <string.h>
#include <signal.h>
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

  cipher__Address a;
  cipher__PubKey p;
  cipher__SecKey s;

  cr_assert(SKY_cipher_GenerateKeyPair(&p, &s) == SKY_OK);
  cr_assert(SKY_cipher_AddressFromPubKey(&p, &a) == SKY_OK);
  cr_assert(SKY_cipher_Address_Verify(&a, &p) == SKY_OK);

  cipher__Address addr;
  GoString strAddr = {SKYCOIN_ADDRESS_VALID, 35};
  cr_assert(SKY_cipher_DecodeBase58Address(strAddr, &addr) == SKY_OK,
            "accept valid address");

  int errorcode;
  char tempStr[50];
  // preceding whitespace is invalid
  strcpy(tempStr, " ");
  strcat(tempStr, SKYCOIN_ADDRESS_VALID);
  strAddr.p = tempStr;
  strAddr.n = strlen(tempStr);
  errorcode = SKY_cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(
      errorcode == SKY_ErrInvalidBase58Char,
      "preceding whitespace is invalid"
  );

  // preceding zeroes are invalid
  strcpy(tempStr, "000");
  strcat(tempStr, SKYCOIN_ADDRESS_VALID);
  strAddr.p = tempStr;
  strAddr.n = strlen(tempStr);
  errorcode = SKY_cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(
      errorcode == SKY_ErrInvalidBase58Char,
      "leading zeroes prefix are invalid"
  );

  // trailing whitespace is invalid
  strcpy(tempStr, SKYCOIN_ADDRESS_VALID);
  strcat(tempStr, " ");
  strAddr.p = tempStr;
  strAddr.n = strlen(tempStr);
  errorcode = SKY_cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(
      errorcode == SKY_ErrInvalidBase58Char,
      "trailing whitespace is invalid"
  );

  // trailing zeroes are invalid
  strcpy(tempStr, SKYCOIN_ADDRESS_VALID);
  strcat(tempStr, "000");
  strAddr.p = tempStr;
  strAddr.n = strlen(tempStr);
  errorcode = SKY_cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(
      errorcode == SKY_ErrInvalidBase58Char,
      "trailing zeroes suffix are invalid"
  );
}

Test(cipher_address, TestAddressFromBytes) {
  GoString strAddr = {SKYCOIN_ADDRESS_VALID, 35};
  cipher__Address addr, addr2;
  GoSlice bytes;

  bytes.data = buff;
  bytes.len = 0;
  bytes.cap = sizeof(buff);

  SKY_cipher_DecodeBase58Address(strAddr, &addr);
  SKY_cipher_Address_BitcoinBytes(&addr, (GoSlice_ *)&bytes);
  cr_assert(bytes.len > 0, "address bytes written");
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2) == SKY_OK,
            "convert bytes to SKY address");

  cr_assert(eq(type(cipher__Address), addr, addr2));

  int bytes_len = bytes.len;

  bytes.len = bytes.len - 2;
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2) == SKY_ErrAddressInvalidLength, "no SKY address due to short bytes length");

  bytes.len = bytes_len;
  ((char *) bytes.data)[bytes.len - 1] = '2';
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2) == SKY_ErrAddressInvalidChecksum, "no SKY address due to corrupted bytes");
}

Test(cipher_address, TestAddressVerify) {

  cipher__PubKey pubkey;
  cipher__SecKey seckey;
  cipher__PubKey pubkey2;
  cipher__SecKey seckey2;
  cipher__Address addr;

  SKY_cipher_GenerateKeyPair(&pubkey, &seckey);
  SKY_cipher_AddressFromPubKey(&pubkey, &addr);

  // Valid pubkey+address
  cr_assert(SKY_cipher_Address_Verify(&addr, &pubkey) == SKY_OK,
            "Valid pubkey + address");

  SKY_cipher_GenerateKeyPair(&pubkey, &seckey2);
  //   // Invalid pubkey
  cr_assert( SKY_cipher_Address_Verify(&addr,&pubkey) == SKY_ErrAddressInvalidPubKey," Invalid pubkey");

  // Bad version
  addr.Version = 0x01;
  cr_assert( SKY_cipher_Address_Verify(&addr,&pubkey) == SKY_ErrAddressInvalidVersion,"  Bad version");
}

Test(cipher_address, TestAddressString) {
  cipher__PubKey p;
  cipher__SecKey s1;
  GoInt result;
  result = SKY_cipher_GenerateKeyPair(&p, &s1);
  cr_assert(result == SKY_OK, "SKY_cipher_GenerateKeyPair failed");

  cipher__Address a;
  result = SKY_cipher_AddressFromPubKey(&p, &a);
  cr_assert(result == SKY_OK, "SKY_cipher_AddressFromPubKey failed");
  char buffer_s[1024];
  GoString_ s = {buffer_s, 0};
  result = SKY_cipher_Address_String(&a, &s);
  cr_assert(result == SKY_OK, "SKY_cipher_Address_String failed");
  cipher__Address a2;
  char buffer_sConvert[1024];
  GoString sConvert = {buffer_sConvert, 0};
  toGoString(&s, &sConvert);
  result = SKY_cipher_DecodeBase58Address(sConvert, &a2);
  cr_assert(result == SKY_OK);
  cr_assert(eq(type(cipher__Address), a, a2));
  char buffer_s2[1024];
  GoString_ s2 = {buffer_s2, 0};
  result = SKY_cipher_Address_String(&a2, &s2);
  cr_assert(result == SKY_OK, "SKY_cipher_Address_String failed");
  cipher__Address a3;
  char buffer_s2Convert[1024];
  GoString s2Convert = {buffer_s2Convert, 0};
  toGoString(&s2, &s2Convert);
  result = SKY_cipher_DecodeBase58Address(s2Convert, &a3);
  cr_assert(result == SKY_OK, "SKY_cipher_DecodeBase58Address failed");
  cr_assert(eq(type(cipher__Address), a2, a3));
}

Test(cipher, TestBitcoinAddress1) {

  cipher__SecKey seckey;
  cipher__PubKey pubkey;

  GoString
      str = {"1111111111111111111111111111111111111111111111111111111111111111",
             64},
      s1, s2;

  unsigned int error;
  error = SKY_cipher_SecKeyFromHex(str, &seckey);
  cr_assert(error == SKY_OK, "Create SecKey from Hex");
  error = SKY_cipher_PubKeyFromSecKey(&seckey, &pubkey);
  cr_assert(error == SKY_OK, "Create PubKey from SecKey");

  GoString pubkeyStr = {
      "034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa", 66};

  SKY_cipher_PubKey_Hex(&pubkey, (GoString_ *)&s1);
  registerMemCleanup((void *)s1.p);
  cr_assert(eq(type(GoString), pubkeyStr, s1));

  GoString bitcoinStr = {"1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9", 34};
  SKY_cipher_BitcoinAddressFromPubkey(&pubkey, (GoString_ *)&s2);
  registerMemCleanup((void *)s2.p);
  cr_assert(eq(type(GoString), bitcoinStr, s2));
}

Test(cipher, TestBitcoinAddress2) {

  cipher__SecKey seckey;
  cipher__PubKey pubkey;
  GoString
      str = {"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
             64},
      s1, s2;

  unsigned int error;
  error = SKY_cipher_SecKeyFromHex(str, &seckey);
  cr_assert(error == SKY_OK, "Create SecKey from Hex");
  error = SKY_cipher_PubKeyFromSecKey(&seckey, &pubkey);
  cr_assert(error == SKY_OK, "Create PubKey from SecKey");

  char strBuff[101];
  GoString pubkeyStr = {
      "02ed83704c95d829046f1ac27806211132102c34e9ac7ffa1b71110658e5b9d1bd", 66};
  SKY_cipher_PubKey_Hex(&pubkey, (GoString_ *)&s1);
  registerMemCleanup((void *)s1.p);
  cr_assert(eq(type(GoString), pubkeyStr, s1));

  GoString bitcoinStr = {"1NKRhS7iYUGTaAfaR5z8BueAJesqaTyc4a", 34};
  SKY_cipher_BitcoinAddressFromPubkey(&pubkey, (GoString_ *)&s2);
  registerMemCleanup((void *)s2.p);
  cr_assert(eq(type(GoString), bitcoinStr, s2));
}

Test(cipher, TestBitcoinAddress3) {

  cipher__SecKey seckey;
  cipher__PubKey pubkey;
  GoString str = {
      "47f7616ea6f9b923076625b4488115de1ef1187f760e65f89eb6f4f7ff04b012", 64};

  unsigned int error;
  error = SKY_cipher_SecKeyFromHex(str, &seckey);
  cr_assert(error == SKY_OK, "Create SecKey from Hex");
  error = SKY_cipher_PubKeyFromSecKey(&seckey, &pubkey);
  cr_assert(error == SKY_OK, "Create PubKey from SecKey");

  char strBuff[101];
  GoString
      pubkeyStr =
          {"032596957532fc37e40486b910802ff45eeaa924548c0e1c080ef804e523ec3ed3",
           66},
      s1, s2;

  SKY_cipher_PubKey_Hex(&pubkey, (GoString_ *)&s1);
  registerMemCleanup((void *)s1.p);
  cr_assert(eq(type(GoString), pubkeyStr, s1));

  GoString bitcoinStr = {"19ck9VKC6KjGxR9LJg4DNMRc45qFrJguvV", 34};
  SKY_cipher_BitcoinAddressFromPubkey(&pubkey, (GoString_ *)&s2);
  registerMemCleanup((void *)s2.p);
  cr_assert(eq(type(GoString), bitcoinStr, s2));
}

Test(cipher_address, TestBitcoinWIPRoundTrio) {

  cipher__SecKey seckey;
  cipher__PubKey pubkey;
  GoSlice slice;
  slice.data = buff;
  slice.cap = sizeof(buff);
  slice.len = 33;

  SKY_cipher_GenerateKeyPair(&pubkey, &seckey);

  GoString_ wip1;

  SKY_cipher_BitcoinWalletImportFormatFromSeckey(&seckey, &wip1);

  cipher__SecKey seckey2;

  unsigned int err;

  err =
      SKY_cipher_SecKeyFromWalletImportFormat((*((GoString *)&wip1)), &seckey2);

  GoString_ wip2;

  SKY_cipher_BitcoinWalletImportFormatFromSeckey(&seckey2, &wip2);

  cr_assert(err == SKY_OK);

  cr_assert(eq(u8[sizeof(cipher__SecKey)], seckey, seckey2));

  GoString_ seckeyhex1;
  GoString_ seckeyhex2;

  SKY_cipher_SecKey_Hex(&seckey, &seckeyhex1);
  SKY_cipher_SecKey_Hex(&seckey2, &seckeyhex2);
  cr_assert(eq(type(GoString), (*(GoString *)&seckeyhex1),
               (*(GoString *)&seckeyhex2)));
  cr_assert(eq(type(GoString), (*(GoString *)&wip1), (*(GoString *)&wip2)));
}

Test(cipher_address, TestBitcoinWIP) {

  // wallet input format string
  GoString wip[3];

  wip[0].p = "KwntMbt59tTsj8xqpqYqRRWufyjGunvhSyeMo3NTYpFYzZbXJ5Hp";
  wip[1].p = "L4ezQvyC6QoBhxB4GVs9fAPhUKtbaXYUn8YTqoeXwbevQq4U92vN";
  wip[2].p = "KydbzBtk6uc7M6dXwEgTEH2sphZxSPbmDSz6kUUHi4eUpSQuhEbq";
  wip[0].n = 52;
  wip[1].n = 52;
  wip[2].n = 52;

  //   // //the expected pubkey to generate
  GoString_ pub[3];

  pub[0].p =
      "034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa";
  pub[1].p =
      "02ed83704c95d829046f1ac27806211132102c34e9ac7ffa1b71110658e5b9d1bd";
  pub[2].p =
      "032596957532fc37e40486b910802ff45eeaa924548c0e1c080ef804e523ec3ed3";

  pub[0].n = 66;
  pub[1].n = 66;
  pub[2].n = 66;

  // //the expected addrss to generate

  GoString addr[3];

  addr[0].p = "1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9";
  addr[1].p = "1NKRhS7iYUGTaAfaR5z8BueAJesqaTyc4a";
  addr[2].p = "19ck9VKC6KjGxR9LJg4DNMRc45qFrJguvV";

  addr[0].n = 34;
  addr[1].n = 34;
  addr[2].n = 34;

  for (int i = 0; i < 3; i++) {
    cipher__SecKey seckey;
    unsigned int err;

    err = SKY_cipher_SecKeyFromWalletImportFormat(wip[i], &seckey);
    cr_assert(err == SKY_OK);

    cipher__PubKey pubkey;

    SKY_cipher_PubKeyFromSecKey(&seckey, &pubkey);

    unsigned char *pubkeyhextmp;
    GoString_ string;

    SKY_cipher_PubKey_Hex(&pubkey, &string);
    cr_assert(
        eq(type(GoString), (*(GoString *)&string), (*(GoString *)&pub[i])));
    GoString bitcoinAddr;
    SKY_cipher_BitcoinAddressFromPubkey(&pubkey, (GoString_ *)&bitcoinAddr);
    cr_assert(eq(type(GoString), addr[i], bitcoinAddr));
  }
}

Test(cipher_address, TestAddressBulk) {

  unsigned char buff[50];
  GoSlice slice = {buff, 0, 50};

  for (int i = 0; i < 1024; ++i) {
    randBytes(&slice, 32);
    cipher__PubKey pubkey;
    cipher__SecKey seckey;
    SKY_cipher_GenerateDeterministicKeyPair(slice, &pubkey, &seckey);
    cipher__Address addr;
    SKY_cipher_AddressFromPubKey(&pubkey, &addr);
    unsigned int err;
    err = SKY_cipher_Address_Verify(&addr, &pubkey);
    cr_assert(err == SKY_OK);
    GoString strAddr;
    SKY_cipher_Address_String(&addr, (GoString_ *)&strAddr);
    registerMemCleanup((void *)strAddr.p);
    cipher__Address addr2;

    err = SKY_cipher_DecodeBase58Address(strAddr, &addr2);
    cr_assert(err == SKY_OK);
    cr_assert(eq(type(cipher__Address), addr, addr2));
  }
}

Test(cipher_address, TestBitcoinAddressFromBytes) {
  cipher__PubKey p;
  cipher__SecKey s;
  GoInt result;
  result = SKY_cipher_GenerateKeyPair(&p, &s);
  cr_assert(result == SKY_OK, "SKY_cipher_GenerateKeyPair failed");
  cipher__Address a;
  result = SKY_cipher_AddressFromPubKey(&p, &a);
  cr_assert(result == SKY_OK, "SKY_cipher_AddressFromPubKey failed");
  cipher__Address a2;
  cipher__PubKeySlice pk = {NULL, 0, 0};
  result = SKY_cipher_Address_BitcoinBytes(&a, &pk);
  cr_assert(result == SKY_OK, "SKY_cipher_Address_BitcoinBytes failed");
  registerMemCleanup(pk.data);
  GoSlice pk_convert = {pk.data, pk.len, pk.cap};
  result = SKY_cipher_BitcoinAddressFromBytes(pk_convert, &a2);
  cr_assert(result == SKY_OK);
  cr_assert(eq(type(cipher__Address), a2, a));

  cipher__PubKeySlice b = {NULL, 0, 0};
  result = SKY_cipher_Address_BitcoinBytes(&a, &b);
  cr_assert(result == SKY_OK, "SKY_cipher_Address_BitcoinBytes");
  registerMemCleanup(b.data);

  GoInt_ b_len = b.len;

  // Invalid number of bytes
  b.len = b.len - 2;
  cipher__Address addr2;
  GoSlice b_convert = {b.data, b.len, b.cap};
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(b_convert, &addr2) ==
                SKY_ErrAddressInvalidLength,
            "Invalid address length");

  // Invalid checksum
  b_convert.len = b_len;
  (((char *)b_convert.data)[b_convert.len - 1])++;
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(b_convert, &addr2) ==
                SKY_ErrAddressInvalidChecksum,
            "Invalid checksum");

  result = SKY_cipher_AddressFromPubKey(&p, &a);
  cr_assert(result == SKY_OK, "SKY_cipher_AddressFromPubKey failed");
  a.Version = 2;
  char *buffer[1024];
  cipher__PubKeySlice b1 = {buffer, 0, 1024};
  result = SKY_cipher_Address_BitcoinBytes(&a, &b1);
  cr_assert(result == SKY_OK, "SKY_cipher_Address_BitcoinBytes failed");
  GoSlice b1_convert = {b1.data, b1.len, b1.cap};
  result = SKY_cipher_BitcoinAddressFromBytes(b1_convert, &addr2);
  cr_assert(result == SKY_ErrAddressInvalidVersion, "Invalid version");
}

Test(cipher_address, TestMustDecodeBase58Address, .signal = ((__linux__) ? SIGABRT : 2)) {

  cipher__PubKey p;
  cipher__SecKey s;
  GoInt result;
  result = SKY_cipher_GenerateKeyPair(&p, &s);
  cr_assert(result == SKY_OK, "SKY_cipher_GenerateKeyPair failed");
  cipher__Address a;
  result = SKY_cipher_AddressFromPubKey(&p, &a);
  cr_assert(result == SKY_OK, "SKY_cipher_AddressFromPubKey failed");

  result = SKY_cipher_Address_Verify(&a, &p);
  cr_assert(result == SKY_OK);
  GoString str = {"", 0};
  cipher__Address addr;
  result = SKY_cipher_MustDecodeBase58Address(str, &addr);
  cr_assert(result == SKY_ERROR);
  str.p = "cascs";
  str.n = 5;
  result = SKY_cipher_MustDecodeBase58Address(str, &addr);
  cr_assert(result == SKY_ERROR);

  char *buff_pks[1024];
  cipher__PubKeySlice b = {buff_pks, 0, 1024};
  result = SKY_cipher_Address_Bytes(&a, &b);
  cr_assert(result == SKY_OK, "SKY_cipher_Address_Bytes failed");
  int b_len = b.len;
  b.len = (int)(b_len / 2);
  GoSlice bConvert = {&b.data, b.len, b.cap};
  char buffer_h[1024];
  GoString_ h = {buffer_h, 0};
  result = SKY_base58_Hex2Base58String(bConvert, &h);
  cr_assert(result == SKY_OK, "SKY_base58_Hex2Base58String failed");
  char buffer_hConvert[1024];
  GoString hConvert = {buffer_hConvert, 0};
  toGoString(&h, &hConvert);
  result = SKY_cipher_MustDecodeBase58Address(hConvert, &addr);
  cr_assert(result == SKY_ERROR);

  b.len = b_len;
  GoSlice b2Convert = {b.data, b.len, b.cap};
  char buffer_h2[1024];
  GoString_ h2 = {buffer_h2, 0};
  result = SKY_base58_Hex2Base58String(b2Convert, &h2);
  cr_assert(result == SKY_OK, "SKY_base58_Hex2Base58String failed");
  char buffer_h2Convert[1024];
  GoString h2Convert = {buffer_h2, 0};
  toGoString(&h2, &h2Convert);
  result = SKY_cipher_MustDecodeBase58Address(h2Convert, &addr);
  cr_assert(result == SKY_OK);

  cipher__Address a2;

  result = SKY_cipher_MustDecodeBase58Address(h2Convert, &a2);
  cr_assert(result == SKY_OK, "SKY_cipher_MustDecodeBase58Address failed");
  cr_assert(eq(type(cipher__Address), a, a2));

  result = SKY_cipher_Address_String(&a, &h2);
  cr_assert(result == SKY_OK, "SKY_cipher_Address_String failed");
  toGoString(&h2, &h2Convert);
  result = SKY_cipher_MustDecodeBase58Address(h2Convert, &a2);
  cr_assert(result == SKY_OK, "SKY_cipher_MustDecodeBase58Address failed");
  cr_assert(eq(type(cipher__Address), a, a2));

  char strbadAddr[1024];
  char buffer_addrStr[1024];
  GoString_ addStr = {buffer_addrStr, 0};
  result = SKY_cipher_Address_String(&a, &addStr);
  cr_assert(result == SKY_OK, "SKY_cipher_Address_String failed");

  // preceding whitespace is invalid
  strcpy(strbadAddr, " ");
  strncat(strbadAddr, addStr.p, addStr.n);
  GoString badAddr = {strbadAddr, strlen(strbadAddr)};
  result = SKY_cipher_MustDecodeBase58Address(badAddr, &addr);
  cr_assert(result == SKY_ERROR);

  // preceding zeroes are invalid
  strcpy(strbadAddr, "000");
  strncat(strbadAddr, addStr.p, addStr.n);
  badAddr.p = strbadAddr;
  badAddr.n = strlen(strbadAddr);
  result = SKY_cipher_MustDecodeBase58Address(badAddr, &addr);
  cr_assert(result == SKY_ERROR);

  // trailing whitespace is invalid
  strcpy(strbadAddr, addStr.p);
  strcat(strbadAddr, " ");
  badAddr.p = strbadAddr;
  badAddr.n = strlen(strbadAddr);
  result = SKY_cipher_MustDecodeBase58Address(badAddr, &addr);
  cr_assert(result == SKY_ERROR);

  // trailing zeroes are invalid
  strcpy(strbadAddr, addStr.p);
  strcat(strbadAddr, "000");
  badAddr.p = strbadAddr;
  badAddr.n = strlen(strbadAddr);
  result = SKY_cipher_MustDecodeBase58Address(badAddr, &addr);
  cr_assert(result == SKY_ERROR);
}

Test(cipher_address, TestAddressRoundtrip) {
  cipher__PubKey p;
  cipher__SecKey s;
  GoInt result;
  result = SKY_cipher_GenerateKeyPair(&p, &s);
  cr_assert(result == SKY_OK, "SKY_cipher_GenerateKeyPair");
  cipher__Address a;
  result = SKY_cipher_AddressFromPubKey(&p, &a);
  cr_assert(result == SKY_OK, "SKY_cipher_AddressFromPubKey failed");
  char buffer_aBytes[1024];
  cipher__PubKeySlice aBytes = {buffer_aBytes, 0, 1024};
  result = SKY_cipher_Address_Bytes(&a, &aBytes);
  cr_assert(result == SKY_OK, "SKY_cipher_Address_Bytes failed");
  GoSlice aBytesSlice = {aBytes.data, aBytes.len, aBytes.cap};
  cipher__Address a2;
  result = SKY_cipher_AddressFromBytes(aBytesSlice, &a2);
  cr_assert(result == SKY_OK, "SKY_cipher_AddressFromBytes failed");

  cr_assert(eq(type(cipher__Address), a, a2));
  char buffer_aString[1024];
  char buffer_a2String[1024];
  GoString_ aString = {buffer_aString, 0};
  GoString_ a2String = {buffer_a2String, 0};
  result = SKY_cipher_Address_String(&a, &aString);
  result = SKY_cipher_Address_String(&a2, &a2String);

  cr_assert(eq(type(GoString_), a2String, aString));
}

Test(cipher_address, TestAddressNull) {
  cipher__Address a;
  memset(&a, 0, sizeof(cipher__Address));
  GoUint32 result;
  GoUint8 isNull;
  result = SKY_cipher_Address_Null(&a, &isNull);
  cr_assert(result == SKY_OK, "SKY_cipher_Address_Null");
  cr_assert(isNull == 1);

  cipher__PubKey p;
  cipher__SecKey s;

  result = SKY_cipher_GenerateKeyPair(&p, &s);
  cr_assert(result == SKY_OK, "SKY_cipher_GenerateKeyPair failed");

  result = SKY_cipher_AddressFromPubKey(&p, &a);
  cr_assert(result == SKY_OK, "SKY_cipher_AddressFromPubKey failed");
  result = SKY_cipher_Address_Null(&a, &isNull);
  cr_assert(result == SKY_OK, "SKY_cipher_Address_Null");
  cr_assert(isNull == 0);
}
