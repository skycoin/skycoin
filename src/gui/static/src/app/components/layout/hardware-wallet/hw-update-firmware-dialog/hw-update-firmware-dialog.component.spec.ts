import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwUpdateFirmwareDialogComponent } from './hw-update-firmware-dialog.component';

describe('HwUpdateFirmwareDialogComponent', () => {
  let component: HwUpdateFirmwareDialogComponent;
  let fixture: ComponentFixture<HwUpdateFirmwareDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwUpdateFirmwareDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwUpdateFirmwareDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
