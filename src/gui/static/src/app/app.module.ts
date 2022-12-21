import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { AppComponent } from './app.component';
import { ApiService } from './services/api.service';
import { WalletsComponent } from './components/pages/wallets/wallets.component';
import { CreateWalletComponent } from './components/pages/wallets/create-wallet/create-wallet.component';
import { ReactiveFormsModule } from '@angular/forms';
import { SendSkycoinComponent } from './components/pages/send-skycoin/send-skycoin.component';
import { DateFromNowPipe } from './pipes/date-from-now.pipe';
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
import { MatLegacyButtonModule as MatButtonModule } from '@angular/material/legacy-button';
import { MatLegacyCardModule as MatCardModule } from '@angular/material/legacy-card';
import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyInputModule as MatInputModule } from '@angular/material/legacy-input';
import { MatLegacyListModule as MatListModule } from '@angular/material/legacy-list';
import { MatLegacyMenuModule as MatMenuModule } from '@angular/material/legacy-menu';
import { MatLegacyProgressBarModule as MatProgressBarModule } from '@angular/material/legacy-progress-bar';
import {
  MatLegacyProgressSpinnerModule as MatProgressSpinnerModule, MAT_LEGACY_PROGRESS_SPINNER_DEFAULT_OPTIONS as MAT_PROGRESS_SPINNER_DEFAULT_OPTIONS,
} from '@angular/material/legacy-progress-spinner';
import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatLegacyTabsModule as MatTabsModule } from '@angular/material/legacy-tabs';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { MatLegacyCheckboxModule as MatCheckboxModule } from '@angular/material/legacy-checkbox';
import { MatLegacyAutocompleteModule as MatAutocompleteModule } from '@angular/material/legacy-autocomplete';
import { MatLegacySliderModule as MatSliderModule } from '@angular/material/legacy-slider';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule } from '@angular/common/http';
import { AppService } from './services/app.service';
import { WizardGuardService } from './services/wizard-guard.service';
import { OnboardingCreateWalletComponent } from './components/pages/onboarding/onboarding-create-wallet/onboarding-create-wallet.component';
import { OnboardingEncryptWalletComponent } from './components/pages/onboarding/onboarding-encrypt-wallet/onboarding-encrypt-wallet.component';
import { DoubleButtonComponent } from './components/layout/double-button/double-button.component';
import { SeedModalComponent } from './components/pages/settings/backup/seed-modal/seed-modal.component';
import { OnboardingComponent } from './components/pages/onboarding/onboarding.component';
import { DontsavepasswordDirective } from './directives/dontsavepassword.directive';
import { FormatNumberDirective } from './directives/format-number.directive';
import { SendVerifyComponent } from './components/pages/send-skycoin/send-preview/send-preview.component';
import { TransactionInfoComponent } from './components/pages/send-skycoin/send-preview/transaction-info/transaction-info.component';
import { SendCoinsFormComponent } from './components/pages/send-skycoin/send-coins-form/send-coins-form.component';
import { NavBarSwitchService } from './services/nav-bar-switch.service';
import { LoadingContentComponent } from './components/layout/loading-content/loading-content.component';
import { NumberOfAddressesComponent } from './components/pages/wallets/number-of-addresses/number-of-addresses';
import { SelectAddressComponent } from './components/layout/select-address/select-address.component';
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
import { OfflineTxPreviewComponent } from './components/pages/send-skycoin/offline-dialogs/offline-tx-preview/offline-tx-preview.component';
import { OfflineDialogsBaseComponent } from './components/pages/send-skycoin/offline-dialogs/offline-dialogs-base.component';
import { ArrowLinkComponent } from './components/layout/arrow-link/arrow-link.component';
import { AddressOptionsComponent } from './components/pages/wallets/wallet-detail/address-options/address-options.component';
import { QrCodeButtonComponent } from './components/layout/qr-code-button/qr-code-button.component';
import { WalletsAndAddressesService } from './services/wallet-operations/wallets-and-addresses.service';
import { SoftwareWalletService } from './services/wallet-operations/software-wallet.service';
import { HardwareWalletService } from './services/wallet-operations/hardware-wallet.service';
import { BalanceAndOutputsService } from './services/wallet-operations/balance-and-outputs.service';
import { SpendingService } from './services/wallet-operations/spending.service';
import { HistoryService } from './services/wallet-operations/history.service';
import { AppRoutingModule } from './app-routing.module';
import { AppTranslationModule } from './app-translation.module';
import { FormFieldErrorDirective } from './directives/form-field-error.directive';
import { EnterLinkComponent } from './components/pages/send-skycoin/enter-link/enter-link.component';
import { DestinationToolsComponent } from './components/pages/send-skycoin/form-parts/form-destination/destination-tools/destination-tools.component';
import { ForceSkywalletWipeComponent } from './components/pages/force-skywallet-wipe/force-skywallet-wipe.component';
import { AssistedSeedFieldComponent } from './components/pages/wallets/create-wallet/create-wallet-form/assisted-seed-field/assisted-seed-field.component';

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
    DoubleButtonComponent,
    PasswordDialogComponent,
    SeedModalComponent,
    OnboardingComponent,
    DontsavepasswordDirective,
    FormatNumberDirective,
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
    OfflineTxPreviewComponent,
    OfflineDialogsBaseComponent,
    ArrowLinkComponent,
    AddressOptionsComponent,
    QrCodeButtonComponent,
    FormFieldErrorDirective,
    EnterLinkComponent,
    DestinationToolsComponent,
    ForceSkywalletWipeComponent,
    AssistedSeedFieldComponent,
  ],
  imports: [
    AppTranslationModule,
    AppRoutingModule,
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
  ],
  providers: [
    ApiService,
    AppService,
    BlockchainService,
    ExchangeService,
    NavBarSwitchService,
    NetworkService,
    PriceService,
    PurchaseService,
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
    WalletsAndAddressesService,
    SoftwareWalletService,
    HardwareWalletService,
    BalanceAndOutputsService,
    SpendingService,
    HistoryService,
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
