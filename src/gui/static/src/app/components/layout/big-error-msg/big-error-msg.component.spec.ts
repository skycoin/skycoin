import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BigErrorMsgComponent } from './big-error-msg.component';

describe('BigErrorMsgComponent', () => {
  let component: BigErrorMsgComponent;
  let fixture: ComponentFixture<BigErrorMsgComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BigErrorMsgComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BigErrorMsgComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
