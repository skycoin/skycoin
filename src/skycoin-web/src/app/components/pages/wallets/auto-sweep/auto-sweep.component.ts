import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogConfig, MatDialogRef } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { Observable, Subscription, catchError, EMPTY, switchMap, timer } from 'rxjs';

import { Wallet } from '../../../../app.datatypes';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import {
  SweepBackend, SweepDestination, SweepMapping, SweepRule, SweepService, SweepState, sweepActive,
} from '../../../../services/wallet/sweep.service';
import { BrowserSweepService, browserWalletId } from '../../../../services/wallet/browser-sweep.service';

/** One destination the form offers, as the rule takes it. */
export interface SweepChoice {
  key: string;
  label: string;
  params?: { [k: string]: string | number };
}

/** Durations an encrypted wallet can be unlocked for, in seconds. 0 is until the node restarts or the page is left. */
export const sweepUnlockTimes = [
  { seconds: 15 * 60, label: 'wallet.sweep.ttl-15m' },
  { seconds: 60 * 60, label: 'wallet.sweep.ttl-1h' },
  { seconds: 24 * 60 * 60, label: 'wallet.sweep.ttl-24h' },
  { seconds: 7 * 24 * 60 * 60, label: 'wallet.sweep.ttl-7d' },
  { seconds: 0, label: 'wallet.sweep.ttl-restart' },
];

/** The same for a wallet kept in the browser, whose password lives in this page. */
export const browserUnlockTimes = [
  { seconds: 15 * 60, label: 'wallet.sweep.ttl-15m' },
  { seconds: 60 * 60, label: 'wallet.sweep.ttl-1h' },
  { seconds: 24 * 60 * 60, label: 'wallet.sweep.ttl-24h' },
  { seconds: 0, label: 'wallet.sweep.ttl-page' },
];

/** How often the status is read while the panel is open. */
export const sweepPollMs = 10000;

/**
 * Sets the rule for moving coins off the exposed addresses of a wallet, unlocks
 * an encrypted one for it, and shows how the sweep is going. The node sweeps a
 * wallet it keeps. This page sweeps one kept in the browser, while it is open.
 */
@Component({
  selector: 'app-auto-sweep',
  templateUrl: './auto-sweep.component.html',
  styleUrls: ['./auto-sweep.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: false
})
export class AutoSweepComponent implements OnInit, OnDestroy {
  /** Who sweeps this wallet. A hardware wallet is never swept. */
  readonly kind: 'node' | 'browser' | 'hardware';
  readonly unlockTimes: { seconds: number, label: string }[];

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
  includeSpent = true;
  maps: SweepMapping[] = [];
  mapFrom = '';
  mapTo = 'address';
  mapToAddress = '';
  password = '';
  ttl = 3600;

  readonly id: string;
  readonly encrypted: boolean;
  readonly bip44: boolean;
  readonly accounts: { index: number, name: string }[];
  readonly otherWallets: { id: string, name: string }[];
  readonly choices: SweepChoice[];
  readonly mapChoices: SweepChoice[];

  private backend: SweepBackend;
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
    sweepService: SweepService,
    browserSweep: BrowserSweepService,
    walletService: WalletService,
    private translate: TranslateService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    const server = walletService.isServerWallet(wallet);
    this.kind = wallet.isHardware ? 'hardware' : server ? 'node' : 'browser';
    this.backend = server ? sweepService : browserSweep;
    this.unlockTimes = server ? sweepUnlockTimes : browserUnlockTimes;
    this.id = server ? wallet.filename! : browserWalletId(wallet);
    this.encrypted = server ? !!wallet.encrypted : !!wallet.encryptedSeed;
    this.bip44 = wallet.walletType === 'bip44';
    this.accounts = server ? (wallet.accounts || []).map(a => ({ index: a.index, name: a.name })) : [];
    if (this.bip44 && server && this.accounts.length === 0) {
      this.accounts = [{ index: 0, name: '' }];
    }
    const all = walletService.wallets.getValue() || [];
    this.otherWallets = server
      ? all.filter(w => w.filename && w.filename !== wallet.filename && !w.temporary && !w.isHardware && w.walletType !== 'xpub')
        .map(w => ({ id: w.filename!, name: w.label || w.filename! }))
      : all.filter(w => !walletService.isServerWallet(w) && !w.temporary && !w.isHardware && w.walletType !== 'xpub' && w !== wallet)
        .map(w => ({ id: browserWalletId(w), name: w.label || browserWalletId(w) }))
        .filter(w => w.id && w.id !== this.id);

    const own: SweepChoice[] = [];
    if (this.bip44 && server) {
      own.push({ key: 'next', label: 'wallet.sweep.to-next' });
      this.accounts.forEach(a => own.push({ key: 'account:' + a.index, label: 'wallet.sweep.to-account', params: { n: a.index, name: a.name } }));
    } else {
      own.push({ key: 'self', label: 'wallet.sweep.to-self' });
    }
    const others = this.otherWallets.map(w => ({ key: 'wallet:' + w.id, label: 'wallet.sweep.to-wallet', params: { name: w.name } }));
    const fixed = { key: 'address', label: 'wallet.sweep.to-address' };
    this.choices = [...own, ...others, fixed];
    this.mapChoices = [fixed, ...own, ...others];
  }

  ngOnInit() {
    if (this.kind === 'hardware') {
      return;
    }
    this.subscriptions.push(this.backend.get(this.id).subscribe({
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
      switchMap(() => this.backend.get(this.id).pipe(catchError(() => EMPTY))),
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

  /** The rule sweeps now. */
  get active(): boolean {
    return !!this.state && this.state.status.active && !this.state.status.paused;
  }

  /** The rule does not sweep yet, so starting it needs the cost confirmed. */
  get starting(): boolean {
    return !(this.state && this.state.status.active);
  }

  /** The advanced section holds something, so it starts open. */
  get advancedInUse(): boolean {
    return this.maps.length > 0 || !this.includeSpent || this.addresses.trim() !== '';
  }

  /** Why the sweep waits. A browser wallet reports it as a code. */
  get pausedText(): string {
    const st = this.state ? this.state.status : null;
    if (!st || !st.paused) {
      return '';
    }
    if (st.paused === 'locked') {
      return this.translate.instant('wallet.sweep.paused-locked', { coins: st.waiting_coins });
    }
    if (st.paused === 'no-hours') {
      return this.translate.instant('wallet.sweep.paused-no-hours', { coins: st.waiting_coins });
    }

    return st.paused;
  }

  toggleShared(index: number, on: boolean) {
    if (on) {
      this.shared.add(index);
    } else {
      this.shared.delete(index);
    }
  }

  sourceText(m: SweepMapping): string {
    return m.from.account !== undefined && m.from.account !== null
      ? this.translate.instant('wallet.sweep.account', { n: m.from.account, name: '' }) : m.from.address || '';
  }

  destText(d: SweepDestination): string {
    if (d.kind === 'address') {
      return d.address || '';
    }
    if (d.kind === 'account') {
      return d.account === undefined || d.account === null
        ? this.translate.instant('wallet.sweep.to-next') : this.translate.instant('wallet.sweep.account', { n: d.account, name: '' });
    }
    if (!d.wallet) {
      return this.translate.instant('wallet.sweep.to-self');
    }
    const w = this.otherWallets.find(o => o.id === d.wallet);

    return w ? w.name : d.wallet;
  }

  /** The rule the form describes, on top of the stored one, so fields it does not show are kept. */
  buildRule(): SweepRule {
    const base = this.state && this.state.rule ? this.state.rule : null;
    const rule: SweepRule = {
      ...(base || {}),
      wallet: this.id,
      enabled: this.enabled,
      destination: choiceDestination(this.destination, this.address),
    };
    rule.shared_accounts = this.bip44 ? Array.from(this.shared).sort((a, b) => a - b) : [];
    rule.addresses = this.addresses.split(/\s+/).filter(a => a);
    rule.mappings = this.maps.slice();
    if (this.includeSpent) {
      delete rule.include_spent;
    } else {
      rule.include_spent = false;
    }

    return rule;
  }

  save() {
    if (this.busy) {
      return;
    }
    if (this.destination === 'address' && !this.address.trim()) {
      this.error = 'wallet.sweep.err-destination';

      return;
    }
    const rule = this.buildRule();
    if (sweepActive(rule) && this.starting && !this.burnOk) {
      this.error = 'wallet.sweep.err-burn';

      return;
    }
    this.run(this.backend.setRule(rule), 'wallet.sweep.done-saved', true);
  }

  addMapping() {
    const from = this.mapFrom.trim();
    const dest = choiceDestination(this.mapTo, this.mapTo === 'address' ? this.mapToAddress : '');
    if (!from || (dest.kind === 'address' && !dest.address)) {
      this.error = 'wallet.sweep.err-mapping';

      return;
    }
    const source = /^\d+$/.test(from) && this.kind === 'node' ? { account: Number(from) } : { address: from };
    const keep = this.maps;
    this.maps = [...this.maps, { from: source, to: dest }];
    this.save();
    if (this.error) {
      this.maps = keep;
    } else {
      this.mapFrom = '';
      this.mapToAddress = '';
    }
  }

  removeMapping(i: number) {
    const keep = this.maps;
    this.maps = this.maps.filter((_, j) => j !== i);
    this.save();
    if (this.error) {
      this.maps = keep;
    }
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
    this.run(this.backend.unlock(this.id, password, Number(this.ttl)),
      this.kind === 'browser' ? 'wallet.sweep.browser-done-unlocked' : 'wallet.sweep.done-unlocked', false);
  }

  lock() {
    if (!this.busy) {
      this.run(this.backend.lock(this.id), this.kind === 'browser' ? 'wallet.sweep.browser-done-locked' : 'wallet.sweep.done-locked', false);
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
        if (ruleChanged && this.state) {
          this.maps = (this.state.rule && this.state.rule.mappings) ? this.state.rule.mappings.slice() : [];
        }
        this.changeDetectorRef.markForCheck();
      },
    }));
  }

  /** Takes the answer, and its rule into the form when it changed. */
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
    this.includeSpent = !rule || rule.include_spent !== false;
    this.maps = rule && rule.mappings ? rule.mappings.slice() : [];
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

    return this.bip44 && this.kind === 'node' ? 'next' : 'self';
  }
}

/** The destination a choice key names, with address for a fixed one. */
export function choiceDestination(choice: string, address: string): SweepDestination {
  const i = choice.indexOf(':');
  const [kind, value] = i < 0 ? [choice, ''] : [choice.substring(0, i), choice.substring(i + 1)];
  if (kind === 'next') {
    return { kind: 'account' };
  } else if (kind === 'account') {
    return { kind: 'account', account: Number(value) };
  } else if (kind === 'self') {
    return { kind: 'wallet' };
  } else if (kind === 'wallet') {
    return { kind: 'wallet', wallet: value };
  }

  return { kind: 'address', address: address.trim() };
}
