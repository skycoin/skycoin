#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

#define NODE_ADDRESS "SKYCOIN_NODE_HOST"
#define NODE_ADDRESS_DEFAULT "http://127.0.0.1:46420"
#define BUFFER_SIZE 1024
#define STABLE 1

#define STRING_SIZE 128
#define JSON_FILE_SIZE 4096
#define JSON_BIG_FILE_SIZE 32768
#define TEST_DATA_DIR "src/api/integration/testdata/"

typedef struct{
	char*  	golden_file;
	int		addresses_count;
	char**	addresses;
	int 	hashes_count;
	char**  hashes;
	int		failure;
}test_out_put;

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
	result = SKY_api_Client_Version( &clientHandle, &versionDataHandle );
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
	
	result = SKY_api_Client_CoinSupply( &clientHandle, &coinSupplyHandle );
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
	test_out_put tests[] = {
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
			result = SKY_api_Client_Outputs(&clientHandle, &outputHandle);
		} else if(tests[i].addresses_count > 0){
			createGoStringSlice(tests[i].addresses, 
					tests[i].addresses_count, &strings);
			result = SKY_api_Client_OutputsForAddresses(&clientHandle, 
														strings, &outputHandle);
		} else if(tests[i].hashes_count > 0){
			createGoStringSlice(tests[i].hashes, 
					tests[i].hashes_count, &strings);
			result = SKY_api_Client_OutputsForHashes(&clientHandle, 
														strings, &outputHandle);
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