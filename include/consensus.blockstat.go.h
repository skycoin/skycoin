typedef struct{
    GoInt_ priority;
    GoInt_ index;
    GoMap_ hash2info;
    GoUint64_ seqno;
    bool frozen;
    GoInt_ accept_count;
    GoMap_ debug_pubkey2count;
    GoInt_ debug_count;
    GoInt_ debug_reject_count;
    GoInt_ debug_neglect_count;
    GoInt_ debug_usage;
} BlockStat;
typedef GoSlice_ PriorityQueue;
typedef struct{
    PriorityQueue queue;
} BlockStatQueue;
