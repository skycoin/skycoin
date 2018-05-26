#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"
#include "json.h"

#define BUFFER_SIZE 1024
#define STRING_SIZE 128
#define JSON_FILE_SIZE 4096
#define JSON_BIG_FILE_SIZE 32768
#define TEST_DATA_DIR "src/cli/integration/testdata/"

TestSuite(api_cli_integration, .init = setup, .fini = teardown);

int useCSRF(){
	GoUint32 errcode;
	
	GoString strCSRFVar = {
		"USE_CSRF",
		8
	};
	GoString_ crsf;
	errcode = SKY_cli_Getenv(strCSRFVar, &crsf);
	cr_assert(errcode == SKY_OK, "SKY_cli_Getenv failed");
	int length = strlen(crsf.p);
	int result = 0;
	if(length == 1){
		result = crsf.p[0] == '1' || crsf.p[0] == 't' || crsf.p[0] == 'T';
	} else {
		result = strcmp(crsf.p, "true") == 0 || 
			strcmp(crsf.p, "True") == 0 ||
			strcmp(crsf.p, "TRUE") == 0;
	}
	free((void*)crsf.p);
	return result;
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

Test(api_cli_integration, TestStableShowConfig) {
	char output[BUFFER_SIZE];
	
	
	Config__Handle configHandle;
	App__Handle appHandle;
	const char* str = "boxfort-worker showConfig";
	GoString showConfigCommand = {str, strlen(str) };
	GoUint32 errcode;
	
	configHandle = SKY_cli_LoadConfig();
	cr_assert(configHandle != 0, "SKY_cli_LoadConfig failed");
	registerHandleClose( configHandle );
	errcode = SKY_cli_NewApp( &configHandle, &appHandle );
	cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
	registerHandleClose( appHandle );
	
	//Redirect standard output to a pipe
	redirectStdOut();
	errcode = SKY_cli_App_Run( &appHandle, showConfigCommand );
	//Get redirected standard output
	getStdOut(output, BUFFER_SIZE);
	cr_assert(errcode == SKY_OK, "SKY_cli_App_Run failed");
	
	//JSON parse output
	json_char* json;
    json_value* value;
	json_value* json_str;
	int result;
	json = (json_char*)output;
	value = json_parse(json, strlen(output));
	cr_assert(value != NULL, "Failed to json parse");
	registerJsonFree(value);
	
	json_value* wallet_dir = json_get_string(value, "wallet_directory");
	cr_assert(wallet_dir != NULL, "Failed to get json value");
	json_value* data_dir = json_get_string(value, "data_directory");
	cr_assert(data_dir != NULL, "Failed to get json value");
	json_value* wallet_name = json_get_string(value, "wallet_name");
	cr_assert(wallet_name != NULL, "Failed to get json value");
	json_value* coin_name = json_get_string(value, "coin");
	cr_assert(coin_name != NULL, "Failed to get json value");
	json_value* rpc_address = json_get_string(value, "rpc_address");
	cr_assert(rpc_address != NULL, "Failed to get json value");
	
	result = string_has_suffix(wallet_dir->u.string.ptr, ".skycoin/wallets");
	cr_assert(result == 1, "Wallet dir must end in .skycoin/wallets");
	result = string_has_suffix(data_dir->u.string.ptr, ".skycoin");
	cr_assert(result == 1, "Data dir must end in .skycoin");
	result = string_has_prefix(wallet_dir->u.string.ptr, data_dir->u.string.ptr);
	cr_assert(result == 1, "Data dir must be prefix of wallet dir");
	
	json_set_string(wallet_dir, "IGNORED/.skycoin/wallets");
	json_set_string(data_dir, "IGNORED/.skycoin");
	//Ignore the rpc address
	json_set_string(rpc_address, "http://127.0.0.1:46420");
	
	const char* golden_file = "show-config.golden";
	if ( useCSRF() ){
		golden_file = "show-config-use-csrf.golden";
	}
	json_value* json_golden = loadGoldenFile(golden_file);
	cr_assert(json_golden != NULL, "loadGoldenFile failed");
	registerJsonFree(json_golden);
	int equal = compareJsonValues(value, json_golden);
	cr_assert(equal, "Output from command different than expected");
}
