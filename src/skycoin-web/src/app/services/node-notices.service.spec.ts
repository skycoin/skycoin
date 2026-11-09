import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { Observable, of, throwError } from 'rxjs';

import { NodeNotices, NodeNoticesService } from './node-notices.service';
import { ApiService } from './api.service';
import { CoinService } from './coin.service';
import { MockCoinService } from '../utils/test-mocks';

describe('NodeNoticesService', () => {
  let http: HttpTestingController;
  let health: Observable<any>;

  function healthOf(version: string, csrf = true, coin = 'skycoin', versionUrl = 'https://version.skycoin.com/skycoin/version.txt') {
    return of({ coin, csrf_enabled: csrf, version: { version }, fiber: { version_url: versionUrl } });
  }

  function notices(): NodeNotices | undefined {
    TestBed.configureTestingModule({
      providers: [
        NodeNoticesService,
        provideHttpClient(),
        provideHttpClientTesting(),
        { provide: CoinService, useClass: MockCoinService },
        { provide: ApiService, useValue: { get: () => health } },
      ],
    });
    http = TestBed.inject(HttpTestingController);
    let result: NodeNotices | undefined;
    TestBed.inject(NodeNoticesService).notices.subscribe(n => result = n);

    return result;
  }

  function latest(): NodeNotices | undefined {
    let result: NodeNotices | undefined;
    TestBed.inject(NodeNoticesService).notices.subscribe(n => result = n);

    return result;
  }

  afterEach(() => http.verify());

  it('tells of a newer skycoin release, with its download page', () => {
    health = healthOf('0.28.6');
    notices();
    http.expectOne(NodeNoticesService.latestReleaseUrl)
      .flush({ tag_name: 'v0.28.7', html_url: 'https://github.com/skycoin/skycoin/releases/tag/v0.28.7' });
    expect(latest()).toEqual({
      csrfDisabled: false,
      latestVersion: '0.28.7',
      downloadUrl: 'https://github.com/skycoin/skycoin/releases/tag/v0.28.7',
    });
  });

  it('says nothing when the node is up to date', () => {
    health = healthOf('0.28.7');
    notices();
    http.expectOne(NodeNoticesService.latestReleaseUrl).flush({ tag_name: 'v0.28.7' });
    expect(latest()!.latestVersion).toBeNull();
  });

  it('fails quietly when the release cannot be fetched', () => {
    health = healthOf('0.28.6', false);
    notices();
    http.expectOne(NodeNoticesService.latestReleaseUrl).error(new ProgressEvent('offline'), { status: 0 });
    expect(latest()).toEqual({ csrfDisabled: true, latestVersion: null, downloadUrl: null });
  });

  it('fails quietly when the node cannot be reached', () => {
    health = throwError(() => new Error('unreachable'));
    expect(notices()).toEqual({ csrfDisabled: false, latestVersion: null, downloadUrl: null });
  });

  it('checks no release for another fiber coin or a node with version checks off', () => {
    health = healthOf('0.1.0', true, 'mycoin');
    expect(notices()!.latestVersion).toBeNull();
    http.expectNone(NodeNoticesService.latestReleaseUrl);

    TestBed.resetTestingModule();
    health = healthOf('0.1.0', true, 'skycoin', '');
    expect(notices()!.latestVersion).toBeNull();
    http.expectNone(NodeNoticesService.latestReleaseUrl);
  });

  it('warns when the node has CSRF turned off', () => {
    health = healthOf('0.28.7', false, 'mycoin');
    expect(notices()!.csrfDisabled).toBeTrue();
  });
});
