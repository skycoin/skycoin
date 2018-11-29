import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwChangePinDialogComponent } from './hw-change-pin-dialog.component';

describe('HwChangePinDialogComponent', () => {
  let component: HwChangePinDialogComponent;
  let fixture: ComponentFixture<HwChangePinDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwChangePinDialogComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwChangePinDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
