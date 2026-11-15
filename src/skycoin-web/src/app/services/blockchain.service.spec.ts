import { TestBed, inject, fakeAsync, tick, discardPeriodicTasks } from '@angular/core/testing';
import { of } from 'rxjs';

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

  beforeEach(inject([BlockchainService, ApiService, BalanceService], (serv: BlockchainService, mockApiService: any, mockBalanceService: any) => {
    service = serv;
    apiService = mockApiService;
    balanceService = mockBalanceService;
  }));

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});

describe('BlockchainService node check', () => {
  it('checks the node again after it had no connections', fakeAsync(() => {
    let connections: any[] = [];
    const api = {
      get: (url: string) => {
        switch (url) {
          case 'network/connections': return of({ connections });
          case 'health': return of({ version: { version: '0.28.7' }, user_verify_transaction: { max_decimals: 3, burn_factor: 10 } });
          case 'blockchain/progress': return of({ current: 5, highest: 5 });
        }
        return of({});
      },
    };
    const globals = new MockGlobalsService();
    globals.setNodeVersion(null as any);
    const service = new BlockchainService(api as any, new MockBalanceService() as any, globals as any, new MockCoinService() as any);
    let synchronized = false;
    service.synchronized.subscribe(s => synchronized = s);

    tick(1000);
    expect(synchronized).toBe(false);

    connections = [{}];
    tick(10000);
    expect(synchronized).toBe(true);
    discardPeriodicTasks();
  }));
});
