import { ComponentFixture, TestBed, fakeAsync, tick, discardPeriodicTasks } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA, Pipe, PipeTransform } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { of, throwError } from 'rxjs';

import { AutoSweepComponent, sweepPollMs } from './auto-sweep.component';
import { MockWalletService } from '../../../../utils/test-mocks';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { TranslateService } from '@ngx-translate/core';

import { SweepService, SweepState, sweepActive } from '../../../../services/wallet/sweep.service';
import { BrowserSweepService } from '../../../../services/wallet/browser-sweep.service';
import { Wallet } from '../../../../app.datatypes';

/** Shows the key, so a test can tell which text is on screen. */
@Pipe({ name: 'translate', standalone: false })
class KeyPipe implements PipeTransform {
  transform(key: string): string {
    return key;
  }
}

function state(over: Partial<SweepState['status']> = {}, rule: SweepState['rule'] = null): SweepState {
  return {
    rule: rule,
    status: {
      wallet: 'reward.wlt', enabled: !!rule && rule.enabled, active: sweepActive(rule), unlocked: false, head_seq: 10,
      exposed_addresses: 21, uncovered_addresses: 0,
      waiting_coins: '0.000000', waiting_hours: 0, contested: 0, confirmed_sweeps: 0, ...over,
    },
  };
}

describe('AutoSweepComponent', () => {
  let fixture: ComponentFixture<AutoSweepComponent>;
  let component: AutoSweepComponent;
  let sweep: jasmine.SpyObj<SweepService>;
  let browser: jasmine.SpyObj<BrowserSweepService>;
  let walletService: MockWalletService;
  let wallet: Wallet;

  const text = () => (fixture.nativeElement as HTMLElement).textContent || '';

  function create(first: SweepState | Error) {
    browser.get.and.returnValue(first instanceof Error ? throwError(() => first) : of(first));
    sweep.get.and.returnValue(first instanceof Error ? throwError(() => first) : of(first));
    fixture = TestBed.createComponent(AutoSweepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }

  beforeEach(() => {
    sweep = jasmine.createSpyObj('SweepService', ['get', 'setRule', 'unlock', 'lock']);
    browser = jasmine.createSpyObj('BrowserSweepService', ['get', 'setRule', 'unlock', 'lock']);
    wallet = {
      label: 'Reward', filename: 'reward.wlt', walletType: 'bip44', encrypted: true, addresses: [], coinId: 1,
      accounts: [
        { index: 0, name: 'default', externalAddresses: [], changeAddresses: [] },
        { index: 1, name: 'safe', externalAddresses: [], changeAddresses: [] },
      ],
    };
    TestBed.configureTestingModule({
      declarations: [AutoSweepComponent, KeyPipe],
      imports: [FormsModule],
      schemas: [NO_ERRORS_SCHEMA],
      providers: [
        { provide: SweepService, useValue: sweep },
        { provide: BrowserSweepService, useValue: browser },
        { provide: TranslateService, useValue: { instant: (k: string, p?: any) => k + (p ? ' ' + JSON.stringify(p) : '') } },
        { provide: WalletService, useClass: MockWalletService },
        { provide: MatDialogRef, useValue: { close: () => {} } },
        { provide: MAT_DIALOG_DATA, useFactory: () => wallet },
      ],
    });
    walletService = TestBed.inject(WalletService) as any;
    walletService.wallets.next([
      wallet,
      { label: 'Cold', filename: 'cold.wlt', walletType: 'deterministic', addresses: [] },
      { label: 'Temp', filename: 'temp.wlt', temporary: true, addresses: [] },
      { label: 'HW', filename: '', isHardware: true, addresses: [] },
    ]);
  });

  it('is off by default and offers the next account, each account and the other wallet', () => {
    create(state());
    expect(sweep.get).toHaveBeenCalledWith('reward.wlt');
    expect(component.enabled).toBeFalse();
    expect(component.destination).toBe('next');
    expect(component.choices.map(c => c.key)).toEqual(['next', 'account:0', 'account:1', 'wallet:cold.wlt', 'address']);
    expect(text()).toContain('wallet.sweep.burn-warning');
    expect(text()).toContain('wallet.sweep.off');
    expect(fixture.nativeElement.querySelector('#sweep-password')).not.toBeNull();
  });

  it('fills the form from the node rule', () => {
    create(state({ unlocked: true }, {
      wallet: 'reward.wlt', enabled: true, shared_accounts: [0], gap: 50, addresses: ['a1'],
      destination: { kind: 'account', account: 1 },
    }));
    expect(component.enabled).toBeTrue();
    expect(component.shared.has(0)).toBeTrue();
    expect(component.destination).toBe('account:1');
    expect(component.addresses).toBe('a1');
    expect(component.active).toBeTrue();
    expect(text()).toContain('wallet.sweep.active');
    expect(fixture.nativeElement.querySelector('#sweep-password')).toBeNull();
    expect(fixture.nativeElement.querySelector('#sweep-lock')).not.toBeNull();
  });

  it('asks to confirm the burn before turning it on', () => {
    create(state());
    component.enabled = true;
    component.save();
    expect(sweep.setRule).not.toHaveBeenCalled();
    expect(component.error).toBe('wallet.sweep.err-burn');
  });

  it('sets the rule, keeping the fields it does not show', () => {
    create(state({}, { wallet: 'reward.wlt', enabled: false, gap: 50, destination: { kind: 'account' } }));
    sweep.setRule.and.returnValue(of(state({}, { wallet: 'reward.wlt', enabled: true, shared_accounts: [0], gap: 50,
      destination: { kind: 'wallet', wallet: 'cold.wlt' } })));
    component.enabled = true;
    component.burnOk = true;
    component.toggleShared(1, true);
    component.toggleShared(0, true);
    component.toggleShared(1, false);
    component.destination = 'wallet:cold.wlt';
    component.addresses = ' a1\n a2 ';
    component.save();
    expect(sweep.setRule).toHaveBeenCalledWith({
      wallet: 'reward.wlt', enabled: true, gap: 50, shared_accounts: [0], addresses: ['a1', 'a2'], mappings: [],
      destination: { kind: 'wallet', wallet: 'cold.wlt' },
    });
    expect(component.done).toBe('wallet.sweep.done-saved');
    expect(component.destination).toBe('wallet:cold.wlt');
  });

  it('builds each destination', () => {
    create(state());
    component.destination = 'account:1';
    expect(component.buildRule().destination).toEqual({ kind: 'account', account: 1 });
    component.destination = 'next';
    expect(component.buildRule().destination).toEqual({ kind: 'account' });
    component.destination = 'address';
    component.address = ' 2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv ';
    expect(component.buildRule().destination).toEqual({ kind: 'address', address: '2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv' });
  });

  it('shows the node refusing a rule', () => {
    create(state());
    sweep.setRule.and.returnValue(throwError(() => new Error('destination account 0 is shared')));
    component.save();
    fixture.detectChanges();
    expect(component.busy).toBeFalse();
    expect(text()).toContain('destination account 0 is shared');
  });

  it('refuses an empty fixed address before asking the node', () => {
    create(state());
    component.destination = 'address';
    component.save();
    expect(sweep.setRule).not.toHaveBeenCalled();
    expect(component.error).toBe('wallet.sweep.err-destination');
  });

  it('unlocks for a time and locks', () => {
    create(state());
    sweep.unlock.and.returnValue(of(state({ unlocked: true })));
    component.unlock();
    expect(sweep.unlock).not.toHaveBeenCalled();
    expect(component.error).toBe('wallet.sweep.err-password');

    component.password = 'pw';
    component.ttl = 0;
    component.unlock();
    expect(sweep.unlock).toHaveBeenCalledWith('reward.wlt', 'pw', 0);
    expect(component.password).toBe('');
    fixture.detectChanges();
    expect(text()).toContain('wallet.sweep.until-restart');

    sweep.lock.and.returnValue(of(state()));
    component.lock();
    expect(sweep.lock).toHaveBeenCalledWith('reward.wlt');
    expect(component.done).toBe('wallet.sweep.done-locked');
  });

  it('polls the status while open and keeps the form', fakeAsync(() => {
    create(state());
    component.destination = 'address';
    sweep.get.and.returnValue(of(state({ exposed_addresses: 30, last_txids: ['txLast'], last_error: 'inject failed' })));
    tick(sweepPollMs);
    fixture.detectChanges();
    expect(sweep.get).toHaveBeenCalledTimes(2);
    expect(component.state!.status.exposed_addresses).toBe(30);
    expect(component.destination).toBe('address');
    expect(text()).toContain('txLast');
    expect(text()).toContain('inject failed');

    fixture.destroy();
    tick(sweepPollMs * 2);
    expect(sweep.get).toHaveBeenCalledTimes(2);
    discardPeriodicTasks();
  }));

  it('shows the error when the status cannot be read', () => {
    create(new Error('sweeping needs the wallet API'));
    expect(text()).toContain('sweeping needs the wallet API');
    expect(fixture.nativeElement.querySelector('#sweep-save')).toBeNull();
  });

  it('offers this wallet for a wallet that is not bip44', () => {
    wallet = { label: 'Cold', filename: 'cold.wlt', walletType: 'deterministic', addresses: [] };
    create(state());
    expect(component.destination).toBe('self');
    expect(component.choices.map(c => c.key)).toEqual(['self', 'wallet:reward.wlt', 'address']);
    expect(component.buildRule().shared_accounts).toEqual([]);
    expect(fixture.nativeElement.querySelector('#sweep-password')).toBeNull();
  });

  it('adds a level 2 mapping with level 1 off, asking for the burn first', () => {
    create(state());
    component.mapFrom = '1';
    component.mapTo = 'wallet:cold.wlt';
    component.addMapping();
    expect(sweep.setRule).not.toHaveBeenCalled();
    expect(component.error).toBe('wallet.sweep.err-burn');
    expect(component.maps).toEqual([]);

    const mapped = { wallet: 'reward.wlt', enabled: false, destination: { kind: 'account' as const },
      mappings: [{ from: { account: 1 }, to: { kind: 'wallet' as const, wallet: 'cold.wlt' } }] };
    sweep.setRule.and.returnValue(of(state({}, mapped)));
    component.burnOk = true;
    component.addMapping();
    expect(sweep.setRule.calls.mostRecent().args[0].mappings).toEqual(mapped.mappings);
    expect(sweep.setRule.calls.mostRecent().args[0].enabled).toBeFalse();
    expect(component.mapFrom).toBe('');
    fixture.detectChanges();
    expect(component.active).toBeTrue();
    expect(fixture.nativeElement.querySelector('#sweep-maps').textContent).toContain('Cold');
    expect(fixture.nativeElement.querySelector('#sweep-advanced').open).toBeTrue();

    sweep.setRule.and.returnValue(of(state({}, { ...mapped, mappings: [] })));
    component.removeMapping(0);
    expect(sweep.setRule.calls.mostRecent().args[0].mappings).toEqual([]);
  });

  it('refuses an incomplete mapping and keeps the list when the node refuses one', () => {
    create(state({}, { wallet: 'reward.wlt', enabled: true, destination: { kind: 'account' } }));
    component.mapFrom = 'a1';
    component.mapTo = 'address';
    component.addMapping();
    expect(component.error).toBe('wallet.sweep.err-mapping');
    component.mapToAddress = 'bad';
    sweep.setRule.and.returnValue(throwError(() => new Error('mapping 1: invalid destination address')));
    component.addMapping();
    expect(sweep.setRule.calls.mostRecent().args[0].mappings).toEqual([{ from: { address: 'a1' }, to: { kind: 'address', address: 'bad' } }]);
    expect(component.maps).toEqual([]);
    expect(component.error).toContain('invalid destination');
  });

  it('refuses a mapping whose destination is not chosen', () => {
    create(state());
    component.mapFrom = '1';
    component.mapTo = '';
    component.addMapping();
    expect(sweep.setRule).not.toHaveBeenCalled();
    expect(component.error).toBe('wallet.sweep.err-mapping');
  });

  it('turns off the has sent inclusion', () => {
    create(state());
    expect(component.includeSpent).toBeTrue();
    expect(component.buildRule().include_spent).toBeUndefined();
    component.includeSpent = false;
    expect(component.buildRule().include_spent).toBeFalse();
  });

  describe('a wallet kept in the browser', () => {
    beforeEach(() => {
      wallet = {
        label: 'Browser', walletType: 'deterministic', encryptedSeed: 'enc', coinId: 1,
        addresses: [{ address: 'b0' }, { address: 'b1' }],
      };
      walletService.wallets.next([wallet, { label: 'Other', walletType: 'bip44', addresses: [],
        accounts: [{ index: 0, name: '', externalAddresses: [{ address: 'o0' }], changeAddresses: [] }] }]);
    });

    it('is swept by this page, which says so', () => {
      create(state({ wallet: 'b0' }));
      expect(browser.get).toHaveBeenCalledWith('b0');
      expect(sweep.get).not.toHaveBeenCalled();
      expect(component.kind).toBe('browser');
      expect(component.choices.map(c => c.key)).toEqual(['self', 'wallet:o0', 'address']);
      expect(component.unlockTimes.map(t => t.label)).toContain('wallet.sweep.ttl-page');
      expect(fixture.nativeElement.querySelector('#sweep-browser-note')).not.toBeNull();
      expect(fixture.nativeElement.querySelector('#sweep-password')).not.toBeNull();
    });

    it('unlocks in the page and shows why it waits', () => {
      create(state({ wallet: 'b0', active: true, paused: 'locked', waiting_coins: '9.000000' },
        { wallet: 'b0', enabled: true, destination: { kind: 'wallet' } }));
      fixture.detectChanges();
      expect(fixture.nativeElement.querySelector('#sweep-paused').textContent).toContain('wallet.sweep.paused-locked');
      expect(fixture.nativeElement.querySelector('#sweep-paused').textContent).toContain('9.000000');
      browser.unlock.and.returnValue(of(state({ wallet: 'b0', active: true, unlocked: true })));
      component.password = 'pw';
      component.ttl = 0;
      component.unlock();
      expect(browser.unlock).toHaveBeenCalledWith('b0', 'pw', 0);
      expect(component.done).toBe('wallet.sweep.browser-done-unlocked');
      fixture.detectChanges();
      expect(fixture.nativeElement.querySelector('#sweep-unlocked').textContent).toContain('wallet.sweep.ttl-page');
    });

    it('takes an address source, never an account', () => {
      create(state({ wallet: 'b0' }));
      browser.setRule.and.returnValue(of(state({ wallet: 'b0' })));
      component.burnOk = true;
      component.mapFrom = '5';
      component.mapTo = 'self';
      component.addMapping();
      expect(browser.setRule.calls.mostRecent().args[0].mappings).toEqual([{ from: { address: '5' }, to: { kind: 'wallet' } }]);
    });
  });

  it('says a hardware wallet is never swept', () => {
    wallet = { label: 'HW', isHardware: true, addresses: [{ address: 'h0' }] };
    create(state());
    expect(component.kind).toBe('hardware');
    expect(sweep.get).not.toHaveBeenCalled();
    expect(browser.get).not.toHaveBeenCalled();
    expect(fixture.nativeElement.querySelector('#sweep-hardware')).not.toBeNull();
    expect(fixture.nativeElement.querySelector('#sweep-save')).toBeNull();
  });
});
