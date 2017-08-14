import { TestBed, inject } from '@angular/core/testing';

import { WalletService } from './wallet.service';

describe('WalletService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [WalletService]
    });
  });

  it('should be created', inject([WalletService], (service: WalletService) => {
    expect(service).toBeTruthy();
  }));
});
