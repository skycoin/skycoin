#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"
#include "base64.h"

#define AX "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
#define AY "483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"
#define AZ "01"
#define EX "7D152C041EA8E1DC2191843D1FA9DB55B68F88FEF695E2C791D40444B365AFC2"
#define EY "56915849F52CC8F76F5FD7E4BF60DB4A43BF633E1B1383F85FE89164BFADCBDB"
#define EZ "9075B4EE4D4788CABB49F7F81C221151FA2F68914D0AA833388FA11FF621A970"



Test(cipher_secp256k1_xyz, TestXYZDouble){	
	GoInt32 error_code;
	secp256k1go__XYZ a; //sample data
	secp256k1go__XYZ r; //result of double
	secp256k1go__XYZ e; //expected
	
	memset(&a, 0, sizeof(secp256k1go__XYZ));
	memset(&e, 0, sizeof(secp256k1go__XYZ));
	memset(&r, 0, sizeof(secp256k1go__XYZ));
	
	GoString strAx = {AX, strlen(AX)};
	GoString strAy = {AY, strlen(AY)};
	GoString strAz = {AZ, strlen(AZ)};
	
	GoString strEx = {EX, strlen(EX)};
	GoString strEy = {EY, strlen(EY)};
	GoString strEz = {EZ, strlen(EZ)};
	
	error_code = SKY_secp256k1go_Field_SetHex(&a.X, strAx);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	error_code = SKY_secp256k1go_Field_SetHex(&a.Y, strAy);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	error_code = SKY_secp256k1go_Field_SetHex(&a.Z, strAz);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	
	error_code = SKY_secp256k1go_Field_SetHex(&e.X, strEx);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	error_code = SKY_secp256k1go_Field_SetHex(&e.Y, strEy);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	error_code = SKY_secp256k1go_Field_SetHex(&e.Z, strEz);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	
	error_code = SKY_secp256k1go_XYZ_Double(&a, &r);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_XYZ_Double failed");
	
	GoInt8 equal = 0;
	error_code = SKY_secp256k1go_XYZ_Equals(&r, &e, &equal);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_XYZ_Equals failed.");
	cr_assert(equal, "SKY_secp256k1go_XYZ_Double failed, result is different than expected.");
}

