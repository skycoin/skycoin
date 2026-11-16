import { EventEmitter } from '@angular/core';
import { BehaviorSubject, Observable, firstValueFrom, of, throwError } from 'rxjs';

import { Wallet } from '../../app.datatypes';
import { convertAsciiToHexa } from '../../utils/converters';
import { loadCipherWasm } from '../../utils/wasm-test-utils';
import { CipherProvider, InitializationResults } from '../cipher.provider';
import { BrowserSweepService, browserSweepRulesKey, browserWalletId } from './browser-sweep.service';

/**
 * Runs the browser sweep against the real wasm cipher, the node's rule code,
 * with the node's API and the wallet service faked.
 */
describe('BrowserSweepService', () => {
  const seed = 'browser sweep spec seed';
  let cipher: CipherProvider;
  let chain: { address: string, secret: string, next: string }[];

  let outputs: any[];
  let outgoing: any[];
  let sentBy: string[];
  let injected: string[];
  let saved: number;
  let wallets: BehaviorSubject<Wallet[]>;
  let service: BrowserSweepService;

  beforeAll(async () => {
    await loadCipherWasm();
    const c = (window as any).SkycoinCipher;
    let state = convertAsciiToHexa(seed);
    chain = [];
    for (let i = 0; i < 4; i++) {
      const a = c.generateAddress(state);
      chain.push({ address: a.address, secret: a.secret, next: a.nextSeed });
      state = a.nextSeed;
    }
  });

  /** Derives the keys of every stored address of w, as the browser backend does. */
  function derive(w: Wallet): Observable<void> {
    w.addresses.forEach((a, i) => {
      if (a.address !== chain[i].address) {
        throw new Error('wrong seed');
      }
      a.secret_key = chain[i].secret;
    });
    w.nextSeed = chain[w.addresses.length - 1].next;

    return of(undefined);
  }

  function hash(n: number): string {
    return (n.toString(16) + '0'.repeat(64)).slice(0, 64);
  }

  beforeEach(() => {
    localStorage.removeItem(browserSweepRulesKey);
    cipher = new CipherProvider(null as any);
    spyOn(cipher, 'initialize').and.returnValue(of(InitializationResults.Ok));
    outputs = [{ hash: hash(1), address: chain[0].address, coins: '5.000000', calculated_hours: 40 }];
    outgoing = [];
    sentBy = [chain[0].address];
    injected = [];
    saved = 0;
    wallets = new BehaviorSubject<Wallet[]>([]);
    const api: any = {
      get: () => of({ head: { seq: 7 } }),
      post: (url: string, body: any) => {
        if (url === 'outputs') {
          const want = body.addrs.split(',');

          return of({ head_outputs: outputs.filter(o => want.includes(o.address)), outgoing_outputs: outgoing });
        }
        if (url === 'transactions') {
          return of(sentBy.map(a => ({ txn: { inputs: [{ owner: a }], outputs: [{ dst: a }] } })));
        }
        if (url === 'injectTransaction') {
          injected.push(body.rawtx);

          return of('txid' + injected.length);
        }

        return throwError(() => new Error('unexpected ' + url));
      },
    };
    const walletService: any = {
      wallets: wallets,
      isServerWallet: (w: Wallet) => !!w.filename,
      unlockWallet: (w: Wallet) => derive(w),
      unlockWalletWithPassword: (w: Wallet, pw: string, _: EventEmitter<number>) =>
        pw === 'pw' ? derive(w) : throwError(() => new Error('wrong password')),
      addAddress: (w: Wallet) => {
        w.addresses.push({ address: chain[w.addresses.length].address });

        return of(undefined);
      },
      saveWallets: () => saved++,
    };
    service = new BrowserSweepService(walletService, cipher, api);
  });

  function wallet(over: Partial<Wallet> = {}): Wallet {
    const w: Wallet = { label: 'B', walletType: 'deterministic', seed: seed, addresses: [{ address: chain[0].address }, { address: chain[1].address }], ...over };
    wallets.next([w]);

    return w;
  }

  it('sweeps an address that sent to the next unused address of the wallet, burning every hour', async () => {
    const w = wallet();
    const id = browserWalletId(w);
    const st = await firstValueFrom(service.setRule({ wallet: id, enabled: true, destination: { kind: 'wallet' } }));
    expect(injected.length).toBe(1);
    expect(st.status.active).toBeTrue();
    expect(st.status.exposed_addresses).toBe(1);
    expect(st.status.destinations).toEqual([chain[1].address]);
    expect(st.status.pending_txids).toEqual(['txid1']);
    expect(st.status.last_error).toBeUndefined();

    expect(JSON.parse(localStorage.getItem(browserSweepRulesKey)!)[id].enabled).toBeTrue();

    // The sweep is in the pool, so the next pass sends nothing new.
    outgoing = [outputs[0]];
    await service.pass(id);
    expect(injected.length).toBe(1);

    // It confirmed.
    outputs = [];
    outgoing = [];
    await service.pass(id);
    const after = await firstValueFrom(service.get(id));
    expect(after.status.confirmed_sweeps).toBe(1);
    expect(after.status.pending_txids).toBeUndefined();
  });

  it('derives and stores a destination before sending to it', async () => {
    const w = wallet({ addresses: [{ address: chain[0].address }] });
    const id = browserWalletId(w);
    const st = await firstValueFrom(service.setRule({ wallet: id, enabled: true, destination: { kind: 'wallet' } }));
    expect(st.status.destinations).toEqual([chain[1].address]);
    expect(w.addresses.map(a => a.address)).toEqual([chain[0].address, chain[1].address]);
    expect(saved).toBe(1);
    expect(injected.length).toBe(1);
  });

  it('waits for an unlock on an encrypted wallet and keeps the password out of storage', async () => {
    const w = wallet({ seed: undefined, encryptedSeed: 'encrypted' });
    const id = browserWalletId(w);
    let st = await firstValueFrom(service.setRule({ wallet: id, enabled: true, destination: { kind: 'wallet' } }));
    expect(injected.length).toBe(0);
    expect(st.status.paused).toBe('locked');
    expect(st.status.waiting_coins).toBe('5.000000');

    await expectAsync(firstValueFrom(service.unlock(id, 'bad', 60))).toBeRejected();
    st = await firstValueFrom(service.unlock(id, 'pw', 60));
    expect(st.status.unlocked).toBeTrue();
    expect(st.status.unlock_expires).toBeDefined();
    expect(injected.length).toBe(1);
    expect(w.addresses.every(a => !a.secret_key)).toBeTrue();
    expect(localStorage.getItem(browserSweepRulesKey)).not.toContain('pw');

    st = await firstValueFrom(service.lock(id));
    expect(st.status.unlocked).toBeFalse();
  });

  it('follows a level 2 mapping with level 1 off and leaves the rest alone', async () => {
    outputs.push({ hash: hash(2), address: chain[1].address, coins: '3.000000', calculated_hours: 9 });
    sentBy = [chain[0].address, chain[1].address];
    const w = wallet({ addresses: [{ address: chain[0].address }, { address: chain[1].address }, { address: chain[2].address }] });
    const id = browserWalletId(w);
    const fixed = chain[3].address;
    const st = await firstValueFrom(service.setRule({
      wallet: id, enabled: false, destination: { kind: 'wallet' },
      mappings: [{ from: { address: chain[1].address }, to: { kind: 'address', address: fixed } }],
    }));
    expect(injected.length).toBe(1);
    expect(st.status.enabled).toBeFalse();
    expect(st.status.active).toBeTrue();
    expect(st.status.uncovered_addresses).toBe(1);
    expect(st.status.destinations).toEqual([fixed]);
  });

  it('refuses what a browser wallet cannot do', async () => {
    const w = wallet();
    const id = browserWalletId(w);
    await expectAsync(firstValueFrom(service.setRule({ wallet: id, enabled: true, destination: { kind: 'account' } })))
      .toBeRejectedWithError('wallet.sweep.err-browser-account');
    await expectAsync(firstValueFrom(service.setRule({ wallet: id, enabled: true, destination: { kind: 'address', address: chain[0].address } })))
      .toBeRejectedWithError(/the destination is exposed/);
    await expectAsync(firstValueFrom(service.setRule({
      wallet: id, enabled: false, destination: { kind: 'wallet' },
      mappings: [{ from: { address: chain[1].address }, to: { kind: 'wallet' } }],
    }))).toBeRejectedWithError(/never sent/);
    expect(localStorage.getItem(browserSweepRulesKey)).toBeNull();
    expect(injected.length).toBe(0);
  });
});
