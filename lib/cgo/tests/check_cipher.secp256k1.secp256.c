
#include <stdio.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

#define BUFFER_SIZE 128
#define TESTS  1
#define SigSize 65

TestSuite(cipher_secp256k1, .init = setup, .fini = teardown);

Test(cipher_secp256k1,Test_Secp256_00){
	unsigned char buff[SigSize];
	visor__ReadableOutputs nonce = {buff,0,64};
	SKY_secp256k1_RandByte(32,&nonce);
	if (nonce.len != 32) cr_fatal();
}


Test(cipher_secp256k1,Test_Secp256_01){

	cipher__PubKey pubkey;
	cipher__SecKey seckey;
	SKY_cipher_GenerateKeyPair(&pubkey,&seckey);
	GoInt errorSecKey;
	char bufferSecKey[101];
	strnhex((unsigned char *)seckey, bufferSecKey, sizeof(cipher__SecKey));
	GoSlice slseckey = { bufferSecKey,sizeof(cipher__SecKey),SigSize  };
	SKY_secp256k1_VerifySeckey(slseckey,&errorSecKey);
	if (!errorSecKey) cr_fatal();

	GoInt errorPubKey;
	GoSlice slpubkey = { &pubkey,sizeof(cipher__PubKey), sizeof(cipher__PubKey) };
	SKY_secp256k1_VerifyPubkey(slpubkey,&errorPubKey);
	if (!errorPubKey) cr_fatal();
}

Test(cipher_secp256k1, TestPubkeyFromSeckey) {

	unsigned char bufferPrivkey[BUFFER_SIZE];
	unsigned char bufferDesiredPubKey[BUFFER_SIZE];
	unsigned char bufferPubKey[BUFFER_SIZE];

	const char* hexPrivkey = "f19c523315891e6e15ae0608a35eec2e00ebd6d1984cf167f46336dabd9b2de4";
	const char* hexDesiredPubKey  = "03fe43d0c2c3daab30f9472beb5b767be020b81c7cc940ed7a7e910f0c1d9feef1";

	int sizePrivkey = hexnstr(hexPrivkey, bufferPrivkey, BUFFER_SIZE);
	int sizeDesiredPubKey = hexnstr(hexDesiredPubKey, bufferDesiredPubKey, BUFFER_SIZE);

	GoSlice privkey = { bufferPrivkey,sizePrivkey,BUFFER_SIZE };
	GoSlice_ desiredPubKey = { bufferDesiredPubKey,sizeDesiredPubKey,BUFFER_SIZE };


	visor__ReadableOutputs pubkey = {bufferPubKey,0,BUFFER_SIZE};

	GoUint32 errocode = SKY_secp256k1_PubkeyFromSeckey(privkey,&pubkey);
	if(errocode) cr_fatal();

	cr_assert(eq(type(GoSlice_),pubkey,desiredPubKey));

}

Test(cipher_secp256k1, Test_UncompressedPubkeyFromSeckey) {

	unsigned char bufferPrivkey[BUFFER_SIZE];
	unsigned char bufferDesiredPubKey[BUFFER_SIZE];
	unsigned char bufferPubKey[BUFFER_SIZE];

	const char* hexPrivkey = "f19c523315891e6e15ae0608a35eec2e00ebd6d1984cf167f46336dabd9b2de4";
	const char* hexDesiredPubKey  = "04fe43d0c2c3daab30f9472beb5b767be020b81c7cc940ed7a7e910f0c1d9feef10fe85eb3ce193405c2dd8453b7aeb6c1752361efdbf4f52ea8bf8f304aab37ab";

	int sizePrivkey = hexnstr(hexPrivkey, bufferPrivkey, BUFFER_SIZE);
	int sizeDesiredPubKey = hexnstr(hexDesiredPubKey, bufferDesiredPubKey, BUFFER_SIZE);

	GoSlice privkey = { bufferPrivkey,sizePrivkey,BUFFER_SIZE };
	GoSlice_ desiredPubKey = { bufferDesiredPubKey,sizeDesiredPubKey,BUFFER_SIZE };


	visor__ReadableOutputs pubkey = {bufferPubKey,0,BUFFER_SIZE};

	GoUint32 errocode = SKY_secp256k1_UncompressedPubkeyFromSeckey(privkey,&pubkey);
	if(errocode) cr_fatal();

	cr_assert(eq(type(GoSlice_),pubkey,desiredPubKey));

}

Test(cipher_secp256k1, Test_SignatureVerifyPubkey){
	unsigned char buff[SigSize];
	char sigBuffer[BUFFER_SIZE];
	cipher__PubKey pubkey;
	cipher__SecKey seckey;
	cipher__PubKey recoveredPubkey;
	GoInt32 error_code;
	GoSlice secKeySlice = {seckey, sizeof(cipher__SecKey), sizeof(cipher__SecKey)};
	GoSlice pubKeySlice = {pubkey, sizeof(cipher__PubKey), sizeof(cipher__PubKey)};

	error_code = SKY_secp256k1_GenerateKeyPair((coin__UxArray*)&pubKeySlice, (coin__UxArray*)&secKeySlice);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateKeyPair failed");

	GoSlice msg = {buff, 0, SigSize};
	SKY_secp256k1_RandByte(32, (coin__UxArray*)&msg);

	GoSlice recoveredPubKeySlice = {recoveredPubkey, 0, sizeof(cipher__PubKey)};
	GoSlice sig = {sigBuffer, 0, BUFFER_SIZE };
	SKY_secp256k1_Sign(msg, secKeySlice, (GoSlice_*)&sig);
	GoInt result = 0;
	error_code = SKY_secp256k1_VerifyPubkey(pubKeySlice, &result);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_VerifyPubkey failed");
	cr_assert(result == 1, "Public key not verified");
	SKY_secp256k1_RecoverPubkey(msg, sig, (coin__UxArray*)&recoveredPubKeySlice);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_RecoverPubkey failed");
	cr_assert(eq(type(GoSlice), recoveredPubKeySlice, pubKeySlice));
}

Test(cipher_secp256k1, Test_verify_functions){
	unsigned char buff[SigSize];
	char sigBuffer[BUFFER_SIZE];
	cipher__PubKey pubkey;
	cipher__SecKey seckey;
	cipher__PubKey recoveredPubkey;
	GoInt32 error_code;
	GoSlice secKeySlice = {seckey, sizeof(cipher__SecKey), sizeof(cipher__SecKey)};
	GoSlice pubKeySlice = {pubkey, sizeof(cipher__PubKey), sizeof(cipher__PubKey)};

	error_code = SKY_secp256k1_GenerateKeyPair((coin__UxArray*)&pubKeySlice, (coin__UxArray*)&secKeySlice);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateKeyPair failed");

	GoSlice msg = {buff, 0, SigSize};
	SKY_secp256k1_RandByte(32, (coin__UxArray*)&msg);

	GoSlice sig = {sigBuffer, 0, BUFFER_SIZE };
	SKY_secp256k1_Sign(msg, secKeySlice, (GoSlice_*)&sig);
	GoInt result = 0;

	error_code = SKY_secp256k1_VerifySeckey(secKeySlice, &result);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_VerifySeckey failed");
	cr_assert(result == 1, "Sec key not verified");

	error_code = SKY_secp256k1_VerifyPubkey(pubKeySlice, &result);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_VerifyPubkey failed");
	cr_assert(result == 1, "Public key not verified");

	error_code = SKY_secp256k1_VerifySignature(msg, sig, pubKeySlice, &result);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_VerifySignature failed");
	cr_assert(result == 1, "Signature not verified");
}

Test(cipher_secp256k1,Test_SignatureVerifySecKey ){
	cipher__PubKey pubkey;
	cipher__SecKey seckey;
	SKY_cipher_GenerateKeyPair(&pubkey,&seckey);
	GoInt errorSecKey;
	char bufferSecKey[101];
	strnhex((unsigned char *)seckey, bufferSecKey, sizeof(cipher__SecKey));
	GoSlice slseckey = { bufferSecKey,sizeof(cipher__SecKey),SigSize  };
	SKY_secp256k1_VerifySeckey(slseckey,&errorSecKey);
	cr_assert(errorSecKey != SKY_OK);
	GoInt errorPubKey;
	GoSlice slpubkey = { &pubkey,sizeof(cipher__PubKey), sizeof(cipher__PubKey) };
	SKY_secp256k1_VerifyPubkey(slpubkey,&errorPubKey);
	cr_assert(errorPubKey != SKY_OK);
}
