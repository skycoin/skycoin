import { Component, Input, OnInit, OnDestroy, ChangeDetectionStrategy, ChangeDetectorRef } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Subscription } from 'rxjs';

import { BaseCoin } from '../../../../../coins/basecoin';
import { CoinService } from '../../../../../services/coin.service';
import { MsgBarService } from '../../../../../services/msg-bar.service';

/**
 * Edits the Ethereum RPC servers and the history indexer the wallet service
 * uses for Ether and every token. The service keeps them beside its wallets.
 */
@Component({
    selector: 'app-eth-endpoints',
    templateUrl: './eth-endpoints.component.html',
    styleUrls: ['./eth-endpoints.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class EthEndpointsComponent implements OnInit, OnDestroy {
  @Input() coin!: BaseCoin;

  rpcText = '';
  indexer = '';
  defaults: { rpc_urls: string[], indexer_url: string } | null = null;
  busy = false;

  private subscription: Subscription | null = null;

  constructor(
    private http: HttpClient,
    private coinService: CoinService,
    private msgBarService: MsgBarService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {}

  ngOnInit() {
    this.subscription = this.http.get(this.endpointsUrl).subscribe(
      (res: any) => {
        this.show(res);
        this.changeDetectorRef.markForCheck();
      },
      () => {
        this.defaults = null;
        this.changeDetectorRef.markForCheck();
      },
    );
  }

  ngOnDestroy() {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }

  get usingDefaults(): boolean {
    return !!this.defaults && this.rpcList().join('\n') === this.defaults.rpc_urls.join('\n') &&
      this.indexer.trim().replace(/\/$/, '') === this.defaults.indexer_url;
  }

  save() {
    this.send({ rpc_urls: this.rpcList(), indexer_url: this.indexer.trim() });
  }

  useDefaults() {
    this.send({});
  }

  private send(body: any) {
    this.busy = true;
    this.subscription = this.http.put(this.endpointsUrl, body).subscribe(
      (res: any) => {
        this.busy = false;
        this.show(res);
        this.msgBarService.showDone('nodes.eth.saved');
        this.changeDetectorRef.markForCheck();
      },
      err => {
        this.busy = false;
        this.msgBarService.showError((err.error && err.error.error) || err.message);
        this.changeDetectorRef.markForCheck();
      },
    );
  }

  private show(res: any) {
    this.rpcText = (res.rpc_urls || []).join('\n');
    this.indexer = res.indexer_url || '';
    this.defaults = { rpc_urls: res.default_rpc_urls || [], indexer_url: res.default_indexer_url || '' };
  }

  private rpcList(): string[] {
    return this.rpcText.split('\n').map(u => u.trim().replace(/\/$/, '')).filter(u => u.length > 0);
  }

  private get endpointsUrl(): string {
    const custom = (this.coinService.customNodeUrls as any)[this.coin.id.toString()];
    const base = (custom || this.coin.nodeUrl).replace(/\/$/, '');

    return base + '/api/v1/eth/endpoints';
  }
}
