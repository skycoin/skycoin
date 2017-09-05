import { TestBed, inject } from '@angular/core/testing';

import { PurchaseService } from './purchase.service';

describe('PurchaseService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [PurchaseService]
    });
  });

  it('should be created', inject([PurchaseService], (service: PurchaseService) => {
    expect(service).toBeTruthy();
  }));
});
