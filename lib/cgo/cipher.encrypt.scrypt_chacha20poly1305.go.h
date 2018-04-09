typedef struct{
    GoInt_ N;
    GoInt_ R;
    GoInt_ P;
    GoInt_ KeyLen;
}ScryptChacha20poly1305;
typedef struct{
    GoInt_ N;
    GoInt_ R;
    GoInt_ P;
    GoInt_ KeyLen;
    GoSlice_ Salt;
    GoSlice_ Nonce;
}meta;
