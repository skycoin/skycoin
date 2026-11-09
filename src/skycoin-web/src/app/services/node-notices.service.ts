import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { catchError, filter, map, shareReplay, switchMap } from 'rxjs';

import { BaseCoin } from '../coins/basecoin';
import { ApiService } from './api.service';
import { CoinService } from './coin.service';

/** What the header warns about the node the wallet uses. */
export interface NodeNotices {
  /** The node answers state changing requests without a CSRF token. */
  csrfDisabled: boolean;
}

const none: NodeNotices = { csrfDisabled: false };

/**
 * Checks the node of the current coin once per coin for the warnings the old
 * desktop wallet showed: CSRF turned off. Every failure, as when offline, just
 * means no warning.
 */
@Injectable()
export class NodeNoticesService {
  readonly notices: Observable<NodeNotices>;

  constructor(private apiService: ApiService, coinService: CoinService) {
    this.notices = coinService.currentCoin.pipe(
      filter(coin => !!coin),
      switchMap(coin => this.check(coin)),
      shareReplay({ bufferSize: 1, refCount: false }),
    );
  }

  private check(coin: BaseCoin): Observable<NodeNotices> {
    if (coin.isGateway()) {
      return of(none);
    }

    return this.apiService.get('health').pipe(
      map((health: any) => ({ csrfDisabled: health.csrf_enabled === false })),
      catchError(() => of(none)),
    );
  }
}
