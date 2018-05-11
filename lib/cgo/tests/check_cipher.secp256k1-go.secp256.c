
#include <stdio.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

#define BUFFER_SIZE 128

Test(cipher_secp256k1, TestPubkeyFromSeckey) {
	unsigned char bufferPrivate[BUFFER_SIZE];
	unsigned char bufferPublic[BUFFER_SIZE];
	unsigned char bufferResult[BUFFER_SIZE];
	
	const char* hexPrivate = "f19c523315891e6e15ae0608a35eec2e00ebd6d1984cf167f46336dabd9b2de4";
	const char* hexPublic  = "03FE43D0C2C3DAAB30F9472BEB5B767BE020B81C7CC940ED7A7E910F0C1D9FEEF1";
	 
	
	int sizePrivate = hexnstr(hexPrivate, bufferPrivate, BUFFER_SIZE);
	int sizePublic = hexnstr(hexPublic, bufferPublic, BUFFER_SIZE);
	GoSlice privateKey = {bufferPrivate, sizePrivate, BUFFER_SIZE};
	GoSlice publicKey = {bufferPublic, sizePublic, BUFFER_SIZE};
	coin__UxArray result = {bufferResult, 0, BUFFER_SIZE};
	
	GoUint32 error_code = SKY_secp256k1_PubkeyFromSeckey(privateKey, &result);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_PubkeyFromSeckey failed.");

	cr_assert(result.len == publicKey.len, "SKY_secp256k1_PubkeyFromSeckey failed. Calculated pub key doesn\'t have expected length");
	int equal = 1;
	for(int i = 0; i < result.len; i++){
		if( ((char*)result.data)[i] != ((char*)publicKey.data)[i] ){
			equal = 0;
			break;
		}
	}
	cr_assert(equal == 1, "SKY_secp256k1_PubkeyFromSeckey failed. Calculated pub key is different than expected.");
}