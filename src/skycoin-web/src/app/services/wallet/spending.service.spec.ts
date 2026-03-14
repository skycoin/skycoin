import { TestBed, fakeAsync } from '@angular/core/testing';
import { Observable, of } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import BigNumber from 'bignumber.js';

import { SpendingService, HoursSelectionTypes } from './spending.service';
import { WalletService } from './wallet.service';
import { ApiService } from '../api.service';
import { CipherProvider } from '../cipher.provider';
import { CoinService } from '../coin.service';
import { MockCoinService, MockWalletService, MockGlobalsService, MockBlockchainService } from '../../utils/test-mocks';
import { createWallet, createAddress } from './wallet.service.spec';
import { GlobalsService } from '../globals.service';
import {
  Wallet,
  TransactionOutput,
  TransactionInput,
  Output,
  GetOutputsRequestOutput } from '../../app.datatypes';
import { BlockchainService } from '../blockchain.service';

describe('SpendingService', () => {
  let store = {};
  let spendingService: SpendingService;
  let walletService:  WalletService;
  let spyApiService:  jasmine.SpyObj<ApiService>;
  let spyTranslateService: jasmine.SpyObj<TranslateService>;
  let spyCipherProvider: jasmine.SpyObj<CipherProvider>;

  beforeEach(() => {
    spyOn(localStorage, 'setItem').and.callFake((key, value) => store[key] = value);
    spyOn(localStorage, 'getItem').and.callFake((key) => store[key]);

    TestBed.configureTestingModule({
      providers: [
        SpendingService,
        { provide: WalletService, useClass: MockWalletService },
        {
          provide: ApiService,
          useValue: jasmine.createSpyObj('ApiService', {
            'get': of([])
          })
        },
        {
          provide: CipherProvider,
          useValue: jasmine.createSpyObj('CipherProvider', ['generateAddress', 'prepareTransaction'])
        },
        {
          provide: TranslateService,
          useValue: jasmine.createSpyObj('TranslateService', ['instant'])
        },
        { provide: CoinService, useClass: MockCoinService },
        { provide: GlobalsService, useClass: MockGlobalsService },
        { provide: BlockchainService, useClass: MockBlockchainService }
      ]
    });

    spendingService = TestBed.inject(SpendingService);
    walletService = TestBed.inject(WalletService);
    spyApiService = TestBed.inject(ApiService);
    spyCipherProvider = TestBed.inject(CipherProvider);
    spyTranslateService = TestBed.inject(TranslateService);
  });

  afterEach(() => {
    store = {};
  });

  it('should be created', () => {
    expect(spendingService).toBeTruthy();
  });

  describe('createTransaction', () => {

    let wallet: Wallet;
    let outputs: Output[];

    beforeEach(() => {
      const addresses = [
        createAddress('address1', 'secretKey1'),
        createAddress('address2', 'secretKey2')
      ];

      wallet = Object.assign(createWallet(), { addresses: addresses });

      outputs = [
        createOutput('address1', 'hash1', new BigNumber(5), new BigNumber(10)),
        createOutput('address2', 'hash2', new BigNumber(10), new BigNumber(20)),
        createOutput('address1', 'hash3', new BigNumber(20), new BigNumber(50)),
        createOutput('address2', 'hash4', new BigNumber(5), new BigNumber(0))
      ];

      spyApiService.get.and.returnValue(of({ head_outputs: outputs }));
      spyCipherProvider.prepareTransaction.and.returnValue(of('preparedTransaction'));
      spyTranslateService.instant.and.callFake((param) => {
        if (param === 'service.wallet.not-enough-hours1') {
          return 'Not enough available';
        }

        if (param === 'service.wallet.not-enough-hours2') {
          return 'to perform transaction!';
        }
      });
    });

    it('should return a correct observable for two outputs', fakeAsync(() => {
      const address = 'address';
      const amount = new BigNumber(26);

      const expectedTxInputs: TransactionInput[] = [
        { hash: 'hash3', secret: 'secretKey1', address: 'address1', calculated_hours: 50, coins: 20 },
        { hash: 'hash4', secret: 'secretKey2', address: 'address2', calculated_hours: 0, coins: 5 },
        { hash: 'hash1', secret: 'secretKey1', address: 'address1', calculated_hours: 10, coins: 5 }
      ];

      const expectedTxOutputs: TransactionOutput[] = [
        {
          address: address,
          coins: amount.toNumber(),
          hours: 15
        },
        {
          address: wallet.addresses[0].address,
          coins: 4,
          hours: 15
        }
      ];

      spendingService.createTransaction(
        wallet,
        null,
        null,
        [{
          address: address,
          coins: amount,
        }],
        {
          type: HoursSelectionTypes.Auto,
          ShareFactor: new BigNumber(0.5),
        },
        null
      ).subscribe((result: any) => {
        expect(result).toEqual({
          inputs: expectedTxInputs,
          outputs: expectedTxOutputs,
          hoursSent: new BigNumber(15),
          hoursBurned: new BigNumber(30),
          encoded: 'preparedTransaction'
        });
      });
    }));

    it('should be return a correct observable for one output', fakeAsync(() => {
      const address = 'address';
      const amount = new BigNumber(40);

      const expectedTxInputs = [
        { hash: 'hash3', secret: 'secretKey1', address: 'address1', calculated_hours: 50, coins: 20 },
        { hash: 'hash4', secret: 'secretKey2', address: 'address2', calculated_hours: 0, coins: 5 },
        { hash: 'hash1', secret: 'secretKey1', address: 'address1', calculated_hours: 10, coins: 5 },
        { hash: 'hash2', secret: 'secretKey2', address: 'address2', calculated_hours: 20, coins: 10 }
      ];

      const expectedTxOutputs = [
        {
          address: address,
          coins: amount.toNumber(),
          hours: 40
        }
      ];

      spendingService.createTransaction(
        wallet,
        null,
        null,
        [{
          address: address,
          coins: amount,
        }],
        {
          type: HoursSelectionTypes.Auto,
          ShareFactor: new BigNumber(0.5),
        },
        null
      ).subscribe((result: any) => {
        expect(result).toEqual({
          inputs: expectedTxInputs,
          outputs: expectedTxOutputs,
          hoursSent: new BigNumber(40),
          hoursBurned: new BigNumber(40),
          encoded: 'preparedTransaction'
        });
      });
    }));

    it('should reject if there are not enough hours', fakeAsync(() => {
      const address = 'address';
      const amount = new BigNumber(41);

      spendingService.createTransaction(
        wallet,
        null,
        null,
        [{
          address: address,
          coins: amount,
        }],
        {
          type: HoursSelectionTypes.Auto,
          ShareFactor: new BigNumber(0.5),
        },
        null
      ).subscribe(() => fail('should be rejected'),
        (error) => expect(error.message).toBe('Not enough available Test Hours to perform transaction!')
      );
    }));

    it('should return a correct observable while using a different share factor', fakeAsync(() => {
      const address1 = 'address_1';
      const amount1 = new BigNumber(10);
      const address2 = 'address_2';
      const amount2 = new BigNumber(2);

      const expectedTxInputs = [
        { hash: 'hash3', secret: 'secretKey1', address: 'address1', calculated_hours: 50, coins: 20 }
      ];

      const expectedTxOutputs = [
        {
          address: address1,
          coins: amount1.toNumber(),
          hours: 17
        },
        {
          address: address2,
          coins: amount2.toNumber(),
          hours: 3
        },
        {
          address: wallet.addresses[0].address,
          coins: 8,
          hours: 5
        }
      ];

      spendingService.createTransaction(
        wallet,
        null,
        null,
        [{
          address: address1,
          coins: amount1,
        },
        {
          address: address2,
          coins: amount2,
        }],
        {
          type: HoursSelectionTypes.Auto,
          ShareFactor: new BigNumber(0.8),
        },
        null
      ).subscribe((result: any) => {
        expect(result).toEqual({
          inputs: expectedTxInputs,
          outputs: expectedTxOutputs,
          hoursSent: new BigNumber(20),
          hoursBurned: new BigNumber(25),
          encoded: 'preparedTransaction'
        });
      });
    }));

    it('should return a correct observable while manually setting how many hours to send', fakeAsync(() => {
      const address1 = 'address_1';
      const amount1 = new BigNumber(1);
      const hours1 = new BigNumber(30);
      const address2 = 'address_2';
      const amount2 = new BigNumber(1);
      const hours2 = new BigNumber(10);

      const expectedTxInputs = [
        { hash: 'hash3', secret: 'secretKey1', address: 'address1', calculated_hours: 50, coins: 20 },
        { hash: 'hash4', secret: 'secretKey2', address: 'address2', calculated_hours: 0, coins: 5 },
        { hash: 'hash1', secret: 'secretKey1', address: 'address1', calculated_hours: 10, coins: 5 },
        { hash: 'hash2', secret: 'secretKey2', address: 'address2', calculated_hours: 20, coins: 10 }
      ];

      const expectedTxOutputs = [
        {
          address: address1,
          coins: amount1.toNumber(),
          hours: 30
        },
        {
          address: address2,
          coins: amount2.toNumber(),
          hours: 10
        },
        {
          address: wallet.addresses[0].address,
          coins: 38,
          hours: 0
        }
      ];

      spendingService.createTransaction(
        wallet,
        null,
        null,
        [{
          address: address1,
          coins: amount1,
          hours: hours1,
        },
        {
          address: address2,
          coins: amount2,
          hours: hours2,
        }],
        {
          type: HoursSelectionTypes.Manual,
        },
        null
      ).subscribe((result: any) => {
        expect(result).toEqual({
          inputs: expectedTxInputs,
          outputs: expectedTxOutputs,
          hoursSent: new BigNumber(40),
          hoursBurned: new BigNumber(40),
          encoded: 'preparedTransaction'
        });
      });
    }));

    it('should reject if there are not enough hours while manually setting how many hours to send', fakeAsync(() => {
      const address1 = 'address_1';
      const amount1 = new BigNumber(1);
      const hours1 = new BigNumber(30);
      const address2 = 'address_2';
      const amount2 = new BigNumber(1);
      const hours2 = new BigNumber(11);

      spendingService.createTransaction(
        wallet,
        null,
        null,
        [{
          address: address1,
          coins: amount1,
          hours: hours1,
        },
        {
          address: address2,
          coins: amount2,
          hours: hours2,
        }],
        {
          type: HoursSelectionTypes.Manual,
        },
        null
      ).subscribe(() => fail('should be rejected'),
        (error) => expect(error.message).toBe('Not enough available Test Hours to perform transaction!')
      );
    }));

    it('should return a correct observable while trying to send from specific addresses', fakeAsync(() => {
      const address1 = 'address_1';
      const amount1 = new BigNumber(1);
      const address2 = 'address_2';
      const amount2 = new BigNumber(1);

      const expectedTxInputs = [
        { hash: 'hash2', secret: 'secretKey2', address: 'address2', calculated_hours: 20, coins: 10 }
      ];

      const expectedTxOutputs = [
        {
          address: address1,
          coins: amount1.toNumber(),
          hours: 3
        },
        {
          address: address2,
          coins: amount2.toNumber(),
          hours: 2
        },
        {
          address: wallet.addresses[0].address,
          coins: 8,
          hours: 5
        }
      ];

      spendingService.createTransaction(
        wallet,
        ['address2'],
        null,
        [{
          address: address1,
          coins: amount1,
        },
        {
          address: address2,
          coins: amount2,
        }],
        {
          type: HoursSelectionTypes.Auto,
          ShareFactor: new BigNumber(0.5),
        },
        null
      ).subscribe((result: any) => {
        expect(result).toEqual({
          inputs: expectedTxInputs,
          outputs: expectedTxOutputs,
          hoursSent: new BigNumber(5),
          hoursBurned: new BigNumber(10),
          encoded: 'preparedTransaction'
        });
      });
    }));

    it('should reject if there are not enough coins while trying to send from specific addresses', fakeAsync(() => {
      const address1 = 'address_1';
      const amount1 = new BigNumber(10);
      const address2 = 'address_2';
      const amount2 = new BigNumber(10);

      spendingService.createTransaction(
        wallet,
        ['address2'],
        null,
        [{
          address: address1,
          coins: amount1,
        },
        {
          address: address2,
          coins: amount2,
        }],
        {
          type: HoursSelectionTypes.Auto,
          ShareFactor: new BigNumber(0.5),
        },
        null
      ).subscribe(() => fail('should be rejected'),
        (error) => expect(error.message).toBe('Not enough available Test Hours to perform transaction!')
      );
    }));

    it('should reject if there are not enough hours while trying to send from specific addresses', fakeAsync(() => {
      const address1 = 'address_1';
      const amount1 = new BigNumber(1);
      const hours1 = new BigNumber(10);
      const address2 = 'address_2';
      const amount2 = new BigNumber(1);
      const hours2 = new BigNumber(10);

      spendingService.createTransaction(
        wallet,
        ['address2'],
        null,
        [{
          address: address1,
          coins: amount1,
          hours: hours1,
        },
        {
          address: address2,
          coins: amount2,
          hours: hours2,
        }],
        {
          type: HoursSelectionTypes.Manual,
        },
        null
      ).subscribe(() => fail('should be rejected'),
        (error) => expect(error.message).toBe('Not enough available Test Hours to perform transaction!')
      );
    }));

    it('should return a correct observable while trying to send from specific outputs', fakeAsync(() => {
      const address1 = 'address_1';
      const amount1 = new BigNumber(3);
      const address2 = 'address_2';
      const amount2 = new BigNumber(3);

      const expectedTxInputs = [
        { hash: 'hash1', secret: 'secretKey1', address: 'address1', calculated_hours: 10, coins: 5 },
        { hash: 'hash4', secret: 'secretKey2', address: 'address2', calculated_hours: 0, coins: 5 }
      ];

      const expectedTxOutputs = [
        {
          address: address1,
          coins: amount1.toNumber(),
          hours: 1
        },
        {
          address: address2,
          coins: amount2.toNumber(),
          hours: 1
        },
        {
          address: wallet.addresses[0].address,
          coins: 4,
          hours: 3
        }
      ];

      spendingService.createTransaction(
        wallet,
        null,
        ['hash4', 'hash1'],
        [{
          address: address1,
          coins: amount1,
        },
        {
          address: address2,
          coins: amount2,
        }],
        {
          type: HoursSelectionTypes.Auto,
          ShareFactor: new BigNumber(0.5),
        },
        null
      ).subscribe((result: any) => {
        expect(result).toEqual({
          inputs: expectedTxInputs,
          outputs: expectedTxOutputs,
          hoursSent: new BigNumber(2),
          hoursBurned: new BigNumber(5),
          encoded: 'preparedTransaction'
        });
      });
    }));

    it('should reject if there are not enough coins while trying to send from specific outputs', fakeAsync(() => {
      const address1 = 'address_1';
      const amount1 = new BigNumber(10);
      const address2 = 'address_2';
      const amount2 = new BigNumber(10);

      spendingService.createTransaction(
        wallet,
        null,
        ['hash4', 'hash1'],
        [{
          address: address1,
          coins: amount1,
        },
        {
          address: address2,
          coins: amount2,
        }],
        {
          type: HoursSelectionTypes.Auto,
          ShareFactor: new BigNumber(0.5),
        },
        null
      ).subscribe(() => fail('should be rejected'),
        (error) => expect(error.message).toBe('Not enough available Test Hours to perform transaction!')
      );
    }));

    it('should reject if there are not enough hours while trying to send from specific outputs', fakeAsync(() => {
      const address1 = 'address_1';
      const amount1 = new BigNumber(10);
      const hours1 = new BigNumber(10);
      const address2 = 'address_2';
      const amount2 = new BigNumber(10);
      const hours2 = new BigNumber(10);

      spendingService.createTransaction(
        wallet,
        null,
        ['hash4', 'hash1'],
        [{
          address: address1,
          coins: amount1,
          hours: hours1,
        },
        {
          address: address2,
          coins: amount2,
          hours: hours2,
        }],
        {
          type: HoursSelectionTypes.Manual,
        },
        null
      ).subscribe(() => fail('should be rejected'),
        (error) => expect(error.message).toBe('Not enough available Test Hours to perform transaction!')
      );
    }));

    it('should add an extra output if it helps the user to save hours', fakeAsync(() => {
      const address = 'address';
      const amount = new BigNumber(20);

      const expectedTxInputs = [
        { hash: 'hash3', secret: 'secretKey1', address: 'address1', calculated_hours: 50, coins: 20 },
        { hash: 'hash4', secret: 'secretKey2', address: 'address2', calculated_hours: 0, coins: 5 },
      ];

      const expectedTxOutputs = [
        {
          address: address,
          coins: amount.toNumber(),
          hours: 12
        },
        {
          address: wallet.addresses[0].address,
          coins: 5,
          hours: 13
        }
      ];

      spendingService.createTransaction(
        wallet,
        null,
        null,
        [{
          address: address,
          coins: amount,
        }],
        {
          type: HoursSelectionTypes.Auto,
          ShareFactor: new BigNumber(0.5),
        },
        null
      ).subscribe((result: any) => {
        expect(result).toEqual({
          inputs: expectedTxInputs,
          outputs: expectedTxOutputs,
          hoursSent: new BigNumber(12),
          hoursBurned: new BigNumber(25),
          encoded: 'preparedTransaction'
        });
      });
    }));
  });

  describe('outputsWithWallets', () => {
    it('should return outputs by wallet', fakeAsync(() => {
      const walletAddress1 = 'address 1';
      const walletAddress2 = 'address 2';
      const wallet1: Wallet = Object.assign(createWallet(), { addresses: [ createAddress(walletAddress1), createAddress(walletAddress2) ] });
      const walletAddress3 = 'address 3';
      const wallet2: Wallet = Object.assign(createWallet(), { addresses: [ createAddress(walletAddress3) ] });

      spyOnProperty(walletService, 'currentWallets', 'get').and.returnValue( of([wallet1, wallet2]) );

      const output1 = createRequestOutput('hash3', walletAddress3, '33', 3 );
      const output2 = createRequestOutput('hash2', walletAddress2, '22', 2 );
      const output3 = createRequestOutput('hash1', walletAddress1, '11', 1 );

      spyApiService.get.and.returnValue(of({ head_outputs: [ output1, output2, output3 ] }));

      const expectedWallet1 = Object.assign(createWallet(), { addresses: [ createAddress(walletAddress1), createAddress(walletAddress2) ] });
      const expectedWallet2 = Object.assign(createWallet(), { addresses: [ createAddress(walletAddress3) ] });

      expectedWallet1.addresses[0].outputs = [ output3 ];
      expectedWallet1.addresses[1].outputs = [ output2 ];
      expectedWallet2.addresses[0].outputs = [ output1 ];

      const expectedOutputs = [ expectedWallet1, expectedWallet2 ];

      spendingService.outputsWithWallets()
        .subscribe((walletOutputs: any[]) => {
          expect(walletOutputs).toEqual(expectedOutputs);
        });
    }));
  });
});

function createOutput(address: string, hash: string, coins: BigNumber = new BigNumber(10), calculated_hours: BigNumber = new BigNumber(100)): Output {
  return {
    address: address,
    coins: coins,
    hash: hash,
    calculated_hours: calculated_hours
  };
}

function createRequestOutput(hash: string, address: string, coins = '10', calculated_hours = 100): GetOutputsRequestOutput {
  return {
    hash: hash,
    src_tx: 'src_tx: string',
    address: address,
    coins: coins,
    calculated_hours: calculated_hours
  };
}
