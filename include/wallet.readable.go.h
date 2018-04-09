typedef struct{
    GoString_ Address;
    GoString_ Public;
    GoString_ Secret;
} ReadableEntry;
typedef GoSlice_ ReadableEntries;
typedef struct{
    GoMap_ Meta;
    ReadableEntries Entries;
} ReadableWallet;
