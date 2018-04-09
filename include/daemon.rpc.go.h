typedef struct{
    GoInt_ ID;
    GoString_ Addr;
    GoInt64_ LastSent;
    GoInt64_ LastReceived;
    bool Outgoing;
    bool Introduced;
    GoUint32_ Mirror;
    GoUint16_ ListenPort;
}Connection;
typedef struct{
    GoSlice_ Connections;
}Connections;
typedef struct{
    GoUint64_ Current;
    GoUint64_ Highest;
    GoSlice_ Peers;
}BlockchainProgress;
typedef struct{
    GoSlice_ Txids;
}ResendResult;
typedef struct{
}RPC;
