
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"
#include "skynumber.h"

#define BUFFER_SIZE 1024

#define R1   "6028b9e3a31c9e725fcbd7d5d16736aaaafcc9bf157dfb4be62bcbcf0969d488"
#define S1 	 "036d4a36fa235b8f9f815aa6f5457a607f956a71a035bf0970d8578bf218bb5a"
#define MSG1 "9cff3da1a4f86caf3683f865232c64992b5ed002af42b321b8d8a48420680487"
#define X1   "56dc5df245955302893d8dda0677cc9865d8011bc678c7803a18b5f6faafec08"
#define Y1 	 "54b5fbdcd8fac6468dac2de88fadce6414f5f3afbb103753e25161bef77705a6"

#define R2   "b470e02f834a3aaafa27bd2b49e07269e962a51410f364e9e195c31351a05e50"
#define S2 	 "560978aed76de9d5d781f87ed2068832ed545f2b21bf040654a2daff694c8b09"
#define MSG2 "9ce428d58e8e4caf619dc6fc7b2c2c28f0561654d1f80f322c038ad5e67ff8a6"
#define X2   "15b7e7d00f024bffcd2e47524bb7b7d3a6b251e23a3a43191ed7f0a418d9a578"
#define Y2 	 "bf29a25e2d1f32c5afb18b41ae60112723278a8af31275965a6ec1d95334e840"

TestSuite(cipher_secp256k1_sig, .init = setup, .fini = teardown);

Test(cipher_secp256k1_sig, TestSigRecover){
	GoUint32 error_code;
	Signature sig;
	Number msg;
	secp256k1go__XY pubKey;
	secp256k1go__XY expected;
	
	memset(&pubKey, 0, sizeof(secp256k1go__XY));
	memset(&expected, 0, sizeof(secp256k1go__XY));
	memset(&sig, 0, sizeof(Signature));
	memset(&msg, 0, sizeof(Number));
	
	GoString R = {R1, strlen(R1)};
	GoString S = {S1, strlen(S1)};
	GoString MSG = {MSG1, strlen(MSG1)};
	GoString X = {X1, strlen(X1)};
	GoString Y = {Y1, strlen(Y1)};
	GoInt rid = 0;
	GoInt8 result;
	
	error_code = SKY_secp256k1go_Number_SetHex(&sig.R, R);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for R failed");
	error_code = SKY_secp256k1go_Number_SetHex(&sig.S, S);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for S failed");
	error_code = SKY_secp256k1go_Number_SetHex(&msg, MSG);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for MSG failed");
	error_code = SKY_secp256k1go_Field_SetHex(&expected.X, X);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for X failed");
	error_code = SKY_secp256k1go_Field_SetHex(&expected.Y, Y);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for Y failed");
	
	error_code = SKY_secp256k1go_Signature_Recover(&sig, &pubKey, &msg, rid, &result);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Signature_Recover failed");
	cr_assert(result, "SKY_secp256k1go_Signature_Recover failed");
	
	cr_assert(cr_user_secp256k1go__Field_eq(&pubKey.X, &expected.X), "SKY_secp256k1go_Signature_Recover Xs different.");
	cr_assert(cr_user_secp256k1go__Field_eq(&pubKey.Y, &expected.Y), "SKY_secp256k1go_Signature_Recover Xs different.");
	
	R.p = R2;
	R.n = strlen(R2);
	S.p = S2;
	S.n = strlen(S2);
	MSG.p = MSG2;
	MSG.n = strlen(MSG2);
	X.p = X2;
	X.n = strlen(X2);
	Y.p = Y2;
	Y.n = strlen(Y2);
	rid = 1;
	
	error_code = SKY_secp256k1go_Number_SetHex(&sig.R, R);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for R failed");
	error_code = SKY_secp256k1go_Number_SetHex(&sig.S, S);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for S failed");
	error_code = SKY_secp256k1go_Number_SetHex(&msg, MSG);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for MSG failed");
	error_code = SKY_secp256k1go_Field_SetHex(&expected.X, X);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for X failed");
	error_code = SKY_secp256k1go_Field_SetHex(&expected.Y, Y);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for Y failed");
	
	error_code = SKY_secp256k1go_Signature_Recover(&sig, &pubKey, &msg, rid, &result);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Signature_Recover failed");
	cr_assert(result, "SKY_secp256k1go_Signature_Recover failed");
	
	GoInt8 equal;
	error_code = SKY_secp256k1go_Field_Equals(&pubKey.X, &expected.X, &equal);
	cr_assert(error_code == SKY_OK && equal, "SKY_secp256k1go_Signature_Recover Xs different.");
	SKY_secp256k1go_Field_Equals(&pubKey.Y, &expected.Y, &equal);
	cr_assert(error_code == SKY_OK && equal, "SKY_secp256k1go_Signature_Recover Ys different.");
}