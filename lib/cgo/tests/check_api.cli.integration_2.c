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
  unsigned char buff[BUFFER_SIZE];
  int readBits;
  // Copy file rf to f
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

/*Test(api_cli_integration, TestGenerateAddresses) {

  int lenStruct = 1;

  struct testStruct {
    unsigned char *name;
    bool encrypted;
    unsigned char *args;
    bool isUsageErr;
    unsigned char *expectOutput;
    unsigned char *goldenFile;
  };

  struct testStruct tt[lenStruct];

  unsigned char buff[BUFFER_SIZE];
  unsigned char cmd[BUFFER_SIZE];
  unsigned char expectCmd[BUFFER_SIZE];

  strcpy(cmd, "boxfort-worker generateAddresses");
  strcpy(expectCmd, "7g3M372kxwNwwQEAmrronu4anXTW8aD1XC\n");
  tt[0].name = "generateAddresses";
  tt[0].encrypted = false;
  tt[0].args = cmd;
  tt[0].expectOutput = expectCmd;
  tt[0].goldenFile = "generate-addresses.golden";

  strcpy(cmd, "boxfort-worker generateAddresses -n 2 -j");
  strcpy(expectCmd, "{\n    \"addresses\": [\n        "
                    "\"7g3M372kxwNwwQEAmrronu4anXTW8aD1XC\",\n        "
                    "\"2EDapDfn1VC6P2hx4nTH2cRUkboGAE16evV\"\n    ]\n}\n");
  tt[1].name = "generateAddresses -n 2 -j";
  tt[1].encrypted = false;
  tt[1].args = cmd;
  tt[1].expectOutput = expectCmd;
  tt[1].goldenFile = "generate-addresses-2.golden";

  for (int i = 0; i < lenStruct; i++) {
    printf("Entrando a la iteracion\n");
    createTempWalletDir(tt[i].encrypted);
    printf("paso el createTempWalletDir\n");
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
    GoString Args = {tt[i].args, strlen(tt[i].args)};
    errcode = SKY_cli_App_Run(&appHandle, Args);
    // Get redirected standard output
    getStdOut(output, BUFFER_SIZE);
    cr_assert(errcode == SKY_OK, "SKY_cli_App_Run failed");

    // printf("Outputs : %s\n", output);
    // printf("expectOutput : %s\n", tt[i].expectOutput);

    cr_assert(eq(u8, output, tt[i].expectOutput));

    json_value *w;
    char buffName[BUFFER_SIZE];
    char buffDir[BUFFER_SIZE];
    GoString WalletName = {"WALLET_NAME", 11};
    GoString WalletPath = {"WALLET_DIR", 10};
    GoString pathname = {buffName, 0};
    GoString Dir = {buffDir, 0};
    SKY_cli_Getenv(WalletPath, ((GoString_ *)&Dir));
    SKY_cli_Getenv(WalletName, ((GoString_ *)&pathname));
    unsigned char filePath[BUFFER_SIZE];
    strcpy(filePath, Dir.p);
    strcat(filePath, "/");
    strcat(filePath, pathname.p);
    w = loadJsonFile(filePath);
    cr_assert(w != NULL, "Failed to json parse");
    registerJsonFree(w);

    json_value *expect;
    unsigned char goldenFileURL[BUFFER_SIZE];
    strcpy(goldenFileURL, TEST_DATA_DIR);
    strcat(goldenFileURL, "/");
    strcat(goldenFileURL, tt[i].goldenFile);
    expect = loadJsonFile(goldenFileURL);
    cr_assert(expect != NULL, "Failed to json parse");
    registerJsonFree(expect);

    cr_assert(compareJsonValues(expect, w) == 1);
  }
};*/
