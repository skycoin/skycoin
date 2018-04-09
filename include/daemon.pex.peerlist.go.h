typedef GoSlice_ Peers;
typedef struct{
    GoMap_ peers;
} peerlist;
typedef Handle Filter;
typedef struct{
    GoString_ Addr;
    GoInterface_ LastSeen;
    bool Private;
    bool Trusted;
    bool * HasIncomePort;
    bool * HasIncomingPort;
} PeerJSON;
