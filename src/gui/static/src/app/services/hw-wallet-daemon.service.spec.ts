import { TestBed, inject } from '@angular/core/testing';

import { HwWalletDaemonService } from './hw-wallet-daemon.service';

describe('HwWalletDaemonService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [HwWalletDaemonService],
    });
  });

  it('should be created', inject([HwWalletDaemonService], (service: HwWalletDaemonService) => {
    expect(service).toBeTruthy();
  }));
});
