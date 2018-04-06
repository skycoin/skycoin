typedef struct{
    GoString_ password;
    GoString_ salt;
    GoInt_ N, r, p;
    GoSlice_ output;
}testVector;
