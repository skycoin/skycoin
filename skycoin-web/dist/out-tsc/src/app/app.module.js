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
var platform_browser_1 = require('@angular/platform-browser');
var core_1 = require('@angular/core');
var forms_1 = require('@angular/forms');
var http_1 = require('@angular/http');
var router_1 = require('@angular/router');
var app_component_1 = require('./app.component');
var skycoin_header_component_1 = require('./components/skycoin-header/skycoin-header.component');
var block_chain_table_component_1 = require('./components/block-chain-table/block-chain-table.component');
var skycoin_pagination_component_1 = require('./components/skycoin-pagination/skycoin-pagination.component');
var skycoin_search_bar_component_1 = require('./components/skycoin-search-bar/skycoin-search-bar.component');
var block_chain_service_1 = require("./components/block-chain-table/block-chain.service");
var footer_component_1 = require('./components/footer/footer.component');
var num_pages_pipe_1 = require('./components/skycoin-pagination/num-pages.pipe');
var skycoin_blockchain_pagination_service_1 = require("./components/skycoin-pagination/skycoin-blockchain-pagination.service");
var common_1 = require("@angular/common");
var angular2_moment_1 = require("angular2-moment");
var block_details_component_1 = require('./components/block-details/block-details.component');
var address_detail_component_1 = require('./components/address-detail/address-detail.component');
var UxOutputs_service_1 = require("./components/address-detail/UxOutputs.service");
var transaction_detail_component_1 = require('./components/transaction-detail/transaction-detail.component');
var transaction_detail_service_1 = require("./components/transaction-detail/transaction-detail.service");
var loading_component_1 = require('./components/loading/loading.component');
var block_chain_coin_supply_component_1 = require('./components/block-chain-coin-supply/block-chain-coin-supply.component');
var coin_supply_service_1 = require("./components/block-chain-coin-supply/coin-supply.service");
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
        component: block_chain_table_component_1.BlockChainTableComponent
    },
    {
        path: 'block/:id',
        component: block_details_component_1.BlockDetailsComponent
    },
    {
        path: 'address/:address',
        component: address_detail_component_1.AddressDetailComponent
    },
    {
        path: 'transaction/:txid',
        component: transaction_detail_component_1.TransactionDetailComponent
    }
];
var AppModule = (function () {
    function AppModule() {
    }
    AppModule = __decorate([
        core_1.NgModule({
            declarations: [
                app_component_1.AppComponent,
                skycoin_header_component_1.SkycoinHeaderComponent,
                block_chain_table_component_1.BlockChainTableComponent,
                skycoin_pagination_component_1.SkycoinPaginationComponent,
                skycoin_search_bar_component_1.SkycoinSearchBarComponent,
                footer_component_1.FooterComponent,
                num_pages_pipe_1.NumPagesPipe,
                block_details_component_1.BlockDetailsComponent,
                address_detail_component_1.AddressDetailComponent,
                transaction_detail_component_1.TransactionDetailComponent,
                loading_component_1.LoadingComponent,
                block_chain_coin_supply_component_1.BlockChainCoinSupplyComponent
            ],
            imports: [
                common_1.CommonModule,
                platform_browser_1.BrowserModule,
                forms_1.FormsModule,
                http_1.HttpModule,
                angular2_moment_1.MomentModule,
                router_1.RouterModule.forRoot(ROUTES)
            ],
            providers: [block_chain_service_1.BlockChainService, skycoin_blockchain_pagination_service_1.SkycoinBlockchainPaginationService, UxOutputs_service_1.UxOutputsService, transaction_detail_service_1.TransactionDetailService, coin_supply_service_1.CoinSupplyService],
            bootstrap: [app_component_1.AppComponent]
        }), 
        __metadata('design:paramtypes', [])
    ], AppModule);
    return AppModule;
}());
exports.AppModule = AppModule;
//# sourceMappingURL=app.module.js.map