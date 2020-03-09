import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { SubscriptionLike } from 'rxjs';
import { first } from 'rxjs/operators';
import { Router } from '@angular/router';

import { CreateWalletComponent } from './create-wallet/create-wallet.component';
import { HwOptionsDialogComponent } from '../../layout/hardware-wallet/hw-options-dialog/hw-options-dialog.component';
import { HwWalletService } from '../../../services/hw-wallet.service';
import { ConfirmationParams, ConfirmationComponent, DefaultConfirmationButtons } from '../../layout/confirmation/confirmation.component';
import { WalletsAndAddressesService } from '../../../services/wallet-operations/wallets-and-addresses.service';
import { BalanceAndOutputsService } from '../../../services/wallet-operations/balance-and-outputs.service';
import { WalletWithBalance } from '../../../services/wallet-operations/wallet-objects';

/**
 * Shows the wallet list and options related to it.
 */
@Component({
  selector: 'app-wallets',
  templateUrl: './wallets.component.html',
  styleUrls: ['./wallets.component.scss'],
})
export class WalletsComponent implements OnInit, OnDestroy {
  /**
   * Allow to know if the page is busy preparing an operation and no other operation must
   * be stated before finishing.
   */
  static busy = false;

  // If the hw wallet options must be shown.
  hwCompatibilityActivated = false;

  // Software wallets to show on the list.
  wallets: WalletWithBalance[] = [];
  // Hardware wallets to show on the list.
  hardwareWallets: WalletWithBalance[] = [];
  // Saves which wallet panels are open.
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

    // Keep the wallet list updated.
    this.subscription = this.balanceAndOutputsService.walletsWithBalance.subscribe(wallets => {
      this.wallets = [];
      this.hardwareWallets = [];

      // Create a map wit the wallets and put each wallet on the appropiate array.
      const walletsMap = new Map<string, boolean>();
      wallets.forEach(value => {
        walletsMap.set(value.id, true);

        if (!value.isHardware) {
          this.wallets.push(value);
        } else {
          this.hardwareWallets.push(value);
        }

        // If it is a new wallet, set it as closed.
        if (!this.walletsOpenedState.has(value.id)) {
          this.walletsOpenedState.set(value.id, false);
        }
      });

      // Remove from walletsOpenedState all the deleted wallets.
      const walletsToRemove: string[] = [];
      this.walletsOpenedState.forEach((value, key) => {
        if (!walletsMap.has(key)) {
          walletsToRemove.push(key);
        }
      });
      walletsToRemove.forEach(walletToRemove => {
        this.walletsOpenedState.delete(walletToRemove);
      });
    });
  }

  ngOnInit(): void {
    // Open the hw wallet options if it was requested before opening the page.
    if (this.hwWalletService.showOptionsWhenPossible) {
      setTimeout(() => {
        this.hwWalletService.showOptionsWhenPossible = false;
        this.adminHwWallet();
      });
    }
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
    WalletsComponent.busy = false;
  }

  // Opens the create wallet modal window, for creating a new wallet (create === true) or loading
  // an old one.
  addWallet(create: boolean) {
    CreateWalletComponent.openDialog(this.dialog, { create });
  }

  // Opens the hw wallet options modal window.
  adminHwWallet() {
    HwOptionsDialogComponent.openDialog(this.dialog, false).afterClosed().subscribe(() => {
      // Check if there are still wallets on the wallet list. If not, go to the wizard.
      this.walletsAndAddressesService.allWallets.pipe(first()).subscribe(wallets => {
        if (wallets.length === 0) {
          setTimeout(() => this.router.navigate(['/wizard']), 500);
        }
      });
    });
  }

  // Opens or closes a wallet panel.
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

      // If there are security warnings related to the hw wallet, ask for confirmation before opening the panel.
      ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          // Do not show the warning in the future and open the panel.
          wallet.stopShowingHwSecurityPopup = true;
          this.walletsAndAddressesService.informValuesUpdated(wallet);
          this.walletsOpenedState.set(wallet.id, true);
        }
      });
    } else {
      // Open or close the panel.
      this.walletsOpenedState.set(wallet.id, !this.walletsOpenedState.get(wallet.id));
    }
  }
}
