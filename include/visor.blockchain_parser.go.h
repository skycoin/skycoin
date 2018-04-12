typedef Handle visor__ParserOption;
typedef struct{
    visor__historyer historyDB;
    GoChan_ blkC;
    GoChan_ quit;
    GoChan_ done;
    visor__Blockchainer bc;
    bool isStart;
} visor__BlockchainParser;
