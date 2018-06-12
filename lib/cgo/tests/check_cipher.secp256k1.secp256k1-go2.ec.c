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

#define AX "0EAEBCD1DF2DF853D66CE0E1B0FDA07F67D1CABEFDE98514AAD795B86A6EA66D"
#define AY "BEB26B67D7A00E2447BAECCC8A4CEF7CD3CAD67376AC1C5785AEEBB4F6441C16"
#define AZ "0000000000000000000000000000000000000000000000000000000000000001"

#define EX "EB6752420B6BDB40A760AC26ADD7E7BBD080BF1DF6C0B009A0D310E4511BDF49"
#define EY "8E8CEB84E1502FC536FFE67967BC44314270A0B38C79865FFED5A85D138DCA6B"
#define EZ "813925AF112AAB8243F8CCBADE4CC7F63DF387263028DE6E679232A73A7F3C31"

#define U1 "B618EBA71EC03638693405C75FC1C9ABB1A74471BAAF1A3A8B9005821491C4B4"
#define U2 "8554470195DE4678B06EDE9F9286545B51FF2D9AA756CE35A39011783563EA60"

#define NONCE "9E3CD9AB0F32911BFDE39AD155F527192CE5ED1F51447D63C4F154C118DA598E"

#define E2X "02D1BF36D37ACD68E4DD00DB3A707FD176A37E42F81AEF9386924032D3428FF0"
#define E2Y "FD52E285D33EC835230EA69F89D9C38673BD5B995716A4063C893AF02F938454"
#define E2Z "4C6ACE7C8C062A1E046F66FD8E3981DC4E8E844ED856B5415C62047129268C1B"

TestSuite(cipher_secp256k1_xyz, .init = setup, .fini = teardown);

Test(cipher_secp256k1_xyz, TestXYZECMult){	

	secp256k1go__XYZ pubkey; //pubkey
	secp256k1go__XYZ pr; 	 //result of ECmult
	secp256k1go__XYZ e; 	 //expected
	Number u1, u2, nonce;
	secp256k1go__Field x, y, z;

	GoInt32 error_code;
	memset(&pubkey, 0, sizeof(secp256k1go__XYZ));
	memset(&pr, 0, sizeof(secp256k1go__XYZ));
	memset(&e, 0, sizeof(secp256k1go__XYZ));
	memset(&u1, 0, sizeof(Number));
	memset(&u2, 0, sizeof(Number));
	memset(&nonce, 0, sizeof(Number));
	memset(&x, 0, sizeof(secp256k1go__Field));
	memset(&y, 0, sizeof(secp256k1go__Field));
	memset(&z, 0, sizeof(secp256k1go__Field));
	
	GoString strAx = {AX, strlen(AX)};
	GoString strAy = {AY, strlen(AY)};
	GoString strAz = {AZ, strlen(AZ)};
	
	GoString strEx = {EX, strlen(EX)};
	GoString strEy = {EY, strlen(EY)};
	GoString strEz = {EZ, strlen(EZ)};
	
	GoString strU1 = {U1, strlen(U1)};
	GoString strU2 = {U2, strlen(U2)};
	
	error_code = SKY_secp256k1go_Field_SetHex(&pubkey.X, strAx);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	error_code = SKY_secp256k1go_Field_SetHex(&pubkey.Y, strAy);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	error_code = SKY_secp256k1go_Field_SetHex(&pubkey.Z, strAz);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	
	error_code = SKY_secp256k1go_Field_SetHex(&e.X, strEx);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	error_code = SKY_secp256k1go_Field_SetHex(&e.Y, strEy);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	error_code = SKY_secp256k1go_Field_SetHex(&e.Z, strEz);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
	
	error_code = SKY_secp256k1go_Number_SetHex(&u1, strU1);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");
	error_code = SKY_secp256k1go_Number_SetHex(&u2, strU2);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");
	
	error_code = SKY_secp256k1go_XYZ_ECmult(&pubkey, &pr, &u2, &u1);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_XYZ_ECmult failed");
	
	GoInt8 equal = 0;
	error_code = SKY_secp256k1go_XYZ_Equals(&pr, &e, &equal);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_XYZ_Equals failed.");
	cr_assert(equal, "SKY_secp256k1go_XYZ_ECmult failed, result is different than expected.");
}

Test(cipher_secp256k1_xyz, TestXYZECMultGen){	
	secp256k1go__XYZ pubkey; //pubkey
	secp256k1go__XYZ pr; 	 //result of ECmult
	secp256k1go__XYZ e; 	 //expected
	Number u1, u2, nonce;
	secp256k1go__Field x, y, z;

	GoInt32 error_code;
	memset(&pubkey, 0, sizeof(secp256k1go__XYZ));
	memset(&pr, 0, sizeof(secp256k1go__XYZ));
	memset(&e, 0, sizeof(secp256k1go__XYZ));
	memset(&u1, 0, sizeof(Number));
	memset(&u2, 0, sizeof(Number));
	memset(&nonce, 0, sizeof(Number));
	memset(&x, 0, sizeof(secp256k1go__Field));
	memset(&y, 0, sizeof(secp256k1go__Field));
	memset(&z, 0, sizeof(secp256k1go__Field));
	
	GoString strNonce = {NONCE, strlen(NONCE)};
	GoString strEx = {E2X, strlen(E2X)};
	GoString strEy = {E2Y, strlen(E2Y)};
	GoString strEz = {E2Z, strlen(E2Z)};
	
	error_code = SKY_secp256k1go_Number_SetHex(&nonce, strNonce);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex failed.");
	error_code = SKY_secp256k1go_Field_SetHex(&x, strEx);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed.");
	error_code = SKY_secp256k1go_Field_SetHex(&y, strEy);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed.");
	error_code = SKY_secp256k1go_Field_SetHex(&z, strEz);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_SetHex failed.");
	
	error_code = SKY_secp256k1go_ECmultGen(&pr, &nonce);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_ECmultGen failed.");
	error_code = SKY_secp256k1go_Field_Normalize(&pr.X);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_Normalize failed.");
	error_code = SKY_secp256k1go_Field_Normalize(&pr.Y);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_Normalize failed.");
	error_code = SKY_secp256k1go_Field_Normalize(&pr.Z);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_Normalize failed.");
	
	GoInt8 equal = 0;
	error_code = SKY_secp256k1go_Field_Equals(&pr.X, &x, &equal);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_Equals failed.");
	cr_assert(equal, "SKY_secp256k1go_ECmultGen failed. X is different than expected");
	error_code = SKY_secp256k1go_Field_Equals(&pr.Y, &y, &equal);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_Equals failed.");
	cr_assert(equal, "SKY_secp256k1go_ECmultGen failed. Y is different than expected");
	error_code = SKY_secp256k1go_Field_Equals(&pr.Z, &z, &equal);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Field_Equals failed.");
	cr_assert(equal, "SKY_secp256k1go_ECmultGen failed. Z is different than expected");
}

