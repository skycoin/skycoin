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
#include <unistd.h>

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

  if (valueMkdir == -1) {
    int errr = system("rm -r build/libskycoin/wallet-data-dir/*.*");
  }

  // Copy the testdata/$stableWalletName to the temporary dir.
  char walletPath[JSON_BIG_FILE_SIZE];
  if (encrypt) {
    strcpy(walletPath, stableEncryptWalletName);
  } else {
    strcpy(walletPath, stableWalletName);
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
  strcpy(walletPath, "");
};

int getCountWord(const char *str) {
  int len = 0;
  do {
    str = strpbrk(str, " "); // find separator
    if (str)
      str += strspn(str, " "); // skip separator
    ++len;                     // increment word count
  } while (str && *str);

  return len;
}

int getCountStringInString(const char *source, const char *str) {
  // TODO:Not implement
}

Test(api_cli_integration, TestGenerateAddresses) {
  int lenStruct = 6;

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
  unsigned char expectCmd[BUFFER_SIZE];

  tt[0].name = "generateAddresses";
  tt[0].encrypted = false;
  tt[0].args = "boxfort-worker generateAddresses";
  tt[0].expectOutput = "7g3M372kxwNwwQEAmrronu4anXTW8aD1XC\n";
  tt[0].goldenFile = "generate-addresses.golden";
  tt[0].isUsageErr = false;

  tt[1].name = "generateAddresses -n 2 -j";
  tt[1].encrypted = false;
  tt[1].args = "boxfort-worker generateAddresses -n 2 -j";
  tt[1].expectOutput = "{\n    \"addresses\": [\n        "
  "\"7g3M372kxwNwwQEAmrronu4anXTW8aD1XC\",\n        "
  "\"2EDapDfn1VC6P2hx4nTH2cRUkboGAE16evV\"\n    ]\n}\n";
  tt[1].goldenFile = "generate-addresses-2.golden";
  tt[1].isUsageErr = false;

  tt[2].name = "generateAddresses -n -2 -j";
  tt[2].encrypted = false;
  tt[2].args = "boxfort-worker generateAddresses -n -2 -j";
  tt[2].expectOutput = "Error: invalid value \"-2\" for flag -n: "
  "strconv.ParseUint: parsing \"-2\": invalid syntax";
  tt[2].goldenFile = "generate-addresses-2.golden";
  tt[2].isUsageErr = true;

  tt[3].name = "generateAddresses in encrypted wallet with invalid password";
  tt[3].encrypted = true;
  tt[3].args = "boxfort-worker generateAddresses -p invalid password -j";
  tt[3].expectOutput = "invalid password\n";
  tt[3].goldenFile = "generate-addresses-encrypted.golden";
  tt[3].isUsageErr = true;

  tt[4].name = "generateAddresses in encrypted wallet";
  tt[4].encrypted = true;
  tt[4].args = "boxfort-worker generateAddresses -p pwd -j";
  tt[4].expectOutput = "{\n    \"addresses\": [\n        "
  "\"7g3M372kxwNwwQEAmrronu4anXTW8aD1XC\"\n    ]\n}\n";
  tt[4].goldenFile = "generate-addresses-encrypted.golden";
  tt[5].isUsageErr = false;

  tt[5].name = "generateAddresses in unencrypted wallet with password";
  tt[5].encrypted = false;
  tt[5].args = "boxfort-worker generateAddresses -p pwd";
  tt[5].expectOutput = "wallet is not encrypted\n";
  tt[5].goldenFile = "generate-addresses.golden";
  tt[5].isUsageErr = true;

  GoUint32 errcode;

  char buffName[BUFFER_SIZE];
  char buffDir[BUFFER_SIZE];
  GoString WalletName = {"WALLET_NAME", 11};
  GoString WalletPath = {"WALLET_DIR", 10};
  GoString pathname = {buffName, 0};
  GoString Dir = {buffDir, 0};
  unsigned char filePath[BUFFER_SIZE];
  unsigned char goldenFileURL[BUFFER_SIZE];
  unsigned char output[JSON_BIG_FILE_SIZE];
  for (int i = 0; i < lenStruct; i++) {
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
    GoString Args = {tt[i].args, strlen(tt[i].args)};
    errcode = SKY_cli_App_Run(&appHandle, Args);
    // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    if (!tt[i].isUsageErr) {
      cr_assert(errcode == SKY_OK, "SKY_cli_App_Run failed in %s", tt[i].name);
      GoSlice Outputs = {output, strlen(tt[i].expectOutput),
       strlen(tt[i].expectOutput)};
       GoSlice expectOutput = {tt[i].expectOutput, strlen(tt[i].expectOutput),
        strlen(tt[i].expectOutput)};
        cr_assert(eq(type(GoSlice), Outputs, expectOutput), "Error in %s",
          tt[i].name);

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

      // Compare value the json
        json_value *meta_w;
        json_value *meta_expect;

        meta_w = get_json_value(w, "meta/coin", json_string);
        meta_expect = get_json_value(expect, "meta/coin", json_string);
        cr_assert(compareJsonValues(meta_expect, meta_w) == 1,
          "compareJsonValues meta/coin in : %s", tt[i].name);

        meta_w = get_json_value(w, "meta/filename", json_string);
        meta_expect = get_json_value(expect, "meta/filename", json_string);
        cr_assert(compareJsonValues(meta_expect, meta_w) == 1,
          "compareJsonValues meta/filename in : %s", tt[i].name);

        meta_w = get_json_value(w, "meta/filename", json_string);
        meta_expect = get_json_value(expect, "meta/filename", json_string);
        cr_assert(compareJsonValues(meta_expect, meta_w) == 1,
          "compareJsonValues meta/filename in : %s", tt[i].name);
        meta_w = get_json_value(w, "meta/seed", json_string);
        meta_expect = get_json_value(expect, "meta/seed", json_string);
        cr_assert(compareJsonValues(meta_expect, meta_w) == 1,
          "compareJsonValues meta/seed in : %s", tt[i].name);
      }
      strcpy(output, "");
    }
  };

  Test(api_cli_integration, TestVerifyAddress) {
    int lenStruct = 3;

    struct testStruct {
      unsigned char *name;
      unsigned char *args;
      int error;
    };

    struct testStruct tt[lenStruct];

    tt[0].name = "valid skycoin address";
    tt[0].args =
    "boxfort-worker verifyAddress 2Kg3eRXUhY6hrDZvNGB99DKahtrPDQ1W9vN";
    tt[0].error = SKY_OK;

    tt[1].name = "valid skycoin address";
    tt[1].args =
    "boxfort-worker verifyAddress 2KG9eRXUhx6hrDZvNGB99DKahtrPDQ1W9vn";
    tt[1].error = SKY_ERROR;

    tt[2].name = "invalid bitcoin address";
    tt[2].args =
    "boxfort-worker verifyAddress 1Dcb9gpaZpBKmjqjCsiBsP3sBW1md2kEM2";
    tt[2].error = SKY_ERROR;
    GoInt32 errcode;
    unsigned char output[JSON_BIG_FILE_SIZE];
    for (int i = 0; i < lenStruct; i++) {
      Config__Handle configHandle;
      App__Handle appHandle;
      errcode = SKY_cli_LoadConfig(&configHandle);
      cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
      registerHandleClose(configHandle);
      errcode = SKY_cli_NewApp(&configHandle, &appHandle);
      cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
      registerHandleClose(appHandle);

    // Redirect standard output to a pipe
      GoString Args = {tt[i].args, strlen(tt[i].args)};
      redirectStdOut();
      errcode = SKY_cli_App_Run(&appHandle, Args);
    // Get redirected standard output
      getStdOut(output, JSON_BIG_FILE_SIZE);
      cr_assert(errcode == tt[i].error, "SKY_cli_App_Run failed in %s",
        tt[i].name);
      cr_assert(errcode == tt[i].error, "Error in %s", tt[i].name);
    }
  }

  Test(cli_integration, TestDecodeRawTransaction) {
    int lenStruct = 2;

    struct testStruct {
      unsigned char *name;
      unsigned char *rawTx;
      unsigned char *errMsg;
      unsigned char *goldenFile;
      unsigned char *args;
      int error;
    };

    struct testStruct tt[lenStruct];

    tt[0].name = "success";
    tt[0].rawTx =
    "2601000000a1d3345ac47f897f24084b1c6b9bd6e03fc92887050d0748bdab5e639c1"
    "fdcd401000000a2a10f07e0e06cf6ba3e793b3186388a126591ee230b3f387617f1cc"
    "b6376a3f18e094bd3f7719aa8191c00764f323872f5192da393852bd85dab70b13409"
    "d2b01010000004d78de698a33abcfff22391c043b57a56bb0efbdc4a5b975bf8e7889"
    "668896bc0400000000bae12bbf671abeb1181fc85f1c01cdfee55deb97980c9c0a000"
    "00000543600000000000000373bb3675cbf3880bba3f3de7eb078925b8a72ad0095ba"
    "0a000000001c12000000000000008829025fe45b48f29795893a642bdaa89b2bb40e4"
    "0d2df03000000001c12000000000000008001532c3a705e7e62bb0bb80630ecc21a87"
    "ec09c0fc9b01000000001b12000000000000";
    tt[0].goldenFile = "decode-raw-transaction.golden";
    tt[0].errMsg = " ";
    tt[0].args =
    "boxfort-worker decodeRawTransaction "
    "2601000000a1d3345ac47f897f24084b1c6b9bd6e03fc92887050d0748bdab5e639c1"
    "fdcd401000000a2a10f07e0e06cf6ba3e793b3186388a126591ee230b3f387617f1cc"
    "b6376a3f18e094bd3f7719aa8191c00764f323872f5192da393852bd85dab70b13409"
    "d2b01010000004d78de698a33abcfff22391c043b57a56bb0efbdc4a5b975bf8e7889"
    "668896bc0400000000bae12bbf671abeb1181fc85f1c01cdfee55deb97980c9c0a000"
    "00000543600000000000000373bb3675cbf3880bba3f3de7eb078925b8a72ad0095ba"
    "0a000000001c12000000000000008829025fe45b48f29795893a642bdaa89b2bb40e4"
    "0d2df03000000001c12000000000000008001532c3a705e7e62bb0bb80630ecc21a87"
    "ec09c0fc9b01000000001b12000000000000";

    tt[0].error = 0;

    tt[1].name = "invalid raw transaction";
    tt[1].rawTx = "2601000000a1d";
    tt[1].goldenFile = "decode-raw-transaction.golden";
    tt[1].errMsg = "invalid raw transaction: encoding/hex: odd length hex "
    "string\nencoding/hex: odd length hex string\n";
    tt[1].args = "boxfort-worker decodeRawTransaction 2601000000a1d";
    tt[1].error = -1;

    GoInt32 errcode;
    char output[JSON_BIG_FILE_SIZE];
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
      unsigned char buffArgs[JSON_BIG_FILE_SIZE];
      GoString Args = {tt[i].args, strlen(tt[i].args)};
      redirectStdOut();
      errcode = SKY_cli_App_Run(&appHandle, Args);

    // Get redirected standard output
      getStdOut(output, JSON_BIG_FILE_SIZE);
      cr_assert(errcode == tt[i].error, "Error  AppRun in %s", tt[i].name);
    // JSON parse output

      if (tt[i].error == 0) {
        json_char *json;
        json_value *value;
        json_value *json_str;
        int result;
        json = (json_char *)output;
        value = json_parse(json, strlen(output));
        cr_assert(value != NULL, "Failed to json parse in %s", tt[i].name);
        registerJsonFree(value);

        unsigned char goldenFileURL[BUFFER_SIZE];
        strcpy(goldenFileURL, TEST_DATA_DIR);
        strcat(goldenFileURL, tt[i].goldenFile);
        json_value *expect = loadJsonFile(goldenFileURL);
        cr_assert(expect != NULL, "Failed to json parse in test :%s", tt[i].name);
        registerJsonFree(expect);

        cr_assert(compareJsonValues(expect, value) == 1,
          "compareJsonValues  in : %s", tt[i].name);
        strcpy(output, "");
      }
    }
  }

  Test(cli_integration, TestAddressGen) {
    int lenStruct = 1;

    unsigned char *name;
    unsigned char *args;
  // addressGen
    name = "addressGen";
    args = "boxfort-worker addressGen";
    GoUint32 errcode;
    char output[JSON_BIG_FILE_SIZE];

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
    GoString Args = {args, strlen(args)};
    errcode = SKY_cli_App_Run(&appHandle, Args);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    json_char *json;
    json_value *value;
    json_value *json_str;
    json_value *meta_coin;
    json_value *meta_seed;
    int result;
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse in %s", name);
    registerJsonFree(value);

    meta_coin = get_json_value(value, "meta/coin", json_string);
    cr_assert(strcmp(meta_coin->u.string.ptr, "skycoin") == 0,
      "Error in compare meta/coin == skycoin %s", name);
    char *seed;
    meta_seed = get_json_value(value, "meta/seed", json_string);
    seed = meta_seed->u.string.ptr;
    cr_assert_str_not_empty(seed);

    unsigned int lenWord = 0;
    lenWord = getCountWord(seed);

    cr_assert(lenWord == 12, "Invalid len seed in %s", name);

  // addressGen --count 2
    name = "addressGen --count 2";
    args = "boxfort-worker addressGen --count 2";

    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args1 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(&appHandle, Args1);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse in %s", name);
    registerJsonFree(value);

    meta_coin = get_json_value(value, "meta/coin", json_string);
    cr_assert(strcmp(meta_coin->u.string.ptr, "skycoin") == 0,
      "Error in compare meta/coin == skycoin %s", name);
    meta_seed = get_json_value(value, "meta/seed", json_string);
    seed = meta_seed->u.string.ptr;
    cr_assert_str_not_empty(seed);

    lenWord = getCountWord(seed);

    cr_assert(lenWord == 12, "Invalid len seed in %s", name);

  // Confirms that the wallet have 2 address
  // TODO: Not implement

    json_value *entries;
    entries = get_json_value(value, "entries", json_array);
  // Confirms the addresses are generated from the seed
    int seedLength = strlen(seed);
    GoSlice seedSlice = {seed, seedLength, seedLength};
    cli__PasswordFromBytes lSeed = {NULL, 0, 0};
    cli__PasswordFromBytes keys = {NULL, 0, 0};
    result =
    SKY_cipher_GenerateDeterministicKeyPairsSeed(seedSlice, 2, &lSeed, &keys);
    cr_assert(result == SKY_OK,
      "SKY_cipher_GenerateDeterministicKeyPairsSeed failed");
    registerMemCleanup((void *)keys.data);
    cipher__SecKey *secKey = (cipher__SecKey *)keys.data;
    json_value *entries_array;
    for (int i = 0; i < 2; i++) {
      entries_array = entries->u.array.values[i];
      char *entries_address =
      entries_array->u.object.values[0].value->u.string.ptr;
      char *entries_public =
      entries_array->u.object.values[1].value->u.string.ptr;
      char *entries_secret_key =
      entries_array->u.object.values[2].value->u.string.ptr;

      char buffer[BUFFER_SIZE];
      strnhexlower((unsigned char *)secKey, buffer, 32);
      cipher__PubKey gpubkey;
      result = SKY_cipher_PubKeyFromSecKey(secKey, &gpubkey);
      cr_assert(result == SKY_OK, "SKY_cipher_PubKeyFromSecKey  failed");
      char bufferPubKey[BUFFER_SIZE];
      strnhexlower(gpubkey, bufferPubKey, sizeof(cipher__PubKey));
      cipher__Address cAddress;
      GoString_ strAddress2;
      memset(&cAddress, 0, sizeof(cipher__Address));
      SKY_cipher_AddressFromSecKey(secKey, &cAddress);
      SKY_cipher_Address_String(&cAddress, &strAddress2);
      cr_assert(strcmp(entries_address, strAddress2.p) == 0);
      cr_assert(strcmp(entries_public, bufferPubKey) == 0);
      cr_assert(strcmp(entries_secret_key, buffer) == 0);
      secKey++;
    }

  // addressGen -c 2
    name = "addressGen -c 2";
    args = "boxfort-worker addressGen -c 2";

    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args2 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(&appHandle, Args2);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse in %s", name);
    registerJsonFree(value);

    meta_coin = get_json_value(value, "meta/coin", json_string);
    cr_assert(strcmp(meta_coin->u.string.ptr, "skycoin") == 0,
      "Error in compare meta/coin == skycoin %s", name);
    meta_seed = get_json_value(value, "meta/seed", json_string);
    seed = meta_seed->u.string.ptr;
    cr_assert_str_not_empty(seed);

    lenWord = getCountWord(seed);

    cr_assert(lenWord == 12, "Invalid len seed in %s", name);

  // Confirms that the wallet have 2 address
  // TODO: Not implement

    entries = get_json_value(value, "entries", json_array);
  // Confirms the addresses are generated from the seed
    seedLength = strlen(seed);
    GoSlice seedSlice1 = {seed, seedLength, seedLength};
    cli__PasswordFromBytes lSeed1 = {NULL, 0, 0};
    cli__PasswordFromBytes keys1 = {NULL, 0, 0};
    result = SKY_cipher_GenerateDeterministicKeyPairsSeed(seedSlice1, 2, &lSeed1,
      &keys1);
    cr_assert(result == SKY_OK,
      "SKY_cipher_GenerateDeterministicKeyPairsSeed failed");
    registerMemCleanup((void *)keys1.data);
    secKey = (cipher__SecKey *)keys1.data;
    for (int i = 0; i < 2; i++) {
      entries_array = entries->u.array.values[i];
      char *entries_address =
      entries_array->u.object.values[0].value->u.string.ptr;
      char *entries_public =
      entries_array->u.object.values[1].value->u.string.ptr;
      char *entries_secret_key =
      entries_array->u.object.values[2].value->u.string.ptr;

      char buffer[BUFFER_SIZE];
      strnhexlower((unsigned char *)secKey, buffer, 32);
      cipher__PubKey gpubkey;
      result = SKY_cipher_PubKeyFromSecKey(secKey, &gpubkey);
      cr_assert(result == SKY_OK, "SKY_cipher_PubKeyFromSecKey  failed");
      char bufferPubKey[BUFFER_SIZE];
      strnhexlower(gpubkey, bufferPubKey, sizeof(cipher__PubKey));
      cipher__Address cAddress1;
      GoString_ strAddress3;
      memset(&cAddress1, 0, sizeof(cipher__Address));
      SKY_cipher_AddressFromSecKey(secKey, &cAddress1);
      SKY_cipher_Address_String(&cAddress1, &strAddress3);
      cr_assert(strcmp(entries_address, strAddress3.p) == 0);
      cr_assert(strcmp(entries_public, bufferPubKey) == 0);
      cr_assert(strcmp(entries_secret_key, buffer) == 0);
      secKey++;
    }

  // addressGen --hide-secret -c 2
    name = "addressGen --hide-secret -c 2";
    args = "boxfort-worker addressGen --hide-secret -c 2";

    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args3 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(&appHandle, Args3);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse in %s", name);
    registerJsonFree(value);

    meta_coin = get_json_value(value, "meta/coin", json_string);
    cr_assert(strcmp(meta_coin->u.string.ptr, "skycoin") == 0,
      "Error in compare meta/coin == skycoin %s", name);

    entries = get_json_value(value, "entries", json_array);
  // Confirms the addresses are generated from the seed
    for (int i = 0; i < 2; i++) {
      entries_array = entries->u.array.values[i];

      char *entries_secret_key =
      entries_array->u.object.values[2].value->u.string.ptr;

      cr_assert(strcmp(entries_secret_key, "") == 0);
    }

  // addressGen -s -c 2
    name = "addressGen -s -c 2";
    args = "boxfort-worker addressGen -s -c 2";

    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args4 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(&appHandle, Args4);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse in %s", name);
    registerJsonFree(value);

    meta_coin = get_json_value(value, "meta/coin", json_string);
    cr_assert(strcmp(meta_coin->u.string.ptr, "skycoin") == 0,
      "Error in compare meta/coin == skycoin %s", name);

    entries = get_json_value(value, "entries", json_array);
  // Confirms the addresses are generated from the seed
    for (int i = 0; i < 2; i++) {
      entries_array = entries->u.array.values[i];

      char *entries_secret_key =
      entries_array->u.object.values[2].value->u.string.ptr;

      cr_assert(strcmp(entries_secret_key, "") == 0);
    }

  // addressGen --bitcoin -c 2
    name = "addressGen --bitcoin -c 2";
    args = "boxfort-worker addressGen --bitcoin -c 2";
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args5 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(&appHandle, Args5);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse in %s", name);
    registerJsonFree(value);

    meta_coin = get_json_value(value, "meta/coin", json_string);
    cr_assert(strcmp(meta_coin->u.string.ptr, "bitcoin") == 0,
      "Error in compare meta/coin == bitcoin %s", name);
    meta_seed = get_json_value(value, "meta/seed", json_string);
    seed = meta_seed->u.string.ptr;
    cr_assert_str_not_empty(seed);

    lenWord = getCountWord(seed);

    cr_assert(lenWord == 12, "Invalid len seed in %s", name);

  // Confirms that the wallet have 2 address
  // TODO: Not implement

    entries = get_json_value(value, "entries", json_array);
  // Confirms the addresses are generated from the seed
    seedLength = strlen(seed);
    GoSlice seedSlice2 = {seed, seedLength, seedLength};
    cli__PasswordFromBytes lSeed2 = {NULL, 0, 0};
    cli__PasswordFromBytes keys2 = {NULL, 0, 0};
    result = SKY_cipher_GenerateDeterministicKeyPairsSeed(seedSlice2, 2, &lSeed2,
      &keys2);
    cr_assert(result == SKY_OK,
      "SKY_cipher_GenerateDeterministicKeyPairsSeed failed");
    registerMemCleanup((void *)keys2.data);
    secKey = (cipher__SecKey *)keys2.data;
    for (int i = 0; i < 2; i++) {
      entries_array = entries->u.array.values[i];
      char *entries_address =
      entries_array->u.object.values[0].value->u.string.ptr;
      char *entries_public =
      entries_array->u.object.values[1].value->u.string.ptr;
      char *entries_secret_key =
      entries_array->u.object.values[2].value->u.string.ptr;

      char secKeyBuffer[BUFFER_SIZE];
      strnhexlower((unsigned char *)secKey, secKeyBuffer, 32);
      cipher__PubKey pk;
      result = SKY_cipher_PubKeyFromSecKey(secKey, &pk);
      cr_assert(result == SKY_OK, "SKY_cipher_PubKeyFromSecKey  failed");
      char bufferPubKey1[BUFFER_SIZE];
      strnhexlower(pk, bufferPubKey1, sizeof(cipher__PubKey));
      cipher__Address cAddressBitcoin;
      GoString_ strAddressBitcoin;
      GoString_ sk;
      memset(&cAddressBitcoin, 0, sizeof(cipher__Address));
      SKY_cipher_BitcoinWalletImportFormatFromSeckey(secKey, &sk);
      SKY_cipher_BitcoinAddressFromPubkey(&pk, &strAddressBitcoin);

      cr_assert(strcmp(strAddressBitcoin.p, entries_address) == 0);
      cr_assert(strcmp(bufferPubKey1, entries_public) == 0);
      cr_assert(strcmp(sk.p, entries_secret_key) == 0);
      secKey++;
    }

  // addressGen -b -c 2
    name = "addressGen -b -c 2";
    args = "boxfort-worker addressGen -b -c 2";
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args6 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(&appHandle, Args6);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse in %s", name);
    registerJsonFree(value);

    meta_coin = get_json_value(value, "meta/coin", json_string);
    cr_assert(strcmp(meta_coin->u.string.ptr, "bitcoin") == 0,
      "Error in compare meta/coin == bitcoin %s", name);
    meta_seed = get_json_value(value, "meta/seed", json_string);
    seed = meta_seed->u.string.ptr;
    cr_assert_str_not_empty(seed);

    lenWord = getCountWord(seed);

    cr_assert(lenWord == 12, "Invalid len seed in %s", name);

  // Confirms that the wallet have 2 address
  // TODO: Not implement

    entries = get_json_value(value, "entries", json_array);
  // Confirms the addresses are generated from the seed
    seedLength = strlen(seed);
    GoSlice seedSlice3 = {seed, seedLength, seedLength};
    cli__PasswordFromBytes lSeed3 = {NULL, 0, 0};
    cli__PasswordFromBytes keys3 = {NULL, 0, 0};
    result = SKY_cipher_GenerateDeterministicKeyPairsSeed(seedSlice3, 2, &lSeed3,
      &keys3);
    cr_assert(result == SKY_OK,
      "SKY_cipher_GenerateDeterministicKeyPairsSeed failed");
    registerMemCleanup((void *)keys3.data);
    secKey = (cipher__SecKey *)keys3.data;
    for (int i = 0; i < 2; i++) {
      entries_array = entries->u.array.values[i];
      char *entries_address =
      entries_array->u.object.values[0].value->u.string.ptr;
      char *entries_public =
      entries_array->u.object.values[1].value->u.string.ptr;
      char *entries_secret_key =
      entries_array->u.object.values[2].value->u.string.ptr;

      char secKeyBuffer[BUFFER_SIZE];
      strnhexlower((unsigned char *)secKey, secKeyBuffer, 32);
      cipher__PubKey pk;
      result = SKY_cipher_PubKeyFromSecKey(secKey, &pk);
      cr_assert(result == SKY_OK, "SKY_cipher_PubKeyFromSecKey  failed");
      char bufferPubKey1[BUFFER_SIZE];
      strnhexlower(pk, bufferPubKey1, sizeof(cipher__PubKey));
      cipher__Address cAddressBitcoin;
      GoString_ strAddressBitcoin;
      GoString_ sk;
      memset(&cAddressBitcoin, 0, sizeof(cipher__Address));
      SKY_cipher_BitcoinWalletImportFormatFromSeckey(secKey, &sk);
      SKY_cipher_BitcoinAddressFromPubkey(&pk, &strAddressBitcoin);

      cr_assert(strcmp(strAddressBitcoin.p, entries_address) == 0);
      cr_assert(strcmp(bufferPubKey1, entries_public) == 0);
      cr_assert(strcmp(sk.p, entries_secret_key) == 0);
      secKey++;
    }

  // addressGen --hex
    name = "addressGen --hex";
    args = "boxfort-worker addressGen --hex";
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args7 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(&appHandle, Args7);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse in %s", name);
    registerJsonFree(value);

    meta_seed = get_json_value(value, "meta/seed", json_string);
    seed = meta_seed->u.string.ptr;
    cr_assert_str_not_empty(seed);
    unsigned char strSeed[BUFFER_SIZE];
    errcode = hexnstr(seed, strSeed, BUFFER_SIZE);
    cr_assert(errcode == sizeof(cipher__SecKey));

  // addressGen --hex
    name = "addressGen -x";
    args = "boxfort-worker addressGen -x";
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args8 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(&appHandle, Args8);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse in %s", name);
    registerJsonFree(value);

    meta_seed = get_json_value(value, "meta/seed", json_string);
    seed = meta_seed->u.string.ptr;
    cr_assert_str_not_empty(seed);
    unsigned char strSeed1[BUFFER_SIZE];
    errcode = hexnstr(seed, strSeed1, BUFFER_SIZE);
    cr_assert(errcode == sizeof(cipher__SecKey));

  // addressGen --only-addr
    name = "addressGen --only-addr";
    args = "boxfort-worker addressGen --only-addr";
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args9 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(&appHandle, Args9);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    cr_assert(getCountWord(output) == 1);

  // addressGen --oa
    name = "addressGen --oa";
    args = "boxfort-worker addressGen --oa";
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args10 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(&appHandle, Args10);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    cr_assert(getCountWord(output) == 1);

  // addressGen --seed=123
    name = "addressGen --seed=123";
    args = "boxfort-worker addressGen --seed=123";

    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args12 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(&appHandle, Args12);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse in %s", name);
    registerJsonFree(value);

  // Confirms that the wallet have 2 address
  // TODO: Not implement
    seed = "123";
    entries = get_json_value(value, "entries", json_array);
  // Confirms the addresses are generated from the seed
    seedLength = strlen(seed);
    GoSlice seedSlice6 = {seed, seedLength, seedLength};
    cli__PasswordFromBytes lSeed6 = {NULL, 0, 0};
    cli__PasswordFromBytes keys6 = {NULL, 0, 0};
    result = SKY_cipher_GenerateDeterministicKeyPairsSeed(seedSlice6, 1, &lSeed6,
      &keys6);
    cr_assert(result == SKY_OK,
      "SKY_cipher_GenerateDeterministicKeyPairsSeed failed");
    registerMemCleanup((void *)keys6.data);
    secKey = (cipher__SecKey *)keys6.data;
    entries_array = entries->u.array.values[0];
    char *entries_address = entries_array->u.object.values[0].value->u.string.ptr;
    char *entries_public = entries_array->u.object.values[1].value->u.string.ptr;
    char *entries_secret_key =
    entries_array->u.object.values[2].value->u.string.ptr;

    char buffer[BUFFER_SIZE];
    strnhexlower((unsigned char *)secKey, buffer, 32);
    cipher__PubKey gpubkey6;
    result = SKY_cipher_PubKeyFromSecKey(secKey, &gpubkey6);
    cr_assert(result == SKY_OK, "SKY_cipher_PubKeyFromSecKey  failed");
    char bufferPubKey[BUFFER_SIZE];
    strnhexlower(gpubkey6, bufferPubKey, sizeof(cipher__PubKey));
    cipher__Address cAddress6;
    GoString_ strAddress9;
    memset(&cAddress6, 0, sizeof(cipher__Address));
    SKY_cipher_AddressFromSecKey(secKey, &cAddress6);
    SKY_cipher_Address_String(&cAddress6, &strAddress9);
    cr_assert(strcmp(entries_address, strAddress9.p) == 0);
    cr_assert(strcmp(entries_public, bufferPubKey) == 0);
    cr_assert(strcmp(entries_secret_key, buffer) == 0);
  }

  Test(api_cli_integracion, TestStableListWallets) {
    createTempWalletDir(false);

    const char *args;
    args = "boxfort-worker listWallets";

    GoUint32 errcode;
    char output[JSON_BIG_FILE_SIZE];

    Config__Handle configHandle;
    App__Handle appHandle;
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    GoString Args = {args, strlen(args)};
    redirectStdOut();
    errcode = SKY_cli_App_Run(&appHandle, Args);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);

    char buffName[BUFFER_SIZE];
    char buffDir[BUFFER_SIZE];
    GoString WalletName = {"WALLET_NAME", 11};
    GoString WalletPath = {"WALLET_DIR", 10};
    GoString pathname = {buffName, 0};
    GoString Dir = {buffDir, 0};

  // JSON parse output
    json_char *json;
    json_value *value;
    json_value *json_str;
    int result;
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse \n%s\n", output);
    registerJsonFree(value);

    unsigned char goldenFileURL[BUFFER_SIZE];
    strcpy(goldenFileURL, TEST_DATA_DIR);
    strcat(goldenFileURL, "list-wallets.golden");
    json_value *expect = loadJsonFile(goldenFileURL);
    cr_assert(expect != NULL, "Failed to json parse in test :%s", goldenFileURL);
    registerJsonFree(expect);

    json_value *wallets_expect;
    wallets_expect = get_json_value(expect, "wallets", json_array);

    wallets_expect = wallets_expect->u.array.values[0];
    char *wallets_expect_name =
    wallets_expect->u.object.values[0].value->u.string.ptr;
    char *wallets_expect_label =
    wallets_expect->u.object.values[1].value->u.string.ptr;
    int64_t wallets_expect_address_num =
    wallets_expect->u.object.values[2].value->u.integer;

    json_value *wallets_value;
    wallets_value = get_json_value(value, "wallets", json_array);

    wallets_value = wallets_value->u.array.values[0];
    char *wallets_value_name =
    wallets_value->u.object.values[0].value->u.string.ptr;
    char *wallets_value_label =
    wallets_value->u.object.values[1].value->u.string.ptr;
    int64_t wallets_value_address_num =
    wallets_value->u.object.values[2].value->u.integer;
    cr_assert(strcmp(wallets_value_name, wallets_expect_name) == 0);
    cr_assert(strcmp(wallets_value_label, wallets_value_label) == 0);
    cr_assert(wallets_value_address_num == wallets_expect_address_num);
  }

  Test(api_cli_integracion, TestStableListAddress) {
    createTempWalletDir(false);

    const char *args;
    args = "boxfort-worker listAddresses";

    GoUint32 errcode;
    char output[JSON_BIG_FILE_SIZE];

    Config__Handle configHandle;
    App__Handle appHandle;
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    GoString Args = {args, strlen(args)};
    redirectStdOut();
    errcode = SKY_cli_App_Run(&appHandle, Args);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);

    char buffName[BUFFER_SIZE];
    char buffDir[BUFFER_SIZE];
    GoString WalletName = {"WALLET_NAME", 11};
    GoString WalletPath = {"WALLET_DIR", 10};
    GoString pathname = {buffName, 0};
    GoString Dir = {buffDir, 0};

  // JSON parse output
    json_char *json;
    json_value *value;
    json_value *json_str;
    int result;
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse \n%s\n", output);
    registerJsonFree(value);

    unsigned char goldenFileURL[BUFFER_SIZE];
    strcpy(goldenFileURL, TEST_DATA_DIR);
    strcat(goldenFileURL, "list-addresses.golden");
    json_value *expect = loadJsonFile(goldenFileURL);
    cr_assert(expect != NULL, "Failed to json parse in test :%s", goldenFileURL);
    registerJsonFree(expect);

    json_value *addresses_expect;
    addresses_expect = get_json_value(expect, "addresses", json_array);

    addresses_expect = addresses_expect->u.array.values[0];

    json_value *addresses_value;
    addresses_value = get_json_value(value, "addresses", json_array);

    addresses_value = addresses_value->u.array.values[0];
    cr_assert(compareJsonValues(addresses_value, addresses_expect) == 1);
  }

  Test(api_cli_integracion, TestStableAddressBalance) {
    const char *args;
    args = "boxfort-worker addressBalance 2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt";

    GoUint32 errcode;
    char output[JSON_BIG_FILE_SIZE];

    Config__Handle configHandle;
    App__Handle appHandle;
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    GoString Args = {args, strlen(args)};
    redirectStdOut();
    errcode = SKY_cli_App_Run(&appHandle, Args);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d", errcode);
    printf("El OUTPUT %s\n", output);
  // JSON parse output
    json_char *json;
    json_value *value;
    json_value *json_str;
    int result;
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse in output \n%s\n", output);
    registerJsonFree(value);

    unsigned char goldenFileURL[BUFFER_SIZE];
    strcpy(goldenFileURL, TEST_DATA_DIR);
    strcat(goldenFileURL, "address-balance.golden");
    json_value *expect = loadJsonFile(goldenFileURL);
    cr_assert(expect != NULL, "Failed to json parse in test :%s", goldenFileURL);
    registerJsonFree(expect);

    cr_assert(compareJsonValues(value, expect) == 1);
  }

  Test(api_cli_integracion, TestStableWalletBalance) {
    createTempWalletDir(false);
    const char *args;
    args = "boxfort-worker walletBalance";

    GoUint32 errcode;
    char output[JSON_BIG_FILE_SIZE];

    Config__Handle configHandle;
    App__Handle appHandle;
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    GoString Args = {args, strlen(args)};
    redirectStdOut();
    errcode = SKY_cli_App_Run(&appHandle, Args);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d", errcode);
    printf("El OUTPUT %s\n", output);
  // JSON parse output
    json_char *json;
    json_value *value;
    json_value *json_str;
    int result;
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse in output \n%s\n", output);
    registerJsonFree(value);

    unsigned char goldenFileURL[BUFFER_SIZE];
    strcpy(goldenFileURL, TEST_DATA_DIR);
    strcat(goldenFileURL, "wallet-balance.golden");
    json_value *expect = loadJsonFile(goldenFileURL);
    cr_assert(expect != NULL, "Failed to json parse in test :%s", goldenFileURL);
    registerJsonFree(expect);

    cr_assert(compareJsonValues(value, expect) == 1);
  }

  Test(api_cli_integracion, TestStableWalletOutputs) {
    createTempWalletDir(false);
    const char *args;
    args = "boxfort-worker walletOutputs";

    GoUint32 errcode;
    char output[JSON_BIG_FILE_SIZE];

    Config__Handle configHandle;
    App__Handle appHandle;
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(&configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    GoString Args = {args, strlen(args)};
    redirectStdOut();
    errcode = SKY_cli_App_Run(&appHandle, Args);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d", errcode);
    printf("El OUTPUT %s\n", output);
  // JSON parse output
    json_char *json;
    json_value *value;
    json_value *json_str;
    int result;
    json = (json_char *)output;
    value = json_parse(json, strlen(output));
    cr_assert(value != NULL, "Failed to json parse in output \n%s\n", output);
    registerJsonFree(value);

    unsigned char goldenFileURL[BUFFER_SIZE];
    strcpy(goldenFileURL, TEST_DATA_DIR);
    strcat(goldenFileURL, "wallet-outputs.golden");
    json_value *expect = loadJsonFile(goldenFileURL);
    cr_assert(expect != NULL, "Failed to json parse in test :%s", goldenFileURL);
    registerJsonFree(expect);

    cr_assert(compareJsonValues(value, expect) == 1);
  }

