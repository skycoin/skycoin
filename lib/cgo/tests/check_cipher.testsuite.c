
#include <stdio.h>

#include "cipher.testsuite.testsuite.go.h"

TestSuite(cipher_testsuite, .init = setup, .fini = teardown);

Test(cipher_testsuite, TestManyAddresses) {
  SeedTestDataJSON dataJSON;
  SeedTestData data;
  GoUint32 err;

  json_value* json = loadGoldenFile(MANY_ADDRESSES_FILENAME);
  cr_assert(json != NULL, "Error loading file");
  SeedTestDataJSON* dataset = jsonToSeedTestData(json, &dataJSON);
  cr_assert(dataset != NULL, "Loaded JSON golden dataset must not be NULL");
  err = SeedTestDataFromJSON(&dataJSON, &data);
  cr_assert(err == SKY_OK, "Deserializing seed test data from JSON ... %d", err);
  ValidateSeedData(&data, NULL);
}

