# REST API Modifications

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
curl http://127.0.0.1:6420/api/v1/transaction/verbose?txid=a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3
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
curl http://127.0.0.1:6420/api/v1/transaction/encoded?txid=a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3
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


### Get unconfirmed transactions of a wallet

API sets: `WALLET`

```
URI: /api/v1/wallet/transactions
Method: GET
Args:
    id: Wallet ID [required]
```

Returns all unconfirmed transactions for all addresses in a given wallet

Example:

```sh
curl http://127.0.0.1:6420/api/v1/wallet/transactions?id=2017_11_25_e5fb.wlt
```

Result:

```json
{
    "transactions": [
        {
            "transaction": {
                "length": 317,
                "type": 0,
                "txid": "76ecbabc53ea2a3be46983058433dda6a3cf7ea0b86ba14d90b932fa97385de7",
                "inner_hash": "5d55837bb0cbda9c9323ff9aafd7c3d31d0d38638346172fbe2d9078ebaa892a",
                "sigs": [
                    "464b7724302178c1cfeacadaaf3556a3b7e5259adf51919476c3acc695747ed244b5ce2187ce7bedb6ad65c71f7f7ff3fa6805e64fe5da3aaa00ad563c7424f600",
                    "1155537b0391d4a6ee5eac07dee5798e953dca3a7c30643403dd2d326582c7d35080a16dc22644782ce1087bfc3bd06c2bf68e9a98e3989d90831646a9be2c9101"
                ],
                "inputs": [
                    "782a8662efb0e933cab7d3ae9429ab53c4208cf44d8cdc07c2fbd7204b6b5cad",
                    "2f6b61a44086588c4eaa56a5dd9f1e0be2528861a6731608fcec38891b95db91"
                ],
                "outputs": [
                    {
                        "uxid": "bd302ef776efa8548183b89f21e90649f21b90fe2d2e90ecc1b880f2d995f226",
                        "dst": "2UXZTg4ZHF6715b6tRhtaqceuQQ3G79GiZg",
                        "coins": "998.000000",
                        "hours": 247538
                    },
                    {
                        "uxid": "31058b6bfb30bfd441aec00929e75782bce47c8a75787ba519dbb268f89d2c4b",
                        "dst": "2awsJ2CR5H6QXCF2hwDjcvcAH9SgyfxCxgz",
                        "coins": "1.000000",
                        "hours": 247538
                    }
                ]
            },
            "received": "2018-03-16T18:03:57.139109904+05:30",
            "checked": "2018-03-16T18:03:57.139109904+05:30",
            "announced": "0001-01-01T00:00:00Z",
            "is_valid": true
        }
    ]
}
```

### Get unconfirmed transactions of a wallet Verbose

API sets: `WALLET`

```
URI: /api/v1/wallet/transactions/verbose
Method: GET
Args:
    id: Wallet ID [required]
```

Returns all unconfirmed transactions for all addresses in a given wallet

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
The calculated hours are based upon the current system time, and are approximately
equal to the hours the output would have if it become confirmed immediately.

Example:

```sh
curl http://127.0.0.1:6420/api/v1/wallet/transactions/verbose?id=2017_11_25_e5fb.wlt
```

Result:

```json
{
    "transactions": [
        {
            "transaction": {
                "length": 317,
                "type": 0,
                "txid": "76ecbabc53ea2a3be46983058433dda6a3cf7ea0b86ba14d90b932fa97385de7",
                "inner_hash": "5d55837bb0cbda9c9323ff9aafd7c3d31d0d38638346172fbe2d9078ebaa892a",
                "sigs": [
                    "464b7724302178c1cfeacadaaf3556a3b7e5259adf51919476c3acc695747ed244b5ce2187ce7bedb6ad65c71f7f7ff3fa6805e64fe5da3aaa00ad563c7424f600",
                    "1155537b0391d4a6ee5eac07dee5798e953dca3a7c30643403dd2d326582c7d35080a16dc22644782ce1087bfc3bd06c2bf68e9a98e3989d90831646a9be2c9101"
                ],
                "fee": 495076,
                "inputs": [
                    {
                        "uxid": "782a8662efb0e933cab7d3ae9429ab53c4208cf44d8cdc07c2fbd7204b6b5cad",
                        "owner": "8C5icxR9zdkYTZZTVV3cCX7QoK4EkLuK4p",
                        "coins": "997.000000",
                        "hours": 880000,
                        "calculated_hours": 990000
                    },
                    {
                        "uxid": "2f6b61a44086588c4eaa56a5dd9f1e0be2528861a6731608fcec38891b95db91",
                        "owner": "23A1EWMZopUFLCwtXMe2CU9xTCbi5Gth643",
                        "coins": "2.000000",
                        "hours": 10,
                        "calculated_hours": 152
                    }
                ],
                "outputs": [
                    {
                        "uxid": "bd302ef776efa8548183b89f21e90649f21b90fe2d2e90ecc1b880f2d995f226",
                        "dst": "2UXZTg4ZHF6715b6tRhtaqceuQQ3G79GiZg",
                        "coins": "998.000000",
                        "hours": 247538
                    },
                    {
                        "uxid": "31058b6bfb30bfd441aec00929e75782bce47c8a75787ba519dbb268f89d2c4b",
                        "dst": "2awsJ2CR5H6QXCF2hwDjcvcAH9SgyfxCxgz",
                        "coins": "1.000000",
                        "hours": 247538
                    }
                ]
            },
            "received": "2018-03-16T18:03:57.139109904+05:30",
            "checked": "2018-03-16T18:03:57.139109904+05:30",
            "announced": "0001-01-01T00:00:00Z",
            "is_valid": true
        }
    ]
}
```


### Get blocks in specific range

API sets: `READ`

```
URI: /api/v1/blocks
Method: GET
Args:
    start: start seq
    end: end seq
    seqs: comma-separated list of block seqs
```

This endpoint has two modes: range and seqs.
The `seqs` parameter cannot be combined with `start`, `end`.

If `start` and/or `end` are provided, returns blocks in the range [`start`, `end`].
Both start and end sequences are included in the returned array of blocks.

If `seqs` is provided, returns blocks matching the specified sequences.
`seqs` must not contain any duplicate values.
If a block does not exist for any of the given sequence numbers, a `404` error is returned.

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
The calculated hours are the hours the transaction had in the block in which it was executed.

Example:

```sh
curl http://127.0.0.1:6420/api/v1/blocks?start=101&end=102
```

Result:

```json
{
    "blocks": [
        {
            "header": {
                "seq": 101,
                "block_hash": "8156057fc823589288f66c91edb60c11ff004465bcbe3a402b1328be7f0d6ce0",
                "previous_block_hash": "725e76907998485d367a847b0fb49f08536c592247762279fcdbd9907fee5607",
                "timestamp": 1429274666,
                "fee": 720335,
                "version": 0,
                "tx_body_hash": "e8fe5290afba3933389fd5860dca2cbcc81821028be9c65d0bb7cf4e8d2c4c18",
                "ux_hash": "348989599d30d3adfaaea98577963caa419ab0276279296e7d194a9cbb8cad04"
            },
            "body": {
                "txns": [
                    {
                        "length": 183,
                        "type": 0,
                        "txid": "e8fe5290afba3933389fd5860dca2cbcc81821028be9c65d0bb7cf4e8d2c4c18",
                        "inner_hash": "45da31b68748eafdb08ef8bf1ebd1c07c0f14fcb0d66759d6cf4642adc956d06",
                        "sigs": [
                            "09bce2c888ceceeb19999005cceb1efdee254cacb60edee118b51ffd740ff6503a8f9cbd60a16c7581bfd64f7529b649d0ecc8adbe913686da97fe8c6543189001"
                        ],
                        "inputs": [
                            "6002f3afc7054c0e1161bcf2b4c1d4d1009440751bc1fe806e0eae33291399f4"
                        ],
                        "outputs": [
                            {
                                "uxid": "f9bffdcbe252acb1c3a8a1e8c99829342ba1963860d5692eebaeb9bcfbcaf274",
                                "dst": "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
                                "coins": "27000.000000",
                                "hours": 102905
                            }
                        ]
                    }
                ]
            },
            "size": 183
        },
        {
            "header": {
                "seq": 102,
                "block_hash": "311f4b83b4fdb9fd1d45648115969cf4b3aab2d1acad9e2aa735829245c525f3",
                "previous_block_hash": "8156057fc823589288f66c91edb60c11ff004465bcbe3a402b1328be7f0d6ce0",
                "timestamp": 1429274686,
                "fee": 710046,
                "version": 0,
                "tx_body_hash": "7b13cab45b52dd2df291ec97cf000bf6ea1b647d6fdf0261a7527578d8b71b9d",
                "ux_hash": "f7512b0718f392c7503f86e69175efd7835ea4c3dd3f71ff65c7ad8873a6a9e8"
            },
            "body": {
                "txns": [
                    {
                        "length": 183,
                        "type": 0,
                        "txid": "7b13cab45b52dd2df291ec97cf000bf6ea1b647d6fdf0261a7527578d8b71b9d",
                        "inner_hash": "73bfee3a7c8d4f8a68657ebcaf69a59639f762bfc1a6f4468f3ca4724bc5b9f8",
                        "sigs": [
                            "c4bcada17604a4a62baf50f929655027f2913639c27b773871f2135b72553c1959737e39d50e8349ffa5a7679de845aa6370999dbaaff4c7f9fd01260818683901"
                        ],
                        "inputs": [
                            "4e75b4bced3404590d38ca06440c275d7fd86618a84966a0a1053fb18164e898"
                        ],
                        "outputs": [
                            {
                                "uxid": "0a5603a1a5aeda575aa498cdaec5a4c893a28669dba84163eba2e90db3d9f39d",
                                "dst": "2JJ8pgq8EDAnrzf9xxBJapE2qkYLefW4uF8",
                                "coins": "26700.000000",
                                "hours": 101435
                            }
                        ]
                    }
                ]
            },
            "size": 183
        }
    ]
}
```

Example (seqs):

```sh
curl http://127.0.0.1:6420/api/v1/blocks?seqs=3,5,7
```

```json
{
    "blocks": [
        {
            "header": {
                "seq": 3,
                "block_hash": "35c3ebbe6feaeeab27ac77c1712051787bdd4bbfb5cdcdebc81f8aac98a2f3f3",
                "previous_block_hash": "01723bc4dc90f1cb857a94fe5e3bb50c02e6689fd998f8147c9cae07fbfa63af",
                "timestamp": 1427927671,
                "fee": 0,
                "version": 0,
                "tx_body_hash": "a6a709e9388a4d67a47d262b11da5f804eddd9d67acc4a3e450f7a567bdc1619"
            },
            "body": {
                "txns": [
                    {
                        "length": 183,
                        "type": 0,
                        "txid": "a6a709e9388a4d67a47d262b11da5f804eddd9d67acc4a3e450f7a567bdc1619",
                        "inner_hash": "ea6adee3180c7f9d73d1e693822d5d1c2bba85067f89a873355bc771a078faa1",
                        "sigs": [
                            "ce8fd47e2044ed17998f92621e90329f673a746c802d67f639ca083705dd199f6ee346781497b44132434922879244d819694b5903093f784570c55d293ab4af01"
                        ],
                        "inputs": [
                            "af0b2c1cc882a56b6c0c06e99e7d2731413b988329a2c47a5c2aa8be589b707a"
                        ],
                        "outputs": [
                            {
                                "uxid": "9eb7954461ba0256c9054fe38c00c66e60428dccf900a62e74b9fe39310aea13",
                                "dst": "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
                                "coins": "10.000000",
                                "hours": 0
                            }
                        ]
                    }
                ]
            },
            "size": 183
        },
        {
            "header": {
                "seq": 5,
                "block_hash": "114fe60587a158428a47e0f9571d764f495912c299aa4e67fc88004cf21b0c24",
                "previous_block_hash": "415e47348a1e642cb2e31d00ee500747d3aed0336aabfff7d783ed21465251c7",
                "timestamp": 1428798821,
                "fee": 2036,
                "version": 0,
                "tx_body_hash": "0579e7727627cd9815a8a8b5e1df86124f45a4132cc0dbd00d2f110e4f409b69"
            },
            "body": {
                "txns": [
                    {
                        "length": 317,
                        "type": 0,
                        "txid": "0579e7727627cd9815a8a8b5e1df86124f45a4132cc0dbd00d2f110e4f409b69",
                        "inner_hash": "fe123ca954a82bb1ce2cc9ef9c56d6b649a4cbaf5b17394b0ffda651ed32327e",
                        "sigs": [
                            "056ed0f74367fb1370d7e98689953983d9cf34eb6669854f1645c8a16c93d85075661e7d4f6df0ce5ca8eb9852eff6a12fbac2caafee03bb8c616f847c61416800",
                            "8aaa7f320a7b01169d3217a600100cb27c55e4ce56cd3455814f56d8e4e65be746e0e20e776087af6f19361f0b898edc2123a5f9bd35d24ef8b8669ca85b142601"
                        ],
                        "inputs": [
                            "9eb7954461ba0256c9054fe38c00c66e60428dccf900a62e74b9fe39310aea13",
                            "706f82c481906108880d79372ab5c126d32ecc98cf3f7c74cf33f5fda49dcf70"
                        ],
                        "outputs": [
                            {
                                "uxid": "fa2b598d233fe434f907f858d5de812eacf50c7b3fd152c77cd6e246fe356a9e",
                                "dst": "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
                                "coins": "999890.000000",
                                "hours": 4073
                            },
                            {
                                "uxid": "dc63c680f408c4e646037966189383a5d50eda34e666c2a0c75c0c6bf13b71a1",
                                "dst": "2fGC7kwAM9yZyEF1QqBqp8uo9RUsF6ENGJF",
                                "coins": "100.000000",
                                "hours": 0
                            }
                        ]
                    }
                ]
            },
            "size": 317
        },
        {
            "header": {
                "seq": 7,
                "block_hash": "6cb71b57c998a5367101e01d48c097eccd4f5abf311c89bcca8ee213581f355f",
                "previous_block_hash": "103949030e90fcebc5d8ca1c9c59f30a31aa71911401d22a2422e4571b035701",
                "timestamp": 1428807671,
                "fee": 0,
                "version": 0,
                "tx_body_hash": "f832428481690fa918d6d29946e191f2c8c89b2388a906e0c53dceee6070a24b"
            },
            "body": {
                "txns": [
                    {
                        "length": 220,
                        "type": 0,
                        "txid": "f832428481690fa918d6d29946e191f2c8c89b2388a906e0c53dceee6070a24b",
                        "inner_hash": "f440c514779522a6387edda9b9d9835f00680fb314546efb7bc9762a17884156",
                        "sigs": [
                            "8fe96f5502270e4efa962b2aef2b81795fe26a8f0c9a494e2ae9c7e624af455c49396270ae7a25b41d439fd56dea9d556a135129122de1b1274b1e2a5d75f2ea01"
                        ],
                        "inputs": [
                            "8ff8a647e4542fab01e078ac467b2c9f2e5f7de55d77ec2711f8abc718e2c91b"
                        ],
                        "outputs": [
                            {
                                "uxid": "17090c40091d009d6a684043d3be2e9cb1dc60a664a9c2e388af1f3a7345724b",
                                "dst": "2fGC7kwAM9yZyEF1QqBqp8uo9RUsF6ENGJF",
                                "coins": "90.000000",
                                "hours": 0
                            },
                            {
                                "uxid": "f9e7a412cdff80e95ddbe1d76fcc73f967cb99d383b0659e1355c8e623f02b62",
                                "dst": "WADSeEwEQVbtUy8CfcVimyxX1KjTRkvfoK",
                                "coins": "5.000000",
                                "hours": 0
                            }
                        ]
                    }
                ]
            },
            "size": 220
        }
    ]
}
```

### Get blocks in specific range Verbose

API sets: `READ`

```
URI: /api/v1/blocks/verbose
Method: GET
Args:
    start: start seq
    end: end seq
    seqs: comma-separated list of block seqs
```

This endpoint has two modes: range and seqs.
The `seqs` parameter cannot be combined with `start`, `end`.

If `start` and/or `end` are provided, returns blocks in the range [`start`, `end`].
Both start and end sequences are included in the returned array of blocks.

If `seqs` is provided, returns blocks matching the specified sequences.
`seqs` must not contain any duplicate values.
If a block does not exist for any of the given sequence numbers, a `404` error is returned.

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
The calculated hours are the hours the transaction had in the block in which it was executed.

Example:

```sh
curl http://127.0.0.1:6420/api/v1/blocks/verbose?start=101&end=102
```

Result:

```json
{
    "blocks": [
        {
            "header": {
                "seq": 101,
                "block_hash": "8156057fc823589288f66c91edb60c11ff004465bcbe3a402b1328be7f0d6ce0",
                "previous_block_hash": "725e76907998485d367a847b0fb49f08536c592247762279fcdbd9907fee5607",
                "timestamp": 1429274666,
                "fee": 720335,
                "version": 0,
                "tx_body_hash": "e8fe5290afba3933389fd5860dca2cbcc81821028be9c65d0bb7cf4e8d2c4c18",
                "ux_hash": "348989599d30d3adfaaea98577963caa419ab0276279296e7d194a9cbb8cad04"
            },
            "body": {
                "txns": [
                    {
                        "length": 183,
                        "type": 0,
                        "txid": "e8fe5290afba3933389fd5860dca2cbcc81821028be9c65d0bb7cf4e8d2c4c18",
                        "inner_hash": "45da31b68748eafdb08ef8bf1ebd1c07c0f14fcb0d66759d6cf4642adc956d06",
                        "fee": 720335,
                        "sigs": [
                            "09bce2c888ceceeb19999005cceb1efdee254cacb60edee118b51ffd740ff6503a8f9cbd60a16c7581bfd64f7529b649d0ecc8adbe913686da97fe8c6543189001"
                        ],
                        "inputs": [
                            {
                                "uxid": "6002f3afc7054c0e1161bcf2b4c1d4d1009440751bc1fe806e0eae33291399f4",
                                "owner": "2M1C5LSZ4Pvu5RWS44bCdY6or3R8grQw7ez",
                                "coins": "27000.000000",
                                "hours": 220,
                                "calculated_hours": 823240
                            }
                        ],
                        "outputs": [
                            {
                                "uxid": "f9bffdcbe252acb1c3a8a1e8c99829342ba1963860d5692eebaeb9bcfbcaf274",
                                "dst": "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
                                "coins": "27000.000000",
                                "hours": 102905
                            }
                        ]
                    }
                ]
            },
            "size": 183
        },
        {
            "header": {
                "seq": 102,
                "block_hash": "311f4b83b4fdb9fd1d45648115969cf4b3aab2d1acad9e2aa735829245c525f3",
                "previous_block_hash": "8156057fc823589288f66c91edb60c11ff004465bcbe3a402b1328be7f0d6ce0",
                "timestamp": 1429274686,
                "fee": 710046,
                "version": 0,
                "tx_body_hash": "7b13cab45b52dd2df291ec97cf000bf6ea1b647d6fdf0261a7527578d8b71b9d",
                "ux_hash": "f7512b0718f392c7503f86e69175efd7835ea4c3dd3f71ff65c7ad8873a6a9e8"
            },
            "body": {
                "txns": [
                    {
                        "length": 183,
                        "type": 0,
                        "txid": "7b13cab45b52dd2df291ec97cf000bf6ea1b647d6fdf0261a7527578d8b71b9d",
                        "inner_hash": "73bfee3a7c8d4f8a68657ebcaf69a59639f762bfc1a6f4468f3ca4724bc5b9f8",
                        "fee": 710046,
                        "sigs": [
                            "c4bcada17604a4a62baf50f929655027f2913639c27b773871f2135b72553c1959737e39d50e8349ffa5a7679de845aa6370999dbaaff4c7f9fd01260818683901"
                        ],
                        "inputs": [
                            {
                                "uxid": "4e75b4bced3404590d38ca06440c275d7fd86618a84966a0a1053fb18164e898",
                                "owner": "2JJ8pgq8EDAnrzf9xxBJapE2qkYLefW4uF8",
                                "coins": "26700.000000",
                                "hours": 54,
                                "calculated_hours": 811481
                            }
                        ],
                        "outputs": [
                            {
                                "uxid": "0a5603a1a5aeda575aa498cdaec5a4c893a28669dba84163eba2e90db3d9f39d",
                                "dst": "2JJ8pgq8EDAnrzf9xxBJapE2qkYLefW4uF8",
                                "coins": "26700.000000",
                                "hours": 101435
                            }
                        ]
                    }
                ]
            },
            "size": 183
        }
    ]
}
```


### Get block by hash or seq

API sets: `READ`

```
URI: /api/v1/block
Method: GET
Args:
    hash: get block by hash
    seq: get block by sequence number
```
Example:

```sh
curl http://127.0.0.1:6420/api/v1/block?hash=6eafd13ab6823223b714246b32c984b56e0043412950faf17defdbb2cbf3fe30
```

or

```sh
curl http://127.0.0.1:6420/api/v1/block?seq=2760
```

Result:

```json
{
    "header": {
        "seq": 2760,
        "block_hash": "6eafd13ab6823223b714246b32c984b56e0043412950faf17defdbb2cbf3fe30",
        "previous_block_hash": "eaccd527ef263573c29000dbfb3c782ee175153c63f42abb671588b7071e877f",
        "timestamp": 1504220821,
        "fee": 196130,
        "version": 0,
        "tx_body_hash": "825ae95b81ae0ce037cdf9f1cda138bac3f3ed41c51b09e0befb71848e0f3bfd",
        "ux_hash": "366af6bd80cfce79ce1ef63b45fb3ae8d9a6afc92a8590f14e18220884bd9d22"
    },
    "body": {
        "txns": [
            {
                "length": 220,
                "type": 0,
                "txid": "825ae95b81ae0ce037cdf9f1cda138bac3f3ed41c51b09e0befb71848e0f3bfd",
                "inner_hash": "312e5dd55e06be5f9a0ee43a00d447f2fea47a7f1fb9669ecb477d2768ab04fd",
                "sigs": [
                    "f0d0eb337e3440af6e8f0c105037ec205f36c83770d26a9e3a0fb4b7ec1a2be64764f4e31cbaf6629933c971613d10d58e6acb592704a7d511f19836441f09fb00"
                ],
                "inputs": [
                    "e7594379c9a6bb111205cbfa6fac908cac1d136e207960eb0429f15fde09ac8c"
                ],
                "outputs": [
                    {
                        "uxid": "840d0ee483c1dc085e6518e1928c68979af61188b809fc74da9fca982e6a61ba",
                        "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                        "coins": "998.000000",
                        "hours": 35390
                    },
                    {
                        "uxid": "38177c437ff42f29dc8d682e2f7c278f2203b6b02f42b1a88f9eb6c2392a7f70",
                        "dst": "2YHKP9yH7baLvkum3U6HCBiJjnAUCLS5Z9U",
                        "coins": "2.000000",
                        "hours": 70780
                    }
                ]
            }
        ]
    },
    "size": 220
}
```

### Get block by hash or seq Verbose

API sets: `READ`

```
URI: /api/v1/block/verbose
Method: GET
Args:
    hash: get block by hash
    seq: get block by sequence number
```

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
The calculated hours are the hours the transaction had in the block in which it was executed.

Example:

```sh
curl http://127.0.0.1:6420/api/v1/block/verbose?hash=6eafd13ab6823223b714246b32c984b56e0043412950faf17defdbb2cbf3fe30
```

or

```sh
curl http://127.0.0.1:6420/api/v1/block/verbose?seq=2760
```

Result:

```json
{
    "header": {
        "seq": 2760,
        "block_hash": "6eafd13ab6823223b714246b32c984b56e0043412950faf17defdbb2cbf3fe30",
        "previous_block_hash": "eaccd527ef263573c29000dbfb3c782ee175153c63f42abb671588b7071e877f",
        "timestamp": 1504220821,
        "fee": 196130,
        "version": 0,
        "tx_body_hash": "825ae95b81ae0ce037cdf9f1cda138bac3f3ed41c51b09e0befb71848e0f3bfd",
        "ux_hash": "366af6bd80cfce79ce1ef63b45fb3ae8d9a6afc92a8590f14e18220884bd9d22"
    },
    "body": {
        "txns": [
            {
                "length": 220,
                "type": 0,
                "txid": "825ae95b81ae0ce037cdf9f1cda138bac3f3ed41c51b09e0befb71848e0f3bfd",
                "inner_hash": "312e5dd55e06be5f9a0ee43a00d447f2fea47a7f1fb9669ecb477d2768ab04fd",
                "fee": 196130,
                "sigs": [
                    "f0d0eb337e3440af6e8f0c105037ec205f36c83770d26a9e3a0fb4b7ec1a2be64764f4e31cbaf6629933c971613d10d58e6acb592704a7d511f19836441f09fb00"
                ],
                "inputs": [
                    {
                        "uxid": "e7594379c9a6bb111205cbfa6fac908cac1d136e207960eb0429f15fde09ac8c",
                        "owner": "kbbzyrUKNVJsJDGFLAjVT5neVcx5SQjFx5",
                        "coins": "1000.000000",
                        "hours": 283123,
                        "calculated_hours": 302300
                    }
                ],
                "outputs": [
                    {
                        "uxid": "840d0ee483c1dc085e6518e1928c68979af61188b809fc74da9fca982e6a61ba",
                        "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                        "coins": "998.000000",
                        "hours": 35390
                    },
                    {
                        "uxid": "38177c437ff42f29dc8d682e2f7c278f2203b6b02f42b1a88f9eb6c2392a7f70",
                        "dst": "2YHKP9yH7baLvkum3U6HCBiJjnAUCLS5Z9U",
                        "coins": "2.000000",
                        "hours": 70780
                    }
                ]
            }
        ]
    },
    "size": 220
}
```

### Get last N blocks

API sets: `READ`

```
URI: /api/v1/last_blocks
Method: GET
Args:
    num: number of most recent blocks to return
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/last_blocks?num=2
```

Result:

```json
{
    "blocks": [
        {
            "header": {
                "seq": 58893,
                "block_hash": "8eca94e7597b87c8587286b66a6b409f6b4bf288a381a56d7fde3594e319c38a",
                "previous_block_hash": "1f042ed976c0cb150ea6b71c9608d65b519e4bc1c507eba9f1146e443a856c2d",
                "timestamp": 1537581594,
                "fee": 970389,
                "version": 0,
                "tx_body_hash": "1bea5cf1279693a0da24828c37b267c702007842b16ca5557ae497574d15aab7",
                "ux_hash": "bf35652af199779bc40cbeb339e8a782ff70673b07779e5c5621d37dfe13b42b"
            },
            "body": {
                "txns": [
                    {
                        "length": 377,
                        "type": 0,
                        "txid": "1bea5cf1279693a0da24828c37b267c702007842b16ca5557ae497574d15aab7",
                        "inner_hash": "a25232405bcef0c007bb2d7d3520f2a389e17e11125c252ab6c00168ec52c08d",
                        "sigs": [
                            "2ff7390c3b66c6b0fbb2b4c59c8e218291d4cbb82a836bb577c7264677f4a8320f6f3ad72d804e3014728baa214c223ecced8725b64be96fe3b51332ad1eda4201",
                            "9e7c715f897b3c987c00ee8c6b14e4b90bb3e4e11d003b481f82042b1795b3c75eaa3d563cd0358cdabdab77cfdbead7323323cf73e781f9c1a8cf6d9b4f8ac100",
                            "5c9748314f2fe0cd442df5ebb8f211087111d22e9463355bf9eee583d44df1bd36addb510eb470cb5dafba0732615f8533072f80ae05fc728c91ce373ada1e7b00"
                        ],
                        "inputs": [
                            "5f634c825b2a53103758024b3cb8578b17d56d422539e23c26b91ea397161703",
                            "16ac52084ffdac2e9169b9e057d44630dec23d18cfb90b9437d28220a3dc585d",
                            "8d3263890d32382e182b86f8772c7685a8f253ed475c05f7d530e9296f692bc9"
                        ],
                        "outputs": [
                            {
                                "uxid": "fb8db3f78928aee3f5cbda8db7fc290df9e64414e8107872a1c5cf83e08e4df7",
                                "dst": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                                "coins": "26.913000",
                                "hours": 970388
                            }
                        ]
                    }
                ]
            },
            "size": 377
        },
        {
            "header": {
                "seq": 58894,
                "block_hash": "3961bea8c4ab45d658ae42effd4caf36b81709dc52a5708fdd4c8eb1b199a1f6",
                "previous_block_hash": "8eca94e7597b87c8587286b66a6b409f6b4bf288a381a56d7fde3594e319c38a",
                "timestamp": 1537581604,
                "fee": 485194,
                "version": 0,
                "tx_body_hash": "c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd",
                "ux_hash": "f7d30ecb49f132283862ad58f691e8747894c9fc241cb3a864fc15bd3e2c83d3"
            },
            "body": {
                "txns": [
                    {
                        "length": 257,
                        "type": 0,
                        "txid": "c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd",
                        "inner_hash": "f7dbd09f7e9f65d87003984640f1977fb9eec95b07ef6275a1ec6261065e68d7",
                        "sigs": [
                            "af5329e77213f34446a0ff41d249fd25bc1dae913390871df359b9bd587c95a10b625a74a3477a05cc7537cb532253b12c03349ead5be066b8e0009e79462b9501"
                        ],
                        "inputs": [
                            "fb8db3f78928aee3f5cbda8db7fc290df9e64414e8107872a1c5cf83e08e4df7"
                        ],
                        "outputs": [
                            {
                                "uxid": "235811602fc96cf8b5b031edb88ee1606830aa641c06e0986681552d8728ec07",
                                "dst": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
                                "coins": "0.500000",
                                "hours": 1
                            },
                            {
                                "uxid": "873da4edc01c0b5184e1f26c4c3471dd407d08e9ab36b018ab93874e7392320b",
                                "dst": "2XBMMDMqTTYmqs2rfjEwYDz8ABd38y9B8r7",
                                "coins": "0.500000",
                                "hours": 1
                            },
                            {
                                "uxid": "42a6f0127f61e1d7bca8e9680027eddcecad772250c5634a03e56a8b1cf5a816",
                                "dst": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                                "coins": "25.913000",
                                "hours": 485192
                            }
                        ]
                    }
                ]
            },
            "size": 257
        }
    ]
}
```

### Get last N blocks Verbose

API sets: `READ`

```
URI: /api/v1/last_blocks/verbose
Method: GET
Args:
    num: number of most recent blocks to return
```

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
The calculated hours are the hours the transaction had in the block in which it was executed.

Example (verbose):

```sh
curl http://127.0.0.1:6420/api/v1/last_blocks/verbose?num=2
```

Result:

```json
{
    "blocks": [
        {
            "header": {
                "seq": 58893,
                "block_hash": "8eca94e7597b87c8587286b66a6b409f6b4bf288a381a56d7fde3594e319c38a",
                "previous_block_hash": "1f042ed976c0cb150ea6b71c9608d65b519e4bc1c507eba9f1146e443a856c2d",
                "timestamp": 1537581594,
                "fee": 970389,
                "version": 0,
                "tx_body_hash": "1bea5cf1279693a0da24828c37b267c702007842b16ca5557ae497574d15aab7",
                "ux_hash": "bf35652af199779bc40cbeb339e8a782ff70673b07779e5c5621d37dfe13b42b"
            },
            "body": {
                "txns": [
                    {
                        "length": 377,
                        "type": 0,
                        "txid": "1bea5cf1279693a0da24828c37b267c702007842b16ca5557ae497574d15aab7",
                        "inner_hash": "a25232405bcef0c007bb2d7d3520f2a389e17e11125c252ab6c00168ec52c08d",
                        "fee": 970389,
                        "sigs": [
                            "2ff7390c3b66c6b0fbb2b4c59c8e218291d4cbb82a836bb577c7264677f4a8320f6f3ad72d804e3014728baa214c223ecced8725b64be96fe3b51332ad1eda4201",
                            "9e7c715f897b3c987c00ee8c6b14e4b90bb3e4e11d003b481f82042b1795b3c75eaa3d563cd0358cdabdab77cfdbead7323323cf73e781f9c1a8cf6d9b4f8ac100",
                            "5c9748314f2fe0cd442df5ebb8f211087111d22e9463355bf9eee583d44df1bd36addb510eb470cb5dafba0732615f8533072f80ae05fc728c91ce373ada1e7b00"
                        ],
                        "inputs": [
                            {
                                "uxid": "5f634c825b2a53103758024b3cb8578b17d56d422539e23c26b91ea397161703",
                                "owner": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                                "coins": "25.910000",
                                "hours": 7745,
                                "calculated_hours": 17458
                            },
                            {
                                "uxid": "16ac52084ffdac2e9169b9e057d44630dec23d18cfb90b9437d28220a3dc585d",
                                "owner": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                                "coins": "1.000000",
                                "hours": 1915246,
                                "calculated_hours": 1915573
                            },
                            {
                                "uxid": "8d3263890d32382e182b86f8772c7685a8f253ed475c05f7d530e9296f692bc9",
                                "owner": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
                                "coins": "0.003000",
                                "hours": 7745,
                                "calculated_hours": 7746
                            }
                        ],
                        "outputs": [
                            {
                                "uxid": "fb8db3f78928aee3f5cbda8db7fc290df9e64414e8107872a1c5cf83e08e4df7",
                                "dst": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                                "coins": "26.913000",
                                "hours": 970388
                            }
                        ]
                    }
                ]
            },
            "size": 377
        },
        {
            "header": {
                "seq": 58894,
                "block_hash": "3961bea8c4ab45d658ae42effd4caf36b81709dc52a5708fdd4c8eb1b199a1f6",
                "previous_block_hash": "8eca94e7597b87c8587286b66a6b409f6b4bf288a381a56d7fde3594e319c38a",
                "timestamp": 1537581604,
                "fee": 485194,
                "version": 0,
                "tx_body_hash": "c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd",
                "ux_hash": "f7d30ecb49f132283862ad58f691e8747894c9fc241cb3a864fc15bd3e2c83d3"
            },
            "body": {
                "txns": [
                    {
                        "length": 257,
                        "type": 0,
                        "txid": "c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd",
                        "inner_hash": "f7dbd09f7e9f65d87003984640f1977fb9eec95b07ef6275a1ec6261065e68d7",
                        "fee": 485194,
                        "sigs": [
                            "af5329e77213f34446a0ff41d249fd25bc1dae913390871df359b9bd587c95a10b625a74a3477a05cc7537cb532253b12c03349ead5be066b8e0009e79462b9501"
                        ],
                        "inputs": [
                            {
                                "uxid": "fb8db3f78928aee3f5cbda8db7fc290df9e64414e8107872a1c5cf83e08e4df7",
                                "owner": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                                "coins": "26.913000",
                                "hours": 970388,
                                "calculated_hours": 970388
                            }
                        ],
                        "outputs": [
                            {
                                "uxid": "235811602fc96cf8b5b031edb88ee1606830aa641c06e0986681552d8728ec07",
                                "dst": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
                                "coins": "0.500000",
                                "hours": 1
                            },
                            {
                                "uxid": "873da4edc01c0b5184e1f26c4c3471dd407d08e9ab36b018ab93874e7392320b",
                                "dst": "2XBMMDMqTTYmqs2rfjEwYDz8ABd38y9B8r7",
                                "coins": "0.500000",
                                "hours": 1
                            },
                            {
                                "uxid": "42a6f0127f61e1d7bca8e9680027eddcecad772250c5634a03e56a8b1cf5a816",
                                "dst": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                                "coins": "25.913000",
                                "hours": 485192
                            }
                        ]
                    }
                ]
            },
            "size": 257
        }
    ]
}
```