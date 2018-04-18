typedef struct{
    GoSlice_  Transactions;
} gui__UnconfirmedTxnsResponse;
typedef struct{
    GoString_ Address;
} gui__WalletFolder;
typedef struct{
    wallet__BalancePair * Balance;
    visor__ReadableTransaction * Transaction;
    GoString_ Error;
} gui__SpendResult;
