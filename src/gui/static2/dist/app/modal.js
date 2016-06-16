System.register(['angular2/core', 'angular2/http', 'rxjs/add/operator/map', 'rxjs/add/operator/catch', './ng2-qrcode.ts'], function(exports_1, context_1) {
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
    var core_1, http_1, http_2, ng2_qrcode_ts_1;
    var Dialog;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (http_1_1) {
                http_1 = http_1_1;
                http_2 = http_1_1;
            },
            function (_1) {},
            function (_2) {},
            function (ng2_qrcode_ts_1_1) {
                ng2_qrcode_ts_1 = ng2_qrcode_ts_1_1;
            }],
        execute: function() {
            Dialog = (function () {
                //Constructor method to load HTTP object
                function Dialog(http) {
                    this.http = http;
                }
                //Show QR code function for view QR popup
                Dialog.prototype.showQR = function (address) {
                    this.QrAddress = address.entries[0].address;
                    this.QrIsVisible = true;
                };
                //Hide QR code function for hide QR popup
                Dialog.prototype.hideQr = function () {
                    this.QrIsVisible = false;
                };
                //Show wallet function for view New wallet popup
                Dialog.prototype.showWallet = function () {
                    this.NewWalletIsVisible = true;
                };
                //Hide wallet function for hide New wallet popup
                Dialog.prototype.hideWallet = function () {
                    this.NewWalletIsVisible = false;
                };
                //Show edit wallet function
                Dialog.prototype.showEditWallet = function (wallet) {
                    this.EditWalletIsVisible = true;
                    this.walletId = wallet.meta.filename;
                };
                //Hide edit wallet function
                Dialog.prototype.hideEditWallet = function () {
                    this.EditWalletIsVisible = false;
                };
                //Add new wallet function for generate new wallet in Skycoin
                Dialog.prototype.generateWallet = function () {
                    var _this = this;
                    //Set http headers
                    var headers = new http_2.Headers();
                    headers.append('Content-Type', 'application/x-www-form-urlencoded');
                    //Post method executed
                    this.http.post('/wallet/create', JSON.stringify({ name: '' }), { headers: headers })
                        .map(function (res) { return res.json(); })
                        .subscribe(function (response) {
                        //Load all wallets after creating new wallet
                        _this.http.post('/wallets', '')
                            .map(function (res) { return res.json(); })
                            .subscribe(
                        //Response from API
                        function (data) {
                            _this.NewWalletIsVisible = false;
                            _this.wallets = data;
                        }, function (err) { return console.error("Error on load wallet: " + err); }, function () { return console.log('Wallet load done'); });
                    }, function (err) { return console.error("Error on create new wallet: " + err); }, function () { return console.log('New wallet create done'); });
                };
                Dialog = __decorate([
                    core_1.Component({
                        selector: 'app-dialog',
                        templateUrl: 'app/templates/modal.html',
                        directives: [ng2_qrcode_ts_1.QRCodeComponent],
                    }), 
                    __metadata('design:paramtypes', [http_1.Http])
                ], Dialog);
                return Dialog;
            }());
            exports_1("Dialog", Dialog);
        }
    }
});
//# sourceMappingURL=modal.js.map