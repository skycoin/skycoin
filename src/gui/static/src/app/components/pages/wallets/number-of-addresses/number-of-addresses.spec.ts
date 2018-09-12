import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NumberOfAddressesComponent } from './number-of-addresses';

describe('NumberOfAddressesComponent', () => {
  let component: NumberOfAddressesComponent;
  let fixture: ComponentFixture<NumberOfAddressesComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NumberOfAddressesComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NumberOfAddressesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
