import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { RouterModule, RouteReuseStrategy, Routes } from '@angular/router';

import { AppComponent } from './app.component';
import { FooterComponent } from './components/layout/footer/footer.component';
import { HeaderComponent } from './components/layout/header/header.component';
import { SearchBarComponent } from './components/layout/search-bar/search-bar.component';
import { LoadingComponent } from './components/layout/loading/loading.component';
import { BlocksComponent } from './components/pages/blocks/blocks.component';
import { UnconfirmedTransactionsComponent } from './components/pages/unconfirmed-transactions/unconfirmed-transactions.component';
import { ApiService } from './services/api/api.service';
import { HttpClientModule } from '@angular/common/http';
import { BlockDetailsComponent } from './components/pages/block-details/block-details.component';
import { TransactionDetailComponent } from './components/pages/transaction-detail/transaction-detail.component';
import { AddressDetailComponent } from './components/pages/address-detail/address-detail.component';
import { ExplorerService } from './services/explorer/explorer.service';
import { QrCodeComponent } from './components/layout/qr-code/qr-code.component';
import { RichlistComponent } from 'app/components/pages/richlist/richlist.component';
import { UnspentOutputsComponent } from 'app/components/pages/unspent-outputs/unspent-outputs.component';
import { CopyButtonComponent } from 'app/components/layout/copy-button/copy-button.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { AppReuseStrategy } from 'app/app.reuse-strategy';
import { TranslateModule, TranslateLoader } from '@ngx-translate/core';
import { AppTranslateLoader } from 'app/app.translate-loader';
import { GenericHeaderComponent } from 'app/components/layout/generic-header/generic-header.component';
import { GenericFooterComponent } from 'app/components/layout/generic-footer/generic-footer.component';
import { DatePipe, DecimalPipe } from '@angular/common';
import { DateFormatterComponent } from 'app/components/layout/date-formatter/date-formatter.component';
import { SearchService } from './services/search/search.service';
import { SearchComponent } from './components/pages/search/search.component';
import { TransactionInfoComponent } from './components/layout/transaction-info/transaction-info.component';
import { LanguageService } from 'app/services/language/language.service';
import { LanguageSelectionComponent } from 'app/components/layout/language-selection/language-selection.component';
import { AmountPipe } from 'app/pipes/amount.pipe';
import { NodeUrlComponent } from './components/pages/node-url.component';


const ROUTES: Routes = [
  {
    path: '',
    redirectTo: 'app/blocks/1',
    pathMatch: 'full'
  },
  {
    path: 'app/blocks',
    component: BlocksComponent
  },
  {
    path: 'app/blocks/:page',
    component: BlocksComponent
  },
  {
    path: 'app/block/:id',
    component: BlockDetailsComponent
  },
  {
    path: 'app/address/:address',
    redirectTo: 'app/address/:address/1',
    pathMatch: 'full'
  },
  {
    path: 'app/address/:address/:page',
    component: AddressDetailComponent
  },
  {
    path: 'app/transaction/:txid',
    component: TransactionDetailComponent
  },
  {
    path: 'app/unconfirmedtransactions',
    component: UnconfirmedTransactionsComponent
  },
  {
    path: 'app/richlist',
    component: RichlistComponent
  },
  {
    path: 'app/unspent/:address',
    component: UnspentOutputsComponent
  },
  {
    path: 'app/search',
    redirectTo: 'app/blocks/1',
    pathMatch: 'full'
  },
  {
    path: 'app/search/:term',
    component: SearchComponent
  },
  {
    path: 'node',
    redirectTo: 'node/null',
    pathMatch: 'full'
  },
  {
    path: 'node/:url',
    component: NodeUrlComponent
  },
];

@NgModule({
  declarations: [
    AddressDetailComponent,
    AppComponent,
    BlockDetailsComponent,
    BlocksComponent,
    UnconfirmedTransactionsComponent,
    FooterComponent,
    HeaderComponent,
    GenericFooterComponent,
    GenericHeaderComponent,
    LoadingComponent,
    QrCodeComponent,
    SearchBarComponent,
    TransactionDetailComponent,
    RichlistComponent,
    UnspentOutputsComponent,
    CopyButtonComponent,
    DateFormatterComponent,
    SearchComponent,
    TransactionInfoComponent,
    LanguageSelectionComponent,
    AmountPipe,
    NodeUrlComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    RouterModule.forRoot(ROUTES, {}),
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useClass: AppTranslateLoader
      }
    })
  ],
  providers: [
    ApiService,
    ExplorerService,
    SearchService,
    LanguageService,
    {provide: RouteReuseStrategy, useClass: AppReuseStrategy},
    DatePipe,
    DecimalPipe,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
