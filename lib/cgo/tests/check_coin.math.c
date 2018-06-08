#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

TestSuite(coin_math, .init = setup, .fini = teardown);

Test(coin_math, TestAddUint64){
    int result;
    GoUint64 r;
    result = SKY_coin_AddUint64(10, 11, &r);
    cr_assert( result == SKY_OK );
    cr_assert(r == 21);
    GoUint64 maxUint64 = 0xFFFFFFFFFFFFFFFF;
    GoUint64 one = 1;
    result = SKY_coin_AddUint64(maxUint64, one, &r);
    cr_assert( result != SKY_OK );
}

typedef struct{
  GoUint64  a;
  GoInt64  b;
  int       failure;
} math_tests;

Test(coin_math, TestUint64ToInt64){
    int result;
    GoInt64 r;
    GoUint64 maxUint64 = 0xFFFFFFFFFFFFFFFF;
    GoInt64 maxInt64 = 0x7FFFFFFFFFFFFFFF;

    math_tests tests[] = {
        {0, 0, 0},
        {1, 1, 0},
        {maxInt64, maxInt64, 0},
        {maxUint64, 0, 1},
        //This is reset to zero in C, and it doesn't fail
        //{maxUint64 + 1, 0, 1},
    };
    int tests_count = sizeof(tests) / sizeof(math_tests);
    for(int i = 0; i < tests_count; i++){
        result = SKY_coin_Uint64ToInt64(tests[i].a, &r);
        if( tests[i].failure ){
          cr_assert(result != SKY_OK, "Failed test # %d", i + 1);
        } else {
          cr_assert(result == SKY_OK, "Failed test # %d", i + 1);
          cr_assert( tests[i].b == r );
        }
    }
}
