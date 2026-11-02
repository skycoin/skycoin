import { EventEmitter } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { provideHttpClient, withInterceptorsFromDi, withXhr } from '@angular/common/http';
import { firstValueFrom, of } from 'rxjs';

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

// The node's bip44 wallet gives these for the bip39 all-zero test mnemonic:
// `skycoin-cli walletCreate -t bip44` then `walletKeyExport -k xpub --path=0`.
describe('BrowserWalletBackend bip44 wallets', () => {
  const mnemonic = 'abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about';
  const node = {
    '': {
      xpub: 'xpub6Cjmbker6mxQGujSPMQvfFgsaSWmF5dL9dki1ZoNS39QWCmkZkguNHHjFxszEwYyVhpxCDkkb9B767LqHdcpEwXGAQpBiub5d532TA4RJGR',
      external: ['28RHxxgAsbCuTv5U9VgWrDGDUpoho2gbh66', '2bjtkT6qDBxhBSLsVK5TmRHfZBy35z7sMVq'],
      change: ['2RKPZ1rhvxjWFX4zWNSbtLW1QLqBTrokP7U'],
    },
    'a seed passphrase': {
      xpub: 'xpub6DRdckbBoazxp4uTz7ABHRVLAxdnBcM97wvrHgoGY1rDSHUPiHcwKiqLnEBprzitbWMPGbR6JuVjWEy4pFxDZ1eeowr1USk3st7JW7qiVvt',
      external: ['2CQQwfehE5ov1EfhDufvxnGBtRzPeBcZXVZ', '2f3fiFEs9wm6eZeHoAdfPqzqwaQVtUAy1ox'],
      change: ['2JpoaYoA9G14wngoR4VAdsnb4Zg3BGZjUJg'],
    },
  };
  const translate = { instant: (key: string) => key };
  let cipher: CipherProvider;
  let api: { get: jasmine.Spy };
  let backend: BrowserWalletBackend;

  beforeAll(async () => {
    await loadCipherWasm();
  }, 60000);

  beforeEach(() => {
    localStorage.clear();
    TestBed.configureTestingModule({
      providers: [CipherProvider, provideHttpClient(withXhr(), withInterceptorsFromDi())],
    });
    cipher = TestBed.inject(CipherProvider);
    api = { get: jasmine.createSpy('get').and.returnValue(of([])) };
    backend = new BrowserWalletBackend(cipher, translate as any, api as any, new EncryptionService());
  });

  afterEach(() => localStorage.clear());

  function create(passphrase: string): Promise<Wallet> {
    return firstValueFrom(backend.create('b', mnemonic, 1, [], 'bip44', passphrase));
  }

  it('gives the node wallet xpub and addresses, with and without a passphrase', async () => {
    for (const passphrase of Object.keys(node) as (keyof typeof node)[]) {
      const wallet = await create(passphrase);
      await firstValueFrom(backend.addAddress(wallet, 0, 0));
      await firstValueFrom(backend.addAddress(wallet, 0, 1));

      expect(wallet.walletType).toBe('bip44');
      expect(wallet.accounts![0].accountXpubKey).toBe(node[passphrase].xpub);
      expect(wallet.accounts![0].externalAddresses.map(a => a.address)).toEqual(node[passphrase].external);
      expect(wallet.accounts![0].changeAddresses.map(a => a.address)).toEqual(node[passphrase].change);
      expect(wallet.addresses.map(a => a.address)).toEqual([...node[passphrase].external, ...node[passphrase].change]);
    }
  });

  it('rejects a mnemonic with a bad checksum', async () => {
    await expectAsync(firstValueFrom(backend.create('b', mnemonic.replace(/about$/, 'abandon'), 1, [], 'bip44'))).toBeRejected();
  });

  it('keeps the xpub and both chains through an encrypted save, and signs from both chains after unlocking', async () => {
    const wallet = await create('a seed passphrase');
    await firstValueFrom(backend.addAddress(wallet, 0, 1));
    await firstValueFrom(backend.setWalletPassword(wallet, 'hunter2'));
    backend.save([wallet]);
    await new Promise(resolve => setTimeout(resolve, 100));

    const stored = localStorage.getItem('wallets')!;
    expect(stored).not.toContain('abandon');
    expect(stored).not.toContain('a seed passphrase');

    const loaded = (await firstValueFrom(backend.load()))![0];
    expect(loaded.accounts![0].accountXpubKey).toBe(node['a seed passphrase'].xpub);
    expect(loaded.addresses.map(a => a.address)).toEqual([node['a seed passphrase'].external[0], node['a seed passphrase'].change[0]]);
    expect(loaded.addresses.every(a => !a.secret_key)).toBeTrue();

    await expectAsync(firstValueFrom(backend.unlockWalletWithPassword(loaded, 'wrong', new EventEmitter<number>(), s => s))).toBeRejected();
    await firstValueFrom(backend.unlockWalletWithPassword(loaded, 'hunter2', new EventEmitter<number>(), s => s));
    expect(loaded.seedPassphrase).toBe('a seed passphrase');
    expect(loaded.addresses.every(a => !!a.secret_key)).toBeTrue();

    const inputs = loaded.addresses.map((a, i) => ({ hash: (i + 1).toString().repeat(64), secret: a.secret_key!, address: a.address }));
    const outputs = [{ address: node[''].external[0], coins: 1000000, hours: 1 }, { address: loaded.addresses[1].address, coins: 1000000, hours: 1 }];
    const tx = await firstValueFrom(cipher.prepareTransaction(inputs, outputs));
    expect(tx).toMatch(/^[0-9a-f]+$/);
  });

  it('refuses a seed that does not give the stored addresses', async () => {
    const wallet = await create('');
    const locked: Wallet = { ...wallet, seed: undefined, seedPassphrase: 'not it' };
    await expectAsync(firstValueFrom(backend.unlockWallet(locked, mnemonic, new EventEmitter<number>()))).toBeRejected();
  });

  it('finds used addresses past the ones it lists, up to the gap limit', async () => {
    const xpub = node[''].xpub;
    const external = (await firstValueFrom(cipher.bip44XPubAddresses(xpub, 0, 0, 30))).map(a => a.address);
    const change = (await firstValueFrom(cipher.bip44XPubAddresses(xpub, 1, 0, 5))).map(a => a.address);
    const used = new Set([external[6], external[24], change[2]]);
    api.get.and.callFake((_url: string, params: any) => of(params.addrs.split(',')
      .filter((a: string) => used.has(a))
      .map((a: string) => ({ txn: { outputs: [{ dst: a }] } }))));

    const wallet = await create('');
    const added = await firstValueFrom(backend.discoverBip44Addresses(wallet));

    expect(added).toBe(24 + 3);
    expect(wallet.accounts![0].externalAddresses.map(a => a.address)).toEqual(external.slice(0, 25));
    expect(wallet.accounts![0].changeAddresses.map(a => a.address)).toEqual(change.slice(0, 3));
    expect(wallet.addresses.length).toBe(28);
  });
});
