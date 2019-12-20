import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { FormSourceSelectionComponent } from './form-source-selection.component';

describe('FormSourceSelectionComponent', () => {
  let component: FormSourceSelectionComponent;
  let fixture: ComponentFixture<FormSourceSelectionComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ FormSourceSelectionComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(FormSourceSelectionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
