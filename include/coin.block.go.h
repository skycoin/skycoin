
typedef struct{
    cipher__SHA256 Hash;
    cipher__SHA256 PreHash;
} coin__HashPair;
typedef struct{
    GoUint32_ Version;
    GoUint64_ Time;
    GoUint64_ BkSeq;
    GoUint64_ Fee;
    cipher__SHA256 PrevHash;
    cipher__SHA256 BodyHash;
    cipher__SHA256 UxHash;
} coin__BlockHeader;
typedef struct{
    coin__Transactions Transactions;
} coin__BlockBody;
typedef struct{
    coin__BlockHeader Head;
    coin__BlockBody Body;
} coin__Block;

typedef struct{
    coin__Block _unnamed;
    cipher__Sig Sig;
} coin__SignedBlock;
