import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs';

import { Wallet } from '../../app.datatypes';
import { BaseCoin } from '../../coins/basecoin';
import { ApiService } from '../api.service';

/** Where swept coins go (src/wallet/sweep/rule.go). */
export interface SweepDestination {
  kind: 'account' | 'wallet' | 'address';
  account?: number;
  wallet?: string;
  address?: string;
}

/** The sweep rule of one wallet. */
export interface SweepRule {
  wallet: string;
  shared_accounts?: number[];
  addresses?: string[];
  destination: SweepDestination;
  enabled: boolean;
  gap?: number;
}

/** What the node reports about the sweep of one wallet. */
export interface SweepStatus {
  wallet: string;
  enabled: boolean;
  paused?: string;
  unlocked: boolean;
  unlock_expires?: string;
  head_seq: number;
  exposed_addresses: number;
  waiting_coins: string;
  waiting_hours: number;
  contested: number;
  destination?: string;
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

/** The node's sweep of exposed addresses, under /api/v2/wallet/sweep (WALLET API set). */
@Injectable()
export class SweepService {
  constructor(private apiService: ApiService) { }

  get(walletId: string): Observable<SweepState> {
    return this.apiService.get('wallet/sweep', { id: walletId }, {}, true).pipe(map(unwrap));
  }

  setRule(rule: SweepRule): Observable<SweepState> {
    return this.apiService.post('wallet/sweep', rule, {}, true).pipe(map(unwrap));
  }

  /** Has the node keep the password of an encrypted wallet in memory, for ttlSeconds or, when 0, until it stops. */
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
 * Says the node can sweep wallet: it keeps the wallet on disk with its keys and
 * runs the coin itself. A hardware wallet's keys never reach the node.
 */
export function canAutoSweep(wallet: Wallet, serverWallet: boolean, coin: BaseCoin | null | undefined): boolean {
  return serverWallet && !wallet.isHardware && !wallet.temporary && wallet.walletType !== 'xpub' && !!coin && !coin.isGateway();
}
