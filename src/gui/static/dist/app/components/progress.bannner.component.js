System.register(["@angular/core", "../services/skycoin.sync.service"], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
        var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
        if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
        else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
        return c > 3 && r && Object.defineProperty(target, key, r), r;
    };
    var __metadata = (this && this.__metadata) || function (k, v) {
        if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
    };
    var core_1, skycoin_sync_service_1;
    var SkycoinSyncWalletBlock;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (skycoin_sync_service_1_1) {
                skycoin_sync_service_1 = skycoin_sync_service_1_1;
            }],
        execute: function() {
            SkycoinSyncWalletBlock = (function () {
                function SkycoinSyncWalletBlock(_syncService) {
                    this._syncService = _syncService;
                    this.currentWalletNumber = this.highestWalletNumber = 0;
                    this.syncDone = false;
                }
                SkycoinSyncWalletBlock.prototype.ngAfterViewInit = function () {
                    var _this = this;
                    this.handlerSync = setInterval(function () {
                        if (_this.highestWalletNumber - _this.currentWalletNumber <= 1 && _this.highestWalletNumber != 0) {
                            clearInterval(_this.handlerSync);
                            _this.syncDone = true;
                        }
                        _this.syncBlocks();
                    }, 2000);
                };
                SkycoinSyncWalletBlock.prototype.syncBlocks = function () {
                    var _this = this;
                    this._syncProgress = this._syncService.getSyncProgress();
                    this._syncProgress.subscribe(function (syncProgress) {
                        _this.currentWalletNumber = syncProgress.current;
                        _this.highestWalletNumber = syncProgress.highest;
                    });
                };
                SkycoinSyncWalletBlock = __decorate([
                    core_1.Component({
                        selector: 'skycoin-block-sync',
                        template: "\n\n        \n             <div class=\"sync-div-container\">\n             \n             <ul class=\"fa-ul\">\n  <li><i class=\"fa-li fa fa-spinner fa-spin\" *ngIf=\"syncDone == false\"></i>\n  <span *ngIf=\"currentWalletNumber>0\">{{currentWalletNumber}} of {{highestWalletNumber}} blocks synced</span>\n  <span *ngIf=\"currentWalletNumber==0\">Syncing wallet</span>\n  </li>\n</ul>\n                \n               \n              </div>\n            \n              ",
                        providers: [skycoin_sync_service_1.BlockSyncService]
                    }), 
                    __metadata('design:paramtypes', [skycoin_sync_service_1.BlockSyncService])
                ], SkycoinSyncWalletBlock);
                return SkycoinSyncWalletBlock;
            }());
            exports_1("SkycoinSyncWalletBlock", SkycoinSyncWalletBlock);
        }
    }
});

//# sourceMappingURL=progress.bannner.component.js.map
