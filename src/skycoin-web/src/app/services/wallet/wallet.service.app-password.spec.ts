import { TestBed } from '@angular/core/testing';
import { of, firstValueFrom } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';

import { WalletService } from './wallet.service';
import { CipherProvider } from '../cipher.provider';
import { CoinService } from '../coin.service';
import { EncryptionService } from '../encryption.service';
import { ApiService } from '../api.service';
import { GlobalsService } from '../globals.service';
import { MockCoinService, MockGlobalsService } from '../../utils/test-mocks';

describe('WalletService app password', () => {
  let store: Record<string, string>;

  const stored = [{ coinId: 1, label: 'savings', addresses: [{ address: 'addr1' }], seed: 'the seed' }];

  function newService(): WalletService {
    TestBed.resetTestingModule();
    TestBed.configureTestingModule({
      providers: [
        WalletService,
        EncryptionService,
        { provide: CipherProvider, useValue: jasmine.createSpyObj('CipherProvider', ['generateAddress']) },
        { provide: TranslateService, useValue: jasmine.createSpyObj('TranslateService', ['instant']) },
        { provide: ApiService, useValue: jasmine.createSpyObj('ApiService', { get: of([]) }) },
        { provide: CoinService, useClass: MockCoinService },
        { provide: GlobalsService, useClass: MockGlobalsService },
      ],
    });

    return TestBed.inject(WalletService);
  }

  beforeEach(() => {
    store = { wallets: JSON.stringify(stored) };
    spyOn(localStorage, 'getItem').and.callFake((key: string) => store[key] ?? null);
    spyOn(localStorage, 'setItem').and.callFake((key: string, value: string) => { store[key] = value; });
    spyOn(localStorage, 'removeItem').and.callFake((key: string) => { delete store[key]; });
  });

  it('stores the wallets encrypted once a password is set', async () => {
    const service = newService();
    expect(service.hasAppPassword).toBeFalse();

    await service.setAppPassword('hunter2');

    expect(service.hasAppPassword).toBeTrue();
    expect(store['wallets']).toBeUndefined();
    expect(store['lockedWallets']).toBeTruthy();
    expect(store['lockedWallets']).not.toContain('savings');
    expect(store['lockedWallets']).not.toContain('the seed');
  });

  it('stays locked until the right password is entered', async () => {
    await newService().setAppPassword('hunter2');

    const service = newService();
    expect(service.locked.value).toBeTrue();
    expect(service.wallets.value).toBeNull();

    await expectAsync(firstValueFrom(service.unlockApp('wrong'))).toBeRejected();
    expect(service.locked.value).toBeTrue();

    await firstValueFrom(service.unlockApp('hunter2'));
    expect(service.locked.value).toBeFalse();
    expect(service.wallets.value.map(w => w.label)).toEqual(['savings']);
  });

  it('keeps saving encrypted after unlocking', async () => {
    await newService().setAppPassword('hunter2');
    const service = newService();
    await firstValueFrom(service.unlockApp('hunter2'));

    service.wallets.value[0].label = 'renamed';
    service.saveWallets();
    await (service as any).browser.saves;

    expect(store['wallets']).toBeUndefined();
    const reopened = newService();
    await firstValueFrom(reopened.unlockApp('hunter2'));
    expect(reopened.wallets.value.map(w => w.label)).toEqual(['renamed']);
  });

  it('removes the password only with the current one', async () => {
    const service = newService();
    await service.setAppPassword('hunter2');

    await expectAsync(service.removeAppPassword('wrong')).toBeRejected();
    expect(store['lockedWallets']).toBeTruthy();

    await service.removeAppPassword('hunter2');
    expect(service.hasAppPassword).toBeFalse();
    expect(store['lockedWallets']).toBeUndefined();
    expect(JSON.parse(store['wallets'])[0].label).toBe('savings');
  });

  it('erases the locked wallets when the password is forgotten', async () => {
    await newService().setAppPassword('hunter2');
    const service = newService();

    service.eraseLockedWallets();

    expect(service.locked.value).toBeFalse();
    expect(service.wallets.value).toEqual([]);
    expect(store['lockedWallets']).toBeUndefined();
    expect(store['wallets']).toBeUndefined();
  });
});
