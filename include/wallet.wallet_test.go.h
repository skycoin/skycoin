typedef GoMap_ mockBalanceGetter;
typedef struct{
    GoString_ name;
    GoUint64_ inputHours;
    GoUint64_ nAddrs;
    bool haveChange;
    GoUint64_ expectChangeHours;
    GoSlice_ expectAddrHours;
}distributeSpendHoursTestCase;
