
#include <criterion/criterion.h>
#include "libskycoin.h"

#define SKYCOIN_ADDRESS_VALID "12345678"
#define SKYCOIN_ADDRESS_WRONG "12345678"

Test(cipher, test_address_valid) {
  GoString strAddr = {
    SKYCOIN_ADDRESS_VALID,
    8
  };
  Address addr;

  int r = DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 1);
}

Test(cipher, test_address_wrong) {
  GoString strAddr = {
    SKYCOIN_ADDRESS_VALID,
    8
  };
  Address addr;

  int r = DecodeBase58Address(strAddr, &addr);
  cr_assert(r == 0);
}

