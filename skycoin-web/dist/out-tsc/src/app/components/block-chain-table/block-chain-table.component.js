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
var block_chain_service_1 = require("./block-chain.service");
var moment = require('moment');
var router_1 = require("@angular/router");
var BlockChainTableComponent = (function () {
    function BlockChainTableComponent(blockService, router) {
        this.blockService = blockService;
        this.router = router;
        this.loading = false;
        this.blocks = [];
    }
    BlockChainTableComponent.prototype.GetBlockAmount = function (txns) {
        var ret = [];
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
    BlockChainTableComponent.prototype.getTime = function (time) {
        return moment.unix(time).format("YYYY-MM-DD HH:mm");
    };
    BlockChainTableComponent.prototype.ngOnInit = function () {
    };
    BlockChainTableComponent.prototype.showDetails = function (block) {
        this.router.navigate(['/block', block.header.seq]);
    };
    BlockChainTableComponent.prototype.handlePageChange = function (pagesData) {
        var _this = this;
        this.totalBlocks = pagesData[1];
        var currentPage = pagesData[0];
        var blockStart = this.totalBlocks - currentPage * 10 + 1;
        var blockEnd = blockStart + 9;
        if (blockEnd >= this.totalBlocks) {
            blockEnd = this.totalBlocks;
        }
        if (blockStart <= 1) {
            blockStart = 1;
        }
        this.loading = true;
        this.blockService.getBlocks(blockStart, blockEnd).subscribe(function (data) {
            var newData = _.sortBy(data, function (block) { return block.header.seq; }).reverse();
            _this.blocks = newData;
            _this.loading = false;
        }, function (err) {
            _this.loading = false;
        });
    };
    BlockChainTableComponent = __decorate([
        core_1.Component({
            selector: 'app-block-chain-table',
            templateUrl: './block-chain-table.component.html',
            styleUrls: ['./block-chain-table.component.css']
        }), 
        __metadata('design:paramtypes', [block_chain_service_1.BlockChainService, router_1.Router])
    ], BlockChainTableComponent);
    return BlockChainTableComponent;
}());
exports.BlockChainTableComponent = BlockChainTableComponent;
//# sourceMappingURL=block-chain-table.component.js.map