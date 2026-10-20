import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ReactiveFormsModule, UntypedFormBuilder } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { of, throwError } from 'rxjs';

import { OfflineTxComponent, OfflineTxData, OfflineTxMode } from './offline-tx.component';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { SpendingService } from '../../../../services/wallet/spending.service';
import { BalanceService } from '../../../../services/wallet/balance.service';
import { ClipboardService } from '../../../../services/clipboard.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { MockTranslatePipe } from '../../../../utils/test-mocks';

describe('OfflineTxComponent', () => {
  let component: OfflineTxComponent;
  let fixture: ComponentFixture<OfflineTxComponent>;
  let dialogRef: jasmine.SpyObj<MatDialogRef<OfflineTxComponent>>;
  let spending: jasmine.SpyObj<SpendingService>;
  let clipboard: jasmine.SpyObj<ClipboardService>;

  const wallets: any[] = [
    { label: 'software', addresses: [], seed: 'seed' },
    { label: 'hardware', addresses: [], isHardware: true },
  ];

  function create(data: OfflineTxData) {
    dialogRef = jasmine.createSpyObj('MatDialogRef', ['close']);
    spending = jasmine.createSpyObj('SpendingService', ['injectTransaction', 'signTransaction']);
    clipboard = jasmine.createSpyObj('ClipboardService', ['copy']);
    clipboard.copy.and.returnValue(Promise.resolve(''));
    TestBed.configureTestingModule({
      declarations: [OfflineTxComponent, MockTranslatePipe],
      imports: [ReactiveFormsModule],
      schemas: [NO_ERRORS_SCHEMA],
      providers: [
        UntypedFormBuilder,
        { provide: MatDialogRef, useValue: dialogRef },
        { provide: MAT_DIALOG_DATA, useValue: data },
        { provide: WalletService, useValue: { currentWallets: of(wallets), isServerWallet: () => false } },
        { provide: SpendingService, useValue: spending },
        { provide: BalanceService, useValue: { startGettingBalances: () => {} } },
        { provide: ClipboardService, useValue: clipboard },
        { provide: CustomMatDialogService, useValue: {} },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(OfflineTxComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }

  it('offers only software wallets for signing', () => {
    create({ mode: OfflineTxMode.Sign });
    expect(component.wallets.map(w => w.label)).toEqual(['software']);
    expect(component.form.value.wallet.label).toBe('software');
  });

  it('broadcasts the pasted transaction and closes', () => {
    create({ mode: OfflineTxMode.Broadcast });
    spending.injectTransaction.and.returnValue(of('txid'));
    component.form.get('rawTx')!.setValue('  signed \n');
    component.broadcast();
    expect(spending.injectTransaction).toHaveBeenCalledWith('signed');
    expect(dialogRef.close).toHaveBeenCalledWith(true);
  });

  it('shows why a broadcast failed and stays open', () => {
    create({ mode: OfflineTxMode.Broadcast });
    spending.injectTransaction.and.returnValue(throwError(() => new Error('Invalid transaction')));
    component.form.get('rawTx')!.setValue('bad');
    component.broadcast();
    expect(component.error).toBe('Invalid transaction');
    expect(dialogRef.close).not.toHaveBeenCalled();
  });

  it('signs with an unlocked browser wallet', () => {
    create({ mode: OfflineTxMode.Sign });
    spending.signTransaction.and.returnValue(throwError(() => new Error('offline-tx.not-in-wallet')));
    component.form.get('rawTx')!.setValue('unsigned');
    component.sign();
    expect(spending.signTransaction).toHaveBeenCalledWith(wallets[0], undefined, 'unsigned');
    expect(component.error).toBe('offline-tx.not-in-wallet');
  });

  it('shows a transaction read only and copies it', fakeAsync(() => {
    create({ mode: OfflineTxMode.Show, rawTx: 'abc', unsigned: true });
    expect(component.title).toBe('offline-tx.unsigned-title');
    const textarea = fixture.nativeElement.querySelector('textarea') as HTMLTextAreaElement;
    expect(textarea.readOnly).toBeTrue();
    component.copy();
    tick();
    expect(clipboard.copy).toHaveBeenCalledWith('abc');
    expect(component.copied).toBeTrue();
  }));
});
