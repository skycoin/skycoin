import { Component, OnDestroy, OnInit } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { CreateWalletComponent } from './create-wallet/create-wallet.component';
import { Wallet, ConfirmationData } from '../../../app.datatypes';
import { HwOptionsDialogComponent } from '../../layout/hardware-wallet/hw-options-dialog/hw-options-dialog.component';
import { SubscriptionLike } from 'rxjs';
import { Router } from '@angular/router';
import { HwWalletService } from '../../../services/hw-wallet.service';
import { showConfirmationModal } from '../../../utils';
import { first } from 'rxjs/operators';

@Component({
  selector: 'app-wallets',
  templateUrl: './wallets.component.html',
  styleUrls: ['./wallets.component.scss'],
})
export class WalletsComponent implements OnInit, OnDestroy {

  hwCompatibilityActivated = false;

  wallets: Wallet[] = [];
  hardwareWallets: Wallet[] = [];

  private subscription: SubscriptionLike;

  constructor(
    private walletService: WalletService,
    private hwWalletService: HwWalletService,
    private dialog: MatDialog,
    private router: Router,
  ) {
    this.hwCompatibilityActivated = this.hwWalletService.hwWalletCompatibilityActivated;

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
      setTimeout(() => {
        this.hwWalletService.showOptionsWhenPossible = false;
        this.adminHwWallet();
      });
    }
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  addWallet(create) {
    CreateWalletComponent.openDialog(this.dialog, { create });
  }

  adminHwWallet() {
    HwOptionsDialogComponent.openDialog(this.dialog, false).afterClosed().subscribe(() => {
      this.walletService.all().pipe(first()).subscribe(wallets => {
        if (wallets.length === 0) {
          setTimeout(() => this.router.navigate(['/wizard']), 500);
        }
      });
    });
  }

  toggleWallet(wallet: Wallet) {
    if (wallet.isHardware && wallet.hasHwSecurityWarnings && !wallet.stopShowingHwSecurityPopup && !wallet.opened) {
      const confirmationData: ConfirmationData = {
        headerText: 'hardware-wallet.security-warning.title',
        text: 'hardware-wallet.security-warning.text',
        checkboxText: 'hardware-wallet.security-warning.check',
        confirmButtonText: 'hardware-wallet.security-warning.continue',
        cancelButtonText: 'hardware-wallet.security-warning.cancel',
        linkText: 'hardware-wallet.security-warning.link',
        linkFunction: this.adminHwWallet.bind(this),
      };

      showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          wallet.stopShowingHwSecurityPopup = true;
          this.walletService.saveHardwareWallets();
          wallet.opened = true;
        }
      });
    } else {
      wallet.opened = !wallet.opened;
    }
  }
}
