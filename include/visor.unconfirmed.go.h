struct _visor__unconfirmedTxns{
};
typedef struct  _visor__unconfirmedTxns visor__unconfirmedTxns;
struct _visor__txUnspents{
};
typedef struct _visor__txUnspents visor__txUnspents;
typedef struct{
    coin__Transaction Txn;
    GoInt64_ Received;
    GoInt64_ Checked;
    GoInt64_ Announced;
    GoInt8_ IsValid;
} visor__UnconfirmedTxn;
