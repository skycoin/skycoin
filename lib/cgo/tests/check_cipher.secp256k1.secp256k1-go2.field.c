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

#define INHEX "813925AF112AAB8243F8CCBADE4CC7F63DF387263028DE6E679232A73A7F3C31"
#define EXPHEX "7F586430EA30F914965770F6098E492699C62EE1DF6CAFFA77681C179FDF3117"

TestSuite(cipher_secp256k1_field, .init = setup, .fini = teardown);

Test(cipher_secp256k1_field, TestFieldInv){
	secp256k1go__Field in;
	secp256k1go__Field out;
	secp256k1go__Field expected;
	
	memset(&in, 0, sizeof(secp256k1go__Field));
	memset(&out, 0, sizeof(secp256k1go__Field));
	memset(&expected, 0, sizeof(secp256k1go__Field));
	
	GoUint32 error_code;
	GoInt8 equal = 0;
	
	GoString InStr = {INHEX, strlen(INHEX)};
	GoString ExpStr = {EXPHEX, strlen(EXPHEX)};
	error_code = SKY_secp256k1go_Field_SetHex(&in, InStr);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	error_code = SKY_secp256k1go_Field_SetHex(&expected, ExpStr);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	error_code = SKY_secp256k1go_Field_Inv(&in, &out);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_Inv failed");
	error_code = SKY_secp256k1go_Field_Equals(&out, &expected, &equal);
	cr_assert(error_code == SKY_OK && equal, "SKY_secp256k1go_Field_Inv failed, result is different than expected.");
}