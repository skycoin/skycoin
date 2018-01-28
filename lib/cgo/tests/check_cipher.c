
#include <criterion/criterion.h>
#include "libskycoin.h"

#define SKYCOIN_ADDRESS_VALID "12345678"
#define SKYCOIN_ADDRESS_WRONG "12345678"

Test(cipher, test_address_valid) {
  GoString strAddr = {
    SKYCOIN_ADDRESS_VALID,
    8
  };

  struct DecodeBase58Address_return r = DecodeBase58Address(strAddr);
  cr_assert(r.r1 == 1);
}

Test(cipher, test_address_wrong) {
  GoString strAddr = {
    SKYCOIN_ADDRESS_VALID,
    8
  };

  struct DecodeBase58Address_return r = DecodeBase58Address(strAddr);
  cr_assert(r.r1 == 0);
}

