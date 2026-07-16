import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs';

import { CoinService } from '../../../services/coin.service';
import { NodeHealthService, CoinHealth } from '../../../services/node-health.service';
import { BaseCoin } from '../../../coins/basecoin';

/**
 * NodeStatusBarComponent is a thin, always-visible strip at the very top of the
 * wallet that reports the current coin's backend connection: reachable? which
 * chain height? It renders on every screen — including the seed-entry /
 * create-wallet screen, where there is no wallet to query yet — so the operator
 * can see whether the node/electrum connection is up (and what the chain tip is)
 * without first unlocking a wallet.
 */
@Component({
  selector: 'app-node-status-bar',
  templateUrl: './node-status-bar.component.html',
  styleUrls: ['./node-status-bar.component.scss'],
  standalone: false,
})
export class NodeStatusBarComponent implements OnInit, OnDestroy {
  health: CoinHealth = null;
  private subscription: Subscription;

  constructor(private coinService: CoinService, private nodeHealthService: NodeHealthService) { }

  ngOnInit() {
    this.subscription = this.coinService.currentCoin.pipe(
      switchMap((coin: BaseCoin) => {
        this.health = null; // reset to the "connecting" state on coin switch
        return this.nodeHealthService.watch(coin);
      }),
    ).subscribe(health => this.health = health);
  }

  ngOnDestroy() {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }
}
