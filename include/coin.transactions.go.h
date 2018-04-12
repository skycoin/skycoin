typedef GoSlice_  coin__Transactions;
typedef struct{
    coin__Transactions Txns;
    GoSlice_  Fees;
    GoSlice_  Hashes;
} coin__SortableTransactions;
typedef Handle coin__FeeCalculator;
typedef struct{
    GoUint32_ Length;
    GoUint8_ Type;
    cipher__SHA256 InnerHash;
    GoSlice_  Sigs;
    GoSlice_  In;
    GoSlice_  Out;
} coin__Transaction;
typedef struct{
    cipher__Address Address;
    GoUint64_ Coins;
    GoUint64_ Hours;
} coin__TransactionOutput;
