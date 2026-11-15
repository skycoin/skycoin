import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogConfig, MatDialogRef } from '@angular/material/dialog';
import { Observable, Subscription, catchError, EMPTY, switchMap, timer } from 'rxjs';

import { Wallet } from '../../../../app.datatypes';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { SweepRule, SweepService, SweepState } from '../../../../services/wallet/sweep.service';

/** One destination the form offers, as the node's rule takes it. */
export interface SweepChoice {
  key: string;
  label: string;
  params?: { [k: string]: string | number };
}

/** Durations an encrypted wallet can be unlocked for, in seconds. 0 is until the node restarts. */
export const sweepUnlockTimes = [
  { seconds: 15 * 60, label: 'wallet.sweep.ttl-15m' },
  { seconds: 60 * 60, label: 'wallet.sweep.ttl-1h' },
  { seconds: 24 * 60 * 60, label: 'wallet.sweep.ttl-24h' },
  { seconds: 7 * 24 * 60 * 60, label: 'wallet.sweep.ttl-7d' },
  { seconds: 0, label: 'wallet.sweep.ttl-restart' },
];

/** How often the status is read while the panel is open. */
export const sweepPollMs = 10000;

/**
 * Sets the node's rule for moving coins off the exposed addresses of a wallet it
 * keeps, unlocks an encrypted one for it, and shows how the sweep is going.
 */
@Component({
  selector: 'app-auto-sweep',
  templateUrl: './auto-sweep.component.html',
  styleUrls: ['./auto-sweep.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: false
})
export class AutoSweepComponent implements OnInit, OnDestroy {
  readonly unlockTimes = sweepUnlockTimes;

  state: SweepState | null = null;
  loadError = '';
  error = '';
  done = '';
  busy = false;

  enabled = false;
  burnOk = false;
  shared = new Set<number>();
  destination = '';
  address = '';
  addresses = '';
  password = '';
  ttl = 3600;

  readonly bip44: boolean;
  readonly accounts: { index: number, name: string }[];
  readonly otherWallets: Wallet[];
  readonly choices: SweepChoice[];

  private subscriptions: Subscription[] = [];

  static openDialog(dialog: CustomMatDialogService, wallet: Wallet): MatDialogRef<AutoSweepComponent> {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.autoFocus = false;
    config.data = wallet;

    return dialog.open(AutoSweepComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) public wallet: Wallet,
    public dialogRef: MatDialogRef<AutoSweepComponent>,
    private sweepService: SweepService,
    walletService: WalletService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    this.bip44 = wallet.walletType === 'bip44';
    this.accounts = (wallet.accounts || []).map(a => ({ index: a.index, name: a.name }));
    if (this.bip44 && this.accounts.length === 0) {
      this.accounts = [{ index: 0, name: '' }];
    }
    this.otherWallets = (walletService.wallets.getValue() || []).filter(w =>
      w.filename && w.filename !== wallet.filename && !w.temporary && !w.isHardware && w.walletType !== 'xpub');

    const choices: SweepChoice[] = [];
    if (this.bip44) {
      choices.push({ key: 'next', label: 'wallet.sweep.to-next' });
      this.accounts.forEach(a => choices.push({ key: 'account:' + a.index, label: 'wallet.sweep.to-account', params: { n: a.index, name: a.name } }));
    } else {
      choices.push({ key: 'self', label: 'wallet.sweep.to-self' });
    }
    this.otherWallets.forEach(w => choices.push({ key: 'wallet:' + w.filename, label: 'wallet.sweep.to-wallet', params: { name: w.label || w.filename! } }));
    choices.push({ key: 'address', label: 'wallet.sweep.to-address' });
    this.choices = choices;
  }

  ngOnInit() {
    this.subscriptions.push(this.sweepService.get(this.wallet.filename!).subscribe({
      next: state => {
        this.apply(state, true);
        this.changeDetectorRef.markForCheck();
      },
      error: (err: Error) => {
        this.loadError = err.message;
        this.changeDetectorRef.markForCheck();
      },
    }));
    this.subscriptions.push(timer(sweepPollMs, sweepPollMs).pipe(
      switchMap(() => this.sweepService.get(this.wallet.filename!).pipe(catchError(() => EMPTY))),
    ).subscribe(state => {
      if (this.state) {
        this.state = { ...this.state, status: state.status };
      }
      this.changeDetectorRef.markForCheck();
    }));
  }

  ngOnDestroy() {
    this.subscriptions.forEach(s => s.unsubscribe());
  }

  closePopup() {
    this.dialogRef.close();
  }

  /** The node sweeps the wallet now. */
  get active(): boolean {
    return !!this.state && this.state.status.enabled && !this.state.status.paused;
  }

  /** Turning the sweep on, which needs the cost confirmed first. */
  get enabling(): boolean {
    return this.enabled && !(this.state && this.state.status.enabled);
  }

  toggleShared(index: number, on: boolean) {
    if (on) {
      this.shared.add(index);
    } else {
      this.shared.delete(index);
    }
  }

  /** The rule the form describes, on top of the node's, so fields it does not show are kept. */
  buildRule(): SweepRule {
    const base = this.state && this.state.rule ? this.state.rule : null;
    const rule: SweepRule = {
      ...(base || {}),
      wallet: this.wallet.filename!,
      enabled: this.enabled,
      destination: { kind: 'account' },
    };
    rule.shared_accounts = this.bip44 ? Array.from(this.shared).sort((a, b) => a - b) : [];
    rule.addresses = this.addresses.split(/\s+/).filter(a => a);
    const [kind, value] = splitChoice(this.destination);
    if (kind === 'next') {
      rule.destination = { kind: 'account' };
    } else if (kind === 'account') {
      rule.destination = { kind: 'account', account: Number(value) };
    } else if (kind === 'self') {
      rule.destination = { kind: 'wallet' };
    } else if (kind === 'wallet') {
      rule.destination = { kind: 'wallet', wallet: value };
    } else {
      rule.destination = { kind: 'address', address: this.address.trim() };
    }

    return rule;
  }

  save() {
    if (this.busy) {
      return;
    }
    if (this.enabling && !this.burnOk) {
      this.error = 'wallet.sweep.err-burn';

      return;
    }
    if (this.destination === 'address' && !this.address.trim()) {
      this.error = 'wallet.sweep.err-destination';

      return;
    }
    this.run(this.sweepService.setRule(this.buildRule()), 'wallet.sweep.done-saved', true);
  }

  unlock() {
    if (this.busy) {
      return;
    }
    if (!this.password) {
      this.error = 'wallet.sweep.err-password';

      return;
    }
    const password = this.password;
    this.password = '';
    this.run(this.sweepService.unlock(this.wallet.filename!, password, Number(this.ttl)), 'wallet.sweep.done-unlocked', false);
  }

  lock() {
    if (!this.busy) {
      this.run(this.sweepService.lock(this.wallet.filename!), 'wallet.sweep.done-locked', false);
    }
  }

  private run(request: Observable<SweepState>, done: string, ruleChanged: boolean) {
    this.busy = true;
    this.error = '';
    this.done = '';
    this.subscriptions.push(request.subscribe({
      next: state => {
        this.busy = false;
        this.done = done;
        this.apply(state, ruleChanged);
        this.changeDetectorRef.markForCheck();
      },
      error: (err: Error) => {
        this.busy = false;
        this.error = err && err.message ? err.message : String(err);
        this.changeDetectorRef.markForCheck();
      },
    }));
  }

  /** Takes the node's answer, and its rule into the form when it changed. */
  private apply(state: SweepState, fillForm: boolean) {
    this.state = state;
    if (!fillForm) {
      return;
    }
    const rule = state.rule;
    this.enabled = !!rule && rule.enabled;
    this.burnOk = false;
    this.shared = new Set(rule && rule.shared_accounts ? rule.shared_accounts : []);
    this.addresses = rule && rule.addresses ? rule.addresses.join('\n') : '';
    this.address = rule && rule.destination.kind === 'address' ? rule.destination.address || '' : '';
    this.destination = this.choiceOf(rule);
  }

  private choiceOf(rule: SweepRule | null): string {
    const d = rule ? rule.destination : null;
    if (d && d.kind === 'account') {
      return d.account !== undefined && d.account !== null ? 'account:' + d.account : 'next';
    }
    if (d && d.kind === 'wallet') {
      return d.wallet ? 'wallet:' + d.wallet : 'self';
    }
    if (d && d.kind === 'address') {
      return 'address';
    }

    return this.bip44 ? 'next' : 'self';
  }
}

function splitChoice(choice: string): [string, string] {
  const i = choice.indexOf(':');

  return i < 0 ? [choice, ''] : [choice.substring(0, i), choice.substring(i + 1)];
}
