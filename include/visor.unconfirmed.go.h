typedef GoMap_ visor__TxnUnspents;
typedef Handle visor__UnspentGetFunc;
typedef struct{
    coin__Transaction Txn;
    GoInt64_ Received;
    GoInt64_ Checked;
    GoInt64_ Announced;
    GoInt8_ IsValid;
} visor__UnconfirmedTxn;
