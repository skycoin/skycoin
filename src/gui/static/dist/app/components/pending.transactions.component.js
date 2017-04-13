System.register(["@angular/core", "../services/pending.transaction.service"], function(exports_1, context_1) {
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
    var core_1, pending_transaction_service_1;
    var PendingTxnsComponent;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (pending_transaction_service_1_1) {
                pending_transaction_service_1 = pending_transaction_service_1_1;
            }],
        execute: function() {
            PendingTxnsComponent = (function () {
                function PendingTxnsComponent(_pendingTxnsService) {
                    this._pendingTxnsService = _pendingTxnsService;
                }
                PendingTxnsComponent.prototype.ngAfterViewInit = function () {
                    this.refreshPendingTxns();
                };
                PendingTxnsComponent.prototype.refreshPendingTxns = function () {
                    var _this = this;
                    this._pendingTxnsService.getPendingTransactions().subscribe(function (pendingTxns) {
                        _this.transactions = pendingTxns;
                    }, function (err) {
                        console.log(err);
                    });
                };
                PendingTxnsComponent.prototype.getDateTimeString = function (ts) {
                    return moment(ts).format("YYYY-MM-DD HH:mm");
                };
                PendingTxnsComponent.prototype.sendExecution = function () {
                    var _this = this;
                    this._pendingTxnsService.resendPendingTxns().subscribe(function () {
                        _this.refreshPendingTxns();
                    });
                };
                PendingTxnsComponent = __decorate([
                    core_1.Component({
                        selector: 'pending-transactions',
                        template: "\n<button class=\"btn btn-default right\" type=\"button\" (click)=\"sendExecution()\" >Resend for execution</button>\n<div class=\"table-responsive\">\n                  <table id=\"pending-table\" class=\"table\">\n                            <tbody>\n                            <tr class=\"dark-row\">\n                                <td>S. No</td>\n                                <td>Time received</td>\n                                <td>Transaction ID</td>\n                                <td>Inputs</td>\n                                <td>Outputs</td>\n                                <td>Amount</td>\n                            </tr>\n                            <tr *ngFor=\"let transaction of transactions;let i=index\">\n                                <td>{{i+1}}</td>\n                                <td>{{getDateTimeString(transaction.received)}}</td>\n                                <td>{{transaction.transaction.txid}}</td>\n                                <td>\n                                <p *ngFor=\"let input of transaction.transaction.inputs\">{{input}},<br></p>\n</td>\n                                <td>\n                                <p *ngFor=\"let output of transaction.transaction.outputs\">{{output.dst}},<br></p>\n</td>\n                                <td><p *ngFor=\"let output of transaction.transaction.outputs\">{{output.coins}},<br></p></td>\n                            </tr>\n                            </tbody>\n                        </table>\n                        </div>\n              ",
                        providers: [pending_transaction_service_1.PendingTransactionService]
                    }), 
                    __metadata('design:paramtypes', [pending_transaction_service_1.PendingTransactionService])
                ], PendingTxnsComponent);
                return PendingTxnsComponent;
            }());
            exports_1("PendingTxnsComponent", PendingTxnsComponent);
        }
    }
});

//# sourceMappingURL=pending.transactions.component.js.map
