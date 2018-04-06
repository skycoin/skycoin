typedef struct{
    bool Confirmed;
    bool Unconfirmed;
    GoUint64_ Height;
    GoUint64_ BlockSeq;
    bool Unknown;
}TransactionStatus;
typedef struct{
    GoString_ Hash;
    GoString_ Address;
    GoString_ Coins;
    GoUint64_ Hours;
}ReadableTransactionOutput;
typedef struct{
    GoString_ Hash;
    GoString_ Address;
    GoString_ Coins;
    GoUint64_ Hours;
}ReadableTransactionInput;
typedef struct{
    GoString_ Hash;
    GoUint64_ Time;
    GoUint64_ BkSeq;
    GoString_ SourceTransaction;
    GoString_ Address;
    GoString_ Coins;
    GoUint64_ Hours;
    GoUint64_ CalculatedHours;
}ReadableOutput;
typedef GoSlice_ ReadableOutputs;
typedef struct{
    GoUint32_ Length;
    GoUint8_ Type;
    GoString_ Hash;
    GoString_ InnerHash;
    GoUint64_ Timestamp;
    GoSlice_ Sigs;
    GoSlice_ In;
    GoSlice_ Out;
}ReadableTransaction;
typedef struct{
    GoUint64_ BkSeq;
    GoString_ BlockHash;
    GoString_ PreviousBlockHash;
    GoUint64_ Time;
    GoUint64_ Fee;
    GoUint32_ Version;
    GoString_ BodyHash;
}ReadableBlockHeader;
typedef struct{
    GoSlice_ Transactions;
}ReadableBlockBody;
typedef struct{
    ReadableBlockHeader Head;
    ReadableBlockBody Body;
}ReadableBlock;
typedef struct{
    GoString_ Hash;
    GoString_ SourceTransaction;
    GoString_ Address;
    GoString_ Coins;
    GoUint64_ Hours;
}TransactionOutputJSON;
typedef struct{
    GoString_ Hash;
    GoString_ InnerHash;
    GoSlice_ Sigs;
    GoSlice_ In;
    GoSlice_ Out;
}TransactionJSON;
typedef struct{
    ReadableBlockHeader Head;
    GoUint64_ Unspents;
    GoUint64_ Unconfirmed;
}BlockchainMetadata;
typedef struct{
    ReadableOutputs HeadOutputs;
    ReadableOutputs OutgoingOutputs;
    ReadableOutputs IncomingOutputs;
}ReadableOutputSet;
