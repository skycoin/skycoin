import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwConfirmAddressDialogComponent } from './hw-confirm-address-dialog.component';

describe('HwConfirmAddressDialogComponent', () => {
  let component: HwConfirmAddressDialogComponent;
  let fixture: ComponentFixture<HwConfirmAddressDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwConfirmAddressDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwConfirmAddressDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
