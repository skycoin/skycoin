import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnDestroy } from '@angular/core';
import { MatDialogConfig } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { of, Subscription, delay, first, mergeMap } from 'rxjs';

import { ConfirmationData, Wallet, Address, Bip44Account } from '../../../../app.datatypes';
import { BaseCoin } from '../../../../coins/basecoin';
import { CoinService } from '../../../../services/coin.service';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { ChangeNameComponent } from '../change-name/change-name.component';
import { openUnlockWalletModal, openQrModal, showConfirmationModal, openDeleteWalletModal } from '../../../../utils/index';
import { WalletOptionsComponent, WalletOptionsResponses } from './wallet-options/wallet-options.component';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { config } from '../../../../app.config';
import { MsgBarService } from '../../../../services/msg-bar.service';

@Component({
    selector: 'app-wallet-detail',
    templateUrl: './wallet-detail.component.html',
    styleUrls: ['./wallet-detail.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class WalletDetailComponent implements OnDestroy {
  @Input() wallet!: Wallet;

  currentCoin!: BaseCoin;
  creatingAddress = false;
  showSlowMobileInfo = false;
  spinnerStyle: any;

  private unlockSubscription!: Subscription;
  private slowInfoSubscription!: Subscription;
  private coinSubscription: Subscription;

  constructor(
    private walletService: WalletService,
    private dialog: CustomMatDialogService,
    private translateService: TranslateService,
    private msgBarService: MsgBarService,
    private hwWalletService: HwWalletService,
    coinService: CoinService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    this.coinSubscription = coinService.currentCoin.subscribe(coin => {
      this.currentCoin = coin;
      this.changeDetectorRef.markForCheck();
    });
  }

  ngOnDestroy() {
    this.msgBarService.hide();
    this.removeUnlockSubscription();
    this.removeSlowInfoSubscription();
    if (this.coinSubscription) {
      this.coinSubscription.unsubscribe();
    }
  }

  onShowQr(address: Address) {
    openQrModal(this.dialog, address.address, true);
  }

  onShowQrIfPointer(component: HTMLDivElement, address: Address) {
    if (getComputedStyle(component).cursor === 'pointer') {
      this.onShowQr(address);
    }
  }

  onEditWallet() {
    const dialogConfig = new MatDialogConfig();
    dialogConfig.width = '566px';
    dialogConfig.data = this.wallet;
    dialogConfig.autoFocus = false;
    this.dialog.open(ChangeNameComponent, dialogConfig);
  }

  onShowOptions() {
    const dialogConfig = new MatDialogConfig();
    dialogConfig.width = '566px';
    dialogConfig.autoFocus = false;
    dialogConfig.data = this.wallet;
    this.dialog.open(WalletOptionsComponent, dialogConfig).afterClosed().subscribe((response: WalletOptionsResponses | null) => {
      if (response != null && response !== undefined) {
        if (response === WalletOptionsResponses.UnlockWallet) {
          openUnlockWalletModal(this.wallet, this.dialog);
        } else if (response === WalletOptionsResponses.AddNewAddress) {
          this.onAddNewAddress();
        } else if (response === WalletOptionsResponses.EditWallet) {
          this.onEditWallet();
        } else if (response === WalletOptionsResponses.DeleteWallet) {
          this.onDeleteWallet();
        } else if (response === WalletOptionsResponses.EncryptWallet) {
          this.onEncryptWallet();
        } else if (response === WalletOptionsResponses.DecryptWallet) {
          this.onDecryptWallet();
        }
      }
      this.changeDetectorRef.markForCheck();
    });
  }

  onAddNewAddress(accountIndex?: number) {
    const addressCount = accountIndex !== undefined
      ? (this.wallet.accounts?.[accountIndex]?.externalAddresses?.length || 0)
      : this.wallet.addresses.length;

    if (addressCount < 5) {
      this.verifyBeforeAddingNewAddress(accountIndex);
    } else {
      const confirmationData: ConfirmationData = {
        text: 'wallet.add-confirmation',
        headerText: 'confirmation.header-text',
        confirmButtonText: 'confirmation.confirm-button',
        cancelButtonText: 'confirmation.cancel-button'
      };

      showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(result => {
        if (result) {
          this.verifyBeforeAddingNewAddress(accountIndex);
        }
        this.changeDetectorRef.markForCheck();
      });
    }
  }

  onCopySuccess(address: Address, interval = 500) {
    if (address.isCopying) {
      return;
    }

    address.isCopying = true;

    // wait for a while and then remove the 'copying' class
    setTimeout(() => {
      address.isCopying = false;
      this.changeDetectorRef.markForCheck();
    }, interval);
  }

  onToggleEmpty() {
    this.wallet.hideEmpty = !this.wallet.hideEmpty;
  }

  get isBip44(): boolean {
    return this.wallet.walletType === 'bip44';
  }

  onToggleXpub(account: Bip44Account) {
    if (account.showXpub) {
      account.showXpub = false;
      return;
    }

    if (account.accountXpubKey) {
      account.showXpub = true;
      return;
    }

    // Fetch all three xpub levels
    this.walletService.getXPubKey(this.wallet, `${account.index}`).subscribe(
      (xpub) => {
        account.accountXpubKey = xpub;
        account.showXpub = true;
        this.changeDetectorRef.markForCheck();
      },
      (error) => {
        this.msgBarService.showError(error.message);
        this.changeDetectorRef.markForCheck();
      }
    );
    this.walletService.getXPubKey(this.wallet, `${account.index}/0`).subscribe(
      (xpub) => {
        account.externalXpubKey = xpub;
        this.changeDetectorRef.markForCheck();
      },
      (error) => {
        this.msgBarService.showError(error.message);
        this.changeDetectorRef.markForCheck();
      }
    );
    this.walletService.getXPubKey(this.wallet, `${account.index}/1`).subscribe(
      (xpub) => {
        account.changeXpubKey = xpub;
        this.changeDetectorRef.markForCheck();
      },
      (error) => {
        this.msgBarService.showError(error.message);
        this.changeDetectorRef.markForCheck();
      }
    );
  }

  onToggleChangeAddresses(account: Bip44Account) {
    account.showChangeAddresses = !account.showChangeAddresses;
  }

  onAddAccount() {
    const name = `Account ${(this.wallet.accounts || []).length}`;
    this.walletService.addAccount(this.wallet, name).subscribe(
      () => { this.changeDetectorRef.markForCheck(); },
      (error) => {
        this.msgBarService.showError(error.message);
        this.changeDetectorRef.markForCheck();
      }
    );
  }

  onEncryptWallet() {
    const password = window.prompt(this.translateService.instant('password.title'));
    if (!password) {
      return;
    }
    const confirm = window.prompt(this.translateService.instant('password.confirm-label'));
    if (password !== confirm) {
      this.msgBarService.showError('Passwords do not match');
      return;
    }

    if (this.currentCoin && this.currentCoin.serverWallets && this.wallet.filename) {
      // Server-managed wallet encryption
      this.walletService.encryptWallet(this.wallet, password).subscribe(
        () => {
          this.wallet.encrypted = true;
          this.msgBarService.showDone('Wallet encrypted successfully');
          this.changeDetectorRef.markForCheck();
        },
        (error) => {
          this.msgBarService.showError(error.message || error);
          this.changeDetectorRef.markForCheck();
        }
      );
    } else {
      // Browser-only wallet: encrypt seed in localStorage
      if (!this.wallet.seed) {
        this.removeUnlockSubscription();
        this.unlockSubscription = openUnlockWalletModal(this.wallet, this.dialog).componentInstance.onWalletUnlocked.pipe(first())
          .subscribe(() => {
            this.encryptBrowserWallet(password);
            this.changeDetectorRef.markForCheck();
          });
      } else {
        this.encryptBrowserWallet(password);
      }
    }
  }

  onDecryptWallet() {
    const password = window.prompt(this.translateService.instant('password.title'));
    if (!password) {
      return;
    }

    if (this.currentCoin && this.currentCoin.serverWallets && this.wallet.filename) {
      // Server-managed wallet decryption
      this.walletService.decryptWallet(this.wallet, password).subscribe(
        () => {
          this.wallet.encrypted = false;
          this.msgBarService.showDone('Wallet decrypted successfully');
          this.changeDetectorRef.markForCheck();
        },
        (error) => {
          this.msgBarService.showError(error.message || error);
          this.changeDetectorRef.markForCheck();
        }
      );
    } else {
      // Browser-only wallet: remove encrypted seed from localStorage
      this.walletService.removeWalletPassword(this.wallet);
      this.msgBarService.showDone('Wallet password removed');
    }
  }

  private encryptBrowserWallet(password: string) {
    this.walletService.setWalletPassword(this.wallet, password).subscribe(
      () => {
        this.msgBarService.showDone('Wallet encrypted successfully');
        this.changeDetectorRef.markForCheck();
      },
      (error) => {
        this.msgBarService.showError(error.message || error);
        this.changeDetectorRef.markForCheck();
      }
    );
  }

  onDeleteWallet() {
    openDeleteWalletModal(this.dialog, this.wallet, this.translateService, this.walletService);
  }

  private verifyBeforeAddingNewAddress(accountIndex?: number) {
    if (this.wallet.isHardware) {
      this.addHwAddress();
    } else if (this.currentCoin && this.currentCoin.serverWallets && this.wallet.filename) {
      // Server-managed wallets don't need seed entry
      this.addNewAddress(accountIndex);
    } else if (!this.wallet.seed || !this.wallet.nextSeed) {
      this.removeUnlockSubscription();

      this.unlockSubscription = openUnlockWalletModal(this.wallet, this.dialog).componentInstance.onWalletUnlocked.pipe(first())
        .subscribe(() => {
          this.addNewAddress();
          this.changeDetectorRef.markForCheck();
        });
    } else {
      this.addNewAddress(accountIndex);
    }
  }

  private addHwAddress() {
    if (this.creatingAddress) {
      this.msgBarService.showError('wallet.already-adding-address-error');
      return;
    }

    this.creatingAddress = true;

    // Verify the correct device is connected, then generate the next address
    this.hwWalletService.checkIfCorrectHwConnected(this.wallet.addresses[0].address).pipe(
      mergeMap(() => this.hwWalletService.getAddresses(1, this.wallet.addresses.length))
    ).subscribe(
      (result) => {
        const address = Array.isArray(result.rawResponse) ? result.rawResponse[0] : result.rawResponse;
        this.wallet.addresses.push({ address: address });
        this.walletService.saveWallets();
        this.creatingAddress = false;
        this.changeDetectorRef.markForCheck();
      },
      (error) => {
        this.creatingAddress = false;
        this.msgBarService.showError(error.translatableErrorMsg || error.message);
        this.changeDetectorRef.markForCheck();
      }
    );
  }

  private addNewAddress(accountIndex?: number) {
    if (this.creatingAddress === true) {
      this.msgBarService.showError('wallet.already-adding-address-error');

      return;
    }

    this.creatingAddress = true;

    this.slowInfoSubscription = of(1).pipe(delay(config.timeBeforeSlowMobileInfo))
      .subscribe(() => {
        this.showSlowMobileInfo = true;
        this.changeDetectorRef.markForCheck();
      });

    setTimeout(() => {
      this.walletService.addAddress(this.wallet, true, accountIndex)
        .subscribe(
          () => {
            this.showSlowMobileInfo = false;
            this.removeSlowInfoSubscription();
            this.creatingAddress = false;
            this.changeDetectorRef.markForCheck();
          },
          (error: Error) => {
            this.onAddAddressError(error);
            this.changeDetectorRef.markForCheck();
          }
        );
      this.changeDetectorRef.markForCheck();
    }, 0);
  }

  onAddNewChangeAddress(accountIndex: number) {
    if (this.creatingAddress) {
      this.msgBarService.showError('wallet.already-adding-address-error');
      return;
    }

    this.creatingAddress = true;

    setTimeout(() => {
      this.walletService.addAddress(this.wallet, true, accountIndex, 1)
        .subscribe(
          () => {
            this.creatingAddress = false;
            this.changeDetectorRef.markForCheck();
          },
          (error: Error) => {
            this.onAddAddressError(error);
            this.changeDetectorRef.markForCheck();
          }
        );
      this.changeDetectorRef.markForCheck();
    }, 0);
  }

  private onAddAddressError(error: Error) {
    this.showSlowMobileInfo = false;
    this.removeSlowInfoSubscription();
    this.creatingAddress = false;
    this.msgBarService.showError(error.message);
  }

  private removeUnlockSubscription() {
    if (this.unlockSubscription) {
      this.unlockSubscription.unsubscribe();
    }
  }

  private removeSlowInfoSubscription() {
    if (this.slowInfoSubscription) {
      this.slowInfoSubscription.unsubscribe();
    }
  }
}
