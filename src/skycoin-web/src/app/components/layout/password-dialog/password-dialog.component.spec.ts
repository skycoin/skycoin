import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ReactiveFormsModule, UntypedFormBuilder } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

import { PasswordDialogComponent, PasswordDialogData } from './password-dialog.component';
import { MockTranslatePipe } from '../../../utils/test-mocks';

describe('PasswordDialogComponent', () => {
  let component: PasswordDialogComponent;
  let fixture: ComponentFixture<PasswordDialogComponent>;
  let dialogRef: jasmine.SpyObj<MatDialogRef<PasswordDialogComponent>>;

  function create(data: PasswordDialogData) {
    dialogRef = jasmine.createSpyObj('MatDialogRef', ['close']);
    TestBed.configureTestingModule({
      declarations: [PasswordDialogComponent, MockTranslatePipe],
      imports: [ReactiveFormsModule],
      schemas: [NO_ERRORS_SCHEMA],
      providers: [
        UntypedFormBuilder,
        { provide: MatDialogRef, useValue: dialogRef },
        { provide: MAT_DIALOG_DATA, useValue: data },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(PasswordDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  }

  it('closes with the password', waitForAsync(() => {
    create({});
    component.form.get('password')!.setValue('secret');
    component.submit();
    expect(dialogRef.close).toHaveBeenCalledWith('secret');
  }));

  it('masks the field', waitForAsync(() => {
    create({});
    const input = fixture.nativeElement.querySelector('#password') as HTMLInputElement;
    expect(input.type).toBe('password');
  }));

  it('does not close with an empty password unless allowed', waitForAsync(() => {
    create({});
    component.submit();
    expect(dialogRef.close).not.toHaveBeenCalled();
  }));

  it('closes with an empty password when allowed', waitForAsync(() => {
    create({ allowEmpty: true });
    component.submit();
    expect(dialogRef.close).toHaveBeenCalledWith('');
  }));

  it('needs the confirmation to match', waitForAsync(() => {
    create({ confirm: true });
    component.form.get('password')!.setValue('secret');
    component.form.get('confirm')!.setValue('other');
    component.form.get('confirm')!.markAsDirty();
    component.submit();
    expect(dialogRef.close).not.toHaveBeenCalled();
    expect(component.mismatch).toBeTrue();

    component.form.get('confirm')!.setValue('secret');
    component.submit();
    expect(dialogRef.close).toHaveBeenCalledWith('secret');
  }));

  it('closes with undefined when cancelled', waitForAsync(() => {
    create({});
    component.cancel();
    expect(dialogRef.close).toHaveBeenCalledWith();
  }));
});
