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
var block_chain_service_1 = require("../block-chain-table/block-chain.service");
var router_1 = require("@angular/router");
var SkycoinSearchBarComponent = (function () {
    function SkycoinSearchBarComponent(service, router) {
        this.service = service;
        this.router = router;
    }
    SkycoinSearchBarComponent.prototype.ngOnInit = function () {
    };
    SkycoinSearchBarComponent.prototype.searchBlockHistory = function (hashVal) {
        if (hashVal.length == 34 || hashVal.length == 35) {
            this.router.navigate(['/address', hashVal]);
            return;
        }
        if (hashVal.length == 64) {
            this.router.navigate(['/transaction', hashVal]);
            return;
        }
        this.router.navigate(['/block', hashVal]);
        return;
        // this.service.getBlockByHash(hashVal).subscribe((block)=>{
        //   this.block = block;
        //   this.router.navigate(['/block', block.header.seq]);
        //
        // });
    };
    SkycoinSearchBarComponent = __decorate([
        core_1.Component({
            selector: 'app-skycoin-search-bar',
            templateUrl: './skycoin-search-bar.component.html',
            styleUrls: ['./skycoin-search-bar.component.css']
        }), 
        __metadata('design:paramtypes', [block_chain_service_1.BlockChainService, router_1.Router])
    ], SkycoinSearchBarComponent);
    return SkycoinSearchBarComponent;
}());
exports.SkycoinSearchBarComponent = SkycoinSearchBarComponent;
//# sourceMappingURL=skycoin-search-bar.component.js.map