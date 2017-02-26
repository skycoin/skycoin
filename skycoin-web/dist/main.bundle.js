webpackJsonp([0,4],{

/***/ 1114:
/***/ (function(module, exports, __webpack_require__) {

module.exports = __webpack_require__(548);


/***/ }),

/***/ 172:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__(233);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_map__ = __webpack_require__(520);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_map___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_map__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_catch__ = __webpack_require__(519);
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
//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/block-chain.service.js.map

/***/ }),

/***/ 363:
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
            template: __webpack_require__(840),
            styles: [__webpack_require__(832)]
        }), 
        __metadata('design:paramtypes', [])
    ], AppComponent);
    return AppComponent;
}());
//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/app.component.js.map

/***/ }),

/***/ 364:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__(233);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs__ = __webpack_require__(848);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_underscore__ = __webpack_require__(277);
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
//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/skycoin-blockchain-pagination.service.js.map

/***/ }),

/***/ 547:
/***/ (function(module, exports) {

function webpackEmptyContext(req) {
	throw new Error("Cannot find module '" + req + "'.");
}
webpackEmptyContext.keys = function() { return []; };
webpackEmptyContext.resolve = webpackEmptyContext;
module.exports = webpackEmptyContext;
webpackEmptyContext.id = 547;


/***/ }),

/***/ 548:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__polyfills_ts__ = __webpack_require__(675);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__ = __webpack_require__(635);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__environments_environment__ = __webpack_require__(674);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__app___ = __webpack_require__(673);





if (__WEBPACK_IMPORTED_MODULE_3__environments_environment__["a" /* environment */].production) {
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_2__angular_core__["enableProdMode"])();
}
__webpack_require__.i(__WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__["a" /* platformBrowserDynamic */])().bootstrapModule(__WEBPACK_IMPORTED_MODULE_4__app___["a" /* AppModule */]);
//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/main.js.map

/***/ }),

/***/ 665:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__ = __webpack_require__(168);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_forms__ = __webpack_require__(626);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_http__ = __webpack_require__(233);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__angular_router__ = __webpack_require__(170);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__app_component__ = __webpack_require__(363);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6__components_skycoin_header_skycoin_header_component__ = __webpack_require__(669);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7__components_block_chain_table_block_chain_table_component__ = __webpack_require__(666);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8__components_skycoin_pagination_skycoin_pagination_component__ = __webpack_require__(671);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9__components_skycoin_search_bar_skycoin_search_bar_component__ = __webpack_require__(672);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10__components_block_chain_table_block_chain_service__ = __webpack_require__(172);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11__components_footer_footer_component__ = __webpack_require__(668);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12__components_skycoin_pagination_num_pages_pipe__ = __webpack_require__(670);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13__components_skycoin_pagination_skycoin_blockchain_pagination_service__ = __webpack_require__(364);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14__angular_common__ = __webpack_require__(98);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_angular2_moment__ = __webpack_require__(676);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_angular2_moment___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_15_angular2_moment__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_16__components_block_details_block_details_component__ = __webpack_require__(667);
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
                __WEBPACK_IMPORTED_MODULE_16__components_block_details_block_details_component__["a" /* BlockDetailsComponent */]
            ],
            imports: [
                __WEBPACK_IMPORTED_MODULE_14__angular_common__["a" /* CommonModule */],
                __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__["a" /* BrowserModule */],
                __WEBPACK_IMPORTED_MODULE_2__angular_forms__["a" /* FormsModule */],
                __WEBPACK_IMPORTED_MODULE_3__angular_http__["a" /* HttpModule */],
                __WEBPACK_IMPORTED_MODULE_15_angular2_moment__["MomentModule"],
                __WEBPACK_IMPORTED_MODULE_4__angular_router__["a" /* RouterModule */].forRoot(ROUTES)
            ],
            providers: [__WEBPACK_IMPORTED_MODULE_10__components_block_chain_table_block_chain_service__["a" /* BlockChainService */], __WEBPACK_IMPORTED_MODULE_13__components_skycoin_pagination_skycoin_blockchain_pagination_service__["a" /* SkycoinBlockchainPaginationService */]],
            bootstrap: [__WEBPACK_IMPORTED_MODULE_5__app_component__["a" /* AppComponent */]]
        }), 
        __metadata('design:paramtypes', [])
    ], AppModule);
    return AppModule;
}());
//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/app.module.js.map

/***/ }),

/***/ 666:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__block_chain_service__ = __webpack_require__(172);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_moment__ = __webpack_require__(2);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_moment___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_moment__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_underscore__ = __webpack_require__(277);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_underscore___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_underscore__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__angular_router__ = __webpack_require__(170);
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
        var _this = this;
        this.blockService.getBlocks(1, 10).subscribe(function (data) {
            _this.blocks = data;
        });
    };
    BlockChainTableComponent.prototype.showDetails = function (block) {
        this.router.navigate(['/block', block.header.seq]);
    };
    BlockChainTableComponent.prototype.handlePageChange = function (currentPage) {
        var _this = this;
        this.blockService.getBlocks((currentPage - 1) * 10 + 1, (currentPage - 1) * 10 + 10).subscribe(function (data) {
            console.log(data);
            _this.blocks = data;
        });
    };
    BlockChainTableComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-block-chain-table',
            template: __webpack_require__(841),
            styles: [__webpack_require__(833)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__block_chain_service__["a" /* BlockChainService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__block_chain_service__["a" /* BlockChainService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_4__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_4__angular_router__["b" /* Router */]) === 'function' && _b) || Object])
    ], BlockChainTableComponent);
    return BlockChainTableComponent;
    var _a, _b;
}());
//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/block-chain-table.component.js.map

/***/ }),

/***/ 667:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_router__ = __webpack_require__(170);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__block_chain_table_block_chain_service__ = __webpack_require__(172);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_underscore__ = __webpack_require__(277);
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
            template: __webpack_require__(842),
            styles: [__webpack_require__(834)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_2__block_chain_table_block_chain_service__["a" /* BlockChainService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__block_chain_table_block_chain_service__["a" /* BlockChainService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__angular_router__["c" /* ActivatedRoute */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_router__["c" /* ActivatedRoute */]) === 'function' && _b) || Object, (typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_1__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__angular_router__["b" /* Router */]) === 'function' && _c) || Object])
    ], BlockDetailsComponent);
    return BlockDetailsComponent;
    var _a, _b, _c;
}());
//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/block-details.component.js.map

/***/ }),

/***/ 668:
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
            template: __webpack_require__(843),
            styles: [__webpack_require__(835)]
        }), 
        __metadata('design:paramtypes', [])
    ], FooterComponent);
    return FooterComponent;
}());
//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/footer.component.js.map

/***/ }),

/***/ 669:
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
            template: __webpack_require__(844),
            styles: [__webpack_require__(836)]
        }), 
        __metadata('design:paramtypes', [])
    ], SkycoinHeaderComponent);
    return SkycoinHeaderComponent;
}());
//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/skycoin-header.component.js.map

/***/ }),

/***/ 670:
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
//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/num-pages.pipe.js.map

/***/ }),

/***/ 671:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__skycoin_blockchain_pagination_service__ = __webpack_require__(364);
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
        this.onChangePage.emit(pageNumber);
        this.currentPage = pageNumber;
        return false;
    };
    SkycoinPaginationComponent.prototype.loadUpcoming = function () {
        this.onChangePage.emit(this.currentPages[0] + this.pagesToShowAtATime);
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
        this.onChangePage.emit(this.currentPages[0] - this.pagesToShowAtATime);
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
            template: __webpack_require__(845),
            styles: [__webpack_require__(837)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__skycoin_blockchain_pagination_service__["a" /* SkycoinBlockchainPaginationService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__skycoin_blockchain_pagination_service__["a" /* SkycoinBlockchainPaginationService */]) === 'function' && _a) || Object])
    ], SkycoinPaginationComponent);
    return SkycoinPaginationComponent;
    var _a;
}());
//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/skycoin-pagination.component.js.map

/***/ }),

/***/ 672:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__(1);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__block_chain_table_block_chain_service__ = __webpack_require__(172);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__(170);
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
        var _this = this;
        this.service.getBlockByHash(hashVal).subscribe(function (block) {
            _this.block = block;
            _this.router.navigate(['/block', block.header.seq]);
        });
    };
    SkycoinSearchBarComponent = __decorate([
        __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
            selector: 'app-skycoin-search-bar',
            template: __webpack_require__(846),
            styles: [__webpack_require__(838)]
        }), 
        __metadata('design:paramtypes', [(typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__block_chain_table_block_chain_service__["a" /* BlockChainService */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_1__block_chain_table_block_chain_service__["a" /* BlockChainService */]) === 'function' && _a) || Object, (typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */] !== 'undefined' && __WEBPACK_IMPORTED_MODULE_2__angular_router__["b" /* Router */]) === 'function' && _b) || Object])
    ], SkycoinSearchBarComponent);
    return SkycoinSearchBarComponent;
    var _a, _b;
}());
//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/skycoin-search-bar.component.js.map

/***/ }),

/***/ 673:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__app_component__ = __webpack_require__(363);
/* unused harmony namespace reexport */
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__app_module__ = __webpack_require__(665);
/* harmony namespace reexport (by used) */ __webpack_require__.d(__webpack_exports__, "a", function() { return __WEBPACK_IMPORTED_MODULE_1__app_module__["a"]; });


//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/index.js.map

/***/ }),

/***/ 674:
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
//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/environment.js.map

/***/ }),

/***/ 675:
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_core_js_es6_symbol__ = __webpack_require__(691);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0_core_js_es6_symbol___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0_core_js_es6_symbol__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_core_js_es6_object__ = __webpack_require__(684);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_core_js_es6_object___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_core_js_es6_object__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_core_js_es6_function__ = __webpack_require__(680);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_core_js_es6_function___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_core_js_es6_function__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_core_js_es6_parse_int__ = __webpack_require__(686);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_core_js_es6_parse_int___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_core_js_es6_parse_int__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_core_js_es6_parse_float__ = __webpack_require__(685);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_core_js_es6_parse_float___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4_core_js_es6_parse_float__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_core_js_es6_number__ = __webpack_require__(683);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_core_js_es6_number___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_5_core_js_es6_number__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6_core_js_es6_math__ = __webpack_require__(682);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6_core_js_es6_math___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_6_core_js_es6_math__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7_core_js_es6_string__ = __webpack_require__(690);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7_core_js_es6_string___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_7_core_js_es6_string__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8_core_js_es6_date__ = __webpack_require__(679);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8_core_js_es6_date___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_8_core_js_es6_date__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9_core_js_es6_array__ = __webpack_require__(678);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9_core_js_es6_array___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_9_core_js_es6_array__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10_core_js_es6_regexp__ = __webpack_require__(688);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10_core_js_es6_regexp___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_10_core_js_es6_regexp__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11_core_js_es6_map__ = __webpack_require__(681);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11_core_js_es6_map___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_11_core_js_es6_map__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12_core_js_es6_set__ = __webpack_require__(689);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12_core_js_es6_set___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_12_core_js_es6_set__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13_core_js_es6_reflect__ = __webpack_require__(687);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13_core_js_es6_reflect___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_13_core_js_es6_reflect__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14_core_js_es7_reflect__ = __webpack_require__(692);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14_core_js_es7_reflect___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_14_core_js_es7_reflect__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_zone_js_dist_zone__ = __webpack_require__(1113);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15_zone_js_dist_zone___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_15_zone_js_dist_zone__);
















//# sourceMappingURL=/usr/local/Cellar/go/1.6/libexec/src/github.com/skycoin/skycoin/skycoin-web/src/polyfills.js.map

/***/ }),

/***/ 830:
/***/ (function(module, exports, __webpack_require__) {

var map = {
	"./af": 409,
	"./af.js": 409,
	"./ar": 415,
	"./ar-dz": 410,
	"./ar-dz.js": 410,
	"./ar-ly": 411,
	"./ar-ly.js": 411,
	"./ar-ma": 412,
	"./ar-ma.js": 412,
	"./ar-sa": 413,
	"./ar-sa.js": 413,
	"./ar-tn": 414,
	"./ar-tn.js": 414,
	"./ar.js": 415,
	"./az": 416,
	"./az.js": 416,
	"./be": 417,
	"./be.js": 417,
	"./bg": 418,
	"./bg.js": 418,
	"./bn": 419,
	"./bn.js": 419,
	"./bo": 420,
	"./bo.js": 420,
	"./br": 421,
	"./br.js": 421,
	"./bs": 422,
	"./bs.js": 422,
	"./ca": 423,
	"./ca.js": 423,
	"./cs": 424,
	"./cs.js": 424,
	"./cv": 425,
	"./cv.js": 425,
	"./cy": 426,
	"./cy.js": 426,
	"./da": 427,
	"./da.js": 427,
	"./de": 429,
	"./de-at": 428,
	"./de-at.js": 428,
	"./de.js": 429,
	"./dv": 430,
	"./dv.js": 430,
	"./el": 431,
	"./el.js": 431,
	"./en-au": 432,
	"./en-au.js": 432,
	"./en-ca": 433,
	"./en-ca.js": 433,
	"./en-gb": 434,
	"./en-gb.js": 434,
	"./en-ie": 435,
	"./en-ie.js": 435,
	"./en-nz": 436,
	"./en-nz.js": 436,
	"./eo": 437,
	"./eo.js": 437,
	"./es": 439,
	"./es-do": 438,
	"./es-do.js": 438,
	"./es.js": 439,
	"./et": 440,
	"./et.js": 440,
	"./eu": 441,
	"./eu.js": 441,
	"./fa": 442,
	"./fa.js": 442,
	"./fi": 443,
	"./fi.js": 443,
	"./fo": 444,
	"./fo.js": 444,
	"./fr": 447,
	"./fr-ca": 445,
	"./fr-ca.js": 445,
	"./fr-ch": 446,
	"./fr-ch.js": 446,
	"./fr.js": 447,
	"./fy": 448,
	"./fy.js": 448,
	"./gd": 449,
	"./gd.js": 449,
	"./gl": 450,
	"./gl.js": 450,
	"./he": 451,
	"./he.js": 451,
	"./hi": 452,
	"./hi.js": 452,
	"./hr": 453,
	"./hr.js": 453,
	"./hu": 454,
	"./hu.js": 454,
	"./hy-am": 455,
	"./hy-am.js": 455,
	"./id": 456,
	"./id.js": 456,
	"./is": 457,
	"./is.js": 457,
	"./it": 458,
	"./it.js": 458,
	"./ja": 459,
	"./ja.js": 459,
	"./jv": 460,
	"./jv.js": 460,
	"./ka": 461,
	"./ka.js": 461,
	"./kk": 462,
	"./kk.js": 462,
	"./km": 463,
	"./km.js": 463,
	"./ko": 464,
	"./ko.js": 464,
	"./ky": 465,
	"./ky.js": 465,
	"./lb": 466,
	"./lb.js": 466,
	"./lo": 467,
	"./lo.js": 467,
	"./lt": 468,
	"./lt.js": 468,
	"./lv": 469,
	"./lv.js": 469,
	"./me": 470,
	"./me.js": 470,
	"./mi": 471,
	"./mi.js": 471,
	"./mk": 472,
	"./mk.js": 472,
	"./ml": 473,
	"./ml.js": 473,
	"./mr": 474,
	"./mr.js": 474,
	"./ms": 476,
	"./ms-my": 475,
	"./ms-my.js": 475,
	"./ms.js": 476,
	"./my": 477,
	"./my.js": 477,
	"./nb": 478,
	"./nb.js": 478,
	"./ne": 479,
	"./ne.js": 479,
	"./nl": 481,
	"./nl-be": 480,
	"./nl-be.js": 480,
	"./nl.js": 481,
	"./nn": 482,
	"./nn.js": 482,
	"./pa-in": 483,
	"./pa-in.js": 483,
	"./pl": 484,
	"./pl.js": 484,
	"./pt": 486,
	"./pt-br": 485,
	"./pt-br.js": 485,
	"./pt.js": 486,
	"./ro": 487,
	"./ro.js": 487,
	"./ru": 488,
	"./ru.js": 488,
	"./se": 489,
	"./se.js": 489,
	"./si": 490,
	"./si.js": 490,
	"./sk": 491,
	"./sk.js": 491,
	"./sl": 492,
	"./sl.js": 492,
	"./sq": 493,
	"./sq.js": 493,
	"./sr": 495,
	"./sr-cyrl": 494,
	"./sr-cyrl.js": 494,
	"./sr.js": 495,
	"./ss": 496,
	"./ss.js": 496,
	"./sv": 497,
	"./sv.js": 497,
	"./sw": 498,
	"./sw.js": 498,
	"./ta": 499,
	"./ta.js": 499,
	"./te": 500,
	"./te.js": 500,
	"./tet": 501,
	"./tet.js": 501,
	"./th": 502,
	"./th.js": 502,
	"./tl-ph": 503,
	"./tl-ph.js": 503,
	"./tlh": 504,
	"./tlh.js": 504,
	"./tr": 505,
	"./tr.js": 505,
	"./tzl": 506,
	"./tzl.js": 506,
	"./tzm": 508,
	"./tzm-latn": 507,
	"./tzm-latn.js": 507,
	"./tzm.js": 508,
	"./uk": 509,
	"./uk.js": 509,
	"./uz": 510,
	"./uz.js": 510,
	"./vi": 511,
	"./vi.js": 511,
	"./x-pseudo": 512,
	"./x-pseudo.js": 512,
	"./yo": 513,
	"./yo.js": 513,
	"./zh-cn": 514,
	"./zh-cn.js": 514,
	"./zh-hk": 515,
	"./zh-hk.js": 515,
	"./zh-tw": 516,
	"./zh-tw.js": 516
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
webpackContext.id = 830;


/***/ }),

/***/ 832:
/***/ (function(module, exports) {

module.exports = ""

/***/ }),

/***/ 833:
/***/ (function(module, exports) {

module.exports = "tr:hover{\n  background-color:lightblue;\n}\n\n.wrap{\n  word-break: break-all;\n}\n"

/***/ }),

/***/ 834:
/***/ (function(module, exports) {

module.exports = ".card-content{\n  background-color: #29b6f6;\n  color:white;\n}\n"

/***/ }),

/***/ 835:
/***/ (function(module, exports) {

module.exports = "footer{\n  background-color: #29b6f6;\n}\n"

/***/ }),

/***/ 836:
/***/ (function(module, exports) {

module.exports = ""

/***/ }),

/***/ 837:
/***/ (function(module, exports) {

module.exports = ".pagination li.active{\n  background-color: #29b6f6;\n  color:white;\n}\n\n.lower-panel{\n  color:black;\n}\n"

/***/ }),

/***/ 838:
/***/ (function(module, exports) {

module.exports = ""

/***/ }),

/***/ 840:
/***/ (function(module, exports) {

module.exports = "<app-skycoin-header></app-skycoin-header>\n<div class=\"container\">\n  <div class=\"row center\">\n    <app-skycoin-search-bar></app-skycoin-search-bar>\n    <router-outlet></router-outlet>\n  </div>\n</div>\n<app-footer></app-footer>\n\n\n\n\n"

/***/ }),

/***/ 841:
/***/ (function(module, exports) {

module.exports = "<h2>Blocks</h2>\n<table>\n  <thead>\n  <tr>\n    <th data-field=\"id\">Time</th>\n    <th data-field=\"name\">Block Number</th>\n    <th data-field=\"price\">Transactions</th>\n    <th data-field=\"id\">Amount Transferred</th>\n    <th data-field=\"price\">block hash</th>\n\n  </tr>\n  </thead>\n\n  <tbody>\n  <tr *ngFor=\"let block of blocks\" >\n    <td>{{getTime(block.header.timestamp)}}</td>\n    <td (click)=\"showDetails(block)\">{{ block.header.seq }}</td>\n    <td>{{ block.body.txns.length }}</td>\n    <td>{{ GetBlockAmount(block.body.txns) }}</td>\n    <td>{{ block.header.block_hash }}</td>\n  </tr>\n\n\n  </tbody>\n</table>\n\n<app-skycoin-pagination (onChangePage)=\"handlePageChange($event)\"></app-skycoin-pagination>\n"

/***/ }),

/***/ 842:
/***/ (function(module, exports) {

module.exports = "<div class=\"row\">\n  <h2>Block details</h2>\n  <div class=\"card\" *ngIf=\"block\">\n    <div class=\"card-content\" >\n      <table>\n        <tbody >\n\n        <tr>\n\n          <td>Height</td><td>{{ block.header.seq }}</td>\n\n        </tr>\n        <tr>\n          <td>Timestamp</td><td>{{getTime(block.header.timestamp)}}</td>\n        </tr>\n        <tr>\n          <td>Hash</td><td>{{ block.header.block_hash }}</td>\n        </tr>\n\n        <tr>\n          <td>Parent Hash</td><td>{{ block.header.previous_block_hash }}</td>\n        </tr>\n\n        <tr>\n          <td>Total Amount</td><td>{{ getAmount(block) }} coins</td>\n        </tr>\n\n        </tbody>\n      </table>\n    </div>\n  </div>\n\n\n    <div class=\"col s12\">Transactions</div>\n    <div class=\"col s12\">\n      <div class=\"col s1\">S.no</div>\n      <div class=\"col s2\">Transaction Id</div>\n      <div class=\"col s3\">Inputs</div>\n      <div class=\"col s4\">Outputs</div>\n      <div class=\"col s2\">Amount</div>\n    </div>\n  <div *ngIf=\"block\">\n    <div  class=\"col s12\" *ngFor=\"let transaction of block.body.txns;let i = index\">\n\n      <div class=\"col s1\">{{i+1}}</div>\n      <div class=\"col s2 wrap\">{{transaction.txid}}</div>\n      <div class=\"col s3 wrap\"><p *ngFor=\"let input of transaction.inputs;\" class=\"input\">{{input}}</p></div>\n      <div class=\"col s6 wrap\"><p *ngFor=\"let output of transaction.outputs;\" class=\"input\">{{output.dst}}  - {{output.coins}}</p></div>\n    </div>\n  </div>\n\n\n\n\n</div>\n"

/***/ }),

/***/ 843:
/***/ (function(module, exports) {

module.exports = "<footer class=\"page-footer\">\n  <div class=\"container\">\n    <div class=\"row\">\n      <div class=\"col l6 s12\">\n        <h5 class=\"white-text\">The Skycoin Wallet allows you to hold and secure skycoin. It not only gives you access to the Skycoin blockchain but it’s fast, secure and easy to use!</h5>\n        <p class=\"grey-text text-lighten-4\">Skycoin is a new form of decentralized digital currency that is created and held electronically. It is a necessary element for operating the Skycoin platform.</p>\n      </div>\n      <div class=\"col l4 offset-l2 s12\">\n        <h5 class=\"white-text\">Links</h5>\n        <ul>\n          <li><a class=\"grey-text text-lighten-3\" href=\"http://skycoin.net/infographics.html\">How it works?</a></li>\n          <li><a class=\"grey-text text-lighten-3\" href=\"http://skycoin.net/downloads.html\">Download wallet</a></li>\n          <li><a class=\"grey-text text-lighten-3\" href=\"http://skycoin.net/whitepapers.html\">White papers</a></li>\n          <li><a class=\"grey-text text-lighten-3\" href=\"http://skycoin.net/faq.html\">FAQs</a></li>\n        </ul>\n      </div>\n    </div>\n  </div>\n  <div class=\"footer-copyright\">\n    <div class=\"container\">\n      © 2016 Skycoin\n      <a class=\"grey-text text-lighten-4 right\" href=\"#!\">More Links</a>\n    </div>\n  </div>\n</footer>\n"

/***/ }),

/***/ 844:
/***/ (function(module, exports) {

module.exports = "<nav class=\"light-blue lighten-1\" role=\"navigation\">\n  <div class=\"nav-wrapper container\"><a id=\"logo-container\" href=\"#\" class=\"brand-logo\">Skycoin</a>\n    <ul class=\"right hide-on-med-and-down\">\n      <li><a href=\"http://skycoin.net/infographics.html\">How it works?</a></li>\n      <li><a href=\"http://skycoin.net/downloads.html\">Download wallet</a></li>\n      <li><a href=\"http://skycoin.net/faq.html\">FAQs</a></li>\n      <li><a href=\"http://skycoin.net/whitepapers.html\">White papers</a></li>\n    </ul>\n\n    <ul id=\"nav-mobile\" class=\"side-nav\">\n      <li><a href=\"http://skycoin.net/infographics.html\">How it works?</a></li>\n      <li><a href=\"http://skycoin.net/downloads.html\">Download wallet</a></li>\n      <li><a href=\"http://skycoin.net/faq.html\">FAQs</a></li>\n      <li><a href=\"http://skycoin.net/whitepapers.html\">White papers</a></li>\n    </ul>\n    <a href=\"#\" data-activates=\"nav-mobile\" class=\"button-collapse\"><i class=\"material-icons\">menu</i></a>\n  </div>\n</nav>\n"

/***/ }),

/***/ 845:
/***/ (function(module, exports) {

module.exports = "<ul class=\"pagination\">\n\n  <li [ngClass]=\"{disabled:currentPage <= 5}\"><a href=\"#!\" (click)=\"loadPrevious()\"><i class=\"material-icons\">chevron_left</i></a></li>\n\n  <li *ngFor='let page of currentPages' [ngClass]=\"{active:currentPage == page}\" class=\"active\">\n\n    <a  (click)=\"changePage(page)\" href=\"#!\">{{page}}</a>\n\n  </li>\n  <li [ngClass]=\"{disabled:currentPage === numberOfBlocks/10}\"><a href=\"#!\" (click)=\"loadUpcoming()\"><i class=\"material-icons\">chevron_right</i></a></li>\n</ul>\n"

/***/ }),

/***/ 846:
/***/ (function(module, exports) {

module.exports = "<div class=\"nav-wrapper\">\n  <form (submit)=\"searchBlockHistory(blockSearchKey.value)\">\n    <div class=\"input-field\">\n      <input #blockSearchKey id=\"search\" type=\"search\" required placeholder=\"blockhash\">\n      <label class=\"label-icon\" for=\"search\"><i class=\"material-icons\">search</i></label>\n      <i class=\"material-icons\">close</i>\n    </div>\n  </form>\n</div>\n"

/***/ })

},[1114]);
//# sourceMappingURL=main.bundle.map