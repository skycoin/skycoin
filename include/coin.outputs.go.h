typedef struct{
    GoUint64_ Time;
    GoUint64_ BkSeq;
} coin__UxHead;

typedef GoSlice_  coin__UxArray;

typedef struct{
    cipher__SHA256 SrcTransaction;
    cipher__Address Address;
    GoUint64_ Coins;
    GoUint64_ Hours;
} coin__UxBody;

typedef struct{
    coin__UxHead Head;
    coin__UxBody Body;
} coin__UxOut;
