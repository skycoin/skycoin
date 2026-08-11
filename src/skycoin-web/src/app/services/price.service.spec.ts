import { TestBed } from '@angular/core/testing';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';

import { PriceService } from './price.service';
import { CoinService } from './coin.service';
import { MockCoinService } from '../utils/test-mocks';

describe('PriceService', () => {
  let priceService: PriceService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        PriceService,
        { provide: CoinService, useClass: MockCoinService },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    });

    priceService = TestBed.inject(PriceService);
  });

  it('should be created', () => {
    expect(priceService).toBeTruthy();
  });
});
