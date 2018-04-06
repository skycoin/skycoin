typedef struct{
    GoString_ CurrentSupply;
    GoString_ TotalSupply;
    GoString_ MaxSupply;
    GoString_ CurrentCoinHourSupply;
    GoString_ TotalCoinHourSupply;
    GoSlice_ UnlockedAddresses;
    GoSlice_ LockedAddresses;
}CoinSupply;
