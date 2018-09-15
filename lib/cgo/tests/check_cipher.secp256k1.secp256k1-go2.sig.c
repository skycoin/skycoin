
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

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

#define forceLowS true

TestSuite(cipher_secp256k1_sig, .init = setup, .fini = teardown);

Test(cipher_secp256k1_sig, TestSigRecover){
	GoUint32 error_code;
	Signature_Handle sig;
	Number_Handle msg;
	secp256k1go__XY pubKey;
	secp256k1go__XY expected;

	memset(&pubKey, 0, sizeof(secp256k1go__XY));
	memset(&expected, 0, sizeof(secp256k1go__XY));
	sig = 0;
	msg = 0;

	GoString R = {R1, strlen(R1)};
	GoString S = {S1, strlen(S1)};
	GoString MSG = {MSG1, strlen(MSG1)};
	GoString X = {X1, strlen(X1)};
	GoString Y = {Y1, strlen(Y1)};
	GoInt rid = 0;
	GoInt8 result;

	error_code = SKY_secp256k1go_Signature_Create(&sig);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Signature_Create failed");
	registerHandleClose(sig);
	Number_Handle r;
	error_code = SKY_secp256k1go_Signature_GetR(sig, &r);
	registerHandleClose(r);
	Number_Handle s;
	error_code = SKY_secp256k1go_Signature_GetS(sig, &s);
	registerHandleClose(s);
	error_code = SKY_secp256k1go_Number_Create(&msg);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_Create failed");
	registerHandleClose(msg);

	error_code = SKY_secp256k1go_Number_SetHex(r, R);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for R failed");
	error_code = SKY_secp256k1go_Number_SetHex(s, S);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for S failed");
	error_code = SKY_secp256k1go_Number_SetHex(msg, MSG);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for MSG failed");
	error_code = SKY_secp256k1go_Field_SetHex(&expected.X, X);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for X failed");
	error_code = SKY_secp256k1go_Field_SetHex(&expected.Y, Y);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for Y failed");

	error_code = SKY_secp256k1go_Signature_Recover(sig, &pubKey, msg, rid, &result);
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

	error_code = SKY_secp256k1go_Number_SetHex(r, R);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for R failed");
	error_code = SKY_secp256k1go_Number_SetHex(s, S);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for S failed");
	error_code = SKY_secp256k1go_Number_SetHex(msg, MSG);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for MSG failed");
	error_code = SKY_secp256k1go_Field_SetHex(&expected.X, X);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for X failed");
	error_code = SKY_secp256k1go_Field_SetHex(&expected.Y, Y);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Number_SetHex for Y failed");

	error_code = SKY_secp256k1go_Signature_Recover(sig, &pubKey, msg, rid, &result);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1go_Signature_Recover failed");
	cr_assert(result, "SKY_secp256k1go_Signature_Recover failed");

	GoInt8 equal;
	error_code = SKY_secp256k1go_Field_Equals(&pubKey.X, &expected.X, &equal);
	cr_assert(error_code == SKY_OK && equal, "SKY_secp256k1go_Signature_Recover Xs different.");
	SKY_secp256k1go_Field_Equals(&pubKey.Y, &expected.Y, &equal);
	cr_assert(error_code == SKY_OK && equal, "SKY_secp256k1go_Signature_Recover Ys different.");
}

Test(cipher_secp256k1_sig, TestSigVerify) {

  Number_Handle msg;
  Signature_Handle sig;
  secp256k1go__XY key;

  msg = 0;
  sig = 0;
  memset(&key, 0, sizeof(secp256k1go__XY));
	GoUint32 result;

	result = SKY_secp256k1go_Signature_Create(&sig);
	cr_assert(result == SKY_OK, "SKY_secp256k1go_Signature_Create failed");
	registerHandleClose(sig);
	Number_Handle r;
	result = SKY_secp256k1go_Signature_GetR(sig, &r);
	registerHandleClose(r);
	Number_Handle s;
	result = SKY_secp256k1go_Signature_GetS(sig, &s);
	registerHandleClose(s);
	result = SKY_secp256k1go_Number_Create(&msg);
	cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_Create failed");
	registerHandleClose(msg);

  GoString str = {
      "D474CBF2203C1A55A411EEC4404AF2AFB2FE942C434B23EFE46E9F04DA8433CA", 64};

  result = SKY_secp256k1go_Number_SetHex(msg, str);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");

  str.p = "98F9D784BA6C5C77BB7323D044C0FC9F2B27BAA0A5B0718FE88596CC56681980";
  result = SKY_secp256k1go_Number_SetHex(r, str);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");

  str.p = "E3599D551029336A745B9FB01566624D870780F363356CEE1425ED67D1294480";
  result = SKY_secp256k1go_Number_SetHex(s, str);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");

  str.p = "7d709f85a331813f9ae6046c56b3a42737abf4eb918b2e7afee285070e968b93";
  result = SKY_secp256k1go_Field_SetHex(&key.X, str);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");

  str.p = "26150d1a63b342986c373977b00131950cb5fc194643cad6ea36b5157eba4602";
  result = SKY_secp256k1go_Field_SetHex(&key.Y, str);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Field_SetHex failed");
  GoUint8 valid;
  result = SKY_secp256k1go_Signature_Verify(sig, &key, msg, &valid);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Signature_Verify failed");
  cr_assert(valid, "sig.Verify 1");

  str.p = "2c43a883f4edc2b66c67a7a355b9312a565bb3d33bb854af36a06669e2028377";
  result = SKY_secp256k1go_Number_SetHex(msg, str);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");

  str.p = "6b2fa9344462c958d4a674c2a42fbedf7d6159a5276eb658887e2e1b3915329b";
  result = SKY_secp256k1go_Number_SetHex(r, str);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");

  str.p = "eddc6ea7f190c14a0aa74e41519d88d2681314f011d253665f301425caf86b86";
  result = SKY_secp256k1go_Number_SetHex(s, str);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");
  char buffer_xy[1024];
  cipher__PubKeySlice xy = {buffer_xy, 0, 1024};
  str.p = "02a60d70cfba37177d8239d018185d864b2bdd0caf5e175fd4454cc006fd2d75ac";
  str.n = 66;
  result = SKY_base58_String2Hex(str, &xy);
  cr_assert(result == SKY_OK, "SKY_base58_String2Hex");
  GoSlice xyConvert = {xy.data, xy.len, xy.cap};
  result = SKY_secp256k1go_XY_ParsePubkey(&key, xyConvert, &valid);
  cr_assert(result == SKY_OK && valid, "SKY_secp256k1go_XY_ParsePubkey failed");
  result = SKY_secp256k1go_Signature_Verify(sig, &key, msg, &valid);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Signature_Verify failed");
  cr_assert(valid, "sig.Verify 2");
}

Test(cipher_secp256k1_sig, TestSigSign) {

  Number_Handle sec;
  Number_Handle msg;
  Number_Handle non;
  Signature_Handle sig;
  GoInt recid;
	GoUint32 result;

  sec = 0;
	msg = 0;
	non = 0;
	sig = 0;

	result = SKY_secp256k1go_Number_Create(&sec);
	cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_Create failed");
	registerHandleClose(msg);

	result = SKY_secp256k1go_Number_Create(&msg);
	cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_Create failed");
	registerHandleClose(msg);

	result = SKY_secp256k1go_Number_Create(&non);
	cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_Create failed");
	registerHandleClose(msg);

  GoString str = {
      "73641C99F7719F57D8F4BEB11A303AFCD190243A51CED8782CA6D3DBE014D146", 64};
  result = SKY_secp256k1go_Number_SetHex(sec, str);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");

  str.p = "D474CBF2203C1A55A411EEC4404AF2AFB2FE942C434B23EFE46E9F04DA8433CA";
  str.n = 64;
  result = SKY_secp256k1go_Number_SetHex(msg, str);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");

  str.p = "9E3CD9AB0F32911BFDE39AD155F527192CE5ED1F51447D63C4F154C118DA598E";
  str.n = 64;
  result = SKY_secp256k1go_Number_SetHex(non, str);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");

	result = SKY_secp256k1go_Signature_Create(&sig);
	cr_assert(result == SKY_OK, "SKY_secp256k1go_Signature_Create failed");
	registerHandleClose(sig);

  GoInt res;
	GoInt8 equal;

  result = SKY_secp256k1go_Signature_Sign(sig, sec, msg, non, &recid, &res);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Signature_Sign failed");
  cr_assert(res == 1, "res failed %d", res);

  if (forceLowS) {
    cr_assert(recid == 0, " recid failed %d", recid);
  } else {
    cr_assert(recid == 1, " recid failed %d", recid);
  }
  str.p = "98f9d784ba6c5c77bb7323d044c0fc9f2b27baa0a5b0718fe88596cc56681980";
  str.n = 64;
  result = SKY_secp256k1go_Number_SetHex(non, str);
  cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");

	Number_Handle r;
	result = SKY_secp256k1go_Signature_GetR(sig, &r);
	cr_assert(result == SKY_OK, "SKY_secp256k1go_Signature_GetR failed");
	registerHandleClose(r);

	equal = 0;
	result = SKY_secp256k1go_Number_IsEqual(r, non, &equal);
	cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_IsEqual failed");
  cr_assert(equal != 0);

  if (forceLowS) {
    str.p = "1ca662aaefd6cc958ba4604fea999db133a75bf34c13334dabac7124ff0cfcc1";
    str.n = 64;
    result = SKY_secp256k1go_Number_SetHex(non, str);
    cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");
  } else {
    str.p = "E3599D551029336A745B9FB01566624D870780F363356CEE1425ED67D1294480";
    str.n = 64;
    result = SKY_secp256k1go_Number_SetHex(non, str);
    cr_assert(result == SKY_OK, "SKY_secp256k1go_Number_SetHex failed");
  }
	Number_Handle s;
	result = SKY_secp256k1go_Signature_GetS(sig, &s);
	cr_assert(result == SKY_OK, "SKY_secp256k1go_Signature_GetS failed");
	registerHandleClose(s);

	equal = 0;
	result = SKY_secp256k1go_Number_IsEqual(s, non, &equal);
	cr_assert(result == SKY_OK, "SKY_secp256k1go_Signature_GetS failed");
  cr_assert(equal != 0);
}
