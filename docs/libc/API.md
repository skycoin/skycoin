# Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`define `[`BUFFER_SIZE`](#cipher_8testsuite_8testsuite_8go_8h_1a6b20d41d6252e9871430c242cb1a56e7)            |
`define `[`STRING_SIZE`](#cipher_8testsuite_8testsuite_8go_8h_1ad78224efe1d3fb39b67ca74ad9d9eec7)            |
`define `[`JSON_FILE_SIZE`](#cipher_8testsuite_8testsuite_8go_8h_1aff447440daa595595664e192e1c01d81)            |
`define `[`JSON_BIG_FILE_SIZE`](#cipher_8testsuite_8testsuite_8go_8h_1a10f4e0e5aa36596ea0886620e02feb49)            |
`define `[`FILEPATH_SEPARATOR`](#cipher_8testsuite_8testsuite_8go_8h_1a6e456d1a7dded40d4dd4fd854c4e81ec)            |
`define `[`TEST_DATA_DIR`](#cipher_8testsuite_8testsuite_8go_8h_1a45050bf269268f85a0a8b2d805b334fc)            |
`define `[`MANY_ADDRESSES_FILENAME`](#cipher_8testsuite_8testsuite_8go_8h_1a6a45cb422542b704977e95e3b843cfba)            |
`define `[`INPUT_HASHES_FILENAME`](#cipher_8testsuite_8testsuite_8go_8h_1a1bbb758f4454d01355b35a36fa3d7b61)            |
`define `[`SEED_FILE_REGEX`](#cipher_8testsuite_8testsuite_8go_8h_1a7ae7c06af79cc79a9ad88e8719df1c9e)            |
`define `[`json_char`](#json_8h_1ae3b21f339690a966e921fe2545939862)            |
`define `[`json_int_t`](#json_8h_1ae8ad072e93f8e6584af231de2f592fc6)            |
`define `[`json_enable_comments`](#json_8h_1a893db2e62d8fbf36b27bcea8654f1105)            |
`define `[`json_error_max`](#json_8h_1a399c15929bed85a9a41bd4cba9703204)            |
`define `[`GO_CGO_EXPORT_PROLOGUE_H`](#libskycoin_8h_1ac91211782906f9494d827fe6e0b2e190)            |
`define `[`GO_CGO_PROLOGUE_H`](#libskycoin_8h_1ad45a58cf8a40d22e35017cb53dd6055a)            |
`define `[`SKY_OK`](#skyerrors_8h_1a5cd9ddcf04c6f149c283c805c7d296da)            |
`define `[`SKY_ERROR`](#skyerrors_8h_1a8405baf075a12e6232d75a8432d44f81)            |
`enum `[`json_type`](#json_8h_1ac75c61993722a9b8aaa44704072ec06c)            |
`public unsigned int `[`b64_int`](#base64_8h_1a0a6be6c96f28086f36d03676296a9372)`(unsigned int ch)`            |
`public unsigned int `[`b64e_size`](#base64_8h_1ae530f943b1ac55252c7ffba9a56fe946)`(unsigned int in_size)`            |
`public unsigned int `[`b64d_size`](#base64_8h_1ae6911453bae790c4ba1933674d51c4cb)`(unsigned int in_size)`            |
`public unsigned int `[`b64_encode`](#base64_8h_1aeddff3b5b68b9080553c10ff2364cc4b)`(const unsigned char * in,unsigned int in_len,unsigned char * out)`            |
`public unsigned int `[`b64_decode`](#base64_8h_1a181a008944edb84bcfd73efacadb41c5)`(const unsigned char * in,unsigned int in_len,unsigned char * out)`            |
`public unsigned int `[`b64_encodef`](#base64_8h_1a2ea67610bb294c8d82deed5d9335d877)`(char * InFile,char * OutFile)`            |
`public unsigned int `[`b64_decodef`](#base64_8h_1ac6582b011d3ebb0af9e47b5ee5d75a2c)`(char * InFile,char * OutFile)`            |
`public `[`json_value`](#struct__json__value)` * `[`loadGoldenFile`](#cipher_8testsuite_8testsuite_8go_8h_1a97f400dcf2127780240374f791ad55cb)`(const char * file)`            |
`public `[`InputTestDataJSON`](#struct_input_test_data_j_s_o_n)` * `[`jsonToInputTestData`](#cipher_8testsuite_8testsuite_8go_8h_1acc7925cd0a944333c5b0efef5926eee5)`(`[`json_value`](#struct__json__value)` * json,`[`InputTestDataJSON`](#struct_input_test_data_j_s_o_n)` * input_data)`            |
`public `[`InputTestData`](#struct_input_test_data)` * `[`registerInputTestDataCleanup`](#cipher_8testsuite_8testsuite_8go_8h_1a4940377b4eca6c3728e034a423d5964f)`(`[`InputTestData`](#struct_input_test_data)` * input_data)`            |
`public `[`InputTestDataJSON`](#struct_input_test_data_j_s_o_n)` * `[`registerInputTestDataJSONCleanup`](#cipher_8testsuite_8testsuite_8go_8h_1a45e3d42f145d2065097b1715037973bf)`(`[`InputTestDataJSON`](#struct_input_test_data_j_s_o_n)` * input_data)`            |
`public void `[`InputTestDataToJSON`](#cipher_8testsuite_8testsuite_8go_8h_1a655f8018783ebb3ff95e7f6fcd392552)`(`[`InputTestData`](#struct_input_test_data)` * input_data,`[`InputTestDataJSON`](#struct_input_test_data_j_s_o_n)` * json_data)`            |
`public GoUint32 `[`InputTestDataFromJSON`](#cipher_8testsuite_8testsuite_8go_8h_1ab3c102e440e98e8cd90a970ce0ade222)`(`[`InputTestDataJSON`](#struct_input_test_data_j_s_o_n)` * json_data,`[`InputTestData`](#struct_input_test_data)` * input_data)`            |
`public `[`KeysTestDataJSON`](#struct_keys_test_data_j_s_o_n)` * `[`jsonToKeysTestData`](#cipher_8testsuite_8testsuite_8go_8h_1adc41ca999e05fb40c7ee8c2b6a59e1bc)`(`[`json_value`](#struct__json__value)` * json,`[`KeysTestDataJSON`](#struct_keys_test_data_j_s_o_n)` * input_data)`            |
`public `[`KeysTestData`](#struct_keys_test_data)` * `[`registerKeysTestDataCleanup`](#cipher_8testsuite_8testsuite_8go_8h_1a28dfd2ed9cfbf047a378e74ba026f053)`(`[`KeysTestData`](#struct_keys_test_data)` * input_data)`            |
`public `[`KeysTestDataJSON`](#struct_keys_test_data_j_s_o_n)` * `[`registerKeysTestDataJSONCleanup`](#cipher_8testsuite_8testsuite_8go_8h_1a333892fc9f11cb43ce5b442d3b03f006)`(`[`KeysTestDataJSON`](#struct_keys_test_data_j_s_o_n)` * input_data)`            |
`public void `[`KeysTestDataToJson`](#cipher_8testsuite_8testsuite_8go_8h_1af9f0285478f247557b6334618d2aa145)`(`[`KeysTestData`](#struct_keys_test_data)` * input_data,`[`KeysTestDataJSON`](#struct_keys_test_data_j_s_o_n)` * json_data)`            |
`public GoUint32 `[`KeysTestDataFromJSON`](#cipher_8testsuite_8testsuite_8go_8h_1a92c62c143ec39d7186c19ab6e43dc92b)`(`[`KeysTestDataJSON`](#struct_keys_test_data_j_s_o_n)` * json_data,`[`KeysTestData`](#struct_keys_test_data)` * input_data)`            |
`public `[`SeedTestDataJSON`](#struct_seed_test_data_j_s_o_n)` * `[`jsonToSeedTestData`](#cipher_8testsuite_8testsuite_8go_8h_1af295a0cdc63ba27fd346852112f28b3c)`(`[`json_value`](#struct__json__value)` * json,`[`SeedTestDataJSON`](#struct_seed_test_data_j_s_o_n)` * input_data)`            |
`public `[`SeedTestData`](#struct_seed_test_data)` * `[`registerSeedTestDataCleanup`](#cipher_8testsuite_8testsuite_8go_8h_1a96bf06d5429ca3ec74e31e279f7f0cf0)`(`[`SeedTestData`](#struct_seed_test_data)` * input_data)`            |
`public `[`SeedTestDataJSON`](#struct_seed_test_data_j_s_o_n)` * `[`registerSeedTestDataJSONCleanup`](#cipher_8testsuite_8testsuite_8go_8h_1ae7d0af21cce697cc8bc097279c3ed398)`(`[`SeedTestDataJSON`](#struct_seed_test_data_j_s_o_n)` * input_data)`            |
`public void `[`SeedTestDataToJson`](#cipher_8testsuite_8testsuite_8go_8h_1abde3615ebe8efb1d0c126d4b71120bd7)`(`[`SeedTestData`](#struct_seed_test_data)` * input_data,`[`SeedTestDataJSON`](#struct_seed_test_data_j_s_o_n)` * json_data)`            |
`public GoUint32 `[`SeedTestDataFromJSON`](#cipher_8testsuite_8testsuite_8go_8h_1a7fb0e2dc54e3623caf580bab823255e5)`(`[`SeedTestDataJSON`](#struct_seed_test_data_j_s_o_n)` * json_data,`[`SeedTestData`](#struct_seed_test_data)` * input_data)`            |
`public void `[`ValidateSeedData`](#cipher_8testsuite_8testsuite_8go_8h_1a2af7891924708bd79f29b5ced351967a)`(`[`SeedTestData`](#struct_seed_test_data)` * seedData,`[`InputTestData`](#struct_input_test_data)` * inputData)`            |
`public `[`json_value`](#struct__json__value)` * `[`json_parse`](#json_8h_1a4dd0cf45ec85a69a6021b6cfe0287b66)`(const json_char * json,size_t length)`            |
`public `[`json_value`](#struct__json__value)` * `[`json_parse_ex`](#json_8h_1ae828aab0174a7e20eec19a40d835d3c1)`(`[`json_settings`](#structjson__settings)` * settings,const json_char * json,size_t length,char * error)`            |
`public void `[`json_value_free`](#json_8h_1a3299652febea64fc59c5917ad47ede28)`(`[`json_value`](#struct__json__value)` *)`            |
`public void `[`json_value_free_ex`](#json_8h_1a129467197843210b7ce2a2c59e92f781)`(`[`json_settings`](#structjson__settings)` * settings,`[`json_value`](#struct__json__value)` *)`            |
`public int `[`DecodeBase58Address`](#libskycoin_8h_1af0bc416968a2873cf5952eecd12f1f92)`(`[`GoString`](#struct___go_string__)` p0,`[`Address`](#struct_address)` * p1)`            |
`public int `[`cr_user_cipher__Address_eq`](#skycriterion_8h_1a00edb99a770d440315ad2c91107a314b)`(`[`cipher__Address`](#structcipher_____address)` * addr1,`[`cipher__Address`](#structcipher_____address)` * addr2)`            |
`public char * `[`cr_user_cipher__Address_tostr`](#skycriterion_8h_1a7684fd986e502ffa967626174ecb121b)`(`[`cipher__Address`](#structcipher_____address)` * addr1)`            |
`public int `[`cr_user_cipher__Address_noteq`](#skycriterion_8h_1aef473b3aaf9517e054136597befabb19)`(`[`cipher__Address`](#structcipher_____address)` * addr1,`[`cipher__Address`](#structcipher_____address)` * addr2)`            |
`public int `[`cr_user_GoString_eq`](#skycriterion_8h_1afde184bfa3d42dadb560478bb384fd0e)`(`[`GoString`](#struct___go_string__)` * string1,`[`GoString`](#struct___go_string__)` * string2)`            |
`public int `[`cr_user_GoString__eq`](#skycriterion_8h_1adc4957c85581c8021d1bc5e1fe68954e)`(`[`GoString_`](#struct_go_string__)` * string1,`[`GoString_`](#struct_go_string__)` * string2)`            |
`public char * `[`cr_user_GoString_tostr`](#skycriterion_8h_1ac49e1ea1279ec23eb1b06fc4cff4346e)`(`[`GoString`](#struct___go_string__)` * string)`            |
`public char * `[`cr_user_GoString__tostr`](#skycriterion_8h_1a8ba00c85c7eede2d955cfe016cb1023d)`(`[`GoString_`](#struct_go_string__)` * string)`            |
`public int `[`cr_user_cipher__SecKey_eq`](#skycriterion_8h_1ac3d286c06a1659717bc392004b857ba0)`(cipher__SecKey * seckey1,cipher__SecKey * seckey2)`            |
`public char * `[`cr_user_cipher__SecKey_tostr`](#skycriterion_8h_1ad6ebdf1335f21b53df8a9606e68889af)`(cipher__SecKey * seckey1)`            |
`public int `[`cr_user_cipher__Ripemd160_noteq`](#skycriterion_8h_1a817a01cd72b552039a61d0add95a14f1)`(cipher__Ripemd160 * rp1,cipher__Ripemd160 * rp2)`            |
`public int `[`cr_user_cipher__Ripemd160_eq`](#skycriterion_8h_1ac843b722e627d29c8450c766b866edd6)`(cipher__Ripemd160 * rp1,cipher__Ripemd160 * rp2)`            |
`public char * `[`cr_user_cipher__Ripemd160_tostr`](#skycriterion_8h_1aeee5ce2262f3cfbdc942aa2ae6c16974)`(cipher__Ripemd160 * rp1)`            |
`public int `[`cr_user_GoSlice_eq`](#skycriterion_8h_1a68e13a153f444839e3dbe06cc14e2348)`(`[`GoSlice`](#struct_go_slice)` * slice1,`[`GoSlice`](#struct_go_slice)` * slice2)`            |
`public char * `[`cr_user_GoSlice_tostr`](#skycriterion_8h_1aa058100c8835ae72f2c609ad2ef1ba85)`(`[`GoSlice`](#struct_go_slice)` * slice1)`            |
`public int `[`cr_user_GoSlice_noteq`](#skycriterion_8h_1a361dbb4ff75151c68df6d37368880b24)`(`[`GoSlice`](#struct_go_slice)` * slice1,`[`GoSlice`](#struct_go_slice)` * slice2)`            |
`public int `[`cr_user_cipher__SHA256_noteq`](#skycriterion_8h_1af9de60cc1f5b338ff07b8a34a0af31d4)`(cipher__SHA256 * sh1,cipher__SHA256 * sh2)`            |
`public int `[`cr_user_cipher__SHA256_eq`](#skycriterion_8h_1a239dd96034613f7c5d187f355d054843)`(cipher__SHA256 * sh1,cipher__SHA256 * sh2)`            |
`public char * `[`cr_user_cipher__SHA256_tostr`](#skycriterion_8h_1a24e2acadd37b726ed2d9b07092e739d9)`(cipher__SHA256 * sh1)`            |
`public void `[`randBytes`](#skystring_8h_1abc646fb4e2f83b9ec86bacd6f8006907)`(`[`GoSlice`](#struct_go_slice)` * bytes,size_t n)`            |
`public void `[`strnhex`](#skystring_8h_1aef6e4f140a965b05589db78792dc3c09)`(unsigned char * buf,char * str,int n)`            |
`public void `[`strhex`](#skystring_8h_1a589986670c6a1cd947da79512078ff05)`(unsigned char * buf,char * str)`            |
`public void `[`fprintbuff`](#skytest_8h_1a1ee45e153c115a9a735b3ccbf992e495)`(FILE * f,void * buff,size_t n)`            |
`public `[`json_value`](#struct__json__value)` * `[`loadJsonFile`](#skytest_8h_1ae9debe21347a5e30565195425a898448)`(const char * filename)`            |
`public void * `[`registerMemCleanup`](#skytest_8h_1a3138ecc83c1c8906c84ef5e0d54cdfbb)`(void * p)`            |
`public void `[`toGoString`](#skytest_8h_1a1bad90cc197623fa8328f71809dda1a3)`(`[`GoString_`](#struct_go_string__)` * s,`[`GoString`](#struct___go_string__)` * r)`            |
`public `[`json_value`](#struct__json__value)` * `[`json_get_string`](#skytest_8h_1aaa0fcc92ec99e3682126c396c7990724)`(`[`json_value`](#struct__json__value)` * value,const char * key)`            |
`public int `[`json_set_string`](#skytest_8h_1a532c9bff6c467be505835bcb9a2fcaa0)`(`[`json_value`](#struct__json__value)` * value,const char * new_string_value)`            |
`public int `[`registerJsonFree`](#skytest_8h_1af298033f76a79ff6945d0a65a83842fe)`(void * p)`            |
`public void `[`freeRegisteredJson`](#skytest_8h_1ae9eb5fd6a792b4db185089a1a14db57a)`(void * p)`            |
`public int `[`compareJsonValues`](#skytest_8h_1abf7dd2e2fa866fc3fe6acedbf9842a3a)`(`[`json_value`](#struct__json__value)` * value1,`[`json_value`](#struct__json__value)` * value2)`            |
`public `[`json_value`](#struct__json__value)` * `[`get_json_value`](#skytest_8h_1ad81ecc74786fda501485b271384b993e)`(`[`json_value`](#struct__json__value)` * node,const char * path,json_type type)`            |
`public `[`json_value`](#struct__json__value)` * `[`get_json_value_not_strict`](#skytest_8h_1aabc1cef09feea6fdb43c33971366b1b5)`(`[`json_value`](#struct__json__value)` * node,const char * path,json_type type,int allow_null)`            |
`public void `[`setup`](#skytest_8h_1a7dfd9b79bc5a37d7df40207afbc5431f)`(void)`            |
`public void `[`teardown`](#skytest_8h_1a75dbff5b7c2c889050e2c49172679905)`(void)`            |
`struct `[`_GoString_`](#struct___go_string__) |
`struct `[`_json_object_entry`](#struct__json__object__entry) |
`struct `[`_json_value`](#struct__json__value) |
`struct `[`Address`](#struct_address) |
`struct `[`cipher__Address`](#structcipher_____address) | Addresses of SKY accounts
`struct `[`cipher__BitcoinAddress`](#structcipher_____bitcoin_address) | Addresses of Bitcoin accounts
`struct `[`cli__SendAmount`](#structcli_____send_amount) | Structure used to specify amounts transferred in a transaction.
`struct `[`coin__Transaction`](#structcoin_____transaction) | Skycoin transaction.
`struct `[`coin__TransactionOutput`](#structcoin_____transaction_output) | Skycoin transaction output.
`struct `[`coin__UxBody`](#structcoin_____ux_body) |
`struct `[`coin__UxHead`](#structcoin_____ux_head) |
`struct `[`coin__UxOut`](#structcoin_____ux_out) |
`struct `[`GoInterface`](#struct_go_interface) |
`struct `[`GoInterface_`](#struct_go_interface__) | Instances of Go interface types.
`struct `[`GoSlice`](#struct_go_slice) |
`struct `[`GoSlice_`](#struct_go_slice__) | Instances of Go slices
`struct `[`GoString_`](#struct_go_string__) | Instances of Go `string` type.
`struct `[`InputTestData`](#struct_input_test_data) |
`struct `[`InputTestDataJSON`](#struct_input_test_data_j_s_o_n) |
`struct `[`json_settings`](#structjson__settings) |
`struct `[`KeysTestData`](#struct_keys_test_data) |
`struct `[`KeysTestDataJSON`](#struct_keys_test_data_j_s_o_n) |
`struct `[`SeedTestData`](#struct_seed_test_data) |
`struct `[`SeedTestDataJSON`](#struct_seed_test_data_j_s_o_n) |
`struct `[`wallet__Entry`](#structwallet_____entry) | Wallet entry.
`struct `[`wallet__UxBalance`](#structwallet_____ux_balance) | Intermediate representation of a UxOut for sorting and spend choosing.
`struct `[`wallet__Wallet`](#structwallet_____wallet) | Internal representation of a Skycoin wallet.

## Members

#### `define `[`BUFFER_SIZE`](#cipher_8testsuite_8testsuite_8go_8h_1a6b20d41d6252e9871430c242cb1a56e7)

#### `define `[`STRING_SIZE`](#cipher_8testsuite_8testsuite_8go_8h_1ad78224efe1d3fb39b67ca74ad9d9eec7)

#### `define `[`JSON_FILE_SIZE`](#cipher_8testsuite_8testsuite_8go_8h_1aff447440daa595595664e192e1c01d81)

#### `define `[`JSON_BIG_FILE_SIZE`](#cipher_8testsuite_8testsuite_8go_8h_1a10f4e0e5aa36596ea0886620e02feb49)

#### `define `[`FILEPATH_SEPARATOR`](#cipher_8testsuite_8testsuite_8go_8h_1a6e456d1a7dded40d4dd4fd854c4e81ec)

#### `define `[`TEST_DATA_DIR`](#cipher_8testsuite_8testsuite_8go_8h_1a45050bf269268f85a0a8b2d805b334fc)

#### `define `[`MANY_ADDRESSES_FILENAME`](#cipher_8testsuite_8testsuite_8go_8h_1a6a45cb422542b704977e95e3b843cfba)

#### `define `[`INPUT_HASHES_FILENAME`](#cipher_8testsuite_8testsuite_8go_8h_1a1bbb758f4454d01355b35a36fa3d7b61)

#### `define `[`SEED_FILE_REGEX`](#cipher_8testsuite_8testsuite_8go_8h_1a7ae7c06af79cc79a9ad88e8719df1c9e)

#### `define `[`json_char`](#json_8h_1ae3b21f339690a966e921fe2545939862)

#### `define `[`json_int_t`](#json_8h_1ae8ad072e93f8e6584af231de2f592fc6)

#### `define `[`json_enable_comments`](#json_8h_1a893db2e62d8fbf36b27bcea8654f1105)

#### `define `[`json_error_max`](#json_8h_1a399c15929bed85a9a41bd4cba9703204)

#### `define `[`GO_CGO_EXPORT_PROLOGUE_H`](#libskycoin_8h_1ac91211782906f9494d827fe6e0b2e190)

#### `define `[`GO_CGO_PROLOGUE_H`](#libskycoin_8h_1ad45a58cf8a40d22e35017cb53dd6055a)

#### `define `[`SKY_OK`](#skyerrors_8h_1a5cd9ddcf04c6f149c283c805c7d296da)

#### `define `[`SKY_ERROR`](#skyerrors_8h_1a8405baf075a12e6232d75a8432d44f81)

#### `enum `[`json_type`](#json_8h_1ac75c61993722a9b8aaa44704072ec06c)

 Values                         | Descriptions
--------------------------------|---------------------------------------------
json_none            |
json_object            |
json_array            |
json_integer            |
json_double            |
json_string            |
json_boolean            |
json_null            |

#### `public unsigned int `[`b64_int`](#base64_8h_1a0a6be6c96f28086f36d03676296a9372)`(unsigned int ch)`

#### `public unsigned int `[`b64e_size`](#base64_8h_1ae530f943b1ac55252c7ffba9a56fe946)`(unsigned int in_size)`

#### `public unsigned int `[`b64d_size`](#base64_8h_1ae6911453bae790c4ba1933674d51c4cb)`(unsigned int in_size)`

#### `public unsigned int `[`b64_encode`](#base64_8h_1aeddff3b5b68b9080553c10ff2364cc4b)`(const unsigned char * in,unsigned int in_len,unsigned char * out)`

#### `public unsigned int `[`b64_decode`](#base64_8h_1a181a008944edb84bcfd73efacadb41c5)`(const unsigned char * in,unsigned int in_len,unsigned char * out)`

#### `public unsigned int `[`b64_encodef`](#base64_8h_1a2ea67610bb294c8d82deed5d9335d877)`(char * InFile,char * OutFile)`

#### `public unsigned int `[`b64_decodef`](#base64_8h_1ac6582b011d3ebb0af9e47b5ee5d75a2c)`(char * InFile,char * OutFile)`

#### `public `[`json_value`](#struct__json__value)` * `[`loadGoldenFile`](#cipher_8testsuite_8testsuite_8go_8h_1a97f400dcf2127780240374f791ad55cb)`(const char * file)`

#### `public `[`InputTestDataJSON`](#struct_input_test_data_j_s_o_n)` * `[`jsonToInputTestData`](#cipher_8testsuite_8testsuite_8go_8h_1acc7925cd0a944333c5b0efef5926eee5)`(`[`json_value`](#struct__json__value)` * json,`[`InputTestDataJSON`](#struct_input_test_data_j_s_o_n)` * input_data)`

#### `public `[`InputTestData`](#struct_input_test_data)` * `[`registerInputTestDataCleanup`](#cipher_8testsuite_8testsuite_8go_8h_1a4940377b4eca6c3728e034a423d5964f)`(`[`InputTestData`](#struct_input_test_data)` * input_data)`

#### `public `[`InputTestDataJSON`](#struct_input_test_data_j_s_o_n)` * `[`registerInputTestDataJSONCleanup`](#cipher_8testsuite_8testsuite_8go_8h_1a45e3d42f145d2065097b1715037973bf)`(`[`InputTestDataJSON`](#struct_input_test_data_j_s_o_n)` * input_data)`

#### `public void `[`InputTestDataToJSON`](#cipher_8testsuite_8testsuite_8go_8h_1a655f8018783ebb3ff95e7f6fcd392552)`(`[`InputTestData`](#struct_input_test_data)` * input_data,`[`InputTestDataJSON`](#struct_input_test_data_j_s_o_n)` * json_data)`

#### `public GoUint32 `[`InputTestDataFromJSON`](#cipher_8testsuite_8testsuite_8go_8h_1ab3c102e440e98e8cd90a970ce0ade222)`(`[`InputTestDataJSON`](#struct_input_test_data_j_s_o_n)` * json_data,`[`InputTestData`](#struct_input_test_data)` * input_data)`

#### `public `[`KeysTestDataJSON`](#struct_keys_test_data_j_s_o_n)` * `[`jsonToKeysTestData`](#cipher_8testsuite_8testsuite_8go_8h_1adc41ca999e05fb40c7ee8c2b6a59e1bc)`(`[`json_value`](#struct__json__value)` * json,`[`KeysTestDataJSON`](#struct_keys_test_data_j_s_o_n)` * input_data)`

#### `public `[`KeysTestData`](#struct_keys_test_data)` * `[`registerKeysTestDataCleanup`](#cipher_8testsuite_8testsuite_8go_8h_1a28dfd2ed9cfbf047a378e74ba026f053)`(`[`KeysTestData`](#struct_keys_test_data)` * input_data)`

#### `public `[`KeysTestDataJSON`](#struct_keys_test_data_j_s_o_n)` * `[`registerKeysTestDataJSONCleanup`](#cipher_8testsuite_8testsuite_8go_8h_1a333892fc9f11cb43ce5b442d3b03f006)`(`[`KeysTestDataJSON`](#struct_keys_test_data_j_s_o_n)` * input_data)`

#### `public void `[`KeysTestDataToJson`](#cipher_8testsuite_8testsuite_8go_8h_1af9f0285478f247557b6334618d2aa145)`(`[`KeysTestData`](#struct_keys_test_data)` * input_data,`[`KeysTestDataJSON`](#struct_keys_test_data_j_s_o_n)` * json_data)`

#### `public GoUint32 `[`KeysTestDataFromJSON`](#cipher_8testsuite_8testsuite_8go_8h_1a92c62c143ec39d7186c19ab6e43dc92b)`(`[`KeysTestDataJSON`](#struct_keys_test_data_j_s_o_n)` * json_data,`[`KeysTestData`](#struct_keys_test_data)` * input_data)`

#### `public `[`SeedTestDataJSON`](#struct_seed_test_data_j_s_o_n)` * `[`jsonToSeedTestData`](#cipher_8testsuite_8testsuite_8go_8h_1af295a0cdc63ba27fd346852112f28b3c)`(`[`json_value`](#struct__json__value)` * json,`[`SeedTestDataJSON`](#struct_seed_test_data_j_s_o_n)` * input_data)`

#### `public `[`SeedTestData`](#struct_seed_test_data)` * `[`registerSeedTestDataCleanup`](#cipher_8testsuite_8testsuite_8go_8h_1a96bf06d5429ca3ec74e31e279f7f0cf0)`(`[`SeedTestData`](#struct_seed_test_data)` * input_data)`

#### `public `[`SeedTestDataJSON`](#struct_seed_test_data_j_s_o_n)` * `[`registerSeedTestDataJSONCleanup`](#cipher_8testsuite_8testsuite_8go_8h_1ae7d0af21cce697cc8bc097279c3ed398)`(`[`SeedTestDataJSON`](#struct_seed_test_data_j_s_o_n)` * input_data)`

#### `public void `[`SeedTestDataToJson`](#cipher_8testsuite_8testsuite_8go_8h_1abde3615ebe8efb1d0c126d4b71120bd7)`(`[`SeedTestData`](#struct_seed_test_data)` * input_data,`[`SeedTestDataJSON`](#struct_seed_test_data_j_s_o_n)` * json_data)`

#### `public GoUint32 `[`SeedTestDataFromJSON`](#cipher_8testsuite_8testsuite_8go_8h_1a7fb0e2dc54e3623caf580bab823255e5)`(`[`SeedTestDataJSON`](#struct_seed_test_data_j_s_o_n)` * json_data,`[`SeedTestData`](#struct_seed_test_data)` * input_data)`

#### `public void `[`ValidateSeedData`](#cipher_8testsuite_8testsuite_8go_8h_1a2af7891924708bd79f29b5ced351967a)`(`[`SeedTestData`](#struct_seed_test_data)` * seedData,`[`InputTestData`](#struct_input_test_data)` * inputData)`

#### `public `[`json_value`](#struct__json__value)` * `[`json_parse`](#json_8h_1a4dd0cf45ec85a69a6021b6cfe0287b66)`(const json_char * json,size_t length)`

#### `public `[`json_value`](#struct__json__value)` * `[`json_parse_ex`](#json_8h_1ae828aab0174a7e20eec19a40d835d3c1)`(`[`json_settings`](#structjson__settings)` * settings,const json_char * json,size_t length,char * error)`

#### `public void `[`json_value_free`](#json_8h_1a3299652febea64fc59c5917ad47ede28)`(`[`json_value`](#struct__json__value)` *)`

#### `public void `[`json_value_free_ex`](#json_8h_1a129467197843210b7ce2a2c59e92f781)`(`[`json_settings`](#structjson__settings)` * settings,`[`json_value`](#struct__json__value)` *)`

#### `public int `[`DecodeBase58Address`](#libskycoin_8h_1af0bc416968a2873cf5952eecd12f1f92)`(`[`GoString`](#struct___go_string__)` p0,`[`Address`](#struct_address)` * p1)`

#### `public int `[`cr_user_cipher__Address_eq`](#skycriterion_8h_1a00edb99a770d440315ad2c91107a314b)`(`[`cipher__Address`](#structcipher_____address)` * addr1,`[`cipher__Address`](#structcipher_____address)` * addr2)`

#### `public char * `[`cr_user_cipher__Address_tostr`](#skycriterion_8h_1a7684fd986e502ffa967626174ecb121b)`(`[`cipher__Address`](#structcipher_____address)` * addr1)`

#### `public int `[`cr_user_cipher__Address_noteq`](#skycriterion_8h_1aef473b3aaf9517e054136597befabb19)`(`[`cipher__Address`](#structcipher_____address)` * addr1,`[`cipher__Address`](#structcipher_____address)` * addr2)`

#### `public int `[`cr_user_GoString_eq`](#skycriterion_8h_1afde184bfa3d42dadb560478bb384fd0e)`(`[`GoString`](#struct___go_string__)` * string1,`[`GoString`](#struct___go_string__)` * string2)`

#### `public int `[`cr_user_GoString__eq`](#skycriterion_8h_1adc4957c85581c8021d1bc5e1fe68954e)`(`[`GoString_`](#struct_go_string__)` * string1,`[`GoString_`](#struct_go_string__)` * string2)`

#### `public char * `[`cr_user_GoString_tostr`](#skycriterion_8h_1ac49e1ea1279ec23eb1b06fc4cff4346e)`(`[`GoString`](#struct___go_string__)` * string)`

#### `public char * `[`cr_user_GoString__tostr`](#skycriterion_8h_1a8ba00c85c7eede2d955cfe016cb1023d)`(`[`GoString_`](#struct_go_string__)` * string)`

#### `public int `[`cr_user_cipher__SecKey_eq`](#skycriterion_8h_1ac3d286c06a1659717bc392004b857ba0)`(cipher__SecKey * seckey1,cipher__SecKey * seckey2)`

#### `public char * `[`cr_user_cipher__SecKey_tostr`](#skycriterion_8h_1ad6ebdf1335f21b53df8a9606e68889af)`(cipher__SecKey * seckey1)`

#### `public int `[`cr_user_cipher__Ripemd160_noteq`](#skycriterion_8h_1a817a01cd72b552039a61d0add95a14f1)`(cipher__Ripemd160 * rp1,cipher__Ripemd160 * rp2)`

#### `public int `[`cr_user_cipher__Ripemd160_eq`](#skycriterion_8h_1ac843b722e627d29c8450c766b866edd6)`(cipher__Ripemd160 * rp1,cipher__Ripemd160 * rp2)`

#### `public char * `[`cr_user_cipher__Ripemd160_tostr`](#skycriterion_8h_1aeee5ce2262f3cfbdc942aa2ae6c16974)`(cipher__Ripemd160 * rp1)`

#### `public int `[`cr_user_GoSlice_eq`](#skycriterion_8h_1a68e13a153f444839e3dbe06cc14e2348)`(`[`GoSlice`](#struct_go_slice)` * slice1,`[`GoSlice`](#struct_go_slice)` * slice2)`

#### `public char * `[`cr_user_GoSlice_tostr`](#skycriterion_8h_1aa058100c8835ae72f2c609ad2ef1ba85)`(`[`GoSlice`](#struct_go_slice)` * slice1)`

#### `public int `[`cr_user_GoSlice_noteq`](#skycriterion_8h_1a361dbb4ff75151c68df6d37368880b24)`(`[`GoSlice`](#struct_go_slice)` * slice1,`[`GoSlice`](#struct_go_slice)` * slice2)`

#### `public int `[`cr_user_cipher__SHA256_noteq`](#skycriterion_8h_1af9de60cc1f5b338ff07b8a34a0af31d4)`(cipher__SHA256 * sh1,cipher__SHA256 * sh2)`

#### `public int `[`cr_user_cipher__SHA256_eq`](#skycriterion_8h_1a239dd96034613f7c5d187f355d054843)`(cipher__SHA256 * sh1,cipher__SHA256 * sh2)`

#### `public char * `[`cr_user_cipher__SHA256_tostr`](#skycriterion_8h_1a24e2acadd37b726ed2d9b07092e739d9)`(cipher__SHA256 * sh1)`

#### `public void `[`randBytes`](#skystring_8h_1abc646fb4e2f83b9ec86bacd6f8006907)`(`[`GoSlice`](#struct_go_slice)` * bytes,size_t n)`

#### `public void `[`strnhex`](#skystring_8h_1aef6e4f140a965b05589db78792dc3c09)`(unsigned char * buf,char * str,int n)`

#### `public void `[`strhex`](#skystring_8h_1a589986670c6a1cd947da79512078ff05)`(unsigned char * buf,char * str)`

#### `public void `[`fprintbuff`](#skytest_8h_1a1ee45e153c115a9a735b3ccbf992e495)`(FILE * f,void * buff,size_t n)`

#### `public `[`json_value`](#struct__json__value)` * `[`loadJsonFile`](#skytest_8h_1ae9debe21347a5e30565195425a898448)`(const char * filename)`

#### `public void * `[`registerMemCleanup`](#skytest_8h_1a3138ecc83c1c8906c84ef5e0d54cdfbb)`(void * p)`

#### `public void `[`toGoString`](#skytest_8h_1a1bad90cc197623fa8328f71809dda1a3)`(`[`GoString_`](#struct_go_string__)` * s,`[`GoString`](#struct___go_string__)` * r)`

#### `public `[`json_value`](#struct__json__value)` * `[`json_get_string`](#skytest_8h_1aaa0fcc92ec99e3682126c396c7990724)`(`[`json_value`](#struct__json__value)` * value,const char * key)`

#### `public int `[`json_set_string`](#skytest_8h_1a532c9bff6c467be505835bcb9a2fcaa0)`(`[`json_value`](#struct__json__value)` * value,const char * new_string_value)`

#### `public int `[`registerJsonFree`](#skytest_8h_1af298033f76a79ff6945d0a65a83842fe)`(void * p)`

#### `public void `[`freeRegisteredJson`](#skytest_8h_1ae9eb5fd6a792b4db185089a1a14db57a)`(void * p)`

#### `public int `[`compareJsonValues`](#skytest_8h_1abf7dd2e2fa866fc3fe6acedbf9842a3a)`(`[`json_value`](#struct__json__value)` * value1,`[`json_value`](#struct__json__value)` * value2)`

#### `public `[`json_value`](#struct__json__value)` * `[`get_json_value`](#skytest_8h_1ad81ecc74786fda501485b271384b993e)`(`[`json_value`](#struct__json__value)` * node,const char * path,json_type type)`

#### `public `[`json_value`](#struct__json__value)` * `[`get_json_value_not_strict`](#skytest_8h_1aabc1cef09feea6fdb43c33971366b1b5)`(`[`json_value`](#struct__json__value)` * node,const char * path,json_type type,int allow_null)`

#### `public void `[`setup`](#skytest_8h_1a7dfd9b79bc5a37d7df40207afbc5431f)`(void)`

#### `public void `[`teardown`](#skytest_8h_1a75dbff5b7c2c889050e2c49172679905)`(void)`

# struct `_GoString_`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public const char * `[`p`](#struct___go_string___1a6bc6b007533335efe02bafff799ec64c) |
`public ptrdiff_t `[`n`](#struct___go_string___1a52d899ae12c13f4df8ff5ee014f3a106) |

## Members

#### `public const char * `[`p`](#struct___go_string___1a6bc6b007533335efe02bafff799ec64c)

#### `public ptrdiff_t `[`n`](#struct___go_string___1a52d899ae12c13f4df8ff5ee014f3a106)

# struct `_json_object_entry`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public json_char * `[`name`](#struct__json__object__entry_1a3c3e575cdb04c92d1bd8e2ffbfd871cc) |
`public unsigned int `[`name_length`](#struct__json__object__entry_1a3f3f10ffae1e364e84a38517a174b01e) |
`public struct `[`_json_value`](#struct__json__value)` * `[`value`](#struct__json__object__entry_1ae4f19c247094dd7870bfbf798d79ac99) |

## Members

#### `public json_char * `[`name`](#struct__json__object__entry_1a3c3e575cdb04c92d1bd8e2ffbfd871cc)

#### `public unsigned int `[`name_length`](#struct__json__object__entry_1a3f3f10ffae1e364e84a38517a174b01e)

#### `public struct `[`_json_value`](#struct__json__value)` * `[`value`](#struct__json__object__entry_1ae4f19c247094dd7870bfbf798d79ac99)

# struct `_json_value`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public struct `[`_json_value`](#struct__json__value)` * `[`parent`](#struct__json__value_1aa86e0e17a210b00b008a001e866050cd) |
`public json_type `[`type`](#struct__json__value_1a5b632686c28261d4d52390dfc8dc0dcd) |
`public int `[`boolean`](#struct__json__value_1afc58e0c0df35925913174accbf1114cb) |
`public json_int_t `[`integer`](#struct__json__value_1a45f8bcb9c0a1417f182d41049a2107e5) |
`public double `[`dbl`](#struct__json__value_1a57291299e530453fdec37a931c728239) |
`public unsigned int `[`length`](#struct__json__value_1ac8d42bcd4a44e078047ccd7291059238) |
`public json_char * `[`ptr`](#struct__json__value_1af340cb5b3b56fa2cc2b043529017fd3a) |
`public struct _json_value::@0::@2 `[`string`](#struct__json__value_1adbf96c673bfce83dd224c60f5a1eedc4) |
`public `[`json_object_entry`](#struct__json__object__entry)` * `[`values`](#struct__json__value_1a35eecbf6405a59adaf2c2001b9b224b0) |
`public struct _json_value::@0::@3 `[`object`](#struct__json__value_1a67174883be078cb852d4748b78aa60ce) |
`public struct `[`_json_value`](#struct__json__value)` ** `[`values`](#struct__json__value_1a6244a16657c988883bbb1fa7f1f0ba55) |
`public struct _json_value::@0::@4 `[`array`](#struct__json__value_1ac837a4233a38fcd3c3d51da2c3a56f0c) |
`public union _json_value::@0 `[`u`](#struct__json__value_1a2e5a72cccd43ae6e6b1dc476fc327d17) |
`public struct `[`_json_value`](#struct__json__value)` * `[`next_alloc`](#struct__json__value_1a6b6c655ef17b09a6ea18d94612a27304) |
`public void * `[`object_mem`](#struct__json__value_1aa166262c3b34bfb69a85f6e625970226) |
`public union _json_value::@1 `[`_reserved`](#struct__json__value_1ac7cdbd31aad7b4a672e38721c1122f35) |

## Members

#### `public struct `[`_json_value`](#struct__json__value)` * `[`parent`](#struct__json__value_1aa86e0e17a210b00b008a001e866050cd)

#### `public json_type `[`type`](#struct__json__value_1a5b632686c28261d4d52390dfc8dc0dcd)

#### `public int `[`boolean`](#struct__json__value_1afc58e0c0df35925913174accbf1114cb)

#### `public json_int_t `[`integer`](#struct__json__value_1a45f8bcb9c0a1417f182d41049a2107e5)

#### `public double `[`dbl`](#struct__json__value_1a57291299e530453fdec37a931c728239)

#### `public unsigned int `[`length`](#struct__json__value_1ac8d42bcd4a44e078047ccd7291059238)

#### `public json_char * `[`ptr`](#struct__json__value_1af340cb5b3b56fa2cc2b043529017fd3a)

#### `public struct _json_value::@0::@2 `[`string`](#struct__json__value_1adbf96c673bfce83dd224c60f5a1eedc4)

#### `public `[`json_object_entry`](#struct__json__object__entry)` * `[`values`](#struct__json__value_1a35eecbf6405a59adaf2c2001b9b224b0)

#### `public struct _json_value::@0::@3 `[`object`](#struct__json__value_1a67174883be078cb852d4748b78aa60ce)

#### `public struct `[`_json_value`](#struct__json__value)` ** `[`values`](#struct__json__value_1a6244a16657c988883bbb1fa7f1f0ba55)

#### `public struct _json_value::@0::@4 `[`array`](#struct__json__value_1ac837a4233a38fcd3c3d51da2c3a56f0c)

#### `public union _json_value::@0 `[`u`](#struct__json__value_1a2e5a72cccd43ae6e6b1dc476fc327d17)

#### `public struct `[`_json_value`](#struct__json__value)` * `[`next_alloc`](#struct__json__value_1a6b6c655ef17b09a6ea18d94612a27304)

#### `public void * `[`object_mem`](#struct__json__value_1aa166262c3b34bfb69a85f6e625970226)

#### `public union _json_value::@1 `[`_reserved`](#struct__json__value_1ac7cdbd31aad7b4a672e38721c1122f35)

# struct `Address`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public unsigned char `[`Version`](#struct_address_1a49fed92a3e4a3cc30678924a13acc19f) |
`public Ripemd160 `[`Key`](#struct_address_1aa7fd9da55c53a8f7a6abe4987a8ea093) |

## Members

#### `public unsigned char `[`Version`](#struct_address_1a49fed92a3e4a3cc30678924a13acc19f)

#### `public Ripemd160 `[`Key`](#struct_address_1aa7fd9da55c53a8f7a6abe4987a8ea093)

# struct `cipher__Address`

Addresses of SKY accounts

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public unsigned char `[`Version`](#structcipher_____address_1a49fed92a3e4a3cc30678924a13acc19f) | [Address](#struct_address) version identifier. Used to differentiate testnet vs mainnet addresses, for instance.
`public cipher__Ripemd160 `[`Key`](#structcipher_____address_1a16586cd3bfc67010c3c185d0da01317c) | [Address](#struct_address) hash identifier.

## Members

#### `public unsigned char `[`Version`](#structcipher_____address_1a49fed92a3e4a3cc30678924a13acc19f)

[Address](#struct_address) version identifier. Used to differentiate testnet vs mainnet addresses, for instance.

#### `public cipher__Ripemd160 `[`Key`](#structcipher_____address_1a16586cd3bfc67010c3c185d0da01317c)

[Address](#struct_address) hash identifier.

# struct `cipher__BitcoinAddress`

Addresses of Bitcoin accounts

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public unsigned char `[`Version`](#structcipher_____bitcoin_address_1a49fed92a3e4a3cc30678924a13acc19f) | [Address](#struct_address) version identifier. Used to differentiate testnet vs mainnet addresses, for instance.
`public cipher__Ripemd160 `[`Key`](#structcipher_____bitcoin_address_1a16586cd3bfc67010c3c185d0da01317c) | [Address](#struct_address) hash identifier.

## Members

#### `public unsigned char `[`Version`](#structcipher_____bitcoin_address_1a49fed92a3e4a3cc30678924a13acc19f)

[Address](#struct_address) version identifier. Used to differentiate testnet vs mainnet addresses, for instance.

#### `public cipher__Ripemd160 `[`Key`](#structcipher_____bitcoin_address_1a16586cd3bfc67010c3c185d0da01317c)

[Address](#struct_address) hash identifier.

# struct `cli__SendAmount`

Structure used to specify amounts transferred in a transaction.

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public `[`GoString_`](#struct_go_string__)` `[`Addr`](#structcli_____send_amount_1af646aec99cf83d30d17fd62e014d19f8) | Sender / receipient address.
`public GoInt64_ `[`Coins`](#structcli_____send_amount_1a7733f16af3115d3cfc712f2f687b73e4) | Amount transferred (e.g. measured in SKY)

## Members

#### `public `[`GoString_`](#struct_go_string__)` `[`Addr`](#structcli_____send_amount_1af646aec99cf83d30d17fd62e014d19f8)

Sender / receipient address.

#### `public GoInt64_ `[`Coins`](#structcli_____send_amount_1a7733f16af3115d3cfc712f2f687b73e4)

Amount transferred (e.g. measured in SKY)

# struct `coin__Transaction`

Skycoin transaction.

Instances of this struct are included in blocks.

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public GoInt32_ `[`Length`](#structcoin_____transaction_1a8a3c288e7b3f7245e0b7916ce322e5f9) | Current transaction's length expressed in bytes.
`public GoInt8_ `[`Type`](#structcoin_____transaction_1a5bf4f40bde41c84f4ab5ff82bc74f744) | Transaction's version. When a node tries to process a transaction, it must verify whether it supports the transaction's type. This is intended to provide a way to update skycoin clients and servers without crashing the network. If the transaction is not compatible with the node, it should not process it.
`public cipher__SHA256 `[`InnerHash`](#structcoin_____transaction_1af5187a82f283be2deca0f1781fa628ad) | It's a SHA256 hash of the inputs and outputs of the transaction. It is used to protect against transaction mutability. This means that the transaction cannot be altered after its creation.
`public `[`GoSlice_`](#struct_go_slice__)` `[`Sigs`](#structcoin_____transaction_1af0a2ba807a16f9ae66dd5682c243b943) | A list of digital signiatures generated by the skycoin client using the private key. It is used by Skycoin servers to verify the authenticy of the transaction. Each input requires a different signature.
`public `[`GoSlice_`](#struct_go_slice__)` `[`In`](#structcoin_____transaction_1a317b73fcfa2b93fd16dfeab7ba228c39) | A list of references to unspent transaction outputs. Unlike other cryptocurrencies, such as Bitcoin, Skycoin unspent transaction outputs (UX) and Skycoin transactions (TX) are separated in the blockchain protocol, allowing for lighter transactions, thus reducing the broadcasting costs across the network.
`public `[`GoSlice_`](#struct_go_slice__)` `[`Out`](#structcoin_____transaction_1a181edcb164c89b192b3838de7792cc89) | Outputs: A list of outputs created by the client, that will be recorded in the blockchain if transactions are confirmed. An output consists of a data structure representing an UTXT, which is composed by a Skycoin address to be sent to, the amount in Skycoin to be sent, and the amount of Coin Hours to be sent, and the SHA256 hash of the previous fields.

## Members

#### `public GoInt32_ `[`Length`](#structcoin_____transaction_1a8a3c288e7b3f7245e0b7916ce322e5f9)

Current transaction's length expressed in bytes.

#### `public GoInt8_ `[`Type`](#structcoin_____transaction_1a5bf4f40bde41c84f4ab5ff82bc74f744)

Transaction's version. When a node tries to process a transaction, it must verify whether it supports the transaction's type. This is intended to provide a way to update skycoin clients and servers without crashing the network. If the transaction is not compatible with the node, it should not process it.

#### `public cipher__SHA256 `[`InnerHash`](#structcoin_____transaction_1af5187a82f283be2deca0f1781fa628ad)

It's a SHA256 hash of the inputs and outputs of the transaction. It is used to protect against transaction mutability. This means that the transaction cannot be altered after its creation.

#### `public `[`GoSlice_`](#struct_go_slice__)` `[`Sigs`](#structcoin_____transaction_1af0a2ba807a16f9ae66dd5682c243b943)

A list of digital signiatures generated by the skycoin client using the private key. It is used by Skycoin servers to verify the authenticy of the transaction. Each input requires a different signature.

#### `public `[`GoSlice_`](#struct_go_slice__)` `[`In`](#structcoin_____transaction_1a317b73fcfa2b93fd16dfeab7ba228c39)

A list of references to unspent transaction outputs. Unlike other cryptocurrencies, such as Bitcoin, Skycoin unspent transaction outputs (UX) and Skycoin transactions (TX) are separated in the blockchain protocol, allowing for lighter transactions, thus reducing the broadcasting costs across the network.

#### `public `[`GoSlice_`](#struct_go_slice__)` `[`Out`](#structcoin_____transaction_1a181edcb164c89b192b3838de7792cc89)

Outputs: A list of outputs created by the client, that will be recorded in the blockchain if transactions are confirmed. An output consists of a data structure representing an UTXT, which is composed by a Skycoin address to be sent to, the amount in Skycoin to be sent, and the amount of Coin Hours to be sent, and the SHA256 hash of the previous fields.

# struct `coin__TransactionOutput`

Skycoin transaction output.

Instances are integral part of transactions included in blocks.

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public `[`cipher__Address`](#structcipher_____address)` `[`Address`](#structcoin_____transaction_output_1a36182709852829827005f90eef8bf78c) | Receipient address.
`public GoInt64_ `[`Coins`](#structcoin_____transaction_output_1a7733f16af3115d3cfc712f2f687b73e4) | Amount sent to the receipient address.
`public GoInt64_ `[`Hours`](#structcoin_____transaction_output_1a7aef551ad5991173b5a6160fd8fe1594) | Amount of Coin Hours sent to the receipient address.

## Members

#### `public `[`cipher__Address`](#structcipher_____address)` `[`Address`](#structcoin_____transaction_output_1a36182709852829827005f90eef8bf78c)

Receipient address.

#### `public GoInt64_ `[`Coins`](#structcoin_____transaction_output_1a7733f16af3115d3cfc712f2f687b73e4)

Amount sent to the receipient address.

#### `public GoInt64_ `[`Hours`](#structcoin_____transaction_output_1a7aef551ad5991173b5a6160fd8fe1594)

Amount of Coin Hours sent to the receipient address.

# struct `coin__UxBody`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public cipher__SHA256 `[`SrcTransaction`](#structcoin_____ux_body_1a94a1b2134fcadb6f35f5a1573f771cde) |
`public `[`cipher__Address`](#structcipher_____address)` `[`Address`](#structcoin_____ux_body_1a36182709852829827005f90eef8bf78c) |
`public GoUint64_ `[`Coins`](#structcoin_____ux_body_1aa8ecc8470e80feb9af67d2e39d01b1eb) |
`public GoUint64_ `[`Hours`](#structcoin_____ux_body_1a527c44c111577ee7ae35d7a53ab11332) |

## Members

#### `public cipher__SHA256 `[`SrcTransaction`](#structcoin_____ux_body_1a94a1b2134fcadb6f35f5a1573f771cde)

#### `public `[`cipher__Address`](#structcipher_____address)` `[`Address`](#structcoin_____ux_body_1a36182709852829827005f90eef8bf78c)

#### `public GoUint64_ `[`Coins`](#structcoin_____ux_body_1aa8ecc8470e80feb9af67d2e39d01b1eb)

#### `public GoUint64_ `[`Hours`](#structcoin_____ux_body_1a527c44c111577ee7ae35d7a53ab11332)

# struct `coin__UxHead`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public GoUint64_ `[`Time`](#structcoin_____ux_head_1a1260435bd13330d2d5c2c61d044b1603) |
`public GoUint64_ `[`BkSeq`](#structcoin_____ux_head_1a5f647c638c8d6a7e2b9ba8b409bf94f2) |

## Members

#### `public GoUint64_ `[`Time`](#structcoin_____ux_head_1a1260435bd13330d2d5c2c61d044b1603)

#### `public GoUint64_ `[`BkSeq`](#structcoin_____ux_head_1a5f647c638c8d6a7e2b9ba8b409bf94f2)

# struct `coin__UxOut`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public `[`coin__UxHead`](#structcoin_____ux_head)` `[`Head`](#structcoin_____ux_out_1ad77c0fb7d1545bcda1657269eab9721c) |
`public `[`coin__UxBody`](#structcoin_____ux_body)` `[`Body`](#structcoin_____ux_out_1a29e9173d148f64e9129d98a8d9bf5ed2) |

## Members

#### `public `[`coin__UxHead`](#structcoin_____ux_head)` `[`Head`](#structcoin_____ux_out_1ad77c0fb7d1545bcda1657269eab9721c)

#### `public `[`coin__UxBody`](#structcoin_____ux_body)` `[`Body`](#structcoin_____ux_out_1a29e9173d148f64e9129d98a8d9bf5ed2)

# struct `GoInterface`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public void * `[`t`](#struct_go_interface_1a6445205ee90f5ff5131595cf7ddfcec0) |
`public void * `[`v`](#struct_go_interface_1a67806b49e20fb1170422969965db6ecb) |

## Members

#### `public void * `[`t`](#struct_go_interface_1a6445205ee90f5ff5131595cf7ddfcec0)

#### `public void * `[`v`](#struct_go_interface_1a67806b49e20fb1170422969965db6ecb)

# struct `GoInterface_`

Instances of Go interface types.

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public void * `[`t`](#struct_go_interface___1a6445205ee90f5ff5131595cf7ddfcec0) | Pointer to the information of the concrete Go type bound to this interface reference.
`public void * `[`v`](#struct_go_interface___1a67806b49e20fb1170422969965db6ecb) | Pointer to the data corresponding to the value bound to this interface type.

## Members

#### `public void * `[`t`](#struct_go_interface___1a6445205ee90f5ff5131595cf7ddfcec0)

Pointer to the information of the concrete Go type bound to this interface reference.

#### `public void * `[`v`](#struct_go_interface___1a67806b49e20fb1170422969965db6ecb)

Pointer to the data corresponding to the value bound to this interface type.

# struct `GoSlice`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public void * `[`data`](#struct_go_slice_1a735984d41155bc1032e09bece8f8d66d) |
`public GoInt `[`len`](#struct_go_slice_1abefd7e3d615fc657a761fd36bcd7296c) |
`public GoInt `[`cap`](#struct_go_slice_1a726ab221ad9e219391b7f4c9a5c5ba33) |

## Members

#### `public void * `[`data`](#struct_go_slice_1a735984d41155bc1032e09bece8f8d66d)

#### `public GoInt `[`len`](#struct_go_slice_1abefd7e3d615fc657a761fd36bcd7296c)

#### `public GoInt `[`cap`](#struct_go_slice_1a726ab221ad9e219391b7f4c9a5c5ba33)

# struct `GoSlice_`

Instances of Go slices

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public void * `[`data`](#struct_go_slice___1a735984d41155bc1032e09bece8f8d66d) | Pointer to buffer containing slice data.
`public GoInt_ `[`len`](#struct_go_slice___1af7b822f9987d08af70f228eb6dd4b7c7) | Number of items stored in slice buffer.
`public GoInt_ `[`cap`](#struct_go_slice___1abb0492bac72ee60cd3cafd152291df94) | Maximum number of items that fits in this slice considering allocated memory and item type's size.

## Members

#### `public void * `[`data`](#struct_go_slice___1a735984d41155bc1032e09bece8f8d66d)

Pointer to buffer containing slice data.

#### `public GoInt_ `[`len`](#struct_go_slice___1af7b822f9987d08af70f228eb6dd4b7c7)

Number of items stored in slice buffer.

#### `public GoInt_ `[`cap`](#struct_go_slice___1abb0492bac72ee60cd3cafd152291df94)

Maximum number of items that fits in this slice considering allocated memory and item type's size.

# struct `GoString_`

Instances of Go `string` type.

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public const char * `[`p`](#struct_go_string___1a6bc6b007533335efe02bafff799ec64c) | Pointer to string characters buffer.
`public GoInt_ `[`n`](#struct_go_string___1aa78f60eaaf1d3eb1661886a694b82b23) | String size not counting trailing `\0` char if at all included.

## Members

#### `public const char * `[`p`](#struct_go_string___1a6bc6b007533335efe02bafff799ec64c)

Pointer to string characters buffer.

#### `public GoInt_ `[`n`](#struct_go_string___1aa78f60eaaf1d3eb1661886a694b82b23)

String size not counting trailing `\0` char if at all included.

# struct `InputTestData`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public `[`GoSlice`](#struct_go_slice)` `[`Hashes`](#struct_input_test_data_1ad0e54a71fc762f7f7c7765162e1a738a) |

## Members

#### `public `[`GoSlice`](#struct_go_slice)` `[`Hashes`](#struct_input_test_data_1ad0e54a71fc762f7f7c7765162e1a738a)

# struct `InputTestDataJSON`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public `[`GoSlice`](#struct_go_slice)` `[`Hashes`](#struct_input_test_data_j_s_o_n_1ad0e54a71fc762f7f7c7765162e1a738a) |

## Members

#### `public `[`GoSlice`](#struct_go_slice)` `[`Hashes`](#struct_input_test_data_j_s_o_n_1ad0e54a71fc762f7f7c7765162e1a738a)

# struct `json_settings`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public unsigned long `[`max_memory`](#structjson__settings_1a0a75ed144583eabc50ba61e0e06abfea) |
`public int `[`settings`](#structjson__settings_1a4bc80b14d6c1a73e1f2ff9f0375dac8b) |
`public void *(* `[`mem_alloc`](#structjson__settings_1adf6905888aecb418419653d2c2ead6d8) |
`public void(* `[`mem_free`](#structjson__settings_1a0618e94c18df6d79aaa78fc57d797020) |
`public void * `[`user_data`](#structjson__settings_1a0f53d287ac7c064d1a49d4bd93ca1cb9) |
`public size_t `[`value_extra`](#structjson__settings_1af493db885ec44977143063289b5275db) |

## Members

#### `public unsigned long `[`max_memory`](#structjson__settings_1a0a75ed144583eabc50ba61e0e06abfea)

#### `public int `[`settings`](#structjson__settings_1a4bc80b14d6c1a73e1f2ff9f0375dac8b)

#### `public void *(* `[`mem_alloc`](#structjson__settings_1adf6905888aecb418419653d2c2ead6d8)

#### `public void(* `[`mem_free`](#structjson__settings_1a0618e94c18df6d79aaa78fc57d797020)

#### `public void * `[`user_data`](#structjson__settings_1a0f53d287ac7c064d1a49d4bd93ca1cb9)

#### `public size_t `[`value_extra`](#structjson__settings_1af493db885ec44977143063289b5275db)

# struct `KeysTestData`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public `[`cipher__Address`](#structcipher_____address)` `[`Address`](#struct_keys_test_data_1a36182709852829827005f90eef8bf78c) |
`public cipher__SecKey `[`Secret`](#struct_keys_test_data_1ab7c6991603c7f6e2f9a7c5b7dd580dac) |
`public cipher__PubKey `[`Public`](#struct_keys_test_data_1a0bb45c3d09a24a1c34d376092a35dfdc) |
`public `[`GoSlice`](#struct_go_slice)` `[`Signatures`](#struct_keys_test_data_1ac97779b64174c30fa4cd2386f399d3b5) |

## Members

#### `public `[`cipher__Address`](#structcipher_____address)` `[`Address`](#struct_keys_test_data_1a36182709852829827005f90eef8bf78c)

#### `public cipher__SecKey `[`Secret`](#struct_keys_test_data_1ab7c6991603c7f6e2f9a7c5b7dd580dac)

#### `public cipher__PubKey `[`Public`](#struct_keys_test_data_1a0bb45c3d09a24a1c34d376092a35dfdc)

#### `public `[`GoSlice`](#struct_go_slice)` `[`Signatures`](#struct_keys_test_data_1ac97779b64174c30fa4cd2386f399d3b5)

# struct `KeysTestDataJSON`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public `[`GoString`](#struct___go_string__)` `[`Address`](#struct_keys_test_data_j_s_o_n_1a9d21e8a82d5c1b3ad4f8ecf6f80e62ad) |
`public `[`GoString`](#struct___go_string__)` `[`Secret`](#struct_keys_test_data_j_s_o_n_1a801564a19fb6529beabfd99e1f4bb12e) |
`public `[`GoString`](#struct___go_string__)` `[`Public`](#struct_keys_test_data_j_s_o_n_1ac024a1d14cfeb43708ca2b526f176389) |
`public `[`GoSlice`](#struct_go_slice)` `[`Signatures`](#struct_keys_test_data_j_s_o_n_1ac97779b64174c30fa4cd2386f399d3b5) |

## Members

#### `public `[`GoString`](#struct___go_string__)` `[`Address`](#struct_keys_test_data_j_s_o_n_1a9d21e8a82d5c1b3ad4f8ecf6f80e62ad)

#### `public `[`GoString`](#struct___go_string__)` `[`Secret`](#struct_keys_test_data_j_s_o_n_1a801564a19fb6529beabfd99e1f4bb12e)

#### `public `[`GoString`](#struct___go_string__)` `[`Public`](#struct_keys_test_data_j_s_o_n_1ac024a1d14cfeb43708ca2b526f176389)

#### `public `[`GoSlice`](#struct_go_slice)` `[`Signatures`](#struct_keys_test_data_j_s_o_n_1ac97779b64174c30fa4cd2386f399d3b5)

# struct `SeedTestData`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public `[`GoSlice`](#struct_go_slice)` `[`Seed`](#struct_seed_test_data_1ac243323ff380d4e4df1741f04f272ad7) |
`public `[`GoSlice`](#struct_go_slice)` `[`Keys`](#struct_seed_test_data_1af7b3b5dfe2d429d7484447fc5135c3b2) |

## Members

#### `public `[`GoSlice`](#struct_go_slice)` `[`Seed`](#struct_seed_test_data_1ac243323ff380d4e4df1741f04f272ad7)

#### `public `[`GoSlice`](#struct_go_slice)` `[`Keys`](#struct_seed_test_data_1af7b3b5dfe2d429d7484447fc5135c3b2)

# struct `SeedTestDataJSON`

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public `[`GoString`](#struct___go_string__)` `[`Seed`](#struct_seed_test_data_j_s_o_n_1a276b95b641aae441b01fe35bf9ffaba1) |
`public `[`GoSlice`](#struct_go_slice)` `[`Keys`](#struct_seed_test_data_j_s_o_n_1af7b3b5dfe2d429d7484447fc5135c3b2) |

## Members

#### `public `[`GoString`](#struct___go_string__)` `[`Seed`](#struct_seed_test_data_j_s_o_n_1a276b95b641aae441b01fe35bf9ffaba1)

#### `public `[`GoSlice`](#struct_go_slice)` `[`Keys`](#struct_seed_test_data_j_s_o_n_1af7b3b5dfe2d429d7484447fc5135c3b2)

# struct `wallet__Entry`

Wallet entry.

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public `[`cipher__Address`](#structcipher_____address)` `[`Address`](#structwallet_____entry_1a36182709852829827005f90eef8bf78c) | Wallet address.
`public cipher__PubKey `[`Public`](#structwallet_____entry_1a0bb45c3d09a24a1c34d376092a35dfdc) | Public key used to generate address.
`public cipher__SecKey `[`Secret`](#structwallet_____entry_1ab7c6991603c7f6e2f9a7c5b7dd580dac) | Secret key used to generate address.

## Members

#### `public `[`cipher__Address`](#structcipher_____address)` `[`Address`](#structwallet_____entry_1a36182709852829827005f90eef8bf78c)

Wallet address.

#### `public cipher__PubKey `[`Public`](#structwallet_____entry_1a0bb45c3d09a24a1c34d376092a35dfdc)

Public key used to generate address.

#### `public cipher__SecKey `[`Secret`](#structwallet_____entry_1ab7c6991603c7f6e2f9a7c5b7dd580dac)

Secret key used to generate address.

# struct `wallet__UxBalance`

Intermediate representation of a UxOut for sorting and spend choosing.

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public cipher__SHA256 `[`Hash`](#structwallet_____ux_balance_1ac47f1b6e05da3c25722b2cc93728ee4d) | Hash of underlying UxOut.
`public GoInt64_ `[`BkSeq`](#structwallet_____ux_balance_1a1c1b05acfa8b1a65809ba4525f14a55b) | moment balance calculation is performed at.
`public `[`cipher__Address`](#structcipher_____address)` `[`Address`](#structwallet_____ux_balance_1a36182709852829827005f90eef8bf78c) | Account holder address.
`public GoInt64_ `[`Coins`](#structwallet_____ux_balance_1a7733f16af3115d3cfc712f2f687b73e4) | Coins amount (e.g. in SKY).
`public GoInt64_ `[`Hours`](#structwallet_____ux_balance_1a7aef551ad5991173b5a6160fd8fe1594) | Balance of Coin Hours generated by underlying UxOut, depending on UxOut's head time.

## Members

#### `public cipher__SHA256 `[`Hash`](#structwallet_____ux_balance_1ac47f1b6e05da3c25722b2cc93728ee4d)

Hash of underlying UxOut.

#### `public GoInt64_ `[`BkSeq`](#structwallet_____ux_balance_1a1c1b05acfa8b1a65809ba4525f14a55b)

moment balance calculation is performed at.

Block height corresponding to the

#### `public `[`cipher__Address`](#structcipher_____address)` `[`Address`](#structwallet_____ux_balance_1a36182709852829827005f90eef8bf78c)

Account holder address.

#### `public GoInt64_ `[`Coins`](#structwallet_____ux_balance_1a7733f16af3115d3cfc712f2f687b73e4)

Coins amount (e.g. in SKY).

#### `public GoInt64_ `[`Hours`](#structwallet_____ux_balance_1a7aef551ad5991173b5a6160fd8fe1594)

Balance of Coin Hours generated by underlying UxOut, depending on UxOut's head time.

# struct `wallet__Wallet`

Internal representation of a Skycoin wallet.

## Summary

 Members                        | Descriptions
--------------------------------|---------------------------------------------
`public GoMap_ `[`Meta`](#structwallet_____wallet_1acdf30a4af55c2c677ebf6fc57b27e740) | Records items that are not deterministic, like filename, lable, wallet type, secrets, etc.
`public `[`GoSlice_`](#struct_go_slice__)` `[`Entries`](#structwallet_____wallet_1a57b718d97f8db7e0bc9d9c755510951b) | Entries field stores the address entries that are deterministically generated from seed.

## Members

#### `public GoMap_ `[`Meta`](#structwallet_____wallet_1acdf30a4af55c2c677ebf6fc57b27e740)

Records items that are not deterministic, like filename, lable, wallet type, secrets, etc.

#### `public `[`GoSlice_`](#struct_go_slice__)` `[`Entries`](#structwallet_____wallet_1a57b718d97f8db7e0bc9d9c755510951b)

Entries field stores the address entries that are deterministically generated from seed.

Generated by [Moxygen](https://sourcey.com/moxygen)
