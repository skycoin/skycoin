typedef struct{
    cipher__PubKey Pubkey;
    cipher__SecKey Seckey;
    consensus__ConnectionManagerInterface pConnectionManager;
    consensus__BlockchainTail block_queue;
    consensus__BlockStatQueue block_stat_queue;
    GoInt_ Incoming_block_count;
} consensus__ConsensusParticipant;
