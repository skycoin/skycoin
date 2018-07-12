typedef struct{
    GoString_ WalletDir;
    GoString_ WalletName;
    GoString_ DataDir;
    GoString_ Coin;
    GoString_ RPCAddress;
    BOOL UseCSRF;
} cli__Config;
typedef struct{
    GoInt32_ _unnamed;
} cli__WalletLoadError;
typedef struct{
    GoInt32_ _unnamed;
} cli__WalletSaveError;
typedef GoSlice_  cli__PasswordFromBytes;
struct _cli__PasswordFromTerm{
	char c[1];
};
typedef struct _cli__PasswordFromTerm cli__PasswordFromTerm;
