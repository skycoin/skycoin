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
  /** Ticker/symbol, e.g. SKY / BTC. */
  coinSymbol?: string;
  isBitcoin: boolean;
  available: boolean;
  /** Blockchain head seq (skycoin/fiber) or electrum tip height (bitcoin). */
  height: number;
  /**
   * The block-publisher public key that signs this fiber chain — the canonical
   * identity that distinguishes one fiber coin's node from another's. Empty for
   * bitcoin (electrum has no such notion). From /api/v1/health blockchain_pubkey.
   */
  blockchainPubkey?: string;
  /** Blockchain head block hash (chain-tip identifier). */
  headHash?: string;
  /** Node/daemon version string, when reported. */
  version?: string;
  /** The effective node/electrum URL actually probed (for display). */
  nodeUrl?: string;
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
 *   - skycoin/fiber: GET <nodeUrl>/api/v1/health  → head seq, coin name,
 *                    blockchain_pubkey, head hash, version (chain identity)
 *   - bitcoin:       GET <electrumUrl>/v1/btc/health → tip_height (electrum)
 */
@Injectable()
export class NodeHealthService {
  constructor(private http: HttpClient, private coinService: CoinService) { }

  /** One-shot health probe for a single coin's backend. Never throws. */
  check(coin: BaseCoin): Observable<CoinHealth> {
    if (!coin) {
      return of({ coinId: 0, coinName: '', isBitcoin: false, available: false, height: 0, error: 'no coin' });
    }
    // The EFFECTIVE backend: a per-coin custom URL set in Settings → Nodes
    // (coinService.customNodeUrls) overrides the coin default, exactly as
    // ApiService resolves it for the wallet's own calls — so the status bar
    // reflects the node/electrum the wallet is actually talking to. For bitcoin
    // that URL is the ssl:// electrum server; the visor's BTC gateway reads it
    // straight from the request origin, so prefixing it here keeps health and
    // the balance/history calls pointed at the same electrum.
    const custom = this.coinService.customNodeUrls && this.coinService.customNodeUrls[coin.id.toString()];
    let url = ((custom || coin.nodeUrl) || '').trim();
    if (url.endsWith('/')) {
      url = url.substring(0, url.length - 1);
    }
    const base: CoinHealth = {
      coinId: coin.id, coinName: coin.coinName, coinSymbol: coin.coinSymbol,
      isBitcoin: coin.isBitcoin(), available: false, height: 0, nodeUrl: url,
    };

    if (coin.isBitcoin()) {
      return this.http.get(url + '/v1/btc/health').pipe(
        map((r: any) => ({ ...base, available: true, height: (r && r.tip_height) || 0 })),
        catchError(err => of({ ...base, error: this.msg(err) })),
      );
    }

    return this.http.get(url + '/api/v1/health').pipe(
      map((r: any) => ({
        ...base,
        available: true,
        height: (r && r.blockchain && r.blockchain.head && r.blockchain.head.seq) || 0,
        blockchainPubkey: (r && r.blockchain_pubkey) || '',
        headHash: (r && r.blockchain && r.blockchain.head && r.blockchain.head.block_hash) || '',
        version: (r && r.version && r.version.version) || '',
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
