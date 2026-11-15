import { ComponentFixture, TestBed, fakeAsync, tick, discardPeriodicTasks } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA, Pipe, PipeTransform } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { of, throwError } from 'rxjs';

import { AutoSweepComponent, sweepPollMs } from './auto-sweep.component';
import { MockWalletService } from '../../../../utils/test-mocks';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { SweepService, SweepState } from '../../../../services/wallet/sweep.service';
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
      wallet: 'reward.wlt', enabled: !!rule && rule.enabled, unlocked: false, head_seq: 10, exposed_addresses: 21,
      waiting_coins: '0.000000', waiting_hours: 0, contested: 0, confirmed_sweeps: 0, ...over,
    },
  };
}

describe('AutoSweepComponent', () => {
  let fixture: ComponentFixture<AutoSweepComponent>;
  let component: AutoSweepComponent;
  let sweep: jasmine.SpyObj<SweepService>;
  let walletService: MockWalletService;
  let wallet: Wallet;

  const text = () => (fixture.nativeElement as HTMLElement).textContent || '';

  function create(first: SweepState | Error) {
    sweep.get.and.returnValue(first instanceof Error ? throwError(() => first) : of(first));
    fixture = TestBed.createComponent(AutoSweepComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }

  beforeEach(() => {
    sweep = jasmine.createSpyObj('SweepService', ['get', 'setRule', 'unlock', 'lock']);
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
      wallet: 'reward.wlt', enabled: true, gap: 50, shared_accounts: [0], addresses: ['a1', 'a2'],
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
});
