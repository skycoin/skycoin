import { TestBed, inject } from '@angular/core/testing';

import { BlockchainService } from './blockchain.service';
import { ApiService } from './api.service';
import { BalanceService } from './wallet/balance.service';
import { MockApiService, MockBalanceService, MockCoinService, MockGlobalsService } from '../utils/test-mocks';
import { CoinService } from './coin.service';
import { GlobalsService } from './globals.service';

describe('BlockchainService', () => {
  let service: BlockchainService;
  let apiService: ApiService;
  let balanceService: BalanceService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        BlockchainService,
        { provide: ApiService, useClass: MockApiService },
        { provide: BalanceService, useClass: MockBalanceService },
        { provide: CoinService, useClass: MockCoinService },
        { provide: GlobalsService, useClass: MockGlobalsService }
      ]
    });
  });

  beforeEach(inject([BlockchainService, ApiService, BalanceService], (serv, mockApiService, mockBalanceService) => {
    service = serv;
    apiService = mockApiService;
    balanceService = mockBalanceService;
  }));

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
