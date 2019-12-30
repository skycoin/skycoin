import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { AppComponent } from './app.component';
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
import { MatProgressSpinnerModule, MAT_PROGRESS_SPINNER_DEFAULT_OPTIONS } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTabsModule } from '@angular/material/tabs';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatSliderModule } from '@angular/material/slider';
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
import { SendVerifyComponent } from './components/pages/send-skycoin/send-preview/send-preview.component';
import { TransactionInfoComponent } from './components/pages/send-skycoin/send-preview/transaction-info/transaction-info.component';
import { SendCoinsFormComponent } from './components/pages/send-skycoin/send-coins-form/send-coins-form.component';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { AppTranslateLoader } from './app.translate-loader';
import { NavBarService } from './services/nav-bar.service';
import { LoadingContentComponent } from './components/layout/loading-content/loading-content.component';
import { NumberOfAddressesComponent } from './components/pages/wallets/number-of-addresses/number-of-addresses';
import { SelectAddressComponent } from './components/layout/select-address/select-address';
import { CreateWalletFormComponent } from './components/pages/wallets/create-wallet/create-wallet-form/create-wallet-form.component';
import { ResetPasswordComponent } from './components/pages/reset-password/reset-password.component';
import { ExchangeComponent } from './components/pages/exchange/exchange.component';
import { ExchangeService } from './services/exchange.service';
import { ExchangeCreateComponent } from './components/pages/exchange/exchange-create/exchange-create.component';
import { ExchangeStatusComponent } from './components/pages/exchange/exchange-status/exchange-status.component';
import { HwWalletService } from './services/hw-wallet.service';
import { HwOptionsDialogComponent } from './components/layout/hardware-wallet/hw-options-dialog/hw-options-dialog.component';
import { HwWipeDialogComponent } from './components/layout/hardware-wallet/hw-wipe-dialog/hw-wipe-dialog.component';
import { HwAddedDialogComponent } from './components/layout/hardware-wallet/hw-added-dialog/hw-added-dialog.component';
import { HwGenerateSeedDialogComponent } from './components/layout/hardware-wallet/hw-generate-seed-dialog/hw-generate-seed-dialog.component';
import { HwBackupDialogComponent } from './components/layout/hardware-wallet/hw-backup-dialog/hw-backup-dialog.component';
import { ConfirmationComponent } from './components/layout/confirmation/confirmation.component';
import { HwMessageComponent } from './components/layout/hardware-wallet/hw-message/hw-message.component';
import { HwPinDialogComponent } from './components/layout/hardware-wallet/hw-pin-dialog/hw-pin-dialog.component';
import { HwChangePinDialogComponent } from './components/layout/hardware-wallet/hw-change-pin-dialog/hw-change-pin-dialog.component';
import { HwPinHelpDialogComponent } from './components/layout/hardware-wallet/hw-pin-help-dialog/hw-pin-help-dialog.component';
import { HwRestoreSeedDialogComponent } from './components/layout/hardware-wallet/hw-restore-seed-dialog/hw-restore-seed-dialog.component';
import { Bip39WordListService } from './services/bip39-word-list.service';
import { HwDialogBaseComponent } from './components/layout/hardware-wallet/hw-dialog-base.component';
import { HwConfirmTxDialogComponent } from './components/layout/hardware-wallet/hw-confirm-tx-dialog/hw-confirm-tx-dialog.component';
import { HwConfirmAddressDialogComponent } from './components/layout/hardware-wallet/hw-confirm-address-dialog/hw-confirm-address-dialog.component';
import { HwWalletDaemonService } from './services/hw-wallet-daemon.service';
import { HwWalletPinService } from './services/hw-wallet-pin.service';
import { HwWalletSeedWordService } from './services/hw-wallet-seed-word.service';
import { LanguageService } from './services/language.service';
import { SelectLanguageComponent } from './components/layout/select-language/select-language.component';
import { ExchangeHistoryComponent } from './components/pages/exchange/exchange-history/exchange-history.component';
import { StorageService } from './services/storage.service';
import { CommonTextPipe } from './pipes/common-text.pipe';
import { AmountPipe } from './pipes/amount.pipe';
import { DecimalPipe } from '@angular/common';
import { HwRemovePinDialogComponent } from './components/layout/hardware-wallet/hw-remove-pin-dialog/hw-remove-pin-dialog.component';
import { HwUpdateFirmwareDialogComponent } from './components/layout/hardware-wallet/hw-update-firmware-dialog/hw-update-firmware-dialog.component';
import { HwUpdateAlertDialogComponent } from './components/layout/hardware-wallet/hw-update-alert-dialog/hw-update-alert-dialog.component';
import { ChangeNoteComponent } from './components/pages/send-skycoin/send-preview/transaction-info/change-note/change-note.component';
import { MsgBarComponent } from './components/layout/msg-bar/msg-bar.component';
import { MsgBarService } from './services/msg-bar.service';
import { SeedWordDialogComponent } from './components/layout/seed-word-dialog/seed-word-dialog.component';
import { MultipleDestinationsDialogComponent } from './components/layout/multiple-destinations-dialog/multiple-destinations-dialog.component';
import { FormSourceSelectionComponent } from './components/pages/send-skycoin/form-parts/form-source-selection/form-source-selection.component';
import { FormDestinationComponent } from './components/pages/send-skycoin/form-parts/form-destination/form-destination.component';
import { CopyRawTxComponent } from './components/pages/send-skycoin/offline-dialogs/implementations/copy-raw-tx.component';
import { SignRawTxComponent } from './components/pages/send-skycoin/offline-dialogs/implementations/sign-raw-tx.component';
import { BroadcastRawTxComponent } from './components/pages/send-skycoin/offline-dialogs/implementations/broadcast-raw-tx.component';


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
    path: 'exchange',
    component: ExchangeComponent,
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
  {
    path: 'reset/:id',
    component: ResetPasswordComponent,
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
    SendVerifyComponent,
    TransactionInfoComponent,
    SendCoinsFormComponent,
    LoadingContentComponent,
    NumberOfAddressesComponent,
    SelectAddressComponent,
    CreateWalletFormComponent,
    ResetPasswordComponent,
    ExchangeComponent,
    ExchangeCreateComponent,
    ExchangeStatusComponent,
    HwOptionsDialogComponent,
    HwWipeDialogComponent,
    HwAddedDialogComponent,
    HwGenerateSeedDialogComponent,
    HwBackupDialogComponent,
    ConfirmationComponent,
    HwMessageComponent,
    HwPinDialogComponent,
    HwChangePinDialogComponent,
    HwPinHelpDialogComponent,
    HwRestoreSeedDialogComponent,
    HwDialogBaseComponent,
    HwConfirmTxDialogComponent,
    HwConfirmAddressDialogComponent,
    SelectLanguageComponent,
    ExchangeHistoryComponent,
    CommonTextPipe,
    AmountPipe,
    HwRemovePinDialogComponent,
    HwUpdateFirmwareDialogComponent,
    HwUpdateAlertDialogComponent,
    ChangeNoteComponent,
    MsgBarComponent,
    SeedWordDialogComponent,
    MultipleDestinationsDialogComponent,
    FormSourceSelectionComponent,
    FormDestinationComponent,
    CopyRawTxComponent,
    SignRawTxComponent,
    BroadcastRawTxComponent,
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
    NumberOfAddressesComponent,
    SelectAddressComponent,
    HwOptionsDialogComponent,
    HwWipeDialogComponent,
    HwAddedDialogComponent,
    HwGenerateSeedDialogComponent,
    HwBackupDialogComponent,
    ConfirmationComponent,
    HwPinDialogComponent,
    HwChangePinDialogComponent,
    HwPinHelpDialogComponent,
    HwRestoreSeedDialogComponent,
    HwConfirmTxDialogComponent,
    HwConfirmAddressDialogComponent,
    SelectLanguageComponent,
    ExchangeHistoryComponent,
    HwRemovePinDialogComponent,
    HwUpdateFirmwareDialogComponent,
    HwUpdateAlertDialogComponent,
    ChangeNoteComponent,
    SeedWordDialogComponent,
    MultipleDestinationsDialogComponent,
    CopyRawTxComponent,
    SignRawTxComponent,
    BroadcastRawTxComponent,
  ],
  imports: [
    BrowserModule,
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
    MatTabsModule,
    MatToolbarModule,
    MatTooltipModule,
    MatCheckboxModule,
    MatSliderModule,
    MatAutocompleteModule,
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
    ExchangeService,
    NavBarService,
    NetworkService,
    PriceService,
    PurchaseService,
    WalletService,
    WizardGuardService,
    HwWalletService,
    Bip39WordListService,
    HwWalletDaemonService,
    HwWalletPinService,
    HwWalletSeedWordService,
    LanguageService,
    StorageService,
    MsgBarService,
    DecimalPipe,
    {
      provide: MAT_PROGRESS_SPINNER_DEFAULT_OPTIONS,
      useValue: {
          _forceAnimations: true,
      },
    },
  ],
  bootstrap: [AppComponent],
})
export class AppModule { }
