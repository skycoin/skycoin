webpackJsonp([0,4],{

/***/ 1121:
/***/ (function(module, exports, __webpack_require__) {

module.exports = __webpack_require__(550);


/***/ }),

/***/ 1125:
/***/ (function(module, exports, __webpack_require__) {

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
var core_1 = __webpack_require__(1);
var QRious = __webpack_require__(1124);
var QRCodeComponent = (function () {
    function QRCodeComponent(elementRef) {
        this.elementRef = elementRef;
        this.background = 'white';
        this.backgroundAlpha = 1.0;
        this.foreground = 'black';
        this.foregroundAlpha = 1.0;
        this.level = 'L';
        this.mime = 'image/png';
        this.padding = null;
        this.size = 100;
        this.value = '';
        this.canvas = false;
    }
    QRCodeComponent.prototype.ngOnChanges = function (changes) {
        if ('background' in changes ||
            'backgroundAlpha' in changes ||
            'foreground' in changes ||
            'foregroundAlpha' in changes ||
            'level' in changes ||
            'mime' in changes ||
            'padding' in changes ||
            'size' in changes ||
            'value' in changes ||
            'canvas' in changes) {
            this.generate();
        }
    };
    QRCodeComponent.prototype.generate = function () {
        try {
            var el = this.elementRef.nativeElement;
            el.innerHTML = '';
            var qr = new QRious({
                background: this.background,
                backgroundAlpha: this.backgroundAlpha,
                foreground: this.foreground,
                foregroundAlpha: this.foregroundAlpha,
                level: this.level,
                mime: this.mime,
                padding: this.padding,
                size: this.size,
                value: this.value
            });
            if (this.canvas) {
                el.appendChild(qr.canvas);
            }
            else {
                el.appendChild(qr.image);
            }
        }
        catch (e) {
            console.error("Could not generate QR Code: " + e.message);
        }
    };
    return QRCodeComponent;
}());
__decorate([
    core_1.Input(),
    __metadata("design:type", String)
], QRCodeComponent.prototype, "background", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", Number)
], QRCodeComponent.prototype, "backgroundAlpha", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", String)
], QRCodeComponent.prototype, "foreground", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", Number)
], QRCodeComponent.prototype, "foregroundAlpha", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", String)
], QRCodeComponent.prototype, "level", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", String)
], QRCodeComponent.prototype, "mime", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", Number)
], QRCodeComponent.prototype, "padding", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", Number)
], QRCodeComponent.prototype, "size", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", String)
], QRCodeComponent.prototype, "value", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", Boolean)
], QRCodeComponent.prototype, "canvas", void 0);
QRCodeComponent = __decorate([
    core_1.Component({
        moduleId: 'module.id',
        selector: 'qr-code',
        template: ""
    }),
    __metadata("design:paramtypes", [core_1.ElementRef])
], QRCodeComponent);
exports.QRCodeComponent = QRCodeComponent;
var QRCodeModule = (function () {
    function QRCodeModule() {
    }
    return QRCodeModule;
}());
QRCodeModule = __decorate([
    core_1.NgModule({
        exports: [QRCodeComponent],
        declarations: [QRCodeComponent],
        entryComponents: [QRCodeComponent]
    }),
    __metadata("design:paramtypes", [])
], QRCodeModule);
exports.QRCodeModule = QRCodeModule;


/***/ }),

/***/ 173:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__(117);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_map__ = __webpack_require__(523);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_map___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_map__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_catch__ = __webpack_require__(522);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_catch___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_catch__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return BlockChainService; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var BlockChainService = (function () {
    function BlockChainService(_http) {
        this._http = _http;
    }
    BlockChainService.prototype.getBlocks = function (startNumber, endNumber) {
        var stringConvert = 'start=' + startNumber + '&end=' + endNumber;
        return this._http.get('/api/blocks?' + stringConvert)
            .map(function (res) {
            console.log(res);
            return res.json();
        })
            .map(function (res) { return res.blocks; })
            .catch(function (error) {
            console.log(error);
            return __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__["Observable"].throw(error || 'Server error');
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
            return __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__["Observable"].throw(error || 'Server error');
        });
    };
    BlockChainService = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Injectable"])(), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */]) === 'function' && _a) || Object])
    ], BlockChainService);
    return BlockChainService;
    var _a;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/block-chain.service.js.map

/***/ }),

/***/ 364:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AppComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};

var AppComponent = (function () {
    function AppComponent() {
        this.title = 'app works!';
    }
    AppComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-root',
            template: __webpack_require__(846),
            styles: [__webpack_require__(836)]
        }), 
        __metadata('design:paramtypes', [])
    ], AppComponent);
    return AppComponent;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/app.component.js.map

/***/ }),

/***/ 365:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__(117);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs__ = __webpack_require__(264);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return UxOutputsService; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



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
            return __WEBPACK_IMPORTED_MODULE_2_rxjs__["Observable"].throw(error || 'Server error');
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
            return __WEBPACK_IMPORTED_MODULE_2_rxjs__["Observable"].throw(error || 'Server error');
        });
    };
    UxOutputsService.prototype.getBlockSource = function (blockNumber) {
        return this._http.get('/api/blocks?start=' + blockNumber + '&end=' + blockNumber)
            .map(function (res) {
            return res.json().blocks[0];
        })
            .catch(function (error) {
            console.log(error);
            return __WEBPACK_IMPORTED_MODULE_2_rxjs__["Observable"].throw(error || 'Server error');
        });
    };
    UxOutputsService = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Injectable"])(), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */]) === 'function' && _a) || Object])
    ], UxOutputsService);
    return UxOutputsService;
    var _a;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/UxOutputs.service.js.map

/***/ }),

/***/ 366:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__(117);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs__ = __webpack_require__(264);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_underscore__ = __webpack_require__(278);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_underscore___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_underscore__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return SkycoinBlockchainPaginationService; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};




var SkycoinBlockchainPaginationService = (function () {
    function SkycoinBlockchainPaginationService(_http) {
        this._http = _http;
    }
    SkycoinBlockchainPaginationService.prototype.fetchNumberOfBlocks = function () {
        return this._http.get('/api/blockchain/metadata')
            .map(function (res) {
            return res.json();
        })
            .map(function (res) { return res.head.seq; })
            .catch(function (error) {
            console.log(error);
            return __WEBPACK_IMPORTED_MODULE_2_rxjs__["Observable"].throw(error || 'Server error');
        });
    };
    SkycoinBlockchainPaginationService.prototype.getPager = function (totalItems, currentPage, pageSize) {
        if (currentPage === void 0) { currentPage = 1; }
        if (pageSize === void 0) { pageSize = 10; }
        // calculate total pages
        var totalPages = Math.ceil(totalItems / pageSize);
        var startPage, endPage;
        if (totalPages <= 10) {
            // less than 10 total pages so show all
            startPage = 1;
            endPage = totalPages;
        }
        else {
            // more than 10 total pages so calculate start and end pages
            if (currentPage <= 6) {
                startPage = 1;
                endPage = 10;
            }
            else if (currentPage + 4 >= totalPages) {
                startPage = totalPages - 9;
                endPage = totalPages;
            }
            else {
                startPage = currentPage - 5;
                endPage = currentPage + 4;
            }
        }
        // calculate start and end item indexes
        var startIndex = (currentPage - 1) * pageSize;
        var endIndex = Math.min(startIndex + pageSize - 1, totalItems - 1);
        // create an array of pages to ng-repeat in the pager control
        var pages = __WEBPACK_IMPORTED_MODULE_3_underscore__["range"](startPage, endPage + 1);
        // return object with all pager properties required by the view
        return {
            totalItems: totalItems,
            currentPage: currentPage,
            pageSize: pageSize,
            totalPages: totalPages,
            startPage: startPage,
            endPage: endPage,
            startIndex: startIndex,
            endIndex: endIndex,
            pages: pages
        };
    };
    SkycoinBlockchainPaginationService = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Injectable"])(), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */]) === 'function' && _a) || Object])
    ], SkycoinBlockchainPaginationService);
    return SkycoinBlockchainPaginationService;
    var _a;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/skycoin-blockchain-pagination.service.js.map

/***/ }),

/***/ 367:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__(117);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs__ = __webpack_require__(264);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return TransactionDetailService; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var TransactionDetailService = (function () {
    function TransactionDetailService(_http) {
        this._http = _http;
    }
    TransactionDetailService.prototype.getTransaction = function (txid) {
        return this._http.get('/api/transaction?txid=' + txid)
            .map(function (res) {
            return res.json();
        })
            .catch(function (error) {
            console.log(error);
            return __WEBPACK_IMPORTED_MODULE_2_rxjs__["Observable"].throw(error || 'Server error');
        });
    };
    TransactionDetailService = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Injectable"])(), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */]) === 'function' && _a) || Object])
    ], TransactionDetailService);
    return TransactionDetailService;
    var _a;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/transaction-detail.service.js.map

/***/ }),

/***/ 549:
/***/ (function(module, exports) {

function webpackEmptyContext(req) {
	throw new Error("Cannot find module '" + req + "'.");
}
webpackEmptyContext.keys = function() { return []; };
webpackEmptyContext.resolve = webpackEmptyContext;
module.exports = webpackEmptyContext;
webpackEmptyContext.id = 549;


/***/ }),

/***/ 550:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__polyfills_ts__ = __webpack_require__(679);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__ = __webpack_require__(637);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__environments_environment__ = __webpack_require__(678);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__app___ = __webpack_require__(677);





if (__WEBPACK_IMPORTED_MODULE_3__environments_environment__["a" /* environment */].production) {
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_2__angular_core__["enableProdMode"])();
}
__webpack_require__.i(__WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__["a" /* platformBrowserDynamic */])().bootstrapModule(__WEBPACK_IMPORTED_MODULE_4__app___["a" /* AppModule */]);
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/main.js.map

/***/ }),

/***/ 667:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__ = __webpack_require__(170);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_forms__ = __webpack_require__(628);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_http__ = __webpack_require__(117);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__angular_router__ = __webpack_require__(90);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__app_component__ = __webpack_require__(364);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6__components_skycoin_header_skycoin_header_component__ = __webpack_require__(672);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7__components_block_chain_table_block_chain_table_component__ = __webpack_require__(669);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8__components_skycoin_pagination_skycoin_pagination_component__ = __webpack_require__(674);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9__components_skycoin_search_bar_skycoin_search_bar_component__ = __webpack_require__(675);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10__components_block_chain_table_block_chain_service__ = __webpack_require__(173);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11__components_footer_footer_component__ = __webpack_require__(671);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12__components_skycoin_pagination_num_pages_pipe__ = __webpack_require__(673);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13__components_skycoin_pagination_skycoin_blockchain_pagination_service__ = __webpack_require__(366);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14__angular_common__ = __webpack_require__(99);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_angular2_moment__ = __webpack_require__(680);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_angular2_moment___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_15_angular2_moment__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_16__components_block_details_block_details_component__ = __webpack_require__(670);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_17__components_address_detail_address_detail_component__ = __webpack_require__(668);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_18__components_address_detail_UxOutputs_service__ = __webpack_require__(365);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_19__components_transaction_detail_transaction_detail_component__ = __webpack_require__(676);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_20__components_transaction_detail_transaction_detail_service__ = __webpack_require__(367);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_21__js_angular2_qrcode__ = __webpack_require__(1125);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_21__js_angular2_qrcode___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_21__js_angular2_qrcode__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AppModule; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};






















var ROUTES = [
    {
        path: '',
        redirectTo: 'blocks',
        pathMatch: 'full'
    },
    {
        path: 'blocks',
        component: __WEBPACK_IMPORTED_MODULE_7__components_block_chain_table_block_chain_table_component__["a" /* BlockChainTableComponent */]
    },
    {
        path: 'block/:id',
        component: __WEBPACK_IMPORTED_MODULE_16__components_block_details_block_details_component__["a" /* BlockDetailsComponent */]
    },
    {
        path: 'address/:address',
        component: __WEBPACK_IMPORTED_MODULE_17__components_address_detail_address_detail_component__["a" /* AddressDetailComponent */]
    },
    {
        path: 'transaction/:txid',
        component: __WEBPACK_IMPORTED_MODULE_19__components_transaction_detail_transaction_detail_component__["a" /* TransactionDetailComponent */]
    }
];
var AppModule = (function () {
    function AppModule() {
    }
    AppModule = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_1__angular_core__["NgModule"])({
            declarations: [
                __WEBPACK_IMPORTED_MODULE_5__app_component__["a" /* AppComponent */],
                __WEBPACK_IMPORTED_MODULE_6__components_skycoin_header_skycoin_header_component__["a" /* SkycoinHeaderComponent */],
                __WEBPACK_IMPORTED_MODULE_7__components_block_chain_table_block_chain_table_component__["a" /* BlockChainTableComponent */],
                __WEBPACK_IMPORTED_MODULE_8__components_skycoin_pagination_skycoin_pagination_component__["a" /* SkycoinPaginationComponent */],
                __WEBPACK_IMPORTED_MODULE_9__components_skycoin_search_bar_skycoin_search_bar_component__["a" /* SkycoinSearchBarComponent */],
                __WEBPACK_IMPORTED_MODULE_11__components_footer_footer_component__["a" /* FooterComponent */],
                __WEBPACK_IMPORTED_MODULE_12__components_skycoin_pagination_num_pages_pipe__["a" /* NumPagesPipe */],
                __WEBPACK_IMPORTED_MODULE_16__components_block_details_block_details_component__["a" /* BlockDetailsComponent */],
                __WEBPACK_IMPORTED_MODULE_17__components_address_detail_address_detail_component__["a" /* AddressDetailComponent */],
                __WEBPACK_IMPORTED_MODULE_19__components_transaction_detail_transaction_detail_component__["a" /* TransactionDetailComponent */]
            ],
            imports: [
                __WEBPACK_IMPORTED_MODULE_14__angular_common__["a" /* CommonModule */],
                __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__["a" /* BrowserModule */],
                __WEBPACK_IMPORTED_MODULE_2__angular_forms__["a" /* FormsModule */],
                __WEBPACK_IMPORTED_MODULE_3__angular_http__["a" /* HttpModule */],
                __WEBPACK_IMPORTED_MODULE_15_angular2_moment__["MomentModule"],
                __WEBPACK_IMPORTED_MODULE_21__js_angular2_qrcode__["QRCodeModule"],
                __WEBPACK_IMPORTED_MODULE_4__angular_router__["a" /* RouterModule */].forRoot(ROUTES)
            ],
            providers: [__WEBPACK_IMPORTED_MODULE_10__components_block_chain_table_block_chain_service__["a" /* BlockChainService */], __WEBPACK_IMPORTED_MODULE_13__components_skycoin_pagination_skycoin_blockchain_pagination_service__["a" /* SkycoinBlockchainPaginationService */], __WEBPACK_IMPORTED_MODULE_18__components_address_detail_UxOutputs_service__["a" /* UxOutputsService */], __WEBPACK_IMPORTED_MODULE_20__components_transaction_detail_transaction_detail_service__["a" /* TransactionDetailService */]],
            bootstrap: [__WEBPACK_IMPORTED_MODULE_5__app_component__["a" /* AppComponent */]]
        }), 
        __metadata('design:paramtypes', [])
    ], AppModule);
    return AppModule;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/app.module.js.map

/***/ }),

/***/ 668:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__UxOutputs_service__ = __webpack_require__(365);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__(90);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AddressDetailComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var AddressDetailComponent = (function () {
    function AddressDetailComponent(service, route, router) {
        this.service = service;
        this.route = route;
        this.router = router;
        this.UxOutputs = null;
        this.transactions = [];
        this.currentAddress = null;
    }
    AddressDetailComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.UxOutputs = this.route.params
            .switchMap(function (params) {
            var address = params['address'];
            _this.currentAddress = address;
            return _this.service.getUxOutputsForAddress(address);
        });
        this.UxOutputs.subscribe(function (uxoutputs) {
            _this.transactions = uxoutputs;
            console.log(uxoutputs);
        });
    };
    AddressDetailComponent.prototype.getCurrentBalance = function () {
        var outputs = this.transactions[this.transactions.length - 1].outputs;
        if (this.currentAddress) {
            for (var i = 0; i < outputs.length; i++) {
                var currentAddress = outputs[i].dst;
                if (currentAddress == this.currentAddress) {
                    return outputs[i].coins;
                }
            }
        }
        return "0";
    };
    AddressDetailComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-address-detail',
            template: __webpack_require__(847),
            styles: [__webpack_require__(837)],
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__UxOutputs_service__["a" /* UxOutputsService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__UxOutputs_service__["a" /* UxOutputsService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["c" /* ActivatedRoute */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__angular_router__["c" /* ActivatedRoute */]) === 'function' && _b) || Object, (typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */]) === 'function' && _c) || Object])
    ], AddressDetailComponent);
    return AddressDetailComponent;
    var _a, _b, _c;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/address-detail.component.js.map

/***/ }),

/***/ 669:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__block_chain_service__ = __webpack_require__(173);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_moment__ = __webpack_require__(2);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_moment___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_moment__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_underscore__ = __webpack_require__(278);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_underscore___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_underscore__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__angular_router__ = __webpack_require__(90);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return BlockChainTableComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var BlockChainTableComponent = (function () {
    function BlockChainTableComponent(blockService, router) {
        this.blockService = blockService;
        this.router = router;
        this.blocks = [];
    }
    BlockChainTableComponent.prototype.GetBlockAmount = function (txns) {
        var ret = [];
        __WEBPACK_IMPORTED_MODULE_3_underscore__["each"](txns, function (o) {
            if (o.outputs) {
                __WEBPACK_IMPORTED_MODULE_3_underscore__["each"](o.outputs, function (_o) {
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
        return __WEBPACK_IMPORTED_MODULE_2_moment__["unix"](time).format("YYYY-MM-DD HH:mm");
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
        this.blockService.getBlocks(blockStart, blockEnd).subscribe(function (data) {
            console.log(data);
            var newData = __WEBPACK_IMPORTED_MODULE_3_underscore__["sortBy"](data, function (block) { return block.header.seq; }).reverse();
            _this.blocks = newData;
        });
    };
    BlockChainTableComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-block-chain-table',
            template: __webpack_require__(848),
            styles: [__webpack_require__(838)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__block_chain_service__["a" /* BlockChainService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__block_chain_service__["a" /* BlockChainService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_4__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_4__angular_router__["b" /* Router */]) === 'function' && _b) || Object])
    ], BlockChainTableComponent);
    return BlockChainTableComponent;
    var _a, _b;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/block-chain-table.component.js.map

/***/ }),

/***/ 670:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_router__ = __webpack_require__(90);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__block_chain_table_block_chain_service__ = __webpack_require__(173);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_underscore__ = __webpack_require__(278);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_underscore___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_underscore__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_moment__ = __webpack_require__(2);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_moment___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4_moment__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return BlockDetailsComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





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
        return __WEBPACK_IMPORTED_MODULE_4_moment__["unix"](time).format();
    };
    BlockDetailsComponent.prototype.getAmount = function (block) {
        var ret = [];
        var txns = block.body.txns;
        __WEBPACK_IMPORTED_MODULE_3_underscore__["each"](txns, function (o) {
            if (o.outputs) {
                __WEBPACK_IMPORTED_MODULE_3_underscore__["each"](o.outputs, function (_o) {
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
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-block-details',
            template: __webpack_require__(849),
            styles: [__webpack_require__(839)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_2__block_chain_table_block_chain_service__["a" /* BlockChainService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__block_chain_table_block_chain_service__["a" /* BlockChainService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__angular_router__["c" /* ActivatedRoute */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_router__["c" /* ActivatedRoute */]) === 'function' && _b) || Object, (typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_1__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_router__["b" /* Router */]) === 'function' && _c) || Object])
    ], BlockDetailsComponent);
    return BlockDetailsComponent;
    var _a, _b, _c;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/block-details.component.js.map

/***/ }),

/***/ 671:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return FooterComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};

var FooterComponent = (function () {
    function FooterComponent() {
    }
    FooterComponent.prototype.ngOnInit = function () {
    };
    FooterComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-footer',
            template: __webpack_require__(850),
            styles: [__webpack_require__(840)]
        }), 
        __metadata('design:paramtypes', [])
    ], FooterComponent);
    return FooterComponent;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/footer.component.js.map

/***/ }),

/***/ 672:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return SkycoinHeaderComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};

var SkycoinHeaderComponent = (function () {
    function SkycoinHeaderComponent() {
    }
    SkycoinHeaderComponent.prototype.ngOnInit = function () {
    };
    SkycoinHeaderComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-skycoin-header',
            template: __webpack_require__(851),
            styles: [__webpack_require__(841)]
        }), 
        __metadata('design:paramtypes', [])
    ], SkycoinHeaderComponent);
    return SkycoinHeaderComponent;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/skycoin-header.component.js.map

/***/ }),

/***/ 673:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return NumPagesPipe; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};

var NumPagesPipe = (function () {
    function NumPagesPipe() {
    }
    NumPagesPipe.prototype.transform = function (value, args) {
        var res = [];
        for (var i = 1; i <= value; i++) {
            res.push(i);
        }
        return res;
    };
    NumPagesPipe = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Pipe"])({
            name: 'numPages'
        }), 
        __metadata('design:paramtypes', [])
    ], NumPagesPipe);
    return NumPagesPipe;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/num-pages.pipe.js.map

/***/ }),

/***/ 674:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__skycoin_blockchain_pagination_service__ = __webpack_require__(366);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return SkycoinPaginationComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var SkycoinPaginationComponent = (function () {
    function SkycoinPaginationComponent(paginationService) {
        this.paginationService = paginationService;
        this.onChangePage = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["EventEmitter"]();
        this.numberOfBlocks = 0;
        this.currentPage = 1;
        this.currentPages = [];
        this.pagesToShowAtATime = 5;
        this.pageStartPointer = this.currentPage;
        this.pageEndPointer = this.currentPage;
    }
    SkycoinPaginationComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.paginationService.fetchNumberOfBlocks().subscribe(function (numberOfBlocks) {
            _this.numberOfBlocks = numberOfBlocks;
            _this.onChangePage.emit([1, _this.numberOfBlocks]);
            _this.pagesToShowAtATime = _this.pagesToShowAtATime < numberOfBlocks ? _this.pagesToShowAtATime : _this.numberOfBlocks;
            _this.currentPages = [];
            for (var i = _this.currentPage; i <= _this.currentPage + 4; i++) {
                _this.currentPages.push(i);
            }
        });
    };
    SkycoinPaginationComponent.prototype.setPage = function (currentPage) {
        if (!(currentPage in this.pages)) {
        }
    };
    SkycoinPaginationComponent.prototype.changePage = function (pageNumber) {
        this.onChangePage.emit([pageNumber, this.numberOfBlocks]);
        this.currentPage = pageNumber;
        return false;
    };
    SkycoinPaginationComponent.prototype.loadUpcoming = function () {
        this.onChangePage.emit([this.currentPages[0] + this.pagesToShowAtATime, this.numberOfBlocks]);
        this.currentPage = this.currentPages[0] + this.pagesToShowAtATime;
        this.currentPages = [];
        for (var i = this.currentPage; i <= this.currentPage + 4; i++) {
            if (this.numberOfBlocks - i * 10 >= 0) {
                this.currentPages.push(i);
            }
            else if (this.numberOfBlocks - i * 10 >= -10) {
                this.currentPages.push(i);
            }
        }
        return false;
    };
    SkycoinPaginationComponent.prototype.loadPrevious = function () {
        if (this.currentPages[0] <= 1) {
            return false;
        }
        this.onChangePage.emit([this.currentPages[0] - this.pagesToShowAtATime, this.numberOfBlocks]);
        this.currentPage = this.currentPages[0] - this.pagesToShowAtATime;
        this.currentPages = [];
        for (var i = this.currentPage; i <= this.currentPage + 4; i++) {
            if (i * 10 <= this.numberOfBlocks) {
                this.currentPages.push(i);
            }
        }
        return false;
    };
    __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Output"])(), 
        __metadata('design:type', Object)
    ], SkycoinPaginationComponent.prototype, "onChangePage", void 0);
    SkycoinPaginationComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-skycoin-pagination',
            template: __webpack_require__(852),
            styles: [__webpack_require__(842)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__skycoin_blockchain_pagination_service__["a" /* SkycoinBlockchainPaginationService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__skycoin_blockchain_pagination_service__["a" /* SkycoinBlockchainPaginationService */]) === 'function' && _a) || Object])
    ], SkycoinPaginationComponent);
    return SkycoinPaginationComponent;
    var _a;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/skycoin-pagination.component.js.map

/***/ }),

/***/ 675:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__block_chain_table_block_chain_service__ = __webpack_require__(173);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__(90);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return SkycoinSearchBarComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var SkycoinSearchBarComponent = (function () {
    function SkycoinSearchBarComponent(service, router) {
        this.service = service;
        this.router = router;
    }
    SkycoinSearchBarComponent.prototype.ngOnInit = function () {
    };
    SkycoinSearchBarComponent.prototype.searchBlockHistory = function (hashVal) {
        if (hashVal.length == 34) {
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
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-skycoin-search-bar',
            template: __webpack_require__(853),
            styles: [__webpack_require__(843)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__block_chain_table_block_chain_service__["a" /* BlockChainService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__block_chain_table_block_chain_service__["a" /* BlockChainService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */]) === 'function' && _b) || Object])
    ], SkycoinSearchBarComponent);
    return SkycoinSearchBarComponent;
    var _a, _b;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/skycoin-search-bar.component.js.map

/***/ }),

/***/ 676:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_router__ = __webpack_require__(90);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__transaction_detail_service__ = __webpack_require__(367);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return TransactionDetailComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



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
            .switchMap(function (params) {
            var txid = params['txid'];
            return _this.service.getTransaction(txid);
        });
        this.transactionObservable.subscribe(function (trans) {
            _this.transaction = trans;
            console.log(trans);
        });
    };
    TransactionDetailComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-transaction-detail',
            template: __webpack_require__(854),
            styles: [__webpack_require__(844)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_2__transaction_detail_service__["a" /* TransactionDetailService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__transaction_detail_service__["a" /* TransactionDetailService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__angular_router__["c" /* ActivatedRoute */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_router__["c" /* ActivatedRoute */]) === 'function' && _b) || Object, (typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_1__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_router__["b" /* Router */]) === 'function' && _c) || Object])
    ], TransactionDetailComponent);
    return TransactionDetailComponent;
    var _a, _b, _c;
}());
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/transaction-detail.component.js.map

/***/ }),

/***/ 677:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__app_component__ = __webpack_require__(364);
/* unused harmony namespace reexport */
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__app_module__ = __webpack_require__(667);
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "a", function() { return __WEBPACK_IMPORTED_MODULE_1__app_module__["a"]; });


//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/index.js.map

/***/ }),

/***/ 678:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return environment; });
// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `angular-cli.json`.
var environment = {
    production: false
};
//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/environment.js.map

/***/ }),

/***/ 679:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_core_js_es6_symbol__ = __webpack_require__(695);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_core_js_es6_symbol___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_core_js_es6_symbol__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_core_js_es6_object__ = __webpack_require__(688);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_core_js_es6_object___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_core_js_es6_object__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_core_js_es6_function__ = __webpack_require__(684);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_core_js_es6_function___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_core_js_es6_function__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_core_js_es6_parse_int__ = __webpack_require__(690);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_core_js_es6_parse_int___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_core_js_es6_parse_int__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_core_js_es6_parse_float__ = __webpack_require__(689);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_core_js_es6_parse_float___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4_core_js_es6_parse_float__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_core_js_es6_number__ = __webpack_require__(687);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_core_js_es6_number___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_5_core_js_es6_number__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6_core_js_es6_math__ = __webpack_require__(686);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6_core_js_es6_math___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_6_core_js_es6_math__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7_core_js_es6_string__ = __webpack_require__(694);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7_core_js_es6_string___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_7_core_js_es6_string__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8_core_js_es6_date__ = __webpack_require__(683);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8_core_js_es6_date___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_8_core_js_es6_date__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9_core_js_es6_array__ = __webpack_require__(682);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9_core_js_es6_array___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_9_core_js_es6_array__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10_core_js_es6_regexp__ = __webpack_require__(692);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10_core_js_es6_regexp___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_10_core_js_es6_regexp__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11_core_js_es6_map__ = __webpack_require__(685);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11_core_js_es6_map___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_11_core_js_es6_map__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12_core_js_es6_set__ = __webpack_require__(693);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12_core_js_es6_set___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_12_core_js_es6_set__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13_core_js_es6_reflect__ = __webpack_require__(691);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13_core_js_es6_reflect___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_13_core_js_es6_reflect__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14_core_js_es7_reflect__ = __webpack_require__(696);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14_core_js_es7_reflect___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_14_core_js_es7_reflect__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_zone_js_dist_zone__ = __webpack_require__(1120);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_zone_js_dist_zone___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_15_zone_js_dist_zone__);
















//# sourceMappingURL=/Users/napandey/coding/skycoin/skycoin-web/src/polyfills.js.map

/***/ }),

/***/ 834:
/***/ (function(module, exports, __webpack_require__) {

var map = {
	"./af": 412,
	"./af.js": 412,
	"./ar": 418,
	"./ar-dz": 413,
	"./ar-dz.js": 413,
	"./ar-ly": 414,
	"./ar-ly.js": 414,
	"./ar-ma": 415,
	"./ar-ma.js": 415,
	"./ar-sa": 416,
	"./ar-sa.js": 416,
	"./ar-tn": 417,
	"./ar-tn.js": 417,
	"./ar.js": 418,
	"./az": 419,
	"./az.js": 419,
	"./be": 420,
	"./be.js": 420,
	"./bg": 421,
	"./bg.js": 421,
	"./bn": 422,
	"./bn.js": 422,
	"./bo": 423,
	"./bo.js": 423,
	"./br": 424,
	"./br.js": 424,
	"./bs": 425,
	"./bs.js": 425,
	"./ca": 426,
	"./ca.js": 426,
	"./cs": 427,
	"./cs.js": 427,
	"./cv": 428,
	"./cv.js": 428,
	"./cy": 429,
	"./cy.js": 429,
	"./da": 430,
	"./da.js": 430,
	"./de": 432,
	"./de-at": 431,
	"./de-at.js": 431,
	"./de.js": 432,
	"./dv": 433,
	"./dv.js": 433,
	"./el": 434,
	"./el.js": 434,
	"./en-au": 435,
	"./en-au.js": 435,
	"./en-ca": 436,
	"./en-ca.js": 436,
	"./en-gb": 437,
	"./en-gb.js": 437,
	"./en-ie": 438,
	"./en-ie.js": 438,
	"./en-nz": 439,
	"./en-nz.js": 439,
	"./eo": 440,
	"./eo.js": 440,
	"./es": 442,
	"./es-do": 441,
	"./es-do.js": 441,
	"./es.js": 442,
	"./et": 443,
	"./et.js": 443,
	"./eu": 444,
	"./eu.js": 444,
	"./fa": 445,
	"./fa.js": 445,
	"./fi": 446,
	"./fi.js": 446,
	"./fo": 447,
	"./fo.js": 447,
	"./fr": 450,
	"./fr-ca": 448,
	"./fr-ca.js": 448,
	"./fr-ch": 449,
	"./fr-ch.js": 449,
	"./fr.js": 450,
	"./fy": 451,
	"./fy.js": 451,
	"./gd": 452,
	"./gd.js": 452,
	"./gl": 453,
	"./gl.js": 453,
	"./he": 454,
	"./he.js": 454,
	"./hi": 455,
	"./hi.js": 455,
	"./hr": 456,
	"./hr.js": 456,
	"./hu": 457,
	"./hu.js": 457,
	"./hy-am": 458,
	"./hy-am.js": 458,
	"./id": 459,
	"./id.js": 459,
	"./is": 460,
	"./is.js": 460,
	"./it": 461,
	"./it.js": 461,
	"./ja": 462,
	"./ja.js": 462,
	"./jv": 463,
	"./jv.js": 463,
	"./ka": 464,
	"./ka.js": 464,
	"./kk": 465,
	"./kk.js": 465,
	"./km": 466,
	"./km.js": 466,
	"./ko": 467,
	"./ko.js": 467,
	"./ky": 468,
	"./ky.js": 468,
	"./lb": 469,
	"./lb.js": 469,
	"./lo": 470,
	"./lo.js": 470,
	"./lt": 471,
	"./lt.js": 471,
	"./lv": 472,
	"./lv.js": 472,
	"./me": 473,
	"./me.js": 473,
	"./mi": 474,
	"./mi.js": 474,
	"./mk": 475,
	"./mk.js": 475,
	"./ml": 476,
	"./ml.js": 476,
	"./mr": 477,
	"./mr.js": 477,
	"./ms": 479,
	"./ms-my": 478,
	"./ms-my.js": 478,
	"./ms.js": 479,
	"./my": 480,
	"./my.js": 480,
	"./nb": 481,
	"./nb.js": 481,
	"./ne": 482,
	"./ne.js": 482,
	"./nl": 484,
	"./nl-be": 483,
	"./nl-be.js": 483,
	"./nl.js": 484,
	"./nn": 485,
	"./nn.js": 485,
	"./pa-in": 486,
	"./pa-in.js": 486,
	"./pl": 487,
	"./pl.js": 487,
	"./pt": 489,
	"./pt-br": 488,
	"./pt-br.js": 488,
	"./pt.js": 489,
	"./ro": 490,
	"./ro.js": 490,
	"./ru": 491,
	"./ru.js": 491,
	"./se": 492,
	"./se.js": 492,
	"./si": 493,
	"./si.js": 493,
	"./sk": 494,
	"./sk.js": 494,
	"./sl": 495,
	"./sl.js": 495,
	"./sq": 496,
	"./sq.js": 496,
	"./sr": 498,
	"./sr-cyrl": 497,
	"./sr-cyrl.js": 497,
	"./sr.js": 498,
	"./ss": 499,
	"./ss.js": 499,
	"./sv": 500,
	"./sv.js": 500,
	"./sw": 501,
	"./sw.js": 501,
	"./ta": 502,
	"./ta.js": 502,
	"./te": 503,
	"./te.js": 503,
	"./tet": 504,
	"./tet.js": 504,
	"./th": 505,
	"./th.js": 505,
	"./tl-ph": 506,
	"./tl-ph.js": 506,
	"./tlh": 507,
	"./tlh.js": 507,
	"./tr": 508,
	"./tr.js": 508,
	"./tzl": 509,
	"./tzl.js": 509,
	"./tzm": 511,
	"./tzm-latn": 510,
	"./tzm-latn.js": 510,
	"./tzm.js": 511,
	"./uk": 512,
	"./uk.js": 512,
	"./uz": 513,
	"./uz.js": 513,
	"./vi": 514,
	"./vi.js": 514,
	"./x-pseudo": 515,
	"./x-pseudo.js": 515,
	"./yo": 516,
	"./yo.js": 516,
	"./zh-cn": 517,
	"./zh-cn.js": 517,
	"./zh-hk": 518,
	"./zh-hk.js": 518,
	"./zh-tw": 519,
	"./zh-tw.js": 519
};
function webpackContext(req) {
	return __webpack_require__(webpackContextResolve(req));
};
function webpackContextResolve(req) {
	var id = map[req];
	if(!(id + 1)) // check for number
		throw new Error("Cannot find module '" + req + "'.");
	return id;
};
webpackContext.keys = function webpackContextKeys() {
	return Object.keys(map);
};
webpackContext.resolve = webpackContextResolve;
module.exports = webpackContext;
webpackContext.id = 834;


/***/ }),

/***/ 836:
/***/ (function(module, exports) {

module.exports = ""

/***/ }),

/***/ 837:
/***/ (function(module, exports) {

module.exports = ".collection{\n  border: none;\n}\n.collection-item{\n  border: none;\n}\n\n.card-title{\n  font-size: 16px;\n}\n.right-arrow{\n  margin-top: 50px;\n}\n\n.address-qr{\n  padding-top: 20px;\n}\n"

/***/ }),

/***/ 838:
/***/ (function(module, exports) {

module.exports = "tr:hover{\n  cursor: pointer;\n}\n\n.wrap{\n  word-break: break-all;\n}\n"

/***/ }),

/***/ 839:
/***/ (function(module, exports) {

module.exports = ".card-content{\n  background-color: #29b6f6;\n  color:white;\n}\n"

/***/ }),

/***/ 840:
/***/ (function(module, exports) {

module.exports = "footer{\n  background-color: #29b6f6;\n}\n"

/***/ }),

/***/ 841:
/***/ (function(module, exports) {

module.exports = ""

/***/ }),

/***/ 842:
/***/ (function(module, exports) {

module.exports = ".pagination li.active{\n  background-color: #29b6f6;\n  color:white;\n}\n\n.lower-panel{\n  color:black;\n}\n"

/***/ }),

/***/ 843:
/***/ (function(module, exports) {

module.exports = ""

/***/ }),

/***/ 844:
/***/ (function(module, exports) {

module.exports = ".collection{\n  border: none;\n}\n.collection-item{\n  border: none;\n}\n\n.card-title{\n  font-size: 16px;\n}\n.right-arrow{\n  margin-top: 50px;\n}\n"

/***/ }),

/***/ 846:
/***/ (function(module, exports) {

module.exports = "<app-skycoin-header></app-skycoin-header>\n<div class=\"container\">\n  <div class=\"row center\">\n    <app-skycoin-search-bar></app-skycoin-search-bar>\n    <router-outlet></router-outlet>\n  </div>\n</div>\n<app-footer></app-footer>\n\n\n\n\n"

/***/ }),

/***/ 847:
/***/ (function(module, exports) {

module.exports = "<div class=\"row\" *ngIf=\"transactions\">\n  <div class=\"col s8\" >\n    <div class=\"card-content\" *ngIf=\"transactions.length>0\">\n      <ul class=\"collection\">\n        <li class=\"collection-item\"><h5>Address</h5><span>{{currentAddress}}</span></li>\n        <li class=\"collection-item\">Number of transactions: {{transactions.length}} </li>\n        <li class=\"collection-item\">Current value: {{getCurrentBalance()}} skycoins</li>\n      </ul>\n    </div>\n  </div>\n\n  <div class=\"card col s4\">\n    <qr-code *ngIf=\"currentAddress!=null\" [value]=\"currentAddress\" [size]=\"300\"></qr-code>\n    <p class=\"address-qr\">{{currentAddress}}</p>\n  </div>\n\n\n  <div class=\" card col s12\" *ngFor=\"let transaction of transactions\">\n    <div class=\"card-title\">Transaction id : <a href=\"/transaction/{{transaction.txid}}\">{{transaction.txid}}</a></div>\n      <div class=\"col s4\">\n        <ul class=\"collection\">\n          <li class=\"collection-item\"><h5>inputs</h5></li>\n          <li class=\"collection-item\" *ngFor=\"let input of transaction.inputs\">\n            <!--<a href=\"/address/{{input.uxid}}\">uxid</a>-->\n            <a href=\"/address/{{input.owner}}\">{{input.owner}}</a></li>\n        </ul>\n      </div>\n      <div class=\"accent-1 col s1 right-arrow center\">\n        <i class=\"material-icons\">trending_flat</i>\n      </div>\n\n      <div class=\"col s5\">\n        <ul class=\"collection\">\n          <li class=\"collection-item\"><h5>outputs</h5></li>\n          <li class=\"collection-item\" *ngFor=\"let output of transaction.outputs\"><a href=\"/address/{{output.dst}}\">{{output.dst}}</a></li>\n        </ul>\n      </div>\n      <div class=\"col s2\">\n        <ul class=\"collection\">\n          <li class=\"collection-item\"><h5>coins</h5></li>\n          <li class=\"collection-item\" *ngFor=\"let output of transaction.outputs\">{{output.coins}}</li>\n        </ul>\n      </div>\n  </div>\n\n\n\n\n\n\n\n</div>\n"

/***/ }),

/***/ 848:
/***/ (function(module, exports) {

module.exports = "<h2>Blockchain</h2>\n<h5>{{totalBlocks}} blocks</h5>\n<table class=\"bordered striped centered responsive-table\">\n  <thead >\n  <tr>\n    <th data-field=\"id\">Time</th>\n    <th data-field=\"name\">Block Number</th>\n    <th data-field=\"price\">Transactions</th>\n    <th data-field=\"id\">Amount Transferred</th>\n    <th data-field=\"price\">Blockhash</th>\n\n  </tr>\n  </thead>\n\n  <tbody>\n  <tr *ngFor=\"let block of blocks\" (click)=\"showDetails(block)\">\n    <td>{{getTime(block.header.timestamp)}}</td>\n    <td >{{ block.header.seq }}</td>\n    <td>{{ block.body.txns.length }}</td>\n    <td>{{ GetBlockAmount(block.body.txns) }}</td>\n    <td>{{ block.header.block_hash }}</td>\n  </tr>\n\n\n  </tbody>\n</table>\n\n<app-skycoin-pagination (onChangePage)=\"handlePageChange($event)\"></app-skycoin-pagination>\n"

/***/ }),

/***/ 849:
/***/ (function(module, exports) {

module.exports = "<div class=\"row\">\n  <h2>Block Details</h2>\n  <div class=\"card\" *ngIf=\"block\">\n    <div class=\"card-content\" >\n      <table>\n        <tbody >\n\n        <tr>\n\n          <td>Height</td><td>{{ block.header.seq }}</td>\n\n        </tr>\n        <tr>\n          <td>Timestamp</td><td>{{getTime(block.header.timestamp)}}</td>\n        </tr>\n        <tr>\n          <td>Hash</td><td>{{ block.header.block_hash }}</td>\n        </tr>\n\n        <tr>\n          <td>Parent Hash</td><td>{{ block.header.previous_block_hash }}</td>\n        </tr>\n\n        <tr>\n          <td>Total Amount</td><td>{{ getAmount(block) }} coins</td>\n        </tr>\n\n        </tbody>\n      </table>\n    </div>\n  </div>\n\n\n    <div class=\"col s12\">Transactions</div>\n    <div class=\"col s12\">\n      <div class=\"col s1\">S.no</div>\n      <div class=\"col s2\">Transaction Id</div>\n      <div class=\"col s3\">Inputs</div>\n      <div class=\"col s4\">Outputs</div>\n      <div class=\"col s2\">Amount</div>\n    </div>\n  <div *ngIf=\"block\">\n    <div  class=\"col s12\" *ngFor=\"let transaction of block.body.txns;let i = index\">\n\n      <div class=\"col s1\">{{i+1}}</div>\n      <div class=\"col s2 wrap\"><a href=\"/transaction/{{transaction.txid}}\">{{transaction.txid}}</a></div>\n      <div class=\"col s3 wrap\"><p *ngFor=\"let input of transaction.inputs;\" class=\"input\">{{input}}</p></div>\n      <div class=\"col s6 wrap\"><p *ngFor=\"let output of transaction.outputs;\" class=\"input\"><a href=\"/address/{{output.dst}}\">{{output.dst}}</a>  - {{output.coins}}</p></div>\n    </div>\n  </div>\n\n\n\n\n</div>\n"

/***/ }),

/***/ 850:
/***/ (function(module, exports) {

module.exports = "<footer class=\"page-footer\">\n  <div class=\"container\">\n    <div class=\"row\">\n      <div class=\"col l6 s12\">\n        <h5 class=\"white-text\">The Skycoin Wallet allows you to hold and secure skycoin. It not only gives you access to the Skycoin blockchain but it’s fast, secure and easy to use!</h5>\n        <p class=\"grey-text text-lighten-4\">Skycoin is a new form of decentralized digital currency that is created and held electronically. It is a necessary element for operating the Skycoin platform.</p>\n      </div>\n      <div class=\"col l4 offset-l2 s12\">\n        <h5 class=\"white-text\">Links</h5>\n        <ul>\n          <li><a class=\"grey-text text-lighten-3\" href=\"http://skycoin.net/infographics.html\">How it works?</a></li>\n          <li><a class=\"grey-text text-lighten-3\" href=\"http://skycoin.net/downloads.html\">Download wallet</a></li>\n          <li><a class=\"grey-text text-lighten-3\" href=\"http://skycoin.net/whitepapers.html\">White papers</a></li>\n          <li><a class=\"grey-text text-lighten-3\" href=\"http://skycoin.net/faq.html\">FAQs</a></li>\n        </ul>\n      </div>\n    </div>\n  </div>\n  <div class=\"footer-copyright\">\n    <div class=\"container\">\n      © 2016 Skycoin\n      <a class=\"grey-text text-lighten-4 right\" href=\"#!\">More Links</a>\n    </div>\n  </div>\n</footer>\n"

/***/ }),

/***/ 851:
/***/ (function(module, exports) {

module.exports = "<nav class=\"light-blue lighten-1\" role=\"navigation\">\n  <div class=\"nav-wrapper container\"><a id=\"logo-container\" href=\"#\" class=\"brand-logo\">Skycoin</a>\n    <ul class=\"right hide-on-med-and-down\">\n      <li><a href=\"http://skycoin.net/infographics.html\">How it works?</a></li>\n      <li><a href=\"http://skycoin.net/downloads.html\">Download wallet</a></li>\n      <li><a href=\"http://skycoin.net/faq.html\">FAQs</a></li>\n      <li><a href=\"http://skycoin.net/whitepapers.html\">White papers</a></li>\n    </ul>\n\n    <ul id=\"nav-mobile\" class=\"side-nav\">\n      <li><a href=\"http://skycoin.net/infographics.html\">How it works?</a></li>\n      <li><a href=\"http://skycoin.net/downloads.html\">Download wallet</a></li>\n      <li><a href=\"http://skycoin.net/faq.html\">FAQs</a></li>\n      <li><a href=\"http://skycoin.net/whitepapers.html\">White papers</a></li>\n    </ul>\n    <a href=\"#\" data-activates=\"nav-mobile\" class=\"button-collapse\"><i class=\"material-icons\">menu</i></a>\n  </div>\n</nav>\n"

/***/ }),

/***/ 852:
/***/ (function(module, exports) {

module.exports = "<ul class=\"pagination\">\n\n  <li [ngClass]=\"{disabled:currentPage <= 5}\"><a href=\"#!\" (click)=\"loadPrevious()\"><i class=\"material-icons\">chevron_left</i></a></li>\n\n  <li *ngFor='let page of currentPages' [ngClass]=\"{active:currentPage == page}\" class=\"active\">\n\n    <a  (click)=\"changePage(page)\" href=\"#!\">{{page}}</a>\n\n  </li>\n  <li [ngClass]=\"{disabled:currentPage === numberOfBlocks/10}\"><a href=\"#!\" (click)=\"loadUpcoming()\"><i class=\"material-icons\">chevron_right</i></a></li>\n</ul>\n"

/***/ }),

/***/ 853:
/***/ (function(module, exports) {

module.exports = "<div class=\"nav-wrapper\">\n  <form (submit)=\"searchBlockHistory(blockSearchKey.value)\">\n    <div class=\"input-field\">\n      <input #blockSearchKey id=\"search\" type=\"search\" required placeholder=\"blockhash, address, block-number, transaction id\">\n      <label class=\"label-icon\" for=\"search\"><i class=\"material-icons\">search</i></label>\n      <i class=\"material-icons\">close</i>\n    </div>\n  </form>\n</div>\n"

/***/ }),

/***/ 854:
/***/ (function(module, exports) {

module.exports = "<div class=\"row\" *ngIf=\"transaction!=null\">\n  <h2>Transaction Details</h2>\n  <div class=\"card\" *ngIf=\"block\">\n    <div class=\"card-content\" >\n      <table>\n        <tbody >\n\n        <tr>\n\n          <td>Status</td><td>{{ transaction.status?\"confirmed\":\"unconfirmed\" }}</td>\n\n        </tr>\n        <tr>\n          <td>Timestamp</td><td>{{getTime(block.header.timestamp)}}</td>\n        </tr>\n        <tr>\n          <td>Block </td><td><a href=\"block/{{ transaction.status.block_seq }}\">{{ transaction.status.block_seq }}</a></td>\n        </tr>\n        </tbody>\n      </table>\n    </div>\n  </div>\n\n\n  <div class=\" card col s12\">\n    <div class=\"card-title\">Transaction id : {{transaction.txn.txid}}</div>\n    <div class=\"col s4\">\n      <ul class=\"collection\">\n        <li class=\"collection-item\"><h5>Inputs</h5></li>\n        <li class=\"collection-item\" *ngFor=\"let input of transaction.txn.inputs\">{{input}}</li>\n\n      </ul>\n    </div>\n    <div class=\"accent-1 col s1 right-arrow center\">\n      <i class=\"material-icons\">trending_flat</i>\n    </div>\n\n    <div class=\"col s5\">\n      <ul class=\"collection\">\n        <li class=\"collection-item\"><h5>Outputs</h5></li>\n        <li class=\"collection-item\" *ngFor=\"let output of transaction.txn.outputs\"><a href=\"/address/{{output.dst}}\">{{output.dst}}</a></li>\n      </ul>\n    </div>\n    <div class=\"col s2\">\n      <ul class=\"collection\">\n        <li class=\"collection-item\"><h5>Coins</h5></li>\n        <li class=\"collection-item\" *ngFor=\"let output of transaction.txn.outputs\">{{output.coins}}</li>\n      </ul>\n    </div>\n  </div>\n</div>\n\n"

/***/ })

},[1121]);
//# sourceMappingURL=main.bundle.map