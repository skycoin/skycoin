import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';

import { NumberOfAddressesComponent } from './number-of-addresses.component';
import { MockMatDialogRef, MockTranslatePipe } from '../../../../utils/test-mocks';

describe('NumberOfAddressesComponent', () => {
  let component: NumberOfAddressesComponent;
  let fixture: ComponentFixture<NumberOfAddressesComponent>;
  let dialogRef: MatDialogRef<NumberOfAddressesComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ NumberOfAddressesComponent, MockTranslatePipe ],
      imports: [ ReactiveFormsModule ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: [
        { provide: MatDialogRef, useClass: MockMatDialogRef },
      ]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NumberOfAddressesComponent);
    component = fixture.componentInstance;
    dialogRef = TestBed.inject(MatDialogRef);
    spyOn(dialogRef, 'close');
    fixture.detectChanges();
  });

  it('starts at one address', () => {
    expect(component.form.value.quantity).toBe(1);
    expect(component.form.valid).toBe(true);
  });

  it('accepts whole numbers from 1 to 100 only', () => {
    const quantity = component.form.get('quantity')!;
    [0, 101, 1.5, -3, '', null].forEach(value => {
      quantity.setValue(value);
      expect(component.form.valid).withContext(String(value)).toBe(false);
    });
    [1, 42, 100, '7'].forEach(value => {
      quantity.setValue(value);
      expect(component.form.valid).withContext(String(value)).toBe(true);
    });
  });

  it('closes with the number entered', () => {
    component.form.get('quantity')!.setValue('12');
    component.continue();

    expect(dialogRef.close).toHaveBeenCalledWith(12);
  });

  it('does not close on an invalid number', () => {
    component.form.get('quantity')!.setValue(500);
    component.continue();

    expect(dialogRef.close).not.toHaveBeenCalled();
  });

  it('closes with nothing when cancelled', () => {
    component.closePopup();

    expect(dialogRef.close).toHaveBeenCalledWith();
  });
});
