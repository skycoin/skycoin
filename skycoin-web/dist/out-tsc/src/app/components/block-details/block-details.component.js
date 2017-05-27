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
var block_chain_service_1 = require("../block-chain-table/block-chain.service");
var moment = require('moment');
var BlockDetailsComponent = (function () {
    function BlockDetailsComponent(service, route, router) {
        this.service = service;
        this.route = route;
        this.router = router;
        this.block = null;
    }
    BlockDetailsComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.blocksObservable = this.route.params
            .switchMap(function (params) {
            var selectedBlock = +params['id'];
            return _this.service.getBlocks(selectedBlock, selectedBlock);
        });
        this.blocksObservable.subscribe(function (blocks) {
            _this.block = blocks[0];
        });
    };
    BlockDetailsComponent.prototype.getTime = function (time) {
        return moment.unix(time).format();
    };
    BlockDetailsComponent.prototype.getAmount = function (block) {
        var ret = [];
        var txns = block.body.txns;
        _.each(txns, function (o) {
            if (o.outputs) {
                _.each(o.outputs, function (_o) {
                    ret.push(_o.coins);
                });
            }
        });
        var totalCoins = ret.reduce(function (memo, coin) {
            return memo + parseInt(coin);
        }, 0);
        return totalCoins;
    };
    BlockDetailsComponent = __decorate([
        core_1.Component({
            selector: 'app-block-details',
            templateUrl: './block-details.component.html',
            styleUrls: ['./block-details.component.css']
        }), 
        __metadata('design:paramtypes', [block_chain_service_1.BlockChainService, router_1.ActivatedRoute, router_1.Router])
    ], BlockDetailsComponent);
    return BlockDetailsComponent;
}());
exports.BlockDetailsComponent = BlockDetailsComponent;
//# sourceMappingURL=block-details.component.js.map