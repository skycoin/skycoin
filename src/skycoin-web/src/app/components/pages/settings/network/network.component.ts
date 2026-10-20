import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';

import { Connection, NetworkService } from '../../../../services/network.service';

/** The peers of the current coin's node. */
@Component({
    selector: 'app-network',
    templateUrl: './network.component.html',
    styleUrls: ['./network.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class NetworkComponent implements OnInit, OnDestroy {
  peers: Connection[] | null = null;
  showError = false;

  private subscription!: Subscription;

  constructor(
    private networkService: NetworkService,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit() {
    this.subscription = this.networkService.connections().subscribe(peers => {
      if (peers === null) {
        this.showError = this.peers === null;
      } else {
        this.peers = peers;
        this.showError = false;
      }
      this.changeDetectorRef.markForCheck();
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  /** How long ago a unix time in seconds was, shortly: "12s", "4m", "3h" or "2d". */
  ago(unixSeconds: number): string {
    if (!unixSeconds) {
      return '-';
    }
    const s = Math.max(0, Math.floor(Date.now() / 1000 - unixSeconds));
    if (s < 60) {
      return s + 's';
    }
    if (s < 3600) {
      return Math.floor(s / 60) + 'm';
    }
    if (s < 86400) {
      return Math.floor(s / 3600) + 'h';
    }

    return Math.floor(s / 86400) + 'd';
  }
}
