"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var core_1 = require('@angular/core');
var router_1 = require("@angular/router");
var rxjs_1 = require("rxjs");
require('rxjs/add/observable/forkJoin');
var transaction_detail_service_1 = require("./transaction-detail.service");
var moment = require('moment');
var TransactionDetailComponent = (function () {
    function TransactionDetailComponent(service, route, router) {
        this.service = service;
        this.route = route;
        this.router = router;
        this.transactionObservable = null;
        this.transaction = null;
    }
    TransactionDetailComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.transactionObservable = this.route.params
            .flatMap(function (params) {
            var txid = params['txid'];
            return _this.service.getTransaction(txid);
        })
            .flatMap(function (trans) {
            var tasks$ = [];
            _this.transaction = trans.txn;
            _this.transaction.status = trans.status.confirmed;
            _this.transaction.block_num = trans.status.block_seq;
            trans = trans.txn;
            for (var i = 0; i < trans.inputs.length; i++) {
                tasks$.push(_this.getAddressOfInput(trans.inputs[i]));
            }
            return rxjs_1.Observable.forkJoin.apply(rxjs_1.Observable, tasks$);
        });
        this.transactionObservable.subscribe(function (trans) {
            for (var i = 0; i < trans.length; i++) {
                _this.transaction.inputs[i] = trans[i].owner_address;
            }
        });
    };
    TransactionDetailComponent.prototype.getAddressOfInput = function (uxid) {
        return this.service.getInputAddress(uxid);
    };
    TransactionDetailComponent.prototype.getTime = function (time) {
        return moment.unix(time).format();
    };
    TransactionDetailComponent = __decorate([
        core_1.Component({
            selector: 'app-transaction-detail',
            templateUrl: './transaction-detail.component.html',
            styleUrls: ['./transaction-detail.component.css']
        }), 
        __metadata('design:paramtypes', [transaction_detail_service_1.TransactionDetailService, router_1.ActivatedRoute, router_1.Router])
    ], TransactionDetailComponent);
    return TransactionDetailComponent;
}());
exports.TransactionDetailComponent = TransactionDetailComponent;
//# sourceMappingURL=transaction-detail.component.js.map