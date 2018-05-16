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
	const char *temp = "/tmp/wallet-data-dir";
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
	strcpy(fullUrl,TEST_DATA_DIR);
	strcat(fullUrl,walletPath);
	rf = fopen(fullUrl,"rb");
	unsigned char buff[BUFFER_SIZE];
	int readBits;

	while ((readBits = fread(buff, 1, BUFFER_SIZE, rf)))
		fwrite(buff, 1, readBits, f);

	fclose(rf);
	fclose(f);

	GoString WalletDir = {"WALLET_DIR", 10};
	GoString Dir = {temp, strlen(temp)};
	SKY_cli_Setenv(WalletDir, Dir);
	GoString WalletPath = {"WALLET_NAME", 11};
	GoString pathname = {walletPath, strlen(walletPath)};
	SKY_cli_Setenv(WalletPath, pathname);
	fclose(f);
	fclose(rf);
};

Test(api_cli_integration, TestGenerateAddresses) {

	int lenStruct = 1;

	struct testStruct {
		unsigned char *name;
		bool encrypted;
		GoString args;
		bool isUsageErr;
		GoSlice expectOutput;
		unsigned char *goldenFile;
	};

	struct testStruct tt[lenStruct];
	GoString args = {"boxfort-worker generateAddresses", 32};
	GoSlice expectOutput = {"7g3M372kxwNwwQEAmrronu4anXTW8aD1XC\n", 36, 36};
	tt[0].name = "generateAddresses";
	tt[0].encrypted = false;
	tt[0].args = args;
	tt[0].expectOutput = expectOutput;
	tt[0].goldenFile = "generate-addresses.golden";

	for (int i = 0; i < lenStruct; i++) {
		createTempWalletDir(tt[i].encrypted);
		char output[BUFFER_SIZE];
		Config__Handle configHandle;
		App__Handle appHandle;
		GoUint32 errcode;

		errcode = SKY_cli_LoadConfig(&configHandle);
		cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
		registerHandleClose(configHandle);
		errcode = SKY_cli_NewApp(&configHandle, &appHandle);
		cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
		registerHandleClose(appHandle);

    // Redirect standard output to a pipe
		redirectStdOut();
		errcode = SKY_cli_App_Run(&appHandle, tt[i].args);
    // Get redirected standard output
		getStdOut(output, BUFFER_SIZE);
		cr_assert(errcode == SKY_OK, "SKY_cli_App_Run failed");

		printf("%s\n", output);

		GoSlice Outputs = {output, strlen(output), strlen(output)};
		cr_assert(eq(type(GoSlice), Outputs, expectOutput));

	}
};
