import { TestBed } from '@angular/core/testing';
import { HttpClientModule } from '@angular/common/http';

import { CipherProvider, InitializationResults } from './cipher.provider';
import { Address, TransactionInput, TransactionOutput } from '../app.datatypes';
import { convertAsciiToHexa } from '../utils/converters';

describe('CipherProvider', () => {
  let cipherProvider: CipherProvider;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [ HttpClientModule ],
      providers: [
        CipherProvider
      ]
    });

    cipherProvider = TestBed.get(CipherProvider);
  });

  it('should be created', () => {
    expect(cipherProvider).toBeTruthy();
  });

  it('should be initialized', done => {
    cipherProvider.initialize().subscribe(response => {
      expect(response === InitializationResults.Ok).toBeTruthy();

      done();
    });
  });

  it('should generate address', done => {
    const expectedAddress: Address = createAddress();

    cipherProvider.generateAddress(convertAsciiToHexa('test')).subscribe(addr => {
      expect(addr.address).toEqual(expectedAddress);

      done();
    });
  });

  it('should throw an error in case of any error generating an address', done => {
    cipherProvider.generateAddress('test').subscribe(addr => {
      fail('The operation did not fail.');
    }, () => done());
  });

  it('should prepare transaction', done => {
    cipherProvider.prepareTransaction([createTransactionInput()], [createTransactionOutput()]).subscribe((rawTx: string) => {
      expect(rawTx.length).toBeGreaterThan(20);

      done();
    });
  });

  it('should throw an error in case of any error preparing a transaction', done => {
    cipherProvider.prepareTransaction([createBadTransactionInput()], [createTransactionOutput()]).subscribe((rawTx: string) => {
      fail('The operation did not fail.');
    }, () => done());
  });
});

function createAddress(): Address {
  return {
    address: '2ccsfuKpQ1EcqgQXPZyYZNW8qEDuKXZh9fj',
    public_key: '024119d3c96e1cb4270fc332ff56b9fef66c4647c6f0b9a287821889e53eb235bc',
    secret_key: '2b98a11bc65fba1dda70dbb61d6f21025c422590716a02697f06ba583e46bcf5'
  };
}

function createTransactionInput(): TransactionInput {
  return {
    hash: '1fe7d0625540d730a396962fe616053f0e1385d10f3f2702f8a490e3c4cee075',
    secret: '001aa9e416aff5f3a3c7f9ae0811757cf54f393d50df861f5c33747954341aa7'
  };
}

function createBadTransactionInput(): TransactionInput {
  return {
    hash: '1',
    secret: '1'
  };
}

function createTransactionOutput(address = '2e1erPpaxNVC37PkEv3n8PESNw2DNr5aJNy', coins = 10000, hours = 0): TransactionOutput {
  return {
    address: address,
    coins: coins,
    hours: hours
  };
}
