import { TestBed, fakeAsync, tick } from '@angular/core/testing';
import { Observable, of, BehaviorSubject } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { BigNumber } from 'bignumber.js';

import { WalletService } from './wallet.service';
import { CipherProvider } from '../cipher.provider';
import { CoinService } from '../coin.service';
import { EncryptionService } from '../encryption.service';
import { EventEmitter } from '@angular/core';
import { MockCoinService, MockGlobalsService } from '../../utils/test-mocks';
import { ApiService } from '../api.service';
import { GlobalsService } from '../globals.service';
import {
  Wallet,
  Address } from '../../app.datatypes';

describe('WalletService', () => {
  let store: Record<string, string> = {};
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
            'get': of([]),
            'post': of({}),
          })
        },
        { provide: CoinService, useClass: MockCoinService },
        { provide: GlobalsService, useClass: MockGlobalsService },
        EncryptionService
      ]
    });

    walletService = TestBed.inject(WalletService);
    spyCipherProvider = TestBed.inject(CipherProvider) as jasmine.SpyObj<CipherProvider>;
    spyTranslateService = TestBed.inject(TranslateService) as jasmine.SpyObj<TranslateService>;
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
      wallet.seed = undefined;
      wallet.nextSeed = undefined;
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
        nextSeed: 'next seed',
        walletType: 'deterministic'
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

  describe('ensureKeys', () => {
    it('derives the keys a reload dropped from the stored seed', fakeAsync(() => {
      const wallet: Wallet = createWallet();
      wallet.addresses = [{ address: 'address' }];

      spyCipherProvider.generateAddress.and.returnValue(of({ address: createAddress(), nextSeed: 'next seed' }));
      walletService.ensureKeys(wallet).subscribe();

      expect(spyCipherProvider.generateAddress).toHaveBeenCalled();
      expect(wallet.addresses[0].secret_key).toBe('secret key');
    }));

    it('leaves a wallet that has its keys alone', fakeAsync(() => {
      walletService.ensureKeys(createWallet()).subscribe();

      expect(spyCipherProvider.generateAddress).not.toHaveBeenCalled();
    }));
  });

  describe('findUsedAddresses', () => {
    let api: jasmine.SpyObj<ApiService>;

    beforeEach(() => {
      api = TestBed.inject(ApiService) as jasmine.SpyObj<ApiService>;
      (TestBed.inject(CoinService) as unknown as MockCoinService).currentCoin.value.serverWallets = true;
    });

    function nodeWallet(type: string, external: string[], change: string[]) {
      return {
        meta: { type: type, filename: 'w.wlt', label: 'w' },
        entries: external.map(address => ({ address })),
        accounts: [{
          name: 'Account 0',
          index: 0,
          external_entries: external.map(address => ({ address })),
          change_entries: change.map(address => ({ address })),
        }],
      };
    }

    function serverWallet(type = 'bip44'): Wallet {
      return {
        label: 'w',
        filename: 'w.wlt',
        walletType: type,
        addresses: [{ address: 'e0' }],
        accounts: [{ name: 'Account 0', index: 0, externalAddresses: [{ address: 'e0' }], changeAddresses: [] }],
      };
    }

    it('asks the node to scan with the gap limit and reads both chains back', fakeAsync(() => {
      const wallet = serverWallet();
      walletService.wallets.next([wallet]);
      api.post.and.returnValue(of({ addresses: ['e1'] }));
      api.get.and.returnValue(of(nodeWallet('bip44', ['e0', 'e1'], ['c0'])));

      let added = -1;
      walletService.findUsedAddresses(wallet).subscribe(n => added = n);
      tick();

      expect(api.post).toHaveBeenCalledWith('wallet/scan', { id: 'w.wlt', num: '20' });
      expect(added).toBe(2);
      expect(wallet.accounts![0].changeAddresses.map(a => a.address)).toEqual(['c0']);
    }));

    it('scans each wallet at most once a minute on its own', fakeAsync(() => {
      walletService.wallets.next([serverWallet()]);
      api.post.and.returnValue(of({ addresses: [] }));
      api.get.and.returnValue(of(nodeWallet('bip44', ['e0'], [])));

      walletService.findAllUsedAddresses();
      walletService.findAllUsedAddresses();
      expect(api.post).toHaveBeenCalledTimes(1);

      tick(60 * 1000);
      walletService.findAllUsedAddresses();
      expect(api.post).toHaveBeenCalledTimes(2);
    }));

    it('leaves a deterministic node wallet alone', fakeAsync(() => {
      const wallet = serverWallet('deterministic');
      walletService.wallets.next([wallet]);

      walletService.findAllUsedAddresses();

      expect(walletService.canFindUsedAddresses(wallet)).toBeFalse();
      expect(api.post).not.toHaveBeenCalled();
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
