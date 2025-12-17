import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { DisclaimerWarningComponent } from './disclaimer-warning.component';
import { MockTranslatePipe } from '../../../utils/test-mocks';

describe('DisclaimerWarningComponent', () => {
  let component: DisclaimerWarningComponent;
  let fixture: ComponentFixture<DisclaimerWarningComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        DisclaimerWarningComponent,
        MockTranslatePipe
      ],
      schemas: [ NO_ERRORS_SCHEMA ],
      providers: []
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DisclaimerWarningComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
