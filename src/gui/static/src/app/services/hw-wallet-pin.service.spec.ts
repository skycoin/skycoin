import { TestBed, inject } from '@angular/core/testing';

import { HwWalletPinService } from './hw-wallet-pin.service';

describe('HwWalletPinService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [HwWalletPinService],
    });
  });

  it('should be created', inject([HwWalletPinService], (service: HwWalletPinService) => {
    expect(service).toBeTruthy();
  }));
});
