import { firstValueFrom, of } from 'rxjs';

import { Wallet } from '../../../app.datatypes';
import { EncryptionService } from '../../encryption.service';
import { BrowserWalletBackend } from './browser-wallet-backend';
import { ServerWalletBackend } from './server-wallet-backend';

describe('temporary wallets', () => {
  const translate = { instant: (key: string) => key };

  function browserBackend(): BrowserWalletBackend {
    return new BrowserWalletBackend(null as any, translate as any, null as any, new EncryptionService());
  }

  function wallet(label: string, fields: Partial<Wallet> = {}): Wallet {
    return { label, seed: label + ' seed', addresses: [{ address: label + '-addr' }], coinId: 1, ...fields };
  }

  // Saves reach localStorage on a promise chain.
  const settle = () => new Promise(resolve => setTimeout(resolve));

  beforeEach(() => localStorage.clear());
  afterEach(() => localStorage.clear());

  it('keeps a temporary browser wallet out of storage but in the loaded list', async () => {
    const backend = browserBackend();
    backend.save([wallet('kept'), wallet('temp', { temporary: true })]);
    await settle();

    const stored = JSON.parse(localStorage.getItem('wallets')!);
    expect(stored.map((w: Wallet) => w.label)).toEqual(['kept']);
    expect(localStorage.getItem('wallets')).not.toContain('temp seed');

    const loaded = await firstValueFrom(backend.load());
    expect(loaded!.map(w => w.label)).toEqual(['kept', 'temp']);
    expect(loaded![1].temporary).toBe(true);
  });

  it('loads a temporary wallet when nothing is stored', async () => {
    const backend = browserBackend();
    backend.save([wallet('temp', { temporary: true })]);
    await settle();

    expect(JSON.parse(localStorage.getItem('wallets')!)).toEqual([]);
    expect((await firstValueFrom(backend.load()))!.map(w => w.label)).toEqual(['temp']);
  });

  it('forgets a temporary wallet when the page loads again', async () => {
    browserBackend().save([wallet('kept'), wallet('temp', { temporary: true })]);
    await settle();

    const reloaded = await firstValueFrom(browserBackend().load());
    expect(reloaded!.map(w => w.label)).toEqual(['kept']);
  });

  it('forgets a temporary wallet once it is deleted from the list', async () => {
    const backend = browserBackend();
    backend.save([wallet('temp', { temporary: true })]);
    backend.save([]);
    await settle();

    expect(await firstValueFrom(backend.load())).toEqual([]);
  });

  it('loads a temporary wallet service wallet with createTemp, without a password', async () => {
    const api = { post: jasmine.createSpy('post').and.returnValue(of({ meta: { label: 't', filename: 't.wlt', temp: true }, entries: [] })) };
    const backend = new ServerWalletBackend(api as any, () => 1);

    const created = await firstValueFrom(backend.create('t', 'seed', 1, [], 'deterministic', undefined, undefined, 'pw', true));

    expect(api.post.calls.argsFor(0)[0]).toBe('wallet/createTemp');
    expect(api.post.calls.argsFor(0)[1].password).toBeUndefined();
    expect(created.temporary).toBe(true);
  });

  it('still stores a wallet service wallet with create', async () => {
    const api = { post: jasmine.createSpy('post').and.returnValue(of({ meta: { label: 'w', filename: 'w.wlt' }, entries: [] })) };
    const backend = new ServerWalletBackend(api as any, () => 1);

    const created = await firstValueFrom(backend.create('w', 'seed', 1, [], 'deterministic'));

    expect(api.post.calls.argsFor(0)[0]).toBe('wallet/create');
    expect(created.temporary).toBeUndefined();
  });
});
