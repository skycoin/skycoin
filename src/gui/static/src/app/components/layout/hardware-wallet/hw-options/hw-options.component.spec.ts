import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HwWalletOptionsComponent } from './hw-options.component';

describe('HwWalletOptionsComponent', () => {
  let component: HwWalletOptionsComponent;
  let fixture: ComponentFixture<HwWalletOptionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ HwWalletOptionsComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HwWalletOptionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
