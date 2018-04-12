typedef struct{
    GoString_ Version;
    GoString_ Commit;
} visor__BuildInfo;
typedef GoInterface_ visor__historyer;
typedef GoInterface_ visor__Blockchainer;
typedef GoInterface_ visor__UnconfirmedTxnPooler;
typedef GoInterface_ visor__TxFilter;
typedef struct{
    Handle f;
} visor__baseFilter;
typedef struct{
    GoSlice_  Addrs;
} visor__addrsFilter;
