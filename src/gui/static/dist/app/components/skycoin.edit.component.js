/**
 * Created by nakul.pandey@gmail.com on 01/01/17.
 */
System.register(['@angular/core', '../services/wallet.service'], function(exports_1, context_1) {
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
    var SkyCoinEditComponent;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (wallet_service_1_1) {
                wallet_service_1 = wallet_service_1_1;
            }],
        execute: function() {
            SkyCoinEditComponent = (function () {
                function SkyCoinEditComponent(el, _walletService) {
                    this._walletService = _walletService;
                    this.onWalletChanged = new core_1.EventEmitter();
                    this.show = false;
                    this.m = 3;
                    this.min = 0;
                    this.max = 100;
                    this.invalid = false;
                    this.el = el;
                }
                SkyCoinEditComponent.prototype.ngOnInit = function () {
                    this.originalText = this.text;
                };
                SkyCoinEditComponent.prototype.validate = function (text) {
                    if (this.regex) {
                        var re = new RegExp('' + this.regex, "ig");
                        if (re.test(text)) {
                            this.invalid = false;
                        }
                        else {
                            this.invalid = true;
                        }
                    }
                    else {
                        if ((text.length <= this.max) && (text.length >= this.min)) {
                            this.invalid = false;
                        }
                        else {
                            this.invalid = true;
                        }
                    }
                };
                SkyCoinEditComponent.prototype.makeEditable = function () {
                    if (this.show == false) {
                        this.show = true;
                    }
                };
                SkyCoinEditComponent.prototype.compareEvent = function (globalEvent) {
                    if (this.tracker != globalEvent && this.show) {
                        this.cancelEditable();
                    }
                };
                SkyCoinEditComponent.prototype.trackEvent = function (newHostEvent) {
                    this.tracker = newHostEvent;
                };
                SkyCoinEditComponent.prototype.cancelEditable = function () {
                    this.show = false;
                    this.invalid = false;
                    this.text = this.originalText;
                };
                SkyCoinEditComponent.prototype.callSave = function () {
                    var _this = this;
                    if (!this.invalid && !this.isDuplicate(this.text)) {
                        var data = {};
                        data["newText"] = this.text;
                        data["walletId"] = this.walletId;
                        this._walletService.updateWallet(data).subscribe(function (response) {
                            _this.onWalletChanged.emit(data);
                            _this.show = false;
                            toastr.info("Wallet name updated");
                        }, function (err) {
                            _this.cancelEditable();
                            toastr.error("Unable to update the name. Please try after some time");
                        });
                    }
                };
                SkyCoinEditComponent.prototype.isDuplicate = function (text) {
                    var old = _.find(this.wallets, function (o) {
                        return (o.meta.label == text);
                    });
                    if (old) {
                        toastr.error("This wallet label is used already");
                        this.cancelEditable();
                        return true;
                    }
                    return false;
                };
                __decorate([
                    core_1.Input('text'), 
                    __metadata('design:type', Object)
                ], SkyCoinEditComponent.prototype, "text", void 0);
                __decorate([
                    core_1.Input('wallets'), 
                    __metadata('design:type', Array)
                ], SkyCoinEditComponent.prototype, "wallets", void 0);
                __decorate([
                    core_1.Input('walletId'), 
                    __metadata('design:type', Object)
                ], SkyCoinEditComponent.prototype, "walletId", void 0);
                __decorate([
                    core_1.Output(), 
                    __metadata('design:type', Object)
                ], SkyCoinEditComponent.prototype, "onWalletChanged", void 0);
                SkyCoinEditComponent = __decorate([
                    core_1.Component({
                        selector: 'skycoin-edit',
                        styles: [
                            " #skycoin-edit-ic {\n        margin-left: 10px;\n        color: #d9d9d9;\n        }\n        .skycoin-edit-comp {\n            padding:6px;\n            border-radius: 3px;\n        }\n        .active-skycoin-edit {\n            background-color: #f0f0f0;\n            border: 1px solid #d9d9d9;\n        }\n        input {\n            border-radius: 5px;\n            box-shadow: none;\n            border: 1px solid #dedede;\n            min-width: 5px;\n        }\n        .skycoin-edit-buttons {\n            background-color: #f0f0f0;\n            border: 1px solid #ccc;\n            border-top: none;\n            border-radius: 0 0 3px 3px;\n            box-shadow: 0 3px 6px rgba(111,111,111,0.2);\n            outline: none;\n            padding: 3px;\n            position: absolute;\n            margin-left: 6px;\n            z-index: 1;\n        }\n        .skycoin-edit-comp:hover {\n            border: 1px solid grey;\n        }\n        .skycoin-edit-comp:hover > skycoin-edit-ic {\n            display:block;\n        }\n        .skycoin-edit-save {\n            margin-right:3px;\n        }\n        .skycoin-edit-active {\n            background-color: #f0f0f0;\n            border: 1px solid #d9d9d9;\n        }\n        .ng-invalid {\n                background: #ffb8b8;\n            }\n        .err-bubble {\n            position: absolute;\n            margin: 16px 100px;\n            border: 1px solid red;\n            font-size: 14px;\n            background: #ffb8b8;\n            padding: 10px;\n            border-radius: 7px;\n        }\n       "
                        ],
                        template: "<span class='skycoin-edit-comp' [ngClass]=\"{'skycoin-edit-active':show}\">\n<input *ngIf='show' [ngClass]=\"{'ng-invalid': invalid}\" (ngModelChange)=\"validate($event)\" type='text' [(ngModel)]='text' />\n<div class='err-bubble' *ngIf=\"invalid\">{{error || \" must contain \" + min + \" to -\" + max +\" chars.\"}}</div>\n<i class=\"fa fa-edit\" (click)='makeEditable()' id='skycoin-edit-ic' *ngIf='!show'></i>\n<span *ngIf='!show' (click)='makeEditable()'>{{text || '-Empty Field-'}}</span>\n</span>\n<div class='skycoin-edit-buttons' *ngIf='show'>\n<button class='btn-x-sm' (click)='callSave()'><i class=\"fa fa-check\"></i></button>\n<button class='btn-x-sm' (click)='cancelEditable()'><i class=\"fa fa-times\"></i></button>\n</div>",
                        host: {
                            "(document: click)": "compareEvent($event)",
                            "(click)": "trackEvent($event)"
                        },
                        providers: [wallet_service_1.WalletService]
                    }), 
                    __metadata('design:paramtypes', [core_1.ElementRef, wallet_service_1.WalletService])
                ], SkyCoinEditComponent);
                return SkyCoinEditComponent;
            }());
            exports_1("SkyCoinEditComponent", SkyCoinEditComponent);
        }
    }
});

//# sourceMappingURL=skycoin.edit.component.js.map
