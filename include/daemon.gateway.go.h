typedef struct{
    GoInt_ BufferSize;
    BOOL EnableWalletAPI;
    BOOL EnableGUI;
} daemon__GatewayConfig;
typedef struct{
    GoSlice_  Txns;
} daemon__TransactionResults;
typedef struct{
    visor__TransactionStatus Status;
    GoUint64_ Time;
    visor__ReadableTransaction Transaction;
} daemon__TransactionResult;