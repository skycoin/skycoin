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
#define PASSWORD1 "pwd"
#define PASSWORD2 "key"
#define PASSWORD3 "9JMkCPphe73NQvGhmab"
#define SHA256XORDATALENGTHSIZE 4
#define SHA256XORBLOCKSIZE 32
#define SHA256XORCHECKSUMSIZE 32
#define SHA256XORNONCESIZE 32


typedef struct{
	int 		dataLength;
	GoSlice* 	pwd;
	int			success;
} TEST_DATA;

Test(cipher_encrypt_sha256xor, TestSha256XorEncrypt){
	unsigned char buff[BUFFER_SIZE];
	unsigned char encryptedBuffer[BUFFER_SIZE];
	unsigned char encryptedText[BUFFER_SIZE];
	GoSlice data = { buff, 0, BUFFER_SIZE };
	GoSlice encrypted = { encryptedBuffer, 0, BUFFER_SIZE };
	GoSlice pwd1 = { PASSWORD1, strlen(PASSWORD1), strlen(PASSWORD1) };
	GoSlice pwd2 = { PASSWORD2, strlen(PASSWORD2), strlen(PASSWORD2) };
	GoSlice pwd3 = { PASSWORD3, strlen(PASSWORD3), strlen(PASSWORD3) };
	GoSlice nullPwd = {NULL, 0, 0};
	GoUint32 errcode;
	
	TEST_DATA test_data[] = {
		{1, &nullPwd, 0},
		{1, &pwd2, 1},
		{2, &pwd1, 1},
		{32, &pwd1, 1},
		{64, &pwd3, 1},
		{65, &pwd3, 1},
	};
	
	encrypt__Sha256Xor encryptSettings = {};
	
	for(int i = 0; i < 6; i++){
		randBytes(&data, test_data[i].dataLength);
		errcode = SKY_encrypt_Sha256Xor_Encrypt(&encryptSettings, data, *(test_data[i].pwd), &encrypted);
		if( test_data[i].success ){
			cr_assert(errcode == SKY_OK, "SKY_encrypt_Sha256Xor_Encrypt failed.");
		} else {
			cr_assert(errcode != SKY_OK, "SKY_encrypt_Sha256Xor_Encrypt with null pwd.");
		}
		if( errcode == SKY_OK ){
			cr_assert(encrypted.cap > 0, "Buffer for encrypted data is too short");
			cr_assert(encrypted.len < BUFFER_SIZE, "Too large encrypted data");
			((char*)encrypted.data)[encrypted.len] = 0;
			
			int n = (SHA256XORDATALENGTHSIZE + test_data[i].dataLength) / SHA256XORBLOCKSIZE;
			int m = (SHA256XORDATALENGTHSIZE + test_data[i].dataLength) % SHA256XORBLOCKSIZE;
			if ( m > 0 ) {
				n++;
			}
			
			int real_size;
			base64_decode_binary((const unsigned char*)encrypted.data, 
				encrypted.len, encryptedText, &real_size, BUFFER_SIZE);
			int totalEncryptedDataLen = SHA256XORCHECKSUMSIZE + SHA256XORNONCESIZE + 32 + n*SHA256XORBLOCKSIZE; // 32 is the hash data length
			
			cr_assert(totalEncryptedDataLen == real_size, "SKY_encrypt_Sha256Xor_Encrypt failed, encrypted data length incorrect.");
			cr_assert(SHA256XORCHECKSUMSIZE == sizeof(cipher__SHA256), "Size of SHA256 struct different than size in constant declaration");
			cipher__SHA256 enc_hash;
			cipher__SHA256 cal_hash;
			for(int j = 0; j < SHA256XORCHECKSUMSIZE; j++){
				enc_hash[j] = (GoUint8_)encryptedText[j];
			}
			int len_minus_checksum = real_size - SHA256XORCHECKSUMSIZE;
			GoSlice slice = {&encryptedText[SHA256XORCHECKSUMSIZE], len_minus_checksum, len_minus_checksum};
			SKY_cipher_SumSHA256(slice, &cal_hash);
			int equal = 1;
			for(int j = 0; j < SHA256XORCHECKSUMSIZE; j++){
				if(enc_hash[j] != cal_hash[j]){
					equal = 0;
					break;
				}
			}
			cr_assert(equal == 1, "SKY_encrypt_Sha256Xor_Encrypt failed, incorrect hash sum.");
		}
	}
	
	for(int i = 33; i <= 64; i++){
		randBytes(&data, i);
		errcode = SKY_encrypt_Sha256Xor_Encrypt(&encryptSettings, data, pwd1, &encrypted);
		cr_assert(errcode == SKY_OK, "SKY_encrypt_Sha256Xor_Encrypt failed.");
		cr_assert(encrypted.cap > 0, "Buffer for encrypted data is too short");
		cr_assert(encrypted.len < BUFFER_SIZE, "Too large encrypted data");
		((char*)encrypted.data)[encrypted.len] = 0;
		
		int n = (SHA256XORDATALENGTHSIZE + i) / SHA256XORBLOCKSIZE;
		int m = (SHA256XORDATALENGTHSIZE + i) % SHA256XORBLOCKSIZE;
		if ( m > 0 ) {
			n++;
		}
		
		int real_size;
		base64_decode_binary((const unsigned char*)encrypted.data, 
			encrypted.len, encryptedText, &real_size, BUFFER_SIZE);
		int totalEncryptedDataLen = SHA256XORCHECKSUMSIZE + SHA256XORNONCESIZE + 32 + n*SHA256XORBLOCKSIZE; // 32 is the hash data length
		
		cr_assert(totalEncryptedDataLen == real_size, "SKY_encrypt_Sha256Xor_Encrypt failed, encrypted data length incorrect.");
		cr_assert(SHA256XORCHECKSUMSIZE == sizeof(cipher__SHA256), "Size of SHA256 struct different than size in constant declaration");
		cipher__SHA256 enc_hash;
		cipher__SHA256 cal_hash;
		for(int j = 0; j < SHA256XORCHECKSUMSIZE; j++){
			enc_hash[j] = (GoUint8_)encryptedText[j];
		}
		int len_minus_checksum = real_size - SHA256XORCHECKSUMSIZE;
		GoSlice slice = {&encryptedText[SHA256XORCHECKSUMSIZE], len_minus_checksum, len_minus_checksum};
		SKY_cipher_SumSHA256(slice, &cal_hash);
		int equal = 1;
		for(int j = 0; j < SHA256XORCHECKSUMSIZE; j++){
			if(enc_hash[j] != cal_hash[j]){
				equal = 0;
				break;
			}
		}
		cr_assert(equal == 1, "SKY_encrypt_Sha256Xor_Encrypt failed, incorrect hash sum.");
		
	}
}