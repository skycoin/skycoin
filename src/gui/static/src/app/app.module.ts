import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { AppComponent } from './app.component';
import { HttpModule } from '@angular/http';
import { ApiService } from './services/api.service';
import { WalletService } from './services/wallet.service';
import { WalletsComponent } from './components/pages/wallets/wallets.component';
import { CreateWalletComponent } from './components/pages/wallets/create-wallet/create-wallet.component';
import { ReactiveFormsModule } from '@angular/forms';
import { SendSkycoinComponent } from './components/pages/send-skycoin/send-skycoin.component';
import { DateFromNowPipe } from './pipes/date-from-now.pipe';
import { RouterModule } from '@angular/router';
import { BlockchainService } from './services/blockchain.service';
import { DateTimePipe } from './pipes/date-time.pipe';
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
import { HeaderComponent } from './components/layout/header/header.component';
import { TopBarComponent } from './components/layout/header/top-bar/top-bar.component';
import { PriceService } from './services/price.service';
import { TransactionListComponent } from './components/pages/transaction-list/transaction-list.component';
import { TransactionDetailComponent } from './components/pages/transaction-list/transaction-detail/transaction-detail.component';
import { NavBarComponent } from './components/layout/header/nav-bar/nav-bar.component';
import { WalletDetailComponent } from './components/pages/wallets/wallet-detail/wallet-detail.component';
import { ModalComponent } from './components/layout/modal/modal.component';
import { PasswordDialogComponent } from './components/layout/password-dialog/password-dialog.component';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatDialogModule } from '@angular/material/dialog';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatListModule } from '@angular/material/list';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTabsModule } from '@angular/material/tabs';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatSliderModule } from '@angular/material';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule } from '@angular/common/http';
import { AppService } from './services/app.service';
import { WizardGuardService } from './services/wizard-guard.service';
import { OnboardingCreateWalletComponent } from './components/pages/onboarding/onboarding-create-wallet/onboarding-create-wallet.component';
import { OnboardingEncryptWalletComponent } from './components/pages/onboarding/onboarding-encrypt-wallet/onboarding-encrypt-wallet.component';
import { OnboardingSafeguardComponent } from './components/pages/onboarding/onboarding-create-wallet/onboarding-safeguard/onboarding-safeguard.component';
import { DoubleButtonComponent } from './components/layout/double-button/double-button.component';
import { SeedModalComponent } from './components/pages/settings/backup/seed-modal/seed-modal.component';
import { OnboardingComponent } from './components/pages/onboarding/onboarding.component';
import { DontsavepasswordDirective } from './directives/dontsavepassword.directive';
import { SendFormComponent } from './components/pages/send-skycoin/send-form/send-form.component';
import { SendVerifyComponent } from './components/pages/send-skycoin/send-preview/send-preview.component';
import { TransactionInfoComponent } from './components/pages/send-skycoin/send-preview/transaction-info/transaction-info.component';
import { SendFormAdvancedComponent } from './components/pages/send-skycoin/send-form-advanced/send-form-advanced.component';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { AppTranslateLoader } from './app.translate-loader';
import { NavBarService } from './services/nav-bar.service';


const ROUTES = [
  {
    path: '',
    redirectTo: 'wallets',
    pathMatch: 'full',
  },
  {
    path: 'wallets',
    component: WalletsComponent,
    canActivate: [WizardGuardService],
  },
  {
    path: 'send',
    component: SendSkycoinComponent,
    canActivate: [WizardGuardService],
  },
  {
    path: 'transactions',
    component: TransactionListComponent,
    canActivate: [WizardGuardService],
  },
  {
    path: 'buy',
    component: BuyComponent,
    canActivate: [WizardGuardService],
  },
  {
    path: 'settings',
    children: [
      {
        path: 'backup',
        component: BackupComponent,
      },
      {
        path: 'blockchain',
        component: BlockchainComponent,
      },
      {
        path: 'network',
        component: NetworkComponent,
      },
      {
        path: 'outputs',
        component: OutputsComponent,
      },
      {
        path: 'pending-transactions',
        component: PendingTransactionsComponent,
      },
    ],
    canActivate: [WizardGuardService],
  },
  {
    path: 'wizard',
    component: OnboardingComponent,
  },
];

@NgModule({
  declarations: [
    AddDepositAddressComponent,
    AppComponent,
    BackupComponent,
    BlockchainComponent,
    BuyComponent,
    ButtonComponent,
    ChangeNameComponent,
    CreateWalletComponent,
    DateFromNowPipe,
    DateTimePipe,
    HeaderComponent,
    NetworkComponent,
    OutputsComponent,
    PendingTransactionsComponent,
    QrCodeComponent,
    SendSkycoinComponent,
    TellerStatusPipe,
    TopBarComponent,
    TransactionDetailComponent,
    TransactionListComponent,
    WalletsComponent,
    NavBarComponent,
    WalletDetailComponent,
    ModalComponent,
    OnboardingCreateWalletComponent,
    OnboardingEncryptWalletComponent,
    OnboardingSafeguardComponent,
    DoubleButtonComponent,
    PasswordDialogComponent,
    SeedModalComponent,
    OnboardingComponent,
    DontsavepasswordDirective,
    SendFormComponent,
    SendVerifyComponent,
    TransactionInfoComponent,
    SendFormAdvancedComponent,
  ],
  entryComponents: [
    AddDepositAddressComponent,
    CreateWalletComponent,
    ChangeNameComponent,
    QrCodeComponent,
    SendSkycoinComponent,
    TransactionDetailComponent,
    OnboardingSafeguardComponent,
    PasswordDialogComponent,
    SeedModalComponent,
  ],
  imports: [
    BrowserModule,
    HttpModule,
    HttpClientModule,
    MatButtonModule,
    MatCardModule,
    MatDialogModule,
    MatExpansionModule,
    MatGridListModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatMenuModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    MatSelectModule,
    MatSnackBarModule,
    MatTabsModule,
    MatToolbarModule,
    MatTooltipModule,
    MatCheckboxModule,
    MatSliderModule,
    NoopAnimationsModule,
    ReactiveFormsModule,
    RouterModule.forRoot(ROUTES, { useHash: true }),
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useClass: AppTranslateLoader,
      },
    }),
  ],
  providers: [
    ApiService,
    AppService,
    BlockchainService,
    NavBarService,
    NetworkService,
    PriceService,
    PurchaseService,
    WalletService,
    WizardGuardService,
  ],
  bootstrap: [AppComponent],
})
export class AppModule { }
