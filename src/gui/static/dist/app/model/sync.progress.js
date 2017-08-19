System.register([], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var SyncProgress;
    return {
        setters:[],
        execute: function() {
            SyncProgress = (function () {
                function SyncProgress(current, highest) {
                    this.current = current;
                    this.highest = highest;
                }
                return SyncProgress;
            }());
            exports_1("SyncProgress", SyncProgress);
        }
    }
});

//# sourceMappingURL=sync.progress.js.map
