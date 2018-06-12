/*
//Requires typemap, contains GoString_
typedef struct{
    GoString_ UxID;
    GoString_ Address;
    GoString_ Coins;
    GoString_ Hours;
} api__CreatedTransactionOutput;

typedef struct{
    GoString_ UxID;
    GoString_ Address;
    GoString_ Coins;
    GoString_ Hours;
    GoString_ CalculatedHours;
    GoUint64_ Time;
    GoUint64_ Block;
    GoString_ TxID;
} api__CreatedTransactionInput;

*/
/*
typedef struct {
	BOOL 		neg;
	GoSlice_ 	nat;
} Number;

typedef struct {
	Number R;
	Number S;
} Signature;
*/
/*

//Contain slices. Should be Handle

typedef struct{
    visor__ReadableOutputSet Outputs;
} webrpc__OutputsResult;


typedef struct{
    cli__Balance Confirmed;
    cli__Balance Spendable;
    cli__Balance Expected;
    GoSlice_  Addresses;
} cli__BalanceResult;

typedef struct{
    wallet__BalancePair * Balance;
    visor__ReadableTransaction * Transaction;
    GoString_ Error;
} api__SpendResult;

typedef struct{
    api__CreatedTransaction Transaction;
    GoString_ EncodedTransaction;
} api__CreateTransactionResponse;

typedef struct{
    GoSlice_  Blocks;
} visor__ReadableBlocks;

typedef GoSlice_  coin__Transactions;


typedef struct{
    GoUint32_ Length;
    GoUint8_ Type;
    GoString_ TxID;
    GoString_ InnerHash;
    GoString_ Fee;
    GoSlice_  Sigs;
    GoSlice_  In;
    GoSlice_  Out;
} api__CreatedTransaction;

typedef struct{
    api__CreatedTransaction Transaction;
    GoString_ EncodedTransaction;
} api__CreateTransactionResponse;

typedef struct{
    coin__Transactions Txns;
    GoSlice_  Fees;
    GoSlice_  Hashes;
} coin__SortableTransactions;

//Should be Handle
typedef struct{
    daemon__TransactionResult * Transaction;
} webrpc__TxnResult;

typedef struct{
    GoInt_ N;
    BOOL IncludeDistribution;
} api__RichlistParams;

typedef struct{
} cli__PasswordFromTerm;

*/
