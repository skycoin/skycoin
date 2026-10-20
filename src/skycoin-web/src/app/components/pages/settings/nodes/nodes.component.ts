import { Component, OnInit, OnDestroy, ChangeDetectionStrategy, ChangeDetectorRef } from '@angular/core';
import { Subscription } from 'rxjs';
import { MatDialogConfig } from '@angular/material/dialog';

import { CoinService } from '../../../../services/coin.service';
import { BaseCoin } from '../../../../coins/basecoin';
import { ChangeNodeURLComponent } from './change-url/change-node-url.component';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { environment } from '../../../../../environments/environment';

@Component({
    selector: 'app-nodes',
    templateUrl: './nodes.component.html',
    styleUrls: ['./nodes.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class NodesComponent implements OnInit, OnDestroy {
  coins: BaseCoin[] = [];
  customNodeUrls: any = {};
  enableCustomNodes = environment.enableCustomNodes;
  // Ether's coin, whose servers every token shares.
  ethCoin: BaseCoin | undefined;

  private coinsLoadedSub!: Subscription;

  constructor(
    private coinService: CoinService,
    private dialog: CustomMatDialogService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {}

  // The page can open before the coin list has loaded, so it waits for it.
  ngOnInit() {
    this.coinsLoadedSub = this.coinService.coinsLoaded.subscribe(() => {
      this.coins = this.coinService.coins;
      this.customNodeUrls = this.coinService.customNodeUrls;
      this.ethCoin = this.coins.find(c => c.coinType === 'ethereum' && c.serverWallets);
      this.changeDetectorRef.markForCheck();
    });
  }

  ngOnDestroy() {
    this.coinsLoadedSub.unsubscribe();
  }

  changeNodeURL(coin: BaseCoin) {
    const customNodeUrl = this.customNodeUrls[coin.id.toString()] ? this.customNodeUrls[coin.id.toString()] : '';

    const config = new MatDialogConfig();
    config.width = '566px';
    config.data = {coinId: coin.id, url: customNodeUrl};
    this.dialog.open(ChangeNodeURLComponent, config);
  }
}
