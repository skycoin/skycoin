webpackJsonp([0,3],{

/***/ 1137:
/***/ (function(module, exports, __webpack_require__) {

module.exports = __webpack_require__(561);


/***/ }),

/***/ 174:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__(91);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_map__ = __webpack_require__(534);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_map___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_map__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_catch__ = __webpack_require__(533);
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
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/block-chain.service.js.map

/***/ }),

/***/ 365:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_router__ = __webpack_require__(78);
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
    function AppComponent(router) {
        var _this = this;
        this.router = router;
        router.events.subscribe(function (event) {
            _this.navigationInterceptor(event);
        });
    }
    // Shows and hides the loading spinner during RouterEvent changes
    AppComponent.prototype.navigationInterceptor = function (event) {
        if (event instanceof __WEBPACK_IMPORTED_MODULE_1__angular_router__["c" /* NavigationStart */]) {
            console.log("Navigation -start");
            this.loading = true;
        }
        if (event instanceof __WEBPACK_IMPORTED_MODULE_1__angular_router__["d" /* NavigationEnd */]) {
            console.log("Navigation -end");
            this.loading = false;
        }
        // Set loading state to false in both of the below events to hide the spinner in case a request fails
        if (event instanceof __WEBPACK_IMPORTED_MODULE_1__angular_router__["e" /* NavigationCancel */]) {
            console.log("Navigation -canceled");
            this.loading = false;
        }
        if (event instanceof __WEBPACK_IMPORTED_MODULE_1__angular_router__["f" /* NavigationError */]) {
            console.log("Navigation -error");
            this.loading = false;
        }
    };
    AppComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-root',
            template: __webpack_require__(862),
            styles: [__webpack_require__(846)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_router__["b" /* Router */]) === 'function' && _a) || Object])
    ], AppComponent);
    return AppComponent;
    var _a;
}());
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/app.component.js.map

/***/ }),

/***/ 366:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__(91);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs__ = __webpack_require__(185);
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
    UxOutputsService.prototype.getCurrentBalanceOfAddress = function (address) {
        return this._http.get('/api/currentBalance?address=' + address)
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
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/UxOutputs.service.js.map

/***/ }),

/***/ 367:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__(91);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs__ = __webpack_require__(185);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return CoinSupplyService; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var CoinSupplyService = (function () {
    function CoinSupplyService(_http) {
        this._http = _http;
    }
    CoinSupplyService.prototype.getCoinSupply = function () {
        return this._http.get('/api/coinSupply')
            .map(function (res) {
            return res.json();
        })
            .catch(function (error) {
            console.log(error);
            return __WEBPACK_IMPORTED_MODULE_2_rxjs__["Observable"].throw(error || 'Server error');
        });
    };
    CoinSupplyService = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Injectable"])(), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Http */]) === 'function' && _a) || Object])
    ], CoinSupplyService);
    return CoinSupplyService;
    var _a;
}());
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/coin-supply.service.js.map

/***/ }),

/***/ 368:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__(91);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs__ = __webpack_require__(185);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs__);
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
        var pages = _.range(startPage, endPage + 1);
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
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/skycoin-blockchain-pagination.service.js.map

/***/ }),

/***/ 369:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__(91);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs__ = __webpack_require__(185);
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
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/transaction-detail.service.js.map

/***/ }),

/***/ 560:
/***/ (function(module, exports) {

function webpackEmptyContext(req) {
	throw new Error("Cannot find module '" + req + "'.");
}
webpackEmptyContext.keys = function() { return []; };
webpackEmptyContext.resolve = webpackEmptyContext;
module.exports = webpackEmptyContext;
webpackEmptyContext.id = 560;


/***/ }),

/***/ 561:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__polyfills_ts__ = __webpack_require__(693);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__ = __webpack_require__(649);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__environments_environment__ = __webpack_require__(692);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__app___ = __webpack_require__(691);





if (__WEBPACK_IMPORTED_MODULE_3__environments_environment__["a" /* environment */].production) {
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_2__angular_core__["enableProdMode"])();
}
__webpack_require__.i(__WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__["a" /* platformBrowserDynamic */])().bootstrapModule(__WEBPACK_IMPORTED_MODULE_4__app___["a" /* AppModule */]);
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/main.js.map

/***/ }),

/***/ 679:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__ = __webpack_require__(171);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_forms__ = __webpack_require__(640);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_http__ = __webpack_require__(91);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__angular_router__ = __webpack_require__(78);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__app_component__ = __webpack_require__(365);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6__components_skycoin_header_skycoin_header_component__ = __webpack_require__(686);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7__components_block_chain_table_block_chain_table_component__ = __webpack_require__(682);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8__components_skycoin_pagination_skycoin_pagination_component__ = __webpack_require__(688);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9__components_skycoin_search_bar_skycoin_search_bar_component__ = __webpack_require__(689);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10__components_block_chain_table_block_chain_service__ = __webpack_require__(174);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11__components_footer_footer_component__ = __webpack_require__(684);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12__components_skycoin_pagination_num_pages_pipe__ = __webpack_require__(687);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13__components_skycoin_pagination_skycoin_blockchain_pagination_service__ = __webpack_require__(368);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14__angular_common__ = __webpack_require__(101);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_angular2_moment__ = __webpack_require__(694);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_angular2_moment___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_15_angular2_moment__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_16__components_block_details_block_details_component__ = __webpack_require__(683);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_17__components_address_detail_address_detail_component__ = __webpack_require__(680);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_18__components_address_detail_UxOutputs_service__ = __webpack_require__(366);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_19__components_transaction_detail_transaction_detail_component__ = __webpack_require__(690);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_20__components_transaction_detail_transaction_detail_service__ = __webpack_require__(369);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_21__components_loading_loading_component__ = __webpack_require__(685);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_22__components_block_chain_coin_supply_block_chain_coin_supply_component__ = __webpack_require__(681);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_23__components_block_chain_coin_supply_coin_supply_service__ = __webpack_require__(367);
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
























// import {QRCodeModule} from "../js/angular2-qrcode";
// import {QRCodeModule} from "angular2-qrcode";
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
                __WEBPACK_IMPORTED_MODULE_19__components_transaction_detail_transaction_detail_component__["a" /* TransactionDetailComponent */],
                __WEBPACK_IMPORTED_MODULE_21__components_loading_loading_component__["a" /* LoadingComponent */],
                __WEBPACK_IMPORTED_MODULE_22__components_block_chain_coin_supply_block_chain_coin_supply_component__["a" /* BlockChainCoinSupplyComponent */]
            ],
            imports: [
                __WEBPACK_IMPORTED_MODULE_14__angular_common__["a" /* CommonModule */],
                __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__["a" /* BrowserModule */],
                __WEBPACK_IMPORTED_MODULE_2__angular_forms__["a" /* FormsModule */],
                __WEBPACK_IMPORTED_MODULE_3__angular_http__["a" /* HttpModule */],
                __WEBPACK_IMPORTED_MODULE_15_angular2_moment__["MomentModule"],
                __WEBPACK_IMPORTED_MODULE_4__angular_router__["a" /* RouterModule */].forRoot(ROUTES)
            ],
            providers: [__WEBPACK_IMPORTED_MODULE_10__components_block_chain_table_block_chain_service__["a" /* BlockChainService */], __WEBPACK_IMPORTED_MODULE_13__components_skycoin_pagination_skycoin_blockchain_pagination_service__["a" /* SkycoinBlockchainPaginationService */], __WEBPACK_IMPORTED_MODULE_18__components_address_detail_UxOutputs_service__["a" /* UxOutputsService */], __WEBPACK_IMPORTED_MODULE_20__components_transaction_detail_transaction_detail_service__["a" /* TransactionDetailService */], __WEBPACK_IMPORTED_MODULE_23__components_block_chain_coin_supply_coin_supply_service__["a" /* CoinSupplyService */]],
            bootstrap: [__WEBPACK_IMPORTED_MODULE_5__app_component__["a" /* AppComponent */]]
        }), 
        __metadata('design:paramtypes', [])
    ], AppModule);
    return AppModule;
}());
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/app.module.js.map

/***/ }),

/***/ 680:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__UxOutputs_service__ = __webpack_require__(366);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__(78);
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
        this.currentBalance = "0";
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
                _this.currentBalance = addressDetails.head_outputs[0].coins;
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
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-address-detail',
            template: __webpack_require__(863),
            styles: [__webpack_require__(847)],
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__UxOutputs_service__["a" /* UxOutputsService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__UxOutputs_service__["a" /* UxOutputsService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["g" /* ActivatedRoute */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__angular_router__["g" /* ActivatedRoute */]) === 'function' && _b) || Object, (typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */]) === 'function' && _c) || Object])
    ], AddressDetailComponent);
    return AddressDetailComponent;
    var _a, _b, _c;
}());
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/address-detail.component.js.map

/***/ }),

/***/ 681:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__coin_supply_service__ = __webpack_require__(367);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return BlockChainCoinSupplyComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var BlockChainCoinSupplyComponent = (function () {
    function BlockChainCoinSupplyComponent(service) {
        this.service = service;
        this.coinSupply = this.coinCap = 0;
    }
    BlockChainCoinSupplyComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.service.getCoinSupply().subscribe(function (supply) {
            _this.coinCap = supply.coinCap;
            _this.coinSupply = supply.coinSupply;
        });
    };
    BlockChainCoinSupplyComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-block-chain-coin-supply',
            template: __webpack_require__(864),
            styles: [__webpack_require__(848)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__coin_supply_service__["a" /* CoinSupplyService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__coin_supply_service__["a" /* CoinSupplyService */]) === 'function' && _a) || Object])
    ], BlockChainCoinSupplyComponent);
    return BlockChainCoinSupplyComponent;
    var _a;
}());
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/block-chain-coin-supply.component.js.map

/***/ }),

/***/ 682:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__block_chain_service__ = __webpack_require__(174);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_moment__ = __webpack_require__(2);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_moment___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_moment__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_router__ = __webpack_require__(78);
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
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-block-chain-table',
            template: __webpack_require__(865),
            styles: [__webpack_require__(849)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__block_chain_service__["a" /* BlockChainService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__block_chain_service__["a" /* BlockChainService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_3__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_3__angular_router__["b" /* Router */]) === 'function' && _b) || Object])
    ], BlockChainTableComponent);
    return BlockChainTableComponent;
    var _a, _b;
}());
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/block-chain-table.component.js.map

/***/ }),

/***/ 683:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_router__ = __webpack_require__(78);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__block_chain_table_block_chain_service__ = __webpack_require__(174);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_moment__ = __webpack_require__(2);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_moment___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_moment__);
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
        return __WEBPACK_IMPORTED_MODULE_3_moment__["unix"](time).format();
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
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-block-details',
            template: __webpack_require__(866),
            styles: [__webpack_require__(850)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_2__block_chain_table_block_chain_service__["a" /* BlockChainService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__block_chain_table_block_chain_service__["a" /* BlockChainService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__angular_router__["g" /* ActivatedRoute */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_router__["g" /* ActivatedRoute */]) === 'function' && _b) || Object, (typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_1__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_router__["b" /* Router */]) === 'function' && _c) || Object])
    ], BlockDetailsComponent);
    return BlockDetailsComponent;
    var _a, _b, _c;
}());
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/block-details.component.js.map

/***/ }),

/***/ 684:
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
            template: __webpack_require__(867),
            styles: [__webpack_require__(851)]
        }), 
        __metadata('design:paramtypes', [])
    ], FooterComponent);
    return FooterComponent;
}());
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/footer.component.js.map

/***/ }),

/***/ 685:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return LoadingComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};

var LoadingComponent = (function () {
    function LoadingComponent() {
    }
    LoadingComponent.prototype.ngOnInit = function () {
        this.loading = false;
    };
    __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"])(), 
        __metadata('design:type', Boolean)
    ], LoadingComponent.prototype, "loading", void 0);
    LoadingComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-loading',
            template: __webpack_require__(868),
            styles: [__webpack_require__(852)]
        }), 
        __metadata('design:paramtypes', [])
    ], LoadingComponent);
    return LoadingComponent;
}());
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/loading.component.js.map

/***/ }),

/***/ 686:
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
            template: __webpack_require__(869),
            styles: [__webpack_require__(853)]
        }), 
        __metadata('design:paramtypes', [])
    ], SkycoinHeaderComponent);
    return SkycoinHeaderComponent;
}());
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/skycoin-header.component.js.map

/***/ }),

/***/ 687:
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
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/num-pages.pipe.js.map

/***/ }),

/***/ 688:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__skycoin_blockchain_pagination_service__ = __webpack_require__(368);
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
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Output"])(), 
        __metadata('design:type', Object)
    ], SkycoinPaginationComponent.prototype, "onChangePage", void 0);
    SkycoinPaginationComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-skycoin-pagination',
            template: __webpack_require__(870),
            styles: [__webpack_require__(854)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__skycoin_blockchain_pagination_service__["a" /* SkycoinBlockchainPaginationService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__skycoin_blockchain_pagination_service__["a" /* SkycoinBlockchainPaginationService */]) === 'function' && _a) || Object])
    ], SkycoinPaginationComponent);
    return SkycoinPaginationComponent;
    var _a;
}());
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/skycoin-pagination.component.js.map

/***/ }),

/***/ 689:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__block_chain_table_block_chain_service__ = __webpack_require__(174);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__(78);
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
            template: __webpack_require__(871),
            styles: [__webpack_require__(855)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__block_chain_table_block_chain_service__["a" /* BlockChainService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__block_chain_table_block_chain_service__["a" /* BlockChainService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */]) === 'function' && _b) || Object])
    ], SkycoinSearchBarComponent);
    return SkycoinSearchBarComponent;
    var _a, _b;
}());
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/skycoin-search-bar.component.js.map

/***/ }),

/***/ 690:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_router__ = __webpack_require__(78);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__transaction_detail_service__ = __webpack_require__(369);
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
            template: __webpack_require__(872),
            styles: [__webpack_require__(856)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_2__transaction_detail_service__["a" /* TransactionDetailService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__transaction_detail_service__["a" /* TransactionDetailService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__angular_router__["g" /* ActivatedRoute */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_router__["g" /* ActivatedRoute */]) === 'function' && _b) || Object, (typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_1__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_router__["b" /* Router */]) === 'function' && _c) || Object])
    ], TransactionDetailComponent);
    return TransactionDetailComponent;
    var _a, _b, _c;
}());
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/transaction-detail.component.js.map

/***/ }),

/***/ 691:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__app_component__ = __webpack_require__(365);
/* unused harmony namespace reexport */
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__app_module__ = __webpack_require__(679);
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "a", function() { return __WEBPACK_IMPORTED_MODULE_1__app_module__["a"]; });


//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/index.js.map

/***/ }),

/***/ 692:
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
//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/environment.js.map

/***/ }),

/***/ 693:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_core_js_es6_symbol__ = __webpack_require__(709);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_core_js_es6_symbol___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_core_js_es6_symbol__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_core_js_es6_object__ = __webpack_require__(702);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_core_js_es6_object___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_core_js_es6_object__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_core_js_es6_function__ = __webpack_require__(698);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_core_js_es6_function___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_core_js_es6_function__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_core_js_es6_parse_int__ = __webpack_require__(704);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_core_js_es6_parse_int___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_core_js_es6_parse_int__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_core_js_es6_parse_float__ = __webpack_require__(703);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_core_js_es6_parse_float___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4_core_js_es6_parse_float__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_core_js_es6_number__ = __webpack_require__(701);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_core_js_es6_number___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_5_core_js_es6_number__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6_core_js_es6_math__ = __webpack_require__(700);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6_core_js_es6_math___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_6_core_js_es6_math__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7_core_js_es6_string__ = __webpack_require__(708);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7_core_js_es6_string___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_7_core_js_es6_string__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8_core_js_es6_date__ = __webpack_require__(697);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8_core_js_es6_date___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_8_core_js_es6_date__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9_core_js_es6_array__ = __webpack_require__(696);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9_core_js_es6_array___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_9_core_js_es6_array__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10_core_js_es6_regexp__ = __webpack_require__(706);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10_core_js_es6_regexp___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_10_core_js_es6_regexp__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11_core_js_es6_map__ = __webpack_require__(699);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11_core_js_es6_map___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_11_core_js_es6_map__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12_core_js_es6_set__ = __webpack_require__(707);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12_core_js_es6_set___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_12_core_js_es6_set__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13_core_js_es6_reflect__ = __webpack_require__(705);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13_core_js_es6_reflect___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_13_core_js_es6_reflect__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14_core_js_es7_reflect__ = __webpack_require__(710);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14_core_js_es7_reflect___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_14_core_js_es7_reflect__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_zone_js_dist_zone__ = __webpack_require__(1136);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_zone_js_dist_zone___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_15_zone_js_dist_zone__);
















//# sourceMappingURL=/Users/napandey/go/src/github.com/skycoin/skycoin/skycoin-web/src/polyfills.js.map

/***/ }),

/***/ 846:
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(40)();
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ 847:
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(40)();
// imports


// module
exports.push([module.i, ".collection{\n  border: none;\n}\na{\n  word-wrap: break-word;\n}\n.collection-item{\n  border: none;\n}\n\n.card-title{\n  font-size: 16px;\n}\n.right-arrow{\n  margin-top: 50px;\n}\n\n.address-qr{\n  padding-top: 20px;\n}\n\n.modalShow{\n  z-index: 1003;\n  display: block;\n  opacity: 1;\n  max-width: 80%;\n  margin: 0 auto;\n  max-height: 25%;\n  -webkit-transform: scaleX(1);\n          transform: scaleX(1);\n  top: 10%;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ 848:
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(40)();
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ 849:
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(40)();
// imports


// module
exports.push([module.i, "tr:hover{\n  cursor: pointer;\n}\n\n.wrap{\n  word-break: break-all;\n}\n\n\n\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ 850:
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(40)();
// imports


// module
exports.push([module.i, ".card{\n  background-color: #29b6f6;\n  color:white;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ 851:
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(40)();
// imports


// module
exports.push([module.i, "footer{\n  background-color: #29b6f6;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ 852:
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(40)();
// imports


// module
exports.push([module.i, ".modalShow{\n  z-index: 1003;\n  display: block;\n  opacity: 1;\n  max-width: 80%;\n  margin: 0 auto;\n  -webkit-transform: scaleX(1);\n          transform: scaleX(1);\n  top: 10%;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ 853:
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(40)();
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ 854:
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(40)();
// imports


// module
exports.push([module.i, ".pagination li.active{\n  background-color: #29b6f6;\n  color:white;\n}\n\n.lower-panel{\n  color:black;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ 855:
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(40)();
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ 856:
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__(40)();
// imports


// module
exports.push([module.i, ".collection{\n  border: none;\n}\n.collection-item{\n  border: none;\n}\n\n.card-title{\n  font-size: 16px;\n}\n.right-arrow{\n  margin-top: 50px;\n}\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ 857:
/***/ (function(module, exports, __webpack_require__) {

var map = {
	"./af": 416,
	"./af.js": 416,
	"./ar": 423,
	"./ar-dz": 417,
	"./ar-dz.js": 417,
	"./ar-kw": 418,
	"./ar-kw.js": 418,
	"./ar-ly": 419,
	"./ar-ly.js": 419,
	"./ar-ma": 420,
	"./ar-ma.js": 420,
	"./ar-sa": 421,
	"./ar-sa.js": 421,
	"./ar-tn": 422,
	"./ar-tn.js": 422,
	"./ar.js": 423,
	"./az": 424,
	"./az.js": 424,
	"./be": 425,
	"./be.js": 425,
	"./bg": 426,
	"./bg.js": 426,
	"./bn": 427,
	"./bn.js": 427,
	"./bo": 428,
	"./bo.js": 428,
	"./br": 429,
	"./br.js": 429,
	"./bs": 430,
	"./bs.js": 430,
	"./ca": 431,
	"./ca.js": 431,
	"./cs": 432,
	"./cs.js": 432,
	"./cv": 433,
	"./cv.js": 433,
	"./cy": 434,
	"./cy.js": 434,
	"./da": 435,
	"./da.js": 435,
	"./de": 438,
	"./de-at": 436,
	"./de-at.js": 436,
	"./de-ch": 437,
	"./de-ch.js": 437,
	"./de.js": 438,
	"./dv": 439,
	"./dv.js": 439,
	"./el": 440,
	"./el.js": 440,
	"./en-au": 441,
	"./en-au.js": 441,
	"./en-ca": 442,
	"./en-ca.js": 442,
	"./en-gb": 443,
	"./en-gb.js": 443,
	"./en-ie": 444,
	"./en-ie.js": 444,
	"./en-nz": 445,
	"./en-nz.js": 445,
	"./eo": 446,
	"./eo.js": 446,
	"./es": 448,
	"./es-do": 447,
	"./es-do.js": 447,
	"./es.js": 448,
	"./et": 449,
	"./et.js": 449,
	"./eu": 450,
	"./eu.js": 450,
	"./fa": 451,
	"./fa.js": 451,
	"./fi": 452,
	"./fi.js": 452,
	"./fo": 453,
	"./fo.js": 453,
	"./fr": 456,
	"./fr-ca": 454,
	"./fr-ca.js": 454,
	"./fr-ch": 455,
	"./fr-ch.js": 455,
	"./fr.js": 456,
	"./fy": 457,
	"./fy.js": 457,
	"./gd": 458,
	"./gd.js": 458,
	"./gl": 459,
	"./gl.js": 459,
	"./gom-latn": 460,
	"./gom-latn.js": 460,
	"./he": 461,
	"./he.js": 461,
	"./hi": 462,
	"./hi.js": 462,
	"./hr": 463,
	"./hr.js": 463,
	"./hu": 464,
	"./hu.js": 464,
	"./hy-am": 465,
	"./hy-am.js": 465,
	"./id": 466,
	"./id.js": 466,
	"./is": 467,
	"./is.js": 467,
	"./it": 468,
	"./it.js": 468,
	"./ja": 469,
	"./ja.js": 469,
	"./jv": 470,
	"./jv.js": 470,
	"./ka": 471,
	"./ka.js": 471,
	"./kk": 472,
	"./kk.js": 472,
	"./km": 473,
	"./km.js": 473,
	"./kn": 474,
	"./kn.js": 474,
	"./ko": 475,
	"./ko.js": 475,
	"./ky": 476,
	"./ky.js": 476,
	"./lb": 477,
	"./lb.js": 477,
	"./lo": 478,
	"./lo.js": 478,
	"./lt": 479,
	"./lt.js": 479,
	"./lv": 480,
	"./lv.js": 480,
	"./me": 481,
	"./me.js": 481,
	"./mi": 482,
	"./mi.js": 482,
	"./mk": 483,
	"./mk.js": 483,
	"./ml": 484,
	"./ml.js": 484,
	"./mr": 485,
	"./mr.js": 485,
	"./ms": 487,
	"./ms-my": 486,
	"./ms-my.js": 486,
	"./ms.js": 487,
	"./my": 488,
	"./my.js": 488,
	"./nb": 489,
	"./nb.js": 489,
	"./ne": 490,
	"./ne.js": 490,
	"./nl": 492,
	"./nl-be": 491,
	"./nl-be.js": 491,
	"./nl.js": 492,
	"./nn": 493,
	"./nn.js": 493,
	"./pa-in": 494,
	"./pa-in.js": 494,
	"./pl": 495,
	"./pl.js": 495,
	"./pt": 497,
	"./pt-br": 496,
	"./pt-br.js": 496,
	"./pt.js": 497,
	"./ro": 498,
	"./ro.js": 498,
	"./ru": 499,
	"./ru.js": 499,
	"./sd": 500,
	"./sd.js": 500,
	"./se": 501,
	"./se.js": 501,
	"./si": 502,
	"./si.js": 502,
	"./sk": 503,
	"./sk.js": 503,
	"./sl": 504,
	"./sl.js": 504,
	"./sq": 505,
	"./sq.js": 505,
	"./sr": 507,
	"./sr-cyrl": 506,
	"./sr-cyrl.js": 506,
	"./sr.js": 507,
	"./ss": 508,
	"./ss.js": 508,
	"./sv": 509,
	"./sv.js": 509,
	"./sw": 510,
	"./sw.js": 510,
	"./ta": 511,
	"./ta.js": 511,
	"./te": 512,
	"./te.js": 512,
	"./tet": 513,
	"./tet.js": 513,
	"./th": 514,
	"./th.js": 514,
	"./tl-ph": 515,
	"./tl-ph.js": 515,
	"./tlh": 516,
	"./tlh.js": 516,
	"./tr": 517,
	"./tr.js": 517,
	"./tzl": 518,
	"./tzl.js": 518,
	"./tzm": 520,
	"./tzm-latn": 519,
	"./tzm-latn.js": 519,
	"./tzm.js": 520,
	"./uk": 521,
	"./uk.js": 521,
	"./ur": 522,
	"./ur.js": 522,
	"./uz": 524,
	"./uz-latn": 523,
	"./uz-latn.js": 523,
	"./uz.js": 524,
	"./vi": 525,
	"./vi.js": 525,
	"./x-pseudo": 526,
	"./x-pseudo.js": 526,
	"./yo": 527,
	"./yo.js": 527,
	"./zh-cn": 528,
	"./zh-cn.js": 528,
	"./zh-hk": 529,
	"./zh-hk.js": 529,
	"./zh-tw": 530,
	"./zh-tw.js": 530
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
webpackContext.id = 857;


/***/ }),

/***/ 862:
/***/ (function(module, exports) {

module.exports = "<app-skycoin-header></app-skycoin-header>\n<div class=\"container\">\n  <div class=\"row center\">\n    <app-skycoin-search-bar></app-skycoin-search-bar>\n    <router-outlet></router-outlet>\n  </div>\n</div>\n<app-loading [loading]=\"loading\"></app-loading>\n<app-footer></app-footer>\n\n\n\n\n"

/***/ }),

/***/ 863:
/***/ (function(module, exports) {

module.exports = "<div class=\"row\" *ngIf=\"transactions\">\n  <div class=\"col s8\" >\n    <div class=\"card-content\" *ngIf=\"transactions.length>0\">\n      <ul class=\"collection\">\n        <li class=\"collection-item\"><h5>Address</h5><span>{{currentAddress}}</span></li>\n        <li class=\"collection-item\">Number of transactions: {{transactions.length}} </li>\n        <li class=\"collection-item\">Current value: <i class=\"material-icons\">ic_account_balance</i> {{currentBalance}} skycoins</li>\n      </ul>\n    </div>\n  </div>\n\n  <div class=\"card col s4\">\n    <!--<qr-code *ngIf=\"currentAddress!=null\" [value]=\"currentAddress\" [size]=\"300\"></qr-code>-->\n    <div id=\"qr-code\"></div>\n    <p class=\"address-qr\">{{currentAddress}}</p>\n  </div>\n\n\n  <div class=\" card col s12\" *ngFor=\"let transaction of transactions\">\n    <div class=\"card-title\">Transaction id : <a href=\"/transaction/{{transaction.txid}}\">{{transaction.txid}}</a></div>\n      <div class=\"col s4\">\n        <ul class=\"collection\">\n          <li class=\"collection-item\"><h5>inputs</h5></li>\n          <li class=\"collection-item\" *ngFor=\"let input of transaction.inputs\">\n            <!--<a href=\"/address/{{input.uxid}}\">uxid</a>-->\n            <a class=\"tooltipped\" data-position=\"bottom\" data-delay=\"250\" attr.data-tooltip=\"UxId: {{input.uxid}}\" href=\"/address/{{input.owner}}\">{{input.owner}}</a>\n            <i (click)=\"showUxId()\" class=\"material-icons\">ic_info_outline</i>\n            <div class=\"modal bottom-sheet\" [ngClass]=\"{modalShow:showUxID}\">\n              <div class=\"modal-content\">\n                <h4>UxId</h4>\n                {{input.uxid}}\n              </div>\n              <div class=\"modal-footer\">\n                <a href=\"#!\" class=\" modal-action modal-close waves-effect waves-green btn-flat\" (click)=\"hideUxId()\">Close</a>\n              </div>\n            </div>\n          </li>\n\n        </ul>\n      </div>\n      <div class=\"accent-1 col s1 right-arrow center\">\n        <i class=\"material-icons\">trending_flat</i>\n      </div>\n\n      <div class=\"col s5\">\n        <ul class=\"collection\">\n          <li class=\"collection-item\"><h5>outputs</h5></li>\n          <li class=\"collection-item\" *ngFor=\"let output of transaction.outputs\"><a href=\"/address/{{output.dst}}\">{{output.dst}}</a></li>\n        </ul>\n      </div>\n      <div class=\"col s2\">\n        <ul class=\"collection\">\n          <li class=\"collection-item\"><h5>coins</h5></li>\n          <li class=\"collection-item\" *ngFor=\"let output of transaction.outputs\">{{output.coins}}</li>\n        </ul>\n      </div>\n  </div>\n</div>\n\n\n"

/***/ }),

/***/ 864:
/***/ (function(module, exports) {

module.exports = "<h6>Coin supply : {{coinSupply}} and Coin Cap is : {{ coinCap}}</h6>\n"

/***/ }),

/***/ 865:
/***/ (function(module, exports) {

module.exports = "\n\n\n<div class=\"row\">\n  <div class=\"col s12 s12\">\n    <div class=\"card blue lighten-1\">\n      <div class=\"card-content white-text\">\n        <span class=\"card-title\">Blockchain</span>\n        <app-block-chain-coin-supply></app-block-chain-coin-supply>\n        <h5>{{totalBlocks}} blocks</h5>\n      </div>\n    </div>\n  </div>\n</div>\n\n\n\n\n\n<table class=\"bordered striped centered responsive-table\">\n  <thead >\n  <tr>\n    <th data-field=\"id\">Time</th>\n    <th data-field=\"name\">Block Number</th>\n    <th data-field=\"price\">Transactions</th>\n    <th data-field=\"id\">Amount Transferred</th>\n    <th data-field=\"price\">Blockhash</th>\n\n  </tr>\n  </thead>\n\n  <tbody>\n  <tr *ngFor=\"let block of blocks\" (click)=\"showDetails(block)\">\n    <td>{{getTime(block.header.timestamp)}}</td>\n    <td >{{ block.header.seq }}</td>\n    <td>{{ block.body.txns.length }}</td>\n    <td>{{ GetBlockAmount(block.body.txns) }}</td>\n    <td>{{ block.header.block_hash }}</td>\n  </tr>\n\n\n  </tbody>\n</table>\n<app-loading [loading]=\"loading\"></app-loading>\n\n\n<app-skycoin-pagination (onChangePage)=\"handlePageChange($event)\"></app-skycoin-pagination>\n"

/***/ }),

/***/ 866:
/***/ (function(module, exports) {

module.exports = "<div class=\"row\">\n  <h2>Block Details</h2>\n  <div class=\"card\" *ngIf=\"block\">\n    <div class=\"card-content\" >\n      <table>\n        <tbody >\n\n        <tr>\n\n          <td>Height</td><td>{{ block.header.seq }}</td>\n\n        </tr>\n        <tr>\n          <td>Timestamp</td><td>{{getTime(block.header.timestamp)}}</td>\n        </tr>\n        <tr>\n          <td>Hash</td><td>{{ block.header.block_hash }}</td>\n        </tr>\n\n        <tr>\n          <td>Parent Hash</td><td>{{ block.header.previous_block_hash }}</td>\n        </tr>\n\n        <tr>\n          <td>Total Amount</td><td>{{ getAmount(block) }} coins</td>\n        </tr>\n\n        </tbody>\n      </table>\n    </div>\n  </div>\n\n\n    <div class=\"col s12\">Transactions</div>\n    <div class=\"col s12\">\n      <div class=\"col s1\">S.no</div>\n      <div class=\"col s2\">Transaction Id</div>\n      <div class=\"col s3\">Inputs</div>\n      <div class=\"col s4\">Outputs</div>\n      <div class=\"col s2\">Amount</div>\n    </div>\n  <div *ngIf=\"block\">\n    <div  class=\"col s12\" *ngFor=\"let transaction of block.body.txns;let i = index\">\n\n      <div class=\"col s1\">{{i+1}}</div>\n      <div class=\"col s2 wrap\"><a [routerLink]=\"['/transaction', transaction.txid ]\" >{{transaction.txid}}</a></div>\n      <div class=\"col s3 wrap\"><p *ngFor=\"let input of transaction.inputs;\" class=\"input\">{{input}}</p></div>\n      <div class=\"col s6 wrap\"><p *ngFor=\"let output of transaction.outputs;\" class=\"input\"><a [routerLink]=\"['/address', output.dst ]\" >{{output.dst}}</a>  - {{output.coins}}</p></div>\n    </div>\n  </div>\n\n\n\n\n</div>\n"

/***/ }),

/***/ 867:
/***/ (function(module, exports) {

module.exports = "<footer class=\"page-footer\">\n  <div class=\"container\">\n    <div class=\"row\">\n      <div class=\"col l6 s12\">\n        <p class=\"grey-text text-lighten-4\">Next generation of digital money.</p>\n      </div>\n      <div class=\"col l4 offset-l2 s12\">\n        <h5 class=\"white-text\">Links</h5>\n        <ul>\n          <li><a class=\"grey-text text-lighten-3\" href=\"http://skycoin.net/infographics.html\">How it works?</a></li>\n          <li><a class=\"grey-text text-lighten-3\" href=\"http://skycoin.net/downloads.html\">Download wallet</a></li>\n          <li><a class=\"grey-text text-lighten-3\" href=\"http://skycoin.net/whitepapers.html\">White papers</a></li>\n          <li><a class=\"grey-text text-lighten-3\" href=\"http://skycoin.net/faq.html\">FAQs</a></li>\n        </ul>\n      </div>\n    </div>\n  </div>\n  <div class=\"footer-copyright\">\n    <div class=\"container\">\n      © 2017 Skycoin\n    </div>\n  </div>\n</footer>\n"

/***/ }),

/***/ 868:
/***/ (function(module, exports) {

module.exports = "<div id=\"loadingModal\" class=\"modal bottom-sheet\" [ngClass]=\"{modalShow:loading}\">\n  <div class=\"modal-content\">\n    <h4>Loading please wait...</h4>\n    <div class=\"preloader-wrapper big active\">\n      <div class=\"spinner-layer spinner-blue-only\">\n        <div class=\"circle-clipper left\">\n          <div class=\"circle\"></div>\n        </div><div class=\"gap-patch\">\n        <div class=\"circle\"></div>\n      </div><div class=\"circle-clipper right\">\n        <div class=\"circle\"></div>\n      </div>\n      </div>\n    </div>\n  </div>\n</div>\n"

/***/ }),

/***/ 869:
/***/ (function(module, exports) {

module.exports = "<nav class=\"light-blue lighten-1\" role=\"navigation\">\n  <div class=\"nav-wrapper container\"><a id=\"logo-container\" href=\"#\" class=\"brand-logo\">Skycoin</a>\n    <ul class=\"right hide-on-med-and-down\">\n      <li><a href=\"http://skycoin.net/infographics.html\">How it works?</a></li>\n      <li><a href=\"http://skycoin.net/downloads.html\">Download wallet</a></li>\n      <li><a href=\"http://skycoin.net/faq.html\">FAQs</a></li>\n      <li><a href=\"http://skycoin.net/whitepapers.html\">White papers</a></li>\n    </ul>\n\n    <ul id=\"nav-mobile\" class=\"side-nav\">\n      <li><a href=\"http://skycoin.net/infographics.html\">How it works?</a></li>\n      <li><a href=\"http://skycoin.net/downloads.html\">Download wallet</a></li>\n      <li><a href=\"http://skycoin.net/faq.html\">FAQs</a></li>\n      <li><a href=\"http://skycoin.net/whitepapers.html\">White papers</a></li>\n    </ul>\n    <a href=\"#\" data-activates=\"nav-mobile\" class=\"button-collapse\"><i class=\"material-icons\">menu</i></a>\n  </div>\n</nav>\n"

/***/ }),

/***/ 870:
/***/ (function(module, exports) {

module.exports = "<ul class=\"pagination\">\n\n  <li [ngClass]=\"{disabled:currentPage <= 5}\"><a href=\"#!\" (click)=\"loadPrevious()\"><i class=\"material-icons\">chevron_left</i></a></li>\n\n  <li *ngFor='let page of currentPages' [ngClass]=\"{active:currentPage == page}\" class=\"active\">\n\n    <a  (click)=\"changePage(page)\" href=\"#!\">{{page}}</a>\n\n  </li>\n  <li [ngClass]=\"{disabled:noUpcoming}\"><a href=\"#!\" (click)=\"loadUpcoming()\"><i class=\"material-icons\">chevron_right</i></a></li>\n</ul>\n"

/***/ }),

/***/ 871:
/***/ (function(module, exports) {

module.exports = "<div class=\"nav-wrapper\">\n  <form (submit)=\"searchBlockHistory(blockSearchKey.value)\">\n    <div class=\"input-field\">\n      <input #blockSearchKey id=\"search\" type=\"search\" required placeholder=\"blockhash, address, block-number, transaction id\">\n      <label class=\"label-icon\" for=\"search\"><i class=\"material-icons\">search</i></label>\n      <i class=\"material-icons\">close</i>\n    </div>\n  </form>\n</div>\n"

/***/ }),

/***/ 872:
/***/ (function(module, exports) {

module.exports = "<div class=\"row\" *ngIf=\"transaction!=null\">\n  <h2>Transaction Details</h2>\n  <div class=\"card\" *ngIf=\"block\">\n    <div class=\"card-content\" >\n      <table>\n        <tbody >\n\n        <tr>\n\n          <td>Status</td><td>{{ transaction.status?\"confirmed\":\"unconfirmed\" }}</td>\n\n        </tr>\n        <tr>\n          <td>Timestamp</td><td>{{getTime(block.header.timestamp)}}</td>\n        </tr>\n        <tr>\n          <td>Block </td><td><a [routerLink]=\"['/block', transaction.status.block_seq ]\">{{ transaction.status.block_seq }}</a></td>\n        </tr>\n        </tbody>\n      </table>\n    </div>\n  </div>\n\n\n  <div class=\" card col s12\">\n    <div class=\"card-title\">Transaction id : {{transaction.txn.txid}}</div>\n    <div class=\"col s4\">\n      <ul class=\"collection\">\n        <li class=\"collection-item\"><h5>Inputs</h5></li>\n        <li class=\"collection-item\" *ngFor=\"let input of transaction.txn.inputs\">{{input}}</li>\n\n      </ul>\n    </div>\n    <div class=\"accent-1 col s1 right-arrow center\">\n      <i class=\"material-icons\">trending_flat</i>\n    </div>\n\n    <div class=\"col s5\">\n      <ul class=\"collection\">\n        <li class=\"collection-item\"><h5>Outputs</h5></li>\n        <li class=\"collection-item\" *ngFor=\"let output of transaction.txn.outputs\"><a [routerLink]=\"['/address', output.dst ]\">{{output.dst}}</a></li>\n      </ul>\n    </div>\n    <div class=\"col s2\">\n      <ul class=\"collection\">\n        <li class=\"collection-item\"><h5>Coins</h5></li>\n        <li class=\"collection-item\" *ngFor=\"let output of transaction.txn.outputs\">{{output.coins}}</li>\n      </ul>\n    </div>\n  </div>\n</div>\n\n"

/***/ })

},[1137]);
//# sourceMappingURL=main.bundle.js.map