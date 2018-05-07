typedef struct{
    GoInt_ priority;
    GoInt_ index;
    GoMap_ hash2info;
    GoUint64_ seqno;
    BOOL frozen;
    GoInt_ accept_count;
    GoMap_ debug_pubkey2count;
    GoInt_ debug_count;
    GoInt_ debug_reject_count;
    GoInt_ debug_neglect_count;
    GoInt_ debug_usage;
} consensus__BlockStat;
typedef GoSlice_  consensus__PriorityQueue;
typedef struct{
    consensus__PriorityQueue queue;
} consensus__BlockStatQueue;
