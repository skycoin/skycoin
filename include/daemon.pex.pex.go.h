typedef struct{
    GoString_ Addr;
    GoInt64_ LastSeen;
    BOOL Private;
    BOOL Trusted;
    BOOL HasIncomingPort;
    GoInt_ RetryTimes;
} pex__Peer;
