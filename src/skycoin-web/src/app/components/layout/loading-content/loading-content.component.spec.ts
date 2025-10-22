import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NO_ERRORS_SCHEMA } from '@angular/core';

import { LoadingContentComponent } from './loading-content.component';
import { MockTranslatePipe } from '../../../utils/test-mocks';

describe('LoadingContentComponent', () => {
  let component: LoadingContentComponent;
  let fixture: ComponentFixture<LoadingContentComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        LoadingContentComponent,
        MockTranslatePipe
      ],
      schemas: [ NO_ERRORS_SCHEMA ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LoadingContentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
