# Summary

 Members                        | Descriptions                                
--------------------------------|---------------------------------------------
`define `[`GO_CGO_EXPORT_PROLOGUE_H`](#libskycoin_8h_1ac91211782906f9494d827fe6e0b2e190)            | 
`define `[`GO_CGO_PROLOGUE_H`](#libskycoin_8h_1ad45a58cf8a40d22e35017cb53dd6055a)            | 
`define `[`SKY_OK`](#skyerrors_8h_1a5cd9ddcf04c6f149c283c805c7d296da)            | 
`define `[`SKY_ERROR`](#skyerrors_8h_1a8405baf075a12e6232d75a8432d44f81)            | 
`define `[`LIBSKY_TESTING_H`](#skytest_8h_1aa31e87416545dcd6dcad132467018e22)            | 
`public GoUint32 `[`SKY_cli_CreateRawTxFromWallet`](#libskycoin_8h_1a20c77077115b9e629b9372dc45052978)`(Handle p0,`[`GoString`](#struct___go_string__)` p1,`[`GoString`](#struct___go_string__)` p2,`[`GoSlice`](#struct_go_slice)` p3,`[`Transaction`](#struct_transaction)` * p4)`            | 
`public GoUint32 `[`SKY_cli_CreateRawTxFromAddress`](#libskycoin_8h_1aa034c786a7dda49ba5caf939787dccd0)`(Handle p0,`[`GoString`](#struct___go_string__)` p1,`[`GoString`](#struct___go_string__)` p2,`[`GoString`](#struct___go_string__)` p3,`[`GoSlice`](#struct_go_slice)` p4,`[`Transaction`](#struct_transaction)` * p5)`            | 
`public void `[`SKY_cli_CreateRawTx`](#libskycoin_8h_1a8c5f5db1256b025fe16334a5bbbd1060)`(Handle p0,`[`Wallet`](#struct_wallet)` * p1,`[`GoSlice`](#struct_go_slice)` p2,`[`GoString`](#struct___go_string__)` p3,`[`GoSlice`](#struct_go_slice)` p4,`[`Transaction`](#struct_transaction)` * p5)`            | 
`public void `[`SKY_cli_NewTransaction`](#libskycoin_8h_1a86e20b22f34804b6ae95d334f0f2a51c)`(`[`GoSlice`](#struct_go_slice)` p0,`[`GoSlice`](#struct_go_slice)` p1,`[`GoSlice`](#struct_go_slice)` p2,`[`Transaction`](#struct_transaction)` * p3)`            | 
`public GoUint32 `[`SKY_cipher_DecodeBase58Address`](#libskycoin_8h_1ac624e50feca30ee4a1d215b55545bee8)`(`[`GoString`](#struct___go_string__)` p0,`[`Address`](#struct_address)` * p1)`            | 
`public void `[`SKY_cipher_AddressFromPubKey`](#libskycoin_8h_1a19502bcf26285130314d51c51034ed81)`(PubKey * p0,`[`Address`](#struct_address)` * p1)`            | 
`public void `[`SKY_cipher_AddressFromSecKey`](#libskycoin_8h_1a43dd635ce2999221eaf53443636a2cc9)`(SecKey * p0,`[`Address`](#struct_address)` * p1)`            | 
`public GoUint32 `[`SKY_cipher_BitcoinDecodeBase58Address`](#libskycoin_8h_1a3dbdee1e58738d24b1eee258971610ff)`(`[`GoString`](#struct___go_string__)` p0,`[`Address`](#struct_address)` * p1)`            | 
`public void `[`SKY_cipher_Address_Bytes`](#libskycoin_8h_1a2158e97d434d32452c4bc483fc42863e)`(`[`Address`](#struct_address)` * p0,`[`PubKeySlice`](#struct_go_slice__)` * p1)`            | 
`public void `[`SKY_cipher_Address_BitcoinBytes`](#libskycoin_8h_1ab6e15c43880cab6f2f7283f985c09b8c)`(`[`Address`](#struct_address)` * p0,`[`PubKeySlice`](#struct_go_slice__)` * p1)`            | 
`public GoUint32 `[`SKY_cipher_Address_Verify`](#libskycoin_8h_1a3c521e58dd6ba4d6c4996c1dd95f445b)`(`[`Address`](#struct_address)` * p0,PubKey * p1)`            | 
`public void `[`SKY_cipher_Address_String`](#libskycoin_8h_1a3fa26ea0e01795b94c8b163bf19677f6)`(`[`Address`](#struct_address)` * p0,`[`GoString_`](#struct_go_string__)` * p1)`            | 
`public void `[`SKY_cipher_Address_BitcoinString`](#libskycoin_8h_1a8f650e9df71fc4fec8a9326bd9ad209a)`(`[`Address`](#struct_address)` * p0,`[`GoString_`](#struct_go_string__)` * p1)`            | 
`public void `[`SKY_cipher_Address_Checksum`](#libskycoin_8h_1a5ccfd64d21d152219b9f6a92ec24099a)`(`[`Address`](#struct_address)` * p0,Checksum * p1)`            | 
`public void `[`SKY_cipher_Address_BitcoinChecksum`](#libskycoin_8h_1a54f6d9269d976d337431869a6025d1b6)`(`[`Address`](#struct_address)` * p0,Checksum * p1)`            | 
`public void `[`SKY_cipher_BitcoinAddressFromPubkey`](#libskycoin_8h_1a3bb30e8687b82b3d2b7a3ecb6ad99d51)`(PubKey * p0,`[`GoString_`](#struct_go_string__)` * p1)`            | 
`public void `[`SKY_cipher_BitcoinWalletImportFormatFromSeckey`](#libskycoin_8h_1a5086015efeb4450facc6e44d01f3c0bf)`(SecKey * p0,`[`GoString_`](#struct_go_string__)` * p1)`            | 
`public GoUint32 `[`SKY_cipher_BitcoinAddressFromBytes`](#libskycoin_8h_1afd7bbc548f9add0fb3846c5537c3e6bb)`(`[`GoSlice`](#struct_go_slice)` p0,`[`Address`](#struct_address)` * p1)`            | 
`public GoUint32 `[`SKY_cipher_SecKeyFromWalletImportFormat`](#libskycoin_8h_1aff9f7e90c09af0fbe68de0c2cb93445b)`(`[`GoString`](#struct___go_string__)` p0,SecKey * p1)`            | 
`public GoInt `[`SKY_cipher_PubKeySlice_Len`](#libskycoin_8h_1a4f8c95cf781be6721227ed34e3e31a82)`(`[`PubKeySlice`](#struct_go_slice__)` * p0)`            | 
`public GoUint8 `[`SKY_cipher_PubKeySlice_Less`](#libskycoin_8h_1aee5e7adb2fb6a981499b67ccd6950f07)`(`[`PubKeySlice`](#struct_go_slice__)` * p0,GoInt p1,GoInt p2)`            | 
`public void `[`SKY_cipher_PubKeySlice_Swap`](#libskycoin_8h_1adec780a8bb7a1e6b06e46725d438802f)`(`[`PubKeySlice`](#struct_go_slice__)` * p0,GoInt p1,GoInt p2)`            | 
`public void `[`SKY_cipher_RandByte`](#libskycoin_8h_1a443396fbe41b5ceca52ca9ce1178da87)`(GoInt p0,`[`PubKeySlice`](#struct_go_slice__)` * p1)`            | 
`public GoUint32 `[`SKY_cipher_NewPubKey`](#libskycoin_8h_1acce44b33fe66eb8e5a474ce5a52ce6ad)`(`[`GoSlice`](#struct_go_slice)` p0,PubKey * p1)`            | 
`public GoUint32 `[`SKY_cipher_PubKeyFromHex`](#libskycoin_8h_1a2941fcfb4b91f1a6831ed7a008f2f492)`(`[`GoString`](#struct___go_string__)` p0,PubKey * p1)`            | 
`public GoUint32 `[`SKY_cipher_PubKeyFromSecKey`](#libskycoin_8h_1afcd38ca547a7dae97d415bbec4d25199)`(SecKey * p0,PubKey * p1)`            | 
`public GoUint32 `[`SKY_cipher_PubKeyFromSig`](#libskycoin_8h_1a92d9b6d12c0f4eba7aa63336c83b920d)`(Sig * p0,SHA256 * p1,PubKey * p2)`            | 
`public GoUint32 `[`SKY_cipher_PubKey_Verify`](#libskycoin_8h_1aa758e062425798924e1eea3cddafb8ef)`(PubKey * p0)`            | 
`public void `[`SKY_cipher_PubKey_Hex`](#libskycoin_8h_1a80a69cd598ab67fc83afa6f8646f4358)`(PubKey * p0,`[`GoString_`](#struct_go_string__)` * p1)`            | 
`public void `[`SKY_cipher_PubKey_ToAddressHash`](#libskycoin_8h_1ad9c96b9e1d8915c6546d4070f2e76cab)`(PubKey * p0,Ripemd160 * p1)`            | 
`public GoUint32 `[`SKY_cipher_NewSecKey`](#libskycoin_8h_1a009cc1bf2f436a1a790db7708f17198c)`(`[`GoSlice`](#struct_go_slice)` p0,SecKey * p1)`            | 
`public GoUint32 `[`SKY_cipher_SecKeyFromHex`](#libskycoin_8h_1a780e839e1ae75fbac8ec674d54927d77)`(`[`GoString`](#struct___go_string__)` p0,SecKey * p1)`            | 
`public GoUint32 `[`SKY_cipher_SecKey_Verify`](#libskycoin_8h_1aa17089f7a830bbd75d11095f561cf39d)`(SecKey * p0)`            | 
`public void `[`SKY_cipher_SecKey_Hex`](#libskycoin_8h_1a645f5f92b38653939297b625d0f8dc21)`(SecKey * p0,`[`GoString_`](#struct_go_string__)` * p1)`            | 
`public void `[`SKY_cipher_ECDH`](#libskycoin_8h_1a23f26a93a05cc2fdb59e6feac5fe5140)`(PubKey * p0,SecKey * p1,`[`PubKeySlice`](#struct_go_slice__)` * p2)`            | 
`public GoUint32 `[`SKY_cipher_NewSig`](#libskycoin_8h_1ae72cdb33ffdd48382414c6125df5592c)`(`[`GoSlice`](#struct_go_slice)` p0,Sig * p1)`            | 
`public GoUint32 `[`SKY_cipher_SigFromHex`](#libskycoin_8h_1acd82a1de9be7d7291f79f7b32add9eec)`(`[`GoString`](#struct___go_string__)` p0,Sig * p1)`            | 
`public void `[`SKY_cipher_Sig_Hex`](#libskycoin_8h_1a0df60afe6e0a6b09a45959c81123cd84)`(Sig * p0,`[`GoString_`](#struct_go_string__)` * p1)`            | 
`public void `[`SKY_cipher_SignHash`](#libskycoin_8h_1a88c349cd3a7f14df8decd6c9c646c7f3)`(SHA256 * p0,SecKey * p1,Sig * p2)`            | 
`public GoUint32 `[`SKY_cipher_ChkSig`](#libskycoin_8h_1af512d40dd13c355ee7222f4bd2085f41)`(`[`Address`](#struct_address)` * p0,SHA256 * p1,Sig * p2)`            | 
`public GoUint32 `[`SKY_cipher_VerifySignedHash`](#libskycoin_8h_1a8a9388cb9ff151f481f4681658381728)`(Sig * p0,SHA256 * p1)`            | 
`public GoUint32 `[`SKY_cipher_VerifySignature`](#libskycoin_8h_1a24949394563c4a49c502549a96809b32)`(PubKey * p0,Sig * p1,SHA256 * p2)`            | 
`public void `[`SKY_cipher_GenerateKeyPair`](#libskycoin_8h_1ab601accdb915f5794554f80801a856d1)`(PubKey * p0,SecKey * p1)`            | 
`public void `[`SKY_cipher_GenerateDeterministicKeyPair`](#libskycoin_8h_1a2e579464b2cf0cdb94701c52b84ad979)`(`[`GoSlice`](#struct_go_slice)` p0,PubKey * p1,SecKey * p2)`            | 
`public void `[`SKY_cipher_DeterministicKeyPairIterator`](#libskycoin_8h_1a43ef1cf5a3d82f8093ab1a8eee043fcb)`(`[`GoSlice`](#struct_go_slice)` p0,`[`PubKeySlice`](#struct_go_slice__)` * p1,PubKey * p2,SecKey * p3)`            | 
`public void `[`SKY_cipher_GenerateDeterministicKeyPairs`](#libskycoin_8h_1a60367b51223730d7d77297309b230be4)`(`[`GoSlice`](#struct_go_slice)` p0,GoInt p1,`[`PubKeySlice`](#struct_go_slice__)` * p2)`            | 
`public void `[`SKY_cipher_GenerateDeterministicKeyPairsSeed`](#libskycoin_8h_1a038ca71fb6cff77b49e41ad69f00a516)`(`[`GoSlice`](#struct_go_slice)` p0,GoInt p1,`[`PubKeySlice`](#struct_go_slice__)` * p2,`[`PubKeySlice`](#struct_go_slice__)` * p3)`            | 
`public GoUint32 `[`SKY_cipher_TestSecKey`](#libskycoin_8h_1ace44f781f8ec684f8e9c59fb7fa3dbb3)`(SecKey * p0)`            | 
`public GoUint32 `[`SKY_cipher_TestSecKeyHash`](#libskycoin_8h_1a6cb852f76a408372d3b0aa5221bcdaed)`(SecKey * p0,SHA256 * p1)`            | 
`public GoUint32 `[`SKY_cipher_Ripemd160_Set`](#libskycoin_8h_1a3ec912a20b9e36b1c12c94f8f5be119e)`(Ripemd160 * p0,`[`GoSlice`](#struct_go_slice)` p1)`            | 
`public void `[`SKY_cipher_HashRipemd160`](#libskycoin_8h_1a54bd9a7ead7b661260abf98d2490153f)`(`[`GoSlice`](#struct_go_slice)` p0,Ripemd160 * p1)`            | 
`public GoUint32 `[`SKY_cipher_SHA256_Set`](#libskycoin_8h_1aba4a6f9df12f6384da9eeb93d7723f73)`(SHA256 * p0,`[`GoSlice`](#struct_go_slice)` p1)`            | 
`public void `[`SKY_cipher_SHA256_Hex`](#libskycoin_8h_1a5c8d433169079581776fe5a5eec1adb4)`(SHA256 * p0,`[`GoString_`](#struct_go_string__)` * p1)`            | 
`public void `[`SKY_cipher_SHA256_Xor`](#libskycoin_8h_1adf30e083c2811ee14eeef58178ceefc6)`(SHA256 * p0,SHA256 * p1,SHA256 * p2)`            | 
`public GoUint32 `[`SKY_cipher_SumSHA256`](#libskycoin_8h_1a7f0f0fa3b1610c7e97dc2950b35a75c4)`(`[`GoSlice`](#struct_go_slice)` p0,SHA256 * p1)`            | 
`public GoUint32 `[`SKY_cipher_SHA256FromHex`](#libskycoin_8h_1a046027f0eb544d6a0b8f1b78a7d189b4)`(`[`GoString`](#struct___go_string__)` p0,SHA256 * p1)`            | 
`public void `[`SKY_cipher_DoubleSHA256`](#libskycoin_8h_1ad3390997c0c9aec4c2cd678b6356d3ea)`(`[`GoSlice`](#struct_go_slice)` p0,SHA256 * p1)`            | 
`public void `[`SKY_cipher_AddSHA256`](#libskycoin_8h_1aa5977d4828735c3bceb1153f8bffc21b)`(SHA256 * p0,SHA256 * p1,SHA256 * p2)`            | 
`public void `[`SKY_cipher_Merkle`](#libskycoin_8h_1a71f945bbf46e4496051c98e390d11460)`(`[`GoSlice`](#struct_go_slice)` * p0,SHA256 * p1)`            | 
`public int `[`cr_user_Address_eq`](#skycriterion_8h_1a5c3dd4cd20db987c789c0a49ba098185)`(`[`Address`](#struct_address)` * addr1,`[`Address`](#struct_address)` * addr2)`            | 
`public char * `[`cr_user_Address_tostr`](#skycriterion_8h_1a232c966bd05993a3e9accc2670af2872)`(`[`Address`](#struct_address)` * addr1)`            | 
`public int `[`cr_user_Address_noteq`](#skycriterion_8h_1a0fc9801f223de4c4dee580e47060ee12)`(`[`Address`](#struct_address)` * addr1,`[`Address`](#struct_address)` * addr2)`            | 
`public int `[`cr_user_GoString_eq`](#skycriterion_8h_1afde184bfa3d42dadb560478bb384fd0e)`(`[`GoString`](#struct___go_string__)` * string1,`[`GoString`](#struct___go_string__)` * string2)`            | 
`public int `[`cr_user_GoString__eq`](#skycriterion_8h_1adc4957c85581c8021d1bc5e1fe68954e)`(`[`GoString_`](#struct_go_string__)` * string1,`[`GoString_`](#struct_go_string__)` * string2)`            | 
`public char * `[`cr_user_GoString_tostr`](#skycriterion_8h_1ac49e1ea1279ec23eb1b06fc4cff4346e)`(`[`GoString`](#struct___go_string__)` * string)`            | 
`public char * `[`cr_user_GoString__tostr`](#skycriterion_8h_1a8ba00c85c7eede2d955cfe016cb1023d)`(`[`GoString_`](#struct_go_string__)` * string)`            | 
`public int `[`cr_user_SecKey_eq`](#skycriterion_8h_1a93763ef6964d4cae39c79a7f46a2f42f)`(SecKey * seckey1,SecKey * seckey2)`            | 
`public char * `[`cr_user_SecKey_tostr`](#skycriterion_8h_1a2285d6f43b6c3d8903444d0983308ad3)`(SecKey * seckey1)`            | 
`public int `[`cr_user_Ripemd160_noteq`](#skycriterion_8h_1ac402ce38ac35394b35bed0b866266f51)`(Ripemd160 * rp1,Ripemd160 * rp2)`            | 
`public int `[`cr_user_Ripemd160_eq`](#skycriterion_8h_1ab77cfde0399d3d261908732b6fd0074e)`(Ripemd160 * rp1,Ripemd160 * rp2)`            | 
`public char * `[`cr_user_Ripemd160_tostr`](#skycriterion_8h_1ab2fe9026270b2ff7589cb44e03bd94c5)`(Ripemd160 * rp1)`            | 
`public int `[`cr_user_GoSlice_eq`](#skycriterion_8h_1a68e13a153f444839e3dbe06cc14e2348)`(`[`GoSlice`](#struct_go_slice)` * slice1,`[`GoSlice`](#struct_go_slice)` * slice2)`            | 
`public char * `[`cr_user_GoSlice_tostr`](#skycriterion_8h_1aa058100c8835ae72f2c609ad2ef1ba85)`(`[`GoSlice`](#struct_go_slice)` * slice1)`            | 
`public int `[`cr_user_GoSlice_noteq`](#skycriterion_8h_1a361dbb4ff75151c68df6d37368880b24)`(`[`GoSlice`](#struct_go_slice)` * slice1,`[`GoSlice`](#struct_go_slice)` * slice2)`            | 
`public int `[`cr_user_SHA256_noteq`](#skycriterion_8h_1ae300f329a27269668d76e11aeeae30e7)`(SHA256 * sh1,SHA256 * sh2)`            | 
`public int `[`cr_user_SHA256_eq`](#skycriterion_8h_1a233276e4a1ceab57a648c997c970bf92)`(SHA256 * sh1,SHA256 * sh2)`            | 
`public char * `[`cr_user_SHA256_tostr`](#skycriterion_8h_1a396e12d7842311b730774d39913d6ab9)`(SHA256 * sh1)`            | 
`public void `[`randBytes`](#skystring_8h_1abc646fb4e2f83b9ec86bacd6f8006907)`(`[`GoSlice`](#struct_go_slice)` * bytes,size_t n)`            | 
`public void `[`strnhex`](#skystring_8h_1aef6e4f140a965b05589db78792dc3c09)`(unsigned char * buf,char * str,int n)`            | 
`public void `[`strhex`](#skystring_8h_1a589986670c6a1cd947da79512078ff05)`(unsigned char * buf,char * str)`            | 
`public void * `[`registerMemCleanup`](#skytest_8h_1a3138ecc83c1c8906c84ef5e0d54cdfbb)`(void * p)`            | 
`public void `[`fprintbuff`](#skytest_8h_1a1ee45e153c115a9a735b3ccbf992e495)`(FILE * f,void * buff,size_t n)`            | 
`public void `[`toGoString`](#skytest_8h_1a1bad90cc197623fa8328f71809dda1a3)`(`[`GoString_`](#struct_go_string__)` * s,`[`GoString`](#struct___go_string__)` * r)`            | 
`struct `[`_GoString_`](#struct___go_string__) | 
`struct `[`Address`](#struct_address) | Addresses of SKY accounts
`struct `[`Entry`](#struct_entry) | [Wallet](#struct_wallet) entry.
`struct `[`GoInterface`](#struct_go_interface) | 
`struct `[`GoInterface_`](#struct_go_interface__) | Instances of Go interface types.
`struct `[`GoSlice`](#struct_go_slice) | 
`struct `[`GoSlice_`](#struct_go_slice__) | Instances of Go slices
`struct `[`GoString_`](#struct_go_string__) | Instances of Go `string` type.
`struct `[`SendAmount`](#struct_send_amount) | Structure used to specify amounts transferred in a transaction.
`struct `[`Transaction`](#struct_transaction) | Skycoin transaction.
`struct `[`TransactionOutput`](#struct_transaction_output) | Skycoin transaction output.
`struct `[`UxBalance`](#struct_ux_balance) | Intermediate representation of a UxOut for sorting and spend choosing.
`struct `[`Wallet`](#struct_wallet) | Internal representation of a Skycoin wallet.

## Members

#### `define `[`GO_CGO_EXPORT_PROLOGUE_H`](#libskycoin_8h_1ac91211782906f9494d827fe6e0b2e190) 

#### `define `[`GO_CGO_PROLOGUE_H`](#libskycoin_8h_1ad45a58cf8a40d22e35017cb53dd6055a) 

#### `define `[`SKY_OK`](#skyerrors_8h_1a5cd9ddcf04c6f149c283c805c7d296da) 

#### `define `[`SKY_ERROR`](#skyerrors_8h_1a8405baf075a12e6232d75a8432d44f81) 

#### `define `[`LIBSKY_TESTING_H`](#skytest_8h_1aa31e87416545dcd6dcad132467018e22) 

#### `public GoUint32 `[`SKY_cli_CreateRawTxFromWallet`](#libskycoin_8h_1a20c77077115b9e629b9372dc45052978)`(Handle p0,`[`GoString`](#struct___go_string__)` p1,`[`GoString`](#struct___go_string__)` p2,`[`GoSlice`](#struct_go_slice)` p3,`[`Transaction`](#struct_transaction)` * p4)` 

#### `public GoUint32 `[`SKY_cli_CreateRawTxFromAddress`](#libskycoin_8h_1aa034c786a7dda49ba5caf939787dccd0)`(Handle p0,`[`GoString`](#struct___go_string__)` p1,`[`GoString`](#struct___go_string__)` p2,`[`GoString`](#struct___go_string__)` p3,`[`GoSlice`](#struct_go_slice)` p4,`[`Transaction`](#struct_transaction)` * p5)` 

#### `public void `[`SKY_cli_CreateRawTx`](#libskycoin_8h_1a8c5f5db1256b025fe16334a5bbbd1060)`(Handle p0,`[`Wallet`](#struct_wallet)` * p1,`[`GoSlice`](#struct_go_slice)` p2,`[`GoString`](#struct___go_string__)` p3,`[`GoSlice`](#struct_go_slice)` p4,`[`Transaction`](#struct_transaction)` * p5)` 

#### `public void `[`SKY_cli_NewTransaction`](#libskycoin_8h_1a86e20b22f34804b6ae95d334f0f2a51c)`(`[`GoSlice`](#struct_go_slice)` p0,`[`GoSlice`](#struct_go_slice)` p1,`[`GoSlice`](#struct_go_slice)` p2,`[`Transaction`](#struct_transaction)` * p3)` 

#### `public GoUint32 `[`SKY_cipher_DecodeBase58Address`](#libskycoin_8h_1ac624e50feca30ee4a1d215b55545bee8)`(`[`GoString`](#struct___go_string__)` p0,`[`Address`](#struct_address)` * p1)` 

#### `public void `[`SKY_cipher_AddressFromPubKey`](#libskycoin_8h_1a19502bcf26285130314d51c51034ed81)`(PubKey * p0,`[`Address`](#struct_address)` * p1)` 

#### `public void `[`SKY_cipher_AddressFromSecKey`](#libskycoin_8h_1a43dd635ce2999221eaf53443636a2cc9)`(SecKey * p0,`[`Address`](#struct_address)` * p1)` 

#### `public GoUint32 `[`SKY_cipher_BitcoinDecodeBase58Address`](#libskycoin_8h_1a3dbdee1e58738d24b1eee258971610ff)`(`[`GoString`](#struct___go_string__)` p0,`[`Address`](#struct_address)` * p1)` 

#### `public void `[`SKY_cipher_Address_Bytes`](#libskycoin_8h_1a2158e97d434d32452c4bc483fc42863e)`(`[`Address`](#struct_address)` * p0,`[`PubKeySlice`](#struct_go_slice__)` * p1)` 

#### `public void `[`SKY_cipher_Address_BitcoinBytes`](#libskycoin_8h_1ab6e15c43880cab6f2f7283f985c09b8c)`(`[`Address`](#struct_address)` * p0,`[`PubKeySlice`](#struct_go_slice__)` * p1)` 

#### `public GoUint32 `[`SKY_cipher_Address_Verify`](#libskycoin_8h_1a3c521e58dd6ba4d6c4996c1dd95f445b)`(`[`Address`](#struct_address)` * p0,PubKey * p1)` 

#### `public void `[`SKY_cipher_Address_String`](#libskycoin_8h_1a3fa26ea0e01795b94c8b163bf19677f6)`(`[`Address`](#struct_address)` * p0,`[`GoString_`](#struct_go_string__)` * p1)` 

#### `public void `[`SKY_cipher_Address_BitcoinString`](#libskycoin_8h_1a8f650e9df71fc4fec8a9326bd9ad209a)`(`[`Address`](#struct_address)` * p0,`[`GoString_`](#struct_go_string__)` * p1)` 

#### `public void `[`SKY_cipher_Address_Checksum`](#libskycoin_8h_1a5ccfd64d21d152219b9f6a92ec24099a)`(`[`Address`](#struct_address)` * p0,Checksum * p1)` 

#### `public void `[`SKY_cipher_Address_BitcoinChecksum`](#libskycoin_8h_1a54f6d9269d976d337431869a6025d1b6)`(`[`Address`](#struct_address)` * p0,Checksum * p1)` 

#### `public void `[`SKY_cipher_BitcoinAddressFromPubkey`](#libskycoin_8h_1a3bb30e8687b82b3d2b7a3ecb6ad99d51)`(PubKey * p0,`[`GoString_`](#struct_go_string__)` * p1)` 

#### `public void `[`SKY_cipher_BitcoinWalletImportFormatFromSeckey`](#libskycoin_8h_1a5086015efeb4450facc6e44d01f3c0bf)`(SecKey * p0,`[`GoString_`](#struct_go_string__)` * p1)` 

#### `public GoUint32 `[`SKY_cipher_BitcoinAddressFromBytes`](#libskycoin_8h_1afd7bbc548f9add0fb3846c5537c3e6bb)`(`[`GoSlice`](#struct_go_slice)` p0,`[`Address`](#struct_address)` * p1)` 

#### `public GoUint32 `[`SKY_cipher_SecKeyFromWalletImportFormat`](#libskycoin_8h_1aff9f7e90c09af0fbe68de0c2cb93445b)`(`[`GoString`](#struct___go_string__)` p0,SecKey * p1)` 

#### `public GoInt `[`SKY_cipher_PubKeySlice_Len`](#libskycoin_8h_1a4f8c95cf781be6721227ed34e3e31a82)`(`[`PubKeySlice`](#struct_go_slice__)` * p0)` 

#### `public GoUint8 `[`SKY_cipher_PubKeySlice_Less`](#libskycoin_8h_1aee5e7adb2fb6a981499b67ccd6950f07)`(`[`PubKeySlice`](#struct_go_slice__)` * p0,GoInt p1,GoInt p2)` 

#### `public void `[`SKY_cipher_PubKeySlice_Swap`](#libskycoin_8h_1adec780a8bb7a1e6b06e46725d438802f)`(`[`PubKeySlice`](#struct_go_slice__)` * p0,GoInt p1,GoInt p2)` 

#### `public void `[`SKY_cipher_RandByte`](#libskycoin_8h_1a443396fbe41b5ceca52ca9ce1178da87)`(GoInt p0,`[`PubKeySlice`](#struct_go_slice__)` * p1)` 

#### `public GoUint32 `[`SKY_cipher_NewPubKey`](#libskycoin_8h_1acce44b33fe66eb8e5a474ce5a52ce6ad)`(`[`GoSlice`](#struct_go_slice)` p0,PubKey * p1)` 

#### `public GoUint32 `[`SKY_cipher_PubKeyFromHex`](#libskycoin_8h_1a2941fcfb4b91f1a6831ed7a008f2f492)`(`[`GoString`](#struct___go_string__)` p0,PubKey * p1)` 

#### `public GoUint32 `[`SKY_cipher_PubKeyFromSecKey`](#libskycoin_8h_1afcd38ca547a7dae97d415bbec4d25199)`(SecKey * p0,PubKey * p1)` 

#### `public GoUint32 `[`SKY_cipher_PubKeyFromSig`](#libskycoin_8h_1a92d9b6d12c0f4eba7aa63336c83b920d)`(Sig * p0,SHA256 * p1,PubKey * p2)` 

#### `public GoUint32 `[`SKY_cipher_PubKey_Verify`](#libskycoin_8h_1aa758e062425798924e1eea3cddafb8ef)`(PubKey * p0)` 

#### `public void `[`SKY_cipher_PubKey_Hex`](#libskycoin_8h_1a80a69cd598ab67fc83afa6f8646f4358)`(PubKey * p0,`[`GoString_`](#struct_go_string__)` * p1)` 

#### `public void `[`SKY_cipher_PubKey_ToAddressHash`](#libskycoin_8h_1ad9c96b9e1d8915c6546d4070f2e76cab)`(PubKey * p0,Ripemd160 * p1)` 

#### `public GoUint32 `[`SKY_cipher_NewSecKey`](#libskycoin_8h_1a009cc1bf2f436a1a790db7708f17198c)`(`[`GoSlice`](#struct_go_slice)` p0,SecKey * p1)` 

#### `public GoUint32 `[`SKY_cipher_SecKeyFromHex`](#libskycoin_8h_1a780e839e1ae75fbac8ec674d54927d77)`(`[`GoString`](#struct___go_string__)` p0,SecKey * p1)` 

#### `public GoUint32 `[`SKY_cipher_SecKey_Verify`](#libskycoin_8h_1aa17089f7a830bbd75d11095f561cf39d)`(SecKey * p0)` 

#### `public void `[`SKY_cipher_SecKey_Hex`](#libskycoin_8h_1a645f5f92b38653939297b625d0f8dc21)`(SecKey * p0,`[`GoString_`](#struct_go_string__)` * p1)` 

#### `public void `[`SKY_cipher_ECDH`](#libskycoin_8h_1a23f26a93a05cc2fdb59e6feac5fe5140)`(PubKey * p0,SecKey * p1,`[`PubKeySlice`](#struct_go_slice__)` * p2)` 

#### `public GoUint32 `[`SKY_cipher_NewSig`](#libskycoin_8h_1ae72cdb33ffdd48382414c6125df5592c)`(`[`GoSlice`](#struct_go_slice)` p0,Sig * p1)` 

#### `public GoUint32 `[`SKY_cipher_SigFromHex`](#libskycoin_8h_1acd82a1de9be7d7291f79f7b32add9eec)`(`[`GoString`](#struct___go_string__)` p0,Sig * p1)` 

#### `public void `[`SKY_cipher_Sig_Hex`](#libskycoin_8h_1a0df60afe6e0a6b09a45959c81123cd84)`(Sig * p0,`[`GoString_`](#struct_go_string__)` * p1)` 

#### `public void `[`SKY_cipher_SignHash`](#libskycoin_8h_1a88c349cd3a7f14df8decd6c9c646c7f3)`(SHA256 * p0,SecKey * p1,Sig * p2)` 

#### `public GoUint32 `[`SKY_cipher_ChkSig`](#libskycoin_8h_1af512d40dd13c355ee7222f4bd2085f41)`(`[`Address`](#struct_address)` * p0,SHA256 * p1,Sig * p2)` 

#### `public GoUint32 `[`SKY_cipher_VerifySignedHash`](#libskycoin_8h_1a8a9388cb9ff151f481f4681658381728)`(Sig * p0,SHA256 * p1)` 

#### `public GoUint32 `[`SKY_cipher_VerifySignature`](#libskycoin_8h_1a24949394563c4a49c502549a96809b32)`(PubKey * p0,Sig * p1,SHA256 * p2)` 

#### `public void `[`SKY_cipher_GenerateKeyPair`](#libskycoin_8h_1ab601accdb915f5794554f80801a856d1)`(PubKey * p0,SecKey * p1)` 

#### `public void `[`SKY_cipher_GenerateDeterministicKeyPair`](#libskycoin_8h_1a2e579464b2cf0cdb94701c52b84ad979)`(`[`GoSlice`](#struct_go_slice)` p0,PubKey * p1,SecKey * p2)` 

#### `public void `[`SKY_cipher_DeterministicKeyPairIterator`](#libskycoin_8h_1a43ef1cf5a3d82f8093ab1a8eee043fcb)`(`[`GoSlice`](#struct_go_slice)` p0,`[`PubKeySlice`](#struct_go_slice__)` * p1,PubKey * p2,SecKey * p3)` 

#### `public void `[`SKY_cipher_GenerateDeterministicKeyPairs`](#libskycoin_8h_1a60367b51223730d7d77297309b230be4)`(`[`GoSlice`](#struct_go_slice)` p0,GoInt p1,`[`PubKeySlice`](#struct_go_slice__)` * p2)` 

#### `public void `[`SKY_cipher_GenerateDeterministicKeyPairsSeed`](#libskycoin_8h_1a038ca71fb6cff77b49e41ad69f00a516)`(`[`GoSlice`](#struct_go_slice)` p0,GoInt p1,`[`PubKeySlice`](#struct_go_slice__)` * p2,`[`PubKeySlice`](#struct_go_slice__)` * p3)` 

#### `public GoUint32 `[`SKY_cipher_TestSecKey`](#libskycoin_8h_1ace44f781f8ec684f8e9c59fb7fa3dbb3)`(SecKey * p0)` 

#### `public GoUint32 `[`SKY_cipher_TestSecKeyHash`](#libskycoin_8h_1a6cb852f76a408372d3b0aa5221bcdaed)`(SecKey * p0,SHA256 * p1)` 

#### `public GoUint32 `[`SKY_cipher_Ripemd160_Set`](#libskycoin_8h_1a3ec912a20b9e36b1c12c94f8f5be119e)`(Ripemd160 * p0,`[`GoSlice`](#struct_go_slice)` p1)` 

#### `public void `[`SKY_cipher_HashRipemd160`](#libskycoin_8h_1a54bd9a7ead7b661260abf98d2490153f)`(`[`GoSlice`](#struct_go_slice)` p0,Ripemd160 * p1)` 

#### `public GoUint32 `[`SKY_cipher_SHA256_Set`](#libskycoin_8h_1aba4a6f9df12f6384da9eeb93d7723f73)`(SHA256 * p0,`[`GoSlice`](#struct_go_slice)` p1)` 

#### `public void `[`SKY_cipher_SHA256_Hex`](#libskycoin_8h_1a5c8d433169079581776fe5a5eec1adb4)`(SHA256 * p0,`[`GoString_`](#struct_go_string__)` * p1)` 

#### `public void `[`SKY_cipher_SHA256_Xor`](#libskycoin_8h_1adf30e083c2811ee14eeef58178ceefc6)`(SHA256 * p0,SHA256 * p1,SHA256 * p2)` 

#### `public GoUint32 `[`SKY_cipher_SumSHA256`](#libskycoin_8h_1a7f0f0fa3b1610c7e97dc2950b35a75c4)`(`[`GoSlice`](#struct_go_slice)` p0,SHA256 * p1)` 

#### `public GoUint32 `[`SKY_cipher_SHA256FromHex`](#libskycoin_8h_1a046027f0eb544d6a0b8f1b78a7d189b4)`(`[`GoString`](#struct___go_string__)` p0,SHA256 * p1)` 

#### `public void `[`SKY_cipher_DoubleSHA256`](#libskycoin_8h_1ad3390997c0c9aec4c2cd678b6356d3ea)`(`[`GoSlice`](#struct_go_slice)` p0,SHA256 * p1)` 

#### `public void `[`SKY_cipher_AddSHA256`](#libskycoin_8h_1aa5977d4828735c3bceb1153f8bffc21b)`(SHA256 * p0,SHA256 * p1,SHA256 * p2)` 

#### `public void `[`SKY_cipher_Merkle`](#libskycoin_8h_1a71f945bbf46e4496051c98e390d11460)`(`[`GoSlice`](#struct_go_slice)` * p0,SHA256 * p1)` 

#### `public int `[`cr_user_Address_eq`](#skycriterion_8h_1a5c3dd4cd20db987c789c0a49ba098185)`(`[`Address`](#struct_address)` * addr1,`[`Address`](#struct_address)` * addr2)` 

#### `public char * `[`cr_user_Address_tostr`](#skycriterion_8h_1a232c966bd05993a3e9accc2670af2872)`(`[`Address`](#struct_address)` * addr1)` 

#### `public int `[`cr_user_Address_noteq`](#skycriterion_8h_1a0fc9801f223de4c4dee580e47060ee12)`(`[`Address`](#struct_address)` * addr1,`[`Address`](#struct_address)` * addr2)` 

#### `public int `[`cr_user_GoString_eq`](#skycriterion_8h_1afde184bfa3d42dadb560478bb384fd0e)`(`[`GoString`](#struct___go_string__)` * string1,`[`GoString`](#struct___go_string__)` * string2)` 

#### `public int `[`cr_user_GoString__eq`](#skycriterion_8h_1adc4957c85581c8021d1bc5e1fe68954e)`(`[`GoString_`](#struct_go_string__)` * string1,`[`GoString_`](#struct_go_string__)` * string2)` 

#### `public char * `[`cr_user_GoString_tostr`](#skycriterion_8h_1ac49e1ea1279ec23eb1b06fc4cff4346e)`(`[`GoString`](#struct___go_string__)` * string)` 

#### `public char * `[`cr_user_GoString__tostr`](#skycriterion_8h_1a8ba00c85c7eede2d955cfe016cb1023d)`(`[`GoString_`](#struct_go_string__)` * string)` 

#### `public int `[`cr_user_SecKey_eq`](#skycriterion_8h_1a93763ef6964d4cae39c79a7f46a2f42f)`(SecKey * seckey1,SecKey * seckey2)` 

#### `public char * `[`cr_user_SecKey_tostr`](#skycriterion_8h_1a2285d6f43b6c3d8903444d0983308ad3)`(SecKey * seckey1)` 

#### `public int `[`cr_user_Ripemd160_noteq`](#skycriterion_8h_1ac402ce38ac35394b35bed0b866266f51)`(Ripemd160 * rp1,Ripemd160 * rp2)` 

#### `public int `[`cr_user_Ripemd160_eq`](#skycriterion_8h_1ab77cfde0399d3d261908732b6fd0074e)`(Ripemd160 * rp1,Ripemd160 * rp2)` 

#### `public char * `[`cr_user_Ripemd160_tostr`](#skycriterion_8h_1ab2fe9026270b2ff7589cb44e03bd94c5)`(Ripemd160 * rp1)` 

#### `public int `[`cr_user_GoSlice_eq`](#skycriterion_8h_1a68e13a153f444839e3dbe06cc14e2348)`(`[`GoSlice`](#struct_go_slice)` * slice1,`[`GoSlice`](#struct_go_slice)` * slice2)` 

#### `public char * `[`cr_user_GoSlice_tostr`](#skycriterion_8h_1aa058100c8835ae72f2c609ad2ef1ba85)`(`[`GoSlice`](#struct_go_slice)` * slice1)` 

#### `public int `[`cr_user_GoSlice_noteq`](#skycriterion_8h_1a361dbb4ff75151c68df6d37368880b24)`(`[`GoSlice`](#struct_go_slice)` * slice1,`[`GoSlice`](#struct_go_slice)` * slice2)` 

#### `public int `[`cr_user_SHA256_noteq`](#skycriterion_8h_1ae300f329a27269668d76e11aeeae30e7)`(SHA256 * sh1,SHA256 * sh2)` 

#### `public int `[`cr_user_SHA256_eq`](#skycriterion_8h_1a233276e4a1ceab57a648c997c970bf92)`(SHA256 * sh1,SHA256 * sh2)` 

#### `public char * `[`cr_user_SHA256_tostr`](#skycriterion_8h_1a396e12d7842311b730774d39913d6ab9)`(SHA256 * sh1)` 

#### `public void `[`randBytes`](#skystring_8h_1abc646fb4e2f83b9ec86bacd6f8006907)`(`[`GoSlice`](#struct_go_slice)` * bytes,size_t n)` 

#### `public void `[`strnhex`](#skystring_8h_1aef6e4f140a965b05589db78792dc3c09)`(unsigned char * buf,char * str,int n)` 

#### `public void `[`strhex`](#skystring_8h_1a589986670c6a1cd947da79512078ff05)`(unsigned char * buf,char * str)` 

#### `public void * `[`registerMemCleanup`](#skytest_8h_1a3138ecc83c1c8906c84ef5e0d54cdfbb)`(void * p)` 

#### `public void `[`fprintbuff`](#skytest_8h_1a1ee45e153c115a9a735b3ccbf992e495)`(FILE * f,void * buff,size_t n)` 

#### `public void `[`toGoString`](#skytest_8h_1a1bad90cc197623fa8328f71809dda1a3)`(`[`GoString_`](#struct_go_string__)` * s,`[`GoString`](#struct___go_string__)` * r)` 

# struct `_GoString_` 

## Summary

 Members                        | Descriptions                                
--------------------------------|---------------------------------------------
`public const char * `[`p`](#struct___go_string___1a6bc6b007533335efe02bafff799ec64c) | 
`public ptrdiff_t `[`n`](#struct___go_string___1a52d899ae12c13f4df8ff5ee014f3a106) | 

## Members

#### `public const char * `[`p`](#struct___go_string___1a6bc6b007533335efe02bafff799ec64c) 

#### `public ptrdiff_t `[`n`](#struct___go_string___1a52d899ae12c13f4df8ff5ee014f3a106) 

# struct `Address` 

Addresses of SKY accounts

## Summary

 Members                        | Descriptions                                
--------------------------------|---------------------------------------------
`public unsigned char `[`Version`](#struct_address_1a49fed92a3e4a3cc30678924a13acc19f) | [Address](#struct_address) version identifier. Used to differentiate testnet vs mainnet addresses, for instance.
`public Ripemd160 `[`Key`](#struct_address_1aa7fd9da55c53a8f7a6abe4987a8ea093) | [Address](#struct_address) hash identifier.

## Members

#### `public unsigned char `[`Version`](#struct_address_1a49fed92a3e4a3cc30678924a13acc19f) 

[Address](#struct_address) version identifier. Used to differentiate testnet vs mainnet addresses, for instance.

#### `public Ripemd160 `[`Key`](#struct_address_1aa7fd9da55c53a8f7a6abe4987a8ea093) 

[Address](#struct_address) hash identifier.

# struct `Entry` 

[Wallet](#struct_wallet) entry.

## Summary

 Members                        | Descriptions                                
--------------------------------|---------------------------------------------
`public `[`Address`](#struct_address)` `[`Address`](#struct_entry_1a0f1895b2f69382d62e5bdc459d30cae1) | [Wallet](#struct_wallet) address.
`public PubKey `[`Public`](#struct_entry_1ac4e11c22d8462b7633be585127f1a8d0) | Public key used to generate address.
`public SecKey `[`Secret`](#struct_entry_1a18296486015d3885702241c1d821d563) | Secret key used to generate address.

## Members

#### `public `[`Address`](#struct_address)` `[`Address`](#struct_entry_1a0f1895b2f69382d62e5bdc459d30cae1) 

[Wallet](#struct_wallet) address.

#### `public PubKey `[`Public`](#struct_entry_1ac4e11c22d8462b7633be585127f1a8d0) 

Public key used to generate address.

#### `public SecKey `[`Secret`](#struct_entry_1a18296486015d3885702241c1d821d563) 

Secret key used to generate address.

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

# struct `SendAmount` 

Structure used to specify amounts transferred in a transaction.

## Summary

 Members                        | Descriptions                                
--------------------------------|---------------------------------------------
`public `[`GoString_`](#struct_go_string__)` `[`Addr`](#struct_send_amount_1af646aec99cf83d30d17fd62e014d19f8) | Sender / receipient address.
`public GoInt64_ `[`Coins`](#struct_send_amount_1a7733f16af3115d3cfc712f2f687b73e4) | Amount transferred (e.g. measured in SKY)

## Members

#### `public `[`GoString_`](#struct_go_string__)` `[`Addr`](#struct_send_amount_1af646aec99cf83d30d17fd62e014d19f8) 

Sender / receipient address.

#### `public GoInt64_ `[`Coins`](#struct_send_amount_1a7733f16af3115d3cfc712f2f687b73e4) 

Amount transferred (e.g. measured in SKY)

# struct `Transaction` 

Skycoin transaction.

Instances of this struct are included in blocks.

## Summary

 Members                        | Descriptions                                
--------------------------------|---------------------------------------------
`public GoInt32_ `[`Length`](#struct_transaction_1a8a3c288e7b3f7245e0b7916ce322e5f9) | Current transaction's length expressed in bytes.
`public GoInt8_ `[`Type`](#struct_transaction_1a5bf4f40bde41c84f4ab5ff82bc74f744) | [Transaction](#struct_transaction)'s version. When a node tries to process a transaction, it must verify whether it supports the transaction's type. This is intended to provide a way to update skycoin clients and servers without crashing the network. If the transaction is not compatible with the node, it should not process it.
`public SHA256 `[`InnerHash`](#struct_transaction_1a62befdbe16b0b7f106cc8ac01d88b51f) | It's a SHA256 hash of the inputs and outputs of the transaction. It is used to protect against transaction mutability. This means that the transaction cannot be altered after its creation.
`public `[`GoSlice_`](#struct_go_slice__)` `[`Sigs`](#struct_transaction_1af0a2ba807a16f9ae66dd5682c243b943) | A list of digital signiatures generated by the skycoin client using the private key. It is used by Skycoin servers to verify the authenticy of the transaction. Each input requires a different signature.
`public `[`GoSlice_`](#struct_go_slice__)` `[`In`](#struct_transaction_1a317b73fcfa2b93fd16dfeab7ba228c39) | A list of references to unspent transaction outputs. Unlike other cryptocurrencies, such as Bitcoin, Skycoin unspent transaction outputs (UX) and Skycoin transactions (TX) are separated in the blockchain protocol, allowing for lighter transactions, thus reducing the broadcasting costs across the network.
`public `[`GoSlice_`](#struct_go_slice__)` `[`Out`](#struct_transaction_1a181edcb164c89b192b3838de7792cc89) | Outputs: A list of outputs created by the client, that will be recorded in the blockchain if transactions are confirmed. An output consists of a data structure representing an UTXT, which is composed by a Skycoin address to be sent to, the amount in Skycoin to be sent, and the amount of Coin Hours to be sent, and the SHA256 hash of the previous fields.

## Members

#### `public GoInt32_ `[`Length`](#struct_transaction_1a8a3c288e7b3f7245e0b7916ce322e5f9) 

Current transaction's length expressed in bytes.

#### `public GoInt8_ `[`Type`](#struct_transaction_1a5bf4f40bde41c84f4ab5ff82bc74f744) 

[Transaction](#struct_transaction)'s version. When a node tries to process a transaction, it must verify whether it supports the transaction's type. This is intended to provide a way to update skycoin clients and servers without crashing the network. If the transaction is not compatible with the node, it should not process it.

#### `public SHA256 `[`InnerHash`](#struct_transaction_1a62befdbe16b0b7f106cc8ac01d88b51f) 

It's a SHA256 hash of the inputs and outputs of the transaction. It is used to protect against transaction mutability. This means that the transaction cannot be altered after its creation.

#### `public `[`GoSlice_`](#struct_go_slice__)` `[`Sigs`](#struct_transaction_1af0a2ba807a16f9ae66dd5682c243b943) 

A list of digital signiatures generated by the skycoin client using the private key. It is used by Skycoin servers to verify the authenticy of the transaction. Each input requires a different signature.

#### `public `[`GoSlice_`](#struct_go_slice__)` `[`In`](#struct_transaction_1a317b73fcfa2b93fd16dfeab7ba228c39) 

A list of references to unspent transaction outputs. Unlike other cryptocurrencies, such as Bitcoin, Skycoin unspent transaction outputs (UX) and Skycoin transactions (TX) are separated in the blockchain protocol, allowing for lighter transactions, thus reducing the broadcasting costs across the network.

#### `public `[`GoSlice_`](#struct_go_slice__)` `[`Out`](#struct_transaction_1a181edcb164c89b192b3838de7792cc89) 

Outputs: A list of outputs created by the client, that will be recorded in the blockchain if transactions are confirmed. An output consists of a data structure representing an UTXT, which is composed by a Skycoin address to be sent to, the amount in Skycoin to be sent, and the amount of Coin Hours to be sent, and the SHA256 hash of the previous fields.

# struct `TransactionOutput` 

Skycoin transaction output.

Instances are integral part of transactions included in blocks.

## Summary

 Members                        | Descriptions                                
--------------------------------|---------------------------------------------
`public `[`Address`](#struct_address)` `[`Address`](#struct_transaction_output_1a0f1895b2f69382d62e5bdc459d30cae1) | Receipient address.
`public GoInt64_ `[`Coins`](#struct_transaction_output_1a7733f16af3115d3cfc712f2f687b73e4) | Amount sent to the receipient address.
`public GoInt64_ `[`Hours`](#struct_transaction_output_1a7aef551ad5991173b5a6160fd8fe1594) | Amount of Coin Hours sent to the receipient address.

## Members

#### `public `[`Address`](#struct_address)` `[`Address`](#struct_transaction_output_1a0f1895b2f69382d62e5bdc459d30cae1) 

Receipient address.

#### `public GoInt64_ `[`Coins`](#struct_transaction_output_1a7733f16af3115d3cfc712f2f687b73e4) 

Amount sent to the receipient address.

#### `public GoInt64_ `[`Hours`](#struct_transaction_output_1a7aef551ad5991173b5a6160fd8fe1594) 

Amount of Coin Hours sent to the receipient address.

# struct `UxBalance` 

Intermediate representation of a UxOut for sorting and spend choosing.

## Summary

 Members                        | Descriptions                                
--------------------------------|---------------------------------------------
`public SHA256 `[`Hash`](#struct_ux_balance_1a26607631022473d778367a1327b77a4c) | Hash of underlying UxOut.
`public GoInt64_ `[`BkSeq`](#struct_ux_balance_1a1c1b05acfa8b1a65809ba4525f14a55b) | moment balance calculation is performed at.
`public `[`Address`](#struct_address)` `[`Address`](#struct_ux_balance_1a0f1895b2f69382d62e5bdc459d30cae1) | Account holder address.
`public GoInt64_ `[`Coins`](#struct_ux_balance_1a7733f16af3115d3cfc712f2f687b73e4) | Coins amount (e.g. in SKY).
`public GoInt64_ `[`Hours`](#struct_ux_balance_1a7aef551ad5991173b5a6160fd8fe1594) | Balance of Coin Hours generated by underlying UxOut, depending on UxOut's head time.

## Members

#### `public SHA256 `[`Hash`](#struct_ux_balance_1a26607631022473d778367a1327b77a4c) 

Hash of underlying UxOut.

#### `public GoInt64_ `[`BkSeq`](#struct_ux_balance_1a1c1b05acfa8b1a65809ba4525f14a55b) 

moment balance calculation is performed at.

Block height corresponding to the

#### `public `[`Address`](#struct_address)` `[`Address`](#struct_ux_balance_1a0f1895b2f69382d62e5bdc459d30cae1) 

Account holder address.

#### `public GoInt64_ `[`Coins`](#struct_ux_balance_1a7733f16af3115d3cfc712f2f687b73e4) 

Coins amount (e.g. in SKY).

#### `public GoInt64_ `[`Hours`](#struct_ux_balance_1a7aef551ad5991173b5a6160fd8fe1594) 

Balance of Coin Hours generated by underlying UxOut, depending on UxOut's head time.

# struct `Wallet` 

Internal representation of a Skycoin wallet.

## Summary

 Members                        | Descriptions                                
--------------------------------|---------------------------------------------
`public GoMap_ `[`Meta`](#struct_wallet_1acdf30a4af55c2c677ebf6fc57b27e740) | Records items that are not deterministic, like filename, lable, wallet type, secrets, etc.
`public `[`GoSlice_`](#struct_go_slice__)` `[`Entries`](#struct_wallet_1a57b718d97f8db7e0bc9d9c755510951b) | Entries field stores the address entries that are deterministically generated from seed.

## Members

#### `public GoMap_ `[`Meta`](#struct_wallet_1acdf30a4af55c2c677ebf6fc57b27e740) 

Records items that are not deterministic, like filename, lable, wallet type, secrets, etc.

#### `public `[`GoSlice_`](#struct_go_slice__)` `[`Entries`](#struct_wallet_1a57b718d97f8db7e0bc9d9c755510951b) 

Entries field stores the address entries that are deterministically generated from seed.

Generated by [Moxygen](https://sourcey.com/moxygen)