# Summary

 Members                        | Descriptions                                
--------------------------------|---------------------------------------------
`define `[`SKY_OK`](#skyerrors_8h_1a5cd9ddcf04c6f149c283c805c7d296da)            | 
`define `[`SKY_ERROR`](#skyerrors_8h_1a8405baf075a12e6232d75a8432d44f81)            | 
`define `[`LIBSKY_TESTING_H`](#skytest_8h_1aa31e87416545dcd6dcad132467018e22)            | 
`public int `[`cr_user_Address_eq`](#skycriterion_8h_1a5c3dd4cd20db987c789c0a49ba098185)`(`[`Address`](#struct_address)` * addr1,`[`Address`](#struct_address)` * addr2)`            | 
`public char * `[`cr_user_Address_tostr`](#skycriterion_8h_1a232c966bd05993a3e9accc2670af2872)`(`[`Address`](#struct_address)` * addr1)`            | 
`public int `[`cr_user_Address_noteq`](#skycriterion_8h_1a0fc9801f223de4c4dee580e47060ee12)`(`[`Address`](#struct_address)` * addr1,`[`Address`](#struct_address)` * addr2)`            | 
`public int `[`cr_user_GoString_eq`](#skycriterion_8h_1afde184bfa3d42dadb560478bb384fd0e)`(GoString * string1,GoString * string2)`            | 
`public int `[`cr_user_GoString__eq`](#skycriterion_8h_1adc4957c85581c8021d1bc5e1fe68954e)`(`[`GoString_`](#struct_go_string__)` * string1,`[`GoString_`](#struct_go_string__)` * string2)`            | 
`public char * `[`cr_user_GoString_tostr`](#skycriterion_8h_1ac49e1ea1279ec23eb1b06fc4cff4346e)`(GoString * string)`            | 
`public char * `[`cr_user_GoString__tostr`](#skycriterion_8h_1a8ba00c85c7eede2d955cfe016cb1023d)`(`[`GoString_`](#struct_go_string__)` * string)`            | 
`public int `[`cr_user_SecKey_eq`](#skycriterion_8h_1a93763ef6964d4cae39c79a7f46a2f42f)`(SecKey * seckey1,SecKey * seckey2)`            | 
`public char * `[`cr_user_SecKey_tostr`](#skycriterion_8h_1a2285d6f43b6c3d8903444d0983308ad3)`(SecKey * seckey1)`            | 
`public int `[`cr_user_Ripemd160_noteq`](#skycriterion_8h_1ac402ce38ac35394b35bed0b866266f51)`(Ripemd160 * rp1,Ripemd160 * rp2)`            | 
`public int `[`cr_user_Ripemd160_eq`](#skycriterion_8h_1ab77cfde0399d3d261908732b6fd0074e)`(Ripemd160 * rp1,Ripemd160 * rp2)`            | 
`public char * `[`cr_user_Ripemd160_tostr`](#skycriterion_8h_1ab2fe9026270b2ff7589cb44e03bd94c5)`(Ripemd160 * rp1)`            | 
`public int `[`cr_user_GoSlice_eq`](#skycriterion_8h_1a68e13a153f444839e3dbe06cc14e2348)`(GoSlice * slice1,GoSlice * slice2)`            | 
`public char * `[`cr_user_GoSlice_tostr`](#skycriterion_8h_1aa058100c8835ae72f2c609ad2ef1ba85)`(GoSlice * slice1)`            | 
`public int `[`cr_user_GoSlice_noteq`](#skycriterion_8h_1a361dbb4ff75151c68df6d37368880b24)`(GoSlice * slice1,GoSlice * slice2)`            | 
`public int `[`cr_user_SHA256_noteq`](#skycriterion_8h_1ae300f329a27269668d76e11aeeae30e7)`(SHA256 * sh1,SHA256 * sh2)`            | 
`public int `[`cr_user_SHA256_eq`](#skycriterion_8h_1a233276e4a1ceab57a648c997c970bf92)`(SHA256 * sh1,SHA256 * sh2)`            | 
`public char * `[`cr_user_SHA256_tostr`](#skycriterion_8h_1a396e12d7842311b730774d39913d6ab9)`(SHA256 * sh1)`            | 
`public void `[`randBytes`](#skystring_8h_1abc646fb4e2f83b9ec86bacd6f8006907)`(GoSlice * bytes,size_t n)`            | 
`public void `[`strnhex`](#skystring_8h_1aef6e4f140a965b05589db78792dc3c09)`(unsigned char * buf,char * str,int n)`            | 
`public void `[`strhex`](#skystring_8h_1a589986670c6a1cd947da79512078ff05)`(unsigned char * buf,char * str)`            | 
`public void * `[`registerMemCleanup`](#skytest_8h_1a3138ecc83c1c8906c84ef5e0d54cdfbb)`(void * p)`            | 
`public void `[`fprintbuff`](#skytest_8h_1a1ee45e153c115a9a735b3ccbf992e495)`(FILE * f,void * buff,size_t n)`            | 
`public void `[`toGoString`](#skytest_8h_1a1bad90cc197623fa8328f71809dda1a3)`(`[`GoString_`](#struct_go_string__)` * s,GoString * r)`            | 
`struct `[`Address`](#struct_address) | Addresses of SKY accounts
`struct `[`Entry`](#struct_entry) | [Wallet](#struct_wallet) entry.
`struct `[`GoInterface_`](#struct_go_interface__) | Instances of Go interface types.
`struct `[`GoSlice_`](#struct_go_slice__) | Instances of Go slices
`struct `[`GoString_`](#struct_go_string__) | Instances of Go `string` type.
`struct `[`SendAmount`](#struct_send_amount) | Structure used to specify amounts transferred in a transaction.
`struct `[`Transaction`](#struct_transaction) | Skycoin transaction.
`struct `[`TransactionOutput`](#struct_transaction_output) | Skycoin transaction output.
`struct `[`UxBalance`](#struct_ux_balance) | Intermediate representation of a UxOut for sorting and spend choosing.
`struct `[`Wallet`](#struct_wallet) | Internal representation of a Skycoin wallet.

## Members

#### `define `[`SKY_OK`](#skyerrors_8h_1a5cd9ddcf04c6f149c283c805c7d296da) 

#### `define `[`SKY_ERROR`](#skyerrors_8h_1a8405baf075a12e6232d75a8432d44f81) 

#### `define `[`LIBSKY_TESTING_H`](#skytest_8h_1aa31e87416545dcd6dcad132467018e22) 

#### `public int `[`cr_user_Address_eq`](#skycriterion_8h_1a5c3dd4cd20db987c789c0a49ba098185)`(`[`Address`](#struct_address)` * addr1,`[`Address`](#struct_address)` * addr2)` 

#### `public char * `[`cr_user_Address_tostr`](#skycriterion_8h_1a232c966bd05993a3e9accc2670af2872)`(`[`Address`](#struct_address)` * addr1)` 

#### `public int `[`cr_user_Address_noteq`](#skycriterion_8h_1a0fc9801f223de4c4dee580e47060ee12)`(`[`Address`](#struct_address)` * addr1,`[`Address`](#struct_address)` * addr2)` 

#### `public int `[`cr_user_GoString_eq`](#skycriterion_8h_1afde184bfa3d42dadb560478bb384fd0e)`(GoString * string1,GoString * string2)` 

#### `public int `[`cr_user_GoString__eq`](#skycriterion_8h_1adc4957c85581c8021d1bc5e1fe68954e)`(`[`GoString_`](#struct_go_string__)` * string1,`[`GoString_`](#struct_go_string__)` * string2)` 

#### `public char * `[`cr_user_GoString_tostr`](#skycriterion_8h_1ac49e1ea1279ec23eb1b06fc4cff4346e)`(GoString * string)` 

#### `public char * `[`cr_user_GoString__tostr`](#skycriterion_8h_1a8ba00c85c7eede2d955cfe016cb1023d)`(`[`GoString_`](#struct_go_string__)` * string)` 

#### `public int `[`cr_user_SecKey_eq`](#skycriterion_8h_1a93763ef6964d4cae39c79a7f46a2f42f)`(SecKey * seckey1,SecKey * seckey2)` 

#### `public char * `[`cr_user_SecKey_tostr`](#skycriterion_8h_1a2285d6f43b6c3d8903444d0983308ad3)`(SecKey * seckey1)` 

#### `public int `[`cr_user_Ripemd160_noteq`](#skycriterion_8h_1ac402ce38ac35394b35bed0b866266f51)`(Ripemd160 * rp1,Ripemd160 * rp2)` 

#### `public int `[`cr_user_Ripemd160_eq`](#skycriterion_8h_1ab77cfde0399d3d261908732b6fd0074e)`(Ripemd160 * rp1,Ripemd160 * rp2)` 

#### `public char * `[`cr_user_Ripemd160_tostr`](#skycriterion_8h_1ab2fe9026270b2ff7589cb44e03bd94c5)`(Ripemd160 * rp1)` 

#### `public int `[`cr_user_GoSlice_eq`](#skycriterion_8h_1a68e13a153f444839e3dbe06cc14e2348)`(GoSlice * slice1,GoSlice * slice2)` 

#### `public char * `[`cr_user_GoSlice_tostr`](#skycriterion_8h_1aa058100c8835ae72f2c609ad2ef1ba85)`(GoSlice * slice1)` 

#### `public int `[`cr_user_GoSlice_noteq`](#skycriterion_8h_1a361dbb4ff75151c68df6d37368880b24)`(GoSlice * slice1,GoSlice * slice2)` 

#### `public int `[`cr_user_SHA256_noteq`](#skycriterion_8h_1ae300f329a27269668d76e11aeeae30e7)`(SHA256 * sh1,SHA256 * sh2)` 

#### `public int `[`cr_user_SHA256_eq`](#skycriterion_8h_1a233276e4a1ceab57a648c997c970bf92)`(SHA256 * sh1,SHA256 * sh2)` 

#### `public char * `[`cr_user_SHA256_tostr`](#skycriterion_8h_1a396e12d7842311b730774d39913d6ab9)`(SHA256 * sh1)` 

#### `public void `[`randBytes`](#skystring_8h_1abc646fb4e2f83b9ec86bacd6f8006907)`(GoSlice * bytes,size_t n)` 

#### `public void `[`strnhex`](#skystring_8h_1aef6e4f140a965b05589db78792dc3c09)`(unsigned char * buf,char * str,int n)` 

#### `public void `[`strhex`](#skystring_8h_1a589986670c6a1cd947da79512078ff05)`(unsigned char * buf,char * str)` 

#### `public void * `[`registerMemCleanup`](#skytest_8h_1a3138ecc83c1c8906c84ef5e0d54cdfbb)`(void * p)` 

#### `public void `[`fprintbuff`](#skytest_8h_1a1ee45e153c115a9a735b3ccbf992e495)`(FILE * f,void * buff,size_t n)` 

#### `public void `[`toGoString`](#skytest_8h_1a1bad90cc197623fa8328f71809dda1a3)`(`[`GoString_`](#struct_go_string__)` * s,GoString * r)` 

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