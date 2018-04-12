typedef GoSlice_  pex__Peers;
typedef struct{
    GoMap_ peers;
} pex__peerlist;
typedef Handle pex__Filter;
typedef struct{
    GoString_ Addr;
    GoInterface_ LastSeen;
    bool Private;
    bool Trusted;
    bool * HasIncomePort;
    bool * HasIncomingPort;
} pex__PeerJSON;
