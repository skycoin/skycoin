import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { CreateWalletComponent } from './create-wallet/create-wallet.component';
import { HwOptionsDialogComponent } from '../../layout/hardware-wallet/hw-options-dialog/hw-options-dialog.component';
import { SubscriptionLike } from 'rxjs';
import { Router } from '@angular/router';
import { HwWalletService } from '../../../services/hw-wallet.service';
import { first } from 'rxjs/operators';
import { ConfirmationParams, ConfirmationComponent, DefaultConfirmationButtons } from '../../layout/confirmation/confirmation.component';
import { WalletsAndAddressesService } from '../../../services/wallet-operations/wallets-and-addresses.service';
import { BalanceAndOutputsService } from '../../../services/wallet-operations/balance-and-outputs.service';
import { WalletWithBalance } from '../../../services/wallet-operations/wallet-objects';

@Component({
  selector: 'app-wallets',
  templateUrl: './wallets.component.html',
  styleUrls: ['./wallets.component.scss'],
})
export class WalletsComponent implements OnInit, OnDestroy {

  hwCompatibilityActivated = false;

  wallets: WalletWithBalance[] = [];
  hardwareWallets: WalletWithBalance[] = [];
  walletsOpenedState = new Map<string, boolean>();

  private subscription: SubscriptionLike;

  constructor(
    private hwWalletService: HwWalletService,
    private dialog: MatDialog,
    private router: Router,
    private walletsAndAddressesService: WalletsAndAddressesService,
    private balanceAndOutputsService: BalanceAndOutputsService,
  ) {
    this.hwCompatibilityActivated = this.hwWalletService.hwWalletCompatibilityActivated;

    this.subscription = this.balanceAndOutputsService.walletsWithBalance.subscribe(wallets => {
      this.wallets = [];
      this.hardwareWallets = [];
      wallets.forEach(value => {
        if (!value.isHardware) {
          this.wallets.push(value);
        } else {
          this.hardwareWallets.push(value);
        }

        if (!this.walletsOpenedState.has(value.id)) {
          this.walletsOpenedState.set(value.id, false);
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
      this.walletsAndAddressesService.allWallets.pipe(first()).subscribe(wallets => {
        if (wallets.length === 0) {
          setTimeout(() => this.router.navigate(['/wizard']), 500);
        }
      });
    });
  }

  toggleWallet(wallet: WalletWithBalance) {
    if (wallet.isHardware && wallet.hasHwSecurityWarnings && !wallet.stopShowingHwSecurityPopup && !this.walletsOpenedState.get(wallet.id)) {
      const confirmationParams: ConfirmationParams = {
        headerText: 'hardware-wallet.security-warning.title',
        text: 'hardware-wallet.security-warning.text',
        checkboxText: 'common.generic-confirmation-check',
        defaultButtons: DefaultConfirmationButtons.ContinueCancel,
        linkText: 'hardware-wallet.security-warning.link',
        linkFunction: this.adminHwWallet.bind(this),
      };

      ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          wallet.stopShowingHwSecurityPopup = true;
          this.walletsAndAddressesService.informValuesUpdated(wallet);
          this.walletsOpenedState.set(wallet.id, true);
        }
      });
    } else {
      this.walletsOpenedState.set(wallet.id, !this.walletsOpenedState.get(wallet.id));
    }
  }
}
