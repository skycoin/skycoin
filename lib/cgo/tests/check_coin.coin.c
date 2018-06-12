
#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"
#include "transutil.h"
#include "time.h"

TestSuite(coin_coin, .init = setup, .fini = teardown);

Test(coin_coin, TestAddress1){
  char* address_hex = "02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8";
  char address[128];
  int result;
  int length = hexnstr(address_hex, address, 128);
  cr_assert(length > 0, "Error decoding hex string");
  GoSlice slice = { address, length, 128 };
  cipher__PubKey pubkey;
  result = SKY_cipher_NewPubKey(slice, &pubkey);
  cr_assert( result == SKY_OK, "SKY_cipher_NewPubKey failed" );
  cipher__Address c_address;
  result = SKY_cipher_AddressFromPubKey( &pubkey, &c_address );
  cr_assert( result == SKY_OK, "SKY_cipher_AddressFromPubKey failed" );
}

Test(coin_coin, TestAddress2){
  char* address_hex = "5a42c0643bdb465d90bf673b99c14f5fa02db71513249d904573d2b8b63d353d";
  char address[128];
  int result;
  int length = hexnstr(address_hex, address, 128);
  cr_assert(length > 0, "Error decoding hex string");
  GoSlice slice = { address, length, 128 };
  cipher__PubKey pubkey;
  cipher__SecKey seckey;
  result = SKY_cipher_NewSecKey(slice, &seckey);
  cr_assert( result == SKY_OK, "SKY_cipher_NewSecKey failed" );
  result = SKY_cipher_PubKeyFromSecKey(&seckey, &pubkey);
  cr_assert( result == SKY_OK, "SKY_cipher_PubKeyFromSecKey failed" );
  cipher__Address c_address;
  result = SKY_cipher_AddressFromPubKey( &pubkey, &c_address );
  cr_assert( result == SKY_OK, "SKY_cipher_AddressFromPubKey failed" );
}

Test(coin_coin, TestCrypto1){
  cipher__PubKey pubkey;
  cipher__SecKey seckey;
  int result;
  for(int i = 0; i < 10; i ++){
    result = SKY_cipher_GenerateKeyPair( &pubkey, &seckey );
    cr_assert( result == SKY_OK, "SKY_cipher_GenerateKeyPair failed" );
    result = SKY_cipher_TestSecKey( &seckey );
    cr_assert( result == SKY_OK, "CRYPTOGRAPHIC INTEGRITY CHECK FAILED" );
  }
}

Test(coin_coin, TestCrypto2){
    char* address_hex = "5a42c0643bdb465d90bf673b99c14f5fa02db71513249d904573d2b8b63d353d";
    char address[128];
    int result;
    int length = hexnstr(address_hex, address, 128);
    cr_assert(length == 32, "Error decoding hex string");

    GoSlice slice = { address, length, 128 };
    cipher__PubKey pubkey;
    cipher__SecKey seckey;
    result = SKY_cipher_NewSecKey(slice, &seckey);
    cr_assert( result == SKY_OK, "SKY_cipher_NewSecKey failed" );
    result = SKY_cipher_PubKeyFromSecKey(&seckey, &pubkey);
    cr_assert( result == SKY_OK, "SKY_cipher_PubKeyFromSecKey failed" );
    cipher__Address c_address;
    result = SKY_cipher_AddressFromPubKey( &pubkey, &c_address );
    cr_assert( result == SKY_OK, "SKY_cipher_AddressFromPubKey failed" );

    char* text = "test message";
    int len = strlen(text);
    GoSlice textslice = {text, len, len};
    cipher__SHA256 hash;
    result = SKY_cipher_SumSHA256(textslice, &hash);
    cr_assert( result == SKY_OK, "SKY_cipher_SumSHA256 failed" );
    result = SKY_cipher_TestSecKeyHash( &seckey, &hash );
    cr_assert( result == SKY_OK, "SKY_cipher_TestSecKeyHash failed" );
}
