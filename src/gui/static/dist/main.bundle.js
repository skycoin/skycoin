webpackJsonp([1],{

/***/ "../../../../../src async recursive":
/***/ (function(module, exports) {

function webpackEmptyContext(req) {
	throw new Error("Cannot find module '" + req + "'.");
}
webpackEmptyContext.keys = function() { return []; };
webpackEmptyContext.resolve = webpackEmptyContext;
module.exports = webpackEmptyContext;
webpackEmptyContext.id = "../../../../../src async recursive";

/***/ }),

/***/ "../../../../../src/app/app.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "md-card {\r\n  max-width: 900px;\r\n  margin-top: 80px;\r\n  margin-right: auto;\r\n  margin-left: auto;\r\n}\r\n\r\n.fill-remaining-space {\r\n  -webkit-box-flex: 1;\r\n      -ms-flex: 1 1 auto;\r\n          flex: 1 1 auto;\r\n}\r\n\r\n.sky-container {\r\n  max-width: 900px;\r\n  margin-top: 20px;\r\n  margin-right: auto;\r\n  margin-left: auto;\r\n  font-family: Roboto, \"Helvetica Neue\", sans-serif;\r\n}\r\n\r\nmd-toolbar span {\r\n  margin: 0 20px;\r\n}\r\n\r\n.search-field {\r\n  border-radius: 8px;\r\n  border: none;\r\n  background-color: #fff;\r\n  padding: 8px;\r\n}\r\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/app.component.html":
/***/ (function(module, exports) {

module.exports = "<md-toolbar color=\"primary\">\n  <span>Skycoin</span>\n  <!--<span><app-breadcrumb></app-breadcrumb></span>-->\n  <!-- This fills the remaining space of the current row -->\n  <span class=\"fill-remaining-space\"></span>\n\n  <span>{{ walletService.sum() | async | sky }}</span>\n  <md-menu #settingsMenu=\"mdMenu\">\n    <button md-menu-item [routerLink]=\"['/settings/network']\"> Networking </button>\n    <button md-menu-item [routerLink]=\"['/settings/blockchain']\"> Blockchain </button>\n    <button md-menu-item [routerLink]=\"['/settings/outputs']\"> Outputs </button>\n    <button md-menu-item [routerLink]=\"['/settings/pending-transactions']\"> Pending Transactions </button>\n    <button md-menu-item [routerLink]=\"['/settings/backup']\"> Back-up wallet </button>\n  </md-menu>\n\n  <button md-button [mdMenuTriggerFor]=\"settingsMenu\">Settings</button>\n</md-toolbar>\n<md-toolbar>\n  <button md-button [routerLink]=\"['/wallets']\" routerLinkActive=\"active\">Wallets</button>\n  <button md-button [routerLink]=\"['/send']\" routerLinkActive=\"active\">Send</button>\n  <button md-button [routerLink]=\"['/history']\" routerLinkActive=\"active\">History</button>\n  <button md-button [routerLink]=\"['/explorer']\" routerLinkActive=\"active\">Explorer</button>\n  <!-- This fills the remaining space of the current row -->\n  <span class=\"fill-remaining-space\"></span>\n\n</md-toolbar>\n<div class=\"sky-container\">\n  <router-outlet></router-outlet>\n</div>\n"

/***/ }),

/***/ "../../../../../src/app/app.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__ = __webpack_require__("../../../../../src/app/services/wallet.service.ts");
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
    function AppComponent(walletService) {
        this.walletService = walletService;
    }
    return AppComponent;
}());
AppComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-root',
        template: __webpack_require__("../../../../../src/app/app.component.html"),
        styles: [__webpack_require__("../../../../../src/app/app.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */]) === "function" && _a || Object])
], AppComponent);

var _a;
//# sourceMappingURL=app.component.js.map

/***/ }),

/***/ "../../../../../src/app/app.module.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__ = __webpack_require__("../../../platform-browser/@angular/platform-browser.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_material__ = __webpack_require__("../../../material/@angular/material.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_platform_browser_animations__ = __webpack_require__("../../../platform-browser/@angular/platform-browser/animations.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__app_component__ = __webpack_require__("../../../../../src/app/app.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5__angular_http__ = __webpack_require__("../../../http/@angular/http.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6__services_api_service__ = __webpack_require__("../../../../../src/app/services/api.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7__services_wallet_service__ = __webpack_require__("../../../../../src/app/services/wallet.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8__components_pages_wallets_wallets_component__ = __webpack_require__("../../../../../src/app/components/pages/wallets/wallets.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_9__components_pages_wallets_address_detail_wallet_detail_component__ = __webpack_require__("../../../../../src/app/components/pages/wallets/address-detail/wallet-detail.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_10__components_pages_wallets_create_wallet_create_wallet_component__ = __webpack_require__("../../../../../src/app/components/pages/wallets/create-wallet/create-wallet.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_11__angular_forms__ = __webpack_require__("../../../forms/@angular/forms.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_12__pipes_sky_pipe__ = __webpack_require__("../../../../../src/app/pipes/sky.pipe.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_13__components_pages_send_skycoin_send_skycoin_component__ = __webpack_require__("../../../../../src/app/components/pages/send-skycoin/send-skycoin.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14__swimlane_ngx_datatable__ = __webpack_require__("../../../../@swimlane/ngx-datatable/release/index.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_14__swimlane_ngx_datatable___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_14__swimlane_ngx_datatable__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_15__components_pages_history_history_component__ = __webpack_require__("../../../../../src/app/components/pages/history/history.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_16__pipes_date_from_now_pipe__ = __webpack_require__("../../../../../src/app/pipes/date-from-now.pipe.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_17__angular_router__ = __webpack_require__("../../../router/@angular/router.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_18__components_layout_breadcrumb_breadcrumb_component__ = __webpack_require__("../../../../../src/app/components/layout/breadcrumb/breadcrumb.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_19__components_pages_transaction_transaction_component__ = __webpack_require__("../../../../../src/app/components/pages/transaction/transaction.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_20__components_layout_back_button_back_button_component__ = __webpack_require__("../../../../../src/app/components/layout/back-button/back-button.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_21__components_pages_explorer_explorer_component__ = __webpack_require__("../../../../../src/app/components/pages/explorer/explorer.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_22__services_blockchain_service__ = __webpack_require__("../../../../../src/app/services/blockchain.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_23__pipes_date_time_pipe__ = __webpack_require__("../../../../../src/app/pipes/date-time.pipe.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_24__pipes_transactions_amount_pipe__ = __webpack_require__("../../../../../src/app/pipes/transactions-amount.pipe.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_25__components_pages_block_block_component__ = __webpack_require__("../../../../../src/app/components/pages/block/block.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_26__components_pages_address_address_component__ = __webpack_require__("../../../../../src/app/components/pages/address/address.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_27__components_pages_settings_pending_transactions_pending_transactions_component__ = __webpack_require__("../../../../../src/app/components/pages/settings/pending-transactions/pending-transactions.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_28__components_pages_settings_outputs_outputs_component__ = __webpack_require__("../../../../../src/app/components/pages/settings/outputs/outputs.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_29__components_pages_settings_blockchain_blockchain_component__ = __webpack_require__("../../../../../src/app/components/pages/settings/blockchain/blockchain.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_30__components_pages_settings_backup_backup_component__ = __webpack_require__("../../../../../src/app/components/pages/settings/backup/backup.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_31__components_pages_settings_network_network_component__ = __webpack_require__("../../../../../src/app/components/pages/settings/network/network.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_32__services_network_service__ = __webpack_require__("../../../../../src/app/services/network.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_33__components_pages_wallets_change_name_change_name_component__ = __webpack_require__("../../../../../src/app/components/pages/wallets/change-name/change-name.component.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_34__components_layout_button_button_component__ = __webpack_require__("../../../../../src/app/components/layout/button/button.component.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AppModule; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};



































var ROUTES = [
    {
        path: '',
        redirectTo: 'wallets',
        pathMatch: 'full'
    },
    {
        path: 'wallets',
        component: __WEBPACK_IMPORTED_MODULE_8__components_pages_wallets_wallets_component__["a" /* WalletsComponent */],
        data: {
            breadcrumb: 'Wallets',
        },
    },
    {
        path: 'send',
        component: __WEBPACK_IMPORTED_MODULE_13__components_pages_send_skycoin_send_skycoin_component__["a" /* SendSkycoinComponent */],
        data: {
            breadcrumb: 'Send Skycoin',
        },
    },
    {
        path: 'history',
        children: [
            {
                path: '',
                component: __WEBPACK_IMPORTED_MODULE_15__components_pages_history_history_component__["a" /* HistoryComponent */],
                data: {
                    breadcrumb: 'History',
                },
            },
            {
                path: ':transaction',
                component: __WEBPACK_IMPORTED_MODULE_19__components_pages_transaction_transaction_component__["a" /* TransactionComponent */],
                data: {
                    breadcrumb: 'Transaction',
                },
            },
        ],
    },
    {
        path: 'explorer',
        children: [
            {
                path: '',
                component: __WEBPACK_IMPORTED_MODULE_21__components_pages_explorer_explorer_component__["a" /* ExplorerComponent */],
                data: {
                    breadcrumb: 'Explorer',
                },
            },
            {
                path: 'address/:address',
                component: __WEBPACK_IMPORTED_MODULE_26__components_pages_address_address_component__["a" /* AddressComponent */],
                data: {
                    breadcrumb: 'Address',
                },
            },
            {
                path: ':block',
                component: __WEBPACK_IMPORTED_MODULE_25__components_pages_block_block_component__["a" /* BlockComponent */],
                data: {
                    breadcrumb: 'Block',
                },
            },
            {
                path: 'transaction/:transaction',
                component: __WEBPACK_IMPORTED_MODULE_19__components_pages_transaction_transaction_component__["a" /* TransactionComponent */],
                data: {
                    breadcrumb: 'Transaction',
                },
            },
        ],
    },
    {
        path: 'settings',
        children: [
            {
                path: 'backup',
                component: __WEBPACK_IMPORTED_MODULE_30__components_pages_settings_backup_backup_component__["a" /* BackupComponent */],
                data: {
                    breadcrumb: 'Backup',
                },
            },
            {
                path: 'blockchain',
                component: __WEBPACK_IMPORTED_MODULE_29__components_pages_settings_blockchain_blockchain_component__["a" /* BlockchainComponent */],
                data: {
                    breadcrumb: 'Blockchain',
                },
            },
            {
                path: 'network',
                component: __WEBPACK_IMPORTED_MODULE_31__components_pages_settings_network_network_component__["a" /* NetworkComponent */],
                data: {
                    breadcrumb: 'Networking',
                },
            },
            {
                path: 'outputs',
                component: __WEBPACK_IMPORTED_MODULE_28__components_pages_settings_outputs_outputs_component__["a" /* OutputsComponent */],
                data: {
                    breadcrumb: 'Outputs',
                },
            },
            {
                path: 'pending-transactions',
                component: __WEBPACK_IMPORTED_MODULE_27__components_pages_settings_pending_transactions_pending_transactions_component__["a" /* PendingTransactionsComponent */],
                data: {
                    breadcrumb: 'Pending transactions',
                },
            },
        ],
    },
];
var AppModule = (function () {
    function AppModule() {
    }
    return AppModule;
}());
AppModule = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_1__angular_core__["NgModule"])({
        declarations: [
            __WEBPACK_IMPORTED_MODULE_4__app_component__["a" /* AppComponent */],
            __WEBPACK_IMPORTED_MODULE_15__components_pages_history_history_component__["a" /* HistoryComponent */],
            __WEBPACK_IMPORTED_MODULE_8__components_pages_wallets_wallets_component__["a" /* WalletsComponent */],
            __WEBPACK_IMPORTED_MODULE_9__components_pages_wallets_address_detail_wallet_detail_component__["a" /* WalletDetailComponent */],
            __WEBPACK_IMPORTED_MODULE_10__components_pages_wallets_create_wallet_create_wallet_component__["a" /* CreateWalletComponent */],
            __WEBPACK_IMPORTED_MODULE_12__pipes_sky_pipe__["a" /* SkyPipe */],
            __WEBPACK_IMPORTED_MODULE_13__components_pages_send_skycoin_send_skycoin_component__["a" /* SendSkycoinComponent */],
            __WEBPACK_IMPORTED_MODULE_16__pipes_date_from_now_pipe__["a" /* DateFromNowPipe */],
            __WEBPACK_IMPORTED_MODULE_18__components_layout_breadcrumb_breadcrumb_component__["a" /* BreadcrumbComponent */],
            __WEBPACK_IMPORTED_MODULE_19__components_pages_transaction_transaction_component__["a" /* TransactionComponent */],
            __WEBPACK_IMPORTED_MODULE_20__components_layout_back_button_back_button_component__["a" /* BackButtonComponent */],
            __WEBPACK_IMPORTED_MODULE_21__components_pages_explorer_explorer_component__["a" /* ExplorerComponent */],
            __WEBPACK_IMPORTED_MODULE_23__pipes_date_time_pipe__["a" /* DateTimePipe */],
            __WEBPACK_IMPORTED_MODULE_24__pipes_transactions_amount_pipe__["a" /* TransactionsAmountPipe */],
            __WEBPACK_IMPORTED_MODULE_25__components_pages_block_block_component__["a" /* BlockComponent */],
            __WEBPACK_IMPORTED_MODULE_26__components_pages_address_address_component__["a" /* AddressComponent */],
            __WEBPACK_IMPORTED_MODULE_27__components_pages_settings_pending_transactions_pending_transactions_component__["a" /* PendingTransactionsComponent */],
            __WEBPACK_IMPORTED_MODULE_28__components_pages_settings_outputs_outputs_component__["a" /* OutputsComponent */],
            __WEBPACK_IMPORTED_MODULE_29__components_pages_settings_blockchain_blockchain_component__["a" /* BlockchainComponent */],
            __WEBPACK_IMPORTED_MODULE_30__components_pages_settings_backup_backup_component__["a" /* BackupComponent */],
            __WEBPACK_IMPORTED_MODULE_31__components_pages_settings_network_network_component__["a" /* NetworkComponent */],
            __WEBPACK_IMPORTED_MODULE_33__components_pages_wallets_change_name_change_name_component__["a" /* ChangeNameComponent */],
            __WEBPACK_IMPORTED_MODULE_34__components_layout_button_button_component__["a" /* ButtonComponent */],
        ],
        entryComponents: [
            __WEBPACK_IMPORTED_MODULE_10__components_pages_wallets_create_wallet_create_wallet_component__["a" /* CreateWalletComponent */],
            __WEBPACK_IMPORTED_MODULE_33__components_pages_wallets_change_name_change_name_component__["a" /* ChangeNameComponent */],
        ],
        imports: [
            __WEBPACK_IMPORTED_MODULE_0__angular_platform_browser__["BrowserModule"],
            __WEBPACK_IMPORTED_MODULE_5__angular_http__["a" /* HttpModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["a" /* MdButtonModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["b" /* MdCardModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["c" /* MdDialogModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["d" /* MdExpansionModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["e" /* MdGridListModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["f" /* MdIconModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["g" /* MdInputModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["h" /* MdListModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["i" /* MdMenuModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["j" /* MdProgressSpinnerModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["k" /* MdSelectModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["l" /* MdTabsModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["m" /* MdToolbarModule */],
            __WEBPACK_IMPORTED_MODULE_2__angular_material__["n" /* MdTooltipModule */],
            __WEBPACK_IMPORTED_MODULE_14__swimlane_ngx_datatable__["NgxDatatableModule"],
            __WEBPACK_IMPORTED_MODULE_3__angular_platform_browser_animations__["a" /* NoopAnimationsModule */],
            __WEBPACK_IMPORTED_MODULE_11__angular_forms__["a" /* ReactiveFormsModule */],
            __WEBPACK_IMPORTED_MODULE_17__angular_router__["a" /* RouterModule */].forRoot(ROUTES),
        ],
        providers: [
            __WEBPACK_IMPORTED_MODULE_6__services_api_service__["a" /* ApiService */],
            __WEBPACK_IMPORTED_MODULE_22__services_blockchain_service__["a" /* BlockchainService */],
            __WEBPACK_IMPORTED_MODULE_32__services_network_service__["a" /* NetworkService */],
            __WEBPACK_IMPORTED_MODULE_7__services_wallet_service__["a" /* WalletService */],
        ],
        bootstrap: [__WEBPACK_IMPORTED_MODULE_4__app_component__["a" /* AppComponent */]]
    })
], AppModule);

//# sourceMappingURL=app.module.js.map

/***/ }),

/***/ "../../../../../src/app/components/layout/back-button/back-button.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "a {\r\n  float: right;\r\n}\r\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/layout/back-button/back-button.component.html":
/***/ (function(module, exports) {

module.exports = "<a md-raised-button color=\"primary\" (click)=\"onClick()\">Back</a>\n"

/***/ }),

/***/ "../../../../../src/app/components/layout/back-button/back-button.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_common__ = __webpack_require__("../../../common/@angular/common.es5.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return BackButtonComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var BackButtonComponent = (function () {
    function BackButtonComponent(location) {
        this.location = location;
    }
    BackButtonComponent.prototype.onClick = function () {
        this.location.back();
    };
    return BackButtonComponent;
}());
BackButtonComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-back-button',
        template: __webpack_require__("../../../../../src/app/components/layout/back-button/back-button.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/layout/back-button/back-button.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_common__["Location"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_common__["Location"]) === "function" && _a || Object])
], BackButtonComponent);

var _a;
//# sourceMappingURL=back-button.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/layout/breadcrumb/breadcrumb.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/layout/breadcrumb/breadcrumb.component.html":
/***/ (function(module, exports) {

module.exports = "<button md-button routerLink=\"\" class=\"breadcrumb\">Skycoin</button>\n\n<button md-button *ngFor=\"let breadcrumb of breadcrumbs\" [routerLink]=\"[breadcrumb.url, breadcrumb.params]\">{{ breadcrumb.label }}</button>\n\n"

/***/ }),

/***/ "../../../../../src/app/components/layout/breadcrumb/breadcrumb.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_router__ = __webpack_require__("../../../router/@angular/router.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_add_operator_filter__ = __webpack_require__("../../../../rxjs/add/operator/filter.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_add_operator_filter___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs_add_operator_filter__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return BreadcrumbComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var BreadcrumbComponent = (function () {
    function BreadcrumbComponent(activatedRoute, router) {
        this.activatedRoute = activatedRoute;
        this.router = router;
        this.breadcrumbs = [];
    }
    BreadcrumbComponent.prototype.ngOnInit = function () {
        var _this = this;
        var ROUTE_DATA_BREADCRUMB = "breadcrumb";
        //subscribe to the NavigationEnd event
        this.router.events.filter(function (event) { return event instanceof __WEBPACK_IMPORTED_MODULE_1__angular_router__["b" /* NavigationEnd */]; }).subscribe(function (event) {
            //set breadcrumbs
            var root = _this.activatedRoute.root;
            _this.breadcrumbs = _this.getBreadcrumbs(root);
        });
    };
    BreadcrumbComponent.prototype.getBreadcrumbs = function (route, url, breadcrumbs) {
        if (url === void 0) { url = ""; }
        if (breadcrumbs === void 0) { breadcrumbs = []; }
        var ROUTE_DATA_BREADCRUMB = "breadcrumb";
        //get the child routes
        var children = route.children;
        //return if there are no more children
        if (children.length === 0) {
            return breadcrumbs;
        }
        //iterate over each children
        for (var _i = 0, children_1 = children; _i < children_1.length; _i++) {
            var child = children_1[_i];
            //verify primary route
            if (child.outlet !== __WEBPACK_IMPORTED_MODULE_1__angular_router__["c" /* PRIMARY_OUTLET */]) {
                continue;
            }
            //verify the custom data property "breadcrumb" is specified on the route
            if (!child.snapshot.data.hasOwnProperty(ROUTE_DATA_BREADCRUMB)) {
                return this.getBreadcrumbs(child, url, breadcrumbs);
            }
            //get the route's URL segment
            var routeURL = child.snapshot.url.map(function (segment) { return segment.path; }).join("/");
            //append route URL to URL
            url += "/" + routeURL;
            //add breadcrumb
            var breadcrumb = {
                label: child.snapshot.data[ROUTE_DATA_BREADCRUMB],
                params: child.snapshot.params,
                url: url
            };
            breadcrumbs.push(breadcrumb);
            //recursive
            return this.getBreadcrumbs(child, url, breadcrumbs);
        }
    };
    return BreadcrumbComponent;
}());
BreadcrumbComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-breadcrumb',
        template: __webpack_require__("../../../../../src/app/components/layout/breadcrumb/breadcrumb.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/layout/breadcrumb/breadcrumb.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_router__["d" /* ActivatedRoute */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_router__["d" /* ActivatedRoute */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__angular_router__["e" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_router__["e" /* Router */]) === "function" && _b || Object])
], BreadcrumbComponent);

var _a, _b;
//# sourceMappingURL=breadcrumb.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/layout/button/button.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "md-icon {\r\n  margin-left: 10px;\r\n  opacity: 0.3;\r\n}\r\n\r\nmd-spinner {\r\n  display: inline-block;\r\n  height: 24px !important;\r\n  width: 24px !important;\r\n  margin-left: 10px;\r\n}\r\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/layout/button/button.component.html":
/***/ (function(module, exports) {

module.exports = "<button type=\"submit\" md-raised-button color=\"primary\" [disabled]=\"disabled()\" [mdTooltip]=\"error ? error : null\">\n  {{ placeholder }}\n  <md-icon *ngIf=\"state === 1\">done</md-icon>\n  <md-icon *ngIf=\"state === 2\">error</md-icon>\n  <md-spinner *ngIf=\"state === 0\" class=\"in-button\"></md-spinner>\n</button>\n"

/***/ }),

/***/ "../../../../../src/app/components/layout/button/button.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return ButtonComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};

var ButtonComponent = (function () {
    function ButtonComponent() {
    }
    ButtonComponent.prototype.setLoading = function () {
        this.state = 0;
    };
    ButtonComponent.prototype.setSuccess = function () {
        var _this = this;
        this.state = 1;
        setTimeout(function () { return _this.state = null; }, 3000);
    };
    ButtonComponent.prototype.setError = function (error) {
        this.error = error['_body'];
        this.state = 2;
    };
    ButtonComponent.prototype.disabled = function () {
        return this.state === 0 || (!(this.form === undefined) && !(this.form && this.form.valid));
    };
    return ButtonComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"])(),
    __metadata("design:type", Object)
], ButtonComponent.prototype, "form", void 0);
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"])(),
    __metadata("design:type", String)
], ButtonComponent.prototype, "placeholder", void 0);
ButtonComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-button',
        template: __webpack_require__("../../../../../src/app/components/layout/button/button.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/layout/button/button.component.css")]
    })
], ButtonComponent);

//# sourceMappingURL=button.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/address/address.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, ".skycoin-details {\r\n  margin-top: 40px;\r\n  background-color: #eee;\r\n  margin-bottom: 20px;\r\n}\r\n\r\n.skycoin-detail-keys {\r\n  display: inline-block;\r\n}\r\n\r\n.skycoin-detail-values {\r\n  display: inline-block;\r\n}\r\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/address/address.component.html":
/***/ (function(module, exports) {

module.exports = "<md-card>\n  <md-card-title>Address <app-back-button></app-back-button></md-card-title>\n  <md-card class=\"skycoin-details\">\n    <div class=\"skycoin-detail-keys\">\n      ID: <br>\n      Transactions: <br>\n      Balance: <br>\n    </div>\n    <div class=\"skycoin-detail-values\" *ngIf=\"transactions\">\n      {{ id }} <br>\n      {{ transactions.length }} <br>\n      {{ transactions | transactionsAmount }}\n    </div>\n  </md-card>\n\n  <h3>Transactions</h3>\n  <md-expansion-panel *ngFor=\"let transaction of transactions\">\n    <md-expansion-panel-header>\n      <md-panel-title>\n        {{ transaction.txid }}\n      </md-panel-title>\n      <md-panel-description>\n        <!--{{ block.header.timestamp | dateTime }}-->\n      </md-panel-description>\n    </md-expansion-panel-header>\n    <md-list *ngIf=\"transaction\">\n      <h3 md-subheader>Inputs</h3>\n      <md-list-item *ngFor=\"let input of transaction.inputs\">\n        <h4 md-line>{{ input.owner }}</h4>\n      </md-list-item>\n      <md-divider></md-divider>\n      <h3 md-subheader>Outputs</h3>\n      <md-list-item *ngFor=\"let output of transaction.outputs\">\n        <h4 md-line>{{ output.dst }} ({{ output.coins }} SKY)</h4>\n      </md-list-item>\n    </md-list>\n    <div class=\"button-line\">\n      <a md-raised-button color=\"primary\" [routerLink]=\"['/explorer/transaction/', transaction.txid]\">Details</a>\n    </div>\n  </md-expansion-panel>\n</md-card>\n\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/address/address.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_blockchain_service__ = __webpack_require__("../../../../../src/app/services/blockchain.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("../../../router/@angular/router.es5.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return AddressComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var AddressComponent = (function () {
    function AddressComponent(blockchainService, route) {
        this.blockchainService = blockchainService;
        this.route = route;
    }
    AddressComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.route.params.switchMap(function (params) {
            _this.id = params.address;
            return _this.blockchainService.addressTransactions(params.address);
        }).subscribe(function (response) {
            _this.transactions = response;
            console.log(response);
        });
        this.route.params.switchMap(function (params) { return _this.blockchainService.addressBalance(params.address); }).subscribe(function (response) {
            _this.balance = response;
            console.log(response);
        });
    };
    return AddressComponent;
}());
AddressComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-address',
        template: __webpack_require__("../../../../../src/app/components/pages/address/address.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/address/address.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_blockchain_service__["a" /* BlockchainService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_blockchain_service__["a" /* BlockchainService */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["d" /* ActivatedRoute */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["d" /* ActivatedRoute */]) === "function" && _b || Object])
], AddressComponent);

var _a, _b;
//# sourceMappingURL=address.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/block/block.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, ".skycoin-details {\r\n  margin-top: 40px;\r\n  background-color: #eee;\r\n  margin-bottom: 20px;\r\n}\r\n\r\n.skycoin-detail-keys {\r\n  display: inline-block;\r\n}\r\n\r\n.skycoin-detail-values {\r\n  display: inline-block;\r\n}\r\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/block/block.component.html":
/***/ (function(module, exports) {

module.exports = "<md-card>\n  <md-card-title>Block {{ block && block.header ? block.header.seq : '..' }} <app-back-button></app-back-button></md-card-title>\n  <md-card class=\"skycoin-details\">\n    <div class=\"skycoin-detail-keys\">\n      ID: <br>\n      Timestamp: <br>\n      Hash: <br>\n      Parent hash: <br>\n      Transactions: <br>\n      Total value: <br>\n    </div>\n    <div class=\"skycoin-detail-values\" *ngIf=\"block\">\n      {{ block && block.header ? block.header.seq : '' }} <br>\n      {{ block.header.timestamp | dateTime }} <br>\n      {{ block && block.header ? block.header.block_hash : '' }} <br>\n      {{ block && block.header ? block.header.previous_block_hash : '' }} <br>\n      {{ block && block.body && block.body.txns ? block.body.txns.length : '' }} <br>\n      {{ block.body.txns | transactionsAmount }}\n    </div>\n  </md-card>\n\n  <h3>Transactions</h3>\n  <div *ngIf=\"block && block.body && block.body.txns\">\n    <md-expansion-panel *ngFor=\"let transaction of block.body.txns\">\n      <md-expansion-panel-header>\n        <md-panel-title>\n          {{ transaction.txid }}\n        </md-panel-title>\n        <md-panel-description>\n          {{ block.header.timestamp | dateTime }}\n        </md-panel-description>\n      </md-expansion-panel-header>\n      <md-list *ngIf=\"transaction\">\n        <h3 md-subheader>Inputs</h3>\n        <md-list-item *ngFor=\"let input of transaction.inputs\">\n          <h4 md-line>{{ input }}</h4>\n        </md-list-item>\n        <md-divider></md-divider>\n        <h3 md-subheader>Outputs</h3>\n        <md-list-item *ngFor=\"let output of transaction.outputs\">\n          <h4 md-line>{{ output.dst }} ({{ output.coins }} SKY)</h4>\n        </md-list-item>\n      </md-list>\n      <div class=\"button-line\">\n        <a md-raised-button color=\"primary\" [routerLink]=\"['/explorer/transaction/', transaction.txid]\">Details</a>\n      </div>\n    </md-expansion-panel>\n  </div>\n</md-card>\n\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/block/block.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_blockchain_service__ = __webpack_require__("../../../../../src/app/services/blockchain.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("../../../router/@angular/router.es5.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return BlockComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var BlockComponent = (function () {
    function BlockComponent(blockchainService, route) {
        this.blockchainService = blockchainService;
        this.route = route;
    }
    BlockComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.route.params.switchMap(function (params) { return _this.blockchainService.block(params.block); }).subscribe(function (response) {
            _this.block = response;
        });
    };
    return BlockComponent;
}());
BlockComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-block',
        template: __webpack_require__("../../../../../src/app/components/pages/block/block.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/block/block.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_blockchain_service__["a" /* BlockchainService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_blockchain_service__["a" /* BlockchainService */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["d" /* ActivatedRoute */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["d" /* ActivatedRoute */]) === "function" && _b || Object])
], BlockComponent);

var _a, _b;
//# sourceMappingURL=block.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/explorer/explorer.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/explorer/explorer.component.html":
/***/ (function(module, exports) {

module.exports = "<ngx-datatable #table\n  class=\"material\"\n  [rows]=\"blocks\"\n  columnMode=\"flex\"\n  [headerHeight]=\"50\"\n  [footerHeight]=\"50\"\n  [rowHeight]=\"50\"\n  [limit]=\"10\"\n  [scrollbarH]=\"true\"\n  (activate)=\"onActivate($event)\">\n  <ngx-datatable-column name=\"Timestamp\" prop=\"header.timestamp\" [flexGrow]=\"1\">\n    <ng-template let-value=\"value\" ngx-datatable-cell-template>\n      {{ value | dateTime }}\n    </ng-template>\n  </ngx-datatable-column>\n  <ngx-datatable-column name=\"Block height\" prop=\"header.seq\" [flexGrow]=\"1\"></ngx-datatable-column>\n  <ngx-datatable-column name=\"Transactions\" prop=\"body.txns\" [flexGrow]=\"1\">\n    <ng-template let-value=\"value\" ngx-datatable-cell-template>\n      {{ value ? value.length : 0 }}\n    </ng-template>\n  </ngx-datatable-column>\n  <ngx-datatable-column name=\"Amount\" prop=\"body.txns\" [flexGrow]=\"1\">\n    <ng-template let-value=\"value\" ngx-datatable-cell-template>\n      {{ value | transactionsAmount }}\n    </ng-template>\n  </ngx-datatable-column>\n</ngx-datatable>\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/explorer/explorer.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_blockchain_service__ = __webpack_require__("../../../../../src/app/services/blockchain.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("../../../router/@angular/router.es5.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return ExplorerComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var ExplorerComponent = (function () {
    function ExplorerComponent(blockchainService, router) {
        this.blockchainService = blockchainService;
        this.router = router;
    }
    ExplorerComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.blockchainService.blocks().subscribe(function (data) { return _this.blocks = data; });
    };
    ExplorerComponent.prototype.onActivate = function (response) {
        if (response.row && response.row.header) {
            this.router.navigate(['/explorer', response.row.header.seq]);
        }
    };
    return ExplorerComponent;
}());
ExplorerComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-explorer',
        template: __webpack_require__("../../../../../src/app/components/pages/explorer/explorer.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/explorer/explorer.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_blockchain_service__["a" /* BlockchainService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_blockchain_service__["a" /* BlockchainService */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["e" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["e" /* Router */]) === "function" && _b || Object])
], ExplorerComponent);

var _a, _b;
//# sourceMappingURL=explorer.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/history/history.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "ngx-datatable {\r\n  margin-top: 20px;\r\n}\r\n\r\nmd-icon {\r\n  cursor: pointer;\r\n}\r\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/history/history.component.html":
/***/ (function(module, exports) {

module.exports = "<ngx-datatable #table\r\n  class=\"material\"\r\n  [rows]=\"walletService.history() | async\"\r\n  columnMode=\"flex\"\r\n  [headerHeight]=\"50\"\r\n  [footerHeight]=\"50\"\r\n  [rowHeight]=\"50\"\r\n  [limit]=\"10\"\r\n  [scrollbarH]=\"true\"\r\n  (activate)=\"onActivate($event)\">\r\n  <ngx-datatable-column name=\"Timestamp\" prop=\"timestamp\" [flexGrow]=\"2\">\r\n    <ng-template let-value=\"value\" ngx-datatable-cell-template>\r\n      {{ value | dateFromNow }}\r\n    </ng-template>\r\n  </ngx-datatable-column>\r\n  <ngx-datatable-column name=\"Amount\" [flexGrow]=\"1\"></ngx-datatable-column>\r\n  <ngx-datatable-column name=\"Transaction ID\" prop=\"txid\" [flexGrow]=\"5\"></ngx-datatable-column>\r\n</ngx-datatable>\r\n\r\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/history/history.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__ = __webpack_require__("../../../../../src/app/services/wallet.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("../../../router/@angular/router.es5.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return HistoryComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var HistoryComponent = (function () {
    function HistoryComponent(router, walletService) {
        this.router = router;
        this.walletService = walletService;
    }
    HistoryComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.walletService.history().subscribe(function (transactions) { return _this.transactions = _this.mapTransactions(transactions); });
    };
    HistoryComponent.prototype.onActivate = function (response) {
        if (response.row && response.row.txid) {
            this.router.navigate(['/history', response.row.txid]);
        }
    };
    HistoryComponent.prototype.mapTransactions = function (transactions) {
        return transactions.map(function (transaction) {
            transaction.amount = transaction.outputs.map(function (output) { return output.coins >= 0 ? output.coins : 0; })
                .reduce(function (a, b) { return a + parseInt(b); }, 0);
            return transaction;
        });
    };
    return HistoryComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["ViewChild"])('table'),
    __metadata("design:type", Object)
], HistoryComponent.prototype, "table", void 0);
HistoryComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-history',
        template: __webpack_require__("../../../../../src/app/components/pages/history/history.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/history/history.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["e" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["e" /* Router */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */]) === "function" && _b || Object])
], HistoryComponent);

var _a, _b;
//# sourceMappingURL=history.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/send-skycoin/send-skycoin.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, ".input-field {\r\n  display: block;\r\n  margin-top: 20px;\r\n  width: 100%;\r\n}\r\n\r\n.send-skycoin-form {\r\n  padding-top: 40px;\r\n}\r\n\r\nmd-select {\r\n  padding-bottom: 1.29688em;\r\n}\r\n\r\nmd-card {\r\n  margin-bottom: 20px;\r\n}\r\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/send-skycoin/send-skycoin.component.html":
/***/ (function(module, exports) {

module.exports = "<md-card [formGroup]=\"form\" class=\"send-skycoin-form\">\r\n  <md-select formControlName=\"wallet_id\" placeholder=\"Wallet\" class=\"input-field\">\r\n    <md-option *ngFor=\"let wallet of walletService.all() | async\" [value]=\"wallet.meta.filename\">\r\n      {{ wallet.meta.label }} ({{ wallet.balance | sky }})\r\n    </md-option>\r\n  </md-select>\r\n  <md-input-container class=\"input-field\">\r\n    <input mdInput formControlName=\"address\" placeholder=\"Recipient address\">\r\n  </md-input-container>\r\n  <md-input-container class=\"input-field\">\r\n    <input mdInput formControlName=\"amount\" placeholder=\"Amount\">\r\n  </md-input-container>\r\n  <div class=\"button-line\">\r\n    <app-button #button placeholder=\"Send\" (click)=\"send()\"></app-button>\r\n  </div>\r\n</md-card>\r\n<!--Time, Status, Address, Amount, Transaction ID-->\r\n<md-card>\r\n  <h3>Recent transactions</h3>\r\n  <ngx-datatable #table\r\n    class=\"material\"\r\n    [rows]=\"records\"\r\n    columnMode=\"flex\"\r\n    [headerHeight]=\"50\"\r\n    [footerHeight]=\"50\"\r\n    [rowHeight]=\"50\"\r\n    [limit]=\"10\"\r\n    [scrollbarH]=\"true\"\r\n    (activate)=\"onActivate($event)\">\r\n    <ngx-datatable-column name=\"Timestamp\" prop=\"txn.timestamp\" [flexGrow]=\"2\">\r\n      <ng-template let-value=\"value\" ngx-datatable-cell-template>\r\n        <strong>{{ value | dateTime }}</strong>\r\n      </ng-template>\r\n    </ngx-datatable-column>\r\n    <ngx-datatable-column name=\"Status\" [flexGrow]=\"2\">\r\n      <ng-template let-value=\"value\" ngx-datatable-cell-template>\r\n        <strong>{{ value.confirmed ? 'Confirmed' : 'Unconfirmed' }}</strong>\r\n      </ng-template>\r\n    </ngx-datatable-column>\r\n    <ngx-datatable-column name=\"Address\" [flexGrow]=\"4\"></ngx-datatable-column>\r\n    <ngx-datatable-column name=\"Amount\" [flexGrow]=\"1\"></ngx-datatable-column>\r\n  </ngx-datatable>\r\n</md-card>\r\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/send-skycoin/send-skycoin.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__ = __webpack_require__("../../../../../src/app/services/wallet.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_forms__ = __webpack_require__("../../../forms/@angular/forms.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_observable_IntervalObservable__ = __webpack_require__("../../../../rxjs/observable/IntervalObservable.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_observable_IntervalObservable___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_rxjs_observable_IntervalObservable__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__angular_router__ = __webpack_require__("../../../router/@angular/router.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_rxjs_add_operator_delay__ = __webpack_require__("../../../../rxjs/add/operator/delay.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_rxjs_add_operator_delay___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_5_rxjs_add_operator_delay__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return SendSkycoinComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};






var SendSkycoinComponent = (function () {
    function SendSkycoinComponent(formBuilder, walletService, router) {
        this.formBuilder = formBuilder;
        this.walletService = walletService;
        this.router = router;
        this.records = [];
        this.transactions = [];
    }
    SendSkycoinComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.initForm();
        __WEBPACK_IMPORTED_MODULE_3_rxjs_observable_IntervalObservable__["IntervalObservable"]
            .create(2500)
            .filter(function () { return !!_this.transactions.length; })
            .flatMap(function () { return _this.walletService.retrieveUpdatedTransactions(_this.transactions); })
            .subscribe(function (transactions) { return _this.records = transactions; });
        this.walletService.recent().subscribe(function (transactions) { return _this.transactions = transactions; });
    };
    SendSkycoinComponent.prototype.onActivate = function (response) {
        if (response.row && response.row.txid) {
            this.router.navigate(['/history', response.row.txid]);
        }
    };
    SendSkycoinComponent.prototype.send = function () {
        var _this = this;
        this.button.setLoading();
        this.walletService.sendSkycoin(this.form.value.wallet_id, this.form.value.address, this.form.value.amount * 1000000)
            .delay(1000)
            .subscribe(function (response) {
            _this.resetForm();
            _this.button.setSuccess();
        }, function (error) { return _this.button.setError(error); });
    };
    SendSkycoinComponent.prototype.initForm = function () {
        this.form = this.formBuilder.group({
            wallet_id: ['', __WEBPACK_IMPORTED_MODULE_2__angular_forms__["h" /* Validators */].required],
            address: ['', __WEBPACK_IMPORTED_MODULE_2__angular_forms__["h" /* Validators */].required],
            amount: ['', __WEBPACK_IMPORTED_MODULE_2__angular_forms__["h" /* Validators */].required],
        });
    };
    SendSkycoinComponent.prototype.resetForm = function () {
        this.form.controls.wallet_id.reset(undefined);
        this.form.controls.address.reset(undefined);
        this.form.controls.amount.reset(undefined);
    };
    return SendSkycoinComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["ViewChild"])('button'),
    __metadata("design:type", Object)
], SendSkycoinComponent.prototype, "button", void 0);
SendSkycoinComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-send-skycoin',
        template: __webpack_require__("../../../../../src/app/components/pages/send-skycoin/send-skycoin.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/send-skycoin/send-skycoin.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_2__angular_forms__["i" /* FormBuilder */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_forms__["i" /* FormBuilder */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_4__angular_router__["e" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_4__angular_router__["e" /* Router */]) === "function" && _c || Object])
], SendSkycoinComponent);

var _a, _b, _c;
//# sourceMappingURL=send-skycoin.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/backup/backup.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "md-card {\r\n  margin-bottom: 20px;\r\n}\r\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/backup/backup.component.html":
/***/ (function(module, exports) {

module.exports = "<md-card>\n  <p>Wallet Directory: {{ folder }}</p>\n\n  <p>BACKUP YOUR SEED. ON PAPER. IN A SAFE PLACE. As long as you have your seed, you can recover your coins.</p>\n</md-card>\n\n<ngx-datatable #table\n  class=\"material\"\n  [rows]=\"walletService.all() | async\"\n  columnMode=\"flex\"\n  [headerHeight]=\"50\"\n  [footerHeight]=\"50\"\n  [rowHeight]=\"50\"\n  [limit]=\"10\"\n  [scrollbarH]=\"true\">\n  <ngx-datatable-column name=\"Wallet label\" prop=\"meta.label\" [flexGrow]=\"1\"></ngx-datatable-column>\n  <ngx-datatable-column name=\"File name\" prop=\"meta.filename\" [flexGrow]=\"1\"></ngx-datatable-column>\n  <ngx-datatable-column name=\"Download\" [flexGrow]=\"1\">\n    <ng-template let-row=\"row\" ngx-datatable-cell-template>\n      <button md-raised-button color=\"primary\" (click)=\"download(row)\">Download</button>\n    </ng-template>\n  </ngx-datatable-column>\n</ngx-datatable>\n\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/backup/backup.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__ = __webpack_require__("../../../../../src/app/services/wallet.service.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return BackupComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var BackupComponent = (function () {
    function BackupComponent(walletService) {
        this.walletService = walletService;
    }
    BackupComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.walletService.folder().subscribe(function (folder) { return _this.folder = folder; });
    };
    BackupComponent.prototype.ngOnDestroy = function () {
        this.walletService.all().subscribe(function (wallets) { return wallets.forEach(function (wallet) { return wallet.visible = false; }); });
    };
    BackupComponent.prototype.download = function (wallet) {
        var blob = new Blob([JSON.stringify({ seed: wallet.meta.seed })], { type: 'application/json' });
        var link = document.createElement('a');
        link.href = window.URL.createObjectURL(blob);
        link['download'] = wallet.meta.filename + '.json';
        link.click();
    };
    return BackupComponent;
}());
BackupComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-backup',
        template: __webpack_require__("../../../../../src/app/components/pages/settings/backup/backup.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/settings/backup/backup.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */]) === "function" && _a || Object])
], BackupComponent);

var _a;
//# sourceMappingURL=backup.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/blockchain/blockchain.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, ".skycoin-details {\r\n  margin-top: 40px;\r\n  background-color: #eee;\r\n  margin-bottom: 20px;\r\n}\r\n\r\n.skycoin-detail-keys {\r\n  display: inline-block;\r\n}\r\n\r\n.skycoin-detail-values {\r\n  display: inline-block;\r\n}\r\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/blockchain/blockchain.component.html":
/***/ (function(module, exports) {

module.exports = "<md-card>\n  <md-card-title>Blockchain</md-card-title>\n  <md-card class=\"skycoin-details\">\n    <div class=\"skycoin-detail-keys\">\n      Number of blocks: <br>\n      Time since last block: <br>\n      Hash of last block: <br>\n    </div>\n    <div class=\"skycoin-detail-values\" *ngIf=\"block && block.header\">\n      {{ block.header.seq }} <br>\n      {{ block.header.timestamp | dateFromNow }} <br>\n      {{ block.header.block_hash }}\n    </div>\n  </md-card>\n</md-card>\n\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/blockchain/blockchain.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_blockchain_service__ = __webpack_require__("../../../../../src/app/services/blockchain.service.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return BlockchainComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var BlockchainComponent = (function () {
    function BlockchainComponent(blockchainService) {
        this.blockchainService = blockchainService;
    }
    BlockchainComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.blockchainService.lastBlock().subscribe(function (block) { return _this.block = block; });
    };
    return BlockchainComponent;
}());
BlockchainComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-blockchain',
        template: __webpack_require__("../../../../../src/app/components/pages/settings/blockchain/blockchain.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/settings/blockchain/blockchain.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_blockchain_service__["a" /* BlockchainService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_blockchain_service__["a" /* BlockchainService */]) === "function" && _a || Object])
], BlockchainComponent);

var _a;
//# sourceMappingURL=blockchain.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/network/network.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "ngx-datatable {\r\n  margin-bottom: 20px;\r\n}\r\n\r\nmd-card {\r\n  margin-bottom: 20px;\r\n}\r\n\r\nmd-card h3 button {\r\n  display: inline;\r\n  float: right;\r\n  margin-top: -8px;\r\n}\r\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/network/network.component.html":
/***/ (function(module, exports) {

module.exports = "<md-card>\n  <h3>Automatic peers</h3>\n\n  <ngx-datatable\n    class=\"material\"\n    [rows]=\"networkService.automatic() | async\"\n    columnMode=\"flex\"\n    [headerHeight]=\"40\"\n    [footerHeight]=\"40\"\n    [rowHeight]=\"40\"\n    [limit]=\"5\">\n    <ngx-datatable-column name=\"ID\" prop=\"id\" [flexGrow]=\"1\"></ngx-datatable-column>\n    <ngx-datatable-column name=\"Address\" prop=\"address\" [flexGrow]=\"1\"></ngx-datatable-column>\n    <ngx-datatable-column name=\"Port\" prop=\"listen_port\" [flexGrow]=\"1\"></ngx-datatable-column>\n  </ngx-datatable>\n</md-card>\n\n<md-card>\n  <h3>Default peers</h3>\n\n  <ngx-datatable\n    class=\"material\"\n    [rows]=\"defaultConnections\"\n    columnMode=\"flex\"\n    [headerHeight]=\"40\"\n    [rowHeight]=\"50\">\n    <ngx-datatable-column name=\"ID\" prop=\"id\" [flexGrow]=\"1\"></ngx-datatable-column>\n    <ngx-datatable-column name=\"Address\" prop=\"address\" [flexGrow]=\"1\"></ngx-datatable-column>\n    <ngx-datatable-column name=\"Port\" prop=\"listen_port\" [flexGrow]=\"1\"></ngx-datatable-column>\n  </ngx-datatable>\n</md-card>\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/network/network.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_network_service__ = __webpack_require__("../../../../../src/app/services/network.service.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return NetworkComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var NetworkComponent = (function () {
    function NetworkComponent(networkService) {
        this.networkService = networkService;
    }
    NetworkComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.networkService.automatic().subscribe(function (output) { return console.log(output); });
        this.networkService.retrieveDefaultConnections().first().subscribe(function (output) { return _this.defaultConnections = output; });
    };
    return NetworkComponent;
}());
NetworkComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-network',
        template: __webpack_require__("../../../../../src/app/components/pages/settings/network/network.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/settings/network/network.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_network_service__["a" /* NetworkService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_network_service__["a" /* NetworkService */]) === "function" && _a || Object])
], NetworkComponent);

var _a;
//# sourceMappingURL=network.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/outputs/outputs.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/outputs/outputs.component.html":
/***/ (function(module, exports) {

module.exports = "<ngx-datatable #table\n               class=\"material\"\n               [rows]=\"outputs\"\n               columnMode=\"flex\"\n               [headerHeight]=\"50\"\n               [footerHeight]=\"50\"\n               [rowHeight]=\"50\"\n               [limit]=\"10\"\n               [scrollbarH]=\"true\">\n  <ngx-datatable-column name=\"Address\" prop=\"address\" [flexGrow]=\"2\"></ngx-datatable-column>\n  <ngx-datatable-column name=\"Coins\" prop=\"coins\" [flexGrow]=\"1\"></ngx-datatable-column>\n  <ngx-datatable-column name=\"Hours\" prop=\"hours\" [flexGrow]=\"1\"></ngx-datatable-column>\n</ngx-datatable>\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/outputs/outputs.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__ = __webpack_require__("../../../../../src/app/services/wallet.service.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return OutputsComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var OutputsComponent = (function () {
    function OutputsComponent(walletService) {
        this.walletService = walletService;
    }
    OutputsComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.walletService.outputs().subscribe(function (outputs) { return _this.outputs = outputs.head_outputs; });
    };
    return OutputsComponent;
}());
OutputsComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-outputs',
        template: __webpack_require__("../../../../../src/app/components/pages/settings/outputs/outputs.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/settings/outputs/outputs.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */]) === "function" && _a || Object])
], OutputsComponent);

var _a;
//# sourceMappingURL=outputs.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/pending-transactions/pending-transactions.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/pending-transactions/pending-transactions.component.html":
/***/ (function(module, exports) {

module.exports = "<ngx-datatable #table\n               class=\"material\"\n               [rows]=\"transactions\"\n               columnMode=\"flex\"\n               [headerHeight]=\"50\"\n               [footerHeight]=\"50\"\n               [rowHeight]=\"50\"\n               [limit]=\"10\"\n               [scrollbarH]=\"true\"\n               (activate)=\"onActivate($event)\">\n  <ngx-datatable-column name=\"Timestamp\" prop=\"timestamp\" [flexGrow]=\"2\">\n    <ng-template let-value=\"value\" ngx-datatable-cell-template>\n      <strong>{{ value | dateTime }}</strong>\n    </ng-template>\n  </ngx-datatable-column>\n  <ngx-datatable-column name=\"Amount\" [flexGrow]=\"1\"></ngx-datatable-column>\n  <ngx-datatable-column name=\"Transaction ID\" prop=\"txid\" [flexGrow]=\"8\"></ngx-datatable-column>\n</ngx-datatable>\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/settings/pending-transactions/pending-transactions.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__ = __webpack_require__("../../../../../src/app/services/wallet.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_moment__ = __webpack_require__("../../../../moment/moment.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_moment___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_moment__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_router__ = __webpack_require__("../../../router/@angular/router.es5.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return PendingTransactionsComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};




var PendingTransactionsComponent = (function () {
    function PendingTransactionsComponent(walletService, router) {
        this.walletService = walletService;
        this.router = router;
    }
    PendingTransactionsComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.walletService.pendingTransactions().subscribe(function (transactions) {
            _this.transactions = _this.mapTransactions(transactions);
        });
    };
    PendingTransactionsComponent.prototype.onActivate = function (response) {
        if (response.row && response.row.txid) {
            this.router.navigate(['/history', response.row.txid]);
        }
    };
    PendingTransactionsComponent.prototype.mapTransactions = function (transactions) {
        return transactions.map(function (transaction) {
            transaction.transaction.timestamp = __WEBPACK_IMPORTED_MODULE_2_moment__(transaction.received).unix();
            return transaction.transaction;
        })
            .map(function (transaction) {
            transaction.amount = transaction.outputs.map(function (output) { return output.coins >= 0 ? output.coins : 0; })
                .reduce(function (a, b) { return a + parseInt(b); }, 0);
            return transaction;
        });
    };
    return PendingTransactionsComponent;
}());
PendingTransactionsComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-pending-transactions',
        template: __webpack_require__("../../../../../src/app/components/pages/settings/pending-transactions/pending-transactions.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/settings/pending-transactions/pending-transactions.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_3__angular_router__["e" /* Router */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_3__angular_router__["e" /* Router */]) === "function" && _b || Object])
], PendingTransactionsComponent);

var _a, _b;
//# sourceMappingURL=pending-transactions.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/transaction/transaction.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, ".skycoin-details {\r\n  margin-top: 40px;\r\n  background-color: #eee;\r\n  margin-bottom: 20px;\r\n}\r\n\r\n.skycoin-detail-keys {\r\n  display: inline-block;\r\n}\r\n\r\n.skycoin-detail-values {\r\n  display: inline-block;\r\n}\r\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/transaction/transaction.component.html":
/***/ (function(module, exports) {

module.exports = "<md-card>\n  <md-card-title>Transaction <app-back-button></app-back-button></md-card-title>\n  <md-card class=\"skycoin-details\">\n    <div class=\"skycoin-detail-keys\">\n      ID: <br>\n      Timestamp: <br>\n      Block: <br>\n      Confirmed: <br>\n      Total value: <br>\n    </div>\n    <div class=\"skycoin-detail-values\">\n      {{ transaction && transaction.txn ? transaction.txn.txid : ''}} <br>\n      {{ transaction && transaction.txn ? transaction.txn.timestamp : ''}} <br>\n      {{ transaction && transaction.txn ? transaction.status.block_seq : ''}}<br>\n      {{ transaction && transaction.status ? (transaction.status.confirmed ? 'True' : 'False') : ''}} <br>\n      {{ total }}\n    </div>\n  </md-card>\n\n  <md-list *ngIf=\"transaction && transaction.txn\">\n    <h3 md-subheader>Inputs</h3>\n    <md-list-item *ngFor=\"let input of transaction.txn.inputs\">\n      <h4 md-line><button md-button [routerLink]=\"['/explorer/address/', input]\">{{ input }}</button></h4>\n    </md-list-item>\n    <md-divider></md-divider>\n    <h3 md-subheader>Outputs</h3>\n    <md-list-item *ngFor=\"let output of transaction.txn.outputs\">\n      <h4 md-line><button md-button [routerLink]=\"['/explorer/address/', output.dst]\">{{ output.dst }} ({{ output.coins }} SKY)</button></h4>\n    </md-list-item>\n  </md-list>\n</md-card>\n\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/transaction/transaction.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__ = __webpack_require__("../../../../../src/app/services/wallet.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_router__ = __webpack_require__("../../../router/@angular/router.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_switchMap__ = __webpack_require__("../../../../rxjs/add/operator/switchMap.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_switchMap___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_switchMap__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return TransactionComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};




var TransactionComponent = (function () {
    function TransactionComponent(route, walletService) {
        this.route = route;
        this.walletService = walletService;
    }
    TransactionComponent.prototype.ngOnInit = function () {
        var _this = this;
        this.route.params.switchMap(function (params) { return _this.walletService.transaction(params.transaction); }).subscribe(function (transaction) {
            _this.transaction = transaction;
            _this.total = transaction.txn.outputs.reduce(function (a, b) { return a + parseInt(b.coins); }, 0);
        });
    };
    return TransactionComponent;
}());
TransactionComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-transaction',
        template: __webpack_require__("../../../../../src/app/components/pages/transaction/transaction.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/transaction/transaction.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_2__angular_router__["d" /* ActivatedRoute */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_router__["d" /* ActivatedRoute */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */]) === "function" && _b || Object])
], TransactionComponent);

var _a, _b;
//# sourceMappingURL=transaction.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/wallets/address-detail/wallet-detail.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/wallets/address-detail/wallet-detail.component.html":
/***/ (function(module, exports) {

module.exports = "<md-list>\n  <h3 md-subheader>Addresses</h3>\n  <md-list-item *ngFor=\"let address of wallet.entries\">\n    <md-icon md-list-icon>note</md-icon>\n    <h4 md-line>{{address.address}} - {{ address.balance | sky }}</h4>\n  </md-list-item>\n  <div class=\"button-line\">\n    <a md-raised-button color=\"primary\" (click)=\"renameWallet()\">Rename wallet</a>\n    <a md-raised-button color=\"primary\" (click)=\"addAddress()\">Add address</a>\n  </div>\n</md-list>\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/wallets/address-detail/wallet-detail.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__ = __webpack_require__("../../../../../src/app/services/wallet.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__models_wallet_model__ = __webpack_require__("../../../../../src/app/models/wallet.model.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__models_wallet_model___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2__models_wallet_model__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_material__ = __webpack_require__("../../../material/@angular/material.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__change_name_change_name_component__ = __webpack_require__("../../../../../src/app/components/pages/wallets/change-name/change-name.component.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return WalletDetailComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var WalletDetailComponent = (function () {
    function WalletDetailComponent(walletService, dialog) {
        this.walletService = walletService;
        this.dialog = dialog;
    }
    WalletDetailComponent.prototype.addAddress = function () {
        var _this = this;
        this.walletService.addAddress(this.wallet).subscribe(function (output) { return _this.wallet.entries.push(output); });
    };
    WalletDetailComponent.prototype.renameWallet = function () {
        var _this = this;
        var config = new __WEBPACK_IMPORTED_MODULE_3__angular_material__["q" /* MdDialogConfig */]();
        config.width = '500px';
        config.data = this.wallet;
        this.dialog.open(__WEBPACK_IMPORTED_MODULE_4__change_name_change_name_component__["a" /* ChangeNameComponent */], config).afterClosed().subscribe(function (result) {
            if (result) {
                _this.wallet.meta.label = result;
            }
        });
    };
    return WalletDetailComponent;
}());
__decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"])(),
    __metadata("design:type", typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_2__models_wallet_model__["WalletModel"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__models_wallet_model__["WalletModel"]) === "function" && _a || Object)
], WalletDetailComponent.prototype, "wallet", void 0);
WalletDetailComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-wallet-detail',
        template: __webpack_require__("../../../../../src/app/components/pages/wallets/address-detail/wallet-detail.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/wallets/address-detail/wallet-detail.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_3__angular_material__["r" /* MdDialog */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_3__angular_material__["r" /* MdDialog */]) === "function" && _c || Object])
], WalletDetailComponent);

var _a, _b, _c;
//# sourceMappingURL=wallet-detail.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/wallets/change-name/change-name.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "md-input-container {\r\n  width: 100%;\r\n}\r\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/wallets/change-name/change-name.component.html":
/***/ (function(module, exports) {

module.exports = "<div [formGroup]=\"form\">\n  <md-input-container>\n    <input mdInput formControlName=\"label\" placeholder=\"Label\">\n  </md-input-container>\n</div>\n<div class=\"button-line\">\n  <a md-raised-button (click)=\"rename()\">Rename</a>\n</div>\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/wallets/change-name/change-name.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__ = __webpack_require__("../../../../../src/app/services/wallet.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_forms__ = __webpack_require__("../../../forms/@angular/forms.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_material__ = __webpack_require__("../../../material/@angular/material.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__models_wallet_model__ = __webpack_require__("../../../../../src/app/models/wallet.model.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4__models_wallet_model___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4__models_wallet_model__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return ChangeNameComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var __param = (this && this.__param) || function (paramIndex, decorator) {
    return function (target, key) { decorator(target, key, paramIndex); }
};





var ChangeNameComponent = (function () {
    function ChangeNameComponent(data, dialogRef, formBuilder, walletService) {
        this.data = data;
        this.dialogRef = dialogRef;
        this.formBuilder = formBuilder;
        this.walletService = walletService;
    }
    ChangeNameComponent.prototype.ngOnInit = function () {
        this.initForm();
    };
    ChangeNameComponent.prototype.rename = function () {
        var _this = this;
        this.walletService.renameWallet(this.data, this.form.value.label)
            .subscribe(function () { return _this.dialogRef.close(_this.form.value.label); });
    };
    ChangeNameComponent.prototype.initForm = function () {
        this.form = this.formBuilder.group({
            label: [this.data.meta.label, __WEBPACK_IMPORTED_MODULE_2__angular_forms__["h" /* Validators */].required],
        });
    };
    return ChangeNameComponent;
}());
ChangeNameComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-change-name',
        template: __webpack_require__("../../../../../src/app/components/pages/wallets/change-name/change-name.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/wallets/change-name/change-name.component.css")]
    }),
    __param(0, __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Inject"])(__WEBPACK_IMPORTED_MODULE_3__angular_material__["o" /* MD_DIALOG_DATA */])),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_4__models_wallet_model__["WalletModel"] !== "undefined" && __WEBPACK_IMPORTED_MODULE_4__models_wallet_model__["WalletModel"]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_3__angular_material__["p" /* MdDialogRef */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_3__angular_material__["p" /* MdDialogRef */]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_2__angular_forms__["i" /* FormBuilder */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_forms__["i" /* FormBuilder */]) === "function" && _c || Object, typeof (_d = typeof __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */]) === "function" && _d || Object])
], ChangeNameComponent);

var _a, _b, _c, _d;
//# sourceMappingURL=change-name.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/wallets/create-wallet/create-wallet.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "md-input-container {\r\n  width: 100%;\r\n}\r\n", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/wallets/create-wallet/create-wallet.component.html":
/***/ (function(module, exports) {

module.exports = "<div [formGroup]=\"form\">\n  <md-input-container>\n    <input mdInput formControlName=\"label\" placeholder=\"Label\">\n  </md-input-container>\n  <md-input-container>\n    <textarea mdInput formControlName=\"seed\" row=\"5\" placeholder=\"Seed\"></textarea>\n  </md-input-container>\n</div>\n<div class=\"button-line\">\n  <a md-raised-button (click)=\"generateSeed()\">New Seed</a>\n  <a md-raised-button color=\"primary\" (click)=\"createWallet()\">Create</a>\n</div>\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/wallets/create-wallet/create-wallet.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_forms__ = __webpack_require__("../../../forms/@angular/forms.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__services_wallet_service__ = __webpack_require__("../../../../../src/app/services/wallet.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__angular_material__ = __webpack_require__("../../../material/@angular/material.es5.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return CreateWalletComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};




var CreateWalletComponent = (function () {
    function CreateWalletComponent(dialogRef, formBuilder, walletService) {
        this.dialogRef = dialogRef;
        this.formBuilder = formBuilder;
        this.walletService = walletService;
    }
    CreateWalletComponent.prototype.ngOnInit = function () {
        this.initForm();
    };
    CreateWalletComponent.prototype.createWallet = function () {
        var _this = this;
        this.walletService.create(this.form.value.label, this.form.value.seed)
            .subscribe(function () { return _this.dialogRef.close(); });
    };
    CreateWalletComponent.prototype.generateSeed = function () {
        var _this = this;
        this.walletService.generateSeed().subscribe(function (seed) { return _this.form.controls.seed.setValue(seed); });
    };
    CreateWalletComponent.prototype.initForm = function () {
        this.form = this.formBuilder.group({
            label: ['', __WEBPACK_IMPORTED_MODULE_1__angular_forms__["h" /* Validators */].required],
            seed: ['', __WEBPACK_IMPORTED_MODULE_1__angular_forms__["h" /* Validators */].required],
        });
        this.generateSeed();
    };
    return CreateWalletComponent;
}());
CreateWalletComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-create-wallet',
        template: __webpack_require__("../../../../../src/app/components/pages/wallets/create-wallet/create-wallet.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/wallets/create-wallet/create-wallet.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_3__angular_material__["p" /* MdDialogRef */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_3__angular_material__["p" /* MdDialogRef */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_1__angular_forms__["i" /* FormBuilder */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_forms__["i" /* FormBuilder */]) === "function" && _b || Object, typeof (_c = typeof __WEBPACK_IMPORTED_MODULE_2__services_wallet_service__["a" /* WalletService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__services_wallet_service__["a" /* WalletService */]) === "function" && _c || Object])
], CreateWalletComponent);

var _a, _b, _c;
//# sourceMappingURL=create-wallet.component.js.map

/***/ }),

/***/ "../../../../../src/app/components/pages/wallets/wallets.component.css":
/***/ (function(module, exports, __webpack_require__) {

exports = module.exports = __webpack_require__("../../../../css-loader/lib/css-base.js")(false);
// imports


// module
exports.push([module.i, "", ""]);

// exports


/*** EXPORTS FROM exports-loader ***/
module.exports = module.exports.toString();

/***/ }),

/***/ "../../../../../src/app/components/pages/wallets/wallets.component.html":
/***/ (function(module, exports) {

module.exports = "<md-expansion-panel *ngFor=\"let wallet of walletService.all() | async\">\n  <md-expansion-panel-header>\n    <md-panel-title>\n      {{ wallet.meta.label }}\n    </md-panel-title>\n    <md-panel-description>\n      {{ wallet.balance | sky }}\n    </md-panel-description>\n  </md-expansion-panel-header>\n  <app-wallet-detail [wallet]=\"wallet\"></app-wallet-detail>\n</md-expansion-panel>\n<div class=\"button-line\">\n  <a md-raised-button color=\"primary\" (click)=\"addWallet()\">Add wallet</a>\n</div>\n"

/***/ }),

/***/ "../../../../../src/app/components/pages/wallets/wallets.component.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__ = __webpack_require__("../../../../../src/app/services/wallet.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__angular_material__ = __webpack_require__("../../../material/@angular/material.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__create_wallet_create_wallet_component__ = __webpack_require__("../../../../../src/app/components/pages/wallets/create-wallet/create-wallet.component.ts");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return WalletsComponent; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};




var WalletsComponent = (function () {
    function WalletsComponent(walletService, dialog) {
        this.walletService = walletService;
        this.dialog = dialog;
    }
    WalletsComponent.prototype.addWallet = function () {
        var config = new __WEBPACK_IMPORTED_MODULE_2__angular_material__["q" /* MdDialogConfig */]();
        config.width = '500px';
        this.dialog.open(__WEBPACK_IMPORTED_MODULE_3__create_wallet_create_wallet_component__["a" /* CreateWalletComponent */], config).afterClosed().subscribe(function (result) {
            //
        });
    };
    return WalletsComponent;
}());
WalletsComponent = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"])({
        selector: 'app-wallets',
        template: __webpack_require__("../../../../../src/app/components/pages/wallets/wallets.component.html"),
        styles: [__webpack_require__("../../../../../src/app/components/pages/wallets/wallets.component.css")]
    }),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__services_wallet_service__["a" /* WalletService */]) === "function" && _a || Object, typeof (_b = typeof __WEBPACK_IMPORTED_MODULE_2__angular_material__["r" /* MdDialog */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_2__angular_material__["r" /* MdDialog */]) === "function" && _b || Object])
], WalletsComponent);

var _a, _b;
//# sourceMappingURL=wallets.component.js.map

/***/ }),

/***/ "../../../../../src/app/models/wallet.model.ts":
/***/ (function(module, exports) {

//# sourceMappingURL=wallet.model.js.map

/***/ }),

/***/ "../../../../../src/app/pipes/date-from-now.pipe.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_moment__ = __webpack_require__("../../../../moment/moment.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_moment___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_moment__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return DateFromNowPipe; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};


var DateFromNowPipe = (function () {
    function DateFromNowPipe() {
    }
    DateFromNowPipe.prototype.transform = function (value) {
        return __WEBPACK_IMPORTED_MODULE_1_moment__["unix"](value).fromNow();
    };
    return DateFromNowPipe;
}());
DateFromNowPipe = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Pipe"])({
        name: 'dateFromNow'
    })
], DateFromNowPipe);

//# sourceMappingURL=date-from-now.pipe.js.map

/***/ }),

/***/ "../../../../../src/app/pipes/date-time.pipe.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_moment__ = __webpack_require__("../../../../moment/moment.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1_moment___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_1_moment__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return DateTimePipe; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};


var DateTimePipe = (function () {
    function DateTimePipe() {
    }
    DateTimePipe.prototype.transform = function (value) {
        return __WEBPACK_IMPORTED_MODULE_1_moment__["unix"](value).format('YYYY-MM-DD HH:mm');
    };
    return DateTimePipe;
}());
DateTimePipe = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Pipe"])({
        name: 'dateTime'
    })
], DateTimePipe);

//# sourceMappingURL=date-time.pipe.js.map

/***/ }),

/***/ "../../../../../src/app/pipes/sky.pipe.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return SkyPipe; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};

var SkyPipe = (function () {
    function SkyPipe() {
    }
    SkyPipe.prototype.transform = function (value) {
        if (value == null || value < 0) {
            return 'loading .. ';
        }
        else {
            return (value ? (value / 1000000) : 0) + ' SKY';
        }
    };
    return SkyPipe;
}());
SkyPipe = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Pipe"])({
        name: 'sky'
    })
], SkyPipe);

//# sourceMappingURL=sky.pipe.js.map

/***/ }),

/***/ "../../../../../src/app/pipes/transactions-amount.pipe.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return TransactionsAmountPipe; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};

var TransactionsAmountPipe = (function () {
    function TransactionsAmountPipe() {
    }
    TransactionsAmountPipe.prototype.transform = function (value) {
        return value.reduce(function (a, b) { return a + b.outputs.reduce(function (c, d) { return c + parseInt(d.coins); }, 0); }, 0);
    };
    return TransactionsAmountPipe;
}());
TransactionsAmountPipe = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Pipe"])({
        name: 'transactionsAmount'
    })
], TransactionsAmountPipe);

//# sourceMappingURL=transactions-amount.pipe.js.map

/***/ }),

/***/ "../../../../../src/app/services/api.service.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_http__ = __webpack_require__("../../../http/@angular/http.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__ = __webpack_require__("../../../../rxjs/Observable.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_observable_throw__ = __webpack_require__("../../../../rxjs/add/observable/throw.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_observable_throw___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_rxjs_add_observable_throw__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_catch__ = __webpack_require__("../../../../rxjs/add/operator/catch.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_catch___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_catch__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_rxjs_add_operator_map__ = __webpack_require__("../../../../rxjs/add/operator/map.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_rxjs_add_operator_map___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_5_rxjs_add_operator_map__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return ApiService; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};






var ApiService = (function () {
    function ApiService(http) {
        this.http = http;
        // private url = 'http://127.0.0.1:6420/'; // production
        this.url = '/api/'; // test
    }
    ApiService.prototype.get = function (url, options) {
        if (options === void 0) { options = null; }
        return this.http.get(this.getUrl(url, options), this.getHeaders())
            .map(function (res) { return res.json(); })
            .catch(function (error) { return __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__["Observable"].throw(error || 'Server error'); });
    };
    ApiService.prototype.post = function (url, options) {
        if (options === void 0) { options = {}; }
        return this.http.post(this.getUrl(url), this.getQueryString(options), this.returnRequestOptions())
            .map(function (res) { return res.json(); })
            .catch(function (error) { return __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__["Observable"].throw(error || 'Server error'); });
    };
    ApiService.prototype.getHeaders = function () {
        var headers = new __WEBPACK_IMPORTED_MODULE_1__angular_http__["b" /* Headers */]();
        headers.append('Content-Type', 'application/x-www-form-urlencoded');
        return headers;
    };
    ApiService.prototype.returnRequestOptions = function () {
        var options = new __WEBPACK_IMPORTED_MODULE_1__angular_http__["c" /* RequestOptions */]();
        options.headers = this.getHeaders();
        return options;
    };
    ApiService.prototype.getQueryString = function (parameters) {
        if (parameters === void 0) { parameters = null; }
        if (!parameters) {
            return '';
        }
        return Object.keys(parameters).reduce(function (array, key) {
            array.push(key + '=' + encodeURIComponent(parameters[key]));
            return array;
        }, []).join('&');
    };
    ApiService.prototype.getUrl = function (url, options) {
        if (options === void 0) { options = null; }
        return this.url + url + '?' + this.getQueryString(options);
    };
    return ApiService;
}());
ApiService = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Injectable"])(),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__angular_http__["d" /* Http */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__angular_http__["d" /* Http */]) === "function" && _a || Object])
], ApiService);

var _a;
//# sourceMappingURL=api.service.js.map

/***/ }),

/***/ "../../../../../src/app/services/blockchain.service.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__api_service__ = __webpack_require__("../../../../../src/app/services/api.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__ = __webpack_require__("../../../../rxjs/Observable.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return BlockchainService; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var BlockchainService = (function () {
    function BlockchainService(apiService) {
        this.apiService = apiService;
    }
    BlockchainService.prototype.addressTransactions = function (id) {
        return this.apiService.get('explorer/address', { address: id });
    };
    BlockchainService.prototype.addressBalance = function (id) {
        return this.apiService.get('outputs', { addrs: id });
    };
    BlockchainService.prototype.block = function (id) {
        var _this = this;
        return this.apiService.get('blocks', { start: id, end: id }).map(function (response) { return response.blocks[0]; }).flatMap(function (block) {
            return __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__["Observable"].forkJoin(block.body.txns.map(function (transaction) {
                return __WEBPACK_IMPORTED_MODULE_2_rxjs_Observable__["Observable"].forkJoin(transaction.inputs.map(function (input) { return _this.retrieveInputAddress(input).map(function (response) {
                    return response.owner_address;
                }); })).map(function (inputs) {
                    transaction.inputs = inputs;
                    return transaction;
                });
            })).map(function (transactions) {
                block.body.txns = transactions;
                return block;
            });
        });
    };
    BlockchainService.prototype.blocks = function (num) {
        if (num === void 0) { num = 100; }
        return this.apiService.get('last_blocks', { num: num }).map(function (response) { return response.blocks.reverse(); });
    };
    BlockchainService.prototype.lastBlock = function () {
        return this.blocks(1).map(function (blocks) { return blocks[0]; });
    };
    BlockchainService.prototype.retrieveInputAddress = function (input) {
        return this.apiService.get('uxout', { uxid: input });
    };
    return BlockchainService;
}());
BlockchainService = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Injectable"])(),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__api_service__["a" /* ApiService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__api_service__["a" /* ApiService */]) === "function" && _a || Object])
], BlockchainService);

var _a;
//# sourceMappingURL=blockchain.service.js.map

/***/ }),

/***/ "../../../../../src/app/services/network.service.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__api_service__ = __webpack_require__("../../../../../src/app/services/api.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_BehaviorSubject__ = __webpack_require__("../../../../rxjs/BehaviorSubject.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_BehaviorSubject___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs_BehaviorSubject__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_observable_IntervalObservable__ = __webpack_require__("../../../../rxjs/observable/IntervalObservable.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_observable_IntervalObservable___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_rxjs_observable_IntervalObservable__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_mergeMap__ = __webpack_require__("../../../../rxjs/add/operator/mergeMap.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_mergeMap___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4_rxjs_add_operator_mergeMap__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return NetworkService; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};





var NetworkService = (function () {
    function NetworkService(apiService) {
        this.apiService = apiService;
        this.automaticPeers = new __WEBPACK_IMPORTED_MODULE_2_rxjs_BehaviorSubject__["BehaviorSubject"]([]);
        this.loadData();
    }
    NetworkService.prototype.automatic = function () {
        return this.automaticPeers.asObservable();
    };
    NetworkService.prototype.retrieveDefaultConnections = function () {
        return this.apiService.post('network/defaultConnections')
            .map(function (output) { return output.map(function (address, index) { return ({
            id: index + 1,
            address: address,
            listen_port: 6000,
        }); }); });
    };
    NetworkService.prototype.loadData = function () {
        var _this = this;
        this.retrieveConnections().subscribe(function (connections) { return _this.automaticPeers.next(connections); });
        __WEBPACK_IMPORTED_MODULE_3_rxjs_observable_IntervalObservable__["IntervalObservable"]
            .create(5000)
            .flatMap(function () { return _this.retrieveConnections(); })
            .subscribe(function (connections) { return _this.automaticPeers.next(connections); });
    };
    NetworkService.prototype.retrieveConnections = function () {
        return this.apiService.post('network/connections')
            .map(function (response) { return response.connections.sort(function (a, b) { return a.id - b.id; }); });
    };
    return NetworkService;
}());
NetworkService = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Injectable"])(),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__api_service__["a" /* ApiService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__api_service__["a" /* ApiService */]) === "function" && _a || Object])
], NetworkService);

var _a;
//# sourceMappingURL=network.service.js.map

/***/ }),

/***/ "../../../../../src/app/services/wallet.service.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__api_service__ = __webpack_require__("../../../../../src/app/services/api.service.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_BehaviorSubject__ = __webpack_require__("../../../../rxjs/BehaviorSubject.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_BehaviorSubject___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_2_rxjs_BehaviorSubject__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_Observable__ = __webpack_require__("../../../../rxjs/Observable.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_Observable___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_3_rxjs_Observable__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_observable_IntervalObservable__ = __webpack_require__("../../../../rxjs/observable/IntervalObservable.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_4_rxjs_observable_IntervalObservable___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_4_rxjs_observable_IntervalObservable__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_rxjs_add_observable_forkJoin__ = __webpack_require__("../../../../rxjs/add/observable/forkJoin.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_5_rxjs_add_observable_forkJoin___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_5_rxjs_add_observable_forkJoin__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6_rxjs_add_operator_do__ = __webpack_require__("../../../../rxjs/add/operator/do.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_6_rxjs_add_operator_do___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_6_rxjs_add_operator_do__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7_rxjs_add_operator_first__ = __webpack_require__("../../../../rxjs/add/operator/first.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_7_rxjs_add_operator_first___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_7_rxjs_add_operator_first__);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8_rxjs_add_operator_mergeMap__ = __webpack_require__("../../../../rxjs/add/operator/mergeMap.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_8_rxjs_add_operator_mergeMap___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_8_rxjs_add_operator_mergeMap__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return WalletService; });
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};









var WalletService = (function () {
    function WalletService(apiService) {
        var _this = this;
        this.apiService = apiService;
        this.recentTransactions = new __WEBPACK_IMPORTED_MODULE_2_rxjs_BehaviorSubject__["BehaviorSubject"]([]);
        this.transactions = new __WEBPACK_IMPORTED_MODULE_2_rxjs_BehaviorSubject__["BehaviorSubject"]([]);
        this.wallets = new __WEBPACK_IMPORTED_MODULE_2_rxjs_BehaviorSubject__["BehaviorSubject"]([]);
        this.loadData();
        __WEBPACK_IMPORTED_MODULE_4_rxjs_observable_IntervalObservable__["IntervalObservable"]
            .create(30000)
            .subscribe(function () { return _this.refreshBalances(); });
    }
    WalletService.prototype.addressesAsString = function () {
        return this.all().map(function (wallets) { return wallets.map(function (wallet) {
            return wallet.entries.reduce(function (a, b) {
                a.push(b.address);
                return a;
            }, []).join(',');
        }).join(','); });
    };
    WalletService.prototype.addAddress = function (wallet) {
        return this.apiService.post('wallet/newAddress', { id: wallet.meta.filename });
    };
    WalletService.prototype.all = function () {
        return this.wallets.asObservable();
    };
    WalletService.prototype.create = function (label, seed) {
        var _this = this;
        return this.apiService.post('wallet/create', { label: label ? label : 'undefined', seed: seed })
            .do(function (wallet) {
            _this.wallets.first().subscribe(function (wallets) {
                wallets.push(wallet);
                _this.wallets.next(wallets);
                _this.refreshBalances();
            });
        });
    };
    WalletService.prototype.folder = function () {
        return this.apiService.get('wallets/folderName').map(function (response) { return response.address; });
    };
    WalletService.prototype.generateSeed = function () {
        return this.apiService.get('wallet/newSeed').map(function (response) { return response.seed; });
    };
    WalletService.prototype.history = function () {
        return this.transactions.asObservable();
    };
    WalletService.prototype.outputs = function () {
        var _this = this;
        return this.addressesAsString()
            .filter(function (addresses) { return !!addresses; })
            .flatMap(function (addresses) { return _this.apiService.get('outputs', { addrs: addresses }); });
    };
    WalletService.prototype.pendingTransactions = function () {
        return this.apiService.get('pendingTxs');
    };
    WalletService.prototype.recent = function () {
        return this.recentTransactions.asObservable();
    };
    WalletService.prototype.renameWallet = function (wallet, label) {
        return this.apiService.post('wallet/update', { id: wallet.meta.filename, label: label });
    };
    WalletService.prototype.retrieveUpdatedTransactions = function (transactions) {
        var _this = this;
        return __WEBPACK_IMPORTED_MODULE_3_rxjs_Observable__["Observable"].forkJoin((transactions.map(function (transaction) {
            return _this.apiService.get('transaction', { txid: transaction.id }).map(function (response) {
                response.amount = transaction.amount;
                response.address = transaction.address;
                return response;
            });
        })));
    };
    WalletService.prototype.sendSkycoin = function (wallet_id, address, amount) {
        var _this = this;
        return this.apiService.post('wallet/spend', { id: wallet_id, dst: address, coins: amount })
            .do(function (output) { return _this.recentTransactions.first().subscribe(function (transactions) {
            var transaction = { id: output.txn.txid, address: address, amount: amount / 1000000 };
            transactions.push(transaction);
            _this.recentTransactions.next(transactions);
        }); });
    };
    WalletService.prototype.sum = function () {
        return this.all().map(function (wallets) { return wallets.map(function (wallet) { return wallet.balance >= 0 ? wallet.balance : 0; }).reduce(function (a, b) { return a + b; }, 0); });
    };
    WalletService.prototype.transaction = function (txid) {
        var _this = this;
        return this.apiService.get('transaction', { txid: txid }).flatMap(function (transaction) {
            return __WEBPACK_IMPORTED_MODULE_3_rxjs_Observable__["Observable"].forkJoin(transaction.txn.inputs.map(function (input) { return _this.retrieveInputAddress(input).map(function (response) {
                return response.owner_address;
            }); })).map(function (inputs) {
                transaction.txn.inputs = inputs;
                return transaction;
            });
        });
    };
    WalletService.prototype.loadData = function () {
        var _this = this;
        this.retrieveWallets().first().subscribe(function (wallets) {
            _this.wallets.next(wallets);
            _this.refreshBalances();
            // this.retrieveHistory();
            _this.retrieveTransactions();
        });
    };
    WalletService.prototype.refreshBalances = function () {
        var _this = this;
        this.wallets.first().subscribe(function (wallets) {
            __WEBPACK_IMPORTED_MODULE_3_rxjs_Observable__["Observable"].forkJoin(wallets.map(function (wallet) { return _this.retrieveWalletBalance(wallet).map(function (response) {
                wallet.entries = response;
                wallet.balance = response.map(function (address) { return address.balance >= 0 ? address.balance : 0; }).reduce(function (a, b) { return a + b; }, 0);
                return wallet;
            }); }))
                .subscribe(function (newWallets) { return _this.wallets.next(newWallets); });
        });
    };
    WalletService.prototype.retrieveAddressBalance = function (address) {
        var addresses = Array.isArray(address) ? address.map(function (address) { return address.address; }).join(',') : address.address;
        return this.apiService.get('balance', { addrs: addresses });
    };
    WalletService.prototype.retrieveAddressTransactions = function (address) {
        return this.apiService.get('explorer/address', { address: address.address });
    };
    WalletService.prototype.retrieveInputAddress = function (input) {
        return this.apiService.get('uxout', { uxid: input });
    };
    WalletService.prototype.retrieveTransactions = function () {
        var _this = this;
        return this.wallets.first().subscribe(function (wallets) {
            __WEBPACK_IMPORTED_MODULE_3_rxjs_Observable__["Observable"].forkJoin(wallets.map(function (wallet) { return _this.retrieveWalletTransactions(wallet); }))
                .map(function (addresses) { return [].concat.apply([], addresses).sort(function (a, b) { return b.timestamp - a.timestamp; }); })
                .subscribe(function (transactions) { return _this.transactions.next(transactions); });
        });
    };
    WalletService.prototype.retrieveWalletBalance = function (wallet) {
        var _this = this;
        return __WEBPACK_IMPORTED_MODULE_3_rxjs_Observable__["Observable"].forkJoin(wallet.entries.map(function (address) { return _this.retrieveAddressBalance(address).map(function (balance) {
            address.balance = balance.confirmed.coins;
            return address;
        }); }));
    };
    WalletService.prototype.retrieveWalletTransactions = function (wallet) {
        var _this = this;
        return __WEBPACK_IMPORTED_MODULE_3_rxjs_Observable__["Observable"].forkJoin(wallet.entries.map(function (address) { return _this.retrieveAddressTransactions(address); }))
            .map(function (addresses) { return [].concat.apply([], addresses); });
    };
    WalletService.prototype.retrieveWalletUnconfirmedTransactions = function (wallet) {
        return this.apiService.get('wallet/transactions', { id: wallet.meta.filename });
    };
    WalletService.prototype.retrieveWallets = function () {
        return this.apiService.get('wallets');
    };
    return WalletService;
}());
WalletService = __decorate([
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["Injectable"])(),
    __metadata("design:paramtypes", [typeof (_a = typeof __WEBPACK_IMPORTED_MODULE_1__api_service__["a" /* ApiService */] !== "undefined" && __WEBPACK_IMPORTED_MODULE_1__api_service__["a" /* ApiService */]) === "function" && _a || Object])
], WalletService);

var _a;
//# sourceMappingURL=wallet.service.js.map

/***/ }),

/***/ "../../../../../src/environments/environment.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "a", function() { return environment; });
// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.
// The file contents for the current environment will overwrite these during build.
var environment = {
    production: false
};
//# sourceMappingURL=environment.js.map

/***/ }),

/***/ "../../../../../src/main.ts":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("../../../core/@angular/core.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__ = __webpack_require__("../../../platform-browser-dynamic/@angular/platform-browser-dynamic.es5.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2__app_app_module__ = __webpack_require__("../../../../../src/app/app.module.ts");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3__environments_environment__ = __webpack_require__("../../../../../src/environments/environment.ts");




if (__WEBPACK_IMPORTED_MODULE_3__environments_environment__["a" /* environment */].production) {
    __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__angular_core__["enableProdMode"])();
}
__webpack_require__.i(__WEBPACK_IMPORTED_MODULE_1__angular_platform_browser_dynamic__["a" /* platformBrowserDynamic */])().bootstrapModule(__WEBPACK_IMPORTED_MODULE_2__app_app_module__["a" /* AppModule */]);
//# sourceMappingURL=main.js.map

/***/ }),

/***/ "../../../../moment/locale recursive ^\\.\\/.*$":
/***/ (function(module, exports, __webpack_require__) {

var map = {
	"./af": "../../../../moment/locale/af.js",
	"./af.js": "../../../../moment/locale/af.js",
	"./ar": "../../../../moment/locale/ar.js",
	"./ar-dz": "../../../../moment/locale/ar-dz.js",
	"./ar-dz.js": "../../../../moment/locale/ar-dz.js",
	"./ar-kw": "../../../../moment/locale/ar-kw.js",
	"./ar-kw.js": "../../../../moment/locale/ar-kw.js",
	"./ar-ly": "../../../../moment/locale/ar-ly.js",
	"./ar-ly.js": "../../../../moment/locale/ar-ly.js",
	"./ar-ma": "../../../../moment/locale/ar-ma.js",
	"./ar-ma.js": "../../../../moment/locale/ar-ma.js",
	"./ar-sa": "../../../../moment/locale/ar-sa.js",
	"./ar-sa.js": "../../../../moment/locale/ar-sa.js",
	"./ar-tn": "../../../../moment/locale/ar-tn.js",
	"./ar-tn.js": "../../../../moment/locale/ar-tn.js",
	"./ar.js": "../../../../moment/locale/ar.js",
	"./az": "../../../../moment/locale/az.js",
	"./az.js": "../../../../moment/locale/az.js",
	"./be": "../../../../moment/locale/be.js",
	"./be.js": "../../../../moment/locale/be.js",
	"./bg": "../../../../moment/locale/bg.js",
	"./bg.js": "../../../../moment/locale/bg.js",
	"./bn": "../../../../moment/locale/bn.js",
	"./bn.js": "../../../../moment/locale/bn.js",
	"./bo": "../../../../moment/locale/bo.js",
	"./bo.js": "../../../../moment/locale/bo.js",
	"./br": "../../../../moment/locale/br.js",
	"./br.js": "../../../../moment/locale/br.js",
	"./bs": "../../../../moment/locale/bs.js",
	"./bs.js": "../../../../moment/locale/bs.js",
	"./ca": "../../../../moment/locale/ca.js",
	"./ca.js": "../../../../moment/locale/ca.js",
	"./cs": "../../../../moment/locale/cs.js",
	"./cs.js": "../../../../moment/locale/cs.js",
	"./cv": "../../../../moment/locale/cv.js",
	"./cv.js": "../../../../moment/locale/cv.js",
	"./cy": "../../../../moment/locale/cy.js",
	"./cy.js": "../../../../moment/locale/cy.js",
	"./da": "../../../../moment/locale/da.js",
	"./da.js": "../../../../moment/locale/da.js",
	"./de": "../../../../moment/locale/de.js",
	"./de-at": "../../../../moment/locale/de-at.js",
	"./de-at.js": "../../../../moment/locale/de-at.js",
	"./de-ch": "../../../../moment/locale/de-ch.js",
	"./de-ch.js": "../../../../moment/locale/de-ch.js",
	"./de.js": "../../../../moment/locale/de.js",
	"./dv": "../../../../moment/locale/dv.js",
	"./dv.js": "../../../../moment/locale/dv.js",
	"./el": "../../../../moment/locale/el.js",
	"./el.js": "../../../../moment/locale/el.js",
	"./en-au": "../../../../moment/locale/en-au.js",
	"./en-au.js": "../../../../moment/locale/en-au.js",
	"./en-ca": "../../../../moment/locale/en-ca.js",
	"./en-ca.js": "../../../../moment/locale/en-ca.js",
	"./en-gb": "../../../../moment/locale/en-gb.js",
	"./en-gb.js": "../../../../moment/locale/en-gb.js",
	"./en-ie": "../../../../moment/locale/en-ie.js",
	"./en-ie.js": "../../../../moment/locale/en-ie.js",
	"./en-nz": "../../../../moment/locale/en-nz.js",
	"./en-nz.js": "../../../../moment/locale/en-nz.js",
	"./eo": "../../../../moment/locale/eo.js",
	"./eo.js": "../../../../moment/locale/eo.js",
	"./es": "../../../../moment/locale/es.js",
	"./es-do": "../../../../moment/locale/es-do.js",
	"./es-do.js": "../../../../moment/locale/es-do.js",
	"./es.js": "../../../../moment/locale/es.js",
	"./et": "../../../../moment/locale/et.js",
	"./et.js": "../../../../moment/locale/et.js",
	"./eu": "../../../../moment/locale/eu.js",
	"./eu.js": "../../../../moment/locale/eu.js",
	"./fa": "../../../../moment/locale/fa.js",
	"./fa.js": "../../../../moment/locale/fa.js",
	"./fi": "../../../../moment/locale/fi.js",
	"./fi.js": "../../../../moment/locale/fi.js",
	"./fo": "../../../../moment/locale/fo.js",
	"./fo.js": "../../../../moment/locale/fo.js",
	"./fr": "../../../../moment/locale/fr.js",
	"./fr-ca": "../../../../moment/locale/fr-ca.js",
	"./fr-ca.js": "../../../../moment/locale/fr-ca.js",
	"./fr-ch": "../../../../moment/locale/fr-ch.js",
	"./fr-ch.js": "../../../../moment/locale/fr-ch.js",
	"./fr.js": "../../../../moment/locale/fr.js",
	"./fy": "../../../../moment/locale/fy.js",
	"./fy.js": "../../../../moment/locale/fy.js",
	"./gd": "../../../../moment/locale/gd.js",
	"./gd.js": "../../../../moment/locale/gd.js",
	"./gl": "../../../../moment/locale/gl.js",
	"./gl.js": "../../../../moment/locale/gl.js",
	"./gom-latn": "../../../../moment/locale/gom-latn.js",
	"./gom-latn.js": "../../../../moment/locale/gom-latn.js",
	"./he": "../../../../moment/locale/he.js",
	"./he.js": "../../../../moment/locale/he.js",
	"./hi": "../../../../moment/locale/hi.js",
	"./hi.js": "../../../../moment/locale/hi.js",
	"./hr": "../../../../moment/locale/hr.js",
	"./hr.js": "../../../../moment/locale/hr.js",
	"./hu": "../../../../moment/locale/hu.js",
	"./hu.js": "../../../../moment/locale/hu.js",
	"./hy-am": "../../../../moment/locale/hy-am.js",
	"./hy-am.js": "../../../../moment/locale/hy-am.js",
	"./id": "../../../../moment/locale/id.js",
	"./id.js": "../../../../moment/locale/id.js",
	"./is": "../../../../moment/locale/is.js",
	"./is.js": "../../../../moment/locale/is.js",
	"./it": "../../../../moment/locale/it.js",
	"./it.js": "../../../../moment/locale/it.js",
	"./ja": "../../../../moment/locale/ja.js",
	"./ja.js": "../../../../moment/locale/ja.js",
	"./jv": "../../../../moment/locale/jv.js",
	"./jv.js": "../../../../moment/locale/jv.js",
	"./ka": "../../../../moment/locale/ka.js",
	"./ka.js": "../../../../moment/locale/ka.js",
	"./kk": "../../../../moment/locale/kk.js",
	"./kk.js": "../../../../moment/locale/kk.js",
	"./km": "../../../../moment/locale/km.js",
	"./km.js": "../../../../moment/locale/km.js",
	"./kn": "../../../../moment/locale/kn.js",
	"./kn.js": "../../../../moment/locale/kn.js",
	"./ko": "../../../../moment/locale/ko.js",
	"./ko.js": "../../../../moment/locale/ko.js",
	"./ky": "../../../../moment/locale/ky.js",
	"./ky.js": "../../../../moment/locale/ky.js",
	"./lb": "../../../../moment/locale/lb.js",
	"./lb.js": "../../../../moment/locale/lb.js",
	"./lo": "../../../../moment/locale/lo.js",
	"./lo.js": "../../../../moment/locale/lo.js",
	"./lt": "../../../../moment/locale/lt.js",
	"./lt.js": "../../../../moment/locale/lt.js",
	"./lv": "../../../../moment/locale/lv.js",
	"./lv.js": "../../../../moment/locale/lv.js",
	"./me": "../../../../moment/locale/me.js",
	"./me.js": "../../../../moment/locale/me.js",
	"./mi": "../../../../moment/locale/mi.js",
	"./mi.js": "../../../../moment/locale/mi.js",
	"./mk": "../../../../moment/locale/mk.js",
	"./mk.js": "../../../../moment/locale/mk.js",
	"./ml": "../../../../moment/locale/ml.js",
	"./ml.js": "../../../../moment/locale/ml.js",
	"./mr": "../../../../moment/locale/mr.js",
	"./mr.js": "../../../../moment/locale/mr.js",
	"./ms": "../../../../moment/locale/ms.js",
	"./ms-my": "../../../../moment/locale/ms-my.js",
	"./ms-my.js": "../../../../moment/locale/ms-my.js",
	"./ms.js": "../../../../moment/locale/ms.js",
	"./my": "../../../../moment/locale/my.js",
	"./my.js": "../../../../moment/locale/my.js",
	"./nb": "../../../../moment/locale/nb.js",
	"./nb.js": "../../../../moment/locale/nb.js",
	"./ne": "../../../../moment/locale/ne.js",
	"./ne.js": "../../../../moment/locale/ne.js",
	"./nl": "../../../../moment/locale/nl.js",
	"./nl-be": "../../../../moment/locale/nl-be.js",
	"./nl-be.js": "../../../../moment/locale/nl-be.js",
	"./nl.js": "../../../../moment/locale/nl.js",
	"./nn": "../../../../moment/locale/nn.js",
	"./nn.js": "../../../../moment/locale/nn.js",
	"./pa-in": "../../../../moment/locale/pa-in.js",
	"./pa-in.js": "../../../../moment/locale/pa-in.js",
	"./pl": "../../../../moment/locale/pl.js",
	"./pl.js": "../../../../moment/locale/pl.js",
	"./pt": "../../../../moment/locale/pt.js",
	"./pt-br": "../../../../moment/locale/pt-br.js",
	"./pt-br.js": "../../../../moment/locale/pt-br.js",
	"./pt.js": "../../../../moment/locale/pt.js",
	"./ro": "../../../../moment/locale/ro.js",
	"./ro.js": "../../../../moment/locale/ro.js",
	"./ru": "../../../../moment/locale/ru.js",
	"./ru.js": "../../../../moment/locale/ru.js",
	"./sd": "../../../../moment/locale/sd.js",
	"./sd.js": "../../../../moment/locale/sd.js",
	"./se": "../../../../moment/locale/se.js",
	"./se.js": "../../../../moment/locale/se.js",
	"./si": "../../../../moment/locale/si.js",
	"./si.js": "../../../../moment/locale/si.js",
	"./sk": "../../../../moment/locale/sk.js",
	"./sk.js": "../../../../moment/locale/sk.js",
	"./sl": "../../../../moment/locale/sl.js",
	"./sl.js": "../../../../moment/locale/sl.js",
	"./sq": "../../../../moment/locale/sq.js",
	"./sq.js": "../../../../moment/locale/sq.js",
	"./sr": "../../../../moment/locale/sr.js",
	"./sr-cyrl": "../../../../moment/locale/sr-cyrl.js",
	"./sr-cyrl.js": "../../../../moment/locale/sr-cyrl.js",
	"./sr.js": "../../../../moment/locale/sr.js",
	"./ss": "../../../../moment/locale/ss.js",
	"./ss.js": "../../../../moment/locale/ss.js",
	"./sv": "../../../../moment/locale/sv.js",
	"./sv.js": "../../../../moment/locale/sv.js",
	"./sw": "../../../../moment/locale/sw.js",
	"./sw.js": "../../../../moment/locale/sw.js",
	"./ta": "../../../../moment/locale/ta.js",
	"./ta.js": "../../../../moment/locale/ta.js",
	"./te": "../../../../moment/locale/te.js",
	"./te.js": "../../../../moment/locale/te.js",
	"./tet": "../../../../moment/locale/tet.js",
	"./tet.js": "../../../../moment/locale/tet.js",
	"./th": "../../../../moment/locale/th.js",
	"./th.js": "../../../../moment/locale/th.js",
	"./tl-ph": "../../../../moment/locale/tl-ph.js",
	"./tl-ph.js": "../../../../moment/locale/tl-ph.js",
	"./tlh": "../../../../moment/locale/tlh.js",
	"./tlh.js": "../../../../moment/locale/tlh.js",
	"./tr": "../../../../moment/locale/tr.js",
	"./tr.js": "../../../../moment/locale/tr.js",
	"./tzl": "../../../../moment/locale/tzl.js",
	"./tzl.js": "../../../../moment/locale/tzl.js",
	"./tzm": "../../../../moment/locale/tzm.js",
	"./tzm-latn": "../../../../moment/locale/tzm-latn.js",
	"./tzm-latn.js": "../../../../moment/locale/tzm-latn.js",
	"./tzm.js": "../../../../moment/locale/tzm.js",
	"./uk": "../../../../moment/locale/uk.js",
	"./uk.js": "../../../../moment/locale/uk.js",
	"./ur": "../../../../moment/locale/ur.js",
	"./ur.js": "../../../../moment/locale/ur.js",
	"./uz": "../../../../moment/locale/uz.js",
	"./uz-latn": "../../../../moment/locale/uz-latn.js",
	"./uz-latn.js": "../../../../moment/locale/uz-latn.js",
	"./uz.js": "../../../../moment/locale/uz.js",
	"./vi": "../../../../moment/locale/vi.js",
	"./vi.js": "../../../../moment/locale/vi.js",
	"./x-pseudo": "../../../../moment/locale/x-pseudo.js",
	"./x-pseudo.js": "../../../../moment/locale/x-pseudo.js",
	"./yo": "../../../../moment/locale/yo.js",
	"./yo.js": "../../../../moment/locale/yo.js",
	"./zh-cn": "../../../../moment/locale/zh-cn.js",
	"./zh-cn.js": "../../../../moment/locale/zh-cn.js",
	"./zh-hk": "../../../../moment/locale/zh-hk.js",
	"./zh-hk.js": "../../../../moment/locale/zh-hk.js",
	"./zh-tw": "../../../../moment/locale/zh-tw.js",
	"./zh-tw.js": "../../../../moment/locale/zh-tw.js"
};
function webpackContext(req) {
	return __webpack_require__(webpackContextResolve(req));
};
function webpackContextResolve(req) {
	var id = map[req];
	if(!(id + 1)) // check for number or string
		throw new Error("Cannot find module '" + req + "'.");
	return id;
};
webpackContext.keys = function webpackContextKeys() {
	return Object.keys(map);
};
webpackContext.resolve = webpackContextResolve;
module.exports = webpackContext;
webpackContext.id = "../../../../moment/locale recursive ^\\.\\/.*$";

/***/ }),

/***/ 0:
/***/ (function(module, exports, __webpack_require__) {

module.exports = __webpack_require__("../../../../../src/main.ts");


/***/ })

},[0]);
//# sourceMappingURL=main.bundle.js.map