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
var UxOutputs_service_1 = require("./UxOutputs.service");
var router_1 = require("@angular/router");
var AddressDetailComponent = (function () {
    function AddressDetailComponent(service, route, router) {
        this.service = service;
        this.route = route;
        this.router = router;
        this.UxOutputs = null;
        this.currentBalance = 0;
        this.transactions = [];
        this.currentAddress = null;
        this.showUxID = false;
    }
    AddressDetailComponent.prototype.ngOnInit = function () {
    };
    AddressDetailComponent.prototype.ngAfterViewInit = function () {
        var _this = this;
        this.UxOutputs = this.route.params
            .switchMap(function (params) {
            var address = params['address'];
            _this.currentAddress = address;
            var qrcode = new QRCode("qr-code");
            qrcode.makeCode(_this.currentAddress);
            return _this.service.getUxOutputsForAddress(address);
        });
        this.UxOutputs.subscribe(function (uxoutputs) {
            _this.transactions = uxoutputs;
            console.log(uxoutputs);
        });
        this.route.params
            .switchMap(function (params) {
            var address = params['address'];
            return _this.service.getCurrentBalanceOfAddress(address);
        }).subscribe(function (addressDetails) {
            if (addressDetails.head_outputs.length > 0) {
                for (var i = 0; i < addressDetails.head_outputs.length; i++) {
                    _this.currentBalance = _this.currentBalance + parseInt(addressDetails.head_outputs[i].coins);
                }
            }
        });
    };
    AddressDetailComponent.prototype.showUxId = function () {
        this.showUxID = true;
        return false;
    };
    AddressDetailComponent.prototype.hideUxId = function () {
        this.showUxID = false;
        return false;
    };
    AddressDetailComponent = __decorate([
        core_1.Component({
            selector: 'app-address-detail',
            templateUrl: './address-detail.component.html',
            styleUrls: ['./address-detail.component.css'],
        }), 
        __metadata('design:paramtypes', [UxOutputs_service_1.UxOutputsService, router_1.ActivatedRoute, router_1.Router])
    ], AddressDetailComponent);
    return AddressDetailComponent;
}());
exports.AddressDetailComponent = AddressDetailComponent;
//# sourceMappingURL=address-detail.component.js.map