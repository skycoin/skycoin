import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MsgBarComponent } from './msg-bar.component';

describe('MsgBarComponent', () => {
  let component: MsgBarComponent;
  let fixture: ComponentFixture<MsgBarComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MsgBarComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MsgBarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
