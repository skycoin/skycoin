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