import { TestBed, fakeAsync } from '@angular/core/testing';
import { Observable, of } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { BigNumber } from 'bignumber.js';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';

import { WalletService } from './wallet.service';
import { SpendingService, HoursSelectionTypes } from './spending.service';
import { ApiService } from '../api.service';
import { CipherProvider, InitializationResults } from '../cipher.provider';
import { Wallet, Address, TransactionOutput, TransactionInput, Output, Balance } from '../../app.datatypes';
import { CoinService } from '../coin.service';
import { MockCoinService, MockGlobalsService } from '../../utils/test-mocks';
import { createWallet } from './wallet.service.spec';
import { GlobalsService } from '../globals.service';
import { BlockchainService } from '../blockchain.service';
import { BalanceService } from './balance.service';

describe('WalletService with cipher:', () => {
  let store = {};
  let walletService: WalletService;
  let spendingService: SpendingService;
  let cipherProvider: CipherProvider;
  let spyApiService:  jasmine.SpyObj<ApiService>;

  beforeEach(() => {
    spyOn(localStorage, 'setItem').and.callFake((key, value) => store[key] = value);
    spyOn(localStorage, 'getItem').and.callFake((key) => store[key]);

    TestBed.configureTestingModule({
    imports: [],
    providers: [
        WalletService,
        CipherProvider,
        SpendingService,
        BlockchainService,
        BalanceService,
        {
            provide: ApiService,
            useValue: jasmine.createSpyObj('ApiService', {
                'getOutputs': of([]),
                'postTransaction': of(''),
                'get': of([])
            })
        },
        {
            provide: TranslateService,
            useValue: jasmine.createSpyObj('TranslateService', ['instant'])
        },
        { provide: CoinService, useClass: MockCoinService },
        { provide: GlobalsService, useClass: MockGlobalsService },
        provideHttpClient(withInterceptorsFromDi())
    ]
});

    walletService = TestBed.inject(WalletService);
    spendingService = TestBed.inject(SpendingService);
    cipherProvider = TestBed.inject(CipherProvider);
    spyApiService = TestBed.inject(ApiService);
  });

  afterEach(() => {
    store = {};
  });

  it('wallet service should be created', () => {
    expect(walletService).toBeTruthy();
  });

  it('spending service should be created', () => {
    expect(spendingService).toBeTruthy();
  });

  it('cipher service should be created', () => {
    expect(cipherProvider).toBeTruthy();
  });

  it('cipher service should be initialized', done => {
    cipherProvider.initialize().subscribe(response => {
      expect(response === InitializationResults.Ok).toBeTruthy();

      done();
    });
  });

  describe('cipher should generate an address', () => {
    it('on add address to wallet', () => {
      const wallet = createWallet();
      wallet.nextSeed = 'f0cecf181dc39bc99d0b9bec3fdce293c22af99af20320974eab1217d0d75b69';
      const expectedWallet = createWallet();
      expectedWallet.nextSeed = '3043a0ac64d63f4eaff7f7c66c30ca5ab544f87f04bbe0f60d189bb8ff125d52';
      const newAddress =
        createAddress(
          '2ewXR5RsixF8yJQSuTSUTKWbjeb7y9brjmE',
          '03841048e4989e456987627083ef3b04626c9afde52fa6137453d6aabe74971b87',
          '0f901b534792e62cf7c4b6584f430fc22f4f2e7f686197441c96a8cafb098c9a'
        );

      expectedWallet.addresses.push(newAddress);

      spyApiService.get.and.callFake(() => {
        return of(createBalance());
      });

      walletService.addAddress(wallet)
        .subscribe(() => {
          expect(wallet).toEqual(expectedWallet);
        });
    });
  });

  describe('cipher createTransaction', () => {
    it('should return the correct transaction inputs and outputs', fakeAsync(() => {
      const amount = new BigNumber(10000);
      const destinationAddress = '2e1erPpaxNVC37PkEv3n8PESNw2DNr5aJNy';
      const addresses = [
        createAddress()
      ];

      const wallet: Wallet = Object.assign(createWallet(), { addresses: addresses });
      const expectedTxInputs: TransactionInput[] = [
        createTransactionInput()
      ];

      const expectedTxOutputs: TransactionOutput[] = [
        createTransactionOutput(destinationAddress)
      ];

      const outputs: Output[] = [
        createOutput(addresses[0].address, amount, new BigNumber(1)),
      ];

      spyApiService.get.and.returnValue(of({ head_outputs: outputs }));

      spendingService.createTransaction(
        wallet,
        null,
        null,
        [{
          address: destinationAddress,
          coins: amount,
        }],
        {
          type: HoursSelectionTypes.Auto,
          ShareFactor: new BigNumber(0.5),
        },
        null
      )
        .subscribe(
          (result: any) => {
            expect(result.inputs).toEqual(expectedTxInputs);
            expect(result.outputs).toEqual(expectedTxOutputs);
          });
    }));

    it('should rejected transaction for an invalid destination address\'', fakeAsync(() => {
      const amount = new BigNumber(10000);
      const wrongDestinationAddress = '2e1erPpaxNVC37PkEv3n8PESNw2DNr5aJNz';
      const addresses = [
        createAddress()
      ];

      const wallet: Wallet = Object.assign(createWallet(), { addresses: addresses });

      const outputs: Output[] = [
        createOutput(addresses[0].address, amount, new BigNumber(1)),
      ];

      spyApiService.get.and.returnValue(of({ head_outputs: outputs }));

      spendingService.createTransaction(
        wallet,
        null,
        null,
        [{
          address: wrongDestinationAddress,
          coins: amount,
        }],
        {
          type: HoursSelectionTypes.Auto,
          ShareFactor: new BigNumber(0.5),
        },
        null
      )
        .subscribe(
          () => {
            fail('should fail');
          },
          (error) => expect(error.message).toEqual('Error: Invalid checksum')
        );
    }));
  });
});

function createAddress(
  address = '2uATq4pdSb8Ka1YKSAAbp6Npehs3QQqTnb',
  public_key = '030a797a31100d3a7b5b403f551975e9a12f93b4d4e9e44b402b84832e0c7b89d2',
  secret_key = '20c3db0e1f3b95d98d1f78d73c134af8f1b5dd34cc05f053da94c20d72558862'
): Address {
  return {
    address: address,
    public_key: public_key,
    secret_key: secret_key
  };
}

function createOutput(address: string, coins: BigNumber = new BigNumber(10), calculated_hours: BigNumber = new BigNumber(100)): Output {
  return {
    address: address,
    coins: coins,
    hash: '1fe7d0625540d730a396962fe616053f0e1385d10f3f2702f8a490e3c4cee075',
    calculated_hours: calculated_hours
  };
}

function createBalance(coins = 0, hours = 0): Balance {
  return {
    confirmed: {
      coins: coins,
      hours: hours
    },
    predicted: {
      coins: coins,
      hours: hours
    },
    addresses: {
      'address': {
        confirmed: {
          coins: coins,
          hours: hours
        },
        predicted: {
          coins: coins,
          hours: hours
        }
      }
    }
  };
}

function createTransactionInput(): TransactionInput {
  return {
    hash: '1fe7d0625540d730a396962fe616053f0e1385d10f3f2702f8a490e3c4cee075',
    secret: '20c3db0e1f3b95d98d1f78d73c134af8f1b5dd34cc05f053da94c20d72558862',
    address: '2uATq4pdSb8Ka1YKSAAbp6Npehs3QQqTnb',
    coins: 10000,
    calculated_hours: 1
  };
}

function createTransactionOutput(address: string, coins = 10000, hours = 0): TransactionOutput {
  return {
    address: address,
    coins: coins,
    hours: hours
  };
}
