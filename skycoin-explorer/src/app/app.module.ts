import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { RouterModule } from '@angular/router';

import { AppComponent } from './app.component';
import { SkycoinHeaderComponent } from './components/skycoin-header/skycoin-header.component';
import { BlockChainTableComponent } from './components/block-chain-table/block-chain-table.component';
import { SkycoinPaginationComponent } from './components/skycoin-pagination/skycoin-pagination.component';
import { SkycoinSearchBarComponent } from './components/skycoin-search-bar/skycoin-search-bar.component';
import {BlockChainService} from "./components/block-chain-table/block-chain.service";
import { FooterComponent } from './components/footer/footer.component';
import { NumPagesPipe } from './components/skycoin-pagination/num-pages.pipe';
import {SkycoinBlockchainPaginationService} from "./components/skycoin-pagination/skycoin-blockchain-pagination.service";
import {CommonModule} from "@angular/common";
import {MomentModule} from "angular2-moment";
import { BlockDetailsComponent } from './components/block-details/block-details.component';
import { AddressDetailComponent } from './components/address-detail/address-detail.component';
import {UxOutputsService} from "./components/address-detail/UxOutputs.service";
import { TransactionDetailComponent } from './components/transaction-detail/transaction-detail.component';
import {TransactionDetailService} from "./components/transaction-detail/transaction-detail.service";
import { LoadingComponent } from './components/loading/loading.component';

import { BlockChainCoinSupplyComponent } from './components/block-chain-coin-supply/block-chain-coin-supply.component';
import {CoinSupplyService} from "./components/block-chain-coin-supply/coin-supply.service";
// import {QRCodeModule} from "../js/angular2-qrcode";
// import {QRCodeModule} from "angular2-qrcode";


const ROUTES = [
  {
    path: '',
    redirectTo: 'blocks',
    pathMatch: 'full'
  },
  {
    path: 'blocks',
    component: BlockChainTableComponent
  },
  {
    path: 'block/:id',
    component: BlockDetailsComponent
  },
  {
    path:'address/:address',
    component: AddressDetailComponent
  }
  ,
  {
    path:'transaction/:txid',
    component: TransactionDetailComponent
  }
];

@NgModule({
  declarations: [
    AppComponent,
    SkycoinHeaderComponent,
    BlockChainTableComponent,
    SkycoinPaginationComponent,
    SkycoinSearchBarComponent,
    FooterComponent,
    NumPagesPipe,
    BlockDetailsComponent,
    AddressDetailComponent,
    TransactionDetailComponent,
    LoadingComponent,
    BlockChainCoinSupplyComponent,

  ],
  imports: [
    CommonModule,
    BrowserModule,
    FormsModule,
    HttpModule,
    MomentModule,
    RouterModule.forRoot(ROUTES)
  ],
  providers: [BlockChainService,SkycoinBlockchainPaginationService,UxOutputsService, TransactionDetailService, CoinSupplyService],
  bootstrap: [AppComponent]
})
export class AppModule { }
