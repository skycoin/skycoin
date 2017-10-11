import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AddDepositAddressComponent } from './add-deposit-address.component';

describe('AddDepositAddressComponent', () => {
  let component: AddDepositAddressComponent;
  let fixture: ComponentFixture<AddDepositAddressComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AddDepositAddressComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AddDepositAddressComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
