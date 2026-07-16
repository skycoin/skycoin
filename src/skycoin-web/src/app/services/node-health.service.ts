import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of, timer } from 'rxjs';
import { catchError, map, switchMap } from 'rxjs';

import { BaseCoin } from '../coins/basecoin';
import { CoinService } from './coin.service';

/**
 * CoinHealth is the connection/chain-status of a single coin's backend, from the
 * wallet's perspective. `height` is the skycoin blockchain head seq or the
 * bitcoin electrum tip height. A probe never throws — failure yields
 * { available: false, error }.
 */
export interface CoinHealth {
  coinId: number;
  coinName: string;
  isBitcoin: boolean;
  available: boolean;
  height: number;
  error?: string;
}

/**
 * NodeHealthService probes whether a coin's backing node (skycoin/fiber) or
 * electrum server (bitcoin) is reachable, and what chain height it reports.
 *
 * This is what lets the wallet show a live connection indicator even at the
 * seed-entry / create-wallet screen (where there is no wallet to query yet) and
 * gate the coin-type selector by which backends are actually available. It uses
 * the same request paths the rest of the wallet does, so the visor's fetch shim
 * routes them over the mesh transparently:
 *   - skycoin/fiber: GET <nodeUrl>/api/v1/health  → blockchain.head.seq
 *   - bitcoin:       GET /v1/btc/health           → tip_height (electrum)
 */
@Injectable()
export class NodeHealthService {
  constructor(private http: HttpClient, private coinService: CoinService) { }

  /** One-shot health probe for a single coin's backend. Never throws. */
  check(coin: BaseCoin): Observable<CoinHealth> {
    if (!coin) {
      return of({ coinId: 0, coinName: '', isBitcoin: false, available: false, height: 0, error: 'no coin' });
    }
    const base: CoinHealth = {
      coinId: coin.id, coinName: coin.coinName, isBitcoin: coin.isBitcoin(), available: false, height: 0,
    };

    if (coin.isBitcoin()) {
      return this.http.get('/v1/btc/health').pipe(
        map((r: any) => ({ ...base, available: true, height: (r && r.tip_height) || 0 })),
        catchError(err => of({ ...base, error: this.msg(err) })),
      );
    }

    // Probe the EFFECTIVE node: a per-coin custom node URL set in Settings →
    // Nodes (coinService.customNodeUrls) overrides the coin default, exactly as
    // ApiService resolves it for the wallet's own calls — so the health bar
    // reflects the node the wallet is actually talking to, not the default.
    const custom = this.coinService.customNodeUrls && this.coinService.customNodeUrls[coin.id.toString()];
    let url = ((custom || coin.nodeUrl) || '').trim();
    if (url.endsWith('/')) {
      url = url.substring(0, url.length - 1);
    }
    return this.http.get(url + '/api/v1/health').pipe(
      map((r: any) => ({
        ...base,
        available: true,
        height: (r && r.blockchain && r.blockchain.head && r.blockchain.head.seq) || 0,
      })),
      catchError(err => of({ ...base, error: this.msg(err) })),
    );
  }

  /**
   * Polls a coin's health every intervalMs (default 20s), emitting immediately
   * then on each tick. Used by the always-visible node-status bar.
   */
  watch(coin: BaseCoin, intervalMs = 20000): Observable<CoinHealth> {
    return timer(0, intervalMs).pipe(switchMap(() => this.check(coin)));
  }

  private msg(err: any): string {
    if (!err) {
      return 'unreachable';
    }
    if (err.status === 0 || err.status === undefined) {
      return 'unreachable';
    }
    if (err.status) {
      return 'HTTP ' + err.status;
    }
    return (err.message || String(err)).substring(0, 120);
  }
}
