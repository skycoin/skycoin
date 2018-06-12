typedef struct{
    GoUint32_  n[10];
} secp256k1go__Field;

typedef struct{
    secp256k1go__Field X;
    secp256k1go__Field Y;
    secp256k1go__Field Z;
    BOOL Infinity;
} secp256k1go__XYZ;

typedef struct{
    secp256k1go__Field X;
    secp256k1go__Field Y;
    BOOL Infinity;
} secp256k1go__XY;
