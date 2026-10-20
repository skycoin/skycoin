import { TestBed } from '@angular/core/testing';
import { TranslateService } from '@ngx-translate/core';
import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';

import { CoinService } from './coin.service';
import { MockTranslateService } from '../utils/test-mocks';

describe('CoinService served by a node', () => {
  let http: HttpTestingController;

  function load(): CoinService {
    TestBed.configureTestingModule({
      providers: [
        CoinService,
        { provide: TranslateService, useClass: MockTranslateService },
        provideHttpClient(),
        provideHttpClientTesting(),
      ],
    });
    const service = TestBed.inject(CoinService);
    http = TestBed.inject(HttpTestingController);

    return service;
  }

  afterEach(() => http.verify());

  it('keeps the wallets in the node when its wallet API answers', () => {
    const service = load();
    http.expectOne('/api/v1/coins').flush('404 page not found', { status: 404, statusText: 'Not Found' });
    http.expectOne('/api/v1/wallets').flush([]);

    expect(service.coins.length).toBe(1);
    expect(service.coins[0].serverWallets).toBeTrue();
    expect(service.coins[0].nodeUrl).toBe('');
  });

  it('keeps the wallets in the browser when the wallet API is off', () => {
    const service = load();
    http.expectOne('/api/v1/coins').flush('404 page not found', { status: 404, statusText: 'Not Found' });
    http.expectOne('/api/v1/wallets').flush('403 Forbidden - Endpoint is disabled', { status: 403, statusText: 'Forbidden' });

    expect(service.coins[0].serverWallets).toBeFalse();
  });

  it('uses the coin list when there is one, without probing', () => {
    const service = load();
    http.expectOne('/api/v1/coins').flush([{ id: 7, coinName: 'Skycoin', nodeUrl: '/coin/7', serverWallets: true }]);
    http.expectNone('/api/v1/wallets');

    expect(service.coins.map(c => c.id)).toEqual([7]);
  });
});
