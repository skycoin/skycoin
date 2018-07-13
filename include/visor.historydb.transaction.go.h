typedef struct {
	char c[1];
} historydb__transactions;
typedef struct{
    coin__Transaction Tx;
    GoUint64_ BlockSeq;
} historydb__Transaction;
