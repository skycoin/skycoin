import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AddressOptionsComponent } from './address-options.component';

describe('AddressOptionsComponent', () => {
  let component: AddressOptionsComponent;
  let fixture: ComponentFixture<AddressOptionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AddressOptionsComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AddressOptionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
