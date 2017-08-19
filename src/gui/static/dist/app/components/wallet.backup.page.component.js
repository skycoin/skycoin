System.register(["@angular/core", "../services/wallet.service"], function(exports_1, context_1) {
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
    var core_1, wallet_service_1;
    var WalletBackupPageComponent;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (wallet_service_1_1) {
                wallet_service_1 = wallet_service_1_1;
            }],
        execute: function() {
            WalletBackupPageComponent = (function () {
                function WalletBackupPageComponent(_service) {
                    this._service = _service;
                }
                WalletBackupPageComponent.prototype.ngAfterViewInit = function () {
                    this.getWalletFolder();
                    this.walletFolder = "";
                };
                WalletBackupPageComponent.prototype.ngOnDestroy = function () {
                    this.wallets.forEach(function (el) {
                        el.showSeed = false;
                    });
                };
                WalletBackupPageComponent.prototype.getWalletFolder = function () {
                    var _this = this;
                    this._service.getWalletFolder().subscribe(function (walletFolder) {
                        _this.walletFolder = walletFolder.address;
                    }, function (err) {
                        console.log(err);
                    });
                };
                WalletBackupPageComponent.prototype.download = function (ev, wallet) {
                    ev.stopImmediatePropagation();
                    ev.stopPropagation();
                    var blob = new Blob([JSON.stringify({ "seed": wallet.meta.seed })], { type: 'application/json' });
                    var link = document.createElement('a');
                    link.href = window.URL.createObjectURL(blob);
                    link['download'] = wallet.meta.filename + '.json';
                    link.click();
                };
                WalletBackupPageComponent.prototype.showOrHideSeed = function (wallet) {
                    wallet.showSeed = !wallet.showSeed;
                };
                __decorate([
                    core_1.Input(), 
                    __metadata('design:type', Array)
                ], WalletBackupPageComponent.prototype, "wallets", void 0);
                WalletBackupPageComponent = __decorate([
                    core_1.Component({
                        selector: 'backup-wallets',
                        template: "\n<h2>Wallet Backup</h2>\n<p> Wallet Directory: <b>{{walletFolder}}</b> </p>\n<p>\n    <b>BACKUP YOUR SEED. ON PAPER. IN A SAFE PLACE.</b> As long as you have your seed, you can recover your coins.\n</p>\n<div class=\"table-responsive\">\n                  <table id=\"wallet-table\" class=\"table\">\n                  <thead>\n                    <tr class=\"dark-row\">\n                                <td>S. No</td>\n                                <td>Wallet Label</td>\n                                <td>File Name</td>\n                                <td>Download</td>\n                                <td>Seed</td>\n          \n                            </tr>\n</thead>\n                            <tbody>\n                      \n                            <tr *ngFor=\"let wallet of wallets;let i=index\">\n                                <td>{{i+1}}</td>\n                                <td>{{wallet.meta.label}}</td>\n                                <td>{{wallet.meta.filename}}</td>\n\n                                <td><a class=\"btn btn-success\" href=\"javascript:void(0);\" (click)=\"download($event,wallet)\">{{wallet.meta.filename}}</a></td>\n                                 <td>\n                                  <a class=\"btn btn-default\" *ngIf=\"!wallet?.showSeed\"  (click)=\"showOrHideSeed(wallet)\">Show Seed</a>\n                                  <p *ngIf=\"wallet?.showSeed\">{{wallet.meta.seed}}<a class=\"btn btn-default btn-margin\" (click)=\"showOrHideSeed(wallet)\">Hide Seed</a></p>\n                                 </td>\n                            </tr>\n                            </tbody>\n                        </table>\n                        </div>\n              ",
                        styles: ["\n    .btn-margin {\n      margin: 0 1rem;\n    }\n  "],
                        providers: [wallet_service_1.WalletService]
                    }), 
                    __metadata('design:paramtypes', [wallet_service_1.WalletService])
                ], WalletBackupPageComponent);
                return WalletBackupPageComponent;
            }());
            exports_1("WalletBackupPageComponent", WalletBackupPageComponent);
        }
    }
});

//# sourceMappingURL=wallet.backup.page.component.js.map
