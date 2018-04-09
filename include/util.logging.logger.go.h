typedef GoInterface_ ExtendedFieldLogger;
typedef struct{
;
    GoString_ module;
    bool allModulesEnabled;
    GoMap_ moduleLoggers;
    GoString_ PriorityKey;
    GoString_ CriticalPriority;
} Logger;
