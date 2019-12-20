import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { FormMultipleDestinationsComponent } from './form-multiple-destinations.component';

describe('FormMultipleDestinationsComponent', () => {
  let component: FormMultipleDestinationsComponent;
  let fixture: ComponentFixture<FormMultipleDestinationsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ FormMultipleDestinationsComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(FormMultipleDestinationsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
