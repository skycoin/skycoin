import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwPinHelpComponent } from './hw-pin-help.component';

describe('HwPinHelpComponent', () => {
  let component: HwPinHelpComponent;
  let fixture: ComponentFixture<HwPinHelpComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwPinHelpComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwPinHelpComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
