
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

Test(cipher_secp256k1,Test_Secp256_00){
	unsigned char buff[SigSize];
	visor__ReadableOutputs nonce = {buff,0,64};
	SKY_cipher_RandByte(32,&nonce);
	if (nonce.len != 32) cr_fatal();
}

//test agreement for highest bit test
// func Test_BitTwiddle(t *testing.T) {
// 	var b byte
// 	for i := 0; i < 512; i++ {
// 		bool1 := ((b >> 7) == 1)
// 		bool2 := ((b & 0x80) == 0x80)
// 		if bool1 != bool2 {
// 			t.Fatal()
// 		}
// 		b++
// 	}
// }
//
Test(cipher_secp256k1,Test_BitTwiddle){
	cr_fatal();
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
