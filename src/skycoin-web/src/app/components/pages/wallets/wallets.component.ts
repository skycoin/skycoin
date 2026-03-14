import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatDialogConfig } from '@angular/material/dialog';
import { Subscription, first } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';

import { Wallet } from '../../../app.datatypes';
import { WalletService } from '../../../services/wallet/wallet.service';
import { CreateWalletComponent } from './create-wallet/create-wallet.component';
import { openUnlockWalletModal, openDeleteWalletModal } from '../../../utils/index';
import { CoinService } from '../../../services/coin.service';
import { BaseCoin } from '../../../coins/basecoin';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';
import { environment } from '../../../../environments/environment';

@Component({
    selector: 'app-wallets',
    templateUrl: './wallets.component.html',
    styleUrls: ['./wallets.component.scss'],
    standalone: false
})
export class WalletsComponent implements OnInit, OnDestroy {

  wallets: Wallet[];
  currentCoin: BaseCoin;
  showLockIcons: boolean;

  private subscriptionsGroup: Subscription[] = [];
  private confirmSeedSubscription: Subscription;
  private deleteWalletSubscription: Subscription;

  constructor(
    private walletService: WalletService,
    private dialog: CustomMatDialogService,
    private coinService: CoinService,
    private translateService: TranslateService
  ) {
    this.showLockIcons = !environment.production;
  }

  ngOnInit() {
    this.subscriptionsGroup.push(this.walletService.currentWallets.subscribe( (wallets) => {
      this.wallets = wallets;
    }));

    this.subscriptionsGroup.push(this.coinService.currentCoin
      .subscribe((coin: BaseCoin) => this.currentCoin = coin)
    );
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
    this.removeConfirmationSuscriptions();
  }

  addWallet(create: boolean) {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.data = { create };
    config.autoFocus = false;
    this.dialog.open(CreateWalletComponent, config);
  }

  unlockWallet(event, wallet: Wallet) {
    if (!wallet.needSeedConfirmation) {
      event.stopPropagation();
      openUnlockWalletModal(wallet, this.dialog);
    }
  }

  toggleWallet(wallet: Wallet) {
    if (wallet.needSeedConfirmation) {
      this.removeConfirmationSuscriptions();

      const unlockDialog = openUnlockWalletModal({wallet: wallet}, this.dialog, false).componentInstance;

      this.confirmSeedSubscription = unlockDialog.onWalletUnlocked.pipe(first()).subscribe(() => {
        wallet.needSeedConfirmation = false;
        this.walletService.saveWallets();
        wallet.opened ? wallet.opened = false : wallet.opened = true;
      });

      this.deleteWalletSubscription = unlockDialog.onDeleteClicked.pipe(first()).subscribe(() => {
        openDeleteWalletModal(this.dialog, wallet, this.translateService, this.walletService);
      });
    } else {
      wallet.opened ? wallet.opened = false : wallet.opened = true;
    }
  }

  private removeConfirmationSuscriptions() {
    if (this.confirmSeedSubscription) {
      this.confirmSeedSubscription.unsubscribe();
    }
    if (this.deleteWalletSubscription) {
      this.deleteWalletSubscription.unsubscribe();
    }
  }
}
