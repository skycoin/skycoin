typedef struct{
    GoString_ Coins;
    GoString_ Hours;
}Balance;
typedef struct{
    Balance Confirmed;
    Balance Spendable;
    Balance Expected;
    GoString_ Address;
}AddressBalance;
typedef struct{
    Balance Confirmed;
    Balance Spendable;
    Balance Expected;
    GoSlice_ Addresses;
}BalanceResult;
