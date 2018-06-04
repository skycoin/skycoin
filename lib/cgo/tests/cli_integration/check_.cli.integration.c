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
#define JSON_BIG_FILE_SIZE 102400
#define TEST_DATA_DIR "src/cli/integration/testdata/"
#define SKYCOIN_NODE_HOST "http://127.0.0.1:6420"
#define stableWalletName "integration-test.wlt"
#define stableEncryptWalletName "integration-test-encrypted.wlt"
#define NODE_ADDRESS "SKYCOIN_NODE_HOST"
#define NODE_ADDRESS_DEFAULT "http://127.0.0.1:46420"
#define STABLE 1
#define NORMAL_TESTS
#define DECRYPTION_TESTS
#define DECRYPT_WALLET_TEST
#define LENARRAY(x) (sizeof(x) / sizeof((*x)))

void useClient() {
  GoString default_webprc = {"RPC_ADDR", 8};
  GoString default_webprc_value = {NODE_ADDRESS_DEFAULT,
   strlen(NODE_ADDRESS_DEFAULT)};
   SKY_cli_Setenv(default_webprc, default_webprc_value);
 }

 int useCSRF() {
  GoUint32 errcode;

  GoString strCSRFVar = {"USE_CSRF", 8};
  char buffercrsf[BUFFER_SIZE];
  GoString_ crsf = {buffercrsf, 0};
  errcode = SKY_cli_Getenv(strCSRFVar, &crsf);
  cr_assert(errcode == SKY_OK, "SKY_cli_Getenv failed");
  int length = strlen(crsf.p);
  int result = 0;
  if (length == 1) {
    result = crsf.p[0] == '1' || crsf.p[0] == 't' || crsf.p[0] == 'T';
  } else {
    result = strcmp(crsf.p, "true") == 0 || strcmp(crsf.p, "True") == 0 ||
    strcmp(crsf.p, "TRUE") == 0;
  }
  free((void *)crsf.p);
  return result;
}
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

int getCountWordWithSeparator(const char *str, const char *sep) {
  int len = 0;
  do {
    str = strpbrk(str, sep); // find separator
    if (str)
      str += strspn(str, sep); // skip separator
    ++len;                     // increment word count
  } while (str && *str);

  return len;
}

int getCountStringInString(const char *source, const char *str) {
  // TODO:Not implement
}

char *getNodeAddress_Cli() {
  if (STABLE) {
    return NODE_ADDRESS_DEFAULT;
  } else {
    GoString_ nodeAddress;
    memset(&nodeAddress, 0, sizeof(GoString_));
    GoString nodeEnvName = {NODE_ADDRESS, strlen(NODE_ADDRESS)};
    int result = SKY_cli_Getenv(nodeEnvName, &nodeAddress);
    cr_assert(result == SKY_OK, "Couldn\'t get node address from enviroment");
    registerMemCleanup((void *)nodeAddress.p);
    if (strcmp(nodeAddress.p, "") == 0) {
      return NODE_ADDRESS_DEFAULT;
    }
    return (char *)nodeAddress.p;
  }
}
json_value *loadGoldenFile_Cli(const char *file) {
  char path[STRING_SIZE];
  if (strlen(TEST_DATA_DIR) + strlen(file) < STRING_SIZE) {
    strcpy(path, TEST_DATA_DIR);
    strcat(path, file);
    return loadJsonFile(path);
  }
  return NULL;
}

int cleanPath(char *path) {
  char cmd[BUFFER_SIZE];
  strcpy(cmd, "rm -r ");
  strcat(cmd, path);
  strcat(cmd, "/*.*");
  int i = system(cmd);
  return i;
}
// createTempWalletDir creates a temporary wallet dir,
// sets the WALLET_DIR environment variable.
// Returns wallet dir path and callback function to clean up the dir.
void createTempWalletDir(bool encrypt) {

  char *pNodeAddress = getNodeAddress_Cli();
  GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
  Client__Handle clientHandle;

  int result = SKY_api_NewClient(nodeAddress, &clientHandle);
  cr_assert(result == SKY_OK, "Couldn\'t create client");
  registerHandleClose(clientHandle);

  char bufferdir[BUFFER_SIZE];
  GoString_ dir_temp_value = {bufferdir, 0};

  int errcode =
  SKY_api_Handle_Client_GetWalletDir(clientHandle, &dir_temp_value);
  cr_assert(errcode == SKY_OK);
  // Copy the testdata/$stableWalletName to the temporary dir.
  char walletPath[JSON_BIG_FILE_SIZE];
  if (encrypt) {
    strcpy(walletPath, stableEncryptWalletName);
  } else {
    strcpy(walletPath, stableWalletName);
  }
  unsigned char pathnameURL[BUFFER_SIZE];
  strcpy(pathnameURL, dir_temp_value.p);
  strcat(pathnameURL, "/");
  strcat(pathnameURL, walletPath);

  cleanPath((char *)dir_temp_value.p);

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
    GoString Dir = {dir_temp_value.p, dir_temp_value.n};
    SKY_cli_Setenv(WalletDir, Dir);
    GoString WalletPath = {"WALLET_NAME", 11};
    GoString pathname = {walletPath, strlen(walletPath)};
    SKY_cli_Setenv(WalletPath, pathname);
  }
  strcpy(walletPath, "");
};

// createTempWalletDir creates a temporary wallet dir,
// sets the WALLET_DIR environment variable.
// Returns wallet dir path and callback function to clean up the dir.
void createWalletDir() {
  char *pNodeAddress = getNodeAddress_Cli();
  GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
  Client__Handle clientHandle;

  int result = SKY_api_NewClient(nodeAddress, &clientHandle);
  cr_assert(result == SKY_OK, "Couldn\'t create client");
  registerHandleClose(clientHandle);

  char bufferdir[BUFFER_SIZE];
  GoString_ dir_temp_value = {bufferdir, 0};

  int errcode =
  SKY_api_Handle_Client_GetWalletDir(clientHandle, &dir_temp_value);
  cr_assert(errcode == SKY_OK);
  cleanPath((char *)dir_temp_value.p);

  GoString WalletDir = {"WALLET_DIR", 10};
  GoString Dir = {dir_temp_value.p, dir_temp_value.n};
  SKY_cli_Setenv(WalletDir, Dir);
};

TestSuite(cli_integration, .init = setup, .fini = teardown);

Test(cli_integration, TestGenerateAddresses) {
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
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

    // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args = {tt[i].args, strlen(tt[i].args)};
    errcode = SKY_cli_App_Run(appHandle, Args);
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

  Test(cli_integration, TestVerifyAddress) {
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
      errcode = SKY_cli_NewApp(configHandle, &appHandle);
      cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
      registerHandleClose(appHandle);

    // Redirect standard output to a pipe
      GoString Args = {tt[i].args, strlen(tt[i].args)};
      redirectStdOut();
      errcode = SKY_cli_App_Run(appHandle, Args);
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
      errcode = SKY_cli_NewApp(configHandle, &appHandle);
      cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
      registerHandleClose(appHandle);

    // Redirect standard output to a pipe
      unsigned char buffArgs[JSON_BIG_FILE_SIZE];
      GoString Args = {tt[i].args, strlen(tt[i].args)};
      redirectStdOut();
      errcode = SKY_cli_App_Run(appHandle, Args);

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
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args = {args, strlen(args)};
    errcode = SKY_cli_App_Run(appHandle, Args);
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
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args1 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(appHandle, Args1);
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
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args2 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(appHandle, Args2);
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
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args3 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(appHandle, Args3);
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
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args4 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(appHandle, Args4);
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
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args5 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(appHandle, Args5);
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
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args6 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(appHandle, Args6);
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
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args7 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(appHandle, Args7);
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
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args8 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(appHandle, Args8);
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
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args9 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(appHandle, Args9);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    cr_assert(getCountWord(output) == 1);

  // addressGen --oa
    name = "addressGen --oa";
    args = "boxfort-worker addressGen --oa";
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args10 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(appHandle, Args10);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    cr_assert(getCountWord(output) == 1);

  // addressGen --seed=123
    name = "addressGen --seed=123";
    args = "boxfort-worker addressGen --seed=123";

    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    redirectStdOut();
    GoString Args12 = {args, strlen(args)};
    errcode = SKY_cli_App_Run(appHandle, Args12);
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
    useClient();
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
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    GoString Args = {args, strlen(args)};
    redirectStdOut();
    errcode = SKY_cli_App_Run(appHandle, Args);
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

    json_value *wallets_value;
    wallets_value = get_json_value(value, "wallets", json_array);

    wallets_value = wallets_value->u.array.values[0];
    cr_assert(compareJsonValues(wallets_expect, wallets_value) == 1);
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
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    GoString Args = {args, strlen(args)};
    redirectStdOut();
    errcode = SKY_cli_App_Run(appHandle, Args);
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
    useClient();
    const char *args;
    args = "boxfort addressBalance 2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt";

    GoUint32 errcode;
    char output[JSON_BIG_FILE_SIZE];

    Config__Handle configHandle;
    App__Handle appHandle;
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);
  // Redirect standard output to a pipe
    GoString Args = {args, strlen(args)};
    redirectStdOut();
    errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d", errcode);
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
    useClient();
    const char *args;
    args = "boxfort-worker walletBalance";

    GoUint32 errcode;
    char output[JSON_BIG_FILE_SIZE];

    Config__Handle configHandle;
    App__Handle appHandle;
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    GoString Args = {args, strlen(args)};
    redirectStdOut();
    errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d", errcode);
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
    useClient();
    const char *args;
    args = "boxfort-worker walletOutputs";

    GoUint32 errcode;
    char output[JSON_BIG_FILE_SIZE];

    Config__Handle configHandle;
    App__Handle appHandle;
    errcode = SKY_cli_LoadConfig(&configHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
    registerHandleClose(configHandle);
    errcode = SKY_cli_NewApp(configHandle, &appHandle);
    cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
    registerHandleClose(appHandle);

  // Redirect standard output to a pipe
    GoString Args = {args, strlen(args)};
    redirectStdOut();
    errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
    getStdOut(output, JSON_BIG_FILE_SIZE);
    cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d", errcode);
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

  Test(cli_integration, TestStableAddressOutputs) {

    typedef struct {
      char *name;
      char *args;
      char *goldenFile;
      int isUsageErr;
    } testStruct;

    testStruct tt[] = {
      {"addressOutputs one address",
      "boxfort-worker addressOutputs 2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt",
      "address-outputs.golden", SKY_OK},
      {"addressOutputs two address",
      "boxfort-worker addressOutputs 2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt "
      "ejJjiCwp86ykmFr5iTJ8LxQXJ2wJPTYmkm",
      "two-addresses-outputs.golden", SKY_OK},
      {"addressOutputs two address one invalid",
      "boxfort-worker addressOutputs 2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt "
      "badaddress",
      "", SKY_ERROR}};

      useClient();
      GoUint32 errcode;
      char output[JSON_BIG_FILE_SIZE];
      Config__Handle configHandle;
      App__Handle appHandle;

      errcode = SKY_cli_LoadConfig(&configHandle);
      cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
      registerHandleClose(configHandle);
      errcode = SKY_cli_NewApp(configHandle, &appHandle);
      cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
      registerHandleClose(appHandle);
      for (int i = 0; i < LENARRAY(tt); i++) {

    // Redirect standard output to a pipe
        GoString Args = {tt[i].args, strlen(tt[i].args)};
        redirectStdOut();
        errcode = SKY_cli_App_Run(appHandle, Args);
    // Get redirected standard output
        getStdOut(output, JSON_BIG_FILE_SIZE);
        cr_assert(errcode == tt[i].isUsageErr,
          "SKY_cli_App_Run failed in return %d", errcode);

        if (tt[i].isUsageErr == SKY_OK) {

      // JSON parse output
          json_char *json;
          json_value *value;
          json_value *json_str;
          int result;
          json = (json_char *)output;
          value = json_parse(json, strlen(output));
          cr_assert(value != NULL, "Failed to json parse in %s\n", tt[i].name);
          registerJsonFree(value);

          unsigned char goldenFileURL[BUFFER_SIZE];
          strcpy(goldenFileURL, TEST_DATA_DIR);
          strcat(goldenFileURL, tt[i].goldenFile);
          json_value *expect = loadJsonFile(goldenFileURL);
          cr_assert(expect != NULL, "Failed to json parse in test :%s",
            tt[i].goldenFile);
          registerJsonFree(expect);

          cr_assert(compareJsonValues(value, expect) == 1);
        }
      }
    }

    Test(cli_integration, TestStableShowConfig) {
      char output[BUFFER_SIZE];

      Config__Handle configHandle;
      App__Handle appHandle;
      const char *str = "boxfort-worker showConfig";
      GoString showConfigCommand = {str, strlen(str)};
      GoUint32 errcode;

      errcode = SKY_cli_LoadConfig(&configHandle);
      cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
      registerHandleClose(configHandle);
      errcode = SKY_cli_NewApp(configHandle, &appHandle);
      cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
      registerHandleClose(appHandle);

  // Redirect standard output to a pipe
      redirectStdOut();
      errcode = SKY_cli_App_Run(appHandle, showConfigCommand);
  // Get redirected standard output
      getStdOut(output, BUFFER_SIZE);
      cr_assert(errcode == SKY_OK, "SKY_cli_App_Run failed");

  // JSON parse output
      json_char *json;
      json_value *value;
      json_value *json_str;
      int result;
      json = (json_char *)output;
      value = json_parse(json, strlen(output));
      cr_assert(value != NULL, "Failed to json parse");
      registerJsonFree(value);

      json_value *wallet_dir = json_get_string(value, "wallet_directory");
      cr_assert(wallet_dir != NULL, "Failed to get json value");
      json_value *data_dir = json_get_string(value, "data_directory");
      cr_assert(data_dir != NULL, "Failed to get json value");
      json_value *wallet_name = json_get_string(value, "wallet_name");
      cr_assert(wallet_name != NULL, "Failed to get json value");
      json_value *coin_name = json_get_string(value, "coin");
      cr_assert(coin_name != NULL, "Failed to get json value");
      json_value *rpc_address = json_get_string(value, "rpc_address");
      cr_assert(rpc_address != NULL, "Failed to get json value");

      result = string_has_suffix(wallet_dir->u.string.ptr, ".skycoin/wallets");
      cr_assert(result == 1, "Wallet dir must end in .skycoin/wallets");
      result = string_has_suffix(data_dir->u.string.ptr, ".skycoin");
      cr_assert(result == 1, "Data dir must end in .skycoin");
      result = string_has_prefix(wallet_dir->u.string.ptr, data_dir->u.string.ptr);
      cr_assert(result == 1, "Data dir must be prefix of wallet dir");

      json_set_string(wallet_dir, "IGNORED/.skycoin/wallets");
      json_set_string(data_dir, "IGNORED/.skycoin");
  // Ignore the rpc address
      json_set_string(rpc_address, "http://127.0.0.1:46420");

      const char *golden_file = "show-config.golden";
      if (useCSRF()) {
        golden_file = "show-config-use-csrf.golden";
      }
      json_value *json_golden = loadGoldenFile_Cli(golden_file);
      cr_assert(json_golden != NULL, "loadGoldenFile_Cli failed");
      registerJsonFree(json_golden);
      int equal = compareJsonValues(value, json_golden);
      cr_assert(equal, "Output from command different than expected");
    }

    Test(api_cli_integracion, TestStableStatus) {
      createTempWalletDir(false);
      useClient();
      const char *args;
      args = "boxfort-worker status";

      GoUint32 errcode;
      char output[JSON_BIG_FILE_SIZE];

      Config__Handle configHandle;
      App__Handle appHandle;
      errcode = SKY_cli_LoadConfig(&configHandle);
      cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
      registerHandleClose(configHandle);
      errcode = SKY_cli_NewApp(configHandle, &appHandle);
      cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
      registerHandleClose(appHandle);

  // Redirect standard output to a pipe
      GoString Args = {args, strlen(args)};
      redirectStdOut();
      errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
      getStdOut(output, JSON_BIG_FILE_SIZE);
      cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d", errcode);

  // JSON parse output
      json_char *json;
      json_value *value;
      json_value *json_str;
      int result;
      json = (json_char *)output;
      value = json_parse(json, strlen(output));
      cr_assert(value != NULL, "Failed to json parse in output \n%s\n", output);
      registerJsonFree(value);

      json_value *TimeSinceLastBlock =
      json_get_string(value, "time_since_last_block");
      cr_assert(TimeSinceLastBlock != NULL,
        "Failed to get json value TimeSinceLastBlock");

      json_set_string(TimeSinceLastBlock, "");

      unsigned char goldenFileURL[BUFFER_SIZE];
      strcpy(goldenFileURL, TEST_DATA_DIR);
      strcat(goldenFileURL, "status.golden");
      json_value *expect = loadJsonFile(goldenFileURL);
      cr_assert(expect != NULL, "Failed to json parse in test :%s", goldenFileURL);
      registerJsonFree(expect);

      cr_assert(compareJsonValues(value, expect) == 1, "NOt eq json");
    }

    Test(cli_integration, TestStableTransaction) {

      typedef struct {
        char *name;
        char *args;
        char *goldenFile;
        int isUsageErr;
      } testStruct;

      testStruct tt[] = {
        {"invalid txid", "boxfort-worker transaction abcd", "", SKY_ERROR},
        {"not exist",
        "boxfort-worker transaction "
        "701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947",
        "", SKY_ERROR},
        {"empty txid", "boxfort-worker transaction ", "", SKY_ERROR},
        {"genesis transaction",
        "boxfort-worker transaction "
        "d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add ",
        "genesis-transaction-cli.golden", SKY_OK}};

        useClient();
        GoUint32 errcode;
        char output[JSON_BIG_FILE_SIZE];
        Config__Handle configHandle;
        App__Handle appHandle;

        errcode = SKY_cli_LoadConfig(&configHandle);
        cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
        registerHandleClose(configHandle);
        errcode = SKY_cli_NewApp(configHandle, &appHandle);
        cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
        registerHandleClose(appHandle);
        for (int i = 0; i < LENARRAY(tt); i++) {

    // Redirect standard output to a pipe
          GoString Args = {tt[i].args, strlen(tt[i].args)};
          redirectStdOut();
          errcode = SKY_cli_App_Run(appHandle, Args);
    // Get redirected standard output
          getStdOut(output, JSON_BIG_FILE_SIZE);
          cr_assert(errcode == tt[i].isUsageErr,
            "SKY_cli_App_Run failed in return %d", errcode);

          if (tt[i].isUsageErr == SKY_OK) {

      // JSON parse output
            json_char *json;
            json_value *value;
            json_value *json_str;
            int result;
            json = (json_char *)output;
            value = json_parse(json, strlen(output));
            cr_assert(value != NULL, "Failed to json parse in %s\n", tt[i].name);
            registerJsonFree(value);

            unsigned char goldenFileURL[BUFFER_SIZE];
            strcpy(goldenFileURL, TEST_DATA_DIR);
            strcat(goldenFileURL, tt[i].goldenFile);
            json_value *expect = loadJsonFile(goldenFileURL);
            cr_assert(expect != NULL, "Failed to json parse in test :%s",
              tt[i].goldenFile);
            registerJsonFree(expect);

            cr_assert(compareJsonValues(value, expect) == 1);
          }
        }
      }

      Test(api_cli_integracion, TestStableBlocks) {
        createTempWalletDir(false);
        useClient();
        const char *args;
        args = "boxfort-worker blocks 180 181";

        GoUint32 errcode;
        char output[JSON_BIG_FILE_SIZE];

        Config__Handle configHandle;
        App__Handle appHandle;
        errcode = SKY_cli_LoadConfig(&configHandle);
        cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
        registerHandleClose(configHandle);
        errcode = SKY_cli_NewApp(configHandle, &appHandle);
        cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
        registerHandleClose(appHandle);

  // Redirect standard output to a pipe
        GoString Args = {args, strlen(args)};
        redirectStdOut();
        errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
        getStdOut(output, JSON_BIG_FILE_SIZE);
        cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d", errcode);

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
        strcat(goldenFileURL, "blocks180.golden");
        json_value *expect = loadJsonFile(goldenFileURL);
        cr_assert(expect != NULL, "Failed to json parse in test :%s", goldenFileURL);
        registerJsonFree(expect);

        cr_assert(compareJsonValues(value, expect) == 1, "NOt eq json");
      }

      Test(cli_integration, testKnownBlocks) {

        typedef struct {
          char *name;
          char *args;
          char *goldenFile;
        } testStruct;

        testStruct tt[] = {
          {"blocks 0", "boxfort-worker blocks 0", "block0.golden"},
          {"blocks 0 5", "boxfort-worker blocks 0 5 ", "blocks0-5.golden"}};

          useClient();
          GoUint32 errcode;
          char output[JSON_BIG_FILE_SIZE];
          Config__Handle configHandle;
          App__Handle appHandle;

          errcode = SKY_cli_LoadConfig(&configHandle);
          cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
          registerHandleClose(configHandle);
          errcode = SKY_cli_NewApp(configHandle, &appHandle);
          cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
          registerHandleClose(appHandle);

          for (int i = 0; i < LENARRAY(tt); i++) {

    // Redirect standard output to a pipe
            GoString Args = {tt[i].args, strlen(tt[i].args)};
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
    // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            cr_assert(errcode == SKY_OK, "SKY_cli_App_Run failed in return %d",
              errcode);

    // JSON parse output
            json_char *json;
            json_value *value;
            json_value *json_str;
            int result;
            json = (json_char *)output;
            value = json_parse(json, strlen(output));
            cr_assert(value != NULL, "Failed to json parse in %s\n", tt[i].name);
            registerJsonFree(value);

            unsigned char goldenFileURL[BUFFER_SIZE];
            strcpy(goldenFileURL, TEST_DATA_DIR);
            strcat(goldenFileURL, tt[i].goldenFile);
            json_value *expect = loadJsonFile(goldenFileURL);
            cr_assert(expect != NULL, "Failed to json parse in test :%s",
              tt[i].goldenFile);
            registerJsonFree(expect);

            cr_assert(compareJsonValues(value, expect) == 1);

            strcpy(output, "");
          }
        }

        Test(cli_integration, TestStableLastBlocks) {

          typedef struct {
            char *name;
            char *args;
            char *goldenFile;
          } testStruct;

          testStruct tt[] = {
            {"lastBlocks 0", "boxfort-worker lastBlocks 0", "last-blocks0.golden"},
            {"lastBlocks 1", "boxfort-worker lastBlocks 1", "last-blocks1.golden"},
            {"lastBlocks 2", "boxfort-worker lastBlocks 2", "last-blocks2.golden"}};

            useClient();
            GoUint32 errcode;
            char output[JSON_BIG_FILE_SIZE];
            Config__Handle configHandle;
            App__Handle appHandle;

            errcode = SKY_cli_LoadConfig(&configHandle);
            cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
            registerHandleClose(configHandle);
            errcode = SKY_cli_NewApp(configHandle, &appHandle);
            cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
            registerHandleClose(appHandle);

            for (int i = 0; i < LENARRAY(tt); i++) {

    // Redirect standard output to a pipe
              GoString Args = {tt[i].args, strlen(tt[i].args)};
              redirectStdOut();
              errcode = SKY_cli_App_Run(appHandle, Args);
    // Get redirected standard output
              getStdOut(output, JSON_BIG_FILE_SIZE);
              cr_assert(errcode == SKY_OK, "SKY_cli_App_Run failed in return %d",
                errcode);

    // JSON parse output
              json_char *json;
              json_value *value;
              json_value *json_str;
              int result;
              json = (json_char *)output;
              value = json_parse(json, strlen(output));
              cr_assert(value != NULL, "Failed to json parse in %s\n", tt[i].name);
              registerJsonFree(value);

              unsigned char goldenFileURL[BUFFER_SIZE];
              strcpy(goldenFileURL, TEST_DATA_DIR);
              strcat(goldenFileURL, tt[i].goldenFile);
              json_value *expect = loadJsonFile(goldenFileURL);
              cr_assert(expect != NULL, "Failed to json parse in test :%s",
                tt[i].goldenFile);
              registerJsonFree(expect);

              cr_assert(compareJsonValues(value, expect) == 1);

              strcpy(output, "");
            }
          }

          Test(api_cli_integracion, TestStableWalletDir) {
            createTempWalletDir(false);
            useClient();
            const char *args;
            args = "boxfort-worker walletDir";

            GoUint32 errcode;
            char output[JSON_BIG_FILE_SIZE];

            Config__Handle configHandle;
            App__Handle appHandle;
            errcode = SKY_cli_LoadConfig(&configHandle);
            cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
            registerHandleClose(configHandle);
            errcode = SKY_cli_NewApp(configHandle, &appHandle);
            cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
            registerHandleClose(appHandle);

  // Redirect standard output to a pipe
            GoString Args = {args, strlen(args)};
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d", errcode);

            GoString wallet_dir = {"WALLET_DIR", 10};
            char bufferdir[BUFFER_SIZE];
            GoString_ wallet_dir_value = {bufferdir, 0};
            errcode = SKY_cli_Getenv(wallet_dir, &wallet_dir_value);
            cr_assert(eq(u8[wallet_dir_value.n], ((char *)wallet_dir_value.p), output));
          }

          Test(api_cli_integracion, TestStableWalletHistory) {
            createTempWalletDir(false);
            useClient();
            const char *args;
            args = "boxfort-worker walletHistory";

            GoUint32 errcode;
            char output[JSON_BIG_FILE_SIZE];

            Config__Handle configHandle;
            App__Handle appHandle;
            errcode = SKY_cli_LoadConfig(&configHandle);
            cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
            registerHandleClose(configHandle);
            errcode = SKY_cli_NewApp(configHandle, &appHandle);
            cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
            registerHandleClose(appHandle);

  // Redirect standard output to a pipe
            GoString Args = {args, strlen(args)};
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d", errcode);

  // JSON parse output
            json_char *json;
            json_value *value;
            json_value *json_str;
            int result;
            json = (json_char *)output;
            value = json_parse(json, strlen(output));
            cr_assert(value != NULL, "Failed to json parse ");
            registerJsonFree(value);

            unsigned char goldenFileURL[BUFFER_SIZE];
            strcpy(goldenFileURL, TEST_DATA_DIR);
            strcat(goldenFileURL, "wallet-history.golden");
            json_value *expect = loadJsonFile(goldenFileURL);
            cr_assert(expect != NULL, "Failed to json parse in test ");
            registerJsonFree(expect);

            cr_assert(compareJsonValues(value, expect) == 1);
          }

          Test(cli_integration, TestStableCheckDB) {

            typedef struct {
              char *name;
              char *args;
              char *result;
              int isUsageErr;
            } testStruct;

            testStruct tt[] = {
              {"no signature",
              "boxfort-worker checkdb src/visor/testdata/data.db.nosig", "",
              SKY_ERROR},
              {"invalid database",
              "boxfort-worker checkdb src/visor/testdata/data.db.garbage", "",
              SKY_ERROR},
              {"valid database",
              "boxfort-worker checkdb src/api/integration/testdata/blockchain-180.db",
              "check db success\n", SKY_OK},
            };

            useClient();
            GoUint32 errcode;
            char output[JSON_BIG_FILE_SIZE];
            Config__Handle configHandle;
            App__Handle appHandle;

            errcode = SKY_cli_LoadConfig(&configHandle);
            cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
            registerHandleClose(configHandle);
            errcode = SKY_cli_NewApp(configHandle, &appHandle);
            cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
            registerHandleClose(appHandle);
            for (int i = 0; i < LENARRAY(tt); i++) {

    // Redirect standard output to a pipe
              GoString Args = {tt[i].args, strlen(tt[i].args)};
              redirectStdOut();
              errcode = SKY_cli_App_Run(appHandle, Args);
    // Get redirected standard output
              getStdOut(output, JSON_BIG_FILE_SIZE);
              cr_assert(errcode == tt[i].isUsageErr,
                "SKY_cli_App_Run failed in return %d", errcode);

              if (tt[i].isUsageErr == SKY_OK) {

                cr_assert(eq(u8[strlen(tt[i].result)], tt[i].result, output));
              }
            }
          }

          Test(api_cli_integracion, TestVersion) {
            createTempWalletDir(false);
            useClient();
            const char *args;
            args = "boxfort-worker version -j";

            GoUint32 errcode;
            char output[JSON_BIG_FILE_SIZE];

            Config__Handle configHandle;
            App__Handle appHandle;
            errcode = SKY_cli_LoadConfig(&configHandle);
            cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
            registerHandleClose(configHandle);
            errcode = SKY_cli_NewApp(configHandle, &appHandle);
            cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
            registerHandleClose(appHandle);

  // Redirect standard output to a pipe
            GoString Args = {args, strlen(args)};
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d", errcode);

  // JSON parse output
            json_char *json;
            json_value *value;
            json_value *json_str;
            int result;
            json = (json_char *)output;
            value = json_parse(json, strlen(output));
            cr_assert(value != NULL, "Failed to json parse ");
            registerJsonFree(value);
            json_value *ver;

            ver = json_get_string(value, "skycoin");
            cr_assert(ver != NULL, "Not get string ");
            cr_assert_str_not_empty(ver->u.string.ptr);

            ver = json_get_string(value, "cli");
            cr_assert(ver != NULL, "Not get string ");
            cr_assert_str_not_empty(ver->u.string.ptr);

            ver = json_get_string(value, "rpc");
            cr_assert(ver != NULL, "Not get string ");
            cr_assert_str_not_empty(ver->u.string.ptr);

            ver = json_get_string(value, "wallet");
            cr_assert(ver != NULL, "Not get string ");
            cr_assert_str_not_empty(ver->u.string.ptr);

            Args.p = "boxfort version";
            Args.n = 15;
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            int count = getCountWordWithSeparator(output, "\n");
            cr_assert(count == 4);
          }

          Test(cli_integration, TestStableGenerateWallet) {
            useClient();
            char *name;
            char *args;

  // generate wallet with -r option
            name = "generate wallet with -r option";
            args = "boxfort-worker generateWallet -r";
            createWalletDir();

            Config__Handle configHandle;
            App__Handle appHandle;
            char output[JSON_BIG_FILE_SIZE];
            GoInt32 errcode;

            SKY_cli_LoadConfig(&configHandle);
            cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed in $s", name);
            registerHandleClose(configHandle);
            errcode = SKY_cli_NewApp(configHandle, &appHandle);
            cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed in %s", name);
            registerHandleClose(appHandle);

  // Redirect standard output to a pipe
            GoString Args = {args, strlen(args)};
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d in %s", errcode,
              name);

  // JSON parse output
            json_char *json;
            json_value *value;

            int result;
            json = (json_char *)output;
            value = json_parse(json, strlen(output));
            cr_assert(value != NULL, "Failed to json parse in %s", name);
            registerJsonFree(value);

            json_value *meta_seed;
            json_value *meta_filename;
            json_value *meta_label;

            meta_filename = get_json_value(value, "meta/filename", json_string);
            cr_assert(strcmp(meta_filename->u.string.ptr, "skycoin_cli.wlt") == 0);

            meta_seed = get_json_value(value, "meta/seed", json_string);
            cr_assert_str_not_empty(meta_seed->u.string.ptr);

            meta_label = get_json_value(value, "meta/label", json_string);
            cr_assert_str_empty(meta_label->u.string.ptr);

  // generate wallet with --rd option
            createWalletDir();
            name = "generate wallet with --rd option";
            args = "boxfort-worker generateWallet --rd";
            Args.p = args;
            Args.n = strlen(args);
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d in %s", errcode,
              name);

  // JSON parse output
            json = (json_char *)output;
            value = json_parse(json, strlen(output));
            cr_assert(value != NULL, "Failed to json parse in %s", name);
            registerJsonFree(value);

            meta_filename = get_json_value(value, "meta/filename", json_string);
            cr_assert(strcmp(meta_filename->u.string.ptr, "skycoin_cli.wlt") == 0);

            meta_seed = get_json_value(value, "meta/seed", json_string);
            getCountWord(meta_seed->u.string.ptr);
            cr_assert(getCountWord(meta_seed->u.string.ptr) == 12);

            meta_label = get_json_value(value, "meta/label", json_string);
            cr_assert_str_empty(meta_label->u.string.ptr);

  // generate wallet with -s option
            createWalletDir();
            name = "generate wallet with -s option";
            args = "boxfort-worker generateWallet -s \"great duck trophy inhale dad pluck include maze smart mechanic ring merge\" ";
            Args.p = args;
            Args.n = strlen(args);
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d in %s", errcode,
              name);
   // JSON parse output
            json = (json_char *)output;
            value = json_parse(json, strlen(output));
            cr_assert(value != NULL, "Failed to json parse in %s", name);
            registerJsonFree(value);

            meta_filename = get_json_value(value, "meta/filename", json_string);
            cr_assert(strcmp(meta_filename->u.string.ptr, "skycoin_cli.wlt") == 0);

            meta_seed = get_json_value(value, "meta/seed", json_string);
            getCountWord(meta_seed->u.string.ptr);
            cr_assert(strcmp(meta_seed->u.string.ptr,
             "great duck trophy inhale dad pluck include maze smart "
             "mechanic ring merge") == 0);

            meta_label = get_json_value(value, "meta/label", json_string);
            cr_assert_str_empty(meta_label->u.string.ptr);

            json_value * entries_array;
            entries_array = get_json_value(value,"entries",json_array);
            entries_array = entries_array->u.array.values[0];
            char *entries_address =
            entries_array->u.object.values[0].value->u.string.ptr;
            char *entries_public =
            entries_array->u.object.values[1].value->u.string.ptr;
            char *entries_secret_key =
            entries_array->u.object.values[2].value->u.string.ptr;

            cr_assert(strcmp(entries_address, "2amA8sxKJhNRp3wfWrE5JfTEUjr9S3C2BaU") == 0);
            cr_assert(strcmp(entries_public, "02b4a4b63f2f8ba56f9508712815eca3c088693333715eaf7a73275d8928e1be5a") == 0);
            cr_assert(strcmp(entries_secret_key, "f4a281d094a6e9e95a84c23701a7d01a0e413c838758e94ad86a10b9b83e0434") == 0);

 // generate wallet with -n option
            createWalletDir();
            name = "generate wallet with -n option";
            args = "boxfort-worker generateWallet -n 5";
            Args.p = args;
            Args.n = strlen(args);
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d in %s", errcode,
              name);
   // JSON parse output
            json = (json_char *)output;
            value = json_parse(json, strlen(output));
            cr_assert(value != NULL, "Failed to json parse in %s", name);
            registerJsonFree(value);

            meta_filename = get_json_value(value, "meta/filename", json_string);
            cr_assert(strcmp(meta_filename->u.string.ptr, "skycoin_cli.wlt") == 0);

            meta_label = get_json_value(value, "meta/label", json_string);
            cr_assert_str_empty(meta_label->u.string.ptr);

            entries_array = get_json_value(value,"entries",json_array);

            cr_assert(entries_array->u.array.length == 5);


 // generate wallet with -f option
            createWalletDir();
            name = "generate wallet with -f option";
            args = "boxfort-worker generateWallet -f integration-cli.wlt";
            Args.p = args;
            Args.n = strlen(args);
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d in %s", errcode,
              name);
   // JSON parse output
            json = (json_char *)output;
            value = json_parse(json, strlen(output));
            cr_assert(value != NULL, "Failed to json parse in %s", name);
            registerJsonFree(value);

            meta_filename = get_json_value(value, "meta/filename", json_string);
            cr_assert(strcmp(meta_filename->u.string.ptr, "integration-cli.wlt") == 0);

            meta_label = get_json_value(value, "meta/label", json_string);
            cr_assert_str_empty(meta_label->u.string.ptr);

             // generate wallet with -l option
            createWalletDir();
            name = "generate wallet with -l option";
            args = "boxfort-worker generateWallet -l integration-cli";
            Args.p = args;
            Args.n = strlen(args);
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d in %s", errcode,
              name);
   // JSON parse output
            json = (json_char *)output;
            value = json_parse(json, strlen(output));
            cr_assert(value != NULL, "Failed to json parse in %s", name);
            registerJsonFree(value);

            meta_filename = get_json_value(value, "meta/filename", json_string);
            cr_assert(strcmp(meta_filename->u.string.ptr, "skycoin_cli.wlt") == 0);

            meta_label = get_json_value(value, "meta/label", json_string);
            cr_assert(strcmp(meta_label->u.string.ptr, "integration-cli") == 0);


             // generate wallet with duplicate wallet name
            createTempWalletDir(false);
            name = "generate wallet with duplicate wallet name";
            args = "boxfort-worker generateWallet";
            Args.p = args;
            Args.n = strlen(args);
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            cr_assert(errcode == SKY_OK, "SKY_cli_App_Run failed in return %d in %s", errcode,
              name);


             // encrypt=true
            createWalletDir();
            name = "encrypt=true";
            args = "boxfort-worker generateWallet -e -p pwd";
            Args.p = args;
            Args.n = strlen(args);
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            cr_assert(errcode == 0, "SKY_cli_App_Run failed in return %d in %s", errcode,
              name);
   // JSON parse output
            json = (json_char *)output;
            value = json_parse(json, strlen(output));
            cr_assert(value != NULL, "Failed to json parse in %s", name);
            registerJsonFree(value);

            meta_filename = get_json_value(value, "meta/filename", json_string);
            cr_assert(strcmp(meta_filename->u.string.ptr, "skycoin_cli.wlt") == 0);

            meta_seed = get_json_value(value, "meta/seed", json_string);
            cr_assert_str_empty(meta_seed->u.string.ptr);

            json_value *meta_lastSeed = get_json_value(value, "meta/lastSeed", json_string);
            cr_assert_str_empty(meta_lastSeed->u.string.ptr);

            entries_array = get_json_value(value,"entries",json_array);

            for(int i=0;i< entries_array->u.array.length;i++ ){
              json_value *entries = entries_array->u.array.values[i];
              char *entries_secret_key =
              entries->u.object.values[2].value->u.string.ptr;
              cr_assert(strlen(entries_secret_key)==0);
            }

             // encrypt=true
            createWalletDir();
            name = "encrypt=false password=pwd";
            args = "boxfort-worker generateWallet  -p pwd";
            Args.p = args;
            Args.n = strlen(args);
            redirectStdOut();
            errcode = SKY_cli_App_Run(appHandle, Args);
  // Get redirected standard output
            getStdOut(output, JSON_BIG_FILE_SIZE);
            cr_assert(errcode == -1, "SKY_cli_App_Run failed in return %d in %s", errcode,
              name);

          }
