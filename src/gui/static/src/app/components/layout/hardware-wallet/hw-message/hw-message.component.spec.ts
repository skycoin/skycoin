import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwMessageComponent } from './hw-message.component';

describe('HwMessageComponent', () => {
  let component: HwMessageComponent;
  let fixture: ComponentFixture<HwMessageComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwMessageComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwMessageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
