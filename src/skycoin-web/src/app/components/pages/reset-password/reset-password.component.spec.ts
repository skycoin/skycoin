import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ReactiveFormsModule, UntypedFormBuilder } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { BehaviorSubject, of, throwError } from 'rxjs';

import { ResetPasswordComponent } from './reset-password.component';
import { Wallet } from '../../../app.datatypes';
import { WalletService } from '../../../services/wallet/wallet.service';
import { MsgBarService } from '../../../services/msg-bar.service';
import { MockTranslatePipe } from '../../../utils/test-mocks';

describe('ResetPasswordComponent', () => {
  let fixture: ComponentFixture<ResetPasswordComponent>;
  let component: ResetPasswordComponent;
  let walletService: jasmine.SpyObj<WalletService>;
  let router: jasmine.SpyObj<Router>;
  let msgBar: jasmine.SpyObj<MsgBarService>;
  const encrypted: Wallet = { label: 'locked', addresses: [{ address: 'addr1' }], encryptedSeed: 'x' };
  const plain: Wallet = { label: 'plain', addresses: [{ address: 'addr2' }], seed: 's' };

  function create(id: string) {
    walletService = jasmine.createSpyObj('WalletService', ['resetPassword', 'isEncrypted'], {
      currentWallets: new BehaviorSubject<Wallet[]>([encrypted, plain]),
    });
    walletService.isEncrypted.and.callFake((w: Wallet) => !!w.encryptedSeed);
    router = jasmine.createSpyObj('Router', ['navigate']);
    msgBar = jasmine.createSpyObj('MsgBarService', ['hide', 'showDone', 'showError']);
    TestBed.configureTestingModule({
      declarations: [ResetPasswordComponent, MockTranslatePipe],
      imports: [ReactiveFormsModule],
      schemas: [NO_ERRORS_SCHEMA],
      providers: [
        UntypedFormBuilder,
        { provide: ActivatedRoute, useValue: { params: of({ id }) } },
        { provide: Router, useValue: router },
        { provide: WalletService, useValue: walletService },
        { provide: MsgBarService, useValue: msgBar },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(ResetPasswordComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    component.resetButton = jasmine.createSpyObj('ButtonComponent', ['setLoading', 'setSuccess', 'setDisabled', 'resetState']);
  }

  it('finds the encrypted wallet named in the route', () => {
    create('addr1');
    expect(component.wallet).toBe(encrypted);
    expect(router.navigate).not.toHaveBeenCalled();
  });

  it('leaves for a wallet that is missing or not encrypted', () => {
    create('addr2');
    expect(component.wallet).toBeNull();
    expect(router.navigate).toHaveBeenCalledWith(['/wallets']);
  });

  it('needs a seed and matching passwords', () => {
    create('addr1');
    expect(component.form.valid).toBeFalse();
    component.form.patchValue({ seed: 'seed', password: 'a', confirm: 'b' });
    expect(component.form.valid).toBeFalse();
    component.form.patchValue({ confirm: 'a' });
    expect(component.form.valid).toBeTrue();
  });

  it('resets with the seed and the new password, then goes back to the wallets', fakeAsync(() => {
    create('addr1');
    walletService.resetPassword.and.returnValue(of(undefined));
    component.form.patchValue({ seed: ' seed words ', password: 'new', confirm: 'new' });
    component.reset();
    expect(walletService.resetPassword).toHaveBeenCalledWith(encrypted, 'seed words', '', 'new');
    expect(msgBar.showDone).toHaveBeenCalledWith('reset.done');
    tick(2000);
    expect(router.navigate).toHaveBeenCalledWith(['/wallets']);
  }));

  it('removes the password when the new one is empty', () => {
    create('addr1');
    walletService.resetPassword.and.returnValue(of(undefined));
    component.form.patchValue({ seed: 'seed' });
    component.reset();
    expect(walletService.resetPassword).toHaveBeenCalledWith(encrypted, 'seed', '', null);
  });

  it('shows the error of a wrong seed and stays', () => {
    create('addr1');
    walletService.resetPassword.and.returnValue(throwError(() => new Error('service.wallet.wrong-seed')));
    component.form.patchValue({ seed: 'wrong' });
    component.reset();
    expect(msgBar.showError).toHaveBeenCalledWith('service.wallet.wrong-seed');
    expect(component.busy).toBeFalse();
    expect(router.navigate).not.toHaveBeenCalled();
  });
});
