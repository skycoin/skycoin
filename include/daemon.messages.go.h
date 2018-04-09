typedef struct{
    GoSlice_ Messages;
} MessagesConfig;
typedef struct{
    MessagesConfig Config;
    GoUint32_ Mirror;
} Messages;
typedef struct{
    GoUint32_ IP;
    GoUint16_ Port;
} IPAddr;
typedef GoInterface_ AsyncMessage;
typedef struct{
    GoString_ addr;
} GetPeersMessage;
typedef struct{
} PongMessage;
