
#include <criterion/criterion.h>
#include "libskycoin.h"

#define SKYCOIN_ADDRESS_VALID "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"
#define SKYCOIN_ADDRESS_WRONG "12345678"

int addr_equal(Address addr1, Address addr2){
  int r = 0;
  if(addr1.Version != addr2.Version) r = 1;
  for (int i = 0; i < sizeof(Ripemd160); ++i) {
    if(addr1.Key[i] != addr2.Key[i]) r = 1;
  }
  return r;
}

Test(cipher, test_address_valid) {
  GoString strAddr = {
    SKYCOIN_ADDRESS_VALID,
    35
  };
  Address addr;

  int r = SKY_Cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 1);
}

Test(cipher, test_address_wrong) {
  GoString strAddr = {
    SKYCOIN_ADDRESS_VALID,
    8
  };
  Address addr;

  int r = SKY_Cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 0, "a1");
  strAddr.p = strcat(" ", SKYCOIN_ADDRESS_VALID);
  strAddr.n = 35;
  r = SKY_Cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 0, "a2");
  strAddr.p = strcat("000", SKYCOIN_ADDRESS_VALID);
  r = SKY_Cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 0, "a3");
  strAddr.p = strcat(SKYCOIN_ADDRESS_VALID, "000");
  r = SKY_Cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 0, "a4");
}

Test(cipher, test_address_frombytes){
  GoString strAddr = {
    SKYCOIN_ADDRESS_VALID,
    35
  };
  Address addr;
  SKY_Cipher_DecodeBase58Address(strAddr, &addr);
  unsigned char* bytes;
  SKY_Cipher_Address_Bytes(&addr, bytes);
  Address addr2;
  size_t size = sizeof(bytes);
  SKY_Cipher_BitcoinAddressFromBytes(bytes, size, &addr2);
  int r = addr_equal(addr, addr2);
  cr_assert(r == 0, "a1");
  size = sizeof(bytes) - 2;
  r = SKY_Cipher_BitcoinAddressFromBytes(bytes, size, &addr2);
  cr_assert(r == 0, "a2");
  bytes[sizeof(bytes) - 1] = '2';
  r = SKY_Cipher_BitcoinAddressFromBytes(bytes, size, &addr2);
  cr_assert(r == 0, "a2");
}



