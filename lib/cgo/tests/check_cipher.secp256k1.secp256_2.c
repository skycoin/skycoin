
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

//test size of messages
Test(cipher_secp256k1, Test_Secp256_02s){
	GoInt32 error_code;
	char bufferPub1[BUFFER_SIZE];
	char bufferSec1[BUFFER_SIZE];
	char bufferSig1[BUFFER_SIZE];
	unsigned char buff[32];
	GoSlice pub1 = {bufferPub1, 0, BUFFER_SIZE};
	GoSlice sec1 = {bufferSec1, 0, BUFFER_SIZE};
	GoSlice sig  = {bufferSig1, 0, BUFFER_SIZE};
	GoSlice msg  = {buff, 0, 32};
	
	error_code = SKY_secp256k1_GenerateKeyPair(
		(coin__UxArray*)&pub1, (coin__UxArray*)&sec1);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateKeyPair failed");
	
	error_code = SKY_secp256k1_RandByte(32, (coin__UxArray*)&msg);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed");
	error_code == SKY_secp256k1_Sign(msg, sec1, (GoSlice_*)&sig);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_Sign failed");
	cr_assert(pub1.len == 33, "Public key should be 33 bytes long.");
	cr_assert(sec1.len == 32, "Private key should be 32 bytes long.");
	cr_assert(sig.len == 65, "Signature should be 65 bytes long.");
	unsigned char last = ((unsigned char*) sig.data)[64]; 
	cr_assert( last <= 4 );
}

//test signing message
Test(cipher_secp256k1, Test_Secp256_02){
	GoInt32 error_code;
	char bufferPub1[BUFFER_SIZE];
	char bufferPub2[BUFFER_SIZE];
	char bufferSec1[BUFFER_SIZE];
	char bufferSig1[BUFFER_SIZE];
	unsigned char buff[32];
	GoSlice pub1 = {bufferPub1, 0, BUFFER_SIZE};
	GoSlice sec1 = {bufferSec1, 0, BUFFER_SIZE};
	GoSlice pub2 = {bufferPub2, 0, BUFFER_SIZE};
	GoSlice sig  = {bufferSig1, 0, BUFFER_SIZE};
	GoSlice msg  = {buff, 0, 32};
	
	error_code = SKY_secp256k1_GenerateKeyPair(
		(coin__UxArray*)&pub1, (coin__UxArray*)&sec1);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateKeyPair failed");
	
	error_code = SKY_secp256k1_RandByte(32, (coin__UxArray*)&msg);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed");
	error_code == SKY_secp256k1_Sign(msg, sec1, (GoSlice_*)&sig);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_Sign failed");
	cr_assert(sig.len == 65, "Signature should be 65 bytes long.");
	
	error_code = SKY_secp256k1_RecoverPubkey(msg, sig, (coin__UxArray*)&pub2);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_RecoverPubkey failed");
	cr_assert(eq(type(GoSlice), pub1, pub2), "Different public keys.");
	
	GoInt result;
	error_code = SKY_secp256k1_VerifySignature(msg, sig, pub1, &result);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_VerifySignature failed");
	cr_assert(result, "Signature invalid");
}

//test pubkey recovery
Test(cipher_secp256k1, Test_Secp256_02a){
	GoInt32 error_code;
	char bufferPub1[BUFFER_SIZE];
	char bufferPub2[BUFFER_SIZE];
	char bufferSec1[BUFFER_SIZE];
	char bufferSig1[BUFFER_SIZE];
	unsigned char buff[32];
	GoSlice pub1 = {bufferPub1, 0, BUFFER_SIZE};
	GoSlice sec1 = {bufferSec1, 0, BUFFER_SIZE};
	GoSlice pub2 = {bufferPub2, 0, BUFFER_SIZE};
	GoSlice sig  = {bufferSig1, 0, BUFFER_SIZE};
	GoSlice msg  = {buff, 0, 32};
	
	error_code = SKY_secp256k1_GenerateKeyPair(
		(coin__UxArray*)&pub1, (coin__UxArray*)&sec1);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateKeyPair failed");
	
	error_code = SKY_secp256k1_RandByte(32, (coin__UxArray*)&msg);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed");
	error_code == SKY_secp256k1_Sign(msg, sec1, (GoSlice_*)&sig);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_Sign failed");
	cr_assert(sig.len == 65, "Signature should be 65 bytes long.");
	GoInt result;
	error_code = SKY_secp256k1_VerifySignature(msg, sig, pub1, &result);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_VerifySignature failed");
	cr_assert(result, "Signature invalid");
	
	error_code = SKY_secp256k1_RecoverPubkey(msg, sig, (coin__UxArray*)&pub2);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_RecoverPubkey failed");
	cr_assert(eq(type(GoSlice), pub1, pub2), "Different public keys.");
}

//test random messages for the same pub/private key
Test(cipher_secp256k1, Test_Secp256_03){
	GoInt32 error_code;
	char bufferPub1[BUFFER_SIZE];
	char bufferPub2[BUFFER_SIZE];
	char bufferSec1[BUFFER_SIZE];
	char bufferSig1[BUFFER_SIZE];
	unsigned char buff[32];
	GoSlice pub1 = {bufferPub1, 0, BUFFER_SIZE};
	GoSlice sec1 = {bufferSec1, 0, BUFFER_SIZE};
	GoSlice pub2 = {bufferPub2, 0, BUFFER_SIZE};
	GoSlice sig  = {bufferSig1, 0, BUFFER_SIZE};
	GoSlice msg  = {buff, 0, 32};
	
	for( int i = 0; i < TESTS; i++ ) {
		error_code = SKY_secp256k1_GenerateKeyPair(
			(coin__UxArray*)&pub1, (coin__UxArray*)&sec1);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateKeyPair failed");
		
		error_code = SKY_secp256k1_RandByte(32, (coin__UxArray*)&msg);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed");
		error_code == SKY_secp256k1_Sign(msg, sec1, (GoSlice_*)&sig);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_Sign failed");
		cr_assert(sig.len == 65, "Signature should be 65 bytes long.");
		((unsigned char*)sig.data)[64] = ((unsigned char*)sig.data)[64] % 4;
		
		error_code = SKY_secp256k1_RecoverPubkey(msg, sig, (coin__UxArray*)&pub2);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_RecoverPubkey failed");
		cr_assert(pub2.len > 0, "Invalid public key");
	}
}

//test random messages for different pub/private keys
Test(cipher_secp256k1, Test_Secp256_04){
	GoInt32 error_code;
	char bufferPub1[BUFFER_SIZE];
	char bufferPub2[BUFFER_SIZE];
	char bufferSec1[BUFFER_SIZE];
	char bufferSig1[BUFFER_SIZE];
	unsigned char buff[32];
	GoSlice pub1 = {bufferPub1, 0, BUFFER_SIZE};
	GoSlice sec1 = {bufferSec1, 0, BUFFER_SIZE};
	GoSlice pub2 = {bufferPub2, 0, BUFFER_SIZE};
	GoSlice sig  = {bufferSig1, 0, BUFFER_SIZE};
	GoSlice msg  = {buff, 0, 32};
	
	for( int i = 0; i < TESTS; i++ ) {
		error_code = SKY_secp256k1_GenerateKeyPair(
			(coin__UxArray*)&pub1, (coin__UxArray*)&sec1);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateKeyPair failed");
		
		error_code = SKY_secp256k1_RandByte(32, (coin__UxArray*)&msg);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed");
		error_code == SKY_secp256k1_Sign(msg, sec1, (GoSlice_*)&sig);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_Sign failed");
		cr_assert(sig.len == 65, "Signature should be 65 bytes long.");
		unsigned char last = ((unsigned char*) sig.data)[64]; 
		cr_assert( last < 4 );
		error_code = SKY_secp256k1_RecoverPubkey(msg, sig, (coin__UxArray*)&pub2);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_RecoverPubkey failed");
		cr_assert(pub2.len > 0, "Invalid public key");
		cr_assert(eq(type(GoSlice), pub1, pub2), "Different public keys.");
	}
}

Test(cipher_secp256k1, Test_Secp256_06a_alt0){
	GoInt32 error_code;
	char bufferPub1[BUFFER_SIZE];
	char bufferPub2[BUFFER_SIZE];
	char bufferSec1[BUFFER_SIZE];
	char bufferSig1[BUFFER_SIZE];
	unsigned char buff[32];
	GoSlice pub1 = {bufferPub1, 0, BUFFER_SIZE};
	GoSlice sec1 = {bufferSec1, 0, BUFFER_SIZE};
	GoSlice pub2 = {bufferPub2, 0, BUFFER_SIZE};
	GoSlice sig  = {bufferSig1, 0, BUFFER_SIZE};
	GoSlice msg  = {buff, 0, 32};
	
	error_code = SKY_secp256k1_GenerateKeyPair(
		(coin__UxArray*)&pub1, (coin__UxArray*)&sec1);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateKeyPair failed");
	error_code = SKY_secp256k1_RandByte(32, (coin__UxArray*)&msg);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed");
	error_code == SKY_secp256k1_Sign(msg, sec1, (GoSlice_*)&sig);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_Sign failed");
	cr_assert(sig.len == 65, "Signature should be 65 bytes long.");
	
	GoInt result;
	for(int i = 0; i < TESTS; i++){
		error_code = SKY_secp256k1_RandByte(65, (coin__UxArray*)&sig);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed");
		cr_assert(sig.len == 65, "Signature should be 65 bytes long.");
		((unsigned char*)sig.data)[32] = ((unsigned char*)sig.data)[32] & 0x70;
		((unsigned char*)sig.data)[64] = ((unsigned char*)sig.data)[64] % 4;
		
		error_code = SKY_secp256k1_RecoverPubkey(msg, sig, (coin__UxArray*)&pub2);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_RecoverPubkey failed");
		cr_assert(cr_user_GoSlice_noteq(&pub1, &pub2), "Public keys must be different.");
		SKY_secp256k1_VerifySignature(msg, sig, pub2, &result);
		cr_assert(pub2.len == 0 || result, "Public key is not valid");
		error_code = SKY_secp256k1_VerifySignature(msg, sig, pub1, &result);
		cr_assert(result == 0, "Public key should not be valid");
	}
}

Test(cipher_secp256k1, Test_Secp256_06b){
	GoInt32 error_code;
	char bufferPub1[BUFFER_SIZE];
	char bufferPub2[BUFFER_SIZE];
	char bufferSec1[BUFFER_SIZE];
	char bufferSig1[BUFFER_SIZE];
	unsigned char buff[32];
	GoSlice pub1 = {bufferPub1, 0, BUFFER_SIZE};
	GoSlice sec1 = {bufferSec1, 0, BUFFER_SIZE};
	GoSlice pub2 = {bufferPub2, 0, BUFFER_SIZE};
	GoSlice sig  = {bufferSig1, 0, BUFFER_SIZE};
	GoSlice msg  = {buff, 0, 32};
	
	error_code = SKY_secp256k1_GenerateKeyPair(
		(coin__UxArray*)&pub1, (coin__UxArray*)&sec1);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateKeyPair failed");
	error_code = SKY_secp256k1_RandByte(32, (coin__UxArray*)&msg);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed");
	error_code == SKY_secp256k1_Sign(msg, sec1, (GoSlice_*)&sig);
	cr_assert(error_code == SKY_OK, "SKY_secp256k1_Sign failed");
	
	GoInt result;
	for(int i = 0; i < TESTS; i++){
		error_code = SKY_secp256k1_RandByte(32, (coin__UxArray*)&msg);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed");
		error_code = SKY_secp256k1_RecoverPubkey(msg, sig, (coin__UxArray*)&pub2);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_RecoverPubkey failed");
		cr_assert(cr_user_GoSlice_noteq(&pub1, &pub2), "Public keys must be different.");
		error_code = SKY_secp256k1_VerifySignature(msg, sig, pub2, &result);
		cr_assert(pub2.len == 0 || result, "Public key is not valid");
		SKY_secp256k1_VerifySignature(msg, sig, pub1, &result);
		cr_assert(result == 0, "Public key should not be valid");
	}
}

Test(cipher_secp256k1, Test_Deterministic_Keypairs_00){
	char bufferSeed[BUFFER_SIZE];
	char bufferHash[BUFFER_SIZE];
	char bufferPub1[BUFFER_SIZE];
	char bufferSec1[BUFFER_SIZE];
	char bufferPub2[BUFFER_SIZE];
	char bufferSec2[BUFFER_SIZE];
	
	GoSlice seed = {bufferSeed, 0, BUFFER_SIZE};
	GoSlice hash = {bufferHash, 0, BUFFER_SIZE};
	GoSlice pub1 = {bufferPub1, 0, BUFFER_SIZE};
	GoSlice sec1 = {bufferSec1, 0, BUFFER_SIZE};
	GoSlice pub2 = {bufferPub2, 0, BUFFER_SIZE};
	GoSlice sec2 = {bufferSec2, 0, BUFFER_SIZE};
	GoInt32 error_code;
	
	for( int i = 0; i < 64; i++){
		error_code = SKY_secp256k1_RandByte( 32, (coin__UxArray*)&seed);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed");
		error_code = SKY_secp256k1_DeterministicKeyPairIterator(seed,
			(coin__UxArray*)&hash, 
			(coin__UxArray*)&pub1, 
			(coin__UxArray*)&sec1);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_DeterministicKeyPairIterator failed");
		error_code = SKY_secp256k1_GenerateDeterministicKeyPair(seed,
			(coin__UxArray*)&pub2, 
			(coin__UxArray*)&sec2);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateDeterministicKeyPair failed");
		cr_assert(eq(type(GoSlice), pub1, pub2), "Different public keys.");
		cr_assert(eq(type(GoSlice), sec1, sec2), "Different private keys.");
	}
}

Test(cipher_secp256k1, Test_Deterministic_Keypairs_01){
	char bufferSeed[BUFFER_SIZE];
	char bufferHash[BUFFER_SIZE];
	char bufferPub1[BUFFER_SIZE];
	char bufferSec1[BUFFER_SIZE];
	char bufferPub2[BUFFER_SIZE];
	char bufferSec2[BUFFER_SIZE];
	
	GoSlice seed = {bufferSeed, 0, BUFFER_SIZE};
	GoSlice hash = {bufferHash, 0, BUFFER_SIZE};
	GoSlice pub1 = {bufferPub1, 0, BUFFER_SIZE};
	GoSlice sec1 = {bufferSec1, 0, BUFFER_SIZE};
	GoSlice pub2 = {bufferPub2, 0, BUFFER_SIZE};
	GoSlice sec2 = {bufferSec2, 0, BUFFER_SIZE};
	GoInt32 error_code;
	
	for( int i = 0; i < 64; i++){
		error_code = SKY_secp256k1_RandByte( 32, (coin__UxArray*)&seed);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed");
		error_code = SKY_secp256k1_DeterministicKeyPairIterator(seed,
			(coin__UxArray*)&hash, 
			(coin__UxArray*)&pub1, 
			(coin__UxArray*)&sec1);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_DeterministicKeyPairIterator failed");
		error_code = SKY_secp256k1_GenerateDeterministicKeyPair(seed,
			(coin__UxArray*)&pub2, 
			(coin__UxArray*)&sec2);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateDeterministicKeyPair failed");
		cr_assert(eq(type(GoSlice), pub1, pub2), "Different public keys.");
		cr_assert(eq(type(GoSlice), sec1, sec2), "Different private keys.");
	}
}

Test(cipher_secp256k1, Test_Deterministic_Keypairs_02){
	char bufferSeed[BUFFER_SIZE];
	char bufferHash[BUFFER_SIZE];
	char bufferPub1[BUFFER_SIZE];
	char bufferSec1[BUFFER_SIZE];
	char bufferPub2[BUFFER_SIZE];
	char bufferSec2[BUFFER_SIZE];
	
	GoSlice seed = {bufferSeed, 0, BUFFER_SIZE};
	GoSlice hash = {bufferHash, 0, BUFFER_SIZE};
	GoSlice pub1 = {bufferPub1, 0, BUFFER_SIZE};
	GoSlice sec1 = {bufferSec1, 0, BUFFER_SIZE};
	GoSlice pub2 = {bufferPub2, 0, BUFFER_SIZE};
	GoSlice sec2 = {bufferSec2, 0, BUFFER_SIZE};
	GoInt32 error_code;
	
	for( int i = 0; i < 64; i++){
		error_code = SKY_secp256k1_RandByte( 32, (coin__UxArray*)&seed);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed");
		error_code = SKY_secp256k1_DeterministicKeyPairIterator(seed,
			(coin__UxArray*)&hash, 
			(coin__UxArray*)&pub1, 
			(coin__UxArray*)&sec1);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_DeterministicKeyPairIterator failed");
		error_code = SKY_secp256k1_GenerateDeterministicKeyPair(seed,
			(coin__UxArray*)&pub2, 
			(coin__UxArray*)&sec2);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_GenerateDeterministicKeyPair failed");
		cr_assert(eq(type(GoSlice), pub1, pub2), "Different public keys.");
		cr_assert(eq(type(GoSlice), sec1, sec2), "Different private keys.");
	}
}

Test(cipher_secp256k1, Test_Deterministic_Keypairs_03){
	int test_count = 16;
	const char* testArray[] = {
		"tQ93w5Aqcunm9SGUfnmF4fJv", "9b8c3e36adce64dedc80d6dfe51ff1742cc1d755bbad457ac01177c5a18a789f",
		"DC7qdQQtbWSSaekXnFmvQgse", "d2deaf4a9ff7a5111fe1d429d6976cbde78811fdd075371a2a4449bb0f4d8bf9",
		"X8EkuUZC7Td7PAXeS7Duc7vR", "cad79b6dcf7bd21891cbe20a51c57d59689ae6e3dc482cd6ec22898ac00cd86b",
		"tVqPYHHNVPRWyEed62v7f23u", "2a386e94e9ffaa409517cbed81b9b2d4e1c5fb4afe3cbd67ce8aba11af0b02fa",
		"kCy4R57HDfLqF3pVhBWxuMcg", "26a7c6d8809c476a56f7455209f58b5ff3f16435fcf208ff2931ece60067f305",
		"j8bjv86ZNjKqzafR6mtSUVCE", "ea5c0f8c9f091a70bf38327adb9b2428a9293e7a7a75119920d759ecfa03a995",
		"qShryAzVY8EtsuD3dsAc7qnG", "331206176509bcae31c881dc51e90a4e82ec33cd7208a5fb4171ed56602017fa",
		"5FGG7ZBa8wVMBJkmzpXj5ESX", "4ea2ad82e7730d30c0c21d01a328485a0cf5543e095139ba613929be7739b52c",
		"f46TZG4xJHXUGWx8ekbNqa9F", "dcddd403d3534c4ef5703cc07a771c107ed49b7e0643c6a2985a96149db26108",
		"XkZdQJ5LT96wshN8JBH8rvEt", "3e276219081f072dff5400ca29a9346421eaaf3c419ff1474ac1c81ad8a9d6e1",
		"GFDqXU4zYymhJJ9UGqRgS8ty", "95be4163085b571e725edeffa83fff8e7a7db3c1ccab19d0f3c6e105859b5e10",
		"tmwZksH2XyvuamnddYxyJ5Lp", "2666dd54e469df56c02e82dffb4d3ea067daafe72c54dc2b4f08c4fb3a7b7e42",
		"EuqZFsbAV5amTzkhgAMgjr7W", "40c325c01f2e4087fcc97fcdbea6c35c88a12259ebf1bce0b14a4d77f075abbf",
		"TW6j8rMffZfmhyDEt2JUCrLB", "e676e0685c5d1afd43ad823b83db5c6100135c35485146276ee0b0004bd6689e",
		"8rvkBnygfhWP8kjX9aXq68CY", "21450a646eed0d4aa50a1736e6c9bf99fff006a470aab813a2eff3ee4d460ae4",
		"phyRfPDuf9JMRFaWdGh7NXPX", "ca7bc04196c504d0e815e125f7f1e086c8ae8c10d5e9df984aeab4b41bf9e398",
	};
	
	GoInt32 error_code;
	char bufferSec1[BUFFER_SIZE];
	char bufferSec2[BUFFER_SIZE];
	char buffer1[BUFFER_SIZE];
	char buffer2[BUFFER_SIZE];
	
	GoSlice seed = {NULL, 0, BUFFER_SIZE};
	GoSlice sec1 = {bufferSec1, 0, BUFFER_SIZE};
	GoSlice sec2 = {bufferSec2, 0, BUFFER_SIZE};
	GoSlice s1   = {buffer1, 0, BUFFER_SIZE};
	GoSlice s2   = {buffer2, 0, BUFFER_SIZE};
	
	for(int i = 0; i < test_count; i++){
		seed.data = (void*)testArray[2 * i];
		seed.len = strlen(testArray[2 * i]);
		seed.cap = seed.len;
		sec1.len = hexnstr(testArray[2 * i + 1], bufferSec1, BUFFER_SIZE);
		error_code = SKY_secp256k1_DeterministicKeyPairIterator(seed, 
			(coin__UxArray*)&s1, (coin__UxArray*)&s2, 
			(coin__UxArray*)&sec2);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_DeterministicKeyPairIterator failed");
		cr_assert(eq(type(GoSlice), sec1, sec2), "Different hashes");
	}
}

Test(cipher_secp256k1, Test_DeterministicWallets1){
	int test_count = 16;
	const char* testArray[] = {
		"90c56f5b8d78a46fb4cddf6fd9c6d88d6d2d7b0ec35917c7dac12c03b04e444e", "94dd1a9de9ffd57b5516b8a7f090da67f142f7d22356fa5d1b894ee4d4fba95b",
		"a3b08ccf8cbae4955c02f223be1f97d2bb41d92b7f0c516eb8467a17da1e6057", "82fba4cc2bc29eef122f116f45d01d82ff488d7ee713f8a95c162a64097239e0",
		"7048eb8fa93cec992b93dc8e93c5543be34aad05239d4c036cf9e587bbcf7654", "44c059496aac871ac168bb6889b9dd3decdb9e1fa082442a95fcbca982643425",
		"6d25375591bbfce7f601fc5eb40e4f3dde2e453dc4bf31595d8ec29e4370cd80", "d709ceb1a6fb906de506ea091c844ca37c65e52778b8d257d1dd3a942ab367fb",
		"7214b4c09f584c5ddff971d469df130b9a3c03e0277e92be159279de39462120", "5fe4986fa964773041e119d2b6549acb392b2277a72232af75cbfb62c357c1a7",
		"b13e78392d5446ae304b5fc9d45b85f26996982b2c0c86138afdac8d2ea9016e", "f784abc2e7f11ee84b4adb72ea4730a6aabe27b09604c8e2b792d8a1a31881ac",
		"9403bff4240a5999e17e0ab4a645d6942c3a7147c7834e092e461a4580249e6e", "d495174b8d3f875226b9b939121ec53f9383bd560d34aa5ca3ac6b257512adf4",
		"2665312a3e3628f4df0b9bc6334f530608a9bcdd4d1eef174ecda99f51a6db94", "1fdc9fbfc6991b9416b3a8385c9942e2db59009aeb2d8de349b73d9f1d389374",
		"6cb37532c80765b7c07698502a49d69351036f57a45a5143e33c57c236d841ca", "c87c85a6f482964db7f8c31720981925b1e357a9fdfcc585bc2164fdef1f54d0",
		"8654a32fa120bfdb7ca02c487469070eba4b5a81b03763a2185fdf5afd756f3c", "e2767d788d1c5620f3ef21d57f2d64559ab203c044f0a5f0730b21984e77019c",
		"66d1945ceb6ef8014b1b6703cb624f058913e722f15d03225be27cb9d8aabe4a", "3fcb80eb1d5b91c491408447ac4e221fcb2254c861adbb5a178337c2750b0846",
		"22c7623bf0e850538329e3e6d9a6f9b1235350824a3feaad2580b7a853550deb", "5577d4be25f1b44487140a626c8aeca2a77507a1fc4fd466dd3a82234abb6785",
		"a5eebe3469d68c8922a1a8b5a0a2b55293b7ff424240c16feb9f51727f734516", "c07275582d0681eb07c7b51f0bca0c48c056d571b7b83d84980ab40ac7d7d720",
		"479ec3b589b14aa7290b48c2e64072e4e5b15ce395d2072a5a18b0a2cf35f3fd", "f10e2b7675dfa557d9e3188469f12d3e953c2d46dce006cd177b6ae7f465cfc0",
		"63952334b731ec91d88c54614925576f82e3610d009657368fc866e7b1efbe73", "0bcbebb39d8fe1cb3eab952c6f701656c234e462b945e2f7d4be2c80b8f2d974",
		"256472ee754ef6af096340ab1e161f58e85fb0cc7ae6e6866b9359a1657fa6c1", "88ba6f6c66fc0ef01c938569c2dd1f05475cb56444f4582d06828e77d54ffbe6",
	};
	GoInt32 error_code;
	char bufferSeed[BUFFER_SIZE];
	char bufferSec1[BUFFER_SIZE];
	char bufferSec2[BUFFER_SIZE];
	char buffer1[BUFFER_SIZE];
	char buffer2[BUFFER_SIZE];
	
	GoSlice seed = {bufferSeed, 0, BUFFER_SIZE};
	GoSlice sec1 = {bufferSec1, 0, BUFFER_SIZE};
	GoSlice sec2 = {bufferSec2, 0, BUFFER_SIZE};
	GoSlice s1   = {buffer1, 0, BUFFER_SIZE};
	GoSlice s2   = {buffer2, 0, BUFFER_SIZE};
	
	for(int i = 0; i < test_count; i++){
		seed.len = hexnstr(testArray[2 * i], bufferSeed, BUFFER_SIZE);
		sec1.len = hexnstr(testArray[2 * i + 1], bufferSec1, BUFFER_SIZE);
		error_code = SKY_secp256k1_DeterministicKeyPairIterator(seed, 
			(coin__UxArray*)&s1, (coin__UxArray*)&s2, 
			(coin__UxArray*)&sec2);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_DeterministicKeyPairIterator failed");
		cr_assert(eq(type(GoSlice), sec1, sec2), "Different hashes");
	}
}

Test(cipher_secp256k1, Test_Secp256k1_Hash){
	int test_count = 16;
	const char* testArray[] = {
		"90c56f5b8d78a46fb4cddf6fd9c6d88d6d2d7b0ec35917c7dac12c03b04e444e", "a70c36286be722d8111e69e910ce4490005bbf9135b0ce8e7a59f84eee24b88b",
		"a3b08ccf8cbae4955c02f223be1f97d2bb41d92b7f0c516eb8467a17da1e6057", "e9db072fe5817325504174253a056be7b53b512f1e588f576f1f5a82cdcad302",
		"7048eb8fa93cec992b93dc8e93c5543be34aad05239d4c036cf9e587bbcf7654", "5e9133e83c4add2b0420d485e1dcda5c00e283c6509388ab8ceb583b0485c13b",
		"6d25375591bbfce7f601fc5eb40e4f3dde2e453dc4bf31595d8ec29e4370cd80", "8d5579cd702c06c40fb98e1d55121ea0d29f3a6c42f5582b902ac243f29b571a",
		"7214b4c09f584c5ddff971d469df130b9a3c03e0277e92be159279de39462120", "3a4e8c72921099a0e6a4e7f979df4c8bced63063097835cdfd5ee94548c9c41a",
		"b13e78392d5446ae304b5fc9d45b85f26996982b2c0c86138afdac8d2ea9016e", "462efa1bf4f639ffaedb170d6fb8ba363efcb1bdf0c5aef0c75afb59806b8053",
		"9403bff4240a5999e17e0ab4a645d6942c3a7147c7834e092e461a4580249e6e", "68dd702ea7c7352632876e9dc2333142fce857a542726e402bb480cad364f260",
		"2665312a3e3628f4df0b9bc6334f530608a9bcdd4d1eef174ecda99f51a6db94", "5db72c31d575c332e60f890c7e68d59bd3d0ac53a832e06e821d819476e1f010",
		"6cb37532c80765b7c07698502a49d69351036f57a45a5143e33c57c236d841ca", "0deb20ec503b4c678213979fd98018c56f24e9c1ec99af3cd84b43c161a9bb5c",
		"8654a32fa120bfdb7ca02c487469070eba4b5a81b03763a2185fdf5afd756f3c", "36f3ede761aa683813013ffa84e3738b870ce7605e0a958ed4ffb540cd3ea504",
		"66d1945ceb6ef8014b1b6703cb624f058913e722f15d03225be27cb9d8aabe4a", "6bcb4819a96508efa7e32ee52b0227ccf5fbe5539687aae931677b24f6d0bbbd",
		"22c7623bf0e850538329e3e6d9a6f9b1235350824a3feaad2580b7a853550deb", "8bb257a1a17fd2233935b33441d216551d5ff1553d02e4013e03f14962615c16",
		"a5eebe3469d68c8922a1a8b5a0a2b55293b7ff424240c16feb9f51727f734516", "d6b780983a63a3e4bcf643ee68b686421079c835a99eeba6962fe41bb355f8da",
		"479ec3b589b14aa7290b48c2e64072e4e5b15ce395d2072a5a18b0a2cf35f3fd", "39c5f108e7017e085fe90acfd719420740e57768ac14c94cb020d87e36d06752",
		"63952334b731ec91d88c54614925576f82e3610d009657368fc866e7b1efbe73", "79f654976732106c0e4a97ab3b6d16f343a05ebfcc2e1d679d69d396e6162a77",
		"256472ee754ef6af096340ab1e161f58e85fb0cc7ae6e6866b9359a1657fa6c1", "387883b86e2acc153aa334518cea48c0c481b573ccaacf17c575623c392f78b2",
	};
	GoInt32 error_code;
	char bufferHash1[BUFFER_SIZE];
	char bufferHash2[BUFFER_SIZE];
	char bufferHash3[BUFFER_SIZE];
	GoSlice hash1 = {bufferHash1, 0, BUFFER_SIZE};
	GoSlice hash2 = {bufferHash2, 0, BUFFER_SIZE};
	GoSlice hash3 = {bufferHash3, 0, BUFFER_SIZE};
	
	for(int i = 0; i < test_count; i++){
		hash1.len = hexnstr(testArray[2 * i], bufferHash1, BUFFER_SIZE);
		hash2.len = hexnstr(testArray[2 * i + 1], bufferHash2, BUFFER_SIZE);
		error_code = SKY_secp256k1_Secp256k1Hash(hash1, (coin__UxArray*)&hash3);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_Secp256k1Hash failed");
		cr_assert(eq(type(GoSlice), hash2, hash3), "Different hashes");
	}
}


Test(cipher_secp256k1, Test_Secp256k1_Equal){
	char bufferSeed[BUFFER_SIZE];
	char bufferHash1[BUFFER_SIZE];
	char bufferHash2[BUFFER_SIZE];
	char bufferPrivate[BUFFER_SIZE];
	char bufferPublic[BUFFER_SIZE];
	
	GoSlice seed = {bufferSeed, 0, BUFFER_SIZE};
	GoSlice hash1 = {bufferHash1, 0, BUFFER_SIZE};
	GoSlice hash2 = {bufferHash2, 0, BUFFER_SIZE};
	GoSlice private = {bufferPrivate, 0, BUFFER_SIZE};
	GoSlice public = {bufferPublic, 0, BUFFER_SIZE};
	GoInt32 error_code;
	
	for(int i = 0; i < 64; i++) {
		error_code = SKY_secp256k1_RandByte( 128, (coin__UxArray*)&seed);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_RandByte failed");
		error_code = SKY_secp256k1_Secp256k1Hash(seed, (coin__UxArray*)&hash1);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_Secp256k1Hash failed");
		error_code = SKY_secp256k1_DeterministicKeyPairIterator(seed, 
			(coin__UxArray*)&hash2, 
			(coin__UxArray*)&public, 
			(coin__UxArray*)&private);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_DeterministicKeyPairIterator failed");
		cr_assert(eq(type(GoSlice), hash1, hash2), "Different hashes");
	}
}

Test(cipher_secp256k1, Test_DeterministicWalletGeneration){
	const char* pSeed = "8654a32fa120bfdb7ca02c487469070eba4b5a81b03763a2185fdf5afd756f3c";
	const char* pSecOut = "10ba0325f1b8633ca463542950b5cd5f97753a9829ba23477c584e7aee9cfbd5";
	const char* pPubOut = "0249964ac7e3fe1b2c182a2f10abe031784e374cc0c665a63bc76cc009a05bc7c6";
	
	char bufferSeed[BUFFER_SIZE];
	char bufferPrivate[BUFFER_SIZE];
	char bufferPublic[BUFFER_SIZE];
	char bufferNewSeed[BUFFER_SIZE];
	char bufferPrivateExpected[BUFFER_SIZE];
	char bufferPublicExpected[BUFFER_SIZE];
	
	GoSlice seed = {bufferSeed, 0, BUFFER_SIZE};
	GoSlice private = {bufferPrivate, 0, BUFFER_SIZE};
	GoSlice public = {bufferPublic, 0, BUFFER_SIZE};
	GoSlice newSeed = {bufferNewSeed, 0, BUFFER_SIZE};
	GoSlice privateExpected = {bufferPrivateExpected, 0, BUFFER_SIZE};
	GoSlice publicExpected = {bufferPublicExpected, 0, BUFFER_SIZE};
	
	strcpy(bufferSeed, pSeed);
	seed.len = strlen(pSeed);
	
	GoInt32 error_code;
	
	for( int i = 0; i < 1024; i++ ) {
		error_code = SKY_secp256k1_DeterministicKeyPairIterator(seed, 
			(coin__UxArray*)&newSeed, 
			(coin__UxArray*)&public, 
			(coin__UxArray*)&private);
		cr_assert(error_code == SKY_OK, "SKY_secp256k1_DeterministicKeyPairIterator failed");
		memcpy( seed.data, newSeed.data, newSeed.len);
		seed.len = newSeed.len;
	}
	
	privateExpected.len = hexnstr(pSecOut, bufferPrivateExpected, BUFFER_SIZE);
	publicExpected.len = hexnstr(pPubOut, bufferPublicExpected, BUFFER_SIZE);
	
	cr_assert(eq(type(GoSlice), privateExpected, private), "Private keyd didn\'t match");
	cr_assert(eq(type(GoSlice), public, publicExpected), "Public keyd didn\'t match");
}

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