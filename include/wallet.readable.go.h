typedef struct{
    GoString_ Address;
    GoString_ Public;
    GoString_ Secret;
} wallet__ReadableEntry;
typedef GoSlice_  wallet__ReadableEntries;
typedef struct{
    GoMap_ Meta;
    wallet__ReadableEntries Entries;
} wallet__ReadableWallet;
