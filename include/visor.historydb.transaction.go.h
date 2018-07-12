struct _historydb__transactions{
	char c[0];
};
typedef struct _historydb__transactions historydb__transactions;
typedef struct{
    coin__Transaction Tx;
    GoUint64_ BlockSeq;
} historydb__Transaction;
