import { Component, OnDestroy } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { CreateWalletComponent } from './create-wallet/create-wallet.component';
import { Wallet } from '../../../app.datatypes';
import { HwWalletOptionsComponent } from '../../layout/hardware-wallet/hw-options/hw-options';
import { ISubscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-wallets',
  templateUrl: './wallets.component.html',
  styleUrls: ['./wallets.component.scss'],
})
export class WalletsComponent implements OnDestroy {

  hwCompatibilityActivated = false;

  wallets: Wallet[] = [];
  hardwareWallets: Wallet[] = [];

  private subscription: ISubscription;

  constructor(
    public walletService: WalletService,
    private dialog: MatDialog,
  ) {
    if (window['isElectron']) {
      this.hwCompatibilityActivated = window['ipcRenderer'].sendSync('hwCompatibilityActivated');
    }

    this.subscription = this.walletService.all().subscribe(wallets => {
      this.wallets = [];
      this.hardwareWallets = [];
      wallets.forEach(value => {
        if (!value.isHardware) {
          this.wallets.push(value);
        } else {
          this.hardwareWallets.push(value);
        }
      });
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  addWallet(create) {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.data = { create };
    this.dialog.open(CreateWalletComponent, config);
  }

  adminHwWallet() {
    const config = new MatDialogConfig();
    config.width = '566px';
    this.dialog.open(HwWalletOptionsComponent, config);
  }

  toggleWallet(wallet: Wallet) {
    wallet.opened = !wallet.opened;
  }
}
