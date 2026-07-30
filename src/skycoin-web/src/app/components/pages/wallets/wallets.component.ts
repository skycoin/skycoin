import { Component, OnInit, OnDestroy, ChangeDetectionStrategy } from '@angular/core';
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
import { HwWalletService } from '../../../services/hw-wallet.service';
import { MsgBarService } from '../../../services/msg-bar.service';
import { environment } from '../../../../environments/environment';

@Component({
    selector: 'app-wallets',
    templateUrl: './wallets.component.html',
    styleUrls: ['./wallets.component.scss'],
    changeDetection: ChangeDetectionStrategy.Eager,
    standalone: false
})
export class WalletsComponent implements OnInit, OnDestroy {

  wallets!: Wallet[];
  currentCoin!: BaseCoin;
  showLockIcons: boolean;

  private subscriptionsGroup: Subscription[] = [];
  private confirmSeedSubscription!: Subscription;
  private deleteWalletSubscription!: Subscription;

  addingHwWallet = false;

  constructor(
    private walletService: WalletService,
    private dialog: CustomMatDialogService,
    private coinService: CoinService,
    private translateService: TranslateService,
    private hwWalletService: HwWalletService,
    private msgBarService: MsgBarService,
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

  unlockWallet(event: any, wallet: Wallet) {
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

  addHardwareWallet() {
    if (this.addingHwWallet) {
      return;
    }

    this.addingHwWallet = true;
    this.msgBarService.hide();

    this.hwWalletService.getAddresses(1, 0).subscribe(
      (result) => {
        const addresses = result.rawResponse;
        const firstAddress = Array.isArray(addresses) ? addresses[0] : addresses;

        // Check if this HW wallet is already added
        if (this.wallets && this.wallets.some(w => w.isHardware && w.addresses.length > 0 && w.addresses[0].address === firstAddress)) {
          this.addingHwWallet = false;
          this.msgBarService.showError('hardware-wallet.errors.already-added');
          return;
        }

        // Get device features for the label
        this.hwWalletService.getFeatures().subscribe(
          (features) => {
            const label = features.rawResponse.label || 'Hardware Wallet';

            const wallet: Wallet = {
              label: label,
              addresses: [{ address: firstAddress }],
              isHardware: true,
              coinId: this.currentCoin ? this.currentCoin.id : 1,
            };

            this.walletService.add(wallet);
            this.addingHwWallet = false;
            this.msgBarService.showDone('hardware-wallet.added');
          },
          () => {
            // If we can't get features, use default label
            const wallet: Wallet = {
              label: 'Hardware Wallet',
              addresses: [{ address: firstAddress }],
              isHardware: true,
              coinId: this.currentCoin ? this.currentCoin.id : 1,
            };

            this.walletService.add(wallet);
            this.addingHwWallet = false;
            this.msgBarService.showDone('hardware-wallet.added');
          }
        );
      },
      (error) => {
        this.addingHwWallet = false;
        this.msgBarService.showError(error.translatableErrorMsg || error.message || 'hardware-wallet.errors.generic-error');
      }
    );
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
