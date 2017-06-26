import { TestBed, inject } from '@angular/core/testing';

import { CoinSupplyService } from './coin-supply.service';

describe('CoinSupplyService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [CoinSupplyService]
    });
  });

  it('should ...', inject([CoinSupplyService], (service: CoinSupplyService) => {
    expect(service).toBeTruthy();
  }));
});
