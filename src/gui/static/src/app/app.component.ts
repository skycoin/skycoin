import { Component, OnInit, ViewChild } from '@angular/core';
import 'rxjs/add/operator/takeWhile';
import { MatDialog } from '@angular/material';

import { AppService } from './services/app.service';
import { WalletService } from './services/wallet.service';
import { HwWalletService } from './services/hw-wallet.service';
import { HwPinDialogComponent } from './components/layout/hardware-wallet/hw-pin-dialog/hw-pin-dialog.component';
import { Bip39WordListService } from './services/bip39-word-list.service';
import { HwConfirmTxDialogComponent } from './components/layout/hardware-wallet/hw-confirm-tx-dialog/hw-confirm-tx-dialog.component';
import { HwWalletPinService } from './services/hw-wallet-pin.service';
import { HwWalletSeedWordService } from './services/hw-wallet-seed-word.service';
import { LanguageService } from './services/language.service';
import { openChangeLanguageModal } from './utils';
import { MsgBarComponent } from './components/layout/msg-bar/msg-bar.component';
import { MsgBarService } from './services/msg-bar.service';
import { SeedWordDialogComponent } from './components/layout/seed-word-dialog/seed-word-dialog.component';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  @ViewChild('msgBar') msgBar: MsgBarComponent;

  constructor(
    private appService: AppService,
    private languageService: LanguageService,
    walletService: WalletService,
    hwWalletService: HwWalletService,
    hwWalletPinService: HwWalletPinService,
    hwWalletSeedWordService: HwWalletSeedWordService,
    private bip38WordList: Bip39WordListService,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
  ) {
    hwWalletPinService.requestPinComponent = HwPinDialogComponent;
    hwWalletSeedWordService.requestWordComponent = SeedWordDialogComponent;
    hwWalletService.signTransactionConfirmationComponent = HwConfirmTxDialogComponent;

    walletService.initialLoadFailed.subscribe(failed => {
      if (failed) {
        // The "?2" part indicates that error number 2 should be displayed.
        window.location.assign('assets/error-alert/index.html?2');
      }
    });
  }

  ngOnInit() {
    this.appService.testBackend();
    this.languageService.loadLanguageSettings();

    const subscription = this.languageService.selectedLanguageLoaded.subscribe(selectedLanguageLoaded => {
      if (!selectedLanguageLoaded) {
        openChangeLanguageModal(this.dialog, true).subscribe(response => {
          if (response) {
            this.languageService.changeLanguage(response);
          }
        });
      }

      subscription.unsubscribe();

      this.msgBarService.msgBarComponent = this.msgBar;
    });
  }
}
