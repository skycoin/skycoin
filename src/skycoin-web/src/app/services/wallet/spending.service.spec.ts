import { TestBed, fakeAsync } from '@angular/core/testing';
import { Observable, of } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import BigNumber from 'bignumber.js';

import { SpendingService, HoursSelectionTypes } from './spending.service';
import { WalletService } from './wallet.service';
import { ApiService } from '../api.service';
import { CipherProvider } from '../cipher.provider';
import { CoinService } from '../coin.service';
import { MockCoinService, MockWalletService, MockGlobalsService, MockBlockchainService, MockHwWalletService } from '../../utils/test-mocks';
import { createWallet, createAddress } from './wallet.service.spec';
import { GlobalsService } from '../globals.service';
import {
  Wallet,
  TransactionOutput,
  TransactionInput,
  Output,
  GetOutputsRequestOutput,
  PreviewTransaction } from '../../app.datatypes';
import { BlockchainService } from '../blockchain.service';
import { HwWalletService } from '../hw-wallet.service';
import { BaseCoin } from '../../coins/basecoin';

describe('SpendingService', () => {
  let store: Record<string, string> = {};
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
        { provide: BlockchainService, useClass: MockBlockchainService },
        { provide: HwWalletService, useClass: MockHwWalletService }
      ]
    });

    spendingService = TestBed.inject(SpendingService);
    walletService = TestBed.inject(WalletService);
    spyApiService = TestBed.inject(ApiService) as jasmine.SpyObj<ApiService>;
    spyCipherProvider = TestBed.inject(CipherProvider) as jasmine.SpyObj<CipherProvider>;
    spyTranslateService = TestBed.inject(TranslateService) as jasmine.SpyObj<TranslateService>;
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

        return undefined;
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

  describe('createTransaction on a server wallet', () => {
    let post: jasmine.Spy;

    beforeEach(() => {
      (TestBed.inject(GlobalsService) as any as MockGlobalsService).setNodeVersion('0.27.0');
      post = jasmine.createSpy('post').and.callFake((url: string) => {
        if (url === 'transaction') {
          return of({ data: {
            encoded_transaction: 'unsigned',
            transaction: {
              fee: '1',
              inputs: [{ uxid: 'hash1', address: 'address1', calculated_hours: '2', coins: '5' }],
              outputs: [{ address: 'dest', coins: '5', hours: '1' }],
            },
          } });
        }

        return of({ data: { encoded_transaction: 'signed' } });
      });
      (spyApiService as any).post = post;
    });

    function send(password?: string) {
      const wallet = Object.assign(createWallet(), { filename: 'w.wlt', addresses: [createAddress('address1', '')] });
      let encoded = '';
      spendingService.createTransaction(
        wallet, null, null,
        [{ address: 'dest', coins: new BigNumber(5) }],
        { type: HoursSelectionTypes.Auto, ShareFactor: new BigNumber(0.5) },
        null, password,
      ).subscribe(tx => encoded = (tx as PreviewTransaction).encoded);

      return encoded;
    }

    it('signs with the wallet service and passes the password', fakeAsync(() => {
      expect(send('secret')).toBe('signed');
      const signCall = post.calls.all().find(c => c.args[0] === 'wallet/transaction/sign')!;
      expect(signCall.args[1]).toEqual({
        wallet_id: 'w.wlt',
        encoded_transaction: 'unsigned',
        input_addresses: ['address1'],
        password: 'secret',
      });
    }));

    it('sends an empty password for an unencrypted wallet', fakeAsync(() => {
      send();
      const signCall = post.calls.all().find(c => c.args[0] === 'wallet/transaction/sign')!;
      expect(signCall.args[1].password).toBe('');
    }));
  });

  describe('offline transactions', () => {
    const created = { data: {
      encoded_transaction: 'unsigned',
      transaction: {
        fee: '1',
        inputs: [{ uxid: 'hash1', address: 'address1', calculated_hours: '2', coins: '5' }],
        outputs: [{ address: 'dest', coins: '1.5', hours: '1' }],
      },
    } };
    let post: jasmine.Spy;

    beforeEach(() => {
      (TestBed.inject(GlobalsService) as any as MockGlobalsService).setNodeVersion('0.27.0');
      spyTranslateService.instant.and.callFake((key: any) => key);
      spyCipherProvider.prepareTransaction.and.returnValue(of('signedInBrowser'));
      post = jasmine.createSpy('post').and.callFake((url: string) =>
        of(url === 'wallet/transaction/sign' ? { data: { encoded_transaction: 'signedOnServer' } } : created));
      (spyApiService as any).post = post;
    });

    it('creates an unsigned transaction without signing or keys', fakeAsync(() => {
      const wallet = Object.assign(createWallet(), { addresses: [createAddress('address1', 'secretKey1')] });
      let tx: any;
      spendingService.createTransaction(
        wallet, null, null, [{ address: 'dest', coins: new BigNumber(1.5) }],
        { type: HoursSelectionTypes.Auto, ShareFactor: new BigNumber(0.5) }, null, undefined, true,
      ).subscribe(t => tx = t);

      expect(tx.encoded).toBe('unsigned');
      expect(tx.inputs[0].secret).toBe('');
      expect(spyCipherProvider.prepareTransaction).not.toHaveBeenCalled();
      expect(post.calls.all().map(c => c.args[0])).toEqual(['transaction']);
    }));

    it('refuses an unsigned transaction on a node older than 0.26', fakeAsync(() => {
      (TestBed.inject(GlobalsService) as any as MockGlobalsService).setNodeVersion('0.25.0');
      let error = '';
      spendingService.createTransaction(
        createWallet(), null, null, [{ address: 'dest', coins: new BigNumber(1) }],
        { type: HoursSelectionTypes.Auto, ShareFactor: new BigNumber(0.5) }, null, undefined, true,
      ).subscribe({ error: e => error = e.message });

      expect(error).toBe('offline-tx.node-too-old');
    }));

    it('signs a pasted transaction with a server wallet', fakeAsync(() => {
      const wallet = Object.assign(createWallet(), { filename: 'w.wlt', addresses: [] });
      let signed = '';
      spendingService.signTransaction(wallet, 'secret', ' unsigned \n').subscribe(s => signed = s);

      expect(signed).toBe('signedOnServer');
      expect(post.calls.argsFor(0)).toEqual(['transaction/verify', { encoded_transaction: 'unsigned', unsigned: true }, { json: true }, true]);
      expect(post.calls.argsFor(1)[1]).toEqual({
        wallet_id: 'w.wlt',
        encoded_transaction: 'unsigned',
        input_addresses: ['address1'],
        password: 'secret',
      });
    }));

    it('signs a pasted transaction with the keys of a browser wallet', fakeAsync(() => {
      const wallet = Object.assign(createWallet(), { addresses: [createAddress('address1', 'secretKey1')] });
      let signed = '';
      spendingService.signTransaction(wallet, undefined, 'unsigned').subscribe(s => signed = s);

      expect(signed).toBe('signedInBrowser');
      expect(spyCipherProvider.prepareTransaction).toHaveBeenCalledWith(
        [{ hash: 'hash1', secret: 'secretKey1', address: 'address1', calculated_hours: '2', coins: '5' } as any],
        [{ address: 'dest', coins: 1500000, hours: 1 }],
      );
    }));

    it('refuses a pasted transaction whose inputs the browser wallet cannot sign', fakeAsync(() => {
      const wallet = Object.assign(createWallet(), { addresses: [createAddress('address1', '')] });
      let error = '';
      spendingService.signTransaction(wallet, undefined, 'unsigned').subscribe({ error: e => error = e.message });

      expect(error).toBe('offline-tx.not-in-wallet');
      expect(spyCipherProvider.prepareTransaction).not.toHaveBeenCalled();
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

  describe('chunkedOutputsGet (via private access)', () => {
    let chunkedOutputsGetFn: (addressesCsv: string) => Observable<any>;

    beforeEach(() => {
      chunkedOutputsGetFn = (spendingService as any).chunkedOutputsGet.bind(spendingService);
    });

    it('should make a single GET request when addresses fit in one chunk', fakeAsync(() => {
      const mockResponse = {
        head_outputs: [createRequestOutput('hash1', 'addr1')],
        outgoing_outputs: [],
        incoming_outputs: [],
      };
      spyApiService.get.and.returnValue(of(mockResponse));

      chunkedOutputsGetFn('addr1,addr2').subscribe((result: any) => {
        expect(result).toEqual(mockResponse);
        expect(spyApiService.get).toHaveBeenCalledWith('outputs', { addrs: 'addr1,addr2' });
      });
    }));

    it('should merge head_outputs from multiple chunks', fakeAsync(() => {
      const response1 = {
        head_outputs: [createRequestOutput('hash1', 'addr1')],
        outgoing_outputs: [createRequestOutput('hash2', 'addr1')],
        incoming_outputs: [],
      };
      const response2 = {
        head_outputs: [createRequestOutput('hash3', 'addr2')],
        outgoing_outputs: [],
        incoming_outputs: [createRequestOutput('hash4', 'addr2')],
      };
      spyApiService.get.and.returnValues(of(response1), of(response2));

      // Create a CSV string long enough to force splitting (> 1800 chars)
      const addr = '2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6';
      const longCsv = Array(60).fill(addr).join(',');

      chunkedOutputsGetFn(longCsv).subscribe((result: any) => {
        expect(result.head_outputs.length).toEqual(2);
        expect(result.outgoing_outputs.length).toEqual(1);
        expect(result.incoming_outputs.length).toEqual(1);
        expect(result.head_outputs[0].hash).toEqual('hash1');
        expect(result.head_outputs[1].hash).toEqual('hash3');
      });
    }));

    it('should handle empty output arrays in responses', fakeAsync(() => {
      const response1 = {
        head_outputs: [],
        outgoing_outputs: [],
        incoming_outputs: [],
      };
      const response2 = {
        head_outputs: [createRequestOutput('hash1', 'addr1')],
        outgoing_outputs: null,
        incoming_outputs: null,
      };
      spyApiService.get.and.returnValues(of(response1), of(response2));

      const addr = '2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6';
      const longCsv = Array(60).fill(addr).join(',');

      chunkedOutputsGetFn(longCsv).subscribe((result: any) => {
        expect(result.head_outputs.length).toEqual(1);
        expect(result.outgoing_outputs.length).toEqual(0);
        expect(result.incoming_outputs.length).toEqual(0);
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

describe('SpendingService gateway sends', () => {
  let spending: SpendingService;
  let post: jasmine.Spy;
  let coin: BaseCoin;

  beforeEach(() => {
    post = jasmine.createSpy('post').and.returnValue(of({ txid: '0xabc' }));
    TestBed.configureTestingModule({
      providers: [
        SpendingService,
        { provide: WalletService, useClass: MockWalletService },
        { provide: ApiService, useValue: { post, get: () => of({}) } },
        { provide: CipherProvider, useValue: {} },
        { provide: TranslateService, useValue: { instant: (k: string) => k } },
        { provide: CoinService, useValue: { currentCoin: { subscribe: (f: (c: BaseCoin) => void) => f(coin) } } },
        { provide: GlobalsService, useClass: MockGlobalsService },
        { provide: BlockchainService, useClass: MockBlockchainService },
        { provide: HwWalletService, useClass: MockHwWalletService },
      ],
    });
  });

  it('sends Ether amounts in wei, as a string', fakeAsync(() => {
    coin = BaseCoin.fromServerData({ id: 1, coinSymbol: 'ETH', coinType: 'ethereum', decimals: 18 });
    spending = TestBed.inject(SpendingService);
    let txid = '';
    spending.sendGateway({ filename: 'e.wlt' } as any, '0xto', new BigNumber('1.5'), 'pw').subscribe(t => txid = t);
    expect(txid).toBe('0xabc');
    expect(post).toHaveBeenCalledWith('eth/send', {
      wallet_id: 'e.wlt', password: 'pw', to: '0xto', amount: '1500000000000000000',
    }, { json: true });
  }));

  it('sends Bitcoin in satoshis with the fee rate', fakeAsync(() => {
    coin = BaseCoin.fromServerData({ id: 2, coinSymbol: 'BTC', coinType: 'bitcoin' });
    spending = TestBed.inject(SpendingService);
    spending.sendGateway({ filename: 'b.wlt' } as any, 'bc1q', new BigNumber('0.0001'), undefined, 7).subscribe();
    expect(post).toHaveBeenCalledWith('btc/send', {
      wallet_id: 'b.wlt', password: '', destinations: [{ address: 'bc1q', coins: 10000 }], fee_rate: 7,
    }, { json: true });
  }));
});
