typedef struct  {
	char c[1];
} visor__unconfirmedTxns;
typedef struct {
	char c[1];
} visor__txUnspents;
typedef struct{
    coin__Transaction Txn;
    GoInt64_ Received;
    GoInt64_ Checked;
    GoInt64_ Announced;
    GoInt8_ IsValid;
} visor__UnconfirmedTxn;
