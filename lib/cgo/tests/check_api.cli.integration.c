#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"
#include "jsmn.h"

#define BUFFER_SIZE 1024
#define STRING_SIZE 1024

Test(api_cli_integration, TestStableShowConfig) {
	char output[BUFFER_SIZE];
	char wallet_dir[STRING_SIZE];
	char data_dir[STRING_SIZE];
	char wallet_name[STRING_SIZE];
	
	Config__Handle configHandle;
	App__Handle appHandle;
	const char* str = "boxfort-worker showConfig";
	GoString showConfigCommand = {str, strlen(str) };
	GoUint32 errcode;
	
	errcode = SKY_cli_LoadConfig( &configHandle );
	cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
	errcode = SKY_cli_NewApp( &configHandle, &appHandle );
	cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
	
	//Redirect standard output to a pipe
	redirectStdOut();
	errcode = SKY_cli_App_Run( &appHandle, showConfigCommand );
	//Get redirected standard output
	getStdOut(output, BUFFER_SIZE);
	cr_assert(errcode == SKY_OK, "SKY_cli_App_Run failed");
	
	//JSON parse output
	jsmn_parser parser;
	jsmntok_t tokens[128];
	jsmn_init(&parser);
	int tokens_count = jsmn_parse(&parser, output, strlen(output), tokens, sizeof(tokens)/sizeof(tokens[0]));
	cr_assert(tokens_count > 0, "Failed to json parse");
	cr_assert(tokens[0].type == JSMN_OBJECT, "Failed to json parse");
	int result = json_get_string(output, "wallet_directory", tokens, tokens_count,
		wallet_dir, STRING_SIZE);
	cr_assert(result == 1, "Failed to get json value");
	result = json_get_string(output, "data_directory", tokens, tokens_count,
		data_dir, STRING_SIZE);
	cr_assert(result == 1, "Failed to get json value");
	result = string_has_suffix(wallet_dir, ".skycoin/wallets");
	cr_assert(result == 1, "Wallet dir must end in .skycoin/wallets");
	result = string_has_suffix(data_dir, ".skycoin");
	cr_assert(result == 1, "Data dir must end in .skycoin");
	result = string_has_prefix(wallet_dir, data_dir);
	cr_assert(result == 1, "Data dir must be prefix of wallet dir");
}