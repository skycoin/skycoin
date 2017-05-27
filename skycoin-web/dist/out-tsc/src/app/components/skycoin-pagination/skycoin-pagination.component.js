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
var skycoin_blockchain_pagination_service_1 = require("./skycoin-blockchain-pagination.service");
var SkycoinPaginationComponent = (function () {
    function SkycoinPaginationComponent(paginationService) {
        this.paginationService = paginationService;
        this.onChangePage = new core_1.EventEmitter();
        this.numberOfBlocks = 0;
        this.currentPage = 1;
        this.currentPages = [];
        this.pagesToShowAtATime = 5;
        this.pageStartPointer = this.currentPage;
        this.pageEndPointer = this.currentPage;
        this.noUpcoming = false;
    }
    SkycoinPaginationComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.paginationService.fetchNumberOfBlocks().subscribe(function (numberOfBlocks) {
            _this.numberOfBlocks = numberOfBlocks;
            _this.onChangePage.emit([1, _this.numberOfBlocks]);
            _this.pagesToShowAtATime = _this.pagesToShowAtATime < numberOfBlocks ? _this.pagesToShowAtATime : _this.numberOfBlocks;
            _this.currentPages = [];
            for (var i = _this.currentPage; i < _this.currentPage + _this.pagesToShowAtATime; i++) {
                _this.currentPages.push(i);
            }
        });
    };
    SkycoinPaginationComponent.prototype.changePage = function (pageNumber) {
        this.onChangePage.emit([pageNumber, this.numberOfBlocks]);
        this.currentPage = pageNumber;
        return false;
    };
    SkycoinPaginationComponent.prototype.loadUpcoming = function () {
        if (this.currentPages[0] * 10 + this.pagesToShowAtATime * 10 >= this.numberOfBlocks) {
            this.noUpcoming = true;
            return false;
        }
        this.onChangePage.emit([this.currentPages[0] + this.pagesToShowAtATime, this.numberOfBlocks]);
        this.currentPage = this.currentPages[0] + this.pagesToShowAtATime;
        this.currentPages = [];
        for (var i = this.currentPage; i < this.currentPage + this.pagesToShowAtATime && i <= this.numberOfBlocks; i++) {
            if (i * 10 - this.numberOfBlocks < 10) {
                this.currentPages.push(i);
            }
            else {
                this.noUpcoming = true;
            }
        }
        return false;
    };
    SkycoinPaginationComponent.prototype.loadPrevious = function () {
        this.noUpcoming = false;
        if (this.currentPages[0] <= 1) {
            return false;
        }
        if (this.currentPages[0] - this.pagesToShowAtATime <= 0) {
            this.currentPages = [];
            this.currentPage = 1;
            this.onChangePage.emit([1, this.numberOfBlocks]);
            for (var i = this.currentPage; i < this.currentPage + this.pagesToShowAtATime; i++) {
                this.currentPages.push(i);
            }
        }
        else {
            this.onChangePage.emit([this.currentPages[0] - this.pagesToShowAtATime, this.numberOfBlocks]);
            this.currentPage = this.currentPages[0] - this.pagesToShowAtATime;
            this.currentPages = [];
            for (var i = this.currentPage; i < this.currentPage + this.pagesToShowAtATime; i++) {
                if (i * 10 <= this.numberOfBlocks) {
                    this.currentPages.push(i);
                }
            }
        }
        return false;
    };
    __decorate([
        core_1.Output(), 
        __metadata('design:type', Object)
    ], SkycoinPaginationComponent.prototype, "onChangePage", void 0);
    SkycoinPaginationComponent = __decorate([
        core_1.Component({
            selector: 'app-skycoin-pagination',
            templateUrl: './skycoin-pagination.component.html',
            styleUrls: ['./skycoin-pagination.component.css']
        }), 
        __metadata('design:paramtypes', [skycoin_blockchain_pagination_service_1.SkycoinBlockchainPaginationService])
    ], SkycoinPaginationComponent);
    return SkycoinPaginationComponent;
}());
exports.SkycoinPaginationComponent = SkycoinPaginationComponent;
//# sourceMappingURL=skycoin-pagination.component.js.map