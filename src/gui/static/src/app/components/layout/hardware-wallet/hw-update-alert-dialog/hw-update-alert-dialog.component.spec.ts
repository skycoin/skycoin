import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwUpdateAlertDialogComponent } from './hw-update-alert-dialog.component';

describe('HwUpdateAlertDialogComponent', () => {
  let component: HwUpdateAlertDialogComponent;
  let fixture: ComponentFixture<HwUpdateAlertDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwUpdateAlertDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwUpdateAlertDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
