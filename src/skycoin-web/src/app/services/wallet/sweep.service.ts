import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs';

import { Wallet } from '../../app.datatypes';
import { BaseCoin } from '../../coins/basecoin';
import { ApiService } from '../api.service';

/** Where swept coins go (src/wallet/sweep/plan/rule.go). */
export interface SweepDestination {
  kind: 'account' | 'wallet' | 'address';
  account?: number;
  wallet?: string;
  address?: string;
}

/** What a mapping sweeps from: one address, or the exposed addresses of one account. */
export interface SweepSource {
  address?: string;
  account?: number;
}

/** A level 2 mapping: its source goes to its own destination. */
export interface SweepMapping {
  from: SweepSource;
  to: SweepDestination;
}

/**
 * The sweep rule of one wallet. enabled is level 1, every exposed address to
 * destination. Mappings sweep whether it is on or not.
 */
export interface SweepRule {
  wallet: string;
  enabled: boolean;
  destination: SweepDestination;
  shared_accounts?: number[];
  addresses?: string[];
  include_spent?: boolean;
  mappings?: SweepMapping[];
  gap?: number;
}

/** What is known about the sweep of one wallet. */
export interface SweepStatus {
  wallet: string;
  enabled: boolean;
  active: boolean;
  paused?: string;
  unlocked: boolean;
  unlock_expires?: string;
  head_seq: number;
  exposed_addresses: number;
  uncovered_addresses: number;
  waiting_coins: string;
  waiting_hours: number;
  contested: number;
  destinations?: string[];
  pending_txids?: string[];
  last_sweep?: string;
  last_txids?: string[];
  confirmed_sweeps: number;
  displaced_txids?: string[];
  last_error?: string;
}

/** The rule of a wallet, null when it has none, and its status. */
export interface SweepState {
  rule: SweepRule | null;
  status: SweepStatus;
}

/** Who sweeps a wallet: the node for one it keeps, this page for one kept in the browser. */
export interface SweepBackend {
  get(walletId: string): Observable<SweepState>;
  setRule(rule: SweepRule): Observable<SweepState>;
  /** Keeps the password for ttlSeconds or, when 0, until the node stops or the page is left. */
  unlock(walletId: string, password: string, ttlSeconds: number): Observable<SweepState>;
  lock(walletId: string): Observable<SweepState>;
}

/** Says whether the rule sweeps anything, through level 1 or a mapping. */
export function sweepActive(rule: SweepRule | null | undefined): boolean {
  return !!rule && (rule.enabled || !!(rule.mappings && rule.mappings.length));
}

/** The node's sweep of exposed addresses, under /api/v2/wallet/sweep (WALLET API set). */
@Injectable()
export class SweepService implements SweepBackend {
  constructor(private apiService: ApiService) { }

  get(walletId: string): Observable<SweepState> {
    return this.apiService.get('wallet/sweep', { id: walletId }, {}, true).pipe(map(unwrap));
  }

  setRule(rule: SweepRule): Observable<SweepState> {
    return this.apiService.post('wallet/sweep', rule, {}, true).pipe(map(unwrap));
  }

  unlock(walletId: string, password: string, ttlSeconds: number): Observable<SweepState> {
    return this.apiService.post('wallet/sweep/unlock', { id: walletId, password: password, ttl_seconds: ttlSeconds }, {}, true)
      .pipe(map(unwrap));
  }

  lock(walletId: string): Observable<SweepState> {
    return this.apiService.post('wallet/sweep/lock', { id: walletId }, {}, true).pipe(map(unwrap));
  }
}

function unwrap(response: any): SweepState {
  return response && response.data ? response.data : response;
}

/**
 * Who can sweep wallet: 'node' for one the node keeps, 'browser' for one kept
 * in this browser, which this page sweeps, 'hardware' for a hardware wallet,
 * whose device must confirm every signature so it is never swept, and null
 * when none applies.
 */
export function sweepKind(wallet: Wallet, serverWallet: boolean, coin: BaseCoin | null | undefined): 'node' | 'browser' | 'hardware' | null {
  if (!coin || coin.isGateway() || wallet.temporary || wallet.walletType === 'xpub') {
    return null;
  }
  if (wallet.isHardware) {
    return 'hardware';
  }

  return serverWallet ? 'node' : 'browser';
}

/** Says the panel applies to wallet. A hardware wallet gets it too, to say why it is not swept. */
export function canAutoSweep(wallet: Wallet, serverWallet: boolean, coin: BaseCoin | null | undefined): boolean {
  return sweepKind(wallet, serverWallet, coin) !== null;
}
