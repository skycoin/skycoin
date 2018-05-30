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

TestSuite(cipher_base64_hex, .init = setup, .fini = teardown);

#define BUFFER_SIZE 1024

Test(cipher_base64_hex, TestBase64Encode){
	unsigned char buff[BUFFER_SIZE];
	unsigned char output[BUFFER_SIZE];
	unsigned char output2[BUFFER_SIZE];
	
	int n;
	GoSlice data =  { buff, 256, BUFFER_SIZE };
	GoSlice data2 = { output2, 0, BUFFER_SIZE };
	for(int i = 0; i < 10; i++){
		randBytes(&data, 256);
		unsigned int encode_result = b64_encode( (const unsigned char*)buff, 256, output );
		cr_assert(encode_result > 0, "b64_encode_string failed");
		int decode_result = b64_decode(output, encode_result, output2);
		cr_assert(decode_result > 0, "base64_decode_string failed");
		data2.len = decode_result;
		cr_assert( eq(type(GoSlice), data, data2) );
	}
	//Testing with invalid characters
	char* invalid_chars = "*<>;@#$()!";
	int chars_len = strlen(invalid_chars);
	for(int i = 0; i < 10; i++){
		randBytes(&data, 256);
		unsigned int encode_result = b64_encode( (const unsigned char*)buff, 256, output );
		cr_assert(encode_result > 0, "b64_encode_string failed");
		n = rand() % 256;
		output[n] = invalid_chars[ i % chars_len ];
		int decode_result = b64_decode(output, encode_result, output2);
		cr_assert(decode_result < 0);
	}
	//Truncating 
	for(int i = 0; i < 10; i++){
		randBytes(&data, 256);
		unsigned int encode_result = b64_encode( (const unsigned char*)buff, 256, output );
		cr_assert(encode_result > 0, "b64_encode_string failed");
		n = rand() % 256;
		output[n] = 0;
		int decode_result = b64_decode(output, encode_result, output2);
		cr_assert(decode_result < 0);
	}
}

Test(cipher_base64_hex, TestHexEncode){
	unsigned char buff[BUFFER_SIZE];
	unsigned char output[BUFFER_SIZE];
	unsigned char output2[BUFFER_SIZE];
	
	GoSlice data =  { buff, 256, BUFFER_SIZE };
	GoSlice data2 = { output2, 0, BUFFER_SIZE };
	for(int i = 0; i < 10; i++){
		randBytes(&data, 256);
		if( i % 2 == 0 )
			strnhex( buff, output, 256 );
		else
			strnhexlower( buff, output, 256 );
		int decode_result = hexnstr(output, output2, BUFFER_SIZE);
		cr_assert(decode_result > 0, "hexnstr failed");
		data2.len = decode_result;
		cr_assert( eq(type(GoSlice), data, data2) );
	}
	int n;
	//Testing with invalid characters
	char* invalid_chars = "*<>;@#$()!PVNQWR/=+-";
	int chars_len = strlen(invalid_chars);
	for(int i = 0; i < 20; i++){
		randBytes(&data, 256);
		if( i % 2 == 0 )
			strnhex( buff, output, 256 );
		else
			strnhexlower( buff, output, 256 );
		n = rand() % 256;
		output[n] = invalid_chars[ i % chars_len ];
		int decode_result = hexnstr(output, output2, BUFFER_SIZE);
		cr_assert(decode_result < 0);
	}
	//Truncating 
	for(int i = 0; i < 20; i++){
		randBytes(&data, 256);
		if( i % 2 == 0 )
			strnhex( buff, output, 256 );
		else
			strnhexlower( buff, output, 256 );
		n = rand() % 128 * 2 + 1;
		output[n] = 0;
		int decode_result = hexnstr(output, output2, BUFFER_SIZE);
		cr_assert(decode_result < 0);
	}
}


