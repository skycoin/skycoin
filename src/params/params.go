package params

/*
CODE GENERATED AUTOMATICALLY WITH FIBER COIN CREATOR
AVOID EDITING THIS MANUALLY
*/

const (
	// Distribution locking parameteres

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
)

var (
	// Transaction verification parameters

	// UserVerifyTxn transaction verification parameters for user-created transactions
	UserVerifyTxn = VerifyTxn{
		// BurnFactor can be overriden with `USER_BURN_FACTOR` env var
		BurnFactor: 2,
		// MaxTransactionSize can be overriden with `USER_MAX_TXN_SIZE` env var
		MaxTransactionSize: 32768, // in bytes
		// MaxDropletPrecision can be overriden with `USER_MAX_DECIMALS` env var
		MaxDropletPrecision: 3,
	}
)

// distributionAddresses are addresses that received coins from the genesis address in the first block,
// used to calculate current and max supply and do distribution timelocking
var distributionAddresses = [DistributionAddressesTotal]string{
	"TkyD4wD64UE6M5BkNQA17zaf7Xcg4AufwX",
	"2PBcLADETphmqWV7sujRZdh3UcabssgKAEB",
	"Ua76nGXLLkTQeEGhgVQRhuewzRNzVw1F91",
	"8zwjnNuETXmgAZYc3pt332Amm8nMRHY55v",
	"2Jnk5YsETRMrTbncBLa451KFWhbDYMibGZy",
	"2bnJybodi4Sjxswq6FkiMpcoDAen9PNZJRo",
	"9khkGALs4MH33QaHjrnWxV9EiAGAWvcLfu",
	"JU5PJvFcyao8u5yewaadK4g6RMrEqaDGHj",
	"5WYS8hdVVeayeRgGEKNCujzve1w2M3PR9Y",
	"JwheRb7GJyB7rU5fsaArS5U2Kq4cpuGHv3",
	"zH1RfGdeBxtuNZ7DLJ4x92EVJmVzXksHMa",
	"2BFUUKWUgcusGc6CiGiy5YjJv7J2BFpHPef",
	"RC1jYstPsnbrZgUKu99ba8fNZw6eG6Kgwa",
	"2PJMA5Bq5T3SZDYL1mJpHhSikP97qcLQr6u",
	"UvHE5vgptjPakvhtrs6QpfK1d4tEK9NVH8",
	"294KAJFmJTzSYDuFezj8PePXoVmHkb693tR",
	"2m8CWuuYwBfFPnbxewXE34MA3VRFwx3skvY",
	"28m2H6zMcbkyuL4wtPjjWcuGq8yMCbGTNMM",
	"ZUv7ZAxcXFhAq26B5yhSvrMKLTHsYeJgnJ",
	"yR28PRxfJ4Tju4sWk5GwKYPsPVMgdpUgzL",
	"R5vGXHJoLSrd3gTrhZtUoumUCNDuqyvYVJ",
	"35gVEEfxZnGZNmc2RhbCg5k2uL6AjZwSe1",
	"2NvybQP9CQyihrWq8CJH7CkS8zu2Fs7sY9w",
	"9vXYTnaQQFuKGB4yaaAUrjfq14D9x3g7pD",
	"2534zD6kyJKSACtxuRMfqsarkGziVeTYcNZ",
	"2Lb4DBPgqr5PUei2Sofx2A6deZf7fWsfBbJ",
	"5bW4xUssBYABUPHYStmzkJ6xAV6gWjfJGq",
	"2FcQhvqXimVDXPT1nUwvNu2mPD4mQcgNxEc",
	"2d64uzJ2TY4sHZQEDHnsEpUzf1Vy5tJ5dMG",
	"29yJCV5KRpGmtXq9uqu6hreVEksiH1DsJhW",
	"oAv19vdjy7Kyoi3n3KaLA32AeUHoxUEip1",
	"RJGP3EX1tt8Hr42E86ZXnwPMQKLrNznnzh",
	"PhpENSBNmt6gsud5FZqSi7jaawxKqf2B1v",
	"2UqA2LwHMUkYXqJU8sspXjzPWgHYmuNZisM",
	"LkLWGp54ui3WYKGitfZx9HTcrNkYaT41EW",
	"25BFvqEVZRLvzTXJZNVL2DSjeXj1mTRGBRK",
	"3xbP76R6ToFS1LPqjm2X6UgqJ5JzDQ2g7U",
	"4LSpB3FzSHHjYxg2T9n8M5mwk36TJusccD",
	"UFohhJEzYw9RWNEaDrADYj9VF1cSS8e2Bc",
	"6KaSuGH2SPCP8HigVPyNwwYY8rtpHPX8E1",
	"vbnskpo3Jkxa67vYjMYTujzQiaxbdY5xem",
	"Ukczjnwi8Q3ugmR4inPMNexXGj1TroG9UN",
	"2KzLraUjhxEYRwg5qzyrx4XJnPye5SkzHVC",
	"BoAqBbzmSyri9XRRUKHnNPDsz7uYp2yM8k",
	"3CaU8dir2r2Eakfz2Ycqi6kYd8uwtDBU4W",
	"ZpSApsPVXRQxfZTb6Rz37KvnG1SAQC9WNC",
	"CU9wBoMn5szwTVqkvyGZn3eECA3Chh5bce",
	"Aw5ejAXbVWRCBWLPFRx8gcEvhi8FnMZJ2X",
	"42DsZHmbt7UuSVSyNm3RJrN5DhEGJ878Ka",
	"2qxvUeMLTGKXy16UJEs7AbBv4VRNhD3BY4",
	"2ULcAggy5CfB3ccRzAC7PSNs9WQAvTXU2D7",
	"2JNeCBDtqh68QePcxHycX3sDFcWWb9t3H8Q",
	"Paj5g4RzcwGKmrCSpEkrbZQpARYtKryaou",
	"uwk8qUok3JGwTeVJsKrPiThLFDJtQUCLJQ",
	"2igXUG4hf3Hkuw8NWFmtRraLaBrp6SVQQQ4",
	"4avsgXGUMDPZC4tU7vDqJjFoqJqCyfi9A4",
	"PkAK7sFDRVrJGqAfU7tqhernkwspDNVsud",
	"2BczL7pXbfhG3s5MBAfqafnLK7jKBpz1L2G",
	"2RhUHBUErQx5Jbjam893HHZDNEygARWjA8o",
	"GzkpicxJud7RrHqkhbU5pGNQLnpLNuA2bJ",
	"2cVqwo4DpDrceDuxoEnAobsjt6RXN1V5HqB",
	"2QZeKZTrK52mowYWVwp1nkR9YCz5ydVbhTW",
	"Z4NUGBA37PoyDnAoY8NAhMkCnpu2mvFCfn",
	"2CSZsiyqS6qaj2WT4LAg1Pew2VwycwetMx6",
	"2QjGLDduBzsf65386BPAA8qcpHSnUAfpuSN",
	"ysUrw6s15qQSfahRy5a2LDofy4VbeJ5Nuw",
	"2EkKpjYVj613enR4VFcpfa4PpwD1evGatYy",
	"2Y9FNVfMMB8A4x7atsnMidzncDT9obu4Wzk",
	"eJL1kJj8Su94cFk5KF9Nhr3axkzmw47TxP",
	"6cE7xcup67ajuMZ98uuNt3zktGu9PohR75",
	"2BTDNu1HUp2GjDfT4TFCMXsD1qhahApSoCj",
	"AHaRFT2aBGmyxJzhJJ1wyDot6tRWL2Dn53",
	"2ATQ6AA2xXoJew6uMjAqGbbwg7YCx9EmVmF",
	"aEDD3o6YXBvTBTC3p8HmuSrGCVe1niaqtv",
	"2fdwHM7jpdEd16iYNpjQiNAksnsDLXGBp72",
	"2PXmFhvv7rtoRy7aGwrEU8DV5NKawf6RV3v",
	"2DeBd7xfLbsBT8xejTvxiSokJiceKvPtBvC",
	"4ttUh7dqKeGcZRxNE9HpQAiTSWk1K22ibh",
	"5T29pJYY73SYajqK6RDu2kCnfY2ViCHNcN",
	"2dF2ZkRaoKismZUD1vWmVkN7RqrScRxdQGR",
	"2Hr7iov8MjwTG1uEoKn8ehXMp6HqRv2ZynV",
	"2ES4km2wCTLcQCgRvWJTU2BsJxr7E5C61ji",
	"2DdYNq84G8ZLgRRKKrmuH4MBevZfyc3hrsp",
	"sWvYCtFLXp5P9NY1GiBJg3ffamTcKHSAvT",
	"7inetACS91yoGZpDTgHWx6EvgXgFgj4jb4",
	"2TrKEZQVrHHemFt79FaRq6GHJyYCFLDfLpE",
	"aFGrjSV2dEdrQUfkqKKkzMtDemtmYMSUz6",
	"YyxzMj91LnpjwGx4Fkg8ZZNrqpm7fHTcce",
	"sdceAbRsMvpfiosNZgpqcXQ1PcHP2GtKGs",
	"jkPYkPf6yse38RuV3azCLmuGWKayj3QyQG",
	"tSwngm1ae4cbPtBt7x2EUZAdcxq6jjNarc",
	"2cviecT9Q4q4vQqYQytULs8zLzXau46mkLM",
	"2At3p84cMd76Y9RgjWsSeSakgSfsAXrMkwb",
	"2J2DmXud1YFjKaMfaEdx9uo5uBpxchpNkTc",
	"uth7PjztD2hR59HrUxzJ4uxVcoMvJ11KoC",
	"2LPkr1o39J3bCQKwozMtaU6WyWgEzTc94Zs",
	"poUaV113NYB2AQkwKuf7HvTwWscvJZ9KC6",
	"2er78M3ZuhSELxy6DEmMyWzQySmVzPB4gXA",
	"2SnuqPeC7Bt6b8shKfXEP3hxpUMVZbWNnps",
	"2DbQ1prnW67sWnJZ4Hrhqd2ttnpZ4xZ1iUA",
}
