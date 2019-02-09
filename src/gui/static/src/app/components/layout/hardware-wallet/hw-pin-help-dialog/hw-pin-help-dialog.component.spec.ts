import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwPinHelpDialogComponent } from './hw-pin-help-dialog.component';

describe('HwPinHelpDialogComponent', () => {
  let component: HwPinHelpDialogComponent;
  let fixture: ComponentFixture<HwPinHelpDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwPinHelpDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwPinHelpDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
