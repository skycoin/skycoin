typedef struct{
    GoSlice_ blockPtr_slice;
    GoMap_ hash_to_blockPtr_map;
}BlockchainTail;
typedef struct{
    GoMap_ pubkey2sig;
    GoMap_ sig2none;
}HashCandidate;
