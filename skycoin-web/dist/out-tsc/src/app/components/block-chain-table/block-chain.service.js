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
var http_1 = require('@angular/http');
var Observable_1 = require('rxjs/Observable');
require('rxjs/add/operator/map');
require('rxjs/add/operator/catch');
var BlockChainService = (function () {
    function BlockChainService(_http) {
        this._http = _http;
    }
    BlockChainService.prototype.getBlocks = function (startNumber, endNumber) {
        var stringConvert = 'start=' + startNumber + '&end=' + endNumber;
        return this._http.get('/api/blocks?' + stringConvert)
            .map(function (res) {
            return res.json();
        })
            .map(function (res) { return res.blocks; })
            .catch(function (error) {
            console.log(error);
            return Observable_1.Observable.throw(error || 'Server error');
        });
    };
    BlockChainService.prototype.getBlockByHash = function (hashNumber) {
        var stringConvert = 'hash=' + hashNumber;
        return this._http.get('/api/block?' + stringConvert)
            .map(function (res) {
            console.log(res);
            return res.json();
        })
            .map(function (res) { return res; })
            .catch(function (error) {
            console.log(error);
            return Observable_1.Observable.throw(error || 'Server error');
        });
    };
    BlockChainService = __decorate([
        core_1.Injectable(), 
        __metadata('design:paramtypes', [http_1.Http])
    ], BlockChainService);
    return BlockChainService;
}());
exports.BlockChainService = BlockChainService;
//# sourceMappingURL=block-chain.service.js.map