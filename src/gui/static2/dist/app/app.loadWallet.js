System.register(['angular2/core', 'angular2/router', 'angular2/http', 'rxjs/add/operator/map', 'rxjs/add/operator/catch', "./modal.ts"], function(exports_1, context_1) {
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
    var core_1, router_1, http_1, http_2, modal_ts_1;
    var loadWalletComponent, DisplayModeEnum;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (router_1_1) {
                router_1 = router_1_1;
            },
            function (http_1_1) {
                http_1 = http_1_1;
                http_2 = http_1_1;
            },
            function (_1) {},
            function (_2) {},
            function (modal_ts_1_1) {
                modal_ts_1 = modal_ts_1_1;
            }],
        execute: function() {
            loadWalletComponent = (function () {
                function loadWalletComponent(http) {
                    this.http = http;
                    this.displayModeEnum = DisplayModeEnum;
                }
                loadWalletComponent.prototype.ngOnInit = function () {
                    this.displayMode = DisplayModeEnum.first;
                    this.loadWallet();
                    this.loadProgress();
                    setInterval(function () {
                        //this.loadWallet();
                        console.log("Refreshing balance");
                    }, 15000);
                };
                loadWalletComponent.prototype.loadWallet = function () {
                    var _this = this;
                    this.http.post('/wallets', '')
                        .map(function (res) { return res.json(); })
                        .subscribe(function (data) {
                        _this.wallets = data;
                        //Load Balance for each wallet
                        var headers = new http_2.Headers();
                        headers.append('Content-Type', 'application/x-www-form-urlencoded');
                        var inc = 0;
                        for (var item in data) {
                            var address = data[inc].meta.filename;
                            _this.http.post('/wallet/balance', JSON.stringify({ id: address }), { headers: headers })
                                .map(function (res) { return res.json(); })
                                .subscribe(function (response) {
                                console.log('load done: ' + inc);
                                _this.wallets[inc].balance = response.confirmed.coins / 1000000;
                                inc++;
                            }, function (err) { return console.error("Error on load balance: " + err); }, function () { return console.log('Balance load done'); });
                        }
                        //Load Balance for each wallet end
                    }, function (err) { return console.error("Error on load wallet: " + err); }, function () { return console.log('Wallet load done'); });
                };
                loadWalletComponent.prototype.loadProgress = function () {
                    var _this = this;
                    this.http.post('/blockchain/progress', '')
                        .map(function (res) { return res.json(); })
                        .subscribe(function (response) { _this.progress = (parseInt(response.current, 10) + 1) / parseInt(response.Highest, 10) * 100; }, function (err) { return console.error("Error on load progress: " + err); }, function () { return console.log('Progress load done'); });
                };
                loadWalletComponent.prototype.switchTab = function (mode) {
                    this.displayMode = mode;
                };
                loadWalletComponent.prototype.showQR = function (wallet) {
                    this._dialog.showQR(wallet);
                };
                loadWalletComponent.prototype.showNewWalletDialog = function (wallet) {
                    this._dialog.showWallet();
                };
                __decorate([
                    core_1.ViewChild(modal_ts_1.Dialog), 
                    __metadata('design:type', modal_ts_1.Dialog)
                ], loadWalletComponent.prototype, "_dialog", void 0);
                loadWalletComponent = __decorate([
                    core_1.Component({
                        selector: 'load-wallet',
                        directives: [router_1.ROUTER_DIRECTIVES, modal_ts_1.Dialog],
                        providers: [],
                        templateUrl: 'app/templates/wallet.html'
                    }), 
                    __metadata('design:paramtypes', [http_1.Http])
                ], loadWalletComponent);
                return loadWalletComponent;
            }());
            exports_1("loadWalletComponent", loadWalletComponent);
            (function (DisplayModeEnum) {
                DisplayModeEnum[DisplayModeEnum["first"] = 0] = "first";
                DisplayModeEnum[DisplayModeEnum["second"] = 1] = "second";
                DisplayModeEnum[DisplayModeEnum["third"] = 2] = "third";
            })(DisplayModeEnum || (DisplayModeEnum = {}));
        }
    }
});
//# sourceMappingURL=app.loadWallet.js.map