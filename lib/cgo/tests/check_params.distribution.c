#include <stdio.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

TestSuite(params_distribution, .init = setup, .fini = teardown);

Test(params_distribution, TestDistributionAddressArrays) {
  GoSlice all = {NULL, 0, 0};
  GoSlice unlocked = {NULL, 0, 0};
  GoSlice locked = {NULL, 0, 0};

  SKY_params_GetDistributionAddresses((GoSlice_ *) &all);
  cr_assert(all.len == 100);

  // At the time of this writing, there should be 25 addresses in the
  // unlocked pool and 75 in the locked pool.
  SKY_params_GetUnlockedDistributionAddresses((GoSlice_ *) &unlocked);
  cr_assert(unlocked.len == 25);
  SKY_params_GetLockedDistributionAddresses((GoSlice_ *) &locked);
  cr_assert(locked.len == 75);

  int i, j, k;
  GoString *iStr, *jStr, *kStr;
  bool notfound;

  for (i = 0, iStr = (GoString *) all.data; i < all.len; ++i, ++iStr) {
    // Check no duplicate address in distribution addresses
    for (j = i + 1, jStr = iStr + 1; j < all.len; ++j, ++jStr) {
      cr_assert(not(eq(type(GoString), *iStr, *jStr)));
    }
  }

  for (i = 0, iStr = (GoString *) unlocked.data; i < unlocked.len; ++i, ++iStr) {
    // Check no duplicate address in unlocked addresses
    for (j = i + 1, jStr = iStr + 1; j < unlocked.len; ++j, ++jStr) {
      cr_assert(not(eq(type(GoString), *iStr, *jStr)));
    }

    // Check unlocked address in set of all addresses
    for (k = 0, notfound = true, kStr = (GoString *) all.data; notfound && (k < all.len); ++k, ++kStr) {
      notfound = !cr_user_GoString__eq((GoString_ *) iStr, (GoString_ *) kStr);
    }
    cr_assert(not(notfound));
  }

  for (i = 0, iStr = (GoString *) locked.data; i < locked.len; ++i, ++iStr) {
    // Check no duplicate address in locked addresses
    for (j = i + 1, jStr = iStr + 1; j < locked.len; ++j, ++jStr) {
      cr_assert(not(eq(type(GoString), *iStr, *jStr)));
    }

    // Check locked address in set of all addresses
    for (k = 0, notfound = true, kStr = (GoString *) all.data; notfound && k < all.len; ++k, ++kStr) {
      notfound = !cr_user_GoString__eq((GoString_ *) iStr, (GoString_ *) kStr);
    }
    cr_assert(not(notfound));

    // Check locked address not in set of unlocked addresses
    for (k = 0, notfound = true, kStr = (GoString *) unlocked.data; notfound && k < unlocked.len; ++k, ++kStr) {
      cr_assert(not(eq(type(GoString), *iStr, *kStr)));
    }
  }
}

