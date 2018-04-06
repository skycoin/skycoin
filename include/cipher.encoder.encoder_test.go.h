typedef struct{
    GoInt32_ X;
    GoInt64_ Y;
    GoUint8_ Z;
    GoSlice_ K;
    bool W;
}TestStruct2;
typedef struct{
    GoInt32_ X;
    GoInt64_ Y;
    GoUint8_ Z;
    GoSlice_ K;
}TestStructIgnore;
typedef struct{
    GoInt32_ X;
    GoInt64_ Y;
    GoSlice_ K;
}TestStructWithoutIgnore;
typedef struct{
    GoInt32_ X;
    GoSlice_ K;
}TestStruct3;
typedef struct{
    GoInt32_ X;
    GoInt32_ Y;
}TestStruct4;
typedef struct{
    GoInt32_ X;
    GoSlice_ A;
}TestStruct5;
typedef struct{
    GoUint32_ X;
    GoUint64_ Y;
    GoSlice_ Bytes;
    GoSlice_ Ints;
}Contained;
typedef struct{
    GoSlice_ Elements;
}Container;
typedef struct{
    GoSlice_ Arr;
}Array;
typedef struct{
    GoUint64_ Test;
}TestStruct5a;
