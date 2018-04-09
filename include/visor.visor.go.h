typedef struct{
    GoString_ Version;
    GoString_ Commit;
}BuildInfo;
typedef GoInterface_ historyer;
typedef GoInterface_ Blockchainer;
typedef GoInterface_ UnconfirmedTxnPooler;
typedef GoInterface_ TxFilter;
typedef struct{
    Handle f;
}baseFilter;
typedef struct{
    GoSlice_ Addrs;
}addrsFilter;
