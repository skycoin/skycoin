typedef struct{
    GoMap_ transactions;
    GoMap_ injectRawTxMap;
    GoMap_ injectedTransactions;
    GoSlice_ addrRecvUxOuts;
    GoSlice_ addrSpentUxOUts;
    GoSlice_ uxouts;
}fakeGateway;
