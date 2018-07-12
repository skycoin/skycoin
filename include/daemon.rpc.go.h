typedef struct{
    GoInt_ ID;
    GoString_ Addr;
    GoInt64_ LastSent;
    GoInt64_ LastReceived;
    BOOL Outgoing;
    BOOL Introduced;
    GoUint32_ Mirror;
    GoUint16_ ListenPort;
} daemon__Connection;
typedef struct{
    GoSlice_  Connections;
} daemon__Connections;
typedef struct{
    GoUint64_ Current;
    GoUint64_ Highest;
    GoSlice_  Peers;
} daemon__BlockchainProgress;
typedef struct{
    GoSlice_  Txids;
} daemon__ResendResult;
struct _daemon__RPC{
	char c[1];
};
typedef struct _daemon__RPC daemon__RPC;
