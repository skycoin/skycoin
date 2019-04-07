import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwRemovePinDialogComponent } from './hw-remove-pin-dialog.component';

describe('HwRemovePinDialogComponent', () => {
  let component: HwRemovePinDialogComponent;
  let fixture: ComponentFixture<HwRemovePinDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwRemovePinDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwRemovePinDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
