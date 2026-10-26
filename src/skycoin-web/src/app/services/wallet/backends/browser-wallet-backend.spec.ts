import { EventEmitter } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { provideHttpClient, withInterceptorsFromDi, withXhr } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';

import { Wallet } from '../../../app.datatypes';
import { EncryptionService } from '../../encryption.service';
import { CipherProvider, pqSeedState } from '../../cipher.provider';
import { convertAsciiToHexa } from '../../../utils/converters';
import { loadCipherWasm } from '../../../utils/wasm-test-utils';
import { BrowserWalletBackend } from './browser-wallet-backend';

describe('BrowserWalletBackend getSeed', () => {
  const encryption = new EncryptionService();
  const translate = { instant: (key: string) => key };
  const backend = new BrowserWalletBackend(null as any, translate as any, null as any, encryption);

  function wallet(fields: Partial<Wallet>): Wallet {
    return { label: 'w', addresses: [], ...fields };
  }

  it('returns a seed it holds without a password', async () => {
    expect(await firstValueFrom(backend.getSeed(wallet({ seed: 'one two three' }), null))).toBe('one two three');
  });

  it('decrypts a seed kept under the wallet password', async () => {
    const encryptedSeed = await encryption.encrypt('one two three', 'hunter2');
    expect(await firstValueFrom(backend.getSeed(wallet({ encryptedSeed }), 'hunter2'))).toBe('one two three');
    await expectAsync(firstValueFrom(backend.getSeed(wallet({ encryptedSeed }), 'wrong'))).toBeRejectedWithError('Incorrect password');
  });

  it('fails for a wallet whose seed it does not hold', async () => {
    await expectAsync(firstValueFrom(backend.getSeed(wallet({}), null))).toBeRejected();
  });
});

describe('BrowserWalletBackend post-quantum wallets', () => {
  const seed = 'abcdefg';
  const seedHex = convertAsciiToHexa(seed);
  const firstAddress = '3rqVbVSdzuWuciwp66RrkjJ51e4vuipzbY';
  const secondAddress = '2UxmE95vcSNuxengyEQBp12mbndgncgRs6a';
  const translate = { instant: (key: string) => key };
  let backend: BrowserWalletBackend;

  beforeAll(async () => {
    await loadCipherWasm();
  }, 60000);

  beforeEach(() => {
    localStorage.clear();
    TestBed.configureTestingModule({
      providers: [CipherProvider, provideHttpClient(withXhr(), withInterceptorsFromDi())],
    });
    backend = new BrowserWalletBackend(TestBed.inject(CipherProvider), translate as any, null as any, new EncryptionService());
  });

  afterEach(() => localStorage.clear());

  function createPQ(): Promise<Wallet> {
    return firstValueFrom(backend.create('pq', seed, 1, [], 'pq'));
  }

  it('creates a wallet whose state is the next pq index', async () => {
    const wallet = await createPQ();

    expect(wallet.walletType).toBe('pq');
    expect(wallet.addresses[0].address).toBe(firstAddress);
    expect(wallet.nextSeed).toBe(pqSeedState(seedHex, 1));
  });

  it('still makes deterministic wallets by default', async () => {
    const wallet = await firstValueFrom(backend.create('det', seed, 1, [], 'deterministic'));

    expect(wallet.walletType).toBe('deterministic');
    expect(wallet.addresses[0].address).not.toBe(firstAddress);
  });

  it('adds the next pq address', async () => {
    const wallet = await createPQ();
    await firstValueFrom(backend.addAddress(wallet));

    expect(wallet.addresses.map(a => a.address)).toEqual([firstAddress, secondAddress]);
    expect(wallet.nextSeed).toBe(pqSeedState(seedHex, 2));
  });

  it('unlocks a wallet that kept only its addresses, and rejects a wrong seed', async () => {
    const wallet = await createPQ();
    await firstValueFrom(backend.addAddress(wallet));
    const locked: Wallet = { label: 'pq', walletType: 'pq', addresses: wallet.addresses.map(a => ({ address: a.address })) };

    await expectAsync(firstValueFrom(backend.unlockWallet(locked, 'wrong', new EventEmitter<number>()))).toBeRejected();
    await firstValueFrom(backend.unlockWallet(locked, seed, new EventEmitter<number>()));

    expect(locked.seed).toBe(seed);
    expect(locked.addresses[1].secret_key).toBeTruthy();
    expect(locked.nextSeed).toBe(pqSeedState(seedHex, 2));
  });

  it('unlocks with the wallet password', async () => {
    const wallet = await createPQ();
    await firstValueFrom(backend.setWalletPassword(wallet, 'hunter2'));
    const locked: Wallet = {
      label: 'pq',
      walletType: 'pq',
      addresses: [{ address: firstAddress }],
      encryptedSeed: wallet.encryptedSeed,
      encryptedNextSeed: wallet.encryptedNextSeed,
    };

    await firstValueFrom(backend.unlockWalletWithPassword(locked, 'hunter2', new EventEmitter<number>(), s => s));

    expect(locked.seed).toBe(seed);
    expect(locked.nextSeed).toBe(pqSeedState(seedHex, 1));
  });

  it('keeps the wallet type across a save and a load', async () => {
    const wallet = await createPQ();
    backend.save([wallet]);
    await new Promise(resolve => setTimeout(resolve, 100));

    const loaded = (await firstValueFrom(backend.load()))!;

    expect(loaded[0].walletType).toBe('pq');
    expect(loaded[0].nextSeed).toBe(pqSeedState(seedHex, 1));
  });
});
