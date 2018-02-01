import { TestBed, inject } from '@angular/core/testing';

import { PriceService } from './price.service';

describe('PriceService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [PriceService]
    });
  });

  it('should be created', inject([PriceService], (service: PriceService) => {
    expect(service).toBeTruthy();
  }));
});
