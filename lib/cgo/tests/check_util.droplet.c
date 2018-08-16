#include <criterion/criterion.h>
#include <criterion/new/assert.h>
#include <signal.h>
#include <stdio.h>
#include <string.h>

#include "libskycoin.h"
#include "skycriterion.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"
#include "transutil.h"

TestSuite(util_droplet, .init = setup, .fini = teardown);
#define BUFFER_SIZE 1024
Test(util_droplet, TestFromString)
{

    typedef struct
    {
        GoString s;
        GoInt64 n;
        GoInt64 e;
    } tmpstruct;

    tmpstruct cases[BUFFER_SIZE];

    cases[0].s.p = "0";
    cases[0].s.n = 1;
    cases[0].n = 0;
    cases[0].e = SKY_OK;

    cases[1].s.p = "0.";
    cases[1].s.n = 2;
    cases[1].n = 0;
    cases[1].e = SKY_OK;

    cases[2].s.p = "0.0";
    cases[2].s.n = 3;
    cases[2].n = 0;
    cases[2].e = SKY_OK;

    cases[3].s.p = "0.000000";
    cases[3].s.n = 8;
    cases[3].n = 0;
    cases[3].e = SKY_OK;

    cases[4].s.p = "0.0000000";
    cases[4].s.n = 9;
    cases[4].n = 0;
    cases[4].e = SKY_OK;

    cases[5].s.p = "0.0000001";
    cases[5].s.n = 9;
    cases[5].n = 0;
    cases[5].e = SKY_ErrTooManyDecimals;

    cases[6].s.p = "0.000001";
    cases[6].s.n = 8;
    cases[6].n = 1;
    cases[6].e = SKY_OK;

    cases[7].s.p = "0.0000010";
    cases[7].s.n = 9;
    cases[7].n = 1;
    cases[7].e = SKY_OK;

    cases[8].s.p = "1";
    cases[8].s.n = 1;
    cases[8].n = 1000000;
    cases[8].e = SKY_OK;

    cases[9].s.p = "1.000001";
    cases[9].s.n = 8;
    cases[9].n = 1000001;
    cases[9].e = SKY_OK;

    cases[10].s.p = "-1";
    cases[10].s.n = 2;
    cases[10].n = 0;
    cases[10].e = SKY_ErrNegativeValue;

    cases[11].s.p = "10000";
    cases[11].s.n = 5;
    cases[11].n = 10000000000;
    cases[11].e = SKY_OK;

    cases[12].s.p = "123456789.123456";
    cases[12].s.n = 16;
    cases[12].n = 123456789123456;
    cases[12].e = SKY_OK;

    cases[13].s.p = "123.000456";
    cases[13].s.n = 10;
    cases[13].n = 123000456;
    cases[13].e = SKY_OK;

    cases[14].s.p = "100SKY";
    cases[14].s.n = 8;
    cases[14].n = 0;
    cases[14].e = SKY_ERROR;

    cases[15].s.p = "";
    cases[15].s.n = 0;
    cases[15].n = 0;
    cases[15].e = SKY_ERROR;

    cases[16].s.p = "999999999999999999999999999999999999999999";
    cases[16].s.n = 42;
    cases[16].n = 0;
    cases[16].e = SKY_ErrTooLarge;

    cases[17].s.p = "9223372036854.775807";
    cases[17].s.n = 20;
    cases[17].n = 9223372036854775807;
    cases[17].e = SKY_OK;

    cases[18].s.p = "-9223372036854.775807";
    cases[18].s.n = 21;
    cases[18].n = 0;
    cases[18].e = SKY_ErrNegativeValue;

    cases[19].s.p = "9223372036854775808";
    cases[19].s.n = 19;
    cases[19].n = 0;
    cases[19].e = SKY_ErrTooLarge;

    cases[20].s.p = "9223372036854775807.000001";
    cases[20].s.n = 26;
    cases[20].n = 0;
    cases[20].e = SKY_ErrTooLarge;

    cases[21].s.p = "9223372036854775807";
    cases[21].s.n = 19;
    cases[21].n = 0;
    cases[21].e = SKY_ErrTooLarge;

    cases[22].s.p = "9223372036854775806.000001";
    cases[22].s.n = 26;
    cases[22].n = 0;
    cases[22].e = SKY_ErrTooLarge;

    cases[23].s.p = "1.1";
    cases[23].s.n = 3;
    cases[23].n = 1100000;
    cases[23].e = SKY_OK;

    cases[24].s.p = "1.01";
    cases[24].s.n = 4;
    cases[24].n = 1010000;
    cases[24].e = SKY_OK;

    cases[25].s.p = "1.001";
    cases[25].s.n = 5;
    cases[25].n = 1001000;
    cases[25].e = SKY_OK;

    cases[26].s.p = "1.0001";
    cases[26].s.n = 6;
    cases[26].n = 1000100;
    cases[26].e = SKY_OK;

    cases[27].s.p = "1.00001";
    cases[27].s.n = 7;
    cases[27].n = 1000010;
    cases[27].e = SKY_OK;

    cases[28].s.p = "1.000001";
    cases[28].s.n = 8;
    cases[28].n = 1000001;
    cases[28].e = SKY_OK;

    cases[29].s.p = "1.0000001";
    cases[29].s.n = 9;
    cases[29].n = 0;
    cases[29].e = SKY_ErrTooManyDecimals;

    int len = 30;
    for (int i = 0; i < len; i++)
    {
        tmpstruct tc = cases[i];
        GoUint64 n;
        int err = SKY_droplet_FromString(tc.s, &n);

        if (tc.e == SKY_OK)
        {
            cr_assert(err == SKY_OK, "SKY_droplet_FromString %d in iter %d and %d", err, i, len);
            cr_assert(tc.n == n, , "result %d in interation %d", n, i);
        }
        else
        {
            cr_assert(err != SKY_OK, "SKY_droplet_FromString %d in iter %d and %d",
                      err, i, len);
            cr_assert(err == tc.e, "Not equal %X != %X in iteration %d", err, tc.e,
                      i);
            // cr_assert(0 == n, "result %d != 0 in iteration %d", n, i);
        }
    }
}
