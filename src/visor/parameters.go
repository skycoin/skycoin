package visor

/*
CODE GENERATED AUTOMATICALLY WITH FIBER COIN CREATOR
AVOID EDITING THIS MANUALLY
*/

const (
	// MaxCoinSupply is the maximum supply of coins
	MaxCoinSupply uint64 = 100000000
	// DistributionAddressesTotal is the number of distribution addresses
	DistributionAddressesTotal uint64 = 100
	// DistributionAddressInitialBalance is the initial balance of each distribution address
	DistributionAddressInitialBalance uint64 = MaxCoinSupply / DistributionAddressesTotal
	// InitialUnlockedCount is the initial number of unlocked addresses
	InitialUnlockedCount uint64 = 25
	// UnlockAddressRate is the number of addresses to unlock per unlock time interval
	UnlockAddressRate uint64 = 5
	// UnlockTimeInterval is the distribution address unlock time interval, measured in seconds
	// Once the InitialUnlockedCount is exhausted,
	// UnlockAddressRate addresses will be unlocked per UnlockTimeInterval
	UnlockTimeInterval uint64 = 31536000 // in seconds
	// MaxDropletPrecision represents the decimal precision of droplets
	MaxDropletPrecision uint64 = 3
	//DefaultMaxBlockSize is max block size
	DefaultMaxBlockSize int = 32768 // in bytes
)

var distributionAddresses = [DistributionAddressesTotal]string{
	"eZjwb8g8N4Vev1yPuzcSnh78AgJV6SGHFK",
	"xSsKwLUAcrFbW97Uc87UQq6zcCvQDmm1gN",
	"rdbWZL2jGGp7thg4YNvZ3okxsYpxG4ztg1",
	"2e6JPoQf4vzjCNjdtoVFQTuaAiCGwH4Y55C",
	"mj6PMjsQDy2CPsqYqjXmT8HUzwFzWJ8rJt",
	"zWHV3z8VRAsvUzgMJG2NmnKmG1DjQTxxa5",
	"FYne2xrqvUwqgey62neJ3muRDuT9D3tLp3",
	"J1t76VTFtCahWyYho2mEykbqvtqp8cAGnR",
	"2JUUijX6WD3NCV4CRY4EhzAo9sMQ1Ua9UwZ",
	"7BbazN2fAxU4zqdq9NcpB2gTjxNDtYQFep",
	"2jvuMseFnaTxKCD7LKqQtbsdTG2ykqFo7YZ",
	"fZbFU9dxwWDotRbnwu8rqMqHnTaybDe7Pr",
	"AHPdavUS2189WDwHQgrGmNPLE5JFfCLzE9",
	"22fEpMTRdVB5g9gFuJjQvhBiZK3h67z3NVZ",
	"un66r3g4jfoGHiFQx8A7Y6G6tjScUUz5MS",
	"22pP5Xz1U8Ri6XkiHxfp52LujiCDXDgRU6R",
	"VCuc9RXPb4iYysynYvVECnqGcZDGzZt9q7",
	"hUKoybfJtTMC9QbT1SSqCGJB1a4UDnW1cB",
	"2EKNvNv2sg2zH6M3fCsP74sDjXu5RGDhVBG",
	"SJYvLbeqjgX4xkQo6stCAbT12NtXSAkssn",
	"2CVbeXD7kfaQ3rCKw1BkvyF2rdAvwLfneGS",
	"233CA5qhMrQL6bUVtyXegJfGYUcXRgjmoiy",
	"2N9j8SdmhuAcUSWPk7spVCcTGLYFjQnXyA2",
	"241d7CGiDMTnBz9kSSwc7kS9cg44ppbBj4o",
	"rZC69jgM3RccLXNHVaz6rvEzJ4Zivn4AkR",
	"VG41kWMQmvpJaxCYA7cThzVpCZ5VyzxzR6",
	"5Jg7K7X6cqzFah8CEsWvrhShAdQTDiXm3D",
	"2NbekABazbCTRUwANiSUCgYsgVNNULLu6Es",
	"oseuGSu8uA5nPjaymznDgtjdemt4w2v1d1",
	"2CqHrgmvdzioz2gA8KARTN3DxaWbHxD8L8Q",
	"GWueTCnhcmotaFsRciJhDdoC4gu1TQwXZT",
	"2Aj3ttHErn6ZTfbVap6daLevxpRSby6s7fQ",
	"GhLT7XudyCP9L7zgB3dypQVXBfur4hQJjR",
	"nTv3PvJvfE4rBxEyNdxSofVkCGZdzFhuFa",
	"cKGyuyweDraXqnV3DDoFL54R1edHKrJ38U",
	"CtxknfpLs8XCJbhFmh6mSHb4D6WYrpqkEp",
	"2T4C6e9aTgn4neZCXAKdZUdTbBidePk3Ars",
	"2UvLUqxwT7UNyNcrzaNCYxbhF5Neze6yc5N",
	"eXLdnMUqXyQGWrsBMgrVGHLHo9LhMzXj6Z",
	"ocomNYyYBeNm3GPVz1RjwZyMqkdN3i6zhw",
	"wPKhiU9SYUiZSK6Yk56Ah94CLcVvqx2QF8",
	"2fmueDWXi3vJ48Ln4NNTSG6eup6CowZCnpQ",
	"yBhTEKUALUqSh3dSvMkDR5FYDVnQs8agd7",
	"2DZPnKFo8LcCLw2z4FScSqfaks3kqGnAatH",
	"cFzsdxqjKjhFWiQNZ97mykAXvcBiNzEac9",
	"275QJFex4oEKKhUR9sVEgwEevSGaMgUoAr2",
	"2BWyjibMwcD4SWzpxJTYe9hYNhQgJfTQb6n",
	"2JA8Ei73EN4bGWoNn4Jdujt8GDWrX5J1yNB",
	"ZMUcfRGHeoLGzDyHwAjxXVC9swbS133SoC",
	"MDkXrv93axi7czJp5hghgSxwyJ1bTCGQhy",
	"g3pDwboPbKHyS4DJ3piZdr2r8KK7qVybXf",
	"ZVzTrvUbxApQkWEkWq3yqfBfztUfW5vKwP",
	"2UznikdMYCBZUiV8eAUygt5UCw613KsFmyi",
	"3nGD5xavnRn87GgwK3GS8HwxvZ88ebXnGc",
	"BiWnLzR5SPmbNYtBNKVBU2Z2BKUDZExwb1",
	"28DUeW3uC3NXCLv3GfMLb5W2vnVN4wkrARp",
	"yhJCQSGuQoVPc2LQSfd1nE1hL8kqTRDRRQ",
	"HdjZsdu3jcmtTCWzZ99jDa3Y2i9w3KVU2P",
	"EuteC1vAuz8q1tXrMqqjyVEYEeDMSuzaFH",
	"SNJkc6eNx6Le1QaE2KP9jsZx7ETp6ATLZS",
	"2GPCcJPPC1Lp977BwckXFSDXAoTADHYq5Kg",
	"nPsFcPsm6GEqjvhKQTymATuABDYLQAig5a",
	"GKQxkaShcij82bAKeZbCRZxzebkUsQTuif",
	"2QJrG5u8tSB9XztHv8H8r6PRNspvSf8EaqR",
	"o1XDGmq2MkBGtdZX3mvw31RGisR3g6diUy",
	"Uu1o7Ridb29gS5FwWrVA8r3hjSGnS2Qpis",
	"QzVZZFbpWgRZRccjRUrgqHNPd8uGH4wWs5",
	"MHivapHxN47JZr7v44U5tL8TSGYPVvRo9o",
	"29zMongabk63DUn1HmVUJhea862DCCs9Kiw",
	"21KxZwc63JYSAfBG2Z7jWULdonXJ27attJP",
	"NtCE1urnwKfwgHk3L8v4aknSxKZAbKNhLC",
	"WvLwEtmoJUEPPaz2DcCYXRy9j32gqfwMUk",
	"2F94Bet5bV4r6AuWMetcdstwH32wVwAXtMq",
	"2iWDqN1jdSTsrQNcJTKfdQCEhoEgJkeyFGE",
	"1EUqEBDYctNuUVHA2riQ49ioJVMfUvv2w2",
	"RZAUKzXxmX9FiRyMgk8nUj2GPUMfTTLkg1",
	"knMTrbWkqMaPanRTuvpQzFYN1ooviZ57wR",
	"2SakgsmRYgiBkunETidzJTyENMtkdZj7MrP",
	"2FEr6xQfNGFneimP7tywhZ3TjX6eZYXVXQa",
	"kw36cLyRPciXW8gmRo12hXoGeVhep4CoQC",
	"sL3RgVKrxPcfF6Ja1TBY3VYLTstAyhjd8v",
	"dopMaLXpee7t4TD1y1Wo4iRBqErikgFGJP",
	"2AHeWqHopEKZT262LXRxepU3ooZF9vYSRJE",
	"cSthsm6JHu7JdMWNrzrzBjsGnkJpYwsEEk",
	"rz9cXngg8eDssR9qfRuF4fzcfe1dktCUZ2",
	"7UzPhWqL9oaxBuk84Zud3WWscoThxyEgvy",
	"TjSsRg8HAWMiCFU91nfw6njZ3c81bdfHAu",
	"ZdBesshBQdBbk5oZTZLN6ekuihBFcyk61G",
	"4FXKFeV38src6hVorswB9kEfneABoa9ae3",
	"JFR2LdLu18b4ywqoPsoqS1cZjW6WNJioBk",
	"2Uj98Q8nfNAeRvVC1CME2pVkDX7znmSm1ZG",
	"4tHyWvqa3pBmzXfJRPRKUZ4YePznvp5Guo",
	"c8woDuthjenfymXiwYZSvqWAB71yPNawu6",
	"22stmS4eyFF1zeNieaAoYKPRiVB1wmxenvm",
	"b8SECsjHFLhrJAvDo92oDq9cLEUXk7EBMH",
	"2BBHUuKaqbutkQreGcvgRMUKxAMJrbsaK8x",
	"2fRCN3WuPDPG3585SqGf3oMR9c5Uj4XqFtv",
	"2Vj8jxsL1uXdLQT3iFsub5Y2Vy4vsJraKCp",
	"2PbZCF1HEfquVavTPxgffDeRgrrYFueQVRG",
	"A1DeRXNpKXmqiaJFRN7Q3gh7kFdspZmpeX",
}
