import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import {
  MdButtonModule, MdCardModule, MdDialogModule, MdExpansionModule, MdGridListModule, MdIconModule, MdInputModule,
  MdListModule, MdMenuModule, MdProgressBarModule, MdProgressSpinnerModule,
  MdSelectModule, MdSnackBarModule, MdTabsModule, MdToolbarModule, MdTooltipModule
} from '@angular/material';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { AppComponent } from './app.component';
import { HttpModule } from '@angular/http';
import { ApiService } from './services/api.service';
import { WalletService } from './services/wallet.service';
import { WalletsComponent } from './components/pages/wallets/wallets.component';
import { WalletDetailComponent } from './components/pages/wallets/address-detail/wallet-detail.component';
import { CreateWalletComponent } from './components/pages/wallets/create-wallet/create-wallet.component';
import { ReactiveFormsModule } from '@angular/forms';
import { SkyPipe } from './pipes/sky.pipe';
import { SendSkycoinComponent } from './components/pages/send-skycoin/send-skycoin.component';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
import { HistoryComponent } from './components/pages/history/history.component';
import { DateFromNowPipe } from './pipes/date-from-now.pipe';
import { RouterModule } from '@angular/router';
import { BreadcrumbComponent } from './components/layout/breadcrumb/breadcrumb.component';
import { TransactionComponent } from './components/pages/transaction/transaction.component';
import { BackButtonComponent } from './components/layout/back-button/back-button.component';
import { ExplorerComponent } from './components/pages/explorer/explorer.component';
import { BlockchainService } from './services/blockchain.service';
import { DateTimePipe } from './pipes/date-time.pipe';
import { TransactionsAmountPipe } from './pipes/transactions-amount.pipe';
import { BlockComponent } from './components/pages/block/block.component';
import { AddressComponent } from './components/pages/address/address.component';
import { PendingTransactionsComponent } from './components/pages/settings/pending-transactions/pending-transactions.component';
import { OutputsComponent } from './components/pages/settings/outputs/outputs.component';
import { BlockchainComponent } from './components/pages/settings/blockchain/blockchain.component';
import { BackupComponent } from './components/pages/settings/backup/backup.component';
import { NetworkComponent } from './components/pages/settings/network/network.component';
import { NetworkService } from './services/network.service';
import { ChangeNameComponent } from './components/pages/wallets/change-name/change-name.component';
import { ButtonComponent } from './components/layout/button/button.component';
import { QrCodeComponent } from './components/layout/qr-code/qr-code.component';
import { BuyComponent } from './components/pages/buy/buy.component';
import { AddDepositAddressComponent } from './components/pages/buy/add-deposit-address/add-deposit-address.component';
import { PurchaseService } from './services/purchase.service';
import { TellerStatusPipe } from './pipes/teller-status.pipe';

const ROUTES = [
  {
    path: '',
    redirectTo: 'wallets',
    pathMatch: 'full'
  },
  {
    path: 'wallets',
    component: WalletsComponent,
    data: {
      breadcrumb: 'Wallets',
    },
  },
  {
    path: 'send',
    component: SendSkycoinComponent,
    data: {
      breadcrumb: 'Send Skycoin',
    },
  },
  {
    path: 'history',
    children: [
      {
        path: '',
        component: HistoryComponent,
        data: {
          breadcrumb: 'History',
        },
      },
      {
        path: ':transaction',
        component: TransactionComponent,
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
        component: ExplorerComponent,
        data: {
          breadcrumb: 'Explorer',
        },
      },
      {
        path: 'address/:address',
        component: AddressComponent,
        data: {
          breadcrumb: 'Address',
        },
      },
      {
        path: ':block',
        component: BlockComponent,
        data: {
          breadcrumb: 'Block',
        },
      },
      {
        path: 'transaction/:transaction',
        component: TransactionComponent,
        data: {
          breadcrumb: 'Transaction',
        },
      },
    ],
  },
  {
    path: 'buy',
    component: BuyComponent,
    data: {
      breadcrumb: 'Buy Skycoin',
    },
  },
  {
    path: 'settings',
    children: [
      {
        path: 'backup',
        component: BackupComponent,
        data: {
          breadcrumb: 'Backup',
        },
      },
      {
        path: 'blockchain',
        component: BlockchainComponent,
        data: {
          breadcrumb: 'Blockchain',
        },
      },
      {
        path: 'network',
        component: NetworkComponent,
        data: {
          breadcrumb: 'Networking',
        },
      },
      {
        path: 'outputs',
        component: OutputsComponent,
        data: {
          breadcrumb: 'Outputs',
        },
      },
      {
        path: 'pending-transactions',
        component: PendingTransactionsComponent,
        data: {
          breadcrumb: 'Pending transactions',
        },
      },
    ],
  },
];

@NgModule({
  declarations: [
    AppComponent,
    HistoryComponent,
    WalletsComponent,
    WalletDetailComponent,
    CreateWalletComponent,
    SkyPipe,
    SendSkycoinComponent,
    DateFromNowPipe,
    BreadcrumbComponent,
    TransactionComponent,
    BackButtonComponent,
    ExplorerComponent,
    DateTimePipe,
    TransactionsAmountPipe,
    BlockComponent,
    AddressComponent,
    PendingTransactionsComponent,
    OutputsComponent,
    BlockchainComponent,
    BackupComponent,
    NetworkComponent,
    ChangeNameComponent,
    ButtonComponent,
    QrCodeComponent,
    BuyComponent,
    AddDepositAddressComponent,
    TellerStatusPipe,
  ],
  entryComponents: [
    AddDepositAddressComponent,
    CreateWalletComponent,
    ChangeNameComponent,
    QrCodeComponent,
  ],
  imports: [
    BrowserModule,
    HttpModule,
    MdButtonModule,
    MdCardModule,
    MdDialogModule,
    MdExpansionModule,
    MdGridListModule,
    MdIconModule,
    MdInputModule,
    MdListModule,
    MdMenuModule,
    MdProgressBarModule,
    MdProgressSpinnerModule,
    MdSelectModule,
    MdSnackBarModule,
    MdTabsModule,
    MdToolbarModule,
    MdTooltipModule,
    NgxDatatableModule,
    NoopAnimationsModule,
    ReactiveFormsModule,
    RouterModule.forRoot(ROUTES, { useHash: true }),
  ],
  providers: [
    ApiService,
    BlockchainService,
    NetworkService,
    PurchaseService,
    WalletService,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
