
#include <criterion/criterion.h>
#include "libskycoin.h"
#include "libsky_util.h"
#include <stdio.h>

#define SKYCOIN_ADDRESS_VALID "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"

// buffer big enough to hold all kind of data needed by test cases
unsigned char buff[1024];

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
  cr_assert( addr_equal(addr1, addr2) == 1);
}

// TODO: Change to write assertion like this cr_assert(not(eq(type(struct Address), &addr1, &addr2)))
void cr_assert_addr_noteq(Address *addr1, Address *addr2, char *msg){
  cr_assert( addr_equal(addr1, addr2) == 0);
}

Test(asserts, TestDecodeBase58Address) {

 GoString strAddr = {
  SKYCOIN_ADDRESS_VALID,
  35
};
Address addr;

cr_assert( SKY_cipher_DecodeBase58Address(strAddr, &addr) == 1, "accept valid address");

// preceding whitespace is invalid
char *worng = join_char(" ",SKYCOIN_ADDRESS_VALID);

GoString strAddrWrong ={
  worng,
  35
};
cr_assert( SKY_cipher_DecodeBase58Address(strAddrWrong, &addr) == 0, "preceding whitespace is invalid");

// preceding zeroes are invalid
strAddrWrong.p=join_char("000",SKYCOIN_ADDRESS_VALID);
cr_assert( SKY_cipher_DecodeBase58Address(strAddrWrong, &addr) == 0, " preceding zeroes are invalid");

// trailing whitespace is invalid

strAddrWrong.p = join_char(SKYCOIN_ADDRESS_VALID," ");
cr_assert( SKY_cipher_DecodeBase58Address(strAddrWrong, &addr) == 0, " trailing whitespace is invalid");

// trailing zeroes are invalid
strAddrWrong.p = join_char(SKYCOIN_ADDRESS_VALID,"000");
cr_assert( SKY_cipher_DecodeBase58Address(strAddrWrong, &addr) == 0, " trailing zeroes are invalid");

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
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2) == 0, "convert bytes to SKY address");
  // cr_assert(eq(type(struct Address), &addr, &addr2));

  int bytes_len = bytes.len;

  bytes.len = bytes.len - 2;
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2) == 1, "no SKY address due to short bytes length");

  bytes.len = bytes_len;
  ((char *) bytes.data)[bytes.len - 1] = '2';
  cr_assert(SKY_cipher_BitcoinAddressFromBytes(bytes, &addr2) == 1, "no SKY address due to corrupted bytes");
}

// Test(cipher, TestAddressRoundtrip){
//  GoString strAddr = {
//     SKYCOIN_ADDRESS_VALID,
//     35
//   };

//   Address addr, addr2;
//   GoSlice bytes;

//   bytes.data = buff;
//   bytes.len = 0;
//   bytes.cap = sizeof(buff);

//   // a2, err := addressFromBytes(a.Bytes())
//   // require.NoError(t, err)
//   // require.Equal(t, a, a2)
//   // require.Equal(t, a.String(), a2.String())
// }

