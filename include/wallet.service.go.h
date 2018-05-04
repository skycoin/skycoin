typedef GoInterface_ wallet__BalanceGetter;
typedef struct{
    GoString_ WalletDir;
    wallet__CryptoType CryptoType;
    BOOL DisableWalletAPI;
} wallet__Config;
