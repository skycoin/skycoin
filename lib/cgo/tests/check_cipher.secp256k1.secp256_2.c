
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

#define BUFFER_SIZE 128
#define TESTS  1
#define SigSize 65

int keys_count = 4;
const char* test_keys[] = {
"08efb79385c9a8b0d1c6f5f6511be0c6f6c2902963d874a3a4bacc18802528d3",
"78298d9ecdc0640c9ae6883201a53f4518055442642024d23c45858f45d0c3e6",
"04e04fe65bfa6ded50a12769a3bd83d7351b2dbff08c9bac14662b23a3294b9e",
"2f5141f1b75747996c5de77c911dae062d16ae48799052c04ead20ccd5afa113",
};

Test(cipher_secp256k1, Test_ECDH) {
	cipher__PubKey pubkey1;
	cipher__SecKey seckey1;
	cipher__PubKey pubkey2;
	cipher__SecKey seckey2;
	unsigned char bufferECDH1[BUFFER_SIZE];
	unsigned char bufferECDH2[BUFFER_SIZE];
	
	GoInt32 error_code;
	GoSlice secKeySlice1 = {seckey1, sizeof(cipher__SecKey), sizeof(cipher__SecKey)};
	GoSlice pubKeySlice1 = {pubkey1, sizeof(cipher__PubKey), sizeof(cipher__PubKey)};
	GoSlice secKeySlice2 = {seckey2, sizeof(cipher__SecKey), sizeof(cipher__SecKey)};
	GoSlice pubKeySlice2 = {pubkey2, sizeof(cipher__PubKey), sizeof(cipher__PubKey)};
	GoSlice ecdh1 = { bufferECDH1, 0, BUFFER_SIZE };
	GoSlice ecdh2 = { bufferECDH2, 0, BUFFER_SIZE };
	
	error_code = SKY_secp256k1_GenerateKeyPair(
		(coin__UxArray*)&pubKeySlice1, (coin__UxArray*)&secKeySlice1);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateKeyPair failed");
	error_code = SKY_secp256k1_GenerateKeyPair(
		(coin__UxArray*)&pubKeySlice2, (coin__UxArray*)&secKeySlice2);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateKeyPair failed");
	
	SKY_secp256k1_ECDH(pubKeySlice1, secKeySlice2, (coin__UxArray*)&ecdh1);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_ECDH failed.");
	SKY_secp256k1_ECDH(pubKeySlice2, secKeySlice1, (coin__UxArray*)&ecdh2);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_ECDH failed.");
	
	cr_assert(eq(type(GoSlice), ecdh1, ecdh2));
}

Test(cipher_secp256k1, Test_ECDH2) {
	cipher__PubKey pubkey1;
	cipher__SecKey seckey1;
	cipher__PubKey pubkey2;
	cipher__SecKey seckey2;
	unsigned char bufferECDH1[BUFFER_SIZE];
	unsigned char bufferECDH2[BUFFER_SIZE];
	
	GoInt32 error_code;
	GoSlice secKeySlice1 = {seckey1, sizeof(cipher__SecKey), sizeof(cipher__SecKey)};
	GoSlice pubKeySlice1 = {pubkey1, sizeof(cipher__PubKey), sizeof(cipher__PubKey)};
	GoSlice secKeySlice2 = {seckey2, sizeof(cipher__SecKey), sizeof(cipher__SecKey)};
	GoSlice pubKeySlice2 = {pubkey2, sizeof(cipher__PubKey), sizeof(cipher__PubKey)};
	GoSlice ecdh1 = { bufferECDH1, 0, BUFFER_SIZE };
	GoSlice ecdh2 = { bufferECDH2, 0, BUFFER_SIZE };
	
	for( int i = 0; i < 32; i++ ) {
		error_code = SKY_secp256k1_GenerateKeyPair(
			(coin__UxArray*)&pubKeySlice1, (coin__UxArray*)&secKeySlice1);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateKeyPair failed");
		error_code = SKY_secp256k1_GenerateKeyPair(
			(coin__UxArray*)&pubKeySlice2, (coin__UxArray*)&secKeySlice2);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateKeyPair failed");
		
		SKY_secp256k1_ECDH(pubKeySlice1, secKeySlice2, (coin__UxArray*)&ecdh1);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_ECDH failed.");
		SKY_secp256k1_ECDH(pubKeySlice2, secKeySlice1, (coin__UxArray*)&ecdh2);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_ECDH failed.");
		
		cr_assert(eq(type(GoSlice), ecdh1, ecdh2));
	}
}

Test(cipher_secp256k1, Test_Abnormal_Keys) {
	char seedBuffer[64];
	GoSlice seed = {seedBuffer, 0, 64};
	unsigned char bufferPrivatekey[BUFFER_SIZE];
	unsigned char bufferPubKey[BUFFER_SIZE];
	GoSlice privatekey = { bufferPrivatekey, 0, BUFFER_SIZE };
	GoSlice pubKey = { bufferPubKey, 0, BUFFER_SIZE };
	GoInt32 error_code;
	
	for( int i = 0; i < 32; i++ ) {
		error_code = SKY_secp256k1_RandByte(32, (coin__UxArray*)&seed);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed.");
		error_code = SKY_secp256k1_GenerateDeterministicKeyPair(seed, 
			(coin__UxArray*)&privatekey, (coin__UxArray*)& pubKey);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateDeterministicKeyPair failed.");
		GoInt verified = 0;
		error_code = SKY_secp256k1_VerifyPubkey(pubKey, &verified);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_VerifyPubkey failed.");
		cr_assert(verified != 0, "Failed verifying key");
	}
}

Test(cipher_secp256k1, Test_Abnormal_Keys2) {
	unsigned char bufferPrivatekey[BUFFER_SIZE];
	unsigned char bufferPubKey[BUFFER_SIZE];
	
	GoSlice privatekey = { bufferPrivatekey, 0, BUFFER_SIZE };
	GoSlice pubKey = { bufferPubKey, 0, BUFFER_SIZE };
	GoInt32 error_code;
	
	for(int i = 0; i < keys_count; i++){
		int sizePrivatekey = hexnstr(test_keys[i], bufferPrivatekey, BUFFER_SIZE);
		privatekey.len = sizePrivatekey;
		error_code = SKY_secp256k1_PubkeyFromSeckey(privatekey, (coin__UxArray*)&pubKey);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_PubkeyFromSeckey failed.");
		cr_assert(pubKey.len > 0, "SKY_secp256k1_PubkeyFromSeckey failed.");
		GoInt verified = 0;
		error_code = SKY_secp256k1_VerifyPubkey(pubKey, &verified);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_VerifyPubkey failed.");
		cr_assert(verified != 0, "Failed verifying key");
	}
}

Test(cipher_secp256k1, Test_Abnormal_Keys3) {
	unsigned char bufferPrivatekey1[BUFFER_SIZE];
	unsigned char bufferPubKey1[BUFFER_SIZE];
	unsigned char bufferPrivatekey2[BUFFER_SIZE];
	unsigned char bufferPubKey2[BUFFER_SIZE];
	unsigned char bufferECDH1[BUFFER_SIZE];
	unsigned char bufferECDH2[BUFFER_SIZE];
	
	int sizePrivatekey1, sizePrivatekey2;
	int sizePubKey1, sizePubKey2;
	GoSlice privatekey1 = { bufferPrivatekey1, 0, BUFFER_SIZE };
	GoSlice privatekey2 = { bufferPrivatekey2, 0, BUFFER_SIZE };
	GoSlice pubKey1 = { bufferPubKey1, 0, BUFFER_SIZE };
	GoSlice pubKey2 = { bufferPubKey2, 0, BUFFER_SIZE };
	GoSlice ecdh1 = { bufferECDH1, 0, BUFFER_SIZE };
	GoSlice ecdh2 = { bufferECDH2, 0, BUFFER_SIZE };
	GoInt32 error_code;
	
	for(int i = 0; i < keys_count; i++){
		int randn = rand() % keys_count;
		sizePrivatekey1 = hexnstr(test_keys[i], bufferPrivatekey1, BUFFER_SIZE);
		sizePrivatekey2 = hexnstr(test_keys[randn], bufferPrivatekey2, BUFFER_SIZE);
		privatekey1.len = sizePrivatekey1;
		privatekey2.len = sizePrivatekey2;
		
		error_code = SKY_secp256k1_PubkeyFromSeckey(privatekey1, (coin__UxArray*)&pubKey1);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_PubkeyFromSeckey failed.");
		cr_assert(pubKey1.len > 0, "SKY_secp256k1_PubkeyFromSeckey failed.");
		error_code = SKY_secp256k1_PubkeyFromSeckey(privatekey2, (coin__UxArray*)&pubKey2);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_PubkeyFromSeckey failed.");
		cr_assert(pubKey2.len > 0, "SKY_secp256k1_PubkeyFromSeckey failed.");
		
		SKY_secp256k1_ECDH(pubKey1, privatekey2, (coin__UxArray*)&ecdh1);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_ECDH failed.");
		SKY_secp256k1_ECDH(pubKey2, privatekey1, (coin__UxArray*)&ecdh2);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_ECDH failed.");
		
		cr_assert(eq(type(GoSlice), ecdh1, ecdh2));
	}
}