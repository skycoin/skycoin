import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwPinDialogComponent } from './hw-pin-dialog.component';

describe('HwPinDialogComponent', () => {
  let component: HwPinDialogComponent;
  let fixture: ComponentFixture<HwPinDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwPinDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwPinDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
