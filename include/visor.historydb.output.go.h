typedef struct{
    GoString_ Uxid;
    GoUint64_ Time;
    GoUint64_ SrcBkSeq;
    GoString_ SrcTx;
    GoString_ OwnerAddress;
    GoUint64_ Coins;
    GoUint64_ Hours;
    GoUint64_ SpentBlockSeq;
    GoString_ SpentTxID;
} historydb__UxOutJSON;
struct _historydb__UxOuts{
	char c[1];
};
typedef struct _historydb__UxOuts historydb__UxOuts;
typedef struct{
    coin__UxOut Out;
    cipher__SHA256 SpentTxID;
    GoUint64_ SpentBlockSeq;
} historydb__UxOut;
