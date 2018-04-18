typedef struct{
    GoString_ Addr;
    bool Solicited;
} daemon__ConnectEvent;
typedef struct{
    GoString_ Addr;
    GoInt32_ Error;
} daemon__ConnectionError;
typedef struct{
    GoString_ Addr;
    gnet__DisconnectReason Reason;
} daemon__DisconnectEvent;
typedef struct{
    daemon__AsyncMessage Message;
    gnet__MessageContext * Context;
} daemon__MessageEvent;
