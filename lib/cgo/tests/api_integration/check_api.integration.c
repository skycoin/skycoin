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
#include <unistd.h>

#define NODE_ADDRESS "SKYCOIN_NODE_HOST"
#define NODE_ADDRESS_DEFAULT "http://127.0.0.1:46420"
#define BUFFER_SIZE 1024
#define STABLE 1

#define STRING_SIZE 128
#define JSON_FILE_SIZE 4096
#define JSON_BIG_FILE_SIZE 32768
#define TEST_DATA_DIR "src/api/integration/testdata/"

#define NORMAL_TESTS
#define DECRYPTION_TESTS
#define DECRYPT_WALLET_TEST

TestSuite(api_integration, .init = setup, .fini = teardown);

char* getNodeAddress(){
	if( STABLE ){
		return NODE_ADDRESS_DEFAULT;
	} else {
		GoString_ nodeAddress;
		memset(&nodeAddress, 0, sizeof(GoString_));
		GoString  nodeEnvName = {NODE_ADDRESS, strlen(NODE_ADDRESS)};
		int result = SKY_cli_Getenv(nodeEnvName, &nodeAddress);
		cr_assert(result == SKY_OK, "Couldn\'t get node address from enviroment");
		registerMemCleanup((void*)nodeAddress.p);
		if( strcmp(nodeAddress.p, "") == 0){
			return NODE_ADDRESS_DEFAULT;
		}
		return (char*)nodeAddress.p;
	}
}

json_value* loadGoldenFile(const char* file){
	char path[STRING_SIZE];
	if(strlen(TEST_DATA_DIR) + strlen(file) < STRING_SIZE){
		strcpy(path, TEST_DATA_DIR);
		strcat(path, file);
		return loadJsonFile(path);
	}
	return NULL;
}

GoString* createGoStringSlice(char** pStrings, int count, GoSlice* slice){
	GoString* goStrings = malloc(sizeof(GoString) * count);
	cr_assert(goStrings != NULL, "Error creating GoString Slice");
	registerMemCleanup( goStrings );
	for(int i = 0; i < count; i++){
		goStrings[i].p = pStrings[i];
		goStrings[i].n = strlen(pStrings[i]);
	}
	slice->data = goStrings;
	slice->len = count;
	slice->cap = count;
	return goStrings;
}

int compareObjectsByHandle(Handle h1, Handle h2){
	GoString_ jsonResult1, jsonResult2;
	int result;
	memset(&jsonResult1, 0, sizeof(GoString_));
	memset(&jsonResult2, 0, sizeof(GoString_));

	result = SKY_JsonEncode_Handle(h1, &jsonResult1);
	cr_assert(result == SKY_OK, "Couldn\'t json encode");
	registerMemCleanup((void*)jsonResult1.p);

	result = SKY_JsonEncode_Handle(h2, &jsonResult2);
	cr_assert(result == SKY_OK, "Couldn\'t json encode");
	registerMemCleanup((void*)jsonResult2.p);

	json_char* json1 = (json_char*)jsonResult1.p;
	json_value* value1 = json_parse(json1, strlen(jsonResult1.p));
	cr_assert(value1 != NULL, "json_parse failed");
	registerJsonFree(value1);

	json_char* json2 = (json_char*)jsonResult2.p;
	json_value* value2 = json_parse(json2, strlen(jsonResult2.p));
	cr_assert(value2 != NULL, "json_parse failed");
	registerJsonFree(value2);

	int equal = compareJsonValues(value1, value2);

	freeRegisteredMemCleanup((void*)jsonResult1.p);
	freeRegisteredMemCleanup((void*)jsonResult2.p);
	freeRegisteredJson(value1);
	freeRegisteredJson(value2);
	return equal;
}

int compareObjectNodeWithGoldenFile(Handle handle,
					const char* golden_file, char* nodePath){
	GoString_ jsonResult;
	int result;
	memset(&jsonResult, 0, sizeof(GoString_));

	result = SKY_JsonEncode_Handle(handle, &jsonResult);
	cr_assert(result == SKY_OK, "Couldn\'t json encode");
	registerMemCleanup((void*)jsonResult.p);

	json_char* json = (json_char*)jsonResult.p;
	json_value* value = json_parse(json, strlen(jsonResult.p));
	cr_assert(value != NULL, "json_parse failed");
	registerJsonFree(value);

	if( nodePath != NULL ){
		value = get_json_value( value, nodePath, json_object );
		cr_assert(value != NULL, "Could\'t find node in json struct");
	}

	json_value* golden_value = loadGoldenFile(golden_file);
	cr_assert(golden_value != NULL, "loadGoldenFile failed");
	registerJsonFree(golden_value);

	int equal = compareJsonValues(value, golden_value);

	freeRegisteredJson(value);
	freeRegisteredJson(golden_value);
	freeRegisteredMemCleanup((void*)jsonResult.p);

	return equal;
}

int compareObjectWithGoldenFile(Handle handle, const char* golden_file){
	return compareObjectNodeWithGoldenFile(handle, golden_file, NULL);
}

int createWallet(Client__Handle clientHandle,
		int encrypt, char* password,
		char* seed, int max_seed_length,
		WalletResponse__Handle* responseHandle){
	char label[10];
	int result;
	if(seed[0] == 0){
		cr_assert(max_seed_length > 64, "Seed buffer is too short");
		unsigned char buff[64];
		GoSlice slice = { buff, 0, 64 };
		SKY_cipher_RandByte( 32, (coin__UxArray*)&slice );
		b64_encode(buff, 32, seed);
	}
	strncpy(label, seed, 6);
	GoString strSeed = {seed, strlen(seed)};
	GoString strLabel = {label, strlen(label)};
	if( encrypt ){
		GoString strPassword = {password, strlen(password)};
		result = SKY_api_Client_CreateEncryptedWallet(
			clientHandle, strSeed, strLabel, strPassword, 0,
			responseHandle);
	} else {
		result = SKY_api_Client_CreateUnencryptedWallet(
			clientHandle, strSeed, strLabel, 0,
			responseHandle);
	}
	return result;
}

#ifdef NORMAL_TESTS

Test(api_integration, TestVersion) {
	GoString_ version;
	memset(&version, 0, sizeof(GoString_));

	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;
	Handle versionDataHandle;

	int result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );
	result = SKY_api_Client_Version( clientHandle, &versionDataHandle );
	cr_assert(result == SKY_OK, "Couldn\'t get version");
	registerHandleClose( versionDataHandle );
	result = SKY_JsonEncode_Handle(versionDataHandle, &version);
	cr_assert(result == SKY_OK, "Couldn\'t json encode version");
	registerMemCleanup((void*)version.p);
	int versionLength = strlen(version.p);
	cr_assert(versionLength > 0, "Invalid version data");
}

Test(api_integration, TestStableCoinSupply) {
	GoString_ jsonResult;
	memset(&jsonResult, 0, sizeof(GoString_));

	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;
	Handle coinSupplyHandle;

	int result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	result = SKY_api_Client_CoinSupply( clientHandle, &coinSupplyHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_CoinSupply failed");
	registerHandleClose( coinSupplyHandle );

	result = SKY_JsonEncode_Handle(coinSupplyHandle, &jsonResult);
	cr_assert(result == SKY_OK, "Couldn\'t json encode");
	registerMemCleanup((void*)jsonResult.p);

	json_char* json = (json_char*)jsonResult.p;
	json_value* value = json_parse(json, strlen(jsonResult.p));
	cr_assert(value != NULL, "json_parse failed");
	registerJsonFree(value);

	json_value* json_golden = loadGoldenFile("coinsupply.golden");
	cr_assert(json_golden != NULL, "loadGoldenFile failed");
	registerJsonFree(json_golden);

	int equal = compareJsonValues(value, json_golden);
	cr_assert(equal, "Output different than expected");
}

typedef struct{
	char*  	golden_file;
	int		addresses_count;
	char**	addresses;
	int 	hashes_count;
	char**  hashes;
	int		failure;
}test_output;

Test(api_integration, TestStableOutputs) {
	int result;
	GoString_ jsonResult;
	memset(&jsonResult, 0, sizeof(GoString_));

	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	char* address_1[] = {
		"ALJVNKYL7WGxFBSriiZuwZKWD4b7fbV1od",
		"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf",
		"qxmeHkwgAMfwXyaQrwv9jq3qt228xMuoT5",
	};

	char* hashes_1[] = {
		"9e53268a18f8d32a44b4fb183033b49bebfe9d0da3bf3ef2ad1d560500aa54c6",
		"d91e07318227651129b715d2db448ae245b442acd08c8b4525a934f0e87efce9",
		"01f9c1d6c83dbc1c993357436cdf7f214acd0bfa107ff7f1466d1b18ec03563e",
		"fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f20",
	};
	int test_cases = 3;
	test_output tests[] = {
		{
			"outputs-noargs.golden",
			0, NULL, 0, NULL, 0
		},
		{
			"outputs-addrs.golden",
			3, address_1,
			0, NULL, 0
		},
		{
			"outputs-hashes.golden",
			0, NULL,
			4, hashes_1, 0
		},
	};
	Handle outputHandle;
	GoSlice strings;

	for(int i = 0; i < test_cases; i++){
		memset(&strings, 0, sizeof(GoSlice));
		cr_assert(tests[i].addresses_count == 0 || tests[i].hashes_count == 0);
		if(tests[i].addresses_count == 0 && tests[i].hashes_count == 0){
			result = SKY_api_Client_Outputs(clientHandle, &outputHandle);
		} else if(tests[i].addresses_count > 0){
			createGoStringSlice(tests[i].addresses,
					tests[i].addresses_count, &strings);
			result = SKY_api_Client_OutputsForAddresses(clientHandle,
											strings, &outputHandle);
		} else if(tests[i].hashes_count > 0){
			createGoStringSlice(tests[i].hashes,
					tests[i].hashes_count, &strings);
			result = SKY_api_Client_OutputsForHashes(clientHandle,
											strings, &outputHandle);
		}

		if( tests[i].failure ){
			cr_assert(result != SKY_OK, "SKY_api_Client_Outputs should have failed");
			continue;
		}
		cr_assert(result == SKY_OK, "SKY_api_Client_Outputs failed");
		registerHandleClose( outputHandle );

		result = SKY_JsonEncode_Handle(outputHandle, &jsonResult);
		cr_assert(result == SKY_OK, "Couldn\'t json encode");
		registerMemCleanup((void*)jsonResult.p);

		json_char* json = (json_char*)jsonResult.p;
		json_value* jsonOutput = json_parse(json, strlen(jsonResult.p));
		cr_assert(jsonOutput != NULL, "json_parse failed");
		registerJsonFree(jsonOutput);

		json_value* json_golden = loadGoldenFile(tests[i].golden_file);
		cr_assert(json_golden != NULL, "loadGoldenFile failed");
		registerJsonFree(json_golden);

		int equal = compareJsonValues(jsonOutput, json_golden);
		cr_assert(equal, "Output different than expected");
	}
}

typedef struct{
	char*  		golden_file;
	char* 		hash;
	GoUint64 	seq;
	int			failure;
}test_block;

Test(api_integration, TestStableBlock) {
	int test_count = 5;
	test_block tests[] = {
		{
			NULL,
			"80744ec25e6233f40074d35bf0bfdbddfac777869b954a96833cb89f44204444",
			0, 1
		},
		{
			"block-hash.golden",
			"70584db7fb8ab88b8dbcfed72ddc42a1aeb8c4882266dbb78439ba3efcd0458d",
			0, 0,
		},
		{
			"block-hash-genesis.golden",
			"0551a1e5af999fe8fff529f6f2ab341e1e33db95135eef1b2be44fe6981349f3",
			0, 0,
		},
		{
			"block-seq-0.golden",
			NULL,
			0, 0,
		},
		{
			"block-seq-100.golden",
			NULL,
			100, 0,
		},
	};

	int result;
	GoString_ jsonResult;
	memset(&jsonResult, 0, sizeof(GoString_));

	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	GoString strHash;
	Handle blockHandle;

	for(int i = 0; i < test_count; i++){
		if( tests[i].hash != NULL ){
			memset( &strHash, 0, sizeof(GoString) );
			strHash.p = tests[i].hash;
			strHash.n = strlen( tests[i].hash );
			result = SKY_api_Client_BlockByHash(clientHandle,
				strHash, &blockHandle);
		} else {
			result = SKY_api_Client_BlockBySeq(clientHandle,
				tests[i].seq, &blockHandle);
		}
		if( tests[i].failure ){
			cr_assert( result != SKY_OK, "Get Block should have failed" );
			continue;
		}
		cr_assert( result == SKY_OK, "Get Block failed" );
		registerHandleClose( blockHandle );

		result = SKY_JsonEncode_Handle(blockHandle, &jsonResult);
		cr_assert(result == SKY_OK, "Couldn\'t json encode");
		registerMemCleanup((void*)jsonResult.p);

		json_char* json = (json_char*)jsonResult.p;
		json_value* jsonOutput = json_parse(json, strlen(jsonResult.p));
		cr_assert(jsonOutput != NULL, "json_parse failed");
		registerJsonFree(jsonOutput);

		json_value* json_golden = loadGoldenFile(tests[i].golden_file);
		cr_assert(json_golden != NULL, "loadGoldenFile failed");
		registerJsonFree(json_golden);

		int equal = compareJsonValues(jsonOutput, json_golden);
		cr_assert(equal, "Output different than expected");
	}

	Handle progressHandle;
	printf("Querying every block in the blockchain");
	result = SKY_api_Client_BlockchainProgress(clientHandle, &progressHandle);
	cr_assert(result == SKY_OK, "SKY_api_Client_BlockchainProgress failed");
	registerHandleClose( progressHandle );
	GoInt64 progress;
	result = SKY_Handle_Progress_GetCurrent( progressHandle, &progress );
	cr_assert(result == SKY_OK, "SKY_Handle_Progress_GetCurrent failed");
	GoInt64 seq, blockSeq;
	Handle prevBlockHandle = 0;
	Handle blockHandle2;
	blockHandle = 0;
	GoString_ hash1, hash2, hash;
	GoString _hash;
	for(seq = 0; seq < progress; seq++){
		result = SKY_api_Client_BlockBySeq(clientHandle,
				seq, &blockHandle);
		cr_assert( result == SKY_OK, "SKY_api_Client_BlockBySeq failed" );
		registerHandleClose( blockHandle );
		result = SKY_Handle_Block_GetHeadSeq( blockHandle, &blockSeq );
		cr_assert( result == SKY_OK, "SKY_Handle_Block_GetHeadSeq failed" );
		cr_assert( seq == blockSeq, "Incorrect block sequence" );
		if( prevBlockHandle ){
			memset( &hash1, 0, sizeof(GoString_) );
			memset( &hash2, 0, sizeof(GoString_) );
			memset( &hash, 0, sizeof(GoString_) );
			result = SKY_Handle_Block_GetHeadHash( prevBlockHandle, &hash1 );
			cr_assert(result == SKY_OK, "SKY_Handle_Block_GetHeadHash failed");
			registerMemCleanup((void*)hash1.p);
			result = SKY_Handle_Block_GetPreviousBlockHash( blockHandle, &hash2 );
			cr_assert(result == SKY_OK, "SKY_Handle_Block_GetPreviousBlockHash failed");
			registerMemCleanup((void*)hash2.p);
			cr_assert(eq(type(GoString_), hash1, hash2));
			freeRegisteredMemCleanup((void*)hash1.p);
			freeRegisteredMemCleanup((void*)hash1.p);
			result = SKY_Handle_Block_GetHeadHash( blockHandle, &hash );
			registerMemCleanup((void*)hash.p);
			_hash.p = hash.p;
			_hash.n = hash.n;
			result = SKY_api_Client_BlockByHash(clientHandle,
				_hash, &blockHandle2);
			cr_assert(result == SKY_OK, "SKY_api_Client_BlockByHash failed");
			registerHandleClose( blockHandle2 );

			int equal = compareObjectsByHandle(blockHandle, blockHandle2);
			cr_assert(equal == 1);
			freeRegisteredMemCleanup((void*)hash.p);
			closeRegisteredHandle( blockHandle2 );
		}
		if( prevBlockHandle ){
			closeRegisteredHandle( prevBlockHandle );
		}
		prevBlockHandle = blockHandle;
	}
	if( blockHandle ){
		closeRegisteredHandle( blockHandle );
	}
}

Test(api_integration, TestStableBlockchainMetadata) {
	int result;
	GoString_ jsonResult;
	memset(&jsonResult, 0, sizeof(GoString_));

	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	Handle metadataHandle;
	result = SKY_api_Client_BlockchainMetadata( clientHandle,
									&metadataHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_BlockchainMetadata failed");
	registerHandleClose( metadataHandle );

	int equal = compareObjectWithGoldenFile(metadataHandle,
									"blockchain-metadata.golden");
	cr_assert(equal, "SKY_api_Client_BlockchainMetadata returned unexpected result");
}

Test(api_integration, TestStableBlockchainProgress) {
	int result;
	GoString_ jsonResult;
	memset(&jsonResult, 0, sizeof(GoString_));

	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	Handle progressHandle;
	result = SKY_api_Client_BlockchainProgress( clientHandle,
									&progressHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_BlockchainMetadata failed");
	registerHandleClose( progressHandle );

	int equal = compareObjectWithGoldenFile(progressHandle,
									"blockchain-progress.golden");
	cr_assert(equal, "SKY_api_Client_BlockchainProgress returned unexpected result");
}

typedef struct{
	int 	addresses_count;
	char** 	addresses;
	char* 	golden_file;
} test_balance;

Test(api_integration, TestStableBalance) {
	char* addr1[] = {
		"prRXwTcDK24hs6AFxj69UuWae3LzhrsPW9"
	};
	char* addr2[] = {
		"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf"
	};
	char* addr3[] = {
		"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf",
		"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf"
	};
	char* addr4[] = {
		"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf",
		"qxmeHkwgAMfwXyaQrwv9jq3qt228xMuoT5"
	};
	int result;
	int tests_count = 4;
	test_balance tests[] = {
		{
			1, addr1, "balance-noaddrs.golden"
		},
		{
			1, addr2, "balance-2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf.golden"
		},
		{
			2, addr3, "balance-2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf.golden"
		},
		{
			2, addr4, "balance-two-addrs.golden"
		},
	};

	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	int i;
	GoSlice strings;
	wallet__BalancePair balance;
	for(i = 0; i < tests_count; i++){
		memset( &strings, 0, sizeof(GoSlice) );
		createGoStringSlice( tests[i].addresses, tests[i].addresses_count,
							&strings);
		result = SKY_api_Client_Balance( clientHandle,
										strings, &balance );
		cr_assert(result == SKY_OK, "SKY_api_Client_BlockchainMetadata failed");
		json_value* json_golden = loadGoldenFile(tests[i].golden_file);
		cr_assert(json_golden != NULL, "loadGoldenFile failed");
		registerJsonFree(json_golden);
		json_value* value;
		value = get_json_value(json_golden,
							"confirmed/coins", json_integer);
		cr_assert(value != NULL, "get_json_value confirmed/coins failed");
		cr_assert(value->u.integer == balance.Confirmed.Coins);
		value = get_json_value(json_golden,
							"confirmed/hours", json_integer);
		cr_assert(value != NULL, "get_json_value confirmed/hours failed");
		cr_assert(value->u.integer == balance.Confirmed.Hours);
		value = get_json_value(json_golden,
							"predicted/coins", json_integer);
		cr_assert(value != NULL, "get_json_value predicted/coins failed");
		cr_assert(value->u.integer == balance.Predicted.Coins);
		value = get_json_value(json_golden,
							"predicted/hours", json_integer);
		cr_assert(value != NULL, "get_json_value predicted/hours failed");
		cr_assert(value->u.integer == balance.Predicted.Hours);
	}
}


Test(api_integration, TestStableUxOut) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	char* golden_file = "uxout.golden";
	char* pUxId = "fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f20";
	GoString strUxId = {pUxId, strlen(pUxId)};
	Handle uxOutHandle;
	result = SKY_api_Client_UxOut( clientHandle, strUxId, &uxOutHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_UxOut failed");
	registerHandleClose( uxOutHandle );

	int equal = compareObjectWithGoldenFile(uxOutHandle, golden_file);
	cr_assert(equal, "SKY_api_Client_UxOut returned unexpected result");
}

typedef struct{
	char* 	address;
	char*	golden_file;
	int 	failure;
}test_address_ux_out;

Test(api_integration, TestStableAddressUxOuts) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );
	int tests_count = 3;
	test_address_ux_out tests[] = {
		{"", NULL, 1},
		{"prRXwTcDK24hs6AFxj69UuWae3LzhrsPW9", "uxout-noaddr.golden", 0},
		{"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf", "uxout-addr.golden", 0},

	};
	GoString addr;
	for(int i = 0; i < tests_count; i++){
		memset(&addr, 0, sizeof(GoString));
		addr.p = tests[i].address;
		addr.n = strlen( tests[i].address );
		Handle outHandle;
		result = SKY_api_Client_AddressUxOuts( clientHandle, addr, &outHandle );
		if( tests[i].failure ){
			cr_assert(result != SKY_OK, "SKY_api_Client_AddressUxOuts should have failed");
			continue;
		} else {
			cr_assert(result == SKY_OK, "SKY_api_Client_AddressUxOuts failed");
		}
		registerHandleClose( outHandle );
		int equal = compareObjectWithGoldenFile( outHandle, tests[i].golden_file );
		cr_assert( equal == 1 );
	}
}

typedef struct{
	char* golden_file;
	GoUint64 start;
	GoUint64 end;
	int 	 failure;
}test_blockn;

Handle testBlocksHandle(Client__Handle clientHandle,
			Handle blocksHandle, GoUint64 start, GoUint64 end,
			int checkIndexes){
	GoUint64 count = 0;
	int result;
	result = SKY_Handle_Blocks_GetCount( blocksHandle, &count );
	cr_assert(result == SKY_OK, "SKY_Handle_Blocks_GetCount failed");
	if( checkIndexes ){
		if( start > end ){
			cr_assert(count == 0);
		} else {
			cr_assert(count == end - start + 1);
		}
	}
	GoUint64 i;
	GoString_ hash1, hash2, hash;
	GoString _hash;
	GoUint64 seq;
	int equal;
	for(i = 0; i < count; i++){
		Handle blockHandle = 0, previousBlockHandle = 0;
		Handle blockHandle2;
		result = SKY_Handle_Blocks_GetAt(blocksHandle, i, &blockHandle);
		cr_assert( result == SKY_OK, "Error getting block from blocks handle" );
		registerHandleClose( blockHandle );
		if( i > 0 ){
			memset(&hash1, 0, sizeof(GoString_));
			memset(&hash2, 0, sizeof(GoString_));
			memset(&hash, 0, sizeof(GoString_));
			result = SKY_Handle_Blocks_GetAt(blocksHandle, i - 1,
								&previousBlockHandle);
			cr_assert( result == SKY_OK, "Error getting previous block from blocks handle" );
			registerHandleClose( previousBlockHandle );

			result = SKY_Handle_Block_GetHeadHash(previousBlockHandle, &hash1);
			cr_assert( result == SKY_OK, "Error getting previous block hash");
			registerMemCleanup( (void*)hash1.p );
			result = SKY_Handle_Block_GetPreviousBlockHash(blockHandle, &hash2);
			cr_assert( result == SKY_OK, "Error getting previous block hash");
			registerMemCleanup( (void*)hash2.p );

			cr_assert(eq(type(GoString_), hash1, hash2));
			freeRegisteredMemCleanup( (void*)hash1.p );
			freeRegisteredMemCleanup( (void*)hash2.p );

		}

		result = SKY_Handle_Block_GetHeadHash(blockHandle, &hash);
		cr_assert( result == SKY_OK, "Error getting previous block hash");
		registerMemCleanup( (void*)hash.p );

		_hash.p = hash.p;
		_hash.n = hash.n;
		result = SKY_api_Client_BlockByHash(clientHandle,
			_hash, &blockHandle2);
		cr_assert( result == SKY_OK, "SKY_api_Client_BlockByHash failed");
		registerHandleClose( blockHandle2 );

		if( checkIndexes ){
			result = SKY_Handle_Block_GetHeadSeq( blockHandle2, &seq );
			cr_assert( result == SKY_OK, "SKY_Handle_Block_GetHeadSeq failed");
			cr_assert(seq == i + start);
		}

		equal = compareObjectsByHandle( blockHandle, blockHandle2 );
		cr_assert( equal == 1);

		freeRegisteredMemCleanup( (void*)hash.p );
		closeRegisteredHandle( blockHandle );
		closeRegisteredHandle( blockHandle2 );
		if( previousBlockHandle > 0 )
			closeRegisteredHandle( previousBlockHandle );
	}
}

Handle testBlocks(Client__Handle clientHandle,
				GoUint64 start, GoUint64 end){
	Handle blocksHandle;
	int result;
	result = SKY_api_Client_Blocks(clientHandle, start, end, &blocksHandle);
	cr_assert(result == SKY_OK, "SKY_api_Client_Blocks failed");
	registerHandleClose( blocksHandle );
	testBlocksHandle( clientHandle, blocksHandle, start, end, 1 );
	return blocksHandle;
}

Test(api_integration, TestStableBlocks) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	Handle progressHandle;
	result = SKY_api_Client_BlockchainProgress(clientHandle, &progressHandle);
	cr_assert(result == SKY_OK, "SKY_api_Client_BlockchainProgress failed");
	registerHandleClose( progressHandle );
	GoUint64 lastNBlocks = 10;
	GoUint64 current;
	result = SKY_Handle_Progress_GetCurrent( progressHandle, &current );
	cr_assert(result == SKY_OK, "SKY_Handle_Progress_GetCurrent failed");
	cr_assert( current > lastNBlocks + 1, "Progress current must be greater than 10" );
	int tests_count = 7;
	test_blockn tests[] = {
		{
			"blocks-first-10.golden", 1, 10, 0
		},
		{
			"blocks-last-10.golden", current - lastNBlocks, current, 0
		},
		{
			"blocks-first-1.golden", 1, 1, 0
		},
		{
			"blocks-all.golden", 0, current, 0
		},
		{
			"blocks-end-less-than-start.golden", 10, 9, 0
		},
		{
			NULL, -10, 9, 1
		},
		{
			NULL, 10, -9, 1
		},
	};
	Handle blocksHandle;
	int equal;
	for(int i = 0; i < tests_count; i++){
		if( tests[i].failure ){
			result = SKY_api_Client_Blocks(clientHandle,
					tests[i].start, tests[i].end, &blocksHandle);
			cr_assert(result != SKY_OK, "SKY_api_Client_Blocks should have failed");
		} else {
			blocksHandle = testBlocks(clientHandle,
					tests[i].start, tests[i].end);
			equal = compareObjectWithGoldenFile(blocksHandle,
										tests[i].golden_file);
			cr_assert(equal == 1, "SKY_api_Client_Blocks returned a value different than expected.");
			closeRegisteredHandle( blocksHandle );
		}
	}
}

Test(api_integration, TestStableLastBlocks) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	Handle blocksHandle;
	result = SKY_api_Client_LastBlocks( clientHandle, 1, &blocksHandle);
	cr_assert(result == SKY_OK, "SKY_api_Client_LastBlocks(1) failed");
	registerHandleClose( blocksHandle );

	equal = compareObjectWithGoldenFile( blocksHandle, "block-last.golden");
	cr_assert(equal == 1, "SKY_api_Client_LastBlocks(1) returned result different than expected");
	closeRegisteredHandle( blocksHandle );

	result = SKY_api_Client_LastBlocks( clientHandle, 10, &blocksHandle);
	cr_assert(result == SKY_OK, "SKY_api_Client_LastBlocks(10) failed");
	registerHandleClose( blocksHandle );
	testBlocksHandle( clientHandle, blocksHandle, 0, 0, 0);

	closeRegisteredHandle( blocksHandle );
}

Test(api_integration, TestStableNetworkConnections) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	Handle connectionsHandle;
	result = SKY_api_Client_NetworkConnections( clientHandle, &connectionsHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_NetworkConnections failed");
	registerHandleClose( connectionsHandle );

	GoUint64 connectionsCount;
	result = SKY_Handle_Connections_GetCount( connectionsHandle, &connectionsCount );
	cr_assert(result == SKY_OK, "SKY_Handle_Connections_GetCount failed");
	cr_assert( connectionsCount == 0 );

	char* pAddress = "127.0.0.1:4444";
	GoString address = { pAddress, strlen(pAddress) };
	Handle connectionHandle;
	result = SKY_api_Client_NetworkConnection( clientHandle, address, &connectionHandle );
	cr_assert(result != SKY_OK, "SKY_api_Client_NetworkConnection should have failed");
}

Test(api_integration, TestNetworkDefaultConnections) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	Strings__Handle connectionsHandle;
	result = SKY_api_Client_NetworkDefaultConnections( clientHandle, &connectionsHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_NetworkDefaultConnections failed");
	registerHandleClose( connectionsHandle );

	result = SKY_Handle_Strings_Sort(connectionsHandle);
	cr_assert(result == SKY_OK);

	equal = compareObjectWithGoldenFile( connectionsHandle,
				"network-default-connections.golden");
	cr_assert(equal == 1);
}

Test(api_integration, TestNetworkTrustedConnections) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	Strings__Handle connectionsHandle;
	result = SKY_api_Client_NetworkTrustedConnections( clientHandle, &connectionsHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_NetworkTrustedConnections failed");
	registerHandleClose( connectionsHandle );

	result = SKY_Handle_Strings_Sort(connectionsHandle);
	cr_assert(result == SKY_OK);

	equal = compareObjectWithGoldenFile( connectionsHandle,
				"network-trusted-connections.golden");
	cr_assert(equal == 1);
}

Test(api_integration, TestStableNetworkExchangeableConnections) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	Strings__Handle connectionsHandle;
	result = SKY_api_Client_NetworkExchangeableConnections( clientHandle, &connectionsHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_NetworkTrustedConnections failed");
	registerHandleClose( connectionsHandle );

	result = SKY_Handle_Strings_Sort(connectionsHandle);
	cr_assert(result == SKY_OK);

	equal = compareObjectWithGoldenFile( connectionsHandle,
				"network-exchangeable-connections.golden");
	cr_assert(equal == 1);
}

typedef struct {
	char* 	golden_file;
	char* 	txId;
	int 	failure;
} test_transaction;

Test(api_integration, TestStableTransaction) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	int tests_count = 4;
	test_transaction tests[] = {
		{
			NULL, "abcd", 1,
		},
		{
			NULL, "701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947", 1,
		},
		{
			NULL, "", 1,
		},
		{
			"genesis-transaction.golden",
			"d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add",
			0,
		},
	};
	GoString txId;
	Handle transactionHandle;
	for(int i = 0; i < tests_count; i++){
		txId.p = tests[i].txId;
		txId.n = strlen(tests[i].txId);
		result = SKY_api_Client_Transaction(
							clientHandle, txId, &transactionHandle);
		if( tests[i].failure ){
			cr_assert(result != SKY_OK, "SKY_api_Client_Transaction should have failed");
			continue;
		}
		cr_assert(result == SKY_OK, "SKY_api_Client_Transaction failed");
		registerHandleClose( transactionHandle );

		equal = compareObjectNodeWithGoldenFile(
				transactionHandle,
				tests[i].golden_file,
				"txn"); //Compare starting from this node
		cr_assert( equal == 1,
			"SKY_api_Client_Transaction returned a value different than expected" );

	}
}

typedef struct {
	char* 	golden_file;
	char** 	addresses;
	int 	addresses_count;
	int 	failure;
} test_transactions;

Test(api_integration, TestStableTransactions) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	char* addrs1[] = {
		"abcd"
	};
	char* addrs2[] = {
		"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"
	};
	char* addrs3[] = {
		"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"
	};
	char* addrs4[] = {
	};
	char* addrs5[] = {
		"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"
	};
	test_transactions tests[] = {
		{
			NULL, addrs1, 1, 1
		},
		{
			NULL, addrs2, 1, 1
		},
		{
			NULL, addrs3, 1, 1
		},
		{
			"empty-addrs.golden", addrs4, 0, 0
		},
		{
			"single-addr.golden", addrs5, 1, 0
		},
	};
	Handle transactionsHandle;
	GoSlice strings;
	int tests_count = sizeof(tests) / sizeof(test_transactions);
	for(int i = 0; i < tests_count; i++){
		memset( &strings, 0, sizeof(GoSlice) );
		createGoStringSlice( tests[i].addresses, tests[i].addresses_count,
							&strings);
		result = SKY_api_Client_Transactions( clientHandle,
						strings, &transactionsHandle);
		if( tests[i].failure ){
			cr_assert( result != SKY_OK, "SKY_api_Client_Transactions should have failed." );
			continue;
		}
		cr_assert( result == SKY_OK, "SKY_api_Client_Transactions failed" );
		registerHandleClose( transactionsHandle );
		equal = compareObjectWithGoldenFile( transactionsHandle,
										tests[i].golden_file );
		cr_assert( equal == 1, "SKY_api_Client_Transactions returned a value different than expected.");
	}
}

Test(api_integration, TestStableConfirmedTransactions) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	char* addrs1[] = {
		"abcd"
	};
	char* addrs2[] = {
		"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"
	};
	char* addrs3[] = {
		"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"
	};
	char* addrs4[] = {
	};
	char* addrs5[] = {
		"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"
	};
	test_transactions tests[] = {
		{
			NULL, addrs1, 1, 1
		},
		{
			NULL, addrs2, 1, 1
		},
		{
			NULL, addrs3, 1, 1
		},
		{
			"empty-addrs.golden", addrs4, 0, 0
		},
		{
			"single-addr.golden", addrs5, 1, 0
		},
	};
	Handle transactionsHandle;
	GoSlice strings;
	int tests_count = sizeof(tests) / sizeof(test_transactions);
	for(int i = 0; i < tests_count; i++){
		memset( &strings, 0, sizeof(GoSlice) );
		createGoStringSlice( tests[i].addresses, tests[i].addresses_count,
							&strings);
		result = SKY_api_Client_ConfirmedTransactions( clientHandle,
						strings, &transactionsHandle);
		if( tests[i].failure ){
			cr_assert( result != SKY_OK, "SKY_api_Client_ConfirmedTransactions should have failed." );
			continue;
		}
		cr_assert( result == SKY_OK, "SKY_api_Client_ConfirmedTransactions failed" );
		registerHandleClose( transactionsHandle );
		equal = compareObjectWithGoldenFile( transactionsHandle,
										tests[i].golden_file );
		cr_assert( equal == 1, "SKY_api_Client_ConfirmedTransactions returned a value different than expected.");
	}
}

Test(api_integration, TestStableUnconfirmedTransactions) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	char* addrs1[] = {
		"abcd"
	};
	char* addrs2[] = {
		"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"
	};
	char* addrs3[] = {
		"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"
	};
	char* addrs4[] = {
	};
	test_transactions tests[] = {
		{
			NULL, addrs1, 1, 1
		},
		{
			NULL, addrs2, 1, 1
		},
		{
			NULL, addrs3, 1, 1
		},
		{
			"empty-addrs-unconfirmed-txs.golden", addrs4, 0, 0
		},
	};
	Handle transactionsHandle;
	GoSlice strings;
	int tests_count = sizeof(tests) / sizeof(test_transactions);
	for(int i = 0; i < tests_count; i++){
		memset( &strings, 0, sizeof(GoSlice) );
		createGoStringSlice( tests[i].addresses, tests[i].addresses_count,
							&strings);
		result = SKY_api_Client_UnconfirmedTransactions( clientHandle,
						strings, &transactionsHandle);
		if( tests[i].failure ){
			cr_assert( result != SKY_OK, "SKY_api_Client_UnconfirmedTransactions should have failed." );
			continue;
		}
		cr_assert( result == SKY_OK, "SKY_api_Client_UnconfirmedTransactions failed" );
		registerHandleClose( transactionsHandle );
		equal = compareObjectWithGoldenFile( transactionsHandle,
										tests[i].golden_file );
		cr_assert( equal == 1, "SKY_api_Client_UnconfirmedTransactions returned a value different than expected.");
	}
}

Test(api_integration, TestStableResendUnconfirmedTransactions) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	Handle resendResultHandle;
	result = SKY_api_Client_ResendUnconfirmedTransactions(
					clientHandle, &resendResultHandle);
	cr_assert( result == SKY_OK, "SKY_api_Client_ResendUnconfirmedTransactions failed" );
	registerHandleClose( resendResultHandle );

	GoString_ jsonResult;
	memset(&jsonResult, 0, sizeof(GoString_));

	result = SKY_JsonEncode_Handle(resendResultHandle, &jsonResult);
	cr_assert(result == SKY_OK, "Couldn\'t json encode");
	registerMemCleanup((void*)jsonResult.p);

	json_value* json = json_parse( (json_char*) jsonResult.p,
							strlen(jsonResult.p) );
	cr_assert(json != NULL, "json_parse failed");
	registerJsonFree( json );
	json_value* json_txtIds =
		get_json_value_not_strict( json, "txids", json_array, 1);
	cr_assert(json_txtIds != NULL, "Error in JSON result from SKY_api_Client_ResendUnconfirmedTransactions");
	int length = 0;
	if ( json_txtIds->type == json_array )//It maybe json_null
		length = json_txtIds->u.array.length;
	cr_assert(length == 0, "SKY_api_Client_ResendUnconfirmedTransactions should have returned an empty or null array of transactions");
}

typedef struct{
	char* 	txId;
	char* 	rawTx;
	int 	failure;
} test_raw_transaction;

Test(api_integration, TestStableRawTransaction) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	test_raw_transaction tests[] = {
		{	//Invalid hex length
			"abcd", NULL, 1,
		},
		{   //Not found
			"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947",
			NULL, 1,
		},
		{   //Odd length hex string
			"abcdeffedca",
			NULL, 1,
		},
		{   //OK
			"d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add",
			"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000f8f9c644772dc5373d85e11094e438df707a42c900407a10f35a000000407a10f35a0000",
			0,
		},
	};
	int tests_count = sizeof(tests) / sizeof(test_raw_transaction);
	GoString txtId;
	GoString_ rawTx;
	GoString_ expected;
	for(int i = 0; i < tests_count; i++){
		memset( &rawTx, 0, sizeof(GoString_) );
		memset( &expected, 0, sizeof(GoString_) );
		txtId.p = tests[i].txId;
		txtId.n = strlen( tests[i].txId );
		result = SKY_api_Client_RawTransaction( clientHandle,
									txtId, &rawTx );
		if( tests[i].failure ){
			cr_assert( result != SKY_OK, "SKY_api_Client_RawTransaction should have failed" );
			continue;
		}
		expected.p = tests[i].rawTx;
		expected.n = strlen( tests[i].rawTx );
		cr_assert(result == SKY_OK, "SKY_api_Client_RawTransaction failed");
		registerMemCleanup( (void*)rawTx.p );
		cr_assert(eq(type(GoString_), rawTx, expected));
	}
}

typedef struct {
	int entropy;
	int words_count;
	int failure;
} test_new_seed;

Test(api_integration, TestWalletNewSeed) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	test_new_seed tests[] = {
		{128, 12, 0},
		{256, 24, 0},
		{100, 0, 1},
	};
	int tests_count = sizeof(tests) / sizeof(test_new_seed);
	GoString_ seed, seed2;
	for( int i = 0; i < tests_count; i++ ){
		memset( &seed, 0, sizeof(GoString_));
		memset( &seed2, 0, sizeof(GoString_));
		result = SKY_api_Client_NewSeed( clientHandle,
				tests[i].entropy, &seed);
		if( tests[i].failure ){
			cr_assert( result != SKY_OK, "SKY_api_Client_NewSeed should have failed" );
			continue;
		}
		cr_assert(result == SKY_OK, "SKY_api_Client_NewSeed failed");
		registerMemCleanup( (void*)seed.p );
		int words = count_words(seed.p, seed.n);
		cr_assert( words == tests[i].words_count, "SKY_api_Client_NewSeed incorrect words count");
		if( seed.n > 0 ){
			cr_assert(seed.p[0] != ' ' && seed.p[seed.n-1] != ' ', "Seed has extra spaces");
		}
		result = SKY_api_Client_NewSeed( clientHandle,
				tests[i].entropy, &seed2);
		cr_assert(result == SKY_OK, "SKY_api_Client_NewSeed failed");
		registerMemCleanup( (void*)seed2.p );
		//Seeds must be different every time
		cr_assert(not(eq(type(GoString_), seed, seed2)));
	}
}

typedef struct {
	char* 	address;
	char* 	golden_file;
	int		failure;
} test_address_transactions;

Test(api_integration, TestStableAddressTransactions) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	test_address_transactions tests[] = {
		{
			"ALJVNKYL7WGxFBSriiZuwZKWD4b7fbV1od",
			"address-transactions-ALJVNKYL7WGxFBSriiZuwZKWD4b7fbV1od.golden",
			0,
		},
		{
			"2b8ourW8fbTkC1yQBSLseVt6srhXvNMHvn9",
			"address-transactions-2b8ourW8fbTkC1yQBSLseVt6srhXvNMHvn9.golden",
			0,
		},
		{
			"prRXwTcDK24hs6AFxj", //Invalid address
			NULL,
			1,
		},
	};
	int tests_count = sizeof(tests) / sizeof(test_address_transactions);
	GoString address;
	Handle transactionsHandle;
	for( int i = 0; i < tests_count; i++){
		address.p = tests[i].address;
		address.n = strlen(tests[i].address);
		result = SKY_api_Client_AddressTransactions(clientHandle,
				address, &transactionsHandle);
		if( tests[i].failure ){
			cr_assert( result != SKY_OK, "SKY_api_Client_AddressTransactions should have failed" );
			continue;
		}
		cr_assert(result == SKY_OK, "SKY_api_Client_AddressTransactions failed");
		registerHandleClose( transactionsHandle );
		equal = compareObjectWithGoldenFile(transactionsHandle,
										tests[i].golden_file);
		cr_assert(equal == 1, "SKY_api_Client_AddressTransactions returned a value different that expected");
	}
}

Test(api_integration, TestStableRichlist) {
	int result, equal;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	Handle richlistHandle;
	result = SKY_api_Client_Richlist(clientHandle, NULL, &richlistHandle);
	cr_assert(result == SKY_OK, "SKY_api_Client_Richlist failed");
	registerHandleClose( richlistHandle );
	equal = compareObjectWithGoldenFile( richlistHandle, "richlist-default.golden" );
	cr_assert( equal , "Richlist default result error");

	api__RichlistParams params;
	params.N = 0;
	params.IncludeDistribution = 0;
	result = SKY_api_Client_Richlist(clientHandle, &params, &richlistHandle);
	cr_assert(result == SKY_OK, "SKY_api_Client_Richlist failed");
	registerHandleClose( richlistHandle );
	equal = compareObjectWithGoldenFile( richlistHandle, "richlist-all.golden" );
	cr_assert( equal , "Richlist all result error");

	params.N = 0;
	params.IncludeDistribution = 1;
	result = SKY_api_Client_Richlist(clientHandle, &params, &richlistHandle);
	cr_assert(result == SKY_OK, "SKY_api_Client_Richlist failed");
	registerHandleClose( richlistHandle );
	equal = compareObjectWithGoldenFile( richlistHandle, "richlist-all-include-distribution.golden" );
	cr_assert( equal , "Richlist all result error");

	params.N = 8;
	params.IncludeDistribution = 0;
	result = SKY_api_Client_Richlist(clientHandle, &params, &richlistHandle);
	cr_assert(result == SKY_OK, "SKY_api_Client_Richlist failed");
	registerHandleClose( richlistHandle );
	equal = compareObjectWithGoldenFile( richlistHandle, "richlist-8.golden" );
	cr_assert( equal , "Richlist result error");

	params.N = 150;
	params.IncludeDistribution = 1;
	result = SKY_api_Client_Richlist(clientHandle, &params, &richlistHandle);
	cr_assert(result == SKY_OK, "SKY_api_Client_Richlist failed");
	registerHandleClose( richlistHandle );
	equal = compareObjectWithGoldenFile( richlistHandle, "richlist-150-include-distribution.golden" );
	cr_assert( equal , "Richlist result error");
}

Test(api_integration, TestStableAddressCount) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	GoUint64 count;
	result = SKY_api_Client_AddressCount( clientHandle, &count );
	cr_assert(result == SKY_OK, "SKY_api_Client_AddressCount failed");
	cr_assert( count == 155 );
}

Test(api_integration, TestStablePendingTransactions) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	Handle txtHandle;
	result = SKY_api_Client_PendingTransactions( clientHandle, &txtHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_PendingTransactions failed");

	GoString_ jsonResult;
	memset(&jsonResult, 0, sizeof(GoString_));

	result = SKY_JsonEncode_Handle(txtHandle, &jsonResult);
	cr_assert(result == SKY_OK, "Couldn\'t json encode");
	registerMemCleanup((void*)jsonResult.p);

	json_value* json = json_parse( (json_char*) jsonResult.p,
							strlen(jsonResult.p) );
	cr_assert(json != NULL, "json_parse failed");
	registerJsonFree( json );
	cr_assert(json->type == json_array, "Transactions result should be slice");
	cr_assert(json->u.array.length == 0, "Transactions should be zero");
}

Test(api_integration, TestCreateWallet) {
	int result, encrypted;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	WalletResponse__Handle responseHandle;
	char seed[128];
	seed[0] = 0;
	result = createWallet( clientHandle, 0, "", seed, 128, &responseHandle );
	cr_assert( result == SKY_OK, "Create Wallet failed" );
	registerHandleClose( responseHandle );
	registerWalletClean( clientHandle, responseHandle );

	GoString_ strFullPath;
	memset( &strFullPath, 0, sizeof(GoString_) );
	result = SKY_api_Handle_Client_GetWalletFullPath(
					clientHandle, responseHandle, &strFullPath);
	cr_assert(result == SKY_OK, "SKY_api_Handle_Client_GetWalletFullPath failed");
	registerMemCleanup( (void*)strFullPath.p );
	cr_assert( access( strFullPath.p, F_OK ) != -1, "Wallet file doesn\'t exist." );

	Wallet__Handle walletHandle;
	GoString strWalletPath = {strFullPath.p, strFullPath.n};
	result = SKY_wallet_Load(strWalletPath, &walletHandle);
	cr_assert(result == SKY_OK, "SKY_wallet_Load failed");
	registerHandleClose( walletHandle );

	GoStringMap_ map;
	GoString_ strSeed, strEncrypted;
	GoString strSeedKey = {"seed", 4};
	GoString_ strCreatedSeed = {seed, strlen(seed)};
	GoString strEncryptedKey = {"encrypted", 9};
	memset(&strSeed, 0, sizeof(GoString_));
	memset(&strEncrypted, 0, sizeof(GoString_));
	result = SKY_api_Handle_GetWalletMeta( walletHandle, &map );
	cr_assert(result == SKY_OK, "SKY_api_Handle_GetWalletMeta failed");
	registerHandleClose((Handle)map);

	result = SKY_map_get(&map, strSeedKey, &strSeed);
	cr_assert(result == SKY_OK, "SKY_map_get failed");
	registerMemCleanup((void*)strSeed.p);
	cr_assert(eq(type(GoString_), strCreatedSeed, strSeed), "Different seeds");

	result = SKY_map_get(&map, strEncryptedKey, &strEncrypted);
	cr_assert(result == SKY_OK, "SKY_map_get failed");
	registerMemCleanup((void*)strEncrypted.p);
	encrypted = parseBoolean( strEncrypted.p, strEncrypted.n );
	cr_assert(encrypted == 0, "It should be not encrypted");

	GoUint32 count1, count2;
	result = SKY_api_Handle_GetWalletEntriesCount(walletHandle, &count1);
	cr_assert(result == SKY_OK, "SKY_api_Handle_GetWalletEntriesCount failed");
	result = SKY_api_Handle_Client_GetWalletResponseEntriesCount(
						responseHandle, &count2);
	cr_assert(result == SKY_OK, "SKY_api_Handle_Client_GetWalletResponseEntriesCount failed");
	cr_assert( count1 == count2, "Entries count are different" );
	cipher__Address address;
	cipher__PubKey pubkey;
	GoString_ strAddress1;
	GoString_ strPubKey1;
	GoString_ strAddress2;
	GoString_ strPubKey2;
	for( GoUint32 i = 0; i < count1; i++){
		memset( &address, 0, sizeof(cipher__Address) );
		memset( &pubkey, 0, sizeof(cipher__PubKey) );
		memset( &strAddress1, 0, sizeof(GoString_) );
		memset( &strPubKey1, 0, sizeof(GoString_) );
		memset( &strAddress2, 0, sizeof(GoString_) );
		memset( &strPubKey2, 0, sizeof(GoString_) );

		result = SKY_api_Handle_WalletGetEntry(walletHandle,
				i, &address, &pubkey);
		cr_assert( result == SKY_OK, "SKY_api_Handle_WalletGetEntry failed" );
		SKY_cipher_Address_String(&address, &strAddress1);
		registerMemCleanup( (void*)strAddress1.p );
		SKY_cipher_PubKey_Hex( &pubkey, &strPubKey1 );
		registerMemCleanup( (void*)strPubKey1.p );

		result = SKY_api_Handle_WalletResponseGetEntry(responseHandle,
				i, &strAddress2, &strPubKey2);
		cr_assert( result == SKY_OK, "SKY_api_Handle_WalletResponseGetEntry failed" );
		registerMemCleanup( (void*)strAddress2.p );
		registerMemCleanup( (void*)strPubKey2.p );

		cr_assert( eq( type(GoString_), strAddress1, strAddress2 ),
				"Entries Addresses should be equal");
		cr_assert( eq( type(GoString_), strPubKey1, strPubKey2 ) ,
				"Entries Public Keys should be equal");
	}

	//cleanRegisteredWallet( clientHandle, responseHandle );
	//Same test again but with wallet encrypted
	seed[0] = 0;
	result = createWallet( clientHandle, 1, "pwd", seed, 128, &responseHandle );
	cr_assert( result == SKY_OK, "Create Wallet failed" );
	registerHandleClose( responseHandle );
	registerWalletClean( clientHandle, responseHandle );
	GoUint8 isEncrypted;
	result = SKY_api_Handle_WalletResponseIsEncrypted( responseHandle, &isEncrypted );
	cr_assert(result == SKY_OK, "SKY_api_Handle_WalletResponseIsEncrypted");
	cr_assert( isEncrypted, "Wallet should be encrypted" );

	memset( &strFullPath, 0, sizeof(GoString_) );
	result = SKY_api_Handle_Client_GetWalletFullPath(
					clientHandle, responseHandle, &strFullPath);
	cr_assert(result == SKY_OK, "SKY_api_Handle_Client_GetWalletFullPath failed");
	registerMemCleanup( (void*)strFullPath.p );

	strWalletPath.p = strFullPath.p;
	strWalletPath.n = strFullPath.n;
	result = SKY_wallet_Load(strWalletPath, &walletHandle);
	cr_assert(result == SKY_OK, "SKY_wallet_Load failed");
	registerHandleClose( walletHandle );


	memset(&strEncrypted, 0, sizeof(GoString_));
	result = SKY_api_Handle_GetWalletMeta( walletHandle, &map );
	cr_assert(result == SKY_OK, "SKY_api_Handle_GetWalletMeta failed");
	registerHandleClose((Handle)map);

	result = SKY_map_get(&map, strEncryptedKey, &strEncrypted);
	cr_assert(result == SKY_OK, "SKY_map_get failed");
	registerMemCleanup((void*)strEncrypted.p);
	encrypted = parseBoolean( strEncrypted.p, strEncrypted.n );
	cr_assert(encrypted == 1, , "Wallet should be encrypted");

	result = SKY_api_Handle_GetWalletEntriesCount(walletHandle, &count1);
	cr_assert(result == SKY_OK, "SKY_api_Handle_GetWalletEntriesCount failed");
	result = SKY_api_Handle_Client_GetWalletResponseEntriesCount(
						responseHandle, &count2);
	cr_assert(result == SKY_OK, "SKY_api_Handle_Client_GetWalletResponseEntriesCount failed");
	cr_assert( count1 == count2, "Entries lengths are different" );

	for( GoUint32 i = 0; i < count1; i++){
		memset( &address, 0, sizeof(cipher__Address) );
		memset( &pubkey, 0, sizeof(cipher__PubKey) );
		memset( &strAddress1, 0, sizeof(GoString_) );
		memset( &strPubKey1, 0, sizeof(GoString_) );
		memset( &strAddress2, 0, sizeof(GoString_) );
		memset( &strPubKey2, 0, sizeof(GoString_) );

		result = SKY_api_Handle_WalletGetEntry(walletHandle,
				i, &address, &pubkey);
		cr_assert( result == SKY_OK, "SKY_api_Handle_WalletGetEntry failed" );
		SKY_cipher_Address_String(&address, &strAddress1);
		registerMemCleanup( (void*)strAddress1.p );
		SKY_cipher_PubKey_Hex( &pubkey, &strPubKey1 );
		registerMemCleanup( (void*)strPubKey1.p );

		result = SKY_api_Handle_WalletResponseGetEntry(responseHandle,
				i, &strAddress2, &strPubKey2);
		cr_assert( result == SKY_OK, "SKY_api_Handle_WalletResponseGetEntry failed" );
		registerMemCleanup( (void*)strAddress2.p );
		registerMemCleanup( (void*)strPubKey2.p );

		cr_assert( eq( type(GoString_), strAddress1, strAddress2 ),
				"Entries Addresses should be equal");
		cr_assert( eq( type(GoString_), strPubKey1, strPubKey2 ) ,
				"Entries Public Keys should be equal");
	}

}

Test(api_integration, TestGetWallet) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	WalletResponse__Handle responseHandle, responseHandle2;
	char seed[128];
	seed[0] = 0;
	result = createWallet( clientHandle, 0, "", seed, 128, &responseHandle );
	cr_assert( result == SKY_OK, "Create Wallet failed" );
	registerHandleClose( responseHandle );
	registerWalletClean( clientHandle, responseHandle );

	GoString_ walletFileName = {NULL, 0};
	result = SKY_api_Handle_Client_GetWalletFileName(
				responseHandle, &walletFileName);
	cr_assert(result == SKY_OK, "SKY_api_Handle_Client_GetWalletFileName failed");
	GoString _walletFileName = { walletFileName.p, walletFileName.n };

	result = SKY_api_Client_Wallet(
				clientHandle, _walletFileName, &responseHandle2 );
	cr_assert(result == SKY_OK, "SKY_api_Client_Wallet failed");
	registerHandleClose( responseHandle2 );
	int equal = compareObjectsByHandle( responseHandle, responseHandle2 );
	cr_assert(equal == 1, "Wallets object should be equal");
}

#define GET_WALLETS_COUNT 2

Test(api_integration, TestGetWallets) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	WalletResponse__Handle original_wallets[GET_WALLETS_COUNT];
	char seed[128];
	for( int i = 0; i < GET_WALLETS_COUNT; i++){
		seed[0] = 0;
		result = createWallet( clientHandle, 0, "", seed, 128, &original_wallets[i] );
		cr_assert( result == SKY_OK, "Create Wallet failed" );
		registerHandleClose( original_wallets[i] );
		registerWalletClean( clientHandle, original_wallets[i] );
	}

	Wallets__Handle walletsHandle;
	result = SKY_api_Client_Wallets( clientHandle, &walletsHandle );
	cr_assert( result == SKY_OK, "SKY_api_Client_Wallets failed" );
	registerHandleClose( walletsHandle ) ;

	GoUint32 count;
	result = SKY_api_Handle_WalletsResponseGetCount( walletsHandle, &count );
	cr_assert( result == SKY_OK, "SKY_api_Handle_WalletsResponseGetCount failed" );
	cr_expect( count >= GET_WALLETS_COUNT,
		"SKY_api_Handle_WalletsResponseGetCount returned incorrect length");

	WalletResponse__Handle w;
	GoString_ name1, name2;
	for( GoUint32 i = 0; i < GET_WALLETS_COUNT; i++){
		memset( &name1, 0, sizeof(GoString_) );
		result = SKY_api_Handle_Client_GetWalletFileName(
				original_wallets[i], &name1);
		cr_assert( result == SKY_OK, "SKY_api_Handle_Client_GetWalletFileName failed" );
		registerMemCleanup( (void*)name1.p );
		int found = 0;
		for( GoUint32 j = 0; j < count; j++ ){
			result = SKY_api_Handle_WalletsResponseGetAt(
				walletsHandle, j, &w);
			cr_assert( result == SKY_OK, "SKY_api_Handle_WalletsResponseGetAt failed" );
			registerHandleClose( w ) ;
			memset( &name2, 0, sizeof(GoString_) );
			result = SKY_api_Handle_Client_GetWalletFileName(w, &name2);
			cr_assert( result == SKY_OK, "SKY_api_Handle_Client_GetWalletFileName failed" );
			registerMemCleanup( (void*)name2.p );
			if( strcmp( name1.p, name2.p ) == 0 ){
				int equal;
				equal = compareObjectsByHandle( w, original_wallets[i]);
				cr_assert( equal, "Wallet is not what expected" );
				found = 1;
				break;
			}
		}
		cr_assert( found == 1, "Wallet not found");
	}
}

Test(api_integration, TestWalletNewAddress) {
	char seed[128];
	GoString pwd;
	for(GoInt i = 1; i <= 3; i++){
		int result;
		memset( &pwd, 0, sizeof(GoString));
		char* pNodeAddress = getNodeAddress();
		GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
		Client__Handle clientHandle;

		result = SKY_api_NewClient(nodeAddress, &clientHandle);
		cr_assert(result == SKY_OK, "Couldn\'t create client");
		registerHandleClose( clientHandle );

		int encrypt = i == 2;

		WalletResponse__Handle w;
		memset(seed, 0, 128);
		result = createWallet( clientHandle, encrypt, "pwd",
				seed, 128, &w );
		cr_assert( result == SKY_OK, "Create Wallet failed" );
		registerHandleClose( w );
		registerWalletClean( clientHandle, w );

		GoString_ _walletFileName = {NULL, 0};
		result = SKY_api_Handle_Client_GetWalletFileName(
					w, &_walletFileName);
		cr_assert( result == SKY_OK, "SKY_api_Handle_Client_GetWalletFileName failed" );
		registerMemCleanup( (void*)_walletFileName.p );
		if( encrypt ){
			pwd.p = "pwd";
			pwd.n = 3;
		}
		Strings__Handle addrsHandle;
		GoString walletFileName = {_walletFileName.p, _walletFileName.n};
		result = SKY_api_Client_NewWalletAddress( clientHandle,
				walletFileName, i, pwd, &addrsHandle);
		cr_assert( result == SKY_OK, "SKY_api_Client_NewWalletAddress failed" );
		registerHandleClose( addrsHandle );

		int seedLength = strlen(seed);
		GoSlice seedSlice = {seed, seedLength, seedLength};
		GoSlice keyPairsSlice = {NULL, 0, 0};
		result = SKY_cipher_GenerateDeterministicKeyPairs(
			seedSlice, i + 1, (coin__UxArray*)&keyPairsSlice);
		cr_assert( result == SKY_OK, "SKY_cipher_GenerateDeterministicKeyPairs failed" );
		registerMemCleanup( keyPairsSlice.data );

		GoUint32 count;
		SKY_Handle_Strings_GetCount( addrsHandle, &count );
		cr_assert( result == SKY_OK, "SKY_Handle_Strings_GetCount failed" );
		cr_assert( count == keyPairsSlice.len - 1,
			"Length mismatch with result for NewWalletAddress and GenerateDeterministicKeyPairs" );
		GoString_ strAddress;
		GoString_ strAddress2;
		cipher__Address cAddress;
		cipher__SecKey* secKey = (cipher__SecKey*)keyPairsSlice.data;


		for(GoUint32 a = 0; a < count; a++){
			memset( &strAddress, 0, sizeof(GoString_));
			memset( &cAddress, 0, sizeof(cipher__Address));
			memset( &strAddress2, 0, sizeof(GoString_));
			secKey++;
			SKY_cipher_AddressFromSecKey(secKey, &cAddress);
			SKY_cipher_Address_String( &cAddress, &strAddress2 );
			registerMemCleanup( (void*)strAddress2.p );
			result = SKY_Handle_Strings_GetAt( addrsHandle, a, &strAddress );
			registerMemCleanup( (void*)strAddress.p );
			cr_assert( eq(type(GoString_), strAddress, strAddress2),
				"Addresses are different");
		}
	}
}

Test(api_integration, TestStableWalletBalance) {
	int result;
	char seed[128];
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	WalletResponse__Handle w;
	strcpy(seed, "casino away claim road artist where blossom warrior demise royal still palm");
	result = createWallet( clientHandle, 0, NULL,
			seed, 128, &w );
	cr_assert( result == SKY_OK, "Create wallet failed" );
	registerHandleClose( w );
	registerWalletClean( clientHandle, w );
	GoString_ _name = {NULL, 0};
	result = SKY_api_Handle_Client_GetWalletFileName(w, &_name);
	cr_assert(result == SKY_OK, "SKY_api_Handle_Client_GetWalletFileName failed");
	GoString name = {_name.p, _name.n };
	wallet__BalancePair balance;
	result = SKY_api_Client_WalletBalance( clientHandle, name, &balance);
	cr_assert(result == SKY_OK, "SKY_api_Client_WalletBalance failed");

	json_value* json_golden = loadGoldenFile("wallet-balance.golden");
	cr_assert(json_golden != NULL, "loadGoldenFile failed");
	registerJsonFree(json_golden);
	json_value* value;
	value = get_json_value(json_golden,
						"confirmed/coins", json_integer);
	cr_assert(value != NULL, "get_json_value confirmed/coins failed");
	cr_assert(value->u.integer == balance.Confirmed.Coins,
						"balance.Confirmed.Coins incorrect");
	value = get_json_value(json_golden,
						"confirmed/hours", json_integer);
	cr_assert(value != NULL, "get_json_value confirmed/hours failed");
	cr_assert(value->u.integer == balance.Confirmed.Hours,
						"balance.Confirmed.Hours incorrect");
	value = get_json_value(json_golden,
						"predicted/coins", json_integer);
	cr_assert(value != NULL, "get_json_value predicted/coins failed");
	cr_assert(value->u.integer == balance.Predicted.Coins,
						"balance.Predicted.Coins incorrect");
	value = get_json_value(json_golden,
						"predicted/hours", json_integer);
	cr_assert(value != NULL, "get_json_value predicted/hours failed");
	cr_assert(value->u.integer == balance.Predicted.Hours,
						"balance.Predicted.Hours incorrect");

	cleanRegisteredWallet( clientHandle, w );

}

Test(api_integration, TestWalletUpdate) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	WalletResponse__Handle w;
	char seed[128];
	seed[0] = 0;
	result = createWallet( clientHandle, 0, "", seed, 128, &w );
	cr_assert( result == SKY_OK, "Create Wallet failed" );
	registerHandleClose( w );
	registerWalletClean( clientHandle, w );

	GoString_ _name;
	memset( &_name, 0, sizeof(GoString_) );
	result = SKY_api_Handle_Client_GetWalletFileName(w, &_name);
	cr_assert( result == SKY_OK, "SKY_api_Handle_Client_GetWalletFileName failed" );
	registerMemCleanup( (void*)_name.p );
	GoString name = { _name.p, _name.n };
	GoString newName = { "new wallet", 10 };
	result = SKY_api_Client_UpdateWallet( clientHandle, name, newName );
	cr_assert( result == SKY_OK, "SKY_api_Client_UpdateWallet failed" );
	WalletResponse__Handle w2;
	result = SKY_api_Client_Wallet( clientHandle, name, &w2 );
	cr_assert( result == SKY_OK, "SKY_api_Client_Wallet failed" );
	registerHandleClose( w2 );
	GoString_ strLabel1 = { NULL, 0 };
	GoString_ strLabel2 = { "new wallet", 10 };
	result = SKY_api_Handle_Client_GetWalletLabel(w2, &strLabel1);
	cr_assert( result == SKY_OK, "SKY_api_Handle_Client_GetWalletLabel failed" );
	registerMemCleanup( (void*)strLabel1.p );
	cr_assert( eq( type(GoString_), strLabel1, strLabel2 ),
			"SKY_api_Handle_Client_GetWalletLabel returned a different label");
}

Test(api_integration, TestStableWalletTransactions) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	WalletResponse__Handle w;
	char seed[128];
	seed[0] = 0;
	result = createWallet( clientHandle, 0, "", seed, 128, &w );
	cr_assert( result == SKY_OK, "Create Wallet failed" );
	registerHandleClose( w );
	registerWalletClean( clientHandle, w );

	GoString_ _name;
	memset( &_name, 0, sizeof(GoString_) );
	result = SKY_api_Handle_Client_GetWalletFileName(w, &_name);
	cr_assert( result == SKY_OK, "SKY_api_Handle_Client_GetWalletFileName failed" );
	registerMemCleanup( (void*)_name.p );
	GoString name = { _name.p, _name.n };

	Handle transHandle;
	result = SKY_api_Client_WalletTransactions( clientHandle,
								name, &transHandle );
	cr_assert( result == SKY_OK, "SKY_api_Client_WalletTransactions failed" );
	registerHandleClose( transHandle );
	int equal = compareObjectWithGoldenFile( transHandle,
					"wallet-transactions.golden" );
	cr_assert( equal, "SKY_api_Client_WalletTransactions returned unexpected value" );
}

Test(api_integration, TestWalletFolderName) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	Handle folderHandle;
	result = SKY_api_Client_WalletFolderName(
				clientHandle, &folderHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_WalletFolderName failed");
	registerHandleClose( folderHandle );
	GoString_ strAddress = {NULL, 0};
	result = SKY_api_Handle_GetWalletFolderAddress(
			folderHandle, &strAddress );
	cr_assert(result == SKY_OK, "Get Wallet Folder Address failed");
	registerMemCleanup( (void*)strAddress. p );
	cr_assert(strAddress.p != NULL, "Folder Address is null");
	cr_assert(strAddress.n > 0, "Folder Address is empty");
}

#endif

#ifdef DECRYPTION_TESTS

Test(api_integration, TestEncryptWallet) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	WalletResponse__Handle w, w2, w3, dw;
	char seed[128];
	seed[0] = 0;
	result = createWallet( clientHandle, 0, "", seed, 128, &w );
	cr_assert( result == SKY_OK, "Create Wallet failed" );
	registerHandleClose( w );
	registerWalletClean( clientHandle, w );

	GoString_ _name;
	memset( &_name, 0, sizeof(GoString_) );
	result = SKY_api_Handle_Client_GetWalletFileName(w, &_name);
	cr_assert( result == SKY_OK, "Error getting wallet file name" );
	registerMemCleanup( (void*)_name.p );
	GoString name = { _name.p, _name.n };
	GoString pwd = { "pwd", 3 };
	result = SKY_api_Client_EncryptWallet( clientHandle, name, pwd, &w2);
	cr_assert( result == SKY_OK, "Error encrypting wallet" );
	registerHandleClose( w2 );

	result = SKY_api_Client_EncryptWallet( clientHandle, name, pwd, &w3);
	cr_expect( result != SKY_OK, "Encrypting wallet twice should fail" );

	Handle folderHandle;
	result = SKY_api_Client_WalletFolderName(
				clientHandle, &folderHandle );
	cr_expect(result == SKY_OK, "SKY_api_Client_WalletFolderName failed");
	registerHandleClose( folderHandle );
	GoString_ strAddress = {NULL, 0};
	result = SKY_api_Handle_GetWalletFolderAddress(
			folderHandle, &strAddress );
	cr_expect(result == SKY_OK, "Get Wallet Folder Address failed");
	registerMemCleanup( (void*)strAddress. p );

	char fullPath[256];
	strncpy( fullPath, strAddress. p, 255 );
	int length = strlen(fullPath);
	if( (length  == 0 || fullPath[length-1] != '/') && length < 254)
		strcat(fullPath, "/");
	if( length + name.n < 254 )
		strcat( fullPath, name.p );

	Wallet__Handle walletHandle;
	GoString strWalletPath = {fullPath, strlen(fullPath)};
	result = SKY_wallet_Load(strWalletPath, &walletHandle);
	cr_expect(result == SKY_OK, "SKY_wallet_Load failed");
	registerHandleClose( walletHandle );

	GoStringMap_ map;
	GoString_ strSeed, strLastSeed, strSecrets;
	GoString strSeedKey = {"seed", 4};
	GoString strSecretsKey = {"secrets", 7};
	GoString strLastSeedKey = {"lastSeed", 8};
	memset(&strSeed, 0, sizeof(GoString_));
	memset(&strLastSeed, 0, sizeof(GoString_));
	memset(&strSecrets, 0, sizeof(GoString_));
	result = SKY_api_Handle_GetWalletMeta( walletHandle, &map );
	cr_expect(result == SKY_OK, "SKY_api_Handle_GetWalletMeta failed");
	registerHandleClose((Handle)map);

	result = SKY_map_get(&map, strSeedKey, &strSeed);
	cr_expect(result == SKY_OK, "SKY_map_get failed");
	registerMemCleanup((void*)strSeed.p);
	cr_expect( strSeed.n == 0, "Seed in encrypted wallet should be empty" );

	result = SKY_map_get(&map, strLastSeedKey, &strLastSeed);
	cr_expect(result == SKY_OK, "SKY_map_get failed");
	registerMemCleanup((void*)strLastSeed.p);
	cr_expect( strSeed.n == 0, "Last Seed in encrypted wallet should be empty" );

	result = SKY_map_get(&map, strSecretsKey, &strSecrets);
	cr_expect(result == SKY_OK, "SKY_map_get failed");
	registerMemCleanup((void*)strSecrets.p);
	cr_expect( strSecrets.n > 0, "Secrets in encrypted wallet shouldn\'t be empty" );

	result = SKY_api_Client_DecryptWallet( clientHandle, name, pwd, &dw);
	cr_assert( result == SKY_OK, "Error decrypting wallet" );
	registerHandleClose( dw );

	int equal = compareObjectsByHandle( w, dw );
	cr_assert( equal, "Decrypted wallet should be equal to the original");
}

#endif

#ifdef DECRYPT_WALLET_TEST

Test(api_integration, TestDecryptWallet) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	WalletResponse__Handle w, w2, w3, dw;
	char seed[128];
	seed[0] = 0;
	result = createWallet( clientHandle, 1, "pwd", seed, 128, &w );
	cr_assert( result == SKY_OK, "Create Wallet failed" );
	registerHandleClose( w );
	registerWalletClean( clientHandle, w );

	GoString_ _name;
	memset( &_name, 0, sizeof(GoString_) );
	result = SKY_api_Handle_Client_GetWalletFileName(w, &_name);
	cr_assert( result == SKY_OK, "Error getting wallet file name" );
	registerMemCleanup( (void*)_name.p );
	GoString name = { _name.p, _name.n };
	GoString pwd = { "pwd", 3 };
	GoString emptyPwd = { "", 0 };
	GoString wrongPwd = { "pwd1", 4 };

	//TODO: Fails if decryption is tried with wrong passwords
	result = SKY_api_Client_DecryptWallet( clientHandle, name, wrongPwd, &dw);
	cr_expect( result != SKY_OK, "Can\'t decrypt wallet with wrong password" );

	result = SKY_api_Client_DecryptWallet( clientHandle, name, emptyPwd, &dw);
	cr_expect( result != SKY_OK, "Can\'t decrypt wallet with empty password" );

	result = SKY_api_Client_DecryptWallet( clientHandle, name, pwd, &dw);
	cr_assert( result == SKY_OK, "Error decrypting wallet" );
	registerHandleClose( dw );

	GoUint8 	isEncrypted;
	GoString_ 	cryptoType = {NULL, 0};
	result = SKY_api_Handle_WalletResponseIsEncrypted( dw, &isEncrypted );
	cr_assert( result == SKY_OK, "SKY_api_Handle_WalletResponseIsEncrypted failed" );
	cr_assert( isEncrypted == 0, "Wallet should not be encrypted" );

	result = SKY_api_Handle_WalletResponseGetCryptoType( dw, &cryptoType );
	cr_assert( result == SKY_OK, "SKY_api_Handle_WalletResponseGetCryptoType failed" );
	cr_assert( cryptoType.n == 0, "CryptoType field should be empty" );

	GoString_ strFullPath = { NULL, 0 };
	result = SKY_api_Handle_Client_GetWalletFullPath(clientHandle, w,
						&strFullPath);
	cr_assert(result == SKY_OK, "SKY_api_Handle_Client_GetWalletFullPath failed");
	registerMemCleanup( (void*)strFullPath.p );

	Wallet__Handle lwHandle;
	GoString strWalletPath = { strFullPath.p, strFullPath.n};
	result = SKY_wallet_Load(strWalletPath, &lwHandle);
	cr_assert(result == SKY_OK, "SKY_wallet_Load failed");
	registerHandleClose( lwHandle );

	GoStringMap_ map;
	GoString_ strSeed;
	GoString strSeedKey = {"seed", 4};
	memset(&strSeed, 0, sizeof(GoString_));
	result = SKY_api_Handle_GetWalletMeta( lwHandle, &map );
	cr_assert(result == SKY_OK, "SKY_api_Handle_GetWalletMeta failed");
	registerHandleClose((Handle)map);

	result = SKY_map_get(&map, strSeedKey, &strSeed);
	cr_assert(result == SKY_OK, "SKY_map_get failed");
	registerMemCleanup((void*)strSeed.p);
	GoString_ prevSeed = {seed, strlen(seed)};
	cr_assert( eq( type(GoString_), prevSeed, strSeed ) );
	GoUint32 count;

	result = SKY_api_Handle_GetWalletEntriesCount( lwHandle, &count );
	cr_assert(result == SKY_OK, "SKY_api_Handle_GetWalletEntriesCount failed");
	cr_assert( count == 1 );

	cipher__Address laddress;
	cipher__PubKey lpubkey;
	memset( &laddress, 0, sizeof(cipher__Address) );
	memset( &lpubkey, 0, sizeof(cipher__PubKey) );
	result = SKY_api_Handle_WalletGetEntry( lwHandle, 0,
					&laddress, &lpubkey);
	cr_assert(result == SKY_OK, "SKY_api_Handle_WalletGetEntry failed");
	int seedLength = strlen(seed);
	GoSlice seedSlice = {seed, seedLength, seedLength};
	GoSlice lSeed = {NULL, 0, 0};
	GoSlice seckeysSliece = { NULL, 0, 0};
	result = SKY_cipher_GenerateDeterministicKeyPairsSeed(
					seedSlice, 1,
					(coin__UxArray*)&lSeed,
					(coin__UxArray*)&seckeysSliece);
	cr_assert(result == SKY_OK, "SKY_cipher_GenerateDeterministicKeyPairsSeed failed");
	registerMemCleanup( (void*)lSeed.data );
	registerMemCleanup( (void*)seckeysSliece.data );
	char hexSeed[256];
	strnhexlower( (unsigned char*)lSeed.data,
		hexSeed, lSeed.len);
	GoString_ lastSeed = {NULL, 0};
	result = SKY_api_Handle_GetWalletLastSeed( lwHandle, &lastSeed );
	cr_assert(result == SKY_OK, "SKY_api_Handle_GetWalletLastSeed failed");
	registerMemCleanup( (void*)lastSeed.p );
	GoString_ strLSeed = { hexSeed, strlen(hexSeed) };
	cr_assert( eq( type(GoString_), strLSeed, lastSeed ), "Seeds different" );
	cipher__PubKey gpubkey;
	result = SKY_cipher_PubKeyFromSecKey(
				(cipher__SecKey*)seckeysSliece.data,
				&gpubkey);
	cr_assert(result == SKY_OK, "SKY_cipher_PubKeyFromSecKey failed");
	cipher__Address gaddress;
	SKY_cipher_AddressFromPubKey(&gpubkey, &gaddress);
	GoString_ strGAddress = {NULL, 0};
	SKY_cipher_Address_String( &gaddress, &strGAddress );
	registerMemCleanup( (void*)strGAddress.p );

	GoString_ strCAddress = {NULL, 0};
	GoString_ strCPubkey = {NULL, 0};

	result = SKY_api_Handle_WalletResponseGetEntry( w, 0,
				&strCAddress, &strCPubkey);
	cr_assert(result == SKY_OK, "SKY_api_Handle_WalletResponseGetEntry failed");
	registerMemCleanup( (void*)strCAddress.p );
	registerMemCleanup( (void*)strCPubkey.p );
	cr_assert( eq( type(GoString_), strGAddress, strCAddress ) );

	GoString_ strLAddress = {NULL, 0};
	SKY_cipher_Address_String( &laddress, &strLAddress );
	registerMemCleanup( (void*)strLAddress.p );
	cr_assert( eq( type(GoString_), strLAddress, strCAddress ) );
}

#endif

//TODO: How to disable wallet seed API
/*
Test(api_integration, TestGetWalletSeedDisabledAPI) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	WalletResponse__Handle w, w2, w3, dw;
	char seed[128];
	seed[0] = 0;
	result = createWallet( clientHandle, 1, "pwd", seed, 128, &w );
	cr_assert( result == SKY_OK, "Create Wallet failed" );
	registerHandleClose( w );
	registerWalletClean( clientHandle, w );

	GoString_ _name;
	memset( &_name, 0, sizeof(GoString_) );
	result = SKY_api_Handle_Client_GetWalletFileName(w, &_name);
	cr_assert( result == SKY_OK, "Error getting wallet file name" );
	registerMemCleanup( (void*)_name.p );
	GoString name = { _name.p, _name.n };
	GoString strPwd = { "pwd", 3 };
	GoString_ strSeed = { NULL, 0 };
	result = SKY_api_Client_GetWalletSeed( &clientHandle, name, strPwd, &strSeed );
	cr_assert( result != SKY_OK );
}*/


#ifdef DECRYPTION_TESTS
Test(api_integration, TestGetWalletSeedDisabledAPI) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	WalletResponse__Handle w;
	char seed[128];
	seed[0] = 0;
	result = createWallet( clientHandle, 1, "pwd", seed, 128, &w );
	cr_assert( result == SKY_OK, "Create Wallet failed" );
	registerHandleClose( w );
	registerWalletClean( clientHandle, w );

	cr_assert( seed[0] != 0, "Generated seed can\'t be empty");

	GoString_ _name;
	memset( &_name, 0, sizeof(GoString_) );
	result = SKY_api_Handle_Client_GetWalletFileName(w, &_name);
	cr_assert( result == SKY_OK, "Error getting wallet file name" );
	registerMemCleanup( (void*)_name.p );
	GoString name = { _name.p, _name.n };
	GoString strPwd = { "pwd", 3 };
	GoString_ strSeed = { NULL, 0 };
	result = SKY_api_Client_GetWalletSeed( clientHandle, name, strPwd, &strSeed );
	cr_assert( result == SKY_OK, "SKY_api_Client_GetWalletSeed failed" );

	GoString_ strPrevSeed = { seed, strlen(seed) };
	cr_assert( eq(type(GoString_), strSeed, strPrevSeed) );

	GoString wrongWallet = { "w.wlt", 5 };
	result = SKY_api_Client_GetWalletSeed( clientHandle, wrongWallet, strPwd, &strSeed );
	cr_assert( result != SKY_OK, "SKY_api_Client_GetWalletSeed must mail with wrong wallet" );

	GoString wrongPassword = { "wrong password", 14 };
	result = SKY_api_Client_GetWalletSeed( clientHandle, name, wrongPassword, &strSeed );
	cr_assert( result != SKY_OK, "SKY_api_Client_GetWalletSeed must mail with wrong password" );

	GoString emptyPassword = { "", 0 };
	result = SKY_api_Client_GetWalletSeed( clientHandle, name, emptyPassword, &strSeed );
	cr_assert( result != SKY_OK, "SKY_api_Client_GetWalletSeed must mail with empty password" );

	WalletResponse__Handle w2;
	seed[0] = 0;
	result = createWallet( clientHandle, 0, NULL, seed, 128, &w2 );
	cr_assert( result == SKY_OK, "Create Wallet failed" );
	registerHandleClose( w2 );
	registerWalletClean( clientHandle, w2 );

	memset( &_name, 0, sizeof(GoString_) );
	result = SKY_api_Handle_Client_GetWalletFileName(w2, &_name);
	cr_assert( result == SKY_OK, "Error getting wallet file name" );
	registerMemCleanup( (void*)_name.p );
	GoString name2 = { _name.p, _name.n };
	memset( &strSeed, 0, sizeof( GoString_ ) ) ;
	result = SKY_api_Client_GetWalletSeed( clientHandle, name2, strPwd, &strSeed );
	cr_assert( result != SKY_OK, "Can\'t get seed on unencrypted wallet failed" );
}

#endif

#ifdef NORMAL_TESTS
Test(api_integration, TestStableHealth) {
	int result;
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;

	result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );

	Handle healthHandle;
	result = SKY_api_Client_Health( clientHandle, &healthHandle );
	cr_assert(result == SKY_OK, "SKY_api_Client_Health failed");
	registerHandleClose( healthHandle );

	GoString_ strJson = { NULL, 0 };
	result = SKY_JsonEncode_Handle( healthHandle, &strJson );
	cr_assert(result == SKY_OK, "Couldn\'t json encode");
	registerMemCleanup( (void*) strJson.p );
	json_value* json = json_parse( (json_char*) strJson.p, strlen(strJson.p) );
	cr_assert(json != NULL, "Couldn\'t json parse");
	registerJsonFree( json );

	//Check JSON values
	int number;
	json_value* value;
	value = get_json_value( json,
				"blockchain/unspents", json_integer);
	cr_assert(value != NULL, "Error getting json value blockchain/unspent");
	cr_assert(value->u.integer > 0, "Health blockchain unspend must be greater than 0");
	value = get_json_value( json,
				"blockchain/head/seq", json_integer);
	cr_assert(value != NULL, "Error getting json value blockchain/head/seq");
	cr_assert(value->u.integer > 0, "Health blockchain head seq must be greater than 0");
	value = get_json_value( json,
				"blockchain/head/timestamp", json_integer);
	cr_assert(value != NULL, "Error getting json value blockchain/head/timestamp");
	cr_assert(value->u.integer > 0, "Health blockchain head time must be non zero");
	value = get_json_value( json,
				"version/version", json_string);
	cr_assert(value != NULL, "Error getting json value version/version");
	cr_assert(value->u.string.length > 0, "Health version must be non empty");
	value = get_json_value( json,
				"uptime", json_string);
	cr_assert(value != NULL, "Error getting json value uptime");
	cr_assert(value->u.string.length > 0, "Health uptime must be non empty");
	number = atoi( value->u.string.ptr );
	cr_assert( number > 0, "Health uptime must be greater than zero");
	value = get_json_value( json,
				"open_connections", json_integer);
	cr_assert(value != NULL, "Error getting json value open_connections");
	cr_assert(value->u.integer == 0, "Health open connections must be zero");
	value = get_json_value( json,
				"blockchain/time_since_last_block", json_string);
	cr_assert(value != NULL, "Error getting json value blockchain/time_since_last_block");
	cr_assert(value->u.string.length > 0, "Health blockchain/time_since_last_block must be non empty");
	number = atoi( value->u.string.ptr );
	cr_assert( number > 0, "Health blockchain/time_since_last_block must be greater than zero");
	value = get_json_value( json,
				"version/commit", json_string);
	cr_assert(value != NULL, "Error getting json value version/commit");
	cr_assert(value->u.string.length > 0, "Health version/commit must be non empty");
	value = get_json_value( json,
				"version/branch", json_string);
	cr_assert(value != NULL, "Error getting json value version/branch");
	cr_assert(value->u.string.length > 0, "Health version/branch must be non empty");
}
#endif
