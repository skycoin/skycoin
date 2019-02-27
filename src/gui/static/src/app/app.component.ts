import { Component, OnInit } from '@angular/core';
import 'rxjs/add/operator/takeWhile';
import { TranslateService } from '@ngx-translate/core';

import { AppService } from './services/app.service';
import { WalletService } from './services/wallet.service';
import { HwWalletService } from './services/hw-wallet.service';
import { HwPinDialogComponent } from './components/layout/hardware-wallet/hw-pin-dialog/hw-pin-dialog.component';
import { HwSeedWordDialogComponent } from './components/layout/hardware-wallet/hw-seed-word-dialog/hw-seed-word-dialog.component';
import { Bip39WordListService } from './services/bip39-word-list.service';
import { HwConfirmTxDialogComponent } from './components/layout/hardware-wallet/hw-confirm-tx-dialog/hw-confirm-tx-dialog.component';
import { HwPassphraseDialogComponent } from './components/layout/hardware-wallet/hw-passphrase-dialog/hw-passphrase-dialog.component';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  constructor(
    private appService: AppService,
    walletService: WalletService,
    translateService: TranslateService,
    hwWalletService: HwWalletService,
    private bip38WordList: Bip39WordListService,
  ) {
    translateService.setDefaultLang('en');
    translateService.use('en');

    hwWalletService.requestPinComponent = HwPinDialogComponent;
    hwWalletService.requestPassphraseComponent = HwPassphraseDialogComponent;
    hwWalletService.requestWordComponent = HwSeedWordDialogComponent;
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
  }
}
