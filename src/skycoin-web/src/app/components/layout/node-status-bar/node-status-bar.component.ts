import { Component, OnDestroy, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { switchMap } from 'rxjs';

import { CoinService } from '../../../services/coin.service';
import { NodeHealthService, CoinHealth } from '../../../services/node-health.service';
import { BaseCoin } from '../../../coins/basecoin';

/**
 * NodeStatusBarComponent is a thin, always-visible strip at the very top of the
 * wallet that reports the current coin's backend connection: reachable? which
 * chain? which height? It renders on every screen — including the seed-entry /
 * create-wallet screen, where there is no wallet to query yet — so the operator
 * can see whether the node/electrum connection is up (and which chain it is)
 * without first unlocking a wallet.
 *
 * It surfaces the coin's identity (name/ticker) and, for fiber coins, the
 * blockchain_pubkey + head hash — so that when several fiber-coin nodes could be
 * configured you can tell WHICH chain the current node serves. The whole strip
 * is a link to Settings → Node, so the node can be inspected/switched right from
 * here — even during onboarding, before any wallet exists.
 */
@Component({
  selector: 'app-node-status-bar',
  templateUrl: './node-status-bar.component.html',
  styleUrls: ['./node-status-bar.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  standalone: false,
})
export class NodeStatusBarComponent implements OnInit, OnDestroy {
  health: CoinHealth | null = null;
  private subscription!: Subscription;

  constructor(
    private coinService: CoinService,
    private nodeHealthService: NodeHealthService,
    private router: Router,
    private translate: TranslateService,
  ) { }

  ngOnInit() {
    this.subscription = this.coinService.currentCoin.pipe(
      switchMap((coin: BaseCoin) => {
        this.health = null; // reset to the "connecting" state on coin switch
        return this.nodeHealthService.watch(coin);
      }),
    ).subscribe(health => this.health = health);
  }

  /** Short form of the blockchain pubkey / head hash for the inline strip. */
  short(hex: string): string {
    if (!hex) {
      return '';
    }
    return hex.length > 14 ? hex.substr(0, 6) + '…' + hex.substr(hex.length - 5) : hex;
  }

  openNodeSettings() {
    this.router.navigate(['/settings/node']);
  }

  /**
   * Full node identity for the strip's hover tooltip: the node/electrum URL, the
   * blockchain pubkey + head hash (fiber only), and the version. Labels are
   * i18n'd; the whole thing is a plain multi-line title string.
   */
  tooltip(): string {
    const h = this.health;
    if (!h || !h.available) {
      return '';
    }
    const t = (k: string) => this.translate.instant(k);
    const lines: string[] = [];
    if (h.nodeUrl) {
      lines.push(t('nodeStatus.tooltip.url') + ' ' + h.nodeUrl);
    }
    if (h.blockchainPubkey) {
      lines.push(t('nodeStatus.tooltip.pubkey') + ' ' + h.blockchainPubkey);
    }
    if (h.headHash) {
      lines.push(t('nodeStatus.tooltip.head') + ' ' + h.headHash);
    }
    if (h.version) {
      lines.push(t('nodeStatus.tooltip.version') + ' ' + h.version);
    }
    lines.push(t('nodeStatus.tooltip.change'));
    return lines.join('\n');
  }

  ngOnDestroy() {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }
}
