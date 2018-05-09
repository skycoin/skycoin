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
#include "jsmn.h"

#define PLAINTEXT "plaintext"
#define PASSWORD "password"
#define PASSWORD2 "pwd"
#define WRONG_PASSWORD "wrong password"
#define ENCRYPTED "dQB7Im4iOjUyNDI4OCwiciI6OCwicCI6MSwia2V5TGVuIjozMiwic2FsdCI6ImpiejUrSFNjTFFLWkI5T0tYblNNRmt2WDBPY3JxVGZ0ZFpDNm9KUFpaeHc9Iiwibm9uY2UiOiJLTlhOQmRQa1ZUWHZYNHdoIn3PQFmOot0ETxTuv//skTG7Q57UVamGCgG5"
#define BUFFER_SIZE 1024
#define SCRYPTCHACHA20METALENGTHSIZE 2

Test(cipher_encrypt_scrypt_chacha20poly1305, TestScryptChacha20poly1305Encrypt){
	GoSlice text = {PLAINTEXT, strlen(PLAINTEXT), strlen(PLAINTEXT)};
	GoSlice password = {PASSWORD, strlen(PASSWORD), strlen(PASSWORD)};
	char buffer[BUFFER_SIZE];
	char str[BUFFER_SIZE];
	char tmp[BUFFER_SIZE];
	unsigned int data[BUFFER_SIZE];
	GoSlice result = {buffer, 0, BUFFER_SIZE};
	GoUint32 errcode;
	unsigned int metalength;
	
	for(int i = 1; i <= 20; i++) {
		encrypt__ScryptChacha20poly1305 encrypt = {1 << i, 8, 1, 32};
		errcode = SKY_encrypt_ScryptChacha20poly1305_Encrypt(
				&encrypt, text, password, (coin__UxArray*)&result);
		cr_assert(errcode == SKY_OK, "SKY_encrypt_ScryptChacha20poly1305_Encrypt failed");
		cr_assert(result.len > SCRYPTCHACHA20METALENGTHSIZE, "SKY_encrypt_ScryptChacha20poly1305_Encrypt failed, result data length too short");
		cr_assert(result.len < BUFFER_SIZE, "SKY_encrypt_ScryptChacha20poly1305_Encrypt failed, result data length greater than buffer");
		((char*)result.data)[result.len] = 0;
		unsigned int decode_len = base64_decode_string((const unsigned char*)result.data, 
				result.len, str, BUFFER_SIZE);
		cr_assert(decode_len >= SCRYPTCHACHA20METALENGTHSIZE, "base64_decode_string failed");
		cr_assert(decode_len < BUFFER_SIZE, "base64_decode_string failed, buffer overflow");
		metalength = (unsigned int)	str[0];
		for(int m = 1; m < SCRYPTCHACHA20METALENGTHSIZE; m++){
			if(str[m] > 0){
				metalength += (((unsigned int)str[m]) << (m * 8));
			}
		}
		cr_assert(metalength + SCRYPTCHACHA20METALENGTHSIZE < decode_len, "SKY_encrypt_ScryptChacha20poly1305_Encrypt failed. Metadata length greater than result lentgh.");
		char* meta = &str[SCRYPTCHACHA20METALENGTHSIZE];
		meta[metalength] = 0;
		
		jsmn_parser parser;
		jsmntok_t tokens[128];
		jsmn_init(&parser);
		int tokens_count = jsmn_parse(&parser, meta, strlen(meta), tokens, sizeof(tokens)/sizeof(tokens[0]));
		cr_assert(tokens_count > 0, "Failed to json parse ScryptChacha20poly1305 encoded string");
		cr_assert(tokens[0].type == JSMN_OBJECT, "Failed to json parse ScryptChacha20poly1305 encoded string");
		int n = json_get_int(meta, "n", tokens, tokens_count);
		int r = json_get_int(meta, "r", tokens, tokens_count);
		int p = json_get_int(meta, "p", tokens, tokens_count);
		int keyLen = json_get_int(meta, "keyLen", tokens, tokens_count);
		cr_assert(n == encrypt.N, "SKY_encrypt_ScryptChacha20poly1305_Encrypt failed. Metadata value N incorrect.");
		cr_assert(r == encrypt.R, "SKY_encrypt_ScryptChacha20poly1305_Encrypt failed. Metadata value R incorrect.");
		cr_assert(p == encrypt.P, "SKY_encrypt_ScryptChacha20poly1305_Encrypt failed. Metadata value P incorrect.");
		cr_assert(keyLen == encrypt.KeyLen, "SKY_encrypt_ScryptChacha20poly1305_Encrypt failed. Metadata value KeyLen incorrect.");
	}
}

Test(cipher_encrypt_scrypt_chacha20poly1305, TestScryptChacha20poly1305Decrypt){
	GoSlice text = {PLAINTEXT, strlen(PLAINTEXT), strlen(PLAINTEXT)};
	GoSlice password = {PASSWORD2, strlen(PASSWORD2), strlen(PASSWORD2)};
	GoSlice wrong_password = {WRONG_PASSWORD, strlen(WRONG_PASSWORD), strlen(WRONG_PASSWORD)};
	GoSlice encrypted = {ENCRYPTED, strlen(ENCRYPTED), strlen(ENCRYPTED)};
	GoSlice nullData = {NULL, 0, 0};
	GoSlice nullPassword = {NULL, 0, 0};
	
	char buffer[BUFFER_SIZE];
	coin__UxArray result = {buffer, 0, BUFFER_SIZE};
	GoUint32 errcode;
	encrypt__ScryptChacha20poly1305 encrypt = {0, 0, 0, 0};
	
	errcode = SKY_encrypt_ScryptChacha20poly1305_Decrypt(&encrypt, encrypted, password, &result);
	cr_assert(errcode == SKY_OK, "SKY_encrypt_ScryptChacha20poly1305_Decrypt failed");
	cr_assert(result.len < BUFFER_SIZE, "SKY_encrypt_ScryptChacha20poly1305_Decrypt failed. Buffer to small.");
	((char*)result.data)[result.len] = 0;
	cr_assert(strcmp((char*)text.data, (char*)result.data) == 0, "SKY_encrypt_ScryptChacha20poly1305_Decrypt failed. Decrypted text incorrect.");
	
	errcode = SKY_encrypt_ScryptChacha20poly1305_Decrypt(&encrypt, encrypted, wrong_password, &result);
	cr_assert(errcode != SKY_OK, "SKY_encrypt_ScryptChacha20poly1305_Decrypt decrypted with wrong password.");
	errcode = SKY_encrypt_ScryptChacha20poly1305_Decrypt(&encrypt, nullData, password, &result);
	cr_assert(errcode != SKY_OK, "SKY_encrypt_ScryptChacha20poly1305_Decrypt decrypted with null encrypted data.");
	errcode = SKY_encrypt_ScryptChacha20poly1305_Decrypt(&encrypt, encrypted, nullPassword, &result);
	cr_assert(errcode != SKY_OK, "SKY_encrypt_ScryptChacha20poly1305_Decrypt decrypted with null password.");
}