typedef struct{
    GoString_ Addr;
    GoInt64_ LastSeen;
    bool Private;
    bool Trusted;
    bool HasIncomingPort;
    GoInt_ RetryTimes;
} Peer;
