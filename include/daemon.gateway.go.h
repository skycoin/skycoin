typedef struct{
    GoInt_ BufferSize;
    BOOL DisableWalletAPI;
} daemon__GatewayConfig;
typedef Handle daemon__OutputsFilter;
typedef struct{
    visor__UnconfirmedTxnPooler uncfm;
    blockdb__UnspentPool unspent;
} daemon__spendValidator;
