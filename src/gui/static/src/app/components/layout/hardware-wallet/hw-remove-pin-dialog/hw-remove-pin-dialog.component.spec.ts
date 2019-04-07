import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwRemovePinComponent } from './hw-remove-pin-dialog.component';

describe('HwRemovePinComponent', () => {
  let component: HwRemovePinComponent;
  let fixture: ComponentFixture<HwRemovePinComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwRemovePinComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwRemovePinComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
