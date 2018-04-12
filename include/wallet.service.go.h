typedef GoInterface_ wallet__BalanceGetter;
typedef struct{
    GoString_ WalletDir;
    wallet__CryptoType CryptoType;
    bool DisableWalletAPI;
} wallet__Config;
