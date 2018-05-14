typedef struct{
    GoString_ Address;
    GoUint64_ Height;
} daemon__PeerBlockchainHeight;
typedef struct{
    GoUint64_ LastBlock;
    GoUint64_ RequestedBlocks;
    gnet__MessageContext * c;
} daemon__GetBlocksMessage;
typedef struct{
    GoSlice_  Blocks;
    gnet__MessageContext * c;
} daemon__GiveBlocksMessage;
typedef struct{
    GoUint64_ MaxBkSeq;
    gnet__MessageContext * c;
} daemon__AnnounceBlocksMessage;
typedef struct{
    GoSlice_  Txns;
    gnet__MessageContext * c;
} daemon__AnnounceTxnsMessage;
typedef struct{
    GoSlice_  Txns;
    gnet__MessageContext * c;
} daemon__GetTxnsMessage;
typedef struct{
    coin__Transactions Txns;
    gnet__MessageContext * c;
} daemon__GiveTxnsMessage;
