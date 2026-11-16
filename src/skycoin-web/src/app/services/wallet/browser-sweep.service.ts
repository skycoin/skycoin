import { EventEmitter, Injectable, OnDestroy } from '@angular/core';
import { BigNumber } from 'bignumber.js';
import { Observable, Subscription, firstValueFrom, from, interval } from 'rxjs';

import { Address, Wallet } from '../../app.datatypes';
import { ApiService } from '../api.service';
import { CipherProvider } from '../cipher.provider';
import { isBip44 } from './backends/browser-wallet-backend';
import { SweepBackend, SweepDestination, SweepRule, SweepState, SweepStatus, sweepActive } from './sweep.service';
import { WalletService, walletAddresses } from './wallet.service';

/** Where the rules of browser wallets are kept, shared with ui2's bw.js. They hold no secret. */
export const browserSweepRulesKey = 'sweepRules';

/** How often the head block is checked. */
export const browserSweepTickMs = 5000;

/** The id a browser wallet's rule is kept under: its first address, as ui2 names it. */
export function browserWalletId(wallet: Wallet): string {
  if (isBip44(wallet)) {
    const ext = wallet.accounts![0].externalAddresses;

    return ext.length ? ext[0].address : '';
  }

  return wallet.addresses.length ? wallet.addresses[0].address : '';
}

interface Session {
  password: string;
  /** When the password is forgotten, 0 for when the page is left. */
  expires: number;
}

interface Pending {
  txid: string;
  inputs: string[];
}

interface Local {
  pending: Pending[];
  confirmed: number;
  last?: string;
  txids: string[];
  dests: string[];
  error: string;
  paused: string;
  exposed: number;
  uncovered: number;
  coins: number;
  hours: number;
  contested: number;
  lastSeq: number;
  retryAt: number;
}

interface Chain {
  outputs: any[];
  outgoing: { [hash: string]: boolean };
  sent: { [address: string]: boolean };
  seen: { [address: string]: boolean };
}

interface Fresh {
  address: string;
  index: number;
  wallet: Wallet;
  own: boolean;
  derived: boolean;
}

/**
 * Sweeps the exposed addresses of wallets kept in this browser, which the node
 * never sees. It runs in this page: on every block while the page is open, and
 * for an encrypted wallet only while its password is held here, in memory, for
 * the time asked. Each pass runs the node's rule code in the wasm cipher
 * (sweepPlan), signs with prepareTransaction as a send does, and hands the
 * transactions to the node.
 */
@Injectable()
export class BrowserSweepService implements SweepBackend, OnDestroy {
  private sessions = new Map<string, Session>();
  private locals = new Map<string, Local>();
  private busy = new Set<string>();
  private ticker: Subscription | null = null;

  constructor(
    private walletService: WalletService,
    private cipherProvider: CipherProvider,
    private apiService: ApiService,
  ) { }

  /** Starts sweeping on every block. The wallet shell calls it once. */
  start() {
    if (!this.ticker) {
      this.ticker = interval(browserSweepTickMs).subscribe(() => this.tick());
      this.tick();
    }
  }

  ngOnDestroy() {
    if (this.ticker) {
      this.ticker.unsubscribe();
    }
    this.sessions.clear();
  }

  get(walletId: string): Observable<SweepState> {
    return from(Promise.resolve(this.state(walletId)));
  }

  setRule(rule: SweepRule): Observable<SweepState> {
    return from(this.saveRule(rule));
  }

  unlock(walletId: string, password: string, ttlSeconds: number): Observable<SweepState> {
    return from(this.doUnlock(walletId, password, ttlSeconds));
  }

  lock(walletId: string): Observable<SweepState> {
    this.sessions.delete(walletId);

    return from(Promise.resolve(this.state(walletId)));
  }

  /** The stored rules of every browser wallet. */
  rules(): { [id: string]: SweepRule } {
    try {
      return JSON.parse(localStorage.getItem(browserSweepRulesKey) || '{}') || {};
    } catch (e) {
      return {};
    }
  }

  private local(id: string): Local {
    let l = this.locals.get(id);
    if (!l) {
      l = { pending: [], confirmed: 0, txids: [], dests: [], error: '', paused: '', exposed: 0, uncovered: 0, coins: 0, hours: 0,
        contested: 0, lastSeq: -1, retryAt: 0 };
      this.locals.set(id, l);
    }

    return l;
  }

  private session(id: string): Session | null {
    const s = this.sessions.get(id);
    if (s && s.expires && Date.now() >= s.expires) {
      this.sessions.delete(id);

      return null;
    }

    return s || null;
  }

  private browserWallets(): Wallet[] {
    return (this.walletService.wallets.getValue() || []).filter(w =>
      !this.walletService.isServerWallet(w) && !w.isHardware && !w.temporary && w.walletType !== 'xpub');
  }

  private wallet(id: string): Wallet {
    const w = this.browserWallets().find(x => browserWalletId(x) === id);
    if (!w) {
      throw new Error('wallet not found');
    }

    return w;
  }

  private state(id: string): SweepState {
    const rule = this.rules()[id] || null;
    const l = this.local(id);
    const s = this.session(id);
    const active = sweepActive(rule);
    const status: SweepStatus = {
      wallet: id,
      enabled: !!rule && rule.enabled,
      active: active,
      paused: active ? l.paused || undefined : 'disabled',
      unlocked: !!s,
      unlock_expires: s && s.expires ? new Date(s.expires).toISOString() : undefined,
      head_seq: l.lastSeq,
      exposed_addresses: l.exposed,
      uncovered_addresses: l.uncovered,
      waiting_coins: new BigNumber(l.coins).dividedBy(1e6).toFixed(6),
      waiting_hours: l.hours,
      contested: l.contested,
      destinations: l.dests.length ? l.dests : undefined,
      pending_txids: l.pending.length ? l.pending.map(p => p.txid) : undefined,
      last_sweep: l.last,
      last_txids: l.txids.length ? l.txids : undefined,
      confirmed_sweeps: l.confirmed,
      last_error: l.error || undefined,
    };

    return { rule: rule, status: status };
  }

  private async saveRule(rule: SweepRule): Promise<SweepState> {
    const w = this.wallet(rule.wallet);
    const dests = [rule.destination].concat((rule.mappings || []).map(m => m.to));
    if (dests.some(d => d.kind === 'account') || (rule.mappings || []).some(m => m.from.account !== undefined && m.from.account !== null)) {
      throw new Error('wallet.sweep.err-browser-account');
    }
    await firstValueFrom(this.cipherProvider.initialize());
    const addrs = sweepAddrs(w);
    const fixed = dests.filter(d => d.kind === 'address').map(d => d.address!);
    const chain = await this.chainOf(addrs.map(a => a.address).concat(fixed));
    await firstValueFrom(this.cipherProvider.sweepCheck({ rule: rule, addresses: addrs, sent: Object.keys(chain.sent), accounts: isBip44(w) ? 1 : 0 }));
    const all = this.rules();
    all[rule.wallet] = rule;
    localStorage.setItem(browserSweepRulesKey, JSON.stringify(all));
    const l = this.local(rule.wallet);
    l.lastSeq = -1;
    await this.pass(rule.wallet);

    return this.state(rule.wallet);
  }

  private async doUnlock(id: string, password: string, ttlSeconds: number): Promise<SweepState> {
    const w = this.wallet(id);
    await this.keysOf(w, { password: password, expires: 0 });
    this.sessions.set(id, { password: password, expires: ttlSeconds ? Date.now() + ttlSeconds * 1000 : 0 });
    this.local(id).lastSeq = -1;
    await this.pass(id);

    return this.state(id);
  }

  /** Runs a pass for every browser wallet with an active rule on a new block, or after a failure. */
  private tick() {
    const rules = this.rules();
    const ids = Object.keys(rules).filter(id => sweepActive(rules[id]) && this.browserWallets().some(w => browserWalletId(w) === id));
    if (!ids.length) {
      return;
    }
    this.apiService.get('blockchain/metadata').subscribe({
      next: (m: any) => ids.forEach(id => {
        const l = this.local(id);
        if (m.head.seq !== l.lastSeq || (l.error && Date.now() >= l.retryAt)) {
          l.lastSeq = m.head.seq;
          this.pass(id);
        }
      }),
      error: () => undefined,
    });
  }

  private async chainOf(addrs: string[]): Promise<Chain> {
    if (!addrs.length) {
      return { outputs: [], outgoing: {}, sent: {}, seen: {} };
    }
    const q = addrs.join(',');
    const [outs, txns] = await Promise.all([
      firstValueFrom(this.apiService.post('outputs', { addrs: q })),
      firstValueFrom(this.apiService.post('transactions', { addrs: q, verbose: 1, confirmed: 1 })),
    ]) as [any, any[]];
    const want = new Set(addrs);
    const chain: Chain = { outputs: outs.head_outputs || [], outgoing: {}, sent: {}, seen: {} };
    (txns || []).forEach(t => {
      (t.txn.inputs || []).forEach((i: any) => {
        if (want.has(i.owner)) {
          chain.sent[i.owner] = true;
          chain.seen[i.owner] = true;
        }
      });
      (t.txn.outputs || []).forEach((o: any) => {
        if (want.has(o.dst)) {
          chain.seen[o.dst] = true;
        }
      });
    });
    (outs.outgoing_outputs || []).forEach((o: any) => chain.outgoing[o.hash] = true);

    return chain;
  }

  /**
   * The secret keys of w's addresses, derived into a copy of it so the app's
   * own unlocked state is left as it was.
   */
  private async keysOf(w: Wallet, s: Session | null): Promise<Wallet> {
    const copy = copyWallet(w);
    if (w.encryptedSeed) {
      if (!s) {
        throw new Error('wallet.sweep.err-password');
      }
      await firstValueFrom(this.walletService.unlockWalletWithPassword(copy, s.password, new EventEmitter<number>()));
    } else if (w.seed) {
      await firstValueFrom(this.walletService.unlockWallet(copy, w.seed, new EventEmitter<number>()));
    } else if (!walletAddresses(w).every(a => a.secret_key)) {
      throw new Error('wallet.sweep.err-password');
    }

    return copy;
  }

  /**
   * The first receive address of target after the last one seen, or the next
   * one to derive. A deterministic wallet derives it from keys, which keyed
   * holds when target is the swept wallet.
   */
  private async nextUnused(target: Wallet, own: boolean, keyed: Wallet | null, seen: { [a: string]: boolean }): Promise<Fresh> {
    const stored = isBip44(target) ? target.accounts![0].externalAddresses : target.addresses;
    let last = -1;
    stored.forEach((a, i) => {
      if (seen[a.address]) {
        last = i;
      }
    });
    if (last + 1 < stored.length) {
      return { address: stored[last + 1].address, index: last + 1, wallet: target, own: own, derived: false };
    }
    if (isBip44(target)) {
      const d = await firstValueFrom(this.cipherProvider.bip44XPubAddresses(target.accounts![0].accountXpubKey!, 0, stored.length, 1));

      return { address: d[0].address, index: stored.length, wallet: target, own: own, derived: true };
    }
    const nextSeed = target.nextSeed || (keyed ? keyed.nextSeed : undefined);
    if (!nextSeed) {
      throw new Error('wallet.sweep.err-dest-full');
    }
    const g = await firstValueFrom(this.cipherProvider.generateAddress(nextSeed));

    return { address: g.address.address, index: stored.length, wallet: target, own: own, derived: true };
  }

  /** Keeps a derived destination in its wallet before anything is sent to it. */
  private async storeFresh(f: Fresh) {
    if (isBip44(f.wallet)) {
      if (f.wallet.accounts![0].externalAddresses.length === f.index) {
        await firstValueFrom(this.walletService.addAddress(f.wallet, true, 0, 0, 1));
      }
    } else if (f.wallet.addresses.length === f.index) {
      if (f.wallet.seed && f.wallet.nextSeed) {
        await firstValueFrom(this.walletService.addAddress(f.wallet));
      } else {
        // The app has not unlocked it, so it has no next seed in memory. An
        // unlock derives the keys of every stored address again.
        f.wallet.addresses.push({ address: f.address } as Address);
        this.walletService.saveWallets();
      }
    }
    if (!walletAddresses(f.wallet).some(a => a.address === f.address)) {
      throw new Error('the destination could not be added to its wallet');
    }
  }

  private async destinationOf(rule: SweepRule, t: number, w: Wallet, keyed: Wallet | null, seen: { [a: string]: boolean }): Promise<Fresh | string> {
    const d: SweepDestination = t === -1 ? rule.destination : rule.mappings![t].to;
    if (d.kind === 'address') {
      return d.address!;
    }
    if (d.kind !== 'wallet') {
      throw new Error('wallet.sweep.err-browser-account');
    }
    if (!d.wallet || d.wallet === rule.wallet) {
      return this.nextUnused(w, true, keyed, seen);
    }
    const other = this.wallet(d.wallet);
    const chain = await this.chainOf(walletAddresses(other).map(a => a.address));

    return this.nextUnused(other, false, null, chain.seen);
  }

  /** One sweep of browser wallet id: everything waiting on its exposed addresses goes out now. */
  async pass(id: string): Promise<void> {
    const rule = this.rules()[id];
    const l = this.local(id);
    if (!sweepActive(rule) || this.busy.has(id)) {
      return;
    }
    this.busy.add(id);
    l.error = '';
    try {
      await this.passOnce(id, rule, l);
    } catch (err: any) {
      l.error = err && err.message ? err.message : String(err);
      l.retryAt = Date.now() + 10000;
    } finally {
      this.busy.delete(id);
    }
  }

  private async passOnce(id: string, rule: SweepRule, l: Local): Promise<void> {
    await firstValueFrom(this.cipherProvider.initialize());
    const w = this.wallet(id);
    const addrs = sweepAddrs(w);
    const chain = await this.chainOf(addrs.map(a => a.address));
    const head = new Set(chain.outputs.map(o => o.hash));
    let covered: string[] = [];
    l.pending = l.pending.filter(p => {
      if (p.inputs.every(h => chain.outgoing[h])) {
        covered = covered.concat(p.inputs);

        return true;
      }
      if (p.inputs.every(h => !head.has(h))) {
        l.confirmed++;
      }

      return false;
    });
    const req: any = {
      rule: rule, addresses: addrs, sent: Object.keys(chain.sent), seen: Object.keys(chain.seen), covered: covered,
      accounts: isBip44(w) ? 1 : 0, destinations: {},
      outputs: chain.outputs.map(o => ({
        hash: o.hash, address: o.address, coins: new BigNumber(o.coins).multipliedBy(1e6).toNumber(), hours: o.calculated_hours,
        contested: !!chain.outgoing[o.hash] && covered.indexOf(o.hash) < 0,
      })),
    };
    let res = await firstValueFrom(this.cipherProvider.sweepPlan(req));
    const gap = res.gap ? res.gap['0'] : 0;
    if (gap && isBip44(w)) {
      await firstValueFrom(this.walletService.addAddress(w, true, 0, 0, gap));

      return this.passOnce(id, rule, l);
    }
    this.record(l, res);
    l.paused = '';
    if (!res.need) {
      return;
    }
    const s = this.session(id);
    if (w.encryptedSeed && !s) {
      l.paused = 'locked';

      return;
    }
    const keyed = await this.keysOf(w, s);
    const fresh: Fresh[] = [];
    for (const t of res.need as number[]) {
      const d = await this.destinationOf(rule, t, w, keyed, chain.seen);
      if (typeof d === 'string') {
        req.destinations[String(t)] = d;
      } else {
        req.destinations[String(t)] = d.address;
        if (d.derived) {
          fresh.push(d);
        }
        if (d.own && d.derived) {
          req.addresses.push({ address: d.address, child: d.index, bip44: isBip44(w) });
        }
      }
    }
    const outside = Object.keys(req.destinations).map(k => req.destinations[k]).filter((a: string) => !addrs.some(x => x.address === a));
    const destChain = await this.chainOf(outside);
    req.sent = req.sent.concat(Object.keys(destChain.sent));
    res = await firstValueFrom(this.cipherProvider.sweepPlan(req));
    this.record(l, res);
    if (res.errors && res.errors.length) {
      l.error = res.errors.join('\n');
    }
    if (res.no_hours && !(res.sweeps && res.sweeps.length)) {
      l.paused = 'no-hours';
    }
    if (!res.sweeps || !res.sweeps.length) {
      return;
    }

    const keys = new Map(walletAddresses(keyed).map(a => [a.address, a.secret_key]));
    const signed: { sw: any, raw: string }[] = [];
    for (const sw of res.sweeps) {
      const inputs = sw.inputs.map((i: any) => {
        const secret = keys.get(i.address);
        if (!secret) {
          throw new Error('address ' + i.address + ' is not in the wallet');
        }

        return { hash: i.hash, address: i.address, secret: secret };
      });
      signed.push({ sw: sw, raw: await firstValueFrom(this.cipherProvider.prepareTransaction(inputs, sw.outputs)) });
    }
    const used = new Set<string>(res.sweeps.map((sw: any) => sw.destination));
    for (const f of fresh.filter(x => used.has(x.address))) {
      await this.storeFresh(f);
    }
    const errs: string[] = [];
    const txids: string[] = [];
    for (const x of signed) {
      try {
        const txid: string = await firstValueFrom(this.apiService.post('injectTransaction', { rawtx: x.raw }, { json: true }));
        txids.push(txid);
        l.pending.push({ txid: txid, inputs: x.sw.inputs.map((i: any) => i.hash) });
        l.coins = Math.max(0, l.coins - x.sw.outputs[0].coins);
        l.hours = Math.max(0, l.hours - x.sw.fee);
      } catch (err: any) {
        errs.push(err && err.message ? err.message : String(err));
      }
    }
    if (txids.length) {
      l.last = new Date().toISOString();
      l.txids = txids;
    }
    l.dests = Array.from(used);
    if (errs.length) {
      throw new Error(errs.join('\n'));
    }
  }

  private record(l: Local, res: any) {
    l.exposed = res.exposed || 0;
    l.uncovered = res.uncovered || 0;
    l.coins = res.coins || 0;
    l.hours = res.hours || 0;
    l.contested = res.contested || 0;
  }
}

/** w's addresses with their place, as sweepPlan takes them. */
export function sweepAddrs(w: Wallet): { address: string, change?: boolean, child?: number, bip44?: boolean }[] {
  if (!isBip44(w)) {
    return w.addresses.map(a => ({ address: a.address }));
  }
  const account = w.accounts![0];

  return account.externalAddresses.map((a, i) => ({ address: a.address, child: i, bip44: true }))
    .concat(account.changeAddresses.map((a, i) => ({ address: a.address, change: true, child: i, bip44: true })) as any[]);
}

/** A copy of w whose addresses can take keys without touching w. */
function copyWallet(w: Wallet): Wallet {
  const copy: Wallet = { ...w, addresses: w.addresses.map(a => ({ ...a })) };
  if (w.accounts) {
    copy.accounts = w.accounts.map(acc => ({
      ...acc,
      externalAddresses: acc.externalAddresses.map(a => ({ ...a })),
      changeAddresses: acc.changeAddresses.map(a => ({ ...a })),
    }));
    if (isBip44(copy)) {
      copy.addresses = copy.accounts[0].externalAddresses.concat(copy.accounts[0].changeAddresses);
    }
  }

  return copy;
}
