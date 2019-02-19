# REST API Modifications


## Transaction APIs

### Get unconfirmed transactions

API sets: `READ`

```
URI: /api/v1/pendingTxs
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/pendingTxs
```

Result:

```json
[
    {
        "transaction": {
            "length": 317,
            "type": 0,
            "txid": "89578005d8730fe1789288ee7dea036160a9bd43234fb673baa6abd91289a48b",
            "inner_hash": "cac977eee019832245724aa643ceff451b9d8b24612b2f6a58177c79e8a4c26f",
            "sigs": [
                "3f084a0c750731dd985d3137200f9b5fc3de06069e62edea0cdd3a91d88e56b95aff5104a3e797ab4d6d417861af0c343efb0fff2e5ba9e7cf88ab714e10f38101",
                "e9a8aa8860d189daf0b1dbfd2a4cc309fc0c7250fa81113aa7258f9603d19727793c1b7533131605db64752aeb9c1f4465198bb1d8dd597213d6406a0a81ed3701"
            ],
            "inputs": [
                "bb89d4ed40d0e6e3a82c12e70b01a4bc240d2cd4f252cfac88235abe61bd3ad0",
                "170d6fd7be1d722a1969cb3f7d45cdf4d978129c3433915dbaf098d4f075bbfc"
            ],
            "outputs": [
                {
                    "uxid": "ec9cf2f6052bab24ec57847c72cfb377c06958a9e04a077d07b6dd5bf23ec106",
                    "dst": "nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq",
                    "coins": "60.000000",
                    "hours": 2458
                },
                {
                    "uxid": "be40210601829ba8653bac1d6ecc4049955d97fb490a48c310fd912280422bd9",
                    "dst": "2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc",
                    "coins": "1.000000",
                    "hours": 2458
                }
            ]
        },
        "received": "2017-05-09T10:11:57.14303834+02:00",
        "checked": "2017-05-09T10:19:58.801315452+02:00",
        "announced": "0001-01-01T00:00:00Z",
        "is_valid": true
    }
]
```

### Get unconfirmed transactions Verbose
```
URI: /api/v1/pendingTxs/verbose
Method: GET
```

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
The calculated hours are calculated based upon the current system time, and provide an approximate
coin hour value of the output if it were to be confirmed at that instant.

Example (verbose):

```sh
curl http://127.0.0.1:6420/api/v1/pendingTxs/verbose
```

Result:

```json
[
    {
        "transaction": {
            "length": 220,
            "type": 0,
            "txid": "d455564dcf1fb666c3846cf579ff33e21c203e2923938c6563fe7fcb8573ba44",
            "inner_hash": "4e73155db8ed04a3bd2b953218efcc9122ebfbf4c55f08f50d1563e48eacf71d",
            "fee": 12855964,
            "sigs": [
                "17330c256a50e2117ddccf51f1980fc14380f0f9476432196ade3043668759847b97e1b209961458745684d9239541f79d9ca9255582864d30a540017ab84f2b01"
            ],
            "inputs": [
                {
                    "uxid": "27e7bc48ceca4d47e806a87100a8a98592b7618702e1cd479bf4c190462a6d09",
                    "owner": "23MjQipM9YsPKkYiuaBmf6m7fD54wrzHxpd",
                    "coins": "7815.000000",
                    "hours": 279089,
                    "calculated_hours": 13101146
                }
            ],
            "outputs": [
                {
                    "uxid": "4b4ebf62acbaece798d0dfc92fcea85768a2874dad8a9b8eb5454288deae468c",
                    "dst": "23MjQipM9YsPKkYiuaBmf6m7fD54wrzHxpd",
                    "coins": "586.000000",
                    "hours": 122591
                },
                {
                    "uxid": "781cfb134d5fdad48f3c937dfcfc66b169a305adc8abdfe92a0ec94c564913f2",
                    "dst": "2ehrG4VKLRuvBNWYz3U7tS75QWvzyWR89Dg",
                    "coins": "7229.000000",
                    "hours": 122591
                }
            ]
        },
        "received": "2018-06-20T14:14:52.415702671+08:00",
        "checked": "2018-08-26T19:47:45.328131142+08:00",
        "announced": "2018-08-26T19:51:47.356083569+08:00",
        "is_valid": true
    }
]
```


### Get transactions for addresses

API sets: `READ`

```
URI: /api/v1/transactions
Method: GET, POST
Args:
    addrs: Comma seperated addresses [optional, returns all transactions if no address is provided]
    confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
```

The `"time"` field at the top level of each object in the response array indicates either the confirmed timestamp of a confirmed
transaction or the last received timestamp of an unconfirmed transaction.

The `POST` method can be used if many addresses need to be queried.

To get confirmed transactions for one or more addresses:

```sh
curl http://127.0.0.1:6420/api/v1/transactions?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY&confirmed=1
```

To get unconfirmed transactions for one or more addresses:

```sh
curl http://127.0.0.1:6420/api/v1/transactions?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY&confirmed=0
```

To get both confirmed and unconfirmed transactions for one or more addresses:

```sh
curl http://127.0.0.1:6420/api/v1/transactions?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY
```

Result:

```json
[
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 10492,
            "block_seq": 1177
        },
        "time": 1494275011,
        "txn": {
            "length": 317,
            "type": 0,
            "txid": "b09cd3a8baef6a449848f50a1b97943006ca92747d4e485d0647a3ea74550eca",
            "inner_hash": "2cb370051c92521a04ba5357e229d8ffa90d9d1741ea223b44dd60a1483ee0e5",
            "timestamp": 1494275011,
            "sigs": [
                "a55155ca15f73f0762f79c15917949a936658cff668647daf82a174eed95703a02622881f9cf6c7495536676f931b2d91d389a9e7b034232b3a1519c8da6fb8800",
                "cc7d7cbd6f31adabd9bde2c0deaa9277c0f3cf807a4ec97e11872817091dc3705841a6adb74acb625ee20ab6d3525350b8663566003276073d94c3bfe22fe48e01"
            ],
            "inputs": [
                "4f4b0078a9cd19b3395e54b3f42af6adc997f77f04e0ca54016c67c4f2384e3c",
                "36f4871646b6564b2f1ab72bd768a67579a1e0242bc68bcbcf1779bc75b3dddd"
            ],
            "outputs": [
                {
                    "uxid": "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191",
                    "dst": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "8.000000",
                    "hours": 7454
                },
                {
                    "uxid": "a1268e9bd2033b49b44afa765d20876467254f51e5515626780467267a65c563",
                    "dst": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
                    "coins": "1.000000",
                    "hours": 7454
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 10491,
            "block_seq": 1178
        },
        "time": 1494275231,
        "txn": {
            "length": 183,
            "type": 0,
            "txid": "a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3",
            "inner_hash": "075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e",
            "timestamp": 1494275231,
            "sigs": [
                "9b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00"
            ],
            "inputs": [
                "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191"
            ],
            "outputs": [
                {
                    "uxid": "70fa9dfb887f9ef55beb4e960f60e4703c56f98201acecf2cad729f5d7e84690",
                    "dst": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
                    "coins": "8.000000",
                    "hours": 931
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 8730,
            "block_seq": 2939
        },
        "time": 1505205561,
        "txn": {
            "length": 474,
            "type": 0,
            "txid": "b45e571988bc07bd0b623c999655fa878fb9bdd24c8cd24fde179bf4b26ae7b7",
            "inner_hash": "393804eca6afadc05db80cfb9e1024ef5761231c70705c406301bad33161f8bf",
            "timestamp": 1505205561,
            "sigs": [
                "fb9dd021cdff51ab56891cca0fd1600877f6e0691136dbe3f8324c3f4f7ee5bc624ded4954c1d70d8cb776ce3454d8f195bbb252e48b0f2cd388f5a733697d9301",
                "0639e61ba87a61f10b0e0114008ddd4e7090d9397370de28da27b7852b231b8e66c36d10fe3424c9b23a41266fd2c50f169233009713b332d6a48ce9c128ccef01",
                "055afe17222aab66c48c8e08e03a406bf2b8719f5221ec54c8e678078033bcd56b66bbc46a866f2be5e3f9ca454e3fbc2021630d0430b72e18c24d02df03c03100",
                "8cf56fb96e11d49bea728cb35ba5953fbc640817fac01b82e62a959ef8d4c3105298f2a6ea127bb07552abd905a667b58f6c79717e9f05258079de08d91f10a500"
            ],
            "inputs": [
                "dea9266aa7b687f4391e92f04436407c51a834274a5a33bc8bcf3189732e82e3",
                "e811bdce52ddac0d952d2546fdca8d1ac4e0ad32f170d3d73b724fb37c802652",
                "e94ccdbc07cc62fb41140b4daa7969438c749837c0808acf20dde113bdf1876b",
                "534afc496a7aee2ec55c71d85abfc27f35d16c56506f663b24d8ee4815583b6e"
            ],
            "outputs": [
                {
                    "uxid": "732e129fc1630aba3f06d833ce0a7a25f05dae5df3e7a135b5f82e99222e8c28",
                    "dst": "2hAjmdPP9R3um9JhKczeVdJUVugY6SPJBDm",
                    "coins": "6.000000",
                    "hours": 204
                }
            ]
        }
    }
]
```

### Get transactions for addresses Verbose

API sets: `READ`

```
URI: /api/v1/transactions/verbose
Method: GET, POST
Args:
    addrs: Comma seperated addresses [optional, returns all transactions if no address is provided]
    confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
```

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
If the transaction is confirmed, the calculated hours are the hours the transaction had in the block in which it was executed.
If the transaction is unconfirmed, the calculated hours are based upon the current system time, and are approximately
equal to the hours the output would have if it become confirmed immediately.

The `"time"` field at the top level of each object in the response array indicates either the confirmed timestamp of a confirmed
transaction or the last received timestamp of an unconfirmed transaction.

The `POST` method can be used if many addresses need to be queried.

Example (verbose):

```sh
curl http://127.0.0.1:6420/api/v1/transactions/verbose?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT
```

Result:

```json
[
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 53207,
            "block_seq": 1131
        },
        "time": 1494192581,
        "txn": {
            "timestamp": 1494192581,
            "length": 220,
            "type": 0,
            "txid": "b785dc57a9b53dbf0390213480dd9dffc32356fb79b82fa622a2607894bfab98",
            "inner_hash": "5279e944502d6bdaff25af7b7fb7c6e503c62ae70a01084031e1cb563afe2e2c",
            "fee": 317021,
            "sigs": [
                "f8cd208acc6674de79fa1192e5177325cda871c26707242dbd6fb9df245bf34b2fbc3dfe32e61eefa0543934556cb073bdeab6e555d7bfe6b7220f1ae575613d01"
            ],
            "inputs": [
                {
                    "uxid": "004d3ef83af64c542701b923ec5c727734de9d88837bcea37a2927a569dd3f0d",
                    "owner": "MbZvwdXHnMUZ1eUFxNDqxPEEHkkffKgq2F",
                    "coins": "904.000000",
                    "hours": 14,
                    "calculated_hours": 422693
                }
            ],
            "outputs": [
                {
                    "uxid": "4047c5cbbaf0ed927caa1391d5456d58e0857ef188f2eec8ee987a30b3f53aed",
                    "dst": "MbZvwdXHnMUZ1eUFxNDqxPEEHkkffKgq2F",
                    "coins": "903.000000",
                    "hours": 52836
                },
                {
                    "uxid": "4f4b0078a9cd19b3395e54b3f42af6adc997f77f04e0ca54016c67c4f2384e3c",
                    "dst": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "1.000000",
                    "hours": 52836
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 53206,
            "block_seq": 1132
        },
        "time": 1494192731,
        "txn": {
            "timestamp": 1494192731,
            "length": 220,
            "type": 0,
            "txid": "dc39c39bea82e5b56a1a77ce8485d9b06fda694e04ddf63af1273351c87dd077",
            "inner_hash": "b8f36a57212a68f4b3ecf9d699f286dafcdb624551e07c35a983832ffd37326c",
            "fee": 39628,
            "sigs": [
                "1005adda19efe31b5cd85caa85b4a42599263f649103fd26761f2261f3ee00460d9693c45406d782b0e04613aa412a5ef6b275c2a665a9f13167912da91777a700"
            ],
            "inputs": [
                {
                    "uxid": "4047c5cbbaf0ed927caa1391d5456d58e0857ef188f2eec8ee987a30b3f53aed",
                    "owner": "MbZvwdXHnMUZ1eUFxNDqxPEEHkkffKgq2F",
                    "coins": "903.000000",
                    "hours": 52836,
                    "calculated_hours": 52836
                }
            ],
            "outputs": [
                {
                    "uxid": "a6662ea872dabee2fae96a4561d67728d16cb3da372d4b7bbc74a18f2bc3fecf",
                    "dst": "MbZvwdXHnMUZ1eUFxNDqxPEEHkkffKgq2F",
                    "coins": "895.000000",
                    "hours": 6604
                },
                {
                    "uxid": "36f4871646b6564b2f1ab72bd768a67579a1e0242bc68bcbcf1779bc75b3dddd",
                    "dst": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "8.000000",
                    "hours": 6604
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 53161,
            "block_seq": 1177
        },
        "time": 1494275011,
        "txn": {
            "timestamp": 1494275011,
            "length": 317,
            "type": 0,
            "txid": "b09cd3a8baef6a449848f50a1b97943006ca92747d4e485d0647a3ea74550eca",
            "inner_hash": "2cb370051c92521a04ba5357e229d8ffa90d9d1741ea223b44dd60a1483ee0e5",
            "fee": 44726,
            "sigs": [
                "a55155ca15f73f0762f79c15917949a936658cff668647daf82a174eed95703a02622881f9cf6c7495536676f931b2d91d389a9e7b034232b3a1519c8da6fb8800",
                "cc7d7cbd6f31adabd9bde2c0deaa9277c0f3cf807a4ec97e11872817091dc3705841a6adb74acb625ee20ab6d3525350b8663566003276073d94c3bfe22fe48e01"
            ],
            "inputs": [
                {
                    "uxid": "4f4b0078a9cd19b3395e54b3f42af6adc997f77f04e0ca54016c67c4f2384e3c",
                    "owner": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "1.000000",
                    "hours": 52836,
                    "calculated_hours": 52857
                },
                {
                    "uxid": "36f4871646b6564b2f1ab72bd768a67579a1e0242bc68bcbcf1779bc75b3dddd",
                    "owner": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "8.000000",
                    "hours": 6604,
                    "calculated_hours": 6777
                }
            ],
            "outputs": [
                {
                    "uxid": "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191",
                    "dst": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "8.000000",
                    "hours": 7454
                },
                {
                    "uxid": "a1268e9bd2033b49b44afa765d20876467254f51e5515626780467267a65c563",
                    "dst": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
                    "coins": "1.000000",
                    "hours": 7454
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 53160,
            "block_seq": 1178
        },
        "time": 1494275231,
        "txn": {
            "timestamp": 1494275231,
            "length": 183,
            "type": 0,
            "txid": "a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3",
            "inner_hash": "075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e",
            "fee": 6523,
            "sigs": [
                "9b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00"
            ],
            "inputs": [
                {
                    "uxid": "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191",
                    "owner": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "8.000000",
                    "hours": 7454,
                    "calculated_hours": 7454
                }
            ],
            "outputs": [
                {
                    "uxid": "70fa9dfb887f9ef55beb4e960f60e4703c56f98201acecf2cad729f5d7e84690",
                    "dst": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
                    "coins": "8.000000",
                    "hours": 931
                }
            ]
        }
    }
]
```


### Get transaction info by id

API sets: `READ`

```
URI: /api/v1/transaction
Method: GET
Args:
    txid: transaction id [required]
```
If the transaction is confirmed, the calculated hours are the hours the transaction had in the block in which it was executed..
If the transaction is unconfirmed, the calculated hours are based upon the current system time, and are approximately
equal to the hours the output would have if it become confirmed immediately.

Example:

```sh
curl http://127.0.0.1:6420/api/v1/transaction?txid=a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3
```

Result:

```json
{
    "status": {
        "confirmed": true,
        "unconfirmed": false,
        "height": 1,
        "block_seq": 1178
    },
    "txn": {
        "length": 183,
        "type": 0,
        "txid": "a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3",
        "inner_hash": "075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e",
        "timestamp": 1494275231,
        "sigs": [
            "9b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00"
        ],
        "inputs": [
            "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191"
        ],
        "outputs": [
            {
                "uxid": "70fa9dfb887f9ef55beb4e960f60e4703c56f98201acecf2cad729f5d7e84690",
                "dst": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
                "coins": "8.000000",
                "hours": 931
            }
        ]
    }
}
```

### Get transaction info by id Verbose

API sets: `READ`

```
URI: /api/v1/transaction/verbose
Method: GET
Args:
    txid: transaction id [required]
```

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
If the transaction is confirmed, the calculated hours are the hours the transaction had in the block in which it was executed..
If the transaction is unconfirmed, the calculated hours are based upon the current system time, and are approximately
equal to the hours the output would have if it become confirmed immediately.

```sh
curl http://127.0.0.1:6420/api/v1/transaction?txid=a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3&verbose=1
```

Result:

```json
{
    "status": {
        "confirmed": true,
        "unconfirmed": false,
        "height": 53107,
        "block_seq": 1178
    },
    "time": 1494275231,
    "txn": {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 53107,
            "block_seq": 1178
        },
        "timestamp": 1494275231,
        "length": 183,
        "type": 0,
        "txid": "a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3",
        "inner_hash": "075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e",
        "fee": 6523,
        "sigs": [
            "9b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00"
        ],
        "inputs": [
            {
                "uxid": "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191",
                "owner": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                "coins": "8.000000",
                "hours": 7454,
                "calculated_hours": 7454
            }
        ],
        "outputs": [
            {
                "uxid": "70fa9dfb887f9ef55beb4e960f60e4703c56f98201acecf2cad729f5d7e84690",
                "dst": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
                "coins": "8.000000",
                "hours": 931
            }
        ]
    }
}
```

### Get transaction info by id Encoded

API sets: `READ`

```
URI: /api/v1/transaction/encoded
Method: GET
Args:
    txid: transaction id [required]
```

Just like `/api/v1/transaction`, only that transaction info is returned encoded.

```sh
curl http://127.0.0.1:6420/api/v1/transaction?txid=a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3&encoded=1
```

Result:

```json
{
    "status": {
        "confirmed": true,
        "unconfirmed": false,
        "height": 53267,
        "block_seq": 1178
    },
    "time": 1494275231,
    "encoded_transaction": "b700000000075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e010000009b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00010000005287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191010000000010722f061aa262381dce35193d43eceb112373c300127a0000000000a303000000000000"
}
```
