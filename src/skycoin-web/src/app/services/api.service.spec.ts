import { TestBed } from '@angular/core/testing';
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
