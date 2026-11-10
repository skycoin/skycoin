import { of } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { ServerWalletBackend } from './server-wallet-backend';
import { Wallet } from '../../../app.datatypes';

describe('ServerWalletBackend', () => {
  it('keeps the address balances when a scan reloads a bip44 wallet', () => {
    const response = {
      meta: { type: 'bip44' },
      entries: [{ address: 'a1' }, { address: 'a2' }, { address: 'a3' }],
      accounts: [{
        name: 'default', index: 0,
        external_entries: [{ address: 'a1' }, { address: 'a2' }, { address: 'a3' }],
        change_entries: [],
      }],
    };
    const api = { post: () => of({}), get: () => of(response) } as any;
    const backend = new ServerWalletBackend(api, () => 0);

    const coins = (n: number) => ({ address: 'a' + n, balance: new BigNumber(n), hours: new BigNumber(n * 10) }) as any;
    const wallet = {
      filename: 'w.wlt',
      addresses: [coins(1), coins(2)],
      accounts: [{ name: 'default', index: 0, externalAddresses: [coins(1), coins(2)], changeAddresses: [] }],
    } as unknown as Wallet;

    let added = -1;
    backend.scanAddresses(wallet, 20).subscribe(n => added = n);

    expect(added).toBe(1);
    const ext = wallet.accounts![0].externalAddresses;
    expect(ext.map(a => a.balance ? a.balance.toNumber() : null)).toEqual([1, 2, null]);
    expect(ext[1].hours!.toNumber()).toBe(20);
    expect(wallet.addresses.map(a => a.balance ? a.balance.toNumber() : null)).toEqual([1, 2, null]);
  });

  it('renames a wallet on the wallet service', () => {
    const api = jasmine.createSpyObj('ApiService', ['post']);
    api.post.and.returnValue(of({}));
    const backend = new ServerWalletBackend(api, () => 0);

    let done = false;
    backend.rename({ filename: 'w.wlt' } as Wallet, 'rewards').subscribe(() => done = true);

    expect(done).toBeTrue();
    expect(api.post).toHaveBeenCalledWith('wallet/update', { id: 'w.wlt', label: 'rewards' });
  });
});
