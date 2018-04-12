typedef struct{
    GoUint64_ Coins;
    GoUint64_ Hours;
} wallet__Balance;
typedef struct{
    wallet__Balance Confirmed;
    wallet__Balance Predicted;
} wallet__BalancePair;
