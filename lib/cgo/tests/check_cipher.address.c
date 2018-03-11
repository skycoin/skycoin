
#include <criterion/criterion.h>
#include "libskycoin.h"

#define SKYCOIN_ADDRESS_VALID "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"
#define SKYCOIN_ADDRESS_WRONG_1 "12345678"
#define SKYCOIN_ADDRESS_WRONG_2 " 2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"
#define SKYCOIN_ADDRESS_WRONG_3 "0002GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"
#define SKYCOIN_ADDRESS_WRONG_4 "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv000"
#define SKYCOIN_ADDRESS_WRONG_5 "abc2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qvdef"

int addr_equal(Address *addr1, Address *addr2){
  if(addr1->Version != addr2->Version)
    return 0;
  for (int i = 0; i < sizeof(Ripemd160); ++i) {
    if(addr1->Key[i] != addr2->Key[i])
      return 0;
  }
  return 1;
}

// TODO: Change to write assertion like this cr_assert(eq(type(struct Address), &addr1, &addr2))
void cr_assert_addr_eq(Address *addr1, Address *addr2, char *msg){
  int r = addr_equal(addr1, addr2);
  cr_assert(r == 1);
}

// TODO: Change to write assertion like this cr_assert(not(eq(type(struct Address), &addr1, &addr2)))
void cr_assert_addr_noteq(Address *addr1, Address *addr2, char *msg){
  int r = addr_equal(addr1, addr2);
  cr_assert(r == 0);
}

Test(cipher, test_address_valid) {
  GoString strAddr = {
    SKYCOIN_ADDRESS_VALID,
    35
  };
  Address addr;

  int r = SKY_cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 1, "accept valid address");

  strAddr.p = SKYCOIN_ADDRESS_WRONG_4;
  r = SKY_cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 1, "accept address with suffix and exact len");
}

Test(cipher, test_address_wrong) {
  GoString strAddr = {
    SKYCOIN_ADDRESS_VALID,
    8
  };
  Address addr;

  int r = SKY_cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 0, "reject shorter strings");

  strAddr.p = SKYCOIN_ADDRESS_WRONG_2;
  strAddr.n = 35;
  r = SKY_cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 0, "reject leading whitespaces");

  strAddr.p = SKYCOIN_ADDRESS_WRONG_3;
  r = SKY_cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 0, "reject unexpected prefix");

  strAddr.p = SKYCOIN_ADDRESS_WRONG_4;
  strAddr.n = 38;
  r = SKY_cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 0, "reject unexpected suffix");

  strAddr.p = SKYCOIN_ADDRESS_WRONG_5;
  strAddr.n = 41;
  r = SKY_cipher_DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 0, "reject unexpected prefix and suffix");
}

Test(cipher, test_address_frombytes){
  GoString strAddr = {
    SKYCOIN_ADDRESS_VALID,
    35
  };
  Address addr, addr2;
  GoSlice bytes;

  SKY_cipher_DecodeBase58Address(strAddr, &addr);
  SKY_cipher_Address_Bytes(&addr, (GoSlice_ *)&bytes);
  int r = SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2);
  cr_assert(r == 1, "convert bytes to SKY address");
  cr_assert_addr_eq(&addr, &addr2, "address from bytes should match original");

  int bytes_len = bytes.len;

  bytes.len = bytes.len - 2;
  r = SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2);
  cr_assert(r == 0, "no SKY address due to short bytes length");

  bytes.len = bytes_len;
  ((char *) bytes.data)[bytes.len - 1] = '2';
  r = SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2);
  cr_assert(r == 0, "no SKY address due to corrupted bytes");
}

