import { TestBed } from '@angular/core/testing';
import { of, throwError } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';

import { ApiService } from './api.service';
import { ApiTransportService } from './api-transport.service';
import { MockTranslateService, MockCoinService } from '../utils/test-mocks';
import { CoinService } from './coin.service';

describe('ApiService', () => {
  let service: ApiService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        ApiService,
        ApiTransportService,
        { provide: TranslateService, useClass: MockTranslateService },
        { provide: CoinService, useClass: MockCoinService },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    });

    service = TestBed.inject(ApiService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});

describe('ApiService CSRF', () => {
  let service: ApiService;
  let transport: jasmine.SpyObj<ApiTransportService>;

  beforeEach(() => {
    transport = jasmine.createSpyObj('ApiTransportService', ['request']);
    TestBed.configureTestingModule({
      providers: [
        ApiService,
        { provide: ApiTransportService, useValue: transport },
        { provide: TranslateService, useClass: MockTranslateService },
        { provide: CoinService, useClass: MockCoinService },
      ],
    });
    service = TestBed.inject(ApiService);
  });

  it('sends the node token with a POST', () => {
    transport.request.and.callFake((method: string, url: string) =>
      of(url.endsWith('v1/csrf') ? { csrf_token: 'tok' } : { ok: true }));
    service.post('wallet/encrypt', { id: 'w' }).subscribe();

    const calls = transport.request.calls.all();
    expect(calls.map(c => c.args[0])).toEqual(['GET', 'POST']);
    expect(calls[1].args[3]!.headers!['X-CSRF-Token']).toBe('tok');
  });

  it('sends a POST without a token when the server has none', () => {
    transport.request.and.callFake((method: string, url: string) =>
      url.endsWith('v1/csrf') ? throwError(() => ({ status: 404 })) : of({ ok: true }));
    let result: any;
    service.post('transaction', {}, {}, true).subscribe(r => result = r);

    const post = transport.request.calls.mostRecent();
    expect(post.args[0]).toBe('POST');
    expect(post.args[3]!.headers!['X-CSRF-Token']).toBeUndefined();
    expect(result).toEqual({ ok: true });
  });

  it('does not fetch a token for a GET', () => {
    transport.request.and.returnValue(of({}));
    service.get('health').subscribe();
    expect(transport.request.calls.count()).toBe(1);
  });
});
