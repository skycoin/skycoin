typedef struct{
    GoSlice_  Messages;
} daemon__MessagesConfig;
typedef struct{
    daemon__MessagesConfig Config;
    GoUint32_ Mirror;
} daemon__Messages;
typedef struct{
    GoUint32_ IP;
    GoUint16_ Port;
} daemon__IPAddr;
typedef struct{
    GoString_ addr;
} daemon__GetPeersMessage;
typedef struct{
} daemon__PongMessage;
typedef struct{
    GoSlice_  Peers;
    gnet__MessageContext * c;
} daemon__GivePeersMessage;
typedef struct{
    GoUint32_ Mirror;
    GoUint16_ Port;
    GoInt32_ Version;
    gnet__MessageContext * c;
    BOOL valid;
} daemon__IntroductionMessage;
typedef struct{
    gnet__MessageContext * c;
} daemon__PingMessage;
