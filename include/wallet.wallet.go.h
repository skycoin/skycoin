typedef GoString_ wallet__CoinType;
typedef struct{
    GoMap_ Meta;
    GoSlice_  Entries;
} wallet__Wallet;
typedef GoInterface_ wallet__Validator;
typedef struct{
    wallet__CoinType Coin;
    GoString_ Label;
    GoString_ Seed;
    bool Encrypt;
    GoSlice_  Password;
    wallet__CryptoType CryptoType;
} wallet__Options;
typedef struct{
    cipher__SHA256 Hash;
    GoUint64_ BkSeq;
    cipher__Address Address;
    GoUint64_ Coins;
    GoUint64_ Hours;
} wallet__UxBalance;
