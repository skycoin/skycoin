import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { UntypedFormArray, UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { Observable, Subscription, first, map, mergeMap, of } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { Output, Wallet } from '../../../app.datatypes';
import { BaseCoin } from '../../../coins/basecoin';
import { CoinService } from '../../../services/coin.service';
import { ClipboardService } from '../../../services/clipboard.service';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';
import { MsgBarService } from '../../../services/msg-bar.service';
import { BalanceService } from '../../../services/wallet/balance.service';
import { SpendingService } from '../../../services/wallet/spending.service';
import { WalletService } from '../../../services/wallet/wallet.service';
import {
  JointDestination, JointTransaction, JointTransactionService, parseHashes,
} from '../../../services/wallet/joint-transaction.service';
import { DoubleButtonActive } from '../../layout/double-button/double-button.component';
import { PasswordDialogComponent } from '../../layout/password-dialog/password-dialog.component';
import { openUnlockWalletModal } from '../../../utils';

/** A typed amount as a number, or null when it is not one. */
function toNumber(value: any): BigNumber | null {
  if (value === null || value === undefined || String(value).trim() === '') {
    return null;
  }
  try {
    const n = new BigNumber(String(value).trim());

    return n.isFinite() ? n : null;
  } catch {
    return null;
  }
}

/**
 * Builds and co-signs joint transactions. Create spends your outputs and the
 * counterparty's with outputs set by hand; Co-sign reviews a pasted transaction
 * and signs only the inputs of the chosen wallet.
 */
@Component({
    selector: 'app-joint-transaction',
    templateUrl: './joint-transaction.component.html',
    styleUrls: ['./joint-transaction.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class JointTransactionComponent implements OnInit, OnDestroy {
  buttons = DoubleButtonActive;
  mode = DoubleButtonActive.LeftButton;
  currentCoin?: BaseCoin;
  wallets: Wallet[] = [];

  createForm: UntypedFormGroup;
  ownOutputs: Output[] = [];
  loadingOwnOutputs = false;
  theirOutputs: Output[] = [];
  missingHashes: string[] = [];
  created = '';

  signForm: UntypedFormGroup;
  tx?: JointTransaction;

  working = false;
  error = '';
  copied = false;

  private subscriptions: Subscription[] = [];
  private operation?: Subscription;
  private ownOutputsSubscription?: Subscription;
  private theirOutputsSubscription?: Subscription;

  constructor(
    private formBuilder: UntypedFormBuilder,
    private jointService: JointTransactionService,
    private walletService: WalletService,
    private spendingService: SpendingService,
    private balanceService: BalanceService,
    private clipboardService: ClipboardService,
    private msgBarService: MsgBarService,
    private dialog: CustomMatDialogService,
    private route: ActivatedRoute,
    coinService: CoinService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    this.createForm = formBuilder.group({
      wallet: [null],
      outputs: [[]],
      theirs: [''],
      destinations: formBuilder.array([this.destinationGroup(), this.destinationGroup()]),
      changeAddress: [''],
    });
    this.signForm = formBuilder.group({
      wallet: [null],
      rawTx: [''],
    });

    this.subscriptions.push(coinService.currentCoin.subscribe(coin => {
      this.currentCoin = coin;
      this.changeDetectorRef.markForCheck();
    }));
  }

  ngOnInit() {
    this.subscriptions.push(this.walletService.currentWallets.subscribe(wallets => {
      this.wallets = wallets.filter(w => !w.isHardware);
      this.changeDetectorRef.markForCheck();
    }));

    this.subscriptions.push(this.createForm.get('wallet')!.valueChanges.subscribe(wallet => {
      this.loadOwnOutputs(wallet);
      this.changeDetectorRef.markForCheck();
    }));
    this.subscriptions.push(this.signForm.get('wallet')!.valueChanges.subscribe(() => {
      if (this.tx) {
        this.review();
      }
      this.changeDetectorRef.markForCheck();
    }));

    const raw = this.route.snapshot.queryParamMap.get('tx');
    if (raw) {
      this.mode = DoubleButtonActive.RightButton;
      this.signForm.get('rawTx')!.setValue(raw);
    }
  }

  ngOnDestroy() {
    this.subscriptions.forEach(s => s.unsubscribe());
    this.operation?.unsubscribe();
    this.ownOutputsSubscription?.unsubscribe();
    this.theirOutputsSubscription?.unsubscribe();
  }

  get destinations(): UntypedFormArray {
    return this.createForm.get('destinations') as UntypedFormArray;
  }

  changeMode(mode: DoubleButtonActive) {
    this.mode = mode;
    this.error = '';
    this.copied = false;
  }

  addDestination() {
    this.destinations.push(this.destinationGroup());
  }

  removeDestination(i: number) {
    this.destinations.removeAt(i);
  }

  outputCompare(a: Output, b: Output): boolean {
    return !!a && !!b && a.hash === b.hash;
  }

  /** Looks up the counterparty's pasted output hashes on the node. */
  lookupTheirOutputs() {
    const hashes = parseHashes(this.createForm.value.theirs);
    this.theirOutputsSubscription?.unsubscribe();
    this.theirOutputsSubscription = this.jointService.outputs(hashes).subscribe(outputs => {
      this.theirOutputs = outputs;
      this.missingHashes = hashes.filter(h => !outputs.some(o => o.hash === h));
      this.changeDetectorRef.markForCheck();
    }, err => {
      this.error = this.message(err);
      this.changeDetectorRef.markForCheck();
    });
  }

  get inputs(): Output[] {
    return [...(this.createForm.value.outputs || []), ...this.theirOutputs];
  }

  get inputCoins(): BigNumber {
    return this.inputs.reduce((sum, o) => sum.plus(o.coins), new BigNumber(0));
  }

  get inputHours(): BigNumber {
    return this.inputs.reduce((sum, o) => sum.plus(o.calculated_hours), new BigNumber(0));
  }

  /** The outputs as typed, or null while one is incomplete or invalid. */
  get parsedDestinations(): JointDestination[] | null {
    const parsed: JointDestination[] = [];
    for (const d of this.destinations.value) {
      const coins = toNumber(d.coins);
      const hours = toNumber(d.hours);
      if (!d.address || !d.address.trim() || !coins || coins.isLessThanOrEqualTo(0)
        || !hours || !hours.isInteger() || hours.isNegative()) {
        return null;
      }
      parsed.push({ address: d.address.trim(), coins, hours });
    }

    return parsed;
  }

  get outputCoins(): BigNumber {
    return (this.parsedDestinations || []).reduce((sum, d) => sum.plus(d.coins), new BigNumber(0));
  }

  get outputHours(): BigNumber {
    return (this.parsedDestinations || []).reduce((sum, d) => sum.plus(d.hours), new BigNumber(0));
  }

  get minimumFee(): BigNumber {
    return this.jointService.minimumFee(this.inputHours);
  }

  get fee(): BigNumber {
    return this.inputHours.minus(this.outputHours);
  }

  get leftoverCoins(): BigNumber {
    return this.inputCoins.minus(this.outputCoins);
  }

  /** Why the transaction cannot be created yet, as a translation key, or ''. */
  get createProblem(): string {
    if (this.inputs.length === 0) {
      return 'joint-tx.need-inputs';
    }
    if (this.missingHashes.length > 0) {
      return 'joint-tx.unknown-outputs';
    }
    if (!this.parsedDestinations || this.destinations.length === 0) {
      return 'joint-tx.need-outputs';
    }
    if (this.leftoverCoins.isNegative()) {
      return 'joint-tx.not-enough-coins';
    }
    if (this.fee.isLessThan(this.minimumFee)) {
      return 'joint-tx.fee-too-low';
    }
    if (this.leftoverCoins.isGreaterThan(0) && !this.createForm.value.changeAddress.trim()) {
      return 'joint-tx.need-change';
    }

    return '';
  }

  create() {
    const unspents = this.inputs.map(o => o.hash);
    const change = this.createForm.value.changeAddress.trim() || null;
    this.run(this.jointService.create(unspents, this.parsedDestinations!, change), encoded => {
      this.created = encoded;
    });
  }

  /** Opens the created transaction in Co-sign, to review it and sign your inputs. */
  cosignCreated() {
    this.signForm.get('rawTx')!.setValue(this.created);
    if (!this.signForm.value.wallet && this.createForm.value.wallet) {
      this.signForm.get('wallet')!.setValue(this.createForm.value.wallet, { emitEvent: false });
    }
    this.changeMode(DoubleButtonActive.RightButton);
    this.review();
  }

  review() {
    const raw = (this.signForm.value.rawTx || '').trim();
    this.run(this.jointService.decode(raw, this.signForm.value.wallet), tx => {
      this.tx = tx;
    }, () => {
      this.tx = undefined;
    });
  }

  get unsignedOwnInputs(): number {
    return this.tx ? this.tx.inputs.filter(i => i.mine && !i.signed).length : 0;
  }

  sign() {
    const wallet: Wallet = this.signForm.value.wallet;
    const tx = this.tx!;
    this.run(this.unlock(wallet).pipe(
      mergeMap(password => password === undefined ? of(undefined) :
        this.jointService.sign(wallet, password || undefined, tx).pipe(
          mergeMap(signed => this.jointService.decode(signed, wallet)))),
    ), signed => {
      this.tx = signed;
      this.signForm.get('rawTx')!.setValue(signed.encoded, { emitEvent: false });
      this.copied = false;
    });
  }

  broadcast() {
    this.run(this.spendingService.injectTransaction(this.tx!.encoded), () => {
      this.balanceService.startGettingBalances();
      this.msgBarService.showDone('joint-tx.broadcast-done');
      this.tx = undefined;
      this.signForm.get('rawTx')!.setValue('');
    });
  }

  copy(text: string) {
    this.clipboardService.copy(text).then(() => {
      this.copied = true;
      this.changeDetectorRef.markForCheck();
    });
  }

  private destinationGroup(): UntypedFormGroup {
    return this.formBuilder.group({ address: [''], coins: [''], hours: ['0'] });
  }

  private loadOwnOutputs(wallet: Wallet | null) {
    this.ownOutputs = [];
    this.createForm.get('outputs')!.setValue([], { emitEvent: false });
    this.ownOutputsSubscription?.unsubscribe();
    if (!wallet) {
      return;
    }
    this.loadingOwnOutputs = true;
    this.ownOutputsSubscription = this.spendingService.getWalletUnspentOutputs(wallet).pipe(first()).subscribe(outputs => {
      this.ownOutputs = outputs;
      this.loadingOwnOutputs = false;
      this.changeDetectorRef.markForCheck();
    }, err => {
      this.loadingOwnOutputs = false;
      this.error = this.message(err);
      this.changeDetectorRef.markForCheck();
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
    const unlockedSeed: Observable<null | undefined> = wallet.seed ? of(null) : this.askSeed(wallet);

    return unlockedSeed.pipe(mergeMap(r => r === undefined ? of(undefined) :
      this.walletService.ensureKeys(wallet).pipe(map(() => null))));
  }

  private askSeed(wallet: Wallet): Observable<null | undefined> {
    const unlocked = openUnlockWalletModal(wallet, this.dialog);

    return new Observable<null | undefined>(subscriber => {
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

  private run<T>(operation: Observable<T | undefined>, onDone: (result: T) => void, onError?: () => void) {
    this.working = true;
    this.error = '';
    this.operation?.unsubscribe();
    this.operation = operation.subscribe(result => {
      this.working = false;
      if (result !== undefined) {
        onDone(result);
      }
      this.changeDetectorRef.markForCheck();
    }, err => {
      this.working = false;
      this.error = this.message(err);
      if (onError) {
        onError();
      }
      this.changeDetectorRef.markForCheck();
    });
  }

  private message(err: any): string {
    return err && err.message ? err.message : String(err);
  }
}
