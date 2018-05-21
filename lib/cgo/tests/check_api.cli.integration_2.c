#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "json.h"
#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"
#include <sys/stat.h>

#define BUFFER_SIZE 1024
#define STRING_SIZE 128
#define JSON_FILE_SIZE 4096
#define JSON_BIG_FILE_SIZE 32768
#define TEST_DATA_DIR "src/cli/integration/testdata/"
#define SKYCOIN_NODE_HOST "http://127.0.0.1:6420"
#define stableWalletName "integration-test.wlt"
#define stableEncryptWalletName "integration-test-encrypted.wlt"


// createTempWalletDir creates a temporary wallet dir,
// sets the WALLET_DIR environment variable.
// Returns wallet dir path and callback function to clean up the dir.
void createTempWalletDir(bool encrypt) {
	const char *temp = "build/libskycoin/wallet-data-dir";
	int valueMkdir = mkdir(temp, S_IRWXU);

  // Copy the testdata/$stableWalletName to the temporary dir.
	unsigned char walletPath[BUFFER_SIZE];
	if (encrypt) {
		strcat(walletPath, stableEncryptWalletName);
	} else {
		strcat(walletPath, stableWalletName);
	}
	unsigned char pathnameURL[BUFFER_SIZE];
	strcpy(pathnameURL, temp);
	strcat(pathnameURL, "/");
	strcat(pathnameURL, walletPath);

	FILE *rf;
	FILE *f;
	f = fopen(pathnameURL, "wb");
	unsigned char fullUrl[BUFFER_SIZE];
	strcpy(fullUrl, TEST_DATA_DIR);
	strcat(fullUrl, walletPath);
	rf = fopen(fullUrl, "rb");
	unsigned char buff[2048];
	int readBits;
  // Copy file rf to f
	if (f && rf) {
		while ((readBits = fread(buff, 1, 2048, rf)))
			fwrite(buff, 1, readBits, f);

		fclose(rf);
		fclose(f);

		GoString WalletDir = {"WALLET_DIR", 10};
		GoString Dir = {temp, strlen(temp)};
		SKY_cli_Setenv(WalletDir, Dir);
		GoString WalletPath = {"WALLET_NAME", 11};
		GoString pathname = {walletPath, strlen(walletPath)};
		SKY_cli_Setenv(WalletPath, pathname);
	}
};

Test(api_cli_integration, TestGenerateAddresses) {

	int lenStruct = 1;

	struct testStruct {
		unsigned char *name;
		bool encrypted;
		char args[BUFFER_SIZE];
		bool isUsageErr;
		unsigned char *expectOutput;
		unsigned char *goldenFile;
	};

	struct testStruct tt[lenStruct];

	unsigned char buff[BUFFER_SIZE];
	unsigned char cmd[BUFFER_SIZE];
	unsigned char expectCmd[BUFFER_SIZE];

	strcpy(cmd, "boxfort-worker generateAddresses");
  // strcpy(expectCmd, "7g3M372kxwNwwQEAmrronu4anXTW8aD1XC\n");
	tt[0].name = "generateAddresses";
	tt[0].encrypted = false;
	strcpy(tt[0].args, cmd);
	tt[0].expectOutput = "7g3M372kxwNwwQEAmrronu4anXTW8aD1XC\n";
	tt[0].goldenFile = "generate-addresses.golden";

	strcpy(expectCmd, "{\n    \"addresses\": [\n        "
		"\"7g3M372kxwNwwQEAmrronu4anXTW8aD1XC\",\n        "
		"\"2EDapDfn1VC6P2hx4nTH2cRUkboGAE16evV\"\n    ]\n}\n");
	tt[1].name = "generateAddresses -n 2 -j";
	tt[1].encrypted = false;
	strcpy(tt[1].args, "boxfort-worker generateAddresses -n 2 -j");
	tt[1].expectOutput = expectCmd;
	tt[1].goldenFile = "generate-addresses-2.golden";

	strcpy(expectCmd, "{\n    \"addresses\": [\n        "
		"\"7g3M372kxwNwwQEAmrronu4anXTW8aD1XC\"\n    ]\n}\n");
	tt[2].name = "generateAddresses in encrypted wallet";
	tt[2].encrypted = true;
	strcpy(tt[2].args, "boxfort-worker generateAddresses -p pwd -j");
	tt[2].expectOutput = expectCmd;
	tt[2].goldenFile = "generate-addresses-encrypted.golden";

	GoUint32 errcode;

	char buffName[BUFFER_SIZE];
	char buffDir[BUFFER_SIZE];
	GoString WalletName = {"WALLET_NAME", 11};
	GoString WalletPath = {"WALLET_DIR", 10};
	GoString pathname = {buffName, 0};
	GoString Dir = {buffDir, 0};
	unsigned char filePath[BUFFER_SIZE];
	unsigned char goldenFileURL[BUFFER_SIZE];

  // for (int i = 0; i < lenStruct; i++) {
	int i = 0;
	char output[BUFFER_SIZE];
	createTempWalletDir(tt[i].encrypted);
	Config__Handle configHandle;
	App__Handle appHandle;
	errcode = SKY_cli_LoadConfig(&configHandle);
	cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
	registerHandleClose(configHandle);
	errcode = SKY_cli_NewApp(&configHandle, &appHandle);
	cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
	registerHandleClose(appHandle);

  // Redirect standard output to a pipe
	redirectStdOut();
	unsigned char buffArgs[BUFFER_SIZE];
	GoString Args = {buffArgs, 0};
	strcpy(Args.p, tt[i].args);
	Args.n = strlen(tt[i].args);
	errcode = SKY_cli_App_Run(&appHandle, Args);
  // Get redirected standard output
	getStdOut(output, BUFFER_SIZE);
	cr_assert(errcode == SKY_OK, "SKY_cli_App_Run failed");

  // printf("Outputs by %s: %s\n",tt[i].name, output);
  // printf("expectOutput by %s : %s\n",tt[i].name, tt[i].expectOutput);

	GoSlice Outputs = {output, strlen(output), strlen(output)};
	GoSlice expectOutput = {tt[i].expectOutput, strlen(tt[i].expectOutput),
		strlen(tt[i].expectOutput)};
		cr_assert(eq(type(GoSlice), Outputs, expectOutput));

		SKY_cli_Getenv(WalletPath, ((GoString_ *)&Dir));
		SKY_cli_Getenv(WalletName, ((GoString_ *)&pathname));

		strcpy(filePath, Dir.p);
		strcat(filePath, "/");
		strcat(filePath, pathname.p);

		json_value *w = loadJsonFile(filePath);
		cr_assert(w != NULL, "Failed to json parse in test :%s", tt[i].name);
		registerJsonFree(w);

		strcpy(goldenFileURL, TEST_DATA_DIR);
		strcat(goldenFileURL, tt[i].goldenFile);
		json_value *expect = loadJsonFile(goldenFileURL);
		cr_assert(expect != NULL, "Failed to json parse in test :%s", tt[i].name);
		registerJsonFree(expect);

		cr_assert(compareJsonValues(expect, w) == 1, "compareJsonValues  in : %s",
			tt[i].name);
		strcpy(output, "");
  // }
	};

	Test(api_cli_integration, TestVerifyAddress) {

		int lenStruct = 3;

		struct testStruct {
			unsigned char *name;
			const char args[BUFFER_SIZE];
			int error;
		};

		struct testStruct tt[lenStruct];

		tt[0].name = "valid skycoin address";
		strcpy(tt[0].args,
			"boxfort-worker verifyAddress 2Kg3eRXUhY6hrDZvNGB99DKahtrPDQ1W9vN");
		tt[0].error = SKY_OK;

		tt[1].name = "valid skycoin address";
		strcpy(tt[1].args,
			"boxfort-worker verifyAddress 2KG9eRXUhx6hrDZvNGB99DKahtrPDQ1W9vn");
		tt[1].error = SKY_ERROR;

		tt[2].name = "invalid bitcoin address";
		strcpy(tt[2].args,
			"boxfort-worker verifyAddress 1Dcb9gpaZpBKmjqjCsiBsP3sBW1md2kEM2");
		tt[2].error = SKY_ERROR;

		GoInt32 errcode;
		char *output;
		for (int i = 0; i < lenStruct; ++i) {

			Config__Handle configHandle;
			App__Handle appHandle;
			errcode = SKY_cli_LoadConfig(&configHandle);
			cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
			registerHandleClose(configHandle);
			errcode = SKY_cli_NewApp(&configHandle, &appHandle);
			cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
			registerHandleClose(appHandle);

    // Redirect standard output to a pipe
			unsigned char buffArgs[BUFFER_SIZE];
			GoString Args = {buffArgs, 0};
			strcpy(Args.p, tt[i].args);
			Args.n = strlen(tt[i].args);
			redirectStdOut();
			errcode = SKY_cli_App_Run(&appHandle, Args);
    // Get redirected standard output
			getStdOut(output, BUFFER_SIZE);
			cr_assert(errcode == tt[i].error, "Error in %s", tt[i].name);
		}
	}


	Test(cli_integration,TestDecodeRawTransaction){


		int lenStruct = 1;

		struct testStruct {
			unsigned char *name;
			char rawTx[BUFFER_SIZE];
			char errMsg[BUFFER_SIZE];
			unsigned char *goldenFile;
			const char args[BUFFER_SIZE];
		};

		struct testStruct tt[lenStruct];

		tt[0].name = "success";
		strcpy(tt[0].rawTx , "2601000000a1d3345ac47f897f24084b1c6b9bd6e03fc92887050d0748bdab5e639c1fdcd401000000a2a10f07e0e06cf6ba3e793b3186388a126591ee230b3f387617f1ccb6376a3f18e094bd3f7719aa8191c00764f323872f5192da393852bd85dab70b13409d2b01010000004d78de698a33abcfff22391c043b57a56bb0efbdc4a5b975bf8e7889668896bc0400000000bae12bbf671abeb1181fc85f1c01cdfee55deb97980c9c0a00000000543600000000000000373bb3675cbf3880bba3f3de7eb078925b8a72ad0095ba0a000000001c12000000000000008829025fe45b48f29795893a642bdaa89b2bb40e40d2df03000000001c12000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec09c0fc9b01000000001b12000000000000");
		tt[0].goldenFile = "decode-raw-transaction.golden";
		strcpy(tt[0].errMsg , "");
		strcpy(tt[0].args,
			"decodeRawTransaction 2601000000a1d3345ac47f897f24084b1c6b9bd6e03fc92887050d0748bdab5e639c1fdcd401000000a2a10f07e0e06cf6ba3e793b3186388a126591ee230b3f387617f1ccb6376a3f18e094bd3f7719aa8191c00764f323872f5192da393852bd85dab70b13409d2b01010000004d78de698a33abcfff22391c043b57a56bb0efbdc4a5b975bf8e7889668896bc0400000000bae12bbf671abeb1181fc85f1c01cdfee55deb97980c9c0a00000000543600000000000000373bb3675cbf3880bba3f3de7eb078925b8a72ad0095ba0a000000001c12000000000000008829025fe45b48f29795893a642bdaa89b2bb40e40d2df03000000001c12000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec09c0fc9b01000000001b12000000000000");

		tt[1].name = "invalid raw transaction";
		strcpy(tt[1].rawTx , "2601000000a1d");
		tt[1].goldenFile = "decode-raw-transaction.golden";
		strcpy(tt[1].errMsg , "invalid raw transaction: encoding/hex: odd length hex string\nencoding/hex: odd length hex string\n");
		strcpy(tt[1].args,
			"boxfort-worker decodeRawTransaction 2601000000a1d");


		GoInt32 errcode;
		char *output;
		for (int i = 0; i < lenStruct; ++i) {

			Config__Handle configHandle;
			App__Handle appHandle;
			errcode = SKY_cli_LoadConfig(&configHandle);
			cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
			registerHandleClose(configHandle);
			errcode = SKY_cli_NewApp(&configHandle, &appHandle);
			cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
			registerHandleClose(appHandle);

    // Redirect standard output to a pipe
			unsigned char buffArgs[BUFFER_SIZE];
			GoString Args = {buffArgs, 0};
			strcpy(Args.p, tt[i].args);
			Args.n = strlen(tt[i].args);
			redirectStdOut();
			errcode = SKY_cli_App_Run(&appHandle, Args);
    // Get redirected standard output
			getStdOut(output, BUFFER_SIZE);
			cr_assert(errcode == SKY_OK, "Error in %s", tt[i].name);

			//JSON parse output
			json_char* json;
			json_value* value;
			json_value* json_str;
			int result;
			json = (json_char*)output;
			printf("The outputs %s in iteration %d\n", output,i);
			printf("The args %s in iteration %d\n", tt[i].args,i);
			value = json_parse(json, strlen(output));
			cr_assert(value != NULL, "Failed to json parse in %s",tt[i].name);
			registerJsonFree(value);
		}
	}
