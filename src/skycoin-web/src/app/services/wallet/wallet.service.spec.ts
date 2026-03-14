import { TestBed, fakeAsync } from '@angular/core/testing';
import { Observable, of, BehaviorSubject } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { BigNumber } from 'bignumber.js';

import { WalletService } from './wallet.service';
import { CipherProvider } from '../cipher.provider';
import { CoinService } from '../coin.service';
import { EventEmitter } from '@angular/core';
import { MockCoinService, MockGlobalsService } from '../../utils/test-mocks';
import { ApiService } from '../api.service';
import { GlobalsService } from '../globals.service';
import {
  Wallet,
  Address } from '../../app.datatypes';

describe('WalletService', () => {
  let store = {};
  let walletService: WalletService;
  let spyTranslateService: jasmine.SpyObj<TranslateService>;
  let spyCipherProvider: jasmine.SpyObj<CipherProvider>;

  beforeEach(() => {
    spyOn(localStorage, 'setItem').and.callFake((key, value) => store[key] = value);
    spyOn(localStorage, 'getItem').and.callFake((key) => store[key]);

    TestBed.configureTestingModule({
      providers: [
        WalletService,
        {
          provide: CipherProvider,
          useValue: jasmine.createSpyObj('CipherProvider', ['generateAddress', 'prepareTransaction'])
        },
        {
          provide: TranslateService,
          useValue: jasmine.createSpyObj('TranslateService', ['instant'])
        },
        {
          provide: ApiService,
          useValue: jasmine.createSpyObj('ApiService', {
            'get': of([])
          })
        },
        { provide: CoinService, useClass: MockCoinService },
        { provide: GlobalsService, useClass: MockGlobalsService }
      ]
    });

    walletService = TestBed.inject(WalletService);
    spyCipherProvider = TestBed.inject(CipherProvider);
    spyTranslateService = TestBed.inject(TranslateService);
  });

  afterEach(() => {
    store = {};
  });

  it('should be created', () => {
    expect(walletService).toBeTruthy();
  });

  describe('addAddress', () => {
    it('should fail adding an address if the wallet seed is null', () => {
      const wallet = createWallet();
      wallet.seed = null;
      wallet.nextSeed = null;
      spyTranslateService.instant.and.returnValue('trying to generate address without seed!');

      expect(() => walletService.addAddress(wallet))
        .toThrowError('trying to generate address without seed!');
    });

    it('saveWallets should be called', () => {
      const wallet = createWallet();
      const expectedWallet = createWallet();
      const newAddress = createAddress('new address');
      expectedWallet.addresses.push(newAddress);

      spyCipherProvider.generateAddress.and.returnValue(of({ address: newAddress, nextSeed: 'next seed' }));

      walletService.addAddress(wallet)
        .subscribe(() => {
          expect(wallet).toEqual(expectedWallet);
        });
    });
  });

  describe('create', () => {
    it('should create a wallet', () => {
      const walletLabel = 'wallet label';
      const walletSeed = 'wallet seed';
      const walletCoinId = 1;
      const walletAddress = createAddress();
      const expectedWallet = {
        label: walletLabel,
        seed: walletSeed,
        needSeedConfirmation: true,
        balance: new BigNumber(0),
        hours: new BigNumber(0),
        addresses: [walletAddress],
        coinId: walletCoinId,
        nextSeed: 'next seed'
      };

      spyCipherProvider.generateAddress.and.returnValue(of({ address: walletAddress, nextSeed: 'next seed' }));

      walletService.create(walletLabel, walletSeed, walletCoinId)
        .subscribe(() => expect(walletService.wallets.value[0]).toEqual(expectedWallet));
    });
  });

  describe('delete', () => {
    it('should delete a wallet', () => {
      const walletToDelete = createWallet();
      walletService.wallets = new BehaviorSubject([walletToDelete]);

      walletService.delete(walletToDelete);
      expect(walletService.wallets.value.length).toEqual(0);
    });
  });

  describe('unlockWallet', () => {
    it('the wallet should be updated', fakeAsync(() => {
      const inputWallet: Wallet = createWallet('wallet', 'no seed');
      const correctSeed = 'seed';

      spyCipherProvider.generateAddress.and.returnValue(of({ address: createAddress(), nextSeed: 'next seed' }));
      walletService.unlockWallet(inputWallet, correctSeed, new EventEmitter<number>())
        .subscribe(() => expect(inputWallet.seed).toEqual(correctSeed));
    }));

    it('should only accept the correct seed', fakeAsync(() => {
      const wallet: Wallet = createWallet();
      const wrongSeedAddress: Address = createAddress('wrong address');

      spyCipherProvider.generateAddress.and.returnValue(of({ address: wrongSeedAddress, nextSeed: 'next seed' }));
      spyTranslateService.instant.and.returnValue('Wrong seed');

      walletService.unlockWallet(wallet, 'wrong seed', new EventEmitter<number>())
        .subscribe(
          () => fail('should have thrown an error'),
          (error) => expect(error.message).toBe('Wrong seed')
        );
    }));
  });
});

export function createWallet(label: string = 'label', seed: string = 'seed', balance: BigNumber = new BigNumber(0), nextSeed = 'next seed'): Wallet {
  return {
    label: label,
    seed: seed,
    nextSeed: nextSeed,
    needSeedConfirmation: true,
    balance: balance,
    hours: new BigNumber(0),
    addresses: [
      createAddress()
    ],
    coinId: 1
  };
}

export function createAddress(address: string = 'address', secretKey: string = 'secret key', publicKey: string = 'public key'): Address {
  return {
    address: address,
    secret_key: secretKey,
    public_key: publicKey,
    balance: new BigNumber(0),
    hours: new BigNumber(0),
    outputs: []
  };
}
