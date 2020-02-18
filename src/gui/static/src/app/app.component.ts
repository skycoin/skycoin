import { Component, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { AppService } from './services/app.service';
import { HwWalletService } from './services/hw-wallet.service';
import { HwPinDialogComponent } from './components/layout/hardware-wallet/hw-pin-dialog/hw-pin-dialog.component';
import { Bip39WordListService } from './services/bip39-word-list.service';
import { HwConfirmTxDialogComponent } from './components/layout/hardware-wallet/hw-confirm-tx-dialog/hw-confirm-tx-dialog.component';
import { HwWalletPinService } from './services/hw-wallet-pin.service';
import { HwWalletSeedWordService } from './services/hw-wallet-seed-word.service';
import { LanguageService } from './services/language.service';
import { MsgBarComponent } from './components/layout/msg-bar/msg-bar.component';
import { MsgBarService } from './services/msg-bar.service';
import { SeedWordDialogComponent } from './components/layout/seed-word-dialog/seed-word-dialog.component';
import { SelectLanguageComponent } from './components/layout/select-language/select-language.component';
import { WalletsAndAddressesService } from './services/wallet-operations/wallets-and-addresses.service';
import { redirectToErrorPage } from './utils/errors';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  @ViewChild('msgBar', { static: false }) msgBar: MsgBarComponent;

  constructor(
    private appService: AppService,
    private languageService: LanguageService,
    hwWalletService: HwWalletService,
    hwWalletPinService: HwWalletPinService,
    hwWalletSeedWordService: HwWalletSeedWordService,
    walletsAndAddressesService: WalletsAndAddressesService,
    private bip38WordList: Bip39WordListService,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
  ) {
    hwWalletPinService.requestPinComponent = HwPinDialogComponent;
    hwWalletSeedWordService.requestWordComponent = SeedWordDialogComponent;
    hwWalletService.signTransactionConfirmationComponent = HwConfirmTxDialogComponent;

    walletsAndAddressesService.errorDuringinitialLoad.subscribe(failed => {
      if (failed) {
        // The error page will show error number 2.
        redirectToErrorPage(2);
      }
    });
  }

  ngOnInit() {
    this.appService.UpdateData();
    this.languageService.initialize();

    const subscription = this.languageService.savedSelectedLanguageLoaded.subscribe(savedSelectedLanguageLoaded => {
      if (!savedSelectedLanguageLoaded) {
        SelectLanguageComponent.openDialog(this.dialog, true);
      }

      subscription.unsubscribe();
    });

    setTimeout(() => {
      this.msgBarService.msgBarComponent = this.msgBar;
    });
  }
}
