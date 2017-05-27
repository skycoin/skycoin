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
var http_1 = require("@angular/http");
var rxjs_1 = require("rxjs");
var UxOutputsService = (function () {
    function UxOutputsService(_http) {
        this._http = _http;
    }
    UxOutputsService.prototype.getUxOutputsForAddress = function (address) {
        return this._http.get('/api/address?address=' + address)
            .map(function (res) {
            return res.json();
        })
            .catch(function (error) {
            console.log(error);
            return rxjs_1.Observable.throw(error || 'Server error');
        });
    };
    UxOutputsService.prototype.getCurrentBalanceOfAddress = function (address) {
        return this._http.get('/api/currentBalance?address=' + address)
            .map(function (res) {
            return res.json();
        })
            .catch(function (error) {
            console.log(error);
            return rxjs_1.Observable.throw(error || 'Server error');
        });
    };
    UxOutputsService.prototype.getAddressFromUxId = function (uxid) {
        return this._http.get('/api/uxout?uxid=' + uxid)
            .map(function (res) {
            return res.json();
        })
            .map(function (res) {
            return res.owner_address;
        })
            .catch(function (error) {
            console.log(error);
            return rxjs_1.Observable.throw(error || 'Server error');
        });
    };
    UxOutputsService.prototype.getBlockSource = function (blockNumber) {
        return this._http.get('/api/blocks?start=' + blockNumber + '&end=' + blockNumber)
            .map(function (res) {
            return res.json().blocks[0];
        })
            .catch(function (error) {
            console.log(error);
            return rxjs_1.Observable.throw(error || 'Server error');
        });
    };
    UxOutputsService = __decorate([
        core_1.Injectable(), 
        __metadata('design:paramtypes', [http_1.Http])
    ], UxOutputsService);
    return UxOutputsService;
}());
exports.UxOutputsService = UxOutputsService;
//# sourceMappingURL=UxOutputs.service.js.map