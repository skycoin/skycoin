import { Component, OnDestroy, OnInit } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { CreateWalletComponent } from './create-wallet/create-wallet.component';
import { Wallet, ConfirmationData } from '../../../app.datatypes';
import { HwWalletOptionsComponent } from '../../layout/hardware-wallet/hw-options/hw-options.component';
import { ISubscription } from 'rxjs/Subscription';
import { Router } from '@angular/router';
import { HwWalletService } from '../../../services/hw-wallet.service';
import { showConfirmationModal } from '../../../utils';

@Component({
  selector: 'app-wallets',
  templateUrl: './wallets.component.html',
  styleUrls: ['./wallets.component.scss'],
})
export class WalletsComponent implements OnInit, OnDestroy {

  hwCompatibilityActivated = false;

  wallets: Wallet[] = [];
  hardwareWallets: Wallet[] = [];

  private subscription: ISubscription;

  constructor(
    public walletService: WalletService,
    public hwWalletService: HwWalletService,
    private dialog: MatDialog,
    private router: Router,
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

  ngOnInit(): void {
    if (this.hwWalletService.showOptionsWhenPossible) {
      this.hwWalletService.showOptionsWhenPossible = false;
      this.adminHwWallet();
    }
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
    config.autoFocus = false;
    this.dialog.open(HwWalletOptionsComponent, config).afterClosed().subscribe(() => {
      this.walletService.all().first().subscribe(wallets => {
        if (wallets.length === 0) {
          setTimeout(() => this.router.navigate(['/wizard']), 500);
        }
      });
    });
  }

  toggleWallet(wallet: Wallet) {
    if (wallet.isHardware && wallet.hasHwSecurityWarnings && !wallet.opened) {
      const confirmationData: ConfirmationData = {
        headerText: 'hardware-wallet.security-warning.title',
        text: 'hardware-wallet.security-warning.text',
        checkboxText: 'hardware-wallet.security-warning.check',
        confirmButtonText: 'hardware-wallet.security-warning.continue',
        cancelButtonText: 'hardware-wallet.security-warning.cancel',
      };

      showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          wallet.hasHwSecurityWarnings = false;
          this.walletService.saveHardwareWallets();
          wallet.opened = true;
        }
      });
    } else {
      wallet.opened = !wallet.opened;
    }
  }
}
