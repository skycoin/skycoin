import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SendCoinsFormComponent } from './send-coins-form.component';

describe('SendCoinsFormComponent', () => {
  let component: SendCoinsFormComponent;
  let fixture: ComponentFixture<SendCoinsFormComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SendCoinsFormComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SendCoinsFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
