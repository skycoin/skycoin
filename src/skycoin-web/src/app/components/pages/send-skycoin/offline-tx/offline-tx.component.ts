import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Observable, Subscription, first, mergeMap, of } from 'rxjs';

import { Wallet } from '../../../../app.datatypes';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { SpendingService } from '../../../../services/wallet/spending.service';
import { BalanceService } from '../../../../services/wallet/balance.service';
import { ClipboardService } from '../../../../services/clipboard.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { openUnlockWalletModal } from '../../../../utils';

export enum OfflineTxMode {
  /** Shows a raw transaction to copy. */
  Show = 'show',
  /** Signs a pasted unsigned transaction with a wallet. */
  Sign = 'sign',
  /** Sends a pasted signed transaction to the network. */
  Broadcast = 'broadcast',
}

export interface OfflineTxData {
  mode: OfflineTxMode;
  /** For Show, the encoded transaction. */
  rawTx?: string;
  /** For Show, whether rawTx is still unsigned. */
  unsigned?: boolean;
}

/**
 * Moves transactions to and from a wallet that is not online: shows one to copy,
 * signs a pasted unsigned one, or broadcasts a pasted signed one.
 */
@Component({
    selector: 'app-offline-tx',
    templateUrl: './offline-tx.component.html',
    styleUrls: ['./offline-tx.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class OfflineTxComponent implements OnInit, OnDestroy {
  modes = OfflineTxMode;
  form: UntypedFormGroup;
  wallets: Wallet[] = [];
  working = false;
  error = '';
  copied = false;

  private subscription?: Subscription;

  static open(dialog: MatDialog, data: OfflineTxData): MatDialogRef<OfflineTxComponent> {
    return dialog.open(OfflineTxComponent, { width: '600px', data });
  }

  constructor(
    public dialogRef: MatDialogRef<OfflineTxComponent>,
    @Inject(MAT_DIALOG_DATA) public data: OfflineTxData,
    formBuilder: UntypedFormBuilder,
    private walletService: WalletService,
    private spendingService: SpendingService,
    private balanceService: BalanceService,
    private clipboardService: ClipboardService,
    private dialog: CustomMatDialogService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    this.form = formBuilder.group({
      wallet: [null],
      rawTx: [data.rawTx || ''],
    });
  }

  ngOnInit() {
    if (this.data.mode === OfflineTxMode.Sign) {
      this.walletService.currentWallets.pipe(first()).subscribe(wallets => {
        this.wallets = wallets.filter(w => !w.isHardware);
        if (this.wallets.length === 1) {
          this.form.get('wallet')!.setValue(this.wallets[0]);
        }
        this.changeDetectorRef.markForCheck();
      });
    }
  }

  ngOnDestroy() {
    this.subscription?.unsubscribe();
  }

  get title(): string {
    if (this.data.mode === OfflineTxMode.Show) {
      return this.data.unsigned ? 'offline-tx.unsigned-title' : 'offline-tx.signed-title';
    }

    return 'offline-tx.' + this.data.mode + '-title';
  }

  get text(): string {
    if (this.data.mode === OfflineTxMode.Show) {
      return this.data.unsigned ? 'offline-tx.unsigned-text' : 'offline-tx.signed-text';
    }

    return 'offline-tx.' + this.data.mode + '-text';
  }

  get canSubmit(): boolean {
    const rawTx = (this.form.value.rawTx || '').trim();

    return !this.working && rawTx.length > 0 && (this.data.mode !== OfflineTxMode.Sign || !!this.form.value.wallet);
  }

  copy() {
    this.clipboardService.copy(this.form.value.rawTx).then(() => {
      this.copied = true;
      this.changeDetectorRef.markForCheck();
    });
  }

  close() {
    this.dialogRef.close();
  }

  sign() {
    const wallet: Wallet = this.form.value.wallet;
    this.run(this.unlock(wallet).pipe(
      mergeMap(password => password === undefined ? of(undefined) :
        this.spendingService.signTransaction(wallet, password || undefined, this.form.value.rawTx)),
    ), signed => {
      this.dialogRef.close();
      OfflineTxComponent.open(this.dialog, { mode: OfflineTxMode.Show, rawTx: signed, unsigned: false });
    });
  }

  broadcast() {
    this.run(this.spendingService.injectTransaction(this.form.value.rawTx.trim()), () => {
      this.balanceService.startGettingBalances();
      this.dialogRef.close(true);
    });
  }

  /**
   * What signing with wallet needs first: the password of an encrypted server
   * wallet, or the seed of a locked browser wallet. Emits undefined when cancelled.
   */
  private unlock(wallet: Wallet): Observable<string | null | undefined> {
    if (this.walletService.isServerWallet(wallet)) {
      return wallet.encrypted ? PasswordDialogComponent.ask(this.dialog, { text: wallet.label }) : of(null);
    }
    if (!wallet.seed) {
      const unlocked = openUnlockWalletModal(wallet, this.dialog);

      return new Observable<string | null | undefined>(subscriber => {
        const done = unlocked.componentInstance.onWalletUnlocked.pipe(first()).subscribe(() => {
          subscriber.next(null);
          subscriber.complete();
          this.changeDetectorRef.markForCheck();
        });
        const closed = unlocked.afterClosed().subscribe(() => {
          if (!wallet.seed) {
            subscriber.next(undefined);
            subscriber.complete();
          }
          this.changeDetectorRef.markForCheck();
        });

        return () => { done.unsubscribe(); closed.unsubscribe(); };
      });
    }

    return of(null);
  }

  private run<T>(operation: Observable<T | undefined>, onDone: (result: T) => void) {
    this.working = true;
    this.error = '';
    this.subscription?.unsubscribe();
    this.subscription = operation.subscribe(result => {
      this.working = false;
      if (result !== undefined) {
        onDone(result);
      }
      this.changeDetectorRef.markForCheck();
    }, err => {
      this.working = false;
      this.error = err && err.message ? err.message : String(err);
      this.changeDetectorRef.markForCheck();
    });
  }
}
