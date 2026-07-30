import { TestBed } from '@angular/core/testing';
import { MockBackend } from '@angular/http/testing';
import { XHRBackend } from '@angular/http';
import { provideHttpClient, withInterceptorsFromDi, withXhr } from '@angular/common/http';

import { PriceService } from './price.service';
import { CoinService } from './coin.service';
import { MockCoinService } from '../utils/test-mocks';

describe('PriceService', () => {
  let priceService: PriceService;

  beforeEach(() => {
    TestBed.configureTestingModule({
    imports: [],
    providers: [
        PriceService,
        { provide: XHRBackend, useClass: MockBackend },
        { provide: CoinService, useClass: MockCoinService },
        provideHttpClient(withXhr(), withInterceptorsFromDi())
    ]
});

    priceService = TestBed.inject(PriceService);
  });

  it('should be created', () => {
    expect(priceService).toBeTruthy();
  });
});
