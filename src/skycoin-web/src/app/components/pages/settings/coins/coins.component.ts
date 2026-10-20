import { Component, OnInit, OnDestroy, ChangeDetectionStrategy, ChangeDetectorRef } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Subscription } from 'rxjs';

import { MsgBarService } from '../../../../services/msg-bar.service';

/** The largest icon the server keeps. */
const maxIconBytes = 96 * 1024;

/**
 * Adds fiber coins and ERC-20 tokens to the ones the wallet service offers,
 * and removes added ones. The service keeps them beside its wallets.
 */
@Component({
    selector: 'app-coins',
    templateUrl: './coins.component.html',
    styleUrls: ['./coins.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class CoinsComponent implements OnInit, OnDestroy {
  added: any[] | null = null;
  unavailable = false;
  busy = false;

  kind: 'fiber' | 'erc20' = 'fiber';
  name = '';
  ticker = '';
  nodeUrl = '';
  contract = '';
  decimals = 18;
  icon = '';

  private subscription: Subscription | null = null;

  constructor(
    private http: HttpClient,
    private msgBarService: MsgBarService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {}

  ngOnInit() {
    this.subscription = this.http.get('/api/v1/coins/custom').subscribe(
      (list: any) => {
        this.added = list || [];
        this.changeDetectorRef.markForCheck();
      },
      () => {
        this.unavailable = true;
        this.changeDetectorRef.markForCheck();
      },
    );
  }

  ngOnDestroy() {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }

  pickIcon(event: Event) {
    const file = (event.target as HTMLInputElement).files?.[0];
    if (!file) {
      return;
    }
    // The data URL is a third larger than the file.
    if (file.size * 4 / 3 > maxIconBytes) {
      this.msgBarService.showError('coins.icon-too-big');
      return;
    }
    const reader = new FileReader();
    reader.onload = () => {
      this.icon = reader.result as string;
      this.changeDetectorRef.markForCheck();
    };
    reader.readAsDataURL(file);
  }

  add() {
    const body: any = { kind: this.kind, name: this.name.trim(), ticker: this.ticker.trim(), icon: this.icon || undefined };
    if (this.kind === 'fiber') {
      body.node_url = this.nodeUrl.trim();
    } else {
      body.contract = this.contract.trim();
      body.decimals = Number(this.decimals);
    }
    this.busy = true;
    this.subscription = this.http.post('/api/v1/coins/custom', body).subscribe(
      () => {
        this.done('coins.added');
        this.changeDetectorRef.markForCheck();
      },
      err => {
        this.failed(err);
        this.changeDetectorRef.markForCheck();
      },
    );
  }

  remove(coin: any) {
    this.busy = true;
    this.subscription = this.http.delete('/api/v1/coins/custom', { params: { id: coin.id } }).subscribe(
      () => {
        this.done('coins.removed');
        this.changeDetectorRef.markForCheck();
      },
      err => {
        this.failed(err);
        this.changeDetectorRef.markForCheck();
      },
    );
  }

  // The coin list is read when the wallet starts, so it starts again.
  private done(message: string) {
    this.msgBarService.showDone(message);
    setTimeout(() => {
      window.location.reload();
      this.changeDetectorRef.markForCheck();
    }, 1500);
  }

  private failed(err: any) {
    this.busy = false;
    this.msgBarService.showError((err.error && err.error.error) || err.message);
  }
}
