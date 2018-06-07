typedef struct{
    cipher__Sig Sig;
    cipher__SHA256 Hash;
    GoUint64_ Seqno;
} consensus__BlockBase;
