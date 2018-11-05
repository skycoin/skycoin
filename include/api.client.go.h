typedef struct{
    GoInt_ N;
    BOOL IncludeDistribution;
} api__RichlistParams;

typedef struct{
    // Comma separated list of connection states ("pending", "connected", "introduced")
    GoString_ States;
    GoString_ Direction;
} api__NetworkConnectionsFilter;

