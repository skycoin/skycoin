typedef struct{
    visor__BlockchainMetadata * _unnamed;
    wh__Duration TimeSinceLastBlock;
} api__BlockchainMetadata;
typedef struct{
    api__BlockchainMetadata BlockchainMetadata;
    visor__BuildInfo Version;
    GoInt_ OpenConnections;
    wh__Duration Uptime;
} api__HealthResponse;
