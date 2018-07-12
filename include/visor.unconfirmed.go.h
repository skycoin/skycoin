struct _visor__unconfirmedTxns{
	char c[1];
};
typedef struct  _visor__unconfirmedTxns visor__unconfirmedTxns;
struct _visor__txUnspents{
	char c[1];
};
typedef struct _visor__txUnspents visor__txUnspents;
typedef struct{
    coin__Transaction Txn;
    GoInt64_ Received;
    GoInt64_ Checked;
    GoInt64_ Announced;
    GoInt8_ IsValid;
} visor__UnconfirmedTxn;
