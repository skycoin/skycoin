#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>

#include "libskycoin.h"
#include "skyerrors.h"

void freshSumRipemd160(GoSlice bytes, Ripemd160* rp160){
	
	SKY_cipher_HashRipemd160(bytes, rp160);
}

void freshSumSHA256(GoSlice bytes, SHA256* sha256){
	
	SKY_cipher_SumSHA256(bytes, sha256);
}