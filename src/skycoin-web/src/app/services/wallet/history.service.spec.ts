import { TestBed, fakeAsync } from '@angular/core/testing';
import 'rxjs/add/observable/of';
import { Observable } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { HistoryService } from './history.service';
import { WalletService } from './wallet.service';
import { ApiService } from '../api.service';
import { MockWalletService, MockGlobalsService } from '../../utils/test-mocks';
import { createAddress } from './wallet.service.spec';
import { GlobalsService } from '../globals.service';
import {
  Address,
  NormalTransaction } from '../../app.datatypes';

describe('HistoryService', () => {
  let store = {};
  let historyService: HistoryService;
  let walletService:  WalletService;
  let spyApiService:  jasmine.SpyObj<ApiService>;

  beforeEach(() => {
    spyOn(localStorage, 'setItem').and.callFake((key, value) => store[key] = value);
    spyOn(localStorage, 'getItem').and.callFake((key) => store[key]);

    TestBed.configureTestingModule({
      providers: [
        HistoryService,
        { provide: WalletService, useClass: MockWalletService },
        { provide: GlobalsService, useClass: MockGlobalsService },
        {
          provide: ApiService,
          useValue: jasmine.createSpyObj('ApiService', {
            'get': Observable.of([])
          })
        }
      ]
    });

    historyService = TestBed.inject(HistoryService);
    walletService = TestBed.inject(WalletService);
    spyApiService = TestBed.inject(ApiService);
  });

  afterEach(() => {
    store = {};
  });

  it('should be created', () => {
    expect(historyService).toBeTruthy();
  });

  describe('retrieveAddressTransactions', () => {
    it('should return a transactions array', fakeAsync(() => {
      const apiResponse = createAddressTransactions('owner address', 'destination address');
      const expectedTransaction = createTransaction([], 'owner address', 'destination address');

      spyApiService.get.and.returnValue( Observable.of([apiResponse]) );

      historyService.retrieveAddressTransactions( createAddress() )
        .subscribe((transactions: NormalTransaction[]) => {
          expect(transactions).toEqual([expectedTransaction]);
        });
    }));
  });

  describe('transactions', () => {
    it('should return an outgoing transaction', fakeAsync(() => {
      const ownerAddress: Address = createAddress('owner address');
      spyOnProperty(walletService, 'addresses', 'get').and.returnValue( Observable.of([ownerAddress]) );

      const destinationAddress = 'destination address';
      const apiResponse = createAddressTransactions(ownerAddress.address, destinationAddress, 13);
      spyApiService.get.and.returnValue( Observable.of([apiResponse]) );

      const expectedTransaction: NormalTransaction = createTransaction([ownerAddress.address], ownerAddress.address, destinationAddress, 13, new BigNumber(-13));
      expectedTransaction['hoursSent'] = new BigNumber(NaN);
      expectedTransaction['hoursBurned'] = new BigNumber(NaN);

      historyService.transactions()
        .subscribe((transactions: any[]) => {
          expect(transactions).toEqual([expectedTransaction]);
        });
    }));

    it('should return an incoming transaction', fakeAsync(() => {
      const destinationAddress: Address = createAddress('destination address');
      spyOnProperty(walletService, 'addresses', 'get').and.returnValue( Observable.of([destinationAddress]) );

      const ownerAddress = 'owner address';
      const apiResponse = createAddressTransactions(ownerAddress, destinationAddress.address, 13);
      spyApiService.get.and.returnValue( Observable.of([apiResponse]) );

      const expectedTransaction: NormalTransaction = createTransaction([destinationAddress.address], ownerAddress, destinationAddress.address, 13, new BigNumber(13));
      expectedTransaction['hoursSent'] = new BigNumber(NaN);
      expectedTransaction['hoursBurned'] = new BigNumber(NaN);

      historyService.transactions()
        .subscribe((transactions: any[]) => {
          expect(transactions).toEqual([expectedTransaction]);
        });
    }));
  });

  describe('getAllPendingTransactions', () => {
    it('should get the transactions', fakeAsync(() => {
      historyService.getAllPendingTransactions().subscribe();
      expect(spyApiService.get).toHaveBeenCalledWith('pendingTxs');
    }));
  });
});

function createAddressTransactions(ownerAddress: string, destinationAddress: string, coins: number = 0) {
  return {
    status: {
      block_seq: 1,
      confirmed: true
    },
    timestamp: 0,
    txid: 'txid',
    inputs: [{
      owner: ownerAddress
    }],
    outputs: [{
      dst: destinationAddress,
      coins: coins
    }]
  };
}

function createTransaction(addresses: string[], ownerAddress: string, destinationAddress: string, coins: number = 0, balance: BigNumber = new BigNumber(0)): NormalTransaction {
  return {
    addresses: addresses,
    balance: balance,
    block: 1,
    confirmed: true,
    inputs: [{
      owner: ownerAddress
    }],
    outputs: [{
      dst: destinationAddress,
      coins: coins
    }],
    timestamp: 0,
    txid: 'txid'
  };
}
