import { ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ReactiveFormsModule, UntypedFormBuilder } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';

import { HwSeedWordDialogComponent } from './hw-seed-word-dialog.component';
import { Bip39WordListService } from '../../../../services/bip39-word-list.service';
import { MockTranslatePipe } from '../../../../utils/test-mocks';

describe('HwSeedWordDialogComponent', () => {
  let component: HwSeedWordDialogComponent;
  let fixture: ComponentFixture<HwSeedWordDialogComponent>;
  let dialogRef: jasmine.SpyObj<MatDialogRef<HwSeedWordDialogComponent>>;

  beforeEach(() => {
    dialogRef = jasmine.createSpyObj('MatDialogRef', ['close']);
    TestBed.configureTestingModule({
      declarations: [HwSeedWordDialogComponent, MockTranslatePipe],
      imports: [ReactiveFormsModule],
      schemas: [NO_ERRORS_SCHEMA],
      providers: [
        UntypedFormBuilder,
        { provide: MatDialogRef, useValue: dialogRef },
        { provide: Bip39WordListService, useValue: { validateWord: (w: string) => w === 'abandon' } },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(HwSeedWordDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('closes with the word, trimmed and in lower case', () => {
    component.form.get('word')!.setValue('  Abandon ');
    component.submit();
    expect(dialogRef.close).toHaveBeenCalledWith('abandon');
  });

  it('refuses a word outside the BIP39 list', () => {
    component.form.get('word')!.setValue('notaword');
    component.submit();
    expect(dialogRef.close).not.toHaveBeenCalled();
    expect(component.error).toBe('hardware-wallet.seed-word.invalid-word');
  });

  it('closes with undefined when cancelled', () => {
    component.cancel();
    expect(dialogRef.close).toHaveBeenCalledWith();
  });
});
