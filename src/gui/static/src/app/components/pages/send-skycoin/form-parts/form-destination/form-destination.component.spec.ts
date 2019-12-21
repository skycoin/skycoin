import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { FormDestinationComponent } from './form-destination.component';

describe('FormDestinationComponent', () => {
  let component: FormDestinationComponent;
  let fixture: ComponentFixture<FormDestinationComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ FormDestinationComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(FormDestinationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
