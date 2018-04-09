typedef GoSlice_ Transactions;
typedef struct{
    Transactions Txns;
    GoSlice_ Fees;
    GoSlice_ Hashes;
}SortableTransactions;
typedef Handle FeeCalculator;
