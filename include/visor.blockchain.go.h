typedef struct{
    BOOL Arbitrating;
    cipher__PubKey Pubkey;
} visor__BlockchainConfig;
typedef struct{
    cipher__Sig sig;
    cipher__SHA256 hash;
} visor__sigHash;
