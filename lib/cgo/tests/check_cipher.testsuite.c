
#include <fts.h>
#include <stdio.h>

#include "cipher.testsuite.testsuite.go.h"

// TODO: File path utils. Move elsewhere.

// Determine if a file name matches pattern for golden dataset
// i.e. matches 'seed-\d+.golden' regex
bool isGoldenFile(const char* filename) {
  if (strncmp(filename, "seed-", 5) != 0)
    return false;
  char* ptr = (char*) filename + 5;
  if (*ptr < '0' || *ptr > '9')
    return false;
  while (*++ptr >= '0' && *ptr <='9') {}
  return strcmp(ptr, ".golden") == 0;
}

TestSuite(cipher_testsuite, .init = setup, .fini = teardown);

Test(cipher_testsuite, TestManyAddresses) {
  SeedTestDataJSON dataJSON;
  SeedTestData data;
  GoUint32 err;

  json_value* json = loadGoldenFile(MANY_ADDRESSES_FILENAME);
  cr_assert(json != NULL, "Error loading file %s", MANY_ADDRESSES_FILENAME);
  registerJsonFree(json);
  SeedTestDataJSON* dataset = jsonToSeedTestData(json, &dataJSON);
  cr_assert(dataset != NULL, "Loaded JSON golden dataset must not be NULL");
  registerSeedTestDataJSONCleanup(&dataJSON);
  err = SeedTestDataFromJSON(&dataJSON, &data);
  registerSeedTestDataCleanup(&data);
  cr_assert(err == SKY_OK, "Deserializing seed test data from JSON ... %d", err);
  ValidateSeedData(&data, NULL);
}

GoUint32 traverseGoldenFiles(const char *path, InputTestData* inputData) {
  char* _path[2];
  _path[0] = (char *) path;
  _path[1] = NULL;
  size_t i = 0;
  FTS* tree = fts_open(_path, FTS_NOCHDIR, NULL);

  if (!tree)
    return 1;
  FTSENT* node;
  while ((node = fts_read(tree))) {
    if ((node->fts_info & FTS_F) && isGoldenFile(node->fts_name)) {
      char fn[FILENAME_MAX];
      fprintf(stderr, "Golden data set %s\n", node->fts_path);
      SeedTestDataJSON seedDataJSON;
      SeedTestData seedData;

      json_value* json = loadGoldenFile(node->fts_name);
      cr_assert(json != NULL, "Error loading file %s", node->fts_name);
      SeedTestDataJSON* dataset = jsonToSeedTestData(json, &seedDataJSON);
      cr_assert(dataset != NULL, "Loaded JSON seed golden dataset must not be NULL");
      GoUint32 err = SeedTestDataFromJSON(&seedDataJSON, &seedData);
      cr_assert(err == SKY_OK, "Deserializing seed test data from JSON ... %d", err);
      ValidateSeedData(&seedData, inputData);
    }
  }
  return 0;
}

Test(cipher_testsuite, TestSeedSignatures) {
  InputTestDataJSON inputDataJSON;
  InputTestData inputData;
  GoUint32 err;

  json_value* json = loadGoldenFile(INPUT_HASHES_FILENAME);
  cr_assert(json != NULL, "Error loading file %s", INPUT_HASHES_FILENAME);
  InputTestDataJSON* dataset = jsonToInputTestData(json, &inputDataJSON);
  cr_assert(dataset != NULL, "Loaded JSON input golden dataset must not be NULL");
  err = InputTestDataFromJSON(&inputDataJSON, &inputData);
  cr_assert(err == SKY_OK, "Deserializing seed test data from JSON ... %d", err);
  err = traverseGoldenFiles(TEST_DATA_DIR, &inputData);
  cr_assert(err == 0);
}

