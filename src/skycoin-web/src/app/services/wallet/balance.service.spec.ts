import { TestBed, fakeAsync, tick } from '@angular/core/testing';
import { NgZone } from '@angular/core';
import { of, BehaviorSubject } from 'rxjs';

import { BalanceService, BalanceStates } from './balance.service';
import { ApiService } from '../api.service';
import { CoinService } from '../coin.service';
import { WalletService } from './wallet.service';
import { GlobalsService } from '../globals.service';
import { MockCoinService, MockWalletService, MockGlobalsService } from '../../utils/test-mocks';
import { Address, Balance } from '../../app.datatypes';

describe('BalanceService', () => {
  let balanceService: BalanceService;
  let spyApiService: jasmine.SpyObj<ApiService>;
  let globalsService: MockGlobalsService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        BalanceService,
        { provide: WalletService, useClass: MockWalletService },
        {
          provide: ApiService,
          useValue: jasmine.createSpyObj('ApiService', ['get', 'post'])
        },
        { provide: CoinService, useClass: MockCoinService },
        { provide: GlobalsService, useClass: MockGlobalsService },
      ]
    });

    balanceService = TestBed.inject(BalanceService);
    spyApiService = TestBed.inject(ApiService) as jasmine.SpyObj<ApiService>;
    globalsService = TestBed.inject(GlobalsService) as unknown as MockGlobalsService;
  });

  it('should be created', () => {
    expect(balanceService).toBeTruthy();
  });

  describe('chunkAddresses (via private access)', () => {
    // Access private method for unit testing chunking logic
    let chunkFn: (addresses: Address[], maxChars?: number) => string[];

    beforeEach(() => {
      chunkFn = (balanceService as any).chunkAddresses.bind(balanceService);
    });

    it('should return a single chunk for a few addresses', () => {
      const addresses = makeAddresses(['addr1', 'addr2', 'addr3']);
      const result = chunkFn(addresses);
      expect(result).toEqual(['addr1,addr2,addr3']);
    });

    it('should return empty array for no addresses', () => {
      const result = chunkFn([]);
      expect(result).toEqual([]);
    });

    it('should return single address without comma', () => {
      const result = chunkFn(makeAddresses(['singleAddr']));
      expect(result).toEqual(['singleAddr']);
    });

    it('should split when combined length exceeds maxChars', () => {
      // Use a small maxChars to force splitting
      const addresses = makeAddresses(['aaaa', 'bbbb', 'cccc', 'dddd']);
      const result = chunkFn(addresses, 10);
      // 'aaaa,bbbb' = 9 chars (fits in 10)
      // 'cccc,dddd' = 9 chars (fits in 10)
      expect(result).toEqual(['aaaa,bbbb', 'cccc,dddd']);
    });

    it('should handle address exactly at the boundary', () => {
      const addresses = makeAddresses(['12345', '67890']);
      // maxChars=11: '12345,67890' = 11 chars, fits exactly
      const result = chunkFn(addresses, 11);
      expect(result).toEqual(['12345,67890']);
    });

    it('should split when adding comma+address exceeds maxChars', () => {
      const addresses = makeAddresses(['12345', '67890']);
      // maxChars=10: '12345,67890' = 11 chars, exceeds 10
      const result = chunkFn(addresses, 10);
      expect(result).toEqual(['12345', '67890']);
    });

    it('should handle many Skycoin-length addresses near 1800 char limit', () => {
      // Skycoin addresses are ~35 chars. 1800 / 36 (addr + comma) = ~50 addresses per chunk
      const addr = '2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6'; // 35 chars
      const addresses = makeAddresses(Array(100).fill(addr));
      const result = chunkFn(addresses);

      // All chunks should be <= 1800 chars
      result.forEach(chunk => {
        expect(chunk.length).toBeLessThanOrEqual(1800);
      });

      // Total addresses across all chunks should equal input
      const totalAddrs = result.reduce((sum, chunk) => sum + chunk.split(',').length, 0);
      expect(totalAddrs).toEqual(100);

      // Should need at least 2 chunks for 100 addresses
      expect(result.length).toBeGreaterThan(1);
    });

    it('should handle a single address longer than maxChars', () => {
      const longAddr = 'a'.repeat(2000);
      const result = chunkFn(makeAddresses([longAddr]), 1800);
      // Should still produce the chunk even though it exceeds maxChars
      // (can't split a single address)
      expect(result).toEqual([longAddr]);
    });
  });

  describe('chunkedBalanceGet (via private access)', () => {
    let chunkedBalanceGetFn: (addresses: Address[]) => any;

    beforeEach(() => {
      chunkedBalanceGetFn = (balanceService as any).chunkedBalanceGet.bind(balanceService);
    });

    it('should make a single GET request when all addresses fit in one chunk', fakeAsync(() => {
      const mockBalance: Balance = {
        confirmed: { coins: 1000000, hours: 100 },
        predicted: { coins: 1000000, hours: 100 },
        addresses: { 'addr1': { confirmed: { coins: 1000000, hours: 100 }, predicted: { coins: 1000000, hours: 100 } } },
      };
      spyApiService.get.and.returnValue(of(mockBalance));

      const addresses = makeAddresses(['addr1', 'addr2']);
      chunkedBalanceGetFn(addresses).subscribe((result: Balance) => {
        expect(result).toEqual(mockBalance);
        expect(spyApiService.get).toHaveBeenCalledTimes(1);
        expect(spyApiService.get).toHaveBeenCalledWith('balance', { addrs: 'addr1,addr2' });
      });
    }));

    it('should merge results from multiple chunks', fakeAsync(() => {
      const balance1: Balance = {
        confirmed: { coins: 1000000, hours: 50 },
        predicted: { coins: 1000000, hours: 50 },
        addresses: { 'addr1': { confirmed: { coins: 1000000, hours: 50 }, predicted: { coins: 1000000, hours: 50 } } },
      };
      const balance2: Balance = {
        confirmed: { coins: 2000000, hours: 100 },
        predicted: { coins: 2000000, hours: 100 },
        addresses: { 'addr2': { confirmed: { coins: 2000000, hours: 100 }, predicted: { coins: 2000000, hours: 100 } } },
      };

      spyApiService.get.and.returnValues(of(balance1), of(balance2));

      // Use addresses that will force splitting (small maxChars via override)
      // We need to force multiple chunks - use many addresses
      const addr = '2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6';
      const addresses = makeAddresses(Array(60).fill(addr));

      chunkedBalanceGetFn(addresses).subscribe((result: Balance) => {
        // Coins and hours should be summed
        expect(result.confirmed.coins).toEqual(3000000);
        expect(result.confirmed.hours).toEqual(150);
        expect(result.predicted.coins).toEqual(3000000);
        expect(result.predicted.hours).toEqual(150);
        // Addresses should be merged
        expect(result.addresses['addr1']).toBeTruthy();
        expect(result.addresses['addr2']).toBeTruthy();
        // Should have made 2 GET requests
        expect(spyApiService.get).toHaveBeenCalledTimes(2);
      });
    }));
  });
});

function makeAddresses(addrs: string[]): Address[] {
  return addrs.map(addr => ({
    address: addr,
    secret_key: '',
    public_key: '',
    outputs: [],
  }));
}
