#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

#define PLAINTEXT "plaintext"
#define PASSWORD "password"

unsigned char buff[1024];

Test(cipher_encrypt_scrypt_chacha20poly1305, TestScryptChacha20poly1305Encrypt){
	GoSlice text = {PLAINTEXT, strlen(PLAINTEXT), strlen(PLAINTEXT)};
	GoSlice password = {PASSWORD, strlen(PLAINTEXT), strlen(PASSWORD)};
	GoSlice result = {NULL, 0, 0};
	GoUint32 errcode;
	 
	encrypt__ScryptChacha20poly1305 encrypt = {1 << 20, 8, 1, 32};
	errcode = SKY_encrypt_ScryptChacha20poly1305_Encrypt(
			&encrypt, text, password, (coin__UxArray*)&result);
	cr_assert(errcode == SKY_OK, "SKY_encrypt_ScryptChacha20poly1305_Encrypt failed");
	{
		FILE *fp;
		fp = fopen("test.txt", "w+");
		fprintf(fp, "rpc: %s",  result.data);
		fclose(fp);
	}
	registerMemCleanup(result.data);
}