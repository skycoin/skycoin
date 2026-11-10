import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError, filter, map, shareReplay, switchMap } from 'rxjs';

import { BaseCoin } from '../coins/basecoin';
import { ApiService } from './api.service';
import { CoinService } from './coin.service';
import { shouldUpgradeVersion } from '../utils/semver';

/** What the header warns about the node the wallet uses. */
export interface NodeNotices {
  /** The node holds wallets and answers state changing requests without a CSRF token. */
  csrfDisabled: boolean;
  /** The newest skycoin release, when it is newer than the node. */
  latestVersion: string | null;
  /** Where that release can be downloaded. */
  downloadUrl: string | null;
}

const none: NodeNotices = { csrfDisabled: false, latestVersion: null, downloadUrl: null };

/**
 * Checks the node of the current coin once per coin for the warnings the old
 * desktop wallet showed: CSRF turned off, and a newer skycoin release. Every
 * failure, as when offline, just means no warning.
 */
@Injectable()
export class NodeNoticesService {
  static readonly latestReleaseUrl = 'https://api.github.com/repos/skycoin/skycoin/releases/latest';

  readonly notices: Observable<NodeNotices>;

  constructor(private http: HttpClient, private apiService: ApiService, coinService: CoinService) {
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
      switchMap((health: any) => {
        // A public node runs without CSRF and without the wallet API, and then CSRF guards nothing.
        const csrfDisabled = health.csrf_enabled === false && health.wallet_api_enabled === true;
        const nodeVersion = (health.version && health.version.version) || '';
        // Only skycoin's own node is a skycoin release, and an empty
        // version_url is how a node turns version checks off.
        if (health.coin !== 'skycoin' || !health.fiber || !health.fiber.version_url || !nodeVersion) {
          return of({ ...none, csrfDisabled });
        }

        return this.latestRelease().pipe(map(release => {
          if (!release || !shouldUpgradeVersion(nodeVersion, release.version)) {
            return { ...none, csrfDisabled };
          }

          return { csrfDisabled, latestVersion: release.version, downloadUrl: release.url };
        }));
      }),
      catchError(() => of(none)),
    );
  }

  private latestRelease(): Observable<{ version: string, url: string } | null> {
    return this.http.get(NodeNoticesService.latestReleaseUrl).pipe(
      map((release: any) => {
        const tag: string = (release && release.tag_name) || '';
        if (!tag) {
          return null;
        }

        return { version: tag.replace(/^v/, ''), url: release.html_url || 'https://github.com/skycoin/skycoin/releases/latest' };
      }),
      catchError(() => of(null)),
    );
  }
}
