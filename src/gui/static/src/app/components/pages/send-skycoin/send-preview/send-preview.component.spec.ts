import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SendVerifyComponent } from './send-preview.component';

describe('SendVerifyComponent', () => {
  let component: SendVerifyComponent;
  let fixture: ComponentFixture<SendVerifyComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SendVerifyComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SendVerifyComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
