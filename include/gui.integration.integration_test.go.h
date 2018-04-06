typedef struct{
    GoInterface_ actual;
    GoInterface_ expected;
}TestData;
typedef struct{
    GoString_ name;
    GoString_ address;
    GoString_ golden;
    GoInt_ errCode;
    GoString_ errMsg;
}addressTransactionsTestCase;
