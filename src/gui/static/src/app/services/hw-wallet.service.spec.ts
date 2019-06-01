import { TestBed, inject } from '@angular/core/testing';

import { HwWalletService } from './hw-wallet.service';

describe('HwWalletService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [HwWalletService],
    });
  });

  it('should be created', inject([HwWalletService], (service: HwWalletService) => {
    expect(service).toBeTruthy();
  }));
});
