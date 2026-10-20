import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';

import { Wallet } from '../../../../app.datatypes';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { MsgBarService } from '../../../../services/msg-bar.service';

/**
 * Backs up the current coin's wallets: their addresses as a CSV file, and their seeds.
 */
@Component({
    selector: 'app-backup',
    templateUrl: './backup.component.html',
    styleUrls: ['./backup.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class BackupComponent implements OnInit, OnDestroy {
  wallets: Wallet[] | null = null;
  // Where a wallet service keeps the files, or null for wallets in this browser.
  folder: string | null = null;
  // The seed shown for each wallet, by label, while it is open.
  shownSeeds = new Map<Wallet, string>();

  private walletsSubscription!: Subscription;
  private folderSubscription!: Subscription;
  private seedSubscription?: Subscription;

  constructor(
    private walletService: WalletService,
    private msgBarService: MsgBarService,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit() {
    this.walletsSubscription = this.walletService.currentWallets.subscribe(wallets => {
      this.wallets = wallets;
      this.changeDetectorRef.markForCheck();
    });
    this.folderSubscription = this.walletService.walletFolder().subscribe(folder => {
      this.folder = folder;
      this.changeDetectorRef.markForCheck();
    }, () => {
      this.folder = '?';
      this.changeDetectorRef.markForCheck();
    });
  }

  ngOnDestroy() {
    this.walletsSubscription.unsubscribe();
    this.folderSubscription.unsubscribe();
    this.seedSubscription?.unsubscribe();
    this.shownSeeds.clear();
    this.msgBarService.hide();
  }

  /** True when a seed can be shown: not for hardware wallets, and on a server only when encrypted. */
  canShowSeed(wallet: Wallet): boolean {
    if (wallet.isHardware) {
      return false;
    }

    return this.walletService.isServerWallet(wallet) ? !!wallet.encrypted : !!(wallet.seed || wallet.encryptedSeed);
  }

  /** The translation key that says why the seed of wallet cannot be shown. */
  seedUnavailableReason(wallet: Wallet): string {
    if (wallet.isHardware) {
      return 'backup.hw-wallet-info';
    }

    return this.walletService.isServerWallet(wallet) ? 'backup.unencrypted-info' : 'backup.no-seed-info';
  }

  /** Asks for the password where the seed is encrypted, then shows it. */
  showSeed(wallet: Wallet) {
    const encrypted = this.walletService.isServerWallet(wallet) ? !!wallet.encrypted : !wallet.seed;
    let password: string | null = null;
    if (encrypted) {
      password = window.prompt(this.passwordPrompt(wallet));
      if (password === null) {
        return;
      }
    }
    this.msgBarService.hide();
    this.seedSubscription?.unsubscribe();
    this.seedSubscription = this.walletService.getSeed(wallet, password).subscribe(seed => {
      this.shownSeeds.set(wallet, seed);
      this.changeDetectorRef.markForCheck();
    }, err => {
      this.msgBarService.showError(err && err.message ? err.message : String(err));
      this.changeDetectorRef.markForCheck();
    });
  }

  hideSeed(wallet: Wallet) {
    this.shownSeeds.delete(wallet);
  }

  /** Downloads the wallet's addresses, one per line. */
  saveAddresses(wallet: Wallet) {
    const blob = new Blob([wallet.addresses.map(a => a.address).join('\n')], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', wallet.label + '.csv');
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  }

  private passwordPrompt(wallet: Wallet): string {
    return 'Password for ' + wallet.label;
  }
}
