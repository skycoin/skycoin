typedef struct{
    GoString_ Addr;
    GoInt_ reqIDCtr;
} webrpc__Client;

typedef struct{
    GoString_ Status;
    GoInt_ StatusCode;
    GoString_ Message;
} webrpc__APIError;