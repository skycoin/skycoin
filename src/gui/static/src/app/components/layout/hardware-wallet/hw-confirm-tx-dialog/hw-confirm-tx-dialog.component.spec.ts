import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwConfirmTxDialogComponent } from './hw-confirm-tx-dialog.component';

describe('HwConfirmTxDialogComponent', () => {
  let component: HwConfirmTxDialogComponent;
  let fixture: ComponentFixture<HwConfirmTxDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwConfirmTxDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwConfirmTxDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
