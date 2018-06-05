import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SendFormAdvancedComponent } from './send-form-advanced.component';

describe('SendFormAdvancedComponent', () => {
  let component: SendFormAdvancedComponent;
  let fixture: ComponentFixture<SendFormAdvancedComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SendFormAdvancedComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SendFormAdvancedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
